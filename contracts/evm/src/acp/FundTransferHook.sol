// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {IHook} from "./IHook.sol";

/// @title FundTransferHook
/// @notice Hook that handles principal entrustment and profit-and-loss settlement.
///
///         When `onApprove` is triggered (job result approved):
///           - The hook releases the job's escrowed budget to the agent (primary share).
///           - A configurable PnL settlement percentage is forwarded to the principal
///             or to a separate settlement address.
///
///         Context encoding for `onApprove`:
///           abi.encode(address agent, address token, uint256 amount, uint256 pnlBps, address settlementAddr)
///
///         Other lifecycle hooks are no-ops (this hook only acts on approval).
contract FundTransferHook is IHook, ReentrancyGuard {
    using SafeERC20 for IERC20;

    // ─────────────────────────────────────────────────────────────
    // Types
    // ─────────────────────────────────────────────────────────────

    struct ApproveContext {
        address agent;
        address token;
        uint256 amount;
        uint256 pnlBps;         // PnL share to retain for settlement (BPS of amount)
        address settlementAddr; // Where the PnL share goes (0 = skip PnL split)
    }

    // ─────────────────────────────────────────────────────────────
    // Constants
    // ─────────────────────────────────────────────────────────────

    uint256 public constant BPS = 10_000;

    // ─────────────────────────────────────────────────────────────
    // Events
    // ─────────────────────────────────────────────────────────────

    /// @notice Emitted when a fund transfer with PnL settlement is executed.
    event FundTransferred(
        uint256 indexed jobId,
        address indexed agent,
        address indexed token,
        uint256 agentAmount,
        uint256 pnlAmount,
        address settlementAddr
    );

    // ─────────────────────────────────────────────────────────────
    // Errors
    // ─────────────────────────────────────────────────────────────

    error InvalidBps();
    error ZeroAddress();
    error ZeroAmount();

    // ─────────────────────────────────────────────────────────────
    // IHook implementation
    // ─────────────────────────────────────────────────────────────

    /// @inheritdoc IHook
    function onAccept(uint256, bytes calldata) external pure override {}

    /// @inheritdoc IHook
    function onSubmit(uint256, bytes calldata) external pure override {}

    /// @notice Execute PnL settlement and fund transfer to agent.
    /// @param jobId   The job identifier.
    /// @param context ABI-encoded ApproveContext.
    function onApprove(uint256 jobId, bytes calldata context) external override nonReentrant {
        ApproveContext memory ctx = abi.decode(context, (ApproveContext));
        if (ctx.agent == address(0)) revert ZeroAddress();
        if (ctx.token == address(0)) revert ZeroAddress();
        if (ctx.amount == 0) revert ZeroAmount();
        if (ctx.pnlBps > BPS) revert InvalidBps();

        uint256 pnlAmount = 0;
        uint256 agentAmount = ctx.amount;

        if (ctx.pnlBps > 0 && ctx.settlementAddr != address(0)) {
            pnlAmount = (ctx.amount * ctx.pnlBps) / BPS;
            agentAmount = ctx.amount - pnlAmount;
        }

        emit FundTransferred(jobId, ctx.agent, ctx.token, agentAmount, pnlAmount, ctx.settlementAddr);

        // Transfer agent share
        IERC20(ctx.token).safeTransferFrom(msg.sender, ctx.agent, agentAmount);

        // Transfer PnL share if applicable
        if (pnlAmount > 0 && ctx.settlementAddr != address(0)) {
            IERC20(ctx.token).safeTransferFrom(msg.sender, ctx.settlementAddr, pnlAmount);
        }
    }

    /// @inheritdoc IHook
    function onReject(uint256, bytes calldata) external pure override {}

    /// @inheritdoc IHook
    function onCancel(uint256, bytes calldata) external pure override {}

    /// @inheritdoc IHook
    function hookName() external pure override returns (string memory) {
        return "FundTransferHook";
    }
}
