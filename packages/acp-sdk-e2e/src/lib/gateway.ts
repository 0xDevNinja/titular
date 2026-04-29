// Spawn `services/gateway-go` as a child process bound to a free port and
// pointed at the test Postgres. We `go run` against the source tree (not a
// pre-built binary) because:
//
//   1. The gateway's release Dockerfile builds the same source via a
//      multi-stage build, so `go run` exercises an identical compile path.
//   2. CI already has Go 1.25 installed for the unit tests (test.yml's
//      `go-test` job), so no extra toolchain step is needed for this
//      workflow.
//   3. Local contributors who already have Go installed don't need to
//      docker-build the gateway image just to run the e2e suite.
//
// The gateway is started with the smallest viable env set: DATABASE_URL +
// LISTEN ADDR. Auth (JWT/Redis) is left disabled so the suite covers the
// public, unauthenticated REST surface. SSE/GraphQL subscriptions need
// NATS, which we deliberately do NOT spin up here — the indexer pipeline
// is exercised by `services/indexer-go/internal/integration/...`; this
// suite focuses on the SDK ↔ gateway REST contract.

import { type ChildProcess, spawn } from "node:child_process";
import { createServer } from "node:net";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));

/** Path to the gateway-go module root (`services/gateway-go`). */
export const GATEWAY_GO_DIR = resolve(here, "../../../../services/gateway-go");

export interface GatewayHandle {
  /** Base URL the SDK's `GatewayClient` should target. */
  readonly baseUrl: string;
  /** Bound port. */
  readonly port: number;
  /** Stop the child process; idempotent. */
  close(): Promise<void>;
}

/** True when `go` is on PATH; the suite uses this to skip gracefully. */
export async function goAvailable(): Promise<boolean> {
  return new Promise((resolve) => {
    const child = spawn("go", ["version"], { stdio: "ignore" });
    child.on("error", () => resolve(false));
    child.on("exit", (code) => resolve(code === 0));
  });
}

export interface StartGatewayParams {
  /** Postgres DSN the gateway should read from. */
  databaseUrl: string;
  /** Optional explicit port; defaults to a kernel-assigned free port. */
  port?: number;
  /** Optional override for the gateway-go module path. */
  goModuleDir?: string;
}

/**
 * `go run ./cmd/gateway` with `GATEWAY_DATABASE_URL` set, wait until the
 * process answers `GET /healthz` (or, fallback, any 4xx so we know the
 * listener is up). Returns a handle whose `close()` SIGTERMs the child.
 */
export async function startGateway(params: StartGatewayParams): Promise<GatewayHandle> {
  const port = params.port ?? (await pickFreePort());
  const cwd = params.goModuleDir ?? GATEWAY_GO_DIR;

  const env: NodeJS.ProcessEnv = {
    ...process.env,
    GATEWAY_ADDR: `:${port}`,
    GATEWAY_DATABASE_URL: params.databaseUrl,
    // Disable rate limiting in the test — we do small bursts of requests
    // synchronously and a 0-RPS bucket would surface as flaky 429s instead
    // of as the assertion the test actually wants to make.
    GATEWAY_RATE_LIMIT_RPS: "0",
    GATEWAY_RATE_LIMIT_BURST: "0",
  };

  const child = spawn("go", ["run", "./cmd/gateway"], {
    cwd,
    env,
    stdio: process.env.GATEWAY_TEST_LOGS ? "inherit" : "ignore",
  });

  const closed = new Promise<void>((resolve) => {
    child.once("exit", () => resolve());
  });

  const baseUrl = `http://127.0.0.1:${port}`;
  await waitForHttpReady(baseUrl, child, 60_000);

  let closing: Promise<void> | undefined;
  return {
    baseUrl,
    port,
    close(): Promise<void> {
      if (closing !== undefined) return closing;
      closing = (async () => {
        if (!child.killed) child.kill("SIGTERM");
        await Promise.race([
          closed,
          new Promise<void>((resolve) =>
            setTimeout(() => {
              if (!child.killed) child.kill("SIGKILL");
              resolve();
            }, 5_000)
          ),
        ]);
      })();
      return closing;
    },
  };
}

async function waitForHttpReady(
  baseUrl: string,
  child: ChildProcess,
  timeoutMs: number
): Promise<void> {
  // `go run` first compiles, which can take 10-30s on a cold runner. We
  // poll a known mounted path (`/healthz`) and treat any HTTP response —
  // including 404 — as "the listener is up", since we only care that the
  // network surface is bound.
  const deadline = Date.now() + timeoutMs;
  let lastErr: unknown;
  while (Date.now() < deadline) {
    if (child.exitCode !== null) {
      throw new Error(`gateway exited prematurely with code ${child.exitCode}`);
    }
    try {
      const r = await fetch(`${baseUrl}/healthz`);
      if (r.status > 0) return;
    } catch (err) {
      lastErr = err;
    }
    await sleep(200);
  }
  throw new Error(
    `gateway never became ready within ${timeoutMs}ms (last error: ${lastErr ?? "n/a"})`
  );
}

async function pickFreePort(): Promise<number> {
  return new Promise((resolve, reject) => {
    const srv = createServer();
    srv.unref();
    srv.on("error", reject);
    srv.listen(0, "127.0.0.1", () => {
      const addr = srv.address();
      srv.close(() => {
        if (addr === null || typeof addr === "string") {
          reject(new Error("pickFreePort: unexpected null address"));
          return;
        }
        resolve(addr.port);
      });
    });
  });
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
