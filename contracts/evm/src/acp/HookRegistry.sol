// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {IHook} from "./IHook.sol";

/// @title HookRegistry
/// @notice Registry of approved IHook implementations.  Job contracts query this
///         registry to validate that a hook address has been whitelisted before
///         calling it, preventing arbitrary code injection via hook fields.
/// @dev    Role model:
///           - DEFAULT_ADMIN_ROLE: register and deregister hook implementations.
///
///         Hooks are identified by their contract address.  The registry stores:
///           - A boolean `approved` flag.
///           - The `hookName()` string (cached at registration time for gas and
///             to surface readable names off-chain without additional calls).
///
///         Any address may query `isApproved(hookAddr)`.
contract HookRegistry is AccessControl {
    // ─────────────────────────────────────────────────────────────
    // Types
    // ─────────────────────────────────────────────────────────────

    struct HookInfo {
        bool approved;
        string name;
    }

    // ─────────────────────────────────────────────────────────────
    // State
    // ─────────────────────────────────────────────────────────────

    /// @notice hook address => HookInfo.
    mapping(address => HookInfo) private _hooks;

    /// @notice Ordered list of all ever-registered hook addresses (may include revoked ones).
    address[] private _hookAddresses;

    // ─────────────────────────────────────────────────────────────
    // Events
    // ─────────────────────────────────────────────────────────────

    /// @notice Emitted when a hook is registered.
    event HookRegistered(address indexed hook, string name);

    /// @notice Emitted when a hook is deregistered (revoked).
    event HookDeregistered(address indexed hook);

    // ─────────────────────────────────────────────────────────────
    // Errors
    // ─────────────────────────────────────────────────────────────

    error ZeroAddress();
    error AlreadyRegistered(address hook);
    error NotRegistered(address hook);

    // ─────────────────────────────────────────────────────────────
    // Constructor
    // ─────────────────────────────────────────────────────────────

    /// @param admin Address granted DEFAULT_ADMIN_ROLE.
    constructor(address admin) {
        if (admin == address(0)) revert ZeroAddress();
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
    }

    // ─────────────────────────────────────────────────────────────
    // Registration
    // ─────────────────────────────────────────────────────────────

    /// @notice Register a hook implementation.  Calls `hookName()` on the hook to
    ///         cache the name.  Reverts if the hook is already registered.
    /// @param  hook Hook contract address (must implement IHook).
    function registerHook(address hook) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (hook == address(0)) revert ZeroAddress();
        if (_hooks[hook].approved) revert AlreadyRegistered(hook);

        string memory name = IHook(hook).hookName();
        _hooks[hook] = HookInfo({approved: true, name: name});
        _hookAddresses.push(hook);

        emit HookRegistered(hook, name);
    }

    /// @notice Deregister (revoke approval of) a hook.
    /// @param  hook Hook contract address.
    function deregisterHook(address hook) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (!_hooks[hook].approved) revert NotRegistered(hook);
        _hooks[hook].approved = false;
        emit HookDeregistered(hook);
    }

    // ─────────────────────────────────────────────────────────────
    // Views
    // ─────────────────────────────────────────────────────────────

    /// @notice Return true iff `hook` is currently approved.
    /// @param  hook Hook address to query.
    /// @return approved Whether the hook is approved.
    function isApproved(address hook) external view returns (bool approved) {
        approved = _hooks[hook].approved;
    }

    /// @notice Return the cached info for a hook.
    /// @param  hook Hook address.
    /// @return info HookInfo struct (approved + name).
    function getHookInfo(address hook) external view returns (HookInfo memory info) {
        info = _hooks[hook];
    }

    /// @notice Return all ever-registered hook addresses (including revoked ones).
    ///         Callers should filter by `isApproved` to get only live hooks.
    /// @return addresses Array of registered addresses.
    function allHooks() external view returns (address[] memory addresses) {
        addresses = _hookAddresses;
    }
}
