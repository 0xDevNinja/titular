// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

/// @title LaunchRadarModule
/// @notice Shared, per-agent trade-gating module. Each agent token configured on this
///         module gets its own (startTime, whitelistOn, whitelist set, agentAdmin)
///         record. Callers such as {BondingCurve} invoke {canTrade} before letting a
///         user swap: trading is blocked until `block.timestamp >= startTime` and,
///         optionally, to a per-agent allowlist of addresses.
/// @dev    One deployment per protocol, many agents. The factory calls
///         {configure} once at launch (owner-gated); after that, mutations run
///         through the agent-specific `agentAdmin` so the protocol owner does not
///         hold fine-grained per-agent keys. Unconfigured agents are permissive so
///         that this module is safe to wire into the trade path before an agent
///         has been registered (makes migrations / backfills trivial). Custom errors
///         only; no reverting string messages.
contract LaunchRadarModule is Ownable {
    // ---------------------------------------------------------------------
    // Types
    // ---------------------------------------------------------------------

    /// @notice Per-agent gating configuration.
    /// @param startTime    Unix timestamp at which trading may begin. `block.timestamp`
    ///                     strictly less than this blocks trading via {canTrade}.
    /// @param whitelistOn  When true, only addresses in {whitelist}[agent] can trade.
    /// @param configured   Set to true on first {configure}; used both to enforce
    ///                     one-shot configuration and to distinguish "unknown agent"
    ///                     (permissive) from "known agent" (gated).
    struct Config {
        uint64 startTime;
        bool whitelistOn;
        bool configured;
    }

    // ---------------------------------------------------------------------
    // Storage
    // ---------------------------------------------------------------------

    /// @notice Per-agent trade-gating record. Unset agents read as the zero struct,
    ///         which {canTrade} treats as permissive.
    mapping(address agent => Config) public configs;

    /// @notice Per-agent, per-user allowlist flag. Only consulted when the agent's
    ///         {Config.whitelistOn} is true.
    mapping(address agent => mapping(address user => bool)) public whitelist;

    /// @notice Per-agent admin authorized to mutate the agent's config and whitelist
    ///         after the protocol owner has bootstrapped the record via {configure}.
    ///         The factory is expected to wire this to a creator-controlled address
    ///         (Safe / EOA / child contract) at launch time.
    mapping(address agent => address) public agentAdmin;

    // ---------------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------------

    /// @dev Emitted once per agent on {configure}. Downstream indexers key off this
    ///      to surface a new agent's gating state.
    event Configured(address indexed agent, uint64 startTime, bool whitelistOn, address agentAdmin);

    /// @dev Emitted on every {setWhitelist} call. `count` is the number of
    ///      `(user, allowed)` pairs the call applied — indexers should re-read the
    ///      mapping rather than reconstruct state from this event alone.
    event WhitelistUpdated(address indexed agent, uint256 count);

    /// @dev Emitted on {setStartTime}. `startTime` is the post-mutation value.
    event StartTimeSet(address indexed agent, uint64 startTime);

    /// @dev Emitted on {setWhitelistOn}. `on` is the post-mutation value.
    event WhitelistToggled(address indexed agent, bool on);

    // ---------------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------------

    /// @dev Thrown when {configure} is called more than once for the same agent.
    error AlreadyConfigured();
    /// @dev Thrown when an agent-scoped mutation is attempted by someone other than
    ///      the agent's configured {agentAdmin}.
    error NotAgentAdmin();
    /// @dev Thrown when {setWhitelist} receives `users` and `allowed` arrays of
    ///      mismatching lengths.
    error LengthMismatch();
    /// @dev Thrown by {setStartTime} when the agent is already live — i.e. the
    ///      existing `startTime` has already elapsed. Prevents a surprise lockout.
    error AlreadyLive();
    /// @dev Thrown when an address argument is the zero address.
    error ZeroAddress();

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    /// @param owner_ Protocol owner authorized to call {configure}. Expected to be
    ///               the launchpad factory or a Safe multisig acting on its behalf.
    constructor(address owner_) Ownable(_validatedOwner(owner_)) {}

    /// @dev Promote OZ's `OwnableInvalidOwner` error into our canonical
    ///      {ZeroAddress} so the ABI surfaces exactly one error for missing
    ///      addresses at construction.
    function _validatedOwner(address owner_) private pure returns (address) {
        if (owner_ == address(0)) revert ZeroAddress();
        return owner_;
    }

    // ---------------------------------------------------------------------
    // Configuration
    // ---------------------------------------------------------------------

    /// @notice Register a brand-new agent on the radar. One-shot per agent: later
    ///         mutations route through {setWhitelist} / {setStartTime} /
    ///         {setWhitelistOn} and are gated on {agentAdmin}.
    /// @param  agent        Agent token address being registered. Must be non-zero.
    /// @param  startTime    Unix timestamp trading unlocks at. Passing a past time is
    ///                      allowed (e.g. backfill of an agent that is already live).
    /// @param  whitelistOn  Initial state of the whitelist gate.
    /// @param  agentAdmin_  Address authorized to mutate this agent's gating post
    ///                      bootstrap. Must be non-zero.
    function configure(address agent, uint64 startTime, bool whitelistOn, address agentAdmin_) external onlyOwner {
        if (agent == address(0) || agentAdmin_ == address(0)) revert ZeroAddress();
        Config storage cfg = configs[agent];
        if (cfg.configured) revert AlreadyConfigured();

        cfg.startTime = startTime;
        cfg.whitelistOn = whitelistOn;
        cfg.configured = true;
        agentAdmin[agent] = agentAdmin_;

        emit Configured(agent, startTime, whitelistOn, agentAdmin_);
    }

    /// @notice Batch-update the per-agent whitelist. Pairs `(users[i], allowed[i])`
    ///         are applied in order; passing `allowed=false` removes.
    /// @dev    Unbounded in array length on purpose — an agent admin is the only
    ///         party that can call this and pays the full gas. Callers are expected
    ///         to chunk updates to fit within a block's gas limit.
    function setWhitelist(address agent, address[] calldata users, bool[] calldata allowed) external {
        _onlyAgentAdmin(agent);
        uint256 n = users.length;
        if (n != allowed.length) revert LengthMismatch();
        mapping(address => bool) storage wl = whitelist[agent];
        for (uint256 i; i < n;) {
            wl[users[i]] = allowed[i];
            unchecked {
                // i is bounded by `users.length` (<= 2^256-1 in theory but <<2^64
                // by gas reality). Checked math is redundant here.
                ++i;
            }
        }
        emit WhitelistUpdated(agent, n);
    }

    /// @notice Update the agent's trading start time. Blocked once the existing
    ///         start time has already passed — at that point trading is live, and
    ///         rewriting the gate would let the admin freeze an active market.
    /// @param  newStart Post-mutation start timestamp. No lower bound: admins are
    ///                  free to bring trading forward.
    function setStartTime(address agent, uint64 newStart) external {
        _onlyAgentAdmin(agent);
        Config storage cfg = configs[agent];
        // block.timestamp is the gate being guarded; tolerance is measured in
        // hours/days so a few seconds of miner drift is immaterial.
        // slither-disable-next-line timestamp
        if (block.timestamp >= cfg.startTime) revert AlreadyLive();
        cfg.startTime = newStart;
        emit StartTimeSet(agent, newStart);
    }

    /// @notice Toggle the whitelist gate for an agent.
    function setWhitelistOn(address agent, bool on) external {
        _onlyAgentAdmin(agent);
        configs[agent].whitelistOn = on;
        emit WhitelistToggled(agent, on);
    }

    // ---------------------------------------------------------------------
    // Views
    // ---------------------------------------------------------------------

    /// @notice Trade-eligibility query used by the bonding curve / router.
    /// @dev    Permissive for unconfigured agents: we deliberately do not revert
    ///         and do not block, so this module can be wired into the trade path
    ///         before every agent has been registered. This is the documented
    ///         contract — integrators relying on positive authorization must
    ///         additionally check {configs}[agent].configured.
    /// @param  agent Agent token being traded.
    /// @param  user  Account attempting the trade.
    /// @return ok     True iff the trade is allowed right now.
    /// @return reason Short human-readable tag when `ok == false`; empty otherwise.
    function canTrade(address agent, address user) external view returns (bool ok, string memory reason) {
        Config storage cfg = configs[agent];
        if (!cfg.configured) return (true, "");
        // Time-based launch gate by design; see contract-level docs.
        // slither-disable-next-line timestamp
        if (block.timestamp < cfg.startTime) return (false, "before start");
        if (cfg.whitelistOn && !whitelist[agent][user]) return (false, "not whitelisted");
        return (true, "");
    }

    // ---------------------------------------------------------------------
    // Internal
    // ---------------------------------------------------------------------

    /// @dev Guard for per-agent mutations. Reverts {NotAgentAdmin} if the caller
    ///      is not the address registered for this agent. An unregistered agent
    ///      has `agentAdmin[agent] == address(0)`, which also fails this check.
    function _onlyAgentAdmin(address agent) private view {
        if (msg.sender != agentAdmin[agent]) revert NotAgentAdmin();
    }
}
