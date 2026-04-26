// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title  ILaunchRadarModule
/// @notice Minimal binding the launchpad factory uses to register an agent on
///         the shared LaunchRadar module at launch.
/// @dev    Mirrors the public {LaunchRadarModule.configure} signature so the
///         factory can dispatch without importing the full OZ-laden module.
///         Keep in lockstep with `src/launchpad/modules/LaunchRadarModule.sol`.
interface ILaunchRadarModule {
    /// @notice Register a brand-new agent on the radar. One-shot per agent.
    /// @param  agent        Agent token address being registered.
    /// @param  startTime    Unix timestamp trading unlocks at.
    /// @param  whitelistOn  Initial state of the whitelist gate.
    /// @param  agentAdmin   Address authorized to mutate this agent's gating
    ///                      after the bootstrap configure call.
    function configure(address agent, uint64 startTime, bool whitelistOn, address agentAdmin) external;
}
