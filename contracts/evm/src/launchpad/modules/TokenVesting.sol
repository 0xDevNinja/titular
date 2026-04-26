// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Initializable} from "@openzeppelin/contracts/proxy/utils/Initializable.sol";

/// @title  TokenVesting
/// @notice Minimal cliff + linear vesting target for EIP-1167 minimal-proxy
///         clones deployed by {PreBuyModule}. One clone per agent: it holds the
///         creator's pre-buy bag and releases it linearly to the creator after
///         the cliff elapses.
/// @dev    Implementation contract — never used directly. {PreBuyModule}
///         deploys a clone via OZ {Clones.cloneDeterministic} and calls
///         {initialize} exactly once. The implementation is locked by calling
///         {_disableInitializers} in the constructor so it cannot itself be
///         hijacked. State layout is intentionally tight: six storage slots.
///         All token movements use {SafeERC20} to tolerate non-standard ERC-20
///         (no return value, USDC-style) implementations.
contract TokenVesting is Initializable {
    using SafeERC20 for IERC20;

    // ---------------------------------------------------------------------
    // Storage
    // ---------------------------------------------------------------------

    /// @notice ERC-20 being vested. Bound at {initialize}.
    IERC20 public token;

    /// @notice Sole recipient of {release} pulls. Bound at {initialize}.
    address public beneficiary;

    /// @notice Wall-clock start of the vesting schedule.
    uint64 public start;

    /// @notice Seconds from {start} before any tokens vest. Pre-cliff
    ///         {releasable} is always zero.
    uint64 public cliff;

    /// @notice Total vesting window in seconds, measured from {start}. After
    ///         `start + duration` the entire {total} grant is vested.
    uint64 public duration;

    /// @notice Total grant size locked at {initialize}. Equals the agent-token
    ///         balance the module transferred into the clone.
    uint256 public total;

    /// @notice Cumulative tokens already pulled by {release}. Monotone
    ///         non-decreasing; bounded above by {total}.
    uint256 public released;

    // ---------------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------------

    /// @dev Emitted exactly once on {initialize}.
    event Initialized(
        address indexed token, address indexed beneficiary, uint256 total, uint64 start, uint64 cliff, uint64 duration
    );

    /// @dev Emitted on every {release} that transfers a non-zero amount.
    event Released(address indexed beneficiary, uint256 amount);

    // ---------------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------------

    /// @dev Thrown when an address argument is zero.
    error ZeroAddress();
    /// @dev Thrown when {initialize} is called with a zero grant.
    error ZeroAmount();
    /// @dev Thrown when `duration == 0` (would divide-by-zero in vesting math).
    error ZeroDuration();
    /// @dev Thrown when `cliff > duration`.
    error CliffExceedsDuration();
    /// @dev Thrown when {release} would transfer zero (idempotent calls revert
    ///      so callers can distinguish "no-op" from successful release).
    error NothingToRelease();

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    /// @dev Lock the implementation contract — clones initialize themselves.
    constructor() {
        _disableInitializers();
    }

    // ---------------------------------------------------------------------
    // Initializer
    // ---------------------------------------------------------------------

    /// @notice Bind the clone's vesting parameters. Must be called by the
    ///         deploying module immediately after {Clones.cloneDeterministic}
    ///         returns the new clone.
    /// @dev    The caller is expected to have already transferred `total_`
    ///         tokens of `token_` into this contract. We do not pull here so
    ///         the deployment + funding + init flow stays atomic in the parent
    ///         transaction without requiring approvals.
    /// @param  token_       ERC-20 to vest.
    /// @param  beneficiary_ Sole recipient of {release}.
    /// @param  total_       Grant size.
    /// @param  start_       Vesting start (unix seconds).
    /// @param  cliff_       Cliff offset from `start_`, in seconds.
    /// @param  duration_    Total vesting window in seconds.
    function initialize(
        IERC20 token_,
        address beneficiary_,
        uint256 total_,
        uint64 start_,
        uint64 cliff_,
        uint64 duration_
    ) external initializer {
        if (address(token_) == address(0) || beneficiary_ == address(0)) revert ZeroAddress();
        if (total_ == 0) revert ZeroAmount();
        if (duration_ == 0) revert ZeroDuration();
        if (cliff_ > duration_) revert CliffExceedsDuration();

        token = token_;
        beneficiary = beneficiary_;
        total = total_;
        start = start_;
        cliff = cliff_;
        duration = duration_;

        emit Initialized(address(token_), beneficiary_, total_, start_, cliff_, duration_);
    }

    // ---------------------------------------------------------------------
    // Release
    // ---------------------------------------------------------------------

    /// @notice Transfer the currently-releasable amount to {beneficiary}.
    /// @dev    Anyone may call (gas arbitrage). CEI: storage update before
    ///         {SafeERC20.safeTransfer}. Reverts on a zero-amount call so the
    ///         caller cannot silently pay gas for a no-op.
    /// @return amount Tokens transferred to {beneficiary} this call.
    /// @dev    Strict `amount == 0` is the documented no-op sentinel; not a
    ///         real `incorrect-equality` finding. The pre-call read is
    ///         time-based — `block.timestamp` drift of seconds is
    ///         immaterial against multi-month vest schedules.
    // slither-disable-next-line incorrect-equality,timestamp
    function release() external returns (uint256 amount) {
        uint256 vestedNow = _vested(uint64(block.timestamp));
        amount = vestedNow - released;
        if (amount == 0) revert NothingToRelease();

        // --- Effects ---
        released = vestedNow;

        // --- Interactions ---
        token.safeTransfer(beneficiary, amount);
        emit Released(beneficiary, amount);
    }

    // ---------------------------------------------------------------------
    // Views
    // ---------------------------------------------------------------------

    /// @notice Cumulative tokens vested as of `block.timestamp`. Monotone
    ///         non-decreasing in time.
    function vested() external view returns (uint256) {
        return _vested(uint64(block.timestamp));
    }

    /// @notice Tokens currently claimable via {release}.
    function releasable() external view returns (uint256) {
        return _vested(uint64(block.timestamp)) - released;
    }

    // ---------------------------------------------------------------------
    // Internal
    // ---------------------------------------------------------------------

    /// @dev Cliff + linear vesting curve.
    ///      - Pre-cliff (`nowTs < start + cliff`) -> 0.
    ///      - Post-end  (`nowTs >= start + duration`) -> {total}.
    ///      - In-window -> `total * (nowTs - start) / duration` (floor).
    ///
    ///      Time-based gating uses `block.timestamp` by design: vesting is
    ///      expressed in wall-clock seconds (cliffs measured in days/years).
    ///      A miner's ~15s timestamp drift is immaterial against grants that
    ///      vest over months. The slither `timestamp` warning is therefore
    ///      not actionable here.
    // slither-disable-next-line timestamp
    function _vested(uint64 nowTs) internal view returns (uint256) {
        uint64 startMem = start;
        uint64 durationMem = duration;
        if (nowTs < startMem + cliff) return 0;
        uint256 end = uint256(startMem) + uint256(durationMem);
        if (uint256(nowTs) >= end) return total;
        // multiplication-before-division; integer floor leaves dust on the
        // contract until the final tick when `nowTs >= end` returns the full
        // `total`.
        uint256 elapsed = uint256(nowTs) - uint256(startMem);
        return (total * elapsed) / uint256(durationMem);
    }
}
