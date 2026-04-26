// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC20Burnable} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import {JobFactory} from "../../src/acp/JobFactory.sol";
import {Job} from "../../src/acp/Job.sol";
import {IJob} from "../../src/acp/IJob.sol";
import {AgentRegistry} from "../../src/acp/AgentRegistry.sol";
import {Escrow} from "../../src/acp/Escrow.sol";
import {FeeSplitter} from "../../src/acp/FeeSplitter.sol";
import {HookRegistry} from "../../src/acp/HookRegistry.sol";
import {IPermit2} from "../../src/acp/IPermit2.sol";
import {BuybackBurner} from "../../src/acp/BuybackBurner.sol";

contract MockERC20 is ERC20 {
    constructor() ERC20("Mock", "MCK") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

contract MockTITU is ERC20Burnable {
    constructor() ERC20("TITU", "TITU") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

/// @dev Minimal Permit2 mock (not actually used by factory tests, just for Escrow init).
contract MockPermit2 {
    function permitTransferFrom(
        IPermit2.PermitTransferFrom calldata permit,
        IPermit2.SignatureTransferDetails calldata transferDetails,
        address owner,
        bytes calldata
    ) external {
        ERC20(permit.permitted.token).transferFrom(owner, transferDetails.to, transferDetails.requestedAmount);
    }

    function permitWitnessTransferFrom(
        IPermit2.PermitTransferFrom calldata permit,
        IPermit2.SignatureTransferDetails calldata transferDetails,
        address owner,
        bytes32,
        string calldata,
        bytes calldata
    ) external {
        ERC20(permit.permitted.token).transferFrom(owner, transferDetails.to, transferDetails.requestedAmount);
    }
}

/// @dev Mock UniV2 router: mints TITU to recipient to simulate a swap.
contract MockUniRouter {
    MockTITU public titu;

    constructor(MockTITU _titu) {
        titu = _titu;
    }

    function swapExactTokensForTokensSupportingFeeOnTransferTokens(
        uint256 amountIn,
        uint256 amountOutMin,
        address[] calldata path,
        address to,
        uint256
    ) external {
        ERC20(path[0]).transferFrom(msg.sender, address(this), amountIn);
        require(amountOutMin == 0 || amountIn >= amountOutMin, "slippage");
        titu.mint(to, amountIn); // 1:1 for test simplicity
    }
}

contract JobFactoryTest is Test {
    AgentRegistry internal registry;
    MockERC20 internal token;
    MockTITU internal titu;
    MockPermit2 internal mockPermit2;
    MockUniRouter internal uniRouter;
    Job internal jobImpl;
    Escrow internal escrow;
    FeeSplitter internal feeSplitter;
    BuybackBurner internal buybackBurner;
    HookRegistry internal hookRegistry;
    JobFactory internal factory;

    address internal admin = makeAddr("admin");
    address internal arbiter = makeAddr("arbiter");
    address internal principal = makeAddr("principal");
    address internal agentCtrl = makeAddr("agentCtrl");
    address internal evaluator = makeAddr("evaluator");
    address internal treasury = makeAddr("treasury");
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

    /// @dev Returns an empty hooks array for convenience.
    function _noHooks() internal pure returns (address[] memory h) {
        h = new address[](0);
    }

    function setUp() public {
        // Deploy registry + agent
        registry = new AgentRegistry(admin);
        vm.prank(agentCtrl);
        agentId0 = registry.register(agentCtrl, "ipfs://QmAgent", 0x1);

        token = new MockERC20();
        titu = new MockTITU();
        mockPermit2 = new MockPermit2();
        uniRouter = new MockUniRouter(titu);

        // Deploy supporting contracts
        escrow = new Escrow(admin, address(mockPermit2));
        address[] memory swapPath = new address[](2);
        swapPath[0] = address(token);
        swapPath[1] = address(titu);
        buybackBurner = new BuybackBurner(admin, address(uniRouter), address(token), address(titu), swapPath, 0, 300);
        feeSplitter = new FeeSplitter(admin, treasury, address(buybackBurner));
        hookRegistry = new HookRegistry(admin);

        jobImpl = new Job();
        factory = new JobFactory(admin, address(jobImpl), registry, arbiter, escrow, feeSplitter, hookRegistry);

        // Grant factory DEFAULT_ADMIN_ROLE on Escrow + FeeSplitter (so it can grant roles to clones)
        bytes32 adminRole = escrow.DEFAULT_ADMIN_ROLE();
        vm.prank(admin);
        escrow.grantRole(adminRole, address(factory));
        bytes32 splitterAdminRole = feeSplitter.DEFAULT_ADMIN_ROLE();
        vm.prank(admin);
        feeSplitter.grantRole(splitterAdminRole, address(factory));

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
        assertEq(address(factory.escrow()), address(escrow));
        assertEq(address(factory.feeSplitter()), address(feeSplitter));
        assertEq(address(factory.hookRegistry()), address(hookRegistry));
    }

    function test_constructor_revert_zeroAdmin() public {
        vm.expectRevert(JobFactory.ZeroAddress.selector);
        new JobFactory(address(0), address(jobImpl), registry, arbiter, escrow, feeSplitter, hookRegistry);
    }

    function test_constructor_revert_zeroImpl() public {
        vm.expectRevert(JobFactory.ZeroAddress.selector);
        new JobFactory(admin, address(0), registry, arbiter, escrow, feeSplitter, hookRegistry);
    }

    function test_constructor_revert_zeroArbiter() public {
        vm.expectRevert(JobFactory.ZeroAddress.selector);
        new JobFactory(admin, address(jobImpl), registry, address(0), escrow, feeSplitter, hookRegistry);
    }

    function test_constructor_revert_zeroEscrow() public {
        vm.expectRevert(JobFactory.ZeroAddress.selector);
        new JobFactory(admin, address(jobImpl), registry, arbiter, Escrow(address(0)), feeSplitter, hookRegistry);
    }

    function test_constructor_revert_zeroFeeSplitter() public {
        vm.expectRevert(JobFactory.ZeroAddress.selector);
        new JobFactory(admin, address(jobImpl), registry, arbiter, escrow, FeeSplitter(address(0)), hookRegistry);
    }

    function test_constructor_revert_zeroHookRegistry() public {
        vm.expectRevert(JobFactory.ZeroAddress.selector);
        new JobFactory(admin, address(jobImpl), registry, arbiter, escrow, feeSplitter, HookRegistry(address(0)));
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
            address(0),
            _noHooks()
        );

        assertEq(jobId, 0);
        assertNotEq(clone, address(0));
        assertEq(factory.getJob(jobId), clone);
        assertEq(factory.totalJobs(), 1);

        // Budget is now in escrow, not in the clone
        assertEq(escrow.getBalance(principal, jobId, address(token)), BUDGET);
        assertEq(token.balanceOf(clone), 0);

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
            address(0),
            _noHooks()
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
            address(0),
            _noHooks()
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
            customArbiter,
            _noHooks()
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
            address(0),
            _noHooks()
        );
        vm.prank(principal);
        (uint256 id1,) = factory.createJob(
            0,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE_DELTA),
            IJob.JobType.Direct,
            address(0),
            address(0),
            _noHooks()
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
            address(0),
            _noHooks()
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
            0,
            address(0),
            BUDGET,
            uint64(block.timestamp + DEADLINE_DELTA),
            IJob.JobType.Direct,
            address(0),
            address(0),
            _noHooks()
        );
    }

    function test_createJob_revert_zeroBudget() public {
        vm.expectRevert(JobFactory.ZeroAmount.selector);
        vm.prank(principal);
        factory.createJob(
            0,
            address(token),
            0,
            uint64(block.timestamp + DEADLINE_DELTA),
            IJob.JobType.Direct,
            address(0),
            address(0),
            _noHooks()
        );
    }

    function test_createJob_revert_pastDeadline() public {
        vm.expectRevert(JobFactory.InvalidDeadline.selector);
        vm.prank(principal);
        factory.createJob(
            0,
            address(token),
            BUDGET,
            uint64(block.timestamp - 1),
            IJob.JobType.Direct,
            address(0),
            address(0),
            _noHooks()
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
            address(0),
            _noHooks()
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
            address(0),
            _noHooks()
        );
        vm.prank(principal);
        (, address clone1) = factory.createJob(
            0,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE_DELTA),
            IJob.JobType.Direct,
            address(0),
            address(0),
            _noHooks()
        );
        assertNotEq(clone0, clone1);
        // State is isolated
        assertEq(Job(clone0).jobId(), 0);
        assertEq(Job(clone1).jobId(), 1);
    }

    function test_implementation_disablesInitializers() public {
        // The implementation itself must revert on initialize (disableInitializers in constructor)
        Job.InitParams memory p = Job.InitParams({
            jobId: 1,
            principal: principal,
            registry: registry,
            targetAgentId: 0,
            token: address(token),
            budget: BUDGET,
            deadline: uint64(block.timestamp + DEADLINE_DELTA),
            jobType: IJob.JobType.Direct,
            evaluator: address(0),
            arbiter: arbiter,
            escrow: escrow,
            feeSplitter: feeSplitter,
            hookRegistry: hookRegistry,
            hooks: _noHooks()
        });
        vm.expectRevert();
        jobImpl.initialize(p);
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
    // End-to-end: create → accept → complete (FeeSplitter receives funds)
    // ─────────────────────────────────────────────────────────────

    function test_e2e_createAndComplete() public {
        vm.prank(principal);
        (uint256 jobId, address clone) = factory.createJob(
            agentId0,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE_DELTA),
            IJob.JobType.Direct,
            address(0),
            address(0),
            _noHooks()
        );

        // Budget held in escrow
        assertEq(escrow.getBalance(principal, jobId, address(token)), BUDGET);

        Job j = Job(clone);
        vm.prank(agentCtrl);
        j.accept(agentId0);
        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmResult");
        vm.prank(principal);
        j.approveResult();

        assertEq(uint8(j.phase()), uint8(IJob.Phase.Completed));
        // 95% of BUDGET goes to agent (Schedule A)
        uint256 expectedAgent = (BUDGET * 9500) / 10_000;
        assertGe(token.balanceOf(agentCtrl), expectedAgent);
        // Escrow balance drained
        assertEq(escrow.getBalance(principal, jobId, address(token)), 0);
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
        (uint256 jobId,) = factory.createJob(
            0,
            address(token),
            budgetAmount,
            uint64(block.timestamp + DEADLINE_DELTA),
            IJob.JobType.Direct,
            address(0),
            address(0),
            _noHooks()
        );
        // Budget is in escrow now
        assertEq(escrow.getBalance(principal, jobId, address(token)), budgetAmount);
    }
}
