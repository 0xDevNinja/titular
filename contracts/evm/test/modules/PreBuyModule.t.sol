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

/// @title  PreBuyModuleSuiteTest
/// @notice Per-contract unit suite for {PreBuyModule}: owner-gated one-shot
///         configure path, deterministic vesting clone deployment, full
///         vesting curve traversal via the clone, sweep, and the companion
///         {TokenVesting} initialize validation.
contract PreBuyModuleSuiteTest is Test {
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
    uint256 internal constant VEST_AMOUNT = 100_000_000e18;
    uint64 internal constant CLIFF = 30 days;
    uint64 internal constant DURATION = 180 days;

    function setUp() public {
        titu = new MockMintableERC20("Titan", "TITU");
        agentA = new MockMintableERC20("AgentA", "AGTA");
        agentB = new MockMintableERC20("AgentB", "AGTB");
        curveA = new MockBondingCurve(address(agentA), address(titu));
        curveB = new MockBondingCurve(address(agentB), address(titu));

        agentA.mint(address(curveA), 1_000_000_000e18);
        agentB.mint(address(curveB), 1_000_000_000e18);

        module = new PreBuyModule(owner);
        titu.mint(address(module), 10 * TITAN_IN);
        vm.warp(1_700_000_000);
    }

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    function test_constructor_zero_owner_reverts() public {
        vm.expectRevert(PreBuyModule.ZeroAddress.selector);
        new PreBuyModule(address(0));
    }

    function test_constructor_deploys_implementation() public view {
        address impl = module.vestingImplementation();
        assertTrue(impl != address(0));
    }

    function test_implementation_locked_against_initialize() public {
        TokenVesting impl = TokenVesting(module.vestingImplementation());
        vm.expectRevert();
        impl.initialize(IERC20(address(titu)), creator, 1, uint64(block.timestamp), 0, 1);
    }

    // ---------------------------------------------------------------------
    // configure — access + validation
    // ---------------------------------------------------------------------

    function test_configure_only_owner() public {
        vm.prank(stranger);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, stranger));
        module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
    }

    function test_configure_zero_agent_reverts() public {
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.ZeroAddress.selector);
        module.configure(address(0), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA)));
    }

    function test_configure_zero_creator_reverts() public {
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.ZeroAddress.selector);
        module.configure(
            address(agentA), address(0), VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
    }

    function test_configure_zero_curve_reverts() public {
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.ZeroAddress.selector);
        module.configure(address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(0)));
    }

    function test_configure_zero_vest_amount_reverts() public {
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.ZeroAmount.selector);
        module.configure(address(agentA), creator, 0, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA)));
    }

    function test_configure_zero_titan_in_reverts() public {
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.ZeroAmount.selector);
        module.configure(address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, 0, IBondingCurve(address(curveA)));
    }

    function test_configure_zero_duration_reverts() public {
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.ZeroAmount.selector);
        module.configure(address(agentA), creator, VEST_AMOUNT, CLIFF, 0, TITAN_IN, IBondingCurve(address(curveA)));
    }

    function test_configure_cliff_above_duration_reverts() public {
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.CliffExceedsDuration.selector);
        module.configure(
            address(agentA), creator, VEST_AMOUNT, DURATION + 1, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
    }

    function test_configure_agent_mismatch_reverts() public {
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.AgentMismatch.selector);
        module.configure(
            address(agentB), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
    }

    function test_configure_one_shot() public {
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

    function test_configure_under_delivery_reverts() public {
        curveA.setPayoutOverride(VEST_AMOUNT - 1);
        vm.prank(owner);
        vm.expectRevert(PreBuyModule.InsufficientBuyOutput.selector);
        module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
    }

    function test_configure_propagates_curve_revert() public {
        curveA.setRevertOnBuy(true);
        vm.prank(owner);
        vm.expectRevert(bytes("curve: forced revert"));
        module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
    }

    // ---------------------------------------------------------------------
    // configure — happy path
    // ---------------------------------------------------------------------

    function test_configure_writes_state_and_funds_clone() public {
        address predicted = module.predictCloneAddress(address(agentA));
        vm.expectEmit(true, true, true, true, address(module));
        emit PreBuyModule.Configured(address(agentA), creator, predicted, VEST_AMOUNT, CLIFF, DURATION);

        vm.prank(owner);
        address clone = module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
        assertEq(clone, predicted);

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

        assertEq(titu.balanceOf(address(curveA)), TITAN_IN);
        assertEq(agentA.balanceOf(clone), VEST_AMOUNT);
        assertEq(agentA.balanceOf(address(module)), 0);

        TokenVesting v = TokenVesting(clone);
        assertEq(address(v.token()), address(agentA));
        assertEq(v.beneficiary(), creator);
        assertEq(v.total(), VEST_AMOUNT);
    }

    function test_configure_no_residual_allowance() public {
        vm.prank(owner);
        module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
        assertEq(titu.allowance(address(module), address(curveA)), 0);
    }

    function test_configure_over_delivery_funds_full_amount() public {
        uint256 over = VEST_AMOUNT + 12_345e18;
        curveA.setPayoutOverride(over);
        vm.prank(owner);
        address clone = module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
        assertEq(agentA.balanceOf(clone), over);
        assertEq(TokenVesting(clone).total(), over);
    }

    // ---------------------------------------------------------------------
    // Vesting curve traversal
    // ---------------------------------------------------------------------

    function _configureA() internal returns (TokenVesting) {
        vm.prank(owner);
        address clone = module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
        return TokenVesting(clone);
    }

    function test_release_pre_cliff_zero_reverts() public {
        TokenVesting v = _configureA();
        vm.warp(block.timestamp + uint256(CLIFF) - 1);
        assertEq(v.releasable(), 0);
        vm.expectRevert(TokenVesting.NothingToRelease.selector);
        v.release();
    }

    function test_release_at_cliff() public {
        TokenVesting v = _configureA();
        vm.warp(block.timestamp + uint256(CLIFF));
        uint256 expected = (VEST_AMOUNT * uint256(CLIFF)) / uint256(DURATION);
        assertEq(v.releasable(), expected);
        v.release();
        assertEq(agentA.balanceOf(creator), expected);

        vm.expectRevert(TokenVesting.NothingToRelease.selector);
        v.release();
    }

    function test_release_at_duration_full() public {
        TokenVesting v = _configureA();
        vm.warp(block.timestamp + uint256(DURATION));
        assertEq(v.releasable(), VEST_AMOUNT);
        v.release();
        assertEq(agentA.balanceOf(creator), VEST_AMOUNT);
    }

    function test_release_far_future_caps_total() public {
        TokenVesting v = _configureA();
        vm.warp(block.timestamp + 10 * uint256(DURATION));
        v.release();
        assertEq(agentA.balanceOf(creator), VEST_AMOUNT);
    }

    function test_release_split_releases_consistent() public {
        TokenVesting v = _configureA();
        vm.warp(block.timestamp + uint256(DURATION) / 4);
        v.release();
        uint256 first = agentA.balanceOf(creator);

        vm.warp(block.timestamp + uint256(DURATION) / 4);
        v.release();
        uint256 totalReleased = agentA.balanceOf(creator);
        assertEq(totalReleased, v.released());
        assertGt(totalReleased, first);
    }

    function test_release_via_module_emits_and_forwards() public {
        TokenVesting v = _configureA();
        vm.warp(block.timestamp + uint256(DURATION));

        vm.expectEmit(true, true, false, true, address(module));
        emit PreBuyModule.Released(address(agentA), address(v), VEST_AMOUNT);
        uint256 amount = module.release(address(agentA));
        assertEq(amount, VEST_AMOUNT);
    }

    function test_release_via_module_unconfigured_reverts() public {
        vm.expectRevert(PreBuyModule.NotConfigured.selector);
        module.release(address(agentB));
    }

    // ---------------------------------------------------------------------
    // Clone independence
    // ---------------------------------------------------------------------

    function test_two_agents_distinct_clones() public {
        vm.startPrank(owner);
        address cA = module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
        address cB = module.configure(
            address(agentB), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveB))
        );
        vm.stopPrank();

        assertTrue(cA != cB);
        assertEq(agentA.balanceOf(cA), VEST_AMOUNT);
        assertEq(agentB.balanceOf(cB), VEST_AMOUNT);
        assertEq(agentA.balanceOf(cB), 0);
        assertEq(agentB.balanceOf(cA), 0);
    }

    function test_predictCloneAddress_matches_deployment() public {
        address predicted = module.predictCloneAddress(address(agentA));
        vm.prank(owner);
        address actual = module.configure(
            address(agentA), creator, VEST_AMOUNT, CLIFF, DURATION, TITAN_IN, IBondingCurve(address(curveA))
        );
        assertEq(predicted, actual);
    }

    // ---------------------------------------------------------------------
    // sweep
    // ---------------------------------------------------------------------

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

    // ---------------------------------------------------------------------
    // TokenVesting initialize validation (bare clones)
    // ---------------------------------------------------------------------

    function _bareClone() internal returns (TokenVesting) {
        return TokenVesting(Clones.clone(module.vestingImplementation()));
    }

    function test_tokenVesting_init_zero_token_reverts() public {
        TokenVesting v = _bareClone();
        vm.expectRevert(TokenVesting.ZeroAddress.selector);
        v.initialize(IERC20(address(0)), creator, 1, uint64(block.timestamp), 0, DURATION);
    }

    function test_tokenVesting_init_zero_beneficiary_reverts() public {
        TokenVesting v = _bareClone();
        vm.expectRevert(TokenVesting.ZeroAddress.selector);
        v.initialize(IERC20(address(agentA)), address(0), 1, uint64(block.timestamp), 0, DURATION);
    }

    function test_tokenVesting_init_zero_total_reverts() public {
        TokenVesting v = _bareClone();
        vm.expectRevert(TokenVesting.ZeroAmount.selector);
        v.initialize(IERC20(address(agentA)), creator, 0, uint64(block.timestamp), 0, DURATION);
    }

    function test_tokenVesting_init_zero_duration_reverts() public {
        TokenVesting v = _bareClone();
        vm.expectRevert(TokenVesting.ZeroDuration.selector);
        v.initialize(IERC20(address(agentA)), creator, 1, uint64(block.timestamp), 0, 0);
    }

    function test_tokenVesting_init_cliff_above_duration_reverts() public {
        TokenVesting v = _bareClone();
        vm.expectRevert(TokenVesting.CliffExceedsDuration.selector);
        v.initialize(IERC20(address(agentA)), creator, 1, uint64(block.timestamp), DURATION + 1, DURATION);
    }

    // ---------------------------------------------------------------------
    // Fuzz
    // ---------------------------------------------------------------------

    /// @dev `vested` is monotone non-decreasing in time and bounded above by total.
    function testFuzz_vested_bounded_monotone(uint64 t1, uint64 t2) public {
        TokenVesting v = _configureA();
        uint256 baseline = block.timestamp;
        uint256 cap = uint256(DURATION) * 2;
        uint256 e1 = bound(uint256(t1), 0, cap);
        uint256 e2 = bound(uint256(t2), 0, cap);
        if (e1 > e2) (e1, e2) = (e2, e1);

        vm.warp(baseline + e1);
        uint256 v1 = v.vested();
        vm.warp(baseline + e2);
        uint256 v2 = v.vested();
        assertLe(v1, v2);
        assertLe(v2, VEST_AMOUNT);
    }
}
