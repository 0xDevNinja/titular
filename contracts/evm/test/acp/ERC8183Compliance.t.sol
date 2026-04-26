// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {AgentRegistry} from "../../src/acp/AgentRegistry.sol";
import {IERC8183} from "../../src/acp/IERC8183.sol";

/// @notice ERC-8183 interface compliance check.
///         Verifies that AgentRegistry satisfies all mandatory IERC8183 function selectors.
contract ERC8183ComplianceTest is Test {
    AgentRegistry internal registry;
    address internal admin = makeAddr("admin");
    address internal scorer = makeAddr("scorer");

    function setUp() public {
        registry = new AgentRegistry(admin);
        bytes32 scorerRole = registry.SCORER_ROLE();
        vm.prank(admin);
        registry.grantRole(scorerRole, scorer);
    }

    // ─────────────────────────────────────────────────────────────
    // Interface checks via supportsInterface-equivalent selector tests
    // ─────────────────────────────────────────────────────────────

    /// @notice `register` exists and returns an agentId.
    function test_erc8183_register() public {
        address ctrl = makeAddr("ctrl");
        vm.prank(ctrl);
        uint256 id = registry.register(ctrl, "ipfs://QmAgent", 0x1);
        assertEq(id, 0);
        assertEq(registry.totalAgents(), 1);
    }

    /// @notice `agentInfo` returns correct packed data (implemented via getAgent mapping).
    function test_erc8183_agentInfo_viaGetAgent() public {
        address ctrl = makeAddr("ctrl");
        vm.prank(ctrl);
        uint256 id = registry.register(ctrl, "ipfs://QmAgent", 0xFF);

        AgentRegistry.AgentInfo memory info = registry.getAgent(id);
        assertEq(info.controller, ctrl);
        assertEq(info.metadataURI, "ipfs://QmAgent");
        assertEq(info.capabilities, 0xFF);
        assertTrue(info.active);
    }

    /// @notice `totalAgents` increments after each registration.
    function test_erc8183_totalAgents() public {
        assertEq(registry.totalAgents(), 0);
        vm.prank(makeAddr("a")); registry.register(makeAddr("a"), "ipfs://A", 0);
        vm.prank(makeAddr("b")); registry.register(makeAddr("b"), "ipfs://B", 0);
        assertEq(registry.totalAgents(), 2);
    }

    /// @notice `postScore` / `reputationScore` round-trip.
    function test_erc8183_reputationRoundTrip() public {
        address ctrl = makeAddr("ctrl");
        vm.prank(ctrl);
        uint256 id = registry.register(ctrl, "ipfs://QmAgent", 0x1);

        vm.prank(scorer);
        registry.postScore(id, 42, 0);
        assertEq(registry.getAgent(id).reputationScore, 42);
    }

    // ─────────────────────────────────────────────────────────────
    // ABI publish: verify key function selectors are stable
    // ─────────────────────────────────────────────────────────────

    function test_selector_register() public pure {
        bytes4 expected = bytes4(keccak256("register(address,string,uint256)"));
        assertEq(AgentRegistry.register.selector, expected);
    }

    function test_selector_postScore() public pure {
        bytes4 expected = bytes4(keccak256("postScore(uint256,int256,uint256)"));
        assertEq(AgentRegistry.postScore.selector, expected);
    }

    function test_selector_getAgent() public pure {
        bytes4 expected = bytes4(keccak256("getAgent(uint256)"));
        assertEq(AgentRegistry.getAgent.selector, expected);
    }

    function test_selector_totalAgents() public pure {
        bytes4 expected = bytes4(keccak256("totalAgents()"));
        assertEq(AgentRegistry.totalAgents.selector, expected);
    }
}
