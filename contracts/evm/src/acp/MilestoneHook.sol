// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {IHook} from "./IHook.sol";

/// @title MilestoneHook
/// @notice Splits a job budget into N equal stages.  Each `onApprove` releases
///         one stage payment to the agent.  The principal may cancel remaining
///         stages via `onCancel`.
///
///         Context for `onAccept`:
///           abi.encode(address agent, address token, uint256 totalBudget, uint8 stages)
///
///         Each `onApprove` call releases `totalBudget / stages` to the agent
///         (final stage gets dust remainder).
///
///         State is tracked per jobId.
contract MilestoneHook is IHook {
    using SafeERC20 for IERC20;

    // ─────────────────────────────────────────────────────────────
    // Types
    // ─────────────────────────────────────────────────────────────

    struct MilestoneState {
        address agent;
        address token;
        uint256 stageAmount; // budget / stages
        uint256 lastStageRemainder; // budget % stages (paid on last stage)
        uint8 totalStages;
        uint8 completedStages;
        bool initialised;
    }

    struct InitContext {
        address agent;
        address token;
        uint256 totalBudget;
        uint8 stages;
    }

    // ─────────────────────────────────────────────────────────────
    // State
    // ─────────────────────────────────────────────────────────────

    /// @notice jobId => MilestoneState.
    mapping(uint256 => MilestoneState) public milestones;

    // ─────────────────────────────────────────────────────────────
    // Events
    // ─────────────────────────────────────────────────────────────

    /// @notice Emitted when a milestone stage is released.
    event MilestoneReleased(
        uint256 indexed jobId, address indexed agent, uint8 stage, uint8 totalStages, uint256 amount
    );

    /// @notice Emitted when remaining milestones are cancelled.
    event MilestoneCancelled(uint256 indexed jobId, uint8 remainingStages);

    // ─────────────────────────────────────────────────────────────
    // Errors
    // ─────────────────────────────────────────────────────────────

    error ZeroAddress();
    error ZeroAmount();
    error InvalidStages();
    error NotInitialised(uint256 jobId);
    error AllStagesComplete(uint256 jobId);

    // ─────────────────────────────────────────────────────────────
    // IHook implementation
    // ─────────────────────────────────────────────────────────────

    /// @notice Initialise milestone state when agent accepts the job.
    /// @param jobId   The job identifier.
    /// @param context ABI-encoded InitContext.
    function onAccept(uint256 jobId, bytes calldata context) external override {
        InitContext memory ctx = abi.decode(context, (InitContext));
        if (ctx.agent == address(0)) revert ZeroAddress();
        if (ctx.token == address(0)) revert ZeroAddress();
        if (ctx.totalBudget == 0) revert ZeroAmount();
        if (ctx.stages == 0) revert InvalidStages();

        // Intentional divide-then-multiply to compute exact per-stage amount and dust remainder.
        // slither-disable-next-line divide-before-multiply
        uint256 stageAmt = ctx.totalBudget / ctx.stages;
        uint256 remainder = ctx.totalBudget % ctx.stages;

        milestones[jobId] = MilestoneState({
            agent: ctx.agent,
            token: ctx.token,
            stageAmount: stageAmt,
            lastStageRemainder: remainder,
            totalStages: ctx.stages,
            completedStages: 0,
            initialised: true
        });
    }

    /// @inheritdoc IHook
    function onSubmit(uint256, bytes calldata) external pure override {}

    /// @notice Release the next milestone stage payment to the agent.
    /// @param jobId The job identifier.
    function onApprove(uint256 jobId, bytes calldata) external override {
        MilestoneState storage ms = milestones[jobId];
        if (!ms.initialised) revert NotInitialised(jobId);
        if (ms.completedStages >= ms.totalStages) revert AllStagesComplete(jobId);

        uint8 stage = ms.completedStages + 1;
        uint256 payout = ms.stageAmount;

        // Last stage gets remainder
        if (stage == ms.totalStages) {
            payout += ms.lastStageRemainder;
        }

        // Effects
        ms.completedStages = stage;
        address agent = ms.agent;
        address token = ms.token;

        emit MilestoneReleased(jobId, agent, stage, ms.totalStages, payout);

        // Interaction
        IERC20(token).safeTransferFrom(msg.sender, agent, payout);
    }

    /// @inheritdoc IHook
    function onReject(uint256, bytes calldata) external pure override {}

    /// @notice Cancel remaining milestone stages when the job is cancelled.
    /// @param jobId The job identifier.
    function onCancel(uint256 jobId, bytes calldata) external override {
        MilestoneState storage ms = milestones[jobId];
        if (!ms.initialised) return;

        uint8 remaining = ms.totalStages - ms.completedStages;
        if (remaining > 0) {
            ms.completedStages = ms.totalStages; // mark all done (cancelled)
            emit MilestoneCancelled(jobId, remaining);
        }
    }

    /// @inheritdoc IHook
    function hookName() external pure override returns (string memory) {
        return "MilestoneHook";
    }
}
