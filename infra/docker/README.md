# infra/docker

Local development stack managed by Docker Compose.

## Services

| Service | Port | Image |
|---------|------|-------|
| postgres | 5432 | postgres:16-alpine |
| redis | 6379 | redis:7-alpine |
| nats | 4222, 8222 | nats:2.10-alpine |
| anvil (EVM) | 8545 | ghcr.io/foundry-rs/foundry:latest |
| solana-test-validator | 8899, 9900 | ghcr.io/anza-xyz/agave:v2.0.17 |

All ports bind to `127.0.0.1` (localhost only).

## Usage

```bash
# Start all services
docker compose -f infra/docker/compose.dev.yml up -d

# Check status
docker compose -f infra/docker/compose.dev.yml ps

# Tail logs
docker compose -f infra/docker/compose.dev.yml logs -f

# Stop
docker compose -f infra/docker/compose.dev.yml down

# Stop + wipe volumes
docker compose -f infra/docker/compose.dev.yml down -v
```

## Environment variables

Services use hard-coded dev credentials. Copy `.env.example` in each service
directory to `.env` for local overrides.
