// SPDX-License-Identifier: MIT
pragma solidity 0.8.25;

/// @notice Minimal mock used only by `@titular/acp-sdk-e2e`. Emits the same
///         AgentRegistered event signature as the real AgentRegistry, with
///         a monotonically increasing `agentId`. The SDK's event-decoder
///         path treats this contract as indistinguishable from production.
///
/// Compiled with `solc 0.8.25 --optimize --bin-runtime`; the runtime hex is
/// pinned in `../mock-bytecode.ts` and re-derivable via
/// `pnpm -F @titular/acp-sdk-e2e run compile:mocks`.
contract MockAgentRegistry {
    event AgentRegistered(
        uint256 indexed agentId,
        address indexed controller,
        string  metadataURI,
        uint256 capabilities
    );

    uint256 public nextId = 1;

    function register(
        address controller,
        string  calldata metadataURI,
        uint256 capabilities
    ) external returns (uint256 agentId) {
        agentId = nextId++;
        emit AgentRegistered(agentId, controller, metadataURI, capabilities);
    }
}
