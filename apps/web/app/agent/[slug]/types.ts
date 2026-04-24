import { z } from "zod";

export const pricePointSchema = z.object({
  t: z.number(),
  price: z.number(),
});

export const agentFixtureSchema = z.object({
  slug: z.string(),
  name: z.string(),
  symbol: z.string(),
  description: z.string(),
  imageURI: z.string().url(),
  creator: z.string(),
  tokenAddress: z.string(),
  curveAddress: z.string(),
  uniswapPair: z.string().optional(),
  moduleBitmap: z.number().int().min(0).max(127),
  modules: z.array(z.string()),
  graduated: z.boolean(),
  virtualTituReserve: z.number(),
  virtualAgentReserve: z.number(),
  realTituReserve: z.number(),
  realAgentReserve: z.number(),
  graduationThreshold: z.number(),
  currentPrice: z.number(),
  marketCap: z.number(),
  priceHistory: z.array(pricePointSchema),
});

export type AgentFixture = z.infer<typeof agentFixtureSchema>;
export type PricePoint = z.infer<typeof pricePointSchema>;

export const MODULE_LABELS: Record<string, string> = {
  AntiSniper: "Anti-Sniper",
  SixtyDays: "60-Day Escrow",
  CapitalFormation: "Capital Formation",
  Airdrop: "Airdrop",
  LaunchRadar: "Launch Radar",
  PreBuy: "Pre-Buy",
  ExistingToken: "Existing Token",
};
