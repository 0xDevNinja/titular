#!/usr/bin/env bash
# Regenerate Go bindings for every M2 + M3 contract the indexer subscribes to.
#
# Usage:
#   ./scripts/abigen.sh            # regenerate (runs forge build first)
#   SKIP_FORGE_BUILD=1 ./scripts/abigen.sh
#
# Env vars:
#   ABIGEN_VERSION   pinned go-ethereum tag (defaults below)
#   FORGE_OUT        path to forge build output (defaults to ../../contracts/evm/out)
#   ABI_OUT          target package dir (defaults to ./internal/abi)
#   SKIP_FORGE_BUILD set to 1 to skip `forge build`
#
# The script is idempotent: re-running on a clean tree must produce no diff.
# CI uses this property as a "stale bindings" gate.

set -euo pipefail

# -----------------------------------------------------------------------------
# Pinned tool version. Bump in lockstep with go-ethereum in go.mod so the
# generated bindings stay ABI-compatible with the runtime library.
# -----------------------------------------------------------------------------
ABIGEN_VERSION="${ABIGEN_VERSION:-v1.17.2}"

# -----------------------------------------------------------------------------
# Paths.
# -----------------------------------------------------------------------------
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVICE_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
REPO_ROOT="$(cd "${SERVICE_DIR}/../.." && pwd)"

FORGE_OUT="${FORGE_OUT:-${REPO_ROOT}/contracts/evm/out}"
ABI_OUT="${ABI_OUT:-${SERVICE_DIR}/internal/abi}"

mkdir -p "${ABI_OUT}"

# -----------------------------------------------------------------------------
# Locate (or install) abigen at the pinned version.
# -----------------------------------------------------------------------------
GOBIN="$(go env GOBIN)"
if [[ -z "${GOBIN}" ]]; then
  GOBIN="$(go env GOPATH)/bin"
fi

ABIGEN_BIN="${GOBIN}/abigen"
need_install=1
if [[ -x "${ABIGEN_BIN}" ]]; then
  if "${ABIGEN_BIN}" --version 2>/dev/null | grep -q "${ABIGEN_VERSION#v}"; then
    need_install=0
  fi
fi

if [[ "${need_install}" -eq 1 ]]; then
  echo "[abigen] installing abigen@${ABIGEN_VERSION} into ${GOBIN}"
  GOBIN="${GOBIN}" go install "github.com/ethereum/go-ethereum/cmd/abigen@${ABIGEN_VERSION}"
fi

echo "[abigen] using $("${ABIGEN_BIN}" --version)"

# -----------------------------------------------------------------------------
# Ensure forge artifacts exist.
# -----------------------------------------------------------------------------
if [[ "${SKIP_FORGE_BUILD:-0}" != "1" ]]; then
  if ! command -v forge >/dev/null 2>&1; then
    echo "[abigen] forge not found on PATH; install foundry or set SKIP_FORGE_BUILD=1" >&2
    exit 1
  fi
  echo "[abigen] running forge build"
  ( cd "${REPO_ROOT}/contracts/evm" && forge build --skip test )
fi

if ! command -v jq >/dev/null 2>&1; then
  echo "[abigen] jq is required to extract ABI from forge artifacts" >&2
  exit 1
fi

# -----------------------------------------------------------------------------
# Contract → binding mapping.
#
# Format: <ContractName>:<go_type>:<go_filename_without_ext>
#   - ContractName  matches contracts/evm/out/<Name>.sol/<Name>.json
#   - go_type       Go struct name produced by abigen (also used as MetaData prefix)
#   - go_filename   target file under ${ABI_OUT}
#
# Keep this list aligned with internal/abi/abi_test.go::TestAllBindingsCompile.
# Interfaces (I*.sol) intentionally excluded — only concrete deployable
# contracts that the indexer subscribes to.
# -----------------------------------------------------------------------------
CONTRACTS=(
  # M2 — launchpad
  "LaunchpadFactory:LaunchpadFactory:launchpad_factory"
  "BondingCurve:BondingCurve:bonding_curve"
  "AgentToken:AgentToken:agent_token"
  "FeeRouter:FeeRouter:fee_router"
  "Graduator:Graduator:graduator"
  "LPLock:LPLock:lp_lock"
  "AirdropModule:AirdropModule:airdrop_module"
  "AntiSniperModule:AntiSniperModule:anti_sniper_module"
  "CapitalFormationModule:CapitalFormationModule:capital_formation_module"
  "ExistingTokenModule:ExistingTokenModule:existing_token_module"
  "LaunchRadarModule:LaunchRadarModule:launch_radar_module"
  "PreBuyModule:PreBuyModule:pre_buy_module"
  "SixtyDaysModule:SixtyDaysModule:sixty_days_module"

  # M1 — token, treasury, governance, vesting, fees
  "TITU:TITU:titu"
  "Treasury:Treasury:treasury"
  "VeTITU:VeTITU:ve_titu"
  "VestingVault:VestingVault:vesting_vault"
  "FeeDistributor:FeeDistributor:fee_distributor"

  # M3 — ACP v2
  "AgentRegistry:AgentRegistry:acp_agent_registry"
  "Job:Job:acp_job"
  "JobFactory:JobFactory:acp_job_factory"
  "Escrow:Escrow:acp_escrow"
  "FeeSplitter:FeeSplitter:acp_fee_splitter"
  "BuybackBurner:BuybackBurner:acp_buyback_burner"
  "HookRegistry:HookRegistry:acp_hook_registry"
  "FundTransferHook:FundTransferHook:acp_fund_transfer_hook"
  "SubscriptionHook:SubscriptionHook:acp_subscription_hook"
  "MilestoneHook:MilestoneHook:acp_milestone_hook"
  "RoyaltyHook:RoyaltyHook:acp_royalty_hook"
)

# -----------------------------------------------------------------------------
# Generate one binding from a forge artifact.
# -----------------------------------------------------------------------------
gen_one() {
  local contract="$1" go_type="$2" go_file="$3"
  local artifact="${FORGE_OUT}/${contract}.sol/${contract}.json"
  local out="${ABI_OUT}/${go_file}.go"

  if [[ ! -f "${artifact}" ]]; then
    echo "[abigen] missing artifact ${artifact} (run forge build)" >&2
    return 1
  fi

  # abigen reads ABI from stdin when --abi=- is given; pass JSON array only.
  jq -c '.abi' "${artifact}" \
    | "${ABIGEN_BIN}" --abi=- --pkg=abi --type="${go_type}" --out="${out}"

  echo "[abigen] ${contract} -> ${out}"
}

# -----------------------------------------------------------------------------
# Drive the table.
# -----------------------------------------------------------------------------
fail=0
for entry in "${CONTRACTS[@]}"; do
  IFS=':' read -r contract go_type go_file <<< "${entry}"
  if ! gen_one "${contract}" "${go_type}" "${go_file}"; then
    fail=1
  fi
done

if [[ "${fail}" -ne 0 ]]; then
  echo "[abigen] one or more bindings failed to generate" >&2
  exit 1
fi

echo "[abigen] done"
