// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/// @title FeeRouter
/// @notice Splits incoming fees (ERC-20 or native) between an agent's creator and
///         the protocol treasury. Default split is 70% creator / 30% treasury but
///         the creator share is configurable per agent by the owner.
/// @dev    Callable by anyone: the caller funds the distribution by approving + pulling
///         ERC-20 via {distribute}, or forwarding native via {distributeNative}. The
///         router never custodies funds across transactions; every call fans out the
///         full amount in the same transaction (after rounding-down the creator side —
///         the dust stays with the treasury). CEI is enforced: state (routes) is
///         validated and splits computed up front, then effects (pulling ERC-20 from
///         the caller or accepting msg.value) settle, then outbound transfers land.
///         `ReentrancyGuard` wraps both distribute entry points because the native
///         path uses raw `call{value:}` and unknown tokens may have callbacks.
///         Immutable `treasury` is the fallback recipient when an agent has no route
///         configured (100% to treasury) and the protocol-side recipient when split.
contract FeeRouter is Ownable, ReentrancyGuard {
    using SafeERC20 for IERC20;

    // ---------------------------------------------------------------------
    // Types
    // ---------------------------------------------------------------------

    /// @notice Route metadata for a single agent token.
    /// @param creator     Recipient of the creator share of fees.
    /// @param creatorBps  Creator share in basis points (0..BPS_DENOMINATOR).
    /// @param configured  True once {setRoute} has been called for this agent.
    struct Route {
        address creator;
        uint16 creatorBps;
        bool configured;
    }

    // ---------------------------------------------------------------------
    // Constants / immutables
    // ---------------------------------------------------------------------

    /// @notice Basis-point denominator used throughout split math.
    uint16 public constant BPS_DENOMINATOR = 10_000;

    /// @notice Default creator share in basis points. 7000 = 70%.
    uint16 public constant DEFAULT_CREATOR_BPS = 7000;

    /// @notice Treasury address — protocol-side recipient. Fixed at construction.
    address public immutable treasury;

    // ---------------------------------------------------------------------
    // Storage
    // ---------------------------------------------------------------------

    /// @notice Per-agent routing table. Key = agent ERC-20 address.
    mapping(address agent => Route) public routes;

    // ---------------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------------

    /// @dev Emitted when the owner configures (or reconfigures) an agent's split.
    event RouteSet(address indexed agent, address indexed creator, uint16 creatorBps);

    /// @dev Emitted when the owner clears an agent's route. Subsequent distributions
    ///      fall back to sending 100% to the treasury.
    event RouteCleared(address indexed agent);

    /// @dev Emitted on every successful distribution. `token` is `address(0)` for
    ///      native distributions.
    event Distributed(address indexed agent, address indexed token, uint256 creatorAmount, uint256 treasuryAmount);

    // ---------------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------------

    /// @dev Thrown on any zero-address argument where a non-zero is required.
    error ZeroAddress();
    /// @dev Thrown when `creatorBps > BPS_DENOMINATOR`. Carries the offending value.
    error BpsTooHigh(uint16 creatorBps);
    /// @dev Thrown when a native transfer (via `call{value:}`) fails. Carries the
    ///      recipient that rejected.
    error NativeTransferFailed(address recipient);
    /// @dev Thrown by {distributeNative} when `msg.value == 0`.
    error InsufficientValue();

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    /// @param treasury_  Protocol treasury address. Non-zero.
    /// @param owner_     Initial owner; gates {setRoute} and {clearRoute}.
    constructor(address treasury_, address owner_) Ownable(_validated(owner_)) {
        if (treasury_ == address(0)) revert ZeroAddress();
        treasury = treasury_;
    }

    /// @dev Hoists the zero-address check above OZ Ownable so callers see our
    ///      canonical {ZeroAddress} error instead of OZ's `OwnableInvalidOwner`.
    function _validated(address addr) private pure returns (address) {
        if (addr == address(0)) revert ZeroAddress();
        return addr;
    }

    // ---------------------------------------------------------------------
    // Admin
    // ---------------------------------------------------------------------

    /// @notice Configure (or reconfigure) the fee split for `agent`.
    /// @param  agent        Agent ERC-20 whose fees are being routed.
    /// @param  creator      Creator recipient of the creator share. Non-zero.
    /// @param  creatorBps   Creator share in basis points; must be ≤ BPS_DENOMINATOR.
    ///                      Treasury receives `BPS_DENOMINATOR - creatorBps`.
    /// @dev    Setting `creatorBps = 0` is legal (all to treasury). Setting
    ///         `creatorBps = BPS_DENOMINATOR` is legal (all to creator). Owner is
    ///         expected to be a Safe multisig.
    function setRoute(address agent, address creator, uint16 creatorBps) external onlyOwner {
        if (agent == address(0) || creator == address(0)) revert ZeroAddress();
        if (creatorBps > BPS_DENOMINATOR) revert BpsTooHigh(creatorBps);

        routes[agent] = Route({creator: creator, creatorBps: creatorBps, configured: true});
        emit RouteSet(agent, creator, creatorBps);
    }

    /// @notice Clear an agent's route. Subsequent distributions send 100% to treasury.
    function clearRoute(address agent) external onlyOwner {
        if (agent == address(0)) revert ZeroAddress();
        delete routes[agent];
        emit RouteCleared(agent);
    }

    // ---------------------------------------------------------------------
    // Distribution
    // ---------------------------------------------------------------------

    /// @notice Pull `amount` of `token` from the caller and fan it out per `agent`'s route.
    /// @dev    CEI:
    ///           1. Compute split from route (or 100% to treasury if unconfigured).
    ///           2. Pull `amount` from `msg.sender`.
    ///           3. Push creator share then treasury share.
    ///         The caller must have approved `amount` to this contract. The router
    ///         never retains an ERC-20 balance across calls.
    /// @param  agent  Agent key for route lookup.
    /// @param  token  ERC-20 being distributed. Non-zero.
    /// @param  amount Amount to distribute. If zero, the call is a no-op revert
    ///                on SafeERC20 for most tokens; callers should size > 0.
    function distribute(address agent, IERC20 token, uint256 amount) external nonReentrant {
        if (address(token) == address(0)) revert ZeroAddress();

        (address creator, uint256 creatorAmount, uint256 treasuryAmount) = _split(agent, amount);

        // --- Interactions ---
        token.safeTransferFrom(msg.sender, address(this), amount);
        if (creatorAmount != 0) token.safeTransfer(creator, creatorAmount);
        if (treasuryAmount != 0) token.safeTransfer(treasury, treasuryAmount);

        emit Distributed(agent, address(token), creatorAmount, treasuryAmount);
    }

    /// @notice Split the attached `msg.value` per `agent`'s route.
    /// @dev    Uses `call{value:}` for both legs to remain compatible with contract
    ///         recipients (treasury is typically a Safe). Reverts if either transfer
    ///         fails — no dust is retained on the router. CEI is trivial here because
    ///         state is only read.
    /// @param  agent Agent key for route lookup.
    function distributeNative(address agent) external payable nonReentrant {
        if (msg.value == 0) revert InsufficientValue();

        (address creator, uint256 creatorAmount, uint256 treasuryAmount) = _split(agent, msg.value);

        // --- Interactions ---
        if (creatorAmount != 0) {
            // slither-disable-next-line arbitrary-send-eth
            (bool ok,) = payable(creator).call{value: creatorAmount}("");
            if (!ok) revert NativeTransferFailed(creator);
        }
        if (treasuryAmount != 0) {
            // slither-disable-next-line arbitrary-send-eth
            (bool ok,) = payable(treasury).call{value: treasuryAmount}("");
            if (!ok) revert NativeTransferFailed(treasury);
        }

        emit Distributed(agent, address(0), creatorAmount, treasuryAmount);
    }

    // ---------------------------------------------------------------------
    // Views
    // ---------------------------------------------------------------------

    /// @notice Preview the split `amount` would produce for `agent`.
    /// @dev    If `agent` has no route configured, `creatorAmount == 0` and the full
    ///         `amount` is attributed to the treasury (matches {distribute} behaviour).
    /// @return creatorAmount  Portion routed to the creator (0 if unconfigured).
    /// @return treasuryAmount Portion routed to the treasury. Includes any rounding dust.
    function getSplit(address agent, uint256 amount)
        external
        view
        returns (uint256 creatorAmount, uint256 treasuryAmount)
    {
        (, creatorAmount, treasuryAmount) = _split(agent, amount);
    }

    // ---------------------------------------------------------------------
    // Internal
    // ---------------------------------------------------------------------

    /// @dev Compute the `(creator, creatorAmount, treasuryAmount)` tuple for a given
    ///      `(agent, amount)`. Rounding: the creator share is floored via integer
    ///      division; the treasury receives the remainder (`amount - creatorAmount`).
    ///      This guarantees the two parts sum to exactly `amount`, and prevents dust
    ///      from escaping the router on-chain. If the agent has no route, the creator
    ///      tuple element is returned as `address(0)` and `creatorAmount == 0`; callers
    ///      must treat `address(0)` as "no creator leg".
    function _split(address agent, uint256 amount)
        internal
        view
        returns (address creator, uint256 creatorAmount, uint256 treasuryAmount)
    {
        Route storage r = routes[agent];
        if (!r.configured) {
            // No route: entire amount to treasury.
            return (address(0), 0, amount);
        }
        creator = r.creator;
        // Creator floored; treasury gets the remainder including dust.
        creatorAmount = (amount * r.creatorBps) / BPS_DENOMINATOR;
        treasuryAmount = amount - creatorAmount;
    }
}
