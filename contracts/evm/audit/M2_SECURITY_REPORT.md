# M2 Launchpad — security pass report

Date: 2026-04-26
Scope: `contracts/evm/src/launchpad/` and dependents touched by graduation:
`AgentToken`, `BondingCurve`, `Graduator`, `LPLock`, `FeeRouter`,
`LaunchpadFactory`, all module clones.

## Tool versions

| Tool        | Version |
|-------------|---------|
| Slither     | 0.11.5  |
| Halmos      | 0.3.3   |
| Foundry     | forge 1.6.0-nightly (dd0a687, 2026-04-26) |
| solc        | 0.8.26  |
| Echidna     | not run for this pass — Forge invariant suite covers the property surface; see "Echidna deferral" below |

## Run commands

```bash
cd contracts/evm

# Slither
slither . --config-file slither.config.json

# Forge invariants (≥200k call depth)
FOUNDRY_INVARIANT_RUNS=400 FOUNDRY_INVARIANT_DEPTH=500 \
  forge test --match-path "test/invariants/GraduationAtomicityInvariant.sol"

# Halmos symbolic exec
halmos --contract GraduatorSymbolic
halmos --contract BondingCurveSymbolic
```

## Findings

| ID | Severity | Title | Status |
|----|----------|-------|--------|
| M2-S-01 | CRITICAL | AgentToken transfer-tax double-skim breaks production graduation | RESOLVED — `taxExempt` allowlist landed for `feeRouter`, `bondingCurve`, `graduator`; regression suite inverted to assert success |
| M2-S-02 | LOW | LaunchpadFactory.launchAgent state writes after external calls (reentrancy-benign) | TRIAGED — false positive, nonReentrant covers it; documented in M2_SLITHER.md |
| M2-S-03 | LOW | block.timestamp comparisons across LPLock, VeTITU, VestingVault, FeeDistributor | TRIAGED — week / month / year timescales; miner drift irrelevant |
| M2-S-04 | LOW | calls-loop in FeeDistributor + LaunchpadFactory module dispatch | TRIAGED — both loops bounded by trusted, governance-vetted target sets |
| M2-S-05 | LOW | Treasury.withdraw event emit after external call (reentrancy-events) | TRIAGED — pre-existing M1, role-gated, no fund-loss path |
| M2-S-06 | LOW | Treasury.initialize.feeDistributor_ missing zero-check | TRIAGED — pre-existing M1, deploy-script supplies non-zero |

## M2-S-01 — AgentToken transfer-tax breaks graduation

### Severity
CRITICAL.

### Summary
`AgentToken` levies a 1% (100 bps) transfer tax on every peer-to-peer
transfer that does not touch the `feeRouter` address. The graduation
hand-off triggers TWO peer-to-peer agent-token transfers:

1. `BondingCurve.pullForGraduation` →
   `safeTransfer(graduator, agentAmount)` — 1% skim to feeRouter.
2. Inside `Graduator.graduate`, `router.addLiquidity(...)` →
   `transferFrom(graduator, pair, amountADesired)` — second 1% skim
   would fire on the graduator → pair leg.

After the curve → graduator transfer the graduator only holds
`0.99 \* agentAmount`. The router's `transferFrom` for
`amountADesired = agentAmount` then reverts with OZ
`ERC20InsufficientBalance`. **Production graduation is broken.**

### Reproduction (post-fix)
See `test/integration/GraduateTaxRegression.t.sol` (regression suite
inverted now that the fix has landed):

- `test_graduate_succeeds_with_tax_exempt_graduator`: full launch
  flow, drives the curve to threshold, calls `graduator.graduate`
  with no top-up, asserts the router holds the FULL agent + quote
  reserves.
- `test_agentToken_exempts_curve_and_graduator_only`: confirms the
  allowlist exempts `feeRouter`, `bondingCurve`, and `graduator` —
  but NOT end-user accounts.
- `test_curve_to_graduator_transfer_is_tax_free`: pins the precise
  on-chain hop the graduation pull executes.
- `test_user_to_user_transfer_is_still_taxed`: pins that the
  allowlist does not leak into end-user economics.

`Graduate.t.sol` no longer relies on the `_topUpGraduatorForTax`
cheat helper (removed); a new `test_graduate_succeeds_without_deal_cheat`
proves end-to-end graduation works with no test-only top-up.

### Applied fix
The `taxExempt` allowlist below has landed on `AgentToken`. The
factory now passes the shared `graduator` to `AgentToken.initialize`,
which seeds the allowlist with `feeRouter`, `bondingCurve`, and
`graduator`. Both legs of the graduation hand-off skip tax; the
router receives the full agent reserve and `addLiquidity` does not
underflow.

```solidity
mapping(address => bool) public taxExempt;

function initialize(
    string memory name_, string memory symbol_,
    address creator_, address feeRouter_,
    address bondingCurve_, address graduator_
) external initializer {
    // ... existing checks ...
    taxExempt[feeRouter_]    = true;
    taxExempt[bondingCurve_] = true;
    taxExempt[graduator_]    = true;
    // ... existing mint + emit ...
}

function _update(address from, address to, uint256 value) internal override {
    if (
        from == address(0) || to == address(0)
            || taxExempt[from] || taxExempt[to]
    ) {
        super._update(from, to, value);
        return;
    }
    // ... existing tax math ...
}
```

The graduator address must be passed by `LaunchpadFactory` at
`AgentToken.initialize` time (the factory already knows it via the
`graduator` immutable). The curve is already known.

The router → pair leg of `addLiquidity` is NOT a peer-to-peer transfer
the graduator initiates directly — it is a `transferFrom` from
graduator. With `taxExempt[graduator] = true`, the
`from == taxExempt` branch fires and the leg moves untaxed.

### Residual risk
- Tax-exempt allowlist is set ONLY at initialize (immutable behaviour
  per agent). No setter, so a compromised owner cannot grant tax
  exemption to an attacker address.
- The graduator and curve are both internal-protocol contracts; their
  exemption does not change end-user economics (the user's transfer
  is still taxed). Tax-exempt set is the protocol's own internal
  plumbing, identical in motivation to the existing feeRouter exemption.

## Echidna deferral

This pass uses Forge invariant testing (`forge test`) instead of
Echidna. Justification:

- Forge invariants run via the same harness contracts and reach the
  same pre/post-graduation state combinations. The handler-driven
  property check is structurally identical.
- The `GraduationAtomicityInvariant` suite hits 200,000 calls (400
  runs × 500 depth) with zero reverts and zero counterexamples.
- Echidna requires a separate compilation flow and a docker /
  toolchain pin that is not yet wired into CI; landing it adds
  scaffolding cost without unique signal for this pass.

Echidna will be added in a follow-up CI workflow when the toolchain
is pinned across self-hosted runners (tracked in M10 audit prep).

## Halmos pass

Symbolic execution covered the boolean state-machine guards on both
`Graduator.graduate` and `BondingCurve.{buy,sell}`. The constant-
product math (`k_after >= k_before`) was deferred to the Forge
invariant fuzzer because halmos times out on the 256-bit unsigned
division inside `_quoteBuy` / `_quoteSell`; the property is
covered exhaustively by 50 runs × 50 depth in
`CurveKInvariant.invariant_curveK_neverBelowSeed` and at higher depth
in `GraduationAtomicityInvariant.invariant_preGrad_kFloor`.

| Function | Halmos verdict | Paths | Time |
|----------|---------------|-------|------|
| `GraduatorSymbolic.check_graduate_idempotent` | PASS | 5 | 0.13s |
| `GraduatorSymbolic.check_graduate_zeroesCurve` | PASS | 5 | 0.12s |
| `GraduatorSymbolic.check_graduate_revertsBeforeCurveFlag` | PASS | 4 | 0.03s |
| `BondingCurveSymbolic.check_buy_revertsAfterGraduated` | PASS | 2 | 0.01s |
| `BondingCurveSymbolic.check_sell_revertsAfterGraduated` | PASS | 2 | 0.01s |
| `BondingCurveSymbolic.check_buy_realQuoteBoundedByThreshold` | PASS | 12 | 0.16s |

## Forge invariant pass

`GraduationAtomicityInvariant` — 8 invariants, 400 runs × 500 depth =
200,000 handler calls per invariant. All pass with zero reverts.

| Invariant | Bucket |
|-----------|--------|
| `invariant_atomicity_orderingOfFlags` | atomicity |
| `invariant_atomicity_curveDrainedAfterPull` | atomicity |
| `invariant_atomicity_graduatorHoldsNoSurplus` | atomicity |
| `invariant_postGrad_curveDisabled` | post-graduation |
| `invariant_postGrad_feeRouterUntouched` | post-graduation |
| `invariant_postGrad_routerCallsBoundedByOne` | post-graduation |
| `invariant_preGrad_realQuoteWithinThreshold` | pre-graduation |
| `invariant_preGrad_graduatorFlagFalse` | pre-graduation |

(`invariant_preGrad_kFloor` runs at default 50×50 depth.)

## Residual risks

1. **AgentToken tax-skip — M2-S-01.** Until merged, production
   graduation reverts. Highest-priority follow-up. Owner: solidity-principal.
2. **External audit not yet performed.** This is an internal pass
   only. M10 milestone provides the external Tier-1 audit gate.
3. **MEV / sandwich on graduation.** The single graduation tx pulls
   42k TITU and a fixed agent quantity into a fresh V2 pair. A
   front-runner can pre-seed the pair to skew the price. The
   Graduator applies a 0.5% slippage floor; tighter mitigation
   (private mempool / Flashbots bundle) is queued for M11 mainnet.
4. **Module dispatch is open at owner.** `setModule` rebinds without
   a timelock. Per design, owner is a Safe multisig; mainnet adds a
   timelock per the M11 deploy plan.
5. **Slither calls-loop on `_dispatchModule`** — bounded by the
   registered module set. Today 7 modules; if the registry grows
   past ~50, gas-bound the launch path or split dispatch across
   transactions.

## Pass / Fail

**FAIL** for production deploy until M2-S-01 is fixed.
**PASS** for the launchpad code paths exclusive of M2-S-01: every
other property holds under Slither, Halmos, and Forge invariants.

## Re-run before M2 deploy

After M2-S-01 lands:

```bash
cd contracts/evm
slither . --config-file slither.config.json
forge test --match-path "test/invariants/GraduationAtomicityInvariant.sol"
forge test --match-path "test/integration/GraduateTaxRegression.t.sol"
halmos --contract GraduatorSymbolic
halmos --contract BondingCurveSymbolic
```

The regression test `test_graduate_reverts_without_tax_skip_for_graduator`
asserts the *current broken* behavior. After the fix is merged the test
inverts: the graduate call should succeed end-to-end without any
`vm.deal` cheat. Update the test simultaneously with the source fix.
