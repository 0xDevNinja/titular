// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {BondingCurve} from "../../../src/launchpad/BondingCurve.sol";

/// @dev Minimal mintable mock used by the curveK invariant suite. Tax-free so the
///      bonding-curve constant-product math stays exact down to the wei.
contract MockToken is ERC20 {
    constructor(string memory n, string memory s) ERC20(n, s) {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @notice Handler that drives randomised buy / sell sequences against a
///         {BondingCurve} from a small pool of actors. Amounts are bounded so
///         the trade can neither push past the graduation threshold nor drain
///         the agent inventory; both edge paths are covered explicitly by the
///         unit tests. The handler stays silent on no-op rejects (zero-out
///         quotes, post-graduation calls, balance-shortfall) so the invariant
///         file can crank `runs * depth` without losing signal to reverts.
contract CurveKHandler is Test {
    BondingCurve public immutable curve;
    MockToken public immutable titu;
    MockToken public immutable agent;
    address[] public actors;

    /// @notice Cap per-buy quote-in well below `graduationThreshold` so the run
    ///         lives entirely in the pre-graduation regime.
    uint256 public constant MAX_BUY = 500e18;
    /// @notice Cap per-sell agent-in by typical inventory so we never starve a
    ///         trader's balance during the depth.
    uint256 public constant MAX_SELL = 5_000_000e18;

    /// @notice Minimum buy size; below this `_quoteBuy` floors `agentOut` to 0
    ///         and the trade reverts inside the curve.
    uint256 public constant MIN_BUY = 1e15;

    constructor(BondingCurve _curve, MockToken _titu, MockToken _agent) {
        curve = _curve;
        titu = _titu;
        agent = _agent;
        actors.push(address(0xA11CE));
        actors.push(address(0xB0B));
        actors.push(address(0xCAFE));
        actors.push(address(0xDEAD));

        for (uint256 i = 0; i < actors.length; ++i) {
            titu.mint(actors[i], 5000e18);
            vm.prank(actors[i]);
            titu.approve(address(curve), type(uint256).max);
            vm.prank(actors[i]);
            agent.approve(address(curve), type(uint256).max);
        }
    }

    function _pick(uint256 i) internal view returns (address) {
        return actors[i % actors.length];
    }

    /// @dev Random buy. Bounded `quoteIn` and pre-flight checks keep the run in
    ///      the safe pre-graduation regime; the call is silently dropped if it
    ///      would revert in the curve.
    function buy(uint256 actorIdx, uint96 amt) external {
        if (curve.graduated()) return;
        address a = _pick(actorIdx);
        uint256 quoteIn = bound(uint256(amt), MIN_BUY, MAX_BUY);

        if (titu.balanceOf(a) < quoteIn) {
            titu.mint(a, quoteIn);
        }

        (uint256 expectedOut, uint256 fee) = curve.quoteBuy(quoteIn);
        if (expectedOut == 0) return;
        if (expectedOut > curve.realAgentReserve()) return;
        if (curve.realQuoteReserve() + (quoteIn - fee) > curve.graduationThreshold()) return;

        vm.prank(a);
        try curve.buy(0, quoteIn) {} catch {}
    }

    /// @dev Random sell using whatever agent inventory the actor holds.
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
}
