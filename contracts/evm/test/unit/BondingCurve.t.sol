// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {BondingCurve} from "../../src/launchpad/BondingCurve.sol";

/// @dev Minimal mintable mock standing in for TITU and the agent token.
///      No transfer tax — keeps the bonding-curve arithmetic exact so we can
///      assert reserve deltas down to the wei without fighting tax bookkeeping.
contract MockERC20 is ERC20 {
    constructor(string memory n, string memory s) ERC20(n, s) {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

contract BondingCurveTest is Test {
    MockERC20 internal titu;
    MockERC20 internal agent;
    BondingCurve internal curve;

    address internal owner = address(0xA0);
    address internal feeRouter = address(0xFEE);
    address internal grad = address(0xBEEF);
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);

    uint256 internal constant VIRTUAL_QUOTE = 30_000e18;
    uint256 internal constant VIRTUAL_AGENT = 1_073_000_000e18;
    uint256 internal constant THRESHOLD = 42_000e18;
    uint256 internal constant INITIAL_AGENT = 1_000_000_000e18;

    function setUp() public virtual {
        titu = new MockERC20("TITU", "TITU");
        agent = new MockERC20("AGENT", "AGT");

        curve = new BondingCurve(
            owner, address(agent), address(titu), feeRouter, VIRTUAL_QUOTE, VIRTUAL_AGENT, THRESHOLD, INITIAL_AGENT
        );

        // Seed the curve with the agent token inventory it claims to hold.
        agent.mint(address(curve), INITIAL_AGENT);

        // Fund traders with TITU; they approve the curve up front.
        titu.mint(alice, 1_000_000e18);
        titu.mint(bob, 1_000_000e18);
        vm.prank(alice);
        titu.approve(address(curve), type(uint256).max);
        vm.prank(bob);
        titu.approve(address(curve), type(uint256).max);
    }

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    function test_constructor_sets_immutables() public view {
        assertEq(curve.agentToken(), address(agent));
        assertEq(curve.quoteToken(), address(titu));
        assertEq(curve.virtualQuoteReserve(), VIRTUAL_QUOTE);
        assertEq(curve.virtualAgentReserve(), VIRTUAL_AGENT);
        assertEq(curve.graduationThreshold(), THRESHOLD);
        assertEq(curve.realAgentReserve(), INITIAL_AGENT);
        assertEq(curve.realQuoteReserve(), 0);
        assertEq(curve.feeRouter(), feeRouter);
        assertEq(curve.owner(), owner);
        assertEq(curve.graduator(), address(0));
        assertFalse(curve.graduated());
        assertEq(uint256(curve.FEE_BPS()), 100);
    }

    function test_constructor_reverts_on_zero_address() public {
        vm.expectRevert(BondingCurve.ZeroAddress.selector);
        new BondingCurve(
            address(0), address(agent), address(titu), feeRouter, VIRTUAL_QUOTE, VIRTUAL_AGENT, THRESHOLD, INITIAL_AGENT
        );
    }

    function test_constructor_reverts_on_zero_value() public {
        vm.expectRevert(BondingCurve.ZeroValue.selector);
        new BondingCurve(owner, address(agent), address(titu), feeRouter, 0, VIRTUAL_AGENT, THRESHOLD, INITIAL_AGENT);
    }

    function test_constructor_reverts_on_zero_agent_token() public {
        vm.expectRevert(BondingCurve.ZeroAddress.selector);
        new BondingCurve(
            owner, address(0), address(titu), feeRouter, VIRTUAL_QUOTE, VIRTUAL_AGENT, THRESHOLD, INITIAL_AGENT
        );
    }

    function test_constructor_reverts_on_zero_quote_token() public {
        vm.expectRevert(BondingCurve.ZeroAddress.selector);
        new BondingCurve(
            owner, address(agent), address(0), feeRouter, VIRTUAL_QUOTE, VIRTUAL_AGENT, THRESHOLD, INITIAL_AGENT
        );
    }

    function test_constructor_reverts_on_zero_fee_router() public {
        vm.expectRevert(BondingCurve.ZeroAddress.selector);
        new BondingCurve(
            owner, address(agent), address(titu), address(0), VIRTUAL_QUOTE, VIRTUAL_AGENT, THRESHOLD, INITIAL_AGENT
        );
    }

    function test_constructor_reverts_zero_virtual_agent() public {
        vm.expectRevert(BondingCurve.ZeroValue.selector);
        new BondingCurve(owner, address(agent), address(titu), feeRouter, VIRTUAL_QUOTE, 0, THRESHOLD, INITIAL_AGENT);
    }

    function test_constructor_reverts_zero_threshold() public {
        vm.expectRevert(BondingCurve.ZeroValue.selector);
        new BondingCurve(
            owner, address(agent), address(titu), feeRouter, VIRTUAL_QUOTE, VIRTUAL_AGENT, 0, INITIAL_AGENT
        );
    }

    function test_constructor_reverts_zero_initial_agent() public {
        vm.expectRevert(BondingCurve.ZeroValue.selector);
        new BondingCurve(owner, address(agent), address(titu), feeRouter, VIRTUAL_QUOTE, VIRTUAL_AGENT, THRESHOLD, 0);
    }

    function test_constructor_reverts_invalid_reserves_equal() public {
        // virtualAgent == initialAgent fails the strict-greater check.
        vm.expectRevert(BondingCurve.InvalidReserves.selector);
        new BondingCurve(
            owner, address(agent), address(titu), feeRouter, VIRTUAL_QUOTE, INITIAL_AGENT, THRESHOLD, INITIAL_AGENT
        );
    }

    function test_constructor_reverts_invalid_reserves_less_than() public {
        vm.expectRevert(BondingCurve.InvalidReserves.selector);
        new BondingCurve(
            owner, address(agent), address(titu), feeRouter, VIRTUAL_QUOTE, INITIAL_AGENT - 1, THRESHOLD, INITIAL_AGENT
        );
    }

    // ---------------------------------------------------------------------
    // Quote math
    // ---------------------------------------------------------------------

    function test_quoteBuy_formula_matches_manual_calc() public view {
        // First trade: X = VIRTUAL_QUOTE, Y_eff = VIRTUAL_AGENT (because real equals initial at launch).
        uint256 quoteIn = 1000e18;
        uint256 expectedFee = (quoteIn * 100) / 10_000; // 10e18
        uint256 dxPrime = quoteIn - expectedFee;
        uint256 x = VIRTUAL_QUOTE;
        uint256 y = VIRTUAL_AGENT;
        uint256 expectedOut = (y * dxPrime) / (x + dxPrime);

        (uint256 agentOut, uint256 fee) = curve.quoteBuy(quoteIn);
        assertEq(fee, expectedFee, "fee mismatch");
        assertEq(agentOut, expectedOut, "agent out mismatch");
    }

    function test_quoteBuy_returns_zero_for_zero_input() public view {
        (uint256 agentOut, uint256 fee) = curve.quoteBuy(0);
        assertEq(agentOut, 0, "agentOut should be 0");
        assertEq(fee, 0, "fee should be 0");
    }

    function test_quoteSell_returns_zero_for_zero_input() public view {
        (uint256 quoteOut, uint256 fee) = curve.quoteSell(0);
        assertEq(quoteOut, 0, "quoteOut should be 0");
        assertEq(fee, 0, "fee should be 0");
    }

    function test_quoteSell_formula_matches_manual_calc() public {
        // Seed state with a prior buy so we can sell.
        vm.prank(alice);
        curve.buy(0, 500e18);

        uint256 agentIn = 1_000_000e18;
        uint256 offset = VIRTUAL_AGENT - INITIAL_AGENT;
        uint256 x = VIRTUAL_QUOTE + curve.realQuoteReserve();
        uint256 y = curve.realAgentReserve() + offset;
        uint256 yPlus = y + agentIn;
        uint256 grossOut = (x * agentIn) / yPlus;
        // Fee folds the bps step into one division to avoid double-flooring.
        uint256 expectedFee = (x * agentIn * 100) / (yPlus * 10_000);
        uint256 expectedNet = grossOut - expectedFee;

        (uint256 quoteOut, uint256 fee) = curve.quoteSell(agentIn);
        assertEq(fee, expectedFee, "fee mismatch");
        assertEq(quoteOut, expectedNet, "quote out mismatch");
    }

    // ---------------------------------------------------------------------
    // Buy
    // ---------------------------------------------------------------------

    function test_buy_updates_reserves_and_transfers() public {
        uint256 quoteIn = 1000e18;
        (uint256 expectedOut, uint256 expectedFee) = curve.quoteBuy(quoteIn);

        uint256 feeRouterBefore = titu.balanceOf(feeRouter);
        uint256 aliceAgentBefore = agent.balanceOf(alice);

        vm.prank(alice);
        curve.buy(expectedOut, quoteIn);

        assertEq(curve.realQuoteReserve(), quoteIn - expectedFee, "real quote reserve");
        assertEq(curve.realAgentReserve(), INITIAL_AGENT - expectedOut, "real agent reserve");
        assertEq(titu.balanceOf(feeRouter), feeRouterBefore + expectedFee, "fee routed");
        assertEq(agent.balanceOf(alice), aliceAgentBefore + expectedOut, "agent delivered");
        assertEq(titu.balanceOf(address(curve)), quoteIn - expectedFee, "pool quote balance");
    }

    function test_buy_reverts_if_min_out_not_met() public {
        (uint256 expectedOut,) = curve.quoteBuy(1000e18);
        vm.prank(alice);
        vm.expectRevert(BondingCurve.InsufficientOutput.selector);
        curve.buy(expectedOut + 1, 1000e18);
    }

    function test_buy_reverts_on_zero_amount() public {
        vm.prank(alice);
        vm.expectRevert(BondingCurve.ZeroAmount.selector);
        curve.buy(0, 0);
    }

    function test_buy_reverts_after_graduation() public {
        _graduate();
        vm.prank(alice);
        vm.expectRevert(BondingCurve.AlreadyGraduated.selector);
        curve.buy(0, 1000e18);
    }

    function test_buy_cannot_cross_graduation_threshold_in_one_tx() public {
        // With 1% fee, to land real quote on THRESHOLD the gross quoteIn must be
        // THRESHOLD * 10000 / 9900. Pushing strictly past should revert.
        uint256 gross = (THRESHOLD * 10_000) / 9900 + 1e18; // well past
        vm.prank(alice);
        vm.expectRevert(BondingCurve.ExceedsGraduationThreshold.selector);
        curve.buy(0, gross);
    }

    function test_buy_can_exactly_hit_threshold_then_graduate() public {
        // Gross that lands realQuoteReserve exactly on THRESHOLD: gross * 0.99 = THRESHOLD.
        uint256 gross = (THRESHOLD * 10_000) / 9900;
        // Round up by one wei if needed so integer floor of fee still leaves >= threshold.
        uint256 feePreview = (gross * 100) / 10_000;
        if (gross - feePreview < THRESHOLD) gross += 1;

        titu.mint(alice, gross);
        vm.prank(alice);
        curve.buy(0, gross);
        assertGe(curve.realQuoteReserve(), THRESHOLD);
        assertLe(curve.realQuoteReserve(), THRESHOLD + 1);

        // Anyone can now call requestGraduation.
        curve.requestGraduation();
        assertTrue(curve.graduated());
    }

    // ---------------------------------------------------------------------
    // Sell
    // ---------------------------------------------------------------------

    function test_sell_updates_reserves() public {
        // First, buy to produce agent inventory on alice.
        vm.prank(alice);
        curve.buy(0, 1000e18);

        uint256 agentIn = agent.balanceOf(alice) / 2;
        vm.prank(alice);
        agent.approve(address(curve), agentIn);

        (uint256 expectedNet, uint256 expectedFee) = curve.quoteSell(agentIn);

        uint256 feeRouterBefore = titu.balanceOf(feeRouter);
        uint256 aliceQuoteBefore = titu.balanceOf(alice);
        uint256 quoteReserveBefore = curve.realQuoteReserve();
        uint256 agentReserveBefore = curve.realAgentReserve();

        vm.prank(alice);
        curve.sell(agentIn, expectedNet);

        assertEq(curve.realQuoteReserve(), quoteReserveBefore - (expectedNet + expectedFee), "quote reserve delta");
        assertEq(curve.realAgentReserve(), agentReserveBefore + agentIn, "agent reserve delta");
        assertEq(titu.balanceOf(feeRouter), feeRouterBefore + expectedFee, "fee routed");
        assertEq(titu.balanceOf(alice), aliceQuoteBefore + expectedNet, "seller paid");
    }

    function test_sell_reverts_if_min_quote_not_met() public {
        vm.prank(alice);
        curve.buy(0, 1000e18);
        uint256 agentIn = agent.balanceOf(alice) / 2;
        vm.prank(alice);
        agent.approve(address(curve), agentIn);
        (uint256 net,) = curve.quoteSell(agentIn);
        vm.prank(alice);
        vm.expectRevert(BondingCurve.InsufficientOutput.selector);
        curve.sell(agentIn, net + 1);
    }

    function test_sell_reverts_on_zero_amount() public {
        vm.prank(alice);
        vm.expectRevert(BondingCurve.ZeroAmount.selector);
        curve.sell(0, 0);
    }

    function test_sell_reverts_after_graduation() public {
        // Drive to graduation first, then cheat a single unit of agent inventory
        // into Alice's wallet so we can exercise the graduated-state revert.
        _graduate();
        agent.mint(alice, 1);
        vm.prank(alice);
        agent.approve(address(curve), 1);
        vm.prank(alice);
        vm.expectRevert(BondingCurve.AlreadyGraduated.selector);
        curve.sell(1, 0);
    }

    // ---------------------------------------------------------------------
    // Graduation lifecycle
    // ---------------------------------------------------------------------

    function test_requestGraduation_reverts_below_threshold() public {
        vm.expectRevert(BondingCurve.BelowThreshold.selector);
        curve.requestGraduation();
    }

    function test_requestGraduation_succeeds_at_threshold() public {
        _reachThreshold();

        vm.expectEmit(false, false, false, false, address(curve));
        emit BondingCurve.Graduated(0, 0, 0);
        curve.requestGraduation();
        assertTrue(curve.graduated());
    }

    function test_requestGraduation_reverts_if_already_graduated() public {
        _graduate();
        vm.expectRevert(BondingCurve.AlreadyGraduated.selector);
        curve.requestGraduation();
    }

    function test_pullForGraduation_reverts_before_graduated() public {
        vm.prank(owner);
        curve.setGraduator(grad);
        vm.prank(grad);
        vm.expectRevert(BondingCurve.NotGraduated.selector);
        curve.pullForGraduation();
    }

    function test_pullForGraduation_onlyGraduator_after_graduated() public {
        _graduate();
        // Graduator is set inside _graduate. Any other caller rejected.
        vm.prank(alice);
        vm.expectRevert(BondingCurve.NotGraduator.selector);
        curve.pullForGraduation();
    }

    function test_pullForGraduation_drains_reserves() public {
        _graduate();
        uint256 expectedQuote = curve.realQuoteReserve();
        uint256 expectedAgent = curve.realAgentReserve();

        uint256 gradQuoteBefore = titu.balanceOf(grad);
        uint256 gradAgentBefore = agent.balanceOf(grad);

        vm.prank(grad);
        (uint256 q, uint256 a) = curve.pullForGraduation();

        assertEq(q, expectedQuote);
        assertEq(a, expectedAgent);
        assertEq(curve.realQuoteReserve(), 0);
        assertEq(curve.realAgentReserve(), 0);
        assertEq(titu.balanceOf(grad), gradQuoteBefore + expectedQuote);
        assertEq(agent.balanceOf(grad), gradAgentBefore + expectedAgent);
    }

    function test_pullForGraduation_is_idempotent_on_empty_balances() public {
        _graduate();
        vm.prank(grad);
        curve.pullForGraduation();
        // Second call returns zeros; nothing reverts because reserves are already 0.
        vm.prank(grad);
        (uint256 q, uint256 a) = curve.pullForGraduation();
        assertEq(q, 0);
        assertEq(a, 0);
    }

    // ---------------------------------------------------------------------
    // Admin
    // ---------------------------------------------------------------------

    function test_setGraduator_once_only() public {
        vm.prank(owner);
        vm.expectEmit(true, true, false, false, address(curve));
        emit BondingCurve.GraduatorSet(address(0), grad);
        curve.setGraduator(grad);
        assertEq(curve.graduator(), grad);

        vm.prank(owner);
        vm.expectRevert(BondingCurve.GraduatorAlreadySet.selector);
        curve.setGraduator(address(0xDEAD));
    }

    function test_setGraduator_reverts_for_non_owner() public {
        vm.prank(alice);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        curve.setGraduator(grad);
    }

    function test_setGraduator_reverts_on_zero() public {
        vm.prank(owner);
        vm.expectRevert(BondingCurve.ZeroAddress.selector);
        curve.setGraduator(address(0));
    }

    function test_setFeeRouter_rebindable() public {
        address newRouter = address(0xFE2);
        vm.prank(owner);
        vm.expectEmit(true, true, false, false, address(curve));
        emit BondingCurve.FeeRouterSet(feeRouter, newRouter);
        curve.setFeeRouter(newRouter);
        assertEq(curve.feeRouter(), newRouter);
    }

    function test_setFeeRouter_reverts_for_non_owner() public {
        vm.prank(alice);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        curve.setFeeRouter(address(0xFE2));
    }

    function test_setFeeRouter_reverts_on_zero() public {
        vm.prank(owner);
        vm.expectRevert(BondingCurve.ZeroAddress.selector);
        curve.setFeeRouter(address(0));
    }

    // ---------------------------------------------------------------------
    // Fee routing
    // ---------------------------------------------------------------------

    function test_fee_transferred_to_feeRouter_on_buy() public {
        uint256 quoteIn = 500e18;
        (, uint256 expectedFee) = curve.quoteBuy(quoteIn);
        uint256 before = titu.balanceOf(feeRouter);
        vm.prank(alice);
        curve.buy(0, quoteIn);
        assertEq(titu.balanceOf(feeRouter) - before, expectedFee);
    }

    function test_fee_transferred_to_feeRouter_on_sell() public {
        vm.prank(alice);
        curve.buy(0, 500e18);
        uint256 agentIn = agent.balanceOf(alice) / 3;
        vm.prank(alice);
        agent.approve(address(curve), agentIn);
        (, uint256 expectedFee) = curve.quoteSell(agentIn);
        uint256 before = titu.balanceOf(feeRouter);
        vm.prank(alice);
        curve.sell(agentIn, 0);
        assertEq(titu.balanceOf(feeRouter) - before, expectedFee);
    }

    // ---------------------------------------------------------------------
    // Fuzz: k non-decreasing across trades
    // ---------------------------------------------------------------------

    /// @dev Fuzzes a buy-then-sell sequence with bounded inputs and asserts `k`
    ///      (over virtual + real reserves) is monotonic non-decreasing. The 1%
    ///      fee leaves the pool, so the pool's notional `k` can grow but never
    ///      shrink modulo integer floor on output sizing.
    function testFuzz_k_non_decreasing_across_trades(uint96 quoteIn, uint96 agentIn) public {
        uint256 qIn = bound(uint256(quoteIn), 1e15, 10_000e18);
        uint256 aIn = bound(uint256(agentIn), 1e15, 1_000_000e18);

        uint256 kBefore = curve.k();
        // Fund Alice with the inputs she might need.
        titu.mint(alice, qIn);

        vm.prank(alice);
        curve.buy(0, qIn);
        uint256 kAfterBuy = curve.k();
        assertGe(kAfterBuy, kBefore, "k decreased on buy");

        // Constrain the sell to inventory Alice actually holds.
        uint256 bal = agent.balanceOf(alice);
        if (bal == 0) return;
        aIn = bound(aIn, 1, bal);
        vm.prank(alice);
        agent.approve(address(curve), aIn);

        (uint256 net,) = curve.quoteSell(aIn);
        if (net == 0) return;
        if (net + (net * 100) / 10_000 > curve.realQuoteReserve()) return;

        vm.prank(alice);
        curve.sell(aIn, 0);
        uint256 kAfterSell = curve.k();
        assertGe(kAfterSell, kAfterBuy, "k decreased on sell");
    }

    // ---------------------------------------------------------------------
    // Helpers
    // ---------------------------------------------------------------------

    /// @dev Advance the curve to the point `realQuoteReserve == THRESHOLD`. Leaves
    ///      the `graduated` flag as false; graduator is NOT set here either.
    function _reachThreshold() internal {
        // Gross quote to deposit so that (gross - fee) == THRESHOLD.
        // fee = gross * 100 / 10_000 = gross / 100, so gross = THRESHOLD * 100 / 99.
        uint256 gross = (THRESHOLD * 10_000) / 9900;
        // Bump by 1 wei if integer flooring left us below.
        uint256 fee = (gross * 100) / 10_000;
        if (gross - fee < THRESHOLD) gross += 1;
        titu.mint(alice, gross);
        vm.prank(alice);
        curve.buy(0, gross);
        assertGe(curve.realQuoteReserve(), THRESHOLD);
    }

    /// @dev Drive the curve all the way to a graduated state with graduator set.
    function _graduate() internal {
        vm.prank(owner);
        curve.setGraduator(grad);
        _reachThreshold();
        curve.requestGraduation();
    }
}

/// @dev ERC-20 mock that re-enters the bonding curve on `transferFrom`. Used to
///      probe the {ReentrancyGuard} on `buy`/`sell`/`pullForGraduation`. The
///      reentry callback is configurable so each test targets a single entry
///      point.
contract ReentrantToken is ERC20 {
    BondingCurve public curve;
    bool public attack;
    bytes public payload;

    constructor() ERC20("R", "R") {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }

    function arm(BondingCurve c, bytes memory p) external {
        curve = c;
        payload = p;
        attack = true;
    }

    function transferFrom(address from, address to, uint256 amt) public override returns (bool) {
        if (attack) {
            attack = false; // single-shot — guards against infinite recursion if guard somehow allows it
            (bool ok, bytes memory ret) = address(curve).call(payload);
            // bubble revert reason if any, so the test sees the guard error.
            if (!ok) {
                assembly {
                    revert(add(ret, 0x20), mload(ret))
                }
            }
        }
        return super.transferFrom(from, to, amt);
    }
}

/// @dev Edge-case suite that constructs intentionally skewed curves to exercise
///      the defensive `agentOut == 0`, `agentOut > realAgentReserve`,
///      `quoteOut == 0`, and `grossQuoteOut > realQuoteReserve` reverts that the
///      production parameter set never reaches. Each test builds a fresh curve
///      tuned to push the math into the relevant guard.
contract BondingCurveEdgeTest is Test {
    MockERC20 internal titu;
    MockERC20 internal agent;

    address internal owner = address(0xA0);
    address internal feeRouter = address(0xFEE);
    address internal alice = address(0xA11CE);

    function setUp() public {
        titu = new MockERC20("TITU", "TITU");
        agent = new MockERC20("AGENT", "AGT");
        titu.mint(alice, type(uint128).max);
        agent.mint(alice, type(uint128).max);
    }

    /// @dev Build a curve with the given parameters and approve transfers from alice.
    function _build(uint256 vQuote, uint256 vAgent, uint256 threshold, uint256 initialAgent)
        internal
        returns (BondingCurve c)
    {
        c = new BondingCurve(owner, address(agent), address(titu), feeRouter, vQuote, vAgent, threshold, initialAgent);
        agent.mint(address(c), initialAgent);
        vm.prank(alice);
        titu.approve(address(c), type(uint256).max);
        vm.prank(alice);
        agent.approve(address(c), type(uint256).max);
    }

    /// @dev With a curve where Y is much smaller than X, a 1-wei buy floor-divides
    ///      to `agentOut == 0`. Asserts the defensive `InsufficientOutput` revert
    ///      fires on this corner.
    function test_buy_reverts_when_agent_out_rounds_to_zero() public {
        // X = 1e30 (huge virtual quote), Y = ~1 (after offset). agentOut = (Y * dx') / (X + dx').
        // For dx' = 1 wei → agentOut = (1 * 1) / (1e30 + 1) = 0.
        BondingCurve c = _build(1e30, 2, 1e30, 1);
        vm.prank(alice);
        vm.expectRevert(BondingCurve.InsufficientOutput.selector);
        c.buy(0, 1);
    }

    /// @dev With a low physical inventory but a much larger pricing-Y (offset),
    ///      a buy can compute an `agentOut` that exceeds `realAgentReserve`.
    ///      The curve must reject this rather than under-deliver.
    function test_buy_reverts_when_agent_out_exceeds_real_reserve() public {
        // virtualAgent = 1e30, initialAgent = 1. Offset = ~1e30.
        // For a non-trivial dx', agentOut ≈ Y_eff (~1e30) ≫ realAgent (1).
        BondingCurve c = _build(1e18, 1e30, 1e30, 1);
        vm.prank(alice);
        vm.expectRevert(BondingCurve.InsufficientOutput.selector);
        c.buy(0, 1e18);
    }

    /// @dev With a tiny `agentIn` against a large `yPlus`, `grossOut` floors to 0
    ///      and the sell defensive guard fires.
    function test_sell_reverts_when_quote_out_rounds_to_zero() public {
        // Curve where realQuote is non-zero so a sell is even reachable.
        BondingCurve c = _build(1e18, 2e30, 1e30, 1e30);
        // Seed real quote so the sell math has something to draw from.
        vm.prank(alice);
        c.buy(0, 1e18);
        // 1 wei agent in: grossOut = (X * 1) / (Y + 1) where Y ≈ 2e30 → 0.
        vm.prank(alice);
        vm.expectRevert(BondingCurve.InsufficientOutput.selector);
        c.sell(1, 0);
    }

    /// @dev With a virtualQuote that dominates realQuote, the gross-out from a
    ///      large sell can exceed `realQuoteReserve`. The curve must refuse to
    ///      pay out more than it physically holds.
    function test_sell_reverts_when_gross_exceeds_real_quote() public {
        // Large virtualQuote vs small realQuote. Pricing yields large grossOut on a
        // big agentIn but only realQuote is available to pay out.
        BondingCurve c = _build(1_000_000e18, 2_000_000e18, 1e30, 1_000_000e18);
        // Buy a small amount so realQuote is small.
        vm.prank(alice);
        c.buy(0, 1e18);
        // Now sell a huge agentIn — gross > realQuote.
        uint256 huge = 500_000e18;
        vm.prank(alice);
        vm.expectRevert(BondingCurve.InsufficientOutput.selector);
        c.sell(huge, 0);
    }

    // ---------------------------------------------------------------------
    // Reentrancy probes
    // ---------------------------------------------------------------------

    /// @dev A malicious quote token re-enters `buy` from inside its own
    ///      `transferFrom`. The {ReentrancyGuard} on `buy` must refuse the
    ///      nested call.
    function test_buy_reentrancy_guarded() public {
        ReentrantToken evilQuote = new ReentrantToken();
        evilQuote.mint(alice, 1_000_000e18);

        BondingCurve c = new BondingCurve(
            owner,
            address(agent),
            address(evilQuote),
            feeRouter,
            30_000e18,
            1_073_000_000e18,
            42_000e18,
            1_000_000_000e18
        );
        agent.mint(address(c), 1_000_000_000e18);

        vm.prank(alice);
        evilQuote.approve(address(c), type(uint256).max);

        // Arm the token to re-enter buy from transferFrom.
        evilQuote.arm(c, abi.encodeWithSelector(BondingCurve.buy.selector, uint256(0), uint256(1e18)));

        vm.prank(alice);
        // OZ ReentrancyGuard reverts with `ReentrancyGuardReentrantCall()`. We don't
        // pin the selector to avoid coupling to OZ internals across versions.
        vm.expectRevert();
        c.buy(0, 1e18);
    }

    /// @dev Same probe targeting `sell` via the agent token's transferFrom hook.
    function test_sell_reentrancy_guarded() public {
        ReentrantToken evilAgent = new ReentrantToken();
        evilAgent.mint(alice, 1_000_000e18);

        BondingCurve c = new BondingCurve(
            owner,
            address(evilAgent),
            address(titu),
            feeRouter,
            30_000e18,
            1_073_000_000e18,
            42_000e18,
            1_000_000_000e18
        );
        evilAgent.mint(address(c), 1_000_000_000e18);

        vm.prank(alice);
        titu.approve(address(c), type(uint256).max);
        // Need realQuote on the curve so sell quote-out math is non-zero.
        vm.prank(alice);
        c.buy(0, 1000e18);

        vm.prank(alice);
        evilAgent.approve(address(c), type(uint256).max);

        evilAgent.arm(c, abi.encodeWithSelector(BondingCurve.sell.selector, uint256(1e18), uint256(0)));

        vm.prank(alice);
        vm.expectRevert();
        c.sell(1e18, 0);
    }
}
