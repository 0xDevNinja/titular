// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {StdInvariant} from "forge-std/StdInvariant.sol";
import {BondingCurve} from "../../src/launchpad/BondingCurve.sol";
import {BondingCurveHandler, MockERC20Inv} from "../handlers/BondingCurveHandler.sol";

/// @notice Invariant harness for {BondingCurve}. Under randomised buy/sell traffic
///         the constant-product `k = (Vq + Rq) * (Va + Ra)` must be monotonic
///         non-decreasing pre-graduation. This is the classic AMM invariant; the
///         1% swap fee leaves the pool on the input side, which if anything grows
///         `k`, and the output side is floored so the invariant holds exactly.
contract BondingCurveInvariant is StdInvariant, Test {
    BondingCurve internal curve;
    MockERC20Inv internal titu;
    MockERC20Inv internal agent;
    BondingCurveHandler internal handler;

    address internal owner = address(0xA0);
    address internal feeRouter = address(0xFEE);

    uint256 internal constant VIRTUAL_QUOTE = 30_000e18;
    uint256 internal constant VIRTUAL_AGENT = 1_073_000_000e18;
    uint256 internal constant THRESHOLD = 42_000e18;
    uint256 internal constant INITIAL_AGENT = 1_000_000_000e18;

    uint256 internal kSeed;

    function setUp() public {
        titu = new MockERC20Inv("TITU", "TITU");
        agent = new MockERC20Inv("AGENT", "AGT");

        curve = new BondingCurve(
            owner, address(agent), address(titu), feeRouter, VIRTUAL_QUOTE, VIRTUAL_AGENT, THRESHOLD, INITIAL_AGENT
        );
        agent.mint(address(curve), INITIAL_AGENT);

        handler = new BondingCurveHandler(curve, titu, agent);
        targetContract(address(handler));

        bytes4[] memory selectors = new bytes4[](2);
        selectors[0] = BondingCurveHandler.buy.selector;
        selectors[1] = BondingCurveHandler.sell.selector;
        targetSelector(FuzzSelector({addr: address(handler), selectors: selectors}));

        kSeed = curve.k();
    }

    /// @notice `k` never decreases across any sequence of handler operations.
    function invariant_curveKNonDecreasing() public view {
        uint256 kNow = curve.k();
        assertGe(kNow, kSeed, "k regressed below seed");
    }

    /// @notice Real quote never pushes past the graduation threshold. The buy
    ///         path explicitly reverts on overshoot; this is a defence-in-depth
    ///         invariant that will catch any future change to the overshoot
    ///         guard.
    function invariant_realQuoteBoundedByThreshold() public view {
        assertLe(curve.realQuoteReserve(), curve.graduationThreshold());
    }

    /// @notice The curve's physical agent balance equals `realAgentReserve`
    ///         exactly. Because the mocks levy no transfer tax, any drift would
    ///         indicate a bookkeeping error.
    function invariant_agentReserveMatchesBalance() public view {
        assertEq(agent.balanceOf(address(curve)), curve.realAgentReserve());
    }
}
