import { cleanup, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it, vi } from "vitest";

// ─── Mock next/navigation ─────────────────────────────────────────────────────
// vi.mock is hoisted to the top of the file, so mockNotFound must be declared
// via vi.hoisted to be available inside the factory.
const { mockNotFound } = vi.hoisted(() => {
  const mockNotFound = vi.fn(() => {
    throw new Error("NEXT_NOT_FOUND");
  });
  return { mockNotFound };
});
vi.mock("next/navigation", () => ({ notFound: mockNotFound }));

import graduatedFixture from "../../fixtures/agents/graduated.json";
// ─── Mock fs so the server component works in jsdom ──────────────────────────
import sampleFixture from "../../fixtures/agents/sample.json";

vi.mock("node:fs", () => ({
  default: {
    existsSync: vi.fn((p: string) => {
      if (p.includes("sample.json")) return true;
      if (p.includes("graduated.json")) return true;
      return false;
    }),
    readFileSync: vi.fn((p: string) => {
      if (p.includes("sample.json")) return JSON.stringify(sampleFixture);
      if (p.includes("graduated.json")) return JSON.stringify(graduatedFixture);
      throw new Error("File not found");
    }),
  },
}));

// ─── Mock server action ───────────────────────────────────────────────────────
vi.mock("../../app/agent/[slug]/actions", () => ({
  submitTrade: vi.fn().mockResolvedValue({
    txHash: "0xmocktxhash000000000000000000000000000000000000000000000000001",
    amountIn: "100",
    amountOut: "2186521",
    isBuy: true,
    gasUsed: 142000,
    blockNumber: 18234567,
  }),
}));

import TradeWidget from "../../app/agent/[slug]/TradeWidget";
import AgentPage from "../../app/agent/[slug]/page";

afterEach(() => cleanup());

// ─── AgentPage: renders from fixture ─────────────────────────────────────────

describe("AgentPage — sample fixture", () => {
  async function renderSample() {
    const jsx = await AgentPage({ params: Promise.resolve({ slug: "sample" }) });
    render(jsx as React.ReactElement);
  }

  it("renders agent name and symbol", async () => {
    await renderSample();
    expect(screen.getByRole("heading", { name: /alphabot/i })).toBeInTheDocument();
    expect(screen.getByText("ALPHA")).toBeInTheDocument();
  });

  it("renders description", async () => {
    await renderSample();
    expect(screen.getByText(/defi research agent/i)).toBeInTheDocument();
  });

  it("renders creator address (truncated)", async () => {
    await renderSample();
    expect(screen.getByText(/0x1234/i)).toBeInTheDocument();
  });

  it("renders module badges", async () => {
    await renderSample();
    expect(screen.getByText("Anti-Sniper")).toBeInTheDocument();
    expect(screen.getByText("60-Day Escrow")).toBeInTheDocument();
  });

  it("renders current price dd element", async () => {
    await renderSample();
    // The <dd> element has aria-label="Current price: … TITU"
    const priceDds = screen.getAllByLabelText(/current price/i);
    expect(priceDds.length).toBeGreaterThan(0);
  });

  it("renders graduation progress bar", async () => {
    await renderSample();
    const bar = screen.getByRole("progressbar");
    expect(bar).toBeInTheDocument();
    // 5000/42000 = 11.904…% → Math.round = 12
    expect(bar).toHaveAttribute("aria-valuenow", "12");
  });

  it("renders the bonding curve chart figure", async () => {
    await renderSample();
    expect(screen.getByRole("img", { name: /bonding curve/i })).toBeInTheDocument();
  });

  it("renders the trade widget", async () => {
    await renderSample();
    expect(screen.getByRole("region", { name: /trade widget/i })).toBeInTheDocument();
  });
});

// ─── AgentPage: graduated fixture ────────────────────────────────────────────

describe("AgentPage — graduated fixture", () => {
  async function renderGraduated() {
    const jsx = await AgentPage({ params: Promise.resolve({ slug: "graduated" }) });
    render(jsx as React.ReactElement);
  }

  it("shows graduated badge", async () => {
    await renderGraduated();
    expect(screen.getByRole("status", { name: /graduated/i })).toBeInTheDocument();
  });

  it("does not render graduation progress bar", async () => {
    await renderGraduated();
    expect(screen.queryByRole("progressbar")).not.toBeInTheDocument();
  });
});

// ─── AgentPage: 404 on unknown slug ──────────────────────────────────────────

describe("AgentPage — unknown slug", () => {
  it("calls notFound() for an unknown slug", async () => {
    await expect(
      AgentPage({ params: Promise.resolve({ slug: "does-not-exist" }) })
    ).rejects.toThrow("NEXT_NOT_FOUND");
    expect(mockNotFound).toHaveBeenCalled();
  });
});

// ─── TradeWidget: calculation + submit flow ───────────────────────────────────

describe("TradeWidget", () => {
  const baseProps = { agent: sampleFixture as Parameters<typeof TradeWidget>[0]["agent"] };

  it("renders buy tab selected by default", () => {
    render(<TradeWidget {...baseProps} />);
    const buyTab = screen.getByRole("tab", { name: /buy/i });
    expect(buyTab).toHaveAttribute("aria-selected", "true");
  });

  it("renders sell tab", () => {
    render(<TradeWidget {...baseProps} />);
    expect(screen.getByRole("tab", { name: /sell/i })).toBeInTheDocument();
  });

  it("switches to sell tab on click", async () => {
    const user = userEvent.setup();
    render(<TradeWidget {...baseProps} />);
    await user.click(screen.getByRole("tab", { name: /sell/i }));
    const sellTab = screen.getByRole("tab", { name: /sell/i });
    expect(sellTab).toHaveAttribute("aria-selected", "true");
  });

  it("shows estimated output after typing amount", async () => {
    const user = userEvent.setup();
    render(<TradeWidget {...baseProps} />);
    await user.type(screen.getByLabelText(/titu to spend/i), "100");
    // estimated output should appear (non-zero)
    const output = screen.getByRole("status", { name: /estimated.*received/i });
    expect(output).toBeInTheDocument();
    expect(output.textContent).not.toContain("—");
  });

  it("shows validation error when submitting with no amount", async () => {
    const user = userEvent.setup();
    render(<TradeWidget {...baseProps} />);
    await user.click(screen.getByRole("button", { name: /^buy$/i }));
    await waitFor(() => {
      expect(screen.getByRole("alert")).toBeInTheDocument();
    });
    expect(screen.getByRole("alert")).toHaveTextContent(/greater than 0/i);
  });

  it("calls submitTrade and shows tx hash on success", async () => {
    const user = userEvent.setup();
    render(<TradeWidget {...baseProps} />);
    await user.type(screen.getByLabelText(/titu to spend/i), "100");
    await user.click(screen.getByRole("button", { name: /^buy$/i }));
    await waitFor(() => {
      expect(screen.getByRole("status", { name: /trade result/i })).toBeInTheDocument();
    });
    expect(screen.getByRole("status", { name: /trade result/i })).toHaveTextContent(
      /0xmocktxhash/i
    );
  });

  it("sell tab shows agent symbol as input label", async () => {
    const user = userEvent.setup();
    render(<TradeWidget {...baseProps} />);
    await user.click(screen.getByRole("tab", { name: /sell/i }));
    expect(screen.getByLabelText(/alpha to sell/i)).toBeInTheDocument();
  });

  it("graduated agent disables the submit button", () => {
    const graduatedAgent = {
      ...sampleFixture,
      graduated: true,
    } as Parameters<typeof TradeWidget>[0]["agent"];
    render(<TradeWidget agent={graduatedAgent} />);
    const btn = screen.getByRole("button", { name: /graduated/i });
    expect(btn).toBeDisabled();
  });

  it("estimated output uses constant-product math", async () => {
    const user = userEvent.setup();
    render(<TradeWidget {...baseProps} />);
    await user.type(screen.getByLabelText(/titu to spend/i), "100");

    // Manually compute expected output
    const { virtualTituReserve, virtualAgentReserve, realTituReserve, realAgentReserve } =
      sampleFixture;
    const effectiveIn = 100 * 0.99;
    const xReserve = virtualTituReserve + realTituReserve;
    const yReserve = virtualAgentReserve + realAgentReserve;
    const k = xReserve * yReserve;
    const newX = xReserve + effectiveIn;
    const expected = yReserve - k / newX;

    const output = screen.getByRole("status", { name: /estimated.*received/i });
    // Strip locale formatting (commas) then compare first 4 digits of integer part.
    const rawText = (output.textContent ?? "").replace(/,/g, "");
    const firstFour = expected.toFixed(0).slice(0, 4);
    expect(rawText).toContain(firstFour);
  });
});
