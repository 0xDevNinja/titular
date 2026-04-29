// SDK -> anvil end-to-end for AcpAgent.createJob.
//
// Mirrors `agent-register.test.ts` but drives the JobFactory path: installs
// the `MockJobFactory` runtime at a fixed address, calls `createJob(...)`,
// and asserts the JobCreated event is decoded into `{ jobId, cloneAddress,
// txHash }` with the values the SDK is expected to surface.
//
// The mock derives its clone address as `0xC1000000 + jobId` (see
// `MockJobFactory.sol`), which the test reproduces verbatim so the
// assertion is independent of the SDK's decoder logic.

import { AcpAgent, JobType } from "@titular/acp-sdk";
import { afterAll, beforeAll, describe, expect, it } from "vitest";
import { newAnvilProvider } from "../lib/anvil-provider.js";
import { type AnvilHandle, anvilAvailable, anvilSetCode, startAnvil } from "../lib/anvil.js";
import { MOCK_AGENT_REGISTRY_RUNTIME, MOCK_JOB_FACTORY_RUNTIME } from "../lib/mock-bytecode.js";
import { E2E_ENABLED, SKIP_REASON } from "./_skip.js";

const REGISTRY_ADDR = "0x000000000000000000000000000000000000ac01" as const;
const FACTORY_ADDR = "0x000000000000000000000000000000000000ac02" as const;
const TEST_TOKEN = "0x000000000000000000000000000000000000c0de" as const;
const ARBITER = "0x000000000000000000000000000000000000a1b1" as const;

/** MockJobFactory.createJob assigns clone = 0xC1000000 + jobId. */
function expectedClone(jobId: bigint): `0x${string}` {
  const numeric = 0xc1000000n + jobId;
  return `0x${numeric.toString(16).padStart(40, "0")}` as `0x${string}`;
}

describe.skipIf(!E2E_ENABLED)("AcpAgent.createJob against anvil", () => {
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
    await anvilSetCode(anvil.rpcUrl, FACTORY_ADDR, MOCK_JOB_FACTORY_RUNTIME);

    agent = new AcpAgent({
      provider: newAnvilProvider(anvil.rpcUrl),
      gatewayUrl: "http://127.0.0.1:1",
      contracts: { agentRegistry: REGISTRY_ADDR, jobFactory: FACTORY_ADDR },
    });
  });

  afterAll(async () => {
    await anvil?.close();
  });

  it("submits a Direct job and decodes JobCreated", async () => {
    if (agent === undefined) return; // suite was skipped
    const deadline = BigInt(Math.floor(Date.now() / 1000) + 3600);
    const result = await agent.createJob({
      token: TEST_TOKEN,
      budget: 1_000_000n,
      deadline,
      jobType: JobType.Direct,
    });
    expect(result.jobId).toBe(1n);
    expect(result.cloneAddress).toBe(expectedClone(1n));
    expect(result.txHash).toMatch(/^0x[0-9a-f]{64}$/);
  });

  it("increments jobId across consecutive createJob calls", async () => {
    if (agent === undefined) return; // suite was skipped
    const deadline = BigInt(Math.floor(Date.now() / 1000) + 3600);
    const first = await agent.createJob({
      token: TEST_TOKEN,
      budget: 100n,
      deadline,
    });
    const second = await agent.createJob({
      token: TEST_TOKEN,
      budget: 200n,
      deadline,
    });
    expect(second.jobId).toBe(first.jobId + 1n);
    expect(second.cloneAddress).toBe(expectedClone(second.jobId));
  });

  it("creates an Evaluated job when an evaluator address is supplied", async () => {
    if (agent === undefined) return; // suite was skipped
    const deadline = BigInt(Math.floor(Date.now() / 1000) + 3600);
    const evaluator = "0x000000000000000000000000000000000000eee0" as const;
    const result = await agent.createJob({
      token: TEST_TOKEN,
      budget: 42n,
      deadline,
      jobType: JobType.Evaluated,
      evaluator,
      arbiter: ARBITER,
    });
    expect(result.jobId).toBeGreaterThan(0n);
    expect(result.cloneAddress).toBe(expectedClone(result.jobId));
  });
});
