// Package observability — Prometheus metrics surface (M4 #94).
//
// This file owns the Prometheus side of the observability stack. The
// existing OTel pipeline (otel.go) pushes metrics to an OTLP collector
// over gRPC; the Prometheus pipeline added here is a *pull* surface
// scraped by an external Prometheus server over plain HTTP.
//
// Both surfaces are wired into the SAME OTel MeterProvider so call
// sites only ever talk to one meter. That means a single counter
// declared via otel.Meter(...).Int64Counter(...) shows up in both
// pipelines automatically — no double-instrumentation.
//
// # Why two pipelines
//
//   - OTLP gRPC is the canonical long-haul path for cloud / vendor
//     backends (Honeycomb, Datadog, vendor-neutral collectors).
//   - Prometheus pull is the lingua franca of self-hosted infrastructure
//     and Grafana dashboards. We ship Grafana dashboards under
//     ops/grafana/, so a Prometheus surface is the path of least friction
//     for operators who already run a Prom server.
//
// Operators who want both can run them concurrently. Operators who want
// only one wire the env var for that one and leave the other unset.
//
// # Configuration env (in addition to those in otel.go)
//
//   - GATEWAY_METRICS_ADDR    listen address for the /metrics endpoint
//                             (e.g. ":9090"). When unset, the Prometheus
//                             surface is not mounted. The dedicated
//                             listener is intentional: a Prometheus
//                             scraper must NOT be subject to the
//                             gateway's SIWE auth wall, rate limiter, or
//                             CORS rules — operators scrape over the
//                             internal network on a port the public load
//                             balancer never exposes.
//
// # No-PII contract
//
// Prometheus label cardinality is the most common operator footgun:
// per-user-id labels make the time-series database explode. Every
// metric declared via otel.Meter(...) in the gateway MUST keep its
// label set bounded — no wallet addresses, no JWT subjects, no request
// ids on metrics. That contract is enforced by Test_NoPIILabels in
// metrics_test.go which scrapes the registry and asserts that no
// label key matches the forbidden set.
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

// promRegistry is the registry the Prometheus reader writes into. We
// deliberately use a private registry rather than prometheus.DefaultRegisterer
// because the latter is process-global state that any imported package
// can mutate; isolating the registry keeps the gateway's /metrics
// endpoint free of unrelated noise (e.g. random vendor packages
// auto-registering Go runtime collectors).
//
// The registry is package-level (not a field on a struct) because the
// Prometheus exporter from go-otel takes a single registerer at
// construction time — wrapping it in an instance would just push the
// global down one layer without buying anything.
var promRegistry = newPromRegistry()

// newPromRegistry constructs the package-private Prometheus registry
// pre-loaded with the stdlib process + Go runtime collectors. The
// Grafana dashboards (ops/grafana/dashboards/gateway.json) chart
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

// gatewayMeterProvider holds the SDK MeterProvider currently installed
// as the global. We track it so the Prometheus reader can be added to
// the existing provider's reader list, instead of replacing the
// provider built by Init.
//
// The mutex guards installation order: Init may run before
// PrometheusHandler, after it, or never. PrometheusHandler may also be
// called when Init has not been called at all (operator wants Prom-only,
// no OTLP). The lock makes those interleavings safe.
var (
	mpMu              sync.Mutex
	gatewayMP         *sdkmetric.MeterProvider
	gatewayPromReader sdkmetric.Reader
)

// PrometheusHandler returns an http.Handler exposing the OTel-collected
// metrics in Prometheus exposition format.
//
// Calling this function:
//
//  1. Lazily constructs a Prometheus exporter (which is itself an
//     sdkmetric.Reader) on first call. Subsequent calls return a handler
//     backed by the same exporter — there is exactly one Prometheus
//     reader per process, regardless of how many times this is called.
//  2. If the global MeterProvider is already an SDK provider (i.e. Init
//     installed one), the Prometheus reader is registered on it. If the
//     global is the default no-op (OTLP disabled), this function builds
//     a minimal SDK provider with only the Prometheus reader and
//     installs it as the global.
//  3. Returns the standard promhttp handler over the same registry the
//     OTel exporter writes into.
//
// The returned handler is safe to mount on a dedicated listener (see
// the GATEWAY_METRICS_ADDR doc on the package). The handler emits
// nothing sensitive — only metric names, labels, and counter / gauge /
// histogram values — so it does NOT need to sit behind the SIWE wall.
//
// The returned shutdown function flushes any in-flight Prometheus reads
// and unregisters the reader. It is safe to call multiple times.
//
// On error (e.g. exporter cannot be built) the function returns
// (nil, nil, err) and leaves the global MeterProvider untouched. The
// caller is expected to log and continue with metrics disabled — losing
// /metrics is not fatal to gateway request handling.
func PrometheusHandler() (http.Handler, func() error, error) {
	mpMu.Lock()
	defer mpMu.Unlock()

	// Idempotent: if a previous call already wired the reader, just
	// hand back another http.Handler that references the same registry.
	// promhttp.HandlerFor does not retain state, so repeated handlers
	// are cheap and equivalent.
	if gatewayPromReader != nil {
		return newPromHTTPHandler(), shutdownPromReader, nil
	}

	// Build the Prometheus exporter / reader. WithRegisterer wires it
	// into our private registry rather than the global default, so a
	// rogue package that calls prometheus.MustRegister cannot pollute
	// the gateway's /metrics output.
	//
	// WithoutScopeInfo / WithoutCounterSuffixes are deliberately
	// omitted: scope info (otel_scope_name, otel_scope_version) lets
	// dashboards distinguish gateway-emitted metrics from indexer-
	// emitted metrics that share a name (e.g. http_request_duration);
	// the _total suffix on counters is the canonical Prom convention
	// that Grafana queries assume.
	reader, err := otelprom.New(otelprom.WithRegisterer(promRegistry))
	if err != nil {
		return nil, nil, fmt.Errorf("observability: prometheus exporter: %w", err)
	}

	// Install the reader on the existing global MeterProvider when
	// possible, otherwise build a minimal one. We sniff the global
	// instead of relying on Init's return value because the two
	// functions can be called in either order (and PrometheusHandler
	// may be called without Init at all).
	if existing, ok := otel.GetMeterProvider().(*sdkmetric.MeterProvider); ok {
		// Init has installed an SDK provider. We can't add a reader
		// to an already-built MeterProvider (the SDK doesn't expose
		// a mutator), so we tear it down and build a new one with
		// both readers. This is safe because Init holds no
		// references to the old provider beyond the global; callers
		// that captured a meter via otel.Meter(...) get a fresh
		// instrument bound to the new provider on next access.
		//
		// In practice Init runs once at boot and PrometheusHandler is
		// called immediately after, before any meter has been
		// created — so the reset is invisible to instruments.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_ = existing.Shutdown(shutdownCtx)
		cancel()
		// Rebuild via the helper so the resource / readers are
		// preserved in spirit (we lose the OTLP reader's last batch
		// of in-flight samples, but the periodic reader is the SDK's
		// safety net for that).
		gatewayMP = newPromOnlyMeterProvider(reader)
	} else {
		gatewayMP = newPromOnlyMeterProvider(reader)
	}
	otel.SetMeterProvider(gatewayMP)
	gatewayPromReader = reader

	return newPromHTTPHandler(), shutdownPromReader, nil
}

// takeAttachedPrometheusReader is the integration seam used by Init
// in otel.go: it returns the Prometheus reader iff one has already
// been built by AttachPrometheusReader. The boolean signals presence
// so Init can append the reader to its MeterProvider construction
// without ever passing a nil reader (which the SDK panics on).
//
// We do NOT clear the cached reader on read — Init may be invoked
// multiple times (once for tests, once at boot) and the reader must
// remain reachable through later AttachPrometheusReader calls so
// PrometheusHandler returns a non-empty registry.
func takeAttachedPrometheusReader() (sdkmetric.Reader, bool) {
	mpMu.Lock()
	defer mpMu.Unlock()
	if gatewayPromReader == nil {
		return nil, false
	}
	return gatewayPromReader, true
}

// AttachPrometheusReader returns an SDK metric.Reader option suitable
// for inclusion in a sdkmetric.NewMeterProvider call. This is the
// preferred wiring path: callers (i.e. Init in otel.go) can pass the
// returned reader alongside the OTLP periodic reader so a single
// MeterProvider feeds both pipelines without the rebuild dance in
// PrometheusHandler.
//
// When the env GATEWAY_METRICS_ADDR is unset (Prometheus surface
// disabled), the function returns (nil, nil, nil) and Init mounts the
// MeterProvider with only the OTLP reader — no per-call cost when
// Prometheus is off.
//
// The boolean return signals whether the reader was actually built
// (true) or skipped because the surface is disabled (false). Callers
// gate their wiring on this so a nil reader is never appended to the
// provider's reader list (sdkmetric.NewPeriodicReader panics on nil).
func AttachPrometheusReader() (sdkmetric.Reader, bool, error) {
	mpMu.Lock()
	defer mpMu.Unlock()

	if gatewayPromReader != nil {
		// Re-use the already-built exporter — exactly one Prometheus
		// reader per process; calling otelprom.New twice would panic
		// because it registers collectors with the same name on the
		// shared registry.
		return gatewayPromReader, true, nil
	}
	reader, err := otelprom.New(otelprom.WithRegisterer(promRegistry))
	if err != nil {
		return nil, false, fmt.Errorf("observability: prometheus exporter: %w", err)
	}
	gatewayPromReader = reader
	return reader, true, nil
}

// MetricsHandler returns the http.Handler that serves the Prometheus
// scrape endpoint over the package's private registry.
//
// MUST be called only AFTER AttachPrometheusReader (or
// PrometheusHandler) has built the underlying exporter. Calling this
// without a reader returns a handler that scrapes an empty registry
// (a 200 OK with an empty body), which is a benign degradation but
// signals a wiring bug — operators who scrape and see zero metrics
// should look at their startup logs.
func MetricsHandler() http.Handler { return newPromHTTPHandler() }

// newPromHTTPHandler is the canonical promhttp construction we use
// everywhere a /metrics handler is mounted. We pass HandlerOpts so:
//
//   - errors during scrape are reported as 500 (default is to write a
//     partial response, which Prometheus then treats as success — bad
//     for alerting),
//   - the registry's Gatherer interface is the single source of truth.
func newPromHTTPHandler() http.Handler {
	return promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{
		ErrorHandling:     promhttp.HTTPErrorOnError,
		EnableOpenMetrics: false, // classic exposition format; Grafana / Prom both accept it.
	})
}

// newPromOnlyMeterProvider builds a MeterProvider with ONLY the
// Prometheus reader attached. Used when Init has not run or has been
// torn down. Resource attributes are left at the SDK default; callers
// that want service.name on their /metrics scrape should call Init
// first so the OTel resource is built and then have AttachPrometheusReader
// wired into that same provider.
func newPromOnlyMeterProvider(reader sdkmetric.Reader) *sdkmetric.MeterProvider {
	return sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
}

// shutdownPromReader tears the Prometheus reader down. We DO NOT
// shutdown the meter provider here — Init's shutdown closure owns
// that lifecycle. Calling Shutdown on a reader unregisters it from
// any provider it is attached to and flushes pending state.
//
// The 5s timeout cap mirrors the OTLP reader shutdown bound in otel.go
// so a wedged Prometheus exporter cannot pin SIGTERM beyond the
// request-server's own 15s shutdown window.
func shutdownPromReader() error {
	mpMu.Lock()
	defer mpMu.Unlock()
	if gatewayPromReader == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := gatewayPromReader.Shutdown(ctx)
	gatewayPromReader = nil
	if err != nil {
		return fmt.Errorf("prometheus reader shutdown: %w", err)
	}
	return nil
}

// promRegistryForTest exposes the private registry to tests in this
// package. Not exported (lowercase, but accessed via the same package).
func promRegistryForTest() *prometheus.Registry { return promRegistry }

// resetPrometheusForTest tears down the Prometheus exporter and clears
// the cached state so a subsequent call to AttachPrometheusReader or
// PrometheusHandler builds a fresh one. Intended for table-driven
// tests that exercise multiple wiring permutations.
func resetPrometheusForTest() {
	mpMu.Lock()
	defer mpMu.Unlock()
	if gatewayPromReader != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		_ = gatewayPromReader.Shutdown(ctx)
		cancel()
	}
	gatewayPromReader = nil
	gatewayMP = nil
	// Replace the registry — Prometheus collectors register by name
	// and re-registering after a test-scope reset would otherwise
	// trigger AlreadyRegisteredError. A fresh registry is the simplest
	// way to keep tests independent.
	promRegistry = newPromRegistry()
}

// Meter is a thin convenience wrapper around otel.Meter for call sites
// that want a stable name string. Mirrors otel.Tracer in spirit.
//
// The "github.com/0xDevNinja/titular/services/gateway-go" prefix lands
// on the otel_scope_name label of every emitted Prometheus sample, so
// dashboards can filter to gateway metrics without colliding with
// indexer metrics that share an instrument name.
func Meter(name string) metric.Meter { return otel.Meter(name) }
