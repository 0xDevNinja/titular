#!/usr/bin/env bash
# Regenerate the gateway OpenAPI artefacts from the swag annotations.
#
# Pipeline:
#
#   1. swag init                 — Swagger 2.0 → docs/swagger.{json,yaml}
#   2. go run ./scripts/openapi-gen — convert + patch → docs/openapi.{json,yaml}
#
# Invoked by:
#
#   * developers, after editing a swag annotation
#   * the gateway-openapi-check CI job, which then runs
#     `git diff --exit-code -- services/gateway-go/docs/`
#
# Determinism: both steps are byte-stable across runs given identical
# input. swag emits sorted JSON / YAML; openapi-gen marshals via
# encoding/json (which sorts map keys) and routes YAML through the
# same JSON object. A drifted PR will produce a non-empty diff and
# fail CI.
#
# Requirements:
#
#   * Go 1.25+
#   * `swag` CLI on PATH OR installed at $(go env GOPATH)/bin/swag.
#     The script will install v1.16.4 in a build-cache directory if
#     swag is not already available.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVICE_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

cd "${SERVICE_DIR}"

# Resolve a swag binary. Prefer one already on PATH; fall back to a
# pinned install so CI doesn't drift from local dev.
SWAG_VERSION="v1.16.4"
SWAG_BIN="$(command -v swag || true)"
if [[ -z "${SWAG_BIN}" ]]; then
  GOPATH_BIN="$(go env GOPATH)/bin"
  if [[ ! -x "${GOPATH_BIN}/swag" ]]; then
    GOBIN="${GOPATH_BIN}" go install "github.com/swaggo/swag/cmd/swag@${SWAG_VERSION}"
  fi
  SWAG_BIN="${GOPATH_BIN}/swag"
fi

# Step 1 — Swagger 2.0 from the annotation comments. The handler
# packages, the auth packages, and the openapi response-type aggregator
# are all in --dir so swag can resolve every $ref.
"${SWAG_BIN}" init \
  -g cmd/gateway/main.go \
  -d ./,./internal/handlers,./internal/auth,./internal/router,./internal/openapi,./internal/store \
  -o ./docs \
  --outputTypes json,yaml \
  --quiet

# Step 2 — Convert + scalar-format patch.
go run ./scripts/openapi-gen
