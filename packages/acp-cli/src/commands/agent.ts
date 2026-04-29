// `acp agent` — register, list, and inspect agents on the gateway.

import type { AcpAgent, Agent, AgentsPage, GatewayClient, RegisterResult } from "@titular/acp-sdk";

import { defaultConfigPath, loadConfig } from "../config.js";
import {
  type OutputContext,
  buildOutputContext,
  errorString,
  writeJson,
  writeKv,
  writeLine,
  writeTable,
} from "../output.js";
import { type ProviderFactory, buildAcpAgent, buildGatewayClient } from "../sdk.js";
import type { GlobalFlags } from "../types.js";

export interface AgentCreateFlags extends GlobalFlags {
  name?: string;
  capabilities?: string;
  metadataUri?: string;
  controller?: string;
}

export interface AgentListFlags extends GlobalFlags {
  kind?: string;
  limit?: number;
  cursor?: string;
}

export interface AgentStatusFlags extends GlobalFlags {
  agentId: string;
  recentJobs?: number;
}

/** Optional dependencies. Tests inject `acpAgentFactory` to skip Alchemy and
 * `gatewayFactory` to mock REST calls. */
export interface AgentDeps {
  acpAgentFactory?: (cfg: import("../types.js").CliConfig) => AcpAgent;
  gatewayFactory?: (cfg: import("../types.js").CliConfig) => GatewayClient;
  providerFactory?: ProviderFactory;
  stdout?: NodeJS.WriteStream;
  stderr?: NodeJS.WriteStream;
  env?: NodeJS.ProcessEnv;
  isTTY?: boolean;
}

/** `acp agent create` — submits AgentRegistry.register on chain. */
export async function runAgentCreate(
  flags: AgentCreateFlags,
  deps: AgentDeps = {}
): Promise<RegisterResult> {
  const { ctx, cfg } = await loadCtx(flags, deps);

  if (flags.metadataUri === undefined || flags.metadataUri.length === 0) {
    throw new AgentError("invalid_param", "agent create: --metadata-uri is required");
  }
  // Capabilities default to 0 (no capabilities advertised) so a quick smoke
  // create works without the operator working out the bitmap. We accept
  // both decimal strings and `0x...` hex.
  const capsRaw = flags.capabilities ?? "0";
  const capabilities = parseCapabilities(capsRaw);

  const agent = (deps.acpAgentFactory ?? ((c) => buildAcpAgent(c, deps.providerFactory)))(cfg);
  const result = await agent.register({
    metadataUri: flags.metadataUri,
    capabilities,
    ...(flags.controller !== undefined ? { controller: flags.controller as `0x${string}` } : {}),
  });

  if (ctx.json) {
    writeJson(ctx, {
      agentId: result.agentId.toString(10),
      txHash: result.txHash,
      controller: result.controller,
      metadataUri: result.metadataUri,
      capabilities: result.capabilities.toString(10),
    });
  } else {
    writeKv(ctx, "agentId", result.agentId.toString(10));
    writeKv(ctx, "txHash", result.txHash);
    writeKv(ctx, "controller", result.controller);
    writeKv(ctx, "metadataUri", result.metadataUri);
    writeKv(ctx, "capabilities", result.capabilities.toString(10));
    if (flags.name !== undefined && flags.name.length > 0) {
      writeKv(ctx, "name", flags.name);
    }
  }
  return result;
}

/** `acp agent list` — registry browse, scoped to ACP agents. */
export async function runAgentList(
  flags: AgentListFlags,
  deps: AgentDeps = {}
): Promise<AgentsPage> {
  const { ctx, cfg } = await loadCtx(flags, deps);
  const gw = (deps.gatewayFactory ?? buildGatewayClient)(cfg);
  const page = await gw.listAgents({
    kind: "acp",
    ...(flags.limit !== undefined ? { limit: flags.limit } : {}),
    ...(flags.cursor !== undefined ? { cursor: flags.cursor } : {}),
  });
  emitAgentsPage(ctx, page);
  return page;
}

/** `acp agent status <agent-id>` — single agent details + recent jobs. */
export async function runAgentStatus(
  flags: AgentStatusFlags,
  deps: AgentDeps = {}
): Promise<Agent> {
  const { ctx, cfg } = await loadCtx(flags, deps);
  if (flags.agentId.length === 0) {
    throw new AgentError("invalid_param", "agent status: <agent-id> is required");
  }
  const gw = (deps.gatewayFactory ?? buildGatewayClient)(cfg);
  const agent = await gw.getAgent(flags.agentId);

  if (ctx.json) {
    writeJson(ctx, agent);
  } else {
    writeKv(ctx, "agent_id", agent.agent_id);
    writeKv(ctx, "kind", agent.kind);
    if (agent.controller !== undefined) writeKv(ctx, "controller", agent.controller);
    if (agent.creator !== undefined) writeKv(ctx, "creator", agent.creator);
    writeKv(ctx, "capabilities", agent.capabilities);
    if (agent.metadata_uri !== undefined) writeKv(ctx, "metadata_uri", agent.metadata_uri);
    writeKv(ctx, "block_number", String(agent.block_number));
    writeKv(ctx, "tx_hash", agent.tx_hash);
    writeKv(ctx, "created_at", agent.created_at);
  }
  return agent;
}

/** Render an AgentsPage as either a table (default) or JSON. Shared by
 * `agent list` and the generic `browse` command. */
export function emitAgentsPage(ctx: OutputContext, page: AgentsPage): void {
  if (ctx.json) {
    writeJson(ctx, page);
    return;
  }
  if (page.data.length === 0) {
    writeLine(ctx, errorString(ctx, "(no agents)"));
  } else {
    const rows = page.data.map((a) => [
      a.agent_id,
      a.kind,
      a.controller ?? "-",
      a.capabilities,
      String(a.block_number),
    ]);
    writeTable(ctx, ["agent_id", "kind", "controller", "capabilities", "block"], rows);
  }
  if (page.next_cursor.length > 0) {
    writeLine(ctx, `next_cursor: ${page.next_cursor}`);
  }
}

/** Shared bootstrap that loads config + builds an output context. */
async function loadCtx(
  flags: GlobalFlags,
  deps: AgentDeps
): Promise<{ ctx: OutputContext; cfg: import("../types.js").CliConfig }> {
  const env = deps.env ?? process.env;
  const path = flags.config ?? defaultConfigPath(env);
  const cfg = await loadConfig(path);
  if (cfg === undefined) {
    throw new AgentError("no_config", `no config at ${path}; run \`acp configure\` first`);
  }
  const ctx = buildOutputContext({
    json: flags.json === true,
    ...(deps.stdout !== undefined ? { stdout: deps.stdout } : {}),
    ...(deps.stderr !== undefined ? { stderr: deps.stderr } : {}),
    env,
    ...(deps.isTTY !== undefined ? { isTTY: deps.isTTY } : {}),
  });
  return { ctx, cfg };
}

function parseCapabilities(raw: string): bigint {
  const trimmed = raw.trim();
  if (/^0x[0-9a-fA-F]+$/.test(trimmed)) {
    return BigInt(trimmed);
  }
  if (/^\d+$/.test(trimmed)) {
    return BigInt(trimmed);
  }
  throw new AgentError(
    "invalid_param",
    `agent create: --capabilities must be a decimal or 0x-hex bitmap (got "${raw}")`
  );
}

export class AgentError extends Error {
  override readonly name = "AgentError";
  readonly code: string;
  constructor(code: string, message: string) {
    super(message);
    this.code = code;
  }
}

export const __testing__ = { parseCapabilities };
