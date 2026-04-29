package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/auth"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/handlers"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/router"
)

// siwsFixtureRouter wires the production router with both SIWE and SIWS
// handler bundles backed by miniredis. Mirrors authFixtureRouter so the
// SIWS tests don't depend on production env wiring.
func siwsFixtureRouter(t *testing.T) (http.Handler, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	c := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = c.Close() })

	signer, err := auth.NewJWTSigner(auth.SignerConfig{
		Secret: []byte(strings.Repeat("k", 32)),
		Issuer: "gateway-test",
		TTL:    time.Hour,
	})
	if err != nil {
		t.Fatalf("NewJWTSigner: %v", err)
	}
	store := auth.NewRedisStore(c)
	siwe, err := auth.NewHandlers(auth.HandlerConfig{
		Store:   store,
		Signer:  signer,
		Domain:  "titular.test",
		ChainID: 84532,
	})
	if err != nil {
		t.Fatalf("NewHandlers: %v", err)
	}
	siws, err := auth.NewSIWSHandlers(auth.SIWSHandlerConfig{
		Store:   store,
		Signer:  signer,
		Domain:  "titular.test",
		Cluster: auth.ClusterDevnet,
	})
	if err != nil {
		t.Fatalf("NewSIWSHandlers: %v", err)
	}

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(file), "..", "handlers", "fixtures")
	hs, err := handlers.NewStoreFromDir(dir)
	if err != nil {
		t.Fatalf("load fixture store: %v", err)
	}
	ah := handlers.NewAgentHandlersWithStore(hs)
	jh := handlers.NewJobHandlersWithStore(hs)

	cfg := router.DefaultConfig()
	cfg.Auth = siwe
	cfg.SIWS = siws
	return router.NewWithConfig(cfg, ah, jh), mr
}

func TestRouter_SIWSNonceMounted(t *testing.T) {
	h, _ := siwsFixtureRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/auth/siws/nonce", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, body=%s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, ok := body["nonce"].(string); !ok {
		t.Fatalf("response missing nonce: %v", body)
	}
}

// SIWS routes must coexist with SIWE routes. A regression where one mount
// shadowed the other would surface as a 404.
func TestRouter_SIWSRoutesNotShadowed(t *testing.T) {
	h, _ := siwsFixtureRouter(t)

	for _, p := range []string{
		"/auth/siwe/nonce", "/auth/siwe/verify", "/auth/logout",
		"/auth/siws/nonce", "/auth/siws/verify",
	} {
		req := httptest.NewRequest(http.MethodPost, p, nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code == http.StatusNotFound {
			t.Errorf("path %s returned 404; mount regressed", p)
		}
	}
}

// When SIWS is unset, /auth/siws/* must NOT be mounted, even if SIWE is
// enabled. This pins the additive nature of the new flag.
func TestRouter_SIWSOptional(t *testing.T) {
	mr := miniredis.RunT(t)
	c := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = c.Close() })

	signer, _ := auth.NewJWTSigner(auth.SignerConfig{
		Secret: []byte(strings.Repeat("k", 32)),
		Issuer: "gateway-test",
		TTL:    time.Hour,
	})
	siwe, _ := auth.NewHandlers(auth.HandlerConfig{
		Store:   auth.NewRedisStore(c),
		Signer:  signer,
		Domain:  "titular.test",
		ChainID: 84532,
	})

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(file), "..", "handlers", "fixtures")
	hs, _ := handlers.NewStoreFromDir(dir)

	cfg := router.DefaultConfig()
	cfg.Auth = siwe
	cfg.SIWS = nil

	h := router.NewWithConfig(
		cfg,
		handlers.NewAgentHandlersWithStore(hs),
		handlers.NewJobHandlersWithStore(hs),
	)

	// SIWE present.
	req := httptest.NewRequest(http.MethodPost, "/auth/siwe/nonce", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("siwe nonce regressed: %d", rec.Code)
	}

	// SIWS absent.
	req = httptest.NewRequest(http.MethodPost, "/auth/siws/nonce", nil)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("siws unset must 404 on /auth/siws/*: got %d", rec.Code)
	}
}
