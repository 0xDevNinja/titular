// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {AgentRegistry} from "../../src/acp/AgentRegistry.sol";

contract AgentRegistryTest is Test {
    AgentRegistry internal registry;

    address internal admin = makeAddr("admin");
    address internal scorer = makeAddr("scorer");
    address internal alice = makeAddr("alice");
    address internal bob = makeAddr("bob");
    address internal charlie = makeAddr("charlie");

    string internal constant URI = "ipfs://QmTest";
    uint256 internal constant CAPS = 0x1;

    event AgentRegistered(uint256 indexed agentId, address indexed controller, string metadataURI, uint256 capabilities);
    event MetadataUpdated(uint256 indexed agentId, string metadataURI);
    event CapabilitiesUpdated(uint256 indexed agentId, uint256 capabilities);
    event ActiveStatusChanged(uint256 indexed agentId, bool active);
    event ControllerTransferProposed(uint256 indexed agentId, address indexed proposed);
    event ControllerTransferAccepted(uint256 indexed agentId, address indexed oldController, address indexed newController);
    event ControllerTransferCancelled(uint256 indexed agentId);
    event ScorePosted(uint256 indexed agentId, address indexed scorer, int256 delta, int256 newTotal, uint256 nonce);

    function setUp() public {
        registry = new AgentRegistry(admin);
        bytes32 scorerRole = registry.SCORER_ROLE();
        vm.prank(admin);
        registry.grantRole(scorerRole, scorer);
    }

    // ─── Registration ────────────────────────────────────────────

    function test_register_success() public {
        vm.expectEmit(true, true, false, true);
        emit AgentRegistered(0, alice, URI, CAPS);

        vm.prank(alice);
        uint256 id = registry.register(alice, URI, CAPS);
        assertEq(id, 0);

        AgentRegistry.AgentInfo memory info = registry.getAgent(0);
        assertEq(info.controller, alice);
        assertEq(info.metadataURI, URI);
        assertEq(info.capabilities, CAPS);
        assertEq(info.reputationScore, 0);
        assertTrue(info.active);
        assertEq(info.registeredAt, uint48(block.timestamp));
    }

    function test_register_incrementsId() public {
        vm.prank(alice);
        uint256 id0 = registry.register(alice, URI, CAPS);
        vm.prank(bob);
        uint256 id1 = registry.register(bob, "ipfs://QmB", 0x2);
        assertEq(id0, 0);
        assertEq(id1, 1);
        assertEq(registry.totalAgents(), 2);
    }

    function test_register_revert_zeroController() public {
        vm.expectRevert(AgentRegistry.ZeroAddress.selector);
        registry.register(address(0), URI, CAPS);
    }

    function test_register_revert_emptyURI() public {
        vm.expectRevert(AgentRegistry.EmptyMetadataURI.selector);
        registry.register(alice, "", CAPS);
    }

    function test_register_revert_whenPaused() public {
        vm.prank(admin);
        registry.pause();
        vm.expectRevert();
        registry.register(alice, URI, CAPS);
    }

    // ─── Metadata / capabilities ──────────────────────────────────

    function test_setMetadata_success() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);

        vm.expectEmit(true, false, false, true);
        emit MetadataUpdated(0, "ipfs://QmNew");

        vm.prank(alice);
        registry.setMetadata(0, "ipfs://QmNew");
        assertEq(registry.getAgent(0).metadataURI, "ipfs://QmNew");
    }

    function test_setMetadata_revert_notController() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);

        vm.expectRevert(abi.encodeWithSelector(AgentRegistry.NotController.selector, 0, bob));
        vm.prank(bob);
        registry.setMetadata(0, "ipfs://QmNew");
    }

    function test_setMetadata_revert_emptyURI() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);

        vm.expectRevert(AgentRegistry.EmptyMetadataURI.selector);
        vm.prank(alice);
        registry.setMetadata(0, "");
    }

    function test_setCapabilities_success() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);

        vm.expectEmit(true, false, false, true);
        emit CapabilitiesUpdated(0, 0xFF);

        vm.prank(alice);
        registry.setCapabilities(0, 0xFF);
        assertEq(registry.getAgent(0).capabilities, 0xFF);
    }

    function test_setCapabilities_revert_notController() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);

        vm.expectRevert(abi.encodeWithSelector(AgentRegistry.NotController.selector, 0, bob));
        vm.prank(bob);
        registry.setCapabilities(0, 0xFF);
    }

    // ─── Active status ───────────────────────────────────────────

    function test_setActive_deactivate() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);

        vm.expectEmit(true, false, false, true);
        emit ActiveStatusChanged(0, false);

        vm.prank(alice);
        registry.setActive(0, false);
        assertFalse(registry.getAgent(0).active);
    }

    function test_setActive_revert_notController() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);

        vm.expectRevert(abi.encodeWithSelector(AgentRegistry.NotController.selector, 0, bob));
        vm.prank(bob);
        registry.setActive(0, false);
    }

    // ─── Two-step controller transfer ────────────────────────────

    function test_controllerTransfer_happyPath() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);

        vm.prank(alice);
        registry.proposeControllerTransfer(0, bob);
        assertEq(registry.pendingController(0), bob);

        vm.expectEmit(true, true, true, false);
        emit ControllerTransferAccepted(0, alice, bob);

        vm.prank(bob);
        registry.acceptControllerTransfer(0);

        assertEq(registry.getAgent(0).controller, bob);
        assertEq(registry.pendingController(0), address(0));
    }

    function test_proposeTransfer_revert_notController() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);

        vm.expectRevert(abi.encodeWithSelector(AgentRegistry.NotController.selector, 0, bob));
        vm.prank(bob);
        registry.proposeControllerTransfer(0, charlie);
    }

    function test_proposeTransfer_revert_zeroProposed() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);

        vm.expectRevert(AgentRegistry.ZeroAddress.selector);
        vm.prank(alice);
        registry.proposeControllerTransfer(0, address(0));
    }

    function test_acceptTransfer_revert_noPending() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);

        vm.expectRevert(abi.encodeWithSelector(AgentRegistry.NoPendingTransfer.selector, 0));
        vm.prank(bob);
        registry.acceptControllerTransfer(0);
    }

    function test_acceptTransfer_revert_notProposed() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);
        vm.prank(alice);
        registry.proposeControllerTransfer(0, bob);

        vm.expectRevert(abi.encodeWithSelector(AgentRegistry.NotProposedController.selector, 0, charlie));
        vm.prank(charlie);
        registry.acceptControllerTransfer(0);
    }

    function test_cancelTransfer() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);
        vm.prank(alice);
        registry.proposeControllerTransfer(0, bob);

        vm.expectEmit(true, false, false, false);
        emit ControllerTransferCancelled(0);

        vm.prank(alice);
        registry.cancelControllerTransfer(0);
        assertEq(registry.pendingController(0), address(0));
    }

    function test_cancelTransfer_revert_noPending() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);

        vm.expectRevert(abi.encodeWithSelector(AgentRegistry.NoPendingTransfer.selector, 0));
        vm.prank(alice);
        registry.cancelControllerTransfer(0);
    }

    // ─── Reputation / scoring ────────────────────────────────────

    function test_postScore_success() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);

        vm.expectEmit(true, true, false, true);
        emit ScorePosted(0, scorer, 10, 10, 0);

        vm.prank(scorer);
        registry.postScore(0, 10, 0);

        assertEq(registry.getAgent(0).reputationScore, 10);
        assertEq(registry.scorerNonce(0), 1);
    }

    function test_postScore_accumulates() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);

        vm.prank(scorer);
        registry.postScore(0, 10, 0);
        vm.prank(scorer);
        registry.postScore(0, -3, 1);

        assertEq(registry.getAgent(0).reputationScore, 7);
        assertEq(registry.scorerNonce(0), 2);
    }

    function test_postScore_revert_wrongNonce() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);

        vm.expectRevert(abi.encodeWithSelector(AgentRegistry.InvalidNonce.selector, 0, 0, 1));
        vm.prank(scorer);
        registry.postScore(0, 10, 1);
    }

    function test_postScore_revert_notScorer() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);

        vm.expectRevert();
        vm.prank(alice);
        registry.postScore(0, 10, 0);
    }

    function test_postScore_revert_agentInactive() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);
        vm.prank(alice);
        registry.setActive(0, false);

        vm.expectRevert(abi.encodeWithSelector(AgentRegistry.AgentInactive.selector, 0));
        vm.prank(scorer);
        registry.postScore(0, 10, 0);
    }

    function test_postScore_revert_agentNotFound() public {
        vm.expectRevert(abi.encodeWithSelector(AgentRegistry.AgentNotFound.selector, 99));
        vm.prank(scorer);
        registry.postScore(99, 10, 0);
    }

    // ─── Views ───────────────────────────────────────────────────

    function test_agentsByController() public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);
        vm.prank(alice);
        registry.register(alice, "ipfs://QmB", 0x2);

        uint256[] memory ids = registry.agentsByController(alice);
        assertEq(ids.length, 2);
        assertEq(ids[0], 0);
        assertEq(ids[1], 1);
    }

    function test_getAgent_revert_notFound() public {
        vm.expectRevert(abi.encodeWithSelector(AgentRegistry.AgentNotFound.selector, 0));
        registry.getAgent(0);
    }

    // ─── Pause / admin ───────────────────────────────────────────

    function test_pause_unpause() public {
        vm.prank(admin);
        registry.pause();
        assertTrue(registry.paused());

        vm.prank(admin);
        registry.unpause();
        assertFalse(registry.paused());
    }

    function test_pause_revert_notPauser() public {
        vm.expectRevert();
        vm.prank(alice);
        registry.pause();
    }

    // ─── Fuzz tests ──────────────────────────────────────────────

    /// @notice Fuzz: any non-zero controller + non-empty URI always succeeds.
    function testFuzz_register(address ctrl, uint256 caps) public {
        vm.assume(ctrl != address(0));
        vm.prank(ctrl);
        uint256 id = registry.register(ctrl, URI, caps);
        AgentRegistry.AgentInfo memory info = registry.getAgent(id);
        assertEq(info.controller, ctrl);
        assertEq(info.capabilities, caps);
    }

    /// @notice Fuzz: sequential nonces always accepted, out-of-order always rejected.
    function testFuzz_postScore_nonceProtection(int256 delta) public {
        vm.prank(alice);
        registry.register(alice, URI, CAPS);

        // Correct nonce — succeeds
        vm.prank(scorer);
        registry.postScore(0, delta, 0);

        // Same nonce again — must revert
        vm.expectRevert(abi.encodeWithSelector(AgentRegistry.InvalidNonce.selector, 0, 1, 0));
        vm.prank(scorer);
        registry.postScore(0, delta, 0);
    }

    /// @notice Fuzz: reputation accumulates correctly for arbitrary deltas.
    function testFuzz_reputation_accumulation(int256 d1, int256 d2) public {
        // Avoid overflow in test expectations
        d1 = bound(d1, type(int128).min, type(int128).max);
        d2 = bound(d2, type(int128).min, type(int128).max);

        vm.prank(alice);
        registry.register(alice, URI, CAPS);

        vm.prank(scorer);
        registry.postScore(0, d1, 0);
        vm.prank(scorer);
        registry.postScore(0, d2, 1);

        int256 expected = d1 + d2;
        assertEq(registry.getAgent(0).reputationScore, expected);
    }
}
