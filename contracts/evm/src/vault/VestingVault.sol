// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";

/// @title VestingVault
/// @notice Per-beneficiary cliff + linear vesting. One grant per address.
/// @dev Non-upgradeable. Admin role (DEFAULT_ADMIN_ROLE) holds grant + revoke power.
///      Tokens are pulled in via `safeTransferFrom` at grant time so the contract
///      always holds collateral for any outstanding grant (invariant enforced by tests).
contract VestingVault is AccessControl {
    using SafeERC20 for IERC20;

    struct Grant {
        uint256 total;
        uint256 released;
        uint64 start;
        uint64 cliff;
        uint64 duration;
    }

    /// @notice Token being vested.
    IERC20 public immutable TOKEN;

    /// @notice Grants keyed by beneficiary.
    mapping(address => Grant) public grants;

    /// @dev Emitted on grant creation.
    event GrantAdded(address indexed beneficiary, uint256 amount, uint64 start, uint64 cliff, uint64 duration);
    /// @dev Emitted on release to the beneficiary.
    event Released(address indexed beneficiary, uint256 amount);
    /// @dev Emitted on revoke. `unvestedReturned` is the amount swept back to admin.
    event Revoked(address indexed beneficiary, uint256 vestedKept, uint256 unvestedReturned);

    error ZeroAddress();
    error ZeroAmount();
    error ZeroDuration();
    error CliffExceedsDuration();
    error GrantExists();
    error NoGrant();
    error NothingToRelease();

    /// @param token IERC20 being vested.
    /// @param admin Address granted DEFAULT_ADMIN_ROLE.
    constructor(IERC20 token, address admin) {
        if (address(token) == address(0) || admin == address(0)) revert ZeroAddress();
        TOKEN = token;
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
    }

    /// @notice Create a grant. Pulls `amount` from caller (who must have approved).
    /// @dev Admin-only. One grant per beneficiary; re-grants after full release or revoke
    ///      are explicitly blocked to keep accounting trivial.
    function addGrant(address beneficiary, uint256 amount, uint64 start, uint64 cliff, uint64 duration)
        external
        onlyRole(DEFAULT_ADMIN_ROLE)
    {
        if (beneficiary == address(0)) revert ZeroAddress();
        if (amount == 0) revert ZeroAmount();
        if (duration == 0) revert ZeroDuration();
        if (cliff > duration) revert CliffExceedsDuration();
        if (grants[beneficiary].total != 0) revert GrantExists();

        grants[beneficiary] = Grant({total: amount, released: 0, start: start, cliff: cliff, duration: duration});
        TOKEN.safeTransferFrom(msg.sender, address(this), amount);
        emit GrantAdded(beneficiary, amount, start, cliff, duration);
    }

    /// @notice Release vested tokens for the caller.
    function release() external {
        _release(msg.sender);
    }

    /// @notice Release vested tokens for `beneficiary`. Anyone may trigger.
    function release(address beneficiary) external {
        _release(beneficiary);
    }

    /// @notice Revoke a grant. Any vested-but-unreleased portion is forwarded to the
    ///         beneficiary; the unvested remainder is returned to the caller (admin).
    function revoke(address beneficiary) external onlyRole(DEFAULT_ADMIN_ROLE) {
        Grant memory g = grants[beneficiary];
        if (g.total == 0) revert NoGrant();

        uint256 vestedNow = _vested(g, uint64(block.timestamp));
        uint256 toBeneficiary = vestedNow - g.released;
        uint256 toAdmin = g.total - vestedNow;

        delete grants[beneficiary];

        if (toBeneficiary > 0) TOKEN.safeTransfer(beneficiary, toBeneficiary);
        if (toAdmin > 0) TOKEN.safeTransfer(msg.sender, toAdmin);
        emit Revoked(beneficiary, toBeneficiary, toAdmin);
    }

    /// @notice View the total amount vested for `beneficiary` as of now.
    function vested(address beneficiary) external view returns (uint256) {
        return _vested(grants[beneficiary], uint64(block.timestamp));
    }

    /// @notice View the currently claimable amount for `beneficiary`.
    function releasable(address beneficiary) external view returns (uint256) {
        Grant memory g = grants[beneficiary];
        uint256 v = _vested(g, uint64(block.timestamp));
        return v - g.released;
    }

    // ---------------------------------------------------------------------
    // Internal
    // ---------------------------------------------------------------------

    /// @dev Strict-equality checks on `total` and `amount` are intentional existence tests.
    // slither-disable-next-line incorrect-equality
    function _release(address beneficiary) internal {
        Grant storage g = grants[beneficiary];
        if (g.total == 0) revert NoGrant();
        uint256 v = _vested(g, uint64(block.timestamp));
        uint256 amount = v - g.released;
        if (amount == 0) revert NothingToRelease();
        g.released = v;
        TOKEN.safeTransfer(beneficiary, amount);
        emit Released(beneficiary, amount);
    }

    /// @dev Strict equality on `g.total == 0` is the "no grant" sentinel.
    // slither-disable-next-line incorrect-equality
    function _vested(Grant memory g, uint64 nowTs) internal pure returns (uint256) {
        if (g.total == 0) return 0;
        if (nowTs < g.start + g.cliff) return 0;
        uint256 end = uint256(g.start) + uint256(g.duration);
        if (nowTs >= end) return g.total;
        uint256 elapsed = uint256(nowTs) - uint256(g.start);
        return (g.total * elapsed) / uint256(g.duration);
    }
}
