// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC20Burnable} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {Clones} from "@openzeppelin/contracts/proxy/Clones.sol";

import {IJob} from "../../src/acp/IJob.sol";
import {Job} from "../../src/acp/Job.sol";
import {AgentRegistry} from "../../src/acp/AgentRegistry.sol";
import {Escrow} from "../../src/acp/Escrow.sol";
import {FeeSplitter} from "../../src/acp/FeeSplitter.sol";
import {HookRegistry} from "../../src/acp/HookRegistry.sol";
import {BuybackBurner} from "../../src/acp/BuybackBurner.sol";
import {IPermit2} from "../../src/acp/IPermit2.sol";

/// @dev Tax-free ERC-20 used by the symbolic harness.
contract HalToken is ERC20 {
    constructor() ERC20("HAL", "HAL") {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

contract HalTITU is ERC20Burnable {
    constructor() ERC20("TITU", "TITU") {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

contract HalPermit2 {
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

contract HalUniRouter {
    HalTITU public titu;

    constructor(HalTITU _titu) {
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
        titu.mint(to, amountIn);
    }
}

/// @title JobLifecycleSymbolic
/// @notice Halmos symbolic execution harness for the Job state machine.
///
///         Run via:
///           halmos --match-contract JobLifecycle
contract JobLifecycleSymbolic is Test {
    Job internal jobImpl;
    Job internal job;
    AgentRegistry internal registry;
    HalToken internal token;
    HalTITU internal titu;
    HalPermit2 internal permit2;
    HalUniRouter internal uniRouter;
    Escrow internal escrow;
    FeeSplitter internal feeSplitter;
    BuybackBurner internal buybackBurner;
    HookRegistry internal hookRegistry;

    address internal constant PRINCIPAL = address(0xA1);
    address internal constant AGENT = address(0xA2);
    address internal constant ARBITER = address(0xA3);
    address internal constant TREASURY = address(0xA4);
    address internal constant ADMIN = address(0xA5);

    uint256 internal constant BUDGET = 1000e18;
    uint256 internal constant JOB_ID = 1;

    uint256 internal agentRegId;

    function setUp() public {
        registry = new AgentRegistry(ADMIN);
        vm.prank(AGENT);
        agentRegId = registry.register(AGENT, "ipfs://hal", 0);

        token = new HalToken();
        titu = new HalTITU();
        permit2 = new HalPermit2();
        uniRouter = new HalUniRouter(titu);

        escrow = new Escrow(ADMIN, address(permit2));
        address[] memory swapPath = new address[](2);
        swapPath[0] = address(token);
        swapPath[1] = address(titu);
        buybackBurner = new BuybackBurner(ADMIN, address(uniRouter), address(token), address(titu), swapPath, 0, 300);
        feeSplitter = new FeeSplitter(ADMIN, TREASURY, address(buybackBurner));
        hookRegistry = new HookRegistry(ADMIN);

        jobImpl = new Job();

        // Create a clone to initialize
        job = Job(Clones.clone(address(jobImpl)));

        // Fund escrow with budget
        token.mint(PRINCIPAL, BUDGET);
        vm.prank(PRINCIPAL);
        token.approve(address(escrow), BUDGET);
        vm.prank(PRINCIPAL);
        escrow.fund(PRINCIPAL, JOB_ID, address(token), BUDGET);

        // Grant clone RELEASER_ROLE and FeeSplitter CALLER_ROLE
        bytes32 releaserRole = escrow.RELEASER_ROLE();
        vm.prank(ADMIN);
        escrow.grantRole(releaserRole, address(job));
        bytes32 callerRole = feeSplitter.CALLER_ROLE();
        vm.prank(ADMIN);
        feeSplitter.grantRole(callerRole, address(job));

        address[] memory hooks = new address[](0);
        job.initialize(
            Job.InitParams({
                jobId: JOB_ID,
                principal: PRINCIPAL,
                registry: registry,
                targetAgentId: 0,
                token: address(token),
                budget: BUDGET,
                deadline: uint64(block.timestamp + 7 days),
                jobType: IJob.JobType.Direct,
                evaluator: address(0),
                arbiter: ARBITER,
                escrow: escrow,
                feeSplitter: feeSplitter,
                hookRegistry: hookRegistry,
                hooks: hooks
            })
        );
    }

    // ─────────────────────────────────────────────────────────────
    // Terminal absorption
    // ─────────────────────────────────────────────────────────────

    /// @notice Once `phase == Completed`, every state-mutating call MUST revert.
    function check_terminal_completed_isAbsorbing() public {
        vm.prank(AGENT);
        job.accept(agentRegId);
        vm.prank(AGENT);
        job.submitResult("ipfs://result");
        vm.prank(PRINCIPAL);
        job.approveResult();
        assertEq(uint8(job.phase()), uint8(IJob.Phase.Completed));

        _assertEveryTransitionReverts();
    }

    /// @notice Once `phase == Cancelled` via principal cancel, every transition reverts.
    function check_terminal_cancelled_byPrincipal_isAbsorbing() public {
        vm.prank(PRINCIPAL);
        job.cancel("changed mind");
        assertEq(uint8(job.phase()), uint8(IJob.Phase.Cancelled));

        _assertEveryTransitionReverts();
    }

    /// @notice Once `phase == Cancelled` via dispute-rejected, every transition reverts.
    function check_terminal_cancelled_byRejection_isAbsorbing() public {
        vm.prank(AGENT);
        job.accept(agentRegId);
        vm.prank(PRINCIPAL);
        job.raiseDispute();
        vm.prank(ARBITER);
        job.resolveDispute(false);
        assertEq(uint8(job.phase()), uint8(IJob.Phase.Cancelled));

        _assertEveryTransitionReverts();
    }

    /// @notice Once `phase == Cancelled` via expiry, every transition reverts.
    function check_terminal_cancelled_byExpiry_isAbsorbing() public {
        vm.warp(uint256(job.deadline()) + 1);
        job.expireJob();
        assertEq(uint8(job.phase()), uint8(IJob.Phase.Cancelled));

        _assertEveryTransitionReverts();
    }

    // ─────────────────────────────────────────────────────────────
    // Budget zeroing and fund conservation
    // ─────────────────────────────────────────────────────────────

    /// @notice After Completed, budget storage MUST be 0; agent received >=95%.
    function check_completed_zeroesBudget_paysAgent() public {
        uint256 origBudget = job.budget();
        vm.prank(AGENT);
        job.accept(agentRegId);
        vm.prank(AGENT);
        job.submitResult("ipfs://result");
        vm.prank(PRINCIPAL);
        job.approveResult();

        assertEq(job.budget(), 0);
        // Agent receives >= 95% (dust goes to agent, 3% treasury, 2% buyback)
        assertGe(token.balanceOf(AGENT), (origBudget * 9500) / 10_000);
        assertEq(escrow.getBalance(PRINCIPAL, JOB_ID, address(token)), 0);
    }

    /// @notice After Cancelled-by-rejection, budget MUST be 0 and principal refunded.
    function check_rejected_zeroesBudget_refundsPrincipal() public {
        uint256 origBudget = job.budget();
        vm.prank(AGENT);
        job.accept(agentRegId);
        vm.prank(PRINCIPAL);
        job.raiseDispute();
        vm.prank(ARBITER);
        job.resolveDispute(false);

        assertEq(job.budget(), 0);
        assertEq(token.balanceOf(PRINCIPAL), origBudget);
        assertEq(escrow.getBalance(PRINCIPAL, JOB_ID, address(token)), 0);
    }

    /// @notice After Cancelled-by-expiry, budget MUST be 0 and principal refunded.
    function check_expired_zeroesBudget_refundsPrincipal() public {
        uint256 origBudget = job.budget();
        vm.warp(uint256(job.deadline()) + 1);
        job.expireJob();

        assertEq(job.budget(), 0);
        assertEq(token.balanceOf(PRINCIPAL), origBudget);
        assertEq(escrow.getBalance(PRINCIPAL, JOB_ID, address(token)), 0);
    }

    // ─────────────────────────────────────────────────────────────
    // expireJob phase-gating
    // ─────────────────────────────────────────────────────────────

    /// @notice expireJob MUST revert from Review phase regardless of timestamp.
    function check_expire_revertsFromReview() public {
        vm.prank(AGENT);
        job.accept(agentRegId);
        vm.prank(AGENT);
        job.submitResult("ipfs://result");
        assertEq(uint8(job.phase()), uint8(IJob.Phase.Review));

        vm.warp(uint256(job.deadline()) + 1);
        bool reverted;
        try job.expireJob() {
            reverted = false;
        } catch {
            reverted = true;
        }
        assertTrue(reverted, "expireJob from Review must revert");
    }

    /// @notice expireJob MUST revert before deadline regardless of phase.
    function check_expire_revertsBeforeDeadline() public {
        bool reverted;
        try job.expireJob() {
            reverted = false;
        } catch {
            reverted = true;
        }
        assertTrue(reverted, "expireJob before deadline must revert");
    }

    // ─────────────────────────────────────────────────────────────
    // Internal helpers
    // ─────────────────────────────────────────────────────────────

    function _assertEveryTransitionReverts() internal {
        bool reverted;

        vm.prank(AGENT);
        try job.accept(agentRegId) {
            reverted = false;
        } catch {
            reverted = true;
        }
        assertTrue(reverted, "accept must revert at terminal");

        vm.prank(AGENT);
        try job.submitResult("ipfs://r") {
            reverted = false;
        } catch {
            reverted = true;
        }
        assertTrue(reverted, "submitResult must revert at terminal");

        vm.prank(PRINCIPAL);
        try job.approveResult() {
            reverted = false;
        } catch {
            reverted = true;
        }
        assertTrue(reverted, "approveResult must revert at terminal");

        try job.release() {
            reverted = false;
        } catch {
            reverted = true;
        }
        assertTrue(reverted, "release must revert at terminal");

        vm.prank(PRINCIPAL);
        try job.cancel("again") {
            reverted = false;
        } catch {
            reverted = true;
        }
        assertTrue(reverted, "cancel must revert at terminal");

        vm.prank(PRINCIPAL);
        try job.raiseDispute() {
            reverted = false;
        } catch {
            reverted = true;
        }
        assertTrue(reverted, "raiseDispute must revert at terminal");

        vm.prank(ARBITER);
        try job.resolveDispute(true) {
            reverted = false;
        } catch {
            reverted = true;
        }
        assertTrue(reverted, "resolveDispute must revert at terminal");

        vm.warp(uint256(job.deadline()) + 1);
        try job.expireJob() {
            reverted = false;
        } catch {
            reverted = true;
        }
        assertTrue(reverted, "expireJob must revert at terminal");
    }
}
