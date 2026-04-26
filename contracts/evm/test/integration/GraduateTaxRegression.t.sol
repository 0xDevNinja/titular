// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
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

/// @title GraduateTaxRegression
/// @notice Inverted regression coverage for the AgentToken transfer-tax
///         interaction with the graduation hand-off. The agent token taxes 1%
///         on every peer-to-peer transfer, BUT the {AgentToken.taxExempt}
///         allowlist — populated at {initialize} with `feeRouter`,
///         `bondingCurve`, and `graduator` — keeps the graduation hand-off
///         lossless. The graduation path triggers two protocol-internal
///         transfers — curve -> graduator and graduator -> Uniswap pair
///         (via router.addLiquidity) — and BOTH skip tax because the source
///         side is allowlisted. The router therefore receives the full
///         pre-pull `realAgentReserve` and `addLiquidity` does not underflow.
///
///         These tests run the FULL launch flow (no `deal()` cheats) and
///         assert two things:
///           (a) graduation succeeds end-to-end against an unmodified flow
///               — production graduation works on Base Sepolia / mainnet;
///           (b) curve -> graduator transfers skip tax, while user-to-user
///               transfers are still taxed in full (the allowlist is narrow).
contract GraduateTaxRegressionTest is Test {
    LaunchpadFactory internal factory;
    AgentToken internal agentTokenImpl;
    BondingCurve internal bondingCurveImpl;
    Graduator internal graduator;
    FeeRouter internal feeRouter;

    MockMintableERC20 internal titu;
    MockUniswapV2Factory internal uniFactory;
    MockUniswapV2Router internal router;

    AgentToken internal agent;
    BondingCurve internal curve;
    LPLock internal lpLock;
    address internal pair;

    address internal owner = address(0xA0);
    address internal treasury = address(0x7);
    address internal creator = address(0xC0E);
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);

    uint256 internal constant VIRTUAL_QUOTE = 30_000e18;
    uint256 internal constant VIRTUAL_AGENT = 1_073_000_000e18;
    uint256 internal constant THRESHOLD = 42_000e18;
    uint256 internal constant INITIAL_AGENT = 1_000_000_000e18;

    function setUp() public {
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

        uint16 defaultCreatorBps = feeRouter.DEFAULT_CREATOR_BPS();
        vm.prank(owner);
        feeRouter.setRoute(address(agent), creator, defaultCreatorBps);

        titu.mint(alice, 60_000e18);
        vm.prank(alice);
        titu.approve(address(curve), type(uint256).max);
    }

    function _driveToThreshold() internal {
        uint256 gross = (THRESHOLD * 10_000) / 9900;
        uint256 feePreview = (gross * 100) / 10_000;
        if (gross - feePreview < THRESHOLD) gross += 1;
        if (titu.balanceOf(alice) < gross) titu.mint(alice, gross);
        vm.prank(alice);
        curve.buy(0, gross);
    }

    /// @notice CRITICAL: Graduation MUST succeed in production. With the
    ///         AgentToken {taxExempt} allowlist in place, both legs of the
    ///         hand-off (curve -> graduator and graduator -> Uniswap router)
    ///         skip the 1% tax. The router receives the full agent reserve
    ///         and `addLiquidity` does not underflow. No `deal()` cheats —
    ///         this is exactly the on-chain mechanic that runs on Base
    ///         Sepolia / mainnet.
    function test_graduate_succeeds_with_tax_exempt_graduator() public {
        _driveToThreshold();
        curve.requestGraduation();

        // Pre-graduation: the graduator holds zero agent tokens. The fix is
        // proven by the absence of any top-up before the graduate() call.
        assertEq(agent.balanceOf(address(graduator)), 0, "graduator pre-state must be zero");

        uint256 quoteReserve = curve.realQuoteReserve();
        uint256 agentReserve = curve.realAgentReserve();
        assertGt(quoteReserve, 0);
        assertGt(agentReserve, 0);

        // Graduation completes without a revert — this is the exact path that
        // was reverting pre-fix on the router's transferFrom underflow.
        graduator.graduate(address(curve));

        // The Uniswap router holds the full reserves; nothing is skimmed
        // to the FeeRouter from the protocol-internal hops.
        assertEq(titu.balanceOf(address(router)), quoteReserve, "router holds full quote leg");
        assertEq(agent.balanceOf(address(router)), agentReserve, "router holds full agent leg");
        assertEq(agent.balanceOf(address(graduator)), 0, "graduator drained on graduation");
        assertTrue(graduator.graduated(address(curve)));
    }

    /// @notice Property check: the AgentToken {taxExempt} allowlist exempts
    ///         `feeRouter`, `bondingCurve`, and `graduator` — but NOT
    ///         end-user accounts. Curve -> graduator transfers skip tax;
    ///         user-to-user transfers (alice -> bob) remain taxed in full.
    function test_agentToken_exempts_curve_and_graduator_only() public view {
        assertTrue(agent.taxExempt(address(curve)), "curve allowlisted");
        assertTrue(agent.taxExempt(address(graduator)), "graduator allowlisted");
        assertTrue(agent.taxExempt(address(feeRouter)), "feeRouter allowlisted");
        assertFalse(agent.taxExempt(creator), "creator NOT allowlisted");
        assertFalse(agent.taxExempt(alice), "alice NOT allowlisted");
        assertFalse(agent.taxExempt(bob), "bob NOT allowlisted");
    }

    /// @notice The curve -> graduator hop, modelled in isolation, must
    ///         deliver the FULL amount with zero skim to the FeeRouter.
    ///         This is the exact transfer path the graduation pull executes.
    function test_curve_to_graduator_transfer_is_tax_free() public {
        // Pull some agent inventory directly from the curve to the graduator,
        // mirroring the {BondingCurve.pullForGraduation} hop. Both endpoints
        // are on the {taxExempt} allowlist so the transfer must be lossless.
        uint256 amount = 1_000_000e18;
        uint256 routerBefore = agent.balanceOf(address(feeRouter));
        uint256 graduatorBefore = agent.balanceOf(address(graduator));

        vm.prank(address(curve));
        agent.transfer(address(graduator), amount);

        assertEq(agent.balanceOf(address(feeRouter)), routerBefore, "no tax skim on curve -> graduator");
        assertEq(agent.balanceOf(address(graduator)) - graduatorBefore, amount, "graduator received full amount");
    }

    /// @notice User-to-user transfers (alice -> bob) MUST still be taxed.
    ///         The narrow allowlist must not leak into end-user economics.
    function test_user_to_user_transfer_is_still_taxed() public {
        // Drive a small buy to give alice agent inventory, then transfer to
        // bob (both non-exempt). The 1% tax must fire as before.
        vm.prank(alice);
        curve.buy(0, 1000e18);
        uint256 aliceBal = agent.balanceOf(alice);
        assertGt(aliceBal, 0);

        uint256 routerBefore = agent.balanceOf(address(feeRouter));
        uint256 sendAmount = aliceBal / 2;
        uint256 expectedTax = (sendAmount * 100) / 10_000;

        vm.prank(alice);
        agent.transfer(bob, sendAmount);

        assertEq(agent.balanceOf(bob), sendAmount - expectedTax, "bob receives net of tax");
        assertEq(agent.balanceOf(address(feeRouter)) - routerBefore, expectedTax, "feeRouter accrues 1%");
    }
}
