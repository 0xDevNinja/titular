// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {FeeSplitter} from "../../src/acp/FeeSplitter.sol";

contract MockERC20 is ERC20 {
    constructor() ERC20("Mock", "MCK") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

contract FeeSplitterTest is Test {
    FeeSplitter internal splitter;
    MockERC20 internal token;

    address internal admin = makeAddr("admin");
    address internal treasury = makeAddr("treasury");
    address internal buyback = makeAddr("buyback");
    address internal caller = makeAddr("caller");
    address internal primary = makeAddr("primary");
    address internal stranger = makeAddr("stranger");

    uint256 internal constant AMOUNT = 10_000e18; // 10 000 tokens for easy BPS arithmetic

    event FeeSplit(
        address indexed token,
        address indexed primary,
        FeeSplitter.Schedule indexed schedule,
        uint256 primaryAmount,
        uint256 treasuryAmount,
        uint256 buybackAmount
    );

    event TreasuryUpdated(address indexed oldTreasury, address indexed newTreasury);
    event BuybackBurnerUpdated(address indexed oldBurner, address indexed newBurner);

    function setUp() public {
        splitter = new FeeSplitter(admin, treasury, buyback);
        token = new MockERC20();

        // Grant caller role
        bytes32 callerRole = splitter.CALLER_ROLE();
        vm.prank(admin);
        splitter.grantRole(callerRole, caller);

        // Mint + approve
        token.mint(caller, AMOUNT * 100);
        vm.prank(caller);
        token.approve(address(splitter), type(uint256).max);
    }

    // ─────────────────────────────────────────────────────────────
    // Schedule A: 95/3/2
    // ─────────────────────────────────────────────────────────────

    function test_split_scheduleA_amounts() public {
        // 10 000 tokens → primary=9500, treasury=300, buyback=200
        uint256 expectedPrimary = 9500e18;
        uint256 expectedTreasury = 300e18;
        uint256 expectedBuyback = 200e18;

        vm.expectEmit(true, true, true, true);
        emit FeeSplit(address(token), primary, FeeSplitter.Schedule.A, expectedPrimary, expectedTreasury, expectedBuyback);

        vm.prank(caller);
        splitter.split(address(token), AMOUNT, primary, FeeSplitter.Schedule.A);

        assertEq(token.balanceOf(primary), expectedPrimary);
        assertEq(token.balanceOf(treasury), expectedTreasury);
        assertEq(token.balanceOf(buyback), expectedBuyback);
        assertEq(token.balanceOf(address(splitter)), 0); // no residual
    }

    // ─────────────────────────────────────────────────────────────
    // Schedule B: 90/5/5
    // ─────────────────────────────────────────────────────────────

    function test_split_scheduleB_amounts() public {
        // 10 000 tokens → primary=9000, treasury=500, buyback=500
        uint256 expectedPrimary = 9000e18;
        uint256 expectedTreasury = 500e18;
        uint256 expectedBuyback = 500e18;

        vm.prank(caller);
        splitter.split(address(token), AMOUNT, primary, FeeSplitter.Schedule.B);

        assertEq(token.balanceOf(primary), expectedPrimary);
        assertEq(token.balanceOf(treasury), expectedTreasury);
        assertEq(token.balanceOf(buyback), expectedBuyback);
    }

    // ─────────────────────────────────────────────────────────────
    // Dust handling
    // ─────────────────────────────────────────────────────────────

    function test_split_dust_goesToPrimary() public {
        // 1 token — 3% treasury = 0, 2% buyback = 0, primary gets all 1
        uint256 tiny = 1;
        token.mint(caller, tiny);
        vm.prank(caller);
        token.approve(address(splitter), type(uint256).max);

        vm.prank(caller);
        splitter.split(address(token), tiny, primary, FeeSplitter.Schedule.A);

        assertEq(token.balanceOf(primary), tiny);
        assertEq(token.balanceOf(treasury), 0);
        assertEq(token.balanceOf(buyback), 0);
    }

    // ─────────────────────────────────────────────────────────────
    // previewSplit
    // ─────────────────────────────────────────────────────────────

    function test_previewSplit_scheduleA() public view {
        (uint256 p, uint256 t, uint256 b) = splitter.previewSplit(AMOUNT, FeeSplitter.Schedule.A);
        assertEq(p, 9500e18);
        assertEq(t, 300e18);
        assertEq(b, 200e18);
        assertEq(p + t + b, AMOUNT);
    }

    function test_previewSplit_scheduleB() public view {
        (uint256 p, uint256 t, uint256 b) = splitter.previewSplit(AMOUNT, FeeSplitter.Schedule.B);
        assertEq(p, 9000e18);
        assertEq(t, 500e18);
        assertEq(b, 500e18);
        assertEq(p + t + b, AMOUNT);
    }

    // ─────────────────────────────────────────────────────────────
    // Revert paths
    // ─────────────────────────────────────────────────────────────

    function test_split_revert_notCaller() public {
        vm.expectRevert();
        vm.prank(stranger);
        splitter.split(address(token), AMOUNT, primary, FeeSplitter.Schedule.A);
    }

    function test_split_revert_zeroToken() public {
        vm.expectRevert(FeeSplitter.ZeroAddress.selector);
        vm.prank(caller);
        splitter.split(address(0), AMOUNT, primary, FeeSplitter.Schedule.A);
    }

    function test_split_revert_zeroPrimary() public {
        vm.expectRevert(FeeSplitter.ZeroAddress.selector);
        vm.prank(caller);
        splitter.split(address(token), AMOUNT, address(0), FeeSplitter.Schedule.A);
    }

    function test_split_revert_zeroAmount() public {
        vm.expectRevert(FeeSplitter.ZeroAmount.selector);
        vm.prank(caller);
        splitter.split(address(token), 0, primary, FeeSplitter.Schedule.A);
    }

    // ─────────────────────────────────────────────────────────────
    // Admin
    // ─────────────────────────────────────────────────────────────

    function test_setTreasury() public {
        address newTreasury = makeAddr("newTreasury");
        vm.expectEmit(true, true, false, false);
        emit TreasuryUpdated(treasury, newTreasury);
        vm.prank(admin);
        splitter.setTreasury(newTreasury);
        assertEq(splitter.treasury(), newTreasury);
    }

    function test_setTreasury_revert_zero() public {
        vm.expectRevert(FeeSplitter.ZeroAddress.selector);
        vm.prank(admin);
        splitter.setTreasury(address(0));
    }

    function test_setTreasury_revert_notAdmin() public {
        vm.expectRevert();
        vm.prank(stranger);
        splitter.setTreasury(makeAddr("x"));
    }

    function test_setBuybackBurner() public {
        address newBurner = makeAddr("newBurner");
        vm.expectEmit(true, true, false, false);
        emit BuybackBurnerUpdated(buyback, newBurner);
        vm.prank(admin);
        splitter.setBuybackBurner(newBurner);
        assertEq(splitter.buybackBurner(), newBurner);
    }

    function test_setBuybackBurner_revert_zero() public {
        vm.expectRevert(FeeSplitter.ZeroAddress.selector);
        vm.prank(admin);
        splitter.setBuybackBurner(address(0));
    }

    // ─────────────────────────────────────────────────────────────
    // Fuzz: sum invariant
    // ─────────────────────────────────────────────────────────────

    function testFuzz_split_sumInvariant_scheduleA(uint256 amount) public {
        amount = bound(amount, 1, type(uint128).max);
        token.mint(caller, amount);
        vm.prank(caller);
        token.approve(address(splitter), type(uint256).max);

        uint256 callerBefore = token.balanceOf(caller);

        vm.prank(caller);
        splitter.split(address(token), amount, primary, FeeSplitter.Schedule.A);

        uint256 totalOut = token.balanceOf(primary) + token.balanceOf(treasury) + token.balanceOf(buyback);
        assertEq(totalOut, amount);
        assertEq(token.balanceOf(caller), callerBefore - amount);
        assertEq(token.balanceOf(address(splitter)), 0);
    }

    function testFuzz_split_sumInvariant_scheduleB(uint256 amount) public {
        amount = bound(amount, 1, type(uint128).max);
        token.mint(caller, amount);
        vm.prank(caller);
        token.approve(address(splitter), type(uint256).max);

        vm.prank(caller);
        splitter.split(address(token), amount, primary, FeeSplitter.Schedule.B);

        uint256 totalOut = token.balanceOf(primary) + token.balanceOf(treasury) + token.balanceOf(buyback);
        assertEq(totalOut, amount);
        assertEq(token.balanceOf(address(splitter)), 0);
    }

    function testFuzz_previewSplit_sumInvariant(uint256 amount) public view {
        amount = bound(amount, 0, type(uint128).max);
        (uint256 p, uint256 t, uint256 b) = splitter.previewSplit(amount, FeeSplitter.Schedule.A);
        assertEq(p + t + b, amount);
        (p, t, b) = splitter.previewSplit(amount, FeeSplitter.Schedule.B);
        assertEq(p + t + b, amount);
    }
}
