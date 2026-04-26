// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC20Burnable} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import {Clones} from "@openzeppelin/contracts/proxy/Clones.sol";
import {Job} from "../../src/acp/Job.sol";
import {IJob} from "../../src/acp/IJob.sol";
import {AgentRegistry} from "../../src/acp/AgentRegistry.sol";
import {Escrow} from "../../src/acp/Escrow.sol";
import {FeeSplitter} from "../../src/acp/FeeSplitter.sol";
import {HookRegistry} from "../../src/acp/HookRegistry.sol";
import {BuybackBurner} from "../../src/acp/BuybackBurner.sol";
import {IPermit2} from "../../src/acp/IPermit2.sol";

/// @dev Minimal ERC-20 mock for testing.
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

/// @dev Minimal Permit2 mock (not needed for Job tests, just for Escrow init).
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
        uint256,
        address[] calldata path,
        address to,
        uint256
    ) external {
        ERC20(path[0]).transferFrom(msg.sender, address(this), amountIn);
        titu.mint(to, amountIn); // 1:1 for test simplicity
    }
}

contract JobTest is Test {
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

    address internal admin = makeAddr("admin");
    address internal scorer = makeAddr("scorer");
    address internal principal = makeAddr("principal");
    address internal agentCtrl = makeAddr("agentCtrl");
    address internal evaluatorAddr = makeAddr("evaluator");
    address internal arbiterAddr = makeAddr("arbiter");
    address internal treasury = makeAddr("treasury");
    address internal stranger = makeAddr("stranger");

    uint256 internal agentId0;

    uint256 internal constant BUDGET = 1000e18;
    uint64 internal constant DEADLINE_DELTA = 7 days;
    uint256 internal constant JOB_ID = 1;

    // ─────────────────────────────────────────────────────────────
    // Helpers
    // ─────────────────────────────────────────────────────────────

    function _deployRegistry() internal {
        registry = new AgentRegistry(admin);
        bytes32 scorerRole = registry.SCORER_ROLE();
        vm.prank(admin);
        registry.grantRole(scorerRole, scorer);
        // register agent
        vm.prank(agentCtrl);
        agentId0 = registry.register(agentCtrl, "ipfs://QmAgent", 0x1);
    }

    bytes32 internal releaserRole;
    bytes32 internal callerRole;

    function _deploySupporting() internal {
        titu = new MockTITU();
        mockPermit2 = new MockPermit2();
        uniRouter = new MockUniRouter(titu);

        escrow = new Escrow(admin, address(mockPermit2));
        address[] memory swapPath = new address[](2);
        swapPath[0] = address(token);
        swapPath[1] = address(titu);
        buybackBurner = new BuybackBurner(admin, address(uniRouter), address(token), address(titu), swapPath, 0, 300);
        feeSplitter = new FeeSplitter(admin, treasury, address(buybackBurner));
        hookRegistry = new HookRegistry(admin);
        jobImpl = new Job();

        // Cache role bytes32 values (avoid consuming vm.prank with staticcalls)
        releaserRole = escrow.RELEASER_ROLE();
        callerRole = feeSplitter.CALLER_ROLE();
    }

    /// @dev Creates a Job clone, funds escrow, and initializes the clone.
    function _newJob(IJob.JobType jType, address _evaluator) internal returns (Job j) {
        j = Job(Clones.clone(address(jobImpl)));

        // Fund escrow for this job (principal sends tokens to escrow)
        token.mint(principal, BUDGET);
        vm.prank(principal);
        token.approve(address(escrow), BUDGET);
        vm.prank(principal);
        escrow.fund(principal, JOB_ID, address(token), BUDGET);

        // Grant the clone RELEASER_ROLE and FeeSplitter CALLER_ROLE
        vm.prank(admin);
        escrow.grantRole(releaserRole, address(j));
        vm.prank(admin);
        feeSplitter.grantRole(callerRole, address(j));

        address[] memory hooks = new address[](0);
        Job.InitParams memory p = Job.InitParams({
            jobId: JOB_ID,
            principal: principal,
            registry: registry,
            targetAgentId: 0,
            token: address(token),
            budget: BUDGET,
            deadline: uint64(block.timestamp + DEADLINE_DELTA),
            jobType: jType,
            evaluator: _evaluator,
            arbiter: arbiterAddr,
            escrow: escrow,
            feeSplitter: feeSplitter,
            hookRegistry: hookRegistry,
            hooks: hooks
        });
        j.initialize(p);
    }

    function setUp() public {
        _deployRegistry();
        token = new MockERC20();
        _deploySupporting();
    }

    // ─────────────────────────────────────────────────────────────
    // Initialize
    // ─────────────────────────────────────────────────────────────

    function test_initialize_direct() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        assertEq(j.principal(), principal);
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Open));
        assertEq(uint8(j.jobType()), uint8(IJob.JobType.Direct));
        assertEq(j.evaluator(), address(0));
        assertEq(j.budget(), BUDGET);
    }

    function test_initialize_evaluated() public {
        Job j = _newJob(IJob.JobType.Evaluated, evaluatorAddr);
        assertEq(j.evaluator(), evaluatorAddr);
    }

    function test_initialize_revert_disabledOnImpl() public {
        // The implementation itself should reject initialize()
        address[] memory hooks = new address[](0);
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
            arbiter: arbiterAddr,
            escrow: escrow,
            feeSplitter: feeSplitter,
            hookRegistry: hookRegistry,
            hooks: hooks
        });
        vm.expectRevert();
        jobImpl.initialize(p);
    }

    function test_initialize_revert_zeroPrincipal() public {
        Job j = Job(Clones.clone(address(jobImpl)));
        address[] memory hooks = new address[](0);
        Job.InitParams memory p = Job.InitParams({
            jobId: 1,
            principal: address(0),
            registry: registry,
            targetAgentId: 0,
            token: address(token),
            budget: BUDGET,
            deadline: uint64(block.timestamp + DEADLINE_DELTA),
            jobType: IJob.JobType.Direct,
            evaluator: address(0),
            arbiter: arbiterAddr,
            escrow: escrow,
            feeSplitter: feeSplitter,
            hookRegistry: hookRegistry,
            hooks: hooks
        });
        vm.expectRevert(Job.ZeroAddress.selector);
        j.initialize(p);
    }

    function test_initialize_revert_zeroToken() public {
        Job j = Job(Clones.clone(address(jobImpl)));
        address[] memory hooks = new address[](0);
        Job.InitParams memory p = Job.InitParams({
            jobId: 1,
            principal: principal,
            registry: registry,
            targetAgentId: 0,
            token: address(0),
            budget: BUDGET,
            deadline: uint64(block.timestamp + DEADLINE_DELTA),
            jobType: IJob.JobType.Direct,
            evaluator: address(0),
            arbiter: arbiterAddr,
            escrow: escrow,
            feeSplitter: feeSplitter,
            hookRegistry: hookRegistry,
            hooks: hooks
        });
        vm.expectRevert(Job.ZeroAddress.selector);
        j.initialize(p);
    }

    function test_initialize_revert_zeroBudget() public {
        Job j = Job(Clones.clone(address(jobImpl)));
        address[] memory hooks = new address[](0);
        Job.InitParams memory p = Job.InitParams({
            jobId: 1,
            principal: principal,
            registry: registry,
            targetAgentId: 0,
            token: address(token),
            budget: 0,
            deadline: uint64(block.timestamp + DEADLINE_DELTA),
            jobType: IJob.JobType.Direct,
            evaluator: address(0),
            arbiter: arbiterAddr,
            escrow: escrow,
            feeSplitter: feeSplitter,
            hookRegistry: hookRegistry,
            hooks: hooks
        });
        vm.expectRevert(Job.ZeroAmount.selector);
        j.initialize(p);
    }

    function test_initialize_revert_pastDeadline() public {
        Job j = Job(Clones.clone(address(jobImpl)));
        address[] memory hooks = new address[](0);
        Job.InitParams memory p = Job.InitParams({
            jobId: 1,
            principal: principal,
            registry: registry,
            targetAgentId: 0,
            token: address(token),
            budget: BUDGET,
            deadline: uint64(block.timestamp - 1),
            jobType: IJob.JobType.Direct,
            evaluator: address(0),
            arbiter: arbiterAddr,
            escrow: escrow,
            feeSplitter: feeSplitter,
            hookRegistry: hookRegistry,
            hooks: hooks
        });
        vm.expectRevert(Job.InvalidDeadline.selector);
        j.initialize(p);
    }

    function test_initialize_revert_evaluatedWithoutEvaluator() public {
        Job j = Job(Clones.clone(address(jobImpl)));
        address[] memory hooks = new address[](0);
        Job.InitParams memory p = Job.InitParams({
            jobId: 1,
            principal: principal,
            registry: registry,
            targetAgentId: 0,
            token: address(token),
            budget: BUDGET,
            deadline: uint64(block.timestamp + DEADLINE_DELTA),
            jobType: IJob.JobType.Evaluated,
            evaluator: address(0),
            arbiter: arbiterAddr,
            escrow: escrow,
            feeSplitter: feeSplitter,
            hookRegistry: hookRegistry,
            hooks: hooks
        });
        vm.expectRevert(Job.EvaluatorRequired.selector);
        j.initialize(p);
    }

    function test_initialize_revert_directWithEvaluator() public {
        Job j = Job(Clones.clone(address(jobImpl)));
        address[] memory hooks = new address[](0);
        Job.InitParams memory p = Job.InitParams({
            jobId: 1,
            principal: principal,
            registry: registry,
            targetAgentId: 0,
            token: address(token),
            budget: BUDGET,
            deadline: uint64(block.timestamp + DEADLINE_DELTA),
            jobType: IJob.JobType.Direct,
            evaluator: evaluatorAddr,
            arbiter: arbiterAddr,
            escrow: escrow,
            feeSplitter: feeSplitter,
            hookRegistry: hookRegistry,
            hooks: hooks
        });
        vm.expectRevert(Job.DirectJobNoEvaluator.selector);
        j.initialize(p);
    }

    function test_initialize_revert_doubleInit() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        address[] memory hooks = new address[](0);
        Job.InitParams memory p = Job.InitParams({
            jobId: 2,
            principal: principal,
            registry: registry,
            targetAgentId: 0,
            token: address(token),
            budget: BUDGET,
            deadline: uint64(block.timestamp + DEADLINE_DELTA),
            jobType: IJob.JobType.Direct,
            evaluator: address(0),
            arbiter: arbiterAddr,
            escrow: escrow,
            feeSplitter: feeSplitter,
            hookRegistry: hookRegistry,
            hooks: hooks
        });
        vm.expectRevert();
        j.initialize(p);
    }

    // ─────────────────────────────────────────────────────────────
    // Accept
    // ─────────────────────────────────────────────────────────────

    function test_accept_success() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        vm.prank(agentCtrl);
        j.accept(agentId0);
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Active));
        assertEq(j.agent(), agentCtrl);
        assertEq(j.agentId(), agentId0);
    }

    function test_accept_revert_wrongPhase() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        vm.prank(agentCtrl);
        j.accept(agentId0);
        // Already active
        vm.expectRevert(abi.encodeWithSelector(Job.WrongPhase.selector, IJob.Phase.Active, IJob.Phase.Open));
        vm.prank(agentCtrl);
        j.accept(agentId0);
    }

    function test_accept_revert_wrongAgent() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        vm.expectRevert(abi.encodeWithSelector(Job.AgentNotFound.selector, agentId0));
        vm.prank(stranger); // stranger doesn't control agentId0
        j.accept(agentId0);
    }

    function test_accept_revert_inactiveAgent() public {
        vm.prank(agentCtrl);
        registry.setActive(agentId0, false);
        Job j = _newJob(IJob.JobType.Direct, address(0));
        vm.expectRevert(abi.encodeWithSelector(Job.AgentInactive.selector, agentId0));
        vm.prank(agentCtrl);
        j.accept(agentId0);
    }

    function test_accept_revert_targetMismatch() public {
        // Create a second agent
        vm.prank(stranger);
        uint256 agentId1 = registry.register(stranger, "ipfs://QmOther", 0x2);

        // Fund escrow for this specific job
        token.mint(principal, BUDGET);
        vm.prank(principal);
        token.approve(address(escrow), BUDGET);
        vm.prank(principal);
        escrow.fund(principal, 99, address(token), BUDGET);

        // Job targeted at agentId1
        Job j = Job(Clones.clone(address(jobImpl)));
        vm.prank(admin);
        escrow.grantRole(releaserRole, address(j));
        vm.prank(admin);
        feeSplitter.grantRole(callerRole, address(j));
        address[] memory hooks = new address[](0);
        Job.InitParams memory p = Job.InitParams({
            jobId: 99,
            principal: principal,
            registry: registry,
            targetAgentId: agentId1,
            token: address(token),
            budget: BUDGET,
            deadline: uint64(block.timestamp + DEADLINE_DELTA),
            jobType: IJob.JobType.Direct,
            evaluator: address(0),
            arbiter: arbiterAddr,
            escrow: escrow,
            feeSplitter: feeSplitter,
            hookRegistry: hookRegistry,
            hooks: hooks
        });
        j.initialize(p);

        // agentCtrl controls agentId0, not agentId1
        vm.expectRevert(abi.encodeWithSelector(Job.AgentNotFound.selector, agentId0));
        vm.prank(agentCtrl);
        j.accept(agentId0);
    }

    // ─────────────────────────────────────────────────────────────
    // Direct job happy path: Open → Active → Review → Completed
    // ─────────────────────────────────────────────────────────────

    function test_directJob_happyPath() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        vm.prank(agentCtrl);
        j.accept(agentId0);

        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmResult");
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Review));

        // Principal approves directly
        vm.prank(principal);
        j.approveResult();
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Completed));
        assertEq(j.budget(), 0);

        // Agent receives 95% (Schedule A)
        uint256 expectedAgent = (BUDGET * 9500) / 10_000;
        assertGe(token.balanceOf(agentCtrl), expectedAgent);
    }

    function test_directJob_release_afterGrace() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        vm.prank(agentCtrl);
        j.accept(agentId0);
        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmResult");

        // Not yet past grace period
        vm.expectRevert(Job.GracePeriodActive.selector);
        j.release();

        // Advance past grace
        vm.warp(block.timestamp + j.RELEASE_GRACE() + 1);
        j.release(); // anyone can call
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Completed));
    }

    // ─────────────────────────────────────────────────────────────
    // Evaluated job happy path
    // ─────────────────────────────────────────────────────────────

    function test_evaluatedJob_happyPath() public {
        Job j = _newJob(IJob.JobType.Evaluated, evaluatorAddr);
        vm.prank(agentCtrl);
        j.accept(agentId0);
        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmResult");

        vm.prank(evaluatorAddr);
        j.approveResult();
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Completed));

        // Agent receives 95%
        assertGe(token.balanceOf(agentCtrl), (BUDGET * 9500) / 10_000);
    }

    function test_evaluatedJob_reject_resubmit() public {
        Job j = _newJob(IJob.JobType.Evaluated, evaluatorAddr);
        vm.prank(agentCtrl);
        j.accept(agentId0);
        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmBad");

        vm.prank(evaluatorAddr);
        j.rejectResult("not good enough");
        // Back to Active; agent can resubmit
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Active));
        assertEq(bytes(j.resultURI()).length, 0);

        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmGood");
        vm.prank(evaluatorAddr);
        j.approveResult();
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Completed));
    }

    function test_evaluatedJob_approveResult_revert_notEvaluator() public {
        Job j = _newJob(IJob.JobType.Evaluated, evaluatorAddr);
        vm.prank(agentCtrl);
        j.accept(agentId0);
        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmResult");

        vm.expectRevert(Job.NotEvaluator.selector);
        vm.prank(principal);
        j.approveResult();
    }

    function test_directJob_approveResult_revert_notPrincipal() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        vm.prank(agentCtrl);
        j.accept(agentId0);
        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmResult");

        vm.expectRevert(Job.NotPrincipal.selector);
        vm.prank(stranger);
        j.approveResult();
    }

    // ─────────────────────────────────────────────────────────────
    // Cancel
    // ─────────────────────────────────────────────────────────────

    function test_cancel_fromOpen() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        uint256 principalBefore = token.balanceOf(principal);
        vm.prank(principal);
        j.cancel("changed mind");
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Cancelled));
        // Principal is refunded via escrow
        assertEq(token.balanceOf(principal), principalBefore + BUDGET);
    }

    function test_cancel_fromActive() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        vm.prank(agentCtrl);
        j.accept(agentId0);
        uint256 principalBefore = token.balanceOf(principal);
        vm.prank(principal);
        j.cancel("urgent");
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Cancelled));
        assertEq(token.balanceOf(principal), principalBefore + BUDGET);
    }

    function test_cancel_revert_notPrincipal() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        vm.expectRevert(Job.NotPrincipal.selector);
        vm.prank(stranger);
        j.cancel("nope");
    }

    function test_cancel_revert_fromCompleted() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        vm.prank(agentCtrl);
        j.accept(agentId0);
        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmResult");
        vm.prank(principal);
        j.approveResult();

        vm.expectRevert();
        vm.prank(principal);
        j.cancel("too late");
    }

    // ─────────────────────────────────────────────────────────────
    // Dispute
    // ─────────────────────────────────────────────────────────────

    function test_dispute_agentFavoured() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        vm.prank(agentCtrl);
        j.accept(agentId0);

        vm.prank(principal);
        j.raiseDispute();
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Disputed));

        vm.prank(arbiterAddr);
        j.resolveDispute(true);
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Completed));
        // Agent receives 95%
        assertGe(token.balanceOf(agentCtrl), (BUDGET * 9500) / 10_000);
    }

    function test_dispute_principalFavoured() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        vm.prank(agentCtrl);
        j.accept(agentId0);

        uint256 principalBefore = token.balanceOf(principal);
        vm.prank(agentCtrl);
        j.raiseDispute();

        vm.prank(arbiterAddr);
        j.resolveDispute(false);
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Cancelled));
        assertEq(token.balanceOf(principal), principalBefore + BUDGET);
    }

    function test_dispute_revert_notArbiter() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        vm.prank(agentCtrl);
        j.accept(agentId0);
        vm.prank(principal);
        j.raiseDispute();

        vm.expectRevert(Job.NotArbiter.selector);
        vm.prank(stranger);
        j.resolveDispute(true);
    }

    function test_raiseDispute_revert_wrongPhase() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        // Open phase — neither principal nor agent (no agent yet anyway)
        vm.expectRevert();
        vm.prank(principal);
        j.raiseDispute();
    }

    function test_raiseDispute_fromReview() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        vm.prank(agentCtrl);
        j.accept(agentId0);
        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmResult");

        vm.prank(agentCtrl);
        j.raiseDispute();
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Disputed));
    }

    // ─────────────────────────────────────────────────────────────
    // Expiry
    // ─────────────────────────────────────────────────────────────

    function test_expireJob_fromOpen() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        uint256 principalBefore = token.balanceOf(principal);
        vm.warp(block.timestamp + DEADLINE_DELTA + 1);
        j.expireJob();
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Cancelled));
        assertEq(token.balanceOf(principal), principalBefore + BUDGET);
    }

    function test_expireJob_fromActive() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        vm.prank(agentCtrl);
        j.accept(agentId0);
        vm.warp(block.timestamp + DEADLINE_DELTA + 1);
        j.expireJob();
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Cancelled));
    }

    function test_expireJob_revert_notExpired() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        vm.expectRevert(Job.JobNotExpired.selector);
        j.expireJob();
    }

    function test_expireJob_revert_fromCompleted() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        vm.prank(agentCtrl);
        j.accept(agentId0);
        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmResult");
        vm.prank(principal);
        j.approveResult();

        vm.warp(block.timestamp + DEADLINE_DELTA + 1);
        vm.expectRevert();
        j.expireJob();
    }

    // ─────────────────────────────────────────────────────────────
    // submitResult reverts
    // ─────────────────────────────────────────────────────────────

    function test_submitResult_revert_notAgent() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        vm.prank(agentCtrl);
        j.accept(agentId0);
        vm.expectRevert(Job.NotAgent.selector);
        vm.prank(stranger);
        j.submitResult("ipfs://QmResult");
    }

    function test_submitResult_revert_wrongPhase() public {
        Job j = _newJob(IJob.JobType.Direct, address(0));
        // Still Open — agent field is unset so modifier fires first
        vm.expectRevert();
        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmResult");
    }

    // ─────────────────────────────────────────────────────────────
    // Fuzz tests
    // ─────────────────────────────────────────────────────────────

    function testFuzz_directJob_fullLifecycle(uint256 budgetAmount) public {
        budgetAmount = bound(budgetAmount, 1, type(uint128).max);

        Job j = Job(Clones.clone(address(jobImpl)));

        token.mint(principal, budgetAmount);
        vm.prank(principal);
        token.approve(address(escrow), budgetAmount);
        vm.prank(principal);
        escrow.fund(principal, 99, address(token), budgetAmount);

        vm.prank(admin);
        escrow.grantRole(releaserRole, address(j));
        vm.prank(admin);
        feeSplitter.grantRole(callerRole, address(j));

        address[] memory hooks = new address[](0);
        Job.InitParams memory p = Job.InitParams({
            jobId: 99,
            principal: principal,
            registry: registry,
            targetAgentId: 0,
            token: address(token),
            budget: budgetAmount,
            deadline: uint64(block.timestamp + 1 days),
            jobType: IJob.JobType.Direct,
            evaluator: address(0),
            arbiter: arbiterAddr,
            escrow: escrow,
            feeSplitter: feeSplitter,
            hookRegistry: hookRegistry,
            hooks: hooks
        });
        j.initialize(p);

        vm.prank(agentCtrl);
        j.accept(agentId0);
        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmFuzz");
        vm.prank(principal);
        j.approveResult();

        // Agent gets >= 95%
        assertGe(token.balanceOf(agentCtrl), (budgetAmount * 9500) / 10_000);
        assertEq(j.budget(), 0);
    }

    function testFuzz_cancel_refundsPrincipal(uint256 budgetAmount) public {
        budgetAmount = bound(budgetAmount, 1, type(uint128).max);

        Job j = Job(Clones.clone(address(jobImpl)));

        token.mint(principal, budgetAmount);
        vm.prank(principal);
        token.approve(address(escrow), budgetAmount);
        vm.prank(principal);
        escrow.fund(principal, 98, address(token), budgetAmount);

        vm.prank(admin);
        escrow.grantRole(releaserRole, address(j));
        vm.prank(admin);
        feeSplitter.grantRole(callerRole, address(j));

        address[] memory hooks = new address[](0);
        Job.InitParams memory p = Job.InitParams({
            jobId: 98,
            principal: principal,
            registry: registry,
            targetAgentId: 0,
            token: address(token),
            budget: budgetAmount,
            deadline: uint64(block.timestamp + 1 days),
            jobType: IJob.JobType.Direct,
            evaluator: address(0),
            arbiter: arbiterAddr,
            escrow: escrow,
            feeSplitter: feeSplitter,
            hookRegistry: hookRegistry,
            hooks: hooks
        });
        j.initialize(p);

        uint256 principalBefore = token.balanceOf(principal);
        vm.prank(principal);
        j.cancel("fuzz cancel");
        assertEq(token.balanceOf(principal), principalBefore + budgetAmount);
    }
}
