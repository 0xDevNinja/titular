// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {StdInvariant} from "forge-std/StdInvariant.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {VeTITU} from "../../src/governance/VeTITU.sol";
import {VeHandler} from "../handlers/VeHandler.sol";

contract MockERC20 is ERC20 {
    constructor() ERC20("T", "T") {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

contract VeInvariant is StdInvariant, Test {
    MockERC20 internal token;
    VeTITU internal ve;
    VeHandler internal handler;

    function setUp() public {
        vm.warp(1_700_000_000);
        token = new MockERC20();
        ve = new VeTITU(IERC20(address(token)));
        handler = new VeHandler(ve, IERC20(address(token)));
        targetContract(address(handler));
    }

    /// @notice Balance views are uint; ensure they never glitch into dust above totalLocked.
    ///         Strongest property we enforce here: each user's `balanceOf` never exceeds
    ///         their locked amount.
    function invariant_veNonNegative() public view {
        // balances are uint256; this checks absence of under/overflow artifacts.
        for (uint256 i; i < 3; ++i) {
            address u = handler.users(i);
            (int128 amt,) = ve.locked(u);
            uint256 cap = amt >= 0 ? uint256(uint128(amt)) : 0;
            assertLe(ve.balanceOf(u), cap);
        }
    }
}
