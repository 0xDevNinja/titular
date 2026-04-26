// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title  IAntiSniperModule
/// @notice Minimal binding the launchpad factory uses to register a per-agent
///         decay curve on the shared AntiSniper module at launch.
/// @dev    Mirrors the public {AntiSniperModule.configure} signature so the
///         factory can dispatch into the module without importing the full
///         OZ-laden implementation. Keep this in lockstep with
///         `src/launchpad/modules/AntiSniperModule.sol` — the module's
///         `configure` is the one-shot, owner-gated bootstrap path.
interface IAntiSniperModule {
    /// @notice Bind a one-shot decay curve to `agent`.
    /// @param  agent      Agent token address (key).
    /// @param  startTime  Unix timestamp at which decay begins.
    /// @param  duration   Decay window in seconds; must be non-zero.
    /// @param  startBps   Tax rate at `startTime` (<= 10_000).
    /// @param  endBps     Tax rate at/after `startTime + duration`; <= startBps.
    function configure(address agent, uint64 startTime, uint32 duration, uint16 startBps, uint16 endBps) external;
}
