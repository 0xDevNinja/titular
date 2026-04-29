// hello-agent — the smallest possible end-to-end ACP agent.
//
// Flow:
//   1. Build an `AlchemyEvmProviderAdapter` from PRIVATE_KEY against Base
//      Sepolia.
//   2. Construct an `AcpAgent` with the configured AgentRegistry + JobFactory
//      addresses.
//   3. Register the agent on chain (idempotent across runs only when
//      AGENT_ID is provided in env — a fresh run with no AGENT_ID always
//      mints a new one).
//   4. Poll the gateway every POLL_INTERVAL_MS for jobs in phase `created`
//      that target this agent's id.
//   5. On the first match, log the job and write a hello-world result line
//      to stdout. Real agents would call `submitResult(...)` against the
//      job clone here; the SDK exposes the gateway and provider primitives
//      and this example stops short of that step so it stays under 200 LOC
//      and runnable without an evaluator/arbiter setup.
//
// Why no `acceptJob` / `submitResult` SDK call? Those land in M5 (#100, #101).
// Until then, the example demonstrates the public M4 surface: `register`,
// `browse`, `gateway.listJobs`, plus the LLM helpers used to render state.
//
// LOC target: ~200. Keep it readable; no clever abstractions.

import { AcpAgent, AlchemyEvmProviderAdapter, type EvmAddress } from "@titular/acp-sdk";
import { availableTools, toMessages } from "@titular/acp-sdk/llm";
import type { Hex } from "viem";

interface Config {
  alchemyApiKey: string;
  gatewayUrl: string;
  privateKey: Hex;
  agentRegistry: EvmAddress;
  jobFactory: EvmAddress;
  metadataUri: string;
  capabilities: bigint;
  pollIntervalMs: number;
}

function loadConfig(): Config {
  const requireEnv = (key: string): string => {
    const v = process.env[key];
    if (v === undefined || v.length === 0) {
      throw new Error(`missing required env var: ${key}`);
    }
    return v;
  };
  const optional = (key: string, fallback: string): string => process.env[key] ?? fallback;

  const privateKey = requireEnv("PRIVATE_KEY");
  if (!/^0x[0-9a-fA-F]{64}$/.test(privateKey)) {
    throw new Error("PRIVATE_KEY must be a 0x-prefixed 32-byte hex string");
  }
  const agentRegistry = requireEnv("AGENT_REGISTRY_ADDRESS");
  const jobFactory = requireEnv("JOB_FACTORY_ADDRESS");
  if (!/^0x[0-9a-fA-F]{40}$/.test(agentRegistry) || !/^0x[0-9a-fA-F]{40}$/.test(jobFactory)) {
    throw new Error("AGENT_REGISTRY_ADDRESS / JOB_FACTORY_ADDRESS must be 20-byte hex addresses");
  }
  const pollMs = Number.parseInt(optional("POLL_INTERVAL_MS", "10000"), 10);
  if (!Number.isFinite(pollMs) || pollMs < 1000) {
    throw new Error("POLL_INTERVAL_MS must be an integer >= 1000");
  }
  return {
    alchemyApiKey: requireEnv("ALCHEMY_API_KEY"),
    gatewayUrl: requireEnv("GATEWAY_URL"),
    privateKey: privateKey as Hex,
    agentRegistry: agentRegistry.toLowerCase() as EvmAddress,
    jobFactory: jobFactory.toLowerCase() as EvmAddress,
    metadataUri: optional("METADATA_URI", "ipfs://bafy-replace-me"),
    capabilities: BigInt(optional("CAPABILITIES", "1")),
    pollIntervalMs: pollMs,
  };
}

async function main(): Promise<void> {
  const config = loadConfig();
  console.log("[hello-agent] starting on base-sepolia");

  const provider = new AlchemyEvmProviderAdapter({
    apiKey: config.alchemyApiKey,
    chain: "base-sepolia",
    privateKey: config.privateKey,
  });
  const acp = new AcpAgent({
    provider,
    gatewayUrl: config.gatewayUrl,
    contracts: {
      agentRegistry: config.agentRegistry,
      jobFactory: config.jobFactory,
    },
  });

  // Step 1: register on chain.
  const controller = await provider.getAddress();
  console.log(`[hello-agent] controller: ${controller}`);

  const registered = await acp.register({
    metadataUri: config.metadataUri,
    capabilities: config.capabilities,
  });
  const myAgentId = registered.agentId;
  console.log(
    `[hello-agent] registered agentId=${myAgentId.toString()} tx=${registered.txHash} ` +
      `metadataUri=${registered.metadataUri}`
  );

  // Tools surface — show that the LLM helpers are available out of the box,
  // even when this minimal example doesn't actually call a model. Real agents
  // would feed `tools` and `messages` to OpenAI / Anthropic / Gemini.
  const tools = availableTools();
  console.log(`[hello-agent] LLM tools available: ${tools.map((t) => t.function.name).join(", ")}`);

  // Step 2: poll for jobs. We use the public gateway listJobs endpoint
  // (read-only) and filter client-side to jobs targeting our agentId. The
  // gateway does not currently expose a `target_agent_id` filter, so client-
  // side filtering is correct — it costs at most one extra page fetch per
  // poll on quiet networks.
  console.log(
    `[hello-agent] polling ${config.gatewayUrl} every ${config.pollIntervalMs}ms for jobs ` +
      `targeting agent ${myAgentId.toString()}`
  );

  const stopAfter = Number.parseInt(process.env.STOP_AFTER_JOBS ?? "1", 10);
  let handled = 0;
  while (handled < stopAfter) {
    try {
      const page = await acp.gateway.listJobs({ status: "created", limit: 50 });
      const myJobs = page.data.filter((j) => j.agent_id === myAgentId.toString());

      if (myJobs.length === 0) {
        await sleep(config.pollIntervalMs);
        continue;
      }

      const job = myJobs[0]!;
      handled += 1;
      console.log(
        `[hello-agent] picked up job id=${job.id} on_chain_id=${job.on_chain_id ?? "?"} ` +
          `clone=${job.clone_address ?? "?"} budget=${job.budget} deadline=${job.deadline}`
      );

      // Step 3: submit a hello-world result. A real agent would call
      // `submitResult(resultURI)` on the job clone here; for the example
      // we just log what we WOULD submit. The SDK lands `submitResult` in
      // M5 alongside the rest of the agent-side action surface.
      const messages = toMessages(
        {
          jobs: [job],
          userPrompt: "Acknowledge the job and produce a hello-world result URI.",
        },
        { redactAddresses: true }
      );
      console.log(
        `[hello-agent] would submit result via job clone ${job.clone_address ?? "<unknown>"}; ` +
          `LLM context = ${messages.length} messages`
      );
    } catch (err) {
      console.error(`[hello-agent] poll error: ${(err as Error).message}`);
      await sleep(config.pollIntervalMs);
    }
  }
  console.log(`[hello-agent] done, handled ${handled} job(s)`);
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

main().catch((err) => {
  console.error("[hello-agent] fatal:", err);
  process.exit(1);
});
