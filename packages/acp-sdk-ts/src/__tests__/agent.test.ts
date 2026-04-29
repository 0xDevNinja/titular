import { describe, expect, it, vi } from "vitest";
import { AcpAgent } from "../agent.js";
import { AcpError } from "../errors.js";
import type {
  ContractWriteParams,
  EvmProviderAdapter,
  ParsedLog,
  TxReceipt,
} from "../providers/types.js";
import { type EvmAddress, JobType, type TxHash } from "../types.js";
import { encodeEventToLog } from "./_helpers.js";

const AGENT_REGISTERED_EVENT = {
  type: "event",
  name: "AgentRegistered",
  inputs: [
    { name: "agentId", type: "uint256", indexed: true },
    { name: "controller", type: "address", indexed: true },
    { name: "metadataURI", type: "string", indexed: false },
    { name: "capabilities", type: "uint256", indexed: false },
  ],
  anonymous: false,
} as const;

const JOB_CREATED_EVENT = {
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
} as const;

const REGISTRY: EvmAddress = "0x1111111111111111111111111111111111111111";
const FACTORY: EvmAddress = "0x2222222222222222222222222222222222222222";
const SIGNER: EvmAddress = "0x3333333333333333333333333333333333333333";
const TOKEN: EvmAddress = "0x4444444444444444444444444444444444444444";
const CLONE: EvmAddress = "0x5555555555555555555555555555555555555555";
const EVALUATOR: EvmAddress = "0x6666666666666666666666666666666666666666";
const ZERO: EvmAddress = "0x0000000000000000000000000000000000000000";

const TX_HASH: TxHash = "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";

/**
 * Build an `AgentRegistered` log shaped exactly the way the EVM emits it.
 */
function encodeAgentRegistered(
  agentId: bigint,
  controller: EvmAddress,
  metadataURI: string,
  capabilities: bigint
): ParsedLog {
  return encodeEventToLog(
    AGENT_REGISTERED_EVENT,
    { agentId, controller, metadataURI, capabilities },
    REGISTRY
  );
}

function encodeJobCreated(
  jobId: bigint,
  clone: EvmAddress,
  principal: EvmAddress,
  jobType: number,
  token: EvmAddress,
  budget: bigint
): ParsedLog {
  return encodeEventToLog(
    JOB_CREATED_EVENT,
    { jobId, clone, principal, jobType, token, budget },
    FACTORY
  );
}

function makeReceipt(status: "success" | "reverted", logs: ParsedLog[] = []): TxReceipt {
  return {
    status,
    transactionHash: TX_HASH,
    blockNumber: 123n,
    logs,
  };
}

interface MockProvider extends EvmProviderAdapter {
  writeContractMock: ReturnType<typeof vi.fn>;
  waitForReceiptMock: ReturnType<typeof vi.fn>;
}

function makeProvider(
  opts: {
    receipt?: TxReceipt;
    address?: EvmAddress;
  } = {}
): MockProvider {
  const writeContractMock = vi
    .fn<(params: ContractWriteParams) => Promise<TxHash>>()
    .mockResolvedValue(TX_HASH);
  const waitForReceiptMock = vi
    .fn<(txHash: TxHash) => Promise<TxReceipt>>()
    .mockResolvedValue(opts.receipt ?? makeReceipt("success"));
  const provider: MockProvider = {
    getChainId: async () => 84532,
    getAddress: async () => opts.address ?? SIGNER,
    readContract: async () => undefined as never,
    writeContract: writeContractMock,
    waitForReceipt: waitForReceiptMock,
    writeContractMock,
    waitForReceiptMock,
  };
  return provider;
}

describe("AcpAgent.constructor", () => {
  it("rejects zero AgentRegistry address", () => {
    expect(
      () =>
        new AcpAgent({
          provider: makeProvider(),
          gatewayUrl: "http://localhost:8080",
          contracts: { agentRegistry: ZERO, jobFactory: FACTORY },
        })
    ).toThrow(AcpError);
  });

  it("rejects zero JobFactory address", () => {
    expect(
      () =>
        new AcpAgent({
          provider: makeProvider(),
          gatewayUrl: "http://localhost:8080",
          contracts: { agentRegistry: REGISTRY, jobFactory: ZERO },
        })
    ).toThrow(AcpError);
  });

  it("uses an injected GatewayClient when provided", () => {
    const fakeGateway = {
      listAgents: vi.fn(),
    } as never;
    const agent = new AcpAgent({
      provider: makeProvider(),
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
      gateway: fakeGateway,
    });
    expect(agent.gateway).toBe(fakeGateway);
  });
});

describe("AcpAgent.register", () => {
  it("submits the tx, decodes AgentRegistered, and returns the new id", async () => {
    const log = encodeAgentRegistered(7n, SIGNER, "ipfs://meta", 5n);
    const provider = makeProvider({ receipt: makeReceipt("success", [log]) });
    const agent = new AcpAgent({
      provider,
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });

    const result = await agent.register({
      metadataUri: "ipfs://meta",
      capabilities: 5n,
    });

    expect(result.agentId).toBe(7n);
    expect(result.txHash).toBe(TX_HASH);
    expect(result.controller).toBe(SIGNER);
    expect(result.metadataUri).toBe("ipfs://meta");
    expect(result.capabilities).toBe(5n);

    expect(provider.writeContractMock).toHaveBeenCalledOnce();
    const call = provider.writeContractMock.mock.calls[0]?.[0];
    expect(call?.address).toBe(REGISTRY);
    expect(call?.functionName).toBe("register");
    expect(call?.args).toEqual([SIGNER, "ipfs://meta", 5n]);
  });

  it("accepts decimal-string capabilities", async () => {
    const log = encodeAgentRegistered(1n, SIGNER, "ipfs://m", 42n);
    const provider = makeProvider({ receipt: makeReceipt("success", [log]) });
    const agent = new AcpAgent({
      provider,
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });
    const result = await agent.register({
      metadataUri: "ipfs://m",
      capabilities: "42",
    });
    expect(result.capabilities).toBe(42n);
  });

  it("accepts safe number capabilities", async () => {
    const log = encodeAgentRegistered(1n, SIGNER, "ipfs://m", 99n);
    const provider = makeProvider({ receipt: makeReceipt("success", [log]) });
    const agent = new AcpAgent({
      provider,
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });
    const result = await agent.register({
      metadataUri: "ipfs://m",
      capabilities: 99,
    });
    expect(result.capabilities).toBe(99n);
  });

  it("rejects empty metadataUri", async () => {
    const agent = new AcpAgent({
      provider: makeProvider(),
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });
    await expect(agent.register({ metadataUri: "", capabilities: 0n })).rejects.toMatchObject({
      code: "invalid_param",
    });
  });

  it("rejects negative number capabilities", async () => {
    const agent = new AcpAgent({
      provider: makeProvider(),
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });
    await expect(
      agent.register({ metadataUri: "ipfs://m", capabilities: -1 })
    ).rejects.toMatchObject({ code: "invalid_param" });
  });

  it("rejects number capabilities > 2^53", async () => {
    const agent = new AcpAgent({
      provider: makeProvider(),
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });
    await expect(
      agent.register({
        metadataUri: "ipfs://m",
        // Number.MAX_SAFE_INTEGER + 1 — but JS coerces to a float, so we go far past.
        capabilities: Number.MAX_SAFE_INTEGER * 2,
      })
    ).rejects.toMatchObject({ code: "invalid_param" });
  });

  it("rejects non-decimal-string capabilities", async () => {
    const agent = new AcpAgent({
      provider: makeProvider(),
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });
    await expect(
      agent.register({ metadataUri: "ipfs://m", capabilities: "0xff" })
    ).rejects.toMatchObject({ code: "invalid_param" });
  });

  it("throws tx_reverted when the receipt failed", async () => {
    const provider = makeProvider({ receipt: makeReceipt("reverted") });
    const agent = new AcpAgent({
      provider,
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });
    await expect(
      agent.register({ metadataUri: "ipfs://m", capabilities: 0n })
    ).rejects.toMatchObject({ code: "tx_reverted" });
  });

  it("throws event_not_found when AgentRegistered is missing", async () => {
    const provider = makeProvider({ receipt: makeReceipt("success", []) });
    const agent = new AcpAgent({
      provider,
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });
    await expect(
      agent.register({ metadataUri: "ipfs://m", capabilities: 0n })
    ).rejects.toMatchObject({ code: "event_not_found" });
  });

  it("uses an explicit controller when provided", async () => {
    const explicit: EvmAddress = "0xabcdef0000000000000000000000000000000001";
    const log = encodeAgentRegistered(2n, explicit, "ipfs://m", 0n);
    const provider = makeProvider({ receipt: makeReceipt("success", [log]) });
    const agent = new AcpAgent({
      provider,
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });
    const result = await agent.register({
      controller: explicit,
      metadataUri: "ipfs://m",
      capabilities: 0n,
    });
    expect(result.controller).toBe(explicit);
    const call = provider.writeContractMock.mock.calls[0]?.[0];
    expect(call?.args[0]).toBe(explicit);
  });
});

describe("AcpAgent.browse", () => {
  it("delegates to the gateway with the given filter", async () => {
    const listAgents = vi.fn().mockResolvedValue({
      data: [],
      next_cursor: "cur",
    });
    const agent = new AcpAgent({
      provider: makeProvider(),
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
      gateway: { listAgents } as never,
    });
    const page = await agent.browse({ kind: "acp", limit: 50 });
    expect(page.next_cursor).toBe("cur");
    expect(listAgents).toHaveBeenCalledWith({ kind: "acp", limit: 50 });
  });

  it("supports a default empty filter", async () => {
    const listAgents = vi.fn().mockResolvedValue({ data: [], next_cursor: "" });
    const agent = new AcpAgent({
      provider: makeProvider(),
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
      gateway: { listAgents } as never,
    });
    await agent.browse();
    expect(listAgents).toHaveBeenCalledWith({});
  });

  it.each([0, 201, 1.5])("rejects invalid limit %s", async (badLimit) => {
    const agent = new AcpAgent({
      provider: makeProvider(),
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
      gateway: { listAgents: vi.fn() } as never,
    });
    await expect(agent.browse({ limit: badLimit })).rejects.toMatchObject({
      code: "invalid_param",
    });
  });
});

describe("AcpAgent.createJob", () => {
  const futureDeadline = BigInt(Math.floor(Date.now() / 1000) + 3_600);

  it("submits a Direct job and decodes JobCreated", async () => {
    const log = encodeJobCreated(11n, CLONE, SIGNER, 0, TOKEN, 1_000_000n);
    const provider = makeProvider({ receipt: makeReceipt("success", [log]) });
    const agent = new AcpAgent({
      provider,
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });
    const result = await agent.createJob({
      token: TOKEN,
      budget: 1_000_000n,
      deadline: futureDeadline,
    });
    expect(result.jobId).toBe(11n);
    expect(result.cloneAddress).toBe(CLONE);
    expect(result.txHash).toBe(TX_HASH);

    const call = provider.writeContractMock.mock.calls[0]?.[0];
    expect(call?.address).toBe(FACTORY);
    expect(call?.functionName).toBe("createJob");
    expect(call?.args).toEqual([
      0n,
      TOKEN,
      1_000_000n,
      futureDeadline,
      JobType.Direct,
      ZERO,
      ZERO,
      [],
    ]);
  });

  it("submits an Evaluated job with required evaluator", async () => {
    const log = encodeJobCreated(12n, CLONE, SIGNER, 1, TOKEN, 100n);
    const provider = makeProvider({ receipt: makeReceipt("success", [log]) });
    const agent = new AcpAgent({
      provider,
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });
    const result = await agent.createJob({
      token: TOKEN,
      budget: 100n,
      deadline: futureDeadline,
      jobType: JobType.Evaluated,
      evaluator: EVALUATOR,
      targetAgentId: 9n,
      hooks: [SIGNER],
    });
    expect(result.jobId).toBe(12n);
    const call = provider.writeContractMock.mock.calls[0]?.[0];
    expect(call?.args[0]).toBe(9n);
    expect(call?.args[4]).toBe(JobType.Evaluated);
    expect(call?.args[5]).toBe(EVALUATOR);
    expect(call?.args[7]).toEqual([SIGNER]);
  });

  it("rejects zero token", async () => {
    const agent = new AcpAgent({
      provider: makeProvider(),
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });
    await expect(
      agent.createJob({ token: ZERO, budget: 1n, deadline: futureDeadline })
    ).rejects.toMatchObject({ code: "invalid_param" });
  });

  it("rejects zero budget", async () => {
    const agent = new AcpAgent({
      provider: makeProvider(),
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });
    await expect(
      agent.createJob({ token: TOKEN, budget: 0n, deadline: futureDeadline })
    ).rejects.toMatchObject({ code: "invalid_param" });
  });

  it("rejects non-positive deadline", async () => {
    const agent = new AcpAgent({
      provider: makeProvider(),
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });
    await expect(agent.createJob({ token: TOKEN, budget: 1n, deadline: 0n })).rejects.toMatchObject(
      { code: "invalid_param" }
    );
  });

  it("rejects Evaluated jobType without evaluator", async () => {
    const agent = new AcpAgent({
      provider: makeProvider(),
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });
    await expect(
      agent.createJob({
        token: TOKEN,
        budget: 1n,
        deadline: futureDeadline,
        jobType: JobType.Evaluated,
      })
    ).rejects.toMatchObject({ code: "invalid_param" });
  });

  it("rejects Evaluated jobType with zero evaluator", async () => {
    const agent = new AcpAgent({
      provider: makeProvider(),
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });
    await expect(
      agent.createJob({
        token: TOKEN,
        budget: 1n,
        deadline: futureDeadline,
        jobType: JobType.Evaluated,
        evaluator: ZERO,
      })
    ).rejects.toMatchObject({ code: "invalid_param" });
  });

  it("throws tx_reverted when the tx fails", async () => {
    const provider = makeProvider({ receipt: makeReceipt("reverted") });
    const agent = new AcpAgent({
      provider,
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });
    await expect(
      agent.createJob({
        token: TOKEN,
        budget: 1n,
        deadline: futureDeadline,
      })
    ).rejects.toMatchObject({ code: "tx_reverted" });
  });

  it("throws event_not_found when JobCreated is missing", async () => {
    const provider = makeProvider({ receipt: makeReceipt("success", []) });
    const agent = new AcpAgent({
      provider,
      gatewayUrl: "http://localhost:8080",
      contracts: { agentRegistry: REGISTRY, jobFactory: FACTORY },
    });
    await expect(
      agent.createJob({
        token: TOKEN,
        budget: 1n,
        deadline: futureDeadline,
      })
    ).rejects.toMatchObject({ code: "event_not_found" });
  });
});
