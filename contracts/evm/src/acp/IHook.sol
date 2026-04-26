// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

/// @title IHook
/// @notice Interface that all ACP v2 hooks must implement.
/// @dev    Hooks are called by Job contracts at specific lifecycle points.
///         A hook may revert to block the transition (blocking hook) or succeed
///         to allow it (permissive hook).  The Job contract is responsible for
///         deciding which hooks to invoke and in what order.
///
///         Hook call context is passed as an abi-encoded `bytes` blob to allow
///         each hook to decode exactly what it needs without requiring the Job
///         contract to know the hook's internal ABI.
interface IHook {
    // ─────────────────────────────────────────────────────────────
    // Lifecycle points
    // ─────────────────────────────────────────────────────────────

    /// @notice Called when an agent accepts a job. Transitions Open → Active.
    /// @param jobId   The job identifier.
    /// @param context ABI-encoded context specific to this hook type.
    function onAccept(uint256 jobId, bytes calldata context) external;

    /// @notice Called when an agent submits a result. Transitions Active → Review.
    /// @param jobId   The job identifier.
    /// @param context ABI-encoded context specific to this hook type.
    function onSubmit(uint256 jobId, bytes calldata context) external;

    /// @notice Called when a result is approved (by evaluator or principal).
    ///         Transitions Review → Completed.
    /// @param jobId   The job identifier.
    /// @param context ABI-encoded context specific to this hook type.
    function onApprove(uint256 jobId, bytes calldata context) external;

    /// @notice Called when a result is rejected by the evaluator.
    ///         Transitions Review → Active (agent may resubmit).
    /// @param jobId   The job identifier.
    /// @param context ABI-encoded context specific to this hook type.
    function onReject(uint256 jobId, bytes calldata context) external;

    /// @notice Called when a job is cancelled by the principal or via expiry.
    /// @param jobId   The job identifier.
    /// @param context ABI-encoded context specific to this hook type.
    function onCancel(uint256 jobId, bytes calldata context) external;

    /// @notice Returns a human-readable name for this hook (used off-chain for UX).
    /// @return name Hook name, e.g. "FundTransferHook".
    function hookName() external view returns (string memory name);
}
