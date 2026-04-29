// End-to-end round-trip: AcpAgent.register on anvil -> mirror row inserted
// into Postgres -> AcpAgent.browse() (gateway) returns the new agent.
//
// Why mirror manually rather than spin up the full indexer pipeline:
//
//   The indexer (subscriber -> decoder -> publisher) is exercised in
//   `services/indexer-go/internal/integration/...` against the same anvil
//   harness this suite uses. Re-running it here would double the CI cost
//   for a property already covered. This test instead asserts the
//   *contract* between the indexer and the gateway: given a registered
//   agent on chain whose row was persisted under the indexer's schema,
//   the SDK's `browse()` call surfaces it with the same fields that came
//   off the chain. The mirror step plays the role of the indexer's
//   `decoder.OnAgentRegistered` handler.

import { AcpAgent } from "@titular/acp-sdk";
import { afterAll, beforeAll, describe, expect, it } from "vitest";
import { newAnvilProvider } from "../lib/anvil-provider.js";
import { type AnvilHandle, anvilAvailable, anvilSetCode, startAnvil } from "../lib/anvil.js";
import { type GatewayHandle, goAvailable, startGateway } from "../lib/gateway.js";
import { MOCK_AGENT_REGISTRY_RUNTIME, MOCK_JOB_FACTORY_RUNTIME } from "../lib/mock-bytecode.js";
import { type PostgresHandle, seedAgent, startPostgres } from "../lib/postgres.js";
import { E2E_ENABLED, SKIP_REASON } from "./_skip.js";

const REGISTRY_ADDR = "0x000000000000000000000000000000000000ac01" as const;
const FACTORY_ADDR = "0x000000000000000000000000000000000000ac02" as const;

interface Fixture {
  anvil: AnvilHandle;
  pg: PostgresHandle;
  gw: GatewayHandle;
  agent: AcpAgent;
}

describe.skipIf(!E2E_ENABLED)("AcpAgent round-trip: register -> indexer-mirror -> browse", () => {
  let fx: Fixture | undefined;

  beforeAll(async (ctx) => {
    if (!E2E_ENABLED) {
      ctx.skip(SKIP_REASON);
      return;
    }
    if (!(await anvilAvailable())) {
      ctx.skip("anvil not on PATH; install foundry to run this suite");
      return;
    }
    if (!(await goAvailable())) {
      ctx.skip("go not on PATH; install go 1.25+ to run this suite");
      return;
    }

    let pg: PostgresHandle | undefined;
    try {
      pg = await startPostgres();
    } catch (err) {
      ctx.skip(`docker daemon unreachable: ${err instanceof Error ? err.message : String(err)}`);
      return;
    }

    const anvil = await startAnvil();
    try {
      await anvilSetCode(anvil.rpcUrl, REGISTRY_ADDR, MOCK_AGENT_REGISTRY_RUNTIME);
      await anvilSetCode(anvil.rpcUrl, FACTORY_ADDR, MOCK_JOB_FACTORY_RUNTIME);

      const gw = await startGateway({ databaseUrl: pg.databaseUrl });

      const acpAgent = new AcpAgent({
        provider: newAnvilProvider(anvil.rpcUrl),
        gatewayUrl: gw.baseUrl,
        contracts: { agentRegistry: REGISTRY_ADDR, jobFactory: FACTORY_ADDR },
      });

      fx = { anvil, pg, gw, agent: acpAgent };
    } catch (err) {
      await anvil.close();
      await pg.close();
      throw err;
    }
  }, 240_000);

  afterAll(async () => {
    if (fx === undefined) return;
    await fx.gw.close();
    await fx.anvil.close();
    await fx.pg.close();
  });

  it("register -> mirror -> browse surfaces the new agent end-to-end", async () => {
    if (fx === undefined) return; // suite was skipped

    // 1. SDK -> anvil: register a fresh agent and capture the on-chain id.
    const registered = await fx.agent.register({
      metadataUri: "ipfs://bafyroundtrip",
      capabilities: 0b101n,
    });
    expect(registered.agentId).toBeGreaterThan(0n);

    // 2. Indexer-mirror step: persist the row the indexer would have
    //    written. Block number / log_index do not need to match the
    //    receipt's exact values for the SDK shape assertions to hold;
    //    the relationship under test is "row in DB <-> SDK Agent shape".
    const blockNumber = 1n;
    const insertedPk = await seedAgent(fx.pg, {
      agentId: registered.agentId,
      controller: registered.controller,
      metadataUri: registered.metadataUri,
      capabilities: registered.capabilities,
      blockNumber,
      txHash: registered.txHash,
      logIndex: 0,
    });

    // 3. SDK -> gateway: browse() must surface the seeded row.
    const page = await fx.agent.browse({ kind: "acp", limit: 50 });
    const found = page.data.find((a) => a.id === insertedPk);
    expect(found).toBeDefined();
    expect(found?.agent_id).toBe(registered.agentId.toString());
    expect(found?.controller).toBe(registered.controller);
    expect(found?.metadata_uri).toBe("ipfs://bafyroundtrip");
    expect(found?.tx_hash).toBe(registered.txHash);
  });
});
