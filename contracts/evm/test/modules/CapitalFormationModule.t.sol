// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

import {CapitalFormationModule} from "../../src/launchpad/modules/CapitalFormationModule.sol";
import {MockUniswapV2Pair} from "../mocks/MockUniswapV2Pair.sol";

/// @dev Plain ERC-20 used as USDC + agent token stand-in.
contract CapFormERC20 is ERC20 {
    constructor(string memory n, string memory s) ERC20(n, s) {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @dev Payout token whose `transfer` re-enters {claim}. Used to verify the
///      reentrancy guard on the milestone-claim path.
contract CapFormReentrantToken is ERC20 {
    CapitalFormationModule public target;
    address public agent;
    uint8 public tier;
    bool public armed;

    constructor() ERC20("Reenter", "RNT") {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }

    function arm(CapitalFormationModule t, address a, uint8 ti) external {
        target = t;
        agent = a;
        tier = ti;
        armed = true;
    }

    function _update(address from, address to, uint256 value) internal override {
        super._update(from, to, value);
        if (armed && from == address(target)) {
            armed = false;
            target.claim(agent, tier);
        }
    }
}

/// @title  CapitalFormationModuleSuiteTest
/// @notice Per-contract unit suite for {CapitalFormationModule}: configure
///         validation (length, ordering, payout, pair binding, funding),
///         snapshot+claim TWAP path, every revert branch on claim, escrow
///         ledger drain, and a reentrancy probe via a malicious payout token.
contract CapitalFormationModuleSuiteTest is Test {
    CapitalFormationModule internal mod;
    MockUniswapV2Pair internal pair;
    CapFormERC20 internal usdc;
    CapFormERC20 internal agentTok;

    address internal owner = address(0xA0);
    address internal attacker = address(0xBAD);
    address internal creator = address(0xC0FFEE);
    address internal agentAdmin = address(0xAD);
    address internal agent = address(0xA1);

    uint256 internal constant T1 = 2_000_000e6;
    uint256 internal constant T2 = 10_000_000e6;
    uint256 internal constant T3 = 50_000_000e6;
    uint256 internal constant T4 = 160_000_000e6;
    uint256 internal constant P1 = 5000e6;
    uint256 internal constant P2 = 25_000e6;
    uint256 internal constant P3 = 100_000e6;
    uint256 internal constant P4 = 500_000e6;
    uint256 internal constant TOTAL_PAY = P1 + P2 + P3 + P4;
    uint256 internal constant SUPPLY = 1_000_000_000e18;

    function setUp() public {
        mod = new CapitalFormationModule(owner);
        usdc = new CapFormERC20("USDC", "USDC");
        agentTok = new CapFormERC20("Agent", "AGT");
        agentTok.mint(address(this), SUPPLY);

        pair = new MockUniswapV2Pair(address(agentTok), address(usdc));
        _syncByRole(100_000_000e18, 100_000e6);
    }

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    function test_constructor_zero_owner_reverts() public {
        vm.expectRevert(CapitalFormationModule.ZeroAddress.selector);
        new CapitalFormationModule(address(0));
    }

    function test_constructor_constants() public view {
        assertEq(uint256(mod.MAX_TIERS()), 4);
        assertEq(uint256(mod.MIN_TWAP_BLOCKS()), 32);
    }

    // ---------------------------------------------------------------------
    // configure
    // ---------------------------------------------------------------------

    function test_configure_writes_state_and_locks_escrow() public {
        usdc.mint(address(this), TOTAL_PAY);
        usdc.transfer(address(mod), TOTAL_PAY);

        vm.expectEmit(true, false, false, true, address(mod));
        emit CapitalFormationModule.Configured(
            agent, address(agentTok), address(usdc), address(pair), creator, agentAdmin, 4, TOTAL_PAY
        );

        vm.prank(owner);
        mod.configure(
            agent, address(agentTok), address(usdc), address(pair), creator, agentAdmin, _thresholds(), _payouts()
        );

        (address aTok, address pTok, address pr, bool isT0, address cr, address adm, bool configured, uint8 bitmap) =
            mod.configs(agent);
        assertEq(aTok, address(agentTok));
        assertEq(pTok, address(usdc));
        assertEq(pr, address(pair));
        assertEq(isT0, pair.token0() == address(agentTok));
        assertEq(cr, creator);
        assertEq(adm, agentAdmin);
        assertTrue(configured);
        assertEq(uint256(bitmap), 0);

        assertEq(mod.tierCount(agent), 4);
        (uint256 thr0, uint256 pay0) = mod.tierAt(agent, 0);
        assertEq(thr0, T1);
        assertEq(pay0, P1);
        assertEq(mod.outstanding(address(usdc)), TOTAL_PAY);
    }

    function test_configure_only_owner() public {
        vm.prank(attacker);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, attacker));
        mod.configure(
            agent, address(agentTok), address(usdc), address(pair), creator, agentAdmin, _thresholds(), _payouts()
        );
    }

    function test_configure_one_shot() public {
        _fundAndConfigure();
        vm.prank(owner);
        vm.expectRevert(CapitalFormationModule.AlreadyConfigured.selector);
        mod.configure(
            agent, address(agentTok), address(usdc), address(pair), creator, agentAdmin, _thresholds(), _payouts()
        );
    }

    function test_configure_zero_address_reverts() public {
        vm.prank(owner);
        vm.expectRevert(CapitalFormationModule.ZeroAddress.selector);
        mod.configure(
            address(0), address(agentTok), address(usdc), address(pair), creator, agentAdmin, _thresholds(), _payouts()
        );
    }

    function test_configure_length_mismatch_reverts() public {
        usdc.mint(address(this), TOTAL_PAY);
        usdc.transfer(address(mod), TOTAL_PAY);
        uint256[] memory thrs = _thresholds();
        uint256[] memory pays = new uint256[](3);
        pays[0] = P1;
        pays[1] = P2;
        pays[2] = P3;
        vm.prank(owner);
        vm.expectRevert(CapitalFormationModule.LengthMismatch.selector);
        mod.configure(agent, address(agentTok), address(usdc), address(pair), creator, agentAdmin, thrs, pays);
    }

    function test_configure_zero_tier_count_reverts() public {
        usdc.mint(address(this), TOTAL_PAY);
        usdc.transfer(address(mod), TOTAL_PAY);
        uint256[] memory empty = new uint256[](0);
        vm.prank(owner);
        vm.expectRevert(CapitalFormationModule.InvalidTierCount.selector);
        mod.configure(agent, address(agentTok), address(usdc), address(pair), creator, agentAdmin, empty, empty);
    }

    function test_configure_too_many_tiers_reverts() public {
        uint256[] memory thrs = new uint256[](5);
        uint256[] memory pays = new uint256[](5);
        for (uint256 i; i < 5; ++i) {
            thrs[i] = (i + 1) * 1e6;
            pays[i] = 1e6;
        }
        vm.prank(owner);
        vm.expectRevert(CapitalFormationModule.InvalidTierCount.selector);
        mod.configure(agent, address(agentTok), address(usdc), address(pair), creator, agentAdmin, thrs, pays);
    }

    function test_configure_thresholds_descending_reverts() public {
        usdc.mint(address(this), TOTAL_PAY);
        usdc.transfer(address(mod), TOTAL_PAY);
        uint256[] memory thrs = new uint256[](2);
        thrs[0] = T2;
        thrs[1] = T1;
        uint256[] memory pays = new uint256[](2);
        pays[0] = P1;
        pays[1] = P2;
        vm.prank(owner);
        vm.expectRevert(CapitalFormationModule.ThresholdsNotIncreasing.selector);
        mod.configure(agent, address(agentTok), address(usdc), address(pair), creator, agentAdmin, thrs, pays);
    }

    function test_configure_first_threshold_zero_reverts() public {
        usdc.mint(address(this), TOTAL_PAY);
        usdc.transfer(address(mod), TOTAL_PAY);
        uint256[] memory thrs = new uint256[](2);
        thrs[0] = 0;
        thrs[1] = T1;
        uint256[] memory pays = new uint256[](2);
        pays[0] = P1;
        pays[1] = P2;
        vm.prank(owner);
        vm.expectRevert(CapitalFormationModule.ThresholdsNotIncreasing.selector);
        mod.configure(agent, address(agentTok), address(usdc), address(pair), creator, agentAdmin, thrs, pays);
    }

    function test_configure_zero_payout_reverts() public {
        usdc.mint(address(this), TOTAL_PAY);
        usdc.transfer(address(mod), TOTAL_PAY);
        uint256[] memory thrs = _thresholds();
        uint256[] memory pays = new uint256[](4);
        pays[0] = P1;
        pays[1] = 0;
        pays[2] = P3;
        pays[3] = P4;
        vm.prank(owner);
        vm.expectRevert(CapitalFormationModule.ZeroPayout.selector);
        mod.configure(agent, address(agentTok), address(usdc), address(pair), creator, agentAdmin, thrs, pays);
    }

    function test_configure_pair_token_mismatch_reverts() public {
        usdc.mint(address(this), TOTAL_PAY);
        usdc.transfer(address(mod), TOTAL_PAY);
        CapFormERC20 stranger = new CapFormERC20("Stranger", "STR");
        vm.prank(owner);
        vm.expectRevert(CapitalFormationModule.PairTokenMismatch.selector);
        mod.configure(
            agent, address(stranger), address(usdc), address(pair), creator, agentAdmin, _thresholds(), _payouts()
        );
    }

    function test_configure_pair_uninitialized_reverts() public {
        usdc.mint(address(this), TOTAL_PAY);
        usdc.transfer(address(mod), TOTAL_PAY);
        MockUniswapV2Pair freshPair = new MockUniswapV2Pair(address(agentTok), address(usdc));
        vm.prank(owner);
        vm.expectRevert(CapitalFormationModule.PairUninitialized.selector);
        mod.configure(
            agent, address(agentTok), address(usdc), address(freshPair), creator, agentAdmin, _thresholds(), _payouts()
        );
    }

    function test_configure_insufficient_funding_reverts() public {
        usdc.mint(address(this), P1);
        usdc.transfer(address(mod), P1);
        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(CapitalFormationModule.InsufficientFunding.selector, TOTAL_PAY, P1));
        mod.configure(
            agent, address(agentTok), address(usdc), address(pair), creator, agentAdmin, _thresholds(), _payouts()
        );
    }

    // ---------------------------------------------------------------------
    // claim — happy paths
    // ---------------------------------------------------------------------

    function test_claim_each_tier_in_order() public {
        _fundAndConfigure();

        _setFdvUsd(3_000_000);
        _runTwap();
        mod.claim(agent, 0);
        assertEq(usdc.balanceOf(creator), P1);
        assertTrue(mod.isClaimed(agent, 0));

        _setFdvUsd(12_000_000);
        _runTwap();
        mod.claim(agent, 1);

        _setFdvUsd(60_000_000);
        _runTwap();
        mod.claim(agent, 2);

        _setFdvUsd(200_000_000);
        _runTwap();
        mod.claim(agent, 3);

        assertEq(usdc.balanceOf(creator), TOTAL_PAY);
        assertEq(mod.outstanding(address(usdc)), 0);
    }

    function test_claim_higher_tier_first() public {
        _fundAndConfigure();
        _setFdvUsd(60_000_000);
        _runTwap();
        mod.claim(agent, 2);
        assertEq(usdc.balanceOf(creator), P3);
        assertTrue(mod.isClaimed(agent, 2));
        assertFalse(mod.isClaimed(agent, 0));
    }

    // ---------------------------------------------------------------------
    // claim — reverts
    // ---------------------------------------------------------------------

    function test_claim_double_reverts() public {
        _fundAndConfigure();
        _setFdvUsd(3_000_000);
        _runTwap();
        mod.claim(agent, 0);

        _setFdvUsd(3_000_000);
        _runTwap();
        vm.expectRevert(CapitalFormationModule.TierAlreadyClaimed.selector);
        mod.claim(agent, 0);
    }

    function test_claim_below_threshold_reverts() public {
        _fundAndConfigure();
        _setFdvUsd(1_500_000);
        _runTwap();
        uint256 expected = _expectedFdvUsd(1_500_000);
        vm.expectRevert(abi.encodeWithSelector(CapitalFormationModule.FdvBelowThreshold.selector, expected, T1));
        mod.claim(agent, 0);
    }

    function test_claim_tier_out_of_range_reverts() public {
        _fundAndConfigure();
        _setFdvUsd(200_000_000);
        _runTwap();
        vm.expectRevert(CapitalFormationModule.TierOutOfRange.selector);
        mod.claim(agent, 4);
    }

    function test_claim_unconfigured_reverts() public {
        vm.expectRevert(CapitalFormationModule.NotConfigured.selector);
        mod.claim(agent, 0);
    }

    function test_claim_no_snapshot_reverts() public {
        _fundAndConfigure();
        _setFdvUsd(3_000_000);
        vm.expectRevert(CapitalFormationModule.SnapshotMissing.selector);
        mod.claim(agent, 0);
    }

    function test_claim_window_too_short_reverts() public {
        _fundAndConfigure();
        _setFdvUsd(3_000_000);

        mod.snapshot(agent);
        vm.roll(block.number + 10);
        vm.warp(block.timestamp + 20);
        _setFdvUsd(3_000_000);

        vm.expectRevert(
            abi.encodeWithSelector(CapitalFormationModule.TwapWindowTooShort.selector, uint64(10), uint64(32))
        );
        mod.claim(agent, 0);
    }

    function test_claim_pair_too_young_reverts() public {
        _fundAndConfigure();
        _setFdvUsd(3_000_000);

        mod.snapshot(agent);
        vm.roll(block.number + 50);
        vm.warp(block.timestamp + 100);

        vm.expectRevert(CapitalFormationModule.PairTooYoung.selector);
        mod.claim(agent, 0);
    }

    // ---------------------------------------------------------------------
    // snapshot
    // ---------------------------------------------------------------------

    function test_snapshot_unconfigured_reverts() public {
        vm.expectRevert(CapitalFormationModule.NotConfigured.selector);
        mod.snapshot(agent);
    }

    function test_snapshot_overwrite() public {
        _fundAndConfigure();
        _setFdvUsd(3_000_000);
        mod.snapshot(agent);
        (,, uint64 first,) = mod.snapshots(agent);
        vm.roll(block.number + 5);
        vm.warp(block.timestamp + 10);
        _setFdvUsd(3_000_000);
        mod.snapshot(agent);
        (,, uint64 second, bool exists) = mod.snapshots(agent);
        assertGt(uint256(second), uint256(first));
        assertTrue(exists);
    }

    // ---------------------------------------------------------------------
    // Reentrancy probe
    // ---------------------------------------------------------------------

    function test_claim_reentrancy_blocked() public {
        CapFormReentrantToken evilUsdc = new CapFormReentrantToken();
        evilUsdc.mint(address(this), TOTAL_PAY);
        evilUsdc.transfer(address(mod), TOTAL_PAY);

        vm.prank(owner);
        mod.configure(
            agent, address(agentTok), address(evilUsdc), address(pair), creator, agentAdmin, _thresholds(), _payouts()
        );

        _setFdvUsd(3_000_000);
        _runTwap();
        evilUsdc.arm(mod, agent, 0);

        vm.expectRevert();
        mod.claim(agent, 0);

        assertEq(evilUsdc.balanceOf(creator), 0);
        assertFalse(mod.isClaimed(agent, 0));
    }

    // ---------------------------------------------------------------------
    // Views
    // ---------------------------------------------------------------------

    function test_isClaimed_out_of_range_false() public view {
        assertFalse(mod.isClaimed(agent, 7));
    }

    function test_tierAt_out_of_range_reverts() public {
        _fundAndConfigure();
        vm.expectRevert(CapitalFormationModule.TierOutOfRange.selector);
        mod.tierAt(agent, 5);
    }

    function test_outstanding_drains_with_claims() public {
        _fundAndConfigure();
        assertEq(mod.outstanding(address(usdc)), TOTAL_PAY);
        _setFdvUsd(3_000_000);
        _runTwap();
        mod.claim(agent, 0);
        assertEq(mod.outstanding(address(usdc)), TOTAL_PAY - P1);
    }

    // ---------------------------------------------------------------------
    // Helpers
    // ---------------------------------------------------------------------

    function _thresholds() internal pure returns (uint256[] memory thrs) {
        thrs = new uint256[](4);
        thrs[0] = T1;
        thrs[1] = T2;
        thrs[2] = T3;
        thrs[3] = T4;
    }

    function _payouts() internal pure returns (uint256[] memory pays) {
        pays = new uint256[](4);
        pays[0] = P1;
        pays[1] = P2;
        pays[2] = P3;
        pays[3] = P4;
    }

    function _fundAndConfigure() internal {
        usdc.mint(address(this), TOTAL_PAY);
        usdc.transfer(address(mod), TOTAL_PAY);
        vm.prank(owner);
        mod.configure(
            agent, address(agentTok), address(usdc), address(pair), creator, agentAdmin, _thresholds(), _payouts()
        );
    }

    function _setFdvUsd(uint256 targetFdvUsdMillionsScale) internal {
        uint256 fdvUsdc = targetFdvUsdMillionsScale * 1e6;
        uint256 agentReserve = 100_000_000e18;
        uint256 usdcReserve = fdvUsdc * agentReserve / SUPPLY;
        _syncByRole(agentReserve, usdcReserve);
    }

    function _expectedFdvUsd(uint256 targetFdvUsdMillionsScale) internal pure returns (uint256) {
        uint256 fdvUsdc = targetFdvUsdMillionsScale * 1e6;
        uint256 agentReserve = 100_000_000e18;
        uint256 usdcReserve = fdvUsdc * agentReserve / SUPPLY;
        uint256 q112 = 0x10000000000000000000000000000;
        uint256 priceQ112 = (usdcReserve * q112) / agentReserve;
        return (priceQ112 * SUPPLY) / q112;
    }

    function _syncByRole(uint256 agentAmt, uint256 usdcAmt) internal {
        if (pair.token0() == address(agentTok)) {
            pair.sync(uint112(agentAmt), uint112(usdcAmt));
        } else {
            pair.sync(uint112(usdcAmt), uint112(agentAmt));
        }
    }

    function _runTwap() internal {
        (uint112 r0, uint112 r1,) = pair.getReserves();
        vm.warp(block.timestamp + 64);
        pair.sync(r0, r1);
        mod.snapshot(agent);
        vm.roll(block.number + 32);
        vm.warp(block.timestamp + 64);
        pair.sync(r0, r1);
    }
}
