// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title  IExistingTokenModule
/// @notice Minimal binding the launchpad factory uses to wrap an existing
///         ERC-20 onto a fresh bonding curve via the shared ExistingToken
///         module.
/// @dev    Mirrors {ExistingTokenModule.configure}. Keep in lockstep with
///         `src/launchpad/modules/ExistingTokenModule.sol`. The configurer
///         MUST pre-deposit `supply` of `externalToken` into the module before
///         this call lands.
interface IExistingTokenModule {
    /// @notice Wrap `externalToken` onto `curve`. One-shot per token.
    /// @param  externalToken  ERC-20 being wrapped (storage key).
    /// @param  curve          Bonding curve to seed with `supply`.
    /// @param  supply         Tokens to forward into the curve.
    /// @param  admin          Creator / admin recorded on the wrap.
    function configure(address externalToken, address curve, uint256 supply, address admin) external;
}
