// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {Initializable} from "@openzeppelin/contracts/proxy/utils/Initializable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {IJob} from "./IJob.sol";
import {AgentRegistry} from "./AgentRegistry.sol";

/// @title Job
/// @notice ACP v2 Job — EIP-1167-clonable state machine for agent work agreements.
/// @dev    Designed to be deployed as a minimal proxy by JobFactory; `initialize`
///         replaces the constructor.  All privileged state transitions gate on role
///         (principal, agent, evaluator, arbiter) and the current Phase.
///
///         Phase graph (allowed transitions):
///           Open      → Active     (agent accepts)
///           Open      → Cancelled  (principal cancels; budget refunded)
///           Open      → Cancelled  (expiry triggers expireJob)
///           Active    → Review     (agent submits result)
///           Active    → Disputed   (principal raises dispute)
///           Active    → Cancelled  (expiry triggers expireJob)
///           Review    → Completed  (principal/evaluator approves; escrow released)
///           Review    → Active     (evaluator rejects; agent may resubmit)
///           Review    → Disputed   (principal or agent raises dispute)
///           Completed → —          (terminal)
///           Cancelled → —          (terminal)
///           Disputed  → Completed  (arbiter resolves agent-favoured; escrow released)
///           Disputed  → Cancelled  (arbiter resolves principal-favoured; budget refunded)
///
///         CEI order is enforced throughout; ReentrancyGuard on all token transfers.
contract Job is IJob, ReentrancyGuard, Initializable {
    using SafeERC20 for IERC20;

    // ─────────────────────────────────────────────────────────────
    // Immutable-ish state (written once in initialize)
    // ─────────────────────────────────────────────────────────────

    /// @notice Monotonic job identifier (set by JobFactory at clone-time via initialize).
    uint256 public jobId;

    /// @notice Address that created and funded this job.
    address public principal;

    /// @notice The registry through which agents are looked up.
    AgentRegistry public registry;

    /// @notice Agent registry ID that this job is targeted at (0 = open to any).
    uint256 public targetAgentId;

    /// @notice ERC-20 token used for payment (address(0) means native ETH — not supported in v1).
    address public token;

    /// @notice Total budget locked in this contract.
    uint256 public budget;

    /// @notice UNIX timestamp after which the job may be expired.
    uint64 public deadline;

    /// @notice Job type (Direct vs Evaluated).
    JobType public jobType;

    /// @notice Arbiter address authorised to resolve disputes.
    address public arbiter;

    // ─────────────────────────────────────────────────────────────
    // Mutable state
    // ─────────────────────────────────────────────────────────────

    /// @notice Current phase.
    Phase public phase;

    /// @notice Address of the accepting agent (set in `accept`).
    address public agent;

    /// @notice Agent registry ID of the accepting agent.
    uint256 public agentId;

    /// @notice Evaluator address (optional; only relevant for Evaluated job type).
    address public evaluator;

    /// @notice URI pointing to the submitted result (set in `submitResult`).
    string public resultURI;

    /// @notice Timestamp at which the result was submitted (used for grace period).
    uint64 public resultSubmittedAt;

    /// @notice Grace period after result submission before principal can force-release (Direct jobs).
    uint64 public constant RELEASE_GRACE = 48 hours;

    /// @notice Grace period for dispute window after result submission.
    uint64 public constant DISPUTE_GRACE = 72 hours;

    // ─────────────────────────────────────────────────────────────
    // Errors
    // ─────────────────────────────────────────────────────────────

    error ZeroAddress();
    error ZeroAmount();
    error InvalidDeadline();
    error WrongPhase(Phase current, Phase required);
    error NotPrincipal();
    error NotAgent();
    error NotEvaluator();
    error NotArbiter();
    error AgentNotFound(uint256 id);
    error AgentInactive(uint256 id);
    error JobNotExpired();
    error GracePeriodActive();
    error AlreadyInitialized();
    error InvalidJobType();
    error EvaluatorRequired();
    error DirectJobNoEvaluator();
    error TokenMismatch();

    // ─────────────────────────────────────────────────────────────
    // Modifiers
    // ─────────────────────────────────────────────────────────────

    modifier onlyPhase(Phase required) {
        if (phase != required) revert WrongPhase(phase, required);
        _;
    }

    modifier onlyPrincipal() {
        if (msg.sender != principal) revert NotPrincipal();
        _;
    }

    modifier onlyAgent() {
        if (msg.sender != agent) revert NotAgent();
        _;
    }

    modifier onlyEvaluator() {
        if (msg.sender != evaluator) revert NotEvaluator();
        _;
    }

    modifier onlyArbiter() {
        if (msg.sender != arbiter) revert NotArbiter();
        _;
    }

    // ─────────────────────────────────────────────────────────────
    // Initializer (replaces constructor for EIP-1167 clones)
    // ─────────────────────────────────────────────────────────────

    // ─────────────────────────────────────────────────────────────
    // Init config struct (avoids stack-too-deep in initializer)
    // ─────────────────────────────────────────────────────────────

    /// @notice Packed init parameters passed to `initialize`.
    struct InitParams {
        uint256 jobId;
        address principal;
        AgentRegistry registry;
        uint256 targetAgentId;
        address token;
        uint256 budget;
        uint64 deadline;
        JobType jobType;
        address evaluator;
        address arbiter;
    }

    /// @notice Initialise a newly-cloned Job contract.
    /// @param p Packed init parameters (see `InitParams`).
    function initialize(InitParams calldata p) external initializer {
        if (p.principal == address(0)) revert ZeroAddress();
        if (address(p.registry) == address(0)) revert ZeroAddress();
        if (p.token == address(0)) revert ZeroAddress();
        if (p.budget == 0) revert ZeroAmount();
        // slither-disable-next-line timestamp
        if (p.deadline <= block.timestamp) revert InvalidDeadline();
        if (p.arbiter == address(0)) revert ZeroAddress();
        if (p.jobType == JobType.Evaluated && p.evaluator == address(0)) revert EvaluatorRequired();
        if (p.jobType == JobType.Direct && p.evaluator != address(0)) revert DirectJobNoEvaluator();

        jobId = p.jobId;
        principal = p.principal;
        registry = p.registry;
        targetAgentId = p.targetAgentId;
        token = p.token;
        budget = p.budget;
        deadline = p.deadline;
        jobType = p.jobType;
        evaluator = p.evaluator;
        arbiter = p.arbiter;

        // phase defaults to Open (Phase(0))

        emit JobInitialised(p.jobId, p.principal, p.targetAgentId, p.jobType, p.budget, p.token, p.deadline);
        if (p.evaluator != address(0)) {
            emit EvaluatorAssigned(p.jobId, p.evaluator);
        }
    }

    // ─────────────────────────────────────────────────────────────
    // Phase transitions
    // ─────────────────────────────────────────────────────────────

    /// @notice Accept the job. Transitions Open → Active.
    /// @dev    If targetAgentId != 0, the caller's agentId must match. The caller must
    ///         control an active agent in the registry.
    /// @param _agentId Registry ID of the accepting agent.
    function accept(uint256 _agentId) external onlyPhase(Phase.Open) nonReentrant {
        AgentRegistry.AgentInfo memory info = registry.getAgent(_agentId);
        if (info.controller != msg.sender) revert AgentNotFound(_agentId);
        if (!info.active) revert AgentInactive(_agentId);
        if (targetAgentId != 0 && targetAgentId != _agentId) revert AgentNotFound(_agentId);

        // Effects
        phase = Phase.Active;
        agent = msg.sender;
        agentId = _agentId;

        emit AgentAccepted(jobId, msg.sender, _agentId);
    }

    /// @notice Submit a result URI. Transitions Active → Review.
    /// @param _resultURI IPFS / HTTPS URI of the deliverable.
    function submitResult(string calldata _resultURI) external onlyPhase(Phase.Active) onlyAgent nonReentrant {
        // Effects
        phase = Phase.Review;
        resultURI = _resultURI;
        // slither-disable-next-line timestamp
        resultSubmittedAt = uint64(block.timestamp);

        emit ResultSubmitted(jobId, msg.sender, _resultURI);
    }

    /// @notice Approve a submitted result. Transitions Review → Completed.
    /// @dev    For Evaluated jobs only the evaluator may call.
    ///         For Direct jobs the principal may call, or the principal may call `release`
    ///         after the grace period.
    function approveResult() external onlyPhase(Phase.Review) nonReentrant {
        if (jobType == JobType.Evaluated) {
            if (msg.sender != evaluator) revert NotEvaluator();
            emit ResultApproved(jobId, msg.sender);
        } else {
            // Direct — principal approves immediately
            if (msg.sender != principal) revert NotPrincipal();
        }

        _completeJob();
    }

    /// @notice Reject a submitted result. Transitions Review → Active (agent may resubmit).
    /// @dev    Only the evaluator may call (Evaluated jobs). For Direct jobs the principal
    ///         should raise a dispute instead.
    /// @param reason Human-readable reason logged to calldata only.
    function rejectResult(string calldata reason) external onlyPhase(Phase.Review) onlyEvaluator {
        // Effects
        phase = Phase.Active;
        resultURI = "";
        resultSubmittedAt = 0;

        emit ResultRejected(jobId, msg.sender, reason);
    }

    /// @notice Release budget to the agent after grace period (Direct jobs) or after
    ///         evaluator approval window. Transitions Review → Completed.
    /// @dev    Callable by anyone once the release grace period has elapsed.
    function release() external onlyPhase(Phase.Review) nonReentrant {
        // slither-disable-next-line timestamp
        if (block.timestamp < resultSubmittedAt + RELEASE_GRACE) revert GracePeriodActive();
        _completeJob();
    }

    /// @notice Cancel the job. Refunds the principal.
    ///         Allowed from Open (any time) or Active/Review (principal only, before agent submits final).
    /// @param reason Human-readable reason.
    function cancel(string calldata reason) external nonReentrant {
        if (msg.sender != principal) revert NotPrincipal();
        Phase current = phase;
        if (current == Phase.Completed || current == Phase.Cancelled || current == Phase.Disputed) {
            revert WrongPhase(current, Phase.Open); // terminal / disputed — cannot cancel
        }

        // Effects before interaction
        phase = Phase.Cancelled;
        emit JobCancelled(jobId, msg.sender, reason);

        // Interaction
        uint256 refund = budget;
        budget = 0;
        IERC20(token).safeTransfer(principal, refund);
    }

    /// @notice Raise a dispute. Transitions Active or Review → Disputed.
    /// @dev    Principal or agent may raise a dispute.
    function raiseDispute() external nonReentrant {
        if (msg.sender != principal && msg.sender != agent) revert NotPrincipal();
        Phase current = phase;
        if (current != Phase.Active && current != Phase.Review) {
            revert WrongPhase(current, Phase.Active);
        }
        // Effects
        phase = Phase.Disputed;
        emit DisputeRaised(jobId, msg.sender);
    }

    /// @notice Resolve a dispute. Transitions Disputed → Completed (agent wins) or Cancelled (principal wins).
    /// @param agentFavoured If true, releases budget to agent; if false, refunds principal.
    function resolveDispute(bool agentFavoured) external onlyPhase(Phase.Disputed) onlyArbiter nonReentrant {
        if (agentFavoured) {
            emit DisputeResolved(jobId, msg.sender, true);
            _completeJob();
        } else {
            phase = Phase.Cancelled;
            emit DisputeResolved(jobId, msg.sender, false);
            emit JobCancelled(jobId, msg.sender, "dispute resolved: principal favoured");

            uint256 refund = budget;
            budget = 0;
            IERC20(token).safeTransfer(principal, refund);
        }
    }

    /// @notice Expire a job that has passed its deadline without being completed.
    ///         Open → Cancelled (refund principal), Active → Cancelled (refund principal).
    function expireJob() external nonReentrant {
        // slither-disable-next-line timestamp
        if (block.timestamp <= deadline) revert JobNotExpired();
        Phase current = phase;
        if (current != Phase.Open && current != Phase.Active) {
            revert WrongPhase(current, Phase.Open);
        }

        phase = Phase.Cancelled;
        emit JobCancelled(jobId, msg.sender, "expired");

        uint256 refund = budget;
        budget = 0;
        IERC20(token).safeTransfer(principal, refund);
    }

    // ─────────────────────────────────────────────────────────────
    // Views
    // ─────────────────────────────────────────────────────────────

    /// @notice Return the current phase of this job.
    /// @return current Current Phase enum value.
    function currentPhase() external view returns (Phase current) {
        current = phase;
    }

    // ─────────────────────────────────────────────────────────────
    // Internal helpers
    // ─────────────────────────────────────────────────────────────

    /// @dev Finalise job: set terminal state, emit event, transfer budget to agent.
    ///      Must be called before external token transfer (CEI).
    function _completeJob() internal {
        phase = Phase.Completed;
        uint256 payout = budget;
        budget = 0;
        address recipient = agent;

        emit JobCompleted(jobId, msg.sender, payout);

        // Interaction last
        IERC20(token).safeTransfer(recipient, payout);
    }
}
