# gateway-go k6 load test

Throughput / latency assertion for the gateway's read surface
(`/api/v1/*` and `/healthz`).

Tracks issue [#92](https://github.com/0xDevNinja/titular/issues/92):

| SLO target            | Threshold              |
| --------------------- | ---------------------- |
| Sustained throughput  | >= 1,000 RPS           |
| P99 latency           | < 200 ms               |
| Error rate (HTTP 5xx) | < 1%                   |

## Layout

```
loadtest/k6/
├── docker-compose.yml   # postgres + migrate + seed + gateway + k6
├── scenarios.js         # k6 scenarios + SLO thresholds
├── script.js            # k6 entry point; mixes endpoints
├── seed.sql             # fixture rows for agents/trades/jobs
└── README.md
```

## Run locally

The supported entry point is `docker compose`. Everything (Postgres,
schema migration, fixture seeding, gateway build, k6 driver) is
brought up by a single command, runs to completion, and exits with
the k6 exit code.

```bash
# from services/gateway-go/loadtest/k6
docker compose up --abort-on-container-exit --exit-code-from k6
```

The `--abort-on-container-exit` flag tears every other container
down as soon as k6 finishes, so the Postgres container does not
linger after the run.

To customise the run:

```bash
# 30-second smoke run at 250 RPS (useful for editing thresholds)
K6_TARGET_RPS=250 K6_DURATION=30s docker compose up \
    --abort-on-container-exit --exit-code-from k6
```

Variables consumed by the compose file (and ultimately by `script.js`
and `scenarios.js`):

| Variable          | Default | Purpose                                                       |
| ----------------- | ------- | ------------------------------------------------------------- |
| `K6_TARGET_RPS`   | `1000`  | Sustained request rate per second during the hold stage.      |
| `K6_WARMUP`       | `30s`   | Ramp window from 0 RPS to `K6_TARGET_RPS` before the hold.    |
| `K6_DURATION`     | `60s`   | Steady-state hold window (Go duration accepted by k6).        |
| `K6_VUS`          | `200`   | Initial size of the VU pool.                                  |
| `K6_AGENT_MAX_ID` | `100`   | Upper bound for the by-id scenario's PK range.                |

### Warmup convention

Each scenario runs as a `ramping-arrival-rate` executor with two
stages:

1. **Warmup** — ramp `0 -> K6_TARGET_RPS` over `K6_WARMUP` (default
   30s). The first ~5s of any cold run is dominated by `go run`
   compile cost, JIT-warming the request pipeline, and Postgres
   priming its plan cache; ramping in keeps those samples out of the
   steady-state percentiles.
2. **Hold** — sustain `K6_TARGET_RPS` for `K6_DURATION` (default
   60s). This is the window the SLO is asserted over.

Threshold `delayAbortEval` is aligned to `K6_WARMUP` so the abort-on-
fail wiring does not fire during the ramp. The `http_reqs` rate
threshold is set to 80% of `K6_TARGET_RPS` because the global rate is
averaged over the *entire* run (warmup + hold); at 30s+60s the mean
is ~83% of target.

Set `K6_WARMUP=0s` to disable the ramp (useful when driving a
pre-warmed gateway from an external runner; not recommended for the
docker-compose path where the gateway is freshly compiled).

## Run against an existing gateway

If you already have `gateway-go` running locally and just want to
drive it with k6, install k6 (`brew install k6`) and run:

```bash
GATEWAY_BASE_URL=http://localhost:8080 \
K6_TARGET_RPS=250 \
K6_DURATION=30s \
k6 run script.js
```

The seed step is the caller's responsibility in this mode. The
fixture rows in `seed.sql` are the simplest way; pipe them into
your local Postgres after the indexer migrations have been applied.

## Why this is gated behind `workflow_dispatch` + cron

The CI workflow at `.github/workflows/gateway-loadtest.yml` is
NOT triggered on every PR. A 60-second 1k-RPS run plus the
postgres + gateway boot would dominate the per-PR CI minutes for
changes that have nothing to do with the gateway.

Instead it runs:

* On `workflow_dispatch` — manual trigger, useful when reviewing a
  PR that touches the gateway hot path.
* Weekly on a cron (Mondays 04:30 UTC) — catches drift introduced
  by dependency or runtime bumps.

Failures here are SLO regressions, not correctness regressions, so
the right place to surface them is a scheduled run + a manual
trigger reviewers can flip on, not a per-PR red light.

## Reading the output

k6 prints a per-threshold summary at the end. The lines to watch
are:

```
http_req_duration..............: p(99)=… (target <200ms)
http_req_failed................: rate=…  (target <1%)
http_reqs......................: rate=…  (target >=950 / 95% of 1k)
```

Per-endpoint percentiles are tagged via `endpoint:` so a regression
that only affects one route surfaces in its own line:

```
http_req_duration{endpoint:agents_list}: p(99)=…
http_req_duration{endpoint:trades_list}: p(99)=…
http_req_duration{endpoint:jobs_list}:   p(99)=…
http_req_duration{endpoint:agent_by_id}: p(99)=…
http_req_duration{endpoint:stats}:       p(99)=…
http_req_duration{endpoint:healthz}:     p(99)=…
```

`/healthz` carries a tighter sub-budget (`p(99) < 50ms`) than the
DB-bound paths because it does not touch Postgres; folding it into
the global p99 would otherwise mask read-path regressions.

## Updating the SLO

Edit `scenarios.js`. Both the thresholds and the per-scenario rate
share live there. Keep `script.js` free of numeric SLO constants so
the assertion surface is one file.
