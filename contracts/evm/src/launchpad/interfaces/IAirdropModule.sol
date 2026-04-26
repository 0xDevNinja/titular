// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title  IAirdropModule
/// @notice Minimal binding the launchpad factory uses to register a per-agent
///         airdrop manifest on the shared Airdrop module.
/// @dev    Mirrors {AirdropModule.configure}. Keep in lockstep with
///         `src/launchpad/modules/AirdropModule.sol`. The configurer (the
///         factory or a creator script) MUST pre-deposit `allocation` of the
///         agent token into the module before this call lands; the module
///         enforces that invariant on `configure`.
interface IAirdropModule {
    /// @notice Bind a merkle-proven airdrop manifest to `agent`. One-shot.
    /// @param  agent          Agent token (the asset airdropped).
    /// @param  root           Merkle root of `(recipient, amount)` allocations.
    /// @param  snapshotBlock  Block at which the indexer snapped veTITU.
    /// @param  allocation     Total airdrop pool, denominated in agent token.
    /// @param  deadline       Unix timestamp after which dust is sweepable.
    /// @param  admin          Address authorised to call `sweep`.
    function configure(
        address agent,
        bytes32 root,
        uint64 snapshotBlock,
        uint128 allocation,
        uint64 deadline,
        address admin
    ) external;
}
