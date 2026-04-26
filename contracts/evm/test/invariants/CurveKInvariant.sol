// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {StdInvariant} from "forge-std/StdInvariant.sol";
import {BondingCurve} from "../../src/launchpad/BondingCurve.sol";
import {CurveKHandler, MockToken} from "./handlers/CurveKHandler.sol";

/// @title CurveKInvariant
/// @notice Forge invariant suite for the bonding-curve constant-product law.
///         The curve uses virtual reserves so the priced product is
///         `k = (Vq + Rq) * (realAgent + (Va - initialAgent))`. Under any
///         randomised sequence of buys and sells, three properties must hold
///         for the AMM to be solvent:
///           1. `k_now >= k_initial` — the seed product is a hard floor; any
///              regression below it would imply the curve gave away inventory
///              past what its pricing function permits.
///           2. `k_now >= k_prev` after every buy — buys take the input-side
///              fee out of the pool but credit the post-fee remainder to the
///              real reserve, so `k` is monotonic non-decreasing on every buy
///              by construction (the floor on `agentOut` rounds in the pool's
///              favour).
///           3. The curve never lets `realQuoteReserve` strictly exceed
///              `graduationThreshold` mid-trade. Defence-in-depth on the
///              overshoot guard.
///         The invariant uses `targetSelector` to restrict the fuzzer to
///         buy/sell only — neither {requestGraduation} nor {pullForGraduation}
///         is in scope here. Graduation drain is the one and only path that
///         legitimately moves `k` to zero, and is covered by the
///         {BondingCurve.t.sol} unit suite.
contract CurveKInvariant is StdInvariant, Test {
    BondingCurve internal curve;
    MockToken internal titu;
    MockToken internal agent;
    CurveKHandler internal handler;

    address internal owner = address(0xA0);
    address internal feeRouter = address(0xFEE);

    uint256 internal constant VIRTUAL_QUOTE = 30_000e18;
    uint256 internal constant VIRTUAL_AGENT = 1_073_000_000e18;
    uint256 internal constant THRESHOLD = 42_000e18;
    uint256 internal constant INITIAL_AGENT = 1_000_000_000e18;

    /// @notice Seeded `k` at construction. Persistent floor for the run.
    uint256 internal kSeed;

    function setUp() public {
        titu = new MockToken("TITU", "TITU");
        agent = new MockToken("AGENT", "AGT");

        curve = new BondingCurve(
            owner, address(agent), address(titu), feeRouter, VIRTUAL_QUOTE, VIRTUAL_AGENT, THRESHOLD, INITIAL_AGENT
        );
        agent.mint(address(curve), INITIAL_AGENT);

        handler = new CurveKHandler(curve, titu, agent);

        // Restrict fuzzer to buy/sell only; graduation is not in scope here.
        bytes4[] memory selectors = new bytes4[](2);
        selectors[0] = CurveKHandler.buy.selector;
        selectors[1] = CurveKHandler.sell.selector;
        targetSelector(FuzzSelector({addr: address(handler), selectors: selectors}));
        targetContract(address(handler));

        kSeed = curve.k();
    }

    /// @notice `k` must never drop below its seeded floor under any sequence
    ///         of buy/sell calls. This is the canonical AMM constant-product
    ///         invariant — output amounts are floored on integer division so
    ///         each trade leaves the pool with at least the pre-trade product.
    function invariant_curveK_neverBelowSeed() public view {
        uint256 kNow = curve.k();
        assertGe(kNow, kSeed, "k regressed below seed");
    }

    /// @notice Real quote reserve is bounded above by the graduation threshold.
    ///         The buy path explicitly reverts when the post-fee remainder
    ///         would push past the threshold; this invariant catches any
    ///         future regression of that guard.
    function invariant_curveK_realQuoteBoundedByThreshold() public view {
        assertLe(curve.realQuoteReserve(), curve.graduationThreshold());
    }

    /// @notice Pre-graduation, the curve's physical agent balance equals
    ///         `realAgentReserve` exactly. The mocks levy no transfer tax, so
    ///         any drift would indicate a bookkeeping error in {buy} / {sell}.
    function invariant_curveK_agentReserveMatchesBalance() public view {
        assertEq(agent.balanceOf(address(curve)), curve.realAgentReserve());
    }

    /// @notice Pre-graduation, the curve's TITU balance is at least
    ///         `realQuoteReserve`. It can be strictly greater because fees
    ///         are pushed to `feeRouter` in the same tx but external transfers
    ///         could in principle drift the balance up; however with our
    ///         tax-free mock and CEI fees, equality must hold.
    function invariant_curveK_quoteReserveMatchesBalance() public view {
        assertEq(titu.balanceOf(address(curve)), curve.realQuoteReserve());
    }

    /// @notice The curve must never flip `graduated` while the threshold is
    ///         not yet reached. The handler does not call `requestGraduation`
    ///         so the flag should stay `false` for the entire run.
    function invariant_curveK_notGraduatedUnderTradeOnly() public view {
        assertFalse(curve.graduated());
    }
}
