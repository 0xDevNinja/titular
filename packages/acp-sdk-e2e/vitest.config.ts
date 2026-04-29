import { defineConfig } from "vitest/config";

// E2E suite is opt-in: vitest runs the file but every `describe` is gated on
// `E2E_INTEGRATION=1`. Without the env var the file boots, prints a single
// "skipped" line per suite, and exits 0 so `pnpm turbo run test` stays green
// in environments that lack docker / anvil / a Go toolchain.
//
// Default `pnpm turbo run test` therefore costs ~constant time on this
// package; CI dedicates a separate scheduled job (see
// .github/workflows/sdk-e2e.yml) that flips the env var on.
export default defineConfig({
  test: {
    environment: "node",
    globals: false,
    include: ["src/__tests__/**/*.test.ts"],
    // Per-test timeouts are bumped from vitest's 5s default because spinning
    // anvil + a Postgres container + a Go service builds is dominated by IO,
    // not by the assertions. Local laptops typically come in under 30s for a
    // warm cache and 90s cold — 120s gives slack on the GitHub runner without
    // hiding genuine hangs.
    testTimeout: 120_000,
    hookTimeout: 180_000,
  },
});
