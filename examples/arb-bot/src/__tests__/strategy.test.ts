import { describe, expect, it } from "vitest";
import {
  PriceBook,
  type PriceQuote,
  type StrategyConfig,
  type TradeEvent,
  compareQuotes,
  detectArb,
  quoteFromTrade,
  spreadBps,
} from "../strategy.js";

const TOKEN_A = "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";
const TOKEN_B = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb";
const TOKEN_C = "0xcccccccccccccccccccccccccccccccccccccccc";

function buy(agentToken: string, quoteIn: string, agentOut: string, blockNumber = 1): TradeEvent {
  return { agentToken, side: "buy", quoteIn, agentOut, blockNumber };
}

function sell(agentToken: string, quoteOut: string, agentIn: string, blockNumber = 1): TradeEvent {
  return { agentToken, side: "sell", quoteOut, agentIn, blockNumber };
}

describe("quoteFromTrade", () => {
  it("derives a price quote from a buy", () => {
    const q = quoteFromTrade(buy(TOKEN_A, "1000000", "500", 42));
    expect(q).not.toBeNull();
    expect(q?.agentToken).toBe(TOKEN_A);
    expect(q?.numerator).toBe(1_000_000n);
    expect(q?.denominator).toBe(500n);
    expect(q?.blockNumber).toBe(42);
  });

  it("derives a price quote from a sell", () => {
    const q = quoteFromTrade(sell(TOKEN_A, "750000", "300", 99));
    expect(q?.numerator).toBe(750_000n);
    expect(q?.denominator).toBe(300n);
  });

  it("lowercases mixed-case agent token addresses", () => {
    const mixed = "0xAaAaAaAaAaAaAaAaAaAaAaAaAaAaAaAaAaAaAaAa";
    const q = quoteFromTrade(buy(mixed, "1", "1"));
    expect(q?.agentToken).toBe(TOKEN_A);
  });

  it("returns null when buy fields are missing", () => {
    const t: TradeEvent = { agentToken: TOKEN_A, side: "buy", blockNumber: 1 };
    expect(quoteFromTrade(t)).toBeNull();
  });

  it("returns null when sell fields are missing", () => {
    const t: TradeEvent = { agentToken: TOKEN_A, side: "sell", blockNumber: 1 };
    expect(quoteFromTrade(t)).toBeNull();
  });

  it("returns null when amounts are zero (would divide by zero)", () => {
    expect(quoteFromTrade(buy(TOKEN_A, "0", "100"))).toBeNull();
    expect(quoteFromTrade(buy(TOKEN_A, "100", "0"))).toBeNull();
  });

  it("returns null on non-decimal junk values", () => {
    const bad: TradeEvent = {
      agentToken: TOKEN_A,
      side: "buy",
      quoteIn: "0xabc",
      agentOut: "10",
      blockNumber: 1,
    };
    expect(quoteFromTrade(bad)).toBeNull();
  });

  it("preserves uint256-scale amounts without precision loss", () => {
    // 2^200 is well above 2^53; coerced as bigint must round-trip exactly.
    const huge = (1n << 200n).toString();
    const q = quoteFromTrade(buy(TOKEN_A, huge, "1"));
    expect(q?.numerator).toBe(1n << 200n);
  });
});

describe("PriceBook", () => {
  it("retains the most recent quote per token", () => {
    const book = new PriceBook();
    const q1: PriceQuote = {
      agentToken: TOKEN_A,
      numerator: 100n,
      denominator: 1n,
      blockNumber: 1,
    };
    const q2: PriceQuote = {
      agentToken: TOKEN_A,
      numerator: 200n,
      denominator: 1n,
      blockNumber: 5,
    };
    book.observe(q1);
    book.observe(q2);
    expect(book.get(TOKEN_A)?.numerator).toBe(200n);
  });

  it("ignores stale observations (lower block number)", () => {
    const book = new PriceBook();
    book.observe({
      agentToken: TOKEN_A,
      numerator: 200n,
      denominator: 1n,
      blockNumber: 5,
    });
    book.observe({
      agentToken: TOKEN_A,
      numerator: 100n,
      denominator: 1n,
      blockNumber: 1,
    });
    expect(book.get(TOKEN_A)?.numerator).toBe(200n);
  });

  it("looks up case-insensitively", () => {
    const book = new PriceBook();
    book.observe({
      agentToken: TOKEN_A,
      numerator: 1n,
      denominator: 1n,
      blockNumber: 1,
    });
    expect(book.get(TOKEN_A.toUpperCase())).toBeDefined();
  });

  it("snapshot returns one entry per token", () => {
    const book = new PriceBook();
    book.observe({
      agentToken: TOKEN_A,
      numerator: 1n,
      denominator: 1n,
      blockNumber: 1,
    });
    book.observe({
      agentToken: TOKEN_B,
      numerator: 2n,
      denominator: 1n,
      blockNumber: 1,
    });
    book.observe({
      agentToken: TOKEN_A,
      numerator: 3n,
      denominator: 1n,
      blockNumber: 2,
    });
    expect(book.snapshot()).toHaveLength(2);
  });
});

describe("spreadBps + compareQuotes", () => {
  it("computes a 0-bps spread for equal prices", () => {
    const a: PriceQuote = {
      agentToken: TOKEN_A,
      numerator: 100n,
      denominator: 1n,
      blockNumber: 1,
    };
    const b: PriceQuote = {
      agentToken: TOKEN_B,
      numerator: 200n,
      denominator: 2n,
      blockNumber: 1,
    };
    expect(spreadBps(a, b)).toBe(0);
    expect(compareQuotes(a, b)).toBe(0);
  });

  it("computes 100 bps for a 1% spread", () => {
    // 100 vs 101 → ~1% spread
    const cheap: PriceQuote = {
      agentToken: TOKEN_A,
      numerator: 100n,
      denominator: 1n,
      blockNumber: 1,
    };
    const rich: PriceQuote = {
      agentToken: TOKEN_B,
      numerator: 101n,
      denominator: 1n,
      blockNumber: 1,
    };
    expect(spreadBps(cheap, rich)).toBe(100);
    expect(compareQuotes(cheap, rich)).toBe(-1);
    expect(compareQuotes(rich, cheap)).toBe(1);
  });

  it("is symmetric", () => {
    const a: PriceQuote = {
      agentToken: TOKEN_A,
      numerator: 90n,
      denominator: 1n,
      blockNumber: 1,
    };
    const b: PriceQuote = {
      agentToken: TOKEN_B,
      numerator: 100n,
      denominator: 1n,
      blockNumber: 1,
    };
    expect(spreadBps(a, b)).toBe(spreadBps(b, a));
  });
});

describe("detectArb", () => {
  const config: StrategyConfig = {
    minSpreadBps: 50,
    watchedTokens: [TOKEN_A, TOKEN_B, TOKEN_C],
  };

  it("returns null when fewer than two tokens are quoted", () => {
    const book = new PriceBook();
    book.observe({
      agentToken: TOKEN_A,
      numerator: 100n,
      denominator: 1n,
      blockNumber: 1,
    });
    expect(detectArb(book, config)).toBeNull();
  });

  it("returns null when the spread is below threshold", () => {
    const book = new PriceBook();
    book.observe({
      agentToken: TOKEN_A,
      numerator: 1000n,
      denominator: 1n,
      blockNumber: 1,
    });
    book.observe({
      agentToken: TOKEN_B,
      numerator: 1001n,
      denominator: 1n,
      blockNumber: 1,
    });
    // spread is 10 bps, threshold is 50 bps → ignored
    expect(detectArb(book, config)).toBeNull();
  });

  it("emits an opportunity when threshold is crossed", () => {
    const book = new PriceBook();
    book.observe({
      agentToken: TOKEN_A,
      numerator: 100n,
      denominator: 1n,
      blockNumber: 1,
    });
    book.observe({
      agentToken: TOKEN_B,
      numerator: 101n,
      denominator: 1n,
      blockNumber: 1,
    });
    const opp = detectArb(book, { ...config, minSpreadBps: 50 });
    expect(opp).not.toBeNull();
    expect(opp?.cheapToken).toBe(TOKEN_A);
    expect(opp?.richToken).toBe(TOKEN_B);
    expect(opp?.spreadBps).toBe(100);
  });

  it("picks the widest spread when multiple pairs qualify", () => {
    const book = new PriceBook();
    book.observe({
      agentToken: TOKEN_A,
      numerator: 100n,
      denominator: 1n,
      blockNumber: 1,
    });
    book.observe({
      agentToken: TOKEN_B,
      numerator: 101n,
      denominator: 1n,
      blockNumber: 1,
    });
    book.observe({
      agentToken: TOKEN_C,
      numerator: 110n,
      denominator: 1n,
      blockNumber: 1,
    });
    const opp = detectArb(book, config);
    expect(opp?.cheapToken).toBe(TOKEN_A);
    expect(opp?.richToken).toBe(TOKEN_C);
  });

  it("ignores tokens not in the watch list", () => {
    const ROGUE = "0xdddddddddddddddddddddddddddddddddddddddd";
    const book = new PriceBook();
    book.observe({
      agentToken: TOKEN_A,
      numerator: 100n,
      denominator: 1n,
      blockNumber: 1,
    });
    book.observe({
      agentToken: ROGUE,
      numerator: 1000n,
      denominator: 1n,
      blockNumber: 1,
    });
    expect(detectArb(book, config)).toBeNull();
  });

  it("is deterministic — same inputs return identical opportunities", () => {
    const book = new PriceBook();
    book.observe({
      agentToken: TOKEN_A,
      numerator: 100n,
      denominator: 1n,
      blockNumber: 1,
    });
    book.observe({
      agentToken: TOKEN_B,
      numerator: 105n,
      denominator: 1n,
      blockNumber: 1,
    });
    const a = detectArb(book, config);
    const b = detectArb(book, config);
    expect(a).toEqual(b);
  });
});
