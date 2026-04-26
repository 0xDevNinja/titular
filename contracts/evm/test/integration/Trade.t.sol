// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import {LaunchpadFactory} from "../../src/launchpad/LaunchpadFactory.sol";
import {AgentToken} from "../../src/launchpad/AgentToken.sol";
import {BondingCurve} from "../../src/launchpad/BondingCurve.sol";
import {Graduator} from "../../src/launchpad/Graduator.sol";
import {FeeRouter} from "../../src/launchpad/FeeRouter.sol";

import {AntiSniperModule} from "../../src/launchpad/modules/AntiSniperModule.sol";

import {IUniswapV2Router02} from "../../src/launchpad/interfaces/IUniswapV2Router02.sol";

import {MockUniswapV2Factory} from "../mocks/MockUniswapV2Factory.sol";
import {MockUniswapV2Router} from "../mocks/MockUniswapV2Router.sol";
import {MockMintableERC20} from "../mocks/MockERC20.sol";

/// @title Trade.t.sol
/// @notice Integration coverage for the trade path post-launch:
///           - buy + sell on the bonding curve through real factory wiring;
///           - the curve's 1% TITU swap fee lands on the FeeRouter and
///             splits 70/30 to (creator, treasury);
///           - the AgentToken's 1% transfer tax fires on the agent leg of
///             every trade (curve -> trader and trader -> curve are both
///             non-router transfers);
///           - the {AntiSniperModule} decay reads correctly across the
///             pre-start / in-window / post-window time bands. Note: the
///             current curve does NOT consume the AntiSniper rate (it is a
///             pure view oracle), so we assert decay via {currentBps} rather
///             than via curve fee deltas — this matches the on-chain spec.
/// @dev    No fork. Uses the in-tree mock V2 factory + router; the trade
///         path itself never touches Uniswap.
contract TradeIntegrationTest is Test {
    // ---------------------------------------------------------------------
    // Shared fleet
    // ---------------------------------------------------------------------

    LaunchpadFactory internal factory;
    AgentToken internal agentTokenImpl;
    BondingCurve internal bondingCurveImpl;
    Graduator internal graduator;
    FeeRouter internal feeRouter;
    AntiSniperModule internal antiSniper;

    MockMintableERC20 internal titu;
    MockUniswapV2Factory internal uniFactory;
    MockUniswapV2Router internal router;

    // ---------------------------------------------------------------------
    // Per-test wiring
    // ---------------------------------------------------------------------

    AgentToken internal agent;
    BondingCurve internal curve;

    // ---------------------------------------------------------------------
    // Actors
    // ---------------------------------------------------------------------

    address internal owner = address(0xA0);
    address internal treasury = address(0x7);
    address internal creator = address(0xC0E);
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);

    // ---------------------------------------------------------------------
    // Spec constants
    // ---------------------------------------------------------------------

    uint256 internal constant VIRTUAL_QUOTE = 30_000e18;
    uint256 internal constant VIRTUAL_AGENT = 1_073_000_000e18;
    uint256 internal constant THRESHOLD = 42_000e18;
    uint256 internal constant INITIAL_AGENT = 1_000_000_000e18;

    // Anti-sniper window: 99% -> 1% over 60 seconds.
    uint16 internal constant START_BPS = 9900;
    uint16 internal constant END_BPS = 100;
    uint32 internal constant DURATION = 60;
    uint64 internal constant LAUNCH_TS = 1_000_000;

    // ---------------------------------------------------------------------
    // Setup
    // ---------------------------------------------------------------------

    function setUp() public {
        // Pin block time well below the launch window so pre-start branches
        // can be exercised without underflow.
        vm.warp(LAUNCH_TS - 100);

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

        // Bind the AntiSniper module so launches with that id dispatch.
        antiSniper = new AntiSniperModule(address(factory));
        bytes32 asId = factory.MODULE_ID_ANTI_SNIPER();
        vm.prank(owner);
        factory.setModule(asId, address(antiSniper));

        // Drive a single launch via `creator`, with the anti-sniper window
        // armed at LAUNCH_TS. The agent token's `creator` field is the
        // address that called `launchAgent`.
        LaunchpadFactory.LaunchParams memory p;
        p.name = "Alpha";
        p.symbol = "AGA";
        p.imageURI = "ipfs://image";
        p.soulURI = "ipfs://soul";
        bytes32[] memory ids = new bytes32[](1);
        ids[0] = asId;
        bytes[] memory data = new bytes[](1);
        data[0] = abi.encode(LAUNCH_TS, DURATION, START_BPS, END_BPS);
        p.enabledModules = ids;
        p.moduleData = data;

        vm.prank(creator);
        (address tok, address crv,) = factory.launchAgent(p);
        agent = AgentToken(tok);
        curve = BondingCurve(crv);

        // Configure the FeeRouter route for this agent: 70/30 to creator/treasury.
        // The agent ERC-20 itself is keyed because the agent's transfer tax
        // sweeps to the FeeRouter and is then distributed by `agent` key.
        // For the curve's quote-side fee (TITU on the FeeRouter), we key on
        // the same agent address so a single `distribute(agent, ...)` call
        // routes both ERC-20 streams identically.
        uint16 defaultCreatorBps = feeRouter.DEFAULT_CREATOR_BPS();
        vm.prank(owner);
        feeRouter.setRoute(address(agent), creator, defaultCreatorBps);

        // Fund traders with TITU + max approvals to the curve.
        titu.mint(alice, 1_000_000e18);
        titu.mint(bob, 1_000_000e18);
        vm.prank(alice);
        titu.approve(address(curve), type(uint256).max);
        vm.prank(bob);
        titu.approve(address(curve), type(uint256).max);
    }

    // ---------------------------------------------------------------------
    // Helpers
    // ---------------------------------------------------------------------

    function _balances(address a) internal view returns (uint256 tituBal, uint256 agentBal) {
        tituBal = titu.balanceOf(a);
        agentBal = agent.balanceOf(a);
    }

    // ---------------------------------------------------------------------
    // Buy
    // ---------------------------------------------------------------------

    function test_buy_credits_full_agent_to_buyer() public {
        uint256 quoteIn = 1000e18;
        (uint256 expectedAgentOut, uint256 expectedSwapFee) = curve.quoteBuy(quoteIn);

        uint256 aliceTituBefore = titu.balanceOf(alice);

        vm.prank(alice);
        curve.buy(0, quoteIn);

        // Quote leg: alice paid the gross `quoteIn`. The swap fee of 1% goes
        // to the FeeRouter; the rest stays on the curve as reserve.
        assertEq(titu.balanceOf(alice), aliceTituBefore - quoteIn, "buyer paid gross quote");
        assertEq(titu.balanceOf(address(feeRouter)), expectedSwapFee, "swap fee on TITU lands on router");
        assertEq(curve.realQuoteReserve(), quoteIn - expectedSwapFee, "real quote reserve = quoteIn - fee");

        // Agent leg: the BondingCurve is on the AgentToken's {taxExempt} allowlist
        // (set at initialize), so curve -> buyer hops skip the 1% tax. Alice
        // receives the full `agentOut`; nothing is skimmed to the FeeRouter on
        // the buy path. End-user transfers (alice -> someone else) remain taxed.
        assertEq(agent.balanceOf(alice), expectedAgentOut, "buyer receives full agentOut tax-free");
        assertEq(agent.balanceOf(address(feeRouter)), 0, "no agent skim on the curve hop");
        assertEq(curve.realAgentReserve(), INITIAL_AGENT - expectedAgentOut);
    }

    function test_buy_emits_Bought_event() public {
        uint256 quoteIn = 500e18;
        (uint256 expectedAgentOut, uint256 expectedFee) = curve.quoteBuy(quoteIn);

        vm.expectEmit(true, false, false, true, address(curve));
        emit BondingCurve.Bought(alice, quoteIn, expectedAgentOut, expectedFee);

        vm.prank(alice);
        curve.buy(0, quoteIn);
    }

    // ---------------------------------------------------------------------
    // Sell
    // ---------------------------------------------------------------------

    function test_sell_returns_titu_minus_swap_fee() public {
        // First buy so alice has agent inventory she can sell back.
        vm.prank(alice);
        curve.buy(0, 1000e18);

        uint256 sellAmount = agent.balanceOf(alice) / 2;
        vm.prank(alice);
        agent.approve(address(curve), sellAmount);

        // Selling pushes agent INTO the curve. AgentToken transfer-tax skims
        // 1% off the in-flight transfer. The curve receives the NET amount,
        // which is what the BondingCurve credits to its real reserve.
        // To assert exact reserve deltas we have to walk the books in the
        // same order the contract does them: tax fires on transferFrom, the
        // curve sees `sellAmount - tax` arrive (and `realAgentReserve += sellAmount`
        // executes BEFORE the transfer per CEI).
        // For an integration assertion we focus on the buyer-visible result:
        // alice's net TITU receipt and the FeeRouter's receipt of the swap fee.

        uint256 aliceTituBefore = titu.balanceOf(alice);
        uint256 routerTituBefore = titu.balanceOf(address(feeRouter));
        (uint256 expectedNetQuote, uint256 expectedSwapFee) = curve.quoteSell(sellAmount);

        vm.prank(alice);
        curve.sell(sellAmount, expectedNetQuote);

        assertEq(titu.balanceOf(alice), aliceTituBefore + expectedNetQuote, "seller TITU credit");
        assertEq(
            titu.balanceOf(address(feeRouter)) - routerTituBefore, expectedSwapFee, "swap fee on TITU lands on router"
        );
    }

    function test_sell_emits_Sold_event() public {
        vm.prank(alice);
        curve.buy(0, 500e18);
        uint256 sellAmount = agent.balanceOf(alice) / 4;
        vm.prank(alice);
        agent.approve(address(curve), sellAmount);

        (uint256 expectedNet, uint256 expectedFee) = curve.quoteSell(sellAmount);

        vm.expectEmit(true, false, false, true, address(curve));
        emit BondingCurve.Sold(alice, sellAmount, expectedNet, expectedFee);

        vm.prank(alice);
        curve.sell(sellAmount, expectedNet);
    }

    // ---------------------------------------------------------------------
    // Fee splits via the FeeRouter
    // ---------------------------------------------------------------------

    function test_fee_router_distributes_swap_fees_70_30() public {
        // Drive a buy to deposit a swap fee on the FeeRouter (in TITU).
        uint256 quoteIn = 10_000e18;
        (, uint256 swapFee) = curve.quoteBuy(quoteIn);
        vm.prank(alice);
        curve.buy(0, quoteIn);

        // Sanity: fee is on the FeeRouter as an ERC-20 balance.
        assertEq(titu.balanceOf(address(feeRouter)), swapFee);

        // Anyone can call `distribute`, but they must hold the funds and
        // approve. Simplest path: the FeeRouter pulls from itself by
        // pre-funding alice with the TITU balance and routing through her.
        // Here we use a clean pattern instead: a third-party caller pulls
        // the same `swapFee` worth of fresh TITU + the router fans it out.
        // This isolates the SPLIT math from the curve's settlement.
        titu.mint(bob, swapFee);
        vm.prank(bob);
        titu.approve(address(feeRouter), swapFee);

        uint256 creatorBefore = titu.balanceOf(creator);
        uint256 treasuryBefore = titu.balanceOf(treasury);

        vm.prank(bob);
        feeRouter.distribute(address(agent), IERC20(address(titu)), swapFee);

        // 70/30 split with floor on the creator side: treasury pockets dust.
        uint256 creatorShare = (swapFee * 7000) / 10_000;
        uint256 treasuryShare = swapFee - creatorShare;
        assertEq(titu.balanceOf(creator) - creatorBefore, creatorShare, "creator gets 70%");
        assertEq(titu.balanceOf(treasury) - treasuryBefore, treasuryShare, "treasury gets 30% + dust");
        // Sums must equal the input (no funds stuck on the router).
        assertEq(creatorShare + treasuryShare, swapFee);
    }

    function test_fee_router_get_split_view_matches_distribute() public view {
        // Walk the same input through `getSplit` and `_split` (via the
        // 70/30 route configured in setUp) and assert the preview matches
        // the actual on-chain split math.
        uint256 amount = 1_234_567e18;
        (uint256 creatorPart, uint256 treasuryPart) = feeRouter.getSplit(address(agent), amount);
        uint256 expectedCreator = (amount * 7000) / 10_000;
        assertEq(creatorPart, expectedCreator);
        assertEq(treasuryPart, amount - expectedCreator);
    }

    function test_fee_router_native_distribute_splits_70_30() public {
        // Re-deploy a treasury that accepts ETH so the native leg succeeds.
        // The default `treasury = address(0x7)` is an EOA at 0x...0007 and
        // accepts native via the empty-code default. The current `creator`
        // (0xC0E) is also an EOA. So both legs of the native distribute
        // should succeed without changes to setUp.
        vm.deal(alice, 1 ether);
        uint256 creatorBefore = creator.balance;
        uint256 treasuryBefore = treasury.balance;

        vm.prank(alice);
        feeRouter.distributeNative{value: 1 ether}(address(agent));

        assertEq(creator.balance - creatorBefore, 0.7 ether, "creator gets 70% of native");
        assertEq(treasury.balance - treasuryBefore, 0.3 ether, "treasury gets 30% of native");
    }

    function test_fee_router_unconfigured_route_routes_all_to_treasury() public {
        // Clear the route configured in setUp so the agent has no entry.
        vm.prank(owner);
        feeRouter.clearRoute(address(agent));

        titu.mint(bob, 1000e18);
        vm.prank(bob);
        titu.approve(address(feeRouter), 1000e18);

        uint256 treasuryBefore = titu.balanceOf(treasury);
        vm.prank(bob);
        feeRouter.distribute(address(agent), IERC20(address(titu)), 1000e18);

        assertEq(titu.balanceOf(treasury) - treasuryBefore, 1000e18, "all to treasury when route cleared");
        assertEq(titu.balanceOf(creator), 0);
    }

    function test_agent_transfer_tax_lands_on_fee_router() public {
        // Pre-fund alice with agent tokens via a buy.
        vm.prank(alice);
        curve.buy(0, 2000e18);
        uint256 aliceAgent = agent.balanceOf(alice);
        uint256 routerAgentBefore = agent.balanceOf(address(feeRouter));

        // Peer-to-peer transfer alice -> bob. Both legs are non-router so
        // the 1% tax must fire.
        uint256 sendAmount = aliceAgent / 4;
        uint256 expectedTax = (sendAmount * 100) / 10_000;

        vm.prank(alice);
        agent.transfer(bob, sendAmount);

        assertEq(agent.balanceOf(bob), sendAmount - expectedTax, "recipient minus tax");
        assertEq(
            agent.balanceOf(address(feeRouter)) - routerAgentBefore,
            expectedTax,
            "FeeRouter received the agent transfer tax"
        );
    }

    function test_agent_transfer_tax_skipped_on_router_legs() public {
        // Mint agent to alice via curve buy.
        vm.prank(alice);
        curve.buy(0, 2000e18);
        uint256 routerBefore = agent.balanceOf(address(feeRouter));

        // alice -> feeRouter: tax must be skipped (one leg is the router
        // itself, otherwise the router would pay tax on its own sweeps).
        uint256 amount = 100e18;
        vm.prank(alice);
        agent.transfer(address(feeRouter), amount);

        assertEq(agent.balanceOf(address(feeRouter)) - routerBefore, amount, "no tax skim on router-to leg");
    }

    // ---------------------------------------------------------------------
    // Anti-sniper decay window
    // ---------------------------------------------------------------------

    function test_anti_sniper_pre_start_returns_start_bps() public view {
        // We seeded LAUNCH_TS = 1_000_000 and warped to LAUNCH_TS - 100 in
        // setUp. The launch happened at LAUNCH_TS - 100, but the AntiSniper
        // window's `startTime` is LAUNCH_TS itself, so right now
        // `block.timestamp < startTime` and the rate is `startBps`.
        assertEq(uint256(antiSniper.currentBps(address(agent))), uint256(START_BPS));
    }

    function test_anti_sniper_window_open_decays_linearly() public {
        // Halfway through the window: rate should be (startBps + endBps) / 2,
        // give or take floor-rounding.
        vm.warp(LAUNCH_TS + DURATION / 2);
        uint16 mid = antiSniper.currentBps(address(agent));

        // Expected linear: startBps - (startBps - endBps) * elapsed / duration.
        uint256 expected =
            uint256(START_BPS) - ((uint256(START_BPS) - uint256(END_BPS)) * (DURATION / 2)) / uint256(DURATION);
        assertEq(uint256(mid), expected, "mid-window decay");
        assertGt(uint256(mid), uint256(END_BPS));
        assertLt(uint256(mid), uint256(START_BPS));
    }

    function test_anti_sniper_post_window_returns_end_bps() public {
        // Past the decay duration: rate clamps to endBps.
        vm.warp(LAUNCH_TS + DURATION + 10);
        assertEq(uint256(antiSniper.currentBps(address(agent))), uint256(END_BPS));
    }

    function test_anti_sniper_decay_monotone_non_increasing() public {
        // Sample 10 evenly-spaced points across the window and assert the
        // rate is monotone non-increasing in time.
        uint16 prev = type(uint16).max;
        for (uint256 i = 0; i <= DURATION; i += DURATION / 10) {
            vm.warp(LAUNCH_TS + i);
            uint16 cur = antiSniper.currentBps(address(agent));
            assertLe(uint256(cur), uint256(prev), "decay not monotone");
            prev = cur;
        }
    }

    function test_anti_sniper_compute_tax_at_window_open() public {
        // At t = startTime, rate = startBps. computeTax must return the
        // matching tax / netAmount split.
        vm.warp(LAUNCH_TS);
        uint256 amount = 10_000e18;
        (uint256 tax, uint256 net) = antiSniper.computeTax(address(agent), amount);
        assertEq(tax, (amount * uint256(START_BPS)) / 10_000);
        assertEq(net, amount - tax);
    }

    // ---------------------------------------------------------------------
    // Multi-trade fee accumulation
    // ---------------------------------------------------------------------

    function test_multiple_trades_accumulate_fees_on_router() public {
        // Walk a sequence of buys + a sell; the FeeRouter's TITU balance
        // must increase monotonically by the swap fees emitted by each call.
        uint256 expectedRouterTitu;

        for (uint256 i = 0; i < 4; ++i) {
            uint256 amt = 500e18;
            (, uint256 fee) = curve.quoteBuy(amt);
            vm.prank(alice);
            curve.buy(0, amt);
            expectedRouterTitu += fee;
            assertEq(titu.balanceOf(address(feeRouter)), expectedRouterTitu, "buy fee accumulation");
        }

        // Sell back a chunk; sell fee adds another increment.
        uint256 sellAmt = agent.balanceOf(alice) / 8;
        vm.prank(alice);
        agent.approve(address(curve), sellAmt);
        (, uint256 sellFee) = curve.quoteSell(sellAmt);
        vm.prank(alice);
        curve.sell(sellAmt, 0);
        expectedRouterTitu += sellFee;
        assertEq(titu.balanceOf(address(feeRouter)), expectedRouterTitu, "sell fee accumulation");
    }

    // ---------------------------------------------------------------------
    // AgentToken transfer tax — dust + router-touch coverage
    // ---------------------------------------------------------------------

    function test_agent_transfer_tax_dust_path_no_skim() public {
        // Pre-fund alice via a buy; then transfer a value so small the
        // floor-divided 1% tax computes to zero. Per spec, the dust path
        // forwards the full amount with no router skim.
        vm.prank(alice);
        curve.buy(0, 1000e18);
        uint256 bobBefore = agent.balanceOf(bob);
        uint256 routerBefore = agent.balanceOf(address(feeRouter));

        // 99 wei: floor(99 * 100 / 10_000) == 0, so the tax branch is bypassed.
        vm.prank(alice);
        agent.transfer(bob, 99);

        assertEq(agent.balanceOf(bob) - bobBefore, 99, "dust transfer is full-amount");
        assertEq(agent.balanceOf(address(feeRouter)) - routerBefore, 0, "no tax skim on dust");
    }

    function test_agent_transfer_tax_skipped_on_router_from_leg() public {
        // The BondingCurve is now on the {taxExempt} allowlist, so a buy
        // skims nothing into the FeeRouter. To set up the FeeRouter with an
        // agent balance we therefore drive a peer-to-peer transfer that DOES
        // tax (alice -> bob), which routes 1% to the FeeRouter.
        vm.prank(alice);
        curve.buy(0, 1000e18);
        uint256 sendAmount = agent.balanceOf(alice) / 2;
        vm.prank(alice);
        agent.transfer(bob, sendAmount);

        uint256 routerAgent = agent.balanceOf(address(feeRouter));
        require(routerAgent > 0, "router must have agent inventory");

        uint256 bobBefore = agent.balanceOf(bob);

        // FeeRouter -> bob: `feeRouter` is on the {taxExempt} allowlist so
        // this hop skips the 1% tax (router doesn't pay tax on its own sweeps).
        vm.prank(address(feeRouter));
        agent.transfer(bob, routerAgent);

        assertEq(agent.balanceOf(bob) - bobBefore, routerAgent, "router-leg sends full amount tax-free");
        assertEq(agent.balanceOf(address(feeRouter)), 0, "router drained");
    }

    function test_curve_k_invariant_non_decreasing_across_trades() public {
        // Standard constant-product invariant: post-trade k >= pre-trade k.
        uint256 kBefore = curve.k();

        vm.prank(alice);
        curve.buy(0, 500e18);
        uint256 kAfterBuy = curve.k();
        assertGe(kAfterBuy, kBefore, "k must not decrease on buy");

        uint256 sellAmt = agent.balanceOf(alice) / 4;
        vm.prank(alice);
        agent.approve(address(curve), sellAmt);
        vm.prank(alice);
        curve.sell(sellAmt, 0);
        uint256 kAfterSell = curve.k();
        assertGe(kAfterSell, kAfterBuy, "k must not decrease on sell");
    }
}
