// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {FeeDistributor} from "../../src/fees/FeeDistributor.sol";
import {IVeTITU} from "../../src/fees/IVeTITU.sol";

contract MockERC20 is ERC20 {
    constructor(string memory n, string memory s) ERC20(n, s) {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @dev Minimal mock exposing the subset FeeDistributor needs.
contract MockVe is IVeTITU {
    mapping(address => mapping(uint256 => uint256)) public balances;
    mapping(uint256 => uint256) public totals;

    function setBalance(address user, uint256 week, uint256 amount) external {
        balances[user][week] = amount;
    }

    function setTotal(uint256 week, uint256 amount) external {
        totals[week] = amount;
    }

    function getPastVotes(address user, uint256 tp) external view override returns (uint256) {
        return balances[user][(tp / 1 weeks) * 1 weeks];
    }

    function getPastTotalSupply(uint256 tp) external view override returns (uint256) {
        return totals[(tp / 1 weeks) * 1 weeks];
    }
}

contract FeeDistributorTest is Test {
    MockERC20 internal token;
    MockERC20 internal usdc;
    MockVe internal ve;
    FeeDistributor internal fd;
    address internal treasury = address(0x7EA);
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);

    uint256 internal constant WEEK = 7 days;

    function setUp() public {
        vm.warp(1_700_000_000);
        token = new MockERC20("R", "R");
        usdc = new MockERC20("USDC", "USDC");
        ve = new MockVe();
        fd = new FeeDistributor(IVeTITU(address(ve)), treasury);
    }

    function _currentWeek() internal view returns (uint256) {
        return (block.timestamp / WEEK) * WEEK;
    }

    function test_checkpointToken_records_delta() public {
        token.mint(address(fd), 1000e18);
        fd.checkpointToken(address(token));
        assertEq(fd.tokensPerWeek(address(token), _currentWeek()), 1000e18);
        // second call with no change = no delta added
        fd.checkpointToken(address(token));
        assertEq(fd.tokensPerWeek(address(token), _currentWeek()), 1000e18);
    }

    function test_claim_before_first_week_ends_pays_zero() public {
        token.mint(address(fd), 1000e18);
        fd.checkpointToken(address(token));
        address[] memory tokens = new address[](1);
        tokens[0] = address(token);
        fd.claim(alice, tokens);
        assertEq(token.balanceOf(alice), 0);
    }

    function test_claim_single_user_full_share() public {
        uint256 w = _currentWeek();
        ve.setBalance(alice, w, 100e18);
        ve.setTotal(w, 100e18);

        token.mint(address(fd), 1000e18);
        fd.checkpointToken(address(token));

        vm.warp(block.timestamp + WEEK);

        address[] memory tokens = new address[](1);
        tokens[0] = address(token);
        fd.claim(alice, tokens);
        assertEq(token.balanceOf(alice), 1000e18);
    }

    function test_claim_proportional_two_users() public {
        uint256 w = _currentWeek();
        ve.setBalance(alice, w, 300e18);
        ve.setBalance(bob, w, 100e18);
        ve.setTotal(w, 400e18);

        token.mint(address(fd), 800e18);
        fd.checkpointToken(address(token));
        vm.warp(block.timestamp + WEEK);

        address[] memory tokens = new address[](1);
        tokens[0] = address(token);
        fd.claim(alice, tokens);
        fd.claim(bob, tokens);
        assertEq(token.balanceOf(alice), 600e18);
        assertEq(token.balanceOf(bob), 200e18);
    }

    function test_claim_multi_token() public {
        uint256 w = _currentWeek();
        ve.setBalance(alice, w, 100e18);
        ve.setTotal(w, 100e18);

        token.mint(address(fd), 500e18);
        usdc.mint(address(fd), 1000e6);
        fd.checkpointToken(address(token));
        fd.checkpointToken(address(usdc));
        vm.warp(block.timestamp + WEEK);

        address[] memory tokens = new address[](2);
        tokens[0] = address(token);
        tokens[1] = address(usdc);
        fd.claim(alice, tokens);
        assertEq(token.balanceOf(alice), 500e18);
        assertEq(usdc.balanceOf(alice), 1000e6);
    }

    function test_claim_idempotent() public {
        uint256 w = _currentWeek();
        ve.setBalance(alice, w, 100e18);
        ve.setTotal(w, 100e18);
        token.mint(address(fd), 500e18);
        fd.checkpointToken(address(token));
        vm.warp(block.timestamp + WEEK);

        address[] memory tokens = new address[](1);
        tokens[0] = address(token);
        fd.claim(alice, tokens);
        uint256 bal = token.balanceOf(alice);
        fd.claim(alice, tokens); // second claim pays nothing
        assertEq(token.balanceOf(alice), bal);
    }

    function test_claim_advances_lastClaimWeek() public {
        uint256 w = _currentWeek();
        ve.setBalance(alice, w, 100e18);
        ve.setTotal(w, 100e18);
        token.mint(address(fd), 500e18);
        fd.checkpointToken(address(token));
        vm.warp(block.timestamp + WEEK);

        address[] memory tokens = new address[](1);
        tokens[0] = address(token);
        fd.claim(alice, tokens);
        assertGt(fd.userNextClaimWeek(alice, address(token)), 0);
    }

    function test_constructor_reverts_zero() public {
        vm.expectRevert(FeeDistributor.ZeroAddress.selector);
        new FeeDistributor(IVeTITU(address(0)), treasury);
    }

    function test_constructor_reverts_zero_treasury() public {
        vm.expectRevert(FeeDistributor.ZeroAddress.selector);
        new FeeDistributor(IVeTITU(address(ve)), address(0));
    }

    function test_claimable_view() public {
        uint256 w = _currentWeek();
        ve.setBalance(alice, w, 100e18);
        ve.setTotal(w, 100e18);
        token.mint(address(fd), 500e18);
        fd.checkpointToken(address(token));
        vm.warp(block.timestamp + WEEK);
        assertEq(fd.claimable(alice, address(token)), 500e18);
    }

    function test_claimable_zero_when_no_stake() public {
        token.mint(address(fd), 500e18);
        fd.checkpointToken(address(token));
        vm.warp(block.timestamp + WEEK);
        assertEq(fd.claimable(alice, address(token)), 0);
    }

    function test_claim_reverts_zero_user() public {
        address[] memory tokens = new address[](1);
        tokens[0] = address(token);
        vm.expectRevert(FeeDistributor.ZeroAddress.selector);
        fd.claim(address(0), tokens);
    }

    function test_checkpointToken_reverts_zero() public {
        vm.expectRevert(FeeDistributor.ZeroAddress.selector);
        fd.checkpointToken(address(0));
    }
}
