// Shared opt-in / skip predicates for the e2e suite. Centralised so each
// test file gets a single read-once decision and the reasons surface in the
// vitest output identically across suites.
//
// Two layers of gating:
//
//   1. `E2E_INTEGRATION=1` — explicit opt-in. Without it every suite is a
//      no-op. Keeps `pnpm turbo run test` cheap on contributor laptops and
//      on PRs that don't touch the SDK / e2e packages.
//   2. Per-prerequisite probes (anvil-on-PATH, go-on-PATH, docker-up). The
//      probes run inside `beforeAll` of each suite and call `ctx.skip()` on
//      failure so vitest reports skipped-not-failed when the host lacks a
//      tool the suite needs.

/** True iff the e2e suite was explicitly opted into via env var. */
export const E2E_ENABLED: boolean = process.env.E2E_INTEGRATION === "1";

/** Reason string surfaced in `describe.skipIf` when the suite is gated off. */
export const SKIP_REASON =
  "E2E_INTEGRATION is not set to 1; rerun with E2E_INTEGRATION=1 to opt in";
