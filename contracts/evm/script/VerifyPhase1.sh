#!/usr/bin/env bash
# Verifies Phase-1 contracts on Basescan.
# Prereqs:
#   - BASESCAN_API_KEY in env
#   - deployments/base-sepolia.json populated (addresses non-null)
#   - jq installed
#   - forge on PATH
#
# Usage: ./script/VerifyPhase1.sh
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
JSON="${ROOT_DIR}/deployments/base-sepolia.json"
CHAIN_ID=84532
VERIFIER_URL="https://api-sepolia.basescan.org/api"

if [[ -z "${BASESCAN_API_KEY:-}" ]]; then
  echo "error: BASESCAN_API_KEY is unset." >&2
  exit 2
fi
if [[ ! -f "$JSON" ]]; then
  echo "error: ${JSON} not found." >&2
  exit 2
fi
if ! command -v jq >/dev/null; then
  echo "error: jq is required." >&2
  exit 2
fi

read_addr() {
  local key="$1"
  local val
  val="$(jq -r ".contracts.${key} // empty" "$JSON")"
  if [[ -z "$val" || "$val" == "null" ]]; then
    echo "error: deployments/base-sepolia.json contracts.${key} is null." >&2
    exit 3
  fi
  echo "$val"
}

TITU_ADDR="$(read_addr TITU)"
TREASURY_IMPL_ADDR="$(read_addr TreasuryImpl)"
TREASURY_ADDR="$(read_addr Treasury)"
VV_ADDR="$(read_addr VestingVault)"
VE_ADDR="$(read_addr VeTITU)"
FD_ADDR="$(read_addr FeeDistributor)"

cd "$ROOT_DIR"

forge verify-contract "$TITU_ADDR"            src/token/TITU.sol:TITU                        --chain-id "$CHAIN_ID" --verifier-url "$VERIFIER_URL" --etherscan-api-key "$BASESCAN_API_KEY" --watch
forge verify-contract "$TREASURY_IMPL_ADDR"   src/treasury/Treasury.sol:Treasury             --chain-id "$CHAIN_ID" --verifier-url "$VERIFIER_URL" --etherscan-api-key "$BASESCAN_API_KEY" --watch
forge verify-contract "$TREASURY_ADDR"        lib/openzeppelin-contracts/contracts/proxy/ERC1967/ERC1967Proxy.sol:ERC1967Proxy --chain-id "$CHAIN_ID" --verifier-url "$VERIFIER_URL" --etherscan-api-key "$BASESCAN_API_KEY" --watch
forge verify-contract "$VV_ADDR"              src/vault/VestingVault.sol:VestingVault        --chain-id "$CHAIN_ID" --verifier-url "$VERIFIER_URL" --etherscan-api-key "$BASESCAN_API_KEY" --watch
forge verify-contract "$VE_ADDR"              src/governance/VeTITU.sol:VeTITU               --chain-id "$CHAIN_ID" --verifier-url "$VERIFIER_URL" --etherscan-api-key "$BASESCAN_API_KEY" --watch
forge verify-contract "$FD_ADDR"              src/fees/FeeDistributor.sol:FeeDistributor     --chain-id "$CHAIN_ID" --verifier-url "$VERIFIER_URL" --etherscan-api-key "$BASESCAN_API_KEY" --watch

echo "Verification submitted for all Phase-1 contracts."
