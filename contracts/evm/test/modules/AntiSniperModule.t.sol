// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

import {AntiSniperModule} from "../../src/launchpad/modules/AntiSniperModule.sol";

/// @title  AntiSniperModuleSuiteTest
/// @notice Per-module unit tests living under test/modules/. Mirrors the
///         coverage in test/unit/AntiSniperModule.t.sol while focusing on the
///         contract-level invariants the composition suite cannot exercise:
///         pure-view monotonicity, owner gating, packed-storage round-trip,
///         and saturation at both ends of the decay curve.
contract AntiSniperModuleSuiteTest is Test {
    AntiSniperModule internal module;

    address internal owner = address(0xA0);
    address internal attacker = address(0xBAD);
    address internal agent = address(0xA1);
    address internal agent2 = address(0xA2);

    uint16 internal constant START_BPS = 9900;
    uint16 internal constant END_BPS = 100;
    uint32 internal constant DURATION = 60;
    uint64 internal constant START_TIME = 1_000_000;

    function setUp() public {
        module = new AntiSniperModule(owner);
        vm.warp(START_TIME - 1000);
    }

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    function test_constructor_zero_owner_reverts() public {
        vm.expectRevert(AntiSniperModule.ZeroAddress.selector);
        new AntiSniperModule(address(0));
    }

    function test_constructor_owner_set() public view {
        assertEq(module.owner(), owner);
        assertEq(uint256(module.MAX_BPS()), 10_000);
    }

    // ---------------------------------------------------------------------
    // configure — access + validation
    // ---------------------------------------------------------------------

    function test_configure_only_owner() public {
        vm.prank(attacker);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, attacker));
        module.configure(agent, START_TIME, DURATION, START_BPS, END_BPS);
    }

    function test_configure_one_shot() public {
        vm.startPrank(owner);
        module.configure(agent, START_TIME, DURATION, START_BPS, END_BPS);
        vm.expectRevert(AntiSniperModule.AlreadyConfigured.selector);
        module.configure(agent, START_TIME, DURATION, START_BPS, END_BPS);
        vm.stopPrank();
    }

    function test_configure_zero_agent_reverts() public {
        vm.prank(owner);
        vm.expectRevert(AntiSniperModule.ZeroAddress.selector);
        module.configure(address(0), START_TIME, DURATION, START_BPS, END_BPS);
    }

    function test_configure_startBps_above_max_reverts() public {
        uint16 bad = 10_001;
        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(AntiSniperModule.InvalidBps.selector, bad));
        module.configure(agent, START_TIME, DURATION, bad, END_BPS);
    }

    function test_configure_endBps_above_startBps_reverts() public {
        uint16 badEnd = START_BPS + 1;
        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(AntiSniperModule.InvalidBps.selector, badEnd));
        module.configure(agent, START_TIME, DURATION, START_BPS, badEnd);
    }

    function test_configure_zero_duration_reverts() public {
        vm.prank(owner);
        vm.expectRevert(AntiSniperModule.InvalidDuration.selector);
        module.configure(agent, START_TIME, 0, START_BPS, END_BPS);
    }

    function test_configure_emits_event_and_persists() public {
        vm.expectEmit(true, false, false, true, address(module));
        emit AntiSniperModule.Configured(agent, START_TIME, DURATION, START_BPS, END_BPS);
        vm.prank(owner);
        module.configure(agent, START_TIME, DURATION, START_BPS, END_BPS);

        (uint64 st, uint32 dur, uint16 sBps, uint16 eBps, bool configured) = module.configs(agent);
        assertEq(st, START_TIME);
        assertEq(uint256(dur), uint256(DURATION));
        assertEq(uint256(sBps), uint256(START_BPS));
        assertEq(uint256(eBps), uint256(END_BPS));
        assertTrue(configured);
    }

    // ---------------------------------------------------------------------
    // currentBps — every branch of the piecewise-linear decay
    // ---------------------------------------------------------------------

    function test_currentBps_unconfigured_zero() public view {
        assertEq(uint256(module.currentBps(agent)), 0);
    }

    function test_currentBps_pre_start_returns_startBps() public {
        _configureDefault();
        // setUp warped to START_TIME - 1_000 so we are pre-start.
        assertEq(uint256(module.currentBps(agent)), uint256(START_BPS));
    }

    function test_currentBps_at_start_returns_startBps() public {
        _configureDefault();
        vm.warp(START_TIME);
        assertEq(uint256(module.currentBps(agent)), uint256(START_BPS));
    }

    function test_currentBps_midway_linear_interpolation() public {
        // 100% -> 0% over 1_000s. 500s in -> exactly 5_000.
        vm.prank(owner);
        module.configure(agent2, 2_000_000, 1000, 10_000, 0);
        vm.warp(2_000_000 + 500);
        assertEq(uint256(module.currentBps(agent2)), 5000);
    }

    function test_currentBps_at_end_saturates_to_endBps() public {
        _configureDefault();
        vm.warp(uint256(START_TIME) + uint256(DURATION));
        assertEq(uint256(module.currentBps(agent)), uint256(END_BPS));
    }

    function test_currentBps_after_end_saturates_to_endBps() public {
        _configureDefault();
        vm.warp(uint256(START_TIME) + uint256(DURATION) + 365 days);
        assertEq(uint256(module.currentBps(agent)), uint256(END_BPS));
    }

    // ---------------------------------------------------------------------
    // computeTax
    // ---------------------------------------------------------------------

    function test_computeTax_unconfigured_returns_zero_tax() public view {
        (uint256 tax, uint256 net) = module.computeTax(agent, 1000);
        assertEq(tax, 0);
        assertEq(net, 1000);
    }

    function test_computeTax_at_start_taxes_startBps() public {
        _configureDefault();
        vm.warp(START_TIME);
        (uint256 tax, uint256 net) = module.computeTax(agent, 10_000 ether);
        assertEq(tax, (10_000 ether * uint256(START_BPS)) / 10_000);
        assertEq(net, 10_000 ether - tax);
    }

    function test_computeTax_zero_amount_yields_zero() public {
        _configureDefault();
        vm.warp(START_TIME);
        (uint256 tax, uint256 net) = module.computeTax(agent, 0);
        assertEq(tax, 0);
        assertEq(net, 0);
    }

    // ---------------------------------------------------------------------
    // Fuzz invariants
    // ---------------------------------------------------------------------

    /// @dev Monotone non-increasing in time for any configured agent.
    function testFuzz_currentBps_monotonic(uint32 e1, uint32 e2) public {
        _configureDefault();
        e1 = uint32(bound(uint256(e1), 0, uint256(DURATION) + 5000));
        e2 = uint32(bound(uint256(e2), 0, uint256(DURATION) + 5000));
        if (e1 > e2) (e1, e2) = (e2, e1);

        vm.warp(uint256(START_TIME) + uint256(e1));
        uint16 b1 = module.currentBps(agent);
        vm.warp(uint256(START_TIME) + uint256(e2));
        uint16 b2 = module.currentBps(agent);

        assertLe(uint256(b2), uint256(b1));
        assertLe(uint256(b1), uint256(START_BPS));
        assertGe(uint256(b2), uint256(END_BPS));
    }

    /// @dev tax + net == amount across the entire feasible amount range.
    function testFuzz_computeTax_conservation(uint256 amount) public {
        _configureDefault();
        amount = bound(amount, 0, type(uint256).max / uint256(module.MAX_BPS()));
        vm.warp(uint256(START_TIME) + uint256(DURATION) / 3);
        (uint256 tax, uint256 net) = module.computeTax(agent, amount);
        assertEq(tax + net, amount);
    }

    function _configureDefault() internal {
        vm.prank(owner);
        module.configure(agent, START_TIME, DURATION, START_BPS, END_BPS);
    }
}
