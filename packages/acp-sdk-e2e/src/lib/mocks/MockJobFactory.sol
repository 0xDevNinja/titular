// SPDX-License-Identifier: MIT
pragma solidity 0.8.25;

/// @notice Minimal mock for `@titular/acp-sdk-e2e`. Emits a JobCreated event
///         shape-identical to the production JobFactory and returns a
///         monotonically increasing `jobId` plus a deterministic clone
///         address. Unused arguments are accepted to keep the ABI exactly
///         in sync with production so the SDK's calldata encoder is
///         exercised end-to-end.
///
/// Compiled with `solc 0.8.25 --optimize --bin-runtime`; the runtime hex is
/// pinned in `../mock-bytecode.ts` and re-derivable via
/// `pnpm -F @titular/acp-sdk-e2e run compile:mocks`.
contract MockJobFactory {
    event JobCreated(
        uint256 indexed jobId,
        address indexed clone,
        address indexed principal,
        uint8           jobType,
        address         token,
        uint256         budget
    );

    uint256 public nextId = 1;

    function createJob(
        uint256 /*targetAgentId*/,
        address token,
        uint256 budget,
        uint64  /*deadline*/,
        uint8   jobType,
        address /*evaluator*/,
        address /*arbiter*/,
        address[] calldata /*hooks*/
    ) external returns (uint256 jobId, address clone) {
        jobId = nextId++;
        // Deterministic faux-clone address; the SDK only verifies the value
        // round-trips in the receipt, so a monotonic mapping is enough.
        clone = address(uint160(0xC1000000 + jobId));
        emit JobCreated(jobId, clone, msg.sender, jobType, token, budget);
    }
}
