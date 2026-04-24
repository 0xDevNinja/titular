// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {AntiSniperModule} from "../../src/launchpad/modules/AntiSniperModule.sol";

/// @dev Unit tests for the shared AntiSniperModule. The module is pure-view:
///      we exercise the owner-gated {configure} path and the piecewise-linear
///      decay read out of {currentBps} / {computeTax}.
contract AntiSniperModuleTest is Test {
    AntiSniperModule internal module;

    address internal owner = address(0xA0);
    address internal attacker = address(0xBAD);
    address internal agent = address(0xA1);
    address internal agent2 = address(0xA2);

    // Representative 99% -> 1% decay over 60 seconds; matches the launch config
    // the factory will bind to each agent.
    uint16 internal constant START_BPS = 9900;
    uint16 internal constant END_BPS = 100;
    uint32 internal constant DURATION = 60;
    uint64 internal constant START_TIME = 1_000_000;

    function setUp() public {
        module = new AntiSniperModule(owner);
        // Put block.timestamp far below any launch window so pre-start branches
        // can be exercised without underflowing.
        vm.warp(START_TIME - 1000);
    }

    // ---------------------------------------------------------------------
    // configure — access control + validation
    // ---------------------------------------------------------------------

    function test_configure_only_owner() public {
        vm.prank(attacker);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, attacker));
        module.configure(agent, START_TIME, DURATION, START_BPS, END_BPS);
    }

    function test_configure_reverts_if_already_configured() public {
        vm.startPrank(owner);
        module.configure(agent, START_TIME, DURATION, START_BPS, END_BPS);
        vm.expectRevert(AntiSniperModule.AlreadyConfigured.selector);
        module.configure(agent, START_TIME, DURATION, START_BPS, END_BPS);
        vm.stopPrank();
    }

    function test_configure_reverts_zero_agent() public {
        vm.prank(owner);
        vm.expectRevert(AntiSniperModule.ZeroAddress.selector);
        module.configure(address(0), START_TIME, DURATION, START_BPS, END_BPS);
    }

    function test_configure_reverts_startBps_over_max() public {
        uint16 bad = 10_001;
        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(AntiSniperModule.InvalidBps.selector, bad));
        module.configure(agent, START_TIME, DURATION, bad, END_BPS);
    }

    function test_configure_reverts_endBps_over_startBps() public {
        uint16 badEnd = START_BPS + 1;
        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(AntiSniperModule.InvalidBps.selector, badEnd));
        module.configure(agent, START_TIME, DURATION, START_BPS, badEnd);
    }

    function test_configure_reverts_zero_duration() public {
        vm.prank(owner);
        vm.expectRevert(AntiSniperModule.InvalidDuration.selector);
        module.configure(agent, START_TIME, 0, START_BPS, END_BPS);
    }

    function test_configure_emits_event_and_writes_storage() public {
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
    // currentBps — piecewise-linear decay
    // ---------------------------------------------------------------------

    function test_currentBps_unconfigured_returns_zero() public view {
        assertEq(uint256(module.currentBps(agent)), 0);
    }

    function test_currentBps_before_start_returns_startBps() public {
        _configureDefault();
        // setUp warped to START_TIME - 1_000; still pre-start.
        assertEq(uint256(module.currentBps(agent)), uint256(START_BPS));
    }

    function test_currentBps_at_start_returns_startBps() public {
        _configureDefault();
        vm.warp(START_TIME);
        assertEq(uint256(module.currentBps(agent)), uint256(START_BPS));
    }

    /// @dev 100% -> 0% over 1_000s. At 500s we expect ~50%.
    ///      With floor division: drop = (10_000 * 500) / 1_000 = 5_000.
    ///      currentBps = 10_000 - 5_000 = 5_000.
    function test_currentBps_midway_linear_interpolation() public {
        uint64 start = 2_000_000;
        uint32 dur = 1000;
        uint16 sBps = 10_000;
        uint16 eBps = 0;

        vm.prank(owner);
        module.configure(agent2, start, dur, sBps, eBps);

        vm.warp(start + 500);
        assertEq(uint256(module.currentBps(agent2)), 5000);
    }

    function test_currentBps_after_end_returns_endBps() public {
        _configureDefault();
        vm.warp(uint256(START_TIME) + uint256(DURATION));
        assertEq(uint256(module.currentBps(agent)), uint256(END_BPS));
        vm.warp(uint256(START_TIME) + uint256(DURATION) + 10_000);
        assertEq(uint256(module.currentBps(agent)), uint256(END_BPS));
    }

    // ---------------------------------------------------------------------
    // computeTax
    // ---------------------------------------------------------------------

    function test_computeTax_matches_currentBps_scale() public {
        _configureDefault();
        vm.warp(uint256(START_TIME) + 30); // halfway — 30/60 of (9900-100) = 4900 drop
        uint16 bps = module.currentBps(agent);
        assertEq(uint256(bps), 9900 - 4900); // = 5_000

        uint256 amount = 1_000_000 ether;
        (uint256 tax, uint256 net) = module.computeTax(agent, amount);
        assertEq(tax, (amount * uint256(bps)) / uint256(module.MAX_BPS()));
        assertEq(net, amount - tax);
    }

    function test_computeTax_unconfigured_is_zero_tax() public view {
        (uint256 tax, uint256 net) = module.computeTax(agent, 1000);
        assertEq(tax, 0);
        assertEq(net, 1000);
    }

    // ---------------------------------------------------------------------
    // Fuzz
    // ---------------------------------------------------------------------

    /// @dev currentBps is monotone non-increasing in time for any configured
    ///      agent. Bound fuzz inputs to the decay window and a small tail.
    function testFuzz_currentBps_monotonic_nonincreasing(uint32 elapsed1, uint32 elapsed2) public {
        _configureDefault();
        // Bound to [0, duration + 10k] so we also cover the post-end saturation.
        elapsed1 = uint32(bound(uint256(elapsed1), 0, uint256(DURATION) + 10_000));
        elapsed2 = uint32(bound(uint256(elapsed2), 0, uint256(DURATION) + 10_000));
        if (elapsed1 > elapsed2) {
            (elapsed1, elapsed2) = (elapsed2, elapsed1);
        }

        vm.warp(uint256(START_TIME) + uint256(elapsed1));
        uint16 b1 = module.currentBps(agent);
        vm.warp(uint256(START_TIME) + uint256(elapsed2));
        uint16 b2 = module.currentBps(agent);

        assertLe(uint256(b2), uint256(b1));
        assertLe(uint256(b1), uint256(START_BPS));
        assertGe(uint256(b2), uint256(END_BPS));
    }

    /// @dev For any amount bounded away from 2^256/MAX_BPS, tax + net == amount.
    function testFuzz_computeTax_sum_equals_amount(uint256 amount) public {
        _configureDefault();
        // Stay below 2^240 so amount * MAX_BPS (~14 bits) cannot overflow.
        amount = bound(amount, 0, type(uint256).max / uint256(module.MAX_BPS()));
        // Pick an arbitrary mid-window timestamp so the rate is non-trivial.
        vm.warp(uint256(START_TIME) + uint256(DURATION) / 3);

        (uint256 tax, uint256 net) = module.computeTax(agent, amount);
        assertEq(tax + net, amount);
    }

    // ---------------------------------------------------------------------
    // Helpers
    // ---------------------------------------------------------------------

    function _configureDefault() internal {
        vm.prank(owner);
        module.configure(agent, START_TIME, DURATION, START_BPS, END_BPS);
    }
}
