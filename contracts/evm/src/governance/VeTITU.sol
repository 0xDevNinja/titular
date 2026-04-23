// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/// @title VeTITU
/// @notice Curve-style vote-escrowed TITU. Locks decay linearly to zero at unlock.
/// @dev Simplified checkpoint model:
///        - Per-user point history recorded on every mutation.
///        - Global history recorded on every mutation.
///        - `getPastVotes(user, t)` evaluates the user's most recent point at-or-before `t`
///          and applies linear decay to `t`. Balance is clamped at zero post-unlock.
///        - `totalSupply()` uses the most recent global point and applies linear decay; we
///          do NOT replay expiring locks across weeks here. For v1 this is documented as
///          a known approximation that is conservative immediately after any state change
///          and is corrected on the next mutation. Full per-week `slope_changes` replay is
///          deferred to M10 audit-hardening.
///      ERC-5805 compat: `clock()` in timestamp mode, `getVotes`, `getPastVotes`.
contract VeTITU is ReentrancyGuard {
    using SafeERC20 for IERC20;

    // ---------------------------------------------------------------------
    // Constants
    // ---------------------------------------------------------------------

    /// @notice One week in seconds; lock ends snap to this grid.
    uint256 public constant WEEK = 7 days;
    /// @notice Maximum lock duration (4 years).
    uint256 public constant MAXTIME = 4 * 365 days;

    /// @notice TITU token being locked.
    IERC20 public immutable TOKEN;

    enum DepositType {
        CREATE_LOCK,
        INCREASE_AMOUNT,
        INCREASE_UNLOCK_TIME
    }

    struct LockedBalance {
        int128 amount;
        uint256 end;
    }

    struct Point {
        int128 bias; // voting power at `ts`
        int128 slope; // rate of decay per second (amount / MAXTIME)
        uint256 ts; // timestamp
    }

    // ---------------------------------------------------------------------
    // Storage
    // ---------------------------------------------------------------------

    /// @notice Current lock per user.
    mapping(address => LockedBalance) public locked;

    /// @notice Per-user point history.
    mapping(address => Point[]) internal _userPoints;

    /// @notice Global point history.
    Point[] internal _globalPoints;

    /// @notice Total TITU locked across all users.
    uint256 public totalLocked;

    // ---------------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------------

    event Deposit(
        address indexed provider, uint256 value, uint256 indexed locktime, DepositType depositType, uint256 ts
    );
    event Withdraw(address indexed provider, uint256 value, uint256 ts);
    event Supply(uint256 prevSupply, uint256 supply);

    // ---------------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------------

    error ZeroAddress();
    error ZeroAmount();
    error LockExists();
    error NoLock();
    error UnlockTimeNotInFuture();
    error UnlockTimeTooShort();
    error UnlockTimeTooLong();
    error UnlockTimeNotLater();
    error LockExpired();
    error LockNotExpired();
    error NoTransfersAllowed();

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    constructor(IERC20 token) {
        if (address(token) == address(0)) revert ZeroAddress();
        TOKEN = token;
        // Seed a zero-state global point at deploy for cleaner historical queries.
        _globalPoints.push(Point({bias: 0, slope: 0, ts: block.timestamp}));
    }

    // ---------------------------------------------------------------------
    // External mutations
    // ---------------------------------------------------------------------

    /// @notice Create a new lock for `amount` TITU ending at `unlockTime` (rounded down to week).
    function createLock(uint256 amount, uint256 unlockTime) external nonReentrant {
        LockedBalance memory old = locked[msg.sender];
        if (old.amount != 0) revert LockExists();

        // Intentional floor-to-week; divide-before-multiply is the exact pattern we need.
        // slither-disable-next-line divide-before-multiply
        uint256 end = (unlockTime / WEEK) * WEEK;
        if (end <= block.timestamp) revert UnlockTimeNotInFuture();
        if (end < block.timestamp + WEEK) revert UnlockTimeTooShort();
        if (end > block.timestamp + MAXTIME) revert UnlockTimeTooLong();
        if (amount == 0) revert ZeroAmount();

        LockedBalance memory newLock = LockedBalance({amount: _toI128(amount), end: end});
        _checkpoint(msg.sender, old, newLock);
        locked[msg.sender] = newLock;

        uint256 prev = totalLocked;
        totalLocked = prev + amount;
        emit Supply(prev, totalLocked);
        emit Deposit(msg.sender, amount, end, DepositType.CREATE_LOCK, block.timestamp);

        TOKEN.safeTransferFrom(msg.sender, address(this), amount);
    }

    /// @notice Add `amount` to the caller's existing lock; unlock time unchanged.
    function increaseAmount(uint256 amount) external nonReentrant {
        LockedBalance memory old = locked[msg.sender];
        if (old.amount == 0) revert NoLock();
        if (old.end <= block.timestamp) revert LockExpired();
        if (amount == 0) revert ZeroAmount();

        LockedBalance memory newLock = LockedBalance({amount: old.amount + _toI128(amount), end: old.end});
        _checkpoint(msg.sender, old, newLock);
        locked[msg.sender] = newLock;

        uint256 prev = totalLocked;
        totalLocked = prev + amount;
        emit Supply(prev, totalLocked);
        emit Deposit(msg.sender, amount, old.end, DepositType.INCREASE_AMOUNT, block.timestamp);

        TOKEN.safeTransferFrom(msg.sender, address(this), amount);
    }

    /// @notice Extend the caller's unlock time (rounds down to week, must be strictly later).
    function increaseUnlockTime(uint256 newUnlockTime) external nonReentrant {
        LockedBalance memory old = locked[msg.sender];
        if (old.amount == 0) revert NoLock();
        if (old.end <= block.timestamp) revert LockExpired();

        // Intentional floor-to-week; divide-before-multiply is the exact pattern we need.
        // slither-disable-next-line divide-before-multiply
        uint256 end = (newUnlockTime / WEEK) * WEEK;
        if (end <= old.end) revert UnlockTimeNotLater();
        if (end > block.timestamp + MAXTIME) revert UnlockTimeTooLong();

        LockedBalance memory newLock = LockedBalance({amount: old.amount, end: end});
        _checkpoint(msg.sender, old, newLock);
        locked[msg.sender] = newLock;

        emit Deposit(msg.sender, 0, end, DepositType.INCREASE_UNLOCK_TIME, block.timestamp);
    }

    /// @notice Withdraw full locked amount after expiry.
    function withdraw() external nonReentrant {
        LockedBalance memory old = locked[msg.sender];
        if (old.amount == 0) revert NoLock();
        if (block.timestamp < old.end) revert LockNotExpired();

        uint256 amount = uint256(uint128(old.amount));
        LockedBalance memory empty = LockedBalance({amount: 0, end: 0});
        _checkpoint(msg.sender, old, empty);
        delete locked[msg.sender];

        uint256 prev = totalLocked;
        totalLocked = prev - amount;
        emit Supply(prev, totalLocked);
        emit Withdraw(msg.sender, amount, block.timestamp);

        TOKEN.safeTransfer(msg.sender, amount);
    }

    // ---------------------------------------------------------------------
    // Views — ERC-5805 + balances
    // ---------------------------------------------------------------------

    /// @notice Voting power of `user` at current time (linear decay, clamped at 0).
    function balanceOf(address user) public view returns (uint256) {
        return _balanceAt(user, block.timestamp);
    }

    /// @notice ERC-5805 alias of `balanceOf`.
    function getVotes(address user) external view returns (uint256) {
        return balanceOf(user);
    }

    /// @notice Voting power of `user` at historical `timepoint`.
    /// @dev Reverts if `timepoint >= block.timestamp` to match ERC-5805 semantics.
    function getPastVotes(address user, uint256 timepoint) external view returns (uint256) {
        require(timepoint < block.timestamp, "VeTITU: future lookup");
        return _balanceAt(user, timepoint);
    }

    /// @notice Total voting supply at current time (global decay approximation).
    function totalSupply() public view returns (uint256) {
        return _totalAt(block.timestamp);
    }

    /// @notice Total voting supply at historical `timepoint`.
    function getPastTotalSupply(uint256 timepoint) external view returns (uint256) {
        require(timepoint < block.timestamp, "VeTITU: future lookup");
        return _totalAt(timepoint);
    }

    /// @notice ERC-6372 clock (timestamp mode).
    function clock() external view returns (uint48) {
        return uint48(block.timestamp);
    }

    // solhint-disable-next-line func-name-mixedcase
    function CLOCK_MODE() external pure returns (string memory) {
        return "mode=timestamp";
    }

    /// @notice Number of user checkpoints.
    function userPointCount(address user) external view returns (uint256) {
        return _userPoints[user].length;
    }

    /// @notice Return the raw user point at index `i`.
    function userPointAt(address user, uint256 i) external view returns (Point memory) {
        return _userPoints[user][i];
    }

    /// @notice Number of global checkpoints.
    function globalPointCount() external view returns (uint256) {
        return _globalPoints.length;
    }

    // ---------------------------------------------------------------------
    // Disallow transfer-like semantics
    // ---------------------------------------------------------------------

    /// @notice ve balance is non-transferable; always reverts.
    function transfer(address, uint256) external pure returns (bool) {
        revert NoTransfersAllowed();
    }

    /// @notice ve balance is non-transferable; always reverts.
    function transferFrom(address, address, uint256) external pure returns (bool) {
        revert NoTransfersAllowed();
    }

    // ---------------------------------------------------------------------
    // Internal
    // ---------------------------------------------------------------------

    /// @dev slope = amount / MAXTIME is the Curve convention; the precision loss is the
    ///      security model (locking dust is pointless). Zero-init locals are intentional.
    // slither-disable-start divide-before-multiply
    // slither-disable-start uninitialized-local
    // slither-disable-start incorrect-equality
    function _checkpoint(address user, LockedBalance memory old, LockedBalance memory newLock) internal {
        int128 oldSlope;
        int128 oldBias;
        if (old.end > block.timestamp && old.amount > 0) {
            oldSlope = old.amount / int128(int256(MAXTIME));
            oldBias = oldSlope * int128(int256(old.end - block.timestamp));
        }
        int128 newSlope;
        int128 newBias;
        if (newLock.end > block.timestamp && newLock.amount > 0) {
            newSlope = newLock.amount / int128(int256(MAXTIME));
            newBias = newSlope * int128(int256(newLock.end - block.timestamp));
        }

        // User point.
        _userPoints[user].push(Point({bias: newBias, slope: newSlope, ts: block.timestamp}));

        // Global delta.
        Point memory last = _globalPoints.length == 0
            ? Point({bias: 0, slope: 0, ts: block.timestamp})
            : _globalPoints[_globalPoints.length - 1];
        // Decay the last global bias forward to now before applying deltas.
        int128 decayed = last.bias - last.slope * int128(int256(block.timestamp - last.ts));
        if (decayed < 0) decayed = 0;
        int128 newGlobalBias = decayed - oldBias + newBias;
        if (newGlobalBias < 0) newGlobalBias = 0;
        int128 newGlobalSlope = last.slope - oldSlope + newSlope;
        if (newGlobalSlope < 0) newGlobalSlope = 0;

        if (_globalPoints.length > 0 && last.ts == block.timestamp) {
            // Same-block update: replace the tail.
            _globalPoints[_globalPoints.length - 1] =
                Point({bias: newGlobalBias, slope: newGlobalSlope, ts: block.timestamp});
        } else {
            _globalPoints.push(Point({bias: newGlobalBias, slope: newGlobalSlope, ts: block.timestamp}));
        }
    }
    // slither-disable-end divide-before-multiply
    // slither-disable-end uninitialized-local
    // slither-disable-end incorrect-equality

    function _balanceAt(address user, uint256 ts) internal view returns (uint256) {
        Point memory p = _latestUserPointAt(user, ts);
        if (p.ts == 0 || p.bias == 0) return 0;
        int128 decayed = p.bias - p.slope * int128(int256(ts - p.ts));
        if (decayed < 0) return 0;
        // Additionally clamp: if the lock ended before `ts`, voting power is zero.
        LockedBalance memory lb = locked[user];
        if (lb.end != 0 && lb.end <= ts && lb.amount == _latestUserLockAmount(user, ts)) {
            // lock matches the point era and has expired relative to `ts`.
            return 0;
        }
        return uint256(uint128(decayed));
    }

    /// @dev Binary search; zero-init `lo` is intentional.
    function _latestUserPointAt(address user, uint256 ts) internal view returns (Point memory) {
        Point[] storage pts = _userPoints[user];
        uint256 n = pts.length;
        if (n == 0) return Point({bias: 0, slope: 0, ts: 0});
        // slither-disable-next-line uninitialized-local
        uint256 lo;
        uint256 hi = n;
        while (lo < hi) {
            uint256 mid = (lo + hi) >> 1;
            if (pts[mid].ts <= ts) lo = mid + 1;
            else hi = mid;
        }
        if (lo == 0) return Point({bias: 0, slope: 0, ts: 0});
        return pts[lo - 1];
    }

    function _latestUserLockAmount(address, uint256) internal pure returns (int128) {
        // Placeholder kept for symmetry with curve impl; we only clamp by lock.end elsewhere.
        return 0;
    }

    /// @dev Binary search; zero-init `lo` is intentional; n==0 empty-list check is not equality risk.
    function _totalAt(uint256 ts) internal view returns (uint256) {
        uint256 n = _globalPoints.length;
        // slither-disable-next-line incorrect-equality
        if (n == 0) return 0;
        // slither-disable-next-line uninitialized-local
        uint256 lo;
        uint256 hi = n;
        while (lo < hi) {
            uint256 mid = (lo + hi) >> 1;
            if (_globalPoints[mid].ts <= ts) lo = mid + 1;
            else hi = mid;
        }
        if (lo == 0) return 0;
        Point memory p = _globalPoints[lo - 1];
        int128 decayed = p.bias - p.slope * int128(int256(ts - p.ts));
        if (decayed < 0) return 0;
        return uint256(uint128(decayed));
    }

    function _toI128(uint256 x) internal pure returns (int128) {
        require(x <= uint256(uint128(type(int128).max)), "VeTITU: i128 overflow");
        return int128(uint128(x));
    }
}
