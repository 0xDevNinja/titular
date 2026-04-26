// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

import {LaunchRadarModule} from "../../src/launchpad/modules/LaunchRadarModule.sol";

/// @title  LaunchRadarModuleSuiteTest
/// @notice Per-contract unit tests for {LaunchRadarModule}: owner-gated
///         bootstrap, agent-admin-gated mutations, batch whitelist updates,
///         and `canTrade` permissive-by-default semantics.
contract LaunchRadarModuleSuiteTest is Test {
    LaunchRadarModule internal radar;

    address internal protocolOwner = address(0xA0);
    address internal agentAdmin = address(0xAD);
    address internal otherAdmin = address(0xBAD);
    address internal agent = address(0xA6E17);
    address internal user1 = address(0xAAA1);
    address internal user2 = address(0xAAA2);

    function setUp() public {
        radar = new LaunchRadarModule(protocolOwner);
    }

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    function test_constructor_zero_owner_reverts() public {
        vm.expectRevert(LaunchRadarModule.ZeroAddress.selector);
        new LaunchRadarModule(address(0));
    }

    // ---------------------------------------------------------------------
    // configure
    // ---------------------------------------------------------------------

    function test_configure_only_owner() public {
        uint64 startAt = uint64(block.timestamp + 1 hours);

        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, address(this)));
        radar.configure(agent, startAt, false, agentAdmin);

        vm.prank(protocolOwner);
        radar.configure(agent, startAt, false, agentAdmin);

        (uint64 gotStart, bool gotWl, bool gotConfigured) = radar.configs(agent);
        assertEq(gotStart, startAt);
        assertFalse(gotWl);
        assertTrue(gotConfigured);
        assertEq(radar.agentAdmin(agent), agentAdmin);
    }

    function test_configure_one_shot() public {
        uint64 startAt = uint64(block.timestamp + 1 hours);
        vm.prank(protocolOwner);
        radar.configure(agent, startAt, false, agentAdmin);

        vm.prank(protocolOwner);
        vm.expectRevert(LaunchRadarModule.AlreadyConfigured.selector);
        radar.configure(agent, startAt, true, otherAdmin);
    }

    function test_configure_emits_event() public {
        uint64 startAt = uint64(block.timestamp + 30 minutes);
        vm.expectEmit(true, false, false, true, address(radar));
        emit LaunchRadarModule.Configured(agent, startAt, true, agentAdmin);
        vm.prank(protocolOwner);
        radar.configure(agent, startAt, true, agentAdmin);
    }

    function test_configure_zero_agent_reverts() public {
        vm.prank(protocolOwner);
        vm.expectRevert(LaunchRadarModule.ZeroAddress.selector);
        radar.configure(address(0), uint64(block.timestamp + 1), false, agentAdmin);
    }

    function test_configure_zero_admin_reverts() public {
        vm.prank(protocolOwner);
        vm.expectRevert(LaunchRadarModule.ZeroAddress.selector);
        radar.configure(agent, uint64(block.timestamp + 1), false, address(0));
    }

    // ---------------------------------------------------------------------
    // setWhitelist
    // ---------------------------------------------------------------------

    function test_setWhitelist_only_admin() public {
        _configureAgent(uint64(block.timestamp + 1 hours), true);

        address[] memory u = new address[](1);
        bool[] memory ok = new bool[](1);
        u[0] = user1;
        ok[0] = true;

        // Even the protocol owner cannot bypass agent admin.
        vm.prank(protocolOwner);
        vm.expectRevert(LaunchRadarModule.NotAgentAdmin.selector);
        radar.setWhitelist(agent, u, ok);

        vm.prank(otherAdmin);
        vm.expectRevert(LaunchRadarModule.NotAgentAdmin.selector);
        radar.setWhitelist(agent, u, ok);

        vm.prank(agentAdmin);
        radar.setWhitelist(agent, u, ok);
        assertTrue(radar.whitelist(agent, user1));
    }

    function test_setWhitelist_batch_and_remove() public {
        _configureAgent(uint64(block.timestamp + 1 hours), true);

        address[] memory u = new address[](3);
        bool[] memory ok = new bool[](3);
        u[0] = user1;
        ok[0] = true;
        u[1] = user2;
        ok[1] = true;
        u[2] = address(0xCAFE);
        ok[2] = false;

        vm.expectEmit(true, false, false, true, address(radar));
        emit LaunchRadarModule.WhitelistUpdated(agent, 3);
        vm.prank(agentAdmin);
        radar.setWhitelist(agent, u, ok);

        assertTrue(radar.whitelist(agent, user1));
        assertTrue(radar.whitelist(agent, user2));

        // Flip user1 off in a follow-up batch.
        address[] memory u2 = new address[](1);
        bool[] memory ok2 = new bool[](1);
        u2[0] = user1;
        ok2[0] = false;
        vm.prank(agentAdmin);
        radar.setWhitelist(agent, u2, ok2);
        assertFalse(radar.whitelist(agent, user1));
    }

    function test_setWhitelist_length_mismatch_reverts() public {
        _configureAgent(uint64(block.timestamp + 1 hours), true);
        address[] memory u = new address[](2);
        bool[] memory ok = new bool[](1);
        vm.prank(agentAdmin);
        vm.expectRevert(LaunchRadarModule.LengthMismatch.selector);
        radar.setWhitelist(agent, u, ok);
    }

    // ---------------------------------------------------------------------
    // setStartTime
    // ---------------------------------------------------------------------

    function test_setStartTime_only_admin() public {
        uint64 startAt = uint64(block.timestamp + 1 hours);
        _configureAgent(startAt, false);

        uint64 newStart = uint64(block.timestamp + 2 hours);
        vm.prank(otherAdmin);
        vm.expectRevert(LaunchRadarModule.NotAgentAdmin.selector);
        radar.setStartTime(agent, newStart);

        vm.expectEmit(true, false, false, true, address(radar));
        emit LaunchRadarModule.StartTimeSet(agent, newStart);
        vm.prank(agentAdmin);
        radar.setStartTime(agent, newStart);
        (uint64 gotStart,,) = radar.configs(agent);
        assertEq(gotStart, newStart);
    }

    function test_setStartTime_post_launch_reverts() public {
        uint64 startAt = uint64(block.timestamp + 1 hours);
        _configureAgent(startAt, false);

        vm.warp(startAt);
        vm.prank(agentAdmin);
        vm.expectRevert(LaunchRadarModule.AlreadyLive.selector);
        radar.setStartTime(agent, uint64(block.timestamp + 1 days));
    }

    // ---------------------------------------------------------------------
    // setWhitelistOn
    // ---------------------------------------------------------------------

    function test_setWhitelistOn_toggle() public {
        _configureAgent(uint64(block.timestamp + 1 hours), false);

        vm.prank(otherAdmin);
        vm.expectRevert(LaunchRadarModule.NotAgentAdmin.selector);
        radar.setWhitelistOn(agent, true);

        vm.expectEmit(true, false, false, true, address(radar));
        emit LaunchRadarModule.WhitelistToggled(agent, true);
        vm.prank(agentAdmin);
        radar.setWhitelistOn(agent, true);
        (, bool wlOn,) = radar.configs(agent);
        assertTrue(wlOn);

        vm.prank(agentAdmin);
        radar.setWhitelistOn(agent, false);
        (, wlOn,) = radar.configs(agent);
        assertFalse(wlOn);
    }

    // ---------------------------------------------------------------------
    // canTrade — every branch
    // ---------------------------------------------------------------------

    function test_canTrade_unconfigured_permissive() public view {
        (bool ok, string memory reason) = radar.canTrade(agent, user1);
        assertTrue(ok);
        assertEq(bytes(reason).length, 0);
    }

    function test_canTrade_pre_start_blocks() public {
        _configureAgent(uint64(block.timestamp + 1 hours), false);
        (bool ok, string memory reason) = radar.canTrade(agent, user1);
        assertFalse(ok);
        assertEq(reason, "before start");
    }

    function test_canTrade_post_start_allows() public {
        uint64 startAt = uint64(block.timestamp + 1 hours);
        _configureAgent(startAt, false);
        vm.warp(startAt);
        (bool ok,) = radar.canTrade(agent, user1);
        assertTrue(ok);
    }

    function test_canTrade_whitelist_blocks_non_member() public {
        uint64 startAt = uint64(block.timestamp + 1 hours);
        _configureAgent(startAt, true);
        vm.warp(startAt);

        (bool ok, string memory reason) = radar.canTrade(agent, user1);
        assertFalse(ok);
        assertEq(reason, "not whitelisted");

        address[] memory u = new address[](1);
        bool[] memory ok_ = new bool[](1);
        u[0] = user2;
        ok_[0] = true;
        vm.prank(agentAdmin);
        radar.setWhitelist(agent, u, ok_);

        (bool ok2,) = radar.canTrade(agent, user1);
        assertFalse(ok2);
        (bool ok3,) = radar.canTrade(agent, user2);
        assertTrue(ok3);
    }

    function test_canTrade_whitelist_off_allows_anyone() public {
        uint64 startAt = uint64(block.timestamp + 1 hours);
        _configureAgent(startAt, false);
        vm.warp(startAt);
        (bool ok, string memory reason) = radar.canTrade(agent, user1);
        assertTrue(ok);
        assertEq(bytes(reason).length, 0);
    }

    // ---------------------------------------------------------------------
    // Helpers
    // ---------------------------------------------------------------------

    function _configureAgent(uint64 startAt, bool whitelistOn) internal {
        vm.prank(protocolOwner);
        radar.configure(agent, startAt, whitelistOn, agentAdmin);
    }
}
