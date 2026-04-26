// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {MerkleProof} from "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import {Hashes} from "@openzeppelin/contracts/utils/cryptography/Hashes.sol";
import {AirdropModule} from "../../src/launchpad/modules/AirdropModule.sol";

/// @dev Plain ERC-20 used as the agent token under airdrop. Caps mint via the
///      constructor so we control total supply for the 5%-cap math.
contract MockAgentToken is ERC20 {
    constructor(uint256 supply, address mintTo) ERC20("Agent", "AGT") {
        _mint(mintTo, supply);
    }
}

/// @dev ERC-20 that re-enters the AirdropModule on every transfer. Used to
///      verify {AirdropModule.claim} and {AirdropModule.sweep} reject reentry.
contract ReenteringToken is ERC20 {
    AirdropModule public module;
    address public agent; // keyed module config
    address public victim;
    uint256 public victimAmount;
    bytes32[] public victimProof;
    bool public reenter;

    constructor(uint256 supply, address mintTo) ERC20("Reenter", "REE") {
        _mint(mintTo, supply);
    }

    function arm(AirdropModule m, address agent_, address victim_, uint256 amount_, bytes32[] memory proof_) external {
        module = m;
        agent = agent_;
        victim = victim_;
        victimAmount = amount_;
        victimProof = proof_;
        reenter = true;
    }

    function disarm() external {
        reenter = false;
    }

    function _update(address from, address to, uint256 value) internal override {
        if (reenter && address(module) != address(0)) {
            // One-shot reentry; flip first to avoid infinite recursion if the
            // guard somehow lets us in.
            reenter = false;
            module.claim(agent, victim, victimAmount, victimProof);
        }
        super._update(from, to, value);
    }
}

/// @notice Unit tests for AirdropModule. Builds a small merkle tree in-memory
///         using OZ's standard double-hash leaf encoding so the on-chain
///         {MerkleProof.verify} accepts our synthetic proofs without needing a
///         JS toolchain.
contract AirdropModuleTest is Test {
    AirdropModule internal module;
    MockAgentToken internal agentToken;

    address internal owner = address(0xA0);
    address internal admin = address(0xAD);
    address internal attacker = address(0xBAD);

    // Three claimants — gives us a 3-leaf tree (one duplicated/padded path)
    // and lets us produce both valid and invalid proofs.
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);
    address internal carol = address(0xCA401);

    uint256 internal aliceAmount = 100 ether;
    uint256 internal bobAmount = 250 ether;
    uint256 internal carolAmount = 50 ether;

    // 5% of 1B token totalSupply = 50M tokens. Alice + Bob + Carol = 400. Well
    // within the cap.
    uint256 internal constant SUPPLY = 1_000_000_000 ether;
    uint128 internal allocation = 1000 ether;
    uint64 internal deadline;

    bytes32 internal root;
    bytes32 internal aliceLeaf;
    bytes32 internal bobLeaf;
    bytes32 internal carolLeaf;

    function setUp() public {
        // Deploy module + agent token. The deployer (this test) owns the
        // initial supply so we can pre-fund the module.
        module = new AirdropModule(owner);
        agentToken = new MockAgentToken(SUPPLY, address(this));

        // Build a simple 3-leaf merkle tree using OZ's commutative-keccak256
        // pairing. Sort leaves so proofs are deterministic.
        aliceLeaf = _leaf(alice, aliceAmount);
        bobLeaf = _leaf(bob, bobAmount);
        carolLeaf = _leaf(carol, carolAmount);

        // Pair (alice, bob), then pair the result with carol.
        bytes32 ab = Hashes.commutativeKeccak256(aliceLeaf, bobLeaf);
        root = Hashes.commutativeKeccak256(ab, carolLeaf);

        // Deadline 30 days out, with `vm.warp` headroom.
        vm.warp(1_000_000);
        deadline = uint64(block.timestamp + 30 days);
    }

    // ---------------------------------------------------------------------
    // constructor
    // ---------------------------------------------------------------------

    function test_constructor_reverts_zero_owner() public {
        vm.expectRevert(AirdropModule.ZeroAddress.selector);
        new AirdropModule(address(0));
    }

    // ---------------------------------------------------------------------
    // configure — happy path
    // ---------------------------------------------------------------------

    function test_configure_happy_path_emits_and_stores() public {
        // Pre-fund the module.
        agentToken.transfer(address(module), allocation);

        vm.expectEmit(true, false, false, true, address(module));
        emit AirdropModule.Configured(address(agentToken), root, 12_345, allocation, deadline, admin);

        vm.prank(owner);
        module.configure(address(agentToken), root, 12_345, allocation, deadline, admin);

        (bytes32 r, uint128 alloc, uint64 snap, uint64 dl, address adm, bool configured) =
            module.configs(address(agentToken));
        assertEq(r, root);
        assertEq(uint256(alloc), uint256(allocation));
        assertEq(uint256(snap), 12_345);
        assertEq(uint256(dl), uint256(deadline));
        assertEq(adm, admin);
        assertTrue(configured);
    }

    // ---------------------------------------------------------------------
    // configure — access control + validation
    // ---------------------------------------------------------------------

    function test_configure_only_owner() public {
        agentToken.transfer(address(module), allocation);
        vm.prank(attacker);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, attacker));
        module.configure(address(agentToken), root, 1, allocation, deadline, admin);
    }

    function test_configure_reverts_zero_agent() public {
        vm.prank(owner);
        vm.expectRevert(AirdropModule.ZeroAddress.selector);
        module.configure(address(0), root, 1, allocation, deadline, admin);
    }

    function test_configure_reverts_zero_admin() public {
        agentToken.transfer(address(module), allocation);
        vm.prank(owner);
        vm.expectRevert(AirdropModule.ZeroAddress.selector);
        module.configure(address(agentToken), root, 1, allocation, deadline, address(0));
    }

    function test_configure_reverts_zero_root() public {
        agentToken.transfer(address(module), allocation);
        vm.prank(owner);
        vm.expectRevert(AirdropModule.ZeroRoot.selector);
        module.configure(address(agentToken), bytes32(0), 1, allocation, deadline, admin);
    }

    function test_configure_reverts_zero_allocation() public {
        vm.prank(owner);
        vm.expectRevert(AirdropModule.ZeroAllocation.selector);
        module.configure(address(agentToken), root, 1, 0, deadline, admin);
    }

    function test_configure_reverts_deadline_in_past() public {
        agentToken.transfer(address(module), allocation);
        uint64 past = uint64(block.timestamp);
        vm.prank(owner);
        vm.expectRevert(AirdropModule.DeadlineInPast.selector);
        module.configure(address(agentToken), root, 1, allocation, past, admin);
    }

    function test_configure_reverts_allocation_over_5pct_cap() public {
        // 5% of 1B = 50M tokens. Ask for 50M + 1 wei.
        uint256 cap = (SUPPLY * 500) / 10_000;
        uint128 over = uint128(cap + 1);
        agentToken.transfer(address(module), over);
        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(AirdropModule.AllocationOverCap.selector, uint256(over), cap));
        module.configure(address(agentToken), root, 1, over, deadline, admin);
    }

    function test_configure_at_5pct_cap_allowed() public {
        uint128 cap = uint128((SUPPLY * 500) / 10_000);
        agentToken.transfer(address(module), cap);
        vm.prank(owner);
        module.configure(address(agentToken), root, 1, cap, deadline, admin);
        (, uint128 stored,,,, bool configured) = module.configs(address(agentToken));
        assertEq(uint256(stored), uint256(cap));
        assertTrue(configured);
    }

    function test_configure_reverts_insufficient_pre_deposit() public {
        // Underfund by 1 wei.
        uint256 short = uint256(allocation) - 1;
        agentToken.transfer(address(module), short);
        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(AirdropModule.InsufficientBalance.selector, short, uint256(allocation)));
        module.configure(address(agentToken), root, 1, allocation, deadline, admin);
    }

    function test_configure_reverts_double_configure() public {
        agentToken.transfer(address(module), allocation);
        vm.prank(owner);
        module.configure(address(agentToken), root, 1, allocation, deadline, admin);
        // Top up so the second {configure} cannot fail on balance — must fail
        // on the AlreadyConfigured guard. Transfer is from this test (deployer)
        // so it does not collide with the `vm.prank(owner)` window.
        agentToken.transfer(address(module), allocation);
        vm.prank(owner);
        vm.expectRevert(AirdropModule.AlreadyConfigured.selector);
        module.configure(address(agentToken), root, 2, allocation, deadline, admin);
    }

    // ---------------------------------------------------------------------
    // claim — happy path + verification
    // ---------------------------------------------------------------------

    function test_claim_valid_proof_transfers_and_marks_nullifier() public {
        _configureDefault();

        // Proof for alice: sibling = bobLeaf, then carolLeaf.
        bytes32[] memory proof = new bytes32[](2);
        proof[0] = bobLeaf;
        proof[1] = carolLeaf;

        vm.expectEmit(true, true, false, true, address(module));
        emit AirdropModule.Claimed(address(agentToken), alice, aliceAmount);

        // Submitter is unrelated to recipient — leaf is keyed on `alice`, not
        // `msg.sender`, so anyone may relay.
        vm.prank(attacker);
        module.claim(address(agentToken), alice, aliceAmount, proof);

        assertEq(agentToken.balanceOf(alice), aliceAmount);
        assertTrue(module.claimed(address(agentToken), alice));
    }

    function test_claim_all_three_recipients_succeed() public {
        _configureDefault();

        // alice
        bytes32[] memory pa = new bytes32[](2);
        pa[0] = bobLeaf;
        pa[1] = carolLeaf;
        module.claim(address(agentToken), alice, aliceAmount, pa);

        // bob — sibling = aliceLeaf, then carolLeaf.
        bytes32[] memory pb = new bytes32[](2);
        pb[0] = aliceLeaf;
        pb[1] = carolLeaf;
        module.claim(address(agentToken), bob, bobAmount, pb);

        // carol — sibling = hash(alice, bob).
        bytes32[] memory pc = new bytes32[](1);
        pc[0] = Hashes.commutativeKeccak256(aliceLeaf, bobLeaf);
        module.claim(address(agentToken), carol, carolAmount, pc);

        assertEq(agentToken.balanceOf(alice), aliceAmount);
        assertEq(agentToken.balanceOf(bob), bobAmount);
        assertEq(agentToken.balanceOf(carol), carolAmount);
    }

    // ---------------------------------------------------------------------
    // claim — revert paths
    // ---------------------------------------------------------------------

    function test_claim_reverts_unconfigured() public {
        bytes32[] memory proof = new bytes32[](0);
        vm.expectRevert(AirdropModule.NotConfigured.selector);
        module.claim(address(agentToken), alice, aliceAmount, proof);
    }

    function test_claim_reverts_invalid_proof_amount() public {
        _configureDefault();
        bytes32[] memory proof = new bytes32[](2);
        proof[0] = bobLeaf;
        proof[1] = carolLeaf;
        // Wrong amount → leaf doesn't match.
        vm.expectRevert(AirdropModule.InvalidProof.selector);
        module.claim(address(agentToken), alice, aliceAmount + 1, proof);
    }

    function test_claim_reverts_invalid_proof_recipient() public {
        _configureDefault();
        bytes32[] memory proof = new bytes32[](2);
        proof[0] = bobLeaf;
        proof[1] = carolLeaf;
        // Right alice proof, wrong claimant — recipient hashed into leaf.
        vm.expectRevert(AirdropModule.InvalidProof.selector);
        module.claim(address(agentToken), attacker, aliceAmount, proof);
    }

    function test_claim_reverts_invalid_proof_garbage_path() public {
        _configureDefault();
        bytes32[] memory proof = new bytes32[](2);
        proof[0] = bytes32(uint256(0xDEAD));
        proof[1] = bytes32(uint256(0xBEEF));
        vm.expectRevert(AirdropModule.InvalidProof.selector);
        module.claim(address(agentToken), alice, aliceAmount, proof);
    }

    function test_claim_reverts_double_claim() public {
        _configureDefault();
        bytes32[] memory proof = new bytes32[](2);
        proof[0] = bobLeaf;
        proof[1] = carolLeaf;
        module.claim(address(agentToken), alice, aliceAmount, proof);
        vm.expectRevert(AirdropModule.AlreadyClaimed.selector);
        module.claim(address(agentToken), alice, aliceAmount, proof);
    }

    function test_claim_reverts_after_deadline() public {
        _configureDefault();
        bytes32[] memory proof = new bytes32[](2);
        proof[0] = bobLeaf;
        proof[1] = carolLeaf;
        // Step past the deadline.
        vm.warp(uint256(deadline) + 1);
        vm.expectRevert(AirdropModule.DeadlinePassed.selector);
        module.claim(address(agentToken), alice, aliceAmount, proof);
    }

    function test_claim_at_exact_deadline_allowed() public {
        _configureDefault();
        bytes32[] memory proof = new bytes32[](2);
        proof[0] = bobLeaf;
        proof[1] = carolLeaf;
        vm.warp(uint256(deadline));
        module.claim(address(agentToken), alice, aliceAmount, proof);
        assertEq(agentToken.balanceOf(alice), aliceAmount);
    }

    // ---------------------------------------------------------------------
    // sweep
    // ---------------------------------------------------------------------

    function test_sweep_reverts_before_deadline() public {
        _configureDefault();
        vm.prank(admin);
        vm.expectRevert(AirdropModule.DeadlineNotPassed.selector);
        module.sweep(address(agentToken));
    }

    function test_sweep_reverts_at_exact_deadline() public {
        _configureDefault();
        vm.warp(uint256(deadline));
        vm.prank(admin);
        vm.expectRevert(AirdropModule.DeadlineNotPassed.selector);
        module.sweep(address(agentToken));
    }

    function test_sweep_reverts_unconfigured() public {
        vm.warp(uint256(deadline) + 1);
        vm.prank(admin);
        vm.expectRevert(AirdropModule.NotConfigured.selector);
        module.sweep(address(agentToken));
    }

    function test_sweep_reverts_not_admin() public {
        _configureDefault();
        vm.warp(uint256(deadline) + 1);
        vm.prank(attacker);
        vm.expectRevert(AirdropModule.NotAgentAdmin.selector);
        module.sweep(address(agentToken));
    }

    function test_sweep_reverts_nothing_to_sweep() public {
        _configureDefault();
        // Drain via claims so balance == 0.
        bytes32[] memory pa = new bytes32[](2);
        pa[0] = bobLeaf;
        pa[1] = carolLeaf;
        module.claim(address(agentToken), alice, aliceAmount, pa);
        bytes32[] memory pb = new bytes32[](2);
        pb[0] = aliceLeaf;
        pb[1] = carolLeaf;
        module.claim(address(agentToken), bob, bobAmount, pb);
        bytes32[] memory pc = new bytes32[](1);
        pc[0] = Hashes.commutativeKeccak256(aliceLeaf, bobLeaf);
        module.claim(address(agentToken), carol, carolAmount, pc);
        // Drain the residual dust manually so balance hits exactly zero.
        // allocation = 1000 ether; claims sum = 400 ether; dust = 600 ether.
        // Forge sweep + reset to zero by simulating someone claiming the rest
        // is overkill — instead reconfigure another agent to exhaust. Simpler:
        // mint zero-balance scenario by setting allocation to exactly the
        // claimed sum. Refactored separate test below; here we verify dust > 0
        // sweep works and then test the empty-balance case by sweeping twice.
        vm.warp(uint256(deadline) + 1);
        vm.prank(admin);
        module.sweep(address(agentToken));
        // Second sweep — nothing left.
        vm.prank(admin);
        vm.expectRevert(AirdropModule.NothingToSweep.selector);
        module.sweep(address(agentToken));
    }

    function test_sweep_after_deadline_returns_dust_to_admin() public {
        _configureDefault();
        // alice claims; bob+carol skip — dust = bobAmount + carolAmount.
        bytes32[] memory pa = new bytes32[](2);
        pa[0] = bobLeaf;
        pa[1] = carolLeaf;
        module.claim(address(agentToken), alice, aliceAmount, pa);

        uint256 dust = uint256(allocation) - aliceAmount;
        assertEq(agentToken.balanceOf(address(module)), dust);

        vm.warp(uint256(deadline) + 1);

        vm.expectEmit(true, false, false, true, address(module));
        emit AirdropModule.Swept(address(agentToken), dust);
        vm.prank(admin);
        module.sweep(address(agentToken));

        assertEq(agentToken.balanceOf(admin), dust);
        assertEq(agentToken.balanceOf(address(module)), 0);
    }

    // ---------------------------------------------------------------------
    // Reentrancy probe
    // ---------------------------------------------------------------------

    function test_claim_reentrancy_blocked() public {
        // Build a fresh tree on a malicious agent token so the leaf->root
        // path matches the reentering claim.
        ReenteringToken evil = new ReenteringToken(SUPPLY, address(this));

        // Single-leaf tree for alice — proof is empty, root = leaf.
        bytes32 evilLeaf = _leaf(alice, aliceAmount);
        bytes32 evilRoot = evilLeaf; // 1-leaf tree

        evil.transfer(address(module), allocation);
        vm.prank(owner);
        module.configure(address(evil), evilRoot, 1, allocation, deadline, admin);

        bytes32[] memory proof = new bytes32[](0);

        // Arm the token to re-enter on the next transfer.
        evil.arm(module, address(evil), alice, aliceAmount, proof);

        // First call enters claim, flips nullifier, then transfer fires which
        // re-enters claim. ReentrancyGuard reverts the inner call which bubbles
        // up the outer call.
        vm.expectRevert();
        module.claim(address(evil), alice, aliceAmount, proof);
    }

    // ---------------------------------------------------------------------
    // Helpers
    // ---------------------------------------------------------------------

    function _leaf(address recipient, uint256 amount) internal pure returns (bytes32) {
        return keccak256(bytes.concat(keccak256(abi.encode(recipient, amount))));
    }

    function _configureDefault() internal {
        agentToken.transfer(address(module), allocation);
        vm.prank(owner);
        module.configure(address(agentToken), root, 12_345, allocation, deadline, admin);
    }
}
