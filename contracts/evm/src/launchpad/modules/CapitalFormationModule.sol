// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

import {IUniswapV2Pair} from "../interfaces/IUniswapV2Pair.sol";

/// @title  CapitalFormationModule
/// @notice Per-agent USDC payouts released to the agent creator when the post-
///         graduation Uniswap V2 pair's TWAP-implied FDV crosses a milestone
///         tier. Up to four tiers per agent (matching the protocol-level
///         $2M / $10M / $50M / $160M FDV ladder); each tier may be claimed
///         exactly once.
/// @dev    Single shared deployment across the agent fleet; the launchpad
///         factory owns it and binds per-agent state via {configure}. The
///         module is fully self-contained: it does not hit any external
///         oracle and does not depend on the rest of the launchpad. It pulls
///         the FDV from one place — the agent's V2 pair — and pays out from
///         a USDC balance pre-deposited by the configurer.
///
///         ## Funding model — pre-deposit
///         The module is funded via direct ERC-20 transfer of the payout
///         token at configure time. {configure}'s caller MUST `transfer`
///         the full sum of `payoutsUsdc[i]` to this contract before (or in
///         the same tx as) the configure call. We deliberately picked
///         pre-deposit over a pull-from-treasury model:
///           - one fewer external integration (no allowance dance, no
///             treasury contract dependency at module deploy time);
///           - escrow solvency is a pure on-chain invariant
///             (`balance >= sum(unclaimed payouts)` per agent ledger);
///           - permissionless {claim} can never be DoS'd by a treasury
///             rugging its allowance.
///         The trade-off is that the operator must front the worst-case
///         payout up front, but that is acceptable for a milestone-payout
///         system whose entire point is committed-stake credibility.
///
///         ## TWAP — two-step claim
///         Spot reserves are sandwich-vulnerable, so {claim} reads a TWAP
///         derived from `priceCumulativeLast` between two snapshots:
///           1. Anyone calls {snapshot(agent)} — stores
///              `(priceCumulative, blockTimestampLast, blockNumber)` from
///              the pair, capped at the current block.
///           2. After at least {MIN_TWAP_BLOCKS} new blocks AND a non-zero
///              elapsed-second window, anyone calls {claim(agent, tier)}.
///         The 32-block / ~64s window on Base raises the cost of pair
///         manipulation to ~32 blocks of held inventory, which any flash
///         loan path cannot bridge. After a successful claim the snapshot
///         is consumed; subsequent claims must take a fresh snapshot.
///
///         ## CEI + reentrancy
///         {claim} sets the tier-claimed bit BEFORE the outbound USDC
///         transfer and is wrapped in {ReentrancyGuard}. {SafeERC20} is
///         used so non-standard ERC-20s (e.g. USDC) cannot silently fail.
contract CapitalFormationModule is Ownable, ReentrancyGuard {
    using SafeERC20 for IERC20;

    // ---------------------------------------------------------------------
    // Types
    // ---------------------------------------------------------------------

    /// @notice Per-agent immutable parameters bound at {configure}.
    /// @param  agentToken     The agent ERC-20 (FDV multiplier).
    /// @param  payoutToken    Token transferred to the creator on each
    ///                        successful claim. USDC in production.
    /// @param  pair           Uniswap V2 pair created at graduation; the
    ///                        sole TWAP source.
    /// @param  agentIsToken0  Cached side of the pair to spare a `token0()`
    ///                        call on every claim.
    /// @param  creator        Recipient of every claim.
    /// @param  agentAdmin     Reserved for future admin operations (unused
    ///                        in v1; stored to mirror the LaunchRadar /
    ///                        Anti-Sniper module ABI shape).
    /// @param  configured     One-shot flag.
    /// @param  claimedBitmap  Four-bit bitmap, one bit per tier index.
    ///                        Bit `i` set => tier `i` already paid.
    struct Config {
        address agentToken;
        address payoutToken;
        address pair;
        bool agentIsToken0;
        address creator;
        address agentAdmin;
        bool configured;
        uint8 claimedBitmap;
    }

    /// @notice Per-agent TWAP cursor.
    /// @param  priceCumulative      Consolidated cumulative at snapshot
    ///                              time, cast to uint256 (UQ112x112
    ///                              integrated over seconds).
    /// @param  blockTimestampLast   Pair-reported timestamp at snapshot.
    /// @param  blockNumber          `block.number` at snapshot. Used to
    ///                              enforce the {MIN_TWAP_BLOCKS} window.
    /// @param  exists               True iff a snapshot is currently held.
    struct Snapshot {
        uint256 priceCumulative;
        uint32 blockTimestampLast;
        uint64 blockNumber;
        bool exists;
    }

    // ---------------------------------------------------------------------
    // Constants
    // ---------------------------------------------------------------------

    /// @notice Hard upper bound on the per-agent tier count. The four-bit
    ///         {Config.claimedBitmap} is sized to this.
    uint8 public constant MAX_TIERS = 4;

    /// @notice Minimum number of blocks that must elapse between a
    ///         {snapshot} and the matching {claim}. ~64 seconds on Base.
    ///         Picked to make a one-block sandwich impossible and to
    ///         price out a multi-block manipulation by requiring an
    ///         attacker to hold a lopsided pair for the full window.
    uint64 public constant MIN_TWAP_BLOCKS = 32;

    /// @dev   2**112. Uniswap V2 cumulatives are in UQ112x112 fixed point;
    ///        dividing by `Q112` returns the raw integer ratio.
    uint256 private constant Q112 = 0x10000000000000000000000000000;

    // ---------------------------------------------------------------------
    // Storage
    // ---------------------------------------------------------------------

    /// @notice Per-agent configuration. Empty until {configure}.
    mapping(address agent => Config) public configs;

    /// @notice Per-agent FDV thresholds, in `payoutToken`-equivalent units
    ///         derived from the pair quote side. Index `i` corresponds to
    ///         claimed-bit `i`. Length is bounded by {MAX_TIERS}.
    mapping(address agent => uint256[]) private _fdvThresholds;

    /// @notice Per-agent payout amounts. Parallel to {_fdvThresholds}.
    mapping(address agent => uint256[]) private _payouts;

    /// @notice Per-agent live TWAP snapshot. Consumed on successful {claim}.
    mapping(address agent => Snapshot) public snapshots;

    // ---------------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------------

    /// @dev Emitted exactly once per agent on {configure}. `totalDeposit`
    ///      is the sum of payouts the configurer must have transferred to
    ///      this contract on or before this call.
    event Configured(
        address indexed agent,
        address agentToken,
        address payoutToken,
        address pair,
        address creator,
        address agentAdmin,
        uint256 tierCount,
        uint256 totalDeposit
    );

    /// @dev Emitted on every {snapshot}. `blockNumber` is the soonest block
    ///      from which a {claim} can succeed (after {MIN_TWAP_BLOCKS} more
    ///      blocks).
    event Snapshotted(address indexed agent, uint256 priceCumulative, uint32 blockTimestampLast, uint64 blockNumber);

    /// @dev Emitted on every successful {claim}. `fdvUsd` is the TWAP-
    ///      derived FDV at claim time; `payoutUsd` is the released amount.
    event MilestoneClaimed(address indexed agent, uint8 indexed tier, uint256 fdvUsd, uint256 payoutUsd);

    // ---------------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------------

    /// @dev {configure} called twice for the same agent.
    error AlreadyConfigured();
    /// @dev A constructor / config argument was the zero address.
    error ZeroAddress();
    /// @dev `agent` was never configured.
    error NotConfigured();
    /// @dev Tier arrays were of different length.
    error LengthMismatch();
    /// @dev Tier count was zero or above {MAX_TIERS}.
    error InvalidTierCount();
    /// @dev FDV thresholds were not strictly increasing.
    error ThresholdsNotIncreasing();
    /// @dev A tier specified a zero payout.
    error ZeroPayout();
    /// @dev `agentToken_` was not one of the pair's tokens.
    error PairTokenMismatch();
    /// @dev `pair_` returned zero reserves at configure time.
    error PairUninitialized();
    /// @dev Caller did not pre-fund the configured payout sum.
    error InsufficientFunding(uint256 expected, uint256 actual);
    /// @dev Tier index was outside `[0, tierCount)`.
    error TierOutOfRange();
    /// @dev Tier was already paid.
    error TierAlreadyClaimed();
    /// @dev {claim} called without a live {snapshot}.
    error SnapshotMissing();
    /// @dev {claim} called too soon after {snapshot}; window not met.
    error TwapWindowTooShort(uint64 elapsedBlocks, uint64 required);
    /// @dev Pair did not advance its `blockTimestampLast` between snapshot
    ///      and claim — TWAP would divide by zero.
    error PairTooYoung();
    /// @dev TWAP-implied FDV is below the tier threshold.
    error FdvBelowThreshold(uint256 fdv, uint256 threshold);

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    /// @param owner_ Initial owner. Expected to be the launchpad factory
    ///               (or a Safe acting as it) so that only the factory
    ///               binds milestone schedules to live agents.
    constructor(address owner_) Ownable(_validatedOwner(owner_)) {}

    /// @dev Surface our canonical {ZeroAddress} instead of OZ's
    ///      `OwnableInvalidOwner` for missing addresses at construction.
    function _validatedOwner(address owner_) private pure returns (address) {
        if (owner_ == address(0)) revert ZeroAddress();
        return owner_;
    }

    // ---------------------------------------------------------------------
    // Admin
    // ---------------------------------------------------------------------

    /// @notice Bind an agent's milestone schedule. One-shot per agent.
    /// @dev    Caller MUST have transferred `sum(payoutsUsdc) - alreadyHeld`
    ///         of `payoutToken_` to this contract before this call (where
    ///         `alreadyHeld` is the contract's pre-existing balance net of
    ///         every other agent's outstanding obligations — but in the
    ///         common case the configurer just transfers the full sum
    ///         atomically). We do NOT track per-agent sub-balances on-chain
    ///         beyond the global `balanceOf(this)`; the invariant is:
    ///
    ///             balanceOf(this) >= sum over all agents of
    ///                 sum over unclaimed tiers of payouts[agent][tier]
    ///
    ///         Off-chain monitors enforce this. The module is owner-gated
    ///         specifically so the operator can guarantee total deposit
    ///         covers total liability across agents.
    /// @param  agent          Agent contract being configured.
    /// @param  agentToken_    The agent ERC-20.
    /// @param  payoutToken_   USDC (or whatever the protocol denominates
    ///                        milestone payouts in).
    /// @param  pair_          Uniswap V2 pair created for this agent at
    ///                        graduation. MUST already hold non-zero
    ///                        reserves on both sides.
    /// @param  creator_       Recipient of every claim.
    /// @param  agentAdmin_    Per-agent admin address. Stored for future
    ///                        ABI parity with sister modules; unused in v1.
    /// @param  fdvThresholds  Strictly increasing FDV thresholds, in pair-
    ///                        quote-token units (the configurer is
    ///                        responsible for picking quote-units that
    ///                        match how the protocol denominates "USD",
    ///                        e.g. 1:1 with USDC if pair quote is USDC,
    ///                        or pre-converted via the launch oracle if
    ///                        pair quote is TITU).
    /// @param  payoutsUsdc    Parallel array of payout amounts, in
    ///                        `payoutToken_` wei.
    function configure(
        address agent,
        address agentToken_,
        address payoutToken_,
        address pair_,
        address creator_,
        address agentAdmin_,
        uint256[] calldata fdvThresholds,
        uint256[] calldata payoutsUsdc
    ) external onlyOwner {
        if (
            agent == address(0) || agentToken_ == address(0) || payoutToken_ == address(0) || pair_ == address(0)
                || creator_ == address(0) || agentAdmin_ == address(0)
        ) revert ZeroAddress();
        if (configs[agent].configured) revert AlreadyConfigured();
        if (fdvThresholds.length != payoutsUsdc.length) revert LengthMismatch();
        if (fdvThresholds.length == 0 || fdvThresholds.length > MAX_TIERS) revert InvalidTierCount();

        bool agentIsToken0 = _validatePair(pair_, agentToken_);
        uint256 totalDeposit = _writeTiers(agent, fdvThresholds, payoutsUsdc);
        _reserveFunding(payoutToken_, totalDeposit);

        configs[agent] = Config({
            agentToken: agentToken_,
            payoutToken: payoutToken_,
            pair: pair_,
            agentIsToken0: agentIsToken0,
            creator: creator_,
            agentAdmin: agentAdmin_,
            configured: true,
            claimedBitmap: 0
        });

        emit Configured(
            agent, agentToken_, payoutToken_, pair_, creator_, agentAdmin_, fdvThresholds.length, totalDeposit
        );
    }

    /// @dev Resolve which side of the pair the agent token is on and
    ///      assert the pair has non-zero reserves. Reverts with
    ///      {PairTokenMismatch} or {PairUninitialized}. The
    ///      `blockTimestampLast` returned by `getReserves()` is intentionally
    ///      ignored — staleness is only relevant on the TWAP path. The
    ///      tuple destructure is the Uniswap V2 idiom; slither's
    ///      unused-return warning is not actionable.
    // slither-disable-next-line unused-return
    function _validatePair(address pair_, address agentToken_) private view returns (bool agentIsToken0) {
        IUniswapV2Pair pair = IUniswapV2Pair(pair_);
        address t0 = pair.token0();
        if (t0 == agentToken_) {
            agentIsToken0 = true;
        } else if (pair.token1() == agentToken_) {
            agentIsToken0 = false;
        } else {
            revert PairTokenMismatch();
        }
        (uint112 r0, uint112 r1,) = pair.getReserves();
        if (r0 == 0 || r1 == 0) revert PairUninitialized();
    }

    /// @dev Validate + persist tier arrays. Returns the sum of payouts so
    ///      the caller can reserve funding atomically.
    function _writeTiers(address agent, uint256[] calldata fdvThresholds, uint256[] calldata payoutsUsdc)
        private
        returns (uint256 totalDeposit)
    {
        // Explicit zero init silences slither's uninitialized-local warning;
        // the strictly-increasing check on i == 0 also rejects threshold == 0.
        uint256 prev = 0;
        uint256 n = fdvThresholds.length;
        for (uint256 i; i < n; ++i) {
            uint256 thr = fdvThresholds[i];
            uint256 pay = payoutsUsdc[i];
            if (i == 0) {
                if (thr == 0) revert ThresholdsNotIncreasing();
            } else if (thr <= prev) {
                revert ThresholdsNotIncreasing();
            }
            if (pay == 0) revert ZeroPayout();
            prev = thr;
            totalDeposit += pay;
            _fdvThresholds[agent].push(thr);
            _payouts[agent].push(pay);
        }
    }

    /// @dev Increment the per-token outstanding-obligation ledger after
    ///      asserting the contract holds enough free balance to cover it.
    function _reserveFunding(address payoutToken_, uint256 totalDeposit) private {
        uint256 free = IERC20(payoutToken_).balanceOf(address(this)) - _outstanding[payoutToken_];
        if (free < totalDeposit) revert InsufficientFunding(totalDeposit, free);
        _outstanding[payoutToken_] += totalDeposit;
    }

    // ---------------------------------------------------------------------
    // Snapshot / claim
    // ---------------------------------------------------------------------

    /// @notice Take a fresh TWAP snapshot for `agent`. Permissionless. The
    ///         next {claim} for this agent uses the difference between
    ///         this snapshot and a fresh on-chain read, so a {snapshot}
    ///         must precede every {claim} by at least {MIN_TWAP_BLOCKS}
    ///         blocks. Calling {snapshot} again before claiming overwrites
    ///         the prior snapshot.
    /// @dev    Reads `priceCumulativeLast` and `blockTimestampLast`
    ///         directly from the pair without extrapolating to
    ///         `block.timestamp`. We deliberately do NOT extrapolate:
    ///         consuming the pair's last-sync timestamp keeps the snapshot
    ///         math symmetric with {claim} and means a stale pair
    ///         (no swaps in the window) is detected via {PairTooYoung}
    ///         rather than silently averaging over the wrong window.
    // slither-disable-next-line unused-return
    function snapshot(address agent) external {
        Config memory cfg = configs[agent];
        if (!cfg.configured) revert NotConfigured();

        IUniswapV2Pair pair = IUniswapV2Pair(cfg.pair);
        // Reserves are intentionally dropped; only the timestamp matters
        // for snapshot freshness. See {_validatePair} for the same pattern.
        (,, uint32 ts) = pair.getReserves();
        uint256 cum = cfg.agentIsToken0 ? pair.price0CumulativeLast() : pair.price1CumulativeLast();

        snapshots[agent] =
            Snapshot({priceCumulative: cum, blockTimestampLast: ts, blockNumber: uint64(block.number), exists: true});

        emit Snapshotted(agent, cum, ts, uint64(block.number));
    }

    /// @notice Release tier `tier`'s payout to the agent's creator.
    /// @dev    Permissionless poke. Funds always route to the creator
    ///         recorded at configure time. Steps:
    ///           1. Validate config + tier index + not-already-claimed.
    ///           2. Require a live snapshot at least {MIN_TWAP_BLOCKS}
    ///              blocks old.
    ///           3. Read fresh `priceCumulative` + `blockTimestampLast`;
    ///              compute the TWAP. The pair's timestamp is used (not
    ///              `block.timestamp`) so that a pair that has not synced
    ///              since the snapshot reverts via {PairTooYoung} instead
    ///              of pretending an updated price.
    ///           4. Compute FDV = `twap * agentTotalSupply / Q112`.
    ///           5. Require FDV >= tier threshold.
    ///           6. Effects: flip claimed bit, decrement outstanding ledger,
    ///              consume snapshot.
    ///           7. Interactions: `safeTransfer` payout to creator.
    function claim(address agent, uint8 tier) external nonReentrant {
        Config storage cfg = configs[agent];
        if (!cfg.configured) revert NotConfigured();

        uint256[] storage thresholds = _fdvThresholds[agent];
        if (tier >= thresholds.length) revert TierOutOfRange();

        uint8 mask = uint8(1) << tier;
        if (cfg.claimedBitmap & mask != 0) revert TierAlreadyClaimed();

        // TWAP + FDV reads pulled into a helper to keep the stack shallow
        // here: solc 0.8.26 hits stack-too-deep otherwise.
        uint256 fdvUsd = _twapFdv(cfg, agent);

        uint256 threshold = thresholds[tier];
        if (fdvUsd < threshold) revert FdvBelowThreshold(fdvUsd, threshold);

        uint256 payout = _payouts[agent][tier];

        // --- Effects ---
        cfg.claimedBitmap |= mask;
        delete snapshots[agent];
        _outstanding[cfg.payoutToken] -= payout;

        // --- Interactions ---
        IERC20(cfg.payoutToken).safeTransfer(cfg.creator, payout);

        emit MilestoneClaimed(agent, tier, fdvUsd, payout);
    }

    /// @dev Compute the TWAP-derived FDV for `agent`. Reverts with
    ///      {SnapshotMissing}, {TwapWindowTooShort}, or {PairTooYoung}
    ///      if the snapshot/pair state cannot produce a valid TWAP. See
    ///      {claim} for the rationale on each revert path.
    /// @dev Uses the pair's own `blockTimestampLast` rather than
    ///      `block.timestamp`, mirroring the official Uniswap V2 oracle
    ///      pattern; slither flags the read but the value is consumed
    ///      only as a divisor and an explicit zero-divisor check below
    ///      handles the freshness case.
    /// @dev Multiple slither warnings are intentional and benign:
    ///      - `unused-return` on `getReserves()`: only the timestamp is
    ///        consumed here; reserves are not needed because the TWAP
    ///        comes from `priceXCumulativeLast`, not the spot ratio.
    ///      - `incorrect-equality` on `elapsedSecs == 0`: we are
    ///        deliberately gating on a precise zero-second delta to
    ///        detect a stale pair (`tsNow == s.blockTimestampLast`),
    ///        which is exactly what {PairTooYoung} signals.
    ///      - `divide-before-multiply` on the FDV math: the cumulative
    ///        delta is a UQ112x112-times-seconds quantity that exceeds
    ///        2^225, so multiplying it by `supply` (~2^90) before
    ///        dividing by `elapsedSecs` overflows uint256. The
    ///        cumulative-divided-then-multiplied ordering is the
    ///        standard Uniswap V2 oracle recipe and the precision loss
    ///        on the second step is the same fixed-point floor the
    ///        protocol accepts elsewhere.
    // slither-disable-next-line timestamp,incorrect-equality,divide-before-multiply,unused-return
    function _twapFdv(Config storage cfg, address agent) private view returns (uint256 fdvUsd) {
        Snapshot memory s = snapshots[agent];
        if (!s.exists) revert SnapshotMissing();

        uint64 elapsedBlocks = uint64(block.number) - s.blockNumber;
        if (elapsedBlocks < MIN_TWAP_BLOCKS) {
            revert TwapWindowTooShort(elapsedBlocks, MIN_TWAP_BLOCKS);
        }

        // Re-read the pair. Use its OWN last-sync timestamp; if the pair
        // has not been touched since the snapshot, ts == s.blockTimestampLast
        // and the divisor is zero → revert {PairTooYoung}.
        IUniswapV2Pair pair = IUniswapV2Pair(cfg.pair);
        (,, uint32 tsNow) = pair.getReserves();
        // uint32 wraparound is allowed by Uniswap V2 spec; mirror with
        // unchecked subtraction so a wrap during the window still yields
        // the correct delta.
        uint32 elapsedSecs;
        unchecked {
            elapsedSecs = tsNow - s.blockTimestampLast;
        }
        if (elapsedSecs == 0) revert PairTooYoung();

        uint256 cumNow = cfg.agentIsToken0 ? pair.price0CumulativeLast() : pair.price1CumulativeLast();
        // Cumulative is allowed to overflow modulo 2^256; mirror with
        // unchecked subtraction so the wrap-around is preserved.
        uint256 cumDelta;
        unchecked {
            cumDelta = cumNow - s.priceCumulative;
        }

        // priceAvgQ112 is "quote per agent" in UQ112x112.
        uint256 priceAvgQ112 = cumDelta / uint256(elapsedSecs);

        // FDV in quote-token wei = priceAvg * totalSupply.
        uint256 supply = IERC20(cfg.agentToken).totalSupply();
        fdvUsd = (priceAvgQ112 * supply) / Q112;
    }

    // ---------------------------------------------------------------------
    // Views
    // ---------------------------------------------------------------------

    /// @notice True iff tier `tier` has been claimed for `agent`.
    function isClaimed(address agent, uint8 tier) external view returns (bool) {
        if (tier >= MAX_TIERS) return false;
        return configs[agent].claimedBitmap & (uint8(1) << tier) != 0;
    }

    /// @notice Number of tiers configured for `agent`.
    function tierCount(address agent) external view returns (uint256) {
        return _fdvThresholds[agent].length;
    }

    /// @notice Read tier `tier`'s threshold + payout for `agent`.
    /// @return threshold FDV threshold in pair-quote units.
    /// @return payout    Payout amount in `payoutToken` wei.
    function tierAt(address agent, uint8 tier) external view returns (uint256 threshold, uint256 payout) {
        uint256[] storage thresholds = _fdvThresholds[agent];
        if (tier >= thresholds.length) revert TierOutOfRange();
        threshold = thresholds[tier];
        payout = _payouts[agent][tier];
    }

    // ---------------------------------------------------------------------
    // Internal — escrow ledger
    // ---------------------------------------------------------------------

    /// @dev Sum of unclaimed payouts denominated in each payout token.
    ///      Lets {configure} reject configurations that would over-commit
    ///      the on-hand balance, while still supporting many agents
    ///      sharing the same payout-token balance.
    mapping(address payoutToken => uint256) internal _outstanding;

    /// @notice Public view of the per-token escrow ledger for off-chain
    ///         solvency checks.
    function outstanding(address payoutToken) external view returns (uint256) {
        return _outstanding[payoutToken];
    }
}
