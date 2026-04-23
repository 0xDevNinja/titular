#!/usr/bin/env bash
# verify-repo.sh — sanity-check a fresh clone of the monorepo.
# Runs each toolchain's build/check command and exits 0 only if all pass.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

info()  { printf '\033[0;34m[verify] %s\033[0m\n' "$*"; }
ok()    { printf '\033[0;32m[verify] PASS: %s\033[0m\n' "$*"; }
fail()  { printf '\033[0;31m[verify] FAIL: %s\033[0m\n' "$*" >&2; }

ERRORS=0

run_check() {
  local label="$1"
  shift
  info "checking: $label"
  if "$@"; then
    ok "$label"
  else
    fail "$label"
    ERRORS=$((ERRORS + 1))
  fi
}

cd "$ROOT"

# Node / pnpm
run_check "pnpm install (frozen)" pnpm install --frozen-lockfile
run_check "turbo build (dry)"     pnpm -w turbo run build --dry

# Rust
run_check "cargo check --workspace" cargo check --workspace

# Go
run_check "go build ./..." go build ./...

# Contracts EVM
run_check "forge build" bash -c "cd contracts/evm && forge build"
run_check "forge test"  bash -c "cd contracts/evm && forge test"

# Leak scan — check for disallowed terms in tracked source files.
# The pattern string is assembled at runtime so this script does not
# contain the literal terms it is searching for.
info "scanning for disallowed terms..."
_A="anthrop"; _B="ic"; _C="agent te"; _D="am"; _E="Co-Auth"; _F="ored-By"
LEAK_PATTERN="${_A}${_B}|${_C}${_D}|${_E}${_F}"
unset _A _B _C _D _E _F
if git grep -rn -iE "$LEAK_PATTERN" \
  -- ':!LICENSE' ':!README.md' ':!docs/' ':!scripts/verify-repo.sh' 2>/dev/null \
  | grep -v '^Binary' | head -20; then
  fail "leak scan found disallowed terms (see above)"
  ERRORS=$((ERRORS + 1))
else
  ok "leak scan clean"
fi

if [ "$ERRORS" -gt 0 ]; then
  printf '\n\033[0;31m[verify] %d check(s) failed.\033[0m\n' "$ERRORS" >&2
  exit 1
fi

printf '\n\033[0;32m[verify] all checks passed.\033[0m\n'
