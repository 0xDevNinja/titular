import { describe, expect, it } from "vitest";
import { toMessages } from "../llm/messages.js";
import type { Agent, Job } from "../types.js";

const AGENT_A: Agent = {
  id: 1,
  agent_id: "1",
  kind: "acp",
  controller: "0x1111111111111111111111111111111111111111",
  capabilities: "5",
  metadata_uri: "ipfs://meta-a",
  block_number: 100,
  log_index: 0,
  tx_hash: "0xaa",
  created_at: "2026-01-01T00:00:00Z",
  updated_at: "2026-01-01T00:00:00Z",
};

const JOB_A: Job = {
  id: 11,
  agent_pk: 1,
  agent_id: "1",
  phase: "active",
  job_type: 0,
  principal: "0x2222222222222222222222222222222222222222",
  clone_address: "0x3333333333333333333333333333333333333333",
  token: "0x4444444444444444444444444444444444444444",
  budget: "1000000",
  deadline: "2026-01-31T00:00:00Z",
  block_number: 200,
  log_index: 1,
  tx_hash: "0xbb",
  created_at: "2026-01-01T00:00:00Z",
  updated_at: "2026-01-01T00:00:00Z",
};

describe("toMessages", () => {
  it("emits a system message with the default prompt when none is provided", () => {
    const messages = toMessages({});
    expect(messages).toHaveLength(1);
    expect(messages[0]?.role).toBe("system");
    expect((messages[0] as { content: string }).content).toMatch(/Titular ACP/);
  });

  it("uses a custom systemPrompt when provided", () => {
    const messages = toMessages({ systemPrompt: "act as a launchpad bot" });
    expect(messages[0]?.role).toBe("system");
    expect((messages[0] as { content: string }).content).toBe("act as a launchpad bot");
  });

  it("falls back to default prompt when systemPrompt is empty", () => {
    const messages = toMessages({ systemPrompt: "" });
    expect((messages[0] as { content: string }).content).toMatch(/Titular ACP/);
  });

  it("emits a JSON state turn when agents are present", () => {
    const messages = toMessages({ agents: [AGENT_A] }, { redactAddresses: false });
    expect(messages).toHaveLength(2);
    expect(messages[1]?.role).toBe("user");
    const content = (messages[1] as { content: string }).content;
    expect(content).toContain("```json");
    const json = extractJson(content);
    expect(json.agents).toHaveLength(1);
    expect(json.agents[0]).toMatchObject({
      agent_id: "1",
      kind: "acp",
      controller: AGENT_A.controller,
      capabilities: "5",
      metadata_uri: "ipfs://meta-a",
    });
    // Indexer-only fields should NOT leak into the model context.
    expect(json.agents[0]?.block_number).toBeUndefined();
    expect(json.agents[0]?.tx_hash).toBeUndefined();
  });

  it("emits jobs in the JSON payload", () => {
    const messages = toMessages({ jobs: [JOB_A] }, { redactAddresses: false });
    expect(messages).toHaveLength(2);
    const json = extractJson((messages[1] as { content: string }).content);
    expect(json.jobs).toHaveLength(1);
    expect(json.jobs[0]).toMatchObject({
      agent_id: "1",
      phase: "active",
      job_type: 0,
      principal: JOB_A.principal,
      clone_address: JOB_A.clone_address,
      token: JOB_A.token,
      budget: "1000000",
    });
    expect(json.jobs[0]?.tx_hash).toBeUndefined();
  });

  it("merges agents + jobs into a single payload turn", () => {
    const messages = toMessages({ agents: [AGENT_A], jobs: [JOB_A] });
    expect(messages).toHaveLength(2);
    const json = extractJson((messages[1] as { content: string }).content);
    expect(json.agents).toHaveLength(1);
    expect(json.jobs).toHaveLength(1);
  });

  it("omits the JSON turn when both agents + jobs are absent or empty", () => {
    const messages = toMessages({ agents: [], jobs: [] });
    expect(messages).toHaveLength(1);
    expect(messages[0]?.role).toBe("system");
  });

  it("appends an optional userPrompt as a final user turn", () => {
    const messages = toMessages({ userPrompt: "what should I do next?" });
    expect(messages).toHaveLength(2);
    expect(messages[1]).toEqual({
      role: "user",
      content: "what should I do next?",
    });
  });

  it("places userPrompt after the JSON state turn", () => {
    const messages = toMessages({
      agents: [AGENT_A],
      userPrompt: "pick the best agent",
    });
    expect(messages.map((m) => m.role)).toEqual(["system", "user", "user"]);
    expect((messages[2] as { content: string }).content).toBe("pick the best agent");
  });

  it("ignores empty userPrompt", () => {
    const messages = toMessages({ userPrompt: "" });
    expect(messages).toHaveLength(1);
  });

  it("drops optional fields cleanly when an agent has no controller / metadata_uri", () => {
    const sparse: Agent = {
      id: AGENT_A.id,
      agent_id: AGENT_A.agent_id,
      kind: AGENT_A.kind,
      capabilities: AGENT_A.capabilities,
      block_number: AGENT_A.block_number,
      log_index: AGENT_A.log_index,
      tx_hash: AGENT_A.tx_hash,
      created_at: AGENT_A.created_at,
      updated_at: AGENT_A.updated_at,
    };
    const messages = toMessages({ agents: [sparse] });
    const json = extractJson((messages[1] as { content: string }).content);
    expect(json.agents[0]?.controller).toBeUndefined();
    expect(json.agents[0]?.metadata_uri).toBeUndefined();
  });

  it("drops optional fields cleanly when a job has no clone_address / result_uri", () => {
    const sparse: Job = {
      id: JOB_A.id,
      agent_pk: JOB_A.agent_pk,
      agent_id: JOB_A.agent_id,
      phase: JOB_A.phase,
      job_type: JOB_A.job_type,
      principal: JOB_A.principal,
      token: JOB_A.token,
      budget: JOB_A.budget,
      deadline: JOB_A.deadline,
      block_number: JOB_A.block_number,
      log_index: JOB_A.log_index,
      tx_hash: JOB_A.tx_hash,
      created_at: JOB_A.created_at,
      updated_at: JOB_A.updated_at,
    };
    const messages = toMessages({ jobs: [sparse] });
    const json = extractJson((messages[1] as { content: string }).content);
    expect(json.jobs[0]?.clone_address).toBeUndefined();
    expect(json.jobs[0]?.result_uri).toBeUndefined();
  });

  it("redacts controller + principal addresses by default", () => {
    const messages = toMessages({ agents: [AGENT_A], jobs: [JOB_A] });
    const json = extractJson((messages[1] as { content: string }).content);

    const controller = json.agents[0]?.controller as string | undefined;
    const principal = json.jobs[0]?.principal as string | undefined;

    // Raw hex addresses MUST NOT appear under the default privacy contract.
    expect(controller).not.toBe(AGENT_A.controller);
    expect(principal).not.toBe(JOB_A.principal);

    // Both should be tokenized to the deterministic short form.
    expect(controller).toMatch(/^addr_[0-9a-f]{8}$/);
    expect(principal).toMatch(/^addr_[0-9a-f]{8}$/);

    // Public protocol fields MUST still flow through unchanged.
    expect(json.jobs[0]?.token).toBe(JOB_A.token);
    expect(json.jobs[0]?.clone_address).toBe(JOB_A.clone_address);
  });

  it("preserves identity equality across messages when redacted", () => {
    // Same wallet appears as both controller (agent) and principal (job).
    const sharedAddr = "0xabababababababababababababababababababab";
    const agent: Agent = { ...AGENT_A, controller: sharedAddr };
    const job: Job = { ...JOB_A, principal: sharedAddr };

    const messages = toMessages({ agents: [agent], jobs: [job] });
    const json = extractJson((messages[1] as { content: string }).content);

    const controller = json.agents[0]?.controller as string;
    const principal = json.jobs[0]?.principal as string;

    expect(controller).toBeDefined();
    expect(controller).toBe(principal);

    // Checksum-cased input collides to the same token (case-insensitive).
    const checksumAgent: Agent = {
      ...AGENT_A,
      controller: `0x${sharedAddr.slice(2).toUpperCase()}` as `0x${string}`,
    };
    const messages2 = toMessages({ agents: [checksumAgent] });
    const json2 = extractJson((messages2[1] as { content: string }).content);
    expect(json2.agents[0]?.controller).toBe(controller);
  });

  it("leaves addresses intact when redactAddresses is false", () => {
    const messages = toMessages({ agents: [AGENT_A], jobs: [JOB_A] }, { redactAddresses: false });
    const json = extractJson((messages[1] as { content: string }).content);
    expect(json.agents[0]?.controller).toBe(AGENT_A.controller);
    expect(json.jobs[0]?.principal).toBe(JOB_A.principal);
  });
});

interface ExtractedPayload {
  agents: Array<Record<string, unknown>>;
  jobs: Array<Record<string, unknown>>;
}

function extractJson(content: string): ExtractedPayload {
  const match = content.match(/```json\n([\s\S]*?)\n```/);
  if (match === null) throw new Error(`no fenced json block in: ${content}`);
  const parsed = JSON.parse(match[1] ?? "") as Partial<ExtractedPayload>;
  return {
    agents: parsed.agents ?? [],
    jobs: parsed.jobs ?? [],
  };
}
