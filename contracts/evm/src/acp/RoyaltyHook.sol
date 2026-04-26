// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {IHook} from "./IHook.sol";

/// @title RoyaltyHook
/// @notice Distributes downstream revenue share to a list of royalty recipients
///         whenever a job result is approved.
///
///         Context for `onApprove`:
///           abi.encode(address token, uint256 amount, RoyaltyRecipient[] recipients)
///
///         The sum of all `recipients[i].shareBps` must equal 10 000 (100%).
///         Any dust from integer division accumulates to the first recipient.
contract RoyaltyHook is IHook {
    using SafeERC20 for IERC20;

    // ─────────────────────────────────────────────────────────────
    // Constants
    // ─────────────────────────────────────────────────────────────

    uint256 public constant BPS = 10_000;

    // ─────────────────────────────────────────────────────────────
    // Types
    // ─────────────────────────────────────────────────────────────

    struct RoyaltyRecipient {
        address recipient;
        uint256 shareBps; // Share in BPS (must sum to BPS across all recipients)
    }

    struct RoyaltyContext {
        address token;
        uint256 amount;
        RoyaltyRecipient[] recipients;
    }

    // ─────────────────────────────────────────────────────────────
    // Events
    // ─────────────────────────────────────────────────────────────

    /// @notice Emitted for each royalty payment.
    event RoyaltyPaid(
        uint256 indexed jobId,
        address indexed token,
        address indexed recipient,
        uint256 amount
    );

    // ─────────────────────────────────────────────────────────────
    // Errors
    // ─────────────────────────────────────────────────────────────

    error ZeroAddress();
    error ZeroAmount();
    error NoRecipients();
    error ShareSumMismatch(uint256 sum);
    error RecipientZeroAddress(uint256 index);

    // ─────────────────────────────────────────────────────────────
    // IHook implementation
    // ─────────────────────────────────────────────────────────────

    /// @inheritdoc IHook
    function onAccept(uint256, bytes calldata) external pure override {}

    /// @inheritdoc IHook
    function onSubmit(uint256, bytes calldata) external pure override {}

    /// @notice Distribute royalties to all recipients.
    /// @param jobId   The job identifier.
    /// @param context ABI-encoded RoyaltyContext.
    function onApprove(uint256 jobId, bytes calldata context) external override {
        RoyaltyContext memory ctx = abi.decode(context, (RoyaltyContext));
        if (ctx.token == address(0)) revert ZeroAddress();
        if (ctx.amount == 0) revert ZeroAmount();
        if (ctx.recipients.length == 0) revert NoRecipients();

        // Validate recipients and sum
        uint256 sum = 0;
        for (uint256 i = 0; i < ctx.recipients.length; ++i) {
            if (ctx.recipients[i].recipient == address(0)) revert RecipientZeroAddress(i);
            sum += ctx.recipients[i].shareBps;
        }
        if (sum != BPS) revert ShareSumMismatch(sum);

        // Distribute — first recipient accumulates dust
        uint256 distributed = 0;
        for (uint256 i = 0; i < ctx.recipients.length; ++i) {
            uint256 payout;
            if (i == ctx.recipients.length - 1) {
                // Last recipient gets remainder (dust)
                payout = ctx.amount - distributed;
            } else {
                payout = (ctx.amount * ctx.recipients[i].shareBps) / BPS;
            }

            if (payout > 0) {
                distributed += payout;
                emit RoyaltyPaid(jobId, ctx.token, ctx.recipients[i].recipient, payout);
                IERC20(ctx.token).safeTransferFrom(msg.sender, ctx.recipients[i].recipient, payout);
            }
        }
    }

    /// @inheritdoc IHook
    function onReject(uint256, bytes calldata) external pure override {}

    /// @inheritdoc IHook
    function onCancel(uint256, bytes calldata) external pure override {}

    /// @inheritdoc IHook
    function hookName() external pure override returns (string memory) {
        return "RoyaltyHook";
    }
}
