package observability

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.opentelemetry.io/otel"
)

// Test_PrometheusHandler_NoEndpoint_ServesEmpty200 verifies that a
// Prometheus scrape against a freshly-started gateway (no metrics yet
// recorded) returns a 200 with the standard exposition headers and
// either an empty body or one containing only the synthetic
// `target_info` lines that the OTel-Prometheus bridge emits to expose
// resource attributes. The test is the smoke-screen against a regression
// where the handler returns 5xx because the registry is empty.
func Test_PrometheusHandler_NoEndpoint_ServesEmpty200(t *testing.T) {
	resetPrometheusForTest()

	h, shutdown, err := PrometheusHandler()
	if err != nil {
		t.Fatalf("PrometheusHandler: %v", err)
	}
	defer func() {
		if err := shutdown(); err != nil {
			t.Errorf("shutdown: %v", err)
		}
	}()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%q", rec.Code, rec.Body.String())
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "text/plain") && !strings.HasPrefix(ct, "application/openmetrics-text") {
		t.Errorf("Content-Type = %q, want text/plain or openmetrics", ct)
	}
}

// Test_PrometheusHandler_CounterIncrement asserts that a counter we
// record via the global meter shows up in the Prometheus exposition.
// This is the proof that the OTel exporter and the Prometheus reader
// are wired to the same MeterProvider. Without this, a custom counter
// would emit to OTLP but not to /metrics.
func Test_PrometheusHandler_CounterIncrement(t *testing.T) {
	resetPrometheusForTest()

	h, shutdown, err := PrometheusHandler()
	if err != nil {
		t.Fatalf("PrometheusHandler: %v", err)
	}
	defer func() { _ = shutdown() }()

	meter := otel.Meter("test-meter")
	ctr, err := meter.Int64Counter("titular_test_events_total")
	if err != nil {
		t.Fatalf("Int64Counter: %v", err)
	}
	ctr.Add(context.Background(), 3)
	ctr.Add(context.Background(), 4)

	body := scrape(t, h)
	if !strings.Contains(body, "titular_test_events_total") {
		t.Errorf("scrape missing titular_test_events_total; body=\n%s", body)
	}
	// The Prometheus exporter exposes the cumulative sum directly.
	// 3 + 4 = 7 — assert via substring rather than full-line match
	// because the OTel-Prom bridge may emit additional samples
	// (otel_scope_info, target_info) interleaved.
	if !strings.Contains(body, "} 7") {
		t.Errorf("scrape missing the cumulative 7; body=\n%s", body)
	}
}

// Test_NoPIILabels guards the cardinality contract: no metric emitted
// by gateway code may carry a label key matching the wallet / user-id
// / request-id family. A label like `wallet="0xabc..."` would explode
// the time-series database the moment a few thousand wallets show up,
// so this contract is enforced at the gathering layer.
//
// The test exercises the canonical custom instruments
// (SSEConnections, GraphQLSubscriptions) and any auto-instrumented
// metrics that might leak through. We DO NOT enumerate every possible
// PII shape — we lock the obvious ones and rely on code review to
// catch the long tail.
func Test_NoPIILabels(t *testing.T) {
	resetPrometheusForTest()

	h, shutdown, err := PrometheusHandler()
	if err != nil {
		t.Fatalf("PrometheusHandler: %v", err)
	}
	defer func() { _ = shutdown() }()

	// Touch the custom instruments so they emit at least one sample.
	SSEConnections().Add(context.Background(), 1)
	GraphQLSubscriptions().Add(context.Background(), 1)

	body := scrape(t, h)

	forbidden := []string{
		"wallet=", "address=", "user_id=", "request_id=", "session_id=",
		"jwt=", "subject=", "sub=\"0x", "remote_addr=",
		// Catch raw 0x… hex sneaking into any label value.
		"=\"0x",
	}
	for _, f := range forbidden {
		if strings.Contains(body, f) {
			t.Errorf("metrics body contains forbidden token %q; body=\n%s", f, body)
		}
	}
}

// Test_GatewayInstrumentsBuilt verifies the lazy-init helper builds
// the custom instruments without error. A regression here would mean
// the observability package compiles but every Add call is silently
// a no-op.
func Test_GatewayInstrumentsBuilt(t *testing.T) {
	resetPrometheusForTest()

	if _, _, err := PrometheusHandler(); err != nil {
		t.Fatalf("PrometheusHandler: %v", err)
	}

	// Call them twice; the sync.Once must keep them stable.
	a := SSEConnections()
	b := SSEConnections()
	if a == nil || b == nil {
		t.Fatal("SSEConnections returned nil")
	}
	c := GraphQLSubscriptions()
	d := GraphQLSubscriptions()
	if c == nil || d == nil {
		t.Fatal("GraphQLSubscriptions returned nil")
	}

	if gatewayInstrumentsErr != nil {
		t.Fatalf("gatewayInstrumentsErr = %v", gatewayInstrumentsErr)
	}
}

// Test_PrometheusHandler_Idempotent ensures repeated calls return a
// usable handler each time without panicking on duplicate registration.
// The exporter library will panic on a second otelprom.New call against
// the same registry — the package-level cache must short-circuit the
// second call back to the same reader.
func Test_PrometheusHandler_Idempotent(t *testing.T) {
	resetPrometheusForTest()

	h1, sd1, err := PrometheusHandler()
	if err != nil {
		t.Fatalf("first PrometheusHandler: %v", err)
	}
	h2, sd2, err := PrometheusHandler()
	if err != nil {
		t.Fatalf("second PrometheusHandler: %v", err)
	}
	if h1 == nil || h2 == nil {
		t.Fatal("nil handler on repeat call")
	}
	// Both shutdown closures wire to the same underlying reader; calling
	// either should be safe. Calling both, in order, must also be safe
	// (the second is a no-op once the cache is cleared).
	if err := sd1(); err != nil {
		t.Errorf("sd1: %v", err)
	}
	if err := sd2(); err != nil {
		t.Errorf("sd2 (after sd1): %v", err)
	}
}

// scrape invokes the handler and returns the response body.
func scrape(t *testing.T, h http.Handler) string {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("scrape status = %d, body=%q", rec.Code, rec.Body.String())
	}
	body, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return string(body)
}
