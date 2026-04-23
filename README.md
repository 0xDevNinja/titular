# Titular

Multi-chain agent launchpad. Monorepo.

## Quickstart

```bash
bash scripts/bootstrap.sh
docker compose -f infra/docker/compose.dev.yml up -d
pnpm -w turbo run build
```

## Stack

| Layer | Technology |
|-------|-----------|
| Frontend | Next.js 14, React 18, TypeScript 5 |
| Backend (API) | Go 1.22 |
| Backend (runtime) | Rust 1.80+ |
| EVM contracts | Solidity 0.8.26, Foundry |
| Solana contracts | Anchor 0.30 |
| Dev infra | Docker Compose (postgres, redis, nats, anvil, solana-test-validator) |
| CI | GitHub Actions |

## Layout

```
apps/          Next.js applications (web, console)
services/      Backend services (Go: gateway, indexer; Rust: runtime, planner, evaluator)
contracts/     Smart contracts (EVM: Foundry; Solana: Anchor)
packages/      Shared TypeScript packages
infra/         Docker Compose, Kubernetes, Terraform, Helm
scripts/       Developer scripts (bootstrap, dev, verify-repo)
docs/          Architecture decisions and runbooks
```

## Development

```bash
# Install all toolchains
bash scripts/bootstrap.sh

# Start dev stack
bash scripts/dev.sh

# Build everything
pnpm -w turbo run build
cargo build --workspace
go build github.com/0xDevNinja/titular/services/gateway-go
go build github.com/0xDevNinja/titular/services/indexer-go

# Run EVM contract tests
(cd contracts/evm && forge test)
```
