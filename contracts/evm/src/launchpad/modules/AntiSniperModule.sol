// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

/// @title  AntiSniperModule
/// @notice Per-agent dynamic buy-tax that decays linearly from {Config.startBps}
///         down to {Config.endBps} over {Config.duration} seconds, beginning at
///         {Config.startTime}. Used by the launchpad's buy path to burn the
///         first-block / first-minute sniper advantage: the earlier a buyer
///         lands after launch, the larger the tax they pay.
/// @dev    Single shared instance across the agent fleet. The launchpad factory
///         owns this module and calls {configure} exactly once per agent as the
///         final step of launch. Per-agent config is one-shot: it can never be
///         mutated after {configure} lands, so a malicious upgrade path cannot
///         retroactively spike the tax on a live curve.
///
///         This contract is a pure view oracle — it holds no tokens, makes no
///         external calls, and emits exactly one event ({Configured}). The
///         consumer (bonding-curve buy path) reads {currentBps} or
///         {computeTax} and deducts the tax itself. Because the module never
///         touches funds, there is no reentrancy surface to guard.
contract AntiSniperModule is Ownable {
    // ---------------------------------------------------------------------
    // Types
    // ---------------------------------------------------------------------

    /// @notice Per-agent decay curve. All fields are immutable once
    ///         {configured} flips to true.
    /// @dev    Packs into a single 256-bit slot:
    ///           startTime  64 bits
    ///           duration   32 bits
    ///           startBps   16 bits
    ///           endBps     16 bits
    ///           configured  8 bits (bool)
    ///         = 136 bits, one SSTORE per {configure} call.
    struct Config {
        uint64 startTime;
        uint32 duration;
        uint16 startBps;
        uint16 endBps;
        bool configured;
    }

    // ---------------------------------------------------------------------
    // Constants
    // ---------------------------------------------------------------------

    /// @notice Basis-point denominator. 10_000 = 100%.
    uint16 public constant MAX_BPS = 10_000;

    // ---------------------------------------------------------------------
    // Storage
    // ---------------------------------------------------------------------

    /// @notice Per-agent decay configuration. Empty struct until the factory
    ///         calls {configure} for that agent.
    mapping(address agent => Config) public configs;

    // ---------------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------------

    /// @dev Emitted exactly once per agent when the factory locks in the decay
    ///      curve at launch. All subsequent {configure} calls for the same
    ///      agent revert with {AlreadyConfigured}.
    event Configured(address indexed agent, uint64 startTime, uint32 duration, uint16 startBps, uint16 endBps);

    // ---------------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------------

    /// @dev Thrown by {configure} if the agent already has a decay curve bound.
    error AlreadyConfigured();
    /// @dev Thrown when `agent` or `owner_` is the zero address.
    error ZeroAddress();
    /// @dev Thrown when a basis-point argument is out of range. `bps` is echoed
    ///      for indexers chasing misconfigured launches.
    error InvalidBps(uint16 bps);
    /// @dev Thrown when `duration == 0` (would divide-by-zero in the linear
    ///      interpolation branch of {currentBps}).
    error InvalidDuration();

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    /// @param owner_ Initial owner. Expected to be the launchpad factory (or a
    ///               Safe multisig acting as it) so that only the factory can
    ///               bind a decay curve to a freshly deployed agent.
    constructor(address owner_) Ownable(_validatedOwner(owner_)) {}

    /// @dev Wrap OZ `Ownable`'s zero-owner check so callers see our canonical
    ///      {ZeroAddress} error instead of OZ's `OwnableInvalidOwner`.
    function _validatedOwner(address owner_) private pure returns (address) {
        if (owner_ == address(0)) revert ZeroAddress();
        return owner_;
    }

    // ---------------------------------------------------------------------
    // Admin
    // ---------------------------------------------------------------------

    /// @notice Bind a decay curve to `agent`. One-shot per agent.
    /// @dev    onlyOwner so only the factory can register. The config is
    ///         immutable once written — there is no setter for individual
    ///         fields and no `reconfigure` path. This is intentional: a live
    ///         curve's sniper-tax schedule is part of its social contract with
    ///         buyers and must not be mutable by the operator.
    /// @param  agent      Agent token address (key).
    /// @param  startTime  Unix timestamp at which decay begins. May be in the
    ///                    past; pre-start reads return {Config.startBps}.
    /// @param  duration   Decay window in seconds. Must be non-zero.
    /// @param  startBps   Tax rate at `startTime`. <= {MAX_BPS}.
    /// @param  endBps     Tax rate at/after `startTime + duration`.
    ///                    Must be <= `startBps` (decay is monotone non-increasing).
    function configure(address agent, uint64 startTime, uint32 duration, uint16 startBps, uint16 endBps)
        external
        onlyOwner
    {
        if (agent == address(0)) revert ZeroAddress();
        if (configs[agent].configured) revert AlreadyConfigured();
        if (startBps > MAX_BPS) revert InvalidBps(startBps);
        if (endBps > startBps) revert InvalidBps(endBps);
        if (duration == 0) revert InvalidDuration();

        configs[agent] =
            Config({startTime: startTime, duration: duration, startBps: startBps, endBps: endBps, configured: true});

        emit Configured(agent, startTime, duration, startBps, endBps);
    }

    // ---------------------------------------------------------------------
    // Views
    // ---------------------------------------------------------------------

    /// @notice Current buy-tax rate for `agent`, in basis points.
    /// @dev    Piecewise-linear:
    ///           - unconfigured           -> 0
    ///           - before `startTime`     -> `startBps`
    ///           - after end (`startTime+duration`) -> `endBps`
    ///           - in-window              -> `startBps - (startBps-endBps)*elapsed/duration`
    ///         Result is monotone non-increasing in time for any fixed agent.
    ///         Overflow is impossible: all terms are bounded by {MAX_BPS} and
    ///         `uint32` max, well inside 256-bit arithmetic.
    function currentBps(address agent) external view returns (uint16) {
        return _currentBps(agent);
    }

    /// @notice Convenience view: split `amount` into `tax` and `netAmount`
    ///         using the live rate from {currentBps}.
    /// @dev    `tax = amount * currentBps / MAX_BPS` (floor).
    ///         `netAmount = amount - tax`. Rounds in the pool's / protocol's
    ///         favour on the tax side, matching the launchpad's other fee math.
    ///         Overflow: `amount * MAX_BPS` overflows above ~1.16e73, which is
    ///         well past any realistic token amount; callers passing absurd
    ///         values will get a standard checked-math revert.
    function computeTax(address agent, uint256 amount) external view returns (uint256 tax, uint256 netAmount) {
        uint16 bps = _currentBps(agent);
        tax = (amount * bps) / MAX_BPS;
        netAmount = amount - tax;
    }

    // ---------------------------------------------------------------------
    // Internal
    // ---------------------------------------------------------------------

    /// @dev Shared implementation behind {currentBps} and {computeTax}; see
    ///      {currentBps} for the piecewise-linear spec.
    /// @dev Uses `block.timestamp` by design — the decay schedule is expressed
    ///      in wall-clock seconds because the launchpad's client-visible
    ///      "snipe window" is seconds from launch, not blocks. A miner who
    ///      shifts `block.timestamp` by the permitted ~15s skew can at most
    ///      advance the decay by a few percent of the typical 60s window; this
    ///      is strictly pro-buyer (lower tax) and cannot be used to extract
    ///      value from the pool, so the standard slither `timestamp` warning
    ///      is not actionable here.
    // slither-disable-next-line timestamp
    function _currentBps(address agent) internal view returns (uint16) {
        Config memory c = configs[agent];
        if (!c.configured) return 0;
        if (block.timestamp < c.startTime) return c.startBps;

        // `elapsed = block.timestamp - c.startTime` fits in 256 bits trivially.
        // When elapsed >= duration we saturate to `endBps` so the linear branch
        // never underflows.
        uint256 elapsed = block.timestamp - c.startTime;
        if (elapsed >= c.duration) return c.endBps;

        // Linear interpolation. `startBps >= endBps` enforced in {configure},
        // so the subtraction is safe. `drop <= startBps <= MAX_BPS` so the
        // final cast back to uint16 cannot truncate.
        uint256 span = uint256(c.startBps) - uint256(c.endBps);
        uint256 drop = (span * elapsed) / uint256(c.duration);
        return uint16(uint256(c.startBps) - drop);
    }
}
