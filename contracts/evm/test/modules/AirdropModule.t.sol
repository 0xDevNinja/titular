// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {Hashes} from "@openzeppelin/contracts/utils/cryptography/Hashes.sol";

import {AirdropModule} from "../../src/launchpad/modules/AirdropModule.sol";

/// @dev Plain ERC-20 stand-in for the agent token under airdrop. Mints the
///      full supply to the deployer at construction.
contract AirdropAgentToken is ERC20 {
    constructor(uint256 supply, address mintTo) ERC20("Agent", "AGT") {
        _mint(mintTo, supply);
    }
}

/// @dev Reentrant ERC-20 — every transfer re-enters {AirdropModule.claim}
///      with the supplied victim arguments. Used to verify the reentrancy
///      guard rejects the inner call.
contract AirdropReentrantToken is ERC20 {
    AirdropModule public module;
    address public agent;
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

    function _update(address from, address to, uint256 value) internal override {
        if (reenter && address(module) != address(0)) {
            reenter = false;
            module.claim(agent, victim, victimAmount, victimProof);
        }
        super._update(from, to, value);
    }
}

/// @title  AirdropModuleSuiteTest
/// @notice Per-contract unit suite for {AirdropModule}: configure validation
///         (cap, root, deadline, balance), merkle claim verification with
///         valid + invalid proofs, double-claim guard, sweep gating, and a
///         reentrancy probe via a malicious agent token.
contract AirdropModuleSuiteTest is Test {
    AirdropModule internal module;
    AirdropAgentToken internal agentToken;

    address internal owner = address(0xA0);
    address internal admin = address(0xAD);
    address internal attacker = address(0xBAD);
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);
    address internal carol = address(0xCA401);

    uint256 internal aliceAmount = 100 ether;
    uint256 internal bobAmount = 250 ether;
    uint256 internal carolAmount = 50 ether;
    uint256 internal constant SUPPLY = 1_000_000_000 ether;
    uint128 internal allocation = 1000 ether;
    uint64 internal deadline;

    bytes32 internal root;
    bytes32 internal aliceLeaf;
    bytes32 internal bobLeaf;
    bytes32 internal carolLeaf;

    function setUp() public {
        module = new AirdropModule(owner);
        agentToken = new AirdropAgentToken(SUPPLY, address(this));

        aliceLeaf = _leaf(alice, aliceAmount);
        bobLeaf = _leaf(bob, bobAmount);
        carolLeaf = _leaf(carol, carolAmount);
        bytes32 ab = Hashes.commutativeKeccak256(aliceLeaf, bobLeaf);
        root = Hashes.commutativeKeccak256(ab, carolLeaf);

        vm.warp(1_000_000);
        deadline = uint64(block.timestamp + 30 days);
    }

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    function test_constructor_zero_owner_reverts() public {
        vm.expectRevert(AirdropModule.ZeroAddress.selector);
        new AirdropModule(address(0));
    }

    function test_constructor_constants() public view {
        assertEq(uint256(module.MAX_BPS()), 10_000);
        assertEq(uint256(module.MAX_ALLOCATION_BPS()), 500);
    }

    // ---------------------------------------------------------------------
    // configure
    // ---------------------------------------------------------------------

    function test_configure_happy_path() public {
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

    function test_configure_only_owner() public {
        agentToken.transfer(address(module), allocation);
        vm.prank(attacker);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, attacker));
        module.configure(address(agentToken), root, 1, allocation, deadline, admin);
    }

    function test_configure_zero_agent_reverts() public {
        vm.prank(owner);
        vm.expectRevert(AirdropModule.ZeroAddress.selector);
        module.configure(address(0), root, 1, allocation, deadline, admin);
    }

    function test_configure_zero_admin_reverts() public {
        agentToken.transfer(address(module), allocation);
        vm.prank(owner);
        vm.expectRevert(AirdropModule.ZeroAddress.selector);
        module.configure(address(agentToken), root, 1, allocation, deadline, address(0));
    }

    function test_configure_zero_root_reverts() public {
        agentToken.transfer(address(module), allocation);
        vm.prank(owner);
        vm.expectRevert(AirdropModule.ZeroRoot.selector);
        module.configure(address(agentToken), bytes32(0), 1, allocation, deadline, admin);
    }

    function test_configure_zero_allocation_reverts() public {
        vm.prank(owner);
        vm.expectRevert(AirdropModule.ZeroAllocation.selector);
        module.configure(address(agentToken), root, 1, 0, deadline, admin);
    }

    function test_configure_deadline_in_past_reverts() public {
        agentToken.transfer(address(module), allocation);
        uint64 past = uint64(block.timestamp);
        vm.prank(owner);
        vm.expectRevert(AirdropModule.DeadlineInPast.selector);
        module.configure(address(agentToken), root, 1, allocation, past, admin);
    }

    function test_configure_above_5pct_cap_reverts() public {
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

    function test_configure_insufficient_pre_deposit_reverts() public {
        uint256 short = uint256(allocation) - 1;
        agentToken.transfer(address(module), short);
        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(AirdropModule.InsufficientBalance.selector, short, uint256(allocation)));
        module.configure(address(agentToken), root, 1, allocation, deadline, admin);
    }

    function test_configure_one_shot() public {
        agentToken.transfer(address(module), allocation);
        vm.prank(owner);
        module.configure(address(agentToken), root, 1, allocation, deadline, admin);
        agentToken.transfer(address(module), allocation);
        vm.prank(owner);
        vm.expectRevert(AirdropModule.AlreadyConfigured.selector);
        module.configure(address(agentToken), root, 2, allocation, deadline, admin);
    }

    // ---------------------------------------------------------------------
    // claim
    // ---------------------------------------------------------------------

    function test_claim_valid_proof_succeeds() public {
        _configureDefault();
        bytes32[] memory proof = _aliceProof();

        vm.expectEmit(true, true, false, true, address(module));
        emit AirdropModule.Claimed(address(agentToken), alice, aliceAmount);

        vm.prank(attacker); // anyone can relay
        module.claim(address(agentToken), alice, aliceAmount, proof);

        assertEq(agentToken.balanceOf(alice), aliceAmount);
        assertTrue(module.claimed(address(agentToken), alice));
    }

    function test_claim_all_three_recipients() public {
        _configureDefault();
        module.claim(address(agentToken), alice, aliceAmount, _aliceProof());
        module.claim(address(agentToken), bob, bobAmount, _bobProof());
        module.claim(address(agentToken), carol, carolAmount, _carolProof());

        assertEq(agentToken.balanceOf(alice), aliceAmount);
        assertEq(agentToken.balanceOf(bob), bobAmount);
        assertEq(agentToken.balanceOf(carol), carolAmount);
    }

    function test_claim_unconfigured_reverts() public {
        bytes32[] memory proof = new bytes32[](0);
        vm.expectRevert(AirdropModule.NotConfigured.selector);
        module.claim(address(agentToken), alice, aliceAmount, proof);
    }

    function test_claim_wrong_amount_reverts() public {
        _configureDefault();
        vm.expectRevert(AirdropModule.InvalidProof.selector);
        module.claim(address(agentToken), alice, aliceAmount + 1, _aliceProof());
    }

    function test_claim_wrong_recipient_reverts() public {
        _configureDefault();
        vm.expectRevert(AirdropModule.InvalidProof.selector);
        module.claim(address(agentToken), attacker, aliceAmount, _aliceProof());
    }

    function test_claim_garbage_proof_reverts() public {
        _configureDefault();
        bytes32[] memory bad = new bytes32[](2);
        bad[0] = bytes32(uint256(0xDEAD));
        bad[1] = bytes32(uint256(0xBEEF));
        vm.expectRevert(AirdropModule.InvalidProof.selector);
        module.claim(address(agentToken), alice, aliceAmount, bad);
    }

    function test_claim_double_claim_reverts() public {
        _configureDefault();
        bytes32[] memory proof = _aliceProof();
        module.claim(address(agentToken), alice, aliceAmount, proof);
        vm.expectRevert(AirdropModule.AlreadyClaimed.selector);
        module.claim(address(agentToken), alice, aliceAmount, proof);
    }

    function test_claim_after_deadline_reverts() public {
        _configureDefault();
        vm.warp(uint256(deadline) + 1);
        vm.expectRevert(AirdropModule.DeadlinePassed.selector);
        module.claim(address(agentToken), alice, aliceAmount, _aliceProof());
    }

    function test_claim_at_exact_deadline_allowed() public {
        _configureDefault();
        vm.warp(uint256(deadline));
        module.claim(address(agentToken), alice, aliceAmount, _aliceProof());
        assertEq(agentToken.balanceOf(alice), aliceAmount);
    }

    // ---------------------------------------------------------------------
    // sweep
    // ---------------------------------------------------------------------

    function test_sweep_pre_deadline_reverts() public {
        _configureDefault();
        vm.prank(admin);
        vm.expectRevert(AirdropModule.DeadlineNotPassed.selector);
        module.sweep(address(agentToken));
    }

    function test_sweep_at_deadline_reverts() public {
        _configureDefault();
        vm.warp(uint256(deadline));
        vm.prank(admin);
        vm.expectRevert(AirdropModule.DeadlineNotPassed.selector);
        module.sweep(address(agentToken));
    }

    function test_sweep_unconfigured_reverts() public {
        vm.warp(uint256(deadline) + 1);
        vm.prank(admin);
        vm.expectRevert(AirdropModule.NotConfigured.selector);
        module.sweep(address(agentToken));
    }

    function test_sweep_non_admin_reverts() public {
        _configureDefault();
        vm.warp(uint256(deadline) + 1);
        vm.prank(attacker);
        vm.expectRevert(AirdropModule.NotAgentAdmin.selector);
        module.sweep(address(agentToken));
    }

    function test_sweep_returns_dust_to_admin() public {
        _configureDefault();
        // alice claims; bob+carol skip — dust = bobAmount + carolAmount + extra.
        module.claim(address(agentToken), alice, aliceAmount, _aliceProof());

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

    function test_sweep_nothing_to_sweep_reverts() public {
        _configureDefault();
        // Claim everything.
        module.claim(address(agentToken), alice, aliceAmount, _aliceProof());
        module.claim(address(agentToken), bob, bobAmount, _bobProof());
        module.claim(address(agentToken), carol, carolAmount, _carolProof());

        vm.warp(uint256(deadline) + 1);
        vm.prank(admin);
        // First sweep returns the residual dust (allocation - claims sum).
        module.sweep(address(agentToken));

        // Second sweep — nothing left.
        vm.prank(admin);
        vm.expectRevert(AirdropModule.NothingToSweep.selector);
        module.sweep(address(agentToken));
    }

    // ---------------------------------------------------------------------
    // Reentrancy probe
    // ---------------------------------------------------------------------

    function test_claim_reentrancy_blocked() public {
        AirdropReentrantToken evil = new AirdropReentrantToken(SUPPLY, address(this));

        bytes32 evilLeaf = _leaf(alice, aliceAmount);
        evil.transfer(address(module), allocation);
        vm.prank(owner);
        module.configure(address(evil), evilLeaf, 1, allocation, deadline, admin);

        bytes32[] memory proof = new bytes32[](0);
        evil.arm(module, address(evil), alice, aliceAmount, proof);

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

    function _aliceProof() internal view returns (bytes32[] memory p) {
        p = new bytes32[](2);
        p[0] = bobLeaf;
        p[1] = carolLeaf;
    }

    function _bobProof() internal view returns (bytes32[] memory p) {
        p = new bytes32[](2);
        p[0] = aliceLeaf;
        p[1] = carolLeaf;
    }

    function _carolProof() internal view returns (bytes32[] memory p) {
        p = new bytes32[](1);
        p[0] = Hashes.commutativeKeccak256(aliceLeaf, bobLeaf);
    }
}
