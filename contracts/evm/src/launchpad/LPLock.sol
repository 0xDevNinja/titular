// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/// @title LPLock
/// @notice Immutable 10-year timelock for Uniswap V2 LP tokens minted at graduation.
///         One instance per graduated agent: the {Graduator} mints LP to this contract
///         (or forwards LP via {deposit}) and only the {beneficiary} — the agent
///         creator — can withdraw the full balance once {unlockTime} has elapsed.
/// @dev    Immutable by design: there is no extend, no early-withdraw, and no owner.
///         Once deployed the {unlockTime} is fixed and the contract can neither be
///         upgraded nor paused. Anyone may top up via {deposit} so liquidity streams
///         (e.g. secondary pairs created later) can be pinned to the same schedule.
///         CEI is enforced on {withdraw}: `withdrawn` flips before the LP transfer.
///         `ReentrancyGuard` is layered on for defence-in-depth against an exotic LP
///         token with transfer hooks.
contract LPLock is ReentrancyGuard {
    using SafeERC20 for IERC20;

    // ---------------------------------------------------------------------
    // Constants
    // ---------------------------------------------------------------------

    /// @notice Duration of the lock, measured from construction. Fixed 10 years in
    ///         seconds (365-day years per spec; leap-day drift is negligible at the
    ///         policy level and keeps the arithmetic auditable).
    uint256 public constant LOCK_DURATION = 365 days * 10;

    // ---------------------------------------------------------------------
    // Immutables
    // ---------------------------------------------------------------------

    /// @notice ERC-20 LP token held by this lock. Typically a Uniswap V2 pair.
    IERC20 public immutable lpToken;

    /// @notice Sole address authorized to {withdraw}. Expected to be the agent creator.
    address public immutable beneficiary;

    /// @notice Earliest timestamp at which {withdraw} may succeed.
    ///         Set to `block.timestamp + LOCK_DURATION` at construction.
    uint256 public immutable unlockTime;

    // ---------------------------------------------------------------------
    // Storage
    // ---------------------------------------------------------------------

    /// @notice Cumulative LP tokens deposited into this lock via {deposit}.
    ///         Monotonic: never decreases, including on {withdraw}. This is an
    ///         accounting counter, not a live balance; the live balance is
    ///         `lpToken.balanceOf(address(this))` and is what {withdraw} transfers.
    uint256 public lockedAmount;

    /// @notice Flips to `true` on the single successful {withdraw}. Prevents any
    ///         follow-on withdrawals, including of LP deposited after unlock.
    bool public withdrawn;

    // ---------------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------------

    /// @dev Emitted on every successful {deposit}. `from` is the depositor; `amount`
    ///      is the LP wei transferred in.
    event Deposited(address indexed from, uint256 amount);

    /// @dev Emitted exactly once on the successful {withdraw}. `to` is the
    ///      beneficiary; `amount` is the full LP balance transferred out.
    event Withdrawn(address indexed to, uint256 amount);

    // ---------------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------------

    /// @dev Thrown when {withdraw} is called by anyone other than {beneficiary}.
    error NotBeneficiary();
    /// @dev Thrown when {withdraw} is called before {unlockTime}. Carries the
    ///      remaining seconds for UX.
    error StillLocked(uint256 secondsRemaining);
    /// @dev Thrown when {withdraw} is called after the one-shot withdrawal has
    ///      already fired.
    error AlreadyWithdrawn();
    /// @dev Thrown when the constructor receives a zero address for the LP token
    ///      or the beneficiary.
    error ZeroAddress();
    /// @dev Thrown when {deposit} is called with a zero amount.
    error ZeroAmount();

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    /// @param lpToken_     ERC-20 LP token to lock. Must be non-zero.
    /// @param beneficiary_ Sole withdrawal recipient. Must be non-zero.
    constructor(IERC20 lpToken_, address beneficiary_) {
        if (address(lpToken_) == address(0)) revert ZeroAddress();
        if (beneficiary_ == address(0)) revert ZeroAddress();

        lpToken = lpToken_;
        beneficiary = beneficiary_;
        unlockTime = block.timestamp + LOCK_DURATION;
    }

    // ---------------------------------------------------------------------
    // External
    // ---------------------------------------------------------------------

    /// @notice Top up the lock with `amount` LP tokens. Anyone may call; the
    ///         depositor must have pre-approved this contract for `amount`.
    /// @dev    Uses {SafeERC20.safeTransferFrom} to handle non-standard ERC-20s.
    ///         `lockedAmount` is incremented regardless of any fee-on-transfer
    ///         behaviour: it is an ingress counter, not a live balance. Reverts if
    ///         `amount` is zero so we never emit a no-op {Deposited}.
    /// @param  amount LP wei to pull from `msg.sender`.
    function deposit(uint256 amount) external {
        if (amount == 0) revert ZeroAmount();

        // --- Effects ---
        lockedAmount += amount;

        // --- Interactions ---
        lpToken.safeTransferFrom(msg.sender, address(this), amount);

        emit Deposited(msg.sender, amount);
    }

    /// @notice Release the full LP balance to the {beneficiary}. One-shot.
    /// @dev    Permissioning: only {beneficiary} may call. Gating:
    ///           - reverts {StillLocked} before {unlockTime};
    ///           - reverts {AlreadyWithdrawn} on a second call.
    ///         CEI: `withdrawn` flips before the LP transfer. `ReentrancyGuard` is
    ///         an additional safety net in case the LP token has transfer hooks.
    function withdraw() external nonReentrant {
        if (msg.sender != beneficiary) revert NotBeneficiary();
        if (withdrawn) revert AlreadyWithdrawn();
        if (block.timestamp < unlockTime) {
            revert StillLocked(unlockTime - block.timestamp);
        }

        uint256 balance = lpToken.balanceOf(address(this));

        // --- Effects ---
        withdrawn = true;

        // --- Interactions ---
        if (balance != 0) {
            lpToken.safeTransfer(beneficiary, balance);
        }

        emit Withdrawn(beneficiary, balance);
    }

    // ---------------------------------------------------------------------
    // Views
    // ---------------------------------------------------------------------

    /// @notice Seconds left until {withdraw} is callable. Returns 0 once unlocked.
    function timeRemaining() external view returns (uint256) {
        if (block.timestamp >= unlockTime) return 0;
        return unlockTime - block.timestamp;
    }
}
