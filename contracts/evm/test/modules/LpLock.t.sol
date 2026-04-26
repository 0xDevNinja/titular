// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import {LPLock} from "../../src/launchpad/LPLock.sol";

/// @dev Mintable ERC-20 used as the LP token under test.
contract LpLockMockLP is ERC20 {
    constructor() ERC20("LP", "LP") {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @title  LpLockModuleSuiteTest
/// @notice Per-contract unit tests for {LPLock}: 10-year timelock semantics,
///         deposit accumulation, beneficiary-only one-shot withdraw, and the
///         time-remaining view. Coverage gate ≥95%.
contract LpLockModuleSuiteTest is Test {
    LPLock internal lock;
    LpLockMockLP internal lp;

    address internal beneficiary = address(0xBEEF);
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);

    uint256 internal constant TEN_YEARS = 365 days * 10;

    function setUp() public {
        lp = new LpLockMockLP();
        lock = new LPLock(IERC20(address(lp)), beneficiary);
        lp.mint(alice, 1_000_000e18);
        lp.mint(bob, 1_000_000e18);
        vm.prank(alice);
        lp.approve(address(lock), type(uint256).max);
        vm.prank(bob);
        lp.approve(address(lock), type(uint256).max);
    }

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    function test_constructor_zero_token_reverts() public {
        vm.expectRevert(LPLock.ZeroAddress.selector);
        new LPLock(IERC20(address(0)), beneficiary);
    }

    function test_constructor_zero_beneficiary_reverts() public {
        vm.expectRevert(LPLock.ZeroAddress.selector);
        new LPLock(IERC20(address(lp)), address(0));
    }

    function test_constructor_pins_immutables() public view {
        assertEq(address(lock.lpToken()), address(lp));
        assertEq(lock.beneficiary(), beneficiary);
        assertEq(lock.unlockTime(), block.timestamp + TEN_YEARS);
        assertEq(lock.LOCK_DURATION(), TEN_YEARS);
        assertEq(lock.lockedAmount(), 0);
        assertFalse(lock.withdrawn());
    }

    // ---------------------------------------------------------------------
    // Deposit
    // ---------------------------------------------------------------------

    function test_deposit_increments_locked_amount() public {
        vm.expectEmit(true, false, false, true, address(lock));
        emit LPLock.Deposited(alice, 1000e18);
        vm.prank(alice);
        lock.deposit(1000e18);

        assertEq(lock.lockedAmount(), 1000e18);
        assertEq(lp.balanceOf(address(lock)), 1000e18);

        vm.prank(alice);
        lock.deposit(500e18);
        assertEq(lock.lockedAmount(), 1500e18);
    }

    function test_deposit_anyone_allowed_accumulates() public {
        vm.prank(alice);
        lock.deposit(100e18);
        vm.prank(bob);
        lock.deposit(250e18);
        assertEq(lock.lockedAmount(), 350e18);
    }

    function test_deposit_zero_amount_reverts() public {
        vm.prank(alice);
        vm.expectRevert(LPLock.ZeroAmount.selector);
        lock.deposit(0);
    }

    // ---------------------------------------------------------------------
    // Withdraw
    // ---------------------------------------------------------------------

    function test_withdraw_pre_unlock_reverts() public {
        vm.prank(alice);
        lock.deposit(1000e18);

        uint256 remaining = lock.unlockTime() - block.timestamp;
        vm.prank(beneficiary);
        vm.expectRevert(abi.encodeWithSelector(LPLock.StillLocked.selector, remaining));
        lock.withdraw();
    }

    function test_withdraw_non_beneficiary_reverts() public {
        vm.prank(alice);
        lock.deposit(1000e18);
        vm.warp(block.timestamp + TEN_YEARS);
        vm.prank(alice);
        vm.expectRevert(LPLock.NotBeneficiary.selector);
        lock.withdraw();
    }

    function test_withdraw_after_unlock_succeeds() public {
        vm.prank(alice);
        lock.deposit(1000e18);
        vm.warp(block.timestamp + TEN_YEARS);

        vm.expectEmit(true, false, false, true, address(lock));
        emit LPLock.Withdrawn(beneficiary, 1000e18);
        vm.prank(beneficiary);
        lock.withdraw();

        assertTrue(lock.withdrawn());
        assertEq(lp.balanceOf(beneficiary), 1000e18);
        assertEq(lp.balanceOf(address(lock)), 0);
    }

    function test_withdraw_sweeps_full_balance_including_dust() public {
        vm.prank(alice);
        lock.deposit(1000e18);
        vm.prank(bob);
        lock.deposit(2000e18);
        // Direct mint into the lock (not through deposit) — still swept.
        lp.mint(address(lock), 500e18);

        uint256 total = lp.balanceOf(address(lock));
        assertEq(total, 3500e18);

        vm.warp(lock.unlockTime());
        vm.prank(beneficiary);
        lock.withdraw();

        assertEq(lp.balanceOf(beneficiary), total);
        assertEq(lock.lockedAmount(), 3000e18);
    }

    function test_withdraw_double_call_reverts() public {
        vm.prank(alice);
        lock.deposit(1000e18);
        vm.warp(lock.unlockTime());
        vm.prank(beneficiary);
        lock.withdraw();
        vm.prank(beneficiary);
        vm.expectRevert(LPLock.AlreadyWithdrawn.selector);
        lock.withdraw();
    }

    function test_withdraw_zero_balance_succeeds_one_shot() public {
        // No deposits — a beneficiary call after unlock should succeed but
        // transfer 0; subsequent calls revert AlreadyWithdrawn.
        vm.warp(lock.unlockTime());
        vm.expectEmit(true, false, false, true, address(lock));
        emit LPLock.Withdrawn(beneficiary, 0);
        vm.prank(beneficiary);
        lock.withdraw();

        assertTrue(lock.withdrawn());

        vm.prank(beneficiary);
        vm.expectRevert(LPLock.AlreadyWithdrawn.selector);
        lock.withdraw();
    }

    // ---------------------------------------------------------------------
    // Views
    // ---------------------------------------------------------------------

    function test_timeRemaining_zero_after_unlock() public {
        uint256 unlock = lock.unlockTime();
        assertEq(lock.timeRemaining(), TEN_YEARS);

        vm.warp(unlock - 7 days);
        assertEq(lock.timeRemaining(), 7 days);

        vm.warp(unlock);
        assertEq(lock.timeRemaining(), 0);

        vm.warp(unlock + 1 days);
        assertEq(lock.timeRemaining(), 0);
    }

    // ---------------------------------------------------------------------
    // Fuzz
    // ---------------------------------------------------------------------

    /// @dev Cumulative deposits track the running sum exactly. Bounded amount
    ///      list so fuzz time stays sane.
    function testFuzz_deposit_accumulation(uint128[] memory amounts) public {
        uint256 cap = amounts.length > 16 ? 16 : amounts.length;
        uint256 running;

        for (uint256 i; i < cap; ++i) {
            uint256 amt = uint256(amounts[i]);
            uint256 remaining = lp.balanceOf(alice);
            if (amt == 0 || amt > remaining) continue;
            vm.prank(alice);
            lock.deposit(amt);
            running += amt;

            assertEq(lock.lockedAmount(), running);
            assertEq(lp.balanceOf(address(lock)), running);
        }
    }
}
