// Scenario definitions and SLO thresholds for the k6 gateway loadtest.
//
// Split out from script.js so the test wiring is small enough to read
// at a glance and so the SLO numbers are not buried inside an executor
// config. Edit this file when issue #92's SLO targets change; do not
// edit thresholds inline anywhere else.
//
// Why ramping-arrival-rate (not constant-arrival-rate / constant-vus):
//
//   * The SLO is *throughput* (1k RPS). The arrival-rate family gives
//     k6 a target rate; if the server slows down, k6 will spawn extra
//     VUs up to `maxVUs` to keep the rate steady. With constant-vus a
//     slow server would silently reduce throughput and the SLO check
//     would pass for the wrong reason.
//   * We use the *ramping* variant so each scenario starts at zero RPS
//     and ramps up to its target over `K6_WARMUP`, then holds at
//     target for `K6_DURATION`. Without the ramp the first ~5s of the
//     run is dominated by `go run` compile cost, JIT-warming the
//     request pipeline, and Postgres priming its plan cache; those
//     samples otherwise pollute the steady-state p(99).
//   * `preAllocatedVUs` sizes the worker pool the executor starts
//     with; `maxVUs` is the hard cap. We size both generously so a
//     temporary p99 spike under 200ms can be absorbed without spawning
//     mid-test (which itself causes a latency wobble).
//
// Mix proportions are realised by giving each scenario its own RPS
// share of the global target. They sum to 1.0 of K6_TARGET_RPS.

const TARGET_RPS = Number.parseInt(__ENV.K6_TARGET_RPS || "1000", 10);
const WARMUP = __ENV.K6_WARMUP || "30s";
const DURATION = __ENV.K6_DURATION || "60s";

// Pool sizes. Sized for ~50ms p99 at 1k RPS: Little's law says
// concurrency ~= RPS * latency, so 1000 * 0.05 = 50 in-flight requests
// is the steady-state floor. We allocate 4x for headroom against
// transient slowdowns and to keep TCP connection reuse warm. maxVUs
// is double preAllocatedVUs so a short p99 spike has spawn budget
// before the executor falls behind.
const PRE_ALLOCATED_VUS = Number.parseInt(__ENV.K6_VUS || "200", 10);
const MAX_VUS = PRE_ALLOCATED_VUS * 2;

// Per-scenario rate share. The four scenarios run as one ramping-
// arrival-rate executor each; k6 will sustain each at the configured
// rate independently, and the *sum* is the global RPS the gateway
// sees.
const LIST_RPS = Math.round(TARGET_RPS * 0.6);
const BY_ID_RPS = Math.round(TARGET_RPS * 0.2);
const STATS_RPS = Math.round(TARGET_RPS * 0.1);
const HEALTH_RPS = TARGET_RPS - LIST_RPS - BY_ID_RPS - STATS_RPS; // remainder

// Common executor config. Each scenario picks one of the function
// exports in script.js via `exec`.
//
// `ramping-arrival-rate` walks `startRate` -> stage[i].target over
// stage[i].duration. We use two stages: a ramp from 0 to the
// scenario's RPS over WARMUP, then a hold at that RPS for DURATION.
// SLO assertions kick in only after WARMUP via `delayAbortEval` on
// the thresholds below.
function arrivalRate(rate, exec) {
  return {
    executor: "ramping-arrival-rate",
    startRate: 0,
    timeUnit: "1s",
    preAllocatedVUs: Math.max(1, Math.ceil((PRE_ALLOCATED_VUS * rate) / TARGET_RPS)),
    maxVUs: Math.max(2, Math.ceil((MAX_VUS * rate) / TARGET_RPS)),
    stages: [
      { duration: WARMUP, target: rate },
      { duration: DURATION, target: rate },
    ],
    exec,
    // Tag every request with the scenario name so the run summary
    // breaks per-endpoint percentiles out separately.
    tags: { scenario: exec },
  };
}

export const scenarios = {
  list_endpoints: arrivalRate(LIST_RPS, "listEndpoints"),
  by_id_reads: arrivalRate(BY_ID_RPS, "byIdReads"),
  stats: arrivalRate(STATS_RPS, "stats"),
  healthz: arrivalRate(HEALTH_RPS, "healthz"),
};

// Thresholds enforce the SLO declared in issue #92.
//
// Each entry is paired with `abortOnFail: true` so a regression fails
// the workflow as soon as the threshold is breached, not at end-of-run.
// That keeps CI minutes spent on a doomed run minimal and surfaces the
// failure point in the log before the steady-state numbers are buried
// under post-failure noise.
//
// `delayAbortEval` is aligned to WARMUP so the threshold abort does
// not fire while the executor is still ramping. During the ramp the
// gateway is JIT-warming the request pipeline, Postgres is priming
// its plan cache, and `go run` is finishing its compile; failing on
// those samples is noisy and not what the SLO measures.
export const thresholds = {
  // Global thresholds — the headline SLO numbers from the issue.
  http_req_failed: [{ threshold: "rate<0.01", abortOnFail: true, delayAbortEval: WARMUP }],
  http_req_duration: [{ threshold: "p(99)<200", abortOnFail: true, delayAbortEval: WARMUP }],
  http_reqs: [
    // The ramping-arrival-rate executors hold at TARGET_RPS during
    // the steady-state stage; the global rate is averaged over the
    // whole run (warmup + hold). With a 30s ramp + 60s hold, mean
    // RPS ~= TARGET * (0.5 * 30 + 60) / 90 = TARGET * 0.833. We
    // assert >= 80% of TARGET to leave slack for k6's wall-clock
    // sampling and the last-fraction-of-a-second drain.
    { threshold: `rate>=${Math.round(TARGET_RPS * 0.8)}` },
  ],
  iteration_duration: [{ threshold: "p(95)<500" }],
  // Per-endpoint sub-budgets. /healthz is much cheaper than the DB-
  // bound paths; folding it into the global p99 hides regressions in
  // the read paths under healthz's near-zero-latency volume. Tighter
  // budgets here flag a regression in the actual read surface even if
  // the global p99 still passes.
  "http_req_duration{endpoint:healthz}": ["p(99)<50"],
  "http_req_duration{endpoint:stats}": ["p(99)<200"],
  "http_req_duration{endpoint:agents_list}": ["p(99)<200"],
  "http_req_duration{endpoint:trades_list}": ["p(99)<200"],
  "http_req_duration{endpoint:jobs_list}": ["p(99)<200"],
  "http_req_duration{endpoint:agent_by_id}": ["p(99)<200"],
  // Custom error counter — should stay flat. Catches non-HTTP failures
  // (body-shape mismatches caught by `check`) that http_req_failed
  // doesn't see.
  titular_ok_rate: ["rate>0.99"],
};
