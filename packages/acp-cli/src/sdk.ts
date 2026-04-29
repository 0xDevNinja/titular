// Adapter that builds the SDK objects each command needs from the persisted
// `CliConfig`. Concentrating this in one place keeps the commands free of
// SDK-construction boilerplate and gives tests a single seam to mock.
//
// Two shapes are exposed:
//   - `buildGatewayClient(cfg)` — read-only commands (list, browse, status).
//   - `buildAcpAgent(cfg)`     — write commands (register, createJob).
//
// `buildAcpAgent` requires an `EvmProviderAdapter`. The default Alchemy
// adapter ships with the SDK; tests inject a stub by overriding the
// `providerFactory` parameter so they never hit the network.

import {
  AcpAgent,
  type AcpContractAddresses,
  AlchemyEvmProviderAdapter,
  type EvmProviderAdapter,
  GatewayClient,
  type SupportedChain as SdkSupportedChain,
} from "@titular/acp-sdk";

import type { CliConfig, SupportedChain } from "./types.js";

/** Default contract addresses per chain, used when the persisted config does
 * not pin its own `agentRegistry` / `jobFactory`. The CLI has to pick *some*
 * addresses or every fresh user gets a "missing contract address" error on
 * their first command; these are the well-known M3 deploys (zero placeholders
 * here — operators targeting a private deploy override via
 * `acp configure --agent-registry`). */
const ZERO_ADDRESS = "0x0000000000000000000000000000000000000000" as const;

const CHAIN_DEFAULTS: Record<
  SupportedChain,
  { contracts: AcpContractAddresses; alchemyChain: SdkSupportedChain | null }
> = {
  "base-sepolia": {
    contracts: { agentRegistry: ZERO_ADDRESS, jobFactory: ZERO_ADDRESS },
    alchemyChain: "base-sepolia",
  },
  "base-mainnet": {
    contracts: { agentRegistry: ZERO_ADDRESS, jobFactory: ZERO_ADDRESS },
    alchemyChain: "base",
  },
  "ethereum-sepolia": {
    // The SDK provider only ships with `base` / `base-sepolia` shorthands in
    // the alpha. Selecting an Ethereum chain forces the user to pass an
    // explicit `--rpc-url` and use a custom adapter; we surface that as a
    // descriptive error rather than a viem stack trace.
    contracts: { agentRegistry: ZERO_ADDRESS, jobFactory: ZERO_ADDRESS },
    alchemyChain: null,
  },
  "ethereum-mainnet": {
    contracts: { agentRegistry: ZERO_ADDRESS, jobFactory: ZERO_ADDRESS },
    alchemyChain: null,
  },
};

/** Build a read-only `GatewayClient` from the persisted config. */
export function buildGatewayClient(cfg: CliConfig, fetchImpl?: typeof fetch): GatewayClient {
  return new GatewayClient({
    gatewayUrl: cfg.gatewayUrl,
    ...(cfg.apiKey !== undefined ? { bearerToken: cfg.apiKey } : {}),
    ...(fetchImpl !== undefined ? { fetch: fetchImpl } : {}),
  });
}

/** Resolve the contract address pair for the configured chain, applying
 * per-config overrides on top of the chain defaults. */
export function resolveContracts(cfg: CliConfig): AcpContractAddresses {
  const defaults = CHAIN_DEFAULTS[cfg.chain].contracts;
  return {
    agentRegistry: cfg.agentRegistry ?? defaults.agentRegistry,
    jobFactory: cfg.jobFactory ?? defaults.jobFactory,
  };
}

/** Resolve the Alchemy `SupportedChain` slug for the configured CLI chain.
 * Throws when the chain is not yet wired in the SDK adapter. */
export function resolveAlchemyChain(cfg: CliConfig): SdkSupportedChain {
  const c = CHAIN_DEFAULTS[cfg.chain].alchemyChain;
  if (c === null) {
    throw new Error(
      `chain ${cfg.chain} is not yet supported by the bundled provider; pass a custom --rpc-url or use base/base-sepolia for the alpha`
    );
  }
  return c;
}

/** Hook so tests can swap the provider construction without touching disk
 * or hitting Alchemy. */
export type ProviderFactory = (cfg: CliConfig) => EvmProviderAdapter;

const defaultProviderFactory: ProviderFactory = (cfg) => {
  const chain = resolveAlchemyChain(cfg);
  if (cfg.apiKey === undefined || cfg.apiKey.length === 0) {
    throw new Error("missing alchemy api key: set `apiKey` in the config or pass --api-key");
  }
  return new AlchemyEvmProviderAdapter({
    chain,
    apiKey: cfg.apiKey,
    ...(cfg.privateKey !== undefined ? { privateKey: cfg.privateKey as `0x${string}` } : {}),
    ...(cfg.rpcUrl !== undefined ? { rpcUrlOverride: cfg.rpcUrl } : {}),
  });
};

/** Build a fully-wired `AcpAgent` for write commands. Tests pass a custom
 * `providerFactory` so the on-chain calls never reach Alchemy. */
export function buildAcpAgent(
  cfg: CliConfig,
  providerFactory: ProviderFactory = defaultProviderFactory,
  fetchImpl?: typeof fetch
): AcpAgent {
  const provider = providerFactory(cfg);
  const contracts = resolveContracts(cfg);
  return new AcpAgent({
    provider,
    gatewayUrl: cfg.gatewayUrl,
    contracts,
    gateway: buildGatewayClient(cfg, fetchImpl),
  });
}

/** Re-exported for tests that want to assert the chain map without
 * pulling private state out of `buildAcpAgent`. */
export const __testing__ = { CHAIN_DEFAULTS };
