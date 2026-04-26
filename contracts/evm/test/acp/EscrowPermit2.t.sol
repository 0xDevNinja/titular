// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {Escrow} from "../../src/acp/Escrow.sol";
import {IPermit2} from "../../src/acp/IPermit2.sol";

contract MockToken is ERC20 {
    constructor() ERC20("Mock", "MCK") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

contract MockPermit2 {
    function permitTransferFrom(
        IPermit2.PermitTransferFrom calldata permit,
        IPermit2.SignatureTransferDetails calldata transferDetails,
        address owner,
        bytes calldata
    ) external {
        ERC20(permit.permitted.token).transferFrom(owner, transferDetails.to, transferDetails.requestedAmount);
    }

    function permitWitnessTransferFrom(
        IPermit2.PermitTransferFrom calldata permit,
        IPermit2.SignatureTransferDetails calldata transferDetails,
        address owner,
        bytes32,
        string calldata,
        bytes calldata
    ) external {
        ERC20(permit.permitted.token).transferFrom(owner, transferDetails.to, transferDetails.requestedAmount);
    }
}

contract EscrowPermit2Test is Test {
    MockPermit2 internal mockPermit2;
    Escrow internal escrow;
    MockToken internal token;

    address internal admin = makeAddr("admin");
    address internal depositor = makeAddr("depositor");
    address internal releaser = makeAddr("releaser");

    uint256 internal constant AMOUNT = 500e18;
    uint256 internal constant JOB_ID = 7;

    event Funded(address indexed depositor, uint256 indexed jobId, address indexed token, uint256 amount);

    function setUp() public {
        mockPermit2 = new MockPermit2();
        escrow = new Escrow(admin, address(mockPermit2));
        token = new MockToken();

        bytes32 releaserRole = escrow.RELEASER_ROLE();
        vm.prank(admin);
        escrow.grantRole(releaserRole, releaser);

        token.mint(depositor, AMOUNT * 10);
        vm.prank(depositor);
        token.approve(address(mockPermit2), type(uint256).max);
    }

    function test_constructor_setsPermit2() public view {
        assertEq(address(escrow.permit2()), address(mockPermit2));
    }

    function test_fundWithPermit2_success() public {
        IPermit2.PermitTransferFrom memory permit = IPermit2.PermitTransferFrom({
            permitted: IPermit2.TokenPermissions({token: address(token), amount: AMOUNT}),
            nonce: 0,
            deadline: block.timestamp + 1 hours
        });
        vm.expectEmit(true, true, true, true);
        emit Funded(depositor, JOB_ID, address(token), AMOUNT);
        escrow.fundWithPermit2(depositor, JOB_ID, address(token), AMOUNT, permit, hex"1234");
        assertEq(escrow.getBalance(depositor, JOB_ID, address(token)), AMOUNT);
        assertEq(token.balanceOf(address(escrow)), AMOUNT);
    }

    function test_fundWithPermit2_revert_zeroDepositor() public {
        IPermit2.PermitTransferFrom memory permit = IPermit2.PermitTransferFrom({
            permitted: IPermit2.TokenPermissions({token: address(token), amount: AMOUNT}),
            nonce: 0,
            deadline: block.timestamp + 1 hours
        });
        vm.expectRevert(Escrow.ZeroAddress.selector);
        escrow.fundWithPermit2(address(0), JOB_ID, address(token), AMOUNT, permit, hex"");
    }

    function test_fundWithPermit2_revert_zeroToken() public {
        IPermit2.PermitTransferFrom memory permit = IPermit2.PermitTransferFrom({
            permitted: IPermit2.TokenPermissions({token: address(0), amount: AMOUNT}),
            nonce: 0,
            deadline: block.timestamp + 1 hours
        });
        vm.expectRevert(Escrow.ZeroAddress.selector);
        escrow.fundWithPermit2(depositor, JOB_ID, address(0), AMOUNT, permit, hex"");
    }

    function test_fundWithPermit2_revert_zeroAmount() public {
        IPermit2.PermitTransferFrom memory permit = IPermit2.PermitTransferFrom({
            permitted: IPermit2.TokenPermissions({token: address(token), amount: 0}),
            nonce: 0,
            deadline: block.timestamp + 1 hours
        });
        vm.expectRevert(Escrow.ZeroAmount.selector);
        escrow.fundWithPermit2(depositor, JOB_ID, address(token), 0, permit, hex"");
    }

    function test_fundWithPermit2_constructor_revert_zeroPermit2() public {
        vm.expectRevert(Escrow.ZeroAddress.selector);
        new Escrow(admin, address(0));
    }

    // ─────────────────────────────────────────────────────────────
    // Cross-token drain protection (#56)
    // ─────────────────────────────────────────────────────────────

    /// @notice Proves that a depositor cannot drain token Y by providing a permit for token X.
    ///         The `permit.permitted.token` must match the `token` argument exactly.
    function test_fundWithPermit2_revert_tokenMismatch() public {
        MockToken otherToken = new MockToken();
        // Permit is for `otherToken` but caller passes `token` as the target
        IPermit2.PermitTransferFrom memory permit = IPermit2.PermitTransferFrom({
            permitted: IPermit2.TokenPermissions({token: address(otherToken), amount: AMOUNT}),
            nonce: 0,
            deadline: block.timestamp + 1 hours
        });
        vm.expectRevert(
            abi.encodeWithSelector(Escrow.Permit2TokenMismatch.selector, address(otherToken), address(token))
        );
        escrow.fundWithPermit2(depositor, JOB_ID, address(token), AMOUNT, permit, hex"1234");
    }

    /// @notice Proves that amount in permit must match the requested amount.
    function test_fundWithPermit2_revert_amountMismatch() public {
        IPermit2.PermitTransferFrom memory permit = IPermit2.PermitTransferFrom({
            permitted: IPermit2.TokenPermissions({token: address(token), amount: AMOUNT / 2}),
            nonce: 0,
            deadline: block.timestamp + 1 hours
        });
        vm.expectRevert(abi.encodeWithSelector(Escrow.Permit2AmountMismatch.selector, AMOUNT / 2, AMOUNT));
        escrow.fundWithPermit2(depositor, JOB_ID, address(token), AMOUNT, permit, hex"1234");
    }

    /// @notice Proves the witness type constants are correctly set.
    function test_witnessTypeHash_isCorrect() public view {
        bytes32 expected = keccak256(bytes(escrow.ESCROW_WITNESS_TYPE()));
        assertEq(escrow.ESCROW_WITNESS_TYPEHASH(), expected);
    }

    /// @notice Standard fund() still works after Permit2 constructor change.
    function test_fund_standardPath_stillWorks() public {
        token.mint(address(this), AMOUNT);
        token.approve(address(escrow), AMOUNT);
        escrow.fund(depositor, JOB_ID, address(token), AMOUNT);
        assertEq(escrow.getBalance(depositor, JOB_ID, address(token)), AMOUNT);
    }

    function testFuzz_fundWithPermit2_amounts(uint256 amount) public {
        amount = bound(amount, 1, type(uint128).max);
        token.mint(depositor, amount);
        IPermit2.PermitTransferFrom memory permit = IPermit2.PermitTransferFrom({
            permitted: IPermit2.TokenPermissions({token: address(token), amount: amount}),
            nonce: 0,
            deadline: block.timestamp + 1 hours
        });
        escrow.fundWithPermit2(depositor, JOB_ID, address(token), amount, permit, hex"1234");
        assertGe(escrow.getBalance(depositor, JOB_ID, address(token)), amount);
    }

    function testFuzz_fundWithPermit2_and_release_invariant(uint256 amount) public {
        amount = bound(amount, 1, type(uint128).max);
        token.mint(depositor, amount);
        IPermit2.PermitTransferFrom memory permit = IPermit2.PermitTransferFrom({
            permitted: IPermit2.TokenPermissions({token: address(token), amount: amount}),
            nonce: 0,
            deadline: block.timestamp + 1 hours
        });
        escrow.fundWithPermit2(depositor, JOB_ID, address(token), amount, permit, hex"1234");

        address recipient = makeAddr("recv");
        Escrow.Recipient[] memory recs = new Escrow.Recipient[](1);
        recs[0] = Escrow.Recipient({to: recipient, amount: amount});
        vm.prank(releaser);
        escrow.release(depositor, JOB_ID, address(token), recs);

        assertEq(token.balanceOf(recipient), amount);
        assertEq(escrow.getBalance(depositor, JOB_ID, address(token)), 0);
    }
}
