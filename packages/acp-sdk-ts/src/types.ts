// Wire types mirroring services/gateway-go/docs/openapi.json (v0.1.0-alpha).
//
// Values that originate as uint256 on chain are encoded as decimal strings on
// the wire — never coerce them to `number`. Addresses are lowercase 0x-hex.
//
// This file is hand-rolled rather than codegen'd to keep the SDK dependency
// graph minimal; the openapi-check CI gate (`gateway-openapi-check` workflow)
// catches drift in the upstream spec, and the `roundtrip` integration tests
// in M4-acp-sdk-e2e (#96) catch drift on the consumer side.

/** Hex-encoded EVM address, lowercase, `0x`-prefixed (20 bytes). */
export type EvmAddress = `0x${string}`;

/** Hex-encoded 32-byte transaction hash, `0x`-prefixed. */
export type TxHash = `0x${string}`;

/** Decimal-string encoding of a uint256 (no `0x` prefix). */
export type BigIntString = string;

/** RFC 3339 timestamp string. */
export type IsoTimestamp = string;

/** Discriminant matching the gateway's `store.AgentKind` enum. */
export type AgentKind = "launchpad" | "acp";

/** Discriminant matching the gateway's `store.JobPhase` enum. */
export type JobPhase =
  | "created"
  | "funded"
  | "active"
  | "completed"
  | "cancelled"
  | "disputed"
  | "released"
  | "resolved";

/** Mirror of `store.Agent` — fields tagged `optional` in the spec are optional here. */
export interface Agent {
  id: number;
  agent_id: BigIntString;
  kind: AgentKind;
  controller?: EvmAddress;
  creator?: EvmAddress;
  capabilities: BigIntString;
  metadata_uri?: string;
  block_number: number;
  log_index: number;
  tx_hash: TxHash;
  created_at: IsoTimestamp;
  updated_at: IsoTimestamp;
}

/** Mirror of `store.Job`. */
export interface Job {
  id: number;
  on_chain_id?: string;
  agent_pk: number;
  agent_id: BigIntString;
  phase: JobPhase;
  job_type: number;
  principal: EvmAddress;
  clone_address?: EvmAddress;
  token: EvmAddress;
  budget: BigIntString;
  deadline: IsoTimestamp;
  result_uri?: string;
  block_number: number;
  log_index: number;
  tx_hash: TxHash;
  created_at: IsoTimestamp;
  updated_at: IsoTimestamp;
}

/** Mirror of `openapi.AgentsPage`. */
export interface AgentsPage {
  data: Agent[];
  next_cursor: string;
}

/** Mirror of `openapi.JobsPage`. */
export interface JobsPage {
  data: Job[];
  next_cursor: string;
}

/** Mirror of `store.Stats`. */
export interface Stats {
  total_agents: number;
  total_jobs: number;
  total_trades: number;
  last_block_indexed: number;
}

/** Auth: shape returned by `/auth/siwe/nonce` and `/auth/siws/nonce`. */
export interface AuthNonceResponse {
  nonce: string;
  expires_at: number;
}

/** Auth: payload for `/auth/siwe/verify` and `/auth/siws/verify`. */
export interface AuthVerifyRequest {
  message: string;
  signature: string;
}

/** Auth: response from `/auth/(siwe|siws)/verify`. */
export interface AuthVerifyResponse {
  address: string;
  expires_at: number;
  /** Bearer token to use as `Authorization: Bearer <token>` on subsequent requests. */
  token?: string;
}

/** Shape of the gateway error envelope. */
export interface GatewayErrorBody {
  code: string;
  message: string;
}

export interface GatewayErrorEnvelope {
  error: GatewayErrorBody;
}

/**
 * Direct vs evaluated job, mirroring `IJob.JobType` on the contracts side.
 * Uses an explicit numeric encoding so the SDK serialises the same uint8 the
 * `JobFactory.createJob` ABI expects.
 */
export const JobType = {
  Direct: 0,
  Evaluated: 1,
} as const;
export type JobTypeValue = (typeof JobType)[keyof typeof JobType];

/** Parameters accepted by `AcpAgent.register`. */
export interface RegisterParams {
  /**
   * Address that will control the agent on chain. If omitted, the SDK uses the
   * connected provider's account address.
   */
  controller?: EvmAddress;
  /** Non-empty URI pointing to off-chain agent metadata (typically `ipfs://...`). */
  metadataUri: string;
  /**
   * Packed capability bitmap (consumer-defined). Accepts both `bigint` (preferred,
   * lossless) and decimal strings; passing `number` is allowed for small bitmaps
   * but disallowed above 2^53 to prevent silent precision loss.
   */
  capabilities: bigint | BigIntString | number;
}

/** Result returned by `AcpAgent.register`. */
export interface RegisterResult {
  agentId: bigint;
  txHash: TxHash;
  controller: EvmAddress;
  metadataUri: string;
  capabilities: bigint;
}

/** Filter accepted by `AcpAgent.browse`. */
export interface BrowseParams {
  kind?: AgentKind;
  /** Page size; the gateway clamps to `[1, 200]` and defaults to 50. */
  limit?: number;
  /** Opaque pagination cursor returned in the previous page's `next_cursor`. */
  cursor?: string;
}

/** Parameters accepted by `AcpAgent.createJob`. */
export interface CreateJobParams {
  /** Agent registry id to target; pass `0n` to leave it open to any active agent. */
  targetAgentId?: bigint;
  /** ERC-20 payment token. */
  token: EvmAddress;
  /** Amount to lock into Escrow (in token base units). */
  budget: bigint;
  /** UNIX seconds after which the job may be expired. Must be in the future. */
  deadline: bigint;
  /** Direct or Evaluated. */
  jobType?: JobTypeValue;
  /** Required iff `jobType === JobType.Evaluated`; ignored otherwise. */
  evaluator?: EvmAddress;
  /** Optional arbiter override; falls back to the factory's default arbiter. */
  arbiter?: EvmAddress;
  /** Optional list of pre-approved hook addresses. */
  hooks?: EvmAddress[];
}

/** Result returned by `AcpAgent.createJob`. */
export interface CreateJobResult {
  jobId: bigint;
  cloneAddress: EvmAddress;
  txHash: TxHash;
}
