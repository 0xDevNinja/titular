/**
 * Anvil fixture: spawns a local EVM node, deploys Phase1 + Phase2 contracts,
 * and returns wired addresses + a funded signer mnemonic for tests.
 *
 * The fixture is used as a Playwright global setup step so the contracts are
 * deployed once per test run (not per test).
 *
 * Required env:
 *   ANVIL_PORT       — defaults to 8545
 *   CONTRACTS_DIR    — path to contracts/evm (default resolved from cwd)
 *
 * Outputs (written to process.env so Playwright workers can read them):
 *   ANVIL_RPC_URL, TITU_ADDR, TREASURY_ADDR, FACTORY_ADDR, GRADUATOR_ADDR,
 *   DEPLOYER_MNEMONIC, ANVIL_CHAIN_ID
 */

import { type ChildProcess, execSync, spawn } from "node:child_process";
import fs from "node:fs";
import path from "node:path";
import { setTimeout as sleep } from "node:timers/promises";

/**
 * Poll http://127.0.0.1:<port> with eth_blockNumber until the JSON-RPC
 * socket is accepting connections. Throws if the node is not ready within
 * the allotted attempts.
 *
 * @param rpcUrl  - full URL, e.g. "http://127.0.0.1:8545"
 * @param retries - number of attempts (default 60)
 * @param delayMs - milliseconds between attempts (default 50)
 */
async function waitForRpc(rpcUrl: string, retries = 60, delayMs = 50): Promise<void> {
  for (let attempt = 0; attempt < retries; attempt++) {
    try {
      const res = await fetch(rpcUrl, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          jsonrpc: "2.0",
          method: "eth_blockNumber",
          params: [],
          id: 1,
        }),
      });
      if (res.ok) {
        const body = (await res.json()) as { result?: string };
        if (body.result !== undefined) return;
      }
    } catch {
      // connection refused — anvil not up yet
    }
    await sleep(delayMs);
  }
  throw new Error(
    `[anvil] RPC at ${rpcUrl} did not respond after ${retries * delayMs} ms. Check that anvil is installed and the port is available.`
  );
}

export const ANVIL_MNEMONIC = "test test test test test test test test test test test junk";
export const DEPLOYER_PRIVATE_KEY =
  "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80";
export const DEPLOYER_ADDRESS = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266";

const ANVIL_PORT = Number(process.env.ANVIL_PORT ?? "8545");
const CONTRACTS_DIR =
  process.env.CONTRACTS_DIR ?? path.resolve(path.join(process.cwd(), "../../contracts/evm"));

let anvilProcess: ChildProcess | null = null;

/**
 * Start anvil and wait until it is accepting connections.
 */
export async function startAnvil(): Promise<string> {
  const rpcUrl = `http://127.0.0.1:${ANVIL_PORT}`;

  // Resolve the anvil binary path explicitly so spawn never relies solely on
  // PATH resolution (which can differ between the shell step and the Node child).
  const anvilBin = (() => {
    try {
      return execSync("which anvil", { stdio: "pipe" }).toString().trim();
    } catch {
      return "anvil"; // fall back to PATH lookup
    }
  })();

  console.log(`[anvil] binary: ${anvilBin}`);

  anvilProcess = spawn(
    anvilBin,
    [
      "--port",
      String(ANVIL_PORT),
      "--mnemonic",
      ANVIL_MNEMONIC,
      "--chain-id",
      "31337",
      // omit --block-time: anvil defaults to instant (on-demand) mining.
      // foundry >=1.5.1 rejects --block-time 0 as invalid.
      "--silent", // suppress banner noise in CI logs
    ],
    // Pipe stdout/stderr so we can surface startup errors. The streams are
    // forwarded to the parent process so CI logs show any anvil crash reason.
    { stdio: ["ignore", "pipe", "pipe"] }
  );

  // Forward anvil output to the parent process so errors are visible in CI.
  anvilProcess.stdout?.pipe(process.stdout);
  anvilProcess.stderr?.pipe(process.stderr);

  anvilProcess.on("error", (err) => {
    console.error("[anvil] spawn error:", err.message);
  });

  anvilProcess.on("close", (code) => {
    if (code !== null && code !== 0) {
      console.error(`[anvil] process exited early with code ${code}`);
    }
  });

  // Poll until the JSON-RPC socket is accepting connections.
  // 300 attempts × 100 ms = 30 s — enough for CI runner cold-start.
  await waitForRpc(rpcUrl, 300, 100);

  return rpcUrl;
}

/**
 * Deploy Phase1 + Phase2 contracts against the running anvil node.
 * Returns the deployed contract addresses extracted from the broadcast
 * artifacts written by forge.
 */
export function deployContracts(rpcUrl: string): AnvilAddresses {
  const deployDir = path.join(CONTRACTS_DIR, "deployments");

  // Phase 1 — TITU, Treasury, VestingVault, VeTITU, FeeDistributor
  execSync(
    `forge script script/DeployPhase1.s.sol:DeployPhase1 --rpc-url ${rpcUrl} --broadcast --private-key ${DEPLOYER_PRIVATE_KEY} --sender ${DEPLOYER_ADDRESS}`,
    {
      cwd: CONTRACTS_DIR,
      env: {
        ...process.env,
        DEPLOYER_PRIVATE_KEY,
        MULTISIG: DEPLOYER_ADDRESS,
      },
      stdio: "pipe",
    }
  );

  // Deploy mock Uniswap V2 infrastructure via a small inline cast call
  // sequence so Phase 2 can reference a real factory + router on anvil.
  // We use the mocks checked in under test/mocks/.
  const mockAddresses = deployMockUniswap(rpcUrl);

  // Phase 2 — LaunchpadFactory, Graduator, modules, etc.
  execSync(
    `forge script script/DeployPhase2.s.sol:DeployPhase2 --rpc-url ${rpcUrl} --broadcast --private-key ${DEPLOYER_PRIVATE_KEY} --sender ${DEPLOYER_ADDRESS}`,
    {
      cwd: CONTRACTS_DIR,
      env: {
        ...process.env,
        DEPLOYER_PRIVATE_KEY,
        MULTISIG: DEPLOYER_ADDRESS,
        UNISWAP_V2_FACTORY: mockAddresses.uniswapFactory,
        UNISWAP_V2_ROUTER: mockAddresses.uniswapRouter,
      },
      stdio: "pipe",
    }
  );

  return parseDeploymentJson(deployDir);
}

/**
 * Deploy MockUniswapV2Factory + MockUniswapV2Router using forge create
 * (they're test-only contracts so the deploy script doesn't cover them).
 */
function deployMockUniswap(rpcUrl: string): {
  uniswapFactory: string;
  uniswapRouter: string;
} {
  // Parse the deployed address from `forge create` text output.
  // The --json flag changed its schema across foundry versions (pre-1.5 had
  // `deployedTo`, 1.5.1 outputs a transaction object with no address field).
  // Parsing the human-readable "Deployed to: 0x..." line is stable across all
  // foundry versions.
  function extractAddress(output: string): string {
    const match = output.match(/Deployed to:\s*(0x[0-9a-fA-F]{40})/);
    if (!match) throw new Error(`forge create output missing "Deployed to:" line:\n${output}`);
    return match[1];
  }

  const factoryOut = execSync(
    `forge create test/mocks/MockUniswapV2Factory.sol:MockUniswapV2Factory --rpc-url ${rpcUrl} --private-key ${DEPLOYER_PRIVATE_KEY}`,
    { cwd: CONTRACTS_DIR, stdio: "pipe" }
  ).toString();

  const factoryAddr = extractAddress(factoryOut);

  const routerOut = execSync(
    `forge create test/mocks/MockUniswapV2Router.sol:MockUniswapV2Router --rpc-url ${rpcUrl} --private-key ${DEPLOYER_PRIVATE_KEY} --constructor-args ${factoryAddr}`,
    { cwd: CONTRACTS_DIR, stdio: "pipe" }
  ).toString();

  const routerAddr = extractAddress(routerOut);

  return { uniswapFactory: factoryAddr, uniswapRouter: routerAddr };
}

/** Addresses surfaced by the fixture for test assertions. */
export interface AnvilAddresses {
  titu: string;
  treasury: string;
  factory: string;
  graduator: string;
  feeRouter: string;
  rpcUrl: string;
  chainId: number;
  deployerAddress: string;
  deployerMnemonic: string;
}

/**
 * Parse the merged deployment JSON written by DeployPhase2.
 */
function parseDeploymentJson(deployDir: string): AnvilAddresses {
  // DeployPhase2 writes <chainId>.json; anvil chain id is 31337.
  const candidates = ["31337.json", "base-sepolia.json"];
  let raw: string | undefined;
  for (const name of candidates) {
    const p = path.join(deployDir, name);
    if (fs.existsSync(p)) {
      raw = fs.readFileSync(p, "utf-8");
      break;
    }
  }
  if (!raw) {
    throw new Error(
      `Deployment JSON not found in ${deployDir}. Looked for: ${candidates.join(", ")}`
    );
  }

  const json: DeploymentJson = JSON.parse(raw);

  // Support both the merged format (phase1/phase2 keys) and legacy flat format.
  const p1 = json.phase1?.contracts ?? json.contracts;
  const p2 = json.phase2?.contracts;

  if (!p1?.TITU || !p1?.Treasury) {
    throw new Error("Phase1 deployment JSON missing TITU or Treasury address");
  }
  if (!p2?.LaunchpadFactory || !p2?.Graduator || !p2?.FeeRouter) {
    throw new Error("Phase2 deployment JSON missing LaunchpadFactory, Graduator, or FeeRouter");
  }

  return {
    titu: p1.TITU,
    treasury: p1.Treasury,
    factory: p2.LaunchpadFactory,
    graduator: p2.Graduator,
    feeRouter: p2.FeeRouter,
    rpcUrl: `http://127.0.0.1:${ANVIL_PORT}`,
    chainId: 31337,
    deployerAddress: DEPLOYER_ADDRESS,
    deployerMnemonic: ANVIL_MNEMONIC,
  };
}

interface DeploymentJson {
  chainId?: number;
  phase1?: {
    deployedAt?: number;
    contracts?: Record<string, string | null>;
  };
  phase2?: {
    deployedAt?: number;
    contracts?: Record<string, string | null>;
  };
  // Legacy flat format (Phase1-only files).
  contracts?: Record<string, string | null>;
}

/**
 * Stop the anvil process started by {startAnvil}.
 */
export function stopAnvil(): void {
  if (anvilProcess) {
    anvilProcess.kill("SIGTERM");
    anvilProcess = null;
  }
}
