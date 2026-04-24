"use server";

import launchResponse from "../../../fixtures/launch-response.json";
import type { LaunchParams } from "./types";

export interface LaunchResult {
  txHash: string;
  agentId: string;
  agentSlug: string;
  tokenAddress: string;
  curveAddress: string;
  gasUsed: number;
  blockNumber: number;
}

export async function submitLaunch(_params: LaunchParams): Promise<LaunchResult> {
  // Fixture stub - no real wallet/chain call.
  // Issues #50 and #51 will wire real Wagmi transaction here.
  await new Promise((resolve) => setTimeout(resolve, 800));
  return launchResponse as LaunchResult;
}
