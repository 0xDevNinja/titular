package observability

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Test_PrometheusHandler_ServesEmpty200 verifies a fresh indexer's
// /metrics endpoint responds 200 with the standard Prometheus
// content-type — even when no instrument has yet been touched.
func Test_PrometheusHandler_ServesEmpty200(t *testing.T) {
	resetPrometheusForTest()

	h, shutdown, err := PrometheusHandler()
	if err != nil {
		t.Fatalf("PrometheusHandler: %v", err)
	}
	defer func() { _ = shutdown() }()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%q", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") &&
		!strings.HasPrefix(ct, "application/openmetrics-text") {
		t.Errorf("unexpected Content-Type %q", ct)
	}
}

// Test_IndexerInstruments_AppearOnScrape exercises every indexer-side
// custom instrument and verifies it shows up in the Prometheus
// exposition. Without this test, an instrument naming change would
// silently break the Grafana dashboards.
func Test_IndexerInstruments_AppearOnScrape(t *testing.T) {
	resetPrometheusForTest()

	h, shutdown, err := PrometheusHandler()
	if err != nil {
		t.Fatalf("PrometheusHandler: %v", err)
	}
	defer func() { _ = shutdown() }()

	ctx := context.Background()
	EventsProcessed().Add(ctx, 5,
		metric.WithAttributes(
			attribute.String("event_name", "Trade"),
			attribute.Int("chain_id", 84532),
		),
	)
	LastBlockIndexed().Record(ctx, 12345,
		metric.WithAttributes(attribute.Int("chain_id", 84532)),
	)
	Reorgs().Add(ctx, 1, metric.WithAttributes(attribute.Int("chain_id", 84532)))
	PublishErrors().Add(ctx, 2,
		metric.WithAttributes(
			attribute.String("error_class", ErrClassTimeout),
			attribute.Int("chain_id", 84532),
		),
	)

	body := scrape(t, h)

	for _, name := range []string{
		"indexer_events_processed_total",
		"indexer_last_block_indexed",
		"indexer_reorgs_total",
		"indexer_publish_errors_total",
	} {
		if !strings.Contains(body, name) {
			t.Errorf("scrape body missing metric %q\n--- body ---\n%s", name, body)
		}
	}
	// Hand-roll a couple of value assertions so we know the values are
	// actually flowing through, not just the metric names showing up
	// from the # HELP / # TYPE prelude.
	if !strings.Contains(body, "indexer_last_block_indexed") ||
		!strings.Contains(body, "12345") {
		t.Errorf("indexer_last_block_indexed sample missing 12345; body=\n%s", body)
	}
}

// Test_NoPIILabels enforces the cardinality contract — block hashes,
// transaction hashes, and contract addresses MUST NOT show up as
// label values. The test deliberately includes the lower-case 0x
// prefix in the forbidden list because every Ethereum address starts
// with 0x; if any label value carries one, the test fails loudly.
func Test_NoPIILabels(t *testing.T) {
	resetPrometheusForTest()

	h, shutdown, err := PrometheusHandler()
	if err != nil {
		t.Fatalf("PrometheusHandler: %v", err)
	}
	defer func() { _ = shutdown() }()

	// Prime the instruments with allowed labels only.
	ctx := context.Background()
	EventsProcessed().Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("event_name", "Trade"),
			attribute.Int("chain_id", 84532),
		),
	)
	PublishErrors().Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("error_class", ErrClassUnmapped),
			attribute.Int("chain_id", 84532),
		),
	)

	body := scrape(t, h)

	forbidden := []string{
		"address=", "tx_hash=", "block_hash=", "wallet=", "from=", "to=",
		"=\"0x", // any 0x-prefixed value sneaking in as a label.
	}
	for _, f := range forbidden {
		if strings.Contains(body, f) {
			t.Errorf("metrics body contains forbidden token %q\n--- body ---\n%s", f, body)
		}
	}
}

// Test_IndexerInstrumentsBuilt verifies the lazy-init helper builds
// every instrument cleanly. A regression here means subsequent
// Add/Record calls silently no-op and dashboards go dark.
func Test_IndexerInstrumentsBuilt(t *testing.T) {
	resetPrometheusForTest()

	if _, _, err := PrometheusHandler(); err != nil {
		t.Fatalf("PrometheusHandler: %v", err)
	}

	if EventsProcessed() == nil {
		t.Fatal("EventsProcessed nil")
	}
	if LastBlockIndexed() == nil {
		t.Fatal("LastBlockIndexed nil")
	}
	if Reorgs() == nil {
		t.Fatal("Reorgs nil")
	}
	if PublishErrors() == nil {
		t.Fatal("PublishErrors nil")
	}
	if indexerInstrumentsErr != nil {
		t.Fatalf("indexerInstrumentsErr = %v", indexerInstrumentsErr)
	}
}

// Test_AttachPrometheusReader_Idempotent ensures repeated calls return
// the same reader without panicking on duplicate exporter registration.
func Test_AttachPrometheusReader_Idempotent(t *testing.T) {
	resetPrometheusForTest()

	r1, ok1, err := AttachPrometheusReader()
	if err != nil || !ok1 || r1 == nil {
		t.Fatalf("first attach: ok=%v err=%v", ok1, err)
	}
	r2, ok2, err := AttachPrometheusReader()
	if err != nil || !ok2 || r2 == nil {
		t.Fatalf("second attach: ok=%v err=%v", ok2, err)
	}
	if r1 != r2 {
		t.Fatal("AttachPrometheusReader returned distinct readers; expected cached singleton")
	}
}

// Test_OTLPAndPrometheus_ShareInstruments simulates the boot path
// where AttachPrometheusReader runs first, then Init. Both readers
// should be wired to a single MeterProvider so a counter recorded
// after Init still reaches the Prometheus surface.
func Test_OTLPAndPrometheus_ShareInstruments(t *testing.T) {
	resetPrometheusForTest()

	// Prom reader first — the boot sequence main.go uses.
	if _, _, err := AttachPrometheusReader(); err != nil {
		t.Fatalf("AttachPrometheusReader: %v", err)
	}

	// Mimic Init by installing a MeterProvider that also has the
	// Prometheus reader. We can't easily run the real Init without an
	// OTLP collector, so we exercise the takeAttachedPrometheusReader
	// + provider construction path directly.
	reader, ok := takeAttachedPrometheusReader()
	if !ok {
		t.Fatal("expected reader to be attached")
	}
	mp := newPromOnlyMeterProvider(reader)
	otel.SetMeterProvider(mp)

	// Now record a sample via the global meter and verify the scrape
	// surface sees it.
	ctr, err := otel.Meter("test").Int64Counter("titular_test_share")
	if err != nil {
		t.Fatalf("Int64Counter: %v", err)
	}
	ctr.Add(context.Background(), 1)

	h := MetricsHandler()
	body := scrape(t, h)
	if !strings.Contains(body, "titular_test_share") {
		t.Errorf("scrape missing titular_test_share; body=\n%s", body)
	}
}

func scrape(t *testing.T, h http.Handler) string {
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
