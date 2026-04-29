// Package observability — gateway-specific custom metrics (M4 #94).
//
// HTTP request count, latency histogram, in-flight count, DB query
// duration and connection-pool stats already flow through the OTel
// auto-instrumentation packages (otelgin, otelpgx) wired in #93, so
// they appear on the /metrics endpoint without any extra code here.
//
// This file owns the *custom* gateway-shaped metrics — the ones the
// auto-instrumentation can't possibly know about because they describe
// gateway-internal state:
//
//   - sse_active_connections_gauge — per-process SSE subscriber count.
//   - graphql_active_subscriptions_gauge — open GraphQL websocket
//     subscriptions across all clients.
//
// All custom metrics MUST keep label cardinality bounded. We
// deliberately do NOT add per-client-IP, per-wallet, or per-request-id
// labels: those would multiply the time-series footprint by every
// unique value the gateway has ever seen, which on a public endpoint
// is unbounded. The Test_NoPIILabels guard in metrics_test.go
// enforces this contract.
package observability

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/metric"
)

// gatewayMeterName is the otel.Meter name that lands on the
// otel_scope_name label of every sample. Stable per service so
// dashboards can filter `{otel_scope_name="github.com/0xDevNinja/.../gateway-go"}`
// to isolate gateway metrics from indexer metrics that share an
// instrument name (e.g. http_server_request_duration_seconds).
const gatewayMeterName = "github.com/0xDevNinja/titular/services/gateway-go"

// gatewayInstruments is the package-level cache of constructed
// instruments. Built lazily on first Get… call so we don't pay the
// instrument-creation cost on a gateway started without observability
// (the no-op MeterProvider returns no-op instruments anyway, but the
// allocation is non-trivial when multiplied by every meter in the
// codebase).
var (
	gatewayInstrumentsOnce sync.Once
	gatewaySSEGauge        metric.Int64UpDownCounter
	gatewayGQLSubGauge     metric.Int64UpDownCounter
	gatewayInstrumentsErr  error
)

// initGatewayInstruments builds the custom-metric instruments. The
// upstream go-otel SDK returns errors lazily when an instrument's name
// violates the spec — we surface them as gatewayInstrumentsErr so the
// caller sees a structured failure rather than a silent no-op.
//
// Both instruments are Int64UpDownCounter (a.k.a. Gauge) because they
// represent a "current number of …" which can both increment and
// decrement over the process lifetime. Cumulative counters (Counter)
// would be wrong here — they only ever increase.
func initGatewayInstruments() {
	gatewayInstrumentsOnce.Do(func() {
		m := Meter(gatewayMeterName)

		var err error
		gatewaySSEGauge, err = m.Int64UpDownCounter(
			"gateway_sse_active_connections",
			metric.WithDescription("Number of currently-open Server-Sent Events subscriber connections."),
			metric.WithUnit("{connection}"),
		)
		if err != nil {
			gatewayInstrumentsErr = err
			return
		}

		gatewayGQLSubGauge, err = m.Int64UpDownCounter(
			"gateway_graphql_active_subscriptions",
			metric.WithDescription("Number of currently-open GraphQL websocket subscriptions across all clients."),
			metric.WithUnit("{subscription}"),
		)
		if err != nil {
			gatewayInstrumentsErr = err
			return
		}
	})
}

// SSEConnections returns the up-down counter that tracks the live
// gateway SSE subscriber count. Call sites in internal/sse increment
// on Subscribe and decrement on Close so the gauge reflects the
// instantaneous fan-out.
//
// On the no-op MeterProvider (observability disabled) the returned
// instrument is a no-op — Add becomes a free function call. Callers
// don't need to guard.
//
// Returns a no-op instrument on internal initialisation error rather
// than nil; this keeps call sites free of nil-checks. The error path
// is observed via Test_GatewayInstrumentsBuilt in metrics_test.go.
func SSEConnections() metric.Int64UpDownCounter {
	initGatewayInstruments()
	if gatewaySSEGauge == nil {
		return noopUpDownCounter{}
	}
	return gatewaySSEGauge
}

// GraphQLSubscriptions returns the up-down counter for in-flight
// GraphQL websocket subscriptions. Increment on subscribe, decrement
// on unsubscribe / connection drop.
func GraphQLSubscriptions() metric.Int64UpDownCounter {
	initGatewayInstruments()
	if gatewayGQLSubGauge == nil {
		return noopUpDownCounter{}
	}
	return gatewayGQLSubGauge
}

// noopUpDownCounter is the fallback instrument we hand out when
// initialisation fails. The interface has only one method (Add) plus
// the embedded marker; both no-op cleanly so call sites can use the
// same code path with or without observability.
//
// We DO NOT use the upstream noop package because that requires a
// MeterProvider switch; here we want the call site to get a working
// instrument back from the package's own helper without touching the
// global state.
type noopUpDownCounter struct{ metric.Int64UpDownCounter }

func (noopUpDownCounter) Add(context.Context, int64, ...metric.AddOption) {}

// ResetForTest tears down the Prometheus reader, MeterProvider, and
// resets the lazy-init guard on the custom gateway instruments so a
// subsequent SSEConnections / GraphQLSubscriptions call rebuilds them
// against a freshly-installed MeterProvider. Intended for use by
// tests in sibling packages that exercise the gauge transitions
// (e.g. internal/sse and internal/graph).
//
// Safe to call from concurrent test goroutines: it serialises through
// the same mpMu used by AttachPrometheusReader / PrometheusHandler.
//
// NOT exported as a stable API — the suffix conveys it's a test seam.
// Production code MUST not call this; it would race with live
// /metrics scrapes and lose every accumulated sample.
func ResetForTest() {
	resetPrometheusForTest()
	gatewayInstrumentsOnce = sync.Once{}
	gatewaySSEGauge = nil
	gatewayGQLSubGauge = nil
	gatewayInstrumentsErr = nil
}
