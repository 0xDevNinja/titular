// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title IUniswapV2Router02 (minimal)
/// @notice Vendor-local subset of the Uniswap V2 router ABI used by the launchpad.
///         Covers only the surface the graduator needs: factory/WETH views and
///         the ERC20/ERC20 `addLiquidity` entry. Keeping this surface narrow
///         avoids pulling the full v2-periphery npm package as a Foundry
///         dependency and insulates compile-time solc ranges.
interface IUniswapV2Router02 {
    /// @notice Address of the canonical V2 factory the router writes through.
    function factory() external view returns (address);

    /// @notice WETH9 address wired into the router at deploy. Unused by the
    ///         graduator today but retained on the interface so fork tests and
    ///         downstream modules can share a single import.
    function WETH() external view returns (address);

    /// @notice ERC20/ERC20 liquidity add. See UniswapV2Router02 natspec for the
    ///         reference semantics. Graduator only calls the ERC20/ERC20 form —
    ///         both sides of an agent pair are ERC-20 by design.
    function addLiquidity(
        address tokenA,
        address tokenB,
        uint256 amountADesired,
        uint256 amountBDesired,
        uint256 amountAMin,
        uint256 amountBMin,
        address to,
        uint256 deadline
    ) external returns (uint256 amountA, uint256 amountB, uint256 liquidity);
}
