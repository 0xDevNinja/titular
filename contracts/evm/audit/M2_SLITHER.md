# M2 Slither pass

Run command (from `contracts/evm`):

```bash
slither . --config-file slither.config.json
```

Slither version: `0.11.5`.
Config: `contracts/evm/slither.config.json` (filters `lib|test|script`,
excludes `solc-version`, `naming-convention`, `pragma`, suppresses
informational only).

## Findings summary

| Impact | Count |
|--------|-------|
| HIGH   | 0     |
| MEDIUM | 0     |
| LOW    | 32    |

All 32 findings are LOW. No HIGH or MEDIUM. Findings break down across
five detector families.

## Findings detail

### 1. `missing-zero-check` — 1

`Treasury.initialize.feeDistributor_` lacks a zero-address check.

**Triage**: pre-existing in M1. Out of M2 scope. Acceptable: the
treasury initializer is owner-gated and is called once at deploy from a
script that supplies a non-zero feeDistributor. Filed under
`docs/security/m1-followups` for solidity-principal to address in a
follow-up if appetite exists.

### 2. `calls-loop` — 11

(a) `FeeDistributor.{_checkpointToken, _computeClaim, claimable}`: external
calls (`balanceOf`, `getPastVotes`, `getPastTotalSupply`) inside the
weekly-checkpoint loop. The loop is bounded by `currentWeek - cursor`
which is at most a few hundred entries; gas cost is bounded and the
external targets are trusted protocol contracts (the TITU + VeTITU
deployment). Acceptable.

(b) `LaunchpadFactory._dispatchModule` (called from a module-iteration
loop in `_dispatchAll`): one external call per enabled module. The
loop is bounded by the registered module set (currently 7, hard-bounded
at deploy time by `setModule`). Each module is a governance-vetted
implementation. Acceptable.

### 3. `reentrancy-benign` — 1

`LaunchpadFactory.launchAgent`: state writes (`_agentIdCounter` and
`_agents`) happen after external calls (`uniV2Factory.createPair`,
`AgentToken.initialize`, `graduator.registerLPLock`,
`curveContract.setGraduator`).

**Triage**: false positive. The function is wrapped in
`nonReentrant` (line 383). The `reentrancy-benign` detector is a
pattern check that does not see the modifier. The state writes after
the external calls are agent-record persistence, ordered AFTER the
deploy + wire steps because the record needs the deployed addresses.
A reentrant call would hit the `nonReentrant` mutex and revert.

**Action**: inline-suppressed via
`// slither-disable-next-line reentrancy-benign`
on the call site, with reason comment in source.

### 4. `reentrancy-events` — 1

`Treasury.withdraw`: `Withdrawn` event emitted after a `call{value:}`
to `to`. Pre-existing in M1.

**Triage**: false positive in practice. The call is gated by
`onlyRole(WITHDRAWER_ROLE)`. The event-after-call pattern is tolerated
when the only callable target is a privileged role. A reentrant call
into `withdraw` would re-fire the role check; multiple events would be
emitted but no fund-loss path exists. Acceptable.

### 5. `timestamp` — 18

`block.timestamp` used in comparisons across `FeeDistributor`,
`VeTITU`, `LPLock`, `VestingVault`. All such comparisons are
threshold-style (`block.timestamp < unlockTime`,
`block.timestamp >= start + cliff`, etc.). The miner-manipulability
window of 12-15s is irrelevant for week-level (`FeeDistributor`),
year-level (`VeTITU` lock end up to 4y, `LPLock` 10y), or
multi-month (`VestingVault`) timescales.

Acceptable. Slither flags every `block.timestamp` use; suppression at
the line level would add noise without value. Project policy: keep
the detector enabled, document the rationale here.

## Suppressions added in this pass

| Source line | Detector | Reason |
|-------------|----------|--------|
| `LaunchpadFactory.sol:397` | `reentrancy-benign` | function is `nonReentrant`; record write after deploy is required because addresses come from the deploy step |

## Re-run policy

- Slither runs on every PR via `.github/workflows/security.yml`.
- New HIGH or MEDIUM findings block merge.
- LOW findings accumulate in this document and triage in the next M2
  monthly review.
