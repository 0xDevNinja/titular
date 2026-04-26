// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {FeeRouter} from "../../../src/launchpad/FeeRouter.sol";

/// @dev Mintable fee-token mock used by the FeeSplit invariant suite.
contract MockFeeToken is ERC20 {
    constructor() ERC20("FEE", "FEE") {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @notice Handler that streams randomised distributions through a
///         pre-configured 70/30 {FeeRouter} route. Tracks the cumulative
///         dispatched amount (`totalDistributed`) so the invariant file can
///         cross-check creator + treasury balances.
contract FeeSplitHandler is Test {
    FeeRouter public immutable router;
    MockFeeToken public immutable token;
    address public immutable agent;
    address public immutable creator;
    address public immutable treasury;
    address[] public actors;

    /// @notice Cumulative wei pushed through {FeeRouter.distribute} since
    ///         construction. Used by the invariant file to confirm the sum
    ///         of recipient balances equals total in.
    uint256 public totalDistributed;

    /// @notice Cumulative wei the handler computed should land at the
    ///         creator side, using the same `floor(amount * 7000 / 10000)`
    ///         recipe as the router.
    uint256 public expectedCreatorTotal;

    /// @notice Cumulative wei the handler computed should land at the
    ///         treasury, including rounding dust.
    uint256 public expectedTreasuryTotal;

    /// @notice Per-call cap so the run stays in 256-bit arithmetic with room
    ///         to spare (cumulative ≤ 4 actors * 1B * MAX_DISTRIBUTE per call).
    uint256 public constant MAX_DISTRIBUTE = 1000e18;
    uint16 public constant CREATOR_BPS = 7000;
    uint16 public constant BPS_DENOMINATOR = 10_000;

    constructor(FeeRouter _router, MockFeeToken _token, address _agent, address _creator, address _treasury) {
        router = _router;
        token = _token;
        agent = _agent;
        creator = _creator;
        treasury = _treasury;
        actors.push(address(0xA11CE));
        actors.push(address(0xB0B));
        actors.push(address(0xCAFE));
        actors.push(address(0xDEAD));

        for (uint256 i = 0; i < actors.length; ++i) {
            token.mint(actors[i], 1_000_000e18);
            vm.prank(actors[i]);
            token.approve(address(router), type(uint256).max);
        }
    }

    function _pick(uint256 i) internal view returns (address) {
        return actors[i % actors.length];
    }

    /// @dev Push a bounded random amount through the router. Mirrors the
    ///      router's split math locally so the invariant file can verify
    ///      both sides land where we expect.
    function distribute(uint256 actorIdx, uint96 amt) external {
        address a = _pick(actorIdx);
        uint256 amount = bound(uint256(amt), 1, MAX_DISTRIBUTE);
        if (token.balanceOf(a) < amount) {
            token.mint(a, amount);
        }

        uint256 creatorAmount = (amount * CREATOR_BPS) / BPS_DENOMINATOR;
        uint256 treasuryAmount = amount - creatorAmount;

        totalDistributed += amount;
        expectedCreatorTotal += creatorAmount;
        expectedTreasuryTotal += treasuryAmount;

        vm.prank(a);
        router.distribute(agent, IERC20(address(token)), amount);
    }
}
