// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {Escrow} from "../../src/acp/Escrow.sol";

contract MockERC20 is ERC20 {
    constructor() ERC20("Mock", "MCK") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

contract EscrowTest is Test {
    Escrow internal escrow;
    MockERC20 internal tokenA;
    MockERC20 internal tokenB;

    address internal admin = makeAddr("admin");
    address internal releaser = makeAddr("releaser");
    address internal depositor = makeAddr("depositor");
    address internal recipient1 = makeAddr("recipient1");
    address internal recipient2 = makeAddr("recipient2");
    address internal stranger = makeAddr("stranger");

    uint256 internal constant JOB_ID = 42;
    uint256 internal constant AMOUNT = 1000e18;

    event Funded(address indexed depositor, uint256 indexed jobId, address indexed token, uint256 amount);
    event Released(address indexed depositor, uint256 indexed jobId, address indexed token, address to, uint256 amount);
    event Refunded(address indexed depositor, uint256 indexed jobId, address indexed token, uint256 amount);

    function setUp() public {
        escrow = new Escrow(admin, address(0x000000000022D473030F116dDEE9F6B43aC78BA3));
        tokenA = new MockERC20();
        tokenB = new MockERC20();

        // Grant releaser role
        bytes32 releaserRole = escrow.RELEASER_ROLE();
        vm.prank(admin);
        escrow.grantRole(releaserRole, releaser);

        // Mint and approve for depositor
        tokenA.mint(depositor, AMOUNT * 10);
        tokenB.mint(depositor, AMOUNT * 10);
        vm.prank(depositor);
        tokenA.approve(address(escrow), type(uint256).max);
        vm.prank(depositor);
        tokenB.approve(address(escrow), type(uint256).max);
    }

    // ─────────────────────────────────────────────────────────────
    // fund
    // ─────────────────────────────────────────────────────────────

    function test_fund_success() public {
        vm.expectEmit(true, true, true, true);
        emit Funded(depositor, JOB_ID, address(tokenA), AMOUNT);

        vm.prank(depositor);
        escrow.fund(depositor, JOB_ID, address(tokenA), AMOUNT);

        assertEq(escrow.getBalance(depositor, JOB_ID, address(tokenA)), AMOUNT);
        assertEq(tokenA.balanceOf(address(escrow)), AMOUNT);
    }

    function test_fund_accumulates() public {
        vm.prank(depositor);
        escrow.fund(depositor, JOB_ID, address(tokenA), AMOUNT);
        vm.prank(depositor);
        escrow.fund(depositor, JOB_ID, address(tokenA), AMOUNT);
        assertEq(escrow.getBalance(depositor, JOB_ID, address(tokenA)), AMOUNT * 2);
    }

    function test_fund_multiToken() public {
        vm.prank(depositor);
        escrow.fund(depositor, JOB_ID, address(tokenA), AMOUNT);
        vm.prank(depositor);
        escrow.fund(depositor, JOB_ID, address(tokenB), AMOUNT / 2);

        assertEq(escrow.getBalance(depositor, JOB_ID, address(tokenA)), AMOUNT);
        assertEq(escrow.getBalance(depositor, JOB_ID, address(tokenB)), AMOUNT / 2);
    }

    function test_fund_differentJobIds() public {
        vm.prank(depositor);
        escrow.fund(depositor, 1, address(tokenA), AMOUNT);
        vm.prank(depositor);
        escrow.fund(depositor, 2, address(tokenA), AMOUNT / 2);

        assertEq(escrow.getBalance(depositor, 1, address(tokenA)), AMOUNT);
        assertEq(escrow.getBalance(depositor, 2, address(tokenA)), AMOUNT / 2);
    }

    function test_fund_revert_zeroDepositor() public {
        vm.expectRevert(Escrow.ZeroAddress.selector);
        escrow.fund(address(0), JOB_ID, address(tokenA), AMOUNT);
    }

    function test_fund_revert_zeroToken() public {
        vm.expectRevert(Escrow.ZeroAddress.selector);
        escrow.fund(depositor, JOB_ID, address(0), AMOUNT);
    }

    function test_fund_revert_zeroAmount() public {
        vm.expectRevert(Escrow.ZeroAmount.selector);
        escrow.fund(depositor, JOB_ID, address(tokenA), 0);
    }

    // ─────────────────────────────────────────────────────────────
    // release
    // ─────────────────────────────────────────────────────────────

    function test_release_singleRecipient() public {
        vm.prank(depositor);
        escrow.fund(depositor, JOB_ID, address(tokenA), AMOUNT);

        Escrow.Recipient[] memory recs = new Escrow.Recipient[](1);
        recs[0] = Escrow.Recipient({to: recipient1, amount: AMOUNT});

        vm.expectEmit(true, true, true, true);
        emit Released(depositor, JOB_ID, address(tokenA), recipient1, AMOUNT);

        vm.prank(releaser);
        escrow.release(depositor, JOB_ID, address(tokenA), recs);

        assertEq(escrow.getBalance(depositor, JOB_ID, address(tokenA)), 0);
        assertEq(tokenA.balanceOf(recipient1), AMOUNT);
    }

    function test_release_multipleRecipients() public {
        vm.prank(depositor);
        escrow.fund(depositor, JOB_ID, address(tokenA), AMOUNT);

        Escrow.Recipient[] memory recs = new Escrow.Recipient[](2);
        recs[0] = Escrow.Recipient({to: recipient1, amount: AMOUNT * 95 / 100});
        recs[1] = Escrow.Recipient({to: recipient2, amount: AMOUNT * 5 / 100});

        vm.prank(releaser);
        escrow.release(depositor, JOB_ID, address(tokenA), recs);

        assertEq(tokenA.balanceOf(recipient1), AMOUNT * 95 / 100);
        assertEq(tokenA.balanceOf(recipient2), AMOUNT * 5 / 100);
        assertEq(escrow.getBalance(depositor, JOB_ID, address(tokenA)), 0);
    }

    function test_release_partialAmount() public {
        vm.prank(depositor);
        escrow.fund(depositor, JOB_ID, address(tokenA), AMOUNT);

        Escrow.Recipient[] memory recs = new Escrow.Recipient[](1);
        recs[0] = Escrow.Recipient({to: recipient1, amount: AMOUNT / 2});

        vm.prank(releaser);
        escrow.release(depositor, JOB_ID, address(tokenA), recs);

        assertEq(escrow.getBalance(depositor, JOB_ID, address(tokenA)), AMOUNT / 2);
        assertEq(tokenA.balanceOf(recipient1), AMOUNT / 2);
    }

    function test_release_revert_notReleaser() public {
        vm.prank(depositor);
        escrow.fund(depositor, JOB_ID, address(tokenA), AMOUNT);

        Escrow.Recipient[] memory recs = new Escrow.Recipient[](1);
        recs[0] = Escrow.Recipient({to: recipient1, amount: AMOUNT});

        vm.expectRevert();
        vm.prank(stranger);
        escrow.release(depositor, JOB_ID, address(tokenA), recs);
    }

    function test_release_revert_emptyRecipients() public {
        vm.prank(depositor);
        escrow.fund(depositor, JOB_ID, address(tokenA), AMOUNT);

        Escrow.Recipient[] memory recs = new Escrow.Recipient[](0);
        vm.expectRevert(Escrow.EmptyRecipients.selector);
        vm.prank(releaser);
        escrow.release(depositor, JOB_ID, address(tokenA), recs);
    }

    function test_release_revert_recipientZeroAddress() public {
        vm.prank(depositor);
        escrow.fund(depositor, JOB_ID, address(tokenA), AMOUNT);

        Escrow.Recipient[] memory recs = new Escrow.Recipient[](1);
        recs[0] = Escrow.Recipient({to: address(0), amount: AMOUNT});

        vm.expectRevert(abi.encodeWithSelector(Escrow.RecipientZeroAddress.selector, 0));
        vm.prank(releaser);
        escrow.release(depositor, JOB_ID, address(tokenA), recs);
    }

    function test_release_revert_recipientZeroAmount() public {
        vm.prank(depositor);
        escrow.fund(depositor, JOB_ID, address(tokenA), AMOUNT);

        Escrow.Recipient[] memory recs = new Escrow.Recipient[](1);
        recs[0] = Escrow.Recipient({to: recipient1, amount: 0});

        vm.expectRevert(abi.encodeWithSelector(Escrow.RecipientZeroAmount.selector, 0));
        vm.prank(releaser);
        escrow.release(depositor, JOB_ID, address(tokenA), recs);
    }

    function test_release_revert_sumExceedsBalance() public {
        vm.prank(depositor);
        escrow.fund(depositor, JOB_ID, address(tokenA), AMOUNT);

        Escrow.Recipient[] memory recs = new Escrow.Recipient[](1);
        recs[0] = Escrow.Recipient({to: recipient1, amount: AMOUNT + 1});

        vm.expectRevert(abi.encodeWithSelector(Escrow.SumExceedsBalance.selector, AMOUNT + 1, AMOUNT));
        vm.prank(releaser);
        escrow.release(depositor, JOB_ID, address(tokenA), recs);
    }

    // ─────────────────────────────────────────────────────────────
    // refund
    // ─────────────────────────────────────────────────────────────

    function test_refund_success() public {
        vm.prank(depositor);
        escrow.fund(depositor, JOB_ID, address(tokenA), AMOUNT);

        uint256 beforeBal = tokenA.balanceOf(depositor);

        vm.expectEmit(true, true, true, true);
        emit Refunded(depositor, JOB_ID, address(tokenA), AMOUNT);

        vm.prank(releaser);
        escrow.refund(depositor, JOB_ID, address(tokenA));

        assertEq(escrow.getBalance(depositor, JOB_ID, address(tokenA)), 0);
        assertEq(tokenA.balanceOf(depositor), beforeBal + AMOUNT);
    }

    function test_refund_revert_zeroBalance() public {
        vm.expectRevert(abi.encodeWithSelector(Escrow.InsufficientBalance.selector, 0, 1));
        vm.prank(releaser);
        escrow.refund(depositor, JOB_ID, address(tokenA));
    }

    function test_refund_revert_notReleaser() public {
        vm.prank(depositor);
        escrow.fund(depositor, JOB_ID, address(tokenA), AMOUNT);

        vm.expectRevert();
        vm.prank(stranger);
        escrow.refund(depositor, JOB_ID, address(tokenA));
    }

    // ─────────────────────────────────────────────────────────────
    // Fuzz
    // ─────────────────────────────────────────────────────────────

    function testFuzz_fund_release_invariant(uint256 amount) public {
        amount = bound(amount, 1, type(uint128).max);
        tokenA.mint(depositor, amount);

        vm.prank(depositor);
        escrow.fund(depositor, JOB_ID, address(tokenA), amount);
        assertEq(escrow.getBalance(depositor, JOB_ID, address(tokenA)), amount);

        Escrow.Recipient[] memory recs = new Escrow.Recipient[](1);
        recs[0] = Escrow.Recipient({to: recipient1, amount: amount});
        vm.prank(releaser);
        escrow.release(depositor, JOB_ID, address(tokenA), recs);

        assertEq(escrow.getBalance(depositor, JOB_ID, address(tokenA)), 0);
        assertEq(tokenA.balanceOf(recipient1), amount);
    }

    function testFuzz_fund_refund_invariant(uint256 amount) public {
        amount = bound(amount, 1, type(uint128).max);
        tokenA.mint(depositor, amount);
        uint256 initialBalance = tokenA.balanceOf(depositor);

        vm.prank(depositor);
        escrow.fund(depositor, JOB_ID, address(tokenA), amount);

        vm.prank(releaser);
        escrow.refund(depositor, JOB_ID, address(tokenA));

        assertEq(tokenA.balanceOf(depositor), initialBalance);
        assertEq(escrow.getBalance(depositor, JOB_ID, address(tokenA)), 0);
    }

    function testFuzz_multiRecipient_split(uint256 part1, uint256 part2) public {
        part1 = bound(part1, 1, AMOUNT - 1);
        part2 = bound(part2, 1, AMOUNT - part1);

        vm.prank(depositor);
        escrow.fund(depositor, JOB_ID, address(tokenA), AMOUNT);

        Escrow.Recipient[] memory recs = new Escrow.Recipient[](2);
        recs[0] = Escrow.Recipient({to: recipient1, amount: part1});
        recs[1] = Escrow.Recipient({to: recipient2, amount: part2});

        vm.prank(releaser);
        escrow.release(depositor, JOB_ID, address(tokenA), recs);

        assertEq(tokenA.balanceOf(recipient1), part1);
        assertEq(tokenA.balanceOf(recipient2), part2);
        assertEq(escrow.getBalance(depositor, JOB_ID, address(tokenA)), AMOUNT - part1 - part2);
    }
}
