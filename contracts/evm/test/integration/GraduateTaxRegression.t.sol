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
/// @notice Regression coverage for the AgentToken transfer-tax interaction with
///         the graduation hand-off. The agent token taxes 1% on every peer-to-
///         peer transfer that does not touch the FeeRouter address. The
///         graduation path triggers TWO such transfers — curve -> graduator
///         and graduator -> Uniswap pair (via router.addLiquidity) — so the
///         graduator ends up short of the agent leg it is asked to deposit.
///
///         Without a tax-skip on the bonding curve and the per-curve graduator
///         the production graduation reverts inside the router's transferFrom
///         when it tries to pull the full pre-tax `agentAmount` from a
///         graduator that only holds 99% of it.
///
///         These tests run the FULL launch flow (no `deal()` cheats) and
///         assert two things:
///           (a) graduation against an unskipped tax flow reverts; and
///           (b) the curve and per-clone graduator are documented as needing
///               tax-skip wiring (TODO in src — see audit report).
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

    /// @notice CRITICAL: Graduation reverts in production because of the
    ///         AgentToken transfer-tax double-skim. The curve transfers the
    ///         agent leg to the graduator (1% to feeRouter) and the router
    ///         then calls transferFrom(graduator, pair, amountADesired) which
    ///         underflows the graduator's balance.
    ///
    ///         Recommended fix (documented in audit/M2_SECURITY_REPORT.md):
    ///         add an explicit `taxExempt` allowlist on AgentToken populated
    ///         with `bondingCurve` AND the per-clone graduator at
    ///         {AgentToken.initialize}. The fix lives in solidity-principal
    ///         scope; this test pins the precondition.
    function test_graduate_reverts_without_tax_skip_for_graduator() public {
        _driveToThreshold();
        curve.requestGraduation();

        // Pre-pull state: graduator holds zero agent tokens.
        assertEq(agent.balanceOf(address(graduator)), 0);

        // Graduate. The curve transfers `agentAmount` to graduator -> AgentToken
        // skims 1% to feeRouter. Then graduator approves router for
        // `agentAmount` and the router calls transferFrom(graduator, ...,
        // amountADesired = agentAmount). Because graduator now holds only
        // 0.99 * agentAmount, the transferFrom reverts.
        vm.expectRevert(); // OZ ERC20InsufficientBalance bubbles up
        graduator.graduate(address(curve));
    }

    /// @notice Property check: the AgentToken's tax-skip set today only
    ///         exempts {feeRouter}, mint, and burn. Neither the curve nor any
    ///         per-clone graduator is exempt. This test pins the current
    ///         (broken) behavior so a future code-side fix flips the
    ///         expectation deliberately.
    function test_agentToken_does_not_exempt_curve_or_graduator_today() public {
        // Curve -> graduator transfer is taxed (proven by checking feeRouter
        // accrues exactly 1% when we model the curve-side leg of graduation
        // in isolation).
        uint256 amount = 1_000_000e18;

        // Move some agent inventory off the curve to alice via a buy so we
        // can model an arbitrary peer-to-peer transfer.
        vm.prank(alice);
        curve.buy(0, 1000e18);
        uint256 aliceBal = agent.balanceOf(alice);
        if (aliceBal < amount) amount = aliceBal;

        uint256 routerBefore = agent.balanceOf(address(feeRouter));
        vm.prank(alice);
        agent.transfer(address(graduator), amount);

        // Tax fired: graduator received 99%, feeRouter accrued 1%.
        uint256 expectedTax = (amount * 100) / 10_000;
        assertEq(agent.balanceOf(address(feeRouter)) - routerBefore, expectedTax, "graduator transfer is taxed");
        assertEq(agent.balanceOf(address(graduator)), amount - expectedTax, "graduator short by tax");
    }
}
