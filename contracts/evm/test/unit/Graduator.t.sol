// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test, Vm} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

import {Graduator} from "../../src/launchpad/Graduator.sol";
import {IUniswapV2Router02} from "../../src/launchpad/interfaces/IUniswapV2Router02.sol";

import {MockUniswapV2Router} from "../mocks/MockUniswapV2Router.sol";
import {MockUniswapV2Factory} from "../mocks/MockUniswapV2Factory.sol";
import {MockBondingCurve} from "../mocks/MockBondingCurve.sol";

/// @dev Lightweight mintable ERC20 used as both the TITU quote and the agent
///      token in these unit tests. No transfer tax — we want the graduator's
///      numbers to match the mock-reserve config exactly.
contract MockERC20 is ERC20 {
    constructor(string memory n, string memory s) ERC20(n, s) {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

contract GraduatorTest is Test {
    MockUniswapV2Factory internal uniFactory;
    MockUniswapV2Router internal router;
    MockBondingCurve internal curve;
    MockERC20 internal quote;
    MockERC20 internal agent;

    Graduator internal grad;

    address internal owner = address(0xA0);
    address internal keeper = address(0xBEEF);
    address internal lpLock = address(0x10C);

    uint256 internal constant QUOTE_RESERVE = 42_000e18;
    uint256 internal constant AGENT_RESERVE = 200_000_000e18;

    function setUp() public {
        uniFactory = new MockUniswapV2Factory();
        router = new MockUniswapV2Router(address(uniFactory));
        quote = new MockERC20("TITU", "TITU");
        agent = new MockERC20("AGENT", "AGT");

        grad = new Graduator(IUniswapV2Router02(address(router)), owner);

        curve = new MockBondingCurve();
        curve.set(address(agent), address(quote), QUOTE_RESERVE, AGENT_RESERVE, false);
        curve.setGraduator(address(grad));

        // Seed the mock curve with the tokens it claims to hold so
        // `pullForGraduation` can actually transfer.
        quote.mint(address(curve), QUOTE_RESERVE);
        agent.mint(address(curve), AGENT_RESERVE);
    }

    // ---------------------------------------------------------------------
    // Helpers
    // ---------------------------------------------------------------------

    function _register(address curve_, address lock_) internal {
        vm.prank(owner);
        grad.registerLPLock(curve_, lock_);
    }

    function _flipCurveGraduated() internal {
        curve.setGraduated(true);
    }

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    function test_constructor_sets_immutables() public view {
        assertEq(address(grad.router()), address(router));
        assertEq(address(grad.uniFactory()), address(uniFactory));
        assertEq(grad.owner(), owner);
    }

    function test_constructor_reverts_zero_router() public {
        vm.expectRevert(Graduator.ZeroAddress.selector);
        new Graduator(IUniswapV2Router02(address(0)), owner);
    }

    function test_constructor_reverts_zero_owner() public {
        vm.expectRevert(Graduator.ZeroAddress.selector);
        new Graduator(IUniswapV2Router02(address(router)), address(0));
    }

    function test_constructor_reverts_router_reports_zero_factory() public {
        // Spin up a router whose `factory()` returns address(0) by pointing it
        // at the zero address. The graduator must bail at construction.
        MockUniswapV2Router zf = new MockUniswapV2Router(address(0));
        vm.expectRevert(Graduator.ZeroAddress.selector);
        new Graduator(IUniswapV2Router02(address(zf)), owner);
    }

    // ---------------------------------------------------------------------
    // registerLPLock
    // ---------------------------------------------------------------------

    function test_registerLPLock_only_owner() public {
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, address(this)));
        grad.registerLPLock(address(curve), lpLock);
    }

    function test_registerLPLock_reverts_zero_curve() public {
        vm.prank(owner);
        vm.expectRevert(Graduator.ZeroAddress.selector);
        grad.registerLPLock(address(0), lpLock);
    }

    function test_registerLPLock_reverts_zero_lock() public {
        vm.prank(owner);
        vm.expectRevert(Graduator.ZeroAddress.selector);
        grad.registerLPLock(address(curve), address(0));
    }

    function test_registerLPLock_stores_and_emits() public {
        vm.expectEmit(true, true, true, true, address(grad));
        emit Graduator.LPLockRegistered(address(curve), lpLock);
        vm.prank(owner);
        grad.registerLPLock(address(curve), lpLock);
        assertEq(grad.lpLocks(address(curve)), lpLock);
    }

    // ---------------------------------------------------------------------
    // graduate — preconditions
    // ---------------------------------------------------------------------

    function test_graduate_reverts_before_curve_graduated() public {
        _register(address(curve), lpLock);
        vm.expectRevert(Graduator.NotYetGraduatedOnCurve.selector);
        vm.prank(keeper);
        grad.graduate(address(curve));
    }

    function test_graduate_reverts_if_already_graduated() public {
        _register(address(curve), lpLock);
        _flipCurveGraduated();

        vm.prank(keeper);
        grad.graduate(address(curve));

        // Second call: curve still reports graduated=true (it was flipped
        // before the first call). Our own idempotency guard must trip.
        vm.expectRevert(Graduator.AlreadyGraduated.selector);
        vm.prank(keeper);
        grad.graduate(address(curve));
    }

    function test_graduate_reverts_if_lp_lock_not_registered() public {
        _flipCurveGraduated();
        vm.expectRevert(Graduator.LPLockNotRegistered.selector);
        vm.prank(keeper);
        grad.graduate(address(curve));
    }

    // ---------------------------------------------------------------------
    // graduate — pair creation
    // ---------------------------------------------------------------------

    function test_graduate_creates_pair_if_missing() public {
        _register(address(curve), lpLock);
        _flipCurveGraduated();

        assertEq(uniFactory.createPairCalls(), 0);

        vm.prank(keeper);
        grad.graduate(address(curve));

        assertEq(uniFactory.createPairCalls(), 1);
        assertTrue(uniFactory.lastPair() != address(0));
    }

    function test_graduate_skips_create_if_pair_exists() public {
        _register(address(curve), lpLock);
        _flipCurveGraduated();

        address preExisting = address(0xC0FFEE);
        uniFactory.setPair(address(agent), address(quote), preExisting);

        vm.prank(keeper);
        grad.graduate(address(curve));

        assertEq(uniFactory.createPairCalls(), 0);
        assertEq(uniFactory.lastPair(), preExisting);
    }

    // ---------------------------------------------------------------------
    // graduate — liquidity path
    // ---------------------------------------------------------------------

    function test_graduate_adds_liquidity_to_lpLock() public {
        _register(address(curve), lpLock);
        _flipCurveGraduated();

        vm.prank(keeper);
        grad.graduate(address(curve));

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
        assertTrue(recorded);
        assertEq(tokenA, address(agent));
        assertEq(tokenB, address(quote));
        assertEq(aDesired, AGENT_RESERVE);
        assertEq(bDesired, QUOTE_RESERVE);
        assertEq(aMin, (AGENT_RESERVE * 9950) / 10_000);
        assertEq(bMin, (QUOTE_RESERVE * 9950) / 10_000);
        assertEq(to, lpLock);
        assertEq(deadline, block.timestamp);
        assertEq(router.addLiquidityCalls(), 1);

        // Balances should have landed in the router via the pull-then-approve
        // flow; graduator must hold no leftover tokens.
        assertEq(quote.balanceOf(address(router)), QUOTE_RESERVE);
        assertEq(agent.balanceOf(address(router)), AGENT_RESERVE);
        assertEq(quote.balanceOf(address(grad)), 0);
        assertEq(agent.balanceOf(address(grad)), 0);
    }

    function test_graduate_marks_graduated_idempotent() public {
        _register(address(curve), lpLock);
        _flipCurveGraduated();

        assertFalse(grad.graduated(address(curve)));
        vm.prank(keeper);
        grad.graduate(address(curve));
        assertTrue(grad.graduated(address(curve)));
    }

    function test_graduate_emits_Graduated() public {
        _register(address(curve), lpLock);
        _flipCurveGraduated();

        vm.recordLogs();
        vm.prank(keeper);
        grad.graduate(address(curve));

        // The last log emitted by the graduator is `Graduated(curve, pair, lp,
        // quote, agent)`. We locate it by topic signature to avoid coupling to
        // the ordering of ERC-20 and approval logs that precede it.
        Vm.Log[] memory entries = vm.getRecordedLogs();
        bytes32 sig = keccak256("Graduated(address,address,uint256,uint256,uint256)");

        bool found;
        for (uint256 i = 0; i < entries.length; i++) {
            Vm.Log memory e = entries[i];
            if (e.emitter == address(grad) && e.topics.length > 0 && e.topics[0] == sig) {
                assertEq(address(uint160(uint256(e.topics[1]))), address(curve));
                assertEq(address(uint160(uint256(e.topics[2]))), uniFactory.lastPair());
                (uint256 liquidity, uint256 quoteAmt, uint256 agentAmt) =
                    abi.decode(e.data, (uint256, uint256, uint256));
                assertEq(liquidity, router.liquidityOut());
                assertEq(quoteAmt, QUOTE_RESERVE);
                assertEq(agentAmt, AGENT_RESERVE);
                found = true;
                break;
            }
        }
        assertTrue(found, "Graduated event not emitted");
    }

    function test_graduate_approves_both_tokens() public {
        _register(address(curve), lpLock);
        _flipCurveGraduated();

        vm.prank(keeper);
        grad.graduate(address(curve));

        // After the router pulled the tokens, the remaining allowance must be
        // zero (router consumes exactly `amountDesired`). This proves the
        // approval flow wasn't skipped.
        assertEq(quote.allowance(address(grad), address(router)), 0);
        assertEq(agent.allowance(address(grad), address(router)), 0);

        // And the router booked a successful add-liquidity.
        assertEq(router.addLiquidityCalls(), 1);
    }

    // ---------------------------------------------------------------------
    // graduate — reentrancy (defence-in-depth)
    // ---------------------------------------------------------------------

    function test_graduate_nonReentrant_blocks_self_call() public {
        // Rebind the curve to a reentrant stub that tries to re-enter
        // graduator.graduate during `pullForGraduation`. The combination of
        // `nonReentrant` on `graduate` and the effect-before-interaction flip
        // of `graduated[curve]` must kill the recursion.
        ReentrantCurve attacker = new ReentrantCurve(grad);
        quote.mint(address(attacker), QUOTE_RESERVE);
        agent.mint(address(attacker), AGENT_RESERVE);
        attacker.configure(address(agent), address(quote), QUOTE_RESERVE, AGENT_RESERVE);

        vm.prank(owner);
        grad.registerLPLock(address(attacker), lpLock);

        vm.prank(keeper);
        // The reentrant attempt reverts inside `pullForGraduation`; that
        // revert propagates up and fails the outer call.
        vm.expectRevert();
        grad.graduate(address(attacker));
    }
}

/// @dev Reenters `Graduator.graduate` during its `pullForGraduation` to prove
///      the guard. The first thing the re-entered call does is check the
///      already-graduated flag — which was just flipped — so it should revert
///      with `AlreadyGraduated`. The outer call bubbles up that revert.
contract ReentrantCurve {
    Graduator public immutable grad;
    address public agentToken;
    address public quoteToken;
    uint256 public realQuoteReserve;
    uint256 public realAgentReserve;
    bool public graduated = true;
    address public graduator;

    constructor(Graduator g) {
        grad = g;
    }

    function configure(address a, address q, uint256 qr, uint256 ar) external {
        agentToken = a;
        quoteToken = q;
        realQuoteReserve = qr;
        realAgentReserve = ar;
    }

    function pullForGraduation() external returns (uint256, uint256) {
        // Attempt reentry. This must revert — either because of the
        // `nonReentrant` guard or because `graduated[curve]` was flipped
        // before this call.
        grad.graduate(address(this));
        return (realQuoteReserve, realAgentReserve);
    }

    function requestGraduation() external {}
}
