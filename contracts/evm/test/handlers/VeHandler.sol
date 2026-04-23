// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {VeTITU} from "../../src/governance/VeTITU.sol";

interface IMint {
    function mint(address to, uint256 amount) external;
}

contract VeHandler is Test {
    VeTITU public immutable ve;
    IERC20 public immutable token;
    address[] public users;

    uint256 public constant WEEK = 7 days;
    uint256 public constant MAXTIME = 4 * 365 days;

    constructor(VeTITU _ve, IERC20 _token) {
        ve = _ve;
        token = _token;
        users.push(address(0xA1));
        users.push(address(0xB2));
        users.push(address(0xC3));
        for (uint256 i; i < users.length; ++i) {
            IMint(address(token)).mint(users[i], 100_000e18);
            vm.prank(users[i]);
            _token.approve(address(ve), type(uint256).max);
        }
    }

    function _pick(uint256 i) internal view returns (address) {
        return users[i % users.length];
    }

    function createLock(uint96 amount, uint32 dur, uint256 idx) external {
        address u = _pick(idx);
        (int128 amt,) = ve.locked(u);
        if (amt != 0) return;
        uint256 a = bound(uint256(amount), 1e18, 10_000e18);
        uint256 d = bound(uint256(dur), 2 * WEEK, MAXTIME);
        vm.prank(u);
        ve.createLock(a, block.timestamp + d);
    }

    function warp(uint32 delta) external {
        vm.warp(block.timestamp + uint256(delta));
    }

    function increaseAmount(uint96 amount, uint256 idx) external {
        address u = _pick(idx);
        (int128 amt, uint256 end) = ve.locked(u);
        if (amt == 0 || end <= block.timestamp) return;
        uint256 a = bound(uint256(amount), 1e18, 10_000e18);
        vm.prank(u);
        ve.increaseAmount(a);
    }
}
