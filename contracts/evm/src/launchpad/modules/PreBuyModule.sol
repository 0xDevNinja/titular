// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {Clones} from "@openzeppelin/contracts/proxy/Clones.sol";

import {IBondingCurve} from "../interfaces/IBondingCurve.sol";
import {TokenVesting} from "./TokenVesting.sol";

/// @title  PreBuyModule
/// @notice Per-agent creator pre-buy with custom vesting. The launchpad
///         factory funds this module with `titanIn` quote (TITU) and calls
///         {configure}; the module routes the buy through the agent's bonding
///         curve, then drops the bought agent tokens into a fresh
///         EIP-1167 minimal-proxy {TokenVesting} clone keyed by the agent
///         address. The clone vests linearly to the creator after a
///         configurable cliff. One-shot per agent.
/// @dev    Single deployment per protocol, many agents. Owner is expected to
///         be the launchpad factory (or the Safe acting as it). The vesting
///         implementation is deployed once at construction; per-agent state
///         lives entirely on cheap minimal-proxy clones. CEI is enforced on
///         {configure}: state writes (config + clone deploy) land before any
///         external buy / transfer / init call, and {ReentrancyGuard} closes
///         the configure path against malicious tokens or curves.
///
///         Funding model: this module never custodies TITU. The owner
///         transfers exactly `titanIn` TITU to this contract immediately
///         before each {configure} call (or atomically in the same tx), and
///         the module forwards it to the curve. Any leftover quote on the
///         contract is swept by the owner via {sweep} so an aborted call
///         cannot strand funds.
contract PreBuyModule is Ownable, ReentrancyGuard {
    using SafeERC20 for IERC20;
    using Clones for address;

    // ---------------------------------------------------------------------
    // Types
    // ---------------------------------------------------------------------

    /// @notice Per-agent pre-buy record. Written exactly once on {configure}.
    /// @param  creator       Beneficiary of the vesting clone.
    /// @param  vestingClone  EIP-1167 minimal proxy for {TokenVesting}.
    /// @param  vestAmount    Agent tokens transferred into the clone (= curve
    ///                       output, with slippage floor `vestAmount`).
    /// @param  cliffSeconds  Cliff measured from {start}, in seconds.
    /// @param  durationSeconds Total vesting window in seconds.
    /// @param  start         Unix timestamp at which the vesting curve begins.
    /// @param  configured    Sentinel set true on first {configure}.
    struct Config {
        address creator;
        address vestingClone;
        uint256 vestAmount;
        uint64 cliffSeconds;
        uint64 durationSeconds;
        uint64 start;
        bool configured;
    }

    // ---------------------------------------------------------------------
    // Immutables
    // ---------------------------------------------------------------------

    /// @notice {TokenVesting} implementation pointed at by every clone.
    address public immutable vestingImplementation;

    // ---------------------------------------------------------------------
    // Storage
    // ---------------------------------------------------------------------

    /// @notice Per-agent pre-buy state. Empty until the factory calls
    ///         {configure} for that agent.
    mapping(address agent => Config) public configs;

    // ---------------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------------

    /// @dev Emitted exactly once per agent on {configure}. Indexers key on
    ///      `agent` to surface the creator's vesting bag and clone address.
    event Configured(
        address indexed agent,
        address indexed creator,
        address indexed vestingClone,
        uint256 vestAmount,
        uint64 cliffSeconds,
        uint64 durationSeconds
    );

    /// @dev Mirror of the clone's own {TokenVesting.Released} event so a
    ///      single subscription on this module is enough to follow every
    ///      creator's pull. Anyone may trigger {release} (gas arbitrage).
    event Released(address indexed agent, address indexed clone, uint256 amount);

    /// @dev Emitted when the owner sweeps stranded quote off the module.
    event QuoteSwept(address indexed token, address indexed to, uint256 amount);

    // ---------------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------------

    /// @dev Thrown when an address argument is the zero address.
    error ZeroAddress();
    /// @dev Thrown when `vestAmount`, `titanIn`, or `durationSeconds` is zero.
    error ZeroAmount();
    /// @dev Thrown when `cliffSeconds > durationSeconds`.
    error CliffExceedsDuration();
    /// @dev Thrown when {configure} is called twice for the same agent.
    error AlreadyConfigured();
    /// @dev Thrown when the curve's `agentToken()` does not match the `agent`
    ///      argument passed to {configure}. Catches a misconfigured factory
    ///      from binding the wrong curve to an agent.
    error AgentMismatch();
    /// @dev Thrown when the curve actually delivered fewer tokens than the
    ///      caller requested as the vest target. The curve's own slippage
    ///      check should pre-empt this, but we re-assert before the clone is
    ///      funded so a buggy curve cannot under-fund the grant.
    error InsufficientBuyOutput();
    /// @dev Thrown when {release} is called for an agent that has not been
    ///      registered via {configure}.
    error NotConfigured();

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    /// @param owner_ Initial owner. Expected to be the launchpad factory or a
    ///               Safe multisig acting on its behalf — only that party is
    ///               trusted to bind a creator's pre-buy bag.
    constructor(address owner_) Ownable(_validatedOwner(owner_)) {
        // Deploy the immutable vesting implementation once. Every per-agent
        // clone is a 45-byte EIP-1167 minimal proxy delegating to this code,
        // keeping per-agent deployment cost ~32k gas vs. ~600k for a fresh
        // copy.
        vestingImplementation = address(new TokenVesting());
    }

    /// @dev Wrap OZ Ownable's zero-owner check so callers see our canonical
    ///      {ZeroAddress} error instead of OZ's `OwnableInvalidOwner`.
    function _validatedOwner(address owner_) private pure returns (address) {
        if (owner_ == address(0)) revert ZeroAddress();
        return owner_;
    }

    // ---------------------------------------------------------------------
    // Configure
    // ---------------------------------------------------------------------

    /// @notice Bind a pre-buy + vesting plan to `agent`. Owner-only, one-shot.
    /// @dev    Pre-conditions:
    ///           - `titanIn` TITU has been transferred into this contract by
    ///             the owner (or atomically in the same tx).
    ///           - `curve.agentToken() == agent` (asserted).
    ///         Effects:
    ///           1. Validate args + look up curve token bindings.
    ///           2. Deploy a deterministic {TokenVesting} clone (salt =
    ///              keccak256(agent)). Storage write (`configs[agent]`)
    ///              lands before any external call so a malicious token cannot
    ///              re-enter into a half-configured agent.
    ///           3. Approve the curve for exactly `titanIn`, call `buy`, then
    ///              clear the approval — leaves no standing allowance on a
    ///              long-lived module.
    ///           4. Re-read this contract's agent-token balance, assert
    ///              `received >= vestAmount`, transfer the actually-received
    ///              amount into the clone, initialize the clone with the
    ///              vesting schedule.
    ///         Reverts (any) leave the module's storage clean and the TITU
    ///         pull-through unspent (owner can {sweep} it).
    /// @param  agent           Agent token (clone key).
    /// @param  creator         Creator address — sole vesting beneficiary.
    /// @param  vestAmount      Slippage floor on agent tokens out of the buy
    ///                         AND nominal grant size. Curve will deliver at
    ///                         least this, often slightly more (curve does
    ///                         not refund excess); the module forwards the
    ///                         entire received balance into the clone.
    /// @param  cliffSeconds    Vesting cliff offset (seconds).
    /// @param  durationSeconds Total vesting window (seconds). Non-zero.
    /// @param  titanIn         Quote (TITU) the curve will pull from this
    ///                         module. Must already be on this contract.
    /// @param  curve           Bonding curve for `agent`.
    /// @return clone           Address of the newly-deployed vesting clone.
    /// @dev    {nonReentrant} closes the curve / token re-entry surface; the
    ///         pre-/post-call balance read is safe because the call is wrapped
    ///         in the guard (slither's reentrancy-balance is a false positive
    ///         in that context).
    // slither-disable-next-line reentrancy-no-eth,reentrancy-benign,reentrancy-events,reentrancy-balance
    function configure(
        address agent,
        address creator,
        uint256 vestAmount,
        uint64 cliffSeconds,
        uint64 durationSeconds,
        uint256 titanIn,
        IBondingCurve curve
    ) external onlyOwner nonReentrant returns (address clone) {
        // -------------------- Validation --------------------
        if (agent == address(0) || creator == address(0) || address(curve) == address(0)) revert ZeroAddress();
        if (vestAmount == 0 || titanIn == 0) revert ZeroAmount();
        if (durationSeconds == 0) revert ZeroAmount();
        if (cliffSeconds > durationSeconds) revert CliffExceedsDuration();
        if (configs[agent].configured) revert AlreadyConfigured();
        if (curve.agentToken() != agent) revert AgentMismatch();

        IERC20 quote = IERC20(curve.quoteToken());
        IERC20 agentToken = IERC20(agent);

        // -------------------- Effects (clone deploy + storage) --------------------
        // Salt = agent address ensures a single canonical clone per agent
        // and lets the FE compute the address ahead of the launch tx.
        clone = vestingImplementation.cloneDeterministic(_salt(agent));

        // Snapshot balance BEFORE the buy so we can compute exactly what the
        // curve delivered, even if the module already held some agent tokens
        // from a prior transfer (defensive — should always be 0 on a
        // freshly-launched agent).
        uint256 balBefore = agentToken.balanceOf(address(this));

        uint64 startTs = uint64(block.timestamp);
        configs[agent] = Config({
            creator: creator,
            vestingClone: clone,
            vestAmount: vestAmount,
            cliffSeconds: cliffSeconds,
            durationSeconds: durationSeconds,
            start: startTs,
            configured: true
        });

        // -------------------- Interactions --------------------
        // Approve the curve for exactly `titanIn`. We use a tight, transient
        // approval rather than `type(uint256).max` per protocol policy.
        quote.forceApprove(address(curve), titanIn);
        // BondingCurve.buy(minAgentOut, quoteIn) — note arg order.
        curve.buy(vestAmount, titanIn);
        // Drop any residual allowance so this module never sits on standing
        // approvals between configures.
        quote.forceApprove(address(curve), 0);

        // Verify what the curve actually delivered. The curve's own slippage
        // floor mirrors `vestAmount`, but we re-check here so a buggy or
        // malicious curve cannot under-fund the grant.
        uint256 received = agentToken.balanceOf(address(this)) - balBefore;
        if (received < vestAmount) revert InsufficientBuyOutput();

        // Forward the full received amount into the clone, then initialize
        // it with the vesting schedule. The clone reads its balance lazily on
        // {release}, so funding-before-init is the safe ordering.
        agentToken.safeTransfer(clone, received);
        TokenVesting(clone).initialize(agentToken, creator, received, startTs, cliffSeconds, durationSeconds);

        emit Configured(agent, creator, clone, received, cliffSeconds, durationSeconds);
    }

    /// @notice Trigger {TokenVesting.release} on `agent`'s clone. Anyone may
    ///         call (gas arbitrage); the released tokens always go to the
    ///         creator stored on the clone.
    /// @dev    Convenience wrapper so a single ABI on this module is enough
    ///         to drive every creator's vest. Reverts with the clone's own
    ///         errors (`NothingToRelease`, etc.) if the call is a no-op.
    /// @param  agent Agent whose clone should release.
    /// @return amount Tokens transferred to the creator this call.
    /// @dev    The clone is a contract we deployed with a known, immutable
    ///         implementation; its only outbound call is a `safeTransfer` of
    ///         the agent token to a stored beneficiary. Emitting after the
    ///         call (slither: reentrancy-events) is safe because no module
    ///         storage is mutated post-call.
    // slither-disable-next-line reentrancy-events
    function release(address agent) external nonReentrant returns (uint256 amount) {
        Config memory cfg = configs[agent];
        if (!cfg.configured) revert NotConfigured();
        amount = TokenVesting(cfg.vestingClone).release();
        emit Released(agent, cfg.vestingClone, amount);
    }

    /// @notice Predict the deterministic clone address for `agent` without
    ///         deploying. Useful for FE pre-flight rendering.
    function predictCloneAddress(address agent) external view returns (address) {
        return Clones.predictDeterministicAddress(vestingImplementation, _salt(agent), address(this));
    }

    // ---------------------------------------------------------------------
    // Recovery
    // ---------------------------------------------------------------------

    /// @notice Owner-only sweep for any quote / dust left on the module after
    ///         an aborted {configure}. Cannot drain a clone — clones hold
    ///         their own balance under their own contract, not this one.
    /// @dev    Whitelisted to ERC-20 only (no native sweep) because the
    ///         module never accepts ETH. We do not restrict the token
    ///         argument: the owner is the factory / Safe and may need to
    ///         recover any token mistakenly sent here.
    function sweep(IERC20 token, address to, uint256 amount) external onlyOwner {
        if (to == address(0)) revert ZeroAddress();
        token.safeTransfer(to, amount);
        emit QuoteSwept(address(token), to, amount);
    }

    // ---------------------------------------------------------------------
    // Internal
    // ---------------------------------------------------------------------

    /// @dev Deterministic salt derivation. Hashing the address (rather than
    ///      casting it directly to bytes32) gives uniform distribution and
    ///      future-proofs us against ever wanting to mix in a chain-id or
    ///      version byte.
    function _salt(address agent) private pure returns (bytes32) {
        return keccak256(abi.encode(agent));
    }
}
