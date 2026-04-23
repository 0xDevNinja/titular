// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {StdInvariant} from "forge-std/StdInvariant.sol";
import {TITU} from "../../src/token/TITU.sol";
import {TituHandler} from "../handlers/TituHandler.sol";

contract TituInvariant is StdInvariant, Test {
    TITU internal titu;
    TituHandler internal handler;
    address internal seed = address(0xBEEF);

    function setUp() public {
        titu = new TITU(seed);
        handler = new TituHandler(titu, seed);
        targetContract(address(handler));
    }

    /// @notice Under non-burn operations the supply must equal the initial mint exactly.
    function invariant_totalSupplyFixed() public view {
        assertEq(titu.totalSupply(), 1_000_000_000 * 1e18);
    }
}
