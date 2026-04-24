// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {LPLock} from "../../src/launchpad/LPLock.sol";

/// @dev Minimal mintable mock standing in for a Uniswap V2 LP ERC-20.
///      No transfer hooks — LPLock correctness does not rely on any quirk beyond
///      the standard ERC-20 surface.
contract MockLP is ERC20 {
    constructor() ERC20("LP", "LP") {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

contract LPLockTest is Test {
    LPLock internal lock;
    MockLP internal lp;

    address internal beneficiary = address(0xBEEF);
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);

    uint256 internal constant TEN_YEARS = 365 days * 10;

    event Deposited(address indexed from, uint256 amount);
    event Withdrawn(address indexed to, uint256 amount);

    function setUp() public {
        lp = new MockLP();
        lock = new LPLock(IERC20(address(lp)), beneficiary);

        // Fund the depositors and approve the lock.
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

    function test_constructor_sets_immutables() public view {
        assertEq(address(lock.lpToken()), address(lp));
        assertEq(lock.beneficiary(), beneficiary);
        assertEq(lock.unlockTime(), block.timestamp + TEN_YEARS);
        assertEq(lock.LOCK_DURATION(), TEN_YEARS);
        assertEq(lock.lockedAmount(), 0);
        assertFalse(lock.withdrawn());
    }

    function test_constructor_reverts_zero_token() public {
        vm.expectRevert(LPLock.ZeroAddress.selector);
        new LPLock(IERC20(address(0)), beneficiary);
    }

    function test_constructor_reverts_zero_beneficiary() public {
        vm.expectRevert(LPLock.ZeroAddress.selector);
        new LPLock(IERC20(address(lp)), address(0));
    }

    // ---------------------------------------------------------------------
    // Deposit
    // ---------------------------------------------------------------------

    function test_deposit_increments_locked_amount() public {
        vm.expectEmit(true, false, false, true, address(lock));
        emit Deposited(alice, 1000e18);
        vm.prank(alice);
        lock.deposit(1000e18);

        assertEq(lock.lockedAmount(), 1000e18);
        assertEq(lp.balanceOf(address(lock)), 1000e18);

        // Second deposit from the same depositor accumulates.
        vm.prank(alice);
        lock.deposit(500e18);
        assertEq(lock.lockedAmount(), 1500e18);
        assertEq(lp.balanceOf(address(lock)), 1500e18);
    }

    function test_deposit_anyone_allowed() public {
        // Two independent depositors should both succeed; totals accumulate.
        vm.prank(alice);
        lock.deposit(100e18);
        vm.prank(bob);
        lock.deposit(250e18);

        assertEq(lock.lockedAmount(), 350e18);
        assertEq(lp.balanceOf(address(lock)), 350e18);
    }

    function test_deposit_reverts_zero_amount() public {
        vm.prank(alice);
        vm.expectRevert(LPLock.ZeroAmount.selector);
        lock.deposit(0);
    }

    // ---------------------------------------------------------------------
    // Withdraw — gating
    // ---------------------------------------------------------------------

    function test_withdraw_reverts_before_unlock() public {
        vm.prank(alice);
        lock.deposit(1000e18);

        uint256 expectedRemaining = lock.unlockTime() - block.timestamp;
        vm.prank(beneficiary);
        vm.expectRevert(abi.encodeWithSelector(LPLock.StillLocked.selector, expectedRemaining));
        lock.withdraw();
    }

    function test_withdraw_reverts_not_beneficiary() public {
        vm.prank(alice);
        lock.deposit(1000e18);

        // Even after the lock expires, only the beneficiary may withdraw.
        vm.warp(block.timestamp + TEN_YEARS);
        vm.prank(alice);
        vm.expectRevert(LPLock.NotBeneficiary.selector);
        lock.withdraw();
    }

    // ---------------------------------------------------------------------
    // Withdraw — happy path
    // ---------------------------------------------------------------------

    function test_withdraw_succeeds_after_unlock() public {
        vm.prank(alice);
        lock.deposit(1000e18);

        vm.warp(block.timestamp + TEN_YEARS);

        vm.expectEmit(true, false, false, true, address(lock));
        emit Withdrawn(beneficiary, 1000e18);
        vm.prank(beneficiary);
        lock.withdraw();

        assertTrue(lock.withdrawn());
        assertEq(lp.balanceOf(beneficiary), 1000e18);
        assertEq(lp.balanceOf(address(lock)), 0);
    }

    function test_withdraw_transfers_all_balance() public {
        // Mix of deposits from different senders + a direct transfer in. All of
        // it should be swept to the beneficiary.
        vm.prank(alice);
        lock.deposit(1000e18);
        vm.prank(bob);
        lock.deposit(2000e18);
        // Direct transfer (not through deposit) — still gets swept out.
        lp.mint(address(lock), 500e18);

        uint256 totalHeld = lp.balanceOf(address(lock));
        assertEq(totalHeld, 3500e18);

        vm.warp(lock.unlockTime());

        vm.expectEmit(true, false, false, true, address(lock));
        emit Withdrawn(beneficiary, totalHeld);
        vm.prank(beneficiary);
        lock.withdraw();

        assertEq(lp.balanceOf(beneficiary), totalHeld);
        assertEq(lp.balanceOf(address(lock)), 0);
        // `lockedAmount` tracks ingress only; the direct mint was not counted.
        assertEq(lock.lockedAmount(), 3000e18);
    }

    function test_withdraw_reverts_if_already_withdrawn() public {
        vm.prank(alice);
        lock.deposit(1000e18);

        vm.warp(lock.unlockTime());
        vm.prank(beneficiary);
        lock.withdraw();

        vm.prank(beneficiary);
        vm.expectRevert(LPLock.AlreadyWithdrawn.selector);
        lock.withdraw();
    }

    // ---------------------------------------------------------------------
    // Views
    // ---------------------------------------------------------------------

    function test_timeRemaining_reports_zero_after_unlock() public {
        uint256 unlock = lock.unlockTime();

        // Pre-unlock: full duration minus elapsed time.
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

    /// @dev Deposit a bounded sequence of amounts from alice and assert
    ///      `lockedAmount` tracks the running sum exactly, and the LP balance on
    ///      the lock matches. Individual amounts are trimmed to the supply cap and
    ///      the list is capped to avoid absurd fuzz runtimes.
    function testFuzz_deposit_multiple_times(uint128[] memory amounts) public {
        uint256 cap = amounts.length > 16 ? 16 : amounts.length;

        uint256 aliceStart = lp.balanceOf(alice);
        uint256 running;

        for (uint256 i = 0; i < cap; i++) {
            uint256 amt = uint256(amounts[i]);
            // Skip zero (reverts) and anything that exceeds remaining balance.
            uint256 remaining = lp.balanceOf(alice);
            if (amt == 0 || amt > remaining) {
                continue;
            }
            vm.prank(alice);
            lock.deposit(amt);
            running += amt;

            assertEq(lock.lockedAmount(), running);
            assertEq(lp.balanceOf(address(lock)), running);
            assertEq(lp.balanceOf(alice), aliceStart - running);
        }
    }
}
