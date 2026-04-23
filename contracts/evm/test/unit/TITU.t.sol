// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {TITU} from "../../src/token/TITU.sol";

contract TITUTest is Test {
    TITU internal titu;
    address internal minter = address(0xBEEF);
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);

    function setUp() public {
        titu = new TITU(minter);
    }

    function test_totalSupply_is_1B() public view {
        assertEq(titu.totalSupply(), 1_000_000_000 * 1e18);
    }

    function test_initialMinter_balance_full_supply() public view {
        assertEq(titu.balanceOf(minter), 1_000_000_000 * 1e18);
    }

    function test_constructor_revert_zeroAddress() public {
        vm.expectRevert(TITU.ZeroAddress.selector);
        new TITU(address(0));
    }

    function test_clock_returns_timestamp() public {
        vm.warp(1_700_000_000);
        assertEq(uint256(titu.clock()), 1_700_000_000);
        assertEq(titu.CLOCK_MODE(), "mode=timestamp");
    }

    function test_burn_reduces_supply() public {
        vm.prank(minter);
        titu.transfer(alice, 1000e18);
        vm.prank(alice);
        titu.burn(400e18);
        assertEq(titu.totalSupply(), 1_000_000_000 * 1e18 - 400e18);
        assertEq(titu.balanceOf(alice), 600e18);
    }

    function test_transfer_updates_checkpoints() public {
        vm.prank(minter);
        titu.delegate(minter);
        vm.prank(minter);
        titu.transfer(alice, 1000e18);
        vm.prank(alice);
        titu.delegate(alice);
        assertEq(titu.getVotes(alice), 1000e18);
    }

    function test_votes_delegate_then_checkpoint() public {
        vm.prank(minter);
        titu.delegate(minter);
        uint256 minted = 1_000_000_000 * 1e18;
        assertEq(titu.getVotes(minter), minted);
        vm.warp(block.timestamp + 10);
        // past lookup valid at an earlier timestamp
        assertEq(titu.getPastVotes(minter, block.timestamp - 1), minted);
    }

    function test_permit_signature_flow() public {
        uint256 ownerPk = 0xA11CE1;
        address owner = vm.addr(ownerPk);
        vm.prank(minter);
        titu.transfer(owner, 500e18);

        uint256 value = 100e18;
        uint256 deadline = block.timestamp + 1 hours;
        bytes32 structHash = keccak256(
            abi.encode(
                keccak256("Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)"),
                owner,
                bob,
                value,
                titu.nonces(owner),
                deadline
            )
        );
        bytes32 digest = keccak256(abi.encodePacked("\x19\x01", titu.DOMAIN_SEPARATOR(), structHash));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(ownerPk, digest);
        titu.permit(owner, bob, value, deadline, v, r, s);
        assertEq(titu.allowance(owner, bob), value);
    }

    function testFuzz_transfer_preserves_total_supply(uint96 amt, address to) public {
        vm.assume(to != address(0) && to != minter);
        vm.assume(amt <= 1_000_000_000 * 1e18);
        vm.prank(minter);
        titu.transfer(to, amt);
        assertEq(titu.totalSupply(), 1_000_000_000 * 1e18);
    }
}
