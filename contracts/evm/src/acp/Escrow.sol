// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {IPermit2} from "./IPermit2.sol";

/// @title Escrow
/// @notice Multi-token escrow with atomic multi-recipient release and Permit2 support.
/// @dev    Funds are held in this contract; balances tracked per (depositor, jobId, token).
///         `release` atomically transfers to N recipients.
///         `fundWithPermit2` allows gasless ERC-20 approval via Permit2 EIP-712 signature.
///         RELEASER_ROLE: addresses allowed to call `release` and `refund`.
///         CEI: all storage writes before external SafeERC20/Permit2 transfers.
///         ReentrancyGuard on `fund`, `fundWithPermit2`, `release`, `refund`.
contract Escrow is AccessControl, ReentrancyGuard {
    using SafeERC20 for IERC20;

    bytes32 public constant RELEASER_ROLE = keccak256("RELEASER_ROLE");

    /// @notice Canonical Permit2 contract (immutable).
    IPermit2 public immutable permit2;

    struct Recipient {
        address to;
        uint256 amount;
    }

    /// @notice depositor => jobId => token => escrowed balance.
    mapping(address => mapping(uint256 => mapping(address => uint256))) public balances;

    event Funded(address indexed depositor, uint256 indexed jobId, address indexed token, uint256 amount);
    event Released(
        address indexed depositor, uint256 indexed jobId, address indexed token, address to, uint256 amount
    );
    event Refunded(address indexed depositor, uint256 indexed jobId, address indexed token, uint256 amount);

    error ZeroAddress();
    error ZeroAmount();
    error InsufficientBalance(uint256 have, uint256 want);
    error EmptyRecipients();
    error RecipientZeroAddress(uint256 index);
    error RecipientZeroAmount(uint256 index);
    error SumExceedsBalance(uint256 sum, uint256 balance);

    /// @param admin    Address granted DEFAULT_ADMIN_ROLE and RELEASER_ROLE.
    /// @param _permit2 Permit2 contract address (canonical: 0x000000000022D473030F116dDEE9F6B43aC78BA3).
    constructor(address admin, address _permit2) {
        if (admin == address(0)) revert ZeroAddress();
        if (_permit2 == address(0)) revert ZeroAddress();
        permit2 = IPermit2(_permit2);
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(RELEASER_ROLE, admin);
    }

    /// @notice Deposit via standard ERC-20 approve+transferFrom.
    /// @param depositor Address whose balance increases.
    /// @param jobId Opaque job identifier.
    /// @param token ERC-20 token address.
    /// @param amount Amount to deposit (must be > 0).
    function fund(address depositor, uint256 jobId, address token, uint256 amount) external nonReentrant {
        if (depositor == address(0)) revert ZeroAddress();
        if (token == address(0)) revert ZeroAddress();
        if (amount == 0) revert ZeroAmount();
        balances[depositor][jobId][token] += amount;
        emit Funded(depositor, jobId, token, amount);
        IERC20(token).safeTransferFrom(msg.sender, address(this), amount);
    }

    /// @notice Deposit via Permit2 signed transfer (no prior approve needed).
    /// @param depositor Address whose balance increases (must be the Permit2 signer).
    /// @param jobId Opaque job identifier.
    /// @param token ERC-20 token address.
    /// @param amount Amount to deposit (must be > 0).
    /// @param permit Permit2 PermitTransferFrom struct.
    /// @param signature EIP-712 signature by depositor over permit.
    function fundWithPermit2(
        address depositor,
        uint256 jobId,
        address token,
        uint256 amount,
        IPermit2.PermitTransferFrom calldata permit,
        bytes calldata signature
    ) external nonReentrant {
        if (depositor == address(0)) revert ZeroAddress();
        if (token == address(0)) revert ZeroAddress();
        if (amount == 0) revert ZeroAmount();
        // Effects before interaction (CEI)
        balances[depositor][jobId][token] += amount;
        emit Funded(depositor, jobId, token, amount);
        // Interaction via Permit2
        permit2.permitTransferFrom(
            permit,
            IPermit2.SignatureTransferDetails({to: address(this), requestedAmount: amount}),
            depositor,
            signature
        );
    }

    /// @notice Atomically release funds to multiple recipients.
    /// @param depositor Owner of the escrowed balance.
    /// @param jobId Job identifier.
    /// @param token ERC-20 token address.
    /// @param recipients Array of (address, amount) pairs.
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
        balances[depositor][jobId][token] = bal - total;
        for (uint256 i = 0; i < recipients.length; ++i) {
            emit Released(depositor, jobId, token, recipients[i].to, recipients[i].amount);
            IERC20(token).safeTransfer(recipients[i].to, recipients[i].amount);
        }
    }

    /// @notice Refund entire escrowed balance back to depositor.
    /// @param depositor Owner of the escrowed balance.
    /// @param jobId Job identifier.
    /// @param token ERC-20 token address.
    function refund(address depositor, uint256 jobId, address token)
        external
        nonReentrant
        onlyRole(RELEASER_ROLE)
    {
        uint256 bal = balances[depositor][jobId][token];
        if (bal == 0) revert InsufficientBalance(0, 1);
        balances[depositor][jobId][token] = 0;
        emit Refunded(depositor, jobId, token, bal);
        IERC20(token).safeTransfer(depositor, bal);
    }

    /// @notice Return the escrowed balance for a (depositor, jobId, token) triple.
    /// @param depositor Address of the depositor.
    /// @param jobId Job identifier.
    /// @param token ERC-20 token address.
    /// @return balance Current escrowed amount.
    function getBalance(address depositor, uint256 jobId, address token) external view returns (uint256 balance) {
        balance = balances[depositor][jobId][token];
    }
}
