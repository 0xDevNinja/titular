// `acp browse` — generic listing across kinds (acp + launchpad).
//
// Mostly thin sugar over `agent list`, but with `kind` defaulting to *no*
// filter (let the gateway return both) and the prettier table including the
// `kind` column. Kept as its own command so users targeting only launchpad
// agents have a one-liner.

import type { AgentKind, AgentsPage, GatewayClient } from "@titular/acp-sdk";

import { defaultConfigPath, loadConfig } from "../config.js";
import { type OutputContext, buildOutputContext } from "../output.js";
import { buildGatewayClient } from "../sdk.js";
import type { GlobalFlags } from "../types.js";
import { AgentError, emitAgentsPage } from "./agent.js";

export interface BrowseFlags extends GlobalFlags {
  kind?: string;
  limit?: number;
  cursor?: string;
}

export interface BrowseDeps {
  gatewayFactory?: (cfg: import("../types.js").CliConfig) => GatewayClient;
  stdout?: NodeJS.WriteStream;
  stderr?: NodeJS.WriteStream;
  env?: NodeJS.ProcessEnv;
  isTTY?: boolean;
}

const KIND_VALUES: readonly AgentKind[] = ["acp", "launchpad"] as const;

function asAgentKind(s: string): AgentKind {
  if ((KIND_VALUES as readonly string[]).includes(s)) return s as AgentKind;
  throw new AgentError(
    "invalid_param",
    `browse: --kind must be one of ${KIND_VALUES.join(", ")} (got "${s}")`
  );
}

export async function runBrowse(flags: BrowseFlags, deps: BrowseDeps = {}): Promise<AgentsPage> {
  const env = deps.env ?? process.env;
  const path = flags.config ?? defaultConfigPath(env);
  const cfg = await loadConfig(path);
  if (cfg === undefined) {
    throw new AgentError("no_config", `no config at ${path}; run \`acp configure\` first`);
  }
  const ctx: OutputContext = buildOutputContext({
    json: flags.json === true,
    ...(deps.stdout !== undefined ? { stdout: deps.stdout } : {}),
    ...(deps.stderr !== undefined ? { stderr: deps.stderr } : {}),
    env,
    ...(deps.isTTY !== undefined ? { isTTY: deps.isTTY } : {}),
  });

  const gw = (deps.gatewayFactory ?? buildGatewayClient)(cfg);
  const page = await gw.listAgents({
    ...(flags.kind !== undefined ? { kind: asAgentKind(flags.kind) } : {}),
    ...(flags.limit !== undefined ? { limit: flags.limit } : {}),
    ...(flags.cursor !== undefined ? { cursor: flags.cursor } : {}),
  });
  emitAgentsPage(ctx, page);
  return page;
}
