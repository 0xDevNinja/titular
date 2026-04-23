// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Initializable} from "@openzeppelin/contracts/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts/proxy/utils/UUPSUpgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

/// @title Treasury
/// @notice Protocol treasury owned by a Safe multisig. Receives native + ERC-20 fees and
///         streams them to the FeeDistributor at the multisig's discretion.
/// @dev UUPS-upgradeable. Authorization for upgrades is gated on `owner`.
///      No pause in v1; an incident-response hook is flagged for post-audit (M10).
contract Treasury is Initializable, OwnableUpgradeable, UUPSUpgradeable {
    using SafeERC20 for IERC20;

    /// @notice Address of the FeeDistributor that receives streamed rewards.
    address public feeDistributor;

    /// @dev Emitted when tokens or native are withdrawn by the owner.
    event Withdrawn(address indexed token, address indexed to, uint256 amount);
    /// @dev Emitted when the owner streams tokens to the fee distributor.
    event StreamedToVe(address indexed token, uint256 amount);
    /// @dev Emitted when the fee distributor address is updated.
    event FeeDistributorSet(address indexed previous, address indexed next);
    /// @dev Emitted on native ETH receipt.
    event NativeReceived(address indexed from, uint256 amount);

    error ZeroAddress();
    error ZeroAmount();
    error NativeTransferFailed();

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @notice Initializes the proxy. Can only run once.
    /// @param safeOwner       Address set as the contract owner; expected to be a Safe multisig.
    /// @param feeDistributor_ Initial fee distributor; may be address(0) to be set later.
    function initialize(address safeOwner, address feeDistributor_) external initializer {
        if (safeOwner == address(0)) revert ZeroAddress();
        __Ownable_init(safeOwner);
        feeDistributor = feeDistributor_;
        if (feeDistributor_ != address(0)) emit FeeDistributorSet(address(0), feeDistributor_);
    }

    /// @notice Accept native ETH.
    receive() external payable {
        emit NativeReceived(msg.sender, msg.value);
    }

    /// @notice Withdraw an ERC-20 (or native ETH if `token == address(0)`) to `to`.
    /// @param token  ERC-20 address, or `address(0)` for native ETH.
    /// @param to     Recipient.
    /// @param amount Amount to withdraw.
    function withdraw(address token, address to, uint256 amount) external onlyOwner {
        if (to == address(0)) revert ZeroAddress();
        if (amount == 0) revert ZeroAmount();
        if (token == address(0)) {
            // slither-disable-next-line arbitrary-send-eth
            (bool ok,) = payable(to).call{value: amount}("");
            if (!ok) revert NativeTransferFailed();
        } else {
            IERC20(token).safeTransfer(to, amount);
        }
        emit Withdrawn(token, to, amount);
    }

    /// @notice Stream ERC-20 rewards to the configured fee distributor.
    /// @param token  ERC-20 address. address(0) is rejected (use `withdraw` for native).
    /// @param amount Amount to stream.
    function streamToVe(address token, uint256 amount) external onlyOwner {
        if (token == address(0)) revert ZeroAddress();
        if (amount == 0) revert ZeroAmount();
        address fd = feeDistributor;
        if (fd == address(0)) revert ZeroAddress();
        IERC20(token).safeTransfer(fd, amount);
        emit StreamedToVe(token, amount);
    }

    /// @notice Update the fee distributor address.
    /// @param next New fee distributor. May not be zero.
    function setFeeDistributor(address next) external onlyOwner {
        if (next == address(0)) revert ZeroAddress();
        address previous = feeDistributor;
        feeDistributor = next;
        emit FeeDistributorSet(previous, next);
    }

    // ---------------------------------------------------------------------
    // UUPS
    // ---------------------------------------------------------------------

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {
        if (newImplementation == address(0)) revert ZeroAddress();
    }
}
