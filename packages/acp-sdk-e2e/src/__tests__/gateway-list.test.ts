// SDK -> gateway-go end-to-end against a real Postgres + go-run gateway.
//
// Boot order:
//
//   1. testcontainers Postgres (16-alpine), pinned to match the indexer's
//      integration tests.
//   2. Apply the indexer's embedded migration set (0001_init,
//      0002_pagination_indexes).
//   3. Seed two ACP agents + one job directly via SQL — the indexer
//      pipeline is exercised by `services/indexer-go/internal/integration`,
//      so this suite focuses on the SDK-to-gateway wire shape only.
//   4. `go run ./cmd/gateway` against the seeded DB.
//   5. Drive `GatewayClient.listAgents`, `listJobs`, `getStats`,
//      `getAgent` (success + 404) and assert the SDK shapes match the
//      seeded values.

import { AcpError, GatewayClient } from "@titular/acp-sdk";
import { afterAll, beforeAll, describe, expect, it } from "vitest";
import { type GatewayHandle, goAvailable, startGateway } from "../lib/gateway.js";
import { type PostgresHandle, seedAgent, seedJob, startPostgres } from "../lib/postgres.js";
import { E2E_ENABLED, SKIP_REASON } from "./_skip.js";

const ALICE = "0x00000000000000000000000000000000000000a1" as const;
const BOB = "0x00000000000000000000000000000000000000b2" as const;
const TOKEN = "0x000000000000000000000000000000000000c0de" as const;
const TX_HASH_AGENT_A =
  "0x1111111111111111111111111111111111111111111111111111111111111111" as const;
const TX_HASH_AGENT_B =
  "0x2222222222222222222222222222222222222222222222222222222222222222" as const;
const TX_HASH_JOB = "0x3333333333333333333333333333333333333333333333333333333333333333" as const;

interface Fixture {
  pg: PostgresHandle;
  gw: GatewayHandle;
  client: GatewayClient;
  alicePk: number;
  bobPk: number;
}

describe.skipIf(!E2E_ENABLED)("GatewayClient against gateway-go + Postgres", () => {
  let fx: Fixture | undefined;

  beforeAll(async (ctx) => {
    if (!E2E_ENABLED) {
      ctx.skip(SKIP_REASON);
      return;
    }
    if (!(await goAvailable())) {
      ctx.skip("go not on PATH; install go 1.25+ to run this suite");
      return;
    }

    let pg: PostgresHandle | undefined;
    try {
      pg = await startPostgres();
    } catch (err) {
      ctx.skip(`docker daemon unreachable: ${err instanceof Error ? err.message : String(err)}`);
      return;
    }

    const alicePk = await seedAgent(pg, {
      agentId: 1n,
      controller: ALICE,
      metadataUri: "ipfs://bafyalice",
      capabilities: 0b1n,
      blockNumber: 100n,
      txHash: TX_HASH_AGENT_A,
      logIndex: 0,
    });
    const bobPk = await seedAgent(pg, {
      agentId: 2n,
      controller: BOB,
      metadataUri: "ipfs://bafybob",
      capabilities: 0b11n,
      blockNumber: 101n,
      txHash: TX_HASH_AGENT_B,
      logIndex: 1,
    });
    await seedJob(pg, {
      onChainId: 7n,
      cloneAddress: "0x00000000000000000000000000000000000000c1" as const,
      agentPk: alicePk,
      agentId: 1n,
      principal: BOB,
      jobType: 0,
      token: TOKEN,
      budget: 5_000n,
      deadline: new Date(Date.now() + 3_600_000),
      blockNumber: 102n,
      txHash: TX_HASH_JOB,
      logIndex: 2,
    });

    let gw: GatewayHandle;
    try {
      gw = await startGateway({ databaseUrl: pg.databaseUrl });
    } catch (err) {
      await pg.close();
      throw err;
    }

    fx = {
      pg,
      gw,
      client: new GatewayClient({ gatewayUrl: gw.baseUrl }),
      alicePk,
      bobPk,
    };
  }, 240_000);

  afterAll(async () => {
    if (fx === undefined) return;
    await fx.gw.close();
    await fx.pg.close();
  });

  it("lists the seeded ACP agents", async () => {
    if (fx === undefined) return; // suite was skipped
    const page = await fx.client.listAgents({ kind: "acp", limit: 10 });
    expect(page.data.length).toBe(2);
    const ids = page.data.map((a) => a.agent_id).sort();
    expect(ids).toEqual(["1", "2"]);
    const alice = page.data.find((a) => a.agent_id === "1");
    expect(alice?.controller).toBe(ALICE);
    expect(alice?.metadata_uri).toBe("ipfs://bafyalice");
    expect(alice?.kind).toBe("acp");
  });

  it("returns the seeded agent by primary key", async () => {
    if (fx === undefined) return; // suite was skipped
    const agent = await fx.client.getAgent(fx.alicePk);
    expect(agent.id).toBe(fx.alicePk);
    expect(agent.agent_id).toBe("1");
    expect(agent.controller).toBe(ALICE);
  });

  it("raises AcpError gateway_error with status 404 for an unknown agent", async () => {
    if (fx === undefined) return; // suite was skipped
    await expect(fx.client.getAgent(999_999_999)).rejects.toMatchObject({
      name: "AcpError",
      code: "gateway_error",
      status: 404,
    });
    // Sanity-check the discriminant works for downstream branching.
    const captured = await fx.client.getAgent(999_999_999).catch((e) => e);
    expect(captured).toBeInstanceOf(AcpError);
  });

  it("lists the seeded job", async () => {
    if (fx === undefined) return; // suite was skipped
    const page = await fx.client.listJobs({ limit: 10 });
    expect(page.data.length).toBe(1);
    const job = page.data[0];
    expect(job?.on_chain_id).toBe("7");
    expect(job?.principal).toBe(BOB);
    expect(job?.budget).toBe("5000");
    expect(job?.token).toBe(TOKEN);
    expect(job?.phase).toBe("created");
  });

  it("returns aggregate stats consistent with the seed set", async () => {
    if (fx === undefined) return; // suite was skipped
    const stats = await fx.client.getStats();
    expect(stats.total_agents).toBeGreaterThanOrEqual(2);
    expect(stats.total_jobs).toBeGreaterThanOrEqual(1);
    expect(typeof stats.last_block_indexed).toBe("number");
  });

  it("raises AcpError gateway_unreachable when the gateway is not listening", async () => {
    // A fresh client pointed at a port nothing is bound to surfaces the
    // network error code path. We keep the timeout short so the assertion
    // does not pad the suite runtime.
    const dead = new GatewayClient({
      gatewayUrl: "http://127.0.0.1:1",
      timeoutMs: 500,
    });
    await expect(dead.listAgents()).rejects.toMatchObject({
      name: "AcpError",
      code: "gateway_unreachable",
    });
  });
});
