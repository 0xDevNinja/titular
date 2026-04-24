import { describe, expect, it } from "vitest";
import {
  GRADUATION_THRESHOLD,
  VIRTUAL_AGENT,
  VIRTUAL_TITU,
  buildCurvePoints,
  getAmountOut,
  getPrice,
  graduationProgress,
} from "../../lib/bondingCurve";

// ─── getAmountOut: buy direction ──────────────────────────────────────────────

describe("getAmountOut — buy (TITU → agent)", () => {
  it("returns 0 for zero input", () => {
    expect(getAmountOut(true, 0, VIRTUAL_TITU, VIRTUAL_AGENT, 0, 0)).toBe(0);
  });

  it("returns 0 for negative input", () => {
    expect(getAmountOut(true, -10, VIRTUAL_TITU, VIRTUAL_AGENT, 0, 0)).toBe(0);
  });

  it("returns a positive amount for positive input", () => {
    const out = getAmountOut(true, 100, VIRTUAL_TITU, VIRTUAL_AGENT, 0, 0);
    expect(out).toBeGreaterThan(0);
  });

  it("applies the 1% fee (effective input = 99% of raw)", () => {
    const xReserve = VIRTUAL_TITU; // realTitu = 0
    const yReserve = VIRTUAL_AGENT;
    const k = xReserve * yReserve;

    const amountIn = 100;
    const effectiveIn = amountIn * 0.99;
    const newX = xReserve + effectiveIn;
    const newY = k / newX;
    const expected = yReserve - newY;

    const result = getAmountOut(true, amountIn, VIRTUAL_TITU, VIRTUAL_AGENT, 0, 0);
    expect(result).toBeCloseTo(expected, 6);
  });

  it("larger input yields more output (monotonic)", () => {
    const out100 = getAmountOut(true, 100, VIRTUAL_TITU, VIRTUAL_AGENT, 5000, 145_800_000);
    const out200 = getAmountOut(true, 200, VIRTUAL_TITU, VIRTUAL_AGENT, 5000, 145_800_000);
    expect(out200).toBeGreaterThan(out100);
  });

  it("output is strictly less without fee than with fee applied", () => {
    // With fee, effectiveIn < amountIn, so output < no-fee output.
    const withFee = getAmountOut(true, 1000, VIRTUAL_TITU, VIRTUAL_AGENT, 0, 0);
    // Manual no-fee version
    const x = VIRTUAL_TITU;
    const y = VIRTUAL_AGENT;
    const k = x * y;
    const noFeeOut = y - k / (x + 1000);
    expect(withFee).toBeLessThan(noFeeOut);
  });
});

// ─── getAmountOut: sell direction ─────────────────────────────────────────────

describe("getAmountOut — sell (agent → TITU)", () => {
  it("returns 0 for zero input", () => {
    expect(getAmountOut(false, 0, VIRTUAL_TITU, VIRTUAL_AGENT, 5000, 145_800_000)).toBe(0);
  });

  it("returns positive TITU for agent tokens in", () => {
    const out = getAmountOut(false, 1_000_000, VIRTUAL_TITU, VIRTUAL_AGENT, 5000, 145_800_000);
    expect(out).toBeGreaterThan(0);
  });

  it("larger sell yields more TITU (monotonic)", () => {
    const out1M = getAmountOut(false, 1_000_000, VIRTUAL_TITU, VIRTUAL_AGENT, 5000, 145_800_000);
    const out2M = getAmountOut(false, 2_000_000, VIRTUAL_TITU, VIRTUAL_AGENT, 5000, 145_800_000);
    expect(out2M).toBeGreaterThan(out1M);
  });
});

// ─── Price function ───────────────────────────────────────────────────────────

describe("getPrice", () => {
  it("returns initial price at zero real reserves", () => {
    const price = getPrice(VIRTUAL_TITU, VIRTUAL_AGENT, 0, 0);
    expect(price).toBeCloseTo(VIRTUAL_TITU / VIRTUAL_AGENT, 10);
  });

  it("returns 0 when total agent reserve (virtual + real) is 0", () => {
    // Both virtual and real agent are 0 → division by zero guard triggers.
    expect(getPrice(0, 0, 100, 0)).toBe(0);
  });

  it("price increases as real TITU increases (monotonic)", () => {
    // Simulate curve: as realTitu increases, realAgent decreases
    const k = VIRTUAL_TITU * VIRTUAL_AGENT;
    const steps = [0, 5000, 15000, 30000, 41000];
    const prices = steps.map((realTitu) => {
      const totalX = VIRTUAL_TITU + realTitu;
      const totalY = k / totalX;
      const realAgent = Math.max(0, totalY - VIRTUAL_AGENT);
      return getPrice(VIRTUAL_TITU, VIRTUAL_AGENT, realTitu, realAgent);
    });
    for (let i = 1; i < prices.length; i++) {
      expect(prices[i]).toBeGreaterThan(prices[i - 1]);
    }
  });
});

// ─── buildCurvePoints ────────────────────────────────────────────────────────

describe("buildCurvePoints", () => {
  it("starts at realTitu=0", () => {
    const pts = buildCurvePoints(VIRTUAL_TITU, VIRTUAL_AGENT, GRADUATION_THRESHOLD, 10);
    expect(pts[0].realTitu).toBe(0);
  });

  it("ends at the graduation threshold", () => {
    const pts = buildCurvePoints(VIRTUAL_TITU, VIRTUAL_AGENT, GRADUATION_THRESHOLD, 10);
    expect(pts[pts.length - 1].realTitu).toBe(GRADUATION_THRESHOLD);
  });

  it("returns steps+1 points", () => {
    const pts = buildCurvePoints(VIRTUAL_TITU, VIRTUAL_AGENT, GRADUATION_THRESHOLD, 50);
    expect(pts).toHaveLength(51);
  });

  it("prices are strictly increasing (monotonic curve)", () => {
    const pts = buildCurvePoints(VIRTUAL_TITU, VIRTUAL_AGENT, GRADUATION_THRESHOLD, 100);
    for (let i = 1; i < pts.length; i++) {
      expect(pts[i].price).toBeGreaterThan(pts[i - 1].price);
    }
  });

  it("all prices are positive", () => {
    const pts = buildCurvePoints(VIRTUAL_TITU, VIRTUAL_AGENT, GRADUATION_THRESHOLD, 50);
    for (const pt of pts) {
      expect(pt.price).toBeGreaterThan(0);
    }
  });
});

// ─── graduationProgress ──────────────────────────────────────────────────────

describe("graduationProgress", () => {
  it("returns 0 at zero realTitu", () => {
    expect(graduationProgress(0)).toBe(0);
  });

  it("returns 0.5 at half threshold", () => {
    expect(graduationProgress(GRADUATION_THRESHOLD / 2)).toBeCloseTo(0.5);
  });

  it("returns 1.0 at threshold", () => {
    expect(graduationProgress(GRADUATION_THRESHOLD)).toBe(1);
  });

  it("clamps to 1.0 above threshold", () => {
    expect(graduationProgress(GRADUATION_THRESHOLD * 2)).toBe(1);
  });

  it("uses custom threshold", () => {
    expect(graduationProgress(10, 100)).toBeCloseTo(0.1);
  });
});

// ─── Graduation threshold math ────────────────────────────────────────────────

describe("graduation threshold math", () => {
  it("at exactly 42,000 TITU raised, progress === 1", () => {
    expect(graduationProgress(42_000, 42_000)).toBe(1);
  });

  it("buy output is never negative even near graduation", () => {
    // Very large buy pushing close to (but not past) the reserve
    const out = getAmountOut(
      true,
      41_999, // just under graduation
      VIRTUAL_TITU,
      VIRTUAL_AGENT,
      0,
      0
    );
    expect(out).toBeGreaterThan(0);
  });

  it("sell output is never negative with maximum agent tokens", () => {
    const out = getAmountOut(
      false,
      VIRTUAL_AGENT, // sell all virtual agent tokens
      VIRTUAL_TITU,
      VIRTUAL_AGENT,
      5000,
      145_800_000
    );
    expect(out).toBeGreaterThanOrEqual(0);
  });
});
