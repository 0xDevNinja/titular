import type { PublicClient, WalletClient } from "viem";
import { describe, expect, it, vi } from "vitest";
import {
  AlchemyEvmProviderAdapter,
  NoSignerError,
  defaultAlchemyRpcUrl,
} from "../providers/alchemy.js";
import type { EvmAddress, TxHash } from "../types.js";

const ADDR: EvmAddress = "0xabcdef0000000000000000000000000000000001";
const TX: TxHash = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb";

/**
 * Build a minimal stub of viem's PublicClient. Only the methods used by
 * AlchemyEvmProviderAdapter are stubbed; the cast is concentrated here so the
 * adapter constructor accepts the stub without losing access to the mock fns
 * in the calling test.
 */
type StubPublicClient = {
  getChainId: ReturnType<typeof vi.fn>;
  readContract: ReturnType<typeof vi.fn>;
  waitForTransactionReceipt: ReturnType<typeof vi.fn>;
};

type StubWalletClient = {
  account: { address: EvmAddress } | undefined;
  getAddresses: ReturnType<typeof vi.fn>;
  writeContract: ReturnType<typeof vi.fn>;
  signMessage: ReturnType<typeof vi.fn>;
};

function makePublicClient(overrides: Partial<StubPublicClient> = {}): StubPublicClient {
  return {
    getChainId: vi.fn().mockResolvedValue(84532),
    readContract: vi.fn().mockResolvedValue("read-result"),
    waitForTransactionReceipt: vi.fn().mockResolvedValue({
      status: "success",
      transactionHash: TX,
      blockNumber: 99n,
      logs: [
        {
          address: ADDR,
          topics: ["0xdead"],
          data: "0xbeef",
          logIndex: 3,
        },
      ],
    }),
    ...overrides,
  };
}

function makeWalletClient(overrides: Partial<StubWalletClient> = {}): StubWalletClient {
  return {
    account: { address: ADDR },
    getAddresses: vi.fn().mockResolvedValue([ADDR]),
    writeContract: vi.fn().mockResolvedValue(TX),
    signMessage: vi.fn().mockResolvedValue("0xsig"),
    ...overrides,
  };
}

/** Cast a stub to viem's branded client types when handing it to the adapter. */
function asPublic(s: StubPublicClient): PublicClient {
  return s as unknown as PublicClient;
}
function asWallet(s: StubWalletClient): WalletClient {
  return s as unknown as WalletClient;
}

describe("defaultAlchemyRpcUrl", () => {
  it("builds the base mainnet URL", () => {
    expect(defaultAlchemyRpcUrl("base", "key")).toBe("https://base-mainnet.g.alchemy.com/v2/key");
  });

  it("builds the base sepolia URL", () => {
    expect(defaultAlchemyRpcUrl("base-sepolia", "key")).toBe(
      "https://base-sepolia.g.alchemy.com/v2/key"
    );
  });

  it("rejects an empty api key", () => {
    expect(() => defaultAlchemyRpcUrl("base", "")).toThrow(/apiKey/);
  });
});

describe("AlchemyEvmProviderAdapter (pre-built clients mode)", () => {
  it("getChainId / getAddress / readContract delegate to the public/wallet clients", async () => {
    const publicClient = makePublicClient();
    const walletClient = makeWalletClient();
    const adapter = new AlchemyEvmProviderAdapter({
      publicClient: asPublic(publicClient),
      walletClient: asWallet(walletClient),
    });

    expect(await adapter.getChainId()).toBe(84532);
    expect(await adapter.getAddress()).toBe(ADDR);
    const r = await adapter.readContract({
      address: ADDR,
      abi: [],
      functionName: "anything",
    });
    expect(r).toBe("read-result");
  });

  it("getAddress falls back to walletClient.getAddresses when account is undefined", async () => {
    const publicClient = makePublicClient();
    const walletClient = makeWalletClient({ account: undefined });
    const adapter = new AlchemyEvmProviderAdapter({
      publicClient: asPublic(publicClient),
      walletClient: asWallet(walletClient),
    });
    expect(await adapter.getAddress()).toBe(ADDR);
  });

  it("getAddress throws NoSignerError when wallet has no addresses", async () => {
    const publicClient = makePublicClient();
    const walletClient = makeWalletClient({
      account: undefined,
      getAddresses: vi.fn().mockResolvedValue([]),
    });
    const adapter = new AlchemyEvmProviderAdapter({
      publicClient: asPublic(publicClient),
      walletClient: asWallet(walletClient),
    });
    await expect(adapter.getAddress()).rejects.toBeInstanceOf(NoSignerError);
  });

  it("writeContract delegates to walletClient.writeContract and returns the hash", async () => {
    const publicClient = makePublicClient();
    const walletClient = makeWalletClient();
    const adapter = new AlchemyEvmProviderAdapter({
      publicClient: asPublic(publicClient),
      walletClient: asWallet(walletClient),
    });
    const hash = await adapter.writeContract({
      address: ADDR,
      abi: [],
      functionName: "register",
      args: [],
      value: 1n,
    });
    expect(hash).toBe(TX);
    expect(walletClient.writeContract).toHaveBeenCalledOnce();
  });

  it("waitForReceipt normalises the receipt shape", async () => {
    const publicClient = makePublicClient();
    const adapter = new AlchemyEvmProviderAdapter({
      publicClient: asPublic(publicClient),
    });
    const receipt = await adapter.waitForReceipt(TX);
    expect(receipt.status).toBe("success");
    expect(receipt.transactionHash).toBe(TX);
    expect(receipt.blockNumber).toBe(99n);
    expect(receipt.logs[0]).toMatchObject({
      address: ADDR,
      topics: ["0xdead"],
      data: "0xbeef",
      logIndex: 3,
    });
  });

  it("waitForReceipt translates non-success receipts into status: reverted", async () => {
    const publicClient = makePublicClient({
      waitForTransactionReceipt: vi.fn().mockResolvedValue({
        status: "reverted",
        transactionHash: TX,
        blockNumber: 1n,
        logs: [],
      }),
    });
    const adapter = new AlchemyEvmProviderAdapter({
      publicClient: asPublic(publicClient),
    });
    const r = await adapter.waitForReceipt(TX);
    expect(r.status).toBe("reverted");
  });

  it("waitForReceipt fills logIndex=0 when missing", async () => {
    const publicClient = makePublicClient({
      waitForTransactionReceipt: vi.fn().mockResolvedValue({
        status: "success",
        transactionHash: TX,
        blockNumber: 1n,
        logs: [{ address: ADDR, topics: [], data: "0x" }],
      }),
    });
    const adapter = new AlchemyEvmProviderAdapter({
      publicClient: asPublic(publicClient),
    });
    const r = await adapter.waitForReceipt(TX);
    expect(r.logs[0]?.logIndex).toBe(0);
  });

  it("signMessage forwards to walletClient.signMessage", async () => {
    const publicClient = makePublicClient();
    const walletClient = makeWalletClient();
    const adapter = new AlchemyEvmProviderAdapter({
      publicClient: asPublic(publicClient),
      walletClient: asWallet(walletClient),
    });
    const sig = await adapter.signMessage("hello");
    expect(sig).toBe("0xsig");
    expect(walletClient.signMessage).toHaveBeenCalledOnce();
  });
});

describe("AlchemyEvmProviderAdapter (read-only)", () => {
  it("write methods throw NoSignerError without a wallet", async () => {
    const publicClient = makePublicClient();
    const adapter = new AlchemyEvmProviderAdapter({
      publicClient: asPublic(publicClient),
    });
    await expect(
      adapter.writeContract({
        address: ADDR,
        abi: [],
        functionName: "x",
        args: [],
      })
    ).rejects.toBeInstanceOf(NoSignerError);
    await expect(adapter.signMessage("x")).rejects.toBeInstanceOf(NoSignerError);
    await expect(adapter.getAddress()).rejects.toBeInstanceOf(NoSignerError);
  });

  it("NoSignerError carries a stable code for branchable handling", async () => {
    const publicClient = makePublicClient();
    const adapter = new AlchemyEvmProviderAdapter({
      publicClient: asPublic(publicClient),
    });
    try {
      await adapter.getAddress();
    } catch (e) {
      expect(e).toBeInstanceOf(NoSignerError);
      expect((e as NoSignerError).code).toBe("no_signer");
      return;
    }
    throw new Error("expected NoSignerError");
  });

  it("signMessage rejects when walletClient has no account", async () => {
    const publicClient = makePublicClient();
    const walletClient = makeWalletClient({ account: undefined });
    const adapter = new AlchemyEvmProviderAdapter({
      publicClient: asPublic(publicClient),
      walletClient: asWallet(walletClient),
    });
    await expect(adapter.signMessage("hi")).rejects.toBeInstanceOf(NoSignerError);
  });
});

describe("AlchemyEvmProviderAdapter (apiKey constructor)", () => {
  it("builds clients without throwing when only apiKey + chain are passed", async () => {
    const adapter = new AlchemyEvmProviderAdapter({
      apiKey: "test-key",
      chain: "base-sepolia",
    });
    // No signer, so write methods throw — but the adapter itself must instantiate.
    await expect(adapter.getAddress()).rejects.toBeInstanceOf(NoSignerError);
  });

  it("builds clients with a private key", () => {
    const adapter = new AlchemyEvmProviderAdapter({
      apiKey: "test-key",
      chain: "base-sepolia",
      privateKey: "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
    });
    expect(adapter).toBeInstanceOf(AlchemyEvmProviderAdapter);
  });

  it("respects rpcUrlOverride (smoke: instance constructs without error)", () => {
    const adapter = new AlchemyEvmProviderAdapter({
      apiKey: "test-key",
      chain: "base",
      rpcUrlOverride: "https://custom.example/rpc",
    });
    expect(adapter).toBeInstanceOf(AlchemyEvmProviderAdapter);
  });
});
