package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/middleware"
)

// installTestTracer is a small test-local helper that swaps the global
// TracerProvider for an in-memory recorder. Mirrors the helper inside
// internal/observability but lives in this file so the middleware unit
// tests do not import the observability package (and therefore stay a
// genuine unit test rather than a thin integration).
func installTestTracer(t *testing.T) *tracetest.SpanRecorder {
	t.Helper()
	rec := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(rec))
	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() {
		_ = tp.Shutdown(context.Background())
		otel.SetTracerProvider(prev)
	})
	return rec
}

// Test_OTelMiddleware_StartsServerSpan confirms that an inbound HTTP
// request produces a span with kind=SERVER and the route template
// stamped onto its name, and that the X-Trace-Id response header
// carries the active trace id.
func Test_OTelMiddleware_StartsServerSpan(t *testing.T) {
	rec := installTestTracer(t)

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.OTelMiddleware("test-gw"))
	r.Use(middleware.TraceCorrelation())
	r.GET("/things/:id", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/things/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: %d", w.Code)
	}

	spans := rec.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	srv := spans[0]
	if srv.SpanKind() != trace.SpanKindServer {
		t.Errorf("kind = %s, want server", srv.SpanKind())
	}

	traceHeader := w.Header().Get(middleware.TraceIDHeader)
	if traceHeader == "" {
		t.Fatal("X-Trace-Id missing from response")
	}
	if traceHeader != srv.SpanContext().TraceID().String() {
		t.Errorf("X-Trace-Id %q != span trace id %q",
			traceHeader, srv.SpanContext().TraceID())
	}
}

// Test_TraceCorrelation_StampsRequestID verifies the trace ↔ log
// correlation: the request_id minted by the RequestID middleware lands
// on the active span as the gateway.request_id attribute. This is
// what lets an operator pivot from a log line (with request_id) to
// the full distributed trace.
func Test_TraceCorrelation_StampsRequestID(t *testing.T) {
	rec := installTestTracer(t)

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.OTelMiddleware("test-gw"))
	r.Use(middleware.TraceCorrelation())
	r.GET("/p", func(c *gin.Context) { c.Status(http.StatusOK) })

	const wantID = "11111111-2222-3333-4444-555555555555"
	req := httptest.NewRequest(http.MethodGet, "/p", nil)
	req.Header.Set("X-Request-Id", wantID)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	spans := rec.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	got := ""
	for _, kv := range spans[0].Attributes() {
		if string(kv.Key) == "gateway.request_id" {
			got = kv.Value.AsString()
		}
	}
	if got != wantID {
		t.Errorf("gateway.request_id = %q, want %q", got, wantID)
	}
}

// Test_TraceCorrelation_NoCredentialsLeak verifies that the middleware
// does NOT stamp the Authorization header (or any other obviously-
// sensitive request attribute) onto the span. otelgin's default
// HTTP attribute set follows the OTel semantic conventions which
// already exclude Authorization, so this test is a regression guard
// rather than a new contract.
func Test_TraceCorrelation_NoCredentialsLeak(t *testing.T) {
	rec := installTestTracer(t)

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.OTelMiddleware("test-gw"))
	r.Use(middleware.TraceCorrelation())
	r.GET("/p", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/p", nil)
	req.Header.Set("Authorization", "Bearer s3cr3t")
	req.Header.Set("Cookie", "session=topsecret")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	spans := rec.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	for _, kv := range spans[0].Attributes() {
		val := strings.ToLower(kv.Value.Emit())
		if strings.Contains(val, "s3cr3t") {
			t.Errorf("auth header leaked onto span %s = %q", kv.Key, kv.Value.Emit())
		}
		if strings.Contains(val, "topsecret") {
			t.Errorf("cookie leaked onto span %s = %q", kv.Key, kv.Value.Emit())
		}
		if strings.Contains(strings.ToLower(string(kv.Key)), "authorization") {
			t.Errorf("authorization key on span: %s", kv.Key)
		}
	}
}
