# @titular/acp-sdk-e2e

End-to-end integration tests for [`@titular/acp-sdk`](../acp-sdk-ts) wired
against a real anvil chain and the gateway-go REST surface.

This package is `private: true` and is **never published**. It lives outside
`acp-sdk-ts` so the SDK's published tarball does not pull `testcontainers`,
`viem` (for chain bring-up), or anything else it does not need at runtime.

## Scope

The suite exercises the public SDK surface against real backends:

- `AcpAgent.register()` — submits `AgentRegistry.register(...)` to anvil,
  waits for inclusion, decodes the `AgentRegistered` event.
- `AcpAgent.createJob()` — submits `JobFactory.createJob(...)` to anvil,
  waits for inclusion, decodes `JobCreated`.
- `GatewayClient.listAgents()` / `listJobs()` / `getStats()` — read-only REST
  calls against a live gateway-go process backed by a Postgres container.
- `GatewayClient` error handling — drives 404 / 500 / network-error paths
  against the real gateway and asserts the `AcpError` codes raised.

The on-chain bring-up uses minimal mock contracts whose runtime bytecode is
embedded in `src/lib/mock-bytecode.ts`. They emit events shape-identical to
the production `AgentRegistry` and `JobFactory` (matching the ABI fragments
in `packages/acp-sdk-ts/src/abi.ts`) — enough to drive the SDK's
event-decoding paths without dragging the full ACP solidity toolchain into
this package.

The Postgres side is seeded directly with rows shaped to the indexer's
`0001_init` schema. The full indexer pipeline (subscriber → decoder →
publisher) is covered by `services/indexer-go/internal/integration/...`;
this suite focuses specifically on the **SDK ↔ gateway** wire shape.

## Running

The whole suite is gated on `E2E_INTEGRATION=1`. Without that env var, every
top-level `describe` short-circuits to `describe.skip`, and `vitest run`
exits 0 in seconds. That keeps `pnpm turbo run test` cheap for contributors
who haven't installed foundry / docker / a Go toolchain.

```sh
# Prereqs: foundry (anvil on PATH), docker daemon running, go 1.25+
export E2E_INTEGRATION=1
pnpm -F @titular/acp-sdk-e2e test
```

When prerequisites are missing the per-suite `beforeAll` calls
`ctx.skip()` with a one-line reason, so `vitest` reports the tests as
skipped (not failed). CI runs the suite in a dedicated workflow
(`.github/workflows/sdk-e2e.yml`) on `workflow_dispatch`, the weekly cron,
and on PRs that touch `packages/acp-sdk-ts/**` or `packages/acp-sdk-e2e/**`.

## Why a separate package, not a build-tag in `acp-sdk-ts`

- Keeps the SDK's `dependencies` and `devDependencies` minimal — `vitest`
  is already there but `testcontainers` (and its transitive `dockerode`)
  would otherwise leak into `pnpm i @titular/acp-sdk`-adjacent dev
  installs in downstream consumers.
- Lets the e2e suite import the SDK as a workspace `dependency` (the
  shape every consumer would use) rather than via a relative path, so a
  publish-time refactor that breaks the public entry surface fails loudly
  here too.
- Mirrors how the Go services split unit vs `//go:build integration`
  suites, but in a way that is idiomatic for the TypeScript half of the
  monorepo (vitest doesn't have build tags).
