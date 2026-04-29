// SessionKeyStore — secure storage for SDK + CLI session secrets.
//
// Wraps `@napi-rs/keyring` so callers can persist EVM private keys, SIWE JWTs,
// and other long-lived session material in the host OS credential store
// (macOS Keychain, Windows Credential Manager, libsecret on Linux) instead of
// a plaintext config file.
//
// The keyring package is a native module — pulling it in unconditionally would
// break consumers that ship to environments without prebuilt binaries (Alpine,
// some serverless runtimes, browser bundlers). We declare it as an
// `optionalDependencies` entry and load it lazily; if the module fails to
// import, the factory returns a clearly-marked in-memory fallback and emits a
// `console.warn` so operators know secrets are NOT being persisted.
//
// Service names are namespaced per-chain (`titular-acp:base-sepolia`) so a
// developer juggling mainnet / testnet keys in the same OS user account does
// not silently overwrite one with the other. Account ids are caller-supplied
// (typically the EVM address or `siwe:<address>`).
//
// Errors NEVER include the secret value or the underlying provider message
// (which on macOS occasionally leaks parts of the prompt) — only a stable
// machine-readable code, so logs at INFO level are safe to ship to telemetry.

import { AcpError } from "./errors.js";

/** Default service name. Callers SHOULD pass an explicit `service` per chain. */
export const DEFAULT_SESSION_KEY_SERVICE = "titular-acp";

/**
 * Common interface implemented by both the OS-keychain-backed store and the
 * in-memory fallback. Consumers should program against this type, not the
 * concrete classes — the factory may return either depending on the host.
 */
export interface SessionKeyStore {
  /** Human-readable description of where secrets live ("os-keychain" / "in-memory"). */
  readonly backend: SessionKeyBackend;
  /** Service namespace this store was constructed with. */
  readonly service: string;
  /** Persist `secret` under `account`. Overwrites any existing value. */
  set(account: string, secret: string): Promise<void>;
  /** Retrieve `secret` for `account`, or `null` if absent. */
  get(account: string): Promise<string | null>;
  /** Remove the entry. Idempotent — succeeds if the entry is already absent. */
  delete(account: string): Promise<void>;
}

/** Discriminator returned by {@link SessionKeyStore.backend}. */
export type SessionKeyBackend = "os-keychain" | "in-memory";

/**
 * Minimal subset of `@napi-rs/keyring`'s `Entry` API the store relies on.
 * Declared explicitly so the SDK does not take a hard type-time dependency
 * on the optional native module.
 */
export interface KeyringEntry {
  setPassword(secret: string): void;
  getPassword(): string | null;
  deletePassword(): boolean;
}

/** Module shape imported lazily from `@napi-rs/keyring`. */
export interface KeyringModule {
  Entry: new (service: string, account: string) => KeyringEntry;
}

/** Options accepted by every store. */
export interface SessionKeyStoreOptions {
  /**
   * Service namespace. Callers SHOULD include the chain id / network name
   * so testnet and mainnet credentials don't collide
   * (e.g. `titular-acp:base-sepolia`). Defaults to {@link DEFAULT_SESSION_KEY_SERVICE}.
   */
  service?: string;
}

/** Options for {@link createSessionKeyStore}. */
export interface CreateSessionKeyStoreOptions extends SessionKeyStoreOptions {
  /**
   * Force a specific backend. When omitted, tries `os-keychain` first and
   * falls back to `in-memory` if the native module cannot be loaded.
   */
  backend?: SessionKeyBackend;
  /**
   * Override the keyring loader — used by tests to inject a mock without
   * relying on `vi.mock`. Must resolve to a module exposing `Entry` (or a
   * `default` namespace wrapping one — CJS-interop is unwrapped automatically).
   */
  loadKeyring?: () => Promise<unknown>;
  /**
   * Override `console.warn`. Used by tests to assert the fallback notice
   * fires, and by hosts that want to route warnings to their logger.
   */
  warn?: (message: string) => void;
}

// --- OS-keychain-backed store ------------------------------------------------

/**
 * Store backed by `@napi-rs/keyring`. Construct directly when you have already
 * resolved the keyring module (e.g. in a CLI that requires persistence and is
 * willing to error out when the native binary is missing). Most callers should
 * use {@link createSessionKeyStore}, which falls back transparently.
 */
export class KeyringSessionKeyStore implements SessionKeyStore {
  readonly backend: SessionKeyBackend = "os-keychain";
  readonly service: string;
  readonly #keyring: KeyringModule;

  constructor(keyring: KeyringModule, options: SessionKeyStoreOptions = {}) {
    this.#keyring = keyring;
    this.service = normalizeService(options.service);
  }

  async set(account: string, secret: string): Promise<void> {
    assertNonEmpty(account, "account");
    assertNonEmpty(secret, "secret");
    const entry = this.#entry(account);
    try {
      entry.setPassword(secret);
    } catch (cause) {
      // Deliberately swallow the underlying message — on macOS the OS surfaces
      // user-facing prompt text that we don't want in our logs.
      throw new AcpError("keyring_unavailable", "failed to write to OS keychain", { cause });
    }
  }

  async get(account: string): Promise<string | null> {
    assertNonEmpty(account, "account");
    const entry = this.#entry(account);
    try {
      const value = entry.getPassword();
      return value ?? null;
    } catch (cause) {
      throw new AcpError("keyring_unavailable", "failed to read from OS keychain", { cause });
    }
  }

  async delete(account: string): Promise<void> {
    assertNonEmpty(account, "account");
    const entry = this.#entry(account);
    try {
      entry.deletePassword();
    } catch (cause) {
      throw new AcpError("keyring_unavailable", "failed to delete from OS keychain", { cause });
    }
  }

  #entry(account: string): KeyringEntry {
    return new this.#keyring.Entry(this.service, account);
  }
}

// --- in-memory fallback ------------------------------------------------------

/**
 * Process-local store used in CI, unit tests, and any runtime where the native
 * keyring is unavailable. Secrets do NOT survive process exit. Callers using
 * this in production must accept that they are running without persistence.
 */
export class InMemorySessionKeyStore implements SessionKeyStore {
  readonly backend: SessionKeyBackend = "in-memory";
  readonly service: string;
  readonly #map = new Map<string, string>();

  constructor(options: SessionKeyStoreOptions = {}) {
    this.service = normalizeService(options.service);
  }

  async set(account: string, secret: string): Promise<void> {
    assertNonEmpty(account, "account");
    assertNonEmpty(secret, "secret");
    this.#map.set(account, secret);
  }

  async get(account: string): Promise<string | null> {
    assertNonEmpty(account, "account");
    return this.#map.get(account) ?? null;
  }

  async delete(account: string): Promise<void> {
    assertNonEmpty(account, "account");
    this.#map.delete(account);
  }

  /**
   * Test helper — wipe all entries. NOT part of the {@link SessionKeyStore}
   * contract because it is meaningless against a real keychain (which holds
   * entries from many services).
   */
  clear(): void {
    this.#map.clear();
  }
}

// --- factory -----------------------------------------------------------------

/**
 * Build a {@link SessionKeyStore}, preferring the OS keychain and falling back
 * to in-memory storage if `@napi-rs/keyring` cannot be loaded (missing native
 * binary, sandboxed runtime, etc.). When falling back, emits a `console.warn`
 * so operators are not silently running without persistence.
 *
 * Pass `backend: "in-memory"` to skip the load attempt entirely — useful for
 * tests and `--non-interactive` CI flows where the keychain would prompt.
 */
export async function createSessionKeyStore(
  options: CreateSessionKeyStoreOptions = {}
): Promise<SessionKeyStore> {
  const service = normalizeService(options.service);

  if (options.backend === "in-memory") {
    return new InMemorySessionKeyStore({ service });
  }

  const loader = options.loadKeyring ?? defaultLoadKeyring;
  try {
    const raw = await loader();
    const mod = normalizeKeyringModule(raw);
    return new KeyringSessionKeyStore(mod, { service });
  } catch (cause) {
    if (options.backend === "os-keychain") {
      // Caller asked explicitly for keychain — surface the failure rather
      // than silently downgrading.
      throw new AcpError("keyring_unavailable", "failed to load @napi-rs/keyring", { cause });
    }
    const warn = options.warn ?? defaultWarn;
    warn(
      "[@titular/acp-sdk] @napi-rs/keyring unavailable; falling back to in-memory session" +
        " key store. Secrets will NOT persist across process restarts."
    );
    return new InMemorySessionKeyStore({ service });
  }
}

async function defaultLoadKeyring(): Promise<unknown> {
  // Dynamic import so the optional dep is only resolved at runtime and does
  // not appear in the static module graph (bundlers can drop it cleanly).
  // The string is split to discourage bundlers from rewriting it.
  const specifier = "@napi-rs/keyring";
  return (await import(/* @vite-ignore */ specifier)) as unknown;
}

/**
 * Coerce a loader result into a {@link KeyringModule}. Tolerates the common
 * CJS-interop pattern where the module is wrapped in a `default` namespace.
 */
function normalizeKeyringModule(raw: unknown): KeyringModule {
  if (isKeyringModule(raw)) return raw;
  if (
    typeof raw === "object" &&
    raw !== null &&
    "default" in raw &&
    isKeyringModule((raw as { default: unknown }).default)
  ) {
    return (raw as { default: KeyringModule }).default;
  }
  throw new Error("@napi-rs/keyring module did not expose Entry");
}

function defaultWarn(message: string): void {
  // Routed through console.warn so the operator sees it on stderr; tests
  // override via `options.warn`.
  console.warn(message);
}

// --- helpers -----------------------------------------------------------------

function normalizeService(service: string | undefined): string {
  if (service === undefined) return DEFAULT_SESSION_KEY_SERVICE;
  if (typeof service !== "string" || service.length === 0) {
    throw new AcpError("invalid_param", "`service` must be a non-empty string");
  }
  return service;
}

function assertNonEmpty(value: unknown, key: string): void {
  if (typeof value !== "string" || value.length === 0) {
    throw new AcpError("invalid_param", `\`${key}\` must be a non-empty string`);
  }
}

function isKeyringModule(value: unknown): value is KeyringModule {
  return (
    typeof value === "object" &&
    value !== null &&
    "Entry" in value &&
    typeof (value as { Entry: unknown }).Entry === "function"
  );
}
