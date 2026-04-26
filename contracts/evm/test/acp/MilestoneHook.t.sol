// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {MilestoneHook} from "../../src/acp/MilestoneHook.sol";

contract MockToken is ERC20 {
    constructor() ERC20("MockMile", "MMILE") {}
    function mint(address to, uint256 amount) external { _mint(to, amount); }
}

contract MilestoneHookTest is Test {
    MilestoneHook internal hook;
    MockToken internal token;

    address internal caller = makeAddr("caller");
    address internal agent = makeAddr("agent");
    uint256 internal constant BUDGET = 1000e18;
    uint256 internal constant JOB_ID = 3;

    event MilestoneReleased(
        uint256 indexed jobId, address indexed agent, uint8 stage, uint8 totalStages, uint256 amount
    );
    event MilestoneCancelled(uint256 indexed jobId, uint8 remainingStages);

    function _initCtx(uint8 stages) internal view returns (bytes memory) {
        return abi.encode(MilestoneHook.InitContext({
            agent: agent, token: address(token), totalBudget: BUDGET, stages: stages
        }));
    }

    function setUp() public {
        hook = new MilestoneHook();
        token = new MockToken();
        token.mint(caller, BUDGET * 10);
        vm.prank(caller);
        token.approve(address(hook), type(uint256).max);
    }

    function test_hookName() public view {
        assertEq(hook.hookName(), "MilestoneHook");
    }

    function test_onAccept_initialisesState() public {
        hook.onAccept(JOB_ID, _initCtx(4));
        (address a,,uint256 sa,,uint8 total,uint8 completed, bool init) = hook.milestones(JOB_ID);
        assertEq(a, agent);
        assertEq(total, 4);
        assertEq(completed, 0);
        assertTrue(init);
        assertEq(sa, BUDGET / 4);
    }

    function test_onApprove_releasesFirstStage() public {
        hook.onAccept(JOB_ID, _initCtx(4));
        uint256 stageAmt = BUDGET / 4;

        vm.expectEmit(true, true, false, true);
        emit MilestoneReleased(JOB_ID, agent, 1, 4, stageAmt);

        vm.prank(caller);
        hook.onApprove(JOB_ID, "");
        assertEq(token.balanceOf(agent), stageAmt);
    }

    function test_onApprove_allStages_sumEqualsTotal() public {
        hook.onAccept(JOB_ID, _initCtx(3));
        for (uint8 i = 0; i < 3; ++i) {
            vm.prank(caller);
            hook.onApprove(JOB_ID, "");
        }
        assertEq(token.balanceOf(agent), BUDGET);
    }

    function test_onApprove_revert_allStagesComplete() public {
        hook.onAccept(JOB_ID, _initCtx(2));
        vm.prank(caller); hook.onApprove(JOB_ID, "");
        vm.prank(caller); hook.onApprove(JOB_ID, "");
        vm.expectRevert(abi.encodeWithSelector(MilestoneHook.AllStagesComplete.selector, JOB_ID));
        vm.prank(caller); hook.onApprove(JOB_ID, "");
    }

    function test_onApprove_revert_notInitialised() public {
        vm.expectRevert(abi.encodeWithSelector(MilestoneHook.NotInitialised.selector, 99));
        vm.prank(caller); hook.onApprove(99, "");
    }

    function test_onCancel_cancelsRemainingStages() public {
        hook.onAccept(JOB_ID, _initCtx(5));
        vm.prank(caller); hook.onApprove(JOB_ID, ""); // 1 of 5

        vm.expectEmit(true, false, false, true);
        emit MilestoneCancelled(JOB_ID, 4);
        hook.onCancel(JOB_ID, "");

        (,,,,, uint8 completed,) = hook.milestones(JOB_ID);
        assertEq(completed, 5); // all marked done
    }

    function test_onCancel_noop_notInitialised() public {
        hook.onCancel(999, ""); // no revert
    }

    function testFuzz_onApprove_dustGoesToLastStage(uint256 totalBudget, uint8 stages) public {
        totalBudget = bound(totalBudget, 1, type(uint128).max);
        stages = uint8(bound(stages, 1, 20));
        token.mint(caller, totalBudget);

        bytes memory ctx = abi.encode(MilestoneHook.InitContext({
            agent: agent, token: address(token), totalBudget: totalBudget, stages: stages
        }));
        hook.onAccept(JOB_ID + uint256(stages), ctx);
        for (uint8 i = 0; i < stages; ++i) {
            vm.prank(caller);
            hook.onApprove(JOB_ID + uint256(stages), "");
        }
        assertEq(token.balanceOf(agent), totalBudget);
    }
}
