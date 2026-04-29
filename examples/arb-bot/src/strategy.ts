// strategy.ts — pure functions implementing the arb-bot's decision logic.
//
// Kept dependency-free (no SDK, no fetch, no chain calls) so it round-trips
// through vitest without any IO. The main loop in `index.ts` is responsible
// for feeding live trade events in and acting on the returned opportunities.
//
// Pricing model:
//   - We treat each `Trade` as a quote-per-agent observation. For a buy,
//     `quoteIn / agentOut` is the effective price the trader paid. For a
//     sell, `quoteOut / agentIn` is the effective price the trader received.
//   - We hold the most recent observation per agent token and compute the
//     spread in basis points between any two tokens. The "rich" leg is the
//     one with the higher price; the "cheap" leg is the lower.
//   - Spread is symmetric (`|a-b| / min(a,b)`) so the order in which two
//     observations arrive doesn't change the signal.
//
// This is intentionally naive — a production bot would model fee tiers,
// slippage, gas, and inventory risk. The goal here is "smallest realistic
// state machine an agent author can read in one sitting".

/** Side of a trade as the gateway emits it (`store.Trade.side`). */
export type TradeSide = "buy" | "sell";

/**
 * Subset of `store.Trade` the strategy actually consumes. Inputs are decimal-
 * encoded uint256 strings; the strategy never coerces them to JS numbers
 * (precision loss above 2^53).
 */
export interface TradeEvent {
  /** Agent-token address, lowercase 0x-hex. */
  agentToken: string;
  side: TradeSide;
  /** Decimal-encoded uint256 amount of the quote asset entering the curve. */
  quoteIn?: string;
  /** Decimal-encoded uint256 amount of the agent token leaving the curve. */
  agentOut?: string;
  /** Decimal-encoded uint256 amount of the agent token entering the curve. */
  agentIn?: string;
  /** Decimal-encoded uint256 amount of the quote asset leaving the curve. */
  quoteOut?: string;
  /** Block number the trade landed in (used as a tiebreaker for "most recent"). */
  blockNumber: number;
}

/** A single price observation derived from a {@link TradeEvent}. */
export interface PriceQuote {
  agentToken: string;
  /** quote-per-agent, as a `bigint` rational (numerator / denominator). */
  numerator: bigint;
  denominator: bigint;
  blockNumber: number;
}

/** A detected arb opportunity between two agent tokens. */
export interface ArbOpportunity {
  cheapToken: string;
  richToken: string;
  /** Spread between the two prices, in basis points (1 bps = 0.01%). */
  spreadBps: number;
  /** The two quotes that produced the signal. */
  cheap: PriceQuote;
  rich: PriceQuote;
}

/** Configuration for the strategy. */
export interface StrategyConfig {
  /** Minimum spread (basis points) to surface as an opportunity. */
  minSpreadBps: number;
  /** Lowercased addresses of the agent tokens the bot is allowed to act on. */
  watchedTokens: readonly string[];
}

/**
 * Convert a {@link TradeEvent} to a {@link PriceQuote}. Returns `null` when
 * the trade carries insufficient data (e.g. zero amounts) — those events
 * MUST NOT update the price book.
 */
export function quoteFromTrade(t: TradeEvent): PriceQuote | null {
  if (t.side === "buy") {
    if (t.quoteIn === undefined || t.agentOut === undefined) return null;
    const num = parseDecimal(t.quoteIn);
    const den = parseDecimal(t.agentOut);
    if (num === null || den === null || den === 0n || num === 0n) return null;
    return {
      agentToken: t.agentToken.toLowerCase(),
      numerator: num,
      denominator: den,
      blockNumber: t.blockNumber,
    };
  }
  // sell
  if (t.quoteOut === undefined || t.agentIn === undefined) return null;
  const num = parseDecimal(t.quoteOut);
  const den = parseDecimal(t.agentIn);
  if (num === null || den === null || den === 0n || num === 0n) return null;
  return {
    agentToken: t.agentToken.toLowerCase(),
    numerator: num,
    denominator: den,
    blockNumber: t.blockNumber,
  };
}

/**
 * Stateful price book — keeps the most recent quote per agent token (most
 * recent meaning highest `blockNumber` seen). Pure data structure; the main
 * loop owns the lifetime.
 */
export class PriceBook {
  readonly #quotes = new Map<string, PriceQuote>();

  /** Replace the cached quote for `q.agentToken` if `q` is newer. */
  observe(q: PriceQuote): void {
    const key = q.agentToken.toLowerCase();
    const prev = this.#quotes.get(key);
    if (prev === undefined || q.blockNumber >= prev.blockNumber) {
      this.#quotes.set(key, q);
    }
  }

  /** Snapshot of all known quotes. Order is insertion order. */
  snapshot(): PriceQuote[] {
    return Array.from(this.#quotes.values());
  }

  /** Latest quote for `token`, or `undefined`. */
  get(token: string): PriceQuote | undefined {
    return this.#quotes.get(token.toLowerCase());
  }
}

/**
 * Scan the price book for an arb opportunity exceeding `minSpreadBps`.
 *
 * Returns the highest-spread pair drawn from `watchedTokens`. When no pair
 * meets the threshold, returns `null`. The function is pure — given the same
 * book + config, it always returns the same opportunity.
 */
export function detectArb(book: PriceBook, config: StrategyConfig): ArbOpportunity | null {
  const watched = config.watchedTokens.map((t) => t.toLowerCase());
  const quotes = watched.map((t) => book.get(t)).filter((q): q is PriceQuote => q !== undefined);
  if (quotes.length < 2) return null;

  let best: ArbOpportunity | null = null;
  for (let i = 0; i < quotes.length; i += 1) {
    for (let j = i + 1; j < quotes.length; j += 1) {
      const a = quotes[i]!;
      const b = quotes[j]!;
      const bps = spreadBps(a, b);
      if (bps < config.minSpreadBps) continue;
      const [cheap, rich] = compareQuotes(a, b) < 0 ? [a, b] : [b, a];
      const candidate: ArbOpportunity = {
        cheapToken: cheap.agentToken,
        richToken: rich.agentToken,
        spreadBps: bps,
        cheap,
        rich,
      };
      if (best === null || candidate.spreadBps > best.spreadBps) {
        best = candidate;
      }
    }
  }
  return best;
}

/**
 * Compute the spread between two price quotes in basis points. Symmetric:
 * `spreadBps(a, b) === spreadBps(b, a)`.
 */
export function spreadBps(a: PriceQuote, b: PriceQuote): number {
  // a = aN/aD, b = bN/bD. |a-b| / min(a,b) in basis points
  // |aN*bD - bN*aD| / min(aN*bD, bN*aD) * 10000
  const aCross = a.numerator * b.denominator;
  const bCross = b.numerator * a.denominator;
  const diff = aCross > bCross ? aCross - bCross : bCross - aCross;
  const minCross = aCross < bCross ? aCross : bCross;
  if (minCross === 0n) return 0;
  // Multiply by 10000 BEFORE the divide to keep precision; then floor.
  const bps = (diff * 10000n) / minCross;
  // bps is bounded by reasonable spreads; safe to coerce to number for the
  // configurable threshold which is itself a JS number.
  if (bps > BigInt(Number.MAX_SAFE_INTEGER)) return Number.MAX_SAFE_INTEGER;
  return Number(bps);
}

/** -1 if a < b, 1 if a > b, 0 if equal. Compares rationals without floating point. */
export function compareQuotes(a: PriceQuote, b: PriceQuote): -1 | 0 | 1 {
  const left = a.numerator * b.denominator;
  const right = b.numerator * a.denominator;
  if (left < right) return -1;
  if (left > right) return 1;
  return 0;
}

function parseDecimal(s: string): bigint | null {
  if (!/^(0|[1-9][0-9]*)$/.test(s)) return null;
  return BigInt(s);
}
