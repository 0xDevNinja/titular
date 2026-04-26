// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import {BondingCurve} from "../../../src/launchpad/BondingCurve.sol";
import {Graduator} from "../../../src/launchpad/Graduator.sol";
import {LPLock} from "../../../src/launchpad/LPLock.sol";

/// @dev Tax-free mintable ERC-20 used by the graduation atomicity invariant.
///      We deliberately avoid the AgentToken's transfer tax here so the
///      invariant tests the GRADUATION STATE MACHINE in isolation; the tax
///      interaction is covered by the regression test in
///      test/integration/GraduateTaxRegression.t.sol.
contract MockToken is ERC20 {
    constructor(string memory n, string memory s) ERC20(n, s) {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @dev Minimal V2-router stub that records calls and pulls both legs from
///      msg.sender via transferFrom — same behavior as the integration mock,
///      replicated here so the invariant suite is self-contained.
contract InvRouter {
    address public immutable factory;
    address public constant WETH = address(0xdeaDDeADDEaDdeaDdEAddEADDEAdDeadDEADDEaD);

    uint256 public liquidityOut = 1_234_567;
    uint256 public addLiquidityCalls;

    constructor(address fac) {
        factory = fac;
    }

    function addLiquidity(
        address tokenA,
        address tokenB,
        uint256 amountADesired,
        uint256 amountBDesired,
        uint256 /*amountAMin*/,
        uint256 /*amountBMin*/,
        address /*to*/,
        uint256 /*deadline*/
    ) external returns (uint256, uint256, uint256) {
        addLiquidityCalls += 1;
        IERC20(tokenA).transferFrom(msg.sender, address(this), amountADesired);
        IERC20(tokenB).transferFrom(msg.sender, address(this), amountBDesired);
        return (amountADesired, amountBDesired, liquidityOut);
    }
}

/// @dev Minimal V2-factory stub. Returns a deterministic address per pair.
contract InvFactory {
    mapping(bytes32 => address) public pairs;
    uint256 public createPairCalls;
    uint256 public _nonce;

    function getPair(address a, address b) external view returns (address) {
        return pairs[_key(a, b)];
    }

    function createPair(address a, address b) external returns (address pair) {
        bytes32 k = _key(a, b);
        if (pairs[k] != address(0)) return pairs[k];
        createPairCalls += 1;
        _nonce += 1;
        pair = address(uint160(uint256(keccak256(abi.encode(a, b, _nonce)))));
        pairs[k] = pair;
    }

    function _key(address a, address b) internal pure returns (bytes32) {
        return a < b ? keccak256(abi.encode(a, b)) : keccak256(abi.encode(b, a));
    }
}

/// @notice Drives randomised buy/sell traffic and graduation-trigger attempts
///         against a single curve+graduator pair. Records:
///           - the pre-graduation k floor;
///           - the moment graduation fires;
///           - the post-graduation reserves snapshot.
contract GraduationAtomicityHandler is Test {
    BondingCurve public immutable curve;
    Graduator public immutable graduator;
    MockToken public immutable titu;
    MockToken public immutable agent;
    LPLock public immutable lpLock;
    InvRouter public immutable router;

    address[] public actors;

    /// @notice Cap per-buy quote-in. Sized so depth-many buys can both stay
    ///         in the pre-graduation regime AND occasionally hit the
    ///         threshold, so the invariant exercises both halves of the
    ///         state machine over a typical 100x100 invariant run.
    uint256 public constant MAX_BUY = 5_000e18;
    uint256 public constant MIN_BUY = 1e15;
    uint256 public constant MAX_SELL = 5_000_000e18;

    bool public graduationFired;
    bool public graduatorPulled;
    uint256 public kAtGraduation;
    uint256 public quoteAtGraduation;
    uint256 public agentAtGraduation;

    constructor(
        BondingCurve _curve,
        Graduator _graduator,
        MockToken _titu,
        MockToken _agent,
        LPLock _lpLock,
        InvRouter _router
    ) {
        curve = _curve;
        graduator = _graduator;
        titu = _titu;
        agent = _agent;
        lpLock = _lpLock;
        router = _router;

        actors.push(address(0xA11CE));
        actors.push(address(0xB0B));
        actors.push(address(0xCAFE));
        actors.push(address(0xDEAD));

        for (uint256 i = 0; i < actors.length; ++i) {
            titu.mint(actors[i], 60_000e18);
            vm.prank(actors[i]);
            titu.approve(address(curve), type(uint256).max);
            vm.prank(actors[i]);
            agent.approve(address(curve), type(uint256).max);
        }
    }

    function _pick(uint256 i) internal view returns (address) {
        return actors[i % actors.length];
    }

    function buy(uint256 actorIdx, uint96 amt) external {
        if (curve.graduated()) return;
        address a = _pick(actorIdx);
        uint256 quoteIn = bound(uint256(amt), MIN_BUY, MAX_BUY);
        if (titu.balanceOf(a) < quoteIn) titu.mint(a, quoteIn);
        (uint256 expectedOut, uint256 fee) = curve.quoteBuy(quoteIn);
        if (expectedOut == 0) return;
        if (expectedOut > curve.realAgentReserve()) return;
        if (curve.realQuoteReserve() + (quoteIn - fee) > curve.graduationThreshold()) return;
        vm.prank(a);
        try curve.buy(0, quoteIn) {} catch {}
    }

    function sell(uint256 actorIdx, uint96 amt) external {
        if (curve.graduated()) return;
        address a = _pick(actorIdx);
        uint256 bal = agent.balanceOf(a);
        if (bal == 0) return;
        uint256 agentIn = bound(uint256(amt), 1, bal > MAX_SELL ? MAX_SELL : bal);
        if (agentIn == 0) return;
        (uint256 net, uint256 fee) = curve.quoteSell(agentIn);
        if (net == 0) return;
        if (net + fee > curve.realQuoteReserve()) return;
        vm.prank(a);
        try curve.sell(agentIn, 0) {} catch {}
    }

    /// @notice Force a single oversize buy that lands the curve at exactly
    ///         the threshold so we can exercise the graduation hand-off.
    ///         The fuzzer's incremental buys may take many runs to reach
    ///         the threshold; this gives the invariant suite a deterministic
    ///         path into the post-graduation regime within a single depth.
    function pushToThreshold(uint256 actorIdx) external {
        if (curve.graduated()) return;
        address a = _pick(actorIdx);
        uint256 threshold = curve.graduationThreshold();
        uint256 current = curve.realQuoteReserve();
        if (current >= threshold) return;
        uint256 needNet = threshold - current;
        // gross such that gross - 1% fee == needNet; bump by 1 wei for floor.
        uint256 gross = (needNet * 10_000) / 9900;
        uint256 fee = (gross * 100) / 10_000;
        if (gross - fee < needNet) gross += 1;
        if (titu.balanceOf(a) < gross) titu.mint(a, gross);
        vm.prank(a);
        try curve.buy(0, gross) {} catch {}
    }

    function requestGraduation() external {
        if (curve.graduated()) return;
        if (curve.realQuoteReserve() < curve.graduationThreshold()) return;
        try curve.requestGraduation() {
            graduationFired = true;
            kAtGraduation = curve.k();
            quoteAtGraduation = curve.realQuoteReserve();
            agentAtGraduation = curve.realAgentReserve();
        } catch {}
    }

    function graduateOnGraduator(uint256 actorIdx) external {
        if (graduator.graduated(address(curve))) return;
        if (!curve.graduated()) return;
        address keeper = _pick(actorIdx);
        vm.prank(keeper);
        try graduator.graduate(address(curve)) {
            graduatorPulled = true;
        } catch {}
    }
}
