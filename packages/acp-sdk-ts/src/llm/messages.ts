// toMessages — render an {@link AgentState} snapshot into a role-tagged
// message stream suitable for an OpenAI / Anthropic / Gemini chat completion.
//
// The output starts with a single `system` message describing how to read the
// state, followed by one `user` message containing the JSON payload itself,
// followed (optionally) by a final `user` turn carrying the human-side prompt.
//
// Why JSON inside a user turn rather than a structured `content` array?
// Because all three providers accept plain string content while only OpenAI
// (currently) accepts the multi-part `content: [...]` shape. JSON-in-string
// keeps the surface portable and round-trippable through any chat completion
// API. Models reliably parse fenced JSON blocks; we wrap the payload in
// `\`\`\`json ... \`\`\`` so the model treats it as data, not prose.

import { keccak256, toBytes } from "viem";
import type { Agent, Job } from "../types.js";
import type { AgentState, Message } from "./types.js";

const DEFAULT_SYSTEM_PROMPT =
  "You are an autonomous agent operating on the Titular ACP. The next message " +
  "contains a JSON snapshot of the on-chain registry and your active jobs. Use " +
  "the provided tools to act on this state; do not invent agentIds, jobIds, or " +
  "addresses that do not appear in the snapshot. All numeric ids are decimal-" +
  "encoded uint256 strings — pass them through verbatim.";

/** Options accepted by {@link toMessages}. */
export interface ToMessagesOptions {
  /**
   * Redact wallet addresses (`controller`, `principal`) before sending to the
   * LLM. Default: `true`.
   *
   * Each address is replaced with a deterministic short token of the form
   * `addr_${first8charsOfKeccak256(addr)}` so the model can still reason about
   * identity equality across messages without seeing the raw wallet.
   *
   * Set to `false` only when the caller has independently established that
   * leaking these addresses to the model provider is acceptable (e.g. a
   * self-hosted model, or an explicit user opt-in).
   */
  redactAddresses?: boolean;
}

/**
 * Render `state` into a role-tagged message stream.
 *
 * The function is synchronous and pure — it does not call the gateway or the
 * chain. Callers are expected to assemble `state` themselves (typically via
 * `agent.browse(...)` for `agents` and `gateway.listJobs(...)` for `jobs`).
 *
 * Output structure:
 *   1. `system` — either {@link AgentState.systemPrompt} or a sensible default.
 *   2. `user`   — JSON-encoded payload of agents + jobs, fenced as ```json ... ```.
 *   3. `user`   — optional, only emitted when {@link AgentState.userPrompt} is set.
 *
 * If `state` carries no agents and no jobs, the JSON payload turn is omitted
 * entirely (so a stand-alone systemPrompt + userPrompt round-trips to a clean
 * 2-message stream).
 *
 * Privacy contract:
 *   - When `opts.redactAddresses` is `true` (default), the `controller` field
 *     on agents and the `principal` field on jobs are replaced with a
 *     deterministic `addr_<8hex>` token. Identity equality is preserved (same
 *     address → same token) across the same call.
 *   - All other on-chain fields (`agent_id`, `chain_id`, `token`,
 *     `clone_address`, `tx_hash`, `metadata_uri`, `result_uri`, etc.) are
 *     fundamentally public and are NOT redacted. Callers who need stricter
 *     redaction must pre-filter `state` themselves.
 *   - The `userPrompt` and `systemPrompt` are passed through verbatim; the
 *     SDK does not attempt to scrub free-form text.
 */
export function toMessages(state: AgentState, opts: ToMessagesOptions = {}): Message[] {
  const redact = opts.redactAddresses ?? true;
  const messages: Message[] = [];

  const systemPrompt =
    state.systemPrompt !== undefined && state.systemPrompt.length > 0
      ? state.systemPrompt
      : DEFAULT_SYSTEM_PROMPT;
  messages.push({ role: "system", content: systemPrompt });

  const hasAgents = state.agents !== undefined && state.agents.length > 0;
  const hasJobs = state.jobs !== undefined && state.jobs.length > 0;

  if (hasAgents || hasJobs) {
    const payload: { agents?: SerializedAgent[]; jobs?: SerializedJob[] } = {};
    if (hasAgents) payload.agents = (state.agents ?? []).map((a) => serializeAgent(a, redact));
    if (hasJobs) payload.jobs = (state.jobs ?? []).map((j) => serializeJob(j, redact));
    const json = JSON.stringify(payload, null, 2);
    messages.push({
      role: "user",
      content: `Current ACP state snapshot:\n\n\`\`\`json\n${json}\n\`\`\``,
    });
  }

  if (state.userPrompt !== undefined && state.userPrompt.length > 0) {
    messages.push({ role: "user", content: state.userPrompt });
  }

  return messages;
}

interface SerializedAgent {
  agent_id: string;
  kind: string;
  controller?: string;
  capabilities: string;
  metadata_uri?: string;
}

interface SerializedJob {
  agent_id: string;
  phase: string;
  job_type: number;
  principal: string;
  clone_address?: string;
  token: string;
  budget: string;
  deadline: string;
  result_uri?: string;
}

/**
 * Project an {@link Agent} down to the fields a model actually needs to make
 * tool-call decisions. We drop the indexer book-keeping columns (`block_number`,
 * `log_index`, `tx_hash`, `created_at`, `updated_at`) so the JSON payload stays
 * compact — every byte counts against the model's context window.
 *
 * When `redact` is `true`, `controller` is replaced with `addr_<8hex>`.
 */
function serializeAgent(a: Agent, redact: boolean): SerializedAgent {
  const out: SerializedAgent = {
    agent_id: a.agent_id,
    kind: a.kind,
    capabilities: a.capabilities,
  };
  if (a.controller !== undefined) {
    out.controller = redact ? redactAddress(a.controller) : a.controller;
  }
  if (a.metadata_uri !== undefined) out.metadata_uri = a.metadata_uri;
  return out;
}

/**
 * Project a {@link Job} down to model-relevant fields. When `redact` is `true`,
 * `principal` is replaced with `addr_<8hex>`. `token` and `clone_address` are
 * left intact — they are protocol-level public data, not user wallets.
 */
function serializeJob(j: Job, redact: boolean): SerializedJob {
  const out: SerializedJob = {
    agent_id: j.agent_id,
    phase: j.phase,
    job_type: j.job_type,
    principal: redact ? redactAddress(j.principal) : j.principal,
    token: j.token,
    budget: j.budget,
    deadline: j.deadline,
  };
  if (j.clone_address !== undefined) out.clone_address = j.clone_address;
  if (j.result_uri !== undefined) out.result_uri = j.result_uri;
  return out;
}

/**
 * Replace a raw EVM address with a short, deterministic token. Same input →
 * same output, so the model can still recognize when two messages reference
 * the same wallet without seeing the wallet itself.
 *
 * Lower-cases the input first so checksum-cased and non-checksum forms collide
 * to the same token.
 */
function redactAddress(addr: string): string {
  const hash = keccak256(toBytes(addr.toLowerCase()));
  // hash is `0x` + 64 hex chars; take the first 8 hex chars after the prefix.
  return `addr_${hash.slice(2, 10)}`;
}
