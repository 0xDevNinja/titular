// Spin up a Postgres container for the e2e suite, seeded with the indexer's
// embedded migration set so the gateway reads against the same schema
// production uses. Uses `testcontainers` rather than docker-compose so the
// suite owns the container lifetime end-to-end and a crashed test cannot
// leak a stray Postgres on the developer machine.

import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import pg from "pg";
import { GenericContainer, Wait } from "testcontainers";

const here = dirname(fileURLToPath(import.meta.url));

/** Resolved path to the indexer-go migrations folder, monorepo-relative. */
const MIGRATIONS_DIR = resolve(here, "../../../../services/indexer-go/internal/db/migrations");

/** Migration files to apply in order. Keep in sync with the indexer's `migrations.go`. */
const MIGRATION_FILES = ["0001_init.up.sql", "0002_pagination_indexes.up.sql"] as const;

export interface PostgresHandle {
  /** `postgresql://...` URL the gateway connects with. */
  readonly databaseUrl: string;
  /** Underlying container host (`127.0.0.1` on Linux runners). */
  readonly host: string;
  /** Mapped port on the host. */
  readonly port: number;
  /** Database name. */
  readonly database: string;
  /** Username. */
  readonly username: string;
  /** Password. */
  readonly password: string;
  /** Run a parameterised query against the seeded database. */
  query<T extends pg.QueryResultRow = pg.QueryResultRow>(
    text: string,
    params?: unknown[]
  ): Promise<pg.QueryResult<T>>;
  /** Stop and remove the container; idempotent. */
  close(): Promise<void>;
}

/**
 * Boot Postgres 16, run the indexer's migration set, return the connection
 * coordinates. Throws on docker daemon errors — callers should catch and
 * `ctx.skip()` so the suite degrades gracefully on hosts without docker.
 */
export async function startPostgres(): Promise<PostgresHandle> {
  // Match the indexer integration tests' postgres:16-alpine pin so the
  // schema-equivalence guarantee is binary-identical, not just SQL-level.
  // testcontainers' `latest` would silently jump majors on reproducer
  // reruns weeks apart, which would defeat the point of pinning.
  //
  // Use `GenericContainer` rather than the v9-era `PostgreSqlContainer`
  // module: in v10+ that lives in a separate `@testcontainers/postgresql`
  // package, and reaching into `testcontainers/build/...` is a private
  // path that breaks across minor releases. The generic shape below is
  // what `PostgreSqlContainer` does internally — env vars + the same
  // log-readiness probe `pg_isready` boils down to.
  const database = "titular_e2e";
  const username = "titular";
  const password = "titular";
  const container = await new GenericContainer("postgres:16-alpine")
    .withEnvironment({
      POSTGRES_DB: database,
      POSTGRES_USER: username,
      POSTGRES_PASSWORD: password,
    })
    .withExposedPorts(5432)
    .withWaitStrategy(
      // Postgres logs the readiness banner twice on first boot: once when
      // the temp init server comes up, again after `init.d` runs. Wait
      // for the second to avoid racing the migration step against
      // initdb's bootstrap.
      Wait.forLogMessage("database system is ready to accept connections", 2)
    )
    .start();

  const host = container.getHost();
  const port = container.getMappedPort(5432);
  const databaseUrl = `postgresql://${username}:${password}@${host}:${port}/${database}`;
  const pool = new pg.Pool({ connectionString: databaseUrl, max: 4 });

  await applyMigrations(pool);

  let closing: Promise<void> | undefined;
  return {
    databaseUrl,
    host,
    port,
    database,
    username,
    password,
    async query(text, params) {
      return pool.query(text, params);
    },
    async close() {
      if (closing !== undefined) return closing;
      closing = (async () => {
        await pool.end().catch(() => {});
        await container.stop({ timeout: 5_000 }).catch(() => {});
      })();
      return closing;
    },
  };
}

async function applyMigrations(pool: pg.Pool): Promise<void> {
  // Run each migration as a single multi-statement query. The indexer's
  // golang-migrate runner wraps each file in a transaction; pg doesn't do
  // that automatically for `query(text)` so we wrap explicitly. A bad
  // migration mid-flight then rolls back to the previous on-disk schema
  // instead of leaving a half-applied object graph.
  for (const filename of MIGRATION_FILES) {
    const sql = readFileSync(resolve(MIGRATIONS_DIR, filename), "utf8");
    const client = await pool.connect();
    try {
      await client.query("BEGIN");
      await client.query(sql);
      await client.query("COMMIT");
    } catch (err) {
      await client.query("ROLLBACK").catch(() => {});
      throw new Error(
        `migration ${filename} failed: ${err instanceof Error ? err.message : String(err)}`
      );
    } finally {
      client.release();
    }
  }
}

/**
 * Insert a hand-crafted `acp` agent row. Mirrors what the indexer's
 * `decoder.OnAgentRegistered` would persist: bytea for controller / tx
 * hash, `numeric(78,0)` for the on-chain uint256, lowercase address.
 */
export interface SeedAgentRow {
  agentId: bigint;
  controller: `0x${string}`;
  metadataUri: string;
  capabilities: bigint;
  blockNumber: bigint;
  txHash: `0x${string}`;
  logIndex: number;
}

export interface SeedJobRow {
  onChainId: bigint;
  cloneAddress: `0x${string}`;
  agentPk: number | null;
  agentId: bigint | null;
  principal: `0x${string}`;
  jobType: number;
  token: `0x${string}`;
  budget: bigint;
  deadline: Date;
  blockNumber: bigint;
  txHash: `0x${string}`;
  logIndex: number;
}

/**
 * Insert one ACP agent row and return its bigserial `id`. Used by the
 * gateway-side test to seed Postgres without spinning the indexer pipeline.
 */
export async function seedAgent(handle: PostgresHandle, row: SeedAgentRow): Promise<number> {
  const result = await handle.query<{ id: string }>(
    `INSERT INTO agents (kind, agent_id, controller, capabilities, metadata_uri, block_number, tx_hash, log_index)
     VALUES ('acp', $1, $2, $3, $4, $5, $6, $7) RETURNING id`,
    [
      row.agentId.toString(),
      hexToBuf(row.controller),
      row.capabilities.toString(),
      row.metadataUri,
      row.blockNumber.toString(),
      hexToBuf(row.txHash),
      row.logIndex,
    ]
  );
  const idStr = result.rows[0]?.id;
  if (idStr === undefined) {
    throw new Error("seedAgent: INSERT returned no row");
  }
  return Number(idStr);
}

/** Insert one job row; returns its bigserial `id`. */
export async function seedJob(handle: PostgresHandle, row: SeedJobRow): Promise<number> {
  const result = await handle.query<{ id: string }>(
    `INSERT INTO jobs (on_chain_id, clone_address, agent_pk, agent_id, principal, job_type, token, budget, deadline, block_number, tx_hash, log_index)
     VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`,
    [
      row.onChainId.toString(),
      hexToBuf(row.cloneAddress),
      row.agentPk,
      row.agentId === null ? null : row.agentId.toString(),
      hexToBuf(row.principal),
      row.jobType,
      hexToBuf(row.token),
      row.budget.toString(),
      row.deadline.toISOString(),
      row.blockNumber.toString(),
      hexToBuf(row.txHash),
      row.logIndex,
    ]
  );
  const idStr = result.rows[0]?.id;
  if (idStr === undefined) {
    throw new Error("seedJob: INSERT returned no row");
  }
  return Number(idStr);
}

/** Convert a `0x...` hex string to a Postgres bytea-friendly Buffer. */
function hexToBuf(hex: `0x${string}`): Buffer {
  const stripped = hex.startsWith("0x") ? hex.slice(2) : hex;
  return Buffer.from(stripped, "hex");
}
