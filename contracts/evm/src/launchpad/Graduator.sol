// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

import {IBondingCurve} from "./interfaces/IBondingCurve.sol";
import {IUniswapV2Factory} from "./interfaces/IUniswapV2Factory.sol";
import {IUniswapV2Router02} from "./interfaces/IUniswapV2Router02.sol";

/// @title Graduator
/// @notice Single shared graduation sink for the launchpad fleet. Every agent
///         bonding curve registers the same `Graduator` instance as its
///         `graduator`; on graduation this contract pulls both reserves from
///         the curve in one atomic hand-off, creates (or reuses) the Uniswap
///         V2 pair for `(agentToken, quoteToken)`, and routes the freshly
///         minted LP tokens directly into a per-curve `LPLock` address that
///         the factory pre-registers via {registerLPLock}.
/// @dev    Owner is the launchpad factory (or a Safe that stands in for it).
///         Only the owner can bind lock targets. `graduate` is permissionless
///         — any keeper can call it once {IBondingCurve.requestGraduation}
///         has flipped the curve. CEI is enforced: the graduated flag is set
///         before the reserves are pulled and before the router interaction,
///         and `ReentrancyGuard` protects the transfer path.
///
///         Design notes:
///         - `router` and `uniFactory` are immutable. Binding the router per
///           network at construction means the factory is implicitly pinned
///           and never needs to be reset — any post-deploy re-binding would
///           be a rug vector across the fleet.
///         - LP tokens are sent directly to the `lpLocks[curve]` address via
///           the router's `to` parameter, so this contract never custodies
///           LP. The lock contract is free to enforce any vesting/burn policy.
///         - A 50 bps (0.5%) slippage floor is applied symmetrically on both
///           sides of `addLiquidity`. The only realistic cause of shortfall
///           is a race with a front-runner pre-seeding the pair; the floor
///           bounds the loss.
contract Graduator is Ownable, ReentrancyGuard {
    using SafeERC20 for IERC20;

    // ---------------------------------------------------------------------
    // Constants
    // ---------------------------------------------------------------------

    /// @notice Slippage floor applied to both sides of `addLiquidity`, in bps
    ///         of the pulled amount. 9_950 / 10_000 = 0.5% tolerance. The
    ///         pair is always freshly seeded by this contract (or empty), so
    ///         the only realistic shortfall is a front-runner seeding the
    ///         pair between {createPair} and {addLiquidity}.
    uint256 public constant MIN_LIQUIDITY_BPS = 9950;

    /// @notice Denominator for basis-point math.
    uint256 public constant BPS_DENOMINATOR = 10_000;

    // ---------------------------------------------------------------------
    // Immutables
    // ---------------------------------------------------------------------

    /// @notice Uniswap V2 router wired at construction. Pinned per deploy.
    IUniswapV2Router02 public immutable router;

    /// @notice Uniswap V2 factory derived from `router.factory()` at deploy.
    ///         Cached so `graduate` never re-hits the router view.
    IUniswapV2Factory public immutable uniFactory;

    // ---------------------------------------------------------------------
    // Storage
    // ---------------------------------------------------------------------

    /// @notice LP token recipient for each curve. Set by the factory via
    ///         {registerLPLock} before the curve graduates.
    mapping(address curve => address lpLock) public lpLocks;

    /// @notice Tracks curves that have already been graduated through this
    ///         instance. Idempotency guard — `graduate` reverts on replay so
    ///         a keeper cannot double-add liquidity.
    mapping(address curve => bool) public graduated;

    // ---------------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------------

    /// @dev Emitted on {registerLPLock}. Indexers can fan out per agent.
    event LPLockRegistered(address indexed curve, address indexed lock);

    /// @dev Emitted once per curve from {graduate}. `pair` is the V2 pair
    ///      address (created here or pre-existing), `liquidity` is the LP
    ///      tokens minted to the lock, and `quoteAmount`/`agentAmount` are
    ///      the physical tokens pushed into the pair (may be slightly lower
    ///      than the pulled amounts if the router rebalanced to match an
    ///      existing reserve ratio).
    event Graduated(
        address indexed curve, address indexed pair, uint256 liquidity, uint256 quoteAmount, uint256 agentAmount
    );

    // ---------------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------------

    /// @dev Thrown when a zero address is supplied where a non-zero is required.
    error ZeroAddress();

    /// @dev Thrown by {graduate} when the target curve has already been processed
    ///      through this graduator.
    error AlreadyGraduated();

    /// @dev Thrown by {graduate} when the curve has not flipped its own
    ///      `graduated` flag yet (i.e. the keeper raced
    ///      {IBondingCurve.requestGraduation}).
    error NotYetGraduatedOnCurve();

    /// @dev Thrown by {graduate} when the factory has not registered an LP lock
    ///      for this curve. Prevents a graduation path with no vesting target.
    error LPLockNotRegistered();

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    /// @param router_  Uniswap V2 router for this network.
    /// @param owner_   Initial owner (launchpad factory or Safe). Gates
    ///                 {registerLPLock}.
    constructor(IUniswapV2Router02 router_, address owner_) Ownable(_validatedOwner(owner_)) {
        if (address(router_) == address(0)) revert ZeroAddress();
        router = router_;
        // `router.factory()` is a view and returns a pinned immutable on the
        // reference UniswapV2Router02. Caching avoids redundant calls and lets
        // us revert fast if a misconfigured router is passed.
        address fac = router_.factory();
        if (fac == address(0)) revert ZeroAddress();
        uniFactory = IUniswapV2Factory(fac);
    }

    /// @dev Hoist the zero-check above the Ownable constructor so callers see
    ///      our {ZeroAddress} error instead of OZ's `OwnableInvalidOwner`.
    function _validatedOwner(address owner_) private pure returns (address) {
        if (owner_ == address(0)) revert ZeroAddress();
        return owner_;
    }

    // ---------------------------------------------------------------------
    // Admin
    // ---------------------------------------------------------------------

    /// @notice Bind the LP-lock recipient for a specific curve. Called by the
    ///         launchpad factory immediately after it deploys the curve and
    ///         the matching lock. Overwriting is permitted by the owner so a
    ///         mis-wired lock can be patched before graduation fires; once
    ///         {graduate} runs, the recipient is immaterial.
    /// @param  curve   Bonding-curve address the lock is tied to.
    /// @param  lpLock  Address that will receive the minted LP tokens.
    function registerLPLock(address curve, address lpLock) external onlyOwner {
        if (curve == address(0) || lpLock == address(0)) revert ZeroAddress();
        lpLocks[curve] = lpLock;
        emit LPLockRegistered(curve, lpLock);
    }

    // ---------------------------------------------------------------------
    // Graduation
    // ---------------------------------------------------------------------

    /// @notice Finalize graduation for a curve: pull reserves, create or
    ///         reuse the V2 pair, and add liquidity into the LP-lock.
    /// @dev    CEI ordering:
    ///           1. Check curve has flipped `graduated = true` on its side.
    ///           2. Check this graduator has not already processed the curve.
    ///           3. Check the factory has wired an LP lock for this curve.
    ///           4. Mark `graduated[curve] = true` (effect BEFORE interaction).
    ///           5. Pull reserves from the curve in a single call.
    ///           6. Approve the router for exactly the pulled amounts (via
    ///              `safeIncreaseAllowance` so we never leave a surplus).
    ///           7. Create the pair if missing.
    ///           8. Call `router.addLiquidity` with 0.5% slippage floors and
    ///              `block.timestamp` as the deadline (executed in the same
    ///              block, so a timestamp deadline is the tightest bound).
    ///           9. Emit `Graduated`.
    ///         The effect-before-interaction step is what makes this safe
    ///         against a hostile token implementation: even if the agent
    ///         token's transfer reenters, the idempotency guard has already
    ///         flipped.
    /// @param  curve  Bonding-curve address to graduate.
    function graduate(address curve) external nonReentrant {
        // --- Checks ---
        if (!IBondingCurve(curve).graduated()) revert NotYetGraduatedOnCurve();
        if (graduated[curve]) revert AlreadyGraduated();

        address lock = lpLocks[curve];
        if (lock == address(0)) revert LPLockNotRegistered();

        // --- Effects ---
        // Flip the idempotency guard BEFORE any external call. If the curve's
        // `pullForGraduation` or the router's `addLiquidity` tries to reenter
        // `graduate` (directly or via an ERC-777-style callback on the agent
        // token), the early `AlreadyGraduated` revert kills the recursion.
        graduated[curve] = true;

        // --- Interactions ---
        address agentToken = IBondingCurve(curve).agentToken();
        address quoteToken = IBondingCurve(curve).quoteToken();
        (uint256 quoteAmount, uint256 agentAmount) = IBondingCurve(curve).pullForGraduation();

        address pair = _ensurePairAndApprove(agentToken, quoteToken, quoteAmount, agentAmount);
        // `agentUsed` and `quoteUsed` are what the router actually pulled into
        // the pair — may be < pulled amounts if the router rebalanced against
        // a front-run seeded reserve. Surplus allowance is dropped after the
        // add, so any unused approval is a no-op for future calls (graduation
        // is one-shot per curve).
        (uint256 agentUsed, uint256 quoteUsed, uint256 liquidity) =
            _addLiquidity(agentToken, quoteToken, quoteAmount, agentAmount, lock);

        emit Graduated(curve, pair, liquidity, quoteUsed, agentUsed);
    }

    /// @dev Issue fresh router approvals sized to the pulled reserves and
    ///      return the V2 pair address, creating it when missing. Split out
    ///      of {graduate} purely to keep the main function's stack under the
    ///      EVM's 16-slot budget when via-ir is disabled.
    function _ensurePairAndApprove(address agentToken, address quoteToken, uint256 quoteAmount, uint256 agentAmount)
        private
        returns (address pair)
    {
        // Additive approvals; graduator starts at zero allowance for every
        // new curve, so this is equivalent to a fresh approve.
        IERC20(quoteToken).safeIncreaseAllowance(address(router), quoteAmount);
        IERC20(agentToken).safeIncreaseAllowance(address(router), agentAmount);

        pair = uniFactory.getPair(agentToken, quoteToken);
        if (pair == address(0)) {
            pair = uniFactory.createPair(agentToken, quoteToken);
        }
    }

    /// @dev Thin wrapper over `router.addLiquidity` with a 0.5% symmetric
    ///      slippage floor. Separated from {graduate} for stack reasons; do
    ///      not inline.
    function _addLiquidity(
        address agentToken,
        address quoteToken,
        uint256 quoteAmount,
        uint256 agentAmount,
        address lock
    ) private returns (uint256 agentUsed, uint256 quoteUsed, uint256 liquidity) {
        (agentUsed, quoteUsed, liquidity) = router.addLiquidity(
            agentToken,
            quoteToken,
            agentAmount,
            quoteAmount,
            (agentAmount * MIN_LIQUIDITY_BPS) / BPS_DENOMINATOR,
            (quoteAmount * MIN_LIQUIDITY_BPS) / BPS_DENOMINATOR,
            lock,
            block.timestamp
        );
    }
}
