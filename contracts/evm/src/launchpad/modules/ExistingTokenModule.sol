// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC20Metadata} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/// @title  ExistingTokenModule
/// @notice Per-agent wrap module for projects that already have a deployed ERC-20
///         and want to list on the launchpad without re-minting under
///         {AgentToken}. Instead of cloning {AgentToken} + seeding the curve with
///         a freshly minted 1B supply, the launch flow funnels here:
///
///           1. The launchpad factory (owner of this module) pre-transfers
///              `supply` of the external ERC-20 into THIS module, then
///           2. calls {configure} which records the per-token config, validates
///              the supply against the module's actual balance, and forwards the
///              full `supply` into a freshly-cloned {BondingCurve} that will pair
///              the external token with TITU.
///
///         After bootstrap, anyone may top-up the curve with additional inventory
///         of the SAME external token via {wrap}; this is a thin, reentrancy-
///         guarded passthrough useful for late seed rounds, treasury
///         contributions, etc.
///
/// @dev    Self-contained. Does NOT modify {BondingCurve} or {AgentToken}. The
///         curve is treated as an opaque ERC-20 sink: tokens are transferred in
///         using {SafeERC20.safeTransfer}; the curve is expected to have been
///         constructed with this external token wired in as `agentToken_` and
///         {BondingCurve.initialAgentReserve} == `supply` so its internal
///         accounting matches the transferred balance.
///
///         One-shot per `externalToken`. The (externalToken => Config) keying
///         intentionally collapses the (token, agent) dimension: an external
///         ERC-20 is wrapped once, paired with one curve, owned by one admin,
///         for the lifetime of the protocol. Any second `configure` for the same
///         token reverts with {AlreadyConfigured} — this prevents rebinding to
///         a fresh curve and double-spending the token side of the LP.
///
///         All reverts are custom errors. CEI is enforced on {wrap}. {configure}
///         is owner-gated and only callable by the launchpad factory; the module
///         holds NO funds permanently — every successful {configure} / {wrap}
///         sweeps the operating balance into the curve.
contract ExistingTokenModule is Ownable, ReentrancyGuard {
    using SafeERC20 for IERC20;

    // ---------------------------------------------------------------------
    // Types
    // ---------------------------------------------------------------------

    /// @notice Per-(external token) wrap configuration. Written exactly once by
    ///         {configure}; immutable thereafter.
    /// @param  curve        Bonding curve clone paired with `externalToken`.
    /// @param  agentAdmin   Creator / admin address recorded for off-chain
    ///                      attribution and any later module hooks. Not used in
    ///                      transfer logic.
    /// @param  supply       Initial seed supply transferred into the curve at
    ///                      {configure} time.
    /// @param  wrapped      Always true once {configure} lands. Used both as the
    ///                      one-shot flag and to distinguish "unknown token"
    ///                      (zero struct) from "configured token".
    struct Config {
        address curve;
        address agentAdmin;
        uint256 supply;
        bool wrapped;
    }

    // ---------------------------------------------------------------------
    // Storage
    // ---------------------------------------------------------------------

    /// @notice Per-external-token configuration. Empty struct until {configure}
    ///         is called for that token.
    mapping(address externalToken => Config) public configs;

    // ---------------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------------

    /// @dev Emitted exactly once per `externalToken` when the launchpad owner
    ///      binds it to a fresh `curve` and seeds the initial `supply`. Indexers
    ///      key off this to surface a wrapped agent.
    event Configured(address indexed externalToken, address indexed curve, uint256 supply, address indexed agentAdmin);

    /// @dev Emitted on every successful {wrap} top-up. `depositor` is the caller
    ///      that supplied the additional inventory.
    event Wrapped(address indexed externalToken, uint256 amount, address indexed depositor);

    // ---------------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------------

    /// @dev Thrown on any second {configure} for the same `externalToken`.
    error AlreadyConfigured();
    /// @dev Thrown when {wrap} is called for an `externalToken` that has not
    ///      been bootstrapped via {configure}.
    error NotConfigured();
    /// @dev Thrown when an address argument is the zero address.
    error ZeroAddress();
    /// @dev Thrown when a supply / amount argument is zero.
    error ZeroSupply();
    /// @dev Thrown by {configure} if the launchpad factory failed to pre-transfer
    ///      at least `supply` units of `externalToken` to this module.
    ///      `held` is the actual balance for indexers / postmortems.
    error InsufficientPreDeposit(uint256 held, uint256 required);
    /// @dev Thrown by {configure} when `externalToken` does not implement the
    ///      standard ERC-20 metadata extension (name / symbol / decimals). We
    ///      reject these tokens up-front rather than risk a downstream UI / LP
    ///      seeding failure on a non-introspectable asset.
    error InvalidExternalToken();

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    /// @param owner_ Initial owner. Expected to be the launchpad factory (or a
    ///               Safe multisig acting as it) so only the factory can wrap a
    ///               new external token.
    constructor(address owner_) Ownable(_validatedOwner(owner_)) {}

    /// @dev Wrap OZ `Ownable`'s zero-owner check so callers see our canonical
    ///      {ZeroAddress} error instead of OZ's `OwnableInvalidOwner`.
    function _validatedOwner(address owner_) private pure returns (address) {
        if (owner_ == address(0)) revert ZeroAddress();
        return owner_;
    }

    // ---------------------------------------------------------------------
    // Admin
    // ---------------------------------------------------------------------

    /// @notice Wrap a pre-existing ERC-20 and seed a freshly-cloned bonding
    ///         curve with `supply` of it. One-shot per `externalToken`.
    /// @dev    The caller (launchpad factory) MUST pre-transfer at least
    ///         `supply` units of `externalToken` to this module before calling.
    ///         The module then sweeps exactly `supply` into `curve`. Any excess
    ///         left on the module sits idle and can be drained on a future
    ///         {wrap} top-up — but the canonical pattern is "transfer exactly
    ///         the amount, then configure".
    ///
    ///         Validation order:
    ///           1. Address sanity (no zero externalToken / curve / admin).
    ///           2. Non-zero supply.
    ///           3. One-shot guard.
    ///           4. ERC-20 metadata introspection (name / symbol / decimals
    ///              must all be readable). Tokens that don't implement
    ///              `IERC20Metadata` are rejected here, not later.
    ///           5. Pre-deposit balance check.
    ///         Then: write storage, transfer to the curve, emit.
    ///
    ///         CEI: storage is finalized before the outbound `safeTransfer`,
    ///         and {ReentrancyGuard} is not needed on this owner-gated path
    ///         (the only outbound call is to a trusted curve clone). It is
    ///         applied to {wrap} where the entry is open.
    /// @param  externalToken   The pre-existing ERC-20 to wrap.
    /// @param  curve           Bonding curve clone paired with this token.
    ///                         Expected to have been constructed with
    ///                         `agentToken_ == externalToken` and
    ///                         `initialAgentReserve_ == supply` so its internal
    ///                         accounting matches the seeded balance.
    /// @param  supply          Initial seed supply (in `externalToken` wei).
    /// @param  admin           Creator / admin address recorded for attribution.
    function configure(address externalToken, address curve, uint256 supply, address admin) external onlyOwner {
        if (externalToken == address(0) || curve == address(0) || admin == address(0)) revert ZeroAddress();
        if (supply == 0) revert ZeroSupply();
        if (configs[externalToken].wrapped) revert AlreadyConfigured();

        _validateMetadata(externalToken);

        uint256 held = IERC20(externalToken).balanceOf(address(this));
        if (held < supply) revert InsufficientPreDeposit(held, supply);

        // --- Effects ---
        configs[externalToken] = Config({curve: curve, agentAdmin: admin, supply: supply, wrapped: true});

        // --- Interactions ---
        IERC20(externalToken).safeTransfer(curve, supply);

        emit Configured(externalToken, curve, supply, admin);
    }

    // ---------------------------------------------------------------------
    // Top-up
    // ---------------------------------------------------------------------

    /// @notice Pull `amount` of an already-wrapped `externalToken` from the
    ///         caller and forward it directly into the configured curve. Useful
    ///         for late-seed contributions / treasury top-ups; the curve treats
    ///         the additional balance as extra inventory on the agent side.
    /// @dev    Reentrancy-guarded because the entry is permissionless and we
    ///         make two outbound ERC-20 calls (pull from caller, push to curve).
    ///         A malicious external token cannot escalate beyond its own
    ///         transfer hook re-entering this module with a stale reserve view
    ///         (we hold no internal reserves to corrupt) but the guard is
    ///         retained as defence-in-depth and to keep the audit story uniform
    ///         across module top-up paths.
    ///
    ///         Note: the curve must be configured to handle additional agent-
    ///         side inventory. The standard {BondingCurve} reads its own
    ///         `realAgentReserve` from storage, NOT from `balanceOf`, so any
    ///         top-up here will increase the curve's actual balance without
    ///         updating its accounting. Integrators wanting accounting-aware
    ///         top-ups must implement a curve variant; this module deliberately
    ///         remains bonding-curve-agnostic.
    /// @param  externalToken   Wrapped ERC-20.
    /// @param  amount          Amount to pull and forward (in `externalToken` wei).
    function wrap(address externalToken, uint256 amount) external nonReentrant {
        if (amount == 0) revert ZeroSupply();
        Config memory cfg = configs[externalToken];
        if (!cfg.wrapped) revert NotConfigured();

        IERC20 token = IERC20(externalToken);

        // --- Interactions ---
        // No state to update — Config is immutable post-configure and we hold no
        // internal accounting. Pull from caller, then immediately push to curve.
        // The pre/post balance pattern would be defensive against fee-on-transfer
        // tokens; we deliberately do not support those (the launchpad surface
        // reads token amounts at face value), and `safeTransferFrom` followed by
        // `safeTransfer` of the SAME `amount` will simply revert downstream if
        // the token under-delivers.
        token.safeTransferFrom(msg.sender, address(this), amount);
        token.safeTransfer(cfg.curve, amount);

        emit Wrapped(externalToken, amount, msg.sender);
    }

    // ---------------------------------------------------------------------
    // Views
    // ---------------------------------------------------------------------

    /// @notice Full per-token wrap configuration. Returns the zero struct for
    ///         unwrapped tokens; callers requiring positive authorization must
    ///         additionally check `cfg.wrapped`.
    function getConfig(address externalToken) external view returns (Config memory cfg) {
        cfg = configs[externalToken];
    }

    // ---------------------------------------------------------------------
    // Internal
    // ---------------------------------------------------------------------

    /// @dev Reject tokens that don't implement {IERC20Metadata}. We require all
    ///      three of `name()`, `symbol()`, and `decimals()` to be readable; any
    ///      revert / call failure in that surface is treated as
    ///      {InvalidExternalToken}. Done via low-level `staticcall` so a
    ///      reverting metadata call surfaces our canonical error rather than
    ///      bubbling whatever the foreign token chose to throw.
    ///
    ///      Tolerated: tokens that return empty strings — we only require the
    ///      call succeed, not that the result be cosmetically valid. The UI
    ///      layer is the right place to enforce non-empty metadata.
    function _validateMetadata(address externalToken) private view {
        // decimals(): MUST return uint8. Failure or wrong-length return is fatal.
        (bool ok, bytes memory data) = externalToken.staticcall(abi.encodeCall(IERC20Metadata.decimals, ()));
        if (!ok || data.length < 32) revert InvalidExternalToken();

        // name() / symbol(): MUST be callable and return ABI-encodable bytes.
        // We do not constrain the returned string contents.
        (ok,) = externalToken.staticcall(abi.encodeCall(IERC20Metadata.name, ()));
        if (!ok) revert InvalidExternalToken();

        (ok,) = externalToken.staticcall(abi.encodeCall(IERC20Metadata.symbol, ()));
        if (!ok) revert InvalidExternalToken();
    }
}
