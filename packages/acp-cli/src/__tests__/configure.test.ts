import { Writable } from "node:stream";
import { describe, expect, it } from "vitest";

import { ConfigureError, runConfigure } from "../commands/configure.js";
import { loadConfig } from "../config.js";

import { CollectingStream, freshConfigPath, readableFromLines } from "./_helpers.js";

class SinkWritable extends Writable {
  override _write(_chunk: unknown, _enc: string, cb: (e?: Error | null) => void): void {
    cb();
  }
}

describe("runConfigure", () => {
  it("writes a fresh config in non-interactive mode", async () => {
    const { path } = await freshConfigPath();
    const stdout = new CollectingStream();
    const cfg = await runConfigure(
      {
        json: true,
        config: path,
        gatewayUrl: "https://api.titular.dev",
        chain: "base-sepolia",
        apiKey: "key123",
        nonInteractive: true,
      },
      {
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: {},
        isTTY: false,
      }
    );
    expect(cfg.gatewayUrl).toBe("https://api.titular.dev");
    expect(cfg.apiKey).toBe("key123");

    const reread = await loadConfig(path);
    expect(reread?.apiKey).toBe("key123");
    expect(stdout.text()).toContain('"saved": true');
    expect(stdout.text()).toContain('"apiKey": "***"');
  });

  it("refuses non-interactive without an existing config and no --gateway-url", async () => {
    const { path } = await freshConfigPath();
    await expect(
      runConfigure({ config: path, nonInteractive: true }, { env: {}, isTTY: false })
    ).rejects.toBeInstanceOf(ConfigureError);
  });

  it("rejects an unknown chain slug", async () => {
    const { path } = await freshConfigPath();
    await expect(
      runConfigure(
        {
          config: path,
          gatewayUrl: "https://x",
          chain: "polygon",
          nonInteractive: true,
        },
        { env: {}, isTTY: false }
      )
    ).rejects.toThrow(/--chain must be one of/);
  });

  it("--print on a fresh path is a typed error", async () => {
    const { path } = await freshConfigPath();
    await expect(
      runConfigure({ config: path, print: true }, { env: {}, isTTY: false })
    ).rejects.toBeInstanceOf(ConfigureError);
  });

  it("--print returns the redacted config", async () => {
    const { path } = await freshConfigPath();
    await runConfigure(
      {
        config: path,
        gatewayUrl: "https://x",
        chain: "base-sepolia",
        apiKey: "secret",
        nonInteractive: true,
      },
      { env: {}, isTTY: false }
    );
    const stdout = new CollectingStream();
    await runConfigure(
      { config: path, print: true, json: true },
      {
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: {},
        isTTY: false,
      }
    );
    expect(stdout.text()).toContain('"apiKey": "***"');
    expect(stdout.text()).not.toContain("secret");
  });

  it("interactive flow consumes prompts and writes the config", async () => {
    const { path } = await freshConfigPath();
    const promptIO = {
      input: readableFromLines([
        "https://api.example.com",
        "1", // chain selection (first option)
        "key-from-prompt",
      ]),
      output: new SinkWritable(),
    };
    const stdout = new CollectingStream();
    const cfg = await runConfigure(
      { config: path, json: true },
      {
        promptIO,
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: {},
        isTTY: false,
      }
    );
    expect(cfg.gatewayUrl).toBe("https://api.example.com");
    expect(cfg.apiKey).toBe("key-from-prompt");
  });

  it("preserves existing fields when re-running with a partial flag set", async () => {
    const { path } = await freshConfigPath();
    await runConfigure(
      {
        config: path,
        gatewayUrl: "https://first",
        chain: "base-sepolia",
        apiKey: "k1",
        nonInteractive: true,
      },
      { env: {}, isTTY: false }
    );
    const cfg = await runConfigure(
      {
        config: path,
        chain: "base-mainnet",
        nonInteractive: true,
      },
      { env: {}, isTTY: false }
    );
    expect(cfg.gatewayUrl).toBe("https://first");
    expect(cfg.chain).toBe("base-mainnet");
    expect(cfg.apiKey).toBe("k1");
  });
});
