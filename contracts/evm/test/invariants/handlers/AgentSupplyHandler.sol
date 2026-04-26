// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {AgentToken} from "../../../src/launchpad/AgentToken.sol";

/// @notice Handler that exercises every public path of {AgentToken} that
///         could conceivably move supply, plus a few that should not. The
///         token's only mint site is {initialize} (one-shot, fired in
///         setUp), and there is no external `burn` exposed — the
///         `to == address(0)` branch in {AgentToken._update} is internal
///         only. The handler therefore drives `transfer` / `transferFrom`
///         from a small actor pool, plus an explicit `tryTransferToZero`
///         that should ALWAYS revert (ERC-20 forbids transfer-to-zero).
///         Any drift in `totalSupply()` over the run is fatal.
contract AgentSupplyHandler is Test {
    AgentToken public immutable token;
    address public immutable bondingCurve;
    address public immutable feeRouter;
    address[] public actors;

    /// @notice Cap per transfer so cumulative drains can't exhaust the
    ///         bonding-curve seed faster than the run depth allows.
    uint256 public constant MAX_TRANSFER = 1_000_000e18;

    constructor(AgentToken _token, address _bondingCurve, address _feeRouter) {
        token = _token;
        bondingCurve = _bondingCurve;
        feeRouter = _feeRouter;
        actors.push(address(0xA11CE));
        actors.push(address(0xB0B));
        actors.push(address(0xCAFE));
        actors.push(address(0xDEAD));

        // Approve a future-proof allowance for transferFrom traffic.
        for (uint256 i = 0; i < actors.length; ++i) {
            for (uint256 j = 0; j < actors.length; ++j) {
                if (i == j) continue;
                vm.prank(actors[i]);
                token.approve(actors[j], type(uint256).max);
            }
        }
    }

    function _pick(uint256 i) internal view returns (address) {
        return actors[i % actors.length];
    }

    /// @dev Seed flow — pull from the bonding curve to an actor. This is
    ///      where actor balances first come from. The curve is pranked
    ///      directly because it is a stand-in EOA, not a real curve.
    function pullFromCurve(uint256 actorIdx, uint96 amt) external {
        address to = _pick(actorIdx);
        uint256 bal = token.balanceOf(bondingCurve);
        if (bal == 0) return;
        uint256 amount = bound(uint256(amt), 1, bal > MAX_TRANSFER ? MAX_TRANSFER : bal);
        if (amount == 0) return;
        vm.prank(bondingCurve);
        token.transfer(to, amount);
    }

    /// @dev Peer-to-peer transfer. Hits the taxed branch in
    ///      {AgentToken._update}; supply still must NOT change because
    ///      the tax is a transfer to `feeRouter`, never a burn.
    function transfer(uint256 fromIdx, uint256 toIdx, uint96 amt) external {
        address from = _pick(fromIdx);
        address to = _pick(toIdx);
        if (from == to) return;
        uint256 bal = token.balanceOf(from);
        if (bal == 0) return;
        uint256 amount = bound(uint256(amt), 1, bal);
        if (amount == 0) return;
        vm.prank(from);
        token.transfer(to, amount);
    }

    /// @dev `transferFrom` flow — exercises the allowance-burn path of OZ
    ///      ERC-20 alongside the tax routing.
    function transferFrom(uint256 fromIdx, uint256 spenderIdx, uint256 toIdx, uint96 amt) external {
        address from = _pick(fromIdx);
        address spender = _pick(spenderIdx);
        address to = _pick(toIdx);
        if (from == to || from == spender) return;
        uint256 bal = token.balanceOf(from);
        if (bal == 0) return;
        uint256 amount = bound(uint256(amt), 1, bal);
        if (amount == 0) return;
        vm.prank(spender);
        token.transferFrom(from, to, amount);
    }

    /// @dev OZ ERC-20 forbids transfer-to-zero; this MUST revert. Used by
    ///      the invariant suite to confirm there is no exposed burn API.
    function tryTransferToZero(uint256 fromIdx, uint96 amt) external {
        address from = _pick(fromIdx);
        uint256 bal = token.balanceOf(from);
        if (bal == 0) return;
        uint256 amount = bound(uint256(amt), 1, bal);
        vm.prank(from);
        try token.transfer(address(0), amount) {
            revert("transfer to zero should revert");
        } catch {}
    }
}
