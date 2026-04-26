// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {JobFactory} from "../../src/acp/JobFactory.sol";
import {Job} from "../../src/acp/Job.sol";
import {IJob} from "../../src/acp/IJob.sol";
import {AgentRegistry} from "../../src/acp/AgentRegistry.sol";

contract MockERC20 is ERC20 {
    constructor() ERC20("Mock", "MCK") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

contract JobFactoryTest is Test {
    AgentRegistry internal registry;
    MockERC20 internal token;
    Job internal jobImpl;
    JobFactory internal factory;

    address internal admin = makeAddr("admin");
    address internal arbiter = makeAddr("arbiter");
    address internal principal = makeAddr("principal");
    address internal agentCtrl = makeAddr("agentCtrl");
    address internal evaluator = makeAddr("evaluator");
    address internal stranger = makeAddr("stranger");

    uint256 internal agentId0;

    uint256 internal constant BUDGET = 1000e18;
    uint64 internal constant DEADLINE_DELTA = 7 days;

    event JobCreated(
        uint256 indexed jobId,
        address indexed clone,
        address indexed principal,
        IJob.JobType jobType,
        address token,
        uint256 budget
    );

    event ImplementationUpdated(address indexed oldImpl, address indexed newImpl);
    event DefaultArbiterUpdated(address indexed oldArbiter, address indexed newArbiter);

    function setUp() public {
        // Deploy registry + agent
        registry = new AgentRegistry(admin);
        vm.prank(agentCtrl);
        agentId0 = registry.register(agentCtrl, "ipfs://QmAgent", 0x1);

        token = new MockERC20();
        jobImpl = new Job();
        factory = new JobFactory(admin, address(jobImpl), registry, arbiter);

        // Mint + approve tokens for principal
        token.mint(principal, BUDGET * 10);
        vm.prank(principal);
        token.approve(address(factory), type(uint256).max);
    }

    // ─────────────────────────────────────────────────────────────
    // Constructor
    // ─────────────────────────────────────────────────────────────

    function test_constructor_setsState() public view {
        assertEq(factory.jobImplementation(), address(jobImpl));
        assertEq(address(factory.registry()), address(registry));
        assertEq(factory.defaultArbiter(), arbiter);
    }

    function test_constructor_revert_zeroAdmin() public {
        vm.expectRevert(JobFactory.ZeroAddress.selector);
        new JobFactory(address(0), address(jobImpl), registry, arbiter);
    }

    function test_constructor_revert_zeroImpl() public {
        vm.expectRevert(JobFactory.ZeroAddress.selector);
        new JobFactory(admin, address(0), registry, arbiter);
    }

    function test_constructor_revert_zeroArbiter() public {
        vm.expectRevert(JobFactory.ZeroAddress.selector);
        new JobFactory(admin, address(jobImpl), registry, address(0));
    }

    // ─────────────────────────────────────────────────────────────
    // createJob — happy paths
    // ─────────────────────────────────────────────────────────────

    function test_createJob_direct() public {
        vm.expectEmit(false, false, true, true);
        emit JobCreated(0, address(0), principal, IJob.JobType.Direct, address(token), BUDGET);

        vm.prank(principal);
        (uint256 jobId, address clone) = factory.createJob(
            0,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE_DELTA),
            IJob.JobType.Direct,
            address(0),
            address(0)
        );

        assertEq(jobId, 0);
        assertNotEq(clone, address(0));
        assertEq(factory.getJob(jobId), clone);
        assertEq(factory.totalJobs(), 1);
        assertEq(token.balanceOf(clone), BUDGET);

        Job j = Job(clone);
        assertEq(j.principal(), principal);
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Open));
        assertEq(uint8(j.jobType()), uint8(IJob.JobType.Direct));
        assertEq(j.budget(), BUDGET);
    }

    function test_createJob_evaluated() public {
        vm.prank(principal);
        (uint256 jobId, address clone) = factory.createJob(
            0,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE_DELTA),
            IJob.JobType.Evaluated,
            evaluator,
            address(0)
        );
        Job j = Job(clone);
        assertEq(j.evaluator(), evaluator);
        assertEq(uint8(j.jobType()), uint8(IJob.JobType.Evaluated));
        assertEq(jobId, 0);
    }

    function test_createJob_usesDefaultArbiter() public {
        vm.prank(principal);
        (, address clone) = factory.createJob(
            0,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE_DELTA),
            IJob.JobType.Direct,
            address(0),
            address(0)
        );
        assertEq(Job(clone).arbiter(), arbiter);
    }

    function test_createJob_customArbiter() public {
        address customArbiter = makeAddr("customArbiter");
        vm.prank(principal);
        (, address clone) = factory.createJob(
            0,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE_DELTA),
            IJob.JobType.Direct,
            address(0),
            customArbiter
        );
        assertEq(Job(clone).arbiter(), customArbiter);
    }

    function test_createJob_multipleJobs_uniqueIds() public {
        vm.prank(principal);
        (uint256 id0,) = factory.createJob(
            0,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE_DELTA),
            IJob.JobType.Direct,
            address(0),
            address(0)
        );
        vm.prank(principal);
        (uint256 id1,) = factory.createJob(
            0,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE_DELTA),
            IJob.JobType.Direct,
            address(0),
            address(0)
        );
        assertEq(id0, 0);
        assertEq(id1, 1);
        assertEq(factory.totalJobs(), 2);
        assertNotEq(factory.getJob(0), factory.getJob(1));
    }

    function test_createJob_tracked_byPrincipal() public {
        vm.prank(principal);
        factory.createJob(
            0,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE_DELTA),
            IJob.JobType.Direct,
            address(0),
            address(0)
        );
        uint256[] memory ids = factory.jobsByPrincipal(principal);
        assertEq(ids.length, 1);
        assertEq(ids[0], 0);
    }

    // ─────────────────────────────────────────────────────────────
    // createJob — revert paths
    // ─────────────────────────────────────────────────────────────

    function test_createJob_revert_zeroToken() public {
        vm.expectRevert(JobFactory.ZeroAddress.selector);
        vm.prank(principal);
        factory.createJob(
            0, address(0), BUDGET, uint64(block.timestamp + DEADLINE_DELTA), IJob.JobType.Direct, address(0), address(0)
        );
    }

    function test_createJob_revert_zeroBudget() public {
        vm.expectRevert(JobFactory.ZeroAmount.selector);
        vm.prank(principal);
        factory.createJob(
            0, address(token), 0, uint64(block.timestamp + DEADLINE_DELTA), IJob.JobType.Direct, address(0), address(0)
        );
    }

    function test_createJob_revert_pastDeadline() public {
        vm.expectRevert(JobFactory.InvalidDeadline.selector);
        vm.prank(principal);
        factory.createJob(
            0, address(token), BUDGET, uint64(block.timestamp - 1), IJob.JobType.Direct, address(0), address(0)
        );
    }

    function test_createJob_revert_paused() public {
        vm.prank(admin);
        factory.pause();
        vm.expectRevert();
        vm.prank(principal);
        factory.createJob(
            0,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE_DELTA),
            IJob.JobType.Direct,
            address(0),
            address(0)
        );
    }

    // ─────────────────────────────────────────────────────────────
    // Clone isolation
    // ─────────────────────────────────────────────────────────────

    function test_cloneIsolation_differentAddresses() public {
        vm.prank(principal);
        (, address clone0) = factory.createJob(
            0,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE_DELTA),
            IJob.JobType.Direct,
            address(0),
            address(0)
        );
        vm.prank(principal);
        (, address clone1) = factory.createJob(
            0,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE_DELTA),
            IJob.JobType.Direct,
            address(0),
            address(0)
        );
        assertNotEq(clone0, clone1);
        // State is isolated
        assertEq(Job(clone0).jobId(), 0);
        assertEq(Job(clone1).jobId(), 1);
    }

    function test_implementation_notInitialized() public view {
        // The implementation contract itself should be uninitialised (version = 0 / default state)
        assertEq(Job(address(jobImpl)).budget(), 0);
        assertEq(Job(address(jobImpl)).principal(), address(0));
    }

    // ─────────────────────────────────────────────────────────────
    // Admin functions
    // ─────────────────────────────────────────────────────────────

    function test_setJobImplementation() public {
        Job newImpl = new Job();
        vm.expectEmit(true, true, false, false);
        emit ImplementationUpdated(address(jobImpl), address(newImpl));
        vm.prank(admin);
        factory.setJobImplementation(address(newImpl));
        assertEq(factory.jobImplementation(), address(newImpl));
    }

    function test_setJobImplementation_revert_zero() public {
        vm.expectRevert(JobFactory.ZeroAddress.selector);
        vm.prank(admin);
        factory.setJobImplementation(address(0));
    }

    function test_setJobImplementation_revert_notAdmin() public {
        vm.expectRevert();
        vm.prank(stranger);
        factory.setJobImplementation(address(jobImpl));
    }

    function test_setDefaultArbiter() public {
        address newArbiter = makeAddr("newArbiter");
        vm.expectEmit(true, true, false, false);
        emit DefaultArbiterUpdated(arbiter, newArbiter);
        vm.prank(admin);
        factory.setDefaultArbiter(newArbiter);
        assertEq(factory.defaultArbiter(), newArbiter);
    }

    function test_setDefaultArbiter_revert_zero() public {
        vm.expectRevert(JobFactory.ZeroAddress.selector);
        vm.prank(admin);
        factory.setDefaultArbiter(address(0));
    }

    function test_pause_unpause() public {
        vm.prank(admin);
        factory.pause();
        assertTrue(factory.paused());
        vm.prank(admin);
        factory.unpause();
        assertFalse(factory.paused());
    }

    function test_pause_revert_notPauser() public {
        vm.expectRevert();
        vm.prank(stranger);
        factory.pause();
    }

    // ─────────────────────────────────────────────────────────────
    // End-to-end: create → accept → complete
    // ─────────────────────────────────────────────────────────────

    function test_e2e_createAndComplete() public {
        vm.prank(principal);
        (, address clone) = factory.createJob(
            agentId0,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE_DELTA),
            IJob.JobType.Direct,
            address(0),
            address(0)
        );

        Job j = Job(clone);
        vm.prank(agentCtrl);
        j.accept(agentId0);
        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmResult");
        vm.prank(principal);
        j.approveResult();

        assertEq(uint8(j.phase()), uint8(IJob.Phase.Completed));
        assertEq(token.balanceOf(agentCtrl), BUDGET);
    }

    // ─────────────────────────────────────────────────────────────
    // Fuzz
    // ─────────────────────────────────────────────────────────────

    function testFuzz_createJob_amounts(uint256 budgetAmount) public {
        budgetAmount = bound(budgetAmount, 1, type(uint128).max);
        token.mint(principal, budgetAmount);
        vm.prank(principal);
        token.approve(address(factory), type(uint256).max);

        vm.prank(principal);
        (, address clone) = factory.createJob(
            0,
            address(token),
            budgetAmount,
            uint64(block.timestamp + DEADLINE_DELTA),
            IJob.JobType.Direct,
            address(0),
            address(0)
        );
        assertEq(token.balanceOf(clone), budgetAmount);
        assertEq(Job(clone).budget(), budgetAmount);
    }
}
