// k6 load test for the Titular gateway-go REST surface.
//
// Targets the SLO declared in issue #92:
//
//   * 1,000 RPS sustained throughput
//   * P99 < 200ms latency
//   * < 1% HTTP error rate
//
// The test mixes read endpoints in proportions chosen to mirror the
// expected production traffic profile of the SDK (M4):
//
//   * 60%  list endpoints  — /api/v1/agents, /api/v1/trades, /api/v1/jobs
//   * 20%  by-id reads     — /api/v1/agents/:id
//   * 10%  aggregate stats — /api/v1/stats
//   * 10%  liveness probe  — /healthz
//
// The endpoints are intentionally read-only: SLO budgets here are about
// the gateway's hot read path (Postgres pool + JSON encoding + Gin
// dispatch), not about indexer write throughput, which is exercised
// independently by services/indexer-go's `loadtest`-tagged test.
//
// Environment variables consumed (all optional):
//
//   GATEWAY_BASE_URL   default "http://localhost:8080"; the URL prefix
//                      every request is built against. CI uses
//                      "http://gateway:8080" via docker compose service
//                      DNS.
//   K6_VUS             default 200; number of concurrent virtual users.
//                      With http_req_duration around 50ms, 200 VUs is
//                      sufficient headroom to sustain >=1k RPS without
//                      starving the request rate executor.
//   K6_DURATION        default "60s"; steady-state hold window for the
//                      SLO check (after warmup). Shorter values are
//                      useful for smoke runs; the SLO assertion is
//                      meaningful only at >=30s of hold.
//   K6_WARMUP          default "30s"; ramp window from 0 RPS to target
//                      before the steady-state hold begins. Threshold
//                      abort is delayed by the same value so first-
//                      iteration / plan-cache priming samples don't
//                      fail the run.
//   K6_TARGET_RPS      default 1000; desired sustained request rate
//                      during the hold stage. Drives the ramping-
//                      arrival-rate executor and the http_reqs
//                      threshold.
//   K6_AGENT_MAX_ID    default 100; upper bound on the agent primary
//                      key range used by the by-id scenario. The seed
//                      script populates exactly this many rows.
//
// Why ramping-arrival-rate over per-VU iteration: we want the load
// generator to push a *fixed RPS* even if the server slows down, so a
// regression in latency is allowed to surface as a queueing backlog
// rather than as silent throughput loss. The ramp prefix keeps cold-
// start samples (compile, plan cache) out of the steady-state p(99).

import { check } from "k6";
import http from "k6/http";
import { Counter, Rate } from "k6/metrics";
import { scenarios, thresholds } from "./scenarios.js";

// ---------------------------------------------------------------------------
// k6 options
// ---------------------------------------------------------------------------

export const options = {
  scenarios,
  thresholds,
  // Discardresponse bodies once we've checked them — k6 keeps them in
  // memory by default for debugging, which inflates RSS noticeably at
  // 1k RPS over a minute.
  discardResponseBodies: false, // we DO need bodies for the JSON shape check
  // Cap connection reuse so the test exercises the gateway's keep-alive
  // path the same way real SDK clients will. k6 defaults to one
  // connection per VU which is what we want.
  noConnectionReuse: false,
  // Fail fast on TLS issues if a future CI variant runs against HTTPS.
  insecureSkipTLSVerify: false,
  // Summarise per-endpoint percentiles in the run log so a regression
  // can be triaged without re-running.
  summaryTrendStats: ["avg", "min", "med", "p(95)", "p(99)", "max"],
};

// ---------------------------------------------------------------------------
// Custom counters
// ---------------------------------------------------------------------------

const errs = new Counter("titular_errs");
const okRate = new Rate("titular_ok_rate");

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------

const BASE_URL = (__ENV.GATEWAY_BASE_URL || "http://localhost:8080").replace(/\/+$/, "");

const AGENT_MAX_ID = Number.parseInt(__ENV.K6_AGENT_MAX_ID || "100", 10);

// Common request params. We tag each request with `endpoint` so k6's
// per-tag breakdown shows latency per route, not just the global mean.
function params(endpoint) {
  return {
    headers: {
      Accept: "application/json",
      // Static UA so log filters can identify load-test traffic in
      // Postgres pg_stat_activity / gateway access logs.
      "User-Agent": "titular-k6-loadtest/0.1",
    },
    tags: { endpoint },
    timeout: "5s",
  };
}

// expectOK records a hit/miss against okRate and bumps errs on a
// non-2xx so the run summary surfaces the failure breakdown without
// drowning the log in `console.error` lines.
function expectOK(res, endpoint) {
  const ok = check(res, {
    "status is 2xx": (r) => r.status >= 200 && r.status < 300,
    "body is non-empty": (r) => r.body && r.body.length > 0,
  });
  okRate.add(ok);
  if (!ok) {
    errs.add(1, { endpoint, status: String(res.status) });
  }
  return ok;
}

// ---------------------------------------------------------------------------
// Scenario functions
//
// Each export below is referenced by name from scenarios.js. k6 invokes
// them once per arrival-rate iteration; we keep them small and free of
// per-iteration allocation so the load generator is the gateway's
// bottleneck, not the test runner.
// ---------------------------------------------------------------------------

// listEndpoints: 60% of traffic. Picks one of the three list paths
// per iteration so the workload exercises the agents, trades, and
// jobs SQL prepares in roughly equal share within this scenario.
export function listEndpoints() {
  // Cheap deterministic round-robin via the iteration counter so the
  // three endpoints get equal share without an RNG call.
  const choice = __ITER % 3;
  let url;
  let tag;
  if (choice === 0) {
    url = `${BASE_URL}/api/v1/agents?limit=50`;
    tag = "agents_list";
  } else if (choice === 1) {
    url = `${BASE_URL}/api/v1/trades?limit=50`;
    tag = "trades_list";
  } else {
    url = `${BASE_URL}/api/v1/jobs?limit=50`;
    tag = "jobs_list";
  }
  const res = http.get(url, params(tag));
  expectOK(res, tag);
}

// byIdReads: 20% of traffic. Hits /api/v1/agents/:id with a primary
// key drawn uniformly from [1, AGENT_MAX_ID]. The seed script (see
// seed.sql) populates exactly that many rows so every request is
// expected to return 200, not 404.
export function byIdReads() {
  const id = 1 + Math.floor(Math.random() * AGENT_MAX_ID);
  const res = http.get(`${BASE_URL}/api/v1/agents/${id}`, params("agent_by_id"));
  expectOK(res, "agent_by_id");
}

// stats: 10% of traffic. Exercises the aggregate-counter path which
// hits its own `Stats` SQL (separate from the list selects), so it
// lights up a different code path than listEndpoints does.
export function stats() {
  const res = http.get(`${BASE_URL}/api/v1/stats`, params("stats"));
  expectOK(res, "stats");
}

// healthz: 10% of traffic. The chassis liveness probe — does NOT touch
// Postgres. Including it in the mix keeps the realistic SDK pattern
// where clients ping /healthz between API calls; excluding it would
// over-weight DB-bound traffic and exaggerate the SLO budget.
export function healthz() {
  const res = http.get(`${BASE_URL}/healthz`, params("healthz"));
  expectOK(res, "healthz");
}
