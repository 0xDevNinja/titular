// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {LaunchRadarModule} from "../../src/launchpad/modules/LaunchRadarModule.sol";

/// @dev Unit coverage for the shared, per-agent launch-radar module.
///      Event expectations use {vm.expectEmit}; revert expectations use the
///      module's custom error selectors exclusively.
contract LaunchRadarModuleTest is Test {
    LaunchRadarModule internal radar;

    address internal protocolOwner = address(0xA0);
    address internal agentAdmin = address(0xAD);
    address internal otherAdmin = address(0xBAD);
    address internal agent = address(0xA6E17);
    address internal user1 = address(0xAAA1);
    address internal user2 = address(0xAAA2);

    // Events re-declared so {vm.expectEmit} can match by signature.
    event Configured(address indexed agent, uint64 startTime, bool whitelistOn, address agentAdmin);
    event WhitelistUpdated(address indexed agent, uint256 count);
    event StartTimeSet(address indexed agent, uint64 startTime);
    event WhitelistToggled(address indexed agent, bool on);

    function setUp() public {
        radar = new LaunchRadarModule(protocolOwner);
    }

    // ---------------------------------------------------------------------
    // configure
    // ---------------------------------------------------------------------

    function test_configure_only_owner() public {
        uint64 startAt = uint64(block.timestamp + 1 hours);

        // Non-owner is rejected.
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, address(this)));
        radar.configure(agent, startAt, false, agentAdmin);

        // Owner succeeds.
        vm.prank(protocolOwner);
        radar.configure(agent, startAt, false, agentAdmin);

        (uint64 gotStart, bool gotWl, bool gotConfigured) = radar.configs(agent);
        assertEq(gotStart, startAt);
        assertFalse(gotWl);
        assertTrue(gotConfigured);
        assertEq(radar.agentAdmin(agent), agentAdmin);
    }

    function test_configure_reverts_if_already_configured() public {
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
        emit Configured(agent, startAt, true, agentAdmin);

        vm.prank(protocolOwner);
        radar.configure(agent, startAt, true, agentAdmin);
    }

    function test_configure_reverts_on_zero_agent() public {
        vm.prank(protocolOwner);
        vm.expectRevert(LaunchRadarModule.ZeroAddress.selector);
        radar.configure(address(0), uint64(block.timestamp + 1), false, agentAdmin);
    }

    function test_configure_reverts_on_zero_agent_admin() public {
        vm.prank(protocolOwner);
        vm.expectRevert(LaunchRadarModule.ZeroAddress.selector);
        radar.configure(agent, uint64(block.timestamp + 1), false, address(0));
    }

    // ---------------------------------------------------------------------
    // setWhitelist
    // ---------------------------------------------------------------------

    function test_setWhitelist_only_agent_admin() public {
        _configureAgent({startAt: uint64(block.timestamp + 1 hours), whitelistOn: true});

        address[] memory users = new address[](1);
        bool[] memory allowed = new bool[](1);
        users[0] = user1;
        allowed[0] = true;

        // Caller that isn't the agent admin is rejected, even the protocol owner.
        vm.prank(protocolOwner);
        vm.expectRevert(LaunchRadarModule.NotAgentAdmin.selector);
        radar.setWhitelist(agent, users, allowed);

        vm.prank(otherAdmin);
        vm.expectRevert(LaunchRadarModule.NotAgentAdmin.selector);
        radar.setWhitelist(agent, users, allowed);

        // Agent admin succeeds.
        vm.prank(agentAdmin);
        radar.setWhitelist(agent, users, allowed);
        assertTrue(radar.whitelist(agent, user1));
    }

    function test_setWhitelist_batch_updates() public {
        _configureAgent({startAt: uint64(block.timestamp + 1 hours), whitelistOn: true});

        address[] memory users = new address[](3);
        bool[] memory allowed = new bool[](3);
        users[0] = user1;
        allowed[0] = true;
        users[1] = user2;
        allowed[1] = true;
        users[2] = address(0xCAFE);
        allowed[2] = false;

        vm.expectEmit(true, false, false, true, address(radar));
        emit WhitelistUpdated(agent, 3);

        vm.prank(agentAdmin);
        radar.setWhitelist(agent, users, allowed);

        assertTrue(radar.whitelist(agent, user1));
        assertTrue(radar.whitelist(agent, user2));
        assertFalse(radar.whitelist(agent, address(0xCAFE)));

        // Second call can flip a previously-allowed user back off.
        address[] memory users2 = new address[](1);
        bool[] memory allowed2 = new bool[](1);
        users2[0] = user1;
        allowed2[0] = false;
        vm.prank(agentAdmin);
        radar.setWhitelist(agent, users2, allowed2);
        assertFalse(radar.whitelist(agent, user1));
    }

    function test_setWhitelist_reverts_length_mismatch() public {
        _configureAgent({startAt: uint64(block.timestamp + 1 hours), whitelistOn: true});

        address[] memory users = new address[](2);
        bool[] memory allowed = new bool[](1);
        users[0] = user1;
        users[1] = user2;
        allowed[0] = true;

        vm.prank(agentAdmin);
        vm.expectRevert(LaunchRadarModule.LengthMismatch.selector);
        radar.setWhitelist(agent, users, allowed);
    }

    // ---------------------------------------------------------------------
    // setStartTime
    // ---------------------------------------------------------------------

    function test_setStartTime_only_agent_admin() public {
        uint64 startAt = uint64(block.timestamp + 1 hours);
        _configureAgent({startAt: startAt, whitelistOn: false});

        uint64 newStart = uint64(block.timestamp + 2 hours);

        vm.prank(otherAdmin);
        vm.expectRevert(LaunchRadarModule.NotAgentAdmin.selector);
        radar.setStartTime(agent, newStart);

        vm.expectEmit(true, false, false, true, address(radar));
        emit StartTimeSet(agent, newStart);

        vm.prank(agentAdmin);
        radar.setStartTime(agent, newStart);

        (uint64 gotStart,,) = radar.configs(agent);
        assertEq(gotStart, newStart);
    }

    function test_setStartTime_reverts_if_already_live() public {
        uint64 startAt = uint64(block.timestamp + 1 hours);
        _configureAgent({startAt: startAt, whitelistOn: false});

        // Advance past the start time: trading is now live and the gate is frozen.
        vm.warp(startAt);

        vm.prank(agentAdmin);
        vm.expectRevert(LaunchRadarModule.AlreadyLive.selector);
        radar.setStartTime(agent, uint64(block.timestamp + 1 days));

        // One second later — still considered live.
        vm.warp(startAt + 1);
        vm.prank(agentAdmin);
        vm.expectRevert(LaunchRadarModule.AlreadyLive.selector);
        radar.setStartTime(agent, uint64(block.timestamp + 1 days));
    }

    // ---------------------------------------------------------------------
    // setWhitelistOn
    // ---------------------------------------------------------------------

    function test_setWhitelistOn_toggle() public {
        _configureAgent({startAt: uint64(block.timestamp + 1 hours), whitelistOn: false});

        vm.prank(otherAdmin);
        vm.expectRevert(LaunchRadarModule.NotAgentAdmin.selector);
        radar.setWhitelistOn(agent, true);

        vm.expectEmit(true, false, false, true, address(radar));
        emit WhitelistToggled(agent, true);
        vm.prank(agentAdmin);
        radar.setWhitelistOn(agent, true);
        (, bool wlOn,) = radar.configs(agent);
        assertTrue(wlOn);

        vm.expectEmit(true, false, false, true, address(radar));
        emit WhitelistToggled(agent, false);
        vm.prank(agentAdmin);
        radar.setWhitelistOn(agent, false);
        (, wlOn,) = radar.configs(agent);
        assertFalse(wlOn);
    }

    // ---------------------------------------------------------------------
    // canTrade
    // ---------------------------------------------------------------------

    function test_canTrade_allowed_before_configure() public view {
        // An agent the radar has never seen is permissive. This is the documented
        // contract: integrators requiring positive auth must additionally check
        // `configs(agent).configured`.
        (bool ok, string memory reason) = radar.canTrade(agent, user1);
        assertTrue(ok);
        assertEq(bytes(reason).length, 0);
    }

    function test_canTrade_blocked_before_start() public {
        uint64 startAt = uint64(block.timestamp + 1 hours);
        _configureAgent({startAt: startAt, whitelistOn: false});

        (bool ok, string memory reason) = radar.canTrade(agent, user1);
        assertFalse(ok);
        assertEq(reason, "before start");
    }

    function test_canTrade_allowed_after_start() public {
        uint64 startAt = uint64(block.timestamp + 1 hours);
        _configureAgent({startAt: startAt, whitelistOn: false});

        vm.warp(startAt);
        (bool ok, string memory reason) = radar.canTrade(agent, user1);
        assertTrue(ok);
        assertEq(bytes(reason).length, 0);
    }

    function test_canTrade_blocked_not_whitelisted_when_on() public {
        uint64 startAt = uint64(block.timestamp + 1 hours);
        _configureAgent({startAt: startAt, whitelistOn: true});
        vm.warp(startAt);

        // user1 not on the list.
        (bool ok1, string memory reason1) = radar.canTrade(agent, user1);
        assertFalse(ok1);
        assertEq(reason1, "not whitelisted");

        // Add user2 to the allowlist; user1 still blocked, user2 now allowed.
        address[] memory users = new address[](1);
        bool[] memory allowed = new bool[](1);
        users[0] = user2;
        allowed[0] = true;
        vm.prank(agentAdmin);
        radar.setWhitelist(agent, users, allowed);

        (bool ok2,) = radar.canTrade(agent, user1);
        assertFalse(ok2);
        (bool ok3, string memory reason3) = radar.canTrade(agent, user2);
        assertTrue(ok3);
        assertEq(bytes(reason3).length, 0);
    }

    function test_canTrade_allowed_when_whitelist_off_regardless() public {
        uint64 startAt = uint64(block.timestamp + 1 hours);
        _configureAgent({startAt: startAt, whitelistOn: false});
        vm.warp(startAt);

        // Not on the (non-existent) whitelist — but whitelistOn is false, so any
        // user can trade.
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
