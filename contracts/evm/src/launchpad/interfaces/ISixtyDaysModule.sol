// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title  ISixtyDaysModule
/// @notice Minimal binding the launchpad factory uses to bootstrap the 60-day
///         creator escrow on the shared SixtyDays module.
/// @dev    Mirrors {SixtyDaysModule.configure}. Keep in lockstep with
///         `src/launchpad/modules/SixtyDaysModule.sol`.
interface ISixtyDaysModule {
    /// @notice Bootstrap an agent on the 60-day module. One-shot per agent.
    /// @param  agent          Agent ERC-20 address (storage key).
    /// @param  escrowToken    Token escrowed (creator-leg fee currency).
    /// @param  bondingCurve   Curve allowed to push accruals + contributions.
    /// @param  creator        Creator allowed to commit or refund.
    /// @param  startTime      Unix seconds the 60-day clock starts at.
    function configure(address agent, address escrowToken, address bondingCurve, address creator, uint64 startTime)
        external;
}
