// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {Pausable} from "@openzeppelin/contracts/utils/Pausable.sol";

/// @title AgentRegistry
/// @notice Composite identity registry for ACP v2 agents.  Each agent maps a controller
///         address to an off-chain metadata URI and a packed capability bitmask.  The
///         registry also accumulates on-chain reputation scores posted by authorised
///         scorers, with EIP-712 nonce-based replay protection.
/// @dev    Role model
///         - DEFAULT_ADMIN_ROLE : can grant / revoke other roles, pause / unpause.
///         - SCORER_ROLE        : authorised to call `postScore`.
///         - PAUSER_ROLE        : can pause the registry.
///         Two-step controller transfer: propose → accept prevents fat-finger takeovers.
contract AgentRegistry is AccessControl, Pausable {
    // ─────────────────────────────────────────────────────────────
    // Roles
    // ─────────────────────────────────────────────────────────────

    bytes32 public constant SCORER_ROLE = keccak256("SCORER_ROLE");
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");

    // ─────────────────────────────────────────────────────────────
    // Types
    // ─────────────────────────────────────────────────────────────

    struct AgentInfo {
        /// @notice Current controller address that manages the agent identity.
        address controller;
        /// @notice IPFS / HTTPS URI pointing to AgentMetadata JSON.
        string metadataURI;
        /// @notice Packed bitmask of agent capabilities (consumer-defined).
        uint256 capabilities;
        /// @notice Cumulative reputation score (sum of all accepted `postScore` calls).
        int256 reputationScore;
        /// @notice Block timestamp of first registration.
        uint48 registeredAt;
        /// @notice Whether this agent is currently active (controller may deactivate).
        bool active;
    }

    // ─────────────────────────────────────────────────────────────
    // Storage
    // ─────────────────────────────────────────────────────────────

    /// @notice Sequential counter, used as agent ID.
    uint256 private _nextId;

    /// @notice agentId => AgentInfo
    mapping(uint256 => AgentInfo) private _agents;

    /// @notice controller address => agentId list (an address may control multiple agents)
    mapping(address => uint256[]) private _controllerAgents;

    /// @notice pending two-step transfer: agentId => proposed new controller
    mapping(uint256 => address) private _pendingController;

    /// @notice EIP-712 replay protection: agentId => nonce (incremented after each `postScore`)
    mapping(uint256 => uint256) public scorerNonce;

    // ─────────────────────────────────────────────────────────────
    // Events
    // ─────────────────────────────────────────────────────────────

    /// @notice Emitted when a new agent identity is registered.
    event AgentRegistered(uint256 indexed agentId, address indexed controller, string metadataURI, uint256 capabilities);

    /// @notice Emitted when an agent's metadata URI is updated.
    event MetadataUpdated(uint256 indexed agentId, string metadataURI);

    /// @notice Emitted when capabilities bitmask is changed.
    event CapabilitiesUpdated(uint256 indexed agentId, uint256 capabilities);

    /// @notice Emitted when an agent is activated or deactivated.
    event ActiveStatusChanged(uint256 indexed agentId, bool active);

    /// @notice Emitted when a controller transfer is initiated.
    event ControllerTransferProposed(uint256 indexed agentId, address indexed proposed);

    /// @notice Emitted when a controller transfer is accepted by the proposed address.
    event ControllerTransferAccepted(uint256 indexed agentId, address indexed oldController, address indexed newController);

    /// @notice Emitted when a controller transfer proposal is cancelled.
    event ControllerTransferCancelled(uint256 indexed agentId);

    /// @notice Emitted when a reputation score is posted for an agent.
    event ScorePosted(uint256 indexed agentId, address indexed scorer, int256 delta, int256 newTotal, uint256 nonce);

    // ─────────────────────────────────────────────────────────────
    // Errors
    // ─────────────────────────────────────────────────────────────

    error ZeroAddress();
    error NotController(uint256 agentId, address caller);
    error AgentNotFound(uint256 agentId);
    error AgentInactive(uint256 agentId);
    error EmptyMetadataURI();
    error NoPendingTransfer(uint256 agentId);
    error NotProposedController(uint256 agentId, address caller);
    error InvalidNonce(uint256 agentId, uint256 expected, uint256 provided);

    // ─────────────────────────────────────────────────────────────
    // Constructor
    // ─────────────────────────────────────────────────────────────

    /// @param admin Address granted DEFAULT_ADMIN_ROLE on deployment.
    constructor(address admin) {
        if (admin == address(0)) revert ZeroAddress();
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(PAUSER_ROLE, admin);
    }

    // ─────────────────────────────────────────────────────────────
    // Registration
    // ─────────────────────────────────────────────────────────────

    /// @notice Register a new agent identity.
    /// @param  controller   Address that will control this agent (cannot be zero).
    /// @param  metadataURI  Non-empty URI pointing to off-chain agent metadata.
    /// @param  capabilities Packed bitmask of declared capabilities (consumer-defined).
    /// @return agentId      Sequential ID assigned to this agent.
    function register(address controller, string calldata metadataURI, uint256 capabilities)
        external
        whenNotPaused
        returns (uint256 agentId)
    {
        if (controller == address(0)) revert ZeroAddress();
        if (bytes(metadataURI).length == 0) revert EmptyMetadataURI();

        agentId = _nextId++;
        _agents[agentId] = AgentInfo({
            controller: controller,
            metadataURI: metadataURI,
            capabilities: capabilities,
            reputationScore: 0,
            registeredAt: uint48(block.timestamp),
            active: true
        });
        _controllerAgents[controller].push(agentId);

        emit AgentRegistered(agentId, controller, metadataURI, capabilities);
    }

    // ─────────────────────────────────────────────────────────────
    // Agent Management (controller-only)
    // ─────────────────────────────────────────────────────────────

    /// @notice Update the metadata URI of an agent.
    /// @param agentId     ID of the agent to update.
    /// @param metadataURI New non-empty URI.
    function setMetadata(uint256 agentId, string calldata metadataURI) external whenNotPaused {
        _assertController(agentId);
        if (bytes(metadataURI).length == 0) revert EmptyMetadataURI();
        _agents[agentId].metadataURI = metadataURI;
        emit MetadataUpdated(agentId, metadataURI);
    }

    /// @notice Update the capabilities bitmask of an agent.
    /// @param agentId      ID of the agent to update.
    /// @param capabilities New capabilities bitmask.
    function setCapabilities(uint256 agentId, uint256 capabilities) external whenNotPaused {
        _assertController(agentId);
        _agents[agentId].capabilities = capabilities;
        emit CapabilitiesUpdated(agentId, capabilities);
    }

    /// @notice Activate or deactivate an agent.
    /// @param agentId ID of the agent.
    /// @param active  Desired active status.
    function setActive(uint256 agentId, bool active) external whenNotPaused {
        _assertController(agentId);
        _agents[agentId].active = active;
        emit ActiveStatusChanged(agentId, active);
    }

    // ─────────────────────────────────────────────────────────────
    // Two-step controller transfer
    // ─────────────────────────────────────────────────────────────

    /// @notice Propose a new controller address for agentId.
    ///         The proposed address must call `acceptControllerTransfer` to complete.
    /// @param agentId  ID of the agent.
    /// @param proposed New controller candidate (cannot be zero).
    function proposeControllerTransfer(uint256 agentId, address proposed) external whenNotPaused {
        _assertController(agentId);
        if (proposed == address(0)) revert ZeroAddress();
        _pendingController[agentId] = proposed;
        emit ControllerTransferProposed(agentId, proposed);
    }

    /// @notice Accept a pending controller transfer initiated by the current controller.
    /// @param agentId ID of the agent.
    function acceptControllerTransfer(uint256 agentId) external whenNotPaused {
        address proposed = _pendingController[agentId];
        if (proposed == address(0)) revert NoPendingTransfer(agentId);
        if (proposed != msg.sender) revert NotProposedController(agentId, msg.sender);

        address old = _agents[agentId].controller;
        _agents[agentId].controller = proposed;
        delete _pendingController[agentId];
        _controllerAgents[proposed].push(agentId);

        emit ControllerTransferAccepted(agentId, old, proposed);
    }

    /// @notice Cancel a pending controller transfer (callable by current controller only).
    /// @param agentId ID of the agent.
    function cancelControllerTransfer(uint256 agentId) external {
        _assertController(agentId);
        if (_pendingController[agentId] == address(0)) revert NoPendingTransfer(agentId);
        delete _pendingController[agentId];
        emit ControllerTransferCancelled(agentId);
    }

    // ─────────────────────────────────────────────────────────────
    // Reputation accumulation (SCORER_ROLE)
    // ─────────────────────────────────────────────────────────────

    /// @notice Post a signed reputation delta for an agent.
    ///         Callers must hold SCORER_ROLE. The nonce must equal the current value
    ///         in `scorerNonce[agentId]` to prevent replay.
    /// @param agentId ID of the agent receiving the score.
    /// @param delta   Signed delta to add to the cumulative score (may be negative).
    /// @param nonce   Must equal `scorerNonce[agentId]`; incremented after acceptance.
    function postScore(uint256 agentId, int256 delta, uint256 nonce)
        external
        whenNotPaused
        onlyRole(SCORER_ROLE)
    {
        AgentInfo storage info = _agents[agentId];
        // slither-disable-next-line incorrect-equality,timestamp
        if (info.registeredAt == 0) revert AgentNotFound(agentId);
        if (!info.active) revert AgentInactive(agentId);

        uint256 expected = scorerNonce[agentId];
        if (nonce != expected) revert InvalidNonce(agentId, expected, nonce);

        scorerNonce[agentId] = expected + 1;
        info.reputationScore += delta;

        emit ScorePosted(agentId, msg.sender, delta, info.reputationScore, nonce);
    }

    // ─────────────────────────────────────────────────────────────
    // Pause
    // ─────────────────────────────────────────────────────────────

    /// @notice Pause the registry (blocks registrations, updates, and score posts).
    function pause() external onlyRole(PAUSER_ROLE) {
        _pause();
    }

    /// @notice Unpause the registry.
    function unpause() external onlyRole(PAUSER_ROLE) {
        _unpause();
    }

    // ─────────────────────────────────────────────────────────────
    // Views
    // ─────────────────────────────────────────────────────────────

    /// @notice Return the full AgentInfo struct for a given agentId.
    /// @param agentId ID of the agent.
    /// @return info   The stored AgentInfo (reverts if never registered).
    function getAgent(uint256 agentId) external view returns (AgentInfo memory info) {
        info = _agents[agentId];
        // slither-disable-next-line incorrect-equality,timestamp
        if (info.registeredAt == 0) revert AgentNotFound(agentId);
    }

    /// @notice Return all agentIds controlled by a given address.
    /// @param controller The controller address to query.
    /// @return ids       Array of agentIds.  May contain deactivated agents.
    function agentsByController(address controller) external view returns (uint256[] memory ids) {
        ids = _controllerAgents[controller];
    }

    /// @notice Total number of agents ever registered (including inactive ones).
    /// @return count The count.
    function totalAgents() external view returns (uint256 count) {
        count = _nextId;
    }

    /// @notice Return the pending proposed controller for agentId (zero if none).
    /// @param agentId ID of the agent.
    /// @return proposed The pending controller address, or address(0) if none.
    function pendingController(uint256 agentId) external view returns (address proposed) {
        proposed = _pendingController[agentId];
    }

    // ─────────────────────────────────────────────────────────────
    // Internal helpers
    // ─────────────────────────────────────────────────────────────

    function _assertController(uint256 agentId) internal view {
        AgentInfo storage info = _agents[agentId];
        // slither-disable-next-line incorrect-equality,timestamp
        if (info.registeredAt == 0) revert AgentNotFound(agentId);
        if (info.controller != msg.sender) revert NotController(agentId, msg.sender);
    }
}
