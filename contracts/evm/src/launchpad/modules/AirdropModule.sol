// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {MerkleProof} from "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";

/// @title  AirdropModule
/// @notice Per-agent merkle-claim airdrop. The launchpad factory snapshots
///         veTITU voting power at an agreed block, an off-chain indexer builds
///         the merkle tree of `(recipient, amount)` allocations, and the root +
///         a one-shot per-agent allocation are committed on-chain via
///         {configure}. Recipients then call {claim} with their merkle proof to
///         pull their share of the agent token from this module.
/// @dev    Single shared instance across the agent fleet. Per-agent funds are
///         **pre-deposited** by the configurer (factory or creator script): the
///         module never mints, never holds approval, never pulls. {configure}
///         enforces `balanceOf(this) >= allocation` so a misfunded launch
///         reverts atomically. The 5%-of-supply cap is enforced at the same
///         site so the module cannot be used to pre-allocate more than the
///         spec allows.
///
///         Storage shape mirrors the other launchpad modules: a `Config` struct
///         per agent + a per-(agent, recipient) nullifier mapping. The
///         nullifier doubles as the double-claim guard.
///
///         Reentrancy: {claim} and {sweep} both transfer ERC-20 (which can be a
///         malicious agent token, since the launchpad supports `ExistingToken`
///         wraps in a future module). Both functions are guarded with
///         {ReentrancyGuard} and follow strict CEI — nullifier flipped before
///         the external transfer.
contract AirdropModule is Ownable, ReentrancyGuard {
    using SafeERC20 for IERC20;

    // ---------------------------------------------------------------------
    // Constants
    // ---------------------------------------------------------------------

    /// @notice Basis-point denominator. 10_000 = 100%.
    uint16 public constant MAX_BPS = 10_000;

    /// @notice Maximum airdrop allocation as a fraction of agent total supply,
    ///         in basis points. 500 = 5%. Hard-coded by protocol spec.
    uint16 public constant MAX_ALLOCATION_BPS = 500;

    // ---------------------------------------------------------------------
    // Types
    // ---------------------------------------------------------------------

    /// @notice Per-agent airdrop configuration.
    /// @param  merkleRoot    Root of the `(recipient, amount)` allocation tree.
    ///                       Leaves are encoded per OZ's standard double-hash
    ///                       scheme: `keccak256(bytes.concat(keccak256(abi.encode(addr, amount))))`.
    /// @param  snapshotBlock Block at which veTITU voting power was sampled
    ///                       off-chain. Stored on-chain as a public commitment
    ///                       so the indexer's snapshot is auditable; never read
    ///                       inside the contract logic.
    /// @param  allocation    Total agent tokens dedicated to this airdrop. Pre-
    ///                       deposited by the configurer; capped at
    ///                       {MAX_ALLOCATION_BPS} of the agent's totalSupply.
    /// @param  deadline      Unix timestamp after which {sweep} unlocks the
    ///                       unclaimed dust to {agentAdmin}. Claims after the
    ///                       deadline revert.
    /// @param  agentAdmin    Address authorized to call {sweep}. Expected to be
    ///                       the agent creator (the recipient of unclaimed dust).
    /// @param  configured    Set on first {configure}; one-shot guard.
    struct Config {
        bytes32 merkleRoot;
        uint128 allocation;
        uint64 snapshotBlock;
        uint64 deadline;
        address agentAdmin;
        bool configured;
    }

    // ---------------------------------------------------------------------
    // Storage
    // ---------------------------------------------------------------------

    /// @notice Per-agent airdrop record. Empty until {configure} is called.
    mapping(address agent => Config) public configs;

    /// @notice Per-(agent, recipient) double-claim nullifier. Flipped to true
    ///         the first time a recipient claims; stays true forever.
    mapping(address agent => mapping(address recipient => bool)) public claimed;

    // ---------------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------------

    /// @dev Emitted exactly once per agent on {configure}. Indexers key the
    ///      merkle tree off this event's `merkleRoot`.
    event Configured(
        address indexed agent,
        bytes32 merkleRoot,
        uint64 snapshotBlock,
        uint128 allocation,
        uint64 deadline,
        address agentAdmin
    );

    /// @dev Emitted on every successful {claim}. `recipient` is the address
    ///      that the merkle leaf was bound to; `amount` is the leaf amount.
    event Claimed(address indexed agent, address indexed recipient, uint256 amount);

    /// @dev Emitted on a successful {sweep}. `amount` is the dust returned to
    ///      the agent admin.
    event Swept(address indexed agent, uint256 amount);

    // ---------------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------------

    /// @dev An address argument was the zero address.
    error ZeroAddress();
    /// @dev `allocation == 0` on {configure} — would lock no tokens and create
    ///      an empty record indistinguishable from "unconfigured".
    error ZeroAllocation();
    /// @dev `merkleRoot == bytes32(0)` on {configure} — empty root would let
    ///      anyone forge proofs.
    error ZeroRoot();
    /// @dev {configure} called twice for the same agent.
    error AlreadyConfigured();
    /// @dev {configure} called for an agent the module is not yet a registered
    ///      configurer for. Surfaced separately from {AlreadyConfigured} for
    ///      indexer clarity.
    error NotConfigured();
    /// @dev `allocation` exceeds {MAX_ALLOCATION_BPS} of the agent's
    ///      totalSupply. `cap` is the computed ceiling for diagnostics.
    error AllocationOverCap(uint256 allocation, uint256 cap);
    /// @dev Pre-deposit balance is below the requested allocation at
    ///      {configure} time.
    error InsufficientBalance(uint256 have, uint256 need);
    /// @dev Claim deadline is at or before `block.timestamp` at {configure}.
    error DeadlineInPast();
    /// @dev Merkle proof did not verify against the stored root.
    error InvalidProof();
    /// @dev Recipient already claimed for this agent.
    error AlreadyClaimed();
    /// @dev Claim attempted after `deadline` (or sweep before).
    error DeadlinePassed();
    error DeadlineNotPassed();
    /// @dev {sweep} caller is not the agent admin.
    error NotAgentAdmin();
    /// @dev {sweep} when there is nothing to return — avoids emitting noise.
    error NothingToSweep();

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    /// @param owner_ Initial owner. Expected to be the launchpad factory (or a
    ///               Safe acting as it) so only the factory can bind a merkle
    ///               root to a freshly deployed agent.
    constructor(address owner_) Ownable(_validatedOwner(owner_)) {}

    /// @dev Promote OZ's `OwnableInvalidOwner` into our canonical
    ///      {ZeroAddress} so the ABI surfaces one error for missing addresses.
    function _validatedOwner(address owner_) private pure returns (address) {
        if (owner_ == address(0)) revert ZeroAddress();
        return owner_;
    }

    // ---------------------------------------------------------------------
    // Admin
    // ---------------------------------------------------------------------

    /// @notice Bind an airdrop merkle root to `agent`. One-shot per agent.
    /// @dev    Caller MUST have transferred at least `allocation` units of the
    ///         agent token to this module **before** calling. The check is
    ///         atomic — partial pre-deposit reverts and the configurer can
    ///         top up and retry on the same agent (since the record is not
    ///         written until all checks pass).
    ///
    ///         The 5% cap is computed against `IERC20(agent).totalSupply()` at
    ///         configure time. AgentToken supply is fixed at deploy
    ///         (1 B units) so this snapshot is authoritative.
    /// @param  agent          Agent token (the asset airdropped).
    /// @param  root           Merkle root of `(recipient, amount)` allocations.
    /// @param  snapshotBlock  Block at which the indexer snapped veTITU.
    /// @param  allocation     Total airdrop pool, denominated in agent token.
    /// @param  deadline       Unix timestamp after which dust is sweepable.
    /// @param  admin          Address authorised to call {sweep}.
    function configure(
        address agent,
        bytes32 root,
        uint64 snapshotBlock,
        uint128 allocation,
        uint64 deadline,
        address admin
    ) external onlyOwner {
        if (agent == address(0) || admin == address(0)) revert ZeroAddress();
        if (root == bytes32(0)) revert ZeroRoot();
        if (allocation == 0) revert ZeroAllocation();
        if (configs[agent].configured) revert AlreadyConfigured();
        // `deadline` is a wall-clock window measured in days (default ~30d);
        // a few seconds of miner drift is immaterial.
        // slither-disable-next-line timestamp
        if (deadline <= block.timestamp) revert DeadlineInPast();

        // Cap = totalSupply * MAX_ALLOCATION_BPS / MAX_BPS.
        // AgentToken caps at 1e27, so totalSupply * 500 fits comfortably below
        // 2^256; checked math would catch any pathological wrap regardless.
        uint256 cap = (IERC20(agent).totalSupply() * uint256(MAX_ALLOCATION_BPS)) / uint256(MAX_BPS);
        if (uint256(allocation) > cap) revert AllocationOverCap(uint256(allocation), cap);

        uint256 bal = IERC20(agent).balanceOf(address(this));
        if (bal < uint256(allocation)) revert InsufficientBalance(bal, uint256(allocation));

        configs[agent] = Config({
            merkleRoot: root,
            allocation: allocation,
            snapshotBlock: snapshotBlock,
            deadline: deadline,
            agentAdmin: admin,
            configured: true
        });

        emit Configured(agent, root, snapshotBlock, allocation, deadline, admin);
    }

    // ---------------------------------------------------------------------
    // Claim
    // ---------------------------------------------------------------------

    /// @notice Claim `amount` of the `agent` token using a merkle proof bound
    ///         to `recipient`. Anyone can submit the proof on behalf of the
    ///         recipient — the leaf is keyed on `recipient`, not `msg.sender`.
    /// @dev    Strict CEI: nullifier flipped before the transfer so a
    ///         malicious agent token cannot reenter and double-claim.
    ///         {nonReentrant} is belt-and-braces against multi-claim reentry
    ///         across recipients.
    ///
    ///         Leaf encoding matches OZ's `merkle-tree` JS library:
    ///         `keccak256(bytes.concat(keccak256(abi.encode(recipient, amount))))`.
    /// @param  agent     Agent token being claimed.
    /// @param  recipient Address the leaf is bound to.
    /// @param  amount    Allocation encoded into the leaf.
    /// @param  proof     Merkle proof from leaf to root.
    function claim(address agent, address recipient, uint256 amount, bytes32[] calldata proof) external nonReentrant {
        Config memory cfg = configs[agent];
        if (!cfg.configured) revert NotConfigured();
        // slither-disable-next-line timestamp
        if (block.timestamp > cfg.deadline) revert DeadlinePassed();
        if (claimed[agent][recipient]) revert AlreadyClaimed();

        bytes32 leaf = keccak256(bytes.concat(keccak256(abi.encode(recipient, amount))));
        if (!MerkleProof.verify(proof, cfg.merkleRoot, leaf)) revert InvalidProof();

        // Effects.
        claimed[agent][recipient] = true;

        // Interaction.
        emit Claimed(agent, recipient, amount);
        IERC20(agent).safeTransfer(recipient, amount);
    }

    // ---------------------------------------------------------------------
    // Sweep
    // ---------------------------------------------------------------------

    /// @notice Return any unclaimed agent-token dust to the agent admin.
    /// @dev    Gated on `block.timestamp > deadline` so the admin cannot rug
    ///         a pending claim. Anyone-callable would be safer (no admin trust
    ///         required) but the spec routes dust back to the creator, so we
    ///         restrict to {agentAdmin} for unambiguous accountability.
    ///         CEI: balance read, transfer last, no state change after.
    /// @param  agent Agent token to sweep.
    function sweep(address agent) external nonReentrant {
        Config memory cfg = configs[agent];
        if (!cfg.configured) revert NotConfigured();
        if (msg.sender != cfg.agentAdmin) revert NotAgentAdmin();
        // slither-disable-next-line timestamp
        if (block.timestamp <= cfg.deadline) revert DeadlineNotPassed();

        uint256 bal = IERC20(agent).balanceOf(address(this));
        // Strict equality with zero is the exact semantic we want — revert on
        // an empty pool so a re-sweep does not emit a noisy zero-amount event.
        // slither-disable-next-line incorrect-equality
        if (bal == 0) revert NothingToSweep();

        emit Swept(agent, bal);
        IERC20(agent).safeTransfer(cfg.agentAdmin, bal);
    }
}
