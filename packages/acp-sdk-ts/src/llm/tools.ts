// availableTools — emit OpenAI-compatible tool definitions for every action
// AcpAgent currently exposes (register, browse, createJob).
//
// Tool names use the `acp_` prefix so models can disambiguate ACP actions from
// any other tools the host application registers.
//
// Schemas are hand-written (rather than generated from Zod or the ABI) because
// the LLM-facing surface is intentionally narrower than the on-chain ABI:
//   - `capabilities` accepts decimal strings only (BigInts can't cross the
//     JSON boundary; numbers > 2^53 lose precision)
//   - `targetAgentId`, `budget`, `deadline` likewise
//   - addresses are constrained to the `0x[0-9a-fA-F]{40}` pattern
//   - jobType is exposed as the literal strings the user thinks in
//     ("direct" | "evaluated") rather than the uint8 the contract expects
//
// The dispatcher in `executor.ts` does the inverse mapping (string -> uint8,
// decimal -> bigint) before calling AcpAgent.

import type { JsonSchema, ToolDefinition } from "./types.js";

const ADDRESS_PATTERN = "^0x[0-9a-fA-F]{40}$";
const DECIMAL_UINT_PATTERN = "^(0|[1-9][0-9]*)$";

const ADDRESS_SCHEMA: JsonSchema = {
  type: "string",
  description: "EVM address as a `0x`-prefixed 20-byte hex string.",
  pattern: ADDRESS_PATTERN,
};

const DECIMAL_UINT_SCHEMA: JsonSchema = {
  type: "string",
  description:
    "Decimal-encoded uint256 (no `0x` prefix, no scientific notation). Required for any " +
    "numeric value that may exceed 2^53.",
  pattern: DECIMAL_UINT_PATTERN,
};

/** Names of every tool emitted by {@link availableTools}. Stable; treat as enum. */
export const TOOL_NAMES = {
  register: "acp_register",
  browse: "acp_browse",
  createJob: "acp_create_job",
} as const;
export type ToolName = (typeof TOOL_NAMES)[keyof typeof TOOL_NAMES];

const REGISTER_TOOL: ToolDefinition = {
  type: "function",
  function: {
    name: TOOL_NAMES.register,
    description:
      "Register a new ACP agent on chain via AgentRegistry.register. Returns the new agentId, " +
      "tx hash, and the controller address that owns the agent. The connected signer pays gas.",
    parameters: {
      type: "object",
      additionalProperties: false,
      required: ["metadataUri", "capabilities"],
      properties: {
        controller: {
          ...ADDRESS_SCHEMA,
          description:
            "Optional. Address that will control the agent on chain. Defaults to the connected " +
            "signer's address.",
        },
        metadataUri: {
          type: "string",
          description:
            "Non-empty URI pointing to off-chain agent metadata. Typically `ipfs://<cid>` or " +
            "`https://...`.",
        },
        capabilities: {
          ...DECIMAL_UINT_SCHEMA,
          description:
            'Packed capability bitmap as a decimal-encoded uint256 (e.g. `"5"` for caps 1+4).',
        },
      },
    },
  },
};

const BROWSE_TOOL: ToolDefinition = {
  type: "function",
  function: {
    name: TOOL_NAMES.browse,
    description:
      "Browse the on-chain agent registry through the gateway. Read-only; does not require a " +
      "signer. Returns a paginated AgentsPage with `data` and `next_cursor`.",
    parameters: {
      type: "object",
      additionalProperties: false,
      properties: {
        kind: {
          type: "string",
          enum: ["launchpad", "acp"],
          description: "Filter by agent kind. Omit to return both.",
        },
        limit: {
          type: "integer",
          minimum: 1,
          maximum: 200,
          description: "Page size in [1, 200]. Defaults to 50 server-side.",
        },
        cursor: {
          type: "string",
          description:
            "Opaque pagination cursor returned in the previous page's `next_cursor`. Omit for " +
            "the first page.",
        },
      },
    },
  },
};

const CREATE_JOB_TOOL: ToolDefinition = {
  type: "function",
  function: {
    name: TOOL_NAMES.createJob,
    description:
      "Create a new ACP job on chain via JobFactory.createJob. The connected signer must have " +
      "already approved the JobFactory for `budget` of `token` (the SDK does not move ERC-20 " +
      "allowance on its own). Returns the new jobId, the deployed clone address, and the tx " +
      "hash.",
    parameters: {
      type: "object",
      additionalProperties: false,
      required: ["token", "budget", "deadline"],
      properties: {
        targetAgentId: {
          ...DECIMAL_UINT_SCHEMA,
          description:
            'Optional agent registry id to target. Pass `"0"` (or omit) to leave the job ' +
            "open to any active agent.",
        },
        token: {
          ...ADDRESS_SCHEMA,
          description: "ERC-20 payment token address. Must not be the zero address.",
        },
        budget: {
          ...DECIMAL_UINT_SCHEMA,
          description: "Amount to lock into Escrow (in token base units). Must be > 0.",
        },
        deadline: {
          ...DECIMAL_UINT_SCHEMA,
          description: "UNIX seconds after which the job may be expired. Must be in the future.",
        },
        jobType: {
          type: "string",
          enum: ["direct", "evaluated"],
          description:
            "Job type. `evaluated` requires `evaluator`. Defaults to `direct` when omitted.",
        },
        evaluator: {
          ...ADDRESS_SCHEMA,
          description:
            "Required iff `jobType` is `evaluated`; ignored otherwise. Must be a non-zero " +
            "address.",
        },
        arbiter: {
          ...ADDRESS_SCHEMA,
          description: "Optional arbiter override; falls back to the factory's default arbiter.",
        },
        hooks: {
          type: "array",
          description: "Optional list of pre-approved hook addresses.",
          items: ADDRESS_SCHEMA,
        },
      },
    },
  },
};

// TODO(M5): add tools for AcpAgent.acceptJob, submitResult, release once those
// methods land on the agent surface (issue #96 deferred per acpAgent scope).
const ALL_TOOLS: readonly ToolDefinition[] = Object.freeze([
  REGISTER_TOOL,
  BROWSE_TOOL,
  CREATE_JOB_TOOL,
]);

/** Filter accepted by {@link availableTools}. */
export interface AvailableToolsParams {
  /**
   * If provided, only the listed tools are returned. Useful when the host
   * application wants to expose, say, only read-only tools to a model that
   * shouldn't be allowed to register agents or spend tokens.
   */
  include?: readonly ToolName[];
}

/**
 * Return OpenAI / Anthropic / Gemini-compatible tool definitions for every
 * action {@link import("../agent.js").AcpAgent} currently exposes.
 *
 * The result is a fresh array on every call so callers can mutate it without
 * affecting subsequent calls. Returned tool objects are deeply frozen.
 *
 * @example
 * ```ts
 * const tools = availableTools();
 * const completion = await openai.chat.completions.create({ model, tools, messages });
 * ```
 */
export function availableTools(params: AvailableToolsParams = {}): ToolDefinition[] {
  if (params.include === undefined) return ALL_TOOLS.slice();
  const allow = new Set<string>(params.include);
  return ALL_TOOLS.filter((t) => allow.has(t.function.name));
}
