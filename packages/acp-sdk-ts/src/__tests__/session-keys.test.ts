import { afterEach, describe, expect, it, vi } from "vitest";
import { AcpError } from "../errors.js";
import {
  DEFAULT_SESSION_KEY_SERVICE,
  InMemorySessionKeyStore,
  type KeyringEntry,
  type KeyringModule,
  KeyringSessionKeyStore,
  createSessionKeyStore,
} from "../session-keys.js";

// Build an in-memory test double that mimics the `@napi-rs/keyring` Entry
// surface. Each (service, account) pair maps to one virtual keychain row.
function makeKeyringMock(): {
  module: KeyringModule;
  rows: Map<string, string>;
} {
  const rows = new Map<string, string>();
  const key = (service: string, account: string) => `${service}::${account}`;

  class FakeEntry implements KeyringEntry {
    constructor(
      private readonly service: string,
      private readonly account: string
    ) {}
    setPassword(secret: string): void {
      rows.set(key(this.service, this.account), secret);
    }
    getPassword(): string | null {
      return rows.get(key(this.service, this.account)) ?? null;
    }
    deletePassword(): boolean {
      return rows.delete(key(this.service, this.account));
    }
  }

  return { module: { Entry: FakeEntry }, rows };
}

describe("InMemorySessionKeyStore", () => {
  it("uses the default service name when none is supplied", () => {
    const store = new InMemorySessionKeyStore();
    expect(store.service).toBe(DEFAULT_SESSION_KEY_SERVICE);
    expect(store.backend).toBe("in-memory");
  });

  it("rejects an empty service name at construction", () => {
    expect(() => new InMemorySessionKeyStore({ service: "" })).toThrow(AcpError);
  });

  it("round-trips a secret", async () => {
    const store = new InMemorySessionKeyStore({
      service: "titular-acp:base-sepolia",
    });
    expect(store.service).toBe("titular-acp:base-sepolia");

    await store.set("0xabc", "deadbeef");
    expect(await store.get("0xabc")).toBe("deadbeef");
  });

  it("returns null for a missing account", async () => {
    const store = new InMemorySessionKeyStore();
    expect(await store.get("0xnobody")).toBeNull();
  });

  it("overwrites on repeated set", async () => {
    const store = new InMemorySessionKeyStore();
    await store.set("a", "first");
    await store.set("a", "second");
    expect(await store.get("a")).toBe("second");
  });

  it("delete is idempotent — succeeds on missing accounts", async () => {
    const store = new InMemorySessionKeyStore();
    await expect(store.delete("never-existed")).resolves.toBeUndefined();
    await store.set("a", "x");
    await store.delete("a");
    expect(await store.get("a")).toBeNull();
  });

  it("clear() wipes all entries (test helper)", async () => {
    const store = new InMemorySessionKeyStore();
    await store.set("a", "1");
    await store.set("b", "2");
    store.clear();
    expect(await store.get("a")).toBeNull();
    expect(await store.get("b")).toBeNull();
  });

  it("rejects empty account / secret on every method", async () => {
    const store = new InMemorySessionKeyStore();
    await expect(store.set("", "x")).rejects.toThrow(AcpError);
    await expect(store.set("a", "")).rejects.toThrow(AcpError);
    await expect(store.get("")).rejects.toThrow(AcpError);
    await expect(store.delete("")).rejects.toThrow(AcpError);
  });

  it("rejects non-string account / secret with invalid_param", async () => {
    const store = new InMemorySessionKeyStore();
    // @ts-expect-error: deliberately passing non-string to verify runtime guard
    await expect(store.set(123, "x")).rejects.toMatchObject({
      code: "invalid_param",
    });
    // @ts-expect-error: deliberately passing non-string to verify runtime guard
    await expect(store.set("a", null)).rejects.toMatchObject({
      code: "invalid_param",
    });
  });
});

describe("KeyringSessionKeyStore", () => {
  it("delegates set/get/delete to the keyring Entry under the configured service", async () => {
    const { module, rows } = makeKeyringMock();
    const store = new KeyringSessionKeyStore(module, {
      service: "titular-acp:base-sepolia",
    });

    expect(store.backend).toBe("os-keychain");
    expect(store.service).toBe("titular-acp:base-sepolia");

    await store.set("0xabc", "secret-jwt");
    expect(rows.get("titular-acp:base-sepolia::0xabc")).toBe("secret-jwt");

    expect(await store.get("0xabc")).toBe("secret-jwt");

    await store.delete("0xabc");
    expect(rows.has("titular-acp:base-sepolia::0xabc")).toBe(false);
    expect(await store.get("0xabc")).toBeNull();
  });

  it("namespaces per service — same account across services does not collide", async () => {
    const { module, rows } = makeKeyringMock();
    const mainnet = new KeyringSessionKeyStore(module, {
      service: "titular-acp:base",
    });
    const testnet = new KeyringSessionKeyStore(module, {
      service: "titular-acp:base-sepolia",
    });

    await mainnet.set("0xabc", "mainnet-key");
    await testnet.set("0xabc", "testnet-key");

    expect(await mainnet.get("0xabc")).toBe("mainnet-key");
    expect(await testnet.get("0xabc")).toBe("testnet-key");
    expect(rows.size).toBe(2);
  });

  it("wraps native errors as AcpError(keyring_unavailable) without leaking the message", async () => {
    const failingModule: KeyringModule = {
      Entry: class {
        setPassword(): void {
          throw new Error("user prompted: passphrase: hunter2");
        }
        getPassword(): string | null {
          throw new Error("backend disconnected: socket detail");
        }
        deletePassword(): boolean {
          throw new Error("os-specific error 0xdeadbeef");
        }
      },
    };
    const store = new KeyringSessionKeyStore(failingModule);

    const setErr = await store.set("a", "x").catch((e) => e);
    expect(setErr).toBeInstanceOf(AcpError);
    expect(setErr.code).toBe("keyring_unavailable");
    expect(setErr.message).not.toContain("hunter2");
    expect(setErr.message).not.toContain("passphrase");

    const getErr = await store.get("a").catch((e) => e);
    expect(getErr).toBeInstanceOf(AcpError);
    expect(getErr.message).not.toContain("socket");

    const delErr = await store.delete("a").catch((e) => e);
    expect(delErr).toBeInstanceOf(AcpError);
    expect(delErr.message).not.toContain("0xdeadbeef");
  });

  it("returns null when the keyring entry is empty", async () => {
    const { module } = makeKeyringMock();
    const store = new KeyringSessionKeyStore(module);
    expect(await store.get("0xnope")).toBeNull();
  });

  it("rejects empty inputs before touching the keyring", async () => {
    const probe = vi.fn();
    const sniffingModule: KeyringModule = {
      Entry: class {
        constructor() {
          probe();
        }
        setPassword(): void {}
        getPassword(): string | null {
          return null;
        }
        deletePassword(): boolean {
          return true;
        }
      },
    };
    const store = new KeyringSessionKeyStore(sniffingModule);
    await expect(store.set("", "x")).rejects.toBeInstanceOf(AcpError);
    expect(probe).not.toHaveBeenCalled();
  });
});

describe("createSessionKeyStore", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("returns the keychain-backed store when the keyring module loads", async () => {
    const { module } = makeKeyringMock();
    const store = await createSessionKeyStore({
      service: "titular-acp:base-sepolia",
      loadKeyring: async () => module,
    });
    expect(store.backend).toBe("os-keychain");
    expect(store.service).toBe("titular-acp:base-sepolia");
  });

  it("unwraps a default-export wrapped module", async () => {
    const { module } = makeKeyringMock();
    const store = await createSessionKeyStore({
      // simulate a CJS interop wrapper
      loadKeyring: async () => ({ default: module }) as unknown as KeyringModule,
    });
    expect(store.backend).toBe("os-keychain");
  });

  it("falls back to in-memory + warns when loader rejects", async () => {
    const warn = vi.fn();
    const store = await createSessionKeyStore({
      loadKeyring: async () => {
        throw new Error("native binary missing");
      },
      warn,
    });
    expect(store.backend).toBe("in-memory");
    expect(warn).toHaveBeenCalledTimes(1);
    const message = warn.mock.calls[0]?.[0] as string;
    expect(message).toContain("@napi-rs/keyring");
    expect(message.toLowerCase()).toContain("in-memory");
  });

  it("falls back when loader resolves to a non-keyring module", async () => {
    const warn = vi.fn();
    const store = await createSessionKeyStore({
      loadKeyring: async () => ({ Wrong: 1 }) as unknown as KeyringModule,
      warn,
    });
    expect(store.backend).toBe("in-memory");
    expect(warn).toHaveBeenCalled();
  });

  it("backend: 'in-memory' skips the loader entirely", async () => {
    const loadKeyring = vi.fn();
    const warn = vi.fn();
    const store = await createSessionKeyStore({
      backend: "in-memory",
      loadKeyring,
      warn,
    });
    expect(store.backend).toBe("in-memory");
    expect(loadKeyring).not.toHaveBeenCalled();
    expect(warn).not.toHaveBeenCalled();
  });

  it("backend: 'os-keychain' surfaces loader failures as AcpError(keyring_unavailable)", async () => {
    await expect(
      createSessionKeyStore({
        backend: "os-keychain",
        loadKeyring: async () => {
          throw new Error("missing");
        },
      })
    ).rejects.toMatchObject({ code: "keyring_unavailable" });
  });

  it("uses console.warn by default when falling back", async () => {
    const spy = vi.spyOn(console, "warn").mockImplementation(() => {});
    const store = await createSessionKeyStore({
      loadKeyring: async () => {
        throw new Error("missing");
      },
    });
    expect(store.backend).toBe("in-memory");
    expect(spy).toHaveBeenCalledTimes(1);
  });

  it("validates service even when forcing in-memory", async () => {
    await expect(createSessionKeyStore({ backend: "in-memory", service: "" })).rejects.toThrow(
      AcpError
    );
  });
});
