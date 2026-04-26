// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {Clones} from "@openzeppelin/contracts/proxy/Clones.sol";

import {PreBuyModule} from "../../src/launchpad/modules/PreBuyModule.sol";
import {TokenVesting} from "../../src/launchpad/modules/TokenVesting.sol";
import {IBondingCurve} from "../../src/launchpad/interfaces/IBondingCurve.sol";
import {MockMintableERC20} from "../mocks/MockERC20.sol";
import {MockBondingCurve} from "../mocks/MockBondingCurve.sol";

/// @dev Unit coverage for {PreBuyModule} + its companion {TokenVesting}
///      clone target. Validates: ownership + one-shot + arg validation on
///      {configure}; correct funding of the deterministic clone; the full
///      vesting curve (pre-cliff zero, cliff pro-rata, post-duration full,
///      idempotent over-release); clone independence across agents.
contract PreBuyModuleTest is Test {
    PreBuyModule internal module;
    MockMintableERC20 internal titu;
    MockMintableERC20 internal agentA;
    MockMintableERC20 internal agentB;
    MockBondingCurve internal curveA;
    MockBondingCurve internal curveB;

    address internal owner = address(0xF1);
    address internal creator = address(0xC0FFEE);
    address internal stranger = address(0xBAD);

    uint256 internal constant TITAN_IN = 1000e18;
    uint256 internal constant VEST_AMOUNT = 100_000_000e18; // 10% of 1B supply
    uint64 internal constant CLIFF = 30 days;
    uint64 internal constant DURATION = 180 days;

    // Re-declare events for `vm.expectEmit` matchers.
    event Configured(
        address indexed agent,
        address indexed creator,
        address indexed vestingClone,
        uint256 vestAmount,
        uint64 cliffSeconds,
        uint64 durationSeconds
    );
    event Released(address indexed agent, address indexed clone, uint256 amount);

    function setUp() public {
        titu = new MockMintableERC20("Titan", "TITU");
        agentA = new MockMintableERC20("AgentA", "AGTA");
        agentB = new MockMintableERC20("AgentB", "AGTB");

        curveA = new MockBondingCurve(address(agentA), address(titu));
        curveB = new MockBondingCurve(address(agentB), address(titu));

        // Pre-fund the mock curves with the agent token they will deliver on
        // buy. The full 1B supply mirrors the real launchpad's curve seed.
        agentA.mint(address(curveA), 1_000_000_000e18);
        agentB.mint(address(curveB), 1_000_000_000e18);

        module = new PreBuyModule(owner);

        // Park enough TITU on the module up front for both planned configures.
        titu.mint(address(module), 10 * TITAN_IN);

        // Stable wall-clock so vesting timestamps are easy to reason about.
        vm.warp(1_700_000_000);
    }

    // ----- Constructor -----

    function test_constructor_reverts_zero_owner() public {
        vm.expectRevert(PreBuyModule.ZeroAddress.selector);
        new PreBuyModule(address(0));
    }

    function test_constructor_deploys_implementation() public view {
        address impl = module.vestingImplementation();
        assertTrue(impl != address(0));
    }

    function test_implementation_is_initialization_locked() public {
        TokenVesting impl = TokenVesting(module.vestingImplementation());
        vm.expectRevert(); // OZ Initializable: InvalidInitialization
        impl.initialize(IERC20(address(titu)), creator, 1, uint64(block.timestamp), 0, 1);
    }

    // ----- configure: access control + validation -----

    function test_configure_only_owner() public {
        vm.prank(stranger);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, stranger));
        module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
    }

    function test_configure_reverts_zero_agent() public {
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.ZeroAddress.selector);
        module.configure(address(0), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA)));
    }

    function test_configure_reverts_zero_creator() public {
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.ZeroAddress.selector);
        module.configure(
            address(agentA), address(0), VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
    }

    function test_configure_reverts_zero_curve() public {
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.ZeroAddress.selector);
        module.configure(address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(0)));
    }

    function test_configure_reverts_zero_vest_amount() public {
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.ZeroAmount.selector);
        module.configure(address(agentA), creator, 0, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA)));
    }

    function test_configure_reverts_zero_titan_in() public {
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.ZeroAmount.selector);
        module.configure(address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, 0, IBondingCurve(address(curveA)));
    }

    function test_configure_reverts_zero_duration() public {
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.ZeroAmount.selector);
        module.configure(address(agentA), creator, VEST_AMOUNT, CLIFF, 0, TITAN_IN, IBondingCurve(address(curveA)));
    }

    function test_configure_reverts_cliff_gt_duration() public {
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.CliffExceedsDuration.selector);
        module.configure(
            address(agentA), creator, VEST_AMOUNT, DURATION + 1, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
    }

    function test_configure_reverts_agent_mismatch() public {
        // Curve's agentToken() reports agentA, but we claim agentB is the agent.
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.AgentMismatch.selector);
        module.configure(
            address(agentB), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
    }

    function test_configure_reverts_double_config() public {
        vm.startPrank(owner);
        module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
        vm.expectRevert(PreBuyModule.AlreadyConfigured.selector);
        module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
        vm.stopPrank();
    }

    function test_configure_reverts_curve_under_funds() public {
        // Force the mock curve to deliver less than the slippage floor so the
        // module's defensive check fires.
        curveA.setPayoutOverride(VEST_AMOUNT - 1);
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.InsufficientBuyOutput.selector);
        module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
    }

    function test_configure_bubbles_curve_revert() public {
        curveA.setRevertOnBuy(true);
        vm.prank(owner);
        vm.expectRevert(bytes("curve: forced revert"));
        module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
    }

    // ----- configure: happy path -----

    function test_configure_happy_path_writes_state_and_funds_clone() public {
        address predicted = module.predictCloneAddress(address(agentA));

        vm.expectEmit(true, true, true, true, address(module));
        emit Configured(address(agentA), creator, predicted, VEST_AMOUNT, CLIFF, DURATION);

        vm.prank(owner);
        address clone = module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );

        assertEq(clone, predicted, "clone deterministic mismatch");

        (
            address gotCreator,
            address gotClone,
            uint256 gotVest,
            uint64 gotCliff,
            uint64 gotDur,
            uint64 gotStart,
            bool gotConfigured
        ) = module.configs(address(agentA));
        assertEq(gotCreator, creator);
        assertEq(gotClone, clone);
        assertEq(gotVest, VEST_AMOUNT);
        assertEq(uint256(gotCliff), uint256(CLIFF));
        assertEq(uint256(gotDur), uint256(DURATION));
        assertEq(uint256(gotStart), block.timestamp);
        assertTrue(gotConfigured);

        // Curve received TITU; clone received agent tokens.
        assertEq(titu.balanceOf(address(curveA)), TITAN_IN);
        assertEq(agentA.balanceOf(clone), VEST_AMOUNT);
        // Module retained no agent tokens.
        assertEq(agentA.balanceOf(address(module)), 0);

        // Clone initialized correctly.
        TokenVesting v = TokenVesting(clone);
        assertEq(address(v.token()), address(agentA));
        assertEq(v.beneficiary(), creator);
        assertEq(v.total(), VEST_AMOUNT);
        assertEq(uint256(v.cliff()), uint256(CLIFF));
        assertEq(uint256(v.duration()), uint256(DURATION));

        // Clone is itself initialization-locked against re-init.
        vm.expectRevert();
        v.initialize(IERC20(address(agentA)), creator, VEST_AMOUNT, uint64(block.timestamp), CLIFF, DURATION);
    }

    function test_configure_no_residual_allowance() public {
        vm.prank(owner);
        module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
        assertEq(titu.allowance(address(module), address(curveA)), 0);
    }

    function test_configure_curve_over_delivers_full_balance_into_clone() public {
        // Curve delivers slightly more than the floor (e.g. real bonding-curve
        // math rounding up). The module forwards the full received amount.
        uint256 over = VEST_AMOUNT + 12_345e18;
        curveA.setPayoutOverride(over);

        vm.prank(owner);
        address clone = module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );

        assertEq(agentA.balanceOf(clone), over);
        TokenVesting v = TokenVesting(clone);
        assertEq(v.total(), over);
    }

    // ----- vesting math via clone -----

    function _configureA() internal returns (TokenVesting) {
        vm.prank(owner);
        address clone = module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
        return TokenVesting(clone);
    }

    function test_release_before_cliff_is_zero() public {
        TokenVesting v = _configureA();
        vm.warp(block.timestamp + uint256(CLIFF) - 1);
        assertEq(v.releasable(), 0);
        vm.expectRevert(TokenVesting.NothingToRelease.selector);
        v.release();
        assertEq(agentA.balanceOf(creator), 0);
    }

    function test_release_at_cliff_pro_rata() public {
        TokenVesting v = _configureA();
        // At exactly `start + cliff` the linear curve has accrued
        // `total * cliff / duration` (clone uses elapsed-from-start basis).
        vm.warp(block.timestamp + uint256(CLIFF));
        uint256 expected = (VEST_AMOUNT * uint256(CLIFF)) / uint256(DURATION);
        assertEq(v.releasable(), expected);

        v.release();
        assertEq(agentA.balanceOf(creator), expected);
        assertEq(v.released(), expected);
        // Immediate re-release is a no-op revert (no time passed).
        vm.expectRevert(TokenVesting.NothingToRelease.selector);
        v.release();
    }

    function test_release_after_duration_full() public {
        TokenVesting v = _configureA();
        vm.warp(block.timestamp + uint256(DURATION));
        assertEq(v.releasable(), VEST_AMOUNT);
        v.release();
        assertEq(agentA.balanceOf(creator), VEST_AMOUNT);
        assertEq(v.released(), VEST_AMOUNT);
    }

    function test_release_after_duration_far_future_still_caps_at_total() public {
        TokenVesting v = _configureA();
        vm.warp(block.timestamp + 10 * uint256(DURATION));
        v.release();
        assertEq(agentA.balanceOf(creator), VEST_AMOUNT);
    }

    function test_release_double_release_safe() public {
        TokenVesting v = _configureA();
        vm.warp(block.timestamp + uint256(DURATION) / 4);
        v.release();
        uint256 first = agentA.balanceOf(creator);

        vm.warp(block.timestamp + uint256(DURATION) / 4);
        v.release();
        uint256 second = agentA.balanceOf(creator) - first;

        // Two quarters released; total should be ~half.
        uint256 totalReleased = agentA.balanceOf(creator);
        assertEq(totalReleased, v.released());
        assertApproxEqAbs(totalReleased, VEST_AMOUNT / 2, 1);
        assertGt(second, 0);
    }

    function test_release_via_module_emits_and_forwards() public {
        TokenVesting v = _configureA();
        vm.warp(block.timestamp + uint256(DURATION));

        vm.expectEmit(true, true, false, true, address(module));
        emit Released(address(agentA), address(v), VEST_AMOUNT);
        uint256 amount = module.release(address(agentA));
        assertEq(amount, VEST_AMOUNT);
        assertEq(agentA.balanceOf(creator), VEST_AMOUNT);
    }

    function test_release_via_module_reverts_unconfigured() public {
        vm.expectRevert(PreBuyModule.NotConfigured.selector);
        module.release(address(agentB));
    }

    // ----- clone independence -----

    function test_configure_two_agents_yields_distinct_clones() public {
        vm.startPrank(owner);
        address cloneA = module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
        address cloneB = module.configure(
            address(agentB), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveB))
        );
        vm.stopPrank();

        assertTrue(cloneA != cloneB);
        // Each clone holds its own bag.
        assertEq(agentA.balanceOf(cloneA), VEST_AMOUNT);
        assertEq(agentB.balanceOf(cloneB), VEST_AMOUNT);
        // Cross-funding is zero.
        assertEq(agentA.balanceOf(cloneB), 0);
        assertEq(agentB.balanceOf(cloneA), 0);

        // A release on cloneA does not touch cloneB.
        vm.warp(block.timestamp + uint256(DURATION));
        TokenVesting(cloneA).release();
        assertEq(agentA.balanceOf(creator), VEST_AMOUNT);
        assertEq(agentB.balanceOf(creator), 0);
        TokenVesting(cloneB).release();
        assertEq(agentB.balanceOf(creator), VEST_AMOUNT);
    }

    function test_predictCloneAddress_matches_deployment() public {
        address predicted = module.predictCloneAddress(address(agentA));
        vm.prank(owner);
        address actual = module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
        assertEq(predicted, actual);
    }

    // ----- sweep -----

    function test_sweep_only_owner() public {
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, address(this)));
        module.sweep(IERC20(address(titu)), address(this), 1);
    }

    function test_sweep_zero_to_reverts() public {
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.ZeroAddress.selector);
        module.sweep(IERC20(address(titu)), address(0), 1);
    }

    function test_sweep_drains_residual_titu() public {
        // The setUp pre-funded the module with 10x TITAN_IN. After one
        // configure, 9x should remain.
        vm.prank(owner);
        module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
        uint256 residual = titu.balanceOf(address(module));
        assertEq(residual, 9 * TITAN_IN);

        vm.prank(owner);
        module.sweep(IERC20(address(titu)), owner, residual);
        assertEq(titu.balanceOf(owner), residual);
        assertEq(titu.balanceOf(address(module)), 0);
    }

    // ----- TokenVesting initialize validation (via fresh clones) -----
    //
    // The implementation contract is initializer-locked, so to exercise
    // initialize's revert branches we deploy bare clones in the test and
    // call them directly.

    function _bareClone() internal returns (TokenVesting) {
        return TokenVesting(Clones.clone(module.vestingImplementation()));
    }

    function test_tokenVesting_initialize_reverts_zero_token() public {
        TokenVesting v = _bareClone();
        vm.expectRevert(TokenVesting.ZeroAddress.selector);
        v.initialize(IERC20(address(0)), creator, 1, uint64(block.timestamp), 0, DURATION);
    }

    function test_tokenVesting_initialize_reverts_zero_beneficiary() public {
        TokenVesting v = _bareClone();
        vm.expectRevert(TokenVesting.ZeroAddress.selector);
        v.initialize(IERC20(address(agentA)), address(0), 1, uint64(block.timestamp), 0, DURATION);
    }

    function test_tokenVesting_initialize_reverts_zero_total() public {
        TokenVesting v = _bareClone();
        vm.expectRevert(TokenVesting.ZeroAmount.selector);
        v.initialize(IERC20(address(agentA)), creator, 0, uint64(block.timestamp), 0, DURATION);
    }

    function test_tokenVesting_initialize_reverts_zero_duration() public {
        TokenVesting v = _bareClone();
        vm.expectRevert(TokenVesting.ZeroDuration.selector);
        v.initialize(IERC20(address(agentA)), creator, 1, uint64(block.timestamp), 0, 0);
    }

    function test_tokenVesting_initialize_reverts_cliff_gt_duration() public {
        TokenVesting v = _bareClone();
        vm.expectRevert(TokenVesting.CliffExceedsDuration.selector);
        v.initialize(IERC20(address(agentA)), creator, 1, uint64(block.timestamp), DURATION + 1, DURATION);
    }

    // ----- vesting view fuzz -----

    /// @dev `vested` is monotone non-decreasing in time and bounded above by
    ///      `total` for any timestamp choice.
    function testFuzz_vested_bounded_and_monotone(uint64 t1, uint64 t2) public {
        TokenVesting v = _configureA();
        uint256 baseline = block.timestamp;
        uint256 cap = uint256(DURATION) * 2;
        uint256 e1 = bound(uint256(t1), 0, cap);
        uint256 e2 = bound(uint256(t2), 0, cap);
        if (e1 > e2) (e1, e2) = (e2, e1);

        vm.warp(baseline + e1);
        uint256 vested1 = v.vested();
        vm.warp(baseline + e2);
        uint256 vested2 = v.vested();

        assertLe(vested1, vested2);
        assertLe(vested2, VEST_AMOUNT);
    }
}
