// CLI-internal types. Public re-exports stay in `index.ts` so the test surface
// can import shared shapes without reaching into command-private modules.

import type { EvmAddress } from "@titular/acp-sdk";

/** Subset of supported chains the CLI knows how to wire. Stays narrow on
 * purpose: the gateway only surfaces these in the alpha. */
export type SupportedChain =
  | "base-sepolia"
  | "base-mainnet"
  | "ethereum-sepolia"
  | "ethereum-mainnet";

/** Persisted CLI configuration, written to `~/.config/titular/config.json`. */
export interface CliConfig {
  /** Base URL of the gateway, e.g. `https://api.titular.dev`. */
  gatewayUrl: string;
  /** Target chain id (slug). */
  chain: SupportedChain;
  /** Optional bearer token minted via `acp login` (SIWE/SIWS). */
  apiKey?: string;
  /** Optional default agent registry address override. */
  agentRegistry?: EvmAddress;
  /** Optional default job factory address override. */
  jobFactory?: EvmAddress;
  /** Optional default RPC URL for write commands; otherwise built from chain. */
  rpcUrl?: string;
  /** Optional 0x-prefixed private key for signing in non-interactive contexts. */
  privateKey?: string;
}

/** Subset of `CliConfig` exposed by `acp configure --print`. The private key
 * is intentionally redacted so a `--json` dump doesn't leak it into shell
 * history or CI logs. */
export interface RedactedConfig extends Omit<CliConfig, "apiKey" | "privateKey"> {
  apiKey?: string;
  privateKey?: string;
}

/** Common per-command flags every subcommand inherits. */
export interface GlobalFlags {
  json?: boolean;
  config?: string;
}
