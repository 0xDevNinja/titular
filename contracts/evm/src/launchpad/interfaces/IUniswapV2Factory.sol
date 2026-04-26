// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title IUniswapV2Factory (minimal)
/// @notice Vendor-local subset of the Uniswap V2 factory ABI used by the
///         launchpad. Only `getPair` and `createPair` are needed — every other
///         entry (pair enumeration, fee toggle, etc.) is out of scope.
interface IUniswapV2Factory {
    /// @notice Returns the existing pair for `(tokenA, tokenB)` or `address(0)`
    ///         if no pair has been deployed yet.
    function getPair(address tokenA, address tokenB) external view returns (address pair);

    /// @notice Deploys the V2 pair for `(tokenA, tokenB)` via CREATE2. Reverts
    ///         inside the router on duplicate deploy; callers should gate on
    ///         {getPair} first.
    function createPair(address tokenA, address tokenB) external returns (address pair);
}
