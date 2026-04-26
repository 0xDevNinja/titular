// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {SixtyDaysModule} from "../../src/launchpad/modules/SixtyDaysModule.sol";

/// @dev Minimal mintable ERC-20 used as the escrow token in the unit tests.
///      No transfer tax so claim math is exact at the wei.
contract TestERC20 is ERC20 {
    constructor(string memory n, string memory s) ERC20(n, s) {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @dev Reentrant token: on every transfer, calls back into a configured target
///      with arbitrary calldata. Lets us prove that the module's CEI + guard
///      stops a re-entrant `claimRefund`.
contract ReentrantToken is ERC20 {
    address public target;
    bytes public reenterData;
    bool public armed;

    constructor() ERC20("REE", "REE") {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }

    function arm(address target_, bytes calldata data_) external {
        target = target_;
        reenterData = data_;
        armed = true;
    }

    function _update(address from, address to, uint256 value) internal override {
        super._update(from, to, value);
        if (armed && to != address(0) && target != address(0)) {
            armed = false; // single shot
            // ignore success — we only care that the inner call's revert bubbles.
            (bool ok, bytes memory ret) = target.call(reenterData);
            if (!ok) {
                // bubble revert reason for clearer test diagnostics
                assembly {
                    revert(add(ret, 32), mload(ret))
                }
            }
        }
    }
}

contract SixtyDaysModuleTest is Test {
    SixtyDaysModule internal module;
    TestERC20 internal escrowToken;

    address internal factory = address(0xFAC);
    address internal agent = address(0xA6E47);
    address internal curve = address(0xC0E);
    address internal creator = address(0xC4EA702);
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);
    address internal carol = address(0xCA501);

    uint64 internal constant START = 1_700_000_000;
    uint64 internal constant WINDOW = 60 days;

    // Phase enum values (must match the constants in the contract)
    uint8 internal constant PHASE_OPEN = 0;
    uint8 internal constant PHASE_COMMITTED = 1;
    uint8 internal constant PHASE_REFUNDING = 2;

    function setUp() public {
        module = new SixtyDaysModule(factory);
        escrowToken = new TestERC20("ESCROW", "ESC");
        vm.warp(START);
    }

    // ---------------------------------------------------------------------
    // Helpers
    // ---------------------------------------------------------------------

    function _configure() internal {
        vm.prank(factory);
        module.configure(agent, address(escrowToken), curve, creator, START);
    }

    function _accrue(uint256 amt) internal {
        escrowToken.mint(curve, amt);
        vm.prank(curve);
        escrowToken.approve(address(module), amt);
        vm.prank(curve);
        module.accrueEscrow(agent, amt);
    }

    function _contribute(address user, uint256 amt) internal {
        vm.prank(curve);
        module.recordContribution(agent, user, amt);
    }

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    function test_constructor_sets_owner_and_window() public view {
        assertEq(module.owner(), factory);
        assertEq(uint256(module.WINDOW()), uint256(WINDOW));
    }

    function test_constructor_reverts_zero_owner() public {
        vm.expectRevert(SixtyDaysModule.ZeroAddress.selector);
        new SixtyDaysModule(address(0));
    }

    // ---------------------------------------------------------------------
    // configure
    // ---------------------------------------------------------------------

    function test_configure_only_owner() public {
        vm.prank(alice);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        module.configure(agent, address(escrowToken), curve, creator, START);
    }

    function test_configure_happy_emits_and_writes() public {
        vm.expectEmit(true, true, true, true, address(module));
        emit SixtyDaysModule.Configured(agent, address(escrowToken), curve, creator, START, START + WINDOW);
        vm.prank(factory);
        module.configure(agent, address(escrowToken), curve, creator, START);

        (
            uint64 startTime,
            uint64 windowEnd,
            address tok,
            address bc,
            address cr,
            uint8 phase,
            bool configured,
            uint256 escrowAccum,
            uint256 pool,
            uint256 totalContrib
        ) = module.configs(agent);
        assertEq(uint256(startTime), uint256(START));
        assertEq(uint256(windowEnd), uint256(START + WINDOW));
        assertEq(tok, address(escrowToken));
        assertEq(bc, curve);
        assertEq(cr, creator);
        assertEq(uint256(phase), uint256(PHASE_OPEN));
        assertTrue(configured);
        assertEq(escrowAccum, 0);
        assertEq(pool, 0);
        assertEq(totalContrib, 0);
    }

    function test_configure_reverts_already_configured() public {
        _configure();
        vm.prank(factory);
        vm.expectRevert(SixtyDaysModule.AlreadyConfigured.selector);
        module.configure(agent, address(escrowToken), curve, creator, START);
    }

    function test_configure_reverts_zero_agent() public {
        vm.prank(factory);
        vm.expectRevert(SixtyDaysModule.ZeroAddress.selector);
        module.configure(address(0), address(escrowToken), curve, creator, START);
    }

    function test_configure_reverts_zero_token() public {
        vm.prank(factory);
        vm.expectRevert(SixtyDaysModule.ZeroAddress.selector);
        module.configure(agent, address(0), curve, creator, START);
    }

    function test_configure_reverts_zero_curve() public {
        vm.prank(factory);
        vm.expectRevert(SixtyDaysModule.ZeroAddress.selector);
        module.configure(agent, address(escrowToken), address(0), creator, START);
    }

    function test_configure_reverts_zero_creator() public {
        vm.prank(factory);
        vm.expectRevert(SixtyDaysModule.ZeroAddress.selector);
        module.configure(agent, address(escrowToken), curve, address(0), START);
    }

    // ---------------------------------------------------------------------
    // accrueEscrow
    // ---------------------------------------------------------------------

    function test_accrueEscrow_only_curve() public {
        _configure();
        escrowToken.mint(alice, 100e18);
        vm.prank(alice);
        escrowToken.approve(address(module), 100e18);
        vm.prank(alice);
        vm.expectRevert(SixtyDaysModule.NotBondingCurve.selector);
        module.accrueEscrow(agent, 100e18);
    }

    function test_accrueEscrow_reverts_zero_amount() public {
        _configure();
        vm.prank(curve);
        vm.expectRevert(SixtyDaysModule.ZeroAmount.selector);
        module.accrueEscrow(agent, 0);
    }

    function test_accrueEscrow_pulls_and_accumulates() public {
        _configure();
        escrowToken.mint(curve, 100e18);
        vm.prank(curve);
        escrowToken.approve(address(module), 100e18);

        vm.expectEmit(true, false, false, true, address(module));
        emit SixtyDaysModule.EscrowAccrued(agent, 30e18, 30e18);
        vm.prank(curve);
        module.accrueEscrow(agent, 30e18);

        vm.expectEmit(true, false, false, true, address(module));
        emit SixtyDaysModule.EscrowAccrued(agent, 70e18, 100e18);
        vm.prank(curve);
        module.accrueEscrow(agent, 70e18);

        (,,,,,,, uint256 escrowAccum,,) = module.configs(agent);
        assertEq(escrowAccum, 100e18);
        assertEq(escrowToken.balanceOf(address(module)), 100e18);
    }

    function test_accrueEscrow_reverts_after_commit() public {
        _configure();
        _accrue(100e18);
        vm.warp(START + WINDOW);
        vm.prank(creator);
        module.commit(agent);

        // further accrual blocked
        escrowToken.mint(curve, 1e18);
        vm.prank(curve);
        escrowToken.approve(address(module), 1e18);
        vm.prank(curve);
        vm.expectRevert(SixtyDaysModule.NotOpen.selector);
        module.accrueEscrow(agent, 1e18);
    }

    function test_accrueEscrow_reverts_after_refund() public {
        _configure();
        _accrue(100e18);
        _contribute(alice, 1e18);
        vm.prank(creator);
        module.refund(agent);

        escrowToken.mint(curve, 1e18);
        vm.prank(curve);
        escrowToken.approve(address(module), 1e18);
        vm.prank(curve);
        vm.expectRevert(SixtyDaysModule.NotOpen.selector);
        module.accrueEscrow(agent, 1e18);
    }

    // ---------------------------------------------------------------------
    // recordContribution
    // ---------------------------------------------------------------------

    function test_recordContribution_only_curve() public {
        _configure();
        vm.prank(alice);
        vm.expectRevert(SixtyDaysModule.NotBondingCurve.selector);
        module.recordContribution(agent, alice, 1e18);
    }

    function test_recordContribution_reverts_zero_amount() public {
        _configure();
        vm.prank(curve);
        vm.expectRevert(SixtyDaysModule.ZeroAmount.selector);
        module.recordContribution(agent, alice, 0);
    }

    function test_recordContribution_reverts_zero_user() public {
        _configure();
        vm.prank(curve);
        vm.expectRevert(SixtyDaysModule.ZeroAddress.selector);
        module.recordContribution(agent, address(0), 1e18);
    }

    function test_recordContribution_accumulates() public {
        _configure();
        _contribute(alice, 30e18);
        _contribute(alice, 20e18);
        _contribute(bob, 50e18);

        assertEq(module.contributed(agent, alice), 50e18);
        assertEq(module.contributed(agent, bob), 50e18);
        (,,,,,,,,, uint256 totalContrib) = module.configs(agent);
        assertEq(totalContrib, 100e18);
    }

    function test_recordContribution_reverts_after_refund() public {
        _configure();
        _contribute(alice, 1e18);
        vm.prank(creator);
        module.refund(agent);
        vm.prank(curve);
        vm.expectRevert(SixtyDaysModule.NotOpen.selector);
        module.recordContribution(agent, bob, 1e18);
    }

    // ---------------------------------------------------------------------
    // commit
    // ---------------------------------------------------------------------

    function test_commit_only_creator() public {
        _configure();
        _accrue(100e18);
        vm.warp(START + WINDOW);
        vm.prank(alice);
        vm.expectRevert(SixtyDaysModule.NotCreator.selector);
        module.commit(agent);
    }

    function test_commit_reverts_too_early() public {
        _configure();
        _accrue(100e18);
        vm.warp(START + WINDOW - 1);
        vm.prank(creator);
        vm.expectRevert(SixtyDaysModule.WindowNotElapsed.selector);
        module.commit(agent);
    }

    function test_commit_happy_at_window_boundary() public {
        _configure();
        _accrue(100e18);
        vm.warp(START + WINDOW);

        vm.expectEmit(true, true, false, true, address(module));
        emit SixtyDaysModule.Committed(agent, creator, 100e18);
        vm.prank(creator);
        module.commit(agent);

        assertEq(escrowToken.balanceOf(creator), 100e18);
        assertEq(escrowToken.balanceOf(address(module)), 0);
        assertEq(uint256(module.phaseOf(agent)), uint256(PHASE_COMMITTED));
    }

    function test_commit_zero_escrow_succeeds() public {
        _configure();
        // No accrual — escrow is zero. Commit still flips the phase.
        vm.warp(START + WINDOW);
        vm.expectEmit(true, true, false, true, address(module));
        emit SixtyDaysModule.Committed(agent, creator, 0);
        vm.prank(creator);
        module.commit(agent);
        assertEq(escrowToken.balanceOf(creator), 0);
        assertEq(uint256(module.phaseOf(agent)), uint256(PHASE_COMMITTED));
    }

    function test_commit_reverts_double_commit() public {
        _configure();
        _accrue(100e18);
        vm.warp(START + WINDOW);
        vm.prank(creator);
        module.commit(agent);
        vm.prank(creator);
        vm.expectRevert(SixtyDaysModule.NotOpen.selector);
        module.commit(agent);
    }

    function test_commit_reverts_after_refund() public {
        _configure();
        _accrue(100e18);
        _contribute(alice, 1e18);
        vm.prank(creator);
        module.refund(agent);
        vm.warp(START + WINDOW);
        vm.prank(creator);
        vm.expectRevert(SixtyDaysModule.NotOpen.selector);
        module.commit(agent);
    }

    // ---------------------------------------------------------------------
    // refund
    // ---------------------------------------------------------------------

    function test_refund_only_creator() public {
        _configure();
        _accrue(100e18);
        vm.prank(alice);
        vm.expectRevert(SixtyDaysModule.NotCreator.selector);
        module.refund(agent);
    }

    function test_refund_happy_before_window_seeds_pool() public {
        _configure();
        _accrue(100e18);
        _contribute(alice, 30e18);
        _contribute(bob, 70e18);

        vm.expectEmit(true, true, false, true, address(module));
        emit SixtyDaysModule.Refunded(agent, creator, 100e18);
        vm.prank(creator);
        module.refund(agent);

        assertEq(uint256(module.phaseOf(agent)), uint256(PHASE_REFUNDING));
        (,,,,,,,, uint256 pool,) = module.configs(agent);
        assertEq(pool, 100e18);
        assertEq(escrowToken.balanceOf(address(module)), 100e18);
    }

    function test_refund_happy_after_window_still_works() public {
        _configure();
        _accrue(100e18);
        _contribute(alice, 100e18);
        vm.warp(START + WINDOW + 100);
        vm.prank(creator);
        module.refund(agent);
        assertEq(uint256(module.phaseOf(agent)), uint256(PHASE_REFUNDING));
    }

    function test_refund_reverts_double() public {
        _configure();
        _accrue(100e18);
        _contribute(alice, 1e18);
        vm.prank(creator);
        module.refund(agent);
        vm.prank(creator);
        vm.expectRevert(SixtyDaysModule.NotOpen.selector);
        module.refund(agent);
    }

    function test_refund_reverts_after_commit() public {
        _configure();
        _accrue(100e18);
        vm.warp(START + WINDOW);
        vm.prank(creator);
        module.commit(agent);
        vm.prank(creator);
        vm.expectRevert(SixtyDaysModule.NotOpen.selector);
        module.refund(agent);
    }

    // ---------------------------------------------------------------------
    // accrueRefundPool
    // ---------------------------------------------------------------------

    function test_accrueRefundPool_only_curve() public {
        _configure();
        _accrue(100e18);
        _contribute(alice, 1e18);
        vm.prank(creator);
        module.refund(agent);

        escrowToken.mint(alice, 50e18);
        vm.prank(alice);
        escrowToken.approve(address(module), 50e18);
        vm.prank(alice);
        vm.expectRevert(SixtyDaysModule.NotBondingCurve.selector);
        module.accrueRefundPool(agent, 50e18);
    }

    function test_accrueRefundPool_reverts_when_not_refunding() public {
        _configure();
        escrowToken.mint(curve, 50e18);
        vm.prank(curve);
        escrowToken.approve(address(module), 50e18);
        vm.prank(curve);
        vm.expectRevert(SixtyDaysModule.NotRefunding.selector);
        module.accrueRefundPool(agent, 50e18);
    }

    function test_accrueRefundPool_reverts_zero_amount() public {
        _configure();
        _accrue(10e18);
        _contribute(alice, 1e18);
        vm.prank(creator);
        module.refund(agent);
        vm.prank(curve);
        vm.expectRevert(SixtyDaysModule.ZeroAmount.selector);
        module.accrueRefundPool(agent, 0);
    }

    function test_accrueRefundPool_pulls_and_emits() public {
        _configure();
        _accrue(100e18); // escrow = 100
        _contribute(alice, 1e18);
        vm.prank(creator);
        module.refund(agent); // pool seeded at 100

        escrowToken.mint(curve, 50e18);
        vm.prank(curve);
        escrowToken.approve(address(module), 50e18);

        vm.expectEmit(true, false, false, true, address(module));
        emit SixtyDaysModule.RefundPoolToppedUp(agent, 50e18, 150e18);
        vm.prank(curve);
        module.accrueRefundPool(agent, 50e18);

        (,,,,,,,, uint256 pool,) = module.configs(agent);
        assertEq(pool, 150e18);
        assertEq(escrowToken.balanceOf(address(module)), 150e18);
    }

    // ---------------------------------------------------------------------
    // claimRefund
    // ---------------------------------------------------------------------

    function _setupRefundScenario(uint256 a, uint256 b, uint256 escrow) internal {
        _configure();
        _accrue(escrow);
        _contribute(alice, a);
        _contribute(bob, b);
        vm.prank(creator);
        module.refund(agent);
    }

    function test_claimRefund_pro_rata_split() public {
        // Alice 30, Bob 70 → escrow 100 → Alice 30, Bob 70.
        _setupRefundScenario(30e18, 70e18, 100e18);

        vm.expectEmit(true, true, false, true, address(module));
        emit SixtyDaysModule.RefundClaimed(agent, alice, 30e18);
        vm.prank(alice);
        uint256 a = module.claimRefund(agent);
        assertEq(a, 30e18);
        assertEq(escrowToken.balanceOf(alice), 30e18);

        vm.prank(bob);
        uint256 b = module.claimRefund(agent);
        assertEq(b, 70e18);
        assertEq(escrowToken.balanceOf(bob), 70e18);

        // Pool fully drained (no rounding dust at these magnitudes).
        assertEq(escrowToken.balanceOf(address(module)), 0);
    }

    function test_claimRefund_includes_topup() public {
        _setupRefundScenario(50e18, 50e18, 100e18);

        // Curve drains 100 of residual reserves into the pool. Total now 200.
        escrowToken.mint(curve, 100e18);
        vm.prank(curve);
        escrowToken.approve(address(module), 100e18);
        vm.prank(curve);
        module.accrueRefundPool(agent, 100e18);

        vm.prank(alice);
        uint256 a = module.claimRefund(agent);
        vm.prank(bob);
        uint256 b = module.claimRefund(agent);
        assertEq(a, 100e18);
        assertEq(b, 100e18);
        assertEq(escrowToken.balanceOf(address(module)), 0);
    }

    function test_claimRefund_reverts_when_not_in_refund_phase() public {
        _configure();
        _contribute(alice, 100e18);
        vm.prank(alice);
        vm.expectRevert(SixtyDaysModule.NotInRefundPhase.selector);
        module.claimRefund(agent);
    }

    function test_claimRefund_reverts_already_claimed() public {
        _setupRefundScenario(50e18, 50e18, 80e18);
        vm.prank(alice);
        module.claimRefund(agent);
        vm.prank(alice);
        vm.expectRevert(SixtyDaysModule.AlreadyClaimed.selector);
        module.claimRefund(agent);
    }

    function test_claimRefund_reverts_no_contribution() public {
        _setupRefundScenario(50e18, 50e18, 80e18);
        vm.prank(carol);
        vm.expectRevert(SixtyDaysModule.NoContribution.selector);
        module.claimRefund(agent);
    }

    function test_claimRefund_reverts_zero_share_floors() public {
        // Single buyer of 1 wei against a refund pool of 0 -> share = 0 -> revert.
        _configure();
        // No accrue — escrow stays 0. Contributor present.
        _contribute(alice, 1e18);
        vm.prank(creator);
        module.refund(agent); // pool = 0

        vm.prank(alice);
        vm.expectRevert(SixtyDaysModule.NoContribution.selector);
        module.claimRefund(agent);
    }

    // ---------------------------------------------------------------------
    // previewRefund
    // ---------------------------------------------------------------------

    function test_previewRefund_matches_claim() public {
        _setupRefundScenario(40e18, 60e18, 100e18);
        uint256 pa = module.previewRefund(agent, alice);
        uint256 pb = module.previewRefund(agent, bob);
        uint256 pc = module.previewRefund(agent, carol);
        assertEq(pa, 40e18);
        assertEq(pb, 60e18);
        assertEq(pc, 0);

        vm.prank(alice);
        uint256 actual = module.claimRefund(agent);
        assertEq(actual, pa);
        assertEq(module.previewRefund(agent, alice), 0);
    }

    function test_previewRefund_zero_when_not_refunding() public {
        _configure();
        _contribute(alice, 100e18);
        assertEq(module.previewRefund(agent, alice), 0);
    }

    // ---------------------------------------------------------------------
    // Reentrancy probe
    // ---------------------------------------------------------------------

    function test_reentrancy_blocked_on_claimRefund() public {
        // Use a reentrant token as the escrow asset; arm it to call
        // claimRefund again from inside the outbound transfer.
        ReentrantToken r = new ReentrantToken();

        SixtyDaysModule m = new SixtyDaysModule(factory);

        vm.prank(factory);
        m.configure(agent, address(r), curve, creator, START);

        // Curve accrues 100 of the reentrant token.
        r.mint(curve, 100e18);
        vm.prank(curve);
        r.approve(address(m), 100e18);
        vm.prank(curve);
        m.accrueEscrow(agent, 100e18);

        // Alice and Bob each contribute 50 quote.
        vm.prank(curve);
        m.recordContribution(agent, alice, 50e18);
        vm.prank(curve);
        m.recordContribution(agent, bob, 50e18);

        // Creator opens refund.
        vm.prank(creator);
        m.refund(agent);

        // Arm the token so the next transfer to alice triggers a re-entrant
        // claimRefund(agent) from alice. The inner call must revert via the
        // ReentrancyGuard, and that revert must bubble up.
        r.arm(address(m), abi.encodeCall(m.claimRefund, (agent)));

        vm.prank(alice);
        // OZ ReentrancyGuard reverts with `ReentrancyGuardReentrantCall()`.
        vm.expectRevert();
        m.claimRefund(agent);

        // State unchanged: alice did not actually claim because the outer call
        // reverted atomically.
        assertEq(r.balanceOf(alice), 0);
        assertFalse(m.claimed(agent, alice));
    }

    // ---------------------------------------------------------------------
    // Fuzz: pull-pattern solvency
    // ---------------------------------------------------------------------

    /// @dev Five buyers with fuzzed contributions, escrow pool fuzzed too. The
    ///      sum of all paid-out claims must never exceed the seeded refund pool
    ///      (rounding floors are in the contract's favour).
    function testFuzz_sum_claims_le_pool(uint96[5] memory contribs, uint96 escrow) public {
        address[5] memory users = [address(0xA1), address(0xA2), address(0xA3), address(0xA4), address(0xA5)];
        _configure();

        uint256 esc = bound(uint256(escrow), 0, 1_000_000e18);
        if (esc != 0) {
            escrowToken.mint(curve, esc);
            vm.prank(curve);
            escrowToken.approve(address(module), esc);
            vm.prank(curve);
            module.accrueEscrow(agent, esc);
        }

        for (uint256 i = 0; i < 5; i++) {
            uint256 a = bound(uint256(contribs[i]), 1, 1_000_000e18);
            vm.prank(curve);
            module.recordContribution(agent, users[i], a);
        }

        vm.prank(creator);
        module.refund(agent);

        uint256 paid;
        for (uint256 i = 0; i < 5; i++) {
            uint256 preview = module.previewRefund(agent, users[i]);
            if (preview == 0) continue;
            vm.prank(users[i]);
            paid += module.claimRefund(agent);
        }
        assertLe(paid, esc, "sum of claims exceeds pool");
        // Module retains only the floor-dust; never goes negative.
        assertEq(escrowToken.balanceOf(address(module)), esc - paid);
    }
}
