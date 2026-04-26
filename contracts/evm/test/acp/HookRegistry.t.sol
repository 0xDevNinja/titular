// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {HookRegistry} from "../../src/acp/HookRegistry.sol";
import {IHook} from "../../src/acp/IHook.sol";

/// @dev Minimal hook stub implementing IHook.
contract MockHook is IHook {
    string private _name;

    constructor(string memory name_) {
        _name = name_;
    }

    function onAccept(uint256, bytes calldata) external pure override {}
    function onSubmit(uint256, bytes calldata) external pure override {}
    function onApprove(uint256, bytes calldata) external pure override {}
    function onReject(uint256, bytes calldata) external pure override {}
    function onCancel(uint256, bytes calldata) external pure override {}

    function hookName() external view override returns (string memory) {
        return _name;
    }
}

contract HookRegistryTest is Test {
    HookRegistry internal registry;
    MockHook internal hookA;
    MockHook internal hookB;

    address internal admin = makeAddr("admin");
    address internal stranger = makeAddr("stranger");

    event HookRegistered(address indexed hook, string name);
    event HookDeregistered(address indexed hook);

    function setUp() public {
        registry = new HookRegistry(admin);
        hookA = new MockHook("FundTransferHook");
        hookB = new MockHook("MilestoneHook");
    }

    // ─────────────────────────────────────────────────────────────
    // Constructor
    // ─────────────────────────────────────────────────────────────

    function test_constructor_revert_zeroAdmin() public {
        vm.expectRevert(HookRegistry.ZeroAddress.selector);
        new HookRegistry(address(0));
    }

    // ─────────────────────────────────────────────────────────────
    // registerHook
    // ─────────────────────────────────────────────────────────────

    function test_registerHook_success() public {
        vm.expectEmit(true, false, false, true);
        emit HookRegistered(address(hookA), "FundTransferHook");

        vm.prank(admin);
        registry.registerHook(address(hookA));

        assertTrue(registry.isApproved(address(hookA)));
        HookRegistry.HookInfo memory info = registry.getHookInfo(address(hookA));
        assertEq(info.name, "FundTransferHook");
        assertTrue(info.approved);
    }

    function test_registerHook_multipleHooks() public {
        vm.prank(admin);
        registry.registerHook(address(hookA));
        vm.prank(admin);
        registry.registerHook(address(hookB));

        address[] memory hooks = registry.allHooks();
        assertEq(hooks.length, 2);
        assertEq(hooks[0], address(hookA));
        assertEq(hooks[1], address(hookB));
    }

    function test_registerHook_revert_notAdmin() public {
        vm.expectRevert();
        vm.prank(stranger);
        registry.registerHook(address(hookA));
    }

    function test_registerHook_revert_zeroAddress() public {
        vm.expectRevert(HookRegistry.ZeroAddress.selector);
        vm.prank(admin);
        registry.registerHook(address(0));
    }

    function test_registerHook_revert_alreadyRegistered() public {
        vm.prank(admin);
        registry.registerHook(address(hookA));

        vm.expectRevert(abi.encodeWithSelector(HookRegistry.AlreadyRegistered.selector, address(hookA)));
        vm.prank(admin);
        registry.registerHook(address(hookA));
    }

    // ─────────────────────────────────────────────────────────────
    // deregisterHook
    // ─────────────────────────────────────────────────────────────

    function test_deregisterHook_success() public {
        vm.prank(admin);
        registry.registerHook(address(hookA));

        vm.expectEmit(true, false, false, false);
        emit HookDeregistered(address(hookA));

        vm.prank(admin);
        registry.deregisterHook(address(hookA));

        assertFalse(registry.isApproved(address(hookA)));
        // Name is retained
        HookRegistry.HookInfo memory info = registry.getHookInfo(address(hookA));
        assertEq(info.name, "FundTransferHook");
        assertFalse(info.approved);
    }

    function test_deregisterHook_revert_notRegistered() public {
        vm.expectRevert(abi.encodeWithSelector(HookRegistry.NotRegistered.selector, address(hookA)));
        vm.prank(admin);
        registry.deregisterHook(address(hookA));
    }

    function test_deregisterHook_revert_alreadyRevoked() public {
        vm.prank(admin);
        registry.registerHook(address(hookA));
        vm.prank(admin);
        registry.deregisterHook(address(hookA));

        vm.expectRevert(abi.encodeWithSelector(HookRegistry.NotRegistered.selector, address(hookA)));
        vm.prank(admin);
        registry.deregisterHook(address(hookA));
    }

    function test_deregisterHook_revert_notAdmin() public {
        vm.prank(admin);
        registry.registerHook(address(hookA));

        vm.expectRevert();
        vm.prank(stranger);
        registry.deregisterHook(address(hookA));
    }

    // ─────────────────────────────────────────────────────────────
    // View: allHooks includes deregistered
    // ─────────────────────────────────────────────────────────────

    function test_allHooks_includesRevoked() public {
        vm.prank(admin);
        registry.registerHook(address(hookA));
        vm.prank(admin);
        registry.registerHook(address(hookB));
        vm.prank(admin);
        registry.deregisterHook(address(hookA));

        address[] memory hooks = registry.allHooks();
        assertEq(hooks.length, 2); // both still listed
        assertFalse(registry.isApproved(hooks[0]));
        assertTrue(registry.isApproved(hooks[1]));
    }

    // ─────────────────────────────────────────────────────────────
    // isApproved returns false for unknown hook
    // ─────────────────────────────────────────────────────────────

    function test_isApproved_unknownHook() public {
        address unknown = makeAddr("unknown");
        assertFalse(registry.isApproved(unknown));
    }
}
