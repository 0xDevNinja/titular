// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {AgentRegistry} from "../../../src/acp/AgentRegistry.sol";

/// @notice Reputation integration tests.
///         Verifies bump-on-accept patterns and decrement-on-reject semantics
///         using AgentRegistry.postScore.
contract ReputationTest is Test {
    AgentRegistry internal registry;

    address internal admin = makeAddr("admin");
    address internal scorer = makeAddr("scorer");
    address internal agentCtrl = makeAddr("agentCtrl");
    address internal stranger = makeAddr("stranger");

    uint256 internal agentId;

    function setUp() public {
        registry = new AgentRegistry(admin);
        bytes32 scorerRole = registry.SCORER_ROLE();
        vm.prank(admin);
        registry.grantRole(scorerRole, scorer);

        vm.prank(agentCtrl);
        agentId = registry.register(agentCtrl, "ipfs://QmAgent", 0x1);
    }

    // ─── Bump on accept (positive delta) ────────────────────────

    function test_bumpOnAccept_positiveScore() public {
        vm.prank(scorer);
        registry.postScore(agentId, 10, 0);
        assertEq(registry.getAgent(agentId).reputationScore, 10);
    }

    function test_bumpOnAccept_accumulates() public {
        vm.prank(scorer);
        registry.postScore(agentId, 10, 0);
        vm.prank(scorer);
        registry.postScore(agentId, 5, 1);
        assertEq(registry.getAgent(agentId).reputationScore, 15);
    }

    // ─── Decrement on reject (negative delta) ───────────────────

    function test_decrementOnReject() public {
        // Start positive
        vm.prank(scorer);
        registry.postScore(agentId, 20, 0);

        // Reject reduces score
        vm.prank(scorer);
        registry.postScore(agentId, -8, 1);
        assertEq(registry.getAgent(agentId).reputationScore, 12);
    }

    function test_decrementToNegative() public {
        vm.prank(scorer);
        registry.postScore(agentId, -5, 0);
        assertEq(registry.getAgent(agentId).reputationScore, -5);
    }

    // ─── Multiple agents ─────────────────────────────────────────

    function test_reputationIsolatedPerAgent() public {
        address ctrl2 = makeAddr("ctrl2");
        vm.prank(ctrl2);
        uint256 agentId2 = registry.register(ctrl2, "ipfs://QmAgent2", 0x2);

        vm.prank(scorer);
        registry.postScore(agentId, 100, 0);
        vm.prank(scorer);
        registry.postScore(agentId2, -50, 0);

        assertEq(registry.getAgent(agentId).reputationScore, 100);
        assertEq(registry.getAgent(agentId2).reputationScore, -50);
    }

    // ─── Nonce replay protection ─────────────────────────────────

    function test_replayProtection_revertOnSameNonce() public {
        vm.prank(scorer);
        registry.postScore(agentId, 10, 0);

        vm.expectRevert(abi.encodeWithSelector(AgentRegistry.InvalidNonce.selector, agentId, 1, 0));
        vm.prank(scorer);
        registry.postScore(agentId, 10, 0);
    }

    function test_nonceIncrementsCorrectly() public {
        for (uint256 i = 0; i < 5; ++i) {
            vm.prank(scorer);
            registry.postScore(agentId, 1, i);
        }
        assertEq(registry.scorerNonce(agentId), 5);
        assertEq(registry.getAgent(agentId).reputationScore, 5);
    }

    // ─── Inactive agent cannot receive score ────────────────────

    function test_inactiveAgent_cannotReceiveScore() public {
        vm.prank(agentCtrl);
        registry.setActive(agentId, false);

        vm.expectRevert(abi.encodeWithSelector(AgentRegistry.AgentInactive.selector, agentId));
        vm.prank(scorer);
        registry.postScore(agentId, 1, 0);
    }

    // ─── Only scorer role ────────────────────────────────────────

    function test_onlyScorer_canPostScore() public {
        vm.expectRevert();
        vm.prank(stranger);
        registry.postScore(agentId, 5, 0);
    }

    // ─── Fuzz: net score equals sum of deltas ────────────────────

    function testFuzz_reputation_netScore(int256 d1, int256 d2, int256 d3) public {
        d1 = bound(d1, type(int64).min, type(int64).max);
        d2 = bound(d2, type(int64).min, type(int64).max);
        d3 = bound(d3, type(int64).min, type(int64).max);

        vm.prank(scorer);
        registry.postScore(agentId, d1, 0);
        vm.prank(scorer);
        registry.postScore(agentId, d2, 1);
        vm.prank(scorer);
        registry.postScore(agentId, d3, 2);

        assertEq(registry.getAgent(agentId).reputationScore, d1 + d2 + d3);
    }
}
