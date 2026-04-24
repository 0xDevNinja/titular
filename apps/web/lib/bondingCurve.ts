/**
 * Bonding curve math matching BondingCurve.sol constant-product formula.
 *
 * Invariant: (virtualTitu + realTitu) * (virtualAgent + realAgent) = k
 *
 * Virtual reserves (Virtuals-style):
 *   virtualTitu  = 30,000  TITU  (18-decimal amounts are passed as plain numbers here)
 *   virtualAgent = 1,073,000,000 agent tokens
 *
 * Graduation threshold: 42,000 real TITU raised.
 *
 * All amounts here are in display units (no wei scaling) to keep the FE math
 * simple. Contract-level calls will apply 1e18 scaling.
 */

export const VIRTUAL_TITU = 30_000;
export const VIRTUAL_AGENT = 1_073_000_000;
export const GRADUATION_THRESHOLD = 42_000;

export interface CurveState {
  virtualTitu: number;
  virtualAgent: number;
  realTitu: number;
  realAgent: number;
}

/**
 * Compute the amount of tokens received for a given input, applying the
 * constant-product formula with a 1% fee taken from the input.
 *
 * Formula (no fee):
 *   k  = (vX + rX) * (vY + rY)
 *   dy = (vY + rY) - k / (vX + rX + dx)
 *
 * 1% fee: effective_dx = dx * 0.99
 *
 * @param isBuy   true  → spend TITU, receive agent tokens
 *                false → spend agent tokens, receive TITU
 */
export function getAmountOut(
  isBuy: boolean,
  amountIn: number,
  virtualTitu: number,
  virtualAgent: number,
  realTitu: number,
  realAgent: number
): number {
  if (amountIn <= 0) return 0;

  const effectiveIn = amountIn * 0.99; // 1% fee

  if (isBuy) {
    // dx = TITU in, dy = agent out
    const xReserve = virtualTitu + realTitu;
    const yReserve = virtualAgent + realAgent;
    const k = xReserve * yReserve;
    const newX = xReserve + effectiveIn;
    const newY = k / newX;
    const amountOut = yReserve - newY;
    return Math.max(0, amountOut);
  }
  // dx = agent in, dy = TITU out
  const xReserve = virtualAgent + realAgent;
  const yReserve = virtualTitu + realTitu;
  const k = xReserve * yReserve;
  const newX = xReserve + effectiveIn;
  const newY = k / newX;
  const amountOut = yReserve - newY;
  return Math.max(0, amountOut);
}

/**
 * Marginal price of agent token denominated in TITU.
 * price = (virtualTitu + realTitu) / (virtualAgent + realAgent)
 */
export function getPrice(
  virtualTitu: number,
  virtualAgent: number,
  realTitu: number,
  realAgent: number
): number {
  const xReserve = virtualTitu + realTitu;
  const yReserve = virtualAgent + realAgent;
  if (yReserve === 0) return 0;
  return xReserve / yReserve;
}

/**
 * Build an array of (realTitu, price) points for a given number of steps,
 * ranging from realTitu=0 to realTitu=graduationThreshold.
 * Used to render the static curve shape.
 */
export function buildCurvePoints(
  virtualTitu: number,
  virtualAgent: number,
  graduationThreshold: number,
  steps = 100
): Array<{ realTitu: number; price: number }> {
  const points: Array<{ realTitu: number; price: number }> = [];
  for (let i = 0; i <= steps; i++) {
    const realTitu = (graduationThreshold * i) / steps;
    // When realTitu increases, realAgent decreases to maintain k.
    // realAgent = k / (virtualTitu + realTitu) - virtualAgent
    const k = virtualTitu * virtualAgent; // k at real=0, real_agent=0
    const newX = virtualTitu + realTitu;
    const totalY = k / newX;
    const realAgent = Math.max(0, totalY - virtualAgent);
    const price = getPrice(virtualTitu, virtualAgent, realTitu, realAgent);
    points.push({ realTitu, price });
  }
  return points;
}

/**
 * Returns progress toward graduation as a fraction [0, 1].
 */
export function graduationProgress(realTitu: number, threshold = GRADUATION_THRESHOLD): number {
  return Math.min(1, realTitu / threshold);
}
