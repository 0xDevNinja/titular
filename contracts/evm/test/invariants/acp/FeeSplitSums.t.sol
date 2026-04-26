// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {StdInvariant} from "forge-std/StdInvariant.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {FeeSplitter} from "../../../src/acp/FeeSplitter.sol";

contract MockToken is ERC20 {
    constructor() ERC20("FeeInv", "FINV") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

contract FeeSplitterHandler is Test {
    FeeSplitter public splitter;
    MockToken public token;

    address internal caller;
    address internal primary;
    address internal treasury;
    address internal buyback;

    uint256 public totalIn;
    uint256 public primaryOut;
    uint256 public treasuryOut;
    uint256 public buybackOut;

    constructor(
        FeeSplitter _splitter,
        MockToken _token,
        address _caller,
        address _primary,
        address _treasury,
        address _buyback
    ) {
        splitter = _splitter;
        token = _token;
        caller = _caller;
        primary = _primary;
        treasury = _treasury;
        buyback = _buyback;
    }

    function splitA(uint256 amount) public {
        amount = bound(amount, 1, 10_000e18);
        token.mint(caller, amount);
        vm.prank(caller);
        token.approve(address(splitter), amount);

        uint256 pBefore = token.balanceOf(primary);
        uint256 tBefore = token.balanceOf(treasury);
        uint256 bBefore = token.balanceOf(buyback);

        vm.prank(caller);
        splitter.split(address(token), amount, primary, FeeSplitter.Schedule.A);

        totalIn += amount;
        primaryOut += token.balanceOf(primary) - pBefore;
        treasuryOut += token.balanceOf(treasury) - tBefore;
        buybackOut += token.balanceOf(buyback) - bBefore;
    }

    function splitB(uint256 amount) public {
        amount = bound(amount, 1, 10_000e18);
        token.mint(caller, amount);
        vm.prank(caller);
        token.approve(address(splitter), amount);

        uint256 pBefore = token.balanceOf(primary);
        uint256 tBefore = token.balanceOf(treasury);
        uint256 bBefore = token.balanceOf(buyback);

        vm.prank(caller);
        splitter.split(address(token), amount, primary, FeeSplitter.Schedule.B);

        totalIn += amount;
        primaryOut += token.balanceOf(primary) - pBefore;
        treasuryOut += token.balanceOf(treasury) - tBefore;
        buybackOut += token.balanceOf(buyback) - bBefore;
    }
}

contract FeeSplitSumsInvariantTest is StdInvariant, Test {
    FeeSplitter internal splitter;
    MockToken internal token;
    FeeSplitterHandler internal handler;

    address internal admin = makeAddr("admin");
    address internal callerAddr = makeAddr("callerAddr");
    address internal primaryAddr = makeAddr("primaryAddr");
    address internal treasuryAddr = makeAddr("treasuryAddr");
    address internal buybackAddr = makeAddr("buybackAddr");

    function setUp() public {
        splitter = new FeeSplitter(admin, treasuryAddr, buybackAddr);
        token = new MockToken();

        bytes32 callerRole = splitter.CALLER_ROLE();
        vm.prank(admin);
        splitter.grantRole(callerRole, callerAddr);

        handler = new FeeSplitterHandler(splitter, token, callerAddr, primaryAddr, treasuryAddr, buybackAddr);
        targetContract(address(handler));
    }

    /// @notice Total output (primary + treasury + buyback) equals total input.
    function invariant_feeSplitSumsEqualInput() public view {
        assertEq(
            handler.primaryOut() + handler.treasuryOut() + handler.buybackOut(),
            handler.totalIn(),
            "fee split sums do not equal total input"
        );
    }

    /// @notice FeeSplitter contract should hold zero balance (nothing retained).
    function invariant_splitterHoldsZeroBalance() public view {
        assertEq(token.balanceOf(address(splitter)), 0, "splitter retains tokens");
    }
}
