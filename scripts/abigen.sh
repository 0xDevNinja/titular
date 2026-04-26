#!/usr/bin/env bash
# Generate Go bindings for all Titular contracts.
# Idempotent: re-running produces the same output.
#
# Usage:
#   scripts/abigen.sh                   # from repo root
#   SKIP_BUILD=1 scripts/abigen.sh      # skip forge build (artifacts must exist)
#
# Requirements:
#   - forge  (foundry, prefer ~/.foundry/bin/forge)
#   - abigen (go install github.com/ethereum/go-ethereum/cmd/abigen@latest)
#   - jq

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CONTRACTS_DIR="${REPO_ROOT}/contracts/evm"
OUT_DIR="${CONTRACTS_DIR}/out"
ABI_PKG_DIR="${REPO_ROOT}/services/indexer-go/internal/abi"
PKG_NAME="abi"

# Locate forge — prefer foundry toolchain binary.
FORGE="${FOUNDRY_BIN:-${HOME}/.foundry/bin/forge}"
if [[ ! -x "${FORGE}" ]]; then
  FORGE="$(command -v forge 2>/dev/null || true)"
fi
if [[ -z "${FORGE}" || ! -x "${FORGE}" ]]; then
  echo "error: forge not found; install via https://getfoundry.sh" >&2
  exit 1
fi

# Locate abigen.
ABIGEN="${ABIGEN_BIN:-${HOME}/go/bin/abigen}"
if [[ ! -x "${ABIGEN}" ]]; then
  ABIGEN="$(command -v abigen 2>/dev/null || true)"
fi
if [[ -z "${ABIGEN}" || ! -x "${ABIGEN}" ]]; then
  echo "error: abigen not found; run: go install github.com/ethereum/go-ethereum/cmd/abigen@latest" >&2
  exit 1
fi

# Require jq.
if ! command -v jq &>/dev/null; then
  echo "error: jq not found; install via your system package manager" >&2
  exit 1
fi

# Step 1: forge build.
if [[ "${SKIP_BUILD:-0}" != "1" ]]; then
  echo "==> forge build"
  "${FORGE}" build --root "${CONTRACTS_DIR}" 2>&1
fi

if [[ ! -d "${OUT_DIR}" ]]; then
  echo "error: forge out/ directory not found at ${OUT_DIR}" >&2
  exit 1
fi

mkdir -p "${ABI_PKG_DIR}"

# Contract list: "ContractName:output_file_stem"
# Format: <ContractName> <snake_case_output_stem>
declare -a CONTRACTS=(
  # Launchpad
  "AgentToken:agent_token"
  "BondingCurve:bonding_curve"
  "LPLock:lp_lock"
  "FeeRouter:fee_router"
  "Graduator:graduator"
  "LaunchpadFactory:launchpad_factory"
  # Modules
  "AirdropModule:airdrop_module"
  "AntiSniperModule:anti_sniper_module"
  "CapitalFormationModule:capital_formation_module"
  "ExistingTokenModule:existing_token_module"
  "LaunchRadarModule:launch_radar_module"
  "PreBuyModule:pre_buy_module"
  "SixtyDaysModule:sixty_days_module"
  # M1 governance / treasury
  "TITU:titu"
  "Treasury:treasury"
  "VestingVault:vesting_vault"
  "VeTITU:ve_titu"
  "FeeDistributor:fee_distributor"
)

echo "==> generating Go bindings in ${ABI_PKG_DIR}"

for entry in "${CONTRACTS[@]}"; do
  contract="${entry%%:*}"
  stem="${entry##*:}"

  artifact="${OUT_DIR}/${contract}.sol/${contract}.json"
  if [[ ! -f "${artifact}" ]]; then
    echo "  warn: artifact not found, skipping: ${artifact}" >&2
    continue
  fi

  abi_json="$(jq -r '.abi' "${artifact}")"
  if [[ -z "${abi_json}" || "${abi_json}" == "null" ]]; then
    echo "  warn: no ABI found in ${artifact}, skipping" >&2
    continue
  fi

  out_file="${ABI_PKG_DIR}/${stem}.go"

  echo "  ${contract} -> ${stem}.go"
  echo "${abi_json}" | "${ABIGEN}" \
    --abi - \
    --pkg "${PKG_NAME}" \
    --type "${contract}" \
    --out "${out_file}"
done

echo "==> done"
