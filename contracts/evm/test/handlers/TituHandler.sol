// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {TITU} from "../../src/token/TITU.sol";

/// @notice Handler that performs only supply-preserving operations against TITU.
///         Transfers and permits only — no burns — so `totalSupply` must remain fixed.
contract TituHandler is Test {
    TITU public immutable titu;
    address public immutable seed;
    address[] public actors;

    constructor(TITU _titu, address _seed) {
        titu = _titu;
        seed = _seed;
        actors.push(address(0xA11CE));
        actors.push(address(0xB0B));
        actors.push(address(0xCAFE));
        actors.push(address(0xDEAD));
    }

    function _pick(uint256 i) internal view returns (address) {
        return actors[i % actors.length];
    }

    function transferFromSeed(uint96 amt, uint256 toIdx) external {
        address to = _pick(toIdx);
        uint256 bal = titu.balanceOf(seed);
        if (bal == 0) return;
        uint256 amount = uint256(amt) % (bal + 1);
        vm.prank(seed);
        titu.transfer(to, amount);
    }

    function transferBetween(uint96 amt, uint256 fromIdx, uint256 toIdx) external {
        address from = _pick(fromIdx);
        address to = _pick(toIdx);
        if (from == to) return;
        uint256 bal = titu.balanceOf(from);
        if (bal == 0) return;
        uint256 amount = uint256(amt) % (bal + 1);
        vm.prank(from);
        titu.transfer(to, amount);
    }
}
