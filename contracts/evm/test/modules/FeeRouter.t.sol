// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

import {FeeRouter} from "../../src/launchpad/FeeRouter.sol";

/// @dev Plain ERC-20 used for the ERC-20 split path.
contract FeeRouterMockToken is ERC20 {
    constructor(string memory n, string memory s) ERC20(n, s) {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @dev Sink that accepts any incoming ETH. Used as creator/treasury for
///      successful native-split paths.
contract FeeRouterAcceptETH {
    // solhint-disable-next-line no-empty-blocks
    receive() external payable {}
}

/// @dev Sink whose `receive()` always reverts; used to drive the
///      `NativeTransferFailed` revert path.
contract FeeRouterRejectETH {
    error Rejected();

    receive() external payable {
        revert Rejected();
    }
}

/// @title  FeeRouterModuleSuiteTest
/// @notice Per-contract unit suite covering every public surface of
///         {FeeRouter}: constructor validation, owner-gated route
///         configuration, ERC-20 + native distribution, dust handling, and
///         `getSplit` parity with `distribute`. Every state-mutating
///         external entry is exercised in both happy + revert directions
///         to satisfy the ≥95% module coverage gate.
contract FeeRouterModuleSuiteTest is Test {
    FeeRouterMockToken internal token;
    FeeRouter internal router;

    address internal owner = address(0xA0);
    address internal treasury = address(0x7EA5);
    address internal creator = address(0xC0E);
    address internal agent = address(0xA6E17);
    address internal alice = address(0xA11CE);

    function setUp() public {
        token = new FeeRouterMockToken("TITU", "TITU");
        router = new FeeRouter(treasury, owner);
        token.mint(alice, 1_000_000e18);
        vm.prank(alice);
        token.approve(address(router), type(uint256).max);
    }

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    function test_constructor_zero_treasury_reverts() public {
        vm.expectRevert(FeeRouter.ZeroAddress.selector);
        new FeeRouter(address(0), owner);
    }

    function test_constructor_zero_owner_reverts() public {
        vm.expectRevert(FeeRouter.ZeroAddress.selector);
        new FeeRouter(treasury, address(0));
    }

    function test_constructor_sets_immutables() public view {
        assertEq(router.treasury(), treasury);
        assertEq(router.owner(), owner);
        assertEq(uint256(router.DEFAULT_CREATOR_BPS()), 7000);
        assertEq(uint256(router.BPS_DENOMINATOR()), 10_000);
    }

    // ---------------------------------------------------------------------
    // setRoute / clearRoute
    // ---------------------------------------------------------------------

    function test_setRoute_only_owner() public {
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        vm.prank(alice);
        router.setRoute(agent, creator, 7000);
    }

    function test_setRoute_bps_too_high_reverts() public {
        vm.expectRevert(abi.encodeWithSelector(FeeRouter.BpsTooHigh.selector, uint16(10_001)));
        vm.prank(owner);
        router.setRoute(agent, creator, 10_001);
    }

    function test_setRoute_zero_creator_reverts() public {
        vm.expectRevert(FeeRouter.ZeroAddress.selector);
        vm.prank(owner);
        router.setRoute(agent, address(0), 7000);
    }

    function test_setRoute_zero_agent_reverts() public {
        vm.expectRevert(FeeRouter.ZeroAddress.selector);
        vm.prank(owner);
        router.setRoute(address(0), creator, 7000);
    }

    function test_setRoute_persists_and_emits() public {
        vm.expectEmit(true, true, false, true, address(router));
        emit FeeRouter.RouteSet(agent, creator, 7000);
        vm.prank(owner);
        router.setRoute(agent, creator, 7000);

        (address storedCreator, uint16 bps, bool configured) = router.routes(agent);
        assertEq(storedCreator, creator);
        assertEq(uint256(bps), 7000);
        assertTrue(configured);
    }

    function test_setRoute_reconfigures() public {
        vm.startPrank(owner);
        router.setRoute(agent, creator, 7000);
        // Reconfigure with different creator + bps — should overwrite.
        address creator2 = address(0xCAFE);
        router.setRoute(agent, creator2, 4000);
        vm.stopPrank();

        (address storedCreator, uint16 bps, bool configured) = router.routes(agent);
        assertEq(storedCreator, creator2);
        assertEq(uint256(bps), 4000);
        assertTrue(configured);
    }

    function test_clearRoute_only_owner() public {
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        vm.prank(alice);
        router.clearRoute(agent);
    }

    function test_clearRoute_zero_agent_reverts() public {
        vm.expectRevert(FeeRouter.ZeroAddress.selector);
        vm.prank(owner);
        router.clearRoute(address(0));
    }

    function test_clearRoute_resets_state() public {
        vm.prank(owner);
        router.setRoute(agent, creator, 7000);

        vm.expectEmit(true, false, false, true, address(router));
        emit FeeRouter.RouteCleared(agent);
        vm.prank(owner);
        router.clearRoute(agent);

        (address storedCreator, uint16 bps, bool configured) = router.routes(agent);
        assertEq(storedCreator, address(0));
        assertEq(uint256(bps), 0);
        assertFalse(configured);
    }

    // ---------------------------------------------------------------------
    // distribute (ERC-20)
    // ---------------------------------------------------------------------

    function test_distribute_default_70_30_split() public {
        uint16 defaultBps = router.DEFAULT_CREATOR_BPS();
        vm.prank(owner);
        router.setRoute(agent, creator, defaultBps);

        vm.expectEmit(true, true, false, true, address(router));
        emit FeeRouter.Distributed(agent, address(token), 700e18, 300e18);
        vm.prank(alice);
        router.distribute(agent, IERC20(address(token)), 1000e18);

        assertEq(token.balanceOf(creator), 700e18);
        assertEq(token.balanceOf(treasury), 300e18);
        assertEq(token.balanceOf(address(router)), 0);
    }

    function test_distribute_custom_split() public {
        vm.prank(owner);
        router.setRoute(agent, creator, 4000);

        vm.prank(alice);
        router.distribute(agent, IERC20(address(token)), 500e18);

        assertEq(token.balanceOf(creator), 200e18);
        assertEq(token.balanceOf(treasury), 300e18);
    }

    function test_distribute_unconfigured_routes_to_treasury() public {
        vm.expectEmit(true, true, false, true, address(router));
        emit FeeRouter.Distributed(agent, address(token), 0, 123e18);
        vm.prank(alice);
        router.distribute(agent, IERC20(address(token)), 123e18);

        assertEq(token.balanceOf(treasury), 123e18);
        assertEq(token.balanceOf(creator), 0);
    }

    function test_distribute_after_clear_routes_to_treasury() public {
        vm.startPrank(owner);
        router.setRoute(agent, creator, 7000);
        router.clearRoute(agent);
        vm.stopPrank();

        vm.prank(alice);
        router.distribute(agent, IERC20(address(token)), 1000e18);

        assertEq(token.balanceOf(treasury), 1000e18);
        assertEq(token.balanceOf(creator), 0);
    }

    function test_distribute_dust_to_treasury() public {
        // 7000 bps of 7 wei = floor(4.9) = 4. Treasury gets 3.
        vm.prank(owner);
        router.setRoute(agent, creator, 7000);

        vm.prank(alice);
        router.distribute(agent, IERC20(address(token)), 7);

        assertEq(token.balanceOf(creator), 4);
        assertEq(token.balanceOf(treasury), 3);
    }

    function test_distribute_full_to_creator() public {
        vm.prank(owner);
        router.setRoute(agent, creator, 10_000);

        vm.prank(alice);
        router.distribute(agent, IERC20(address(token)), 1000e18);

        assertEq(token.balanceOf(creator), 1000e18);
        assertEq(token.balanceOf(treasury), 0);
    }

    function test_distribute_zero_creator_share_all_to_treasury() public {
        vm.prank(owner);
        router.setRoute(agent, creator, 0);

        vm.prank(alice);
        router.distribute(agent, IERC20(address(token)), 1000e18);

        assertEq(token.balanceOf(creator), 0);
        assertEq(token.balanceOf(treasury), 1000e18);
    }

    function test_distribute_zero_token_reverts() public {
        vm.prank(alice);
        vm.expectRevert(FeeRouter.ZeroAddress.selector);
        router.distribute(agent, IERC20(address(0)), 1);
    }

    // ---------------------------------------------------------------------
    // distributeNative
    // ---------------------------------------------------------------------

    function test_distributeNative_split() public {
        FeeRouterAcceptETH cs = new FeeRouterAcceptETH();
        FeeRouterAcceptETH ts = new FeeRouterAcceptETH();
        FeeRouter r = new FeeRouter(address(ts), owner);
        vm.prank(owner);
        r.setRoute(agent, address(cs), 7000);

        vm.deal(alice, 10 ether);
        vm.expectEmit(true, true, false, true, address(r));
        emit FeeRouter.Distributed(agent, address(0), 0.7 ether, 0.3 ether);
        vm.prank(alice);
        r.distributeNative{value: 1 ether}(agent);

        assertEq(address(cs).balance, 0.7 ether);
        assertEq(address(ts).balance, 0.3 ether);
        assertEq(address(r).balance, 0);
    }

    function test_distributeNative_unconfigured_all_to_treasury() public {
        FeeRouterAcceptETH ts = new FeeRouterAcceptETH();
        FeeRouter r = new FeeRouter(address(ts), owner);

        vm.deal(alice, 10 ether);
        vm.prank(alice);
        r.distributeNative{value: 1 ether}(agent);

        assertEq(address(ts).balance, 1 ether);
        assertEq(address(r).balance, 0);
    }

    function test_distributeNative_treasury_revert_propagates() public {
        FeeRouterRejectETH ts = new FeeRouterRejectETH();
        FeeRouter r = new FeeRouter(address(ts), owner);
        vm.prank(owner);
        r.setRoute(agent, creator, 0); // 100% to treasury

        vm.deal(alice, 10 ether);
        vm.expectRevert(abi.encodeWithSelector(FeeRouter.NativeTransferFailed.selector, address(ts)));
        vm.prank(alice);
        r.distributeNative{value: 1 ether}(agent);
    }

    function test_distributeNative_creator_revert_propagates() public {
        FeeRouterAcceptETH ts = new FeeRouterAcceptETH();
        FeeRouterRejectETH cs = new FeeRouterRejectETH();
        FeeRouter r = new FeeRouter(address(ts), owner);
        vm.prank(owner);
        r.setRoute(agent, address(cs), 7000);

        vm.deal(alice, 10 ether);
        vm.expectRevert(abi.encodeWithSelector(FeeRouter.NativeTransferFailed.selector, address(cs)));
        vm.prank(alice);
        r.distributeNative{value: 1 ether}(agent);
    }

    function test_distributeNative_zero_value_reverts() public {
        vm.expectRevert(FeeRouter.InsufficientValue.selector);
        vm.prank(alice);
        router.distributeNative(agent);
    }

    // ---------------------------------------------------------------------
    // getSplit
    // ---------------------------------------------------------------------

    function test_getSplit_matches_distribute() public {
        vm.prank(owner);
        router.setRoute(agent, creator, 7000);

        uint256 amount = 12_345e18 + 7;
        (uint256 pc, uint256 pt) = router.getSplit(agent, amount);

        vm.prank(alice);
        router.distribute(agent, IERC20(address(token)), amount);

        assertEq(token.balanceOf(creator), pc);
        assertEq(token.balanceOf(treasury), pt);
        assertEq(pc + pt, amount);
    }

    function test_getSplit_unconfigured() public view {
        (uint256 pc, uint256 pt) = router.getSplit(agent, 1000);
        assertEq(pc, 0);
        assertEq(pt, 1000);
    }

    // ---------------------------------------------------------------------
    // Fuzz invariants
    // ---------------------------------------------------------------------

    /// @dev Splits sum to exactly the input, no dust escapes.
    function testFuzz_split_conservation(uint256 amount, uint16 bps) public {
        amount = bound(amount, 0, type(uint128).max);
        bps = uint16(bound(uint256(bps), 0, 10_000));

        vm.prank(owner);
        router.setRoute(agent, creator, bps);

        (uint256 c, uint256 t) = router.getSplit(agent, amount);
        assertEq(c + t, amount);
        assertEq(c, (amount * bps) / 10_000);
    }

    /// @dev End-to-end ERC-20 distribute fuzz: balances reconcile, router
    ///      retains zero.
    function testFuzz_distribute_erc20(uint128 amount, uint16 bps) public {
        vm.assume(amount > 0);
        bps = uint16(bound(uint256(bps), 0, 10_000));

        vm.prank(owner);
        router.setRoute(agent, creator, bps);
        token.mint(alice, amount);
        vm.prank(alice);
        router.distribute(agent, IERC20(address(token)), amount);

        uint256 expC = (uint256(amount) * bps) / 10_000;
        assertEq(token.balanceOf(creator), expC);
        assertEq(token.balanceOf(treasury), uint256(amount) - expC);
        assertEq(token.balanceOf(address(router)), 0);
    }
}
