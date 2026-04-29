import { describe, expect, it, vi } from "vitest";

import {
  AgentError,
  __testing__,
  runAgentCreate,
  runAgentList,
  runAgentStatus,
} from "../commands/agent.js";
import { saveConfig } from "../config.js";

import { CollectingStream, freshConfigPath } from "./_helpers.js";

const BASE_CFG = {
  gatewayUrl: "https://api.titular.dev",
  chain: "base-sepolia" as const,
  apiKey: "k",
};

describe("parseCapabilities", () => {
  const { parseCapabilities } = __testing__;
  it("parses decimal", () => {
    expect(parseCapabilities("42")).toBe(42n);
  });
  it("parses hex", () => {
    expect(parseCapabilities("0xff")).toBe(255n);
  });
  it("rejects garbage", () => {
    expect(() => parseCapabilities("nope")).toThrow(/decimal or 0x-hex/);
  });
});

describe("runAgentCreate", () => {
  it("requires a metadata URI", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    await expect(
      runAgentCreate({ config: path }, { env: {}, isTTY: false })
    ).rejects.toBeInstanceOf(AgentError);
  });

  it("calls AcpAgent.register and prints jsonified result", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const register = vi.fn().mockResolvedValue({
      agentId: 7n,
      txHash: `0x${"ab".repeat(32)}`,
      controller: `0x${"cd".repeat(20)}`,
      metadataUri: "ipfs://x",
      capabilities: 0n,
    });
    const stdout = new CollectingStream();
    await runAgentCreate(
      {
        config: path,
        json: true,
        metadataUri: "ipfs://x",
        capabilities: "0",
      },
      {
        acpAgentFactory: () => ({ register }) as never,
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: {},
        isTTY: false,
      }
    );
    expect(register).toHaveBeenCalledWith({
      metadataUri: "ipfs://x",
      capabilities: 0n,
    });
    expect(stdout.text()).toContain('"agentId": "7"');
  });

  it("requires config to exist", async () => {
    const { path } = await freshConfigPath();
    await expect(
      runAgentCreate({ config: path, metadataUri: "ipfs://x" }, { env: {}, isTTY: false })
    ).rejects.toThrow(/no config/);
  });
});

describe("runAgentList", () => {
  it("renders a pretty table for non-empty pages", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const listAgents = vi.fn().mockResolvedValue({
      data: [
        {
          id: 1,
          agent_id: "1",
          kind: "acp",
          controller: `0x${"ab".repeat(20)}`,
          capabilities: "0",
          block_number: 100,
          log_index: 0,
          tx_hash: `0x${"cd".repeat(32)}`,
          created_at: "2026-01-01T00:00:00Z",
          updated_at: "2026-01-01T00:00:00Z",
        },
      ],
      next_cursor: "",
    });
    const stdout = new CollectingStream();
    const page = await runAgentList(
      { config: path },
      {
        gatewayFactory: () => ({ listAgents }) as never,
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: { NO_COLOR: "1" },
        isTTY: false,
      }
    );
    expect(page.data).toHaveLength(1);
    expect(stdout.text()).toContain("agent_id");
    expect(stdout.text()).toContain("acp");
  });

  it("emits empty marker when there are no agents", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const listAgents = vi.fn().mockResolvedValue({ data: [], next_cursor: "" });
    const stdout = new CollectingStream();
    await runAgentList(
      { config: path },
      {
        gatewayFactory: () => ({ listAgents }) as never,
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: { NO_COLOR: "1" },
        isTTY: false,
      }
    );
    expect(stdout.text()).toContain("(no agents)");
  });

  it("forwards limit + cursor to the gateway", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const listAgents = vi.fn().mockResolvedValue({ data: [], next_cursor: "" });
    await runAgentList(
      { config: path, limit: 5, cursor: "abc" },
      {
        gatewayFactory: () => ({ listAgents }) as never,
        env: { NO_COLOR: "1" },
        isTTY: false,
      }
    );
    expect(listAgents).toHaveBeenCalledWith({
      kind: "acp",
      limit: 5,
      cursor: "abc",
    });
  });
});

describe("runAgentStatus", () => {
  it("prints agent details", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const getAgent = vi.fn().mockResolvedValue({
      id: 1,
      agent_id: "42",
      kind: "acp",
      controller: `0x${"ab".repeat(20)}`,
      capabilities: "0",
      block_number: 100,
      log_index: 0,
      tx_hash: `0x${"cd".repeat(32)}`,
      created_at: "2026-01-01T00:00:00Z",
      updated_at: "2026-01-01T00:00:00Z",
    });
    const stdout = new CollectingStream();
    await runAgentStatus(
      { config: path, agentId: "42", json: true },
      {
        gatewayFactory: () => ({ getAgent }) as never,
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: {},
        isTTY: false,
      }
    );
    expect(getAgent).toHaveBeenCalledWith("42");
    expect(stdout.text()).toContain('"agent_id": "42"');
  });

  it("rejects an empty agentId", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    await expect(
      runAgentStatus(
        { config: path, agentId: "" },
        {
          gatewayFactory: () => ({ getAgent: vi.fn() }) as never,
          env: {},
          isTTY: false,
        }
      )
    ).rejects.toBeInstanceOf(AgentError);
  });

  it("prints pretty kv output by default and includes optional creator + metadata fields", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const getAgent = vi.fn().mockResolvedValue({
      id: 1,
      agent_id: "42",
      kind: "acp",
      controller: `0x${"ab".repeat(20)}`,
      creator: `0x${"cd".repeat(20)}`,
      capabilities: "0",
      metadata_uri: "ipfs://m",
      block_number: 100,
      log_index: 0,
      tx_hash: `0x${"ef".repeat(32)}`,
      created_at: "2026-01-01T00:00:00Z",
      updated_at: "2026-01-01T00:00:00Z",
    });
    const stdout = new CollectingStream();
    await runAgentStatus(
      { config: path, agentId: "42" },
      {
        gatewayFactory: () => ({ getAgent }) as never,
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: { NO_COLOR: "1" },
        isTTY: false,
      }
    );
    const text = stdout.text();
    expect(text).toContain("agent_id:");
    expect(text).toContain("creator:");
    expect(text).toContain("metadata_uri:");
    expect(text).toContain("ipfs://m");
  });
});

describe("runAgentCreate pretty mode", () => {
  it("prints kv output and the optional name when --json is off", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const register = vi.fn().mockResolvedValue({
      agentId: 7n,
      txHash: `0x${"ab".repeat(32)}`,
      controller: `0x${"cd".repeat(20)}`,
      metadataUri: "ipfs://x",
      capabilities: 0n,
    });
    const stdout = new CollectingStream();
    await runAgentCreate(
      {
        config: path,
        metadataUri: "ipfs://x",
        capabilities: "0xff",
        name: "alpha-bot",
        controller: `0x${"cd".repeat(20)}`,
      },
      {
        acpAgentFactory: () => ({ register }) as never,
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: { NO_COLOR: "1" },
        isTTY: false,
      }
    );
    expect(register).toHaveBeenCalledWith({
      metadataUri: "ipfs://x",
      capabilities: 255n,
      controller: `0x${"cd".repeat(20)}`,
    });
    const text = stdout.text();
    expect(text).toContain("agentId:");
    expect(text).toContain("name:");
    expect(text).toContain("alpha-bot");
  });

  it("rejects malformed capability bitmaps", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    await expect(
      runAgentCreate(
        {
          config: path,
          metadataUri: "ipfs://x",
          capabilities: "garbage",
        },
        { env: {}, isTTY: false }
      )
    ).rejects.toBeInstanceOf(AgentError);
  });
});
