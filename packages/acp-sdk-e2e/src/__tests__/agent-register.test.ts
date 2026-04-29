// SDK -> anvil end-to-end for AcpAgent.register.
//
// Spins up anvil, installs the runtime bytecode of `MockAgentRegistry` at a
// fixed address via `anvil_setCode`, builds an `AcpAgent` whose provider
// adapter wraps viem against that anvil, calls `register(...)`, and asserts:
//
//   * the receipt is success
//   * the AgentRegistered event was decoded
//   * `agentId`, `controller`, `metadataUri`, `capabilities` round-trip
//   * monotonic id increment across two consecutive registrations
//
// The full SDK->gateway interaction is covered by `gateway-list.test.ts`.

import { AcpAgent } from "@titular/acp-sdk";
import { afterAll, beforeAll, describe, expect, it } from "vitest";
import { newAnvilProvider } from "../lib/anvil-provider.js";
import {
  ANVIL_DEV_ADDRESS,
  type AnvilHandle,
  anvilAvailable,
  anvilSetCode,
  startAnvil,
} from "../lib/anvil.js";
import { MOCK_AGENT_REGISTRY_RUNTIME, MOCK_JOB_FACTORY_RUNTIME } from "../lib/mock-bytecode.js";
import { E2E_ENABLED, SKIP_REASON } from "./_skip.js";

const REGISTRY_ADDR = "0x000000000000000000000000000000000000ac01" as const;
const FACTORY_ADDR = "0x000000000000000000000000000000000000ac02" as const;

describe.skipIf(!E2E_ENABLED)("AcpAgent.register against anvil", () => {
  let anvil: AnvilHandle | undefined;
  let agent: AcpAgent | undefined;

  beforeAll(async (ctx) => {
    if (!E2E_ENABLED) {
      ctx.skip(SKIP_REASON);
      return;
    }
    if (!(await anvilAvailable())) {
      ctx.skip("anvil not on PATH; install foundry to run this suite");
      return;
    }
    anvil = await startAnvil();
    await anvilSetCode(anvil.rpcUrl, REGISTRY_ADDR, MOCK_AGENT_REGISTRY_RUNTIME);
    // Install the JobFactory mock too so AcpAgent.constructor's non-zero
    // address invariant holds; the test doesn't call createJob here.
    await anvilSetCode(anvil.rpcUrl, FACTORY_ADDR, MOCK_JOB_FACTORY_RUNTIME);

    agent = new AcpAgent({
      provider: newAnvilProvider(anvil.rpcUrl),
      // No gateway needed for register(); pass a non-empty placeholder so
      // the constructor's `gatewayUrl` invariant holds.
      gatewayUrl: "http://127.0.0.1:1",
      contracts: { agentRegistry: REGISTRY_ADDR, jobFactory: FACTORY_ADDR },
    });
  });

  afterAll(async () => {
    await anvil?.close();
  });

  it("submits, decodes AgentRegistered, and returns the new agentId", async () => {
    if (agent === undefined) return; // suite was skipped
    const result = await agent.register({
      metadataUri: "ipfs://bafytestregister1",
      capabilities: 0b1011n,
    });
    expect(result.agentId).toBe(1n);
    expect(result.controller).toBe(ANVIL_DEV_ADDRESS);
    expect(result.metadataUri).toBe("ipfs://bafytestregister1");
    expect(result.capabilities).toBe(0b1011n);
    expect(result.txHash).toMatch(/^0x[0-9a-f]{64}$/);
  });

  it("increments agentId monotonically on subsequent registrations", async () => {
    if (agent === undefined) return; // suite was skipped
    const first = await agent.register({
      metadataUri: "ipfs://bafytestregister-mono-a",
      capabilities: 1n,
    });
    const second = await agent.register({
      metadataUri: "ipfs://bafytestregister-mono-b",
      capabilities: 2n,
    });
    expect(second.agentId).toBe(first.agentId + 1n);
    expect(second.controller).toBe(ANVIL_DEV_ADDRESS);
  });

  it("honours an explicit controller override", async () => {
    if (agent === undefined) return; // suite was skipped
    const otherController = "0x000000000000000000000000000000000000beef" as const;
    const result = await agent.register({
      controller: otherController,
      metadataUri: "ipfs://bafytestregister-controller",
      capabilities: 0n,
    });
    expect(result.controller).toBe(otherController);
  });
});
