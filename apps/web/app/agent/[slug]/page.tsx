import fs from "node:fs";
import path from "node:path";
import { notFound } from "next/navigation";
import { graduationProgress } from "../../../lib/bondingCurve";
import BondingCurveChart from "./BondingCurveChart";
import TradeWidget from "./TradeWidget";
import { MODULE_LABELS, agentFixtureSchema } from "./types";

interface PageProps {
  params: Promise<{ slug: string }>;
}

function loadAgent(slug: string) {
  const fixtureDir = path.join(process.cwd(), "fixtures", "agents");
  const filePath = path.join(fixtureDir, `${slug}.json`);

  // Validate the slug resolves to a file inside the fixture directory
  // (prevents path traversal — canonical path must stay inside fixtureDir).
  const resolved = path.resolve(filePath);
  if (!resolved.startsWith(path.resolve(fixtureDir) + path.sep)) {
    return null;
  }

  if (!fs.existsSync(resolved)) return null;

  let raw: unknown;
  try {
    raw = JSON.parse(fs.readFileSync(resolved, "utf-8"));
  } catch {
    return null;
  }

  const parsed = agentFixtureSchema.safeParse(raw);
  return parsed.success ? parsed.data : null;
}

function formatPrice(p: number): string {
  if (p === 0) return "0";
  if (p < 0.0001) return p.toExponential(3);
  return p.toFixed(6);
}

function truncateAddress(addr: string): string {
  return `${addr.slice(0, 6)}…${addr.slice(-4)}`;
}

export async function generateMetadata({ params }: PageProps) {
  const { slug } = await params;
  const agent = loadAgent(slug);
  if (!agent) return { title: "Agent not found — Titular" };
  return {
    title: `${agent.name} (${agent.symbol}) — Titular`,
    description: agent.description,
  };
}

export default async function AgentPage({ params }: PageProps) {
  const { slug } = await params;
  const agent = loadAgent(slug);

  if (!agent) {
    notFound();
  }

  const progress = graduationProgress(agent.realTituReserve, agent.graduationThreshold);
  const progressPct = Math.round(progress * 100);

  return (
    <main style={{ maxWidth: 820, margin: "0 auto", padding: "40px 24px" }}>
      {/* ── Agent header ─────────────────────────────────────────── */}
      <header style={{ marginBottom: 32 }}>
        <div style={{ display: "flex", alignItems: "flex-start", gap: 20 }}>
          {/* Agent image — rendered via next/image-compatible <img> with explicit dimensions */}
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img
            src={agent.imageURI}
            alt={`${agent.name} logo`}
            width={80}
            height={80}
            style={{ borderRadius: 12, objectFit: "cover", flexShrink: 0 }}
          />

          <div style={{ flex: 1, minWidth: 0 }}>
            <div style={{ display: "flex", alignItems: "center", gap: 12, flexWrap: "wrap" }}>
              <h1 style={{ margin: 0, fontSize: 28, fontWeight: 800 }}>{agent.name}</h1>
              <span
                style={{
                  padding: "2px 10px",
                  background: "#f1f5f9",
                  borderRadius: 20,
                  fontSize: 13,
                  fontWeight: 600,
                  color: "#475569",
                }}
              >
                {agent.symbol}
              </span>
              {agent.graduated && (
                <output
                  aria-label='Graduated agent'
                  style={{
                    padding: "2px 10px",
                    background: "#d1fae5",
                    borderRadius: 20,
                    fontSize: 12,
                    fontWeight: 700,
                    color: "#065f46",
                  }}
                >
                  Graduated
                </output>
              )}
            </div>

            <p style={{ margin: "8px 0 12px", color: "#475569", fontSize: 14 }}>
              {agent.description}
            </p>

            <p style={{ margin: 0, fontSize: 13, color: "#6b7280" }}>
              Creator: <code style={{ fontSize: 12 }}>{truncateAddress(agent.creator)}</code>
            </p>
          </div>
        </div>

        {/* Module badges */}
        {agent.modules.length > 0 && (
          <ul
            aria-label='Enabled modules'
            style={{
              display: "flex",
              gap: 8,
              flexWrap: "wrap",
              marginTop: 16,
              listStyle: "none",
              padding: 0,
              margin: "16px 0 0",
            }}
          >
            {agent.modules.map((mod) => (
              <li
                key={mod}
                style={{
                  padding: "3px 10px",
                  background: "#eff6ff",
                  border: "1px solid #bfdbfe",
                  borderRadius: 20,
                  fontSize: 12,
                  color: "#1d4ed8",
                  fontWeight: 500,
                }}
              >
                {MODULE_LABELS[mod] ?? mod}
              </li>
            ))}
          </ul>
        )}
      </header>

      {/* ── Price & market cap ───────────────────────────────────── */}
      <section aria-labelledby='price-heading' style={{ marginBottom: 32 }}>
        <h2 id='price-heading' style={{ fontSize: 16, fontWeight: 700, marginBottom: 12 }}>
          Price &amp; Market Cap
        </h2>
        <dl
          style={{ display: "flex", gap: 32, flexWrap: "wrap", margin: 0 }}
          aria-label='Price and market cap stats'
        >
          <div>
            <dt style={{ fontSize: 12, color: "#6b7280", marginBottom: 2 }}>Current price</dt>
            <dd
              style={{ fontSize: 22, fontWeight: 700, margin: 0 }}
              aria-live='polite'
              aria-label={`Current price: ${formatPrice(agent.currentPrice)} TITU`}
            >
              {formatPrice(agent.currentPrice)}{" "}
              <span style={{ fontSize: 14, fontWeight: 400, color: "#6b7280" }}>TITU</span>
            </dd>
          </div>
          <div>
            <dt style={{ fontSize: 12, color: "#6b7280", marginBottom: 2 }}>Market cap</dt>
            <dd style={{ fontSize: 22, fontWeight: 700, margin: 0 }}>
              {agent.marketCap.toLocaleString()}{" "}
              <span style={{ fontSize: 14, fontWeight: 400, color: "#6b7280" }}>TITU</span>
            </dd>
          </div>
          <div>
            <dt style={{ fontSize: 12, color: "#6b7280", marginBottom: 2 }}>TITU raised</dt>
            <dd style={{ fontSize: 22, fontWeight: 700, margin: 0 }}>
              {agent.realTituReserve.toLocaleString()}{" "}
              <span style={{ fontSize: 14, fontWeight: 400, color: "#6b7280" }}>
                / {agent.graduationThreshold.toLocaleString()}
              </span>
            </dd>
          </div>
        </dl>
      </section>

      {/* ── Graduation progress ──────────────────────────────────── */}
      {!agent.graduated && (
        <section aria-labelledby='grad-heading' style={{ marginBottom: 32 }}>
          <h2 id='grad-heading' style={{ fontSize: 16, fontWeight: 700, marginBottom: 8 }}>
            Progress to graduation
          </h2>
          <div
            role='progressbar'
            aria-valuenow={progressPct}
            aria-valuemin={0}
            aria-valuemax={100}
            aria-label={`${progressPct}% of graduation threshold reached`}
            tabIndex={0}
            style={{
              height: 12,
              background: "#e2e8f0",
              borderRadius: 6,
              overflow: "hidden",
            }}
          >
            <div
              style={{
                height: "100%",
                width: `${progressPct}%`,
                background: progressPct >= 80 ? "#16a34a" : "#2563eb",
                borderRadius: 6,
                transition: "width 0.3s ease",
              }}
            />
          </div>
          <p style={{ margin: "6px 0 0", fontSize: 13, color: "#475569" }}>
            {agent.realTituReserve.toLocaleString()} / {agent.graduationThreshold.toLocaleString()}{" "}
            TITU ({progressPct}%)
          </p>
        </section>
      )}

      {/* ── Bonding curve chart ──────────────────────────────────── */}
      <section aria-labelledby='chart-heading' style={{ marginBottom: 40 }}>
        <h2 id='chart-heading' style={{ fontSize: 16, fontWeight: 700, marginBottom: 12 }}>
          Bonding curve
        </h2>
        <BondingCurveChart agent={agent} />
      </section>

      {/* ── Trade widget ─────────────────────────────────────────── */}
      <section aria-labelledby='trade-heading' style={{ marginBottom: 40 }}>
        <h2 id='trade-heading' style={{ fontSize: 16, fontWeight: 700, marginBottom: 12 }}>
          Trade
        </h2>
        <TradeWidget agent={agent} />
      </section>
    </main>
  );
}
