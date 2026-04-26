// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {StdInvariant} from "forge-std/StdInvariant.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {LPLock} from "../../src/launchpad/LPLock.sol";
import {LpLockHandler, MockLPToken} from "./handlers/LpLockHandler.sol";

/// @title LpLockInvariant
/// @notice Forge invariant suite for {LPLock}. Drives randomised deposit and
///         withdraw-attempt traffic from a small actor pool while time stays
///         parked just below `unlockTime` (the handler never warps the clock).
///         Properties:
///           1. The beneficiary CANNOT withdraw before the unlock — the live
///              LP balance must equal the cumulative ingress at all times.
///           2. Non-beneficiaries cannot withdraw, ever.
///           3. `lockedAmount` (the cumulative ingress counter) is monotonic
///              non-decreasing, and equals the sum of all `deposit` amounts
///              the handler has ever issued.
///           4. The `withdrawn` flag stays `false` for the entire run.
///         The `unlockTime + 1` block warp is left to a unit test
///         ({LPLock.t.sol}); this suite is the pre-unlock safety net.
contract LpLockInvariant is StdInvariant, Test {
    LPLock internal lock;
    MockLPToken internal lp;
    LpLockHandler internal handler;

    address internal beneficiary = address(0xBEEF);
    uint256 internal unlockAt;

    function setUp() public {
        lp = new MockLPToken();
        lock = new LPLock(IERC20(address(lp)), beneficiary);
        unlockAt = lock.unlockTime();

        handler = new LpLockHandler(lock, lp, beneficiary);
        targetContract(address(handler));
    }

    /// @notice Pre-unlock the live LP balance MUST equal the cumulative ingress.
    ///         If a withdrawal had succeeded the balance would drop below
    ///         `lockedAmount`; if a deposit had been double-counted it would
    ///         drift above. Either case is a fatal break.
    function invariant_lpLock_balanceMatchesLockedAmount() public view {
        // Sanity: the handler does not warp time, so the lock must still be
        // engaged at every check.
        assertLt(block.timestamp, unlockAt);

        assertEq(lp.balanceOf(address(lock)), lock.lockedAmount());
    }

    /// @notice The handler tracks every deposit it issues; the lock's own
    ///         counter must agree exactly. Rules out hidden state mutation
    ///         by attempted withdrawals.
    function invariant_lpLock_lockedAmountEqualsHandlerCounter() public view {
        assertEq(lock.lockedAmount(), handler.totalDepositedByHandler());
    }

    /// @notice The lock has not flipped its one-shot withdrawal latch.
    function invariant_lpLock_notWithdrawn() public view {
        assertFalse(lock.withdrawn());
    }

    /// @notice The unlock timestamp is immutable for the run. Defence-in-depth
    ///         on the lock's "no extend, no early-unlock" promise.
    function invariant_lpLock_unlockTimeImmutable() public view {
        assertEq(lock.unlockTime(), unlockAt);
    }
}
