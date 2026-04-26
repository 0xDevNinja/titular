// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC20Burnable} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import {Job} from "../../../src/acp/Job.sol";
import {IJob} from "../../../src/acp/IJob.sol";
import {JobFactory} from "../../../src/acp/JobFactory.sol";
import {AgentRegistry} from "../../../src/acp/AgentRegistry.sol";
import {Escrow} from "../../../src/acp/Escrow.sol";
import {FeeSplitter} from "../../../src/acp/FeeSplitter.sol";
import {HookRegistry} from "../../../src/acp/HookRegistry.sol";
import {BuybackBurner} from "../../../src/acp/BuybackBurner.sol";
import {IPermit2} from "../../../src/acp/IPermit2.sol";

/// @notice End-to-end integration test for the full ACP v2 job lifecycle:
///         factory clone → fund escrow → accept → submit result → complete →
///         FeeSplitter splits (Schedule A) → BuybackBurner receives 2%.
///
///         This test covers issue #55 acceptance criteria:
///           - JobFactory creates clone and funds Escrow
///           - Job _completeJob releases from Escrow → FeeSplitter
///           - FeeSplitter distributes 95% agent, 3% treasury, 2% buyback
///           - BuybackBurner receives its 2% share
contract MockPaymentToken is ERC20 {
    constructor() ERC20("USDC", "USDC") {}

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

/// @dev Mock UniV2 router: simulates a swap by transferring in and minting TITU to recipient.
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
        // Simulate 1:1 swap for simplicity
        titu.mint(to, amountIn);
    }
}

/// @title JobLifecycleE2E
/// @notice Full end-to-end test: JobFactory clone → Escrow funding → lifecycle →
///         FeeSplitter Schedule A split → BuybackBurner receives 2%.
contract JobLifecycleE2ETest is Test {
    // Contracts
    AgentRegistry internal registry;
    MockPaymentToken internal token;
    MockTITU internal titu;
    MockPermit2 internal mockPermit2;
    MockUniRouter internal uniRouter;
    Job internal jobImpl;
    Escrow internal escrow;
    FeeSplitter internal feeSplitter;
    BuybackBurner internal buybackBurner;
    HookRegistry internal hookRegistry;
    JobFactory internal factory;

    // Actors
    address internal admin = makeAddr("admin");
    address internal arbiter = makeAddr("arbiter");
    address internal principal = makeAddr("principal");
    address internal agentCtrl = makeAddr("agentCtrl");
    address internal treasury = makeAddr("treasury");

    // Schedule A basis points
    uint256 internal constant BPS = 10_000;
    uint256 internal constant TREASURY_BPS = 300; // 3%
    uint256 internal constant BUYBACK_BPS = 200; // 2%
    uint256 internal constant AGENT_BPS = BPS - TREASURY_BPS - BUYBACK_BPS; // 95%

    uint256 internal agentId;
    uint256 internal constant BUDGET = 10_000e18; // 10,000 USDC
    uint64 internal constant DEADLINE = 7 days;

    address[] internal noHooks;

    function setUp() public {
        // 1. Registry
        registry = new AgentRegistry(admin);
        vm.prank(agentCtrl);
        agentId = registry.register(agentCtrl, "ipfs://QmAgent", 0x1);

        // 2. Tokens
        token = new MockPaymentToken();
        titu = new MockTITU();
        mockPermit2 = new MockPermit2();
        uniRouter = new MockUniRouter(titu);

        // 3. Escrow
        escrow = new Escrow(admin, address(mockPermit2));

        // 4. BuybackBurner
        address[] memory swapPath = new address[](2);
        swapPath[0] = address(token);
        swapPath[1] = address(titu);
        buybackBurner = new BuybackBurner(admin, address(uniRouter), address(token), address(titu), swapPath, 0, 300);

        // 5. FeeSplitter (references BuybackBurner)
        feeSplitter = new FeeSplitter(admin, treasury, address(buybackBurner));

        // 6. HookRegistry
        hookRegistry = new HookRegistry(admin);

        // 7. Job implementation + Factory
        jobImpl = new Job();
        factory = new JobFactory(admin, address(jobImpl), registry, arbiter, escrow, feeSplitter, hookRegistry);

        // 8. Wire: grant factory admin rights on Escrow + FeeSplitter
        bytes32 escrowAdmin = escrow.DEFAULT_ADMIN_ROLE();
        vm.prank(admin);
        escrow.grantRole(escrowAdmin, address(factory));
        bytes32 splitterAdmin = feeSplitter.DEFAULT_ADMIN_ROLE();
        vm.prank(admin);
        feeSplitter.grantRole(splitterAdmin, address(factory));

        // 9. Fund principal
        token.mint(principal, BUDGET * 5);
        vm.prank(principal);
        token.approve(address(factory), type(uint256).max);

        noHooks = new address[](0);
    }

    // ─────────────────────────────────────────────────────────────
    // E2E: factory clone → fund escrow → accept → submit → complete
    // → FeeSplitter splits → BuybackBurner receives 2%
    // ─────────────────────────────────────────────────────────────

    function test_e2e_fullLifecycle_scheduleA_split() public {
        // Step 1: Principal creates a job
        uint256 principalBalanceBefore = token.balanceOf(principal);
        vm.prank(principal);
        (uint256 jobId, address clone) = factory.createJob(
            agentId,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE),
            IJob.JobType.Direct,
            address(0),
            address(0),
            noHooks
        );

        // Budget deducted from principal, held in escrow
        assertEq(token.balanceOf(principal), principalBalanceBefore - BUDGET);
        assertEq(escrow.getBalance(principal, jobId, address(token)), BUDGET);

        Job j = Job(clone);
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Open));

        // Step 2: Agent accepts
        vm.prank(agentCtrl);
        j.accept(agentId);
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Active));

        // Step 3: Agent submits result
        vm.prank(agentCtrl);
        j.submitResult("ipfs://QmE2EResult");
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Review));

        // Step 4: Principal approves → _completeJob → Escrow.release → FeeSplitter.split
        uint256 agentBalanceBefore = token.balanceOf(agentCtrl);
        uint256 treasuryBalanceBefore = token.balanceOf(treasury);
        uint256 buybackBalanceBefore = token.balanceOf(address(buybackBurner));

        vm.prank(principal);
        j.approveResult();

        assertEq(uint8(j.phase()), uint8(IJob.Phase.Completed));
        assertEq(j.budget(), 0);
        assertEq(escrow.getBalance(principal, jobId, address(token)), 0);

        // Step 5: Verify Schedule A distribution
        uint256 expectedTreasury = (BUDGET * TREASURY_BPS) / BPS; // 300 USDC
        uint256 expectedBuyback = (BUDGET * BUYBACK_BPS) / BPS; // 200 USDC
        uint256 expectedAgent = BUDGET - expectedTreasury - expectedBuyback; // 9500 USDC

        assertEq(token.balanceOf(agentCtrl), agentBalanceBefore + expectedAgent, "agent did not receive 95%");
        assertEq(token.balanceOf(treasury), treasuryBalanceBefore + expectedTreasury, "treasury did not receive 3%");
        assertEq(
            token.balanceOf(address(buybackBurner)),
            buybackBalanceBefore + expectedBuyback,
            "buybackBurner did not receive 2%"
        );
    }

    // ─────────────────────────────────────────────────────────────
    // E2E: cancel refunds principal via Escrow
    // ─────────────────────────────────────────────────────────────

    function test_e2e_cancel_refundsPrincipalViaEscrow() public {
        uint256 principalBefore = token.balanceOf(principal);
        vm.prank(principal);
        (, address clone) = factory.createJob(
            0,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE),
            IJob.JobType.Direct,
            address(0),
            address(0),
            noHooks
        );

        // Budget in escrow
        assertEq(token.balanceOf(principal), principalBefore - BUDGET);

        vm.prank(principal);
        Job(clone).cancel("changed mind");

        // Full refund via Escrow.refund
        assertEq(token.balanceOf(principal), principalBefore, "principal not refunded");
        assertEq(token.balanceOf(agentCtrl), 0, "agent should receive nothing");
    }

    // ─────────────────────────────────────────────────────────────
    // E2E: dispute arbiter-resolved (agent wins) → FeeSplitter
    // ─────────────────────────────────────────────────────────────

    function test_e2e_dispute_agentWins_scheduleA() public {
        vm.prank(principal);
        (, address clone) = factory.createJob(
            agentId,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE),
            IJob.JobType.Direct,
            address(0),
            address(0),
            noHooks
        );
        Job j = Job(clone);

        vm.prank(agentCtrl);
        j.accept(agentId);
        vm.prank(principal);
        j.raiseDispute();
        assertEq(uint8(j.phase()), uint8(IJob.Phase.Disputed));

        uint256 agentBefore = token.balanceOf(agentCtrl);
        uint256 buybackBefore = token.balanceOf(address(buybackBurner));

        vm.prank(arbiter);
        j.resolveDispute(true); // agent wins

        assertEq(uint8(j.phase()), uint8(IJob.Phase.Completed));

        // Agent gets 95%, buyback gets 2%
        uint256 expectedAgent = (BUDGET * AGENT_BPS) / BPS;
        uint256 expectedBuyback = (BUDGET * BUYBACK_BPS) / BPS;
        assertGe(token.balanceOf(agentCtrl) - agentBefore, expectedAgent);
        assertEq(token.balanceOf(address(buybackBurner)) - buybackBefore, expectedBuyback);
    }

    // ─────────────────────────────────────────────────────────────
    // E2E: expiry refunds principal via Escrow
    // ─────────────────────────────────────────────────────────────

    function test_e2e_expiry_refundsPrincipal() public {
        uint256 principalBefore = token.balanceOf(principal);
        vm.prank(principal);
        (, address clone) = factory.createJob(
            0,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE),
            IJob.JobType.Direct,
            address(0),
            address(0),
            noHooks
        );

        vm.warp(block.timestamp + DEADLINE + 1);
        Job(clone).expireJob();

        assertEq(token.balanceOf(principal), principalBefore, "principal not refunded after expiry");
    }

    // ─────────────────────────────────────────────────────────────
    // E2E: multiple independent jobs, escrow isolation
    // ─────────────────────────────────────────────────────────────

    function test_e2e_multipleJobs_escrowIsolation() public {
        // Create two jobs
        vm.prank(principal);
        (uint256 id0, address c0) = factory.createJob(
            agentId,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE),
            IJob.JobType.Direct,
            address(0),
            address(0),
            noHooks
        );
        vm.prank(principal);
        (uint256 id1, address c1) = factory.createJob(
            agentId,
            address(token),
            BUDGET,
            uint64(block.timestamp + DEADLINE),
            IJob.JobType.Direct,
            address(0),
            address(0),
            noHooks
        );

        // Both have independent escrow balances
        assertEq(escrow.getBalance(principal, id0, address(token)), BUDGET);
        assertEq(escrow.getBalance(principal, id1, address(token)), BUDGET);

        // Complete job 0
        vm.prank(agentCtrl);
        Job(c0).accept(agentId);
        vm.prank(agentCtrl);
        Job(c0).submitResult("ipfs://QmA");
        vm.prank(principal);
        Job(c0).approveResult();

        // Job 1 escrow untouched
        assertEq(escrow.getBalance(principal, id0, address(token)), 0);
        assertEq(escrow.getBalance(principal, id1, address(token)), BUDGET);
        assertEq(uint8(Job(c1).phase()), uint8(IJob.Phase.Open));
    }

    // ─────────────────────────────────────────────────────────────
    // Fuzz: verify conservation — total budget == agent + treasury + buyback
    // ─────────────────────────────────────────────────────────────

    function testFuzz_e2e_conservation(uint256 budget) public {
        budget = bound(budget, 1000, type(uint128).max);

        token.mint(principal, budget);
        vm.prank(principal);
        token.approve(address(factory), budget);

        uint256 totalBefore =
            token.balanceOf(agentCtrl) + token.balanceOf(treasury) + token.balanceOf(address(buybackBurner));

        vm.prank(principal);
        (, address clone) = factory.createJob(
            agentId,
            address(token),
            budget,
            uint64(block.timestamp + DEADLINE),
            IJob.JobType.Direct,
            address(0),
            address(0),
            noHooks
        );

        vm.prank(agentCtrl);
        Job(clone).accept(agentId);
        vm.prank(agentCtrl);
        Job(clone).submitResult("ipfs://QmFuzz");
        vm.prank(principal);
        Job(clone).approveResult();

        uint256 totalAfter =
            token.balanceOf(agentCtrl) + token.balanceOf(treasury) + token.balanceOf(address(buybackBurner));

        // Conservation: all funds distributed; total increase == budget
        assertEq(totalAfter - totalBefore, budget, "fund conservation violated");
    }
}
