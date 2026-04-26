// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

/// @title IJob
/// @notice Interface for the ACP v2 Job state machine.
interface IJob {
    // ─────────────────────────────────────────────────────────────
    // Types
    // ─────────────────────────────────────────────────────────────

    /// @notice Ordered job lifecycle phases.
    enum Phase {
        Open, // 0 — awaiting agent acceptance
        Active, // 1 — agent accepted; work in progress
        Review, // 2 — result submitted; in evaluator / principal review
        Completed, // 3 — payment released
        Cancelled, // 4 — cancelled before completion (by principal before acceptance, or by dispute resolution)
        Disputed // 5 — dispute raised; awaiting arbitration
    }

    /// @notice Job type: whether an external evaluator is required before release.
    enum JobType {
        Direct, // No evaluator — principal releases on result submission after grace period
        Evaluated // Evaluator role is required to approve / reject before release
    }

    // ─────────────────────────────────────────────────────────────
    // Events
    // ─────────────────────────────────────────────────────────────

    event JobInitialised(
        uint256 indexed jobId,
        address indexed principal,
        uint256 indexed agentId,
        JobType jobType,
        uint256 budget,
        address token,
        uint64 deadline
    );

    event AgentAccepted(uint256 indexed jobId, address indexed agent, uint256 indexed agentId);

    event ResultSubmitted(uint256 indexed jobId, address indexed agent, string resultURI);

    event JobCompleted(uint256 indexed jobId, address indexed releasedBy, uint256 amount);

    event JobCancelled(uint256 indexed jobId, address indexed cancelledBy, string reason);

    event DisputeRaised(uint256 indexed jobId, address indexed raisedBy);

    event DisputeResolved(uint256 indexed jobId, address indexed resolver, bool agentFavoured);

    event EvaluatorAssigned(uint256 indexed jobId, address indexed evaluator);

    event ResultApproved(uint256 indexed jobId, address indexed evaluator);

    event ResultRejected(uint256 indexed jobId, address indexed evaluator, string reason);

    // ─────────────────────────────────────────────────────────────
    // State transitions (phase function signatures)
    // ─────────────────────────────────────────────────────────────

    function accept(uint256 agentId) external;
    function submitResult(string calldata resultURI) external;
    function approveResult() external;
    function rejectResult(string calldata reason) external;
    function release() external;
    function cancel(string calldata reason) external;
    function raiseDispute() external;
    function resolveDispute(bool agentFavoured) external;
    function expireJob() external;
}
