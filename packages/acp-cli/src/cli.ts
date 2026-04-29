#!/usr/bin/env node
// Entry point for the `acp` binary. Wires the per-command runners onto a
// commander program. The runners themselves are unit-tested directly; this
// file is exercised indirectly via a small smoke test that asserts the help
// surface and version flag.

import { Command } from "commander";

import { runAgentCreate, runAgentList, runAgentStatus } from "./commands/agent.js";
import { runBrowse } from "./commands/browse.js";
import { runConfigure } from "./commands/configure.js";
import { runEventsListen } from "./commands/events.js";
import { runJobCreate, runJobList } from "./commands/job.js";

/** Build the commander program. Exported so a test can introspect the help
 * surface without forking a process. */
export function buildProgram(): Command {
  const program = new Command();
  program
    .name("acp")
    .description("Titular ACP command-line interface")
    .version("0.1.0-alpha.0")
    .option("--json", "emit machine-readable JSON instead of pretty output")
    .option("--config <path>", "override the config file path");

  // ---- configure ----
  program
    .command("configure")
    .description("set up gateway URL, chain, and api key")
    .option("--gateway-url <url>", "gateway base URL")
    .option("--chain <slug>", "chain slug (base-sepolia, base-mainnet, ...)")
    .option("--api-key <key>", "bearer / alchemy api key")
    .option("--agent-registry <addr>", "AgentRegistry contract address")
    .option("--job-factory <addr>", "JobFactory contract address")
    .option("--rpc-url <url>", "override RPC URL for write commands")
    .option("--non-interactive", "fail rather than prompting; for CI flows")
    .option("--print", "print the existing config and exit")
    .action(async (opts) => {
      const globals = program.opts();
      await runConfigure({
        json: globals.json === true,
        ...(globals.config !== undefined ? { config: globals.config } : {}),
        ...(opts.gatewayUrl !== undefined ? { gatewayUrl: opts.gatewayUrl } : {}),
        ...(opts.chain !== undefined ? { chain: opts.chain } : {}),
        ...(opts.apiKey !== undefined ? { apiKey: opts.apiKey } : {}),
        ...(opts.agentRegistry !== undefined ? { agentRegistry: opts.agentRegistry } : {}),
        ...(opts.jobFactory !== undefined ? { jobFactory: opts.jobFactory } : {}),
        ...(opts.rpcUrl !== undefined ? { rpcUrl: opts.rpcUrl } : {}),
        ...(opts.nonInteractive === true ? { nonInteractive: true } : {}),
        ...(opts.print === true ? { print: true } : {}),
      });
    });

  // ---- agent ----
  const agent = program.command("agent").description("agent registry operations");
  agent
    .command("create")
    .description("register a new ACP agent on chain")
    .option("--name <name>", "human-readable name (printed only)")
    .option("--capabilities <bitmap>", "capability bitmap (decimal or 0x-hex)")
    .option("--metadata-uri <uri>", "ipfs:// or https:// metadata URI")
    .option("--controller <addr>", "controller address; defaults to signer")
    .action(async (opts) => {
      const globals = program.opts();
      await runAgentCreate({
        json: globals.json === true,
        ...(globals.config !== undefined ? { config: globals.config } : {}),
        ...(opts.name !== undefined ? { name: opts.name } : {}),
        ...(opts.capabilities !== undefined ? { capabilities: opts.capabilities } : {}),
        ...(opts.metadataUri !== undefined ? { metadataUri: opts.metadataUri } : {}),
        ...(opts.controller !== undefined ? { controller: opts.controller } : {}),
      });
    });
  agent
    .command("list")
    .description("list ACP agents in the registry")
    .option("--limit <n>", "page size", parseIntFlag)
    .option("--cursor <c>", "pagination cursor")
    .action(async (opts) => {
      const globals = program.opts();
      await runAgentList({
        json: globals.json === true,
        ...(globals.config !== undefined ? { config: globals.config } : {}),
        ...(opts.limit !== undefined ? { limit: opts.limit } : {}),
        ...(opts.cursor !== undefined ? { cursor: opts.cursor } : {}),
      });
    });
  agent
    .command("status <agentId>")
    .description("show details for a single agent")
    .action(async (agentId, _opts) => {
      const globals = program.opts();
      await runAgentStatus({
        agentId,
        json: globals.json === true,
        ...(globals.config !== undefined ? { config: globals.config } : {}),
      });
    });

  // ---- browse ----
  program
    .command("browse")
    .description("browse all agents (acp + launchpad)")
    .option("--kind <kind>", "filter by kind (acp|launchpad)")
    .option("--limit <n>", "page size", parseIntFlag)
    .option("--cursor <c>", "pagination cursor")
    .action(async (opts) => {
      const globals = program.opts();
      await runBrowse({
        json: globals.json === true,
        ...(globals.config !== undefined ? { config: globals.config } : {}),
        ...(opts.kind !== undefined ? { kind: opts.kind } : {}),
        ...(opts.limit !== undefined ? { limit: opts.limit } : {}),
        ...(opts.cursor !== undefined ? { cursor: opts.cursor } : {}),
      });
    });

  // ---- job ----
  const job = program.command("job").description("job factory operations");
  job
    .command("create")
    .description("create a new ACP job on chain")
    .requiredOption("--token <addr>", "ERC-20 payment token address")
    .requiredOption("--budget <amount>", "budget in token base units (decimal)")
    .requiredOption("--deadline <unix>", "deadline as a UNIX timestamp")
    .option("--target <agentId>", "target agent registry id (0 = open)")
    .option("--terms <text>", "free-form job terms (printed only)")
    .option("--job-type <type>", "direct or evaluated", "direct")
    .option("--evaluator <addr>", "evaluator address (required for evaluated)")
    .option("--arbiter <addr>", "arbiter address override")
    .action(async (opts) => {
      const globals = program.opts();
      await runJobCreate({
        json: globals.json === true,
        ...(globals.config !== undefined ? { config: globals.config } : {}),
        token: opts.token,
        budget: opts.budget,
        deadline: opts.deadline,
        ...(opts.target !== undefined ? { target: opts.target } : {}),
        ...(opts.terms !== undefined ? { terms: opts.terms } : {}),
        ...(opts.jobType !== undefined ? { jobType: opts.jobType } : {}),
        ...(opts.evaluator !== undefined ? { evaluator: opts.evaluator } : {}),
        ...(opts.arbiter !== undefined ? { arbiter: opts.arbiter } : {}),
      });
    });
  job
    .command("list")
    .description("list jobs from the gateway")
    .option("--status <phase>", "filter by phase (created|funded|active|...)")
    .option("--limit <n>", "page size", parseIntFlag)
    .option("--cursor <c>", "pagination cursor")
    .action(async (opts) => {
      const globals = program.opts();
      await runJobList({
        json: globals.json === true,
        ...(globals.config !== undefined ? { config: globals.config } : {}),
        ...(opts.status !== undefined ? { status: opts.status } : {}),
        ...(opts.limit !== undefined ? { limit: opts.limit } : {}),
        ...(opts.cursor !== undefined ? { cursor: opts.cursor } : {}),
      });
    });

  // ---- events ----
  const events = program.command("events").description("gateway event stream operations");
  events
    .command("listen")
    .description("subscribe to the gateway SSE event stream")
    .option("--subjects <csv>", "comma-separated subject filter")
    .option("--last-event-id <id>", "resume from a previous event id")
    .option("--max-events <n>", "stop after N events (debug)", parseIntFlag)
    .action(async (opts) => {
      const globals = program.opts();
      await runEventsListen({
        json: globals.json === true,
        ...(globals.config !== undefined ? { config: globals.config } : {}),
        ...(opts.subjects !== undefined ? { subjects: opts.subjects } : {}),
        ...(opts.lastEventId !== undefined ? { lastEventId: opts.lastEventId } : {}),
        ...(opts.maxEvents !== undefined ? { maxEvents: opts.maxEvents } : {}),
      });
    });

  return program;
}

function parseIntFlag(raw: string): number {
  const n = Number.parseInt(raw, 10);
  if (!Number.isInteger(n) || n < 0) {
    throw new Error(`expected non-negative integer, got "${raw}"`);
  }
  return n;
}

/** Bin entry. Exported so the smoke test can call `main(["--help"])`
 * directly. */
export async function main(argv: readonly string[] = process.argv.slice(2)): Promise<number> {
  const program = buildProgram();
  program.exitOverride();
  try {
    await program.parseAsync(argv, { from: "user" });
    return 0;
  } catch (err) {
    const e = err as { code?: string; exitCode?: number; message?: string };
    if (e.code === "commander.helpDisplayed" || e.code === "commander.help") return 0;
    if (e.code === "commander.version") return 0;
    if (typeof e.exitCode === "number") return e.exitCode;
    process.stderr.write(`${e.message ?? String(err)}\n`);
    return 1;
  }
}

// Top-level await is fine under "module": "ESNext" + Node 22.
const isDirectInvocation =
  typeof process !== "undefined" &&
  Array.isArray(process.argv) &&
  process.argv[1] !== undefined &&
  /acp-cli\/dist\/cli\.js$/.test(process.argv[1]);

if (isDirectInvocation) {
  // eslint-disable-next-line unicorn/no-process-exit -- bin entry
  process.exit(await main());
}
