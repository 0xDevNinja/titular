/**
 * Playwright global setup: starts anvil and deploys contracts once before
 * any worker spins up. Addresses are written to process.env so every worker
 * can read them via globalThis / env in the test files.
 *
 * Global teardown (`global-teardown.ts`) stops anvil after all tests finish.
 */

import { type AnvilAddresses, deployContracts, startAnvil } from "./fixtures/anvil.js";

// Augment the global object so the addresses survive across workers.
declare global {
  var __anvilAddresses: AnvilAddresses | undefined;
}

export default async function globalSetup(): Promise<void> {
  const rpcUrl = await startAnvil();
  const addresses = deployContracts(rpcUrl);
  globalThis.__anvilAddresses = addresses;

  // Expose to process.env so tests that run in separate processes can read them.
  process.env.ANVIL_RPC_URL = addresses.rpcUrl;
  process.env.ANVIL_CHAIN_ID = String(addresses.chainId);
  process.env.TITU_ADDR = addresses.titu;
  process.env.TREASURY_ADDR = addresses.treasury;
  process.env.FACTORY_ADDR = addresses.factory;
  process.env.GRADUATOR_ADDR = addresses.graduator;
  process.env.FEE_ROUTER_ADDR = addresses.feeRouter;
}
