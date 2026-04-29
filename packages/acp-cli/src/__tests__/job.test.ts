import { describe, expect, it, vi } from "vitest";

import { JobError, __testing__, runJobCreate, runJobList } from "../commands/job.js";
import { saveConfig } from "../config.js";

import { CollectingStream, freshConfigPath } from "./_helpers.js";

const BASE_CFG = {
  gatewayUrl: "https://api.titular.dev",
  chain: "base-sepolia" as const,
  apiKey: "k",
};

const TOKEN = `0x${"ab".repeat(20)}` as const;

describe("parseBigInt / parseJobPhase / parseJobType", () => {
  const { parseBigInt, parseJobPhase, parseJobType } = __testing__;

  it("parses a decimal string", () => {
    expect(parseBigInt("123", "--budget")).toBe(123n);
  });

  it("rejects non-decimal", () => {
    expect(() => parseBigInt("0x1", "--budget")).toThrow(/decimal/);
  });

  it("accepts known phases", () => {
    expect(parseJobPhase("created")).toBe("created");
    expect(parseJobPhase("FUNDED")).toBe("funded");
  });

  it("rejects unknown phases", () => {
    expect(() => parseJobPhase("nope")).toThrow(/--status/);
  });

  it("maps job type aliases", () => {
    expect(parseJobType(undefined)).toBe(0);
    expect(parseJobType("direct")).toBe(0);
    expect(parseJobType("evaluated")).toBe(1);
    expect(parseJobType("0")).toBe(0);
    expect(parseJobType("1")).toBe(1);
  });

  it("rejects unknown job type", () => {
    expect(() => parseJobType("blast")).toThrow(/--job-type/);
  });
});

describe("runJobCreate", () => {
  it("validates --token", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    await expect(
      runJobCreate(
        { config: path, budget: "1", deadline: "1", token: "0xnope" },
        { env: {}, isTTY: false }
      )
    ).rejects.toBeInstanceOf(JobError);
  });

  it("rejects budget == 0", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    await expect(
      runJobCreate(
        { config: path, budget: "0", deadline: "1", token: TOKEN },
        { env: {}, isTTY: false }
      )
    ).rejects.toThrow(/--budget/);
  });

  it("rejects deadline == 0", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    await expect(
      runJobCreate(
        { config: path, budget: "1", deadline: "0", token: TOKEN },
        { env: {}, isTTY: false }
      )
    ).rejects.toThrow(/--deadline/);
  });

  it("invokes AcpAgent.createJob and emits JSON when requested", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const createJob = vi.fn().mockResolvedValue({
      jobId: 11n,
      cloneAddress: `0x${"cd".repeat(20)}`,
      txHash: `0x${"ef".repeat(32)}`,
    });
    const stdout = new CollectingStream();
    await runJobCreate(
      {
        config: path,
        json: true,
        token: TOKEN,
        budget: "1000",
        deadline: "9999999999",
        target: "5",
        jobType: "direct",
      },
      {
        acpAgentFactory: () => ({ createJob }) as never,
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: {},
        isTTY: false,
      }
    );
    expect(createJob).toHaveBeenCalledWith({
      targetAgentId: 5n,
      token: TOKEN,
      budget: 1000n,
      deadline: 9999999999n,
      jobType: 0,
    });
    expect(stdout.text()).toContain('"jobId": "11"');
  });

  it("requires config to exist", async () => {
    const { path } = await freshConfigPath();
    await expect(
      runJobCreate(
        { config: path, budget: "1", deadline: "1", token: TOKEN },
        { env: {}, isTTY: false }
      )
    ).rejects.toThrow(/no config/);
  });

  it("renders pretty kv output when --json is off", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const createJob = vi.fn().mockResolvedValue({
      jobId: 1n,
      cloneAddress: `0x${"a".repeat(40)}`,
      txHash: `0x${"b".repeat(64)}`,
    });
    const stdout = new CollectingStream();
    await runJobCreate(
      {
        config: path,
        token: TOKEN,
        budget: "1",
        deadline: "1",
        terms: "bake a cake",
      },
      {
        acpAgentFactory: () => ({ createJob }) as never,
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: { NO_COLOR: "1" },
        isTTY: false,
      }
    );
    const text = stdout.text();
    expect(text).toContain("jobId:");
    expect(text).toContain("terms:");
    expect(text).toContain("bake a cake");
  });
});

describe("runJobList", () => {
  it("renders an empty marker when no jobs", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const listJobs = vi.fn().mockResolvedValue({ data: [], next_cursor: "" });
    const stdout = new CollectingStream();
    await runJobList(
      { config: path },
      {
        gatewayFactory: () => ({ listJobs }) as never,
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: { NO_COLOR: "1" },
        isTTY: false,
      }
    );
    expect(stdout.text()).toContain("(no jobs)");
  });

  it("forwards status + pagination to the gateway", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const listJobs = vi.fn().mockResolvedValue({ data: [], next_cursor: "abc" });
    const stdout = new CollectingStream();
    await runJobList(
      { config: path, status: "funded", limit: 25, cursor: "c" },
      {
        gatewayFactory: () => ({ listJobs }) as never,
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: { NO_COLOR: "1" },
        isTTY: false,
      }
    );
    expect(listJobs).toHaveBeenCalledWith({
      status: "funded",
      limit: 25,
      cursor: "c",
    });
    expect(stdout.text()).toContain("next_cursor: abc");
  });

  it("renders a table for non-empty pages", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const listJobs = vi.fn().mockResolvedValue({
      data: [
        {
          id: 1,
          on_chain_id: "1",
          agent_id: "0",
          phase: "created",
          principal: `0x${"a".repeat(40)}`,
          token: TOKEN,
          budget: "1000",
          deadline: "999",
          job_type: 0,
          block_number: 100,
          log_index: 0,
          tx_hash: `0x${"b".repeat(64)}`,
          created_at: "2026-01-01T00:00:00Z",
          updated_at: "2026-01-01T00:00:00Z",
        },
      ],
      next_cursor: "",
    });
    const stdout = new CollectingStream();
    await runJobList(
      { config: path },
      {
        gatewayFactory: () => ({ listJobs }) as never,
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: { NO_COLOR: "1" },
        isTTY: false,
      }
    );
    expect(stdout.text()).toContain("phase");
    expect(stdout.text()).toContain("created");
  });
});
