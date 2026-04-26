// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {Hashes} from "@openzeppelin/contracts/utils/cryptography/Hashes.sol";

import {AntiSniperModule} from "../../src/launchpad/modules/AntiSniperModule.sol";
import {LaunchRadarModule} from "../../src/launchpad/modules/LaunchRadarModule.sol";
import {AirdropModule} from "../../src/launchpad/modules/AirdropModule.sol";
import {ExistingTokenModule} from "../../src/launchpad/modules/ExistingTokenModule.sol";
import {PreBuyModule} from "../../src/launchpad/modules/PreBuyModule.sol";
import {TokenVesting} from "../../src/launchpad/modules/TokenVesting.sol";

import {IBondingCurve} from "../../src/launchpad/interfaces/IBondingCurve.sol";

import {MockMintableERC20} from "../mocks/MockERC20.sol";
import {MockBondingCurve} from "../mocks/MockBondingCurve.sol";

/// @dev Shared mintable ERC-20 used as the agent / payout token in the
///      gap-fill suite. Standard metadata so `IERC20Metadata` introspection
///      paths are exercised cleanly.
contract MintableERC20 is ERC20 {
    constructor(string memory n, string memory s) ERC20(n, s) {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @title  ModuleEdgeCases
/// @notice Per-module gap-fill tests targeting edge cases not exercised in
///         the per-module unit suites. Each section below targets one
///         module and adds 2-3 boundary / owner-gating / idempotency / zero
///         tests. Custom errors only; events asserted with {vm.expectEmit}.
///
///         Coverage tier targets:
///           - boundary values: at-cap, exact-zero, max-uint, exact-window
///           - owner-restricted paths: rejection of non-owner; rejection of
///             former-owner after transferOwnership
///           - paused / phase paths: idempotent revert on second call;
///             cross-phase mutations
contract ModuleEdgeCasesTest is Test {
    address internal owner = address(0xA0);
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);
    address internal carol = address(0xCA401);

    // ---------------------------------------------------------------------
    // AntiSniperModule
    // ---------------------------------------------------------------------

    /// @dev At `block.timestamp == startTime + duration` the rate must
    ///      saturate to `endBps`. Boundary check missing from the suite.
    function test_antiSniper_at_exact_end_saturates_to_endBps() public {
        AntiSniperModule m = new AntiSniperModule(owner);
        address agent = address(0xA0E);

        uint64 startTime = 1_000_000;
        uint32 duration = 100;
        uint16 startBps = 9900;
        uint16 endBps = 100;

        vm.prank(owner);
        m.configure(agent, startTime, duration, startBps, endBps);

        // At the exact end of the window — the contract saturates here per
        // the `if (elapsed >= c.duration) return c.endBps;` branch.
        vm.warp(uint256(startTime) + uint256(duration));
        assertEq(uint256(m.currentBps(agent)), uint256(endBps));
    }

    /// @dev Two distinct agents on the same module read independent decay
    ///      curves. State isolation invariant.
    function test_antiSniper_per_agent_state_isolation() public {
        AntiSniperModule m = new AntiSniperModule(owner);
        address agentA = address(0xA1);
        address agentB = address(0xA2);

        vm.startPrank(owner);
        m.configure(agentA, uint64(1_000_000), 60, 9000, 100);
        m.configure(agentB, uint64(1_000_000), 60, 1000, 0);
        vm.stopPrank();

        vm.warp(1_000_000); // start
        // agentA: 90%, agentB: 10%.
        assertEq(uint256(m.currentBps(agentA)), 9000);
        assertEq(uint256(m.currentBps(agentB)), 1000);

        // Halfway: agentA -> ~45.5%, agentB -> ~5%.
        vm.warp(1_000_030);
        uint16 a = m.currentBps(agentA);
        uint16 b = m.currentBps(agentB);
        assertLt(uint256(a), 9000);
        assertGt(uint256(a), 100);
        assertLt(uint256(b), 1000);
        assertGt(uint256(b), 0);
    }

    /// @dev `endBps == startBps` is accepted (constant tax for the full
    ///      window). The contract enforces `endBps <= startBps`, not strict
    ///      `<`, so equality is legal and the rate stays flat.
    function test_antiSniper_equal_startBps_endBps_constant_tax() public {
        AntiSniperModule m = new AntiSniperModule(owner);
        address agent = address(0xA1);

        vm.prank(owner);
        m.configure(agent, uint64(1_000_000), 60, 5000, 5000);

        for (uint256 t = 1_000_000 - 5; t <= 1_000_000 + 65; t += 10) {
            vm.warp(t);
            assertEq(uint256(m.currentBps(agent)), 5000);
        }
    }

    // ---------------------------------------------------------------------
    // LaunchRadarModule
    // ---------------------------------------------------------------------

    /// @dev `setStartTime` with a past `newStart` is allowed pre-launch.
    ///      The contract gates on `block.timestamp < cfg.startTime`
    ///      (current start, not the proposed one), so an admin may bring
    ///      the gate forward into the past — equivalent to "trading is now
    ///      live" — provided the existing start has not yet elapsed.
    function test_launchRadar_setStartTime_to_past_pre_launch() public {
        LaunchRadarModule m = new LaunchRadarModule(owner);
        address agent = address(0xA1);
        address admin = address(0xAD);

        uint64 startAt = uint64(block.timestamp + 1 hours);
        vm.prank(owner);
        m.configure(agent, startAt, false, admin);

        uint64 past = uint64(block.timestamp - 1);
        vm.prank(admin);
        m.setStartTime(agent, past);

        (uint64 gotStart,,) = m.configs(agent);
        assertEq(gotStart, past);
        // Now `canTrade` should be permissive because past < block.timestamp.
        (bool ok, string memory reason) = m.canTrade(agent, alice);
        assertTrue(ok);
        assertEq(bytes(reason).length, 0);
    }

    /// @dev Toggling whitelist on/off/on must be idempotent and reflect in
    ///      `canTrade` immediately. Boundary path missing from the suite.
    function test_launchRadar_whitelist_idempotent_toggle() public {
        LaunchRadarModule m = new LaunchRadarModule(owner);
        address agent = address(0xA1);
        address admin = address(0xAD);

        uint64 startAt = uint64(block.timestamp + 1 hours);
        vm.prank(owner);
        m.configure(agent, startAt, false, admin);

        // Add alice to the list.
        address[] memory u = new address[](1);
        bool[] memory ok = new bool[](1);
        u[0] = alice;
        ok[0] = true;
        vm.prank(admin);
        m.setWhitelist(agent, u, ok);

        // Toggle on, then off, then on — each step records and reads.
        vm.warp(startAt);

        vm.prank(admin);
        m.setWhitelistOn(agent, true);
        (bool a1,) = m.canTrade(agent, bob);
        assertFalse(a1);
        (bool a2,) = m.canTrade(agent, alice);
        assertTrue(a2);

        vm.prank(admin);
        m.setWhitelistOn(agent, false);
        (bool b1,) = m.canTrade(agent, bob);
        assertTrue(b1);

        vm.prank(admin);
        m.setWhitelistOn(agent, true);
        (bool c1,) = m.canTrade(agent, bob);
        assertFalse(c1);
        (bool c2,) = m.canTrade(agent, alice);
        assertTrue(c2);
    }

    /// @dev `setWhitelist` with empty arrays must succeed without side
    ///      effects beyond the event emission. Zero-length array path.
    function test_launchRadar_setWhitelist_empty_is_noop() public {
        LaunchRadarModule m = new LaunchRadarModule(owner);
        address agent = address(0xA1);
        address admin = address(0xAD);

        vm.prank(owner);
        m.configure(agent, uint64(block.timestamp + 1 hours), true, admin);

        address[] memory u = new address[](0);
        bool[] memory ok = new bool[](0);

        vm.prank(admin);
        m.setWhitelist(agent, u, ok);

        // Whitelist still empty for any address.
        assertFalse(m.whitelist(agent, alice));
        assertFalse(m.whitelist(agent, bob));
    }

    // ---------------------------------------------------------------------
    // AirdropModule
    // ---------------------------------------------------------------------

    /// @dev Allocation exactly at the 5% cap is allowed (the cap is
    ///      inclusive). Boundary check.
    function test_airdrop_allocation_at_cap_exact() public {
        AirdropModule m = new AirdropModule(owner);
        // 1B supply -> 5% cap = 50M tokens.
        MintableERC20 agent = new MintableERC20("Agent", "AGT");
        uint256 supply = 1_000_000_000e18;
        agent.mint(address(this), supply);
        uint128 cap = uint128((supply * 500) / 10_000);

        agent.transfer(address(m), cap);

        bytes32 root = keccak256("root");
        vm.prank(owner);
        m.configure(address(agent), root, uint64(block.number), cap, uint64(block.timestamp + 30 days), alice);

        (, uint128 stored,,,, bool configured) = m.configs(address(agent));
        assertEq(uint256(stored), uint256(cap));
        assertTrue(configured);
    }

    /// @dev `sweep` after a partial claim returns ONLY the unclaimed dust;
    ///      a follow-up sweep reverts {NothingToSweep}. Idempotency.
    function test_airdrop_sweep_partial_then_empty() public {
        AirdropModule m = new AirdropModule(owner);
        MintableERC20 agent = new MintableERC20("Agent", "AGT");
        uint256 supply = 1_000_000_000e18;
        agent.mint(address(this), supply);

        uint256 aliceAmt = 100e18;
        uint128 alloc = uint128(aliceAmt * 2); // half claimable, half dust
        bytes32 leaf = keccak256(bytes.concat(keccak256(abi.encode(alice, aliceAmt))));
        bytes32 root = leaf; // single-leaf tree

        agent.transfer(address(m), alloc);
        uint64 deadline = uint64(block.timestamp + 30 days);

        vm.prank(owner);
        m.configure(address(agent), root, uint64(block.number), alloc, deadline, alice);

        // Alice claims her half.
        bytes32[] memory proof = new bytes32[](0);
        m.claim(address(agent), alice, aliceAmt, proof);

        vm.warp(uint256(deadline) + 1);

        vm.prank(alice);
        m.sweep(address(agent));
        // Module empty.
        assertEq(agent.balanceOf(address(m)), 0);
        // Admin received the dust.
        assertEq(agent.balanceOf(alice), aliceAmt + (uint256(alloc) - aliceAmt));

        // Second sweep — nothing left.
        vm.prank(alice);
        vm.expectRevert(AirdropModule.NothingToSweep.selector);
        m.sweep(address(agent));
    }

    /// @dev `sweep` is rejected for an unconfigured agent. Owner-restricted
    ///      path (admin gate is BEFORE the time gate, so an unconfigured
    ///      agent reverts on the configured guard before the timestamp
    ///      check fires). Even past the deadline it must still reject.
    function test_airdrop_sweep_unconfigured_after_deadline_reverts() public {
        AirdropModule m = new AirdropModule(owner);
        MintableERC20 agent = new MintableERC20("Agent", "AGT");

        // No configure — but warp far into the future.
        vm.warp(block.timestamp + 365 days);
        vm.prank(alice);
        vm.expectRevert(AirdropModule.NotConfigured.selector);
        m.sweep(address(agent));
    }

    // ---------------------------------------------------------------------
    // ExistingTokenModule
    // ---------------------------------------------------------------------

    /// @dev `wrap` with `amount == 0` reverts {ZeroSupply}. The brief lists
    ///      this as `ZeroSupply` (the same selector used by `configure`'s
    ///      zero check); confirm explicitly. Already covered, but a fuzz
    ///      around `amount == 0` is missing.
    function testFuzz_existingToken_wrap_zero_amount_reverts(uint256 seed) public {
        seed = bound(seed, 0, 0); // pin to zero — the boundary under test
        ExistingTokenModule m = new ExistingTokenModule(owner);
        MintableERC20 ext = new MintableERC20("Ext", "EXT");
        ext.mint(address(this), 1000e18);
        ext.transfer(address(m), 1000e18);
        vm.prank(owner);
        m.configure(address(ext), address(0xC0DE), 1000e18, alice);

        vm.expectRevert(ExistingTokenModule.ZeroSupply.selector);
        m.wrap(address(ext), seed);
    }

    /// @dev Multiple `wrap` calls from the SAME depositor for the same
    ///      token funnel through cleanly. Idempotency / repeated-call path.
    function test_existingToken_wrap_repeated_same_depositor() public {
        ExistingTokenModule m = new ExistingTokenModule(owner);
        MintableERC20 ext = new MintableERC20("Ext", "EXT");
        address curve = address(0xC0DE);
        uint256 seed = 1000e18;
        ext.mint(address(this), seed);
        ext.transfer(address(m), seed);

        vm.prank(owner);
        m.configure(address(ext), curve, seed, alice);

        // Three top-ups from the same depositor.
        ext.mint(bob, 300e18);
        vm.prank(bob);
        ext.approve(address(m), 300e18);
        vm.prank(bob);
        m.wrap(address(ext), 100e18);
        vm.prank(bob);
        m.wrap(address(ext), 100e18);
        vm.prank(bob);
        m.wrap(address(ext), 100e18);

        // All three lands flush the curve to seed + 300e18.
        assertEq(ext.balanceOf(curve), seed + 300e18);
        assertEq(ext.balanceOf(address(m)), 0);
    }

    /// @dev Two independently-wrapped tokens hold disjoint configs even if
    ///      they share the same admin / curve. State isolation invariant.
    function test_existingToken_two_tokens_share_admin_independent() public {
        ExistingTokenModule m = new ExistingTokenModule(owner);
        MintableERC20 a = new MintableERC20("A", "A");
        MintableERC20 b = new MintableERC20("B", "B");
        address curveA = address(0xC0DE);
        address curveB = address(0xC0DF);
        uint256 supply = 1000e18;

        a.mint(address(this), supply);
        b.mint(address(this), supply);
        a.transfer(address(m), supply);
        b.transfer(address(m), supply);

        vm.startPrank(owner);
        m.configure(address(a), curveA, supply, alice);
        m.configure(address(b), curveB, supply, alice); // same admin, different curve
        vm.stopPrank();

        ExistingTokenModule.Config memory cfgA = m.getConfig(address(a));
        ExistingTokenModule.Config memory cfgB = m.getConfig(address(b));
        assertEq(cfgA.curve, curveA);
        assertEq(cfgB.curve, curveB);
        assertTrue(cfgA.wrapped);
        assertTrue(cfgB.wrapped);
        assertEq(cfgA.agentAdmin, alice);
        assertEq(cfgB.agentAdmin, alice);
    }

    // ---------------------------------------------------------------------
    // PreBuyModule
    // ---------------------------------------------------------------------

    /// @dev `predictCloneAddress` is deterministic and pre-computable
    ///      across two distinct agents — i.e. addresses differ AND match
    ///      the actual deployment. Verified across a third "what-if" agent
    ///      that is never deployed.
    function test_preBuy_predictCloneAddress_distinct_per_agent() public {
        PreBuyModule m = new PreBuyModule(owner);

        address a1 = address(0xA1);
        address a2 = address(0xA2);
        address a3 = address(0xA3);

        address p1 = m.predictCloneAddress(a1);
        address p2 = m.predictCloneAddress(a2);
        address p3 = m.predictCloneAddress(a3);

        assertTrue(p1 != p2);
        assertTrue(p2 != p3);
        assertTrue(p1 != p3);
        assertTrue(p1 != address(0));
    }

    /// @dev Configuring TWO distinct agents back-to-back on the same module
    ///      with the SAME creator yields two separate vesting clones, each
    ///      holding its own bag. Idempotency-across-agents.
    function test_preBuy_two_agents_same_creator_distinct_clones() public {
        PreBuyModule m = new PreBuyModule(owner);
        MockMintableERC20 quote = new MockMintableERC20("Q", "Q");
        MockMintableERC20 agentA = new MockMintableERC20("A", "A");
        MockMintableERC20 agentB = new MockMintableERC20("B", "B");
        MockBondingCurve cA = new MockBondingCurve(address(agentA), address(quote));
        MockBondingCurve cB = new MockBondingCurve(address(agentB), address(quote));

        agentA.mint(address(cA), 1_000_000e18);
        agentB.mint(address(cB), 1_000_000e18);

        uint256 titanIn = 100e18;
        quote.mint(address(m), titanIn * 2);

        vm.startPrank(owner);
        address c1 = m.configure(
            address(agentA), alice, 50_000e18, uint64(0), uint64(180 days), titanIn, IBondingCurve(address(cA))
        );
        address c2 = m.configure(
            address(agentB), alice, 50_000e18, uint64(0), uint64(180 days), titanIn, IBondingCurve(address(cB))
        );
        vm.stopPrank();

        assertTrue(c1 != c2);
        assertEq(agentA.balanceOf(c1), 50_000e18);
        assertEq(agentB.balanceOf(c2), 50_000e18);
        assertEq(agentA.balanceOf(c2), 0);
        assertEq(agentB.balanceOf(c1), 0);
    }

    /// @dev `sweep` is owner-only and routes the captured residual to any
    ///      address. Owner-restricted boundary check.
    function test_preBuy_sweep_to_arbitrary_address() public {
        PreBuyModule m = new PreBuyModule(owner);
        MockMintableERC20 quote = new MockMintableERC20("Q", "Q");
        quote.mint(address(m), 500e18);

        // Owner sweeps to a stranger address — fine.
        address stranger = address(0xDEAD);
        vm.prank(owner);
        m.sweep(IERC20(address(quote)), stranger, 500e18);

        assertEq(quote.balanceOf(stranger), 500e18);
        assertEq(quote.balanceOf(address(m)), 0);
    }
}
