// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {VestingVault} from "../../src/vault/VestingVault.sol";

contract VestingHandler is Test {
    VestingVault public immutable vault;
    IERC20 public immutable token;
    address public immutable admin;
    address[] public beneficiaries;
    uint256 public totalDeposited;
    uint256 public totalReleased;

    constructor(VestingVault _vault, IERC20 _token, address _admin) {
        vault = _vault;
        token = _token;
        admin = _admin;
        beneficiaries.push(address(0xAA));
        beneficiaries.push(address(0xBB));
        beneficiaries.push(address(0xCC));
    }

    function _pick(uint256 i) internal view returns (address) {
        return beneficiaries[i % beneficiaries.length];
    }

    function addGrant(uint96 amount, uint32 dur, uint256 idx) external {
        address b = _pick(idx);
        (uint256 total,,,,) = vault.grants(b);
        if (total != 0) return;
        uint256 a = bound(uint256(amount), 1, 10_000e18);
        uint64 d = uint64(bound(uint256(dur), 1 days, 4 * 365 days));
        vm.startPrank(admin);
        token.approve(address(vault), a);
        vault.addGrant(b, a, uint64(block.timestamp), 0, d);
        vm.stopPrank();
        totalDeposited += a;
    }

    function warp(uint32 delta) external {
        vm.warp(block.timestamp + uint256(delta));
    }

    function release(uint256 idx) external {
        address b = _pick(idx);
        uint256 before = token.balanceOf(b);
        try vault.release(b) {
            totalReleased += token.balanceOf(b) - before;
        } catch {}
    }
}
