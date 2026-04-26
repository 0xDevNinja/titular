import { defineConfig, devices } from "@playwright/test";

/**
 * Playwright E2E configuration.
 *
 * - Boots a Next.js dev server on port 3100 (avoids clash with default 3000).
 * - globalSetup spawns anvil + deploys contracts before any worker starts.
 * - globalTeardown stops anvil after all workers finish.
 * - Tests run sequentially (workers: 1) to avoid anvil port conflicts.
 */
export default defineConfig({
  testDir: ".",
  testMatch: "**/*.spec.ts",
  timeout: 60_000,
  retries: process.env.CI ? 1 : 0,
  workers: 1,
  reporter: process.env.CI ? "github" : "list",

  globalSetup: "./global-setup.ts",
  globalTeardown: "./global-teardown.ts",

  use: {
    baseURL: "http://localhost:3100",
    trace: "on-first-retry",
    screenshot: "only-on-failure",
  },

  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
  ],

  webServer: {
    command: "pnpm dev --port 3100",
    url: "http://localhost:3100",
    reuseExistingServer: !process.env.CI,
    timeout: 60_000,
  },
});
