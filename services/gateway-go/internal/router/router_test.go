package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/handlers"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/middleware"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/router"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// fixtureRouter returns a router wired to the M2 fixture store. We reuse the
// fixture loader rather than copy-pasting JSON; the goal of these tests is to
// exercise the chassis (middleware + Gin engine + chi sub-mount), not the
// fixture content.
func fixtureRouter(t *testing.T, cfg router.Config) http.Handler {
	t.Helper()
	_, file, _, _ := runtime.Caller(0)
	// internal/router/router_test.go -> ../handlers/fixtures
	dir := filepath.Join(filepath.Dir(file), "..", "handlers", "fixtures")

	store, err := handlers.NewStoreFromDir(dir)
	if err != nil {
		t.Fatalf("load fixture store: %v", err)
	}
	ah := handlers.NewAgentHandlersWithStore(store)
	jh := handlers.NewJobHandlersWithStore(store)
	return router.NewWithConfig(cfg, ah, jh)
}

func TestRouter_HealthzReturnsOK(t *testing.T) {
	h := fixtureRouter(t, router.DefaultConfig())

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status: got %q, want ok", body["status"])
	}
}

func TestRouter_RequestIDHeaderSet(t *testing.T) {
	h := fixtureRouter(t, router.DefaultConfig())

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	got := rec.Header().Get(middleware.RequestIDHeader)
	if got == "" {
		t.Fatal("expected X-Request-Id on response, got empty")
	}
	if _, err := uuid.Parse(got); err != nil {
		t.Fatalf("expected UUID, got %q: %v", got, err)
	}
}

func TestRouter_LogsHaveRequestID(t *testing.T) {
	buf := &bytes.Buffer{}
	cfg := router.DefaultConfig()
	cfg.Logger = zerolog.New(buf)

	h := fixtureRouter(t, cfg)

	// Hit a non-skipped path so a log line is emitted (the default config
	// skips /healthz).
	req := httptest.NewRequest(http.MethodGet, "/v1/agents", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}

	line := strings.TrimRight(buf.String(), "\n")
	if line == "" {
		t.Fatal("expected a log line, got none")
	}
	var rec2 map[string]any
	if err := json.Unmarshal([]byte(line), &rec2); err != nil {
		t.Fatalf("decode log line %q: %v", line, err)
	}
	if id, ok := rec2["request_id"].(string); !ok || id == "" {
		t.Fatalf("expected request_id in log, got %v", rec2["request_id"])
	}
	hdr := rec.Header().Get(middleware.RequestIDHeader)
	if rec2["request_id"] != hdr {
		t.Fatalf("log id %v != header id %q", rec2["request_id"], hdr)
	}
}

func TestRouter_LegacyV1AgentsStillReachable(t *testing.T) {
	h := fixtureRouter(t, router.DefaultConfig())

	req := httptest.NewRequest(http.MethodGet, "/v1/agents", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}
}

func TestRouter_CORSHeadersOnPreflight(t *testing.T) {
	cfg := router.DefaultConfig()
	cfg.CORS.AllowedOrigins = []string{"https://app.titular.xyz"}

	h := fixtureRouter(t, cfg)

	req := httptest.NewRequest(http.MethodOptions, "/v1/agents", nil)
	req.Header.Set("Origin", "https://app.titular.xyz")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("preflight status: got %d, want 204", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://app.titular.xyz" {
		t.Fatalf("Allow-Origin: got %q", got)
	}
}

func TestRouter_RateLimitEnforced(t *testing.T) {
	cfg := router.DefaultConfig()
	cfg.RateLimit = middleware.RateLimitConfig{
		RPS:   0.0001,
		Burst: 1,
		// All requests in the test share the same fake key so the bucket
		// caps everything beyond the first hit.
		KeyFunc: func(_ *gin.Context) string { return "fixed" },
	}
	h := fixtureRouter(t, cfg)

	// First request: allowed.
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("warmup: got %d, want 200", rec.Code)
	}

	// Second request: rate-limited.
	req = httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("got %d, want 429", rec.Code)
	}
}

func TestRouter_RecoveryConvertsPanicTo500(t *testing.T) {
	// Inject a route that panics by hand-rolling a router with the same
	// middleware stack as production but a controlled handler.
	buf := &bytes.Buffer{}
	logger := zerolog.New(buf)

	r := gin.New()
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.RequestID())
	r.Use(middleware.Log(middleware.LogConfig{Logger: logger, Service: "gateway"}))
	r.GET("/boom", func(_ *gin.Context) { panic("explode") })

	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "internal_server_error") {
		t.Fatalf("body missing canonical code: %q", body)
	}
	if strings.Contains(body, "explode") {
		t.Fatalf("response leaked panic value: %q", body)
	}
}
