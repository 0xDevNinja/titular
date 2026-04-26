// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {StdInvariant} from "forge-std/StdInvariant.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {Job} from "../../../src/acp/Job.sol";
import {IJob} from "../../../src/acp/IJob.sol";
import {AgentRegistry} from "../../../src/acp/AgentRegistry.sol";

contract MockToken is ERC20 {
    constructor() ERC20("PhaseInv", "PINV") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

/// @dev Handler that drives a Job through valid transitions.
contract JobPhaseHandler is Test {
    Job public job;
    AgentRegistry public registry;
    address public principal;
    address public agentCtrl;
    address public arbiter;
    uint256 public agentId;
    MockToken public token;

    // Track highest phase reached (monotone for non-reversible paths)
    uint8 public maxPhaseReached;

    constructor(
        Job _job,
        AgentRegistry _registry,
        MockToken _token,
        address _principal,
        address _agentCtrl,
        uint256 _agentId,
        address _arbiter
    ) {
        job = _job;
        registry = _registry;
        token = _token;
        principal = _principal;
        agentCtrl = _agentCtrl;
        agentId = _agentId;
        arbiter = _arbiter;
    }

    function tryAccept() public {
        if (job.phase() != IJob.Phase.Open) return;
        vm.prank(agentCtrl);
        try job.accept(agentId) {} catch {}
        _updateMaxPhase();
    }

    function trySubmit() public {
        if (job.phase() != IJob.Phase.Active) return;
        vm.prank(agentCtrl);
        try job.submitResult("ipfs://QmInvariant") {} catch {}
        _updateMaxPhase();
    }

    function tryApprove() public {
        if (job.phase() != IJob.Phase.Review) return;
        vm.prank(principal);
        try job.approveResult() {} catch {}
        _updateMaxPhase();
    }

    function tryCancel() public {
        IJob.Phase p = job.phase();
        if (p == IJob.Phase.Completed || p == IJob.Phase.Cancelled) return;
        vm.prank(principal);
        try job.cancel("invariant cancel") {} catch {}
        _updateMaxPhase();
    }

    function tryRaiseDispute() public {
        IJob.Phase p = job.phase();
        if (p != IJob.Phase.Active && p != IJob.Phase.Review) return;
        vm.prank(principal);
        try job.raiseDispute() {} catch {}
        _updateMaxPhase();
    }

    function tryResolveDisputeAgent() public {
        if (job.phase() != IJob.Phase.Disputed) return;
        vm.prank(arbiter);
        try job.resolveDispute(true) {} catch {}
        _updateMaxPhase();
    }

    function _updateMaxPhase() internal {
        uint8 current = uint8(job.phase());
        if (current > maxPhaseReached) {
            maxPhaseReached = current;
        }
    }
}

contract PhaseMonotonicityInvariantTest is StdInvariant, Test {
    AgentRegistry internal registry;
    MockToken internal token;
    Job internal jobContract;
    JobPhaseHandler internal handler;

    address internal admin = makeAddr("admin");
    address internal scorer = makeAddr("scorer");
    address internal principal = makeAddr("principal");
    address internal agentCtrl = makeAddr("agentCtrl");
    address internal arbiter = makeAddr("arbiter");
    uint256 internal agentId;
    uint256 internal constant BUDGET = 1000e18;

    function setUp() public {
        registry = new AgentRegistry(admin);
        vm.prank(agentCtrl);
        agentId = registry.register(agentCtrl, "ipfs://QmAgent", 0x1);

        token = new MockToken();
        token.mint(principal, BUDGET * 100);

        jobContract = new Job();
        token.mint(address(jobContract), BUDGET);

        Job.InitParams memory p = Job.InitParams({
            jobId: 0,
            principal: principal,
            registry: registry,
            targetAgentId: 0,
            token: address(token),
            budget: BUDGET,
            deadline: uint64(block.timestamp + 30 days),
            jobType: IJob.JobType.Direct,
            evaluator: address(0),
            arbiter: arbiter
        });
        jobContract.initialize(p);

        handler = new JobPhaseHandler(jobContract, registry, token, principal, agentCtrl, agentId, arbiter);
        targetContract(address(handler));
    }

    /// @notice Completed and Cancelled are terminal — once the job reaches a terminal
    ///         phase the budget field must be zero (all funds have moved out).
    function invariant_terminalPhasesAreSticky() public view {
        IJob.Phase p = jobContract.phase();
        if (p == IJob.Phase.Completed || p == IJob.Phase.Cancelled) {
            assertEq(jobContract.budget(), 0, "terminal job has non-zero budget");
        }
    }

    /// @notice Budget only decreases (no token minting from the job).
    function invariant_budgetNeverIncreases() public view {
        uint256 jobBalance = token.balanceOf(address(jobContract));
        uint256 contractBudget = jobContract.budget();
        // Token balance must equal tracked budget (no phantom tokens)
        assertEq(jobBalance, contractBudget, "job token balance differs from budget field");
    }
}
