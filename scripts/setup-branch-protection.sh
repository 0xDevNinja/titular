#!/usr/bin/env bash
# setup-branch-protection.sh
# Configure branch protection rules via the GitHub API.
# Idempotent: safe to run multiple times.
#
# Usage:
#   ./scripts/setup-branch-protection.sh [--repo OWNER/REPO] [--dry-run]
#
# Requires: gh CLI authenticated with repo admin scope.
set -euo pipefail

# ---------------------------------------------------------------------------
# defaults
# ---------------------------------------------------------------------------
REPO="0xDevNinja/titular"
DRY_RUN=false

# ---------------------------------------------------------------------------
# argument parsing
# ---------------------------------------------------------------------------
while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo)
      REPO="$2"
      shift 2
      ;;
    --dry-run)
      DRY_RUN=true
      shift
      ;;
    *)
      printf 'Unknown argument: %s\n' "$1" >&2
      exit 1
      ;;
  esac
done

info()  { printf '\033[0;34m[branch-protect] %s\033[0m\n' "$*"; }
ok()    { printf '\033[0;32m[branch-protect] %s\033[0m\n' "$*"; }
dry()   { printf '\033[0;33m[branch-protect] DRY-RUN: %s\033[0m\n' "$*"; }

# ---------------------------------------------------------------------------
# required status check contexts (must match exact job names in CI)
# ---------------------------------------------------------------------------
REQUIRED_CONTEXTS='["lint","test","security","pr-check / all checks"]'

# ---------------------------------------------------------------------------
# apply_protection <branch>
# ---------------------------------------------------------------------------
apply_protection() {
  local branch="$1"
  local encoded_branch
  encoded_branch="$(python3 -c "import urllib.parse,sys; print(urllib.parse.quote(sys.argv[1],safe=''))" "$branch" 2>/dev/null || printf '%s' "$branch" | sed 's|/|%2F|g')"

  local payload
  payload="$(cat <<JSON
{
  "required_status_checks": {
    "strict": false,
    "contexts": ${REQUIRED_CONTEXTS}
  },
  "enforce_admins": false,
  "required_pull_request_reviews": null,
  "restrictions": null,
  "required_linear_history": true,
  "allow_force_pushes": false,
  "allow_deletions": false,
  "block_creations": false,
  "required_conversation_resolution": true,
  "lock_branch": false,
  "allow_fork_syncing": false
}
JSON
)"

  if "$DRY_RUN"; then
    dry "PUT /repos/${REPO}/branches/${branch}/protection"
    dry "payload: ${payload}"
    return
  fi

  info "applying protection to branch: ${branch}"
  gh api -X PUT \
    "repos/${REPO}/branches/${encoded_branch}/protection" \
    --input - <<< "$payload"
  ok "protected: ${branch}"
}

# ---------------------------------------------------------------------------
# main
# ---------------------------------------------------------------------------
info "repo: ${REPO}"
info "dry-run: ${DRY_RUN}"

# Protect main
apply_protection "main"

# Protect milestone/* branches (list current ones from API)
info "fetching milestone/* branches..."
MILESTONE_BRANCHES=$(
  gh api "repos/${REPO}/branches" --paginate \
    --jq '.[] | select(.name | startswith("milestone/")) | .name' \
    2>/dev/null || true
)

if [[ -z "$MILESTONE_BRANCHES" ]]; then
  info "no milestone/* branches found yet — skipping"
else
  while IFS= read -r branch; do
    [[ -z "$branch" ]] && continue
    apply_protection "$branch"
  done <<< "$MILESTONE_BRANCHES"
fi

ok "branch protection setup complete."
info "note: label-based merge gates are enforced via required_conversation_resolution"
info "      and external label-check workflow (not gh api branch protection)."
