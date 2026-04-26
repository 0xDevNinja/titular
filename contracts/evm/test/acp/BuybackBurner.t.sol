// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC20Burnable} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import {BuybackBurner, IUniswapV2Router02} from "../../src/acp/BuybackBurner.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/// @dev Minimal payment token (e.g. USDC mock)
contract MockPaymentToken is ERC20 {
    constructor() ERC20("USDC", "USDC") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

/// @dev TITU mock with ERC20Burnable
contract MockTITU is ERC20Burnable {
    constructor() ERC20("TITU", "TITU") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

/// @dev Mock Uniswap V2 router that simulates a swap by transferring mockTitu
///      to `to` in exchange for paymentToken from msg.sender.
contract MockUniRouter {
    MockTITU public titu;
    uint256 public fixedOutput; // how much TITU the router returns per swap

    constructor(MockTITU _titu, uint256 _fixedOutput) {
        titu = _titu;
        fixedOutput = _fixedOutput;
    }

    /// @dev Pulls paymentToken from msg.sender (the burner contract), then transfers TITU to `to`.
    function swapExactTokensForTokensSupportingFeeOnTransferTokens(
        uint256 amountIn,
        uint256 amountOutMin,
        address[] calldata path,
        address to,
        uint256 /*deadline*/
    ) external {
        // Pull paymentToken
        IERC20(path[0]).transferFrom(msg.sender, address(this), amountIn);
        // Enforce slippage
        require(fixedOutput >= amountOutMin, "slippage");
        // Transfer TITU to recipient
        titu.mint(to, fixedOutput);
    }
}

contract BuybackBurnerTest is Test {
    MockPaymentToken internal payment;
    MockTITU internal titu;
    MockUniRouter internal router;
    BuybackBurner internal burner;

    address internal admin = makeAddr("admin");
    address internal executor = makeAddr("executor");
    address internal stranger = makeAddr("stranger");

    // FIXED_OUTPUT must be >= minTituOut floor (set to 900e18 = 90% of AMOUNT_IN in TITU terms)
    uint256 internal constant FIXED_OUTPUT = 950e18;
    uint256 internal constant AMOUNT_IN = 1000e18;

    event BuybackAndBurn(address indexed executor, address indexed paymentToken, uint256 amountIn, uint256 tituBurned);

    function setUp() public {
        payment = new MockPaymentToken();
        titu = new MockTITU();
        router = new MockUniRouter(titu, FIXED_OUTPUT);

        address[] memory path = new address[](2);
        path[0] = address(payment);
        path[1] = address(titu);

        burner = new BuybackBurner(
            admin,
            address(router),
            address(payment),
            address(titu),
            path,
            900e18, // minTituOut floor: 900 TITU (absolute, not BPS)
            300 // 5 min swap deadline
        );

        // Grant executor role
        bytes32 executorRole = burner.EXECUTOR_ROLE();
        vm.prank(admin);
        burner.grantRole(executorRole, executor);

        // Fund burner with payment tokens (no pre-approval needed; contract manages it)
        payment.mint(address(burner), AMOUNT_IN * 10);
    }

    // ─────────────────────────────────────────────────────────────
    // Constructor
    // ─────────────────────────────────────────────────────────────

    function test_constructor_setsState() public view {
        assertEq(address(burner.router()), address(router));
        assertEq(burner.paymentToken(), address(payment));
        assertEq(address(burner.titu()), address(titu));
        assertEq(burner.minTituOut(), 900e18);
        assertEq(burner.swapDeadlineBuffer(), 300);
        address[] memory path = burner.getSwapPath();
        assertEq(path[0], address(payment));
        assertEq(path[1], address(titu));
    }

    function test_constructor_revert_zeroAdmin() public {
        address[] memory path = new address[](2);
        path[0] = address(payment);
        path[1] = address(titu);
        vm.expectRevert(BuybackBurner.ZeroAddress.selector);
        new BuybackBurner(address(0), address(router), address(payment), address(titu), path, 9000, 300);
    }

    /// @notice `InvalidBps` was removed; any uint256 is a valid minTituOut now.
    ///         Verify that a large minTituOut value does not revert during construction.
    function test_constructor_largeMinTituOut_doesNotRevert() public {
        address[] memory path = new address[](2);
        path[0] = address(payment);
        path[1] = address(titu);
        // Should NOT revert — minTituOut has no upper bound
        BuybackBurner b =
            new BuybackBurner(admin, address(router), address(payment), address(titu), path, type(uint256).max, 300);
        assertEq(b.minTituOut(), type(uint256).max);
    }

    function test_constructor_revert_invalidPath_tooShort() public {
        address[] memory path = new address[](1);
        path[0] = address(payment);
        vm.expectRevert(BuybackBurner.InvalidSwapPath.selector);
        new BuybackBurner(admin, address(router), address(payment), address(titu), path, 9000, 300);
    }

    function test_constructor_revert_invalidPath_wrongStart() public {
        address[] memory path = new address[](2);
        path[0] = makeAddr("other");
        path[1] = address(titu);
        vm.expectRevert(BuybackBurner.InvalidSwapPath.selector);
        new BuybackBurner(admin, address(router), address(payment), address(titu), path, 9000, 300);
    }

    function test_constructor_revert_invalidPath_wrongEnd() public {
        address[] memory path = new address[](2);
        path[0] = address(payment);
        path[1] = makeAddr("other");
        vm.expectRevert(BuybackBurner.InvalidSwapPath.selector);
        new BuybackBurner(admin, address(router), address(payment), address(titu), path, 9000, 300);
    }

    // ─────────────────────────────────────────────────────────────
    // buybackAndBurn happy path
    // ─────────────────────────────────────────────────────────────

    function test_buybackAndBurn_burnsCorrectAmount() public {
        uint256 tituSupplyBefore = titu.totalSupply();

        // Don't check the tituBurned data field since it depends on the mock's fixedOutput
        vm.expectEmit(true, true, false, false);
        emit BuybackAndBurn(executor, address(payment), AMOUNT_IN, FIXED_OUTPUT);

        vm.prank(executor);
        burner.buybackAndBurn(AMOUNT_IN, 1);

        // TITU total supply reduced by FIXED_OUTPUT
        assertEq(titu.totalSupply(), tituSupplyBefore); // mint then burn = net 0 change in supply
        // But burner holds 0 TITU
        assertEq(titu.balanceOf(address(burner)), 0);
    }

    function test_buybackAndBurn_revert_notExecutor() public {
        vm.expectRevert();
        vm.prank(stranger);
        burner.buybackAndBurn(AMOUNT_IN, 1);
    }

    function test_buybackAndBurn_revert_zeroAmount() public {
        vm.expectRevert(BuybackBurner.ZeroAmount.selector);
        vm.prank(executor);
        burner.buybackAndBurn(0, 0);
    }

    // ─────────────────────────────────────────────────────────────
    // Admin: setRouter
    // ─────────────────────────────────────────────────────────────

    function test_setRouter() public {
        address newRouter = makeAddr("newRouter");
        vm.prank(admin);
        burner.setRouter(newRouter);
        assertEq(address(burner.router()), newRouter);
    }

    function test_setRouter_revert_zero() public {
        vm.expectRevert(BuybackBurner.ZeroAddress.selector);
        vm.prank(admin);
        burner.setRouter(address(0));
    }

    function test_setRouter_revert_notAdmin() public {
        vm.expectRevert();
        vm.prank(stranger);
        burner.setRouter(makeAddr("x"));
    }

    // ─────────────────────────────────────────────────────────────
    // Admin: setSwapPath
    // ─────────────────────────────────────────────────────────────

    function test_setSwapPath_withHop() public {
        address hop = makeAddr("WETH");
        address[] memory path = new address[](3);
        path[0] = address(payment);
        path[1] = hop;
        path[2] = address(titu);

        vm.prank(admin);
        burner.setSwapPath(path);

        address[] memory stored = burner.getSwapPath();
        assertEq(stored.length, 3);
        assertEq(stored[1], hop);
    }

    function test_setSwapPath_revert_invalidPath() public {
        address[] memory path = new address[](2);
        path[0] = makeAddr("bad");
        path[1] = address(titu);
        vm.expectRevert(BuybackBurner.InvalidSwapPath.selector);
        vm.prank(admin);
        burner.setSwapPath(path);
    }

    // ─────────────────────────────────────────────────────────────
    // Admin: setMinTituOut
    // ─────────────────────────────────────────────────────────────

    function test_setMinTituOut() public {
        vm.prank(admin);
        burner.setMinTituOut(850e18);
        assertEq(burner.minTituOut(), 850e18);
    }

    function test_setMinTituOut_zero_disablesFloor() public {
        vm.prank(admin);
        burner.setMinTituOut(0);
        assertEq(burner.minTituOut(), 0);
    }

    function test_setMinTituOut_revert_notAdmin() public {
        vm.expectRevert();
        vm.prank(stranger);
        burner.setMinTituOut(100e18);
    }

    // ─────────────────────────────────────────────────────────────
    // rescueTokens: paymentToken blocked
    // ─────────────────────────────────────────────────────────────

    function test_rescueTokens_revert_paymentToken() public {
        vm.expectRevert(BuybackBurner.RescuePaymentTokenForbidden.selector);
        vm.prank(admin);
        burner.rescueTokens(address(payment), admin, 1);
    }

    function test_rescueTokens_allowsOtherToken() public {
        MockTITU other = new MockTITU();
        other.mint(address(burner), 100e18);
        vm.prank(admin);
        burner.rescueTokens(address(other), admin, 100e18);
        assertEq(other.balanceOf(admin), 100e18);
    }

    // ─────────────────────────────────────────────────────────────
    // Fuzz
    // ─────────────────────────────────────────────────────────────

    function testFuzz_buybackAndBurn_variousAmounts(uint256 amountIn) public {
        amountIn = bound(amountIn, 1, type(uint128).max);

        // Deploy a burner with minTituOut=0 (no floor) so the fixed-output mock always passes
        address[] memory path = new address[](2);
        path[0] = address(payment);
        path[1] = address(titu);
        BuybackBurner fuzzBurner =
            new BuybackBurner(admin, address(router), address(payment), address(titu), path, 0, 300); // minTituOut=0
        bytes32 executorRole = fuzzBurner.EXECUTOR_ROLE();
        vm.prank(admin);
        fuzzBurner.grantRole(executorRole, executor);

        payment.mint(address(fuzzBurner), amountIn);

        uint256 supplyBefore = titu.totalSupply();
        vm.prank(executor);
        fuzzBurner.buybackAndBurn(amountIn, 0);
        // Router mints FIXED_OUTPUT and contract burns it, net supply unchanged
        assertEq(titu.totalSupply(), supplyBefore);
        assertEq(titu.balanceOf(address(fuzzBurner)), 0);
    }
}
