// Provider adapter interface used by AcpAgent. The SDK ships an
// AlchemyEvmProviderAdapter (see ./alchemy.ts) but consumers can swap in their
// own adapter — Coinbase Smart Wallet, MetaMask, viem WalletClient, etc — by
// implementing this interface.

import type { EvmAddress, TxHash } from "../types.js";

/**
 * Minimal contract-call shape. Mirrors viem's `WriteContractParameters` but
 * keeps the surface small so the adapter doesn't have to depend on viem
 * directly. `args` is unknown[] rather than `readonly unknown[]` so consumers
 * can pass mutable arrays without a cast.
 */
export interface ContractWriteParams {
  /** Target contract address. */
  address: EvmAddress;
  /** Human-readable ABI fragment describing only the function being called. */
  abi: readonly object[];
  /** Function name to invoke. */
  functionName: string;
  /** Positional arguments matching `abi[functionName]`. */
  args: readonly unknown[];
  /** Optional ETH value to attach (wei). */
  value?: bigint;
}

/**
 * Minimal read-only contract-call shape. Same conventions as
 * {@link ContractWriteParams} but for `eth_call`.
 */
export interface ContractReadParams {
  address: EvmAddress;
  abi: readonly object[];
  functionName: string;
  args?: readonly unknown[];
}

/** Decoded log returned by waitForReceipt. */
export interface ParsedLog {
  address: EvmAddress;
  topics: readonly `0x${string}`[];
  data: `0x${string}`;
  logIndex: number;
}

/** Minimal receipt returned by waitForReceipt. */
export interface TxReceipt {
  status: "success" | "reverted";
  transactionHash: TxHash;
  blockNumber: bigint;
  logs: ParsedLog[];
}

/**
 * EVM provider adapter consumed by AcpAgent.
 *
 * Implementations MUST:
 *   - return the connected account address from {@link getAddress}
 *   - submit transactions via {@link writeContract} and return the tx hash
 *   - block until inclusion via {@link waitForReceipt}, returning logs in
 *     emission order so AcpAgent can decode `AgentRegistered` / `JobCreated`
 *
 * Implementations MUST NOT:
 *   - cache nonces across calls (the underlying client should handle that)
 *   - swallow or rewrap revert reasons — surface them so the SDK can map them
 *     to its `AcpError` taxonomy
 */
export interface EvmProviderAdapter {
  /** Numeric chain id the adapter is connected to. */
  getChainId(): Promise<number>;

  /** Address of the connected signer / smart account. */
  getAddress(): Promise<EvmAddress>;

  /** Read-only `eth_call`. */
  readContract<T = unknown>(params: ContractReadParams): Promise<T>;

  /** Submit a state-changing transaction; return the tx hash. */
  writeContract(params: ContractWriteParams): Promise<TxHash>;

  /** Block until `txHash` is included; return the parsed receipt. */
  waitForReceipt(txHash: TxHash): Promise<TxReceipt>;

  /**
   * Sign an arbitrary message string for SIWE/SIWS auth against the gateway.
   * Optional — adapters that cannot sign (e.g. read-only RPCs) may omit this,
   * but `AcpAgent.authenticate` will throw `provider_cannot_sign` if called.
   */
  signMessage?(message: string): Promise<`0x${string}`>;
}
