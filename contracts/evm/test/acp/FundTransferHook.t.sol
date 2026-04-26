// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {FundTransferHook} from "../../src/acp/FundTransferHook.sol";

contract MockToken is ERC20 {
    constructor() ERC20("Mock", "MCK") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

contract FundTransferHookTest is Test {
    FundTransferHook internal hook;
    MockToken internal token;

    address internal caller = makeAddr("caller");
    address internal agent = makeAddr("agent");
    address internal settlement = makeAddr("settlement");

    uint256 internal constant AMOUNT = 1000e18;
    uint256 internal constant JOB_ID = 1;

    event FundTransferred(
        uint256 indexed jobId,
        address indexed agent,
        address indexed token,
        uint256 agentAmount,
        uint256 pnlAmount,
        address settlementAddr
    );

    function setUp() public {
        hook = new FundTransferHook();
        token = new MockToken();
        token.mint(caller, AMOUNT * 10);
        vm.prank(caller);
        token.approve(address(hook), type(uint256).max);
    }

    function test_hookName() public view {
        assertEq(hook.hookName(), "FundTransferHook");
    }

    function test_onApprove_noSplit() public {
        bytes memory ctx = abi.encode(
            FundTransferHook.ApproveContext({
                agent: agent, token: address(token), amount: AMOUNT, pnlBps: 0, settlementAddr: address(0)
            })
        );
        vm.expectEmit(true, true, true, true);
        emit FundTransferred(JOB_ID, agent, address(token), AMOUNT, 0, address(0));
        vm.prank(caller);
        hook.onApprove(JOB_ID, ctx);
        assertEq(token.balanceOf(agent), AMOUNT);
    }

    function test_onApprove_withPnlSplit() public {
        // 10% PnL split
        bytes memory ctx = abi.encode(
            FundTransferHook.ApproveContext({
                agent: agent, token: address(token), amount: AMOUNT, pnlBps: 1000, settlementAddr: settlement
            })
        );
        vm.prank(caller);
        hook.onApprove(JOB_ID, ctx);
        assertEq(token.balanceOf(agent), AMOUNT * 90 / 100);
        assertEq(token.balanceOf(settlement), AMOUNT * 10 / 100);
    }

    function test_onApprove_revert_invalidBps() public {
        bytes memory ctx = abi.encode(
            FundTransferHook.ApproveContext({
                agent: agent, token: address(token), amount: AMOUNT, pnlBps: 10_001, settlementAddr: settlement
            })
        );
        vm.expectRevert(FundTransferHook.InvalidBps.selector);
        vm.prank(caller);
        hook.onApprove(JOB_ID, ctx);
    }

    function test_onApprove_revert_zeroAgent() public {
        bytes memory ctx = abi.encode(
            FundTransferHook.ApproveContext({
                agent: address(0), token: address(token), amount: AMOUNT, pnlBps: 0, settlementAddr: address(0)
            })
        );
        vm.expectRevert(FundTransferHook.ZeroAddress.selector);
        vm.prank(caller);
        hook.onApprove(JOB_ID, ctx);
    }

    function test_onApprove_revert_zeroAmount() public {
        bytes memory ctx = abi.encode(
            FundTransferHook.ApproveContext({
                agent: agent, token: address(token), amount: 0, pnlBps: 0, settlementAddr: address(0)
            })
        );
        vm.expectRevert(FundTransferHook.ZeroAmount.selector);
        vm.prank(caller);
        hook.onApprove(JOB_ID, ctx);
    }

    function test_noOp_onAccept() public {
        vm.prank(caller);
        hook.onAccept(JOB_ID, "");
    }

    function test_noOp_onSubmit() public {
        vm.prank(caller);
        hook.onSubmit(JOB_ID, "");
    }

    function testFuzz_onApprove_sumInvariant(uint256 amount, uint256 pnlBps) public {
        amount = bound(amount, 1, type(uint128).max);
        pnlBps = bound(pnlBps, 0, 10_000);
        token.mint(caller, amount);

        bytes memory ctx = abi.encode(
            FundTransferHook.ApproveContext({
                agent: agent, token: address(token), amount: amount, pnlBps: pnlBps, settlementAddr: settlement
            })
        );
        uint256 agentBefore = token.balanceOf(agent);
        uint256 settlementBefore = token.balanceOf(settlement);
        vm.prank(caller);
        hook.onApprove(JOB_ID, ctx);

        uint256 totalOut = token.balanceOf(agent) - agentBefore + token.balanceOf(settlement) - settlementBefore;
        assertEq(totalOut, amount);
    }
}
