// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {LPLock} from "../../../src/launchpad/LPLock.sol";

/// @dev Mintable LP-token mock used by the LpLock invariant suite.
contract MockLPToken is ERC20 {
    constructor() ERC20("LP", "LP") {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @notice Handler that exercises every path on {LPLock} except the time-warp.
///         Specifically it tries to:
///           * deposit from arbitrary actors (legal pre- and post-unlock),
///           * call {LPLock.withdraw} from the beneficiary AND from non-beneficiaries,
///         while time stays parked just shy of `unlockTime`. Every withdrawal
///         attempt is expected to revert; the invariant file asserts the
///         lock's storage never drifted from a post-deposit state.
contract LpLockHandler is Test {
    LPLock public immutable lock;
    MockLPToken public immutable lp;
    address public immutable beneficiary;
    address[] public actors;

    /// @notice Cumulative LP wei the handler has tried to deposit. The lock's
    ///         own counter must equal this exactly across the run.
    uint256 public totalDepositedByHandler;

    /// @notice Per-deposit cap. Bounded so we don't blow past uint96.
    uint256 public constant MAX_DEPOSIT = 100_000e18;

    constructor(LPLock _lock, MockLPToken _lp, address _beneficiary) {
        lock = _lock;
        lp = _lp;
        beneficiary = _beneficiary;
        actors.push(address(0xA11CE));
        actors.push(address(0xB0B));
        actors.push(address(0xCAFE));
        actors.push(_beneficiary);

        for (uint256 i = 0; i < actors.length; ++i) {
            lp.mint(actors[i], 1_000_000e18);
            vm.prank(actors[i]);
            lp.approve(address(lock), type(uint256).max);
        }
    }

    function _pick(uint256 i) internal view returns (address) {
        return actors[i % actors.length];
    }

    /// @dev Deposit a bounded random amount from a random actor.
    function deposit(uint256 actorIdx, uint96 amt) external {
        address a = _pick(actorIdx);
        uint256 amount = bound(uint256(amt), 1, MAX_DEPOSIT);
        if (lp.balanceOf(a) < amount) {
            lp.mint(a, amount);
        }

        totalDepositedByHandler += amount;
        vm.prank(a);
        lock.deposit(amount);
    }

    /// @dev Try to withdraw from the beneficiary while still locked. MUST
    ///      revert with {StillLocked}. Silently swallow the revert so the
    ///      fuzzer keeps moving — the invariant catches any state change.
    function tryWithdrawAsBeneficiary() external {
        vm.prank(beneficiary);
        try lock.withdraw() {
            revert("withdraw should have failed");
        } catch {}
    }

    /// @dev Try to withdraw from a non-beneficiary. MUST revert with
    ///      {NotBeneficiary}. Same swallow-pattern as above.
    function tryWithdrawAsAttacker(uint256 actorIdx) external {
        address a = _pick(actorIdx);
        if (a == beneficiary) return;
        vm.prank(a);
        try lock.withdraw() {
            revert("attacker should not withdraw");
        } catch {}
    }
}
