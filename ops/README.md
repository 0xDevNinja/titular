# ops/ — Titular operations stack

Operator-facing tooling for running the Titular gateway and indexer
locally or in shared infrastructure. Today this directory hosts the
M4 observability stack (#94); future milestones add deployment
manifests and operator runbooks here.

## Layout

| Path                                      | Owner          | Purpose                                                            |
|-------------------------------------------|----------------|--------------------------------------------------------------------|
| `prometheus/prometheus.yml`               | M4 #94         | Prometheus scrape config for the gateway + indexer `/metrics`.     |
| `grafana/provisioning/datasources/`       | M4 #94         | Auto-installs the Prometheus datasource at first Grafana boot.     |
| `grafana/provisioning/dashboards/`        | M4 #94         | Tells Grafana to load every JSON file in `grafana/dashboards/`.    |
| `grafana/dashboards/gateway.json`         | M4 #94         | Gateway HTTP / SSE / GraphQL dashboards.                           |
| `grafana/dashboards/indexer.json`         | M4 #94         | Indexer head, reorg, and event-throughput dashboards.              |
| `docker-compose.observability.yml`        | M4 #94         | One-liner Prometheus + Grafana stack for local development.        |

## Quick start

1. Run the gateway and / or indexer with the metrics surface enabled:

   ```bash
   GATEWAY_METRICS_ADDR=:9090 go run ./services/gateway-go/cmd/gateway
   ```

   The indexer follows the same pattern via `INDEXER_METRICS_ADDR=:9091`
   when its daemon entry point lands.

2. Bring up Prometheus + Grafana:

   ```bash
   docker compose -f ops/docker-compose.observability.yml up
   ```

3. Open Grafana at `http://localhost:3000`. Default credentials are
   `admin / admin`. The Titular dashboards auto-provision under the
   `Titular` folder.

## Metric naming

Both services emit metrics through OpenTelemetry's Prometheus exporter
(`go.opentelemetry.io/otel/exporters/prometheus`). The exporter
prefixes the OTel scope onto the `otel_scope_name` label so dashboards
can filter `{otel_scope_name="github.com/0xDevNinja/titular/services/gateway-go"}`
to isolate gateway samples from indexer samples that share an
instrument name (e.g. `http_server_request_duration_seconds`).

Custom instruments owned by the gateway:

- `gateway_sse_active_connections` — current SSE subscriber count.
- `gateway_graphql_active_subscriptions` — current GraphQL websocket
  subscription count.

Custom instruments owned by the indexer:

- `indexer_events_processed_total` — counter of decoded + published
  events, labelled by `event_name` and `chain_id`.
- `indexer_last_block_indexed` — gauge of the highest confirmed block
  the subscriber has acknowledged.
- `indexer_reorgs_total` — counter of reorg detections.
- `indexer_publish_errors_total` — counter of NATS publish failures,
  labelled by bounded `error_class`.

## No-PII contract

Label cardinality is the most common operator footgun in a Prometheus
deployment. Both services enforce a "no PII labels" contract via
`Test_NoPIILabels` in their respective observability packages:

- No wallet addresses, transaction hashes, block hashes, or contract
  addresses appear as label values.
- No request IDs, JWT subjects, or session IDs appear as label
  values.
- The forbidden-label test asserts the absence of `=\"0x` substrings
  on the scrape body so any 0x-prefixed leak fails CI loudly.

When adding a new metric, keep its labels bounded by the contract
ABI, the chain ID, or a finite hand-curated set (see
`ErrClassTimeout` etc. in `services/indexer-go/internal/observability/indexer_metrics.go`).
