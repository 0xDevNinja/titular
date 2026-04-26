// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import {IUniswapV2Router02} from "../../src/launchpad/interfaces/IUniswapV2Router02.sol";
import {MockUniswapV2Factory} from "./MockUniswapV2Factory.sol";

/// @title MockUniswapV2Router
/// @notice Minimal router stub that:
///         - reports a configurable `factory` / `WETH`;
///         - records the last `addLiquidity` call parameters so tests can
///           assert the graduator passed through its pulled reserves with the
///           right slippage floors, recipient, and deadline;
///         - actually pulls both tokens from the caller (so allowance and
///           balance checks on the graduator are exercised) and returns a
///           deterministic LP amount so `Graduated` event assertions are
///           stable.
contract MockUniswapV2Router is IUniswapV2Router02 {
    address public immutable override factory;

    /// @dev Returned from `WETH()`. Not used by the graduator; pinned to a
    ///      sentinel so any accidental use blows up fast.
    address public constant override WETH = address(0xdeaDDeADDEaDdeaDdEAddEADDEAdDeadDEADDEaD);

    struct AddLiquidityCall {
        address tokenA;
        address tokenB;
        uint256 amountADesired;
        uint256 amountBDesired;
        uint256 amountAMin;
        uint256 amountBMin;
        address to;
        uint256 deadline;
        bool recorded;
    }

    /// @notice Most recent `addLiquidity` parameter set. `recorded == true`
    ///         after the first call; tests can key on it as a presence check.
    AddLiquidityCall public lastCall;

    /// @notice Number of `addLiquidity` invocations against this router.
    uint256 public addLiquidityCalls;

    /// @notice LP token quantity the router reports as minted. Deterministic
    ///         stub so the `Graduated` event and downstream accounting are
    ///         assertable. Can be tweaked via {setLiquidityOut}.
    uint256 public liquidityOut = 1_234_567;

    constructor(address factory_) {
        factory = factory_;
    }

    /// @notice Override the LP quantity reported by subsequent `addLiquidity`
    ///         calls. Useful for edge-case tests (e.g. slippage clamp).
    function setLiquidityOut(uint256 v) external {
        liquidityOut = v;
    }

    /// @inheritdoc IUniswapV2Router02
    function addLiquidity(
        address tokenA,
        address tokenB,
        uint256 amountADesired,
        uint256 amountBDesired,
        uint256 amountAMin,
        uint256 amountBMin,
        address to,
        uint256 deadline
    ) external override returns (uint256 amountA, uint256 amountB, uint256 liquidity) {
        lastCall = AddLiquidityCall({
            tokenA: tokenA,
            tokenB: tokenB,
            amountADesired: amountADesired,
            amountBDesired: amountBDesired,
            amountAMin: amountAMin,
            amountBMin: amountBMin,
            to: to,
            deadline: deadline,
            recorded: true
        });
        addLiquidityCalls += 1;

        // Pull both tokens from the caller to exercise the allowance path on
        // the graduator. Reverts here mean the graduator forgot to approve.
        IERC20(tokenA).transferFrom(msg.sender, address(this), amountADesired);
        IERC20(tokenB).transferFrom(msg.sender, address(this), amountBDesired);

        amountA = amountADesired;
        amountB = amountBDesired;
        liquidity = liquidityOut;
    }
}
