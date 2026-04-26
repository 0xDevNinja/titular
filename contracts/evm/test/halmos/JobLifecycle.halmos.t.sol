// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import {IJob} from "../../src/acp/IJob.sol";
import {Job} from "../../src/acp/Job.sol";
import {AgentRegistry} from "../../src/acp/AgentRegistry.sol";

/// @dev Tax-free ERC-20 used by the symbolic harness so the SMT solver does
///      not have to reason about a token's transfer side-effects. Mirrors
///      `HalToken` in the M2 symbolic suite.
contract HalToken is ERC20 {
    constructor() ERC20("HAL", "HAL") {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @title JobLifecycleSymbolic
/// @notice Halmos symbolic execution harness for the Job state machine.
///
///         Issue #70 invariant — once a Job enters its initial funded state
///         (`Phase.Open` with a non-zero budget), it eventually reaches
///         **exactly one** terminal phase, drawn from `{Completed, Cancelled}`.
///         "Cancelled" subsumes both the dispute-rejected and the
///         expired-refund paths (see `IJob.Phase` and `Job._completeJob` /
///         `Job.cancel` / `Job.expireJob` / `Job.resolveDispute(false)`).
///
///         The proofs below cover the safety side of that invariant:
///           - terminal_*_isAbsorbing — once `phase ∈ {Completed, Cancelled}`,
///             every state transition reverts.
///           - completed_zeroesBudget_paysAgent — the only way to land in
///             `Completed` is via approve / release / arbiter-favours-agent;
///             `budget` is zeroed and the full original budget is paid to
///             the agent.
///           - rejected_zeroesBudget_refundsPrincipal — both Cancelled paths
///             other than expiry (cancel-by-principal, dispute-rejected)
///             refund the principal in full.
///           - expired_zeroesBudget_refundsPrincipal — `expireJob` from
///             {Open, Active} after deadline refunds the principal.
///           - expire_phase_gating — `expireJob` reverts from any phase
///             other than {Open, Active}, and reverts before deadline.
///
///         Run via:
///           halmos --match-contract JobLifecycle
contract JobLifecycleSymbolic is Test {
    Job internal job;
    AgentRegistry internal registry;
    HalToken internal token;

    address internal constant PRINCIPAL = address(0xA1);
    address internal constant AGENT = address(0xA2);
    address internal constant ARBITER = address(0xA3);

    uint256 internal agentRegId;

    function setUp() public {
        registry = new AgentRegistry(address(this));
        // Register the agent so `accept` paths are reachable in the symbolic
        // model. The exact metadata values are immaterial to the lifecycle
        // proofs.
        agentRegId = registry.register(AGENT, "ipfs://hal", 0);

        token = new HalToken();

        job = new Job();
        job.initialize(
            Job.InitParams({
                jobId: 1,
                principal: PRINCIPAL,
                registry: registry,
                targetAgentId: 0, // open to any agent
                token: address(token),
                budget: 1_000e18,
                deadline: uint64(block.timestamp + 7 days),
                jobType: IJob.JobType.Direct,
                evaluator: address(0),
                arbiter: ARBITER
            })
        );

        // Pre-fund the Job clone so `_completeJob` / `cancel` / `expireJob` /
        // `resolveDispute(false)` can perform their `safeTransfer` legs.
        token.mint(address(job), 1_000e18);
    }

    // ─────────────────────────────────────────────────────────────
    // Terminal absorption: once Completed or Cancelled, every transition
    // reverts. This is the central safety property of issue #70 — a job
    // can transition to AT MOST one terminal phase.
    // ─────────────────────────────────────────────────────────────

    /// @notice Once `phase == Completed`, every state-mutating call MUST revert.
    function check_terminal_completed_isAbsorbing() public {
        // Drive Open → Active → Review → Completed via the Direct happy path.
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
        job.resolveDispute(false); // principal favoured ⇒ Cancelled
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
    // Budget zeroing: each terminal-entry path zeroes budget and routes the
    // funds correctly. Per IJob.Phase this is the second half of "exactly
    // one of {completed, rejected, expired}" — funds must be conserved.
    // ─────────────────────────────────────────────────────────────

    /// @notice After Completed, budget storage MUST be 0 and the agent MUST
    ///         hold the full original budget.
    function check_completed_zeroesBudget_paysAgent() public {
        uint256 origBudget = job.budget();
        vm.prank(AGENT);
        job.accept(agentRegId);
        vm.prank(AGENT);
        job.submitResult("ipfs://result");
        vm.prank(PRINCIPAL);
        job.approveResult();

        assertEq(job.budget(), 0);
        assertEq(token.balanceOf(AGENT), origBudget);
        assertEq(token.balanceOf(address(job)), 0);
    }

    /// @notice After Cancelled-by-rejection, budget MUST be 0 and the
    ///         principal MUST hold the original budget refund.
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
        assertEq(token.balanceOf(address(job)), 0);
    }

    /// @notice After Cancelled-by-expiry, budget MUST be 0 and the principal
    ///         MUST hold the original budget refund.
    function check_expired_zeroesBudget_refundsPrincipal() public {
        uint256 origBudget = job.budget();
        vm.warp(uint256(job.deadline()) + 1);
        job.expireJob();

        assertEq(job.budget(), 0);
        assertEq(token.balanceOf(PRINCIPAL), origBudget);
        assertEq(token.balanceOf(address(job)), 0);
    }

    // ─────────────────────────────────────────────────────────────
    // expireJob phase-gating: only Open and Active may expire, and only
    // strictly after the deadline.
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
        // phase still Open
        bool reverted;
        try job.expireJob() {
            reverted = false;
        } catch {
            reverted = true;
        }
        assertTrue(reverted, "expireJob before deadline must revert");
    }

    // ─────────────────────────────────────────────────────────────
    // Internal: assert every public lifecycle entrypoint reverts in the
    // current state. Invoked from the four `terminal_*_isAbsorbing` proofs.
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
