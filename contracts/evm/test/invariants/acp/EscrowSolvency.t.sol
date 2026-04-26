// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {StdInvariant} from "forge-std/StdInvariant.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {Escrow} from "../../../src/acp/Escrow.sol";

contract MockToken is ERC20 {
    constructor() ERC20("InvTest", "INV") {}
    function mint(address to, uint256 amount) external { _mint(to, amount); }
}

/// @dev Handler contract that wraps Escrow with bounded random calls for invariant testing.
contract EscrowHandler is Test {
    Escrow public escrow;
    MockToken public token;

    address internal admin;
    address internal depositor;
    address internal releaser;
    uint256 internal jobIdCounter;

    // Track expected total escrowed balance
    uint256 public totalDeposited;
    uint256 public totalReleased;
    uint256 public totalRefunded;

    constructor(Escrow _escrow, MockToken _token, address _admin, address _depositor, address _releaser) {
        escrow = _escrow;
        token = _token;
        admin = _admin;
        depositor = _depositor;
        releaser = _releaser;
    }

    function fund(uint256 amount, uint256 jobSeed) public {
        amount = bound(amount, 1, 1000e18);
        uint256 jobId = bound(jobSeed, 0, 9);
        token.mint(depositor, amount);
        vm.prank(depositor);
        token.approve(address(escrow), amount);
        vm.prank(depositor);
        escrow.fund(depositor, jobId, address(token), amount);
        totalDeposited += amount;
    }

    function refund(uint256 jobSeed) public {
        uint256 jobId = bound(jobSeed, 0, 9);
        uint256 bal = escrow.getBalance(depositor, jobId, address(token));
        if (bal == 0) return;
        vm.prank(releaser);
        escrow.refund(depositor, jobId, address(token));
        totalRefunded += bal;
    }

    function release(uint256 jobSeed, uint256 amount) public {
        uint256 jobId = bound(jobSeed, 0, 9);
        uint256 bal = escrow.getBalance(depositor, jobId, address(token));
        if (bal == 0) return;
        amount = bound(amount, 1, bal);

        address recipient = makeAddr("recipient");
        Escrow.Recipient[] memory recs = new Escrow.Recipient[](1);
        recs[0] = Escrow.Recipient({to: recipient, amount: amount});
        vm.prank(releaser);
        escrow.release(depositor, jobId, address(token), recs);
        totalReleased += amount;
    }
}

contract EscrowSolvencyInvariantTest is StdInvariant, Test {
    Escrow internal escrow;
    MockToken internal token;
    EscrowHandler internal handler;

    address internal admin = makeAddr("admin");
    address internal depositor = makeAddr("depositor");
    address internal releaser = makeAddr("releaser");

    function setUp() public {
        token = new MockToken();
        escrow = new Escrow(admin, makeAddr("permit2"));

        bytes32 releaserRole = escrow.RELEASER_ROLE();
        vm.prank(admin);
        escrow.grantRole(releaserRole, releaser);

        handler = new EscrowHandler(escrow, token, admin, depositor, releaser);
        targetContract(address(handler));
    }

    /// @notice The ERC-20 balance of the Escrow contract must always equal
    ///         totalDeposited - totalReleased - totalRefunded.
    function invariant_escrowContractBalanceMatchesAccounting() public view {
        uint256 contractBalance = token.balanceOf(address(escrow));
        uint256 expected = handler.totalDeposited() - handler.totalReleased() - handler.totalRefunded();
        assertEq(contractBalance, expected, "escrow contract balance mismatch");
    }

    /// @notice Released + refunded never exceeds deposited (no minting).
    function invariant_noFundsCreatedFromThinAir() public view {
        assertLe(
            handler.totalReleased() + handler.totalRefunded(),
            handler.totalDeposited(),
            "released+refunded exceeds deposited"
        );
    }
}
