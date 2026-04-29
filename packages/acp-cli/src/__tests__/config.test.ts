import { mkdir, readFile, writeFile } from "node:fs/promises";
import { dirname, join } from "node:path";
import { describe, expect, it } from "vitest";

import {
  configExists,
  defaultConfigPath,
  isSupportedChain,
  loadConfig,
  mergeConfig,
  redactConfig,
  saveConfig,
  validateConfig,
} from "../config.js";

import { freshConfigPath } from "./_helpers.js";

describe("defaultConfigPath", () => {
  it("honours TITULAR_CONFIG override", () => {
    expect(defaultConfigPath({ TITULAR_CONFIG: "/tmp/x.json" })).toBe("/tmp/x.json");
  });

  it("uses XDG_CONFIG_HOME when set on POSIX", () => {
    if (process.platform === "win32") return;
    const path = defaultConfigPath({ XDG_CONFIG_HOME: "/tmp/xdg" });
    expect(path).toBe("/tmp/xdg/titular/config.json");
  });

  it("falls back to ~/.config on POSIX when XDG_CONFIG_HOME is unset", () => {
    if (process.platform === "win32") return;
    const path = defaultConfigPath({});
    expect(path.endsWith("/.config/titular/config.json")).toBe(true);
  });
});

describe("isSupportedChain", () => {
  it.each([
    ["base-sepolia", true],
    ["base-mainnet", true],
    ["ethereum-sepolia", true],
    ["ethereum-mainnet", true],
    ["bsc", false],
    ["", false],
  ])("isSupportedChain(%s) === %s", (slug, expected) => {
    expect(isSupportedChain(slug)).toBe(expected);
  });
});

describe("validateConfig", () => {
  it("accepts a minimal config", () => {
    const cfg = validateConfig(
      { gatewayUrl: "https://api.titular.dev", chain: "base-sepolia" },
      "test"
    );
    expect(cfg.gatewayUrl).toBe("https://api.titular.dev");
    expect(cfg.chain).toBe("base-sepolia");
  });

  it("rejects missing gatewayUrl", () => {
    expect(() => validateConfig({ chain: "base-sepolia" }, "p")).toThrow(/gatewayUrl/);
  });

  it("rejects unknown chain", () => {
    expect(() => validateConfig({ gatewayUrl: "https://x", chain: "polygon" }, "p")).toThrow(
      /chain/
    );
  });

  it("rejects malformed agentRegistry", () => {
    expect(() =>
      validateConfig(
        {
          gatewayUrl: "https://x",
          chain: "base-sepolia",
          agentRegistry: "0xnope",
        },
        "p"
      )
    ).toThrow(/agentRegistry/);
  });

  it("rejects malformed privateKey", () => {
    expect(() =>
      validateConfig({ gatewayUrl: "https://x", chain: "base-sepolia", privateKey: "0x123" }, "p")
    ).toThrow(/privateKey/);
  });

  it("accepts a valid privateKey", () => {
    const pk = `0x${"ab".repeat(32)}`;
    const cfg = validateConfig(
      { gatewayUrl: "https://x", chain: "base-sepolia", privateKey: pk },
      "p"
    );
    expect(cfg.privateKey).toBe(pk);
  });

  it("rejects non-object input", () => {
    expect(() => validateConfig("nope", "p")).toThrow(/JSON object/);
  });
});

describe("loadConfig + saveConfig", () => {
  it("returns undefined when the file is absent", async () => {
    const { path } = await freshConfigPath();
    expect(await loadConfig(path)).toBeUndefined();
  });

  it("round-trips a saved config", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, {
      gatewayUrl: "https://api.titular.dev",
      chain: "base-sepolia",
      apiKey: "k1",
    });
    const reread = await loadConfig(path);
    expect(reread?.gatewayUrl).toBe("https://api.titular.dev");
    expect(reread?.apiKey).toBe("k1");
  });

  it("creates parent directories", async () => {
    const { dir } = await freshConfigPath();
    const path = join(dir, "nested", "deep", "config.json");
    await saveConfig(path, { gatewayUrl: "https://x", chain: "base-sepolia" });
    const stat = await readFile(path, "utf8");
    expect(stat).toContain("https://x");
  });

  it("writes with mode 0600", async () => {
    if (process.platform === "win32") return;
    const { path } = await freshConfigPath();
    await saveConfig(path, { gatewayUrl: "https://x", chain: "base-sepolia" });
    const { stat } = await import("node:fs/promises");
    const s = await stat(path);
    // Mask off file-type bits.
    expect(s.mode & 0o777).toBe(0o600);
  });

  it("rejects malformed JSON", async () => {
    const { path } = await freshConfigPath();
    await mkdir(dirname(path), { recursive: true });
    await writeFile(path, "{ not json", "utf8");
    await expect(loadConfig(path)).rejects.toThrow(/not valid JSON/);
  });
});

describe("mergeConfig", () => {
  it("requires gatewayUrl on a fresh config", () => {
    expect(() => mergeConfig(undefined, { chain: "base-sepolia" })).toThrow(/gatewayUrl/);
  });

  it("preserves existing fields when patch is partial", () => {
    const merged = mergeConfig(
      { gatewayUrl: "https://old", chain: "base-sepolia", apiKey: "k" },
      { chain: "base-mainnet" }
    );
    expect(merged.gatewayUrl).toBe("https://old");
    expect(merged.chain).toBe("base-mainnet");
    expect(merged.apiKey).toBe("k");
  });
});

describe("redactConfig", () => {
  it("masks apiKey and privateKey", () => {
    const r = redactConfig({
      gatewayUrl: "https://x",
      chain: "base-sepolia",
      apiKey: "secret",
      privateKey: `0x${"ab".repeat(32)}`,
    });
    expect(r.apiKey).toBe("***");
    expect(r.privateKey).toBe("***");
  });

  it("preserves non-secret fields", () => {
    const addr = `0x${"ab".repeat(20)}` as `0x${string}`;
    const r = redactConfig({
      gatewayUrl: "https://x",
      chain: "base-sepolia",
      agentRegistry: addr,
    });
    expect(r.agentRegistry).toBe(addr);
    expect(r.apiKey).toBeUndefined();
  });
});

describe("configExists", () => {
  it("returns false for missing file", async () => {
    const { path } = await freshConfigPath();
    expect(await configExists(path)).toBe(false);
  });

  it("returns true after saveConfig", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, { gatewayUrl: "https://x", chain: "base-sepolia" });
    expect(await configExists(path)).toBe(true);
  });
});
