"use server";

import tradeResponse from "../../../fixtures/trade-response.json";

export interface TradeResult {
  txHash: string;
  amountIn: string;
  amountOut: string;
  isBuy: boolean;
  gasUsed: number;
  blockNumber: number;
}

export interface TradeParams {
  slug: string;
  isBuy: boolean;
  amountIn: number;
  minAmountOut: number;
  slippageBps: number;
}

export async function submitTrade(_params: TradeParams): Promise<TradeResult> {
  // Fixture stub — no real wallet or chain call.
  // Real Wagmi tx wiring follows in a later iteration.
  await new Promise((resolve) => setTimeout(resolve, 600));
  return tradeResponse as TradeResult;
}
