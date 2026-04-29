package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/middleware"
)

func newCORSEngine(cfg middleware.CORSConfig) *gin.Engine {
	r := gin.New()
	r.Use(middleware.CORS(cfg))
	r.GET("/x", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	r.POST("/x", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	return r
}

func TestCORS_AllowsConfiguredOrigin(t *testing.T) {
	cfg := middleware.DefaultCORSConfig()
	cfg.AllowedOrigins = []string{"https://app.titular.xyz"}
	r := newCORSEngine(cfg)

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Origin", "https://app.titular.xyz")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://app.titular.xyz" {
		t.Fatalf("Allow-Origin: got %q, want https://app.titular.xyz", got)
	}
	if got := rec.Header().Get("Vary"); !strings.Contains(got, "Origin") {
		t.Fatalf("Vary should include Origin, got %q", got)
	}
}

func TestCORS_DoesNotAllowUnconfiguredOrigin(t *testing.T) {
	cfg := middleware.DefaultCORSConfig()
	cfg.AllowedOrigins = []string{"https://app.titular.xyz"}
	r := newCORSEngine(cfg)

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Origin", "https://evil.example.com")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no Allow-Origin for disallowed origin, got %q", got)
	}
}

func TestCORS_HandlesPreflight(t *testing.T) {
	cfg := middleware.DefaultCORSConfig()
	cfg.AllowedOrigins = []string{"https://app.titular.xyz"}
	r := newCORSEngine(cfg)

	req := httptest.NewRequest(http.MethodOptions, "/x", nil)
	req.Header.Set("Origin", "https://app.titular.xyz")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	req.Header.Set("Access-Control-Request-Headers", "content-type")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status: got %d, want 204", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); !strings.Contains(got, http.MethodPost) {
		t.Fatalf("Allow-Methods should include POST, got %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Headers"); got == "" {
		t.Fatal("Allow-Headers should be set on preflight")
	}
	if got := rec.Header().Get("Access-Control-Max-Age"); got == "" {
		t.Fatal("Max-Age should be set on preflight")
	}
}

func TestCORS_WildcardWithoutCredentials(t *testing.T) {
	cfg := middleware.DefaultCORSConfig()
	cfg.AllowedOrigins = []string{"*"}
	cfg.AllowCredentials = false
	r := newCORSEngine(cfg)

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Origin", "https://anything.example.com")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("Allow-Origin: got %q, want *", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "" {
		t.Fatalf("Allow-Credentials must not be set with wildcard, got %q", got)
	}
}

func TestCORS_WildcardWithCredentialsEchoesOrigin(t *testing.T) {
	cfg := middleware.DefaultCORSConfig()
	cfg.AllowedOrigins = []string{"*"}
	cfg.AllowCredentials = true
	// Opt in to the credentialed-reflector mode so CORS does not refuse the
	// configuration. Production deployments should not enable this.
	cfg.AllowReflectedOrigins = true
	r := newCORSEngine(cfg)

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Origin", "https://app.example.com")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example.com" {
		t.Fatalf("Allow-Origin: got %q, want echoed origin", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("Allow-Credentials: got %q, want true", got)
	}
}

func TestCORS_RejectsWildcardWithCredentialsByDefault(t *testing.T) {
	cfg := middleware.DefaultCORSConfig()
	cfg.AllowedOrigins = []string{"*"}
	cfg.AllowCredentials = true
	// AllowReflectedOrigins left at its zero value (false).

	if _, err := middleware.NewCORS(cfg); err == nil {
		t.Fatal("expected NewCORS to refuse wildcard+credentials, got nil err")
	} else if !errors.Is(err, middleware.ErrCORSWildcardWithCredentials) {
		t.Fatalf("expected ErrCORSWildcardWithCredentials, got %v", err)
	}

	// CORS (the panicking variant) should also refuse — recover and assert.
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected CORS to panic on unsafe config, got no panic")
		}
	}()
	_ = middleware.CORS(cfg)
}

func TestCORS_NoOriginHeaderPassesThrough(t *testing.T) {
	cfg := middleware.DefaultCORSConfig()
	cfg.AllowedOrigins = []string{"*"}
	r := newCORSEngine(cfg)

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("Allow-Origin should be empty when no Origin sent, got %q", got)
	}
}
