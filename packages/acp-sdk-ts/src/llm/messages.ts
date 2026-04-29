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

import type { Agent, Job } from "../types.js";
import type { AgentState, Message } from "./types.js";

const DEFAULT_SYSTEM_PROMPT =
  "You are an autonomous agent operating on the Titular ACP. The next message " +
  "contains a JSON snapshot of the on-chain registry and your active jobs. Use " +
  "the provided tools to act on this state; do not invent agentIds, jobIds, or " +
  "addresses that do not appear in the snapshot. All numeric ids are decimal-" +
  "encoded uint256 strings — pass them through verbatim.";

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
 */
export function toMessages(state: AgentState): Message[] {
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
    if (hasAgents) payload.agents = (state.agents ?? []).map(serializeAgent);
    if (hasJobs) payload.jobs = (state.jobs ?? []).map(serializeJob);
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
 */
function serializeAgent(a: Agent): SerializedAgent {
  const out: SerializedAgent = {
    agent_id: a.agent_id,
    kind: a.kind,
    capabilities: a.capabilities,
  };
  if (a.controller !== undefined) out.controller = a.controller;
  if (a.metadata_uri !== undefined) out.metadata_uri = a.metadata_uri;
  return out;
}

function serializeJob(j: Job): SerializedJob {
  const out: SerializedJob = {
    agent_id: j.agent_id,
    phase: j.phase,
    job_type: j.job_type,
    principal: j.principal,
    token: j.token,
    budget: j.budget,
    deadline: j.deadline,
  };
  if (j.clone_address !== undefined) out.clone_address = j.clone_address;
  if (j.result_uri !== undefined) out.result_uri = j.result_uri;
  return out;
}
