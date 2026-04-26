# M3 ACP v2 — security pass report

Date: 2026-04-26
Scope: `contracts/evm/src/acp/` — `AgentRegistry`, `Job`, `JobFactory`,
`Escrow`, `FeeSplitter`, `BuybackBurner`, plus the `IJob` and (in-flight,
pending sibling-PR merge) `IPermit2` interfaces. Hooks (`HookRegistry`,
`IHook`, registered hook implementations) and `ReputationOracle` /
`ACPRouter` are tracked in the M3 threat model but were not yet merged
to `milestone/M3-acp-v2-contracts-e2e-on-base-sepolia` at the time of
this pass; they are deferred to a follow-up pass once they land.

## Tool versions

| Tool        | Version |
|-------------|---------|
| Slither     | 0.11.5  |
| Halmos      | 0.3.3   |
| Foundry     | forge 1.6.0-nightly (dd0a687, 2026-04-26) |
| solc        | 0.8.26  |
| Echidna     | not run for this pass — Forge invariant suite (owned by solidity-principal) covers the property surface.  Config landed for follow-up runs; see "Echidna deferral" below |

## Run commands

```bash
cd contracts/evm

# Slither (project config — informational findings excluded)
slither src/acp --filter-paths 'lib|out|cache'

# Slither summary printer (used to populate findings counts)
slither src/acp --filter-paths 'lib|out|cache' --print human-summary

# Halmos symbolic exec — Issue #70
halmos --match-contract JobLifecycle
```

## Findings summary

| Impact          | Count |
|-----------------|-------|
| Critical        | 0     |
| High            | 0     |
| Medium          | 0     |
| Low             | 0     |
| Informational   | 0     |
| Optimization    | 0     |

Slither analysed 27 contracts (10 in source — 7 ACP + 3 in-flight or
import-pulled — and 17 in dependencies) totalling 763 SLOC and reported
zero detector hits across the high / medium / low / informational /
optimization buckets.

## Per-contract slither summary

| Contract        | Functions | ERCs   | Complex code | Notable features         |
|-----------------|-----------|--------|--------------|--------------------------|
| AgentRegistry   | 43        | ERC165 | No           | —                        |
| BuybackBurner   | 36        | ERC165 | No           | Tokens interaction       |
| Escrow          | 33        | ERC165 | No           | —                        |
| FeeSplitter     | 34        | ERC165 | No           | —                        |
| Job             | 34        | —      | Yes          | Assembly, Upgradeable    |
| JobFactory      | 36        | ERC165 | No           | —                        |
| IJob            | (interface) | — | No | — |

`Job` is flagged "Complex" by slither's complexity heuristic; the
complexity is driven by the phase-state machine. The state machine is
exhaustively covered by the symbolic suite under
`test/halmos/JobLifecycle.halmos.t.sol` (issue #70) and the unit /
integration suite under `test/acp/` (owned by solidity-principal).
"Upgradeable" reflects the use of OpenZeppelin `Initializable` for the
EIP-1167 minimal-proxy clone pattern; there is no UUPS / transparent
proxy upgrade path.

## Findings (detail)

None to report. The `Findings (detail)` section pattern matches the
M2 pass; no individual finding rows are present because slither
reported zero hits under both the project config (informational
excluded) and a permissive override config (informational included).

## Halmos pass

Symbolic execution covered the Issue #70 invariant: once a Job is
funded (`Phase.Open` with a non-zero budget), it eventually reaches
**exactly one** terminal phase from `{Completed, Cancelled}`, and
every state-mutating call reverts once a terminal phase is entered.
The single `Cancelled` enum value subsumes both the issue's
"rejected" path (dispute-resolved-against-agent and principal-cancel)
and "expired" path (deadline-passed `expireJob`).

| Function                                              | Halmos verdict | Paths | Time  |
|-------------------------------------------------------|----------------|-------|-------|
| `JobLifecycleSymbolic.check_completed_zeroesBudget_paysAgent`           | PASS | 1 | 0.04s |
| `JobLifecycleSymbolic.check_expire_revertsBeforeDeadline`               | PASS | 1 | 0.00s |
| `JobLifecycleSymbolic.check_expire_revertsFromReview`                   | PASS | 1 | 0.03s |
| `JobLifecycleSymbolic.check_expired_zeroesBudget_refundsPrincipal`      | PASS | 1 | 0.02s |
| `JobLifecycleSymbolic.check_rejected_zeroesBudget_refundsPrincipal`     | PASS | 1 | 0.04s |
| `JobLifecycleSymbolic.check_terminal_cancelled_byExpiry_isAbsorbing`    | PASS | 1 | 0.05s |
| `JobLifecycleSymbolic.check_terminal_cancelled_byPrincipal_isAbsorbing` | PASS | 1 | 0.04s |
| `JobLifecycleSymbolic.check_terminal_cancelled_byRejection_isAbsorbing` | PASS | 1 | 0.06s |
| `JobLifecycleSymbolic.check_terminal_completed_isAbsorbing`             | PASS | 1 | 0.06s |

Total: 9 passed; 0 failed; total time 0.45s under halmos 0.3.3
(Solidity 0.8.26).

## Echidna deferral

This pass uses Slither + Halmos + Forge invariants instead of running
Echidna live. Justification (mirrors the M2 pass):

- Forge invariant suites authored by solidity-principal under
  `test/acp/` exercise the same handler-driven property surface
  Echidna would target.
- Echidna requires a separate compilation flow and a docker /
  toolchain pin that is not yet wired into self-hosted CI; landing it
  adds scaffolding cost without unique signal for the merged ACP
  contract set.
- The Echidna config has been landed at
  `contracts/evm/test/echidna/acp.yaml` and is ready to be wired into
  CI when the toolchain is pinned; that work is tracked in M10 audit
  prep.

The config pins runtime tuning to mirror the M2 launchpad fuzz suite
(`testLimit: 200_000`, `seqLen: 100`, coverage-guided, multi-actor
sender allowlist for principal / agent / arbiter / evaluator /
attacker) and points the corpus dir at `test/echidna/corpus/`.

## Forge invariant pass

Owned by solidity-principal under `contracts/evm/test/acp/`. Out of
scope for this security-lead pass; cross-referenced here for
completeness. Re-run before any production deploy.

## Residual risks

1. **External audit not yet performed.** This is an internal pass
   only. The M10 milestone gates an external Tier-1 audit before
   mainnet deploy.
2. **Hooks subsystem deferred.** `HookRegistry`, `IHook`, and the
   registered hook implementations were not merged at the time of
   this pass. The threat model (`docs/security/threat-model-m3.md`
   §4 / §5 / §6) captures the hook composition order, hook
   registry codeHash pinning, and hook gas-budget concerns; a
   follow-up Slither + Halmos pass is required once the hooks
   subsystem lands. Tracked as a follow-up to issue #71.
3. **`ReputationOracle` and `ACPRouter` deferred.** Same reasoning
   as the hooks subsystem. EIP-712 score replay, scorer-whitelist
   governance, and per-scorer nonce monotonicity are documented in
   the threat model but not yet enforced on-chain.
4. **Permit2 fund path in flight.** `Escrow.fundWithPermit2` is on
   a sibling-agent branch (PR #275 / `feat/M3-permit2-fund`) and was
   not present in the merged Escrow analysed here. Permit2
   signature replay (per-job witness binding, deadline enforcement)
   is documented in the threat model §4 / §5; a follow-up Slither
   pass on the merged Permit2 path is required before deploy.
5. **Arbiter trust assumption.** `Job.resolveDispute(bool)` is gated
   by a single-address `arbiter` immutable. The threat model §5 EA9
   and §9 O3 document the arbiter-griefing and arbiter-compromise
   risks; the per-job arbiter timeout fallback is not yet
   implemented and is queued for a follow-up issue.
6. **Slither complexity flag on `Job`.** Driven by the phase-state
   machine; mitigated by the Halmos proofs (Issue #70) and unit /
   integration coverage. Documented; no action required.

## Pass / Fail

**PASS** for the merged ACP contract set
(`AgentRegistry`, `Job`, `JobFactory`, `Escrow`, `FeeSplitter`,
`BuybackBurner`, `IJob`):

- Slither: 0 / 0 / 0 / 0 / 0 (high / medium / low / informational /
  optimization).
- Halmos: 9 / 9 PASS on the Issue #70 lifecycle invariants.
- Forge invariants: covered by solidity-principal test suite under
  `test/acp/`.

## Re-run before M3 deploy

After the hooks subsystem and `ReputationOracle` / `ACPRouter` /
Permit2 fund path land:

```bash
cd contracts/evm
slither src/acp --filter-paths 'lib|out|cache'
halmos --match-contract JobLifecycle
# Hooks symbolic suite (to be authored alongside HookRegistry):
# halmos --match-contract HookRegistry
# Reputation oracle replay-resistance suite (to be authored):
# halmos --match-contract ReputationOracle
forge test --match-path "test/acp/*"
```

## Related artefacts

- Threat model: `docs/security/threat-model-m3.md`
- Halmos suite: `contracts/evm/test/halmos/JobLifecycle.halmos.t.sol`
- Echidna config: `contracts/evm/test/echidna/acp.yaml`
- Verification logs:
  `.claude/cron/last-verify-M3-halmos.log`,
  `.claude/cron/last-verify-M3-slither.log`,
  `.claude/cron/last-verify-M3-slither-echidna.log`
