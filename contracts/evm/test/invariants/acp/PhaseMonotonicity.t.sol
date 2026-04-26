// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {StdInvariant} from "forge-std/StdInvariant.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC20Burnable} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import {Clones} from "@openzeppelin/contracts/proxy/Clones.sol";
import {Job} from "../../../src/acp/Job.sol";
import {IJob} from "../../../src/acp/IJob.sol";
import {AgentRegistry} from "../../../src/acp/AgentRegistry.sol";
import {Escrow} from "../../../src/acp/Escrow.sol";
import {FeeSplitter} from "../../../src/acp/FeeSplitter.sol";
import {HookRegistry} from "../../../src/acp/HookRegistry.sol";
import {BuybackBurner} from "../../../src/acp/BuybackBurner.sol";
import {IPermit2} from "../../../src/acp/IPermit2.sol";

contract MockToken is ERC20 {
    constructor() ERC20("PhaseInv", "PINV") {}

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
        titu.mint(to, amountIn);
    }
}

/// @dev Handler that drives a Job through valid transitions.
contract JobPhaseHandler is Test {
    Job public job;
    AgentRegistry public registry;
    Escrow public escrow;
    address public principal;
    address public agentCtrl;
    address public arbiter;
    uint256 public agentId;
    MockToken public token;
    uint256 public jobId;

    // Track highest phase reached (monotone for non-reversible paths)
    uint8 public maxPhaseReached;

    constructor(
        Job _job,
        AgentRegistry _registry,
        MockToken _token,
        Escrow _escrow,
        address _principal,
        address _agentCtrl,
        uint256 _agentId,
        address _arbiter,
        uint256 _jobId
    ) {
        job = _job;
        registry = _registry;
        token = _token;
        escrow = _escrow;
        principal = _principal;
        agentCtrl = _agentCtrl;
        agentId = _agentId;
        arbiter = _arbiter;
        jobId = _jobId;
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
    MockTITU internal titu;
    MockPermit2 internal permit2;
    MockUniRouter internal uniRouter;
    Job internal jobImpl;
    Job internal jobContract;
    Escrow internal escrow;
    FeeSplitter internal feeSplitter;
    BuybackBurner internal buybackBurner;
    HookRegistry internal hookRegistry;
    JobPhaseHandler internal handler;

    address internal admin = makeAddr("admin");
    address internal scorer = makeAddr("scorer");
    address internal principal = makeAddr("principal");
    address internal agentCtrl = makeAddr("agentCtrl");
    address internal arbiter = makeAddr("arbiter");
    address internal treasury = makeAddr("treasury");
    uint256 internal agentId;
    uint256 internal constant BUDGET = 1000e18;
    uint256 internal constant JOB_ID = 0;

    function setUp() public {
        registry = new AgentRegistry(admin);
        vm.prank(agentCtrl);
        agentId = registry.register(agentCtrl, "ipfs://QmAgent", 0x1);

        token = new MockToken();
        titu = new MockTITU();
        permit2 = new MockPermit2();
        uniRouter = new MockUniRouter(titu);

        escrow = new Escrow(admin, address(permit2));
        address[] memory swapPath = new address[](2);
        swapPath[0] = address(token);
        swapPath[1] = address(titu);
        buybackBurner = new BuybackBurner(admin, address(uniRouter), address(token), address(titu), swapPath, 0, 300);
        feeSplitter = new FeeSplitter(admin, treasury, address(buybackBurner));
        hookRegistry = new HookRegistry(admin);

        jobImpl = new Job();
        jobContract = Job(Clones.clone(address(jobImpl)));

        // Fund escrow with budget
        token.mint(principal, BUDGET * 100);
        vm.prank(principal);
        token.approve(address(escrow), BUDGET);
        vm.prank(principal);
        escrow.fund(principal, JOB_ID, address(token), BUDGET);

        // Grant clone RELEASER_ROLE and FeeSplitter CALLER_ROLE
        bytes32 releaserRole = escrow.RELEASER_ROLE();
        vm.prank(admin);
        escrow.grantRole(releaserRole, address(jobContract));
        bytes32 callerRole = feeSplitter.CALLER_ROLE();
        vm.prank(admin);
        feeSplitter.grantRole(callerRole, address(jobContract));

        address[] memory hooks = new address[](0);
        Job.InitParams memory p = Job.InitParams({
            jobId: JOB_ID,
            principal: principal,
            registry: registry,
            targetAgentId: 0,
            token: address(token),
            budget: BUDGET,
            deadline: uint64(block.timestamp + 30 days),
            jobType: IJob.JobType.Direct,
            evaluator: address(0),
            arbiter: arbiter,
            escrow: escrow,
            feeSplitter: feeSplitter,
            hookRegistry: hookRegistry,
            hooks: hooks
        });
        jobContract.initialize(p);

        handler =
            new JobPhaseHandler(jobContract, registry, token, escrow, principal, agentCtrl, agentId, arbiter, JOB_ID);
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

    /// @notice Escrow balance mirrors job budget until terminal phase.
    function invariant_escrowBalanceMirrorsJobBudget() public view {
        IJob.Phase p = jobContract.phase();
        if (p == IJob.Phase.Completed || p == IJob.Phase.Cancelled) {
            // After terminal: both should be 0
            assertEq(escrow.getBalance(principal, JOB_ID, address(token)), 0, "escrow not drained at terminal");
            assertEq(jobContract.budget(), 0, "budget not zeroed at terminal");
        } else {
            // Before terminal: escrow balance must equal job budget
            assertEq(
                escrow.getBalance(principal, JOB_ID, address(token)),
                jobContract.budget(),
                "escrow balance differs from job budget"
            );
        }
    }
}
