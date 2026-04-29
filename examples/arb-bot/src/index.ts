// arb-bot — connects to the Titular gateway SSE stream, watches `trades.*`
// events, and opens an ACP job whenever the spread between two configured
// agent tokens exceeds the configured threshold.
//
// State machine:
//
//                ┌──────────┐
//                │   idle   │  ◄── boot
//                └──────────┘
//                     │ first SSE event
//                     ▼
//                ┌──────────┐
//                │ scanning │  ◄────────────────┐
//                └──────────┘                   │
//                     │ detectArb >= threshold  │
//                     ▼                         │
//                ┌──────────┐                   │
//                │executing │  createJob (1 tx) │
//                └──────────┘                   │
//                     │ tx confirmed            │
//                     ▼                         │
//                ┌──────────┐ cooldown elapsed  │
//                │ cooldown │ ──────────────────┘
//                └──────────┘
//
// We deliberately consume the gateway SSE stream rather than connecting
// directly to NATS so the example mirrors what an external developer can
// reach without infra access. The stream format is documented in
// services/gateway-go/internal/sse/handler.go (subjects=trades.*).
//
// What this example demonstrates:
//   - SSE consumption against the gateway, with reconnect-on-error.
//   - Pure-function strategy module that's testable in isolation.
//   - SDK usage: `AcpAgent.createJob(...)` against Base Sepolia, with
//     evaluator/arbiter left to factory defaults.
//   - State-machine discipline (no concurrent createJob calls).
//
// What it does NOT do:
//   - Move ERC-20 allowance — the bot's wallet must have already approved
//     the JobFactory for QUOTE_TOKEN.
//   - Actually execute the arb on chain. The created job carries the
//     intent; a downstream agent (or the bot's own future M5 surface)
//     would `acceptJob` + `submitResult` to settle.

import {
  AcpAgent,
  AlchemyEvmProviderAdapter,
  type CreateJobResult,
  type EvmAddress,
} from "@titular/acp-sdk";
import type { Hex } from "viem";
import {
  type ArbOpportunity,
  PriceBook,
  type StrategyConfig,
  type TradeEvent,
  detectArb,
  quoteFromTrade,
} from "./strategy.js";

interface Config {
  alchemyApiKey: string;
  gatewayUrl: string;
  privateKey: Hex;
  agentRegistry: EvmAddress;
  jobFactory: EvmAddress;
  quoteToken: EvmAddress;
  agentTokens: readonly string[];
  minSpreadBps: number;
  arbBudget: bigint;
  cooldownMs: number;
  jobDeadlineSecs: number;
}

type State = "idle" | "scanning" | "executing" | "cooldown";

function loadConfig(): Config {
  const requireEnv = (key: string): string => {
    const v = process.env[key];
    if (v === undefined || v.length === 0) {
      throw new Error(`missing required env var: ${key}`);
    }
    return v;
  };
  const optional = (key: string, fallback: string): string => process.env[key] ?? fallback;
  const requireAddress = (key: string): EvmAddress => {
    const v = requireEnv(key);
    if (!/^0x[0-9a-fA-F]{40}$/.test(v)) {
      throw new Error(`${key} must be a 0x-prefixed 20-byte hex address`);
    }
    return v.toLowerCase() as EvmAddress;
  };
  const privateKey = requireEnv("PRIVATE_KEY");
  if (!/^0x[0-9a-fA-F]{64}$/.test(privateKey)) {
    throw new Error("PRIVATE_KEY must be a 0x-prefixed 32-byte hex string");
  }
  const tokensRaw = requireEnv("AGENT_TOKENS");
  const agentTokens = tokensRaw
    .split(",")
    .map((s) => s.trim().toLowerCase())
    .filter((s) => s.length > 0);
  if (agentTokens.length < 2) {
    throw new Error("AGENT_TOKENS must list at least 2 comma-separated 0x-hex addresses");
  }
  for (const t of agentTokens) {
    if (!/^0x[0-9a-fA-F]{40}$/.test(t)) {
      throw new Error(`AGENT_TOKENS contains an invalid address: ${t}`);
    }
  }
  return {
    alchemyApiKey: requireEnv("ALCHEMY_API_KEY"),
    gatewayUrl: requireEnv("GATEWAY_URL").replace(/\/+$/, ""),
    privateKey: privateKey as Hex,
    agentRegistry: requireAddress("AGENT_REGISTRY_ADDRESS"),
    jobFactory: requireAddress("JOB_FACTORY_ADDRESS"),
    quoteToken: requireAddress("QUOTE_TOKEN_ADDRESS"),
    agentTokens,
    minSpreadBps: Number.parseInt(optional("MIN_SPREAD_BPS", "50"), 10),
    arbBudget: BigInt(optional("ARB_BUDGET", "1000000")),
    cooldownMs: Number.parseInt(optional("COOLDOWN_MS", "30000"), 10),
    jobDeadlineSecs: Number.parseInt(optional("JOB_DEADLINE_SECS", "3600"), 10),
  };
}

interface SseTradeEnvelope {
  subject?: string;
  payload?: {
    curve?: string;
    side?: "buy" | "sell";
    quote_in?: string;
    agent_out?: string;
    agent_in?: string;
    quote_out?: string;
    block_number?: number;
    agent_token?: string;
  };
}

/** Minimal SSE line parser. Splits on double-newline events; ignores comments. */
async function* readSseEvents(response: Response): AsyncGenerator<SseTradeEnvelope> {
  if (response.body === null) {
    throw new Error("SSE response has no body");
  }
  const reader = response.body.getReader();
  const decoder = new TextDecoder("utf-8");
  let buffer = "";
  while (true) {
    const { value, done } = await reader.read();
    if (done) return;
    buffer += decoder.decode(value, { stream: true });
    let sep = buffer.indexOf("\n\n");
    while (sep !== -1) {
      const block = buffer.slice(0, sep);
      buffer = buffer.slice(sep + 2);
      sep = buffer.indexOf("\n\n");
      const dataLine = block
        .split("\n")
        .map((l) => l.trim())
        .filter((l) => l.startsWith("data:"))
        .map((l) => l.slice("data:".length).trim())
        .join("");
      if (dataLine.length === 0) continue;
      try {
        yield JSON.parse(dataLine) as SseTradeEnvelope;
      } catch {
        // ignore malformed events; the gateway only emits JSON but a
        // misbehaving proxy might inject keepalive comments etc.
      }
    }
  }
}

/** Translate a gateway SSE envelope into a {@link TradeEvent}, or `null` if it doesn't apply. */
function envelopeToTrade(env: SseTradeEnvelope): TradeEvent | null {
  if (env.subject === undefined || !env.subject.startsWith("trades.")) return null;
  const p = env.payload;
  if (p === undefined) return null;
  const agentToken = (p.agent_token ?? p.curve ?? "").toLowerCase();
  if (!/^0x[0-9a-fA-F]{40}$/.test(agentToken)) return null;
  const side = p.side;
  if (side !== "buy" && side !== "sell") return null;
  const trade: TradeEvent = {
    agentToken,
    side,
    blockNumber: p.block_number ?? 0,
  };
  if (p.quote_in !== undefined) trade.quoteIn = p.quote_in;
  if (p.quote_out !== undefined) trade.quoteOut = p.quote_out;
  if (p.agent_in !== undefined) trade.agentIn = p.agent_in;
  if (p.agent_out !== undefined) trade.agentOut = p.agent_out;
  return trade;
}

async function executeArb(
  acp: AcpAgent,
  config: Config,
  opp: ArbOpportunity
): Promise<CreateJobResult> {
  const deadline = BigInt(Math.floor(Date.now() / 1000) + config.jobDeadlineSecs);
  console.log(
    `[arb-bot] executing arb cheap=${opp.cheapToken} rich=${opp.richToken} ` +
      `spread=${opp.spreadBps}bps budget=${config.arbBudget.toString()} deadline=${deadline.toString()}`
  );
  const result = await acp.createJob({
    token: config.quoteToken,
    budget: config.arbBudget,
    deadline,
  });
  console.log(
    `[arb-bot] job created jobId=${result.jobId.toString()} clone=${result.cloneAddress} ` +
      `tx=${result.txHash}`
  );
  return result;
}

async function streamTrades(
  config: Config,
  onEvent: (env: SseTradeEnvelope) => void,
  abort: AbortSignal
): Promise<void> {
  const url = `${config.gatewayUrl}/events?subjects=trades.*`;
  const res = await fetch(url, {
    headers: { Accept: "text/event-stream" },
    signal: abort,
  });
  if (!res.ok) {
    throw new Error(`SSE connect failed: ${res.status} ${res.statusText}`);
  }
  for await (const env of readSseEvents(res)) {
    if (abort.aborted) return;
    onEvent(env);
  }
}

async function main(): Promise<void> {
  const config = loadConfig();
  console.log(
    `[arb-bot] starting on base-sepolia; watching ${config.agentTokens.length} tokens; ` +
      `min spread=${config.minSpreadBps}bps`
  );

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

  const book = new PriceBook();
  const strategyConfig: StrategyConfig = {
    minSpreadBps: config.minSpreadBps,
    watchedTokens: config.agentTokens,
  };
  let state: State = "idle";
  let cooldownUntil = 0;

  const ac = new AbortController();
  const onShutdown = (): void => {
    console.log("[arb-bot] shutdown signal received");
    ac.abort();
  };
  process.once("SIGINT", onShutdown);
  process.once("SIGTERM", onShutdown);

  // Reconnect loop. SSE drops are normal (proxy rotations, gateway
  // redeploys); we backoff exponentially up to 30 s.
  let backoffMs = 1000;
  while (!ac.signal.aborted) {
    try {
      state = "scanning";
      console.log(`[arb-bot] connecting to ${config.gatewayUrl}/events`);
      await streamTrades(
        config,
        (env) => {
          const trade = envelopeToTrade(env);
          if (trade === null) return;
          if (!config.agentTokens.includes(trade.agentToken)) return;
          const quote = quoteFromTrade(trade);
          if (quote === null) return;
          book.observe(quote);
          if (state !== "scanning") return;
          if (Date.now() < cooldownUntil) return;
          const opp = detectArb(book, strategyConfig);
          if (opp === null) return;
          state = "executing";
          executeArb(acp, config, opp)
            .catch((err: unknown) => {
              console.error(`[arb-bot] arb execution failed: ${(err as Error).message}`);
            })
            .finally(() => {
              cooldownUntil = Date.now() + config.cooldownMs;
              state = "cooldown";
              console.log(`[arb-bot] cooling down for ${config.cooldownMs}ms`);
              setTimeout(() => {
                if (!ac.signal.aborted) state = "scanning";
              }, config.cooldownMs);
            });
        },
        ac.signal
      );
      // streamTrades returned cleanly (gateway closed the connection); reset
      // backoff and reconnect.
      backoffMs = 1000;
    } catch (err) {
      if (ac.signal.aborted) break;
      console.error(
        `[arb-bot] sse error, reconnecting in ${backoffMs}ms: ${(err as Error).message}`
      );
      await sleep(backoffMs);
      backoffMs = Math.min(backoffMs * 2, 30_000);
    }
  }
  console.log("[arb-bot] exited cleanly");
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

main().catch((err) => {
  console.error("[arb-bot] fatal:", err);
  process.exit(1);
});
