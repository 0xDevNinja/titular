// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

/// @dev Minimal interface for Uniswap V2 router used by this contract.
interface IUniswapV2Router02 {
    function swapExactTokensForTokensSupportingFeeOnTransferTokens(
        uint256 amountIn,
        uint256 amountOutMin,
        address[] calldata path,
        address to,
        uint256 deadline
    ) external;
}

/// @dev ERC-20 with burn capability.
interface IERC20Burnable is IERC20 {
    /// @dev Burns `amount` of the caller's own tokens.
    function burn(uint256 amount) external;
}

/// @title BuybackBurner
/// @notice Receives a payment token, swaps it for TITU on Uniswap V2, then burns
///         the TITU using ERC20Burnable.burnFrom.
/// @dev    Role model:
///           - DEFAULT_ADMIN_ROLE: update router, path, and slippage parameters.
///           - EXECUTOR_ROLE: allowed to call `buybackAndBurn`.
///
///         Slippage: caller supplies `amountOutMin`; the contract enforces a
///         global floor `minOutBps` (in BPS of input amount) as an extra safety
///         net. Set `minOutBps = 0` to disable the floor (not recommended).
///
///         CEI: safeApprove the router, then swap, then burn.  No state changes
///         occur after the external calls, so ReentrancyGuard covers residual risk.
contract BuybackBurner is AccessControl, ReentrancyGuard {
    using SafeERC20 for IERC20;

    // ─────────────────────────────────────────────────────────────
    // Roles
    // ─────────────────────────────────────────────────────────────

    bytes32 public constant EXECUTOR_ROLE = keccak256("EXECUTOR_ROLE");

    // ─────────────────────────────────────────────────────────────
    // Constants
    // ─────────────────────────────────────────────────────────────

    uint256 public constant BPS = 10_000;

    // ─────────────────────────────────────────────────────────────
    // State
    // ─────────────────────────────────────────────────────────────

    /// @notice Uniswap V2 compatible router.
    IUniswapV2Router02 public router;

    /// @notice Payment token received from FeeSplitter (e.g. USDC).
    address public immutable paymentToken;

    /// @notice TITU token that is bought and burned.
    IERC20Burnable public immutable titu;

    /// @notice Swap path from paymentToken → TITU.  Must start with paymentToken
    ///         and end with address(titu).  May include intermediate hops.
    address[] public swapPath;

    /// @notice Minimum output in BPS of the input amount.  Acts as a global floor.
    ///         E.g. 9000 = 90%. Set to 0 to disable.
    uint256 public minOutBps;

    /// @notice Swap deadline extension in seconds added to block.timestamp.
    uint256 public immutable swapDeadlineBuffer;

    // ─────────────────────────────────────────────────────────────
    // Events
    // ─────────────────────────────────────────────────────────────

    /// @notice Emitted when a buyback-and-burn is executed.
    event BuybackAndBurn(address indexed executor, address indexed paymentToken, uint256 amountIn, uint256 tituBurned);

    /// @notice Emitted when router is updated.
    event RouterUpdated(address indexed oldRouter, address indexed newRouter);

    /// @notice Emitted when swap path is updated.
    event SwapPathUpdated(address[] path);

    /// @notice Emitted when min-out BPS floor is updated.
    event MinOutBpsUpdated(uint256 oldBps, uint256 newBps);

    // ─────────────────────────────────────────────────────────────
    // Errors
    // ─────────────────────────────────────────────────────────────

    error ZeroAddress();
    error ZeroAmount();
    error InvalidSwapPath();
    error InvalidBps();
    error InsufficientAmountOut(uint256 amountIn, uint256 minOut, uint256 actualOut);

    // ─────────────────────────────────────────────────────────────
    // Constructor
    // ─────────────────────────────────────────────────────────────

    /// @param admin             Address granted DEFAULT_ADMIN_ROLE.
    /// @param _router           Uniswap V2 router address.
    /// @param _paymentToken     ERC-20 received and swapped.
    /// @param _titu             TITU token address (must implement burnFrom).
    /// @param _swapPath         Initial swap path (must start with _paymentToken, end with _titu).
    /// @param _minOutBps        Minimum output BPS (0–10 000).
    /// @param _swapDeadlineBuffer Seconds added to block.timestamp for swap deadline.
    constructor(
        address admin,
        address _router,
        address _paymentToken,
        address _titu,
        address[] memory _swapPath,
        uint256 _minOutBps,
        uint256 _swapDeadlineBuffer
    ) {
        if (admin == address(0)) revert ZeroAddress();
        if (_router == address(0)) revert ZeroAddress();
        if (_paymentToken == address(0)) revert ZeroAddress();
        if (_titu == address(0)) revert ZeroAddress();
        if (_minOutBps > BPS) revert InvalidBps();
        _validatePath(_swapPath, _paymentToken, _titu);

        router = IUniswapV2Router02(_router);
        paymentToken = _paymentToken;
        titu = IERC20Burnable(_titu);
        swapPath = _swapPath;
        minOutBps = _minOutBps;
        swapDeadlineBuffer = _swapDeadlineBuffer;

        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(EXECUTOR_ROLE, admin);
    }

    // ─────────────────────────────────────────────────────────────
    // Core action
    // ─────────────────────────────────────────────────────────────

    /// @notice Swap all held `paymentToken` for TITU via Uniswap V2 and burn the TITU.
    /// @dev    Caller must hold EXECUTOR_ROLE.
    /// @param  amountIn    Amount of paymentToken to swap.
    /// @param  amountOutMin Minimum TITU out (caller's slippage tolerance).
    function buybackAndBurn(uint256 amountIn, uint256 amountOutMin) external nonReentrant onlyRole(EXECUTOR_ROLE) {
        if (amountIn == 0) revert ZeroAmount();

        // Enforce global floor if configured
        if (minOutBps > 0) {
            uint256 floor = (amountIn * minOutBps) / BPS;
            if (amountOutMin < floor) {
                amountOutMin = floor;
            }
        }

        address _paymentToken = paymentToken;
        address _titu = address(titu);

        // Record TITU balance before swap to calculate actual amount received
        uint256 tituBefore = titu.balanceOf(address(this));

        // Set exact approval for router (forceApprove resets to 0 first, safe for USDT-like tokens)
        IERC20(_paymentToken).forceApprove(address(router), amountIn);

        // Swap paymentToken → TITU; tokens land in this contract
        // slither-disable-next-line timestamp
        router.swapExactTokensForTokensSupportingFeeOnTransferTokens(
            amountIn, amountOutMin, swapPath, address(this), block.timestamp + swapDeadlineBuffer
        );

        uint256 tituReceived = titu.balanceOf(address(this)) - tituBefore;

        // Reset router allowance to 0 (defensive)
        IERC20(_paymentToken).forceApprove(address(router), 0);

        // Burn — ERC20Burnable.burn burns caller's own tokens (no allowance needed)
        titu.burn(tituReceived);

        emit BuybackAndBurn(msg.sender, _paymentToken, amountIn, tituReceived);
    }

    // ─────────────────────────────────────────────────────────────
    // Admin configuration
    // ─────────────────────────────────────────────────────────────

    /// @notice Update the Uniswap V2 router address.
    /// @param  newRouter New router address (must be non-zero).
    function setRouter(address newRouter) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (newRouter == address(0)) revert ZeroAddress();
        address old = address(router);
        router = IUniswapV2Router02(newRouter);
        emit RouterUpdated(old, newRouter);
    }

    /// @notice Update the swap path.
    /// @param  path New swap path; must start with paymentToken and end with titu.
    function setSwapPath(address[] calldata path) external onlyRole(DEFAULT_ADMIN_ROLE) {
        _validatePath(path, paymentToken, address(titu));
        swapPath = path;
        emit SwapPathUpdated(path);
    }

    /// @notice Update the minimum output BPS floor.
    /// @param  newBps New BPS value (0–10 000).
    function setMinOutBps(uint256 newBps) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (newBps > BPS) revert InvalidBps();
        uint256 old = minOutBps;
        minOutBps = newBps;
        emit MinOutBpsUpdated(old, newBps);
    }

    /// @notice Rescue ERC-20 tokens accidentally sent to this contract (excludes paymentToken mid-swap).
    /// @param  tokenAddr Token to rescue.
    /// @param  to        Recipient.
    /// @param  amount    Amount to rescue.
    function rescueTokens(address tokenAddr, address to, uint256 amount) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (to == address(0)) revert ZeroAddress();
        IERC20(tokenAddr).safeTransfer(to, amount);
    }

    // ─────────────────────────────────────────────────────────────
    // Views
    // ─────────────────────────────────────────────────────────────

    /// @notice Return the full swap path array.
    /// @return path Swap path.
    function getSwapPath() external view returns (address[] memory path) {
        path = swapPath;
    }

    // ─────────────────────────────────────────────────────────────
    // Internal helpers
    // ─────────────────────────────────────────────────────────────

    function _validatePath(address[] memory path, address _paymentToken, address _titu) internal pure {
        if (path.length < 2) revert InvalidSwapPath();
        if (path[0] != _paymentToken) revert InvalidSwapPath();
        if (path[path.length - 1] != _titu) revert InvalidSwapPath();
    }
}
