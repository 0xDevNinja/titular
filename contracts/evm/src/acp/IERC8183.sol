// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

/// @title IERC-8183
/// @notice Agent Communication Protocol v2 — on-chain interface standard.
/// @dev    ERC-8183 defines the minimal interface an ACP v2 agent registry must implement
///         to be composable with tooling (indexers, SDKs, routers).
///
///         Inspired by ERC-6551 (token-bound accounts) and ERC-7484 (registry adapters).
///         This interface is informally proposed pending formal EIP submission.
interface IERC8183 {
    // ─────────────────────────────────────────────────────────────
    // Agent identity
    // ─────────────────────────────────────────────────────────────

    /// @notice Register a new agent identity.
    /// @param  controller   Address that controls this agent.
    /// @param  metadataURI  URI pointing to agent metadata JSON (IPFS or HTTPS).
    /// @param  capabilities Packed bitmask of declared capabilities.
    /// @return agentId      Sequential ID assigned to this agent.
    function register(address controller, string calldata metadataURI, uint256 capabilities)
        external
        returns (uint256 agentId);

    /// @notice Return the full info for a given agentId.
    /// @param agentId The agent identifier.
    /// @return controller   Current controller address.
    /// @return metadataURI  Metadata URI.
    /// @return capabilities Capabilities bitmask.
    /// @return active       Whether the agent is active.
    function agentInfo(uint256 agentId)
        external
        view
        returns (address controller, string memory metadataURI, uint256 capabilities, bool active);

    /// @notice Total number of registered agents (may include inactive ones).
    /// @return count Total agent count.
    function totalAgents() external view returns (uint256 count);

    // ─────────────────────────────────────────────────────────────
    // Reputation
    // ─────────────────────────────────────────────────────────────

    /// @notice Post a signed reputation score for an agent.
    /// @param agentId The agent to score.
    /// @param delta   Signed score delta.
    /// @param nonce   Replay-protection nonce.
    function postScore(uint256 agentId, int256 delta, uint256 nonce) external;

    /// @notice Return the cumulative reputation score for an agent.
    /// @param agentId The agent to query.
    /// @return score Cumulative reputation score.
    function reputationScore(uint256 agentId) external view returns (int256 score);
}
