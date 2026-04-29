package graph

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/observability"
)

// Test_NATSBus_GraphQLSubscriptionGaugeTransitions verifies the M4
// dashboard contract for the gateway_graphql_active_subscriptions
// gauge: opening a subscription must increment it from 0 to 1, and
// cancelling the subscription's context must decrement it back to 0.
// Without this guard the GraphQL subscription panel on
// ops/grafana/dashboards/gateway.json would either flatline at zero
// or grow unbounded as real subscriptions come and go.
//
// We probe the gauge through the Prometheus exposition layer rather
// than reading the OTel SDK directly — that is the same path operators
// scrape, so the assertion catches both a wiring break and a metric-
// rename regression at the exposition boundary.
func Test_NATSBus_GraphQLSubscriptionGaugeTransitions(t *testing.T) {
	observability.ResetForTest()

	h, shutdown, err := observability.PrometheusHandler()
	if err != nil {
		t.Fatalf("PrometheusHandler: %v", err)
	}
	defer func() { _ = shutdown() }()

	nc := runEmbeddedNATS(t)
	bus, err := NewNATSBus(nc)
	if err != nil {
		t.Fatalf("NewNATSBus: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	if _, err := bus.SubscribeTrades(ctx, nil); err != nil {
		t.Fatalf("SubscribeTrades: %v", err)
	}

	if !waitForGaugeValue(t, h, "gateway_graphql_active_subscriptions", "1", 2*time.Second) {
		t.Fatalf("expected gauge=1 after SubscribeTrades; body=\n%s", scrapeText(t, h))
	}

	cancel()

	if !waitForGaugeValue(t, h, "gateway_graphql_active_subscriptions", "0", 2*time.Second) {
		t.Fatalf("expected gauge=0 after ctx cancel; body=\n%s", scrapeText(t, h))
	}
}

// waitForGaugeValue polls the scrape endpoint until the named gauge
// equals `want`. See the SSE-side twin for the rationale on the poll
// shape.
func waitForGaugeValue(t *testing.T, h http.Handler, name, want string, timeout time.Duration) bool {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if hasGaugeValue(scrapeText(t, h), name, want) {
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return false
}

// hasGaugeValue scans Prometheus exposition for a sample of `name`
// (with any label set) whose value equals `want`. Tolerant of label
// permutations so a future reordering by the OTel-Prom bridge does
// not break the test.
func hasGaugeValue(body, name, want string) bool {
	for _, line := range strings.Split(body, "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.HasPrefix(line, name) {
			continue
		}
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
