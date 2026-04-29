// Package observability — Prometheus metrics surface (M4 #94, indexer side).
//
// The indexer's observability stack already pushes metrics to an OTLP
// gRPC collector via Init in otel.go. This file adds the parallel
// Prometheus *pull* surface so the same instruments are also scraped
// by a Prometheus server (and visualised in the Grafana dashboards
// shipped under ops/grafana/).
//
// Both pipelines feed off the same OTel MeterProvider, so call sites
// only ever talk to one meter and a single Int64Counter shows up in
// both pipelines automatically.
//
// # Configuration env (in addition to those in otel.go)
//
//   - INDEXER_METRICS_ADDR    listen address for the /metrics endpoint
//                             (e.g. ":9091"). When unset, the Prometheus
//                             surface is not mounted. Bound to a
//                             dedicated listener so internal scrapers
//                             reach the surface without going through
//                             the indexer's NATS / RPC interfaces.
//
// # No-PII contract
//
// Same contract as the gateway twin: per-block-hash, per-tx-hash,
// per-address labels are out of bounds. We ship `chain_id` (low
// cardinality) and event-name labels (bounded by the contract ABI)
// only. Tested in metrics_test.go.
package observability

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// promRegistry holds the indexer's private Prometheus registry. We do
// NOT use prometheus.DefaultRegisterer because that is process-global
// and any imported package can clobber it — the isolated registry
// keeps /metrics output to indexer-emitted samples only.
var promRegistry = newPromRegistry()

// newPromRegistry constructs the package-private Prometheus registry
// pre-loaded with the stdlib process + Go runtime collectors. The
// Grafana dashboard (ops/grafana/dashboards/indexer.json) charts
// `process_resident_memory_bytes` directly; the OTel-Prom bridge does
// NOT export those samples on its own, so we MUST register the
// collectors on the same registry the bridge writes into.
//
// MustRegister is safe here because the registry is freshly minted on
// every call — collisions are impossible. The same path is reused by
// resetPrometheusForTest so unit tests start with the same baseline.
func newPromRegistry() *prometheus.Registry {
	r := prometheus.NewRegistry()
	r.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	r.MustRegister(collectors.NewGoCollector())
	return r
}

var (
	mpMu              sync.Mutex
	indexerMP         *sdkmetric.MeterProvider
	indexerPromReader sdkmetric.Reader
)

// PrometheusHandler returns the http.Handler that serves the
// Prometheus scrape endpoint. See the gateway twin for the full
// contract; the indexer-side lifecycle is identical.
func PrometheusHandler() (http.Handler, func() error, error) {
	mpMu.Lock()
	defer mpMu.Unlock()

	if indexerPromReader != nil {
		return newPromHTTPHandler(), shutdownPromReader, nil
	}

	reader, err := otelprom.New(otelprom.WithRegisterer(promRegistry))
	if err != nil {
		return nil, nil, fmt.Errorf("observability: prometheus exporter: %w", err)
	}

	if existing, ok := otel.GetMeterProvider().(*sdkmetric.MeterProvider); ok {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_ = existing.Shutdown(shutdownCtx)
		cancel()
		indexerMP = newPromOnlyMeterProvider(reader)
	} else {
		indexerMP = newPromOnlyMeterProvider(reader)
	}
	otel.SetMeterProvider(indexerMP)
	indexerPromReader = reader

	return newPromHTTPHandler(), shutdownPromReader, nil
}

// AttachPrometheusReader builds (or returns the cached) Prometheus
// reader so Init in otel.go can append it to its MeterProvider's
// reader list. This is the preferred wiring path: a single
// MeterProvider feeds both OTLP and Prometheus so the
// PrometheusHandler rebuild dance is avoided.
func AttachPrometheusReader() (sdkmetric.Reader, bool, error) {
	mpMu.Lock()
	defer mpMu.Unlock()

	if indexerPromReader != nil {
		return indexerPromReader, true, nil
	}
	reader, err := otelprom.New(otelprom.WithRegisterer(promRegistry))
	if err != nil {
		return nil, false, fmt.Errorf("observability: prometheus exporter: %w", err)
	}
	indexerPromReader = reader
	return reader, true, nil
}

// takeAttachedPrometheusReader is the seam used by Init.
func takeAttachedPrometheusReader() (sdkmetric.Reader, bool) {
	mpMu.Lock()
	defer mpMu.Unlock()
	if indexerPromReader == nil {
		return nil, false
	}
	return indexerPromReader, true
}

// MetricsHandler returns the http.Handler that serves the Prometheus
// scrape endpoint over the package's private registry.
func MetricsHandler() http.Handler { return newPromHTTPHandler() }

func newPromHTTPHandler() http.Handler {
	return promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{
		ErrorHandling:     promhttp.HTTPErrorOnError,
		EnableOpenMetrics: false,
	})
}

func newPromOnlyMeterProvider(reader sdkmetric.Reader) *sdkmetric.MeterProvider {
	return sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
}

// shutdownPromReader tears the Prometheus reader down. The 5s cap
// matches the OTLP reader's shutdown bound in otel.go.
func shutdownPromReader() error {
	mpMu.Lock()
	defer mpMu.Unlock()
	if indexerPromReader == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := indexerPromReader.Shutdown(ctx)
	indexerPromReader = nil
	if err != nil {
		return fmt.Errorf("prometheus reader shutdown: %w", err)
	}
	return nil
}

// promRegistryForTest exposes the private registry to tests in this
// package.
func promRegistryForTest() *prometheus.Registry { return promRegistry }

// resetPrometheusForTest is the test-only reset used by
// metrics_test.go.
func resetPrometheusForTest() {
	mpMu.Lock()
	defer mpMu.Unlock()
	if indexerPromReader != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		_ = indexerPromReader.Shutdown(ctx)
		cancel()
	}
	indexerPromReader = nil
	indexerMP = nil
	promRegistry = newPromRegistry()
}

// Meter returns an OTel meter via the global provider. Mirrors
// otel.Tracer in spirit and centralises the meter scope name on the
// caller — see the indexerMeterName constant in indexer_metrics.go.
func Meter(name string) metric.Meter { return otel.Meter(name) }
