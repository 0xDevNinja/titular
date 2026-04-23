// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {VestingVault} from "../../src/vault/VestingVault.sol";

contract MockERC20 is ERC20 {
    constructor() ERC20("MCK", "MCK") {
        _mint(msg.sender, 100_000_000e18);
    }
}

contract VestingVaultTest is Test {
    MockERC20 internal token;
    VestingVault internal vault;
    address internal admin = address(this);
    address internal beneficiary = address(0xB0B);
    address internal stranger = address(0xBAD);

    uint64 internal constant START = 1_700_000_000;
    uint64 internal constant CLIFF = 365 days;
    uint64 internal constant DURATION = 4 * 365 days;

    function setUp() public {
        token = new MockERC20();
        vault = new VestingVault(IERC20(address(token)), admin);
        token.approve(address(vault), type(uint256).max);
        vm.warp(START);
    }

    function test_addGrant_requires_admin() public {
        vm.prank(stranger);
        vm.expectRevert();
        vault.addGrant(beneficiary, 1000e18, START, CLIFF, DURATION);
    }

    function test_addGrant_pulls_tokens() public {
        vault.addGrant(beneficiary, 1000e18, START, CLIFF, DURATION);
        assertEq(token.balanceOf(address(vault)), 1000e18);
    }

    function test_addGrant_reverts_existing() public {
        vault.addGrant(beneficiary, 1000e18, START, CLIFF, DURATION);
        vm.expectRevert(VestingVault.GrantExists.selector);
        vault.addGrant(beneficiary, 1, START, CLIFF, DURATION);
    }

    function test_addGrant_reverts_zero_duration() public {
        vm.expectRevert(VestingVault.ZeroDuration.selector);
        vault.addGrant(beneficiary, 1000e18, START, 0, 0);
    }

    function test_addGrant_reverts_cliff_gt_duration() public {
        vm.expectRevert(VestingVault.CliffExceedsDuration.selector);
        vault.addGrant(beneficiary, 1000e18, START, DURATION + 1, DURATION);
    }

    function test_vested_before_cliff_returns_zero() public {
        vault.addGrant(beneficiary, 1000e18, START, CLIFF, DURATION);
        vm.warp(START + CLIFF - 1);
        assertEq(vault.vested(beneficiary), 0);
    }

    function test_vested_linear_after_cliff() public {
        vault.addGrant(beneficiary, 1000e18, START, CLIFF, DURATION);
        vm.warp(START + 2 * 365 days);
        assertEq(vault.vested(beneficiary), 500e18);
    }

    function test_vested_full_after_duration() public {
        vault.addGrant(beneficiary, 1000e18, START, CLIFF, DURATION);
        vm.warp(START + DURATION);
        assertEq(vault.vested(beneficiary), 1000e18);
    }

    function test_release_transfers_and_updates() public {
        vault.addGrant(beneficiary, 1000e18, START, CLIFF, DURATION);
        vm.warp(START + 2 * 365 days);
        vault.release(beneficiary);
        assertEq(token.balanceOf(beneficiary), 500e18);
        // second release with no time passed reverts
        vm.expectRevert(VestingVault.NothingToRelease.selector);
        vault.release(beneficiary);
    }

    function test_revoke_returns_unvested() public {
        vault.addGrant(beneficiary, 1000e18, START, CLIFF, DURATION);
        vm.warp(START + 2 * 365 days);
        uint256 adminBefore = token.balanceOf(admin);
        vault.revoke(beneficiary);
        assertEq(token.balanceOf(beneficiary), 500e18);
        assertEq(token.balanceOf(admin) - adminBefore, 500e18);
    }

    function test_release_self() public {
        vault.addGrant(beneficiary, 1000e18, START, CLIFF, DURATION);
        vm.warp(START + 2 * 365 days);
        vm.prank(beneficiary);
        vault.release();
        assertEq(token.balanceOf(beneficiary), 500e18);
    }

    function test_release_reverts_no_grant() public {
        vm.expectRevert(VestingVault.NoGrant.selector);
        vault.release(beneficiary);
    }

    function test_revoke_reverts_no_grant() public {
        vm.expectRevert(VestingVault.NoGrant.selector);
        vault.revoke(beneficiary);
    }

    function test_releasable_view() public {
        vault.addGrant(beneficiary, 1000e18, START, CLIFF, DURATION);
        vm.warp(START + 2 * 365 days);
        assertEq(vault.releasable(beneficiary), 500e18);
    }

    function test_constructor_reverts_zero() public {
        vm.expectRevert(VestingVault.ZeroAddress.selector);
        new VestingVault(IERC20(address(0)), admin);
        vm.expectRevert(VestingVault.ZeroAddress.selector);
        new VestingVault(IERC20(address(token)), address(0));
    }

    function test_addGrant_reverts_zero_beneficiary() public {
        vm.expectRevert(VestingVault.ZeroAddress.selector);
        vault.addGrant(address(0), 1000e18, START, CLIFF, DURATION);
    }

    function test_addGrant_reverts_zero_amount() public {
        vm.expectRevert(VestingVault.ZeroAmount.selector);
        vault.addGrant(beneficiary, 0, START, CLIFF, DURATION);
    }

    function testFuzz_release_never_exceeds_total(uint64 elapsed) public {
        elapsed = uint64(bound(elapsed, 0, 10 * 365 days));
        vault.addGrant(beneficiary, 1000e18, START, CLIFF, DURATION);
        vm.warp(uint256(START) + uint256(elapsed));
        if (uint256(elapsed) >= CLIFF) {
            try vault.release(beneficiary) {} catch {}
        }
        assertLe(token.balanceOf(beneficiary), 1000e18);
    }
}
