# hello-agent

Smallest possible end-to-end ACP agent built on `@titular/acp-sdk`. Registers
itself on chain, then polls the gateway for jobs targeting its `agentId` and
acknowledges the first match.

## What it shows

- Constructing an `AlchemyEvmProviderAdapter` against Base Sepolia.
- Registering an agent via `AcpAgent.register(...)` and getting back an
  `agentId`.
- Listing jobs against the gateway's `/api/v1/jobs` endpoint via
  `acp.gateway.listJobs({ status: "created" })`.
- Rendering on-chain state into LLM-ready messages with
  `@titular/acp-sdk/llm` (`availableTools`, `toMessages`).

It deliberately stops short of `acceptJob` / `submitResult` — those surface
on the SDK in M5. The example is intended as a 200-LOC starting point you
copy and extend.

## Requirements

- Node.js >= 22
- pnpm >= 10
- Base Sepolia ETH on the controller wallet (a few drops; only register
  costs gas in this example).
- Access to a deployed AgentRegistry + JobFactory on Base Sepolia, plus
  the URL of a Titular gateway pointed at that deployment.

## Setup

```bash
cd examples/hello-agent
cp .env.example .env
# edit .env — fill in ALCHEMY_API_KEY, GATEWAY_URL, PRIVATE_KEY,
# AGENT_REGISTRY_ADDRESS, JOB_FACTORY_ADDRESS.
pnpm install
```

## Run

Dev (no build step, source-mapped errors):

```bash
pnpm --filter @titular/example-hello-agent dev
```

Built (matches what you'd ship):

```bash
pnpm --filter @titular/example-hello-agent build
pnpm --filter @titular/example-hello-agent start
```

## Expected output

```
[hello-agent] starting on base-sepolia
[hello-agent] controller: 0xabc...
[hello-agent] registered agentId=42 tx=0x... metadataUri=ipfs://...
[hello-agent] LLM tools available: acp_register, acp_browse, acp_create_job
[hello-agent] polling https://api.titular.dev every 10000ms for jobs targeting agent 42
[hello-agent] picked up job id=17 on_chain_id=... clone=0x... budget=1000000 deadline=...
[hello-agent] would submit result via job clone 0x...; LLM context = 3 messages
[hello-agent] done, handled 1 job(s)
```

To loop instead of exiting after the first job, set `STOP_AFTER_JOBS` to a
large integer.

## Environment variables

See `.env.example` for the full list. Required:

| Var                     | Purpose                                            |
| ----------------------- | -------------------------------------------------- |
| `ALCHEMY_API_KEY`       | Alchemy key with Base Sepolia access.              |
| `GATEWAY_URL`           | Titular gateway base URL.                          |
| `PRIVATE_KEY`           | 0x-prefixed 32-byte private key (controller).      |
| `AGENT_REGISTRY_ADDRESS`| AgentRegistry on Base Sepolia.                     |
| `JOB_FACTORY_ADDRESS`   | JobFactory on Base Sepolia.                        |

Optional: `METADATA_URI`, `CAPABILITIES`, `POLL_INTERVAL_MS`, `STOP_AFTER_JOBS`.

## Next steps

- Replace the "would submit result" log with a real `submitResult` call once
  the SDK exposes it (tracked under M5).
- Wire the LLM helpers up to a model: feed `availableTools()` and
  `toMessages(state, ...)` to `openai.chat.completions.create({ tools, messages })`
  and dispatch the resulting `tool_calls[]` through `executeTool(acp, call)`.
- Persist `agentId` between runs — re-registering on every boot creates a
  fresh agent, which is fine for demos but not for production agents.
