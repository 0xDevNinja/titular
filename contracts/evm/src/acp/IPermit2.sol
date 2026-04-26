// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

/// @title IPermit2
/// @notice Minimal Permit2 interface used by Escrow.fundWithPermit2.
interface IPermit2 {
    struct TokenPermissions {
        address token;
        uint256 amount;
    }

    struct PermitTransferFrom {
        TokenPermissions permitted;
        uint256 nonce;
        uint256 deadline;
    }

    struct SignatureTransferDetails {
        address to;
        uint256 requestedAmount;
    }
    /// @notice Transfer tokens using a signed Permit2 message.
    /// @param permit Permit data (token, amount, nonce, deadline).
    /// @param transferDetails Transfer destination and amount.
    /// @param owner The account whose tokens will be transferred.
    /// @param signature EIP-712 signature over the permit.
    function permitTransferFrom(
        PermitTransferFrom calldata permit,
        SignatureTransferDetails calldata transferDetails,
        address owner,
        bytes calldata signature
    ) external;
}
