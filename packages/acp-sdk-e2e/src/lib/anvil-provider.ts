// `EvmProviderAdapter` backed by viem's `walletClient` + `publicClient`,
// pointed at an in-process anvil daemon. Used by the e2e suite to drive the
// SDK's on-chain code paths (register, createJob) against a real EVM with
// real receipts and real log decoding.
//
// We do NOT add a viem-based adapter to the published SDK in this PR —
// `@titular/acp-sdk` keeps its single shipping adapter (Alchemy) to avoid
// pulling viem into every consumer's runtime graph. The adapter here is
// test-only and lives behind the e2e package boundary.

import type {
  ContractReadParams,
  ContractWriteParams,
  EvmProviderAdapter,
  ParsedLog,
  TxReceipt,
} from "@titular/acp-sdk";
import type { EvmAddress, TxHash } from "@titular/acp-sdk";
import {
  http,
  type Account,
  type Chain,
  type PublicClient,
  type WalletClient,
  createPublicClient,
  createWalletClient,
  defineChain,
} from "viem";
import { privateKeyToAccount } from "viem/accounts";
import { ANVIL_CHAIN_ID, ANVIL_DEV_PRIVATE_KEY } from "./anvil.js";

/** Local anvil chain definition shared across the suite. */
export const anvilChain: Chain = defineChain({
  id: ANVIL_CHAIN_ID,
  name: "anvil-local",
  nativeCurrency: { name: "Ether", symbol: "ETH", decimals: 18 },
  rpcUrls: {
    default: { http: ["http://127.0.0.1"] },
  },
});

/**
 * Build an `EvmProviderAdapter` for the SDK that signs with anvil's account 0
 * and submits over HTTP RPC at `rpcUrl`. The adapter implements the same
 * contract `AlchemyEvmProviderAdapter` does in production code — same return
 * shapes, same lowercase-address invariant, same revert mapping.
 */
export function newAnvilProvider(rpcUrl: string): EvmProviderAdapter {
  const account: Account = privateKeyToAccount(ANVIL_DEV_PRIVATE_KEY);

  const publicClient: PublicClient = createPublicClient({
    chain: anvilChain,
    transport: http(rpcUrl),
  });
  const walletClient: WalletClient = createWalletClient({
    account,
    chain: anvilChain,
    transport: http(rpcUrl),
  });

  return {
    async getChainId(): Promise<number> {
      return publicClient.getChainId();
    },

    async getAddress(): Promise<EvmAddress> {
      return account.address.toLowerCase() as EvmAddress;
    },

    async readContract<T = unknown>(params: ContractReadParams): Promise<T> {
      const result = await publicClient.readContract({
        address: params.address,
        // viem's typings demand a narrow Abi; the SDK passes a wider readonly
        // object[] to keep adapters viem-agnostic. The cast is sound because
        // viem only inspects `abi` shape at runtime and the SDK guarantees
        // the fragment matches `functionName`.
        abi: params.abi as never,
        functionName: params.functionName,
        args: params.args as readonly unknown[] | undefined,
      });
      return result as T;
    },

    async writeContract(params: ContractWriteParams): Promise<TxHash> {
      const hash = await walletClient.writeContract({
        account,
        chain: anvilChain,
        address: params.address,
        abi: params.abi as never,
        functionName: params.functionName,
        args: params.args as readonly unknown[],
        value: params.value,
      });
      return hash as TxHash;
    },

    async waitForReceipt(txHash: TxHash): Promise<TxReceipt> {
      const receipt = await publicClient.waitForTransactionReceipt({
        hash: txHash,
        // Anvil mines on submission with auto-mining on; one confirmation is
        // already final on a single-node devnet. Larger numbers would only
        // add latency on the test runner.
        confirmations: 1,
      });
      const logs: ParsedLog[] = receipt.logs.map((log) => ({
        address: log.address.toLowerCase() as EvmAddress,
        topics: log.topics as readonly `0x${string}`[],
        data: log.data,
        logIndex: Number(log.logIndex),
      }));
      return {
        status: receipt.status === "success" ? "success" : "reverted",
        transactionHash: receipt.transactionHash as TxHash,
        blockNumber: receipt.blockNumber,
        logs,
      };
    },

    async signMessage(message: string): Promise<`0x${string}`> {
      return walletClient.signMessage({ account, message });
    },
  };
}
