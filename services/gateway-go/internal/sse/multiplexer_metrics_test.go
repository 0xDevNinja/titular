package sse

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/observability"
)

// Test_Multiplexer_SSEGaugeTransitions verifies the M4 dashboard
// contract for the gateway_sse_active_connections gauge: a fresh
// Subscribe must increment it from 0 to 1, and a matching Close must
// decrement it back to 0. Without this guard the SSE clients panel on
// ops/grafana/dashboards/gateway.json would either flatline at zero
// or grow unbounded as real connections come and go.
//
// We probe the gauge through the Prometheus exposition layer rather
// than reading the OTel SDK directly — that is the same path operators
// scrape, so the assertion catches both a wiring break and a metric-
// rename regression at the exposition boundary.
func Test_Multiplexer_SSEGaugeTransitions(t *testing.T) {
	observability.ResetForTest()

	h, shutdown, err := observability.PrometheusHandler()
	if err != nil {
		t.Fatalf("PrometheusHandler: %v", err)
	}
	defer func() { _ = shutdown() }()

	nc := runEmbeddedNATS(t)
	mux, err := NewMultiplexer(Config{NATS: nc})
	if err != nil {
		t.Fatalf("NewMultiplexer: %v", err)
	}
	if err := mux.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer mux.Stop()

	sub, err := mux.Subscribe(SubscribeOptions{})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	if !waitForGaugeValue(t, h, "gateway_sse_active_connections", "1", 2*time.Second) {
		t.Fatalf("expected gauge=1 after Subscribe; body=\n%s", scrapeText(t, h))
	}

	sub.Close()

	if !waitForGaugeValue(t, h, "gateway_sse_active_connections", "0", 2*time.Second) {
		t.Fatalf("expected gauge=0 after Close; body=\n%s", scrapeText(t, h))
	}
}

// waitForGaugeValue polls the scrape endpoint until the named gauge
// equals `want`, returning true on success. The poll interval (50ms)
// keeps the budget bounded: the OTel Prometheus reader is collect-on-
// scrape, so we expect convergence within a single scrape, but
// scheduler delay on the test box can push a Subscribe's Add through
// after the first scrape returns.
func waitForGaugeValue(t *testing.T, h http.Handler, name, want string, timeout time.Duration) bool {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		body := scrapeText(t, h)
		if hasGaugeValue(body, name, want) {
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return false
}

// hasGaugeValue scans Prometheus exposition for a sample of `name`
// (with any label set) whose value equals `want`. The scanner is
// deliberately tolerant of label permutations because the OTel-Prom
// bridge interleaves otel_scope_name with the user-supplied labels —
// a strict line match would fail spuriously when the bridge changes
// label ordering between releases.
func hasGaugeValue(body, name, want string) bool {
	for _, line := range strings.Split(body, "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.HasPrefix(line, name) {
			continue
		}
		// Expected shape: `metric{labels...} value` or `metric value`.
		// The value is the last whitespace-separated token on the line.
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if fields[len(fields)-1] == want {
			return true
		}
	}
	return false
}

// scrapeText invokes the handler and returns its body as a string.
func scrapeText(t *testing.T, h http.Handler) string {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("scrape status=%d body=%q", rec.Code, rec.Body.String())
	}
	body, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return string(body)
}
