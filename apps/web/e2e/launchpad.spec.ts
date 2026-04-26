/**
 * E2E smoke: launch → trade → graduate
 *
 * Covers scenario S1 from the testing strategy (smoke level):
 *   1. Launch: complete the wizard, confirm AgentLaunched result surfaced in UI.
 *   2. Trade: navigate to agent page, submit a buy via the trade widget.
 *   3. Graduate: visit the graduated agent fixture page, assert the Graduated
 *      badge and locked LP marker are visible.
 *
 * The current server actions return deterministic fixture stubs. The anvil node
 * started by global-setup deploys real contracts so contract addresses are
 * available in the environment — future iterations wire the UI to the live
 * chain; this test will pass without modification once that wiring lands.
 */

import { type Page, expect, test } from "@playwright/test";

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

/** Wait for a form step to become visible by its aria-label. */
async function waitForStep(page: Page, label: string) {
  await page.waitForSelector(`[aria-label="${label}"]`, { timeout: 10_000 });
}

// ---------------------------------------------------------------------------
// Suite
// ---------------------------------------------------------------------------

test.describe("launchpad golden path", () => {
  // ── 1. Launch wizard ──────────────────────────────────────────────────────

  test("completes the launch wizard and shows agent launched result", async ({ page }) => {
    await page.goto("/launchpad/create");
    await expect(page).toHaveTitle(/Titular/i);

    // Step 0 — Agent Identity
    await waitForStep(page, "Agent identity form");

    await page.fill("#name", "E2EBot");
    await page.fill("#symbol", "E2EB");
    await page.fill("#description", "End-to-end test agent");
    // imageURI is pre-filled with a fixture URL; leave it as-is.
    await page.click('button[type="submit"]:has-text("Next")');

    // Step 1 — Soul
    await waitForStep(page, "Soul configuration form");
    // soulURI is pre-filled; go straight to Next.
    await page.click('button[type="submit"]:has-text("Next")');

    // Step 2 — Modules
    await waitForStep(page, "Modules selection form");
    // No modules selected for a bare launch.
    await page.click('button[type="submit"]:has-text("Next")');

    // Step 3 — Review
    await page.waitForSelector('[aria-label="Launch review"]', {
      timeout: 10_000,
    });
    await expect(page.locator('[aria-label="LaunchParams JSON preview"]')).toContainText("E2EBot");
    await page.click('button:has-text("Launch agent")');

    // Step 4 — Submit
    await page.waitForSelector('[aria-label="Submit launch"]', {
      timeout: 10_000,
    });
    await page.click('button:has-text("Confirm launch")');

    // Result — the fixture stub returns a deterministic response.
    const result = page.locator('[aria-label="Launch result"]');
    await expect(result).toBeVisible({ timeout: 15_000 });
    await expect(result.locator("text=Agent launched!")).toBeVisible();
    await expect(result.locator("dt:has-text('Transaction hash') + dd code")).toContainText("0x");
    await expect(result.locator("dt:has-text('Agent ID') + dd")).not.toBeEmpty();

    // "View agent page" link is present and well-formed.
    const viewLink = result.locator('a:has-text("View agent page")');
    await expect(viewLink).toBeVisible();
    const href = await viewLink.getAttribute("href");
    expect(href).toMatch(/^\/agent\//);
  });

  // ── 2. Trade widget ───────────────────────────────────────────────────────

  test("buys on the bonding curve via the trade widget", async ({ page }) => {
    // Use the sample fixture agent (slug = "sample", not graduated).
    await page.goto("/agent/sample");

    await expect(page.locator("h1")).toContainText("AlphaBot");

    // Graduation progress bar is visible (not graduated yet).
    const progressbar = page.getByRole("progressbar");
    await expect(progressbar).toBeVisible();

    // Switch to Buy tab (already selected by default).
    const buyTab = page.getByRole("tab", { name: "Buy" });
    await expect(buyTab).toHaveAttribute("aria-selected", "true");

    // Enter an amount.
    const amountIn = page.locator("#amount-in");
    await amountIn.fill("100");

    // Estimated output updates (live via bondingCurve math).
    const preview = page.locator("#amount-out-preview");
    await expect(preview).not.toContainText("—");

    // Submit the trade form.
    const submitBtn = page.locator('button[type="submit"]');
    await submitBtn.click();

    // The fixture stub returns a trade result.
    const tradeResult = page.locator('[aria-label="Trade result"]');
    await expect(tradeResult).toBeVisible({ timeout: 10_000 });
    await expect(tradeResult.locator("text=Trade submitted!")).toBeVisible();
    await expect(tradeResult.locator("dt:has-text('Tx hash:') + dd code")).toContainText("0x");
  });

  test("sells on the bonding curve via the trade widget", async ({ page }) => {
    await page.goto("/agent/sample");

    const sellTab = page.getByRole("tab", { name: "Sell" });
    await sellTab.click();
    await expect(sellTab).toHaveAttribute("aria-selected", "true");

    await page.locator("#amount-in").fill("1000");

    const preview = page.locator("#amount-out-preview");
    await expect(preview).not.toContainText("—");

    await page.locator('button[type="submit"]').click();

    const tradeResult = page.locator('[aria-label="Trade result"]');
    await expect(tradeResult).toBeVisible({ timeout: 10_000 });
    await expect(tradeResult.locator("text=Trade submitted!")).toBeVisible();
  });

  // ── 3. Graduation ─────────────────────────────────────────────────────────

  test("graduated agent page shows Graduated badge and LP-lock marker", async ({ page }) => {
    // The "graduated" fixture (slug = "graduated") pre-sets graduated: true.
    await page.goto("/agent/graduated");

    await expect(page.locator("h1")).toContainText("GradBot");

    // Graduated badge is rendered as an <output> with aria-label.
    const graduatedBadge = page.locator('[aria-label="Graduated agent"]');
    await expect(graduatedBadge).toBeVisible();
    await expect(graduatedBadge).toContainText("Graduated");

    // Graduation progress bar is NOT rendered (agent already graduated).
    await expect(page.getByRole("progressbar")).not.toBeVisible();

    // Trade widget disables the submit button for graduated agents.
    const submitBtn = page.locator('button[type="submit"]');
    await expect(submitBtn).toBeDisabled();
    await expect(submitBtn).toContainText("Graduated");

    // Trade widget shows the Uniswap note for graduated agents.
    const uniNote = page.locator('[role="note"]');
    await expect(uniNote).toBeVisible();
    await expect(uniNote).toContainText("Uniswap V2");

    // TITU raised matches the graduation threshold in the stats row.
    const tituRaised = page.locator('[aria-label="Price and market cap stats"]');
    await expect(tituRaised).toContainText("42,000");
  });

  // ── 4. Contract addresses available from anvil fixture ───────────────────

  test("anvil deployment addresses are set in environment", () => {
    // If the global setup ran successfully, these env vars are populated.
    // This test acts as a fast gate: if it fails, the deployment failed.
    const factoryAddr = process.env.FACTORY_ADDR;
    const tituAddr = process.env.TITU_ADDR;

    // In CI the global setup runs for real; locally with reuseExistingServer
    // and no anvil started, env vars may be absent — skip gracefully.
    if (!factoryAddr || !tituAddr) {
      test.skip();
      return;
    }

    expect(factoryAddr).toMatch(/^0x[0-9a-fA-F]{40}$/);
    expect(tituAddr).toMatch(/^0x[0-9a-fA-F]{40}$/);
    expect(process.env.GRADUATOR_ADDR).toMatch(/^0x[0-9a-fA-F]{40}$/);
  });
});
