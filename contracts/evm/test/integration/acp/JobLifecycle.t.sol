// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {Job} from "../../../src/acp/Job.sol";
import {IJob} from "../../../src/acp/IJob.sol";
import {JobFactory} from "../../../src/acp/JobFactory.sol";
import {AgentRegistry} from "../../../src/acp/AgentRegistry.sol";

/// @notice Full 4-phase lifecycle integration test covering both job types
///         (Direct and Evaluated), with and without an evaluator.
contract MockToken is ERC20 {
    constructor() ERC20("IntTest", "INT") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

contract JobLifecycleTest is Test {
    AgentRegistry internal registry;
    MockToken internal token;
    Job internal jobImpl;
    JobFactory internal factory;

    address internal admin = makeAddr("admin");
    address internal arbiter = makeAddr("arbiter");
    address internal principal = makeAddr("principal");
    address internal agentCtrl = makeAddr("agentCtrl");
    address internal evaluator = makeAddr("evaluator");

    uint256 internal agentId;
    uint256 internal constant BUDGET = 2000e18;
    uint64 internal constant DEADLINE = 7 days;

    function setUp() public {
        registry = new AgentRegistry(admin);
        vm.prank(agentCtrl);
        agentId = registry.register(agentCtrl, "ipfs://QmAgent", 0x1);

        token = new MockToken();
        jobImpl = new Job();
        factory = new JobFactory(admin, address(jobImpl), registry, arbiter);

        token.mint(principal, BUDGET * 20);
        vm.prank(principal);
        token.approve(address(factory), type(uint256).max);
    }

    // ─── Direct Job: full Open → Active → Review → Completed ────

    function test_directJob_fullLifecycle_principalApproves() public {
        vm.prank(principal);
        (uint256 jobId, address clone) = factory.createJob(
            agentId,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE),
            IJob.JobType.Direct,
            address(0),
            address(0)
        );
        Job j = Job(clone);

        // Phase 1: Accept
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Open));
        vm.prank(agentCtrl);
        j.accept(agentId);
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Active));
        assertEq(j.agent(), agentCtrl);
        assertEq(j.agentId(), agentId);

        // Phase 2: Submit
        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmResult1");
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Review));
        assertEq(j.resultURI(), "ipfs://QmResult1");

        // Phase 3: Approve (principal)
        vm.prank(principal);
        j.approveResult();
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Completed));
        assertEq(j.budget(), 0);
        assertEq(token.balanceOf(agentCtrl), BUDGET);

        assertEq(jobId, 0);
    }

    function test_directJob_release_afterGracePeriod() public {
        vm.prank(principal);
        (, address clone) = factory.createJob(
            0, address(token), BUDGET, uint64(block.timestamp + DEADLINE), IJob.JobType.Direct, address(0), address(0)
        );
        Job j = Job(clone);

        vm.prank(agentCtrl);
        j.accept(agentId);
        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmResult");

        // Before grace period — cannot release
        vm.expectRevert(Job.GracePeriodActive.selector);
        j.release();

        // After grace period — anyone can release
        vm.warp(block.timestamp + j.RELEASE_GRACE() + 1);
        address stranger = makeAddr("stranger");
        vm.prank(stranger);
        j.release();
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Completed));
    }

    // ─── Evaluated Job: Open → Active → Review → Completed ──────

    function test_evaluatedJob_fullLifecycle_evaluatorApproves() public {
        vm.prank(principal);
        (, address clone) = factory.createJob(
            agentId,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE),
            IJob.JobType.Evaluated,
            evaluator,
            address(0)
        );
        Job j = Job(clone);

        vm.prank(agentCtrl);
        j.accept(agentId);
        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmEval");
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Review));

        // Evaluator approves
        vm.prank(evaluator);
        j.approveResult();
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Completed));
        assertEq(token.balanceOf(agentCtrl), BUDGET);
    }

    function test_evaluatedJob_rejectAndResubmit() public {
        vm.prank(principal);
        (, address clone) = factory.createJob(
            agentId,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE),
            IJob.JobType.Evaluated,
            evaluator,
            address(0)
        );
        Job j = Job(clone);

        vm.prank(agentCtrl);
        j.accept(agentId);
        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmBad");

        // Evaluator rejects → back to Active
        vm.prank(evaluator);
        j.rejectResult("needs improvement");
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Active));
        assertEq(bytes(j.resultURI()).length, 0);

        // Agent resubmits
        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmGood");
        vm.prank(evaluator);
        j.approveResult();
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Completed));
    }

    // ─── Cancellation path ───────────────────────────────────────

    function test_cancel_fromOpen_refundsPrincipal() public {
        uint256 before = token.balanceOf(principal);
        vm.prank(principal);
        (, address clone) = factory.createJob(
            0, address(token), BUDGET, uint64(block.timestamp + DEADLINE), IJob.JobType.Direct, address(0), address(0)
        );
        Job j = Job(clone);

        vm.prank(principal);
        j.cancel("changed mind");
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Cancelled));
        assertEq(token.balanceOf(principal), before);
    }

    // ─── Dispute path ────────────────────────────────────────────

    function test_dispute_agentFavoured() public {
        vm.prank(principal);
        (, address clone) = factory.createJob(
            agentId,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE),
            IJob.JobType.Direct,
            address(0),
            address(0)
        );
        Job j = Job(clone);

        vm.prank(agentCtrl);
        j.accept(agentId);
        vm.prank(principal);
        j.raiseDispute();
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Disputed));

        vm.prank(arbiter);
        j.resolveDispute(true);
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Completed));
        assertEq(token.balanceOf(agentCtrl), BUDGET);
    }

    function test_dispute_principalFavoured_refunds() public {
        vm.prank(principal);
        (, address clone) = factory.createJob(
            agentId,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE),
            IJob.JobType.Direct,
            address(0),
            address(0)
        );
        Job j = Job(clone);

        vm.prank(agentCtrl);
        j.accept(agentId);
        vm.prank(agentCtrl);
        j.raiseDispute();

        uint256 before = token.balanceOf(principal);
        vm.prank(arbiter);
        j.resolveDispute(false);
        assertEq(token.balanceOf(principal), before + BUDGET);
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Cancelled));
    }

    // ─── Expiry path ─────────────────────────────────────────────

    function test_expiry_fromOpen() public {
        uint256 before = token.balanceOf(principal);
        vm.prank(principal);
        (, address clone) = factory.createJob(
            0, address(token), BUDGET, uint64(block.timestamp + DEADLINE), IJob.JobType.Direct, address(0), address(0)
        );
        Job j = Job(clone);

        vm.warp(block.timestamp + DEADLINE + 1);
        j.expireJob();
        assertEq(token.balanceOf(principal), before);
    }

    // ─── Multiple factories and clone isolation ──────────────────

    function test_multipleJobs_isolatedState() public {
        vm.prank(principal);
        (, address c0) = factory.createJob(
            agentId,
            address(token),
            BUDGET / 2,
            uint64(block.timestamp + DEADLINE),
            IJob.JobType.Direct,
            address(0),
            address(0)
        );
        vm.prank(principal);
        (, address c1) = factory.createJob(
            agentId,
            address(token),
            BUDGET / 2,
            uint64(block.timestamp + DEADLINE),
            IJob.JobType.Direct,
            address(0),
            address(0)
        );

        // Complete job 0
        vm.prank(agentCtrl);
        Job(c0).accept(agentId);
        vm.prank(agentCtrl);
        Job(c0).submitResult("ipfs://QmA");
        vm.prank(principal);
        Job(c0).approveResult();

        // Job 1 still Open
        assertEq(uint8(Job(c0).phase()), uint8(IJob.Phase.Completed));
        assertEq(uint8(Job(c1).phase()), uint8(IJob.Phase.Open));
    }
}
