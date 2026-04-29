// `acp job` — create and list jobs.

import type {
  AcpAgent,
  CreateJobResult,
  GatewayClient,
  JobPhase,
  JobsPage,
} from "@titular/acp-sdk";

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
import type { CliConfig, GlobalFlags } from "../types.js";

export interface JobCreateFlags extends GlobalFlags {
  target?: string;
  budget: string;
  deadline: string;
  terms?: string;
  token?: string;
  jobType?: string;
  evaluator?: string;
  arbiter?: string;
}

export interface JobListFlags extends GlobalFlags {
  status?: string;
  limit?: number;
  cursor?: string;
}

export interface JobDeps {
  acpAgentFactory?: (cfg: CliConfig) => AcpAgent;
  gatewayFactory?: (cfg: CliConfig) => GatewayClient;
  providerFactory?: ProviderFactory;
  stdout?: NodeJS.WriteStream;
  stderr?: NodeJS.WriteStream;
  env?: NodeJS.ProcessEnv;
  isTTY?: boolean;
}

const JOB_PHASES: readonly JobPhase[] = [
  "created",
  "funded",
  "active",
  "completed",
  "cancelled",
  "disputed",
  "released",
  "resolved",
] as const;

export async function runJobCreate(
  flags: JobCreateFlags,
  deps: JobDeps = {}
): Promise<CreateJobResult> {
  const { ctx, cfg } = await loadCtx(flags, deps);

  if (flags.token === undefined || !/^0x[0-9a-fA-F]{40}$/.test(flags.token)) {
    throw new JobError(
      "invalid_param",
      "job create: --token must be a 0x-prefixed 20-byte address"
    );
  }
  const budget = parseBigInt(flags.budget, "--budget");
  const deadline = parseBigInt(flags.deadline, "--deadline");
  if (budget <= 0n) {
    throw new JobError("invalid_param", "job create: --budget must be > 0");
  }
  if (deadline <= 0n) {
    throw new JobError("invalid_param", "job create: --deadline must be a positive UNIX timestamp");
  }

  const targetAgentId = flags.target !== undefined ? parseBigInt(flags.target, "--target") : 0n;

  const jobType = parseJobType(flags.jobType);

  const agent = (deps.acpAgentFactory ?? ((c) => buildAcpAgent(c, deps.providerFactory)))(cfg);
  const result = await agent.createJob({
    targetAgentId,
    token: flags.token as `0x${string}`,
    budget,
    deadline,
    jobType,
    ...(flags.evaluator !== undefined ? { evaluator: flags.evaluator as `0x${string}` } : {}),
    ...(flags.arbiter !== undefined ? { arbiter: flags.arbiter as `0x${string}` } : {}),
  });

  if (ctx.json) {
    writeJson(ctx, {
      jobId: result.jobId.toString(10),
      cloneAddress: result.cloneAddress,
      txHash: result.txHash,
    });
  } else {
    writeKv(ctx, "jobId", result.jobId.toString(10));
    writeKv(ctx, "cloneAddress", result.cloneAddress);
    writeKv(ctx, "txHash", result.txHash);
    if (flags.terms !== undefined && flags.terms.length > 0) {
      writeKv(ctx, "terms", flags.terms);
    }
  }
  return result;
}

export async function runJobList(flags: JobListFlags, deps: JobDeps = {}): Promise<JobsPage> {
  const { ctx, cfg } = await loadCtx(flags, deps);
  const gw = (deps.gatewayFactory ?? buildGatewayClient)(cfg);
  const page = await gw.listJobs({
    ...(flags.status !== undefined ? { status: parseJobPhase(flags.status) } : {}),
    ...(flags.limit !== undefined ? { limit: flags.limit } : {}),
    ...(flags.cursor !== undefined ? { cursor: flags.cursor } : {}),
  });
  emitJobsPage(ctx, page);
  return page;
}

function emitJobsPage(ctx: OutputContext, page: JobsPage): void {
  if (ctx.json) {
    writeJson(ctx, page);
    return;
  }
  if (page.data.length === 0) {
    writeLine(ctx, errorString(ctx, "(no jobs)"));
  } else {
    const rows = page.data.map((j) => [
      j.on_chain_id ?? String(j.id),
      j.agent_id,
      j.phase,
      j.principal,
      j.token,
      j.budget,
      j.deadline,
    ]);
    writeTable(
      ctx,
      ["job_id", "agent_id", "phase", "principal", "token", "budget", "deadline"],
      rows
    );
  }
  if (page.next_cursor.length > 0) {
    writeLine(ctx, `next_cursor: ${page.next_cursor}`);
  }
}

async function loadCtx(
  flags: GlobalFlags,
  deps: JobDeps
): Promise<{ ctx: OutputContext; cfg: CliConfig }> {
  const env = deps.env ?? process.env;
  const path = flags.config ?? defaultConfigPath(env);
  const cfg = await loadConfig(path);
  if (cfg === undefined) {
    throw new JobError("no_config", `no config at ${path}; run \`acp configure\` first`);
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

function parseBigInt(raw: string, label: string): bigint {
  const t = raw.trim();
  if (!/^\d+$/.test(t)) {
    throw new JobError("invalid_param", `${label}: must be a non-negative decimal integer`);
  }
  return BigInt(t);
}

function parseJobType(raw: string | undefined): 0 | 1 {
  if (raw === undefined) return 0;
  const t = raw.trim().toLowerCase();
  if (t === "direct" || t === "0") return 0;
  if (t === "evaluated" || t === "1") return 1;
  throw new JobError("invalid_param", `--job-type must be "direct" or "evaluated" (got "${raw}")`);
}

function parseJobPhase(raw: string): JobPhase {
  const t = raw.trim().toLowerCase();
  if ((JOB_PHASES as readonly string[]).includes(t)) return t as JobPhase;
  throw new JobError(
    "invalid_param",
    `--status must be one of ${JOB_PHASES.join(", ")} (got "${raw}")`
  );
}

export class JobError extends Error {
  override readonly name = "JobError";
  readonly code: string;
  constructor(code: string, message: string) {
    super(message);
    this.code = code;
  }
}

export const __testing__ = { parseBigInt, parseJobPhase, parseJobType };
