// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title IBondingCurve
/// @notice Minimal view + graduation surface of `BondingCurve` that downstream
///         modules (graduator, fee router, launch-radar indexer) bind against
///         without importing the full implementation and its OZ dependencies.
/// @dev    Mirrors the subset of public/external signatures on
///         `src/launchpad/BondingCurve.sol` that live callers consume. Keep
///         signatures and return ordering in lockstep.
interface IBondingCurve {
    /// @notice Agent ERC-20 traded by this curve.
    function agentToken() external view returns (address);

    /// @notice Quote ERC-20 (TITU) used for trading and the LP seed.
    function quoteToken() external view returns (address);

    /// @notice Flipped once `requestGraduation` has fired. Trading halts and
    ///         the graduator is permitted to pull reserves exactly once.
    function graduated() external view returns (bool);

    /// @notice Contract authorized to drain the curve post-graduation.
    function graduator() external view returns (address);

    /// @notice Current real TITU reserve held by the curve.
    function realQuoteReserve() external view returns (uint256);

    /// @notice Current real agent-token inventory held by the curve.
    function realAgentReserve() external view returns (uint256);

    /// @notice Drain both reserves to the graduator. Graduator-only, one-shot.
    /// @return quoteAmount Quote tokens transferred.
    /// @return agentAmount Agent tokens transferred.
    function pullForGraduation() external returns (uint256 quoteAmount, uint256 agentAmount);

    /// @notice Flip the curve to the graduated state once the real quote reserve
    ///         has reached the graduation threshold. Keeper-style trigger.
    function requestGraduation() external;
}
