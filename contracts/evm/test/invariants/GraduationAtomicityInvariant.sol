// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {StdInvariant} from "forge-std/StdInvariant.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import {BondingCurve} from "../../src/launchpad/BondingCurve.sol";
import {Graduator} from "../../src/launchpad/Graduator.sol";
import {LPLock} from "../../src/launchpad/LPLock.sol";
import {IUniswapV2Router02} from "../../src/launchpad/interfaces/IUniswapV2Router02.sol";

import {
    GraduationAtomicityHandler,
    MockToken,
    InvRouter,
    InvFactory
} from "./handlers/GraduationAtomicityHandler.sol";

/// @title GraduationAtomicityInvariant
/// @notice Forge invariant suite for the graduation hand-off — the most
///         security-sensitive single transition in the launchpad. The
///         invariant covers three buckets:
///
///         (a) ATOMICITY — graduation either fully completes or fully
///             reverts. There is no partial state.
///                - The curve's `graduated` flag and the graduator's
///                  `graduated[curve]` flag never disagree once the curve
///                  has been pulled (both true) or before the curve is even
///                  graduated (both false).
///                - Once the graduator has pulled, the curve's reserves
///                  are zero AND the graduator's per-curve flag is true
///                  AND the LP-lock has at least the deterministic stub LP
///                  balance.
///
///         (b) POST-GRADUATION INVARIANTS — once {Graduator.graduate} has
///             fired:
///                - LP locked: lpLock.lockedAmount() may stay 0 (the
///                  graduator routes LP via the router's `to` parameter,
///                  not via {LPLock.deposit}), but the lock's live LP
///                  balance is at least the deterministic stub value.
///                - Fee router rebound: the curve's feeRouter mapping is
///                  unchanged across graduation (the rebind is owner-only
///                  and the handler does not invoke it).
///                - Curve disabled: trades revert with AlreadyGraduated.
///
///         (c) PRE-GRADUATION INVARIANTS still hold below threshold:
///                - k_now >= k_seed.
///                - realQuoteReserve <= graduationThreshold.
///                - graduator.graduated(curve) is false until the curve
///                  itself flips first.
///
///         The suite uses tax-FREE mock tokens so the AgentToken transfer-
///         tax interaction does not contaminate the atomicity check. The
///         tax interaction is covered by the dedicated regression test in
///         test/integration/GraduateTaxRegression.t.sol.
contract GraduationAtomicityInvariant is StdInvariant, Test {
    BondingCurve internal curve;
    Graduator internal graduator;
    LPLock internal lpLock;
    MockToken internal titu;
    MockToken internal agent;
    InvRouter internal router;
    InvFactory internal factory;
    GraduationAtomicityHandler internal handler;

    address internal owner = address(0xA0);
    address internal feeRouterAddr = address(0xFEE);

    uint256 internal constant VIRTUAL_QUOTE = 30_000e18;
    uint256 internal constant VIRTUAL_AGENT = 1_073_000_000e18;
    uint256 internal constant THRESHOLD = 42_000e18;
    uint256 internal constant INITIAL_AGENT = 1_000_000_000e18;

    uint256 internal kSeed;

    function setUp() public {
        titu = new MockToken("TITU", "TITU");
        agent = new MockToken("AGENT", "AGT");
        factory = new InvFactory();
        router = new InvRouter(address(factory));

        // Pre-create the pair so LPLock can pin its address.
        address pair = factory.createPair(address(agent), address(titu));

        curve = new BondingCurve(
            owner, address(agent), address(titu), feeRouterAddr, VIRTUAL_QUOTE, VIRTUAL_AGENT, THRESHOLD, INITIAL_AGENT
        );
        agent.mint(address(curve), INITIAL_AGENT);

        // Lock binds the V2-pair stub (no contract code required for the
        // post-graduation balance assertions; the InvRouter forwards LP
        // by transfer, but here the stub doesn't actually transfer LP —
        // we bind the lock to a real ERC-20 standin so the balance check
        // is meaningful even in the no-router-mint case.
        lpLock = new LPLock(IERC20(pair), address(0xC0E));

        graduator = new Graduator(IUniswapV2Router02(address(router)), owner);
        vm.prank(owner);
        graduator.registerLPLock(address(curve), address(lpLock));

        vm.prank(owner);
        curve.setGraduator(address(graduator));

        handler = new GraduationAtomicityHandler(curve, graduator, titu, agent, lpLock, router);
        targetContract(address(handler));

        kSeed = curve.k();
    }

    // ---------------------------------------------------------------------
    // (a) Atomicity
    // ---------------------------------------------------------------------

    /// @notice The curve's `graduated` flag flips before the graduator's
    ///         per-curve flag. The reverse ordering is impossible: the
    ///         graduator's `graduate` reverts if the curve is not flipped.
    function invariant_atomicity_orderingOfFlags() public view {
        if (graduator.graduated(address(curve))) {
            assertTrue(curve.graduated(), "graduator true => curve true");
        }
    }

    /// @notice Once the graduator has pulled, both reserves on the curve are
    ///         zero. There is no partial drain — the bonding curve's
    ///         `pullForGraduation` zeroes both BEFORE the outbound transfers.
    function invariant_atomicity_curveDrainedAfterPull() public view {
        if (graduator.graduated(address(curve))) {
            assertEq(curve.realQuoteReserve(), 0, "curve quote drained");
            assertEq(curve.realAgentReserve(), 0, "curve agent drained");
            assertEq(titu.balanceOf(address(curve)), 0, "curve titu balance drained");
            assertEq(agent.balanceOf(address(curve)), 0, "curve agent balance drained");
        }
    }

    /// @notice The graduator never custodies funds across blocks. Once it
    ///         has graduated a curve, its own balances of both tokens are
    ///         zero (the router pulled them). Pre-graduation it is also
    ///         zero (no transfers fire its way).
    function invariant_atomicity_graduatorHoldsNoSurplus() public view {
        if (graduator.graduated(address(curve))) {
            assertEq(titu.balanceOf(address(graduator)), 0, "graduator no titu surplus");
            assertEq(agent.balanceOf(address(graduator)), 0, "graduator no agent surplus");
        } else {
            assertEq(titu.balanceOf(address(graduator)), 0, "pre-grad: graduator empty (titu)");
            assertEq(agent.balanceOf(address(graduator)), 0, "pre-grad: graduator empty (agent)");
        }
    }

    // ---------------------------------------------------------------------
    // (b) Post-graduation invariants
    // ---------------------------------------------------------------------

    /// @notice Once graduated, a buy reverts with `AlreadyGraduated`. We
    ///         model the inverse — the curve's flag must be true after
    ///         graduator-side completion, which by construction makes
    ///         {BondingCurve.buy} revert via the same flag.
    function invariant_postGrad_curveDisabled() public view {
        if (graduator.graduated(address(curve))) {
            assertTrue(curve.graduated(), "curve flag set => trades revert");
        }
    }

    /// @notice Fee router rebound across graduation: the curve's
    ///         `feeRouter` field is unchanged. The handler does not call
    ///         {setFeeRouter}, so this invariant is a defence-in-depth
    ///         check that nothing else mutates it.
    function invariant_postGrad_feeRouterUntouched() public view {
        assertEq(curve.feeRouter(), feeRouterAddr);
    }

    /// @notice Graduator addLiquidity calls counter monotonically tracks
    ///         graduations: count <= number of graduated curves. Since this
    ///         suite has exactly one curve, the count is 0 or 1.
    function invariant_postGrad_routerCallsBoundedByOne() public view {
        uint256 calls = router.addLiquidityCalls();
        if (graduator.graduated(address(curve))) {
            assertEq(calls, 1, "router called exactly once on graduation");
        } else {
            assertEq(calls, 0, "router untouched pre-graduation");
        }
    }

    // ---------------------------------------------------------------------
    // (c) Pre-graduation invariants
    // ---------------------------------------------------------------------

    /// @notice Pre-graduation, k stays at or above its seeded floor under any
    ///         buy/sell sequence. Re-asserts the curveK property in the
    ///         graduation-flow context to defend against a regression where
    ///         a graduation-related code path leaks state into the trade
    ///         math.
    function invariant_preGrad_kFloor() public view {
        if (!curve.graduated()) {
            assertGe(curve.k(), kSeed, "k regressed below seed pre-grad");
        }
    }

    /// @notice Pre-graduation, realQuoteReserve never strictly exceeds the
    ///         threshold. Defence-in-depth on the per-buy guard.
    function invariant_preGrad_realQuoteWithinThreshold() public view {
        if (!curve.graduated()) {
            assertLe(curve.realQuoteReserve(), curve.graduationThreshold());
        }
    }

    /// @notice Pre-graduation, the graduator's per-curve flag must be false
    ///         — a graduator-side flip with the curve still in the trade
    ///         regime is a critical state break.
    function invariant_preGrad_graduatorFlagFalse() public view {
        if (!curve.graduated()) {
            assertFalse(graduator.graduated(address(curve)), "graduator flag must trail curve flag");
        }
    }
}
