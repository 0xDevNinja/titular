package auth_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/auth"
)

// protectedEngine wires RequireAuth in front of a single echo route that
// reflects the resolved address back to the caller. We use this rather than
// the production router to keep the middleware test surface minimal.
func protectedEngine(h *auth.Handlers) *gin.Engine {
	r := gin.New()
	r.GET("/me", auth.RequireAuth(h), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"address": auth.AddressFromContext(c),
			"jti":     auth.JTIFromContext(c),
		})
	})
	return r
}

func TestRequireAuth_NoHeader_401(t *testing.T) {
	h, _, _ := newHandlers(t, "titular.test", 84532)
	r := protectedEngine(h)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: got %d, want 401", rec.Code)
	}
}

func TestRequireAuth_GarbageToken_401(t *testing.T) {
	h, _, _ := newHandlers(t, "titular.test", 84532)
	r := protectedEngine(h)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-jwt")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: got %d, want 401", rec.Code)
	}
}

func TestRequireAuth_ValidTokenButNoSession_401(t *testing.T) {
	// Mint a token without seeding the matching session in Redis. RequireAuth
	// must reject because the server-side session is the source of truth.
	h, _, signer := newHandlers(t, "titular.test", 84532)
	r := protectedEngine(h)

	tok, _, err := signer.Sign("0xabc")
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: got %d, want 401", rec.Code)
	}
}

func TestRequireAuth_HappyPath_PublishesIdentity(t *testing.T) {
	h, store, signer := newHandlers(t, "titular.test", 84532)
	r := protectedEngine(h)

	tok, jti, err := signer.Sign("0xABC")
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if err := store.PutSession(t.Context(), jti, "0xabc", time.Hour); err != nil {
		t.Fatalf("PutSession: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"address":"0xabc"`) {
		t.Fatalf("body missing address: %s", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"jti":"`+jti+`"`) {
		t.Fatalf("body missing jti: %s", rec.Body.String())
	}
}

func TestRequireAuth_LogoutRevokes(t *testing.T) {
	// Mint -> use successfully -> logout -> reuse must 401. This is the
	// session-backed revocation behaviour that motivates the Redis check.
	h, store, signer := newHandlers(t, "titular.test", 84532)
	r := protectedEngine(h)

	tok, jti, err := signer.Sign("0xabc")
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	_ = store.PutSession(t.Context(), jti, "0xabc", time.Hour)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("warmup: got %d", rec.Code)
	}

	if err := store.DeleteSession(t.Context(), jti); err != nil {
		t.Fatalf("DeleteSession: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("post-logout reuse: got %d, want 401", rec.Code)
	}
}

func TestRequireAuth_MalformedAuthHeader_401(t *testing.T) {
	h, _, _ := newHandlers(t, "titular.test", 84532)
	r := protectedEngine(h)

	cases := []string{
		"",
		"Bearer",
		"Token abc",
		"basic dXNlcjpwYXNz",
	}
	for _, hdr := range cases {
		hdr := hdr
		t.Run("hdr="+hdr, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/me", nil)
			if hdr != "" {
				req.Header.Set("Authorization", hdr)
			}
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("hdr=%q: got %d, want 401", hdr, rec.Code)
			}
		})
	}
}

func TestAddressFromContext_AbsentReturnsEmpty(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	if got := auth.AddressFromContext(c); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
	if got := auth.JTIFromContext(c); got != "" {
		t.Fatalf("expected empty jti, got %q", got)
	}
}
