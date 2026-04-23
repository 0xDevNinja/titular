// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {IVeTITU} from "./IVeTITU.sol";

/// @title FeeDistributor
/// @notice Weekly-epoch fee distribution for veTITU holders. Multi-token.
/// @dev Per-token accounting:
///        - `checkpointToken(token)` credits `(balanceOf(this) - lastTokenBalance)` into
///          the current week bucket. Anyone may call.
///        - `claim(user, tokens)` pays out the user's share for fully-elapsed weeks
///          only (never the in-progress week), using `ve.getPastVotes(user, weekStart)`
///          and `ve.getPastTotalSupply(weekStart)`.
///      CEI + ReentrancyGuard on claim. SafeERC20 throughout.
contract FeeDistributor is ReentrancyGuard {
    using SafeERC20 for IERC20;

    uint256 public constant WEEK = 7 days;

    /// @notice Contract providing historical veTITU balances.
    IVeTITU public immutable VE;
    /// @notice Treasury (informational; not authoritative for transfers here).
    address public immutable TREASURY;
    /// @notice Week bucket at which this distributor was constructed.
    uint256 public immutable START_WEEK;

    /// @notice `tokensPerWeek[token][weekStart]` — amount available for that week.
    mapping(address => mapping(uint256 => uint256)) public tokensPerWeek;
    /// @notice Last known balance of `token` after a checkpoint, used to compute delta.
    mapping(address => uint256) public lastTokenBalance;
    /// @notice Last week at which `token` was checkpointed.
    mapping(address => uint256) public lastTokenCheckpoint;
    /// @notice Per-(token,user) next-week-to-claim cursor.
    mapping(address => mapping(address => uint256)) public userNextClaimWeek;

    event TokenCheckpointed(address indexed token, uint256 indexed week, uint256 amount);
    event Claimed(address indexed user, address indexed token, uint256 amount, uint256 weeksAdvanced);

    error ZeroAddress();
    error NothingToClaim();

    constructor(IVeTITU ve, address treasury) {
        if (address(ve) == address(0) || treasury == address(0)) revert ZeroAddress();
        VE = ve;
        TREASURY = treasury;
        // Intentional floor-to-week; divide-before-multiply is the exact pattern we need.
        // slither-disable-next-line divide-before-multiply
        START_WEEK = (block.timestamp / WEEK) * WEEK;
    }

    /// @notice Distribute the delta of this contract's balance for `token` into the current week bucket.
    /// @dev Callable by anyone. If no new tokens arrived the call is a cheap no-op.
    function checkpointToken(address token) external {
        _checkpointToken(token);
    }

    /// @notice Claim rewards for `user` across `tokens`. Pays only fully-elapsed weeks.
    /// @param user   Beneficiary whose ve balance determines the share.
    /// @param tokens Reward tokens to claim.
    function claim(address user, address[] calldata tokens) external nonReentrant {
        if (user == address(0)) revert ZeroAddress();
        // Intentional floor-to-week; divide-before-multiply is the exact pattern we need.
        // slither-disable-next-line divide-before-multiply
        uint256 currentWeek = (block.timestamp / WEEK) * WEEK;
        uint256 len = tokens.length;
        for (uint256 i; i < len; ++i) {
            address token = tokens[i];
            _checkpointToken(token);
            uint256 owed = _computeClaim(user, token, currentWeek);
            if (owed > 0) {
                lastTokenBalance[token] -= owed;
                IERC20(token).safeTransfer(user, owed);
            }
        }
    }

    /// @notice View the amount `user` would receive for `token` at the current time.
    function claimable(address user, address token) external view returns (uint256) {
        // Intentional floor-to-week; divide-before-multiply is the exact pattern we need.
        // slither-disable-next-line divide-before-multiply
        uint256 currentWeek = (block.timestamp / WEEK) * WEEK;
        uint256 cursor = userNextClaimWeek[user][token];
        if (cursor == 0) cursor = START_WEEK;
        uint256 total;
        while (cursor < currentWeek) {
            uint256 veBal = VE.getPastVotes(user, cursor);
            uint256 veTotal = VE.getPastTotalSupply(cursor);
            if (veBal > 0 && veTotal > 0) {
                total += (tokensPerWeek[token][cursor] * veBal) / veTotal;
            }
            cursor += WEEK;
        }
        return total;
    }

    // ---------------------------------------------------------------------
    // Internal
    // ---------------------------------------------------------------------

    function _checkpointToken(address token) internal {
        if (token == address(0)) revert ZeroAddress();
        uint256 balance = IERC20(token).balanceOf(address(this));
        uint256 last = lastTokenBalance[token];
        if (balance <= last) {
            lastTokenCheckpoint[token] = block.timestamp;
            return;
        }
        uint256 delta = balance - last;
        // Intentional floor-to-week; divide-before-multiply is the exact pattern we need.
        // slither-disable-next-line divide-before-multiply
        uint256 currentWeek = (block.timestamp / WEEK) * WEEK;
        tokensPerWeek[token][currentWeek] += delta;
        lastTokenBalance[token] = balance;
        lastTokenCheckpoint[token] = block.timestamp;
        emit TokenCheckpointed(token, currentWeek, delta);
    }

    function _computeClaim(address user, address token, uint256 currentWeek) internal returns (uint256 total) {
        uint256 cursor = userNextClaimWeek[user][token];
        if (cursor == 0) cursor = START_WEEK;
        uint256 weeksAdvanced;
        while (cursor < currentWeek) {
            uint256 veBal = VE.getPastVotes(user, cursor);
            uint256 veTotal = VE.getPastTotalSupply(cursor);
            if (veBal > 0 && veTotal > 0) {
                total += (tokensPerWeek[token][cursor] * veBal) / veTotal;
            }
            cursor += WEEK;
            ++weeksAdvanced;
        }
        userNextClaimWeek[user][token] = cursor;
        if (total > 0) emit Claimed(user, token, total, weeksAdvanced);
    }
}
