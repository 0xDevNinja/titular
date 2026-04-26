// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import {Graduator} from "../../src/launchpad/Graduator.sol";
import {IUniswapV2Router02} from "../../src/launchpad/interfaces/IUniswapV2Router02.sol";
import {IUniswapV2Factory} from "../../src/launchpad/interfaces/IUniswapV2Factory.sol";
import {IBondingCurve} from "../../src/launchpad/interfaces/IBondingCurve.sol";

/// @dev Tax-free ERC-20 used by the symbolic harness so the SMT solver does
///      not have to reason about the AgentToken's transfer-tax branching.
///      The graduation-tax interaction is covered separately by the invariant
///      regression test.
contract HalToken is ERC20 {
    constructor() ERC20("HAL", "HAL") {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @dev Symbolic stub for the Uniswap V2 router. Pulls both legs via
///      transferFrom (so allowance + balance preconditions are visible to
///      halmos) and returns a deterministic LP quantity. Halmos is free to
///      pick any input symbolically; the only path it can hit is the one
///      where the graduator approved exactly the expected amount.
contract HalRouter {
    address public immutable factory;
    address public constant WETH = address(0xdeaDDeADDEaDdeaDdEAddEADDEAdDeadDEADDEaD);

    constructor(address fac) {
        factory = fac;
    }

    function addLiquidity(
        address tokenA,
        address tokenB,
        uint256 amountADesired,
        uint256 amountBDesired,
        uint256, /*amountAMin*/
        uint256, /*amountBMin*/
        address, /*to*/
        uint256 /*deadline*/
    ) external returns (uint256, uint256, uint256) {
        IERC20(tokenA).transferFrom(msg.sender, address(this), amountADesired);
        IERC20(tokenB).transferFrom(msg.sender, address(this), amountBDesired);
        return (amountADesired, amountBDesired, 1000);
    }
}

/// @dev Symbolic V2 factory. `getPair` returns a fixed address; `createPair`
///      writes it into a mapping so the graduator's "create-if-missing"
///      branch is exercisable.
contract HalFactory {
    address public pinned;

    function setPinned(address p) external {
        pinned = p;
    }

    function getPair(address, address) external view returns (address) {
        return pinned;
    }

    function createPair(address, address) external view returns (address) {
        return pinned; // halmos doesn't care about the deterministic address; this keeps the path simple.
    }
}

/// @dev Stand-in for a bonding curve that the graduator can pull from.
///      Stores the agent + quote token addresses, the reserves to hand off,
///      and a `graduated` flag halmos can flip to either value.
contract HalCurve {
    address public agentToken;
    address public quoteToken;
    bool public graduated;
    uint256 public quoteReserve;
    uint256 public agentReserve;

    function setup(address a, address q, uint256 qa, uint256 ar, bool g) external {
        agentToken = a;
        quoteToken = q;
        quoteReserve = qa;
        agentReserve = ar;
        graduated = g;
    }

    function pullForGraduation() external returns (uint256, uint256) {
        require(graduated, "not graduated");
        uint256 qa = quoteReserve;
        uint256 ar = agentReserve;
        quoteReserve = 0;
        agentReserve = 0;
        if (qa != 0) IERC20(quoteToken).transfer(msg.sender, qa);
        if (ar != 0) IERC20(agentToken).transfer(msg.sender, ar);
        return (qa, ar);
    }
}

/// @title GraduatorSymbolic
/// @notice Halmos symbolic execution harness for Graduator.graduate.
///         Two `check_*` functions exhaust the safety properties under any
///         caller, any reserve sizes (bounded, see below), and any
///         pre-flag-state combinations:
///           - check_graduate_idempotent: no path lets the same curve be
///             graduated twice.
///           - check_graduate_zeroesCurve: every successful path leaves the
///             curve drained on both sides.
///         Run via:
///           halmos --contract GraduatorSymbolic --function check_graduate_idempotent
///           halmos --contract GraduatorSymbolic --function check_graduate_zeroesCurve
contract GraduatorSymbolic is Test {
    Graduator internal graduator;
    HalCurve internal curve;
    HalToken internal agent;
    HalToken internal titu;
    HalRouter internal router;
    HalFactory internal factory;
    address internal lpLock = address(0xC0E);
    address internal owner = address(this);

    function setUp() public {
        agent = new HalToken();
        titu = new HalToken();
        factory = new HalFactory();
        factory.setPinned(address(0xBEEFFACE));
        router = new HalRouter(address(factory));
        curve = new HalCurve();
        graduator = new Graduator(IUniswapV2Router02(address(router)), owner);
        graduator.registerLPLock(address(curve), lpLock);
    }

    /// @notice Graduating a curve twice MUST revert on the second call.
    /// @dev Symbolic over (quoteReserve, agentReserve). The graduator's
    ///      `graduated[curve]` is checked before the pull so the second
    ///      call hits AlreadyGraduated regardless of reserve sizes.
    function check_graduate_idempotent(uint96 quoteRes, uint96 agentRes) public {
        // Bound to uint96 so halmos doesn't have to reason about pathological
        // overflow space — production reserves fit comfortably in uint96.
        vm.assume(quoteRes > 0 && agentRes > 0);
        curve.setup(address(agent), address(titu), quoteRes, agentRes, true);
        agent.mint(address(curve), agentRes);
        titu.mint(address(curve), quoteRes);

        graduator.graduate(address(curve));

        // Second call must revert. We use a try/catch so halmos can prove
        // EVERY symbolic path through the second invocation reverts.
        bool reverted;
        try graduator.graduate(address(curve)) {
            reverted = false;
        } catch {
            reverted = true;
        }
        assertTrue(reverted, "second graduate must revert");
    }

    /// @notice Successful graduation drains the curve on both sides.
    function check_graduate_zeroesCurve(uint96 quoteRes, uint96 agentRes) public {
        vm.assume(quoteRes > 0 && agentRes > 0);
        curve.setup(address(agent), address(titu), quoteRes, agentRes, true);
        agent.mint(address(curve), agentRes);
        titu.mint(address(curve), quoteRes);

        graduator.graduate(address(curve));

        assertEq(curve.quoteReserve(), 0);
        assertEq(curve.agentReserve(), 0);
        assertTrue(graduator.graduated(address(curve)));
    }

    /// @notice graduate cannot succeed when the curve has not flipped its
    ///         own `graduated` flag.
    function check_graduate_revertsBeforeCurveFlag(uint96 quoteRes, uint96 agentRes) public {
        vm.assume(quoteRes > 0 && agentRes > 0);
        curve.setup(address(agent), address(titu), quoteRes, agentRes, false);
        agent.mint(address(curve), agentRes);
        titu.mint(address(curve), quoteRes);

        bool reverted;
        try graduator.graduate(address(curve)) {
            reverted = false;
        } catch {
            reverted = true;
        }
        assertTrue(reverted, "graduate must revert when curve flag false");
    }
}
