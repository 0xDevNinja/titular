// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Clones} from "@openzeppelin/contracts/proxy/Clones.sol";
import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {Pausable} from "@openzeppelin/contracts/utils/Pausable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Job} from "./Job.sol";
import {IJob} from "./IJob.sol";
import {AgentRegistry} from "./AgentRegistry.sol";

/// @title JobFactory
/// @notice Deploys EIP-1167 minimal-proxy clones of the `Job` implementation and
///         transfers the budget ERC-20 tokens into each clone before calling `initialize`.
/// @dev    Access control:
///           - DEFAULT_ADMIN_ROLE : set implementation address, pause/unpause.
///           - PAUSER_ROLE        : pause/unpause.
///         Any address may call `createJob`; the caller becomes the `principal`.
///
///         Token flow:
///           1. Caller approves JobFactory for `budget` of `token`.
///           2. `createJob` transfers `budget` from caller → new Job clone.
///           3. Job clone is initialised; budget is locked inside the clone.
contract JobFactory is AccessControl, Pausable {
    using SafeERC20 for IERC20;
    using Clones for address;

    // ─────────────────────────────────────────────────────────────
    // Roles
    // ─────────────────────────────────────────────────────────────

    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");

    // ─────────────────────────────────────────────────────────────
    // State
    // ─────────────────────────────────────────────────────────────

    /// @notice Current Job implementation address used for EIP-1167 cloning.
    address public jobImplementation;

    /// @notice AgentRegistry used to validate agent IDs at job creation time.
    AgentRegistry public immutable registry;

    /// @notice Default arbiter for dispute resolution (configurable by admin).
    address public defaultArbiter;

    /// @notice Monotonic counter; every clone gets a unique jobId.
    uint256 private _nextJobId;

    /// @notice jobId => clone address.
    mapping(uint256 => address) public jobs;

    /// @notice principal => list of jobIds created by that principal.
    mapping(address => uint256[]) private _principalJobs;

    // ─────────────────────────────────────────────────────────────
    // Events
    // ─────────────────────────────────────────────────────────────

    /// @notice Emitted when a new Job clone is deployed and funded.
    event JobCreated(
        uint256 indexed jobId,
        address indexed clone,
        address indexed principal,
        IJob.JobType jobType,
        address token,
        uint256 budget
    );

    /// @notice Emitted when the job implementation is updated.
    event ImplementationUpdated(address indexed oldImpl, address indexed newImpl);

    /// @notice Emitted when the default arbiter is updated.
    event DefaultArbiterUpdated(address indexed oldArbiter, address indexed newArbiter);

    // ─────────────────────────────────────────────────────────────
    // Errors
    // ─────────────────────────────────────────────────────────────

    error ZeroAddress();
    error ZeroAmount();
    error InvalidDeadline();

    // ─────────────────────────────────────────────────────────────
    // Constructor
    // ─────────────────────────────────────────────────────────────

    /// @param admin           Address granted DEFAULT_ADMIN_ROLE.
    /// @param _jobImpl        Initial Job implementation for cloning.
    /// @param _registry       AgentRegistry reference.
    /// @param _defaultArbiter Default arbiter for all created jobs.
    constructor(address admin, address _jobImpl, AgentRegistry _registry, address _defaultArbiter) {
        if (admin == address(0)) revert ZeroAddress();
        if (_jobImpl == address(0)) revert ZeroAddress();
        if (address(_registry) == address(0)) revert ZeroAddress();
        if (_defaultArbiter == address(0)) revert ZeroAddress();

        jobImplementation = _jobImpl;
        registry = _registry;
        defaultArbiter = _defaultArbiter;

        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(PAUSER_ROLE, admin);
    }

    // ─────────────────────────────────────────────────────────────
    // Job creation
    // ─────────────────────────────────────────────────────────────

    /// @notice Create a new Job clone, transfer the budget into it, and initialise it.
    /// @dev    Caller must have pre-approved this contract for `budget` of `token`.
    /// @param  targetAgentId Agent registry ID to target (0 = open to any active agent).
    /// @param  token         ERC-20 token address for payment.
    /// @param  budget        Amount to lock into the escrow.
    /// @param  deadline      UNIX timestamp after which the job may be expired.
    /// @param  jobType       Direct (no evaluator) or Evaluated.
    /// @param  evaluator     Evaluator address; must be non-zero iff jobType == Evaluated.
    /// @param  arbiter       Arbiter override; if address(0), uses defaultArbiter.
    /// @return jobId         The monotonic ID assigned to this job.
    /// @return clone         The address of the deployed Job clone.
    function createJob(
        uint256 targetAgentId,
        address token,
        uint256 budget,
        uint64 deadline,
        IJob.JobType jobType,
        address evaluator,
        address arbiter
    ) external whenNotPaused returns (uint256 jobId, address clone) {
        if (token == address(0)) revert ZeroAddress();
        if (budget == 0) revert ZeroAmount();
        // slither-disable-next-line timestamp
        if (deadline <= block.timestamp) revert InvalidDeadline();

        address resolvedArbiter = arbiter == address(0) ? defaultArbiter : arbiter;

        jobId = _nextJobId++;
        clone = jobImplementation.clone();

        jobs[jobId] = clone;
        _principalJobs[msg.sender].push(jobId);

        // Emit before external calls (satisfies CEI for event ordering)
        emit JobCreated(jobId, clone, msg.sender, jobType, token, budget);

        // Transfer budget from principal → clone before initialise
        IERC20(token).safeTransferFrom(msg.sender, clone, budget);

        // Initialise the clone
        Job.InitParams memory p = Job.InitParams({
            jobId: jobId,
            principal: msg.sender,
            registry: registry,
            targetAgentId: targetAgentId,
            token: token,
            budget: budget,
            deadline: deadline,
            jobType: jobType,
            evaluator: evaluator,
            arbiter: resolvedArbiter
        });
        Job(clone).initialize(p);
    }

    // ─────────────────────────────────────────────────────────────
    // Admin
    // ─────────────────────────────────────────────────────────────

    /// @notice Update the Job implementation used for future clones.
    /// @param  newImpl New implementation address (must be non-zero).
    function setJobImplementation(address newImpl) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (newImpl == address(0)) revert ZeroAddress();
        address old = jobImplementation;
        jobImplementation = newImpl;
        emit ImplementationUpdated(old, newImpl);
    }

    /// @notice Update the default arbiter.
    /// @param  newArbiter New default arbiter (must be non-zero).
    function setDefaultArbiter(address newArbiter) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (newArbiter == address(0)) revert ZeroAddress();
        address old = defaultArbiter;
        defaultArbiter = newArbiter;
        emit DefaultArbiterUpdated(old, newArbiter);
    }

    /// @notice Pause job creation.
    function pause() external onlyRole(PAUSER_ROLE) {
        _pause();
    }

    /// @notice Unpause job creation.
    function unpause() external onlyRole(PAUSER_ROLE) {
        _unpause();
    }

    // ─────────────────────────────────────────────────────────────
    // Views
    // ─────────────────────────────────────────────────────────────

    /// @notice Return the clone address for a given jobId.
    /// @param  jobId The job ID to look up.
    /// @return clone Address of the Job clone, or address(0) if not found.
    function getJob(uint256 jobId) external view returns (address clone) {
        clone = jobs[jobId];
    }

    /// @notice Return all jobIds created by a given principal.
    /// @param  principal The principal address.
    /// @return jobIds    Array of job IDs created by this principal.
    function jobsByPrincipal(address principal) external view returns (uint256[] memory jobIds) {
        jobIds = _principalJobs[principal];
    }

    /// @notice Total number of jobs ever created.
    /// @return count Total job count.
    function totalJobs() external view returns (uint256 count) {
        count = _nextJobId;
    }
}
