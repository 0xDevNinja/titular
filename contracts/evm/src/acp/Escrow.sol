// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

/// @title Escrow
/// @notice Multi-token escrow with atomic multi-recipient release.
/// @dev    Design decisions:
///           - Each escrow entry is identified by a (depositor, jobId) pair.  The jobId
///             is an opaque uint256 supplied by the caller (typically the JobFactory
///             counter); it need not be unique across depositors.
///           - Funds are held in this contract; individual balances are tracked per
///             (depositor, jobId, token) to allow multiple tokens per job.
///           - `release` atomically transfers to N recipients in a single call.
///           - `refund` returns all tokens to the depositor (dispute resolved in
///             depositor's favour, or admin override).
///           - RELEASER_ROLE: addresses allowed to call `release` and `refund`
///             (typically the Job contract or an arbiter).
///           - ADMIN_ROLE: can grant/revoke roles, emergency-pause.
///
///         CEI: all storage writes happen before external SafeERC20 transfers.
///         ReentrancyGuard on `fund`, `release`, `refund`.
contract Escrow is AccessControl, ReentrancyGuard {
    using SafeERC20 for IERC20;

    // ─────────────────────────────────────────────────────────────
    // Roles
    // ─────────────────────────────────────────────────────────────

    bytes32 public constant RELEASER_ROLE = keccak256("RELEASER_ROLE");

    // ─────────────────────────────────────────────────────────────
    // Types
    // ─────────────────────────────────────────────────────────────

    /// @notice A single recipient in a multi-recipient release.
    struct Recipient {
        address to;
        uint256 amount;
    }

    // ─────────────────────────────────────────────────────────────
    // Storage
    // ─────────────────────────────────────────────────────────────

    /// @notice depositor => jobId => token => escrowed balance.
    mapping(address => mapping(uint256 => mapping(address => uint256))) public balances;

    // ─────────────────────────────────────────────────────────────
    // Events
    // ─────────────────────────────────────────────────────────────

    /// @notice Emitted when tokens are deposited into escrow.
    event Funded(address indexed depositor, uint256 indexed jobId, address indexed token, uint256 amount);

    /// @notice Emitted for each recipient in an atomic release.
    event Released(
        address indexed depositor, uint256 indexed jobId, address indexed token, address to, uint256 amount
    );

    /// @notice Emitted when the full balance is refunded to the depositor.
    event Refunded(address indexed depositor, uint256 indexed jobId, address indexed token, uint256 amount);

    // ─────────────────────────────────────────────────────────────
    // Errors
    // ─────────────────────────────────────────────────────────────

    error ZeroAddress();
    error ZeroAmount();
    error InsufficientBalance(uint256 have, uint256 want);
    error EmptyRecipients();
    error RecipientZeroAddress(uint256 index);
    error RecipientZeroAmount(uint256 index);
    error SumExceedsBalance(uint256 sum, uint256 balance);

    // ─────────────────────────────────────────────────────────────
    // Constructor
    // ─────────────────────────────────────────────────────────────

    /// @param admin Address granted DEFAULT_ADMIN_ROLE and RELEASER_ROLE on deployment.
    constructor(address admin) {
        if (admin == address(0)) revert ZeroAddress();
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(RELEASER_ROLE, admin);
    }

    // ─────────────────────────────────────────────────────────────
    // Fund
    // ─────────────────────────────────────────────────────────────

    /// @notice Deposit `amount` of `token` into escrow for a specific job.
    /// @dev    Caller must have approved this contract for at least `amount` of `token`.
    /// @param  depositor Address whose balance increases (usually msg.sender; can be different
    ///                   for forwarded deposits, e.g. from JobFactory on behalf of principal).
    /// @param  jobId     Opaque identifier tying this deposit to a specific job.
    /// @param  token     ERC-20 token address.
    /// @param  amount    Amount to deposit (must be > 0).
    function fund(address depositor, uint256 jobId, address token, uint256 amount) external nonReentrant {
        if (depositor == address(0)) revert ZeroAddress();
        if (token == address(0)) revert ZeroAddress();
        if (amount == 0) revert ZeroAmount();

        // Effects
        balances[depositor][jobId][token] += amount;

        emit Funded(depositor, jobId, token, amount);

        // Interaction (after state update)
        IERC20(token).safeTransferFrom(msg.sender, address(this), amount);
    }

    // ─────────────────────────────────────────────────────────────
    // Release
    // ─────────────────────────────────────────────────────────────

    /// @notice Atomically release funds to multiple recipients in a single call.
    /// @dev    The sum of all `recipient.amount` values must be <= the escrowed balance.
    ///         Callers must hold RELEASER_ROLE.
    /// @param  depositor  Owner of the escrowed balance.
    /// @param  jobId      Job identifier.
    /// @param  token      ERC-20 token address.
    /// @param  recipients Array of (address, amount) pairs.  Must be non-empty.
    function release(address depositor, uint256 jobId, address token, Recipient[] calldata recipients)
        external
        nonReentrant
        onlyRole(RELEASER_ROLE)
    {
        if (recipients.length == 0) revert EmptyRecipients();

        uint256 total = 0;
        for (uint256 i = 0; i < recipients.length; ++i) {
            if (recipients[i].to == address(0)) revert RecipientZeroAddress(i);
            if (recipients[i].amount == 0) revert RecipientZeroAmount(i);
            total += recipients[i].amount;
        }

        uint256 bal = balances[depositor][jobId][token];
        if (total > bal) revert SumExceedsBalance(total, bal);

        // Effects — deduct full sum atomically before any transfer
        balances[depositor][jobId][token] = bal - total;

        // Interactions
        for (uint256 i = 0; i < recipients.length; ++i) {
            emit Released(depositor, jobId, token, recipients[i].to, recipients[i].amount);
            IERC20(token).safeTransfer(recipients[i].to, recipients[i].amount);
        }
    }

    // ─────────────────────────────────────────────────────────────
    // Refund
    // ─────────────────────────────────────────────────────────────

    /// @notice Refund the entire escrowed balance for a job back to the depositor.
    /// @dev    Callers must hold RELEASER_ROLE.
    /// @param  depositor  Owner of the escrowed balance.
    /// @param  jobId      Job identifier.
    /// @param  token      ERC-20 token address.
    function refund(address depositor, uint256 jobId, address token)
        external
        nonReentrant
        onlyRole(RELEASER_ROLE)
    {
        uint256 bal = balances[depositor][jobId][token];
        if (bal == 0) revert InsufficientBalance(0, 1);

        // Effects
        balances[depositor][jobId][token] = 0;

        emit Refunded(depositor, jobId, token, bal);

        // Interaction
        IERC20(token).safeTransfer(depositor, bal);
    }

    // ─────────────────────────────────────────────────────────────
    // Views
    // ─────────────────────────────────────────────────────────────

    /// @notice Return the escrowed balance for a (depositor, jobId, token) triple.
    /// @param  depositor Address of the depositor.
    /// @param  jobId     Job identifier.
    /// @param  token     ERC-20 token address.
    /// @return balance   Current escrowed amount.
    function getBalance(address depositor, uint256 jobId, address token) external view returns (uint256 balance) {
        balance = balances[depositor][jobId][token];
    }
}
