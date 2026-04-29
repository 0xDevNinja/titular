// Types for the LLM helpers exposed at `@titular/acp-sdk/llm`.
//
// The shapes here intentionally mirror OpenAI's Chat Completions tool / message
// surface (`{type: "function", function: {name, description, parameters}}` for
// tools; role-tagged objects for messages) because Anthropic's `messages` API
// and Gemini's function-calling surface both accept the same shape via official
// adapters today (Anthropic auto-translates OpenAI-format tools; Gemini exposes
// `function_declarations` that map 1:1).
//
// Keeping the public surface vendor-neutral means consumers can pass our output
// straight to `openai.chat.completions.create({tools, messages})` or to
// `anthropic.messages.create({tools, messages})` without an adapter layer.
//
// JSON Schema is hand-written rather than generated from Zod because the SDK
// avoids adding optional runtime deps; the surface is small (a handful of
// tools) and the schemas are checked against the `AcpAgent` method signatures
// in unit tests (`tools.test.ts`, `executor.test.ts`).

import type { Agent, Job } from "../types.js";

/**
 * JSON Schema (draft-07 subset) used in tool `parameters`. Restricted to the
 * subset the OpenAI / Anthropic / Gemini tool-calling adapters all accept.
 */
export interface JsonSchema {
  type?: "object" | "string" | "number" | "integer" | "boolean" | "array" | "null";
  description?: string;
  properties?: Record<string, JsonSchema>;
  required?: string[];
  additionalProperties?: boolean | JsonSchema;
  items?: JsonSchema;
  enum?: readonly (string | number | boolean | null)[];
  /** Inclusive lower bound for numeric types. */
  minimum?: number;
  /** Inclusive upper bound for numeric types. */
  maximum?: number;
  /** Regex pattern (ECMA-262) for string types. */
  pattern?: string;
  /** Minimum item count for array types. */
  minItems?: number;
  /** Maximum item count for array types. */
  maxItems?: number;
}

/**
 * OpenAI-compatible tool definition. Anthropic and Gemini both accept this
 * exact shape (Anthropic since 2024-04, Gemini via the `tools` param on the
 * v1beta `generateContent` endpoint).
 */
export interface ToolDefinition {
  type: "function";
  function: {
    /** Stable, snake_case tool name. Used as the dispatch key in `executeTool`. */
    name: string;
    /** Single-paragraph description shown to the model. */
    description: string;
    /** JSON Schema for the function's argument object. */
    parameters: JsonSchema;
  };
}

/** Role-tagged message. Mirrors `OpenAI.Chat.Completions.ChatCompletionMessageParam`. */
export type Message =
  | { role: "system"; content: string }
  | { role: "user"; content: string }
  | { role: "assistant"; content: string };

/**
 * Snapshot of agent-side state to render into LLM context. All fields are
 * optional so callers can render partial views (e.g. only owned jobs without
 * the full registry).
 */
export interface AgentState {
  /** Optional system-prompt prefix prepended to the message stream. */
  systemPrompt?: string;
  /** Agents the model should know about. Typically the result of `agent.browse(...)`. */
  agents?: readonly Agent[];
  /** Jobs the model should know about. Typically the result of `gateway.listJobs(...)`. */
  jobs?: readonly Job[];
  /**
   * Optional human-side turn appended after the state snapshot. Useful when
   * the SDK is feeding state plus an instruction in a single render.
   */
  userPrompt?: string;
}

/**
 * Tool call as emitted by the model. Matches the OpenAI Chat Completions
 * `tool_calls[]` element shape; consumers typically pass these through
 * verbatim.
 */
export interface ToolCall {
  /** Provider-supplied call id. Echoed back in {@link ToolResult.tool_call_id}. */
  id: string;
  type: "function";
  function: {
    name: string;
    /**
     * JSON-encoded argument string. Most providers stringify even when the
     * model emits a JSON object, so the SDK parses it lazily.
     */
    arguments: string;
  };
}

/**
 * Result of dispatching a {@link ToolCall} via {@link executeTool}. Shape is
 * compatible with the `role: "tool"` message both OpenAI and Anthropic accept.
 */
export interface ToolResult {
  role: "tool";
  /** Echo of {@link ToolCall.id} so the model can correlate the result. */
  tool_call_id: string;
  /** Tool name (echoed for Anthropic compatibility — OpenAI ignores it). */
  name: string;
  /**
   * JSON-encoded payload. On success this is the SDK return value with
   * BigInts coerced to decimal strings. On failure this is `{error: {code, message}}`
   * — the error is returned, not thrown, so the model can recover.
   */
  content: string;
  /** True iff the dispatch raised. */
  isError: boolean;
}
