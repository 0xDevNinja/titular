"use client";

import { buildCurvePoints } from "../../../lib/bondingCurve";
import type { AgentFixture, PricePoint } from "./types";

interface Props {
  agent: AgentFixture;
}

const WIDTH = 560;
const HEIGHT = 200;
const PAD_LEFT = 56;
const PAD_RIGHT = 16;
const PAD_TOP = 12;
const PAD_BOTTOM = 32;

const CHART_W = WIDTH - PAD_LEFT - PAD_RIGHT;
const CHART_H = HEIGHT - PAD_TOP - PAD_BOTTOM;

function toSvgCoords(
  points: Array<{ x: number; y: number }>,
  xMin: number,
  xMax: number,
  yMin: number,
  yMax: number
): string {
  if (points.length === 0) return "";
  const xRange = xMax - xMin || 1;
  const yRange = yMax - yMin || 1;
  return points
    .map(({ x, y }, i) => {
      const cx = PAD_LEFT + ((x - xMin) / xRange) * CHART_W;
      const cy = PAD_TOP + CHART_H - ((y - yMin) / yRange) * CHART_H;
      return `${i === 0 ? "M" : "L"}${cx.toFixed(2)},${cy.toFixed(2)}`;
    })
    .join(" ");
}

function formatPrice(p: number): string {
  if (p === 0) return "0";
  if (p < 0.0001) return p.toExponential(2);
  return p.toFixed(6);
}

export default function BondingCurveChart({ agent }: Props) {
  const { virtualTituReserve, virtualAgentReserve, realTituReserve, graduationThreshold } = agent;

  // Static bonding curve shape (theoretical full curve)
  const curvePoints = buildCurvePoints(
    virtualTituReserve,
    virtualAgentReserve,
    graduationThreshold,
    120
  );

  // Price history scatter/line
  const history: PricePoint[] = agent.priceHistory;

  // Compute domains
  const allPrices = curvePoints.map((p) => p.price).concat(history.map((h) => h.price));
  const priceMin = 0;
  const priceMax = Math.max(...allPrices) * 1.05;

  // Curve path (price vs realTitu)
  const curveSvgPoints = curvePoints.map((p) => ({ x: p.realTitu, y: p.price }));
  const curvePath = toSvgCoords(curveSvgPoints, 0, graduationThreshold, priceMin, priceMax);

  // History path (price vs time, normalised to same x domain)
  const tMin = history.length ? history[0].t : 0;
  const tMax = history.length ? history[history.length - 1].t : 1;
  const histSvgPoints = history.map((h) => ({ x: h.t, y: h.price }));
  const histPath = toSvgCoords(histSvgPoints, tMin, tMax, priceMin, priceMax);

  // Current position marker on the curve
  const currentCurveX = PAD_LEFT + (realTituReserve / graduationThreshold) * CHART_W;
  const currentCurveYFrac = (agent.currentPrice - priceMin) / (priceMax - priceMin || 1);
  const currentCurveY = PAD_TOP + CHART_H - currentCurveYFrac * CHART_H;

  // Y-axis labels (3 ticks)
  const yTicks = [0, 0.5, 1].map((frac) => ({
    y: PAD_TOP + CHART_H - frac * CHART_H,
    label: formatPrice(priceMin + frac * (priceMax - priceMin)),
  }));

  // X-axis labels (graduation threshold)
  const xTicks = [0, 0.5, 1].map((frac) => ({
    x: PAD_LEFT + frac * CHART_W,
    label: `${Math.round(frac * graduationThreshold).toLocaleString()}`,
  }));

  return (
    <figure aria-label='Bonding curve price chart' style={{ margin: 0 }}>
      <svg
        viewBox={`0 0 ${WIDTH} ${HEIGHT}`}
        width='100%'
        style={{ display: "block", maxWidth: WIDTH }}
        role='img'
        aria-label={`${agent.name} bonding curve. Current real reserve: ${realTituReserve} TITU. Graduation at ${graduationThreshold} TITU.`}
      >
        {/* Grid lines */}
        {yTicks.map((tick) => (
          <line
            key={tick.label}
            x1={PAD_LEFT}
            x2={WIDTH - PAD_RIGHT}
            y1={tick.y}
            y2={tick.y}
            stroke='#e2e8f0'
            strokeWidth={1}
          />
        ))}

        {/* Theoretical curve (grey) */}
        <path d={curvePath} fill='none' stroke='#94a3b8' strokeWidth={1.5} />

        {/* Price history (blue) */}
        {histPath && <path d={histPath} fill='none' stroke='#2563eb' strokeWidth={2} />}

        {/* Current position dot */}
        <circle
          cx={currentCurveX}
          cy={currentCurveY}
          r={5}
          fill='#2563eb'
          stroke='#fff'
          strokeWidth={2}
          aria-label={`Current price: ${formatPrice(agent.currentPrice)} TITU`}
        />

        {/* Y-axis labels */}
        {yTicks.map((tick) => (
          <text
            key={tick.label}
            x={PAD_LEFT - 4}
            y={tick.y + 4}
            textAnchor='end'
            fontSize={10}
            fill='#64748b'
          >
            {tick.label}
          </text>
        ))}

        {/* X-axis labels */}
        {xTicks.map((tick) => (
          <text
            key={tick.label}
            x={tick.x}
            y={HEIGHT - 4}
            textAnchor='middle'
            fontSize={10}
            fill='#64748b'
          >
            {tick.label}
          </text>
        ))}

        {/* Axis labels */}
        <text x={PAD_LEFT - 4} y={PAD_TOP - 2} textAnchor='end' fontSize={9} fill='#94a3b8'>
          TITU
        </text>
        <text x={WIDTH - PAD_RIGHT} y={HEIGHT - 4} textAnchor='end' fontSize={9} fill='#94a3b8'>
          raised TITU
        </text>
      </svg>

      <figcaption style={{ fontSize: 12, color: "#64748b", marginTop: 4, textAlign: "center" }}>
        Price (TITU) vs. TITU raised — <span style={{ color: "#94a3b8" }}>grey = curve shape</span>,{" "}
        <span style={{ color: "#2563eb" }}>blue = history</span>
      </figcaption>
    </figure>
  );
}
