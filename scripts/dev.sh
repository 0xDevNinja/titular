#!/usr/bin/env bash
# dev.sh — start the local dev stack and tail key service logs.
set -euo pipefail

COMPOSE_FILE="$(cd "$(dirname "$0")/.." && pwd)/infra/docker/compose.dev.yml"

info() { printf '\033[0;34m[dev] %s\033[0m\n' "$*"; }
ok()   { printf '\033[0;32m[dev] %s\033[0m\n' "$*"; }

if ! command -v docker &>/dev/null; then
  printf 'ERROR: docker not found. Install Docker Desktop first.\n' >&2
  exit 1
fi

info "starting dev stack..."
docker compose -f "$COMPOSE_FILE" up -d

ok "stack started. tailing logs (ctrl-c to detach)..."
docker compose -f "$COMPOSE_FILE" logs -f postgres redis nats anvil
