// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

/// @title  IUniswapV2Pair
/// @notice Minimal Uniswap V2 pair interface. Only the read methods consumed
///         by launchpad modules are declared.
/// @dev    The pair stores `price0CumulativeLast` and `price1CumulativeLast`
///         as UQ112x112 fixed-point cumulatives integrated over time. To
///         derive a Time-Weighted Average Price (TWAP), a consumer takes two
///         snapshots `(cumulative_t0, t0)` and `(cumulative_t1, t1)` and
///         computes:
///             priceAvg = (cumulative_t1 - cumulative_t0) / (t1 - t0)
///         The cumulatives are intentionally allowed to overflow modulo
///         2^256: the subtraction wraps cleanly so consumers MUST use
///         unchecked subtraction when reading them.
interface IUniswapV2Pair {
    /// @notice Address of `token0`, the lexicographically smaller token in
    ///         the pair.
    function token0() external view returns (address);

    /// @notice Address of `token1`, the lexicographically larger token in
    ///         the pair.
    function token1() external view returns (address);

    /// @notice Current reserves and the block timestamp of the last sync.
    /// @return reserve0           Reserve of `token0`.
    /// @return reserve1           Reserve of `token1`.
    /// @return blockTimestampLast uint32-wrapped timestamp of the last
    ///                            state-mutating interaction with the pair.
    function getReserves() external view returns (uint112 reserve0, uint112 reserve1, uint32 blockTimestampLast);

    /// @notice Cumulative `price0 = reserve1/reserve0` integrated over time
    ///         in UQ112x112 fixed-point. Wraps modulo 2^256 by design.
    function price0CumulativeLast() external view returns (uint256);

    /// @notice Cumulative `price1 = reserve0/reserve1` integrated over time
    ///         in UQ112x112 fixed-point. Wraps modulo 2^256 by design.
    function price1CumulativeLast() external view returns (uint256);
}
