// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IUniswapV2Factory} from "../../src/launchpad/interfaces/IUniswapV2Factory.sol";

/// @title MockUniswapV2Factory
/// @notice Minimal stand-in for the Uniswap V2 factory used in Graduator unit tests.
///         Records create-pair calls, returns a deterministic pair address per
///         unordered token tuple, and supports pre-seeding a pair address to
///         exercise the "pair already exists" code path.
contract MockUniswapV2Factory is IUniswapV2Factory {
    /// @dev Canonical ordering: (sorted_token0, sorted_token1) -> pair.
    mapping(address => mapping(address => address)) private _pairs;

    /// @notice Counter bumped on every `createPair` call. Used by tests to
    ///         assert the graduator skipped a redundant deploy.
    uint256 public createPairCalls;

    /// @notice Last pair deployed (or pre-seeded). Convenience accessor.
    address public lastPair;

    /// @notice Pre-seed a pair address for `(tokenA, tokenB)` to exercise the
    ///         already-deployed branch.
    function setPair(address tokenA, address tokenB, address pair) external {
        (address t0, address t1) = _sort(tokenA, tokenB);
        _pairs[t0][t1] = pair;
        lastPair = pair;
    }

    /// @inheritdoc IUniswapV2Factory
    function getPair(address tokenA, address tokenB) external view override returns (address) {
        (address t0, address t1) = _sort(tokenA, tokenB);
        return _pairs[t0][t1];
    }

    /// @inheritdoc IUniswapV2Factory
    function createPair(address tokenA, address tokenB) external override returns (address pair) {
        (address t0, address t1) = _sort(tokenA, tokenB);
        // Deterministic, non-zero stub address derived from the token tuple so
        // tests can assert the exact pair the router received.
        pair = address(uint160(uint256(keccak256(abi.encodePacked(t0, t1, "mock-pair")))));
        _pairs[t0][t1] = pair;
        lastPair = pair;
        createPairCalls += 1;
    }

    function _sort(address a, address b) private pure returns (address, address) {
        return a < b ? (a, b) : (b, a);
    }
}
