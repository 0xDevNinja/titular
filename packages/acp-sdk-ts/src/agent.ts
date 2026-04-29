// AcpAgent — the public, top-level SDK class.
//
// Surface:
//   - `register(meta)`     -> AgentRegistry.register on chain, returns the new agentId
//   - `browse(filter)`     -> GET /api/v1/agents on the gateway, returns AgentsPage
//   - `createJob(params)`  -> JobFactory.createJob on chain, returns jobId + clone
//
// AcpAgent never imports viem at the type-API surface; everything talks to
// chain through {@link EvmProviderAdapter} and to the gateway through
// {@link GatewayClient}. That keeps the public surface stable across provider
// implementations (Alchemy default, Coinbase Smart Wallet, EIP-1193, mocks).

import { AGENT_REGISTRY_ABI, JOB_FACTORY_ABI } from "./abi.js";
import { AcpError } from "./errors.js";
import { GatewayClient, type GatewayClientParams } from "./gateway-client.js";
import { decodeFirstEvent } from "./log-decoder.js";
import type { EvmProviderAdapter } from "./providers/types.js";
import {
  type AgentsPage,
  type BrowseParams,
  type CreateJobParams,
  type CreateJobResult,
  type EvmAddress,
  JobType,
  type RegisterParams,
  type RegisterResult,
  type TxHash,
} from "./types.js";

/** On-chain contract addresses required by AcpAgent. */
export interface AcpContractAddresses {
  agentRegistry: EvmAddress;
  jobFactory: EvmAddress;
}

/** Constructor parameters for {@link AcpAgent}. */
export interface AcpAgentParams {
  /** EVM provider adapter (signer + RPC). */
  provider: EvmProviderAdapter;
  /** Gateway base URL (e.g. `https://api.titular.dev`). */
  gatewayUrl: string;
  /** AgentRegistry + JobFactory addresses for the target chain. */
  contracts: AcpContractAddresses;
  /** Optional pre-built gateway client; mostly used in tests. */
  gateway?: GatewayClient;
  /** Extra options forwarded to the auto-built gateway client. */
  gatewayClientOptions?: Omit<GatewayClientParams, "gatewayUrl">;
}

const ZERO_ADDRESS: EvmAddress = "0x0000000000000000000000000000000000000000";
const TWO_53 = BigInt(Number.MAX_SAFE_INTEGER);

/**
 * High-level entry point used by 3rd-party developers building agents on the
 * Titular ACP. Combines the on-chain AgentRegistry + JobFactory with the
 * read-only gateway REST surface behind a single ergonomic API.
 */
export class AcpAgent {
  readonly provider: EvmProviderAdapter;
  readonly gateway: GatewayClient;
  readonly contracts: AcpContractAddresses;

  constructor(params: AcpAgentParams) {
    if (params.contracts.agentRegistry === ZERO_ADDRESS) {
      throw new AcpError("invalid_param", "AcpAgent: contracts.agentRegistry cannot be zero");
    }
    if (params.contracts.jobFactory === ZERO_ADDRESS) {
      throw new AcpError("invalid_param", "AcpAgent: contracts.jobFactory cannot be zero");
    }

    this.provider = params.provider;
    this.contracts = params.contracts;
    this.gateway =
      params.gateway ??
      new GatewayClient({
        gatewayUrl: params.gatewayUrl,
        ...params.gatewayClientOptions,
      });
  }

  /**
   * Register a new ACP agent on chain.
   *
   * Submits `AgentRegistry.register(controller, metadataURI, capabilities)`,
   * waits for inclusion, and decodes the resulting `AgentRegistered` event.
   *
   * @throws {AcpError} `invalid_param` if `metadataUri` is empty.
   * @throws {AcpError} `tx_reverted` if the transaction reverts.
   * @throws {AcpError} `event_not_found` if the receipt is missing the event.
   */
  async register(meta: RegisterParams): Promise<RegisterResult> {
    if (meta.metadataUri.length === 0) {
      throw new AcpError("invalid_param", "register: `metadataUri` cannot be empty");
    }
    const capabilities = normalizeBigInt(meta.capabilities, "capabilities");
    const controller = (
      meta.controller ?? (await this.provider.getAddress())
    ).toLowerCase() as EvmAddress;

    const txHash = await this.provider.writeContract({
      address: this.contracts.agentRegistry,
      abi: AGENT_REGISTRY_ABI,
      functionName: "register",
      args: [controller, meta.metadataUri, capabilities],
    });

    const receipt = await this.provider.waitForReceipt(txHash);
    assertSuccess(receipt.status, txHash);

    const decoded = decodeFirstEvent<{
      agentId: bigint;
      controller: `0x${string}`;
      metadataURI: string;
      capabilities: bigint;
    }>(AGENT_REGISTRY_ABI, "AgentRegistered", receipt.logs);
    if (decoded === null) {
      throw new AcpError(
        "event_not_found",
        `register: AgentRegistered event missing from receipt for tx ${txHash}`
      );
    }

    return {
      agentId: decoded.agentId,
      txHash,
      controller: decoded.controller.toLowerCase() as EvmAddress,
      metadataUri: decoded.metadataURI,
      capabilities: decoded.capabilities,
    };
  }

  /**
   * Browse the agent registry through the gateway. This is a read-only call
   * against the gateway's Postgres-backed read model and does not require
   * authentication or a connected signer.
   */
  async browse(filter: BrowseParams = {}): Promise<AgentsPage> {
    if (filter.limit !== undefined) {
      if (!Number.isInteger(filter.limit) || filter.limit < 1 || filter.limit > 200) {
        throw new AcpError("invalid_param", "browse: `limit` must be an integer in [1, 200]");
      }
    }
    return this.gateway.listAgents(filter);
  }

  /**
   * Create a new ACP job on chain.
   *
   * Submits `JobFactory.createJob(...)` with the provided params, waits for
   * inclusion, and decodes the resulting `JobCreated` event. The caller MUST
   * have already approved the JobFactory for `params.budget` of `params.token`
   * (the SDK does not silently move ERC-20 allowance for the user).
   *
   * @throws {AcpError} `invalid_param` for zero budget, past deadline, or
   *   `Evaluated` jobType without an `evaluator`.
   */
  async createJob(params: CreateJobParams): Promise<CreateJobResult> {
    if (params.token === ZERO_ADDRESS) {
      throw new AcpError("invalid_param", "createJob: `token` cannot be zero");
    }
    if (params.budget <= 0n) {
      throw new AcpError("invalid_param", "createJob: `budget` must be > 0");
    }
    if (params.deadline <= 0n) {
      throw new AcpError(
        "invalid_param",
        "createJob: `deadline` must be a positive UNIX timestamp"
      );
    }
    const jobType = params.jobType ?? JobType.Direct;
    if (jobType === JobType.Evaluated) {
      if (params.evaluator === undefined || params.evaluator === ZERO_ADDRESS) {
        throw new AcpError(
          "invalid_param",
          "createJob: `evaluator` must be a non-zero address when jobType === Evaluated"
        );
      }
    }
    const evaluator = params.evaluator ?? ZERO_ADDRESS;
    const arbiter = params.arbiter ?? ZERO_ADDRESS;
    const hooks = params.hooks ?? [];
    const targetAgentId = params.targetAgentId ?? 0n;

    const txHash = await this.provider.writeContract({
      address: this.contracts.jobFactory,
      abi: JOB_FACTORY_ABI,
      functionName: "createJob",
      args: [
        targetAgentId,
        params.token,
        params.budget,
        params.deadline,
        jobType,
        evaluator,
        arbiter,
        hooks,
      ],
    });

    const receipt = await this.provider.waitForReceipt(txHash);
    assertSuccess(receipt.status, txHash);

    const decoded = decodeFirstEvent<{
      jobId: bigint;
      clone: `0x${string}`;
      principal: `0x${string}`;
      jobType: number;
      token: `0x${string}`;
      budget: bigint;
    }>(JOB_FACTORY_ABI, "JobCreated", receipt.logs);
    if (decoded === null) {
      throw new AcpError(
        "event_not_found",
        `createJob: JobCreated event missing from receipt for tx ${txHash}`
      );
    }

    return {
      jobId: decoded.jobId,
      cloneAddress: decoded.clone.toLowerCase() as EvmAddress,
      txHash,
    };
  }
}

/**
 * Normalise a capability bitmap to a `bigint`. Accepts `bigint`, decimal
 * strings, and `number` (rejected above 2^53 to avoid silent precision loss).
 */
function normalizeBigInt(value: bigint | string | number, label: string): bigint {
  if (typeof value === "bigint") return value;
  if (typeof value === "number") {
    if (!Number.isInteger(value) || value < 0) {
      throw new AcpError("invalid_param", `${label}: must be a non-negative integer`);
    }
    const big = BigInt(value);
    if (big > TWO_53) {
      throw new AcpError(
        "invalid_param",
        `${label}: value > 2^53 must be passed as a bigint or decimal string to avoid precision loss`
      );
    }
    return big;
  }
  if (!/^\d+$/.test(value)) {
    throw new AcpError("invalid_param", `${label}: decimal-string form must match /^\\d+$/`);
  }
  return BigInt(value);
}

function assertSuccess(status: "success" | "reverted", txHash: TxHash): void {
  if (status !== "success") {
    throw new AcpError("tx_reverted", `transaction reverted: ${txHash}`);
  }
}
