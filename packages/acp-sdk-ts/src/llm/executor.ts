// executeTool — dispatch a model-emitted {@link ToolCall} to the matching
// {@link AcpAgent} method, marshal the return value back into a tool-result
// envelope the model can consume.
//
// Errors are caught and returned in the result envelope (`isError: true`)
// rather than thrown, because the LLM tool-calling protocols expect every
// `tool_call_id` to be answered exactly once. Throwing here would force the
// caller to invent a synthetic error response for the model, which is both
// boilerplate and easy to get wrong.
//
// All BigInts in the SDK return values are coerced to decimal-encoded strings
// before serialisation — JSON cannot represent BigInt natively, and JS engines
// throw `TypeError: Do not know how to serialize a BigInt` if you try
// `JSON.stringify({n: 1n})`. Decimal strings round-trip cleanly through both
// the OpenAI tool-result protocol and our own `availableTools` schemas.

import type { AcpAgent } from "../agent.js";
import { AcpError } from "../errors.js";
import {
  type BrowseParams,
  type CreateJobParams,
  type EvmAddress,
  JobType,
  type RegisterParams,
} from "../types.js";
import { TOOL_NAMES } from "./tools.js";
import type { ToolCall, ToolResult } from "./types.js";

const ADDRESS_RE = /^0x[0-9a-fA-F]{40}$/;
const DECIMAL_UINT_RE = /^(0|[1-9][0-9]*)$/;

/**
 * Dispatch a single tool call. The return value is always a {@link ToolResult}
 * — failures are encoded as `{isError: true, content: JSON({error: ...})}`
 * rather than thrown.
 */
export async function executeTool(agent: AcpAgent, call: ToolCall): Promise<ToolResult> {
  const name = call.function.name;
  let args: Record<string, unknown>;
  try {
    args = parseArgs(call.function.arguments);
  } catch (err) {
    return errorResult(call, name, "invalid_param", (err as Error).message);
  }

  try {
    switch (name) {
      case TOOL_NAMES.register: {
        const params = parseRegisterArgs(args);
        const result = await agent.register(params);
        return successResult(call, name, {
          agentId: result.agentId.toString(),
          txHash: result.txHash,
          controller: result.controller,
          metadataUri: result.metadataUri,
          capabilities: result.capabilities.toString(),
        });
      }
      case TOOL_NAMES.browse: {
        const params = parseBrowseArgs(args);
        const page = await agent.browse(params);
        return successResult(call, name, page);
      }
      case TOOL_NAMES.createJob: {
        const params = parseCreateJobArgs(args);
        const result = await agent.createJob(params);
        return successResult(call, name, {
          jobId: result.jobId.toString(),
          cloneAddress: result.cloneAddress,
          txHash: result.txHash,
        });
      }
      default:
        return errorResult(call, name, "invalid_param", `unknown tool: ${name}`);
    }
  } catch (err) {
    if (err instanceof AcpError) {
      return errorResult(call, name, err.code, err.message);
    }
    return errorResult(call, name, "invalid_param", (err as Error).message ?? "unknown error");
  }
}

// -- argument parsing -------------------------------------------------------

function parseArgs(raw: string): Record<string, unknown> {
  if (raw === "" || raw === undefined) return {};
  let parsed: unknown;
  try {
    parsed = JSON.parse(raw);
  } catch (cause) {
    throw new Error(`tool arguments are not valid JSON: ${(cause as Error).message}`);
  }
  if (parsed === null || typeof parsed !== "object" || Array.isArray(parsed)) {
    throw new Error("tool arguments must be a JSON object");
  }
  return parsed as Record<string, unknown>;
}

function parseRegisterArgs(args: Record<string, unknown>): RegisterParams {
  const metadataUri = requireString(args, "metadataUri");
  const capabilities = requireDecimalUint(args, "capabilities");
  const params: RegisterParams = { metadataUri, capabilities };
  if (args.controller !== undefined) {
    params.controller = requireAddress(args, "controller");
  }
  return params;
}

function parseBrowseArgs(args: Record<string, unknown>): BrowseParams {
  const params: BrowseParams = {};
  if (args.kind !== undefined) {
    if (args.kind !== "launchpad" && args.kind !== "acp") {
      throw new Error("`kind` must be 'launchpad' or 'acp'");
    }
    params.kind = args.kind;
  }
  if (args.limit !== undefined) {
    if (typeof args.limit !== "number" || !Number.isInteger(args.limit)) {
      throw new Error("`limit` must be an integer");
    }
    params.limit = args.limit;
  }
  if (args.cursor !== undefined) {
    params.cursor = requireString(args, "cursor");
  }
  return params;
}

function parseCreateJobArgs(args: Record<string, unknown>): CreateJobParams {
  const token = requireAddress(args, "token");
  const budget = requireBigUint(args, "budget");
  const deadline = requireBigUint(args, "deadline");
  const params: CreateJobParams = { token, budget, deadline };
  if (args.targetAgentId !== undefined) {
    params.targetAgentId = requireBigUint(args, "targetAgentId");
  }
  if (args.jobType !== undefined) {
    if (args.jobType === "direct") {
      params.jobType = JobType.Direct;
    } else if (args.jobType === "evaluated") {
      params.jobType = JobType.Evaluated;
    } else {
      throw new Error("`jobType` must be 'direct' or 'evaluated'");
    }
  }
  if (args.evaluator !== undefined) {
    params.evaluator = requireAddress(args, "evaluator");
  }
  if (args.arbiter !== undefined) {
    params.arbiter = requireAddress(args, "arbiter");
  }
  if (args.hooks !== undefined) {
    if (!Array.isArray(args.hooks)) {
      throw new Error("`hooks` must be an array of addresses");
    }
    params.hooks = args.hooks.map((h, i) => coerceAddress(h, `hooks[${i}]`));
  }
  return params;
}

// -- coercion primitives ----------------------------------------------------

function requireString(args: Record<string, unknown>, key: string): string {
  const v = args[key];
  if (typeof v !== "string" || v.length === 0) {
    throw new Error(`\`${key}\` must be a non-empty string`);
  }
  return v;
}

function requireAddress(args: Record<string, unknown>, key: string): EvmAddress {
  return coerceAddress(args[key], key);
}

function coerceAddress(v: unknown, key: string): EvmAddress {
  if (typeof v !== "string" || !ADDRESS_RE.test(v)) {
    throw new Error(`\`${key}\` must be a 0x-prefixed 20-byte hex address`);
  }
  return v.toLowerCase() as EvmAddress;
}

function requireDecimalUint(args: Record<string, unknown>, key: string): string {
  const v = args[key];
  if (typeof v !== "string" || !DECIMAL_UINT_RE.test(v)) {
    throw new Error(`\`${key}\` must be a decimal-encoded uint256 string`);
  }
  return v;
}

function requireBigUint(args: Record<string, unknown>, key: string): bigint {
  return BigInt(requireDecimalUint(args, key));
}

// -- result envelope helpers ------------------------------------------------

function successResult(call: ToolCall, name: string, payload: unknown): ToolResult {
  return {
    role: "tool",
    tool_call_id: call.id,
    name,
    content: JSON.stringify(payload),
    isError: false,
  };
}

function errorResult(call: ToolCall, name: string, code: string, message: string): ToolResult {
  return {
    role: "tool",
    tool_call_id: call.id,
    name,
    content: JSON.stringify({ error: { code, message } }),
    isError: true,
  };
}
