// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/// @title BondingCurve
/// @notice Pre-graduation AMM for a single agent token. Prices the agent token against
///         the quote asset (TITU) using a constant-product `x * y = k` curve with
///         virtual reserves so the initial slope is bounded. Charges a fixed 1% swap
///         fee on both sides, routed to a shared `feeRouter`. Trading is halted once
///         `realQuoteReserve` reaches `graduationThreshold`; a graduator contract then
///         drains both sides in a single atomic hand-off.
/// @dev    Not cloneable — each agent deploys its own instance via the launchpad
///         factory's `new BondingCurve(...)`. Owner is the factory deployer and gates
///         `setGraduator` (one-shot) and `setFeeRouter` (rebindable). CEI is enforced
///         on every external entry point: reserve deltas land before any outbound
///         ERC-20 call, and `ReentrancyGuard` is in place for the `buy`, `sell`, and
///         `pullForGraduation` paths. Custom errors are used for every revert; no
///         string messages. All constant-product math uses multiplication-before-
///         division and keeps `k` non-decreasing across trades modulo integer floor
///         on the output side.
contract BondingCurve is Ownable, ReentrancyGuard {
    using SafeERC20 for IERC20;

    // ---------------------------------------------------------------------
    // Immutables / constants
    // ---------------------------------------------------------------------

    /// @notice Agent ERC-20 traded by this curve. Fixed at construction.
    address public immutable agentToken;

    /// @notice Quote ERC-20 (TITU). Fixed at construction.
    address public immutable quoteToken;

    /// @notice Virtual quote reserve used for price discovery. Seeded to prevent a
    ///         degenerate zero-reserve price at launch.
    uint256 public immutable virtualQuoteReserve;

    /// @notice Virtual agent reserve used for price discovery. This is the *initial*
    ///         effective Y used in the constant-product formula — i.e. the notional
    ///         agent inventory the curve "starts with" for pricing purposes. The
    ///         physical balance of the curve starts at {initialAgentReserve}; the
    ///         difference `virtualAgentReserve - initialAgentReserve` acts as a
    ///         permanent offset on the agent side so that at graduation some agent
    ///         inventory is always left in the curve to seed the Uniswap LP.
    uint256 public immutable virtualAgentReserve;

    /// @notice Physical agent token amount deposited into the curve at construction.
    ///         Stored immutably so the constant-product math can reconstruct the
    ///         agent-side offset `virtualAgentReserve - initialAgentReserve` on the
    ///         fly without loading state. See {_effectiveAgentReserve}.
    uint256 public immutable initialAgentReserve;

    /// @notice Real TITU raised at which the curve halts trading and becomes eligible
    ///         for graduation. Virtual reserves are excluded from this threshold.
    uint256 public immutable graduationThreshold;

    /// @notice Swap fee in basis points applied to the input side of both buy and
    ///         sell. 100 bps = 1%.
    uint16 public constant FEE_BPS = 100;

    /// @notice Basis-point denominator used throughout fee math.
    uint16 public constant BPS_DENOMINATOR = 10_000;

    // ---------------------------------------------------------------------
    // Storage
    // ---------------------------------------------------------------------

    /// @notice TITU balance owed to the curve (monotonically increases on buy,
    ///         decreases on sell). Does NOT include virtual reserve or collected fees.
    uint256 public realQuoteReserve;

    /// @notice Agent tokens still held by this curve for trading.
    uint256 public realAgentReserve;

    /// @notice Set to true once `requestGraduation` has fired. Disables trading.
    bool public graduated;

    /// @notice Contract authorized to drain the curve at graduation. Set once by owner
    ///         via {setGraduator} and cannot be rebound.
    address public graduator;

    /// @notice Recipient of all swap fees. Rebindable by owner via {setFeeRouter} to
    ///         support fee-routing migrations across the agent fleet.
    address public feeRouter;

    // ---------------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------------

    /// @dev Emitted on every successful {buy}. `fee` is the quote amount diverted to
    ///      {feeRouter}; `agentOut` is the tokens delivered to `buyer`.
    event Bought(address indexed buyer, uint256 quoteIn, uint256 agentOut, uint256 fee);

    /// @dev Emitted on every successful {sell}. `fee` is the quote amount diverted to
    ///      {feeRouter}; `quoteOut` is the net delivered to `seller`.
    event Sold(address indexed seller, uint256 agentIn, uint256 quoteOut, uint256 fee);

    /// @dev Emitted exactly once when {requestGraduation} succeeds. `timestamp` is the
    ///      block timestamp at which trading halted.
    event Graduated(uint256 quoteReserve, uint256 agentReserve, uint256 timestamp);

    /// @dev Emitted on {setGraduator}. `prev` is always address(0) because the setter
    ///      is one-shot; included for a uniform ABI with {FeeRouterSet}.
    event GraduatorSet(address indexed prev, address indexed next);

    /// @dev Emitted on {setFeeRouter}. Includes the previous value for indexers.
    event FeeRouterSet(address indexed prev, address indexed next);

    /// @dev Emitted on {pullForGraduation}. `to` is the graduator.
    event Pulled(address indexed to, uint256 quoteAmount, uint256 agentAmount);

    // ---------------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------------

    /// @dev Thrown when a constructor argument is the zero address.
    error ZeroAddress();
    /// @dev Thrown when a constructor argument is a zero numeric value where a
    ///      non-zero value is required (reserves / threshold).
    error ZeroValue();
    /// @dev Thrown when the virtual agent reserve is not strictly greater than the
    ///      physical agent inventory (would leave the curve unable to graduate with
    ///      non-zero LP seed).
    error InvalidReserves();
    /// @dev Thrown by trade entry points after graduation has fired.
    error AlreadyGraduated();
    /// @dev Thrown by {requestGraduation} when called below threshold.
    error BelowThreshold();
    /// @dev Thrown by {pullForGraduation} when called before graduation.
    error NotGraduated();
    /// @dev Thrown on zero-valued trade inputs.
    error ZeroAmount();
    /// @dev Thrown by trade entry points when computed output undershoots caller's
    ///      slippage floor.
    error InsufficientOutput();
    /// @dev Thrown by {buy} when the trade would push `realQuoteReserve` strictly
    ///      past {graduationThreshold}. The caller must size the trade to land at or
    ///      below the threshold, then invoke {requestGraduation}.
    error ExceedsGraduationThreshold();
    /// @dev Thrown by {setGraduator} on any call after the graduator has been set.
    error GraduatorAlreadySet();
    /// @dev Thrown when a non-graduator address tries to call {pullForGraduation}.
    error NotGraduator();

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    /// @param owner_                Initial owner; gates `setGraduator` and `setFeeRouter`.
    /// @param agentToken_           Agent ERC-20 traded by this curve.
    /// @param quoteToken_           Quote ERC-20 (TITU).
    /// @param feeRouter_            Initial fee recipient.
    /// @param virtualQuoteReserve_  Virtual quote reserve (e.g. 30_000e18 TITU).
    /// @param virtualAgentReserve_  Virtual agent reserve (e.g. 1_073_000_000e18).
    /// @param graduationThreshold_  Real quote threshold that freezes trading
    ///                              (e.g. 42_000e18 TITU).
    /// @param initialAgentReserve_  Real agent tokens already deposited into this curve
    ///                              and available for distribution. The factory is
    ///                              expected to transfer exactly this amount before
    ///                              any trade lands; the constructor does NOT pull.
    constructor(
        address owner_,
        address agentToken_,
        address quoteToken_,
        address feeRouter_,
        uint256 virtualQuoteReserve_,
        uint256 virtualAgentReserve_,
        uint256 graduationThreshold_,
        uint256 initialAgentReserve_
    ) Ownable(_validatedOwner(owner_)) {
        if (agentToken_ == address(0) || quoteToken_ == address(0) || feeRouter_ == address(0)) revert ZeroAddress();
        if (
            virtualQuoteReserve_ == 0 || virtualAgentReserve_ == 0 || graduationThreshold_ == 0
                || initialAgentReserve_ == 0
        ) revert ZeroValue();
        // Virtuals-style: virtual agent must be strictly greater than the physical
        // inventory so there is always LP seed left when the curve graduates.
        if (virtualAgentReserve_ <= initialAgentReserve_) revert InvalidReserves();

        agentToken = agentToken_;
        quoteToken = quoteToken_;
        feeRouter = feeRouter_;
        virtualQuoteReserve = virtualQuoteReserve_;
        virtualAgentReserve = virtualAgentReserve_;
        graduationThreshold = graduationThreshold_;
        initialAgentReserve = initialAgentReserve_;
        realAgentReserve = initialAgentReserve_;

        emit FeeRouterSet(address(0), feeRouter_);
    }

    /// @dev Lift the zero-address check above the `Ownable` constructor so callers
    ///      see our canonical {ZeroAddress} error instead of OZ's `OwnableInvalidOwner`.
    function _validatedOwner(address owner_) private pure returns (address) {
        if (owner_ == address(0)) revert ZeroAddress();
        return owner_;
    }

    // ---------------------------------------------------------------------
    // Trading
    // ---------------------------------------------------------------------

    /// @notice Buy agent tokens with `quoteIn` TITU.
    /// @dev    CEI:
    ///           1. Validate preconditions and compute `agentOut` / `fee` against
    ///              pre-trade virtual+real reserves.
    ///           2. Update storage so reserves reflect the post-trade state before
    ///              any external call.
    ///           3. Pull quote from caller; push fee to `feeRouter`; push agent tokens
    ///              to caller.
    ///         Reverts if the post-trade `realQuoteReserve` would exceed the
    ///         graduation threshold; callers must size the terminal buy exactly.
    /// @param  minAgentOut  Slippage floor (agent wei).
    /// @param  quoteIn      Quote in (TITU wei, pre-fee).
    function buy(uint256 minAgentOut, uint256 quoteIn) external nonReentrant {
        if (graduated) revert AlreadyGraduated();
        if (quoteIn == 0) revert ZeroAmount();

        (uint256 agentOut, uint256 fee) = _quoteBuy(quoteIn);
        if (agentOut < minAgentOut) revert InsufficientOutput();
        if (agentOut == 0) revert InsufficientOutput();
        if (agentOut > realAgentReserve) revert InsufficientOutput();

        uint256 quoteInAfterFee = quoteIn - fee;
        uint256 newRealQuote = realQuoteReserve + quoteInAfterFee;
        if (newRealQuote > graduationThreshold) revert ExceedsGraduationThreshold();

        // --- Effects ---
        realQuoteReserve = newRealQuote;
        realAgentReserve -= agentOut;

        // --- Interactions ---
        // Pull the gross quote-in from the buyer. Fee is diverted to `feeRouter`; the
        // remainder stays on this contract as pool reserve.
        IERC20 qt = IERC20(quoteToken);
        qt.safeTransferFrom(msg.sender, address(this), quoteIn);
        if (fee != 0) qt.safeTransfer(feeRouter, fee);
        IERC20(agentToken).safeTransfer(msg.sender, agentOut);

        emit Bought(msg.sender, quoteIn, agentOut, fee);
    }

    /// @notice Sell `agentIn` agent tokens for TITU.
    /// @dev    Mirrors {buy}. Fee is taken from the quote-out side so the seller bears
    ///         the swap tax. CEI ordering preserved.
    /// @param  agentIn       Agent tokens to sell (agent wei).
    /// @param  minQuoteOut   Slippage floor on the net quote delivered to the seller.
    function sell(uint256 agentIn, uint256 minQuoteOut) external nonReentrant {
        if (graduated) revert AlreadyGraduated();
        if (agentIn == 0) revert ZeroAmount();

        (uint256 quoteOut, uint256 fee) = _quoteSell(agentIn);
        if (quoteOut < minQuoteOut) revert InsufficientOutput();
        if (quoteOut == 0) revert InsufficientOutput();
        uint256 grossQuoteOut = quoteOut + fee;
        if (grossQuoteOut > realQuoteReserve) revert InsufficientOutput();

        // --- Effects ---
        realQuoteReserve -= grossQuoteOut;
        realAgentReserve += agentIn;

        // --- Interactions ---
        IERC20(agentToken).safeTransferFrom(msg.sender, address(this), agentIn);
        IERC20 qt = IERC20(quoteToken);
        if (fee != 0) qt.safeTransfer(feeRouter, fee);
        qt.safeTransfer(msg.sender, quoteOut);

        emit Sold(msg.sender, agentIn, quoteOut, fee);
    }

    // ---------------------------------------------------------------------
    // Views / quoting
    // ---------------------------------------------------------------------

    /// @notice Preview a buy without mutating state.
    /// @return agentOut Net agent tokens the caller would receive.
    /// @return fee      Quote wei diverted to `feeRouter`.
    function quoteBuy(uint256 quoteIn) external view returns (uint256 agentOut, uint256 fee) {
        if (quoteIn == 0) return (0, 0);
        return _quoteBuy(quoteIn);
    }

    /// @notice Preview a sell without mutating state.
    /// @return quoteOut Net quote wei the caller would receive.
    /// @return fee      Quote wei diverted to `feeRouter`.
    function quoteSell(uint256 agentIn) external view returns (uint256 quoteOut, uint256 fee) {
        if (agentIn == 0) return (0, 0);
        return _quoteSell(agentIn);
    }

    /// @notice Current constant-product `k = (virtualQuote + realQuote) * Y_effective`,
    ///         where `Y_effective = virtualAgentReserve + realAgentReserve - initialAgentReserve`.
    ///         Equivalently `Y_effective = realAgentReserve + agentOffset` where
    ///         `agentOffset = virtualAgentReserve - initialAgentReserve` is a
    ///         construction-time constant.
    function k() external view returns (uint256) {
        return (virtualQuoteReserve + realQuoteReserve) * _effectiveAgentReserve(realAgentReserve);
    }

    // ---------------------------------------------------------------------
    // Graduation
    // ---------------------------------------------------------------------

    /// @notice Freeze trading once the curve has raised >= `graduationThreshold` in
    ///         real quote. Callable by anyone — this is a keeper-style trigger that
    ///         runs after the final buy lands. Does not move funds; the graduator
    ///         pulls them in a separate tx via {pullForGraduation}.
    function requestGraduation() external {
        if (graduated) revert AlreadyGraduated();
        if (realQuoteReserve < graduationThreshold) revert BelowThreshold();

        graduated = true;
        emit Graduated(realQuoteReserve, realAgentReserve, block.timestamp);
    }

    /// @notice Drain both reserves to the graduator. Only the configured graduator
    ///         may call, and only after {requestGraduation} has flipped the flag.
    /// @dev    Zeroes both reserves BEFORE transferring out (CEI) so any reentrant
    ///         callback into the curve sees a drained state.
    /// @return quoteAmount Quote transferred to the graduator.
    /// @return agentAmount Agent tokens transferred to the graduator.
    function pullForGraduation() external nonReentrant returns (uint256 quoteAmount, uint256 agentAmount) {
        if (!graduated) revert NotGraduated();
        if (msg.sender != graduator) revert NotGraduator();

        quoteAmount = realQuoteReserve;
        agentAmount = realAgentReserve;

        // --- Effects ---
        realQuoteReserve = 0;
        realAgentReserve = 0;

        // --- Interactions ---
        if (quoteAmount != 0) IERC20(quoteToken).safeTransfer(msg.sender, quoteAmount);
        if (agentAmount != 0) IERC20(agentToken).safeTransfer(msg.sender, agentAmount);

        emit Pulled(msg.sender, quoteAmount, agentAmount);
    }

    // ---------------------------------------------------------------------
    // Admin
    // ---------------------------------------------------------------------

    /// @notice Bind the graduator address. One-shot — cannot be changed once set.
    /// @dev    Guarded by Ownable. The `graduator` is trusted with drain rights on
    ///         both reserves post-graduation, so re-binding would be a protocol-
    ///         level rug vector. Owner is expected to be a Safe multisig.
    function setGraduator(address next) external onlyOwner {
        if (next == address(0)) revert ZeroAddress();
        if (graduator != address(0)) revert GraduatorAlreadySet();

        graduator = next;
        emit GraduatorSet(address(0), next);
    }

    /// @notice Rebind the swap-fee recipient. Rebindable because FeeRouter may be
    ///         upgraded across the fleet and we want to retarget existing curves
    ///         without redeploys.
    function setFeeRouter(address next) external onlyOwner {
        if (next == address(0)) revert ZeroAddress();
        address prev = feeRouter;
        feeRouter = next;
        emit FeeRouterSet(prev, next);
    }

    // ---------------------------------------------------------------------
    // Internal math
    // ---------------------------------------------------------------------

    /// @dev Constant-product buy math. Fee is taken from `quoteIn` up front; the
    ///      remainder lands on the curve. `agentOut` is floored via integer division
    ///      so post-trade `k` is >= pre-trade `k`.
    ///
    ///      Let
    ///        X = virtualQuoteReserve + realQuoteReserve
    ///        Y = virtualAgentReserve + realAgentReserve
    ///        k = X * Y
    ///        dx' = quoteIn - fee           (post-fee quote going into the pool)
    ///      Then
    ///        agentOut = Y - ceil(k / (X + dx'))  // to ensure k' >= k with floor div
    ///                 = Y - (k + (X + dx') - 1) / (X + dx')
    ///      Equivalently (standard Uniswap V2 fee-less form, which also floors out):
    ///        agentOut = (Y * dx') / (X + dx')
    ///      We use the latter: it is the same floor behaviour and marginally cheaper.
    ///
    ///      Overflow: `Y * dx'` is bounded by
    ///        Y  <= virtualAgentReserve + initialAgentReserve <= ~2.1e27 (< 2^91)
    ///        dx' <= graduationThreshold + virtualQuote  <= ~1e23 (< 2^77)
    ///      product < 2^168. Well below 2^256.
    function _quoteBuy(uint256 quoteIn) internal view returns (uint256 agentOut, uint256 fee) {
        fee = (quoteIn * FEE_BPS) / BPS_DENOMINATOR;
        uint256 dxPrime = quoteIn - fee;
        uint256 x = virtualQuoteReserve + realQuoteReserve;
        uint256 y = _effectiveAgentReserve(realAgentReserve);
        // multiplication-before-division, Uniswap V2 floor form.
        agentOut = (y * dxPrime) / (x + dxPrime);
    }

    /// @dev Constant-product sell math. `agentIn` is the exact tokens the seller
    ///      hands over; the gross quote-out is computed against the pool, then a
    ///      fee is peeled off for the router.
    ///
    ///      grossOut = (X * agentIn) / (Y + agentIn)   // floor
    ///      fee      = grossOut * FEE_BPS / BPS_DENOMINATOR
    ///      netOut   = grossOut - fee
    function _quoteSell(uint256 agentIn) internal view returns (uint256 quoteOut, uint256 fee) {
        uint256 x = virtualQuoteReserve + realQuoteReserve;
        uint256 y = _effectiveAgentReserve(realAgentReserve);
        uint256 yPlus = y + agentIn;
        // To avoid Slither's `divide-before-multiply` warning and stay precise, compute
        // the fee directly from the un-floored rational `x*agentIn / yPlus` by folding
        // the bps division into one step:
        //   fee = floor( (x * agentIn * FEE_BPS) / (yPlus * BPS_DENOMINATOR) )
        // Then gross = floor(x*agentIn / yPlus); quoteOut = gross - fee.
        // This rounds in favour of the pool on both terms and keeps `k_after >= k_before`.
        // Overflow: x <= virtualQuote + graduationThreshold (<2^77) and agentIn <=
        // initialAgentReserve (<2^91), so the numerator fits comfortably in 256 bits.
        uint256 grossOut = (x * agentIn) / yPlus;
        fee = (x * agentIn * FEE_BPS) / (yPlus * BPS_DENOMINATOR);
        quoteOut = grossOut - fee;
    }

    /// @dev Effective Y used in the constant-product formula. The agent side is
    ///      tracked as a physical balance (`realAgentReserve` starts at
    ///      `initialAgentReserve` and ticks down on buy / up on sell). The
    ///      pricing Y starts at `virtualAgentReserve` which, by construction,
    ///      equals `initialAgentReserve + (virtualAgentReserve - initialAgentReserve)`.
    ///      Rearranging, Y_eff(realAgent) = realAgent + (virtualAgentReserve -
    ///      initialAgentReserve). The subtraction cannot underflow — the
    ///      constructor asserts `virtualAgentReserve > initialAgentReserve` — and
    ///      Solidity 0.8 checked math would catch it regardless.
    function _effectiveAgentReserve(uint256 realAgent) internal view returns (uint256) {
        unchecked {
            // virtualAgentReserve > initialAgentReserve enforced in constructor, so
            // the offset is strictly positive. realAgent + offset fits in 2^256 for
            // the parameter ranges we support (see overflow note on {_quoteBuy}).
            return realAgent + (virtualAgentReserve - initialAgentReserve);
        }
    }
}
