import { describe, expect, it, vi } from "vitest";
import type { AcpAgent } from "../agent.js";
import { AcpError } from "../errors.js";
import { executeTool } from "../llm/executor.js";
import { TOOL_NAMES } from "../llm/tools.js";
import type { ToolCall } from "../llm/types.js";
import {
  type CreateJobResult,
  type EvmAddress,
  JobType,
  type RegisterResult,
  type TxHash,
} from "../types.js";

const SIGNER: EvmAddress = "0x1111111111111111111111111111111111111111";
const TOKEN: EvmAddress = "0x2222222222222222222222222222222222222222";
const CLONE: EvmAddress = "0x3333333333333333333333333333333333333333";
const EVALUATOR: EvmAddress = "0x4444444444444444444444444444444444444444";
const TX: TxHash = "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";

interface MockAgent {
  register: ReturnType<typeof vi.fn>;
  browse: ReturnType<typeof vi.fn>;
  createJob: ReturnType<typeof vi.fn>;
}

function makeAgent(): MockAgent {
  return {
    register: vi.fn(),
    browse: vi.fn(),
    createJob: vi.fn(),
  };
}

function call(name: string, args: unknown): ToolCall {
  return {
    id: "call_1",
    type: "function",
    function: {
      name,
      arguments: typeof args === "string" ? args : JSON.stringify(args),
    },
  };
}

function asAcpAgent(m: MockAgent): AcpAgent {
  return m as unknown as AcpAgent;
}

describe("executeTool — register", () => {
  const okResult: RegisterResult = {
    agentId: 7n,
    txHash: TX,
    controller: SIGNER,
    metadataUri: "ipfs://meta",
    capabilities: 5n,
  };

  it("dispatches register and serialises BigInts as decimal strings", async () => {
    const agent = makeAgent();
    agent.register.mockResolvedValue(okResult);

    const res = await executeTool(
      asAcpAgent(agent),
      call(TOOL_NAMES.register, {
        metadataUri: "ipfs://meta",
        capabilities: "5",
      })
    );

    expect(res.isError).toBe(false);
    expect(res.role).toBe("tool");
    expect(res.tool_call_id).toBe("call_1");
    expect(res.name).toBe(TOOL_NAMES.register);
    const payload = JSON.parse(res.content);
    expect(payload.agentId).toBe("7");
    expect(payload.capabilities).toBe("5");
    expect(payload.txHash).toBe(TX);
    expect(payload.controller).toBe(SIGNER);
    expect(agent.register).toHaveBeenCalledWith({
      metadataUri: "ipfs://meta",
      capabilities: "5",
    });
  });

  it("forwards an explicit controller and lowercases it", async () => {
    const agent = makeAgent();
    agent.register.mockResolvedValue(okResult);

    await executeTool(
      asAcpAgent(agent),
      call(TOOL_NAMES.register, {
        metadataUri: "ipfs://meta",
        capabilities: "5",
        controller: "0xABCDEF0000000000000000000000000000000001",
      })
    );

    expect(agent.register).toHaveBeenCalledWith({
      metadataUri: "ipfs://meta",
      capabilities: "5",
      controller: "0xabcdef0000000000000000000000000000000001",
    });
  });

  it("returns an error envelope for non-decimal capabilities", async () => {
    const agent = makeAgent();
    const res = await executeTool(
      asAcpAgent(agent),
      call(TOOL_NAMES.register, {
        metadataUri: "ipfs://meta",
        capabilities: "0xff",
      })
    );
    expect(res.isError).toBe(true);
    expect(JSON.parse(res.content).error.message).toMatch(/decimal/);
    expect(agent.register).not.toHaveBeenCalled();
  });

  it("returns an error envelope when metadataUri is missing", async () => {
    const agent = makeAgent();
    const res = await executeTool(
      asAcpAgent(agent),
      call(TOOL_NAMES.register, { capabilities: "5" })
    );
    expect(res.isError).toBe(true);
  });

  it("returns an error envelope when controller is malformed", async () => {
    const agent = makeAgent();
    const res = await executeTool(
      asAcpAgent(agent),
      call(TOOL_NAMES.register, {
        metadataUri: "ipfs://meta",
        capabilities: "5",
        controller: "not-an-address",
      })
    );
    expect(res.isError).toBe(true);
    expect(JSON.parse(res.content).error.message).toMatch(/address/);
  });

  it("propagates AcpError code from the agent", async () => {
    const agent = makeAgent();
    agent.register.mockRejectedValue(new AcpError("tx_reverted", "tx reverted: 0xabc"));
    const res = await executeTool(
      asAcpAgent(agent),
      call(TOOL_NAMES.register, { metadataUri: "ipfs://m", capabilities: "0" })
    );
    expect(res.isError).toBe(true);
    expect(JSON.parse(res.content).error.code).toBe("tx_reverted");
  });

  it("propagates non-AcpError as invalid_param fallback", async () => {
    const agent = makeAgent();
    agent.register.mockRejectedValue(new Error("network down"));
    const res = await executeTool(
      asAcpAgent(agent),
      call(TOOL_NAMES.register, { metadataUri: "ipfs://m", capabilities: "0" })
    );
    expect(res.isError).toBe(true);
    expect(JSON.parse(res.content).error.code).toBe("invalid_param");
    expect(JSON.parse(res.content).error.message).toBe("network down");
  });
});

describe("executeTool — browse", () => {
  it("dispatches browse with full filter", async () => {
    const agent = makeAgent();
    agent.browse.mockResolvedValue({ data: [], next_cursor: "next" });
    const res = await executeTool(
      asAcpAgent(agent),
      call(TOOL_NAMES.browse, { kind: "acp", limit: 10, cursor: "abc" })
    );
    expect(res.isError).toBe(false);
    expect(JSON.parse(res.content)).toEqual({ data: [], next_cursor: "next" });
    expect(agent.browse).toHaveBeenCalledWith({
      kind: "acp",
      limit: 10,
      cursor: "abc",
    });
  });

  it("dispatches browse with no filter", async () => {
    const agent = makeAgent();
    agent.browse.mockResolvedValue({ data: [], next_cursor: "" });
    await executeTool(asAcpAgent(agent), call(TOOL_NAMES.browse, {}));
    expect(agent.browse).toHaveBeenCalledWith({});
  });

  it("rejects non-enum kind", async () => {
    const agent = makeAgent();
    const res = await executeTool(asAcpAgent(agent), call(TOOL_NAMES.browse, { kind: "other" }));
    expect(res.isError).toBe(true);
    expect(agent.browse).not.toHaveBeenCalled();
  });

  it("rejects non-integer limit", async () => {
    const agent = makeAgent();
    const res = await executeTool(asAcpAgent(agent), call(TOOL_NAMES.browse, { limit: 1.5 }));
    expect(res.isError).toBe(true);
  });
});

describe("executeTool — createJob", () => {
  const okResult: CreateJobResult = {
    jobId: 11n,
    cloneAddress: CLONE,
    txHash: TX,
  };

  it("dispatches a Direct job and coerces strings to bigints", async () => {
    const agent = makeAgent();
    agent.createJob.mockResolvedValue(okResult);
    const res = await executeTool(
      asAcpAgent(agent),
      call(TOOL_NAMES.createJob, {
        token: TOKEN,
        budget: "1000000",
        deadline: "2000000000",
      })
    );
    expect(res.isError).toBe(false);
    const payload = JSON.parse(res.content);
    expect(payload.jobId).toBe("11");
    expect(payload.cloneAddress).toBe(CLONE);
    expect(agent.createJob).toHaveBeenCalledWith({
      token: TOKEN,
      budget: 1000000n,
      deadline: 2000000000n,
    });
  });

  it("dispatches an Evaluated job with evaluator + targetAgentId + hooks", async () => {
    const agent = makeAgent();
    agent.createJob.mockResolvedValue(okResult);
    await executeTool(
      asAcpAgent(agent),
      call(TOOL_NAMES.createJob, {
        token: TOKEN,
        budget: "100",
        deadline: "2000000000",
        jobType: "evaluated",
        evaluator: EVALUATOR,
        targetAgentId: "9",
        hooks: [SIGNER],
      })
    );
    expect(agent.createJob).toHaveBeenCalledWith({
      token: TOKEN,
      budget: 100n,
      deadline: 2000000000n,
      jobType: JobType.Evaluated,
      evaluator: EVALUATOR,
      targetAgentId: 9n,
      hooks: [SIGNER],
    });
  });

  it("translates jobType=direct to JobType.Direct", async () => {
    const agent = makeAgent();
    agent.createJob.mockResolvedValue(okResult);
    await executeTool(
      asAcpAgent(agent),
      call(TOOL_NAMES.createJob, {
        token: TOKEN,
        budget: "1",
        deadline: "2000000000",
        jobType: "direct",
      })
    );
    expect(agent.createJob.mock.calls[0]?.[0].jobType).toBe(JobType.Direct);
  });

  it("rejects unknown jobType", async () => {
    const agent = makeAgent();
    const res = await executeTool(
      asAcpAgent(agent),
      call(TOOL_NAMES.createJob, {
        token: TOKEN,
        budget: "1",
        deadline: "2000000000",
        jobType: "garbage",
      })
    );
    expect(res.isError).toBe(true);
  });

  it("rejects malformed token", async () => {
    const agent = makeAgent();
    const res = await executeTool(
      asAcpAgent(agent),
      call(TOOL_NAMES.createJob, {
        token: "0xnope",
        budget: "1",
        deadline: "2000000000",
      })
    );
    expect(res.isError).toBe(true);
  });

  it("rejects malformed hook entry", async () => {
    const agent = makeAgent();
    const res = await executeTool(
      asAcpAgent(agent),
      call(TOOL_NAMES.createJob, {
        token: TOKEN,
        budget: "1",
        deadline: "2000000000",
        hooks: ["not-an-address"],
      })
    );
    expect(res.isError).toBe(true);
  });

  it("rejects non-array hooks", async () => {
    const agent = makeAgent();
    const res = await executeTool(
      asAcpAgent(agent),
      call(TOOL_NAMES.createJob, {
        token: TOKEN,
        budget: "1",
        deadline: "2000000000",
        hooks: SIGNER,
      })
    );
    expect(res.isError).toBe(true);
  });

  it("propagates AcpError from createJob", async () => {
    const agent = makeAgent();
    agent.createJob.mockRejectedValue(new AcpError("invalid_param", "budget must be > 0"));
    const res = await executeTool(
      asAcpAgent(agent),
      call(TOOL_NAMES.createJob, {
        token: TOKEN,
        budget: "1",
        deadline: "2000000000",
      })
    );
    expect(res.isError).toBe(true);
    expect(JSON.parse(res.content).error.code).toBe("invalid_param");
  });
});

describe("executeTool — argument parsing", () => {
  it("returns an error envelope for invalid JSON arguments", async () => {
    const agent = makeAgent();
    const res = await executeTool(asAcpAgent(agent), call(TOOL_NAMES.register, "{not json"));
    expect(res.isError).toBe(true);
    expect(JSON.parse(res.content).error.message).toMatch(/JSON/);
  });

  it("returns an error envelope when arguments are an array", async () => {
    const agent = makeAgent();
    const res = await executeTool(asAcpAgent(agent), call(TOOL_NAMES.register, [1, 2, 3]));
    expect(res.isError).toBe(true);
    expect(JSON.parse(res.content).error.message).toMatch(/object/);
  });

  it("treats empty argument string as empty object", async () => {
    const agent = makeAgent();
    const res = await executeTool(asAcpAgent(agent), {
      id: "x",
      type: "function",
      function: { name: TOOL_NAMES.register, arguments: "" },
    });
    // metadataUri is required, so this surfaces as an arg-parsing error,
    // but importantly NOT as a "not valid JSON" error.
    expect(res.isError).toBe(true);
    expect(JSON.parse(res.content).error.message).not.toMatch(/JSON/);
  });

  it("returns an error envelope for unknown tool name", async () => {
    const agent = makeAgent();
    const res = await executeTool(asAcpAgent(agent), call("acp_does_not_exist", {}));
    expect(res.isError).toBe(true);
    expect(JSON.parse(res.content).error.message).toMatch(/unknown tool/);
  });
});
