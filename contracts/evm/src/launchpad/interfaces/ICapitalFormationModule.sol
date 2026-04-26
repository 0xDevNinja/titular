// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title  ICapitalFormationModule
/// @notice Minimal binding the launchpad factory uses to register an agent's
///         FDV-milestone payout schedule on the shared CapitalFormation module.
/// @dev    Mirrors {CapitalFormationModule.configure}. Keep in lockstep with
///         `src/launchpad/modules/CapitalFormationModule.sol`. The configurer
///         MUST pre-deposit the sum of `payoutsUsdc` of `payoutToken` into the
///         module before this call lands.
interface ICapitalFormationModule {
    /// @notice Register an FDV-milestone payout schedule for `agent`.
    /// @param  agent          Storage key (typically the agent ERC-20).
    /// @param  agentToken     Agent ERC-20 used to compute FDV.
    /// @param  payoutToken    Token paid out at each milestone (e.g. USDC).
    /// @param  pair           Uniswap V2 pair used as the FDV oracle.
    /// @param  creator        Recipient of the payouts.
    /// @param  agentAdmin     Per-agent admin authorized post-bootstrap.
    /// @param  fdvThresholds  Sorted FDV thresholds (in USD wei).
    /// @param  payoutsUsdc    Parallel array of payout amounts.
    function configure(
        address agent,
        address agentToken,
        address payoutToken,
        address pair,
        address creator,
        address agentAdmin,
        uint256[] calldata fdvThresholds,
        uint256[] calldata payoutsUsdc
    ) external;
}
