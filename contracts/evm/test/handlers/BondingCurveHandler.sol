// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {BondingCurve} from "../../src/launchpad/BondingCurve.sol";

/// @dev Mintable mock used by the invariant harness. Identical to the unit-test
///      mock; redeclared here so the invariant file can instantiate it without
///      importing unit-test scaffolding into the invariant target.
contract MockERC20Inv is ERC20 {
    constructor(string memory n, string memory s) ERC20(n, s) {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @notice Fuzz handler for BondingCurve. Drives randomised buy/sell sequences
///         from a small set of actors; amounts are bounded so the curve can
///         neither graduate nor drain to zero during the run (we are asserting
///         pre-graduation `k` monotonicity in the invariant file).
contract BondingCurveHandler is Test {
    BondingCurve public immutable curve;
    MockERC20Inv public immutable titu;
    MockERC20Inv public immutable agent;
    address[] public actors;

    uint256 public constant MAX_BUY = 1000e18; // keeps sum well under 42k threshold
    uint256 public constant MAX_SELL = 5_000_000e18;

    constructor(BondingCurve _curve, MockERC20Inv _titu, MockERC20Inv _agent) {
        curve = _curve;
        titu = _titu;
        agent = _agent;
        actors.push(address(0xA11CE));
        actors.push(address(0xB0B));
        actors.push(address(0xCAFE));
        actors.push(address(0xDEAD));

        for (uint256 i = 0; i < actors.length; ++i) {
            titu.mint(actors[i], 10_000e18);
            vm.prank(actors[i]);
            titu.approve(address(curve), type(uint256).max);
            vm.prank(actors[i]);
            agent.approve(address(curve), type(uint256).max);
        }
    }

    function _pick(uint256 i) internal view returns (address) {
        return actors[i % actors.length];
    }

    /// @dev Random buy. Input is bounded and the handler silently returns if the
    ///      trade would exceed the graduation threshold — the invariant does not
    ///      care about that path (a separate unit test already covers it).
    function buy(uint256 actorIdx, uint96 amt) external {
        address a = _pick(actorIdx);
        uint256 quoteIn = bound(uint256(amt), 1e15, MAX_BUY);
        if (titu.balanceOf(a) < quoteIn) {
            titu.mint(a, quoteIn);
        }

        (uint256 expectedOut, uint256 fee) = curve.quoteBuy(quoteIn);
        if (expectedOut == 0) return;
        if (curve.realQuoteReserve() + (quoteIn - fee) > curve.graduationThreshold()) return;
        if (expectedOut > curve.realAgentReserve()) return;
        if (curve.graduated()) return;

        vm.prank(a);
        try curve.buy(0, quoteIn) {} catch {}
    }

    /// @dev Random sell using whatever agent inventory the actor currently holds.
    function sell(uint256 actorIdx, uint96 amt) external {
        address a = _pick(actorIdx);
        uint256 bal = agent.balanceOf(a);
        if (bal == 0) return;
        uint256 agentIn = bound(uint256(amt), 1, bal);
        if (agentIn == 0) return;

        (uint256 net, uint256 fee) = curve.quoteSell(agentIn);
        if (net == 0) return;
        if (net + fee > curve.realQuoteReserve()) return;
        if (curve.graduated()) return;

        vm.prank(a);
        try curve.sell(agentIn, 0) {} catch {}
    }
}
