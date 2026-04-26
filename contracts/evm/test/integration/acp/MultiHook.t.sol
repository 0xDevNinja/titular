// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {FundTransferHook} from "../../../src/acp/FundTransferHook.sol";
import {MilestoneHook} from "../../../src/acp/MilestoneHook.sol";
import {HookRegistry} from "../../../src/acp/HookRegistry.sol";
import {AgentRegistry} from "../../../src/acp/AgentRegistry.sol";

/// @notice Composition test: FundTransferHook + MilestoneHook on the same job.
///         Simulates a principal that uses MilestoneHook for staged budget release
///         and FundTransferHook for a PnL settlement on the final stage.
contract MockToken is ERC20 {
    constructor() ERC20("MultiHook", "MH") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

contract MultiHookTest is Test {
    HookRegistry internal hookRegistry;
    FundTransferHook internal fundHook;
    MilestoneHook internal milestoneHook;
    MockToken internal token;

    address internal admin = makeAddr("admin");
    address internal principal = makeAddr("principal");
    address internal agent = makeAddr("agent");
    address internal settlement = makeAddr("settlement");

    uint256 internal constant BUDGET = 3000e18;
    uint256 internal constant JOB_ID = 42;
    uint8 internal constant STAGES = 3;

    function setUp() public {
        hookRegistry = new HookRegistry(admin);
        fundHook = new FundTransferHook();
        milestoneHook = new MilestoneHook();
        token = new MockToken();

        // Register both hooks
        vm.prank(admin);
        hookRegistry.registerHook(address(fundHook));
        vm.prank(admin);
        hookRegistry.registerHook(address(milestoneHook));

        token.mint(principal, BUDGET * 10);
        vm.prank(principal);
        token.approve(address(fundHook), type(uint256).max);
        vm.prank(principal);
        token.approve(address(milestoneHook), type(uint256).max);
    }

    /// @notice Both hooks are approved in HookRegistry.
    function test_bothHooksRegistered() public view {
        assertTrue(hookRegistry.isApproved(address(fundHook)));
        assertTrue(hookRegistry.isApproved(address(milestoneHook)));
    }

    /// @notice MilestoneHook: 3-stage job, each stage releases BUDGET/3.
    function test_milestoneHook_threeStages() public {
        bytes memory initCtx = abi.encode(
            MilestoneHook.InitContext({agent: agent, token: address(token), totalBudget: BUDGET, stages: STAGES})
        );
        milestoneHook.onAccept(JOB_ID, initCtx);

        uint256 stageAmt = BUDGET / STAGES;
        for (uint8 i = 0; i < STAGES; ++i) {
            vm.prank(principal);
            milestoneHook.onApprove(JOB_ID, "");
        }
        assertEq(token.balanceOf(agent), BUDGET);
    }

    /// @notice FundTransferHook: final stage with 5% PnL settlement.
    function test_fundTransferHook_withPnlOnFinalStage() public {
        uint256 finalStageAmt = BUDGET / STAGES;
        bytes memory ctx = abi.encode(
            FundTransferHook.ApproveContext({
                agent: agent,
                token: address(token),
                amount: finalStageAmt,
                pnlBps: 500, // 5% PnL
                settlementAddr: settlement
            })
        );
        vm.prank(principal);
        fundHook.onApprove(JOB_ID, ctx);

        assertEq(token.balanceOf(agent), finalStageAmt * 95 / 100);
        assertEq(token.balanceOf(settlement), finalStageAmt * 5 / 100);
    }

    /// @notice Combined: milestone pays first 2 stages normally, fund-transfer handles stage 3 with PnL.
    function test_combined_milestoneAndFundTransfer() public {
        bytes memory initCtx = abi.encode(
            MilestoneHook.InitContext({agent: agent, token: address(token), totalBudget: BUDGET, stages: STAGES})
        );
        milestoneHook.onAccept(JOB_ID, initCtx);

        uint256 stageAmt = BUDGET / STAGES;

        // Stages 1 and 2 via MilestoneHook
        vm.prank(principal);
        milestoneHook.onApprove(JOB_ID, "");
        vm.prank(principal);
        milestoneHook.onApprove(JOB_ID, "");

        assertEq(token.balanceOf(agent), stageAmt * 2);

        // Stage 3 via FundTransferHook with 10% PnL settlement
        bytes memory ftCtx = abi.encode(
            FundTransferHook.ApproveContext({
                agent: agent,
                token: address(token),
                amount: stageAmt + (BUDGET % STAGES), // stage amount + dust
                pnlBps: 1000, // 10%
                settlementAddr: settlement
            })
        );
        uint256 finalAmt = stageAmt + (BUDGET % STAGES);
        vm.prank(principal);
        fundHook.onApprove(JOB_ID, ftCtx);

        uint256 totalExpected = stageAmt * 2 + finalAmt * 90 / 100;
        assertEq(token.balanceOf(agent), totalExpected);
        assertEq(token.balanceOf(settlement), finalAmt * 10 / 100);
    }

    /// @notice Deregistering a hook prevents it from being used (logical guard at caller layer).
    function test_deregisterHook_isNoLongerApproved() public {
        vm.prank(admin);
        hookRegistry.deregisterHook(address(fundHook));
        assertFalse(hookRegistry.isApproved(address(fundHook)));
        // MilestoneHook unaffected
        assertTrue(hookRegistry.isApproved(address(milestoneHook)));
    }
}
