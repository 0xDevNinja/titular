"use client";

import { useState } from "react";
import { getAmountOut } from "../../../lib/bondingCurve";
import type { TradeResult } from "./actions";
import { submitTrade } from "./actions";
import type { AgentFixture } from "./types";

interface Props {
  agent: AgentFixture;
}

type Tab = "buy" | "sell";

function formatAmount(n: number, decimals = 4): string {
  if (n === 0) return "0";
  if (n < 0.0001) return n.toExponential(3);
  return n.toLocaleString("en-US", { maximumFractionDigits: decimals });
}

export default function TradeWidget({ agent }: Props) {
  const [tab, setTab] = useState<Tab>("buy");
  const [amountIn, setAmountIn] = useState("");
  const [slippage, setSlippage] = useState("1");
  const [pending, setPending] = useState(false);
  const [result, setResult] = useState<TradeResult | null>(null);
  const [error, setError] = useState<string | null>(null);

  const isBuy = tab === "buy";
  const numIn = Number(amountIn) || 0;
  const slippageBps = Math.round((Number(slippage) || 1) * 100);

  const estimatedOut = getAmountOut(
    isBuy,
    numIn,
    agent.virtualTituReserve,
    agent.virtualAgentReserve,
    agent.realTituReserve,
    agent.realAgentReserve
  );

  const minAmountOut = estimatedOut * (1 - slippageBps / 10_000);

  const inputLabel = isBuy ? "TITU to spend" : `${agent.symbol} to sell`;
  const outputLabel = isBuy ? agent.symbol : "TITU";

  function handleTabChange(next: Tab) {
    setTab(next);
    setAmountIn("");
    setResult(null);
    setError(null);
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (numIn <= 0) {
      setError("Enter an amount greater than 0.");
      return;
    }
    setPending(true);
    setError(null);
    setResult(null);
    try {
      const res = await submitTrade({
        slug: agent.slug,
        isBuy,
        amountIn: numIn,
        minAmountOut,
        slippageBps,
      });
      setResult(res);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error");
    } finally {
      setPending(false);
    }
  }

  return (
    <section aria-label='Trade widget' style={{ maxWidth: 420 }}>
      {/* Tabs */}
      <div
        role='tablist'
        aria-label='Buy or sell'
        style={{
          display: "flex",
          gap: 0,
          marginBottom: 16,
          borderRadius: 8,
          overflow: "hidden",
          border: "1px solid #e2e8f0",
        }}
      >
        {(["buy", "sell"] as const).map((t) => (
          <button
            key={t}
            role='tab'
            id={`tab-${t}`}
            aria-selected={tab === t}
            aria-controls='tabpanel-trade'
            type='button'
            onClick={() => handleTabChange(t)}
            style={{
              flex: 1,
              padding: "10px 0",
              fontWeight: tab === t ? 700 : 400,
              background: tab === t ? (t === "buy" ? "#16a34a" : "#dc2626") : "#f8fafc",
              color: tab === t ? "#fff" : "#374151",
              border: "none",
              cursor: "pointer",
              fontSize: 14,
              letterSpacing: 0.3,
            }}
          >
            {t === "buy" ? "Buy" : "Sell"}
          </button>
        ))}
      </div>

      <div id='tabpanel-trade' role='tabpanel' aria-labelledby={`tab-${tab}`}>
        {agent.graduated ? (
          <p
            style={{
              padding: "10px 14px",
              background: "#fef3c7",
              borderRadius: 6,
              fontSize: 13,
              color: "#92400e",
              marginBottom: 16,
            }}
            role='note'
          >
            This agent has graduated. Trades now route through Uniswap V2 (wiring in next
            iteration).
          </p>
        ) : null}

        <form onSubmit={handleSubmit} aria-label={`${isBuy ? "Buy" : "Sell"} form`} noValidate>
          {/* Input amount */}
          <div style={{ marginBottom: 16 }}>
            <label
              htmlFor='amount-in'
              style={{ display: "block", fontWeight: 600, marginBottom: 4, fontSize: 14 }}
            >
              {inputLabel} <span aria-hidden='true'>*</span>
            </label>
            <input
              id='amount-in'
              type='number'
              min='0'
              step='any'
              aria-required='true'
              aria-describedby='amount-out-preview'
              value={amountIn}
              onChange={(e) => {
                setAmountIn(e.target.value);
                setResult(null);
                setError(null);
              }}
              placeholder='0'
              disabled={pending}
              style={{ width: "100%", padding: "8px 12px", boxSizing: "border-box", fontSize: 15 }}
            />
          </div>

          {/* Estimated output */}
          <output
            id='amount-out-preview'
            aria-live='polite'
            aria-label={`Estimated ${outputLabel} received`}
            style={{
              display: "block",
              padding: "10px 14px",
              background: "#f1f5f9",
              borderRadius: 6,
              marginBottom: 16,
              fontSize: 13,
            }}
          >
            Estimated {outputLabel}:{" "}
            <strong>
              {numIn > 0 ? formatAmount(estimatedOut) : "—"} {outputLabel}
            </strong>
          </output>

          {/* Slippage */}
          <div style={{ marginBottom: 20 }}>
            <label
              htmlFor='slippage'
              style={{ display: "block", fontWeight: 600, marginBottom: 4, fontSize: 13 }}
            >
              Slippage tolerance (%)
            </label>
            <input
              id='slippage'
              type='number'
              min='0.1'
              max='50'
              step='0.1'
              value={slippage}
              onChange={(e) => setSlippage(e.target.value)}
              disabled={pending}
              style={{ width: 100, padding: "6px 10px" }}
            />
            {numIn > 0 && (
              <p style={{ fontSize: 11, color: "#6b7280", marginTop: 4 }}>
                Min received: {formatAmount(minAmountOut)} {outputLabel}
              </p>
            )}
          </div>

          {/* Error */}
          {error && (
            <p role='alert' style={{ color: "#dc2626", fontSize: 13, marginBottom: 12 }}>
              {error}
            </p>
          )}

          {/* Submit */}
          <button
            type='submit'
            disabled={pending || agent.graduated}
            aria-busy={pending}
            aria-disabled={agent.graduated}
            style={{
              width: "100%",
              padding: "12px 0",
              fontWeight: 700,
              fontSize: 15,
              background: agent.graduated ? "#94a3b8" : isBuy ? "#16a34a" : "#dc2626",
              color: "#fff",
              border: "none",
              borderRadius: 8,
              cursor: agent.graduated || pending ? "not-allowed" : "pointer",
              opacity: pending ? 0.7 : 1,
            }}
          >
            {pending ? "Submitting…" : agent.graduated ? "Graduated" : isBuy ? "Buy" : "Sell"}
          </button>
        </form>

        {/* Result */}
        {result && (
          <output
            aria-live='polite'
            aria-label='Trade result'
            style={{
              display: "block",
              marginTop: 16,
              padding: "12px 14px",
              background: "#f0fdf4",
              borderRadius: 8,
              fontSize: 13,
            }}
          >
            <p style={{ fontWeight: 700, color: "#16a34a", marginBottom: 8 }}>Trade submitted!</p>
            <dl style={{ margin: 0, lineHeight: 2 }}>
              <dt style={{ fontWeight: 600, display: "inline" }}>Tx hash: </dt>
              <dd style={{ display: "inline", wordBreak: "break-all" }}>
                <code>{result.txHash}</code>
              </dd>
            </dl>
          </output>
        )}
      </div>
    </section>
  );
}
