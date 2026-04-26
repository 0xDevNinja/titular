// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {Clones} from "@openzeppelin/contracts/proxy/Clones.sol";
import {AgentToken} from "../../src/launchpad/AgentToken.sol";

contract AgentTokenTest is Test {
    AgentToken internal impl;
    AgentToken internal token;

    address internal creator = address(0xC0FFEE);
    address internal feeRouter = address(0xFEE);
    address internal bondingCurve = address(0xB0A2D);
    address internal graduator = address(0x6BADD);
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);

    string internal constant NAME = "Agent Smith";
    string internal constant SYMBOL = "SMITH";

    event AgentTokenInitialized(
        address indexed creator, address indexed feeRouter, address indexed bondingCurve, string name, string symbol
    );
    event TaxExemptSet(address indexed addr, bool exempt);
    event TaxCollected(address indexed from, address indexed to, uint256 taxAmount);
    event Transfer(address indexed from, address indexed to, uint256 value);

    function setUp() public {
        impl = new AgentToken();
        token = _deployClone(NAME, SYMBOL, creator, feeRouter, bondingCurve, graduator);
    }

    function _deployClone(
        string memory name_,
        string memory symbol_,
        address creator_,
        address feeRouter_,
        address bondingCurve_,
        address graduator_
    ) internal returns (AgentToken cloned) {
        cloned = AgentToken(Clones.clone(address(impl)));
        cloned.initialize(name_, symbol_, creator_, feeRouter_, bondingCurve_, graduator_);
    }

    // ---------------------------------------------------------------------
    // initialize
    // ---------------------------------------------------------------------

    function test_initialize_sets_state_and_mints_total_supply() public view {
        assertEq(token.name(), NAME);
        assertEq(token.symbol(), SYMBOL);
        assertEq(token.decimals(), 18);
        assertEq(token.creator(), creator);
        assertEq(token.feeRouter(), feeRouter);
        assertEq(token.bondingCurve(), bondingCurve);
        assertEq(token.owner(), creator);
        assertEq(token.totalSupply(), token.TOTAL_SUPPLY());
        assertEq(token.balanceOf(bondingCurve), token.TOTAL_SUPPLY());
        assertEq(token.TRADE_TAX_BPS(), 100);
        assertEq(token.BPS_DENOMINATOR(), 10_000);
    }

    function test_initialize_populates_tax_exempt_allowlist() public view {
        // The three protocol-internal contracts must land on the allowlist; nothing
        // else may. End-user accounts (creator, alice, bob) must NOT be exempt.
        assertTrue(token.taxExempt(feeRouter), "feeRouter exempt");
        assertTrue(token.taxExempt(bondingCurve), "bondingCurve exempt");
        assertTrue(token.taxExempt(graduator), "graduator exempt");
        assertFalse(token.taxExempt(creator), "creator not exempt");
        assertFalse(token.taxExempt(alice), "alice not exempt");
        assertFalse(token.taxExempt(address(0)), "zero address not exempt");
    }

    function test_initialize_emits_event_and_mint_transfer() public {
        AgentToken fresh = AgentToken(Clones.clone(address(impl)));

        vm.expectEmit(true, true, false, true, address(fresh));
        emit Transfer(address(0), bondingCurve, fresh.TOTAL_SUPPLY());
        vm.expectEmit(true, false, false, true, address(fresh));
        emit TaxExemptSet(feeRouter, true);
        vm.expectEmit(true, false, false, true, address(fresh));
        emit TaxExemptSet(bondingCurve, true);
        vm.expectEmit(true, false, false, true, address(fresh));
        emit TaxExemptSet(graduator, true);
        vm.expectEmit(true, true, true, true, address(fresh));
        emit AgentTokenInitialized(creator, feeRouter, bondingCurve, NAME, SYMBOL);

        fresh.initialize(NAME, SYMBOL, creator, feeRouter, bondingCurve, graduator);
    }

    function test_initialize_revert_when_called_twice() public {
        vm.expectRevert();
        token.initialize(NAME, SYMBOL, creator, feeRouter, bondingCurve, graduator);
    }

    function test_initialize_revert_on_implementation_directly() public {
        vm.expectRevert();
        impl.initialize(NAME, SYMBOL, creator, feeRouter, bondingCurve, graduator);
    }

    function test_initialize_revert_zero_creator() public {
        AgentToken fresh = AgentToken(Clones.clone(address(impl)));
        vm.expectRevert(AgentToken.ZeroAddress.selector);
        fresh.initialize(NAME, SYMBOL, address(0), feeRouter, bondingCurve, graduator);
    }

    function test_initialize_revert_zero_feeRouter() public {
        AgentToken fresh = AgentToken(Clones.clone(address(impl)));
        vm.expectRevert(AgentToken.ZeroAddress.selector);
        fresh.initialize(NAME, SYMBOL, creator, address(0), bondingCurve, graduator);
    }

    function test_initialize_revert_zero_bondingCurve() public {
        AgentToken fresh = AgentToken(Clones.clone(address(impl)));
        vm.expectRevert(AgentToken.ZeroAddress.selector);
        fresh.initialize(NAME, SYMBOL, creator, feeRouter, address(0), graduator);
    }

    function test_initialize_revert_zero_graduator() public {
        AgentToken fresh = AgentToken(Clones.clone(address(impl)));
        vm.expectRevert(AgentToken.ZeroAddress.selector);
        fresh.initialize(NAME, SYMBOL, creator, feeRouter, bondingCurve, address(0));
    }

    function test_initialize_revert_empty_name() public {
        AgentToken fresh = AgentToken(Clones.clone(address(impl)));
        vm.expectRevert(AgentToken.EmptyMetadata.selector);
        fresh.initialize("", SYMBOL, creator, feeRouter, bondingCurve, graduator);
    }

    function test_initialize_revert_empty_symbol() public {
        AgentToken fresh = AgentToken(Clones.clone(address(impl)));
        vm.expectRevert(AgentToken.EmptyMetadata.selector);
        fresh.initialize(NAME, "", creator, feeRouter, bondingCurve, graduator);
    }

    // ---------------------------------------------------------------------
    // transfer tax
    // ---------------------------------------------------------------------

    function test_transfer_routes_one_percent_to_feeRouter() public {
        // Seed alice via a regular user, NOT via the bondingCurve (which is
        // exempt from tax). Use a curve -> alice -> bob hop so the seed transfer
        // is taxed and the subsequent alice -> bob transfer is taxed too.
        // Curve is exempt, so the first hop is untaxed; we instead route the
        // seed via a non-exempt mid-actor.
        uint256 seed = 10_000e18;
        // Curve is exempt: curve -> mid is untaxed and lands the full seed on `mid`.
        address mid = address(0x111D);
        vm.prank(bondingCurve);
        token.transfer(mid, seed);

        // Now mid -> alice is taxed (neither party is exempt).
        uint256 expectedTax = (seed * 100) / 10_000; // 1% of seed
        // Send the full seed mid -> alice and assert the router skim.
        vm.prank(mid);
        token.transfer(alice, seed);
        assertEq(token.balanceOf(alice), seed - expectedTax);
        assertEq(token.balanceOf(feeRouter), expectedTax);

        // Now alice -> bob should also be taxed.
        uint256 amount = 1000e18;
        uint256 tax = (amount * 100) / 10_000;
        uint256 routerBefore = token.balanceOf(feeRouter);

        vm.expectEmit(true, true, false, true, address(token));
        emit TaxCollected(alice, bob, tax);
        vm.prank(alice);
        token.transfer(bob, amount);

        assertEq(token.balanceOf(bob), amount - tax);
        assertEq(token.balanceOf(feeRouter), routerBefore + tax);
    }

    function test_transfer_preserves_supply() public {
        uint256 supplyBefore = token.totalSupply();
        // Curve -> alice is exempt (curve allowlisted), then alice -> bob is taxed.
        vm.prank(bondingCurve);
        token.transfer(alice, 500e18);

        vm.prank(alice);
        token.transfer(bob, 100e18);

        assertEq(token.totalSupply(), supplyBefore);
        assertEq(
            token.balanceOf(bondingCurve) + token.balanceOf(alice) + token.balanceOf(bob) + token.balanceOf(feeRouter),
            supplyBefore
        );
    }

    function test_transferFrom_is_taxed() public {
        // Pre-seed alice through the curve (exempt -> alice is untaxed; alice gets the
        // full 5000e18). Then alice -> bob is taxed on the transferFrom path.
        uint256 seed = 5000e18;
        vm.prank(bondingCurve);
        token.transfer(alice, seed);

        uint256 amount = 1000e18;
        uint256 tax = (amount * 100) / 10_000;

        vm.prank(alice);
        token.approve(bob, amount);

        uint256 routerBefore = token.balanceOf(feeRouter);
        vm.prank(bob);
        token.transferFrom(alice, bob, amount);

        assertEq(token.balanceOf(bob), amount - tax);
        assertEq(token.balanceOf(feeRouter), routerBefore + tax);
        assertEq(token.allowance(alice, bob), 0);
    }

    // ---------------------------------------------------------------------
    // tax-skip paths
    // ---------------------------------------------------------------------

    function test_mint_is_not_taxed() public view {
        // The constructor-time mint in initialize must not have been taxed.
        assertEq(token.balanceOf(feeRouter), 0);
        assertEq(token.balanceOf(bondingCurve), token.TOTAL_SUPPLY());
    }

    function test_transfer_from_bondingCurve_is_not_taxed() public {
        // BondingCurve is allowlisted: curve -> user delivers the FULL amount,
        // no skim to the FeeRouter.
        uint256 amount = 1000e18;
        vm.prank(bondingCurve);
        token.transfer(alice, amount);

        assertEq(token.balanceOf(alice), amount, "curve transfer untaxed");
        assertEq(token.balanceOf(feeRouter), 0, "no skim from curve hop");
    }

    function test_transfer_from_graduator_is_not_taxed() public {
        // Graduator is allowlisted: graduator -> any delivers full amount. Critical
        // for the graduation hand-off (graduator -> Uniswap router -> pair).
        uint256 amount = 500_000e18;
        vm.prank(bondingCurve);
        token.transfer(graduator, amount);
        // bondingCurve is also exempt, so graduator received the full `amount`.
        assertEq(token.balanceOf(graduator), amount);

        vm.prank(graduator);
        token.transfer(alice, amount);
        assertEq(token.balanceOf(alice), amount, "graduator hop untaxed");
        assertEq(token.balanceOf(feeRouter), 0, "no skim from graduator hop");
    }

    function test_transfer_from_feeRouter_is_not_taxed() public {
        // Seed the router via a non-exempt mid-actor so the router accrues some balance.
        uint256 seed = 1000e18;
        // bondingCurve -> mid is untaxed (curve exempt); mid -> alice gets taxed.
        address mid = address(0x222D);
        vm.prank(bondingCurve);
        token.transfer(mid, seed);
        vm.prank(mid);
        token.transfer(alice, seed);

        uint256 routerBalance = token.balanceOf(feeRouter);
        assertGt(routerBalance, 0, "router accrued seed tax");

        // Router -> alice must deliver full amount, no tax loop.
        uint256 aliceBefore = token.balanceOf(alice);
        vm.prank(feeRouter);
        token.transfer(alice, routerBalance);

        assertEq(token.balanceOf(alice), aliceBefore + routerBalance);
        assertEq(token.balanceOf(feeRouter), 0);
    }

    function test_transfer_to_feeRouter_is_not_taxed() public {
        // Seed alice via a non-exempt mid-actor so alice's balance is what it gets,
        // then alice -> feeRouter must NOT levy tax (recipient is exempt).
        uint256 seed = 2000e18;
        address mid = address(0x333D);
        vm.prank(bondingCurve);
        token.transfer(mid, seed);
        vm.prank(mid);
        token.transfer(alice, seed);

        uint256 aliceBalance = token.balanceOf(alice);
        uint256 routerBefore = token.balanceOf(feeRouter);

        // Direct transfer to router — full amount lands, no double-tax.
        vm.prank(alice);
        token.transfer(feeRouter, aliceBalance);

        assertEq(token.balanceOf(alice), 0);
        assertEq(token.balanceOf(feeRouter), routerBefore + aliceBalance);
    }

    function test_transfer_to_graduator_is_not_taxed() public {
        // Curve -> graduator (the actual graduation pull path) must be untaxed.
        // This is the precise on-chain hop that issue #265 was breaking.
        uint256 amount = 750_000e18;
        uint256 routerBefore = token.balanceOf(feeRouter);

        vm.prank(bondingCurve);
        token.transfer(graduator, amount);

        assertEq(token.balanceOf(graduator), amount, "graduator received full amount");
        assertEq(token.balanceOf(feeRouter), routerBefore, "no tax skim on graduation pull");
    }

    function test_zero_value_transfer_emits_no_tax() public {
        // bondingCurve -> alice is exempt (curve allowlisted), so alice gets full 1000e18.
        // Then alice -> bob with value 0 still flows through `_update` and must not skim.
        vm.prank(bondingCurve);
        token.transfer(alice, 1000e18);

        uint256 routerBefore = token.balanceOf(feeRouter);
        vm.prank(alice);
        token.transfer(bob, 0);
        assertEq(token.balanceOf(feeRouter), routerBefore);
    }

    function test_dust_value_rounds_tax_to_zero() public {
        // Seed alice (curve exempt -> full transfer).
        vm.prank(bondingCurve);
        token.transfer(alice, 1000e18);

        // `value = 99` -> tax = 99 * 100 / 10_000 = 0. No tax collected, no dust loss.
        uint256 routerBefore = token.balanceOf(feeRouter);
        uint256 bobBefore = token.balanceOf(bob);

        vm.prank(alice);
        token.transfer(bob, 99);

        assertEq(token.balanceOf(feeRouter), routerBefore);
        assertEq(token.balanceOf(bob), bobBefore + 99);
    }

    // ---------------------------------------------------------------------
    // supply conservation
    // ---------------------------------------------------------------------

    function test_totalSupply_immutable_across_transfers() public {
        uint256 expected = token.TOTAL_SUPPLY();
        assertEq(token.totalSupply(), expected);

        vm.prank(bondingCurve);
        token.transfer(alice, 10_000e18);
        assertEq(token.totalSupply(), expected);

        vm.prank(alice);
        token.transfer(bob, 500e18);
        assertEq(token.totalSupply(), expected);

        uint256 routerBalance = token.balanceOf(feeRouter);
        vm.prank(feeRouter);
        token.transfer(alice, routerBalance);
        assertEq(token.totalSupply(), expected);
    }

    function testFuzz_tax_math_never_exceeds_value(uint256 value) public {
        value = bound(value, 0, token.TOTAL_SUPPLY());
        // Curve is exempt, so curve -> alice ships the full `value` to alice with no skim.
        // Then route alice -> bob (both non-exempt) and assert the tax math.
        vm.prank(bondingCurve);
        token.transfer(alice, value);
        assertEq(token.balanceOf(alice), value, "curve hop is untaxed");

        vm.prank(alice);
        token.transfer(bob, value);

        uint256 tax = (value * 100) / 10_000;
        assertEq(token.balanceOf(bob), value - tax);
        assertEq(token.balanceOf(feeRouter), tax);
    }

    function testFuzz_router_sweeps_without_tax(uint256 seed) public {
        seed = bound(seed, 1e18, token.TOTAL_SUPPLY());
        // bondingCurve -> alice is now untaxed (curve exempt). Push alice -> bob through
        // a non-exempt mid actor so the FeeRouter accrues the tax for the sweep test.
        vm.prank(bondingCurve);
        token.transfer(alice, seed);
        vm.prank(alice);
        token.transfer(bob, seed);

        uint256 routerBalance = token.balanceOf(feeRouter);
        vm.prank(feeRouter);
        token.transfer(address(0xDEAD), routerBalance);

        assertEq(token.balanceOf(feeRouter), 0);
        assertEq(token.balanceOf(address(0xDEAD)), routerBalance); // no tax skimmed
    }
}
