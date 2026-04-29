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
	// middleware order as production: RequestID → Recovery → Log. This
	// guarantees Recovery's panic log line carries the request id assigned
	// upstream by RequestID rather than logging an empty string.
	buf := &bytes.Buffer{}
	logger := zerolog.New(buf)

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Recovery(logger))
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
	// The X-Request-Id header must round-trip even on panic — assert that
	// Recovery did not strip the response header that RequestID set.
	if got := rec.Header().Get(middleware.RequestIDHeader); got == "" {
		t.Fatal("expected X-Request-Id header on 500 response, got empty")
	}

	// And every log line in the buffer must carry a non-empty request_id.
	for _, line := range strings.Split(strings.TrimRight(buf.String(), "\n"), "\n") {
		if line == "" {
			continue
		}
		var rec2 map[string]any
		if err := json.Unmarshal([]byte(line), &rec2); err != nil {
			t.Fatalf("decode log line %q: %v", line, err)
		}
		id, _ := rec2["request_id"].(string)
		if id == "" {
			t.Fatalf("log line missing request_id: %q", line)
		}
	}
}

// TestRouter_PrepareJobReachable locks in the chi-wrap contract for POST
// routes. The chi sub-router lives behind gin.WrapH, so the gin engine has
// to forward the body and method correctly for fixture POSTs to work.
func TestRouter_PrepareJobReachable(t *testing.T) {
	h := fixtureRouter(t, router.DefaultConfig())

	body := `{
        "agent_id": "agent-001",
        "metadata_uri": "ipfs://QmTest",
        "amount": "1000000000000000000",
        "deadline_seconds": 4102444800
    }`

	req := httptest.NewRequest(http.MethodPost, "/v1/jobs/prepare", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d (body=%q), want 200", rec.Code, rec.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	for _, k := range []string{"to", "calldata", "value"} {
		if _, ok := resp[k]; !ok {
			t.Errorf("response missing %q: %v", k, resp)
		}
	}
}

// TestRouter_GetAgentByIDReachable locks in the chi-wrap contract for path
// params. chi.URLParam reads the parameter from the request context — if the
// gin → chi handoff dropped the route context the handler would 404.
func TestRouter_GetAgentByIDReachable(t *testing.T) {
	h := fixtureRouter(t, router.DefaultConfig())

	req := httptest.NewRequest(http.MethodGet, "/v1/agents/agent-001", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d (body=%q), want 200", rec.Code, rec.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if resp["id"] != "agent-001" {
		t.Fatalf("id: got %v, want agent-001", resp["id"])
	}
}

// TestRouter_NoTrustedProxies_IgnoresXForwardedFor verifies the safe
// default: with no GATEWAY_TRUSTED_PROXIES set, an inbound X-Forwarded-For
// must NOT influence the rate-limit key (or any other ClientIP-derived
// behaviour). We exercise this by routing through a KeyFunc that captures
// what the limiter sees, then assert the captured value is the unspoofable
// RemoteAddr peer rather than the spoofed header value.
func TestRouter_NoTrustedProxies_IgnoresXForwardedFor(t *testing.T) {
	var captured string

	cfg := router.DefaultConfig()
	cfg.RateLimit = middleware.RateLimitConfig{
		RPS:   1000, // permissive so no 429 leaks into the assertion
		Burst: 1000,
		KeyFunc: func(c *gin.Context) string {
			captured = c.ClientIP()
			return captured
		},
	}
	// TrustedProxies left nil — the production-default that should refuse
	// to honour X-Forwarded-For.

	h := fixtureRouter(t, cfg)

	req := httptest.NewRequest(http.MethodGet, "/v1/agents", nil)
	req.RemoteAddr = "192.0.2.10:54321" // synthetic peer; must win
	req.Header.Set("X-Forwarded-For", "1.1.1.1")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}
	if captured == "" {
		t.Fatal("KeyFunc was never called; rate limiter wiring broken")
	}
	if captured == "1.1.1.1" {
		t.Fatalf("rate-limit key honoured X-Forwarded-For: got %q", captured)
	}
	// We do not pin the exact value (gin may strip the port) but it must be
	// the peer host, not the forwarded header.
	if !strings.Contains(captured, "192.0.2.10") {
		t.Fatalf("rate-limit key %q does not match peer 192.0.2.10", captured)
	}
}
