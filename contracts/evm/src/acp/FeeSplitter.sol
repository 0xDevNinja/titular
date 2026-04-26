// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

/// @title FeeSplitter
/// @notice Distributes protocol fees according to two split schedules:
///
///         Schedule A — "95/5" (agent job payment):
///           - 95% → agent (primary recipient)
///           -  3% → treasury
///           -  2% → buyback-burner
///           Total: 100%
///
///         Schedule B — "90/5/5" (platform or liquidity fee):
///           - 90% → primary recipient
///           -  5% → treasury
///           -  5% → buyback-burner
///           Total: 100%
///
///         Basis points (BPS) are used throughout (10 000 = 100%).
///         Dust (remainder from integer division) stays with primary recipient.
///
///         CALLER_ROLE: addresses allowed to call `split` (e.g. Job contracts).
///         DEFAULT_ADMIN_ROLE: configure treasury + buyback addresses, grant/revoke roles.
contract FeeSplitter is AccessControl, ReentrancyGuard {
    using SafeERC20 for IERC20;

    // ─────────────────────────────────────────────────────────────
    // Constants
    // ─────────────────────────────────────────────────────────────

    uint256 public constant BPS = 10_000;

    /// @notice Schedule A: treasury share = 3%, buyback share = 2%, primary = 95%.
    uint256 public constant SCHEDULE_A_TREASURY_BPS = 300;
    uint256 public constant SCHEDULE_A_BUYBACK_BPS = 200;

    /// @notice Schedule B: treasury share = 5%, buyback share = 5%, primary = 90%.
    uint256 public constant SCHEDULE_B_TREASURY_BPS = 500;
    uint256 public constant SCHEDULE_B_BUYBACK_BPS = 500;

    // ─────────────────────────────────────────────────────────────
    // Roles
    // ─────────────────────────────────────────────────────────────

    bytes32 public constant CALLER_ROLE = keccak256("CALLER_ROLE");

    // ─────────────────────────────────────────────────────────────
    // Split schedule enum
    // ─────────────────────────────────────────────────────────────

    enum Schedule {
        A, // 95 / 3 / 2
        B // 90 / 5 / 5
    }

    // ─────────────────────────────────────────────────────────────
    // State
    // ─────────────────────────────────────────────────────────────

    /// @notice Treasury address receives the treasury share.
    address public treasury;

    /// @notice BuybackBurner address receives the buyback share.
    address public buybackBurner;

    // ─────────────────────────────────────────────────────────────
    // Events
    // ─────────────────────────────────────────────────────────────

    /// @notice Emitted when a fee split is executed.
    event FeeSplit(
        address indexed token,
        address indexed primary,
        Schedule indexed schedule,
        uint256 primaryAmount,
        uint256 treasuryAmount,
        uint256 buybackAmount
    );

    /// @notice Emitted when treasury address is updated.
    event TreasuryUpdated(address indexed oldTreasury, address indexed newTreasury);

    /// @notice Emitted when buyback burner address is updated.
    event BuybackBurnerUpdated(address indexed oldBurner, address indexed newBurner);

    // ─────────────────────────────────────────────────────────────
    // Errors
    // ─────────────────────────────────────────────────────────────

    error ZeroAddress();
    error ZeroAmount();

    // ─────────────────────────────────────────────────────────────
    // Constructor
    // ─────────────────────────────────────────────────────────────

    /// @param admin         Address granted DEFAULT_ADMIN_ROLE.
    /// @param _treasury     Initial treasury address.
    /// @param _buybackBurner Initial BuybackBurner address.
    constructor(address admin, address _treasury, address _buybackBurner) {
        if (admin == address(0)) revert ZeroAddress();
        if (_treasury == address(0)) revert ZeroAddress();
        if (_buybackBurner == address(0)) revert ZeroAddress();

        treasury = _treasury;
        buybackBurner = _buybackBurner;

        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(CALLER_ROLE, admin);
    }

    // ─────────────────────────────────────────────────────────────
    // Split
    // ─────────────────────────────────────────────────────────────

    /// @notice Split `amount` of `token` according to `schedule`.
    /// @dev    Caller must hold CALLER_ROLE.
    ///         Caller must have pre-approved this contract for `amount` of `token`.
    ///         The primary recipient receives any dust from integer division.
    /// @param  token     ERC-20 token to split.
    /// @param  amount    Total amount to split (must be > 0).
    /// @param  primary   Primary recipient (agent or LPs) — receives the largest share.
    /// @param  schedule  Split schedule (A or B).
    function split(address token, uint256 amount, address primary, Schedule schedule)
        external
        nonReentrant
        onlyRole(CALLER_ROLE)
    {
        if (token == address(0)) revert ZeroAddress();
        if (amount == 0) revert ZeroAmount();
        if (primary == address(0)) revert ZeroAddress();

        (uint256 treasuryBps, uint256 buybackBps) = _bps(schedule);

        uint256 treasuryAmt = (amount * treasuryBps) / BPS;
        uint256 buybackAmt = (amount * buybackBps) / BPS;
        // Dust goes to primary (primary gets floor; treasury/buyback get floor)
        uint256 primaryAmt = amount - treasuryAmt - buybackAmt;

        emit FeeSplit(token, primary, schedule, primaryAmt, treasuryAmt, buybackAmt);

        // Pull total from caller first
        IERC20(token).safeTransferFrom(msg.sender, address(this), amount);

        // Then distribute — CEI order for each sub-transfer
        IERC20(token).safeTransfer(primary, primaryAmt);
        IERC20(token).safeTransfer(treasury, treasuryAmt);
        IERC20(token).safeTransfer(buybackBurner, buybackAmt);
    }

    // ─────────────────────────────────────────────────────────────
    // Admin
    // ─────────────────────────────────────────────────────────────

    /// @notice Update the treasury address.
    /// @param  newTreasury New treasury address (must be non-zero).
    function setTreasury(address newTreasury) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (newTreasury == address(0)) revert ZeroAddress();
        address old = treasury;
        treasury = newTreasury;
        emit TreasuryUpdated(old, newTreasury);
    }

    /// @notice Update the buyback burner address.
    /// @param  newBurner New buyback burner address (must be non-zero).
    function setBuybackBurner(address newBurner) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (newBurner == address(0)) revert ZeroAddress();
        address old = buybackBurner;
        buybackBurner = newBurner;
        emit BuybackBurnerUpdated(old, newBurner);
    }

    // ─────────────────────────────────────────────────────────────
    // Views
    // ─────────────────────────────────────────────────────────────

    /// @notice Preview the split amounts for a given amount and schedule without transferring.
    /// @param  amount   Total amount to split.
    /// @param  schedule Split schedule.
    /// @return primaryAmt  Amount going to the primary recipient.
    /// @return treasuryAmt Amount going to treasury.
    /// @return buybackAmt  Amount going to buyback burner.
    function previewSplit(uint256 amount, Schedule schedule)
        external
        pure
        returns (uint256 primaryAmt, uint256 treasuryAmt, uint256 buybackAmt)
    {
        (uint256 treasuryBps, uint256 buybackBps) = _bps(schedule);
        treasuryAmt = (amount * treasuryBps) / BPS;
        buybackAmt = (amount * buybackBps) / BPS;
        primaryAmt = amount - treasuryAmt - buybackAmt;
    }

    // ─────────────────────────────────────────────────────────────
    // Internal
    // ─────────────────────────────────────────────────────────────

    function _bps(Schedule schedule) internal pure returns (uint256 treasuryBps, uint256 buybackBps) {
        if (schedule == Schedule.A) {
            treasuryBps = SCHEDULE_A_TREASURY_BPS;
            buybackBps = SCHEDULE_A_BUYBACK_BPS;
        } else {
            treasuryBps = SCHEDULE_B_TREASURY_BPS;
            buybackBps = SCHEDULE_B_BUYBACK_BPS;
        }
    }
}
