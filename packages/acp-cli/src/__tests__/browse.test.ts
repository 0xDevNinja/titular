import { describe, expect, it, vi } from "vitest";

import { runBrowse } from "../commands/browse.js";
import { saveConfig } from "../config.js";

import { CollectingStream, freshConfigPath } from "./_helpers.js";

const BASE_CFG = {
  gatewayUrl: "https://api.titular.dev",
  chain: "base-sepolia" as const,
};

describe("runBrowse", () => {
  it("calls listAgents with no kind by default", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const listAgents = vi.fn().mockResolvedValue({ data: [], next_cursor: "" });
    await runBrowse(
      { config: path },
      {
        gatewayFactory: () => ({ listAgents }) as never,
        env: { NO_COLOR: "1" },
        isTTY: false,
      }
    );
    expect(listAgents).toHaveBeenCalledWith({});
  });

  it("forwards a valid kind", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const listAgents = vi.fn().mockResolvedValue({ data: [], next_cursor: "" });
    await runBrowse(
      { config: path, kind: "launchpad" },
      {
        gatewayFactory: () => ({ listAgents }) as never,
        env: { NO_COLOR: "1" },
        isTTY: false,
      }
    );
    expect(listAgents).toHaveBeenCalledWith({ kind: "launchpad" });
  });

  it("rejects an unknown kind", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    await expect(
      runBrowse(
        { config: path, kind: "spaceship" },
        {
          gatewayFactory: () => ({ listAgents: vi.fn() }) as never,
          env: { NO_COLOR: "1" },
          isTTY: false,
        }
      )
    ).rejects.toThrow(/--kind must be one of/);
  });

  it("emits next_cursor when set", async () => {
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
      next_cursor: "abc",
    });
    const stdout = new CollectingStream();
    await runBrowse(
      { config: path },
      {
        gatewayFactory: () => ({ listAgents }) as never,
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: { NO_COLOR: "1" },
        isTTY: false,
      }
    );
    expect(stdout.text()).toContain("next_cursor: abc");
  });
});
