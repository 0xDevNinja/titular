// `acp configure` — interactive (and non-interactive) gateway/chain/key
// setup. Persists to `~/.config/titular/config.json` (or wherever the user
// pointed XDG_CONFIG_HOME / TITULAR_CONFIG).

import {
  SUPPORTED_CHAINS,
  configExists,
  defaultConfigPath,
  isSupportedChain,
  loadConfig,
  mergeConfig,
  redactConfig,
  saveConfig,
} from "../config.js";
import { buildOutputContext, errorString, writeJson, writeKv, writeLine } from "../output.js";
import { type PromptIO, ask, select } from "../prompt.js";
import type { CliConfig, GlobalFlags, SupportedChain } from "../types.js";

/** Flags accepted by `acp configure`. */
export interface ConfigureFlags extends GlobalFlags {
  gatewayUrl?: string;
  chain?: string;
  apiKey?: string;
  agentRegistry?: string;
  jobFactory?: string;
  rpcUrl?: string;
  nonInteractive?: boolean;
  print?: boolean;
}

/** Optional dependencies — exposed so tests can inject. */
export interface ConfigureDeps {
  promptIO?: PromptIO;
  stdout?: NodeJS.WriteStream;
  stderr?: NodeJS.WriteStream;
  env?: NodeJS.ProcessEnv;
  isTTY?: boolean;
}

/** Run the `configure` command. Returns the resolved config so tests can
 * assert on it without round-tripping through disk. */
export async function runConfigure(
  flags: ConfigureFlags,
  deps: ConfigureDeps = {}
): Promise<CliConfig> {
  const env = deps.env ?? process.env;
  const path = flags.config ?? defaultConfigPath(env);
  const ctx = buildOutputContext({
    json: flags.json === true,
    ...(deps.stdout !== undefined ? { stdout: deps.stdout } : {}),
    ...(deps.stderr !== undefined ? { stderr: deps.stderr } : {}),
    env,
    ...(deps.isTTY !== undefined ? { isTTY: deps.isTTY } : {}),
  });

  // `--print` is a read-only escape hatch — useful for `acp configure --print
  // --json | jq` and for shell completions to discover the configured
  // gateway URL.
  if (flags.print === true) {
    const existing = await loadConfig(path);
    if (existing === undefined) {
      ctx.stderr.write(`${errorString(ctx, `no config found at ${path}`)}\n`);
      throw new ConfigureError("no_config", `no config at ${path}`);
    }
    if (ctx.json) {
      writeJson(ctx, redactConfig(existing));
    } else {
      printConfigPretty(ctx, existing, path);
    }
    return existing;
  }

  const existing = await loadConfig(path);
  const patch: Partial<CliConfig> = {};

  // --- gatewayUrl ---
  if (flags.gatewayUrl !== undefined) {
    patch.gatewayUrl = flags.gatewayUrl;
  } else if (flags.nonInteractive === true) {
    if (existing?.gatewayUrl === undefined) {
      throw new ConfigureError(
        "missing_gateway",
        "configure --non-interactive requires --gateway-url when no config exists"
      );
    }
  } else {
    patch.gatewayUrl = await ask(
      "gateway URL",
      existing?.gatewayUrl ?? "https://api.titular.dev",
      deps.promptIO
    );
  }

  // --- chain ---
  if (flags.chain !== undefined) {
    if (!isSupportedChain(flags.chain)) {
      throw new ConfigureError(
        "invalid_chain",
        `--chain must be one of ${SUPPORTED_CHAINS.join(", ")} (got ${flags.chain})`
      );
    }
    patch.chain = flags.chain;
  } else if (flags.nonInteractive !== true) {
    const chosen = await select(
      "chain",
      SUPPORTED_CHAINS as readonly string[],
      existing?.chain ?? "base-sepolia",
      deps.promptIO
    );
    if (!isSupportedChain(chosen)) {
      throw new ConfigureError("invalid_chain", `unknown chain ${chosen}`);
    }
    patch.chain = chosen as SupportedChain;
  }

  // --- apiKey (optional) ---
  if (flags.apiKey !== undefined) {
    patch.apiKey = flags.apiKey;
  } else if (flags.nonInteractive !== true && existing?.apiKey === undefined) {
    // We *only* ask the first time; subsequent runs show "***" and skip.
    const reply = await ask("api key (leave blank to skip)", "", deps.promptIO);
    if (reply.length > 0) patch.apiKey = reply;
  }

  // --- optional address overrides — flag-only, no prompt to keep the
  // interactive flow short. Operators wiring a private deploy pass them
  // explicitly. ---
  if (flags.agentRegistry !== undefined) {
    patch.agentRegistry = flags.agentRegistry as `0x${string}`;
  }
  if (flags.jobFactory !== undefined) {
    patch.jobFactory = flags.jobFactory as `0x${string}`;
  }
  if (flags.rpcUrl !== undefined) {
    patch.rpcUrl = flags.rpcUrl;
  }

  const merged = mergeConfig(existing, patch);
  await saveConfig(path, merged);

  const wasNew = !(await configExists(path)) ? false : existing === undefined;
  if (ctx.json) {
    writeJson(ctx, {
      saved: true,
      path,
      created: wasNew,
      config: redactConfig(merged),
    });
  } else {
    writeLine(ctx, `saved ${path}`);
    printConfigPretty(ctx, merged, path);
  }
  return merged;
}

function printConfigPretty(
  ctx: ReturnType<typeof buildOutputContext>,
  cfg: CliConfig,
  path: string
): void {
  const r = redactConfig(cfg);
  writeKv(ctx, "path", path);
  writeKv(ctx, "gatewayUrl", r.gatewayUrl);
  writeKv(ctx, "chain", r.chain);
  if (r.apiKey !== undefined) writeKv(ctx, "apiKey", r.apiKey);
  if (r.agentRegistry !== undefined) writeKv(ctx, "agentRegistry", r.agentRegistry);
  if (r.jobFactory !== undefined) writeKv(ctx, "jobFactory", r.jobFactory);
  if (r.rpcUrl !== undefined) writeKv(ctx, "rpcUrl", r.rpcUrl);
  if (r.privateKey !== undefined) writeKv(ctx, "privateKey", r.privateKey);
}

/** Typed error class so the test surface can assert on `code` rather than
 * string-matching messages. */
export class ConfigureError extends Error {
  override readonly name = "ConfigureError";
  readonly code: string;
  constructor(code: string, message: string) {
    super(message);
    this.code = code;
  }
}
