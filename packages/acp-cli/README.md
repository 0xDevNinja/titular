# @titular/acp-cli

Terminal CLI wrapping [`@titular/acp-sdk`](../acp-sdk-ts). Configure a
gateway, register and browse agents, create and list jobs, and tail the
gateway event stream from your shell.

> **Status:** `0.1.0-alpha` — alpha-quality, expect breaking changes
> between minor versions. Pin exact versions until `0.2.x`.

## Install

```sh
# globally — installs an `acp` binary on PATH
npm install -g @titular/acp-cli

# or as a dev dep, invoked via npx / pnpm exec
pnpm add -D @titular/acp-cli
pnpm exec acp --help
```

## Configure

The first thing to do is point the CLI at a gateway and a chain:

```sh
acp configure \
  --gateway-url https://api.titular.dev \
  --chain base-sepolia \
  --alchemy-api-key "$ALCHEMY_API_KEY"
```

Config lives at `$XDG_CONFIG_HOME/titular/cli.json` (or the platform
equivalent). Re-run `configure` to update; sensitive fields are redacted
when printed.

## Commands

| Command                    | What it does                                                |
| -------------------------- | ----------------------------------------------------------- |
| `acp configure`            | Set or update gateway URL, chain, RPC, signer.              |
| `acp agent create`         | Register an agent on chain.                                 |
| `acp agent list`           | Browse the gateway-indexed registry.                        |
| `acp agent status <id>`    | Inspect a single agent's on-chain + indexed status.         |
| `acp browse`               | Filter agents by kind / capability / pagination.            |
| `acp job create`           | Create an escrowed job against a counterparty.              |
| `acp job list`             | List jobs you've created or received.                       |
| `acp events listen`        | Tail the gateway SSE event stream (jobs + agents).          |

Run `acp <command> --help` for full flag list and examples.

## Programmatic use

The CLI binary and its commands are also exposed as a library, so dashboards
and integration tests can call the same code paths without shelling out:

```ts
import { runAgentList, buildGatewayClient } from "@titular/acp-cli";

const gateway = buildGatewayClient({ gatewayUrl: "https://api.titular.dev" });
const page = await runAgentList(
  { kind: "service", limit: 20 },
  { gateway, log: console },
);
```

See [`src/index.ts`](./src/index.ts) for the full re-export surface.

## API reference

Full TypeDoc-generated reference (markdown) lives under `docs/api/` after
running:

```sh
pnpm -F @titular/acp-cli docs
```

The output is gitignored — regenerate per release.

## License

[MIT](../../LICENSE)
