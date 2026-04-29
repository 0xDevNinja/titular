package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/middleware"
)

// installTestMeterProvider builds an SDK MeterProvider with a
// Prometheus reader writing into a fresh prometheus.Registry, swaps it
// in as the global, and returns an http.Handler that scrapes the same
// registry. The deferred Cleanup restores the previous global and
// shuts the provider down so this test does not leak meter state into
// any other package-level test.
func installTestMeterProvider(t *testing.T) http.Handler {
	t.Helper()
	reg := prometheus.NewRegistry()
	reader, err := otelprom.New(otelprom.WithRegisterer(reg))
	if err != nil {
		t.Fatalf("otelprom.New: %v", err)
	}
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	prev := otel.GetMeterProvider()
	otel.SetMeterProvider(mp)
	t.Cleanup(func() {
		_ = mp.Shutdown(t.Context())
		otel.SetMeterProvider(prev)
	})
	return promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
}

// Test_HTTPMetrics_RecordsHistogramOnScrape exercises the canonical
// happy path: hit a couple of routes through gin, scrape the global
// MeterProvider's Prometheus surface, assert the histogram's
// _bucket / _sum / _count families are emitted with the
// http_route + http_request_method + http_response_status_code
// labels the Grafana dashboards key off.
func Test_HTTPMetrics_RecordsHistogramOnScrape(t *testing.T) {
	scrape := installTestMeterProvider(t)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.HTTPMetrics())
	r.GET("/agents/:id", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.POST("/jobs", func(c *gin.Context) {
		c.Status(http.StatusCreated)
	})

	for _, req := range []*http.Request{
		httptest.NewRequest(http.MethodGet, "/agents/abc", nil),
		httptest.NewRequest(http.MethodGet, "/agents/def", nil),
		httptest.NewRequest(http.MethodPost, "/jobs", nil),
	} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}

	body := scrapeBody(t, scrape)

	wantTokens := []string{
		"http_server_request_duration_seconds_bucket",
		"http_server_request_duration_seconds_sum",
		"http_server_request_duration_seconds_count",
		`http_route="/agents/:id"`,
		`http_route="/jobs"`,
		`http_request_method="GET"`,
		`http_request_method="POST"`,
		`http_response_status_code="200"`,
		`http_response_status_code="201"`,
		// SLO boundary at 200ms must appear as an explicit bucket so
		// histogram_quantile can interpolate at the SLO threshold.
		`le="0.2"`,
	}
	for _, tok := range wantTokens {
		if !strings.Contains(body, tok) {
			t.Errorf("scrape missing %q\n--- body ---\n%s", tok, body)
		}
	}
}

// Test_HTTPMetrics_UnmatchedRouteFallsBackToUnknown guards the
// cardinality contract: a 404 against an unregistered path must NOT
// stamp the raw URL onto the http_route label (it is unbounded), and
// MUST instead use the fixed "unknown" bucket. Without this guard a
// drive-by scanner walking thousands of random paths would explode
// the time-series database.
func Test_HTTPMetrics_UnmatchedRouteFallsBackToUnknown(t *testing.T) {
	scrape := installTestMeterProvider(t)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.HTTPMetrics())
	// No routes — every request 404s.

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/totally/unknown/path", nil))

	body := scrapeBody(t, scrape)
	if !strings.Contains(body, `http_route="unknown"`) {
		t.Errorf("expected http_route=\"unknown\" on 404; body=\n%s", body)
	}
	if strings.Contains(body, `http_route="/totally/unknown/path"`) {
		t.Errorf("raw URL leaked onto http_route label; body=\n%s", body)
	}
}

func scrapeBody(t *testing.T, h http.Handler) string {
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
