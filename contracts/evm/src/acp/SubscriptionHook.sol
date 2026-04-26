// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {IHook} from "./IHook.sol";

/// @title SubscriptionHook
/// @notice Hook that implements recurring subscription renewal logic.
///
///         Tracks subscription state per jobId. On each `onApprove`, if the
///         subscription is still active and within the renewal window, it
///         debits the subscriber and credits the provider.
///
///         Context for `onApprove`:
///           abi.encode(address subscriber, address provider, address token,
///                      uint256 periodAmount, uint256 periodDuration)
///
///         Subscription is "active" while block.timestamp < nextRenewalAt.
///         The hook does NOT initiate payments on its own; it acts when `onApprove`
///         is called by the Job contract at the end of each billing period.
contract SubscriptionHook is IHook {
    using SafeERC20 for IERC20;

    // ─────────────────────────────────────────────────────────────
    // Types
    // ─────────────────────────────────────────────────────────────

    struct SubscriptionState {
        uint64 nextRenewalAt;
        bool active;
    }

    struct RenewContext {
        address subscriber;
        address provider;
        address token;
        uint256 periodAmount;
        uint64 periodDuration;
    }

    // ─────────────────────────────────────────────────────────────
    // State
    // ─────────────────────────────────────────────────────────────

    /// @notice jobId => subscription state.
    mapping(uint256 => SubscriptionState) public subscriptions;

    // ─────────────────────────────────────────────────────────────
    // Events
    // ─────────────────────────────────────────────────────────────

    /// @notice Emitted when a subscription is renewed.
    event SubscriptionRenewed(
        uint256 indexed jobId,
        address indexed subscriber,
        address indexed provider,
        uint256 amount,
        uint64 nextRenewalAt
    );

    /// @notice Emitted when a subscription is cancelled via `onCancel`.
    event SubscriptionCancelled(uint256 indexed jobId);

    // ─────────────────────────────────────────────────────────────
    // Errors
    // ─────────────────────────────────────────────────────────────

    error ZeroAddress();
    error ZeroAmount();
    error InvalidDuration();
    error SubscriptionNotActive(uint256 jobId);

    // ─────────────────────────────────────────────────────────────
    // IHook implementation
    // ─────────────────────────────────────────────────────────────

    /// @inheritdoc IHook
    function onAccept(uint256 jobId, bytes calldata context) external override {
        // Initialise subscription state on first accept
        RenewContext memory ctx = abi.decode(context, (RenewContext));
        if (ctx.periodDuration == 0) revert InvalidDuration();
        // slither-disable-next-line timestamp
        subscriptions[jobId] = SubscriptionState({
            nextRenewalAt: uint64(block.timestamp) + ctx.periodDuration,
            active: true
        });
    }

    /// @inheritdoc IHook
    function onSubmit(uint256, bytes calldata) external pure override {}

    /// @notice Renew the subscription: debit subscriber, credit provider.
    /// @param jobId   The job identifier.
    /// @param context ABI-encoded RenewContext.
    function onApprove(uint256 jobId, bytes calldata context) external override {
        SubscriptionState storage sub = subscriptions[jobId];
        if (!sub.active) revert SubscriptionNotActive(jobId);

        RenewContext memory ctx = abi.decode(context, (RenewContext));
        if (ctx.subscriber == address(0)) revert ZeroAddress();
        if (ctx.provider == address(0)) revert ZeroAddress();
        if (ctx.token == address(0)) revert ZeroAddress();
        if (ctx.periodAmount == 0) revert ZeroAmount();
        if (ctx.periodDuration == 0) revert InvalidDuration();

        // Effects before interaction
        uint64 next = uint64(block.timestamp) + ctx.periodDuration;
        sub.nextRenewalAt = next;

        emit SubscriptionRenewed(jobId, ctx.subscriber, ctx.provider, ctx.periodAmount, next);

        // Interaction — ctx.subscriber is the pre-approved payer; only the Job contract
        // (a trusted caller registered in HookRegistry) supplies this context.
        // slither-disable-next-line arbitrary-send-erc20
        IERC20(ctx.token).safeTransferFrom(ctx.subscriber, ctx.provider, ctx.periodAmount);
    }

    /// @inheritdoc IHook
    function onReject(uint256, bytes calldata) external pure override {}

    /// @notice Cancel the subscription when the job is cancelled.
    /// @param jobId The job identifier.
    function onCancel(uint256 jobId, bytes calldata) external override {
        if (subscriptions[jobId].active) {
            subscriptions[jobId].active = false;
            emit SubscriptionCancelled(jobId);
        }
    }

    /// @inheritdoc IHook
    function hookName() external pure override returns (string memory) {
        return "SubscriptionHook";
    }
}
