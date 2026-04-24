// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {FeeRouter} from "../../src/launchpad/FeeRouter.sol";

/// @dev Minimal mintable mock used for both TITU and the agent token. No tax so
///      split math stays exact down to the wei.
contract TestERC20 is ERC20 {
    constructor(string memory n, string memory s) ERC20(n, s) {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @dev Helper that happily accepts native ETH. Used as creator/treasury stand-in
///      for native-split tests where both legs must succeed.
contract ReceiveAccept {
    // solhint-disable-next-line no-empty-blocks
    receive() external payable {}
}

/// @dev Helper whose `receive()` always reverts. Used to assert that the router
///      surfaces failures from contract recipients rather than silently dropping
///      funds on the floor.
contract ReceiveRevert {
    error Rejected();

    receive() external payable {
        revert Rejected();
    }
}

contract FeeRouterTest is Test {
    TestERC20 internal token;
    FeeRouter internal router;

    address internal owner = address(0xA0);
    address internal treasury = address(0x7EA5);
    address internal creator = address(0xC0E);
    address internal agent = address(0xA6E17);
    address internal alice = address(0xA11CE);

    function setUp() public virtual {
        token = new TestERC20("TITU", "TITU");
        router = new FeeRouter(treasury, owner);
        token.mint(alice, 1_000_000e18);
        vm.prank(alice);
        token.approve(address(router), type(uint256).max);
    }

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    function test_constructor_reverts_zero_treasury() public {
        vm.expectRevert(FeeRouter.ZeroAddress.selector);
        new FeeRouter(address(0), owner);
    }

    function test_constructor_reverts_zero_owner() public {
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

    function test_setRoute_reverts_bps_too_high() public {
        vm.expectRevert(abi.encodeWithSelector(FeeRouter.BpsTooHigh.selector, uint16(10_001)));
        vm.prank(owner);
        router.setRoute(agent, creator, 10_001);
    }

    function test_setRoute_reverts_zero_creator() public {
        vm.expectRevert(FeeRouter.ZeroAddress.selector);
        vm.prank(owner);
        router.setRoute(agent, address(0), 7000);
    }

    function test_setRoute_reverts_zero_agent() public {
        vm.expectRevert(FeeRouter.ZeroAddress.selector);
        vm.prank(owner);
        router.setRoute(address(0), creator, 7000);
    }

    function test_setRoute_stores_and_emits() public {
        vm.expectEmit(true, true, false, true, address(router));
        emit FeeRouter.RouteSet(agent, creator, 7000);
        vm.prank(owner);
        router.setRoute(agent, creator, 7000);

        (address storedCreator, uint16 bps, bool configured) = router.routes(agent);
        assertEq(storedCreator, creator);
        assertEq(uint256(bps), 7000);
        assertTrue(configured);
    }

    function test_clearRoute_only_owner() public {
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        vm.prank(alice);
        router.clearRoute(agent);
    }

    function test_clearRoute_resets_state_and_emits() public {
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

    function test_distribute_erc20_splits_70_30_default() public {
        uint16 defaultBps = router.DEFAULT_CREATOR_BPS();
        vm.prank(owner);
        router.setRoute(agent, creator, defaultBps);

        uint256 amount = 1000e18;
        vm.expectEmit(true, true, false, true, address(router));
        emit FeeRouter.Distributed(agent, address(token), 700e18, 300e18);
        vm.prank(alice);
        router.distribute(agent, IERC20(address(token)), amount);

        assertEq(token.balanceOf(creator), 700e18);
        assertEq(token.balanceOf(treasury), 300e18);
        assertEq(token.balanceOf(address(router)), 0);
    }

    function test_distribute_erc20_splits_configured_ratio() public {
        // 40% creator / 60% treasury.
        vm.prank(owner);
        router.setRoute(agent, creator, 4000);

        uint256 amount = 500e18;
        vm.prank(alice);
        router.distribute(agent, IERC20(address(token)), amount);

        assertEq(token.balanceOf(creator), 200e18);
        assertEq(token.balanceOf(treasury), 300e18);
        assertEq(token.balanceOf(address(router)), 0);
    }

    function test_distribute_erc20_unconfigured_all_to_treasury() public {
        uint256 amount = 123e18;
        vm.expectEmit(true, true, false, true, address(router));
        emit FeeRouter.Distributed(agent, address(token), 0, amount);
        vm.prank(alice);
        router.distribute(agent, IERC20(address(token)), amount);

        assertEq(token.balanceOf(treasury), amount);
        assertEq(token.balanceOf(creator), 0);
        assertEq(token.balanceOf(address(router)), 0);
    }

    function test_distribute_erc20_cleared_route_all_to_treasury() public {
        vm.prank(owner);
        router.setRoute(agent, creator, 7000);
        vm.prank(owner);
        router.clearRoute(agent);

        vm.prank(alice);
        router.distribute(agent, IERC20(address(token)), 1000e18);

        assertEq(token.balanceOf(treasury), 1000e18);
        assertEq(token.balanceOf(creator), 0);
    }

    function test_distribute_erc20_rounding_dust_to_treasury() public {
        // 7000 bps of 7 wei = floor(7 * 7000 / 10_000) = floor(4.9) = 4 → treasury gets 3.
        vm.prank(owner);
        router.setRoute(agent, creator, 7000);

        uint256 amount = 7;
        vm.prank(alice);
        router.distribute(agent, IERC20(address(token)), amount);

        assertEq(token.balanceOf(creator), 4);
        assertEq(token.balanceOf(treasury), 3);
        assertEq(token.balanceOf(creator) + token.balanceOf(treasury), amount);
    }

    function test_distribute_reverts_zero_token() public {
        vm.expectRevert(FeeRouter.ZeroAddress.selector);
        vm.prank(alice);
        router.distribute(agent, IERC20(address(0)), 1);
    }

    // ---------------------------------------------------------------------
    // distributeNative
    // ---------------------------------------------------------------------

    function test_distributeNative_splits_correctly() public {
        ReceiveAccept creatorSink = new ReceiveAccept();
        ReceiveAccept treasurySink = new ReceiveAccept();
        FeeRouter r = new FeeRouter(address(treasurySink), owner);
        vm.prank(owner);
        r.setRoute(agent, address(creatorSink), 7000);

        vm.deal(alice, 10 ether);
        vm.expectEmit(true, true, false, true, address(r));
        emit FeeRouter.Distributed(agent, address(0), 0.7 ether, 0.3 ether);
        vm.prank(alice);
        r.distributeNative{value: 1 ether}(agent);

        assertEq(address(creatorSink).balance, 0.7 ether);
        assertEq(address(treasurySink).balance, 0.3 ether);
        assertEq(address(r).balance, 0);
    }

    function test_distributeNative_unconfigured_all_to_treasury() public {
        ReceiveAccept treasurySink = new ReceiveAccept();
        FeeRouter r = new FeeRouter(address(treasurySink), owner);

        vm.deal(alice, 10 ether);
        vm.prank(alice);
        r.distributeNative{value: 1 ether}(agent);

        assertEq(address(treasurySink).balance, 1 ether);
        assertEq(address(r).balance, 0);
    }

    function test_distributeNative_reverts_if_treasury_rejects() public {
        ReceiveRevert treasurySink = new ReceiveRevert();
        FeeRouter r = new FeeRouter(address(treasurySink), owner);
        // Route with zero creator share so the entire amount tries to land at the
        // reverting treasury.
        vm.prank(owner);
        r.setRoute(agent, creator, 0);

        vm.deal(alice, 10 ether);
        vm.expectRevert(abi.encodeWithSelector(FeeRouter.NativeTransferFailed.selector, address(treasurySink)));
        vm.prank(alice);
        r.distributeNative{value: 1 ether}(agent);
    }

    function test_distributeNative_reverts_if_creator_rejects() public {
        ReceiveAccept treasurySink = new ReceiveAccept();
        ReceiveRevert creatorSink = new ReceiveRevert();
        FeeRouter r = new FeeRouter(address(treasurySink), owner);
        vm.prank(owner);
        r.setRoute(agent, address(creatorSink), 7000);

        vm.deal(alice, 10 ether);
        vm.expectRevert(abi.encodeWithSelector(FeeRouter.NativeTransferFailed.selector, address(creatorSink)));
        vm.prank(alice);
        r.distributeNative{value: 1 ether}(agent);
    }

    function test_distributeNative_reverts_zero_value() public {
        vm.expectRevert(FeeRouter.InsufficientValue.selector);
        vm.prank(alice);
        router.distributeNative(agent);
    }

    // ---------------------------------------------------------------------
    // getSplit
    // ---------------------------------------------------------------------

    function test_getSplit_matches_actual_distribute() public {
        vm.prank(owner);
        router.setRoute(agent, creator, 7000);

        uint256 amount = 12_345e18 + 7;
        (uint256 previewCreator, uint256 previewTreasury) = router.getSplit(agent, amount);

        vm.prank(alice);
        router.distribute(agent, IERC20(address(token)), amount);

        assertEq(token.balanceOf(creator), previewCreator);
        assertEq(token.balanceOf(treasury), previewTreasury);
        assertEq(previewCreator + previewTreasury, amount);
    }

    function test_getSplit_unconfigured_returns_treasury_only() public view {
        (uint256 previewCreator, uint256 previewTreasury) = router.getSplit(agent, 1000);
        assertEq(previewCreator, 0);
        assertEq(previewTreasury, 1000);
    }

    // ---------------------------------------------------------------------
    // Fuzz
    // ---------------------------------------------------------------------

    /// @dev Any split must sum to the exact input amount (no dust escapes).
    function testFuzz_split_exact_sum(uint256 amount, uint16 creatorBps) public {
        amount = bound(amount, 0, type(uint128).max);
        creatorBps = uint16(bound(uint256(creatorBps), 0, 10_000));

        vm.prank(owner);
        router.setRoute(agent, creator, creatorBps);

        (uint256 creatorAmount, uint256 treasuryAmount) = router.getSplit(agent, amount);
        assertEq(creatorAmount + treasuryAmount, amount);
        // Creator share is floored; treasury owns the dust.
        assertLe(creatorAmount, (amount * creatorBps) / 10_000);
        assertGe(creatorAmount, (amount * creatorBps) / 10_000);
    }

    /// @dev End-to-end fuzz: configure a route and exercise the actual ERC-20 pull +
    ///      push path. Balances must reconcile exactly and the router never retains
    ///      a wei.
    function testFuzz_distribute_erc20(uint128 amount, uint16 creatorBps) public {
        vm.assume(amount > 0);
        uint16 bps = uint16(bound(uint256(creatorBps), 0, 10_000));

        vm.prank(owner);
        router.setRoute(agent, creator, bps);

        token.mint(alice, amount);
        vm.prank(alice);
        router.distribute(agent, IERC20(address(token)), amount);

        uint256 expectedCreator = (uint256(amount) * bps) / 10_000;
        uint256 expectedTreasury = uint256(amount) - expectedCreator;
        assertEq(token.balanceOf(creator), expectedCreator);
        assertEq(token.balanceOf(treasury), expectedTreasury);
        assertEq(token.balanceOf(address(router)), 0);
    }
}
