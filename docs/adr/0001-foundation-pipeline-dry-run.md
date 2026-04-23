# ADR 0001 — Foundation Pipeline Dry Run

**Status**: proposed
**Date**: 2026-04-24
**Deciders**: @0xDevNinja

---

## Context

With the foundation monorepo scaffolded (workspaces, CI, branch protection, Docker dev
stack), we need to verify that the end-to-end PR workflow operates correctly before
accepting production feature work. A controlled dry-run on a low-risk change provides
confidence that each gate fires in the expected order.

## Decision

Run the dry-run against a `fix/readme-typo` PR:

1. Create branch `fix/readme-typo` from `main`.
2. Commit a single-word typo fix in `README.md`.
3. Open a PR targeting `main`.
4. Confirm all three workflow groups trigger: `lint`, `test`, `security`.
5. Confirm `pr-check / all checks` aggregator reports green.
6. Confirm branch protection blocks direct merge until status checks pass.
7. Merge via GitHub UI (squash) once all checks are green.
8. Confirm the merged commit appears on `main` with linear history.

### Expected outcomes

| Gate | Expected result |
|------|-----------------|
| biome check | pass (no JS/TS changes) |
| go vet | pass |
| cargo clippy | pass |
| forge fmt | pass |
| forge test | 3 tests pass |
| cargo test | pass |
| gitleaks | no secrets found |
| trivy fs | no critical/high CVEs in stub code |
| cargo-audit | no advisories |
| govulncheck | no known vulnerabilities |
| pr-check aggregate | green |
| force-push attempt | blocked by branch protection |

## Consequences

- **Positive**: validates the full CI pipeline before any feature work lands.
- **Positive**: catches misconfigured required status check names early.
- **Negative**: minor: consumes one CI run per attempt.
- **Follow-up**: after dry-run passes, run `scripts/setup-branch-protection.sh`
  to apply production protection rules on `main` and any existing `milestone/*` branches.
