// Persisted CLI configuration.
//
// Storage layout:
//   - `XDG_CONFIG_HOME` / `~/.config/titular/config.json` (POSIX)
//   - `%APPDATA%\titular\config.json` (Windows, via APPDATA)
//
// The file is plain JSON so a user can edit it by hand. Unknown keys are
// preserved on round-trip so a future CLI version writing extra fields does
// not lose data when an older CLI saves over them. The `mode` is set to 0600
// because the file may hold an `apiKey` and a private key — narrower than the
// default umask, but standard for credential-bearing dotfiles.
//
// Loading is forgiving: a missing file is treated as "no config yet" rather
// than a hard error so `acp configure` works on a fresh machine. A malformed
// JSON file is a hard error because silent fallback would mask user mistakes
// (a config they thought they saved but typoed).

import { mkdir, readFile, stat, writeFile } from "node:fs/promises";
import { homedir, platform } from "node:os";
import { dirname, join } from "node:path";

import type { CliConfig, SupportedChain } from "./types.js";

/** Map of supported chain slugs (single source of truth, used by the
 * `acp configure` validator and by `agent.ts` when building defaults). */
export const SUPPORTED_CHAINS: readonly SupportedChain[] = [
  "base-sepolia",
  "base-mainnet",
  "ethereum-sepolia",
  "ethereum-mainnet",
] as const;

/** Returns true iff `s` is a known chain slug. Narrows for callers. */
export function isSupportedChain(s: string): s is SupportedChain {
  return (SUPPORTED_CHAINS as readonly string[]).includes(s);
}

/** Resolve the default config path for the current host. The result is a
 * fully-qualified path even when the relevant env vars are unset, so callers
 * can always `dirname()` it for `mkdir -p`. */
export function defaultConfigPath(env: NodeJS.ProcessEnv = process.env): string {
  const override = env.TITULAR_CONFIG;
  if (override !== undefined && override.length > 0) return override;

  if (platform() === "win32") {
    const appdata = env.APPDATA;
    const base =
      appdata !== undefined && appdata.length > 0 ? appdata : join(homedir(), "AppData", "Roaming");
    return join(base, "titular", "config.json");
  }

  const xdg = env.XDG_CONFIG_HOME;
  const base = xdg !== undefined && xdg.length > 0 ? xdg : join(homedir(), ".config");
  return join(base, "titular", "config.json");
}

/** Load config from disk; returns `undefined` when the file does not exist.
 * Throws on malformed JSON, unreadable files, or invalid shape. */
export async function loadConfig(path: string): Promise<CliConfig | undefined> {
  let raw: string;
  try {
    raw = await readFile(path, "utf8");
  } catch (err) {
    if (isFsErrno(err) && err.code === "ENOENT") return undefined;
    throw err;
  }

  let parsed: unknown;
  try {
    parsed = JSON.parse(raw);
  } catch (cause) {
    const msg = cause instanceof Error ? cause.message : String(cause);
    throw new Error(`config file ${path} is not valid JSON: ${msg}`);
  }

  return validateConfig(parsed, path);
}

/** Persist `cfg` to `path`, creating parent directories as needed. The file
 * is written with 0600 permissions because it may carry secrets. */
export async function saveConfig(path: string, cfg: CliConfig): Promise<void> {
  await mkdir(dirname(path), { recursive: true });
  const json = `${JSON.stringify(cfg, null, 2)}\n`;
  await writeFile(path, json, { mode: 0o600 });
}

/** Validates `cfg` and throws a descriptive error when fields are missing or
 * mistyped. Used by `loadConfig` and by tests that want to assert the same
 * normalisation rules without going through the disk. */
export function validateConfig(value: unknown, path: string): CliConfig {
  if (typeof value !== "object" || value === null) {
    throw new Error(`config file ${path}: expected JSON object, got ${typeof value}`);
  }
  const obj = value as Record<string, unknown>;

  const gatewayUrl = obj.gatewayUrl;
  if (typeof gatewayUrl !== "string" || gatewayUrl.length === 0) {
    throw new Error(
      `config file ${path}: \`gatewayUrl\` is required and must be a non-empty string`
    );
  }

  const chain = obj.chain;
  if (typeof chain !== "string" || !isSupportedChain(chain)) {
    throw new Error(`config file ${path}: \`chain\` must be one of ${SUPPORTED_CHAINS.join(", ")}`);
  }

  const out: CliConfig = {
    gatewayUrl,
    chain,
  };

  if (obj.apiKey !== undefined) {
    if (typeof obj.apiKey !== "string") {
      throw new Error(`config file ${path}: \`apiKey\` must be a string`);
    }
    out.apiKey = obj.apiKey;
  }
  if (obj.agentRegistry !== undefined) {
    if (typeof obj.agentRegistry !== "string" || !isHexAddress(obj.agentRegistry)) {
      throw new Error(
        `config file ${path}: \`agentRegistry\` must be a 0x-prefixed 20-byte address`
      );
    }
    out.agentRegistry = obj.agentRegistry as `0x${string}`;
  }
  if (obj.jobFactory !== undefined) {
    if (typeof obj.jobFactory !== "string" || !isHexAddress(obj.jobFactory)) {
      throw new Error(`config file ${path}: \`jobFactory\` must be a 0x-prefixed 20-byte address`);
    }
    out.jobFactory = obj.jobFactory as `0x${string}`;
  }
  if (obj.rpcUrl !== undefined) {
    if (typeof obj.rpcUrl !== "string" || obj.rpcUrl.length === 0) {
      throw new Error(`config file ${path}: \`rpcUrl\` must be a non-empty string`);
    }
    out.rpcUrl = obj.rpcUrl;
  }
  if (obj.privateKey !== undefined) {
    if (typeof obj.privateKey !== "string" || !/^0x[0-9a-fA-F]{64}$/.test(obj.privateKey)) {
      throw new Error(
        `config file ${path}: \`privateKey\` must be a 0x-prefixed 32-byte hex string`
      );
    }
    out.privateKey = obj.privateKey;
  }
  return out;
}

/** Merge a partial update over an existing config (or build one from scratch
 * when `existing` is undefined). Centralised so `configure` and any future
 * sub-command that mutates state share the same defaults. */
export function mergeConfig(existing: CliConfig | undefined, patch: Partial<CliConfig>): CliConfig {
  const base: CliConfig = existing ?? {
    gatewayUrl: "",
    chain: "base-sepolia",
  };
  const merged: CliConfig = { ...base, ...patch };
  if (merged.gatewayUrl.length === 0) {
    throw new Error("mergeConfig: `gatewayUrl` is required");
  }
  return merged;
}

/** Strip the secret-bearing fields from a config so the result is safe to
 * log or pass to `--json`. */
export function redactConfig(cfg: CliConfig): CliConfig {
  const out: CliConfig = { gatewayUrl: cfg.gatewayUrl, chain: cfg.chain };
  if (cfg.agentRegistry !== undefined) out.agentRegistry = cfg.agentRegistry;
  if (cfg.jobFactory !== undefined) out.jobFactory = cfg.jobFactory;
  if (cfg.rpcUrl !== undefined) out.rpcUrl = cfg.rpcUrl;
  if (cfg.apiKey !== undefined) out.apiKey = "***";
  if (cfg.privateKey !== undefined) out.privateKey = "***";
  return out;
}

function isHexAddress(s: string): boolean {
  return /^0x[0-9a-fA-F]{40}$/.test(s);
}

function isFsErrno(value: unknown): value is NodeJS.ErrnoException {
  return value instanceof Error && typeof (value as { code?: unknown }).code === "string";
}

/** Test helper — does the config file exist? Used by integration tests so we
 * don't reach into `node:fs` from the test surface. */
export async function configExists(path: string): Promise<boolean> {
  try {
    await stat(path);
    return true;
  } catch (err) {
    if (isFsErrno(err) && err.code === "ENOENT") return false;
    throw err;
  }
}
