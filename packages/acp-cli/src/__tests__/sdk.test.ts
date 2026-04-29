import { describe, expect, it } from "vitest";

import {
  __testing__,
  buildAcpAgent,
  buildGatewayClient,
  resolveAlchemyChain,
  resolveContracts,
} from "../sdk.js";
import type { CliConfig } from "../types.js";

const ADDR_A = `0x${"a".repeat(40)}` as const;
const ADDR_B = `0x${"b".repeat(40)}` as const;

describe("resolveContracts", () => {
  it("falls back to chain defaults when not pinned", () => {
    const cfg: CliConfig = {
      gatewayUrl: "https://x",
      chain: "base-sepolia",
    };
    const c = resolveContracts(cfg);
    expect(c.agentRegistry).toBe(
      __testing__.CHAIN_DEFAULTS["base-sepolia"].contracts.agentRegistry
    );
  });

  it("respects per-config overrides", () => {
    const cfg: CliConfig = {
      gatewayUrl: "https://x",
      chain: "base-sepolia",
      agentRegistry: ADDR_A,
      jobFactory: ADDR_B,
    };
    const c = resolveContracts(cfg);
    expect(c.agentRegistry).toBe(ADDR_A);
    expect(c.jobFactory).toBe(ADDR_B);
  });
});

describe("resolveAlchemyChain", () => {
  it("returns the SDK slug for base-sepolia", () => {
    expect(resolveAlchemyChain({ gatewayUrl: "https://x", chain: "base-sepolia" })).toBe(
      "base-sepolia"
    );
  });

  it("returns 'base' for base-mainnet", () => {
    expect(resolveAlchemyChain({ gatewayUrl: "https://x", chain: "base-mainnet" })).toBe("base");
  });

  it("throws for unsupported chains", () => {
    expect(() =>
      resolveAlchemyChain({
        gatewayUrl: "https://x",
        chain: "ethereum-sepolia",
      })
    ).toThrow(/not yet supported/);
  });
});

describe("buildGatewayClient", () => {
  it("constructs without a bearer when apiKey is absent", () => {
    const gw = buildGatewayClient({
      gatewayUrl: "https://api.titular.dev",
      chain: "base-sepolia",
    });
    expect(gw.getBearerToken()).toBeUndefined();
  });

  it("threads apiKey as bearer token", () => {
    const gw = buildGatewayClient({
      gatewayUrl: "https://api.titular.dev",
      chain: "base-sepolia",
      apiKey: "k",
    });
    expect(gw.getBearerToken()).toBe("k");
  });
});

describe("buildAcpAgent", () => {
  it("uses overridden contracts and a stub provider factory", () => {
    const stubProvider = {
      getAddress: () => Promise.resolve(`0x${"f".repeat(40)}` as const),
      readContract: () => Promise.resolve(undefined),
      writeContract: () => Promise.resolve(`0x${"f".repeat(64)}` as const),
      waitForReceipt: () =>
        Promise.resolve({
          status: 1 as const,
          blockNumber: 1n,
          transactionHash: `0x${"f".repeat(64)}` as const,
          logs: [],
        }),
    };
    const agent = buildAcpAgent(
      {
        gatewayUrl: "https://api.titular.dev",
        chain: "base-sepolia",
        agentRegistry: ADDR_A,
        jobFactory: ADDR_B,
      },
      () => stubProvider as never
    );
    expect(agent.contracts.agentRegistry).toBe(ADDR_A);
    expect(agent.contracts.jobFactory).toBe(ADDR_B);
  });
});
