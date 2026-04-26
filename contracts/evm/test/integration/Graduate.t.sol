// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test, Vm} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import {LaunchpadFactory} from "../../src/launchpad/LaunchpadFactory.sol";
import {AgentToken} from "../../src/launchpad/AgentToken.sol";
import {BondingCurve} from "../../src/launchpad/BondingCurve.sol";
import {LPLock} from "../../src/launchpad/LPLock.sol";
import {Graduator} from "../../src/launchpad/Graduator.sol";
import {FeeRouter} from "../../src/launchpad/FeeRouter.sol";

import {IUniswapV2Router02} from "../../src/launchpad/interfaces/IUniswapV2Router02.sol";

import {MockUniswapV2Factory} from "../mocks/MockUniswapV2Factory.sol";
import {MockUniswapV2Router} from "../mocks/MockUniswapV2Router.sol";
import {MockMintableERC20} from "../mocks/MockERC20.sol";

/// @title Graduate.t.sol
/// @notice End-to-end integration coverage for the graduation flow:
///           - the curve halts at exactly the 42k TITU threshold;
///           - {Graduator.graduate} pulls both reserves atomically and
///             dispatches them into Uniswap V2 `addLiquidity` with the
///             LP-lock as the `to` recipient;
///           - the LP lock binds the V2 pair as its LP token, the creator
///             as the sole beneficiary, and a 10-year unlock anchored to
///             launch — withdraw before unlock reverts; after unlock and
///             with a real LP balance, the full balance is released to
///             the creator;
///           - the FeeRouter splits both pre-graduation curve fees AND
///             the agent transfer-tax stream 70/30 to creator/treasury.
/// @dev    No fork. Uses the deterministic in-tree mocks for V2 factory +
///         router. The mock router records the `addLiquidity` parameters
///         instead of minting real LP — for the post-mint LP-lock semantics
///         test we use a separately deployed LPLock against a real ERC-20
///         standin so we can exercise withdraw / withdraw-before-unlock
///         paths against a real balance.
contract GraduateIntegrationTest is Test {
    // ---------------------------------------------------------------------
    // Shared fleet
    // ---------------------------------------------------------------------

    LaunchpadFactory internal factory;
    AgentToken internal agentTokenImpl;
    BondingCurve internal bondingCurveImpl;
    Graduator internal graduator;
    FeeRouter internal feeRouter;

    MockMintableERC20 internal titu;
    MockUniswapV2Factory internal uniFactory;
    MockUniswapV2Router internal router;

    // Per-test artifacts pulled from the launch.
    AgentToken internal agent;
    BondingCurve internal curve;
    LPLock internal lpLock;
    address internal pair;

    // ---------------------------------------------------------------------
    // Actors
    // ---------------------------------------------------------------------

    address internal owner = address(0xA0);
    address internal treasury = address(0x7);
    address internal creator = address(0xC0E);
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);

    // ---------------------------------------------------------------------
    // Spec constants — mirror LaunchpadFactory immutables.
    // ---------------------------------------------------------------------

    uint256 internal constant VIRTUAL_QUOTE = 30_000e18;
    uint256 internal constant VIRTUAL_AGENT = 1_073_000_000e18;
    uint256 internal constant THRESHOLD = 42_000e18;
    uint256 internal constant INITIAL_AGENT = 1_000_000_000e18;
    uint256 internal constant TEN_YEARS = 365 days * 10;

    // ---------------------------------------------------------------------
    // Setup
    // ---------------------------------------------------------------------

    function setUp() public {
        // Pin block time to a stable anchor so the unlock window math is
        // easy to reason about.
        vm.warp(1_000_000);

        titu = new MockMintableERC20("TITU", "TITU");
        uniFactory = new MockUniswapV2Factory();
        router = new MockUniswapV2Router(address(uniFactory));

        agentTokenImpl = new AgentToken();
        bondingCurveImpl = new BondingCurve(
            address(this),
            address(uint160(uint256(keccak256("agent")))),
            address(uint160(uint256(keccak256("titu")))),
            address(uint160(uint256(keccak256("router")))),
            VIRTUAL_QUOTE,
            VIRTUAL_AGENT,
            THRESHOLD,
            INITIAL_AGENT
        );

        feeRouter = new FeeRouter(treasury, owner);
        graduator = new Graduator(IUniswapV2Router02(address(router)), owner);

        factory = new LaunchpadFactory(
            address(titu),
            treasury,
            address(agentTokenImpl),
            address(bondingCurveImpl),
            address(feeRouter),
            address(graduator),
            address(uniFactory),
            address(router),
            owner
        );

        vm.prank(owner);
        graduator.transferOwnership(address(factory));

        // Configure 70/30 fee route on the FeeRouter for the agent ERC-20
        // BEFORE we know the agent address — bind it inside `_launch` so
        // the route is keyed on the real address.

        // Launch a single agent, no modules.
        _launch();

        // Fund traders with TITU (enough to reach the 42k threshold + slack)
        // and approve the curve.
        titu.mint(alice, 50_000e18);
        titu.mint(bob, 50_000e18);
        vm.prank(alice);
        titu.approve(address(curve), type(uint256).max);
        vm.prank(bob);
        titu.approve(address(curve), type(uint256).max);
    }

    // ---------------------------------------------------------------------
    // Helpers
    // ---------------------------------------------------------------------

    function _launch() internal {
        LaunchpadFactory.LaunchParams memory p;
        p.name = "Alpha";
        p.symbol = "AGA";
        p.imageURI = "ipfs://image";
        p.soulURI = "ipfs://soul";
        p.enabledModules = new bytes32[](0);
        p.moduleData = new bytes[](0);

        vm.prank(creator);
        (address tok, address crv,) = factory.launchAgent(p);
        agent = AgentToken(tok);
        curve = BondingCurve(crv);

        LaunchpadFactory.Agent memory rec = factory.getAgent(1);
        lpLock = LPLock(rec.lpLock);
        pair = rec.pair;

        // Now wire the FeeRouter route for this agent.
        uint16 defaultCreatorBps = feeRouter.DEFAULT_CREATOR_BPS();
        vm.prank(owner);
        feeRouter.setRoute(address(agent), creator, defaultCreatorBps);
    }

    /// @dev Drive the curve to exactly the graduation threshold by sizing
    ///      a single buy so `realQuoteReserve` >= THRESHOLD by the smallest
    ///      possible margin. Fee = 1%, so gross input must be `THRESHOLD *
    ///      10_000 / 9_900`. Bump by one wei if rounding bites.
    function _driveToThreshold() internal {
        uint256 gross = (THRESHOLD * 10_000) / 9900;
        // Floor of 1% may leave us 1 wei short of the threshold — bump.
        uint256 feePreview = (gross * 100) / 10_000;
        if (gross - feePreview < THRESHOLD) gross += 1;

        // Top up alice if 50k TITU isn't enough.
        if (titu.balanceOf(alice) < gross) {
            titu.mint(alice, gross);
        }

        vm.prank(alice);
        curve.buy(0, gross);
    }

    // ---------------------------------------------------------------------
    // Threshold gating
    // ---------------------------------------------------------------------

    function test_request_graduation_reverts_below_threshold() public {
        vm.prank(alice);
        curve.buy(0, 1000e18);
        assertLt(curve.realQuoteReserve(), THRESHOLD);

        vm.expectRevert(BondingCurve.BelowThreshold.selector);
        curve.requestGraduation();
    }

    function test_buy_cannot_cross_threshold_in_a_single_tx() public {
        // A gross buy that would push real reserve strictly past the
        // threshold must revert. Sizing: anything above THRESHOLD * 10000 /
        // 9900 by more than rounding leaves the fee-net amount above
        // THRESHOLD and trips ExceedsGraduationThreshold.
        uint256 gross = (THRESHOLD * 10_000) / 9900 + 1000e18; // well past
        titu.mint(alice, gross);
        vm.prank(alice);
        vm.expectRevert(BondingCurve.ExceedsGraduationThreshold.selector);
        curve.buy(0, gross);
    }

    function test_request_graduation_succeeds_at_threshold() public {
        _driveToThreshold();
        assertGe(curve.realQuoteReserve(), THRESHOLD);

        // Anyone can now flip the graduation flag; we use bob as a stand-in
        // keeper to prove the keeper-style permissionless trigger.
        vm.prank(bob);
        curve.requestGraduation();
        assertTrue(curve.graduated());
    }

    function test_buy_reverts_after_graduation_request() public {
        _driveToThreshold();
        curve.requestGraduation();

        vm.prank(alice);
        vm.expectRevert(BondingCurve.AlreadyGraduated.selector);
        curve.buy(0, 1e18);
    }

    function test_sell_reverts_after_graduation_request() public {
        // Drive straight to threshold so the very last buy lands the curve
        // EXACTLY at `realQuoteReserve == THRESHOLD`. A subsequent buy would
        // overshoot, so we use the threshold buy itself to give alice agent
        // inventory she could in principle sell back.
        _driveToThreshold();
        uint256 sellAmount = agent.balanceOf(alice);
        assertGt(sellAmount, 0, "alice must hold agent to sell");

        curve.requestGraduation();

        vm.prank(alice);
        agent.approve(address(curve), sellAmount);
        vm.prank(alice);
        vm.expectRevert(BondingCurve.AlreadyGraduated.selector);
        curve.sell(sellAmount, 0);
    }

    // ---------------------------------------------------------------------
    // Atomic graduation: pull + addLiquidity + LP-lock target
    //
    // Tax accounting (post-issue #265 fix): the AgentToken's {taxExempt}
    // allowlist — populated at {initialize} with `feeRouter`, `bondingCurve`,
    // and `graduator` — keeps the graduation hand-off lossless. The two
    // protocol-internal hops (curve -> graduator, graduator -> Uniswap
    // router) BOTH skip the 1% transfer tax because the source side
    // (`bondingCurve` then `graduator`) is allowlisted. The Uniswap router
    // therefore receives the curve's full pre-pull `realAgentReserve` and
    // `addLiquidity` does not underflow. End-user transfers (curve -> buyer
    // is exempt because the curve is allowlisted; buyer -> buyer is taxed
    // because neither side is allowlisted) remain on the original economics.
    // ---------------------------------------------------------------------

    function test_graduate_pulls_both_reserves_to_router() public {
        _driveToThreshold();
        curve.requestGraduation();

        uint256 quoteReserve = curve.realQuoteReserve();
        uint256 agentReserve = curve.realAgentReserve();
        assertGt(quoteReserve, 0);
        assertGt(agentReserve, 0);

        vm.prank(bob);
        graduator.graduate(address(curve));

        // Reserves are zeroed on the curve.
        assertEq(curve.realQuoteReserve(), 0);
        assertEq(curve.realAgentReserve(), 0);
        assertEq(titu.balanceOf(address(curve)), 0);
        assertEq(agent.balanceOf(address(curve)), 0);

        // Quote is tax-free (TITU is a vanilla ERC-20 in this test fixture);
        // the router holds the full quote reserve.
        assertEq(titu.balanceOf(address(router)), quoteReserve, "router holds the quote leg");

        // Agent leg: with the {taxExempt} allowlist populated at AgentToken
        // initialization, BOTH legs of the graduation hand-off (curve ->
        // graduator and graduator -> router) skip the 1% tax. The router
        // receives the full pre-pull `realAgentReserve` and zero is skimmed
        // to the FeeRouter from the graduation path.
        assertEq(agent.balanceOf(address(router)), agentReserve, "router holds the full agent leg");

        // Graduator must hold no surplus after the cross-tax round-trip.
        assertEq(titu.balanceOf(address(graduator)), 0);
        assertEq(agent.balanceOf(address(graduator)), 0, "graduator drained on graduation");

        // Idempotency flag flips.
        assertTrue(graduator.graduated(address(curve)));
    }

    /// @dev End-to-end "real-world" graduation that does NOT use any cheat
    ///      code (no `vm.deal`, no balance pokes). Production graduation on
    ///      Base Sepolia / mainnet must succeed under exactly these mechanics
    ///      — drive to the threshold, request graduation, then graduate with
    ///      no top-up. This is the primary regression for issue #265.
    function test_graduate_succeeds_without_deal_cheat() public {
        _driveToThreshold();
        curve.requestGraduation();

        uint256 quoteReserve = curve.realQuoteReserve();
        uint256 agentReserve = curve.realAgentReserve();
        assertGt(quoteReserve, 0);
        assertGt(agentReserve, 0);

        // Pre-condition: the graduator has zero agent balance — i.e. there is
        // no off-curve top-up. In the unfixed code path the subsequent
        // `graduate` would revert on the router's `transferFrom` underflow.
        assertEq(agent.balanceOf(address(graduator)), 0, "graduator pre-graduation balance must be zero");

        vm.prank(bob);
        graduator.graduate(address(curve));

        // Both legs land on the router with no skim.
        assertEq(titu.balanceOf(address(router)), quoteReserve, "router holds the full quote leg");
        assertEq(agent.balanceOf(address(router)), agentReserve, "router holds the full agent leg");
        assertTrue(graduator.graduated(address(curve)), "graduation flag flipped");
    }

    function test_graduate_calls_addLiquidity_with_lpLock_as_recipient() public {
        _driveToThreshold();
        curve.requestGraduation();

        uint256 quoteReserve = curve.realQuoteReserve();
        uint256 agentReserve = curve.realAgentReserve();

        vm.prank(bob);
        graduator.graduate(address(curve));

        (
            address tokenA,
            address tokenB,
            uint256 aDesired,
            uint256 bDesired,
            uint256 aMin,
            uint256 bMin,
            address to,
            uint256 deadline,
            bool recorded
        ) = router.lastCall();

        assertTrue(recorded, "router must record the addLiquidity call");
        // Graduator orders args as (agentToken, quoteToken, agentAmount, quoteAmount).
        assertEq(tokenA, address(agent));
        assertEq(tokenB, address(titu));
        assertEq(aDesired, agentReserve);
        assertEq(bDesired, quoteReserve);
        // 0.5% slippage floor on both sides.
        assertEq(aMin, (agentReserve * 9950) / 10_000);
        assertEq(bMin, (quoteReserve * 9950) / 10_000);
        // LP must be minted to the LPLock the factory registered at launch.
        assertEq(to, address(lpLock), "LP recipient is the LP lock");
        // Same-block deadline.
        assertEq(deadline, block.timestamp);
        assertEq(router.addLiquidityCalls(), 1);
    }

    function test_graduate_creates_pair_only_once() public {
        // The factory pre-creates the pair at launch, so the graduator
        // should NOT need to create it again.
        uint256 createCallsBefore = uniFactory.createPairCalls();
        assertEq(createCallsBefore, 1, "pair pre-created by factory");

        _driveToThreshold();
        curve.requestGraduation();
        vm.prank(bob);
        graduator.graduate(address(curve));

        // No extra createPair on the graduation path.
        assertEq(uniFactory.createPairCalls(), 1, "graduator must reuse the pre-created pair");
        assertEq(uniFactory.lastPair(), pair, "pair address unchanged");
    }

    function test_graduate_emits_Graduated_event() public {
        _driveToThreshold();
        curve.requestGraduation();

        uint256 quoteReserve = curve.realQuoteReserve();
        uint256 agentReserve = curve.realAgentReserve();

        vm.recordLogs();
        vm.prank(bob);
        graduator.graduate(address(curve));

        Vm.Log[] memory logs = vm.getRecordedLogs();
        bytes32 sig = keccak256("Graduated(address,address,uint256,uint256,uint256)");
        bool found;
        for (uint256 i = 0; i < logs.length; ++i) {
            if (logs[i].emitter != address(graduator)) continue;
            if (logs[i].topics.length == 0 || logs[i].topics[0] != sig) continue;
            // Indexed: curve, pair.
            assertEq(address(uint160(uint256(logs[i].topics[1]))), address(curve));
            assertEq(address(uint160(uint256(logs[i].topics[2]))), pair);
            (uint256 liquidity, uint256 quoteAmt, uint256 agentAmt) =
                abi.decode(logs[i].data, (uint256, uint256, uint256));
            assertEq(liquidity, router.liquidityOut());
            assertEq(quoteAmt, quoteReserve);
            assertEq(agentAmt, agentReserve);
            found = true;
            break;
        }
        assertTrue(found, "Graduated event missing");
    }

    function test_graduate_is_idempotent() public {
        _driveToThreshold();
        curve.requestGraduation();
        vm.prank(bob);
        graduator.graduate(address(curve));

        // Second call must trip the graduator's `AlreadyGraduated` guard.
        vm.expectRevert(Graduator.AlreadyGraduated.selector);
        vm.prank(bob);
        graduator.graduate(address(curve));
    }

    function test_graduate_reverts_before_curve_request() public {
        _driveToThreshold();
        // We intentionally skip `requestGraduation` here.
        vm.expectRevert(Graduator.NotYetGraduatedOnCurve.selector);
        vm.prank(bob);
        graduator.graduate(address(curve));
    }

    // ---------------------------------------------------------------------
    // 10-year LP lock — wiring + permissioning
    // ---------------------------------------------------------------------

    function test_lp_lock_wired_to_pair_with_creator_beneficiary() public view {
        assertEq(address(lpLock.lpToken()), pair, "LP token is the V2 pair");
        assertEq(lpLock.beneficiary(), creator, "beneficiary is the agent creator");
        assertFalse(lpLock.withdrawn());
        assertEq(lpLock.lockedAmount(), 0, "no deposits accounted yet");
    }

    function test_lp_lock_unlock_time_is_ten_years_from_launch() public view {
        // The lock was constructed during `factory.launchAgent` in setUp,
        // which executed at `block.timestamp` of the setUp anchor.
        uint256 unlockAt = lpLock.unlockTime();
        // The lock anchors at the per-launch block; we know that block was
        // setUp's `vm.warp(1_000_000)` since no warps fired between then
        // and the launch.
        assertEq(unlockAt, 1_000_000 + TEN_YEARS, "10-year lock from launch block");
    }

    function test_lp_lock_withdraw_reverts_before_unlock() public {
        // No LP balance is needed to exercise the time gate; the gate fires
        // before the balance load.
        bytes memory expected =
            abi.encodeWithSelector(LPLock.StillLocked.selector, lpLock.unlockTime() - block.timestamp);
        vm.prank(creator);
        vm.expectRevert(expected);
        lpLock.withdraw();
    }

    function test_lp_lock_withdraw_reverts_for_non_beneficiary() public {
        vm.warp(lpLock.unlockTime() + 1);
        vm.prank(alice);
        vm.expectRevert(LPLock.NotBeneficiary.selector);
        lpLock.withdraw();
    }

    /// @dev The factory-registered LP token is a deterministic stub (no
    ///      contract code at the address). To exercise withdraw with a real
    ///      LP balance, deploy a separate LPLock against a real ERC-20 and
    ///      walk the full flow. This isolates the LOCK SEMANTICS from the
    ///      mock-router's no-mint behaviour without sacrificing coverage.
    function test_lp_lock_withdraw_releases_balance_after_unlock() public {
        // Spin up a real ERC-20 standin for the LP token and a fresh lock
        // bound to it. The 10-year clock starts at construction.
        MockMintableERC20 lpToken = new MockMintableERC20("AGA-TITU LP", "AGA-LP");
        LPLock isolated = new LPLock(IERC20(address(lpToken)), creator);
        uint256 lpAmount = router.liquidityOut();
        // Mint LP directly to the lock to model the post-graduation state
        // a real V2 router would have produced.
        lpToken.mint(address(isolated), lpAmount);

        // Pre-unlock: revert.
        bytes memory expected =
            abi.encodeWithSelector(LPLock.StillLocked.selector, isolated.unlockTime() - block.timestamp);
        vm.prank(creator);
        vm.expectRevert(expected);
        isolated.withdraw();

        // Warp past the lock and withdraw.
        vm.warp(isolated.unlockTime() + 1);

        uint256 creatorBefore = lpToken.balanceOf(creator);
        vm.prank(creator);
        isolated.withdraw();

        assertEq(lpToken.balanceOf(creator), creatorBefore + lpAmount, "creator received full LP balance");
        assertEq(lpToken.balanceOf(address(isolated)), 0, "lock drained");
        assertTrue(isolated.withdrawn(), "withdrawn flag set");
    }

    function test_lp_lock_withdraw_one_shot() public {
        MockMintableERC20 lpToken = new MockMintableERC20("AGA-TITU LP", "AGA-LP");
        LPLock isolated = new LPLock(IERC20(address(lpToken)), creator);
        lpToken.mint(address(isolated), 1000e18);

        vm.warp(isolated.unlockTime() + 1);
        vm.prank(creator);
        isolated.withdraw();

        // A second call must trip AlreadyWithdrawn even if more LP arrived.
        lpToken.mint(address(isolated), 500e18);
        vm.prank(creator);
        vm.expectRevert(LPLock.AlreadyWithdrawn.selector);
        isolated.withdraw();
    }

    // ---------------------------------------------------------------------
    // Fee router 70/30 — pre-graduation + agent-tax stream
    // ---------------------------------------------------------------------

    function test_curve_swap_fees_distribute_70_30() public {
        // A buy deposits a 1% TITU swap fee on the FeeRouter.
        uint256 quoteIn = 5000e18;
        (, uint256 swapFee) = curve.quoteBuy(quoteIn);
        vm.prank(alice);
        curve.buy(0, quoteIn);
        assertEq(titu.balanceOf(address(feeRouter)), swapFee, "TITU swap fee on router");

        // Have a fresh keeper pull `swapFee` worth of TITU and call distribute
        // on the FeeRouter — isolates the SPLIT math from the curve settlement.
        titu.mint(bob, swapFee);
        vm.prank(bob);
        titu.approve(address(feeRouter), swapFee);

        uint256 creatorBefore = titu.balanceOf(creator);
        uint256 treasuryBefore = titu.balanceOf(treasury);
        vm.prank(bob);
        feeRouter.distribute(address(agent), IERC20(address(titu)), swapFee);

        uint256 expectedCreator = (swapFee * 7000) / 10_000;
        uint256 expectedTreasury = swapFee - expectedCreator;
        assertEq(titu.balanceOf(creator) - creatorBefore, expectedCreator, "creator 70%");
        assertEq(titu.balanceOf(treasury) - treasuryBefore, expectedTreasury, "treasury 30% + dust");
    }

    function test_agent_transfer_tax_fee_distributes_70_30() public {
        // Buy to give alice agent inventory. The bondingCurve is on the
        // {taxExempt} allowlist, so curve -> alice delivers the FULL amount
        // and NOTHING is skimmed to the FeeRouter on the buy hop. The 1%
        // agent tax is exclusively levied on user-to-user transfers.
        vm.prank(alice);
        curve.buy(0, 2000e18);
        assertEq(agent.balanceOf(address(feeRouter)), 0, "curve buy is exempt: no agent skim");

        // Peer-to-peer transfer that triggers the 1% agent tax.
        uint256 sendAmount = 100e18;
        uint256 expectedTax = (sendAmount * 100) / 10_000;
        vm.prank(alice);
        agent.transfer(bob, sendAmount);
        assertEq(agent.balanceOf(address(feeRouter)), expectedTax, "p2p transfer adds 1% tax to router");

        // Snapshot the FULL feeRouter agent balance (just the p2p tax) and
        // route it through `distribute(...)`. The mock router accepts arbitrary
        // amounts. To feed the FeeRouter's `distribute` (which pulls via
        // transferFrom from the caller), route the balance through bob using
        // the FeeRouter -> bob hop, which is allowlisted (FeeRouter exempt) so
        // no tax fires.
        uint256 stream = agent.balanceOf(address(feeRouter));
        vm.prank(address(feeRouter));
        agent.transfer(bob, stream);

        vm.prank(bob);
        agent.approve(address(feeRouter), stream);

        uint256 creatorBefore = agent.balanceOf(creator);
        uint256 treasuryBefore = agent.balanceOf(treasury);
        vm.prank(bob);
        feeRouter.distribute(address(agent), IERC20(address(agent)), stream);

        // The transferFrom inside distribute pulls `stream` from bob into
        // the FeeRouter (router-leg tax-skip applies). Then the router pushes
        // the creator and treasury legs out. Both legs are FROM the router,
        // so they also skip the tax — the recipient deltas equal the floor
        // of the BPS split exactly.
        uint256 expectedCreator = (stream * 7000) / 10_000;
        uint256 expectedTreasury = stream - expectedCreator;
        assertEq(agent.balanceOf(creator) - creatorBefore, expectedCreator, "creator 70% on agent stream");
        assertEq(agent.balanceOf(treasury) - treasuryBefore, expectedTreasury, "treasury 30% on agent stream");
    }

    function test_fee_router_treasury_is_immutable() public view {
        // Sanity wiring: the treasury we passed at construction is what the
        // router still reports — the constructor pinned it as immutable.
        assertEq(feeRouter.treasury(), treasury);
    }

    // ---------------------------------------------------------------------
    // LPLock — view surface
    // ---------------------------------------------------------------------

    function test_lp_lock_time_remaining_pre_unlock() public view {
        // setUp anchored at block.timestamp = 1_000_000; unlockTime = + 10y.
        // `timeRemaining` should report exactly the lock duration since no
        // time has elapsed since launch.
        assertEq(lpLock.timeRemaining(), TEN_YEARS);
    }

    function test_lp_lock_time_remaining_after_unlock_is_zero() public {
        vm.warp(lpLock.unlockTime() + 1);
        assertEq(lpLock.timeRemaining(), 0);
    }

    // ---------------------------------------------------------------------
    // BondingCurve admin — setFeeRouter rebind
    // ---------------------------------------------------------------------

    function test_curve_setFeeRouter_rebinds_recipient() public {
        // The factory is the curve's owner. Build a fresh FeeRouter and
        // rebind the curve's swap-fee recipient via the factory's prank.
        FeeRouter alt = new FeeRouter(treasury, owner);
        address curveOwner = curve.owner();

        vm.expectEmit(true, true, false, false, address(curve));
        emit BondingCurve.FeeRouterSet(address(feeRouter), address(alt));

        vm.prank(curveOwner);
        curve.setFeeRouter(address(alt));
        assertEq(curve.feeRouter(), address(alt));

        // Subsequent buys route the swap fee to the new router.
        uint256 quoteIn = 1000e18;
        (, uint256 fee) = curve.quoteBuy(quoteIn);
        uint256 oldRouterBefore = titu.balanceOf(address(feeRouter));
        uint256 newRouterBefore = titu.balanceOf(address(alt));

        vm.prank(alice);
        curve.buy(0, quoteIn);

        assertEq(titu.balanceOf(address(feeRouter)) - oldRouterBefore, 0, "old router untouched");
        assertEq(titu.balanceOf(address(alt)) - newRouterBefore, fee, "new router collects fee");
    }

    // ---------------------------------------------------------------------
    // Coverage — direct lock deposit path
    // ---------------------------------------------------------------------

    function test_lp_lock_deposit_increments_locked_amount() public {
        MockMintableERC20 lpToken = new MockMintableERC20("AGA-TITU LP", "AGA-LP");
        LPLock isolated = new LPLock(IERC20(address(lpToken)), creator);

        lpToken.mint(alice, 1000e18);
        vm.prank(alice);
        lpToken.approve(address(isolated), 1000e18);

        vm.prank(alice);
        isolated.deposit(1000e18);

        assertEq(isolated.lockedAmount(), 1000e18);
        assertEq(lpToken.balanceOf(address(isolated)), 1000e18);
    }
}
