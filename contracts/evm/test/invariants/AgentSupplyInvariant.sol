// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {StdInvariant} from "forge-std/StdInvariant.sol";
import {Clones} from "@openzeppelin/contracts/proxy/Clones.sol";
import {AgentToken} from "../../src/launchpad/AgentToken.sol";
import {AgentSupplyHandler} from "./handlers/AgentSupplyHandler.sol";

/// @title AgentSupplyInvariant
/// @notice Forge invariant suite for the {AgentToken} fixed supply.
///         {AgentToken} mints exactly `TOTAL_SUPPLY` (1B * 1e18) into the
///         bonding curve in {initialize}, and exposes no public mint or
///         burn surface. The internal `_update` allows the burn branch
///         (`to == address(0)`) but no external function ever routes there
///         — OZ's ERC-20 forbids transfer-to-zero. This invariant asserts:
///           1. `totalSupply()` equals the constant `TOTAL_SUPPLY` at
///              every reachable state — neither mint nor burn ever fires
///              after {initialize}.
///           2. The sum of every actor balance, plus the bonding curve,
///              plus the fee router (where transfer tax accrues), equals
///              `totalSupply()`. Conservation under all transfer paths.
///           3. {AgentToken.transfer} to `address(0)` ALWAYS reverts —
///              there is no exposed burn API. The handler's
///              `tryTransferToZero` proves this dynamically.
contract AgentSupplyInvariant is StdInvariant, Test {
    AgentToken internal impl;
    AgentToken internal token;
    AgentSupplyHandler internal handler;

    address internal creator = address(0xC0FFEE);
    address internal feeRouter = address(0xFEE);
    address internal bondingCurve = address(0xB0A2D);

    address[] internal actors;

    string internal constant NAME = "Agent Smith";
    string internal constant SYMBOL = "SMITH";

    function setUp() public {
        impl = new AgentToken();
        token = AgentToken(Clones.clone(address(impl)));
        token.initialize(NAME, SYMBOL, creator, feeRouter, bondingCurve);

        actors.push(address(0xA11CE));
        actors.push(address(0xB0B));
        actors.push(address(0xCAFE));
        actors.push(address(0xDEAD));

        handler = new AgentSupplyHandler(token, bondingCurve, feeRouter);
        targetContract(address(handler));
    }

    /// @notice `totalSupply()` is fixed at `TOTAL_SUPPLY` for the entire run.
    ///         Catches any future regression that adds an external mint or
    ///         burn path.
    function invariant_agentSupply_totalSupplyFixed() public view {
        assertEq(token.totalSupply(), token.TOTAL_SUPPLY());
    }

    /// @notice Conservation: the sum of every reachable holder's balance
    ///         equals `totalSupply()`. Holders are the four handler actors,
    ///         the bonding curve (initial mint recipient), and the fee router
    ///         (transfer-tax sink).
    function invariant_agentSupply_balancesSumToTotal() public view {
        uint256 sum = token.balanceOf(bondingCurve) + token.balanceOf(feeRouter);
        for (uint256 i = 0; i < actors.length; ++i) {
            sum += token.balanceOf(actors[i]);
        }
        assertEq(sum, token.totalSupply());
    }

    /// @notice The bonding curve never receives more than its initial
    ///         allocation. Catches any rogue mint that targets the curve.
    function invariant_agentSupply_bondingCurveNeverGrows() public view {
        assertLe(token.balanceOf(bondingCurve), token.TOTAL_SUPPLY());
    }
}
