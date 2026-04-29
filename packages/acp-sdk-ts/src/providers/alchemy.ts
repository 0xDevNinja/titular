// Default EVM provider adapter — wraps a viem WalletClient + PublicClient
// against an Alchemy-hosted Base RPC. Consumers MAY pass their own
// pre-built clients (e.g. for testing, or to use account-abstraction
// bundlers like @account-kit/aa); see the second constructor overload.
//
// The adapter intentionally keeps the dependency surface to viem only.
// The "Alchemy" naming is about the default RPC URL and key handling, not
// a hard dependency on `alchemy-sdk`.

import {
  http,
  type Account,
  type Address,
  type Chain,
  type Hex,
  type PublicClient,
  type Transport,
  type WalletClient,
  createPublicClient,
  createWalletClient,
  custom,
} from "viem";
import { privateKeyToAccount } from "viem/accounts";
import { base, baseSepolia } from "viem/chains";
import type { EvmAddress, TxHash } from "../types.js";
import type {
  ContractReadParams,
  ContractWriteParams,
  EvmProviderAdapter,
  ParsedLog,
  TxReceipt,
} from "./types.js";

/** Supported chain shorthand. */
export type SupportedChain = "base" | "base-sepolia";

const CHAIN_BY_NAME: Record<SupportedChain, Chain> = {
  base,
  "base-sepolia": baseSepolia,
};

/**
 * Constructor params for {@link AlchemyEvmProviderAdapter}.
 *
 * Mutually exclusive auth modes:
 *   - `apiKey` + `chain`: SDK builds its own viem clients hitting Alchemy.
 *     A signer (`privateKey` or `account`) is required for write methods.
 *   - `transport`: caller provides a viem transport (e.g. EIP-1193 from a
 *     browser wallet). Useful for dapps where the user signs in MetaMask.
 *   - `publicClient` + `walletClient`: caller provides pre-built viem
 *     clients. Used in tests and advanced setups (account-abstraction
 *     bundlers, custom transports, etc).
 */
export type AlchemyEvmProviderAdapterParams =
  | {
      apiKey: string;
      chain: SupportedChain;
      privateKey?: Hex;
      account?: Account;
      rpcUrlOverride?: string;
    }
  | {
      transport: Transport;
      chain: SupportedChain;
      account?: Account;
    }
  | {
      publicClient: PublicClient;
      walletClient?: WalletClient;
    };

/** Thrown when a write method is called on an adapter without a signer. */
export class NoSignerError extends Error {
  override readonly name = "NoSignerError";
  readonly code = "no_signer" as const;
  constructor() {
    super(
      "AlchemyEvmProviderAdapter is read-only: provide a `privateKey`, `account`, " +
        "or pre-built `walletClient` to send transactions."
    );
  }
}

/**
 * Default {@link EvmProviderAdapter} backed by viem.
 *
 * The adapter normalises three setup modes (Alchemy-hosted RPC + private
 * key, browser-wallet EIP-1193, fully custom viem clients) into a single
 * surface and delegates to viem for ABI encoding, gas estimation, and
 * receipt waiting. AcpAgent never imports viem directly; everything goes
 * through this interface so a consumer can swap in their own adapter.
 */
export class AlchemyEvmProviderAdapter implements EvmProviderAdapter {
  readonly #publicClient: PublicClient;
  readonly #walletClient?: WalletClient;

  constructor(params: AlchemyEvmProviderAdapterParams) {
    if ("publicClient" in params) {
      this.#publicClient = params.publicClient;
      if (params.walletClient !== undefined) {
        this.#walletClient = params.walletClient;
      }
      return;
    }

    if ("transport" in params) {
      const chain = CHAIN_BY_NAME[params.chain];
      this.#publicClient = createPublicClient({
        chain,
        transport: params.transport,
      });
      const walletParams: Parameters<typeof createWalletClient>[0] = {
        chain,
        transport: params.transport,
      };
      if (params.account !== undefined) {
        walletParams.account = params.account;
      }
      this.#walletClient = createWalletClient(walletParams);
      return;
    }

    // apiKey path
    const chain = CHAIN_BY_NAME[params.chain];
    const url = params.rpcUrlOverride ?? defaultAlchemyRpcUrl(params.chain, params.apiKey);
    const transport = http(url);
    this.#publicClient = createPublicClient({ chain, transport });

    const account =
      params.account ?? (params.privateKey ? privateKeyToAccount(params.privateKey) : undefined);
    if (account !== undefined) {
      this.#walletClient = createWalletClient({ account, chain, transport });
    }
  }

  async getChainId(): Promise<number> {
    return this.#publicClient.getChainId();
  }

  async getAddress(): Promise<EvmAddress> {
    const wallet = this.#requireWallet();
    const account = wallet.account;
    if (account === undefined) {
      const [addr] = await wallet.getAddresses();
      if (addr === undefined) throw new NoSignerError();
      return addr.toLowerCase() as EvmAddress;
    }
    return account.address.toLowerCase() as EvmAddress;
  }

  async readContract<T = unknown>(params: ContractReadParams): Promise<T> {
    const result = await this.#publicClient.readContract({
      address: params.address as Address,
      abi: params.abi as readonly object[],
      functionName: params.functionName,
      args: params.args ?? [],
      // viem's `readContract` returns `unknown` for arbitrary ABIs; the SDK
      // is responsible for matching the call to the expected return shape.
      // biome-ignore lint/suspicious/noExplicitAny: heterogeneous ABI surface
    } as any);
    return result as T;
  }

  async writeContract(params: ContractWriteParams): Promise<TxHash> {
    const wallet = this.#requireWallet();
    const txParams = {
      address: params.address as Address,
      abi: params.abi as readonly object[],
      functionName: params.functionName,
      args: params.args,
      ...(params.value !== undefined ? { value: params.value } : {}),
      // biome-ignore lint/suspicious/noExplicitAny: heterogeneous ABI surface
    } as any;
    // walletClient.writeContract requires `account` + `chain`; viem will pull
    // them from the configured wallet client when omitted.
    const hash = await wallet.writeContract(txParams);
    return hash as TxHash;
  }

  async waitForReceipt(txHash: TxHash): Promise<TxReceipt> {
    const r = await this.#publicClient.waitForTransactionReceipt({
      hash: txHash,
    });
    const logs: ParsedLog[] = r.logs.map((log) => ({
      address: log.address.toLowerCase() as EvmAddress,
      topics: log.topics,
      data: log.data,
      logIndex: log.logIndex ?? 0,
    }));
    return {
      status: r.status === "success" ? "success" : "reverted",
      transactionHash: r.transactionHash as TxHash,
      blockNumber: r.blockNumber,
      logs,
    };
  }

  async signMessage(message: string): Promise<Hex> {
    const wallet = this.#requireWallet();
    const account = wallet.account;
    if (account === undefined) throw new NoSignerError();
    return wallet.signMessage({ account, message });
  }

  #requireWallet(): WalletClient {
    if (this.#walletClient === undefined) throw new NoSignerError();
    return this.#walletClient;
  }
}

/**
 * Build the default Alchemy RPC URL for the given chain. Exported so tests
 * can assert URL composition without pulling in the network layer.
 */
export function defaultAlchemyRpcUrl(chain: SupportedChain, apiKey: string): string {
  if (apiKey.length === 0) {
    throw new Error("AlchemyEvmProviderAdapter: `apiKey` cannot be empty");
  }
  const subdomain = chain === "base" ? "base-mainnet" : "base-sepolia";
  return `https://${subdomain}.g.alchemy.com/v2/${apiKey}`;
}

/** Re-export viem's `custom` transport for adapters that need EIP-1193. */
export { custom };
