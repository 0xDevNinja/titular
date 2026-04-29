// Minimal ABI fragments needed by AcpAgent. Hand-trimmed from
// contracts/evm/src/acp/AgentRegistry.sol and JobFactory.sol so the SDK
// doesn't drag the full forge artifact set into its npm tarball.
//
// If a function signature here drifts from the contract, the M4-acp-sdk-e2e
// integration tests (#96) will fail at registration / job creation time.

export const AGENT_REGISTRY_ABI = [
  {
    type: "function",
    name: "register",
    stateMutability: "nonpayable",
    inputs: [
      { name: "controller", type: "address" },
      { name: "metadataURI", type: "string" },
      { name: "capabilities", type: "uint256" },
    ],
    outputs: [{ name: "agentId", type: "uint256" }],
  },
  {
    type: "event",
    name: "AgentRegistered",
    inputs: [
      { name: "agentId", type: "uint256", indexed: true },
      { name: "controller", type: "address", indexed: true },
      { name: "metadataURI", type: "string", indexed: false },
      { name: "capabilities", type: "uint256", indexed: false },
    ],
    anonymous: false,
  },
] as const;

export const JOB_FACTORY_ABI = [
  {
    type: "function",
    name: "createJob",
    stateMutability: "nonpayable",
    inputs: [
      { name: "targetAgentId", type: "uint256" },
      { name: "token", type: "address" },
      { name: "budget", type: "uint256" },
      { name: "deadline", type: "uint64" },
      { name: "jobType", type: "uint8" },
      { name: "evaluator", type: "address" },
      { name: "arbiter", type: "address" },
      { name: "jobHooks", type: "address[]" },
    ],
    outputs: [
      { name: "jobId", type: "uint256" },
      { name: "clone", type: "address" },
    ],
  },
  {
    type: "event",
    name: "JobCreated",
    inputs: [
      { name: "jobId", type: "uint256", indexed: true },
      { name: "clone", type: "address", indexed: true },
      { name: "principal", type: "address", indexed: true },
      { name: "jobType", type: "uint8", indexed: false },
      { name: "token", type: "address", indexed: false },
      { name: "budget", type: "uint256", indexed: false },
    ],
    anonymous: false,
  },
] as const;
