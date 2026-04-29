// Scenario definitions and SLO thresholds for the k6 gateway loadtest.
//
// Split out from script.js so the test wiring is small enough to read
// at a glance and so the SLO numbers are not buried inside an executor
// config. Edit this file when issue #92's SLO targets change; do not
// edit thresholds inline anywhere else.
//
// Why constant-arrival-rate (not constant-vus / ramping-vus):
//
//   * The SLO is *throughput* (1k RPS). constant-arrival-rate gives k6
//     a target rate; if the server slows down, k6 will spawn extra VUs
//     up to `maxVUs` to keep the rate steady. With constant-vus a slow
//     server would silently reduce throughput and the SLO check would
//     pass for the wrong reason.
//   * `preAllocatedVUs` sizes the worker pool the executor starts
//     with; `maxVUs` is the hard cap. We size both generously so a
//     temporary p99 spike under 200ms can be absorbed without spawning
//     mid-test (which itself causes a latency wobble).
//
// Mix proportions are realised by giving each scenario its own RPS
// share of the global target. They sum to 1.0 of K6_TARGET_RPS.

const TARGET_RPS = Number.parseInt(__ENV.K6_TARGET_RPS || "1000", 10);
const DURATION = __ENV.K6_DURATION || "60s";

// Pool sizes. Sized for ~50ms p99 at 1k RPS: Little's law says
// concurrency ~= RPS * latency, so 1000 * 0.05 = 50 in-flight requests
// is the steady-state floor. We allocate 4x for headroom against
// transient slowdowns and to keep TCP connection reuse warm. maxVUs
// is double preAllocatedVUs so a short p99 spike has spawn budget
// before the executor falls behind.
const PRE_ALLOCATED_VUS = Number.parseInt(__ENV.K6_VUS || "200", 10);
const MAX_VUS = PRE_ALLOCATED_VUS * 2;

// Per-scenario rate share. The four scenarios run as one constant-
// arrival-rate executor each; k6 will sustain each at the configured
// rate independently, and the *sum* is the global RPS the gateway
// sees.
const LIST_RPS = Math.round(TARGET_RPS * 0.6);
const BY_ID_RPS = Math.round(TARGET_RPS * 0.2);
const STATS_RPS = Math.round(TARGET_RPS * 0.1);
const HEALTH_RPS = TARGET_RPS - LIST_RPS - BY_ID_RPS - STATS_RPS; // remainder

// Common executor config. Each scenario picks one of the function
// exports in script.js via `exec`.
function arrivalRate(rate, exec) {
  return {
    executor: "constant-arrival-rate",
    rate,
    timeUnit: "1s",
    duration: DURATION,
    preAllocatedVUs: Math.max(1, Math.ceil((PRE_ALLOCATED_VUS * rate) / TARGET_RPS)),
    maxVUs: Math.max(2, Math.ceil((MAX_VUS * rate) / TARGET_RPS)),
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
// `delayAbortEval` gives the test 10 seconds to warm up before
// thresholds start tripping. The first second of any constant-arrival-
// rate run sees k6 spinning up its VU pool, JIT-warming the request
// pipeline, and Postgres priming its plan cache; failing on those is
// noisy.
export const thresholds = {
  // Global thresholds — the headline SLO numbers from the issue.
  http_req_failed: [{ threshold: "rate<0.01", abortOnFail: true, delayAbortEval: "10s" }],
  http_req_duration: [{ threshold: "p(99)<200", abortOnFail: true, delayAbortEval: "10s" }],
  http_reqs: [
    // The constant-arrival-rate executors are configured to sum to
    // TARGET_RPS; we assert the gateway *received* at least 95% of
    // that. Five percent slack covers the 1-second warmup and the
    // last-fraction-of-a-second drain — k6's sampling is wall-clock
    // and the rate executor cannot land its last iteration exactly
    // on the duration boundary.
    { threshold: `rate>=${Math.round(TARGET_RPS * 0.95)}` },
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
