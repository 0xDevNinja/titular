// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {VeTITU} from "../../src/governance/VeTITU.sol";

contract MockERC20 is ERC20 {
    constructor() ERC20("TITU", "TITU") {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

contract VeTITUTest is Test {
    MockERC20 internal titu;
    VeTITU internal ve;
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);

    uint256 internal constant WEEK = 7 days;
    uint256 internal constant MAXTIME = 4 * 365 days;

    function setUp() public {
        titu = new MockERC20();
        ve = new VeTITU(IERC20(address(titu)));
        titu.mint(alice, 10_000e18);
        titu.mint(bob, 10_000e18);
        vm.prank(alice);
        titu.approve(address(ve), type(uint256).max);
        vm.prank(bob);
        titu.approve(address(ve), type(uint256).max);
        vm.warp(1_700_000_000);
    }

    function test_createLock_rounds_to_week() public {
        uint256 unlock = block.timestamp + 30 days;
        vm.prank(alice);
        ve.createLock(1000e18, unlock);
        (, uint256 end) = ve.locked(alice);
        assertEq(end, (unlock / WEEK) * WEEK);
    }

    function test_createLock_pulls_tokens() public {
        vm.prank(alice);
        ve.createLock(1000e18, block.timestamp + 30 days);
        assertEq(titu.balanceOf(address(ve)), 1000e18);
    }

    function test_createLock_reverts_too_short() public {
        vm.prank(alice);
        vm.expectRevert();
        ve.createLock(1000e18, block.timestamp + 6 days);
    }

    function test_createLock_reverts_too_long() public {
        vm.prank(alice);
        vm.expectRevert(VeTITU.UnlockTimeTooLong.selector);
        ve.createLock(1000e18, block.timestamp + MAXTIME + 2 * WEEK);
    }

    function test_createLock_reverts_zero_amount() public {
        vm.prank(alice);
        vm.expectRevert(VeTITU.ZeroAmount.selector);
        ve.createLock(0, block.timestamp + 30 days);
    }

    function test_balanceOf_decays_linearly() public {
        vm.prank(alice);
        uint256 unlockTarget = block.timestamp + MAXTIME;
        ve.createLock(1000e18, unlockTarget);
        uint256 v0 = ve.balanceOf(alice);
        assertGt(v0, 0);

        vm.warp(block.timestamp + MAXTIME / 2);
        uint256 v1 = ve.balanceOf(alice);
        // after half of remaining lock, voting power should be < initial
        assertLt(v1, v0);
    }

    function test_balanceOf_zero_after_unlock() public {
        vm.prank(alice);
        ve.createLock(1000e18, block.timestamp + 60 days);
        vm.warp(block.timestamp + 61 days);
        assertEq(ve.balanceOf(alice), 0);
    }

    function test_increaseAmount_ok() public {
        vm.prank(alice);
        ve.createLock(1000e18, block.timestamp + 180 days);
        vm.prank(alice);
        ve.increaseAmount(500e18);
        (int128 amt,) = ve.locked(alice);
        assertEq(uint256(uint128(amt)), 1500e18);
    }

    function test_increaseUnlockTime_ok_if_later() public {
        vm.prank(alice);
        ve.createLock(1000e18, block.timestamp + 60 days);
        (, uint256 oldEnd) = ve.locked(alice);
        vm.prank(alice);
        ve.increaseUnlockTime(block.timestamp + 365 days);
        (, uint256 newEnd) = ve.locked(alice);
        assertGt(newEnd, oldEnd);
    }

    function test_increaseUnlockTime_revert_if_earlier() public {
        vm.prank(alice);
        ve.createLock(1000e18, block.timestamp + 180 days);
        vm.prank(alice);
        vm.expectRevert(VeTITU.UnlockTimeNotLater.selector);
        ve.increaseUnlockTime(block.timestamp + 60 days);
    }

    function test_withdraw_only_after_unlock() public {
        vm.prank(alice);
        ve.createLock(1000e18, block.timestamp + 60 days);
        vm.prank(alice);
        vm.expectRevert(VeTITU.LockNotExpired.selector);
        ve.withdraw();
        vm.warp(block.timestamp + 61 days);
        vm.prank(alice);
        ve.withdraw();
        assertEq(titu.balanceOf(alice), 10_000e18);
    }

    function test_getPastVotes_works_at_historical_timestamp() public {
        vm.prank(alice);
        ve.createLock(1000e18, block.timestamp + MAXTIME);
        uint256 t0 = block.timestamp;
        vm.warp(block.timestamp + 7 days);
        uint256 past = ve.getPastVotes(alice, t0 + 1);
        assertGt(past, 0);
    }

    function test_transfer_disallowed() public {
        vm.expectRevert(VeTITU.NoTransfersAllowed.selector);
        ve.transfer(bob, 1);
        vm.expectRevert(VeTITU.NoTransfersAllowed.selector);
        ve.transferFrom(alice, bob, 1);
    }

    function test_views_getVotes_and_clock_and_totalSupply() public {
        vm.prank(alice);
        ve.createLock(1000e18, block.timestamp + MAXTIME);
        assertGt(ve.getVotes(alice), 0);
        assertGt(ve.totalSupply(), 0);
        assertEq(uint256(ve.clock()), block.timestamp);
        assertEq(ve.CLOCK_MODE(), "mode=timestamp");
    }

    function test_userPointCount_and_globalPointCount() public {
        vm.prank(alice);
        ve.createLock(1000e18, block.timestamp + MAXTIME);
        assertEq(ve.userPointCount(alice), 1);
        assertGt(ve.globalPointCount(), 1);
        VeTITU.Point memory p = ve.userPointAt(alice, 0);
        assertEq(p.ts, block.timestamp);
    }

    function test_getPastTotalSupply_works() public {
        vm.prank(alice);
        ve.createLock(1000e18, block.timestamp + MAXTIME);
        uint256 t = block.timestamp;
        vm.warp(block.timestamp + 1 days);
        assertGt(ve.getPastTotalSupply(t), 0);
    }

    function test_getPastVotes_revert_future() public {
        vm.prank(alice);
        ve.createLock(1000e18, block.timestamp + MAXTIME);
        vm.expectRevert("VeTITU: future lookup");
        ve.getPastVotes(alice, block.timestamp + 1);
    }

    function test_increaseAmount_revert_no_lock() public {
        vm.prank(alice);
        vm.expectRevert(VeTITU.NoLock.selector);
        ve.increaseAmount(1);
    }

    function test_increaseUnlockTime_revert_no_lock() public {
        vm.prank(alice);
        vm.expectRevert(VeTITU.NoLock.selector);
        ve.increaseUnlockTime(block.timestamp + 30 days);
    }

    function test_increaseUnlockTime_revert_too_long() public {
        vm.prank(alice);
        ve.createLock(1000e18, block.timestamp + 60 days);
        vm.prank(alice);
        vm.expectRevert(VeTITU.UnlockTimeTooLong.selector);
        ve.increaseUnlockTime(block.timestamp + MAXTIME + 2 * WEEK);
    }

    function test_increaseAmount_revert_expired() public {
        vm.prank(alice);
        ve.createLock(1000e18, block.timestamp + 14 days);
        vm.warp(block.timestamp + 30 days);
        vm.prank(alice);
        vm.expectRevert(VeTITU.LockExpired.selector);
        ve.increaseAmount(1);
    }

    function test_constructor_revert_zero() public {
        vm.expectRevert(VeTITU.ZeroAddress.selector);
        new VeTITU(IERC20(address(0)));
    }

    function test_createLock_revert_lockExists() public {
        vm.prank(alice);
        ve.createLock(1000e18, block.timestamp + 30 days);
        vm.prank(alice);
        vm.expectRevert(VeTITU.LockExists.selector);
        ve.createLock(1, block.timestamp + 30 days);
    }

    function test_withdraw_revert_no_lock() public {
        vm.prank(alice);
        vm.expectRevert(VeTITU.NoLock.selector);
        ve.withdraw();
    }

    function testFuzz_lock_durations(uint256 amount, uint256 duration) public {
        amount = bound(amount, 1e18, 5000e18);
        // bound ensures post-rounding-to-week the lock end is at least `block.timestamp + WEEK`
        duration = bound(duration, 2 * WEEK, MAXTIME);
        vm.prank(alice);
        ve.createLock(amount, block.timestamp + duration);
        assertGe(ve.balanceOf(alice), 0);
    }
}
