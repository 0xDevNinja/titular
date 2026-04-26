// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

import {CapitalFormationModule} from "../../src/launchpad/modules/CapitalFormationModule.sol";
import {MockUniswapV2Pair} from "../mocks/MockUniswapV2Pair.sol";

/// @dev Plain ERC-20 stand-in for USDC / agent token in unit tests.
contract TestERC20 is ERC20 {
    constructor(string memory n, string memory s) ERC20(n, s) {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @dev ERC-20 whose `transfer` re-enters the module. Used to confirm
///      {ReentrancyGuard} stops a malicious payout token from re-entering
///      {claim} during the outbound transfer.
contract ReentrantPayoutToken is ERC20 {
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
            armed = false; // single-shot to avoid recursion blow-up
            target.claim(agent, tier);
        }
    }
}

/// @dev Unit tests for `CapitalFormationModule`. Exercises every external
///      entrypoint, every revert path, the TWAP window, and the escrow
///      ledger. Coverage target ≥ 90%.
contract CapitalFormationModuleTest is Test {
    CapitalFormationModule internal mod;
    MockUniswapV2Pair internal pair;
    TestERC20 internal usdc;
    TestERC20 internal agentTok;

    address internal owner = address(0xA0);
    address internal attacker = address(0xBAD);
    address internal creator = address(0xC0FFEE);
    address internal agentAdmin = address(0xAD);
    address internal agent = address(0xA1);

    // FDV milestones: $2M, $10M, $50M, $160M (USDC = 6 decimals).
    uint256 internal constant T1 = 2_000_000e6;
    uint256 internal constant T2 = 10_000_000e6;
    uint256 internal constant T3 = 50_000_000e6;
    uint256 internal constant T4 = 160_000_000e6;

    // Per-tier USDC payouts (illustrative).
    uint256 internal constant P1 = 5000e6;
    uint256 internal constant P2 = 25_000e6;
    uint256 internal constant P3 = 100_000e6;
    uint256 internal constant P4 = 500_000e6;

    uint256 internal constant TOTAL_PAY = P1 + P2 + P3 + P4;

    // 1B agent supply, 18 decimals.
    uint256 internal constant SUPPLY = 1_000_000_000e18;

    function setUp() public {
        mod = new CapitalFormationModule(owner);
        usdc = new TestERC20("USD Coin", "USDC");
        agentTok = new TestERC20("Agent", "AGT");
        agentTok.mint(address(this), SUPPLY); // owner-of-supply doesn't matter; we only read totalSupply.

        // Build a pair with reserves seeded so {configure} passes the
        // PairUninitialized check. Spot ratio ≈ $0.001 per agent token →
        // FDV ≈ $1M (below T1 by default).
        pair = new MockUniswapV2Pair(address(agentTok), address(usdc));
        // Seed the pair with non-zero reserves so {configure} passes the
        // PairUninitialized check. Spot ratio matches a tiny "$1M FDV"
        // baseline; tests that need a specific FDV call {_setFdvUsd}.
        // Reserves must fit in uint112; well within bounds.
        _syncByRole(100_000_000e18, 100_000e6);
    }

    // ---------------------------------------------------------------------
    // configure — happy + reverts
    // ---------------------------------------------------------------------

    function test_configure_happy_path_records_state_and_locks_escrow() public {
        // Pre-fund.
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
        // token0 is whichever address is numerically smaller; we check
        // against the pair's own report.
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

    function test_configure_reverts_already_configured() public {
        _fundAndConfigure();
        vm.prank(owner);
        vm.expectRevert(CapitalFormationModule.AlreadyConfigured.selector);
        mod.configure(
            agent, address(agentTok), address(usdc), address(pair), creator, agentAdmin, _thresholds(), _payouts()
        );
    }

    function test_configure_reverts_zero_address() public {
        vm.prank(owner);
        vm.expectRevert(CapitalFormationModule.ZeroAddress.selector);
        mod.configure(
            address(0), address(agentTok), address(usdc), address(pair), creator, agentAdmin, _thresholds(), _payouts()
        );
    }

    function test_configure_reverts_length_mismatch() public {
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

    function test_configure_reverts_tier_count_zero() public {
        usdc.mint(address(this), TOTAL_PAY);
        usdc.transfer(address(mod), TOTAL_PAY);
        uint256[] memory empty = new uint256[](0);
        vm.prank(owner);
        vm.expectRevert(CapitalFormationModule.InvalidTierCount.selector);
        mod.configure(agent, address(agentTok), address(usdc), address(pair), creator, agentAdmin, empty, empty);
    }

    function test_configure_reverts_tier_count_exceeds_max() public {
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

    function test_configure_reverts_thresholds_not_increasing() public {
        usdc.mint(address(this), TOTAL_PAY);
        usdc.transfer(address(mod), TOTAL_PAY);
        uint256[] memory thrs = new uint256[](2);
        thrs[0] = T2;
        thrs[1] = T1; // descending
        uint256[] memory pays = new uint256[](2);
        pays[0] = P1;
        pays[1] = P2;
        vm.prank(owner);
        vm.expectRevert(CapitalFormationModule.ThresholdsNotIncreasing.selector);
        mod.configure(agent, address(agentTok), address(usdc), address(pair), creator, agentAdmin, thrs, pays);
    }

    function test_configure_reverts_first_threshold_zero() public {
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

    function test_configure_reverts_zero_payout() public {
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

    function test_configure_reverts_pair_token_mismatch() public {
        usdc.mint(address(this), TOTAL_PAY);
        usdc.transfer(address(mod), TOTAL_PAY);
        // Wrong agent token: pass an unrelated address.
        TestERC20 stranger = new TestERC20("Stranger", "STR");
        vm.prank(owner);
        vm.expectRevert(CapitalFormationModule.PairTokenMismatch.selector);
        mod.configure(
            agent, address(stranger), address(usdc), address(pair), creator, agentAdmin, _thresholds(), _payouts()
        );
    }

    function test_configure_reverts_pair_uninitialized() public {
        usdc.mint(address(this), TOTAL_PAY);
        usdc.transfer(address(mod), TOTAL_PAY);
        MockUniswapV2Pair freshPair = new MockUniswapV2Pair(address(agentTok), address(usdc));
        // Reserves are zero by default.
        vm.prank(owner);
        vm.expectRevert(CapitalFormationModule.PairUninitialized.selector);
        mod.configure(
            agent, address(agentTok), address(usdc), address(freshPair), creator, agentAdmin, _thresholds(), _payouts()
        );
    }

    function test_configure_reverts_insufficient_funding() public {
        // Underfund: only P1 deposited, but configuring all four tiers.
        usdc.mint(address(this), P1);
        usdc.transfer(address(mod), P1);
        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(CapitalFormationModule.InsufficientFunding.selector, TOTAL_PAY, P1));
        mod.configure(
            agent, address(agentTok), address(usdc), address(pair), creator, agentAdmin, _thresholds(), _payouts()
        );
    }

    // ---------------------------------------------------------------------
    // claim — happy paths per tier
    // ---------------------------------------------------------------------

    function test_claim_happy_each_tier_in_sequence() public {
        _fundAndConfigure();

        // Drive the pair to ~$3M FDV → tier 0 only.
        _setFdvUsd(3_000_000);
        _runTwap();
        vm.expectEmit(true, true, false, true, address(mod));
        emit CapitalFormationModule.MilestoneClaimed(agent, 0, _expectedFdvUsd(3_000_000), P1);
        mod.claim(agent, 0);
        assertEq(usdc.balanceOf(creator), P1);
        assertTrue(mod.isClaimed(agent, 0));

        // Drive to $12M → tier 1.
        _setFdvUsd(12_000_000);
        _runTwap();
        mod.claim(agent, 1);
        assertEq(usdc.balanceOf(creator), P1 + P2);
        assertTrue(mod.isClaimed(agent, 1));

        // Drive to $60M → tier 2.
        _setFdvUsd(60_000_000);
        _runTwap();
        mod.claim(agent, 2);
        assertEq(usdc.balanceOf(creator), P1 + P2 + P3);
        assertTrue(mod.isClaimed(agent, 2));

        // Drive to $200M → tier 3.
        _setFdvUsd(200_000_000);
        _runTwap();
        mod.claim(agent, 3);
        assertEq(usdc.balanceOf(creator), TOTAL_PAY);
        assertTrue(mod.isClaimed(agent, 3));

        // Escrow drained.
        assertEq(mod.outstanding(address(usdc)), 0);
        assertEq(usdc.balanceOf(address(mod)), 0);
    }

    function test_claim_higher_tier_first_pays_only_that_tier() public {
        _fundAndConfigure();
        _setFdvUsd(60_000_000); // covers T1, T2, T3 but tier 2 only this call.
        _runTwap();
        mod.claim(agent, 2);
        assertEq(usdc.balanceOf(creator), P3);
        assertTrue(mod.isClaimed(agent, 2));
        assertFalse(mod.isClaimed(agent, 0));
        assertFalse(mod.isClaimed(agent, 1));
    }

    function test_claim_reverts_double_claim() public {
        _fundAndConfigure();
        _setFdvUsd(3_000_000);
        _runTwap();
        mod.claim(agent, 0);

        // Need a fresh snapshot for the second attempt.
        _setFdvUsd(3_000_000);
        _runTwap();
        vm.expectRevert(CapitalFormationModule.TierAlreadyClaimed.selector);
        mod.claim(agent, 0);
    }

    function test_claim_reverts_fdv_below_threshold() public {
        _fundAndConfigure();
        // FDV $1.5M, below T1 $2M.
        _setFdvUsd(1_500_000);
        _runTwap();
        uint256 expectedFdv = _expectedFdvUsd(1_500_000);
        vm.expectRevert(abi.encodeWithSelector(CapitalFormationModule.FdvBelowThreshold.selector, expectedFdv, T1));
        mod.claim(agent, 0);
    }

    function test_claim_reverts_tier_out_of_range() public {
        _fundAndConfigure();
        _setFdvUsd(200_000_000);
        _runTwap();
        vm.expectRevert(CapitalFormationModule.TierOutOfRange.selector);
        mod.claim(agent, 4);
    }

    function test_claim_reverts_not_configured() public {
        vm.expectRevert(CapitalFormationModule.NotConfigured.selector);
        mod.claim(agent, 0);
    }

    function test_claim_reverts_snapshot_missing() public {
        _fundAndConfigure();
        _setFdvUsd(3_000_000);
        // No snapshot taken.
        vm.expectRevert(CapitalFormationModule.SnapshotMissing.selector);
        mod.claim(agent, 0);
    }

    function test_claim_reverts_window_too_short() public {
        _fundAndConfigure();
        _setFdvUsd(3_000_000);

        mod.snapshot(agent);
        // Advance only 10 blocks (< MIN_TWAP_BLOCKS).
        vm.roll(block.number + 10);
        vm.warp(block.timestamp + 20);
        // Re-sync the pair to update its block timestamp.
        _setFdvUsd(3_000_000);

        vm.expectRevert(
            abi.encodeWithSelector(CapitalFormationModule.TwapWindowTooShort.selector, uint64(10), uint64(32))
        );
        mod.claim(agent, 0);
    }

    function test_claim_reverts_pair_too_young_no_sync_after_snapshot() public {
        _fundAndConfigure();
        _setFdvUsd(3_000_000);

        mod.snapshot(agent);
        // Advance blocks past the window without syncing the pair so the
        // pair's blockTimestampLast does not advance.
        vm.roll(block.number + 50);
        vm.warp(block.timestamp + 100);

        vm.expectRevert(CapitalFormationModule.PairTooYoung.selector);
        mod.claim(agent, 0);
    }

    // ---------------------------------------------------------------------
    // snapshot
    // ---------------------------------------------------------------------

    function test_snapshot_reverts_not_configured() public {
        vm.expectRevert(CapitalFormationModule.NotConfigured.selector);
        mod.snapshot(agent);
    }

    function test_snapshot_overwrites_prior() public {
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
    // Reentrancy
    // ---------------------------------------------------------------------

    function test_claim_reentrancy_blocked() public {
        // Use a payout token that re-enters {claim}.
        ReentrantPayoutToken evilUsdc = new ReentrantPayoutToken();
        evilUsdc.mint(address(this), TOTAL_PAY);
        evilUsdc.transfer(address(mod), TOTAL_PAY);

        vm.prank(owner);
        mod.configure(
            agent, address(agentTok), address(evilUsdc), address(pair), creator, agentAdmin, _thresholds(), _payouts()
        );

        _setFdvUsd(3_000_000);
        _runTwap();
        evilUsdc.arm(mod, agent, 0);

        // The outer call should bubble the inner ReentrancyGuardReentrantCall
        // (OZ v5 selector) up through the malicious _update.
        vm.expectRevert();
        mod.claim(agent, 0);

        // No payout, no claimed bit.
        assertEq(evilUsdc.balanceOf(creator), 0);
        assertFalse(mod.isClaimed(agent, 0));
    }

    // ---------------------------------------------------------------------
    // Views
    // ---------------------------------------------------------------------

    function test_isClaimed_out_of_range_returns_false() public view {
        assertFalse(mod.isClaimed(agent, 7));
    }

    function test_tierAt_reverts_out_of_range() public {
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
    // Constructor
    // ---------------------------------------------------------------------

    function test_constructor_reverts_zero_owner() public {
        vm.expectRevert(CapitalFormationModule.ZeroAddress.selector);
        new CapitalFormationModule(address(0));
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

    /// @dev Drive the pair reserves so that the implied FDV in USD is
    ///      `targetFdvUsd` (pre-decimal). Keeps the agent reserve constant
    ///      and adjusts the USDC reserve so that
    ///         usdcReserve / agentReserve = priceUsdcPerAgent.
    function _setFdvUsd(uint256 targetFdvUsdMillionsScale) internal {
        // We treat `targetFdvUsdMillionsScale` as the raw USD figure (no
        // decimals). Convert into USDC wei by multiplying by 1e6.
        // Price in USDC-wei per agent-wei = (FDV_usdc) / supply.
        uint256 fdvUsdc = targetFdvUsdMillionsScale * 1e6;
        // We choose agent reserve = SUPPLY_PORTION_FOR_PAIR (constant) so
        // supply / agentReserve is fixed; then usdcReserve must equal
        //   priceUsdcPerAgent * agentReserve = (fdvUsdc / SUPPLY) * agentReserve
        // = fdvUsdc * (agentReserve / SUPPLY).
        uint256 agentReserve = 100_000_000e18; // 10% of supply in pair
        uint256 usdcReserve = fdvUsdc * agentReserve / SUPPLY;
        // Both reserves must fit in uint112: max = 2^112 ≈ 5.19e33.
        // agentReserve = 1e26 — fine. usdcReserve up to 1.6e13 — fine.
        _syncByRole(agentReserve, usdcReserve);
    }

    /// @dev Expected FDV the module will compute for a given target USD
    ///      figure. Mirrors the contract's math after lossy integer ops:
    ///      we re-derive from the same reserves the contract reads.
    function _expectedFdvUsd(uint256 targetFdvUsdMillionsScale) internal view returns (uint256) {
        // The contract reads cumulative-derived TWAP. With our `_runTwap`
        // helper holding constant reserves over the full window, TWAP =
        // spot, so the expected FDV equals the spot FDV computed exactly
        // the way the contract does:
        //   priceQ112 = (reserveQuote * Q112) / reserveAgent;
        //   fdv       = priceQ112 * supply / Q112.
        uint256 fdvUsdc = targetFdvUsdMillionsScale * 1e6;
        uint256 agentReserve = 100_000_000e18;
        uint256 usdcReserve = fdvUsdc * agentReserve / SUPPLY;
        uint256 q112 = 0x10000000000000000000000000000;
        uint256 priceQ112 = (usdcReserve * q112) / agentReserve;
        return (priceQ112 * SUPPLY) / q112;
    }

    /// @dev Sync helper that arranges reserves into the correct (token0,
    ///      token1) slots regardless of which token is which.
    function _syncByRole(uint256 agentAmt, uint256 usdcAmt) internal {
        if (pair.token0() == address(agentTok)) {
            pair.sync(uint112(agentAmt), uint112(usdcAmt));
        } else {
            pair.sync(uint112(usdcAmt), uint112(agentAmt));
        }
    }

    /// @dev Walk through a TWAP window:
    ///   1. Advance time so the cumulative grows from the LAST {sync} (the
    ///      one driven by {_setFdvUsd}) at the held price.
    ///   2. Take a {snapshot}.
    ///   3. Roll {MIN_TWAP_BLOCKS} blocks + matching seconds.
    ///   4. Re-sync the pair at the SAME reserves so the cumulative bumps
    ///      with our held TWAP price.
    function _runTwap() internal {
        // Read held reserves out of the pair.
        (uint112 r0, uint112 r1,) = pair.getReserves();

        // Advance time so the cumulative grows.
        vm.warp(block.timestamp + 64);
        // Re-sync: this rolls the cumulative forward by `priceHeld * 64`.
        pair.sync(r0, r1);

        // Snapshot.
        mod.snapshot(agent);

        // Now wait the TWAP window.
        vm.roll(block.number + 32);
        vm.warp(block.timestamp + 64);
        // Re-sync at SAME reserves so cumulative grows by `priceHeld * 64`.
        pair.sync(r0, r1);
    }
}
