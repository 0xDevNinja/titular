// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/// @title  SixtyDaysModule
/// @notice Reversible-launch escrow shared across every agent launch. The 70% creator
///         leg of the swap fee (routed by the launchpad's `FeeRouter`) is held in
///         this module for the first 60 days of an agent's life. During that window
///         the agent's `creator` (registered at configure-time) has exactly one
///         binary decision to make:
///             * `commit()` — sealed only after `startTime + WINDOW`. Unlocks the
///               full escrow balance to the creator. One-shot, irreversible.
///             * `refund()` — sealed before `commit()` and never after. Opens a
///               pull-pattern claim phase: each contributor (per the per-buyer
///               ledger pushed in by the bonding curve) can claim a pro-rata share
///               of the escrow plus any pre-graduation reserves the curve drains
///               into the module via `accrueRefundPool`. One-shot, irreversible.
/// @dev    Single deployment per protocol; per-agent state lives in `configs`.
///         The protocol owner (typically the launchpad factory or a Safe) calls
///         `configure` once per agent to bootstrap. After that the only privileged
///         actors per agent are:
///             * the registered `bondingCurve` — pushes contribution ledger entries
///               via {accrueEscrow}+{recordContribution} during the trade path, and
///               drains residual reserves via {accrueRefundPool} on `refund()`.
///             * the registered `creator` — calls {commit} or {refund}.
///         CEI is enforced on every state-mutating external entry that moves
///         tokens. `claimed[agent][user]` flips before the outbound transfer so a
///         re-entrant ERC-20 cannot double-pay; `ReentrancyGuard` layers on top.
///         Only custom errors; no string reverts. Funds are never permanently
///         locked: the module exposes either an unlock-to-creator path or an
///         unlock-to-buyers path, and both phases are reachable from any state
///         the module is ever placed in by an honest caller.
contract SixtyDaysModule is Ownable, ReentrancyGuard {
    using SafeERC20 for IERC20;

    // ---------------------------------------------------------------------
    // Constants
    // ---------------------------------------------------------------------

    /// @notice Reversibility window measured from `startTime`. Hard-coded to 60
    ///         days so it cannot be lengthened per-agent post-bootstrap.
    uint64 public constant WINDOW = 60 days;

    /// @notice Phase enum encoded as `uint8` to keep `Config` packed.
    uint8 internal constant PHASE_OPEN = 0;
    uint8 internal constant PHASE_COMMITTED = 1;
    uint8 internal constant PHASE_REFUNDING = 2;

    // ---------------------------------------------------------------------
    // Types
    // ---------------------------------------------------------------------

    /// @notice Per-agent escrow + lifecycle state.
    /// @param startTime                Unix seconds the 60-day clock starts at.
    /// @param windowEnd                `startTime + WINDOW`. `commit` is gated on
    ///                                 `block.timestamp >= windowEnd`.
    /// @param escrowToken              Token escrowed (typically the curve's quote
    ///                                 asset, TITU). Same asset is used to pay out
    ///                                 commit + refund.
    /// @param bondingCurve             Sole address allowed to push escrow accruals,
    ///                                 contribution-ledger entries, and refund-pool
    ///                                 top-ups for this agent.
    /// @param creator                  Sole address allowed to call {commit} or
    ///                                 {refund}.
    /// @param phase                    {PHASE_OPEN | PHASE_COMMITTED | PHASE_REFUNDING}.
    /// @param configured               One-shot guard; flipped on first `configure`.
    /// @param escrowAmountAccumulated  Running sum of escrowed token wei pushed in by
    ///                                 the curve (creator-leg fees). On commit this
    ///                                 is paid out gross to the creator. On refund
    ///                                 it is added to the pull-pattern pool.
    /// @param refundPool               Total claimable wei after `refund()` has
    ///                                 fired: equals `escrowAmountAccumulated` plus
    ///                                 any extra reserves the curve drains via
    ///                                 {accrueRefundPool}. Pro-rata denominator on
    ///                                 the contributor-ledger side stays
    ///                                 `totalContributed`.
    /// @param totalContributed         Sum of buyer-side quote contributions during
    ///                                 the open phase. Used as the pro-rata
    ///                                 denominator on {claimRefund}.
    struct Config {
        uint64 startTime;
        uint64 windowEnd;
        address escrowToken;
        address bondingCurve;
        address creator;
        uint8 phase;
        bool configured;
        uint256 escrowAmountAccumulated;
        uint256 refundPool;
        uint256 totalContributed;
    }

    // ---------------------------------------------------------------------
    // Storage
    // ---------------------------------------------------------------------

    /// @notice Per-agent configuration + accumulated state.
    mapping(address agent => Config) public configs;

    /// @notice Per-agent, per-user cumulative buy-side quote contribution during
    ///         the open phase. Pro-rata numerator on {claimRefund}.
    mapping(address agent => mapping(address user => uint256)) public contributed;

    /// @notice Per-agent, per-user one-shot claim flag. Set BEFORE the outbound
    ///         token transfer in {claimRefund} (CEI). Once true, the user cannot
    ///         claim again for this agent.
    mapping(address agent => mapping(address user => bool)) public claimed;

    // ---------------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------------

    /// @notice Emitted exactly once when the protocol owner registers an agent.
    event Configured(
        address indexed agent,
        address indexed escrowToken,
        address indexed bondingCurve,
        address creator,
        uint64 startTime,
        uint64 windowEnd
    );

    /// @notice Emitted on every {accrueEscrow} push from the bonding curve.
    event EscrowAccrued(address indexed agent, uint256 amount, uint256 totalEscrow);

    /// @notice Emitted on every {recordContribution} push from the bonding curve.
    event Contributed(address indexed agent, address indexed user, uint256 amount, uint256 totalContributed);

    /// @notice Emitted on every {accrueRefundPool} top-up after refund has fired.
    event RefundPoolToppedUp(address indexed agent, uint256 amount, uint256 totalPool);

    /// @notice Emitted exactly once per agent when {commit} succeeds.
    event Committed(address indexed agent, address indexed creator, uint256 amount);

    /// @notice Emitted exactly once per agent when {refund} succeeds.
    event Refunded(address indexed agent, address indexed creator, uint256 escrowAmount);

    /// @notice Emitted on every successful {claimRefund}.
    event RefundClaimed(address indexed agent, address indexed user, uint256 amount);

    // ---------------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------------

    /// @dev Thrown when `configure` is called more than once for the same agent.
    error AlreadyConfigured();
    /// @dev Thrown by accrual / contribution paths when the caller is not the
    ///      bonding curve registered for the agent in {configure}.
    error NotBondingCurve();
    /// @dev Thrown by {commit} / {refund} when the caller is not the registered
    ///      creator.
    error NotCreator();
    /// @dev Thrown by {commit} when invoked before `startTime + WINDOW`.
    error WindowNotElapsed();
    /// @dev Thrown by {commit} / {refund} / accrual paths when the agent is not
    ///      in {PHASE_OPEN} (the only state from which commit or refund can fire,
    ///      and from which the curve can keep accruing).
    error NotOpen();
    /// @dev Thrown by {accrueRefundPool} when the agent is not in
    ///      {PHASE_REFUNDING}.
    error NotRefunding();
    /// @dev Thrown by {claimRefund} when the agent has not yet entered
    ///      {PHASE_REFUNDING}.
    error NotInRefundPhase();
    /// @dev Thrown by {claimRefund} when the caller has already claimed.
    error AlreadyClaimed();
    /// @dev Thrown by {claimRefund} when the caller never contributed (or their
    ///      pro-rata share floors to zero).
    error NoContribution();
    /// @dev Thrown when an `address` argument required to be non-zero is zero.
    error ZeroAddress();
    /// @dev Thrown when a `uint256` amount required to be non-zero is zero.
    error ZeroAmount();

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    /// @param owner_ Protocol owner authorized to call {configure}. Expected to
    ///               be the launchpad factory or a Safe multisig acting on its
    ///               behalf. Reverts via OZ Ownable's `OwnableInvalidOwner` if
    ///               zero — surfaced as our canonical {ZeroAddress} below.
    constructor(address owner_) Ownable(_validatedOwner(owner_)) {}

    /// @dev Promote OZ's `OwnableInvalidOwner` into our canonical {ZeroAddress}
    ///      so the ABI exposes a single error type for missing addresses at
    ///      construction.
    function _validatedOwner(address owner_) private pure returns (address) {
        if (owner_ == address(0)) revert ZeroAddress();
        return owner_;
    }

    // ---------------------------------------------------------------------
    // Configuration
    // ---------------------------------------------------------------------

    /// @notice Bootstrap an agent on this module. One-shot per agent.
    /// @dev    Owner-gated. The factory wires the agent's bonding curve as the
    ///         sole accrual / contribution source and the creator as the sole
    ///         party authorized to call {commit} or {refund}.
    /// @param  agent          Agent ERC-20 address (storage key).
    /// @param  escrowToken_   Token escrowed (creator-leg fee currency).
    /// @param  bondingCurve_  Curve allowed to push accruals + contributions.
    /// @param  creator_       Creator allowed to commit or refund.
    /// @param  startTime_     Unix seconds the 60-day clock starts at.
    function configure(address agent, address escrowToken_, address bondingCurve_, address creator_, uint64 startTime_)
        external
        onlyOwner
    {
        if (agent == address(0) || escrowToken_ == address(0) || bondingCurve_ == address(0) || creator_ == address(0))
        {
            revert ZeroAddress();
        }
        Config storage cfg = configs[agent];
        if (cfg.configured) revert AlreadyConfigured();

        uint64 windowEnd_ = startTime_ + WINDOW;
        cfg.startTime = startTime_;
        cfg.windowEnd = windowEnd_;
        cfg.escrowToken = escrowToken_;
        cfg.bondingCurve = bondingCurve_;
        cfg.creator = creator_;
        cfg.configured = true;
        cfg.phase = PHASE_OPEN;

        emit Configured(agent, escrowToken_, bondingCurve_, creator_, startTime_, windowEnd_);
    }

    // ---------------------------------------------------------------------
    // Curve hooks (open phase)
    // ---------------------------------------------------------------------

    /// @notice Push the creator-leg fee for `agent` into escrow. Pulled from the
    ///         curve via {SafeERC20-safeTransferFrom}; the curve must have
    ///         approved the module for `amount` of `escrowToken`.
    /// @dev    Curve-gated. Allowed only while the agent is in {PHASE_OPEN};
    ///         after commit / refund the escrow accounting is frozen and any
    ///         further fees should be routed elsewhere by the integrator.
    ///         CEI: the running counter updates before the inbound pull so a
    ///         malicious token cannot reenter and double-count itself; the
    ///         outbound transfer is the inbound pull (no further hops).
    /// @param  agent   Agent the escrow is for.
    /// @param  amount  Amount of `escrowToken` to pull. Non-zero.
    function accrueEscrow(address agent, uint256 amount) external nonReentrant {
        if (amount == 0) revert ZeroAmount();
        Config storage cfg = configs[agent];
        if (msg.sender != cfg.bondingCurve) revert NotBondingCurve();
        if (cfg.phase != PHASE_OPEN) revert NotOpen();

        uint256 newTotal = cfg.escrowAmountAccumulated + amount;
        cfg.escrowAmountAccumulated = newTotal;

        IERC20(cfg.escrowToken).safeTransferFrom(msg.sender, address(this), amount);
        emit EscrowAccrued(agent, amount, newTotal);
    }

    /// @notice Record `user`'s buy-side contribution to `agent`'s curve.
    /// @dev    Curve-gated. Used as the pro-rata numerator on {claimRefund}.
    ///         No tokens move here — the curve already holds the buyer's quote
    ///         in its reserves; this is pure accounting.
    /// @param  agent   Agent the contribution is for.
    /// @param  user    Buyer credited.
    /// @param  amount  Net quote (post-fee) contributed by `user`.
    function recordContribution(address agent, address user, uint256 amount) external {
        if (amount == 0) revert ZeroAmount();
        if (user == address(0)) revert ZeroAddress();
        Config storage cfg = configs[agent];
        if (msg.sender != cfg.bondingCurve) revert NotBondingCurve();
        if (cfg.phase != PHASE_OPEN) revert NotOpen();

        contributed[agent][user] += amount;
        cfg.totalContributed += amount;
        emit Contributed(agent, user, amount, cfg.totalContributed);
    }

    // ---------------------------------------------------------------------
    // Creator decisions
    // ---------------------------------------------------------------------

    /// @notice Unlock the full escrow balance for `agent` to the registered
    ///         creator. Callable only after `startTime + WINDOW` and only when
    ///         the agent is still in {PHASE_OPEN} (i.e. the creator has not
    ///         already chosen `refund`).
    /// @dev    One-shot: the phase is flipped to {PHASE_COMMITTED} BEFORE the
    ///         outbound transfer (CEI). `nonReentrant` belt-and-braces.
    ///         Idempotency: a second call reverts {NotOpen}.
    /// @param  agent  Agent committing.
    function commit(address agent) external nonReentrant {
        Config storage cfg = configs[agent];
        if (msg.sender != cfg.creator) revert NotCreator();
        if (cfg.phase != PHASE_OPEN) revert NotOpen();
        // Time gate: 60 days from `startTime` must have elapsed. The granularity
        // here is days — miner timestamp drift of seconds is irrelevant.
        // slither-disable-next-line timestamp
        if (block.timestamp < cfg.windowEnd) revert WindowNotElapsed();

        uint256 amount = cfg.escrowAmountAccumulated;
        cfg.phase = PHASE_COMMITTED;

        if (amount != 0) {
            IERC20(cfg.escrowToken).safeTransfer(cfg.creator, amount);
        }
        emit Committed(agent, cfg.creator, amount);
    }

    /// @notice Open the pull-pattern refund phase for `agent`. Callable by the
    ///         registered creator any time before {commit} fires — including
    ///         strictly before `startTime + WINDOW`. The full escrow balance is
    ///         seeded into the refund pool; the bonding curve may also drain
    ///         residual pre-graduation reserves into the pool via
    ///         {accrueRefundPool}.
    /// @dev    One-shot: the phase flips to {PHASE_REFUNDING} and any later
    ///         {commit} or {refund} call reverts {NotOpen}. No tokens leave the
    ///         module on this call — the escrow already lives here. Buyers pull
    ///         their share via {claimRefund}.
    /// @param  agent  Agent opening refund.
    function refund(address agent) external {
        Config storage cfg = configs[agent];
        if (msg.sender != cfg.creator) revert NotCreator();
        if (cfg.phase != PHASE_OPEN) revert NotOpen();

        uint256 escrowAmount = cfg.escrowAmountAccumulated;
        cfg.phase = PHASE_REFUNDING;
        cfg.refundPool = escrowAmount;

        emit Refunded(agent, cfg.creator, escrowAmount);
    }

    // ---------------------------------------------------------------------
    // Curve hooks (refund phase)
    // ---------------------------------------------------------------------

    /// @notice Top up the refund pool for `agent` from the bonding curve.
    /// @dev    Curve-gated. Used by the curve to drain residual quote reserves
    ///         into the module after {refund} has fired so that buyers receive
    ///         not just the escrowed creator-leg fees but also their share of
    ///         the unspent curve reserves. The curve must have approved the
    ///         module for `amount`. Allowed only in {PHASE_REFUNDING}.
    ///         CEI: the running counter updates before the inbound pull.
    /// @param  agent   Agent receiving the top-up.
    /// @param  amount  Amount of `escrowToken` pulled from the curve. Non-zero.
    function accrueRefundPool(address agent, uint256 amount) external nonReentrant {
        if (amount == 0) revert ZeroAmount();
        Config storage cfg = configs[agent];
        if (msg.sender != cfg.bondingCurve) revert NotBondingCurve();
        if (cfg.phase != PHASE_REFUNDING) revert NotRefunding();

        uint256 newTotal = cfg.refundPool + amount;
        cfg.refundPool = newTotal;

        IERC20(cfg.escrowToken).safeTransferFrom(msg.sender, address(this), amount);
        emit RefundPoolToppedUp(agent, amount, newTotal);
    }

    /// @notice Pull the caller's pro-rata share of the refund pool for `agent`.
    /// @dev    CEI: `claimed` is set BEFORE the outbound transfer and a second
    ///         call from the same address reverts {AlreadyClaimed}.
    ///         `nonReentrant` belt-and-braces against ERC-20 callbacks.
    ///         Pro-rata: `share = userContribution * refundPool / totalContributed`,
    ///         floored. The dust left by the floor stays in the contract until
    ///         the last claimer takes a slightly-truncated share (no recipient
    ///         is favoured; the floor strictly under-pays so the module remains
    ///         solvent against `refundPool`).
    /// @param  agent  Agent the caller is claiming against.
    /// @return amount Tokens transferred to the caller (>= 1 wei or revert).
    function claimRefund(address agent) external nonReentrant returns (uint256 amount) {
        Config storage cfg = configs[agent];
        if (cfg.phase != PHASE_REFUNDING) revert NotInRefundPhase();
        if (claimed[agent][msg.sender]) revert AlreadyClaimed();

        uint256 userContribution = contributed[agent][msg.sender];
        if (userContribution == 0) revert NoContribution();

        // `totalContributed` is non-zero whenever any user has a non-zero
        // `userContribution` because the latter only grows under the same
        // {recordContribution} path that grows the former by the same delta.
        amount = (userContribution * cfg.refundPool) / cfg.totalContributed;
        if (amount == 0) revert NoContribution();

        claimed[agent][msg.sender] = true;

        IERC20(cfg.escrowToken).safeTransfer(msg.sender, amount);
        emit RefundClaimed(agent, msg.sender, amount);
    }

    // ---------------------------------------------------------------------
    // Views
    // ---------------------------------------------------------------------

    /// @notice Read-only twin of {claimRefund}. Returns zero for users who have
    ///         already claimed, never contributed, or whose agent is not in the
    ///         refund phase.
    /// @param  agent  Agent to preview against.
    /// @param  user   Address to preview for.
    /// @return amount Tokens the user would receive if they claimed now.
    function previewRefund(address agent, address user) external view returns (uint256 amount) {
        Config storage cfg = configs[agent];
        if (cfg.phase != PHASE_REFUNDING) return 0;
        if (claimed[agent][user]) return 0;
        uint256 userContribution = contributed[agent][user];
        if (userContribution == 0) return 0;
        amount = (userContribution * cfg.refundPool) / cfg.totalContributed;
    }

    /// @notice Convenience accessor for the agent's current phase.
    /// @return phase 0 = open, 1 = committed, 2 = refunding.
    function phaseOf(address agent) external view returns (uint8 phase) {
        return configs[agent].phase;
    }
}
