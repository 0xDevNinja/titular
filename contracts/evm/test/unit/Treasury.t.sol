// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {Treasury} from "../../src/treasury/Treasury.sol";

contract MockERC20 is ERC20 {
    constructor() ERC20("MCK", "MCK") {
        _mint(msg.sender, 1_000_000e18);
    }
}

contract TreasuryV2 is Treasury {
    function pong() external pure returns (uint256) {
        return 42;
    }
}

contract TreasuryTest is Test {
    Treasury internal treasury;
    MockERC20 internal token;
    address internal safe = address(0x5AFE);
    address internal fd = address(0xFEE);
    address internal stranger = address(0xBAD);

    event Withdrawn(address indexed token, address indexed to, uint256 amount);
    event StreamedToVe(address indexed token, uint256 amount);
    event FeeDistributorSet(address indexed previous, address indexed next);

    function setUp() public {
        Treasury impl = new Treasury();
        bytes memory initData = abi.encodeCall(Treasury.initialize, (safe, fd));
        ERC1967Proxy proxy = new ERC1967Proxy(address(impl), initData);
        treasury = Treasury(payable(address(proxy)));
        token = new MockERC20();
        token.transfer(address(treasury), 10_000e18);
    }

    function test_initialize_sets_owner_and_feeDistributor() public view {
        assertEq(treasury.owner(), safe);
        assertEq(treasury.feeDistributor(), fd);
    }

    function test_initialize_revert_zero_addresses() public {
        Treasury impl = new Treasury();
        bytes memory initData = abi.encodeCall(Treasury.initialize, (address(0), fd));
        vm.expectRevert(Treasury.ZeroAddress.selector);
        new ERC1967Proxy(address(impl), initData);
    }

    function test_initialize_cannot_be_called_twice() public {
        vm.expectRevert();
        treasury.initialize(safe, fd);
    }

    function test_withdraw_erc20_owner_only() public {
        vm.prank(safe);
        vm.expectEmit(true, true, false, true);
        emit Withdrawn(address(token), safe, 1000e18);
        treasury.withdraw(address(token), safe, 1000e18);
        assertEq(token.balanceOf(safe), 1000e18);
    }

    function test_withdraw_native_owner_only() public {
        vm.deal(address(treasury), 5 ether);
        vm.prank(safe);
        treasury.withdraw(address(0), safe, 2 ether);
        assertEq(safe.balance, 2 ether);
    }

    function test_withdraw_reverts_not_owner() public {
        vm.prank(stranger);
        vm.expectRevert();
        treasury.withdraw(address(token), stranger, 1);
    }

    function test_withdraw_reverts_zero_amount() public {
        vm.prank(safe);
        vm.expectRevert(Treasury.ZeroAmount.selector);
        treasury.withdraw(address(token), safe, 0);
    }

    function test_receive_native_ok() public {
        (bool ok,) = address(treasury).call{value: 3 ether}("");
        assertTrue(ok);
        assertEq(address(treasury).balance, 3 ether);
    }

    function test_streamToVe_transfers_to_feeDistributor() public {
        vm.prank(safe);
        vm.expectEmit(true, false, false, true);
        emit StreamedToVe(address(token), 500e18);
        treasury.streamToVe(address(token), 500e18);
        assertEq(token.balanceOf(fd), 500e18);
    }

    function test_streamToVe_reverts_zero_token() public {
        vm.prank(safe);
        vm.expectRevert(Treasury.ZeroAddress.selector);
        treasury.streamToVe(address(0), 1);
    }

    function test_setFeeDistributor_emits_and_updates() public {
        address next = address(0xFEED2);
        vm.prank(safe);
        vm.expectEmit(true, true, false, true);
        emit FeeDistributorSet(fd, next);
        treasury.setFeeDistributor(next);
        assertEq(treasury.feeDistributor(), next);
    }

    function test_setFeeDistributor_reverts_zero() public {
        vm.prank(safe);
        vm.expectRevert(Treasury.ZeroAddress.selector);
        treasury.setFeeDistributor(address(0));
    }

    function test_upgrade_authorized_owner_only() public {
        TreasuryV2 v2 = new TreasuryV2();
        vm.prank(stranger);
        vm.expectRevert();
        treasury.upgradeToAndCall(address(v2), "");

        vm.prank(safe);
        treasury.upgradeToAndCall(address(v2), "");
        assertEq(TreasuryV2(payable(address(treasury))).pong(), 42);
    }
}
