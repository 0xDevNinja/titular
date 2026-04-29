// Spawn a foundry `anvil` instance on a free port, expose its HTTP RPC URL,
// and tear it down on close(). Mirrors the Go-side `anvil_test.go` helper in
// `services/indexer-go/internal/integration/` so a contributor familiar with
// one codebase can read the other.
//
// We deliberately use a real anvil child process (not a JS-only EVM like
// hardhat-network or ganache-core) because the SDK's call path leans on
// transaction-receipt semantics — gas, nonce, log decoding — that the
// in-process EVMs occasionally diverge on. Running anvil makes the e2e
// suite a high-fidelity proxy for a Base Sepolia bring-up.

import { type ChildProcess, spawn } from "node:child_process";
import { createServer } from "node:net";

/** Default chain id anvil exposes on `--chain-id 31337`. */
export const ANVIL_CHAIN_ID = 31337;

/**
 * Anvil's account 0 private key. Stable across versions and safe to embed:
 * it is the public dev key shipped with foundry, never used outside
 * ephemeral test instances. Address: 0xf39Fd6e51aad88F6F4ce6aB8827279cfFFb92266.
 */
// gitleaks:allow
export const ANVIL_DEV_PRIVATE_KEY =
  "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" as const;

/** Address derived from {@link ANVIL_DEV_PRIVATE_KEY}. */
export const ANVIL_DEV_ADDRESS = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266" as const;

export interface AnvilHandle {
  /** HTTP JSON-RPC URL, e.g. `http://127.0.0.1:42017`. */
  readonly rpcUrl: string;
  /** Bound port. */
  readonly port: number;
  /** Stop the child process; idempotent. */
  close(): Promise<void>;
}

/** True when `anvil` is on PATH; the suite uses this to skip gracefully. */
export async function anvilAvailable(): Promise<boolean> {
  return new Promise((resolve) => {
    const child = spawn("anvil", ["--version"], { stdio: "ignore" });
    child.on("error", () => resolve(false));
    child.on("exit", (code) => resolve(code === 0));
  });
}

/**
 * Start anvil on a kernel-assigned port and wait until `eth_blockNumber`
 * answers. Throws if the binary is missing — callers should pre-flight with
 * {@link anvilAvailable} and skip the suite when false.
 */
export async function startAnvil(): Promise<AnvilHandle> {
  const port = await pickFreePort();

  // --no-mining: every block must come from an explicit `anvil_mine` (or
  // an implicit one via `eth_sendRawTransaction` with auto-mining set).
  // We keep auto-mining on (default) here because the SDK assumes its
  // submitted tx will be mined; the indexer-go harness disables it because
  // it drives reorg simulations directly. Different goal, different knob.
  const args = [
    "--port",
    String(port),
    "--host",
    "127.0.0.1",
    "--chain-id",
    String(ANVIL_CHAIN_ID),
    "--accounts",
    "5",
    "--silent",
  ];
  const child = spawn("anvil", args, {
    stdio: process.env.ANVIL_TEST_LOGS ? "inherit" : "ignore",
  });

  const closed = new Promise<void>((resolve) => {
    child.once("exit", () => resolve());
  });

  const rpcUrl = `http://127.0.0.1:${port}`;
  await waitForRpcReady(rpcUrl, child, 10_000);

  let closing: Promise<void> | undefined;
  return {
    rpcUrl,
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
            }, 2_000)
          ),
        ]);
      })();
      return closing;
    },
  };
}

async function waitForRpcReady(
  rpcUrl: string,
  child: ChildProcess,
  timeoutMs: number
): Promise<void> {
  const deadline = Date.now() + timeoutMs;
  let lastErr: unknown;
  while (Date.now() < deadline) {
    if (child.exitCode !== null) {
      throw new Error(`anvil exited prematurely with code ${child.exitCode}`);
    }
    try {
      const r = await fetch(rpcUrl, {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({
          jsonrpc: "2.0",
          id: 1,
          method: "eth_blockNumber",
          params: [],
        }),
      });
      if (r.ok) return;
    } catch (err) {
      lastErr = err;
    }
    await sleep(50);
  }
  throw new Error(
    `anvil never answered eth_blockNumber within ${timeoutMs}ms (last error: ${lastErr ?? "n/a"})`
  );
}

async function pickFreePort(): Promise<number> {
  // Bind 0.0.0.0:0, ask the kernel for the assigned port, close, hand the
  // number to anvil. There is a tiny TOCTOU window between close and anvil
  // re-binding; on a developer laptop this is well below the noise floor
  // and `cmd.Start` would fail loudly on collision anyway.
  return new Promise((resolve, reject) => {
    const srv = createServer();
    srv.unref();
    srv.on("error", reject);
    srv.listen(0, "127.0.0.1", () => {
      const addr = srv.address();
      srv.close(() => {
        if (addr === null || typeof addr === "string") {
          reject(new Error("pickFreePort: unexpected null address from net.Server"));
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

/**
 * Install raw runtime bytecode at `address` via the `anvil_setCode` cheat.
 * Skips the constructor entirely — pass the runtime form (post-constructor)
 * not the deploy form.
 */
export async function anvilSetCode(
  rpcUrl: string,
  address: `0x${string}`,
  runtimeHex: `0x${string}`
): Promise<void> {
  const r = await fetch(rpcUrl, {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({
      jsonrpc: "2.0",
      id: 1,
      method: "anvil_setCode",
      params: [address, runtimeHex],
    }),
  });
  if (!r.ok) {
    throw new Error(`anvil_setCode @ ${address}: HTTP ${r.status}`);
  }
  const body = (await r.json()) as { error?: { message?: string } };
  if (body.error !== undefined) {
    throw new Error(`anvil_setCode @ ${address}: ${body.error.message ?? "rpc error"}`);
  }
}
