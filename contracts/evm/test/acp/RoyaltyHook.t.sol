// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {RoyaltyHook} from "../../src/acp/RoyaltyHook.sol";

contract MockToken is ERC20 {
    constructor() ERC20("MockRoy", "MROY") {}
    function mint(address to, uint256 amount) external { _mint(to, amount); }
}

contract RoyaltyHookTest is Test {
    RoyaltyHook internal hook;
    MockToken internal token;

    address internal caller = makeAddr("caller");
    address internal r1 = makeAddr("r1");
    address internal r2 = makeAddr("r2");
    address internal r3 = makeAddr("r3");
    uint256 internal constant AMOUNT = 1000e18;
    uint256 internal constant JOB_ID = 10;

    event RoyaltyPaid(
        uint256 indexed jobId, address indexed token, address indexed recipient, uint256 amount
    );

    function setUp() public {
        hook = new RoyaltyHook();
        token = new MockToken();
        token.mint(caller, AMOUNT * 10);
        vm.prank(caller);
        token.approve(address(hook), type(uint256).max);
    }

    function test_hookName() public view {
        assertEq(hook.hookName(), "RoyaltyHook");
    }

    function test_onApprove_twoRecipients() public {
        RoyaltyHook.RoyaltyRecipient[] memory recs = new RoyaltyHook.RoyaltyRecipient[](2);
        recs[0] = RoyaltyHook.RoyaltyRecipient({recipient: r1, shareBps: 7000});
        recs[1] = RoyaltyHook.RoyaltyRecipient({recipient: r2, shareBps: 3000});
        bytes memory ctx = abi.encode(RoyaltyHook.RoyaltyContext({
            token: address(token), amount: AMOUNT, recipients: recs
        }));
        vm.prank(caller);
        hook.onApprove(JOB_ID, ctx);
        assertEq(token.balanceOf(r1), 700e18);
        assertEq(token.balanceOf(r2), 300e18);
    }

    function test_onApprove_threeRecipients_dustToLast() public {
        // 10 000 shares into 3 splits of 3333/3333/3334
        RoyaltyHook.RoyaltyRecipient[] memory recs = new RoyaltyHook.RoyaltyRecipient[](3);
        recs[0] = RoyaltyHook.RoyaltyRecipient({recipient: r1, shareBps: 3333});
        recs[1] = RoyaltyHook.RoyaltyRecipient({recipient: r2, shareBps: 3333});
        recs[2] = RoyaltyHook.RoyaltyRecipient({recipient: r3, shareBps: 3334});
        bytes memory ctx = abi.encode(RoyaltyHook.RoyaltyContext({
            token: address(token), amount: AMOUNT, recipients: recs
        }));
        vm.prank(caller);
        hook.onApprove(JOB_ID, ctx);
        uint256 totalOut = token.balanceOf(r1) + token.balanceOf(r2) + token.balanceOf(r3);
        assertEq(totalOut, AMOUNT);
    }

    function test_onApprove_revert_sumMismatch() public {
        RoyaltyHook.RoyaltyRecipient[] memory recs = new RoyaltyHook.RoyaltyRecipient[](2);
        recs[0] = RoyaltyHook.RoyaltyRecipient({recipient: r1, shareBps: 5000});
        recs[1] = RoyaltyHook.RoyaltyRecipient({recipient: r2, shareBps: 5001}); // 10001
        bytes memory ctx = abi.encode(RoyaltyHook.RoyaltyContext({
            token: address(token), amount: AMOUNT, recipients: recs
        }));
        vm.expectRevert(abi.encodeWithSelector(RoyaltyHook.ShareSumMismatch.selector, 10_001));
        vm.prank(caller);
        hook.onApprove(JOB_ID, ctx);
    }

    function test_onApprove_revert_noRecipients() public {
        RoyaltyHook.RoyaltyRecipient[] memory recs = new RoyaltyHook.RoyaltyRecipient[](0);
        bytes memory ctx = abi.encode(RoyaltyHook.RoyaltyContext({
            token: address(token), amount: AMOUNT, recipients: recs
        }));
        vm.expectRevert(RoyaltyHook.NoRecipients.selector);
        vm.prank(caller);
        hook.onApprove(JOB_ID, ctx);
    }

    function test_onApprove_revert_zeroToken() public {
        RoyaltyHook.RoyaltyRecipient[] memory recs = new RoyaltyHook.RoyaltyRecipient[](1);
        recs[0] = RoyaltyHook.RoyaltyRecipient({recipient: r1, shareBps: 10_000});
        bytes memory ctx = abi.encode(RoyaltyHook.RoyaltyContext({
            token: address(0), amount: AMOUNT, recipients: recs
        }));
        vm.expectRevert(RoyaltyHook.ZeroAddress.selector);
        vm.prank(caller);
        hook.onApprove(JOB_ID, ctx);
    }

    function test_noOp_hooks() public {
        hook.onAccept(JOB_ID, "");
        hook.onSubmit(JOB_ID, "");
        hook.onReject(JOB_ID, "");
        hook.onCancel(JOB_ID, "");
    }

    function testFuzz_onApprove_sumInvariant(uint256 amount) public {
        amount = bound(amount, 2, type(uint128).max);
        token.mint(caller, amount);
        RoyaltyHook.RoyaltyRecipient[] memory recs = new RoyaltyHook.RoyaltyRecipient[](2);
        recs[0] = RoyaltyHook.RoyaltyRecipient({recipient: r1, shareBps: 6000});
        recs[1] = RoyaltyHook.RoyaltyRecipient({recipient: r2, shareBps: 4000});
        bytes memory ctx = abi.encode(RoyaltyHook.RoyaltyContext({
            token: address(token), amount: amount, recipients: recs
        }));
        vm.prank(caller);
        hook.onApprove(JOB_ID, ctx);
        assertEq(token.balanceOf(r1) + token.balanceOf(r2), amount);
    }
}
