// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IUniswapV2Pair} from "../../src/launchpad/interfaces/IUniswapV2Pair.sol";

/// @title  MockUniswapV2Pair
/// @notice Minimal Uniswap V2 pair stand-in for launchpad unit tests.
///         Exposes the `IUniswapV2Pair` read surface plus setters so tests
///         can deterministically drive reserves AND the
///         `priceXCumulativeLast` accumulators consumed by TWAP-based
///         consumers like `CapitalFormationModule`.
/// @dev    Cumulative semantics mirror the real pair: each call to
///         {sync} advances the cumulatives by `priceX * timeElapsed` where
///         `timeElapsed` is `currentTimestamp - blockTimestampLast`.
///         Tests can either:
///           (a) call {setReserves} for spot-only consumers that ignore
///               the cumulatives;
///           (b) call {sync} repeatedly across `vm.warp` jumps to grow
///               the cumulatives the way a live pair would.
///         Does NOT implement swap / mint / burn.
contract MockUniswapV2Pair is IUniswapV2Pair {
    uint256 private constant Q112 = 0x10000000000000000000000000000;

    address private _token0;
    address private _token1;
    uint112 private _reserve0;
    uint112 private _reserve1;
    uint32 private _blockTimestampLast;
    uint256 private _price0CumulativeLast;
    uint256 private _price1CumulativeLast;

    constructor(address tokenA, address tokenB) {
        // Mirror the real pair's lexicographic ordering: token0 is the
        // numerically smaller address.
        (address t0, address t1) = tokenA < tokenB ? (tokenA, tokenB) : (tokenB, tokenA);
        _token0 = t0;
        _token1 = t1;
    }

    /// @notice Hard-set reserves WITHOUT advancing cumulatives. Used by
    ///         spot-only tests; callers that need TWAP semantics should
    ///         use {sync} instead.
    /// @dev    Stamps `block.timestamp` as last-sync time so the returned
    ///         tuple looks fresh.
    function setReserves(uint112 reserve0_, uint112 reserve1_) external {
        _reserve0 = reserve0_;
        _reserve1 = reserve1_;
        _blockTimestampLast = uint32(block.timestamp);
    }

    /// @notice Update reserves AND advance the cumulatives over the time
    ///         elapsed since the last {sync} / {setReserves} using the
    ///         PRE-update reserves — i.e. the price the pair held during
    ///         the elapsed window.
    /// @dev    Models the real pair's `_update`. Cumulatives wrap modulo
    ///         2^256 the same way: we use unchecked addition.
    function sync(uint112 reserve0_, uint112 reserve1_) external {
        uint32 blockTimestamp = uint32(block.timestamp);
        uint32 timeElapsed;
        unchecked {
            // uint32 wrap mirrors Uniswap V2 behaviour.
            timeElapsed = blockTimestamp - _blockTimestampLast;
        }
        if (timeElapsed > 0 && _reserve0 != 0 && _reserve1 != 0) {
            // price0 = reserve1/reserve0 in UQ112x112 → (reserve1 << 112) / reserve0.
            // price1 = reserve0/reserve1 in UQ112x112 → (reserve0 << 112) / reserve1.
            unchecked {
                _price0CumulativeLast += (uint256(_reserve1) * Q112 / uint256(_reserve0)) * uint256(timeElapsed);
                _price1CumulativeLast += (uint256(_reserve0) * Q112 / uint256(_reserve1)) * uint256(timeElapsed);
            }
        }
        _reserve0 = reserve0_;
        _reserve1 = reserve1_;
        _blockTimestampLast = blockTimestamp;
    }

    /// @notice Direct-write the cumulatives. Convenience for tests that
    ///         want to skip the `sync`-warp dance and inject a precomputed
    ///         delta straight into the accumulator.
    function setCumulatives(uint256 p0, uint256 p1, uint32 ts) external {
        _price0CumulativeLast = p0;
        _price1CumulativeLast = p1;
        _blockTimestampLast = ts;
    }

    function token0() external view returns (address) {
        return _token0;
    }

    function token1() external view returns (address) {
        return _token1;
    }

    function getReserves() external view returns (uint112 reserve0, uint112 reserve1, uint32 blockTimestampLast) {
        return (_reserve0, _reserve1, _blockTimestampLast);
    }

    function price0CumulativeLast() external view returns (uint256) {
        return _price0CumulativeLast;
    }

    function price1CumulativeLast() external view returns (uint256) {
        return _price1CumulativeLast;
    }
}
