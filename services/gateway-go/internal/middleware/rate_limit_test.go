package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/middleware"
)

// fixedKey forces all requests onto the same bucket so the bucket capacity
// determines the response sequence regardless of httptest.ResponseRecorder's
// fake remote address.
func fixedKey(_ *gin.Context) string { return "test-key" }

func newRateLimitEngine(cfg middleware.RateLimitConfig) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RateLimit(cfg))
	r.GET("/x", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	return r
}

func TestRateLimit_AllowsBurstThenBlocks(t *testing.T) {
	r := newRateLimitEngine(middleware.RateLimitConfig{
		RPS:     0.0001, // effectively no refill within the test duration
		Burst:   2,
		KeyFunc: fixedKey,
	})

	doRequest := func() int {
		req := httptest.NewRequest(http.MethodGet, "/x", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		return rec.Code
	}

	// First two should pass — bucket has burst=2 tokens.
	for i := 0; i < 2; i++ {
		if got := doRequest(); got != http.StatusOK {
			t.Fatalf("request %d: got %d, want 200", i+1, got)
		}
	}
	// Third should be limited.
	if got := doRequest(); got != http.StatusTooManyRequests {
		t.Fatalf("third request: got %d, want 429", got)
	}
}

func TestRateLimit_429BodyAndHeaders(t *testing.T) {
	r := newRateLimitEngine(middleware.RateLimitConfig{
		RPS:     0.0001,
		Burst:   1,
		KeyFunc: fixedKey,
	})

	// First request consumes the only token.
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("warmup: got %d, want 200", rec.Code)
	}

	// Second is rate-limited.
	req = httptest.NewRequest(http.MethodGet, "/x", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("got %d, want 429", rec.Code)
	}
	if got := rec.Header().Get("Retry-After"); got == "" {
		t.Fatal("expected Retry-After header on 429")
	}
	body := rec.Body.String()
	if !contains(body, `"rate_limited"`) {
		t.Fatalf("body missing canonical error code: %q", body)
	}
}

func TestRateLimit_DistinctKeysAreIndependent(t *testing.T) {
	keys := []string{"alpha", "beta"}
	idx := 0
	keyFn := func(_ *gin.Context) string {
		k := keys[idx%len(keys)]
		idx++
		return k
	}

	r := gin.New()
	r.Use(middleware.RateLimit(middleware.RateLimitConfig{
		RPS:     0.0001,
		Burst:   1,
		KeyFunc: keyFn,
	}))
	r.GET("/x", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/x", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("alternating-keys request %d: got %d, want 200", i+1, rec.Code)
		}
	}
}

func TestRateLimit_DisabledByZeroConfig(t *testing.T) {
	r := newRateLimitEngine(middleware.RateLimitConfig{}) // RPS=0, Burst=0

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/x", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("disabled limiter blocked request %d (got %d)", i+1, rec.Code)
		}
	}
}

func TestRateLimit_EmptyKeyFailsOpen(t *testing.T) {
	r := gin.New()
	r.Use(middleware.RateLimit(middleware.RateLimitConfig{
		RPS:     0.0001,
		Burst:   1,
		KeyFunc: func(_ *gin.Context) string { return "" },
	}))
	r.GET("/x", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/x", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("empty-key request %d: got %d, want 200", i+1, rec.Code)
		}
	}
}

// contains is a tiny helper that avoids importing strings just for one call.
func contains(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
