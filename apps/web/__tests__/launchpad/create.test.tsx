import { cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it, vi } from "vitest";
import LaunchpadCreatePage from "../../app/launchpad/create/page";
import { buildLaunchParams, computeBitmap } from "../../app/launchpad/create/types";

// ─── Mock server action ───────────────────────────────────────────────────────
vi.mock("../../app/launchpad/create/actions", () => ({
  submitLaunch: vi.fn().mockResolvedValue({
    txHash: "0xdeadbeef",
    agentId: "test-agent-001",
    agentSlug: "test-agent",
    tokenAddress: "0xAAAA",
    curveAddress: "0xBBBB",
    gasUsed: 350000,
    blockNumber: 99,
  }),
}));

afterEach(() => cleanup());

// ─── Unit: computeBitmap ──────────────────────────────────────────────────────
describe("computeBitmap", () => {
  it("returns 0 when no modules enabled", () => {
    expect(computeBitmap([])).toBe(0);
  });

  it("sets bit 0 for AntiSniper", () => {
    expect(computeBitmap(["AntiSniper"])).toBe(0b0000001);
  });

  it("sets bit 1 for SixtyDays", () => {
    expect(computeBitmap(["SixtyDays"])).toBe(0b0000010);
  });

  it("combines multiple modules correctly", () => {
    // AntiSniper=bit0, SixtyDays=bit1, CapitalFormation=bit2
    expect(computeBitmap(["AntiSniper", "SixtyDays", "CapitalFormation"])).toBe(0b0000111);
  });

  it("sets all 7 module bits", () => {
    const all = [
      "AntiSniper",
      "SixtyDays",
      "CapitalFormation",
      "Airdrop",
      "LaunchRadar",
      "PreBuy",
      "ExistingToken",
    ];
    expect(computeBitmap(all)).toBe(0b1111111);
  });

  it("ignores unknown module ids", () => {
    expect(computeBitmap(["UnknownModule"])).toBe(0);
  });
});

// ─── Unit: buildLaunchParams ──────────────────────────────────────────────────
describe("buildLaunchParams", () => {
  it("builds params with empty modules", () => {
    const params = buildLaunchParams({
      identity: {
        name: "Alpha",
        symbol: "ALP",
        description: "Test",
        imageURI: "https://a.com/img.png",
      },
      soul: { soulURI: "https://a.com/soul.json" },
      modules: { enabled: [], configs: {} },
    });
    expect(params.name).toBe("Alpha");
    expect(params.symbol).toBe("ALP");
    expect(params.moduleFlags).toBe(0);
    expect(params.moduleConfigs).toHaveLength(0);
  });

  it("encodes enabled module configs as JSON strings", () => {
    const params = buildLaunchParams({
      identity: { name: "Beta", symbol: "BET", description: "d", imageURI: "https://b.com/i.png" },
      soul: { soulURI: "https://b.com/s.json" },
      modules: {
        enabled: ["AntiSniper"],
        configs: { AntiSniper: { decayBlocks: 200 } },
      },
    });
    expect(params.moduleFlags).toBe(1);
    expect(params.moduleConfigs).toHaveLength(1);
    expect(JSON.parse(params.moduleConfigs[0])).toEqual({ decayBlocks: 200 });
  });
});

// ─── Integration: wizard renders and step nav works ───────────────────────────
describe("LaunchpadCreatePage", () => {
  it("renders step 1 on mount", () => {
    render(<LaunchpadCreatePage />);
    expect(screen.getByRole("heading", { name: /launch your agent/i })).toBeInTheDocument();
    expect(screen.getByLabelText(/agent identity form/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/agent name/i)).toBeInTheDocument();
  });

  it("shows validation error for empty name", async () => {
    const user = userEvent.setup();
    render(<LaunchpadCreatePage />);
    await user.click(screen.getByRole("button", { name: /next/i }));
    await waitFor(() => {
      expect(screen.getByText(/name is required/i)).toBeInTheDocument();
    });
  });

  it("shows validation error for invalid ERC-20 symbol", async () => {
    const user = userEvent.setup();
    render(<LaunchpadCreatePage />);

    await user.type(screen.getByLabelText(/agent name/i), "My Agent");
    await user.clear(screen.getByLabelText(/token symbol/i));
    await user.type(screen.getByLabelText(/token symbol/i), "lowercase!");
    await user.click(screen.getByRole("button", { name: /next/i }));

    await waitFor(() => {
      expect(screen.getByText(/uppercase letters\/digits/i)).toBeInTheDocument();
    });
  });

  it("navigates from step 1 to step 2 with valid identity", async () => {
    const user = userEvent.setup();
    render(<LaunchpadCreatePage />);

    await user.type(screen.getByLabelText(/agent name/i), "Alpha Bot");
    await user.clear(screen.getByLabelText(/token symbol/i));
    await user.type(screen.getByLabelText(/token symbol/i), "ALPHA");
    await user.type(screen.getByLabelText(/description/i), "A test agent description");
    // imageURI is pre-filled with fixture URL — leave it

    await user.click(screen.getByRole("button", { name: /next/i }));

    await waitFor(() => {
      expect(screen.getByLabelText(/soul configuration form/i)).toBeInTheDocument();
    });
  });

  it("navigates back from step 2 to step 1", async () => {
    const user = userEvent.setup();
    render(<LaunchpadCreatePage />);

    // advance to step 2
    await user.type(screen.getByLabelText(/agent name/i), "Alpha Bot");
    await user.clear(screen.getByLabelText(/token symbol/i));
    await user.type(screen.getByLabelText(/token symbol/i), "ALPHA");
    await user.type(screen.getByLabelText(/description/i), "A test agent description");
    await user.click(screen.getByRole("button", { name: /next/i }));
    await waitFor(() => {
      expect(screen.getByLabelText(/soul configuration form/i)).toBeInTheDocument();
    });

    // go back
    await user.click(screen.getByRole("button", { name: /back/i }));
    await waitFor(() => {
      expect(screen.getByLabelText(/agent identity form/i)).toBeInTheDocument();
    });
  });

  it("reaches modules step and shows module toggles", async () => {
    const user = userEvent.setup();
    render(<LaunchpadCreatePage />);

    // step 1
    await user.type(screen.getByLabelText(/agent name/i), "Alpha Bot");
    await user.clear(screen.getByLabelText(/token symbol/i));
    await user.type(screen.getByLabelText(/token symbol/i), "ALPHA");
    await user.type(screen.getByLabelText(/description/i), "A test agent");
    await user.click(screen.getByRole("button", { name: /next/i }));
    await waitFor(() => screen.getByLabelText(/soul configuration form/i));

    // step 2
    await user.click(screen.getByRole("button", { name: /next/i }));
    await waitFor(() => screen.getByLabelText(/modules selection form/i));

    // all 7 module checkboxes present
    expect(screen.getByLabelText(/enable anti-sniper module/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/enable 60-day escrow module/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/enable capital formation module/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/enable airdrop module/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/enable launch radar module/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/enable pre-buy module/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/enable existing token module/i)).toBeInTheDocument();
  });

  it("module bitmap updates when a module is toggled", async () => {
    const user = userEvent.setup();
    render(<LaunchpadCreatePage />);

    // advance to modules step
    await user.type(screen.getByLabelText(/agent name/i), "Bot");
    await user.clear(screen.getByLabelText(/token symbol/i));
    await user.type(screen.getByLabelText(/token symbol/i), "BOT");
    await user.type(screen.getByLabelText(/description/i), "desc");
    await user.click(screen.getByRole("button", { name: /next/i }));
    await waitFor(() => screen.getByLabelText(/soul configuration form/i));
    await user.click(screen.getByRole("button", { name: /next/i }));
    await waitFor(() => screen.getByLabelText(/modules selection form/i));

    // bitmap starts at 0
    expect(screen.getByRole("status")).toHaveTextContent("0b0000000 (0)");

    // enable AntiSniper (bit 0 → value 1)
    await user.click(screen.getByLabelText(/enable anti-sniper module/i));
    await waitFor(() => {
      expect(screen.getByRole("status")).toHaveTextContent("0b0000001 (1)");
    });

    // enable SixtyDays (bit 1 → total 3)
    await user.click(screen.getByLabelText(/enable 60-day escrow module/i));
    await waitFor(() => {
      expect(screen.getByRole("status")).toHaveTextContent("0b0000011 (3)");
    });
  });

  it("reaches review step and shows LaunchParams JSON", async () => {
    const user = userEvent.setup();
    render(<LaunchpadCreatePage />);

    // step 1
    await user.type(screen.getByLabelText(/agent name/i), "ReviewBot");
    await user.clear(screen.getByLabelText(/token symbol/i));
    await user.type(screen.getByLabelText(/token symbol/i), "RBOT");
    await user.type(screen.getByLabelText(/description/i), "Review desc");
    await user.click(screen.getByRole("button", { name: /next/i }));
    await waitFor(() => screen.getByLabelText(/soul configuration form/i));

    // step 2
    await user.click(screen.getByRole("button", { name: /next/i }));
    await waitFor(() => screen.getByLabelText(/modules selection form/i));

    // step 3
    await user.click(screen.getByRole("button", { name: /next/i }));
    await waitFor(() => screen.getByLabelText(/launch review/i));

    const preview = screen.getByLabelText(/launchparams json preview/i);
    expect(preview).toBeInTheDocument();
    const json = JSON.parse(preview.textContent ?? "{}");
    expect(json.name).toBe("ReviewBot");
    expect(json.symbol).toBe("RBOT");
    expect(typeof json.moduleFlags).toBe("number");
  });

  it("submit step calls action and shows mock tx hash", async () => {
    const user = userEvent.setup();
    render(<LaunchpadCreatePage />);

    // step 1
    await user.type(screen.getByLabelText(/agent name/i), "SubmitBot");
    await user.clear(screen.getByLabelText(/token symbol/i));
    await user.type(screen.getByLabelText(/token symbol/i), "SBOT");
    await user.type(screen.getByLabelText(/description/i), "Submit desc");
    await user.click(screen.getByRole("button", { name: /next/i }));
    await waitFor(() => screen.getByLabelText(/soul configuration form/i));

    // step 2
    await user.click(screen.getByRole("button", { name: /next/i }));
    await waitFor(() => screen.getByLabelText(/modules selection form/i));

    // step 3
    await user.click(screen.getByRole("button", { name: /next/i }));
    await waitFor(() => screen.getByLabelText(/launch review/i));

    // step 4 via review button
    await user.click(screen.getByRole("button", { name: /launch agent/i }));
    await waitFor(() => screen.getByLabelText(/submit launch/i));

    // confirm launch
    await user.click(screen.getByRole("button", { name: /confirm launch/i }));
    await waitFor(() => {
      expect(screen.getByRole("status")).toHaveTextContent("0xdeadbeef");
    });
    expect(screen.getByText(/view agent page/i)).toBeInTheDocument();
  });

  it("step indicator marks current step with aria-current=step", () => {
    render(<LaunchpadCreatePage />);
    const stepItems = screen.getAllByRole("listitem");
    expect(stepItems[0]).toHaveAttribute("aria-current", "step");
    expect(stepItems[1]).not.toHaveAttribute("aria-current");
  });

  it("progress nav has accessible label", () => {
    render(<LaunchpadCreatePage />);
    expect(screen.getByRole("navigation", { name: /wizard progress/i })).toBeInTheDocument();
  });
});

// ─── Regression: fireEvent path for bitmap (faster alternative) ──────────────
describe("module bitmap via fireEvent", () => {
  it("enabling ExistingToken sets bit 6", async () => {
    const user = userEvent.setup();
    render(<LaunchpadCreatePage />);

    await user.type(screen.getByLabelText(/agent name/i), "X");
    await user.clear(screen.getByLabelText(/token symbol/i));
    await user.type(screen.getByLabelText(/token symbol/i), "X");
    await user.type(screen.getByLabelText(/description/i), "x");
    await user.click(screen.getByRole("button", { name: /next/i }));
    await waitFor(() => screen.getByLabelText(/soul configuration form/i));
    await user.click(screen.getByRole("button", { name: /next/i }));
    await waitFor(() => screen.getByLabelText(/modules selection form/i));

    fireEvent.click(screen.getByLabelText(/enable existing token module/i));
    await waitFor(() => {
      // ExistingToken is bit 6 → 64
      expect(screen.getByRole("status")).toHaveTextContent("0b1000000 (64)");
    });
  });
});
