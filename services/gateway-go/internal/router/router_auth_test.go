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

// authFixtureRouter wires the production router with a real auth.Handlers
// backed by miniredis. We use the production router here on purpose so the
// auth/v1 ordering (auth must not be shadowed by /v1/*) is locked in by an
// end-to-end exercise rather than a unit test.
func authFixtureRouter(t *testing.T) (http.Handler, *miniredis.Miniredis) {
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
	h, err := auth.NewHandlers(auth.HandlerConfig{
		Store:   auth.NewRedisStore(c),
		Signer:  signer,
		Domain:  "titular.test",
		ChainID: 84532,
	})
	if err != nil {
		t.Fatalf("NewHandlers: %v", err)
	}

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(file), "..", "handlers", "fixtures")
	store, err := handlers.NewStoreFromDir(dir)
	if err != nil {
		t.Fatalf("load fixture store: %v", err)
	}
	ah := handlers.NewAgentHandlersWithStore(store)
	jh := handlers.NewJobHandlersWithStore(store)

	cfg := router.DefaultConfig()
	cfg.Auth = h
	return router.NewWithConfig(cfg, ah, jh), mr
}

func TestRouter_AuthNonceMounted(t *testing.T) {
	h, _ := authFixtureRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/auth/siwe/nonce", nil)
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

// Auth routes must be reachable independently of /v1/*: a regression where
// the chi catch-all swallowed /auth would surface as a 404 here.
func TestRouter_AuthRoutesNotShadowedByV1(t *testing.T) {
	h, _ := authFixtureRouter(t)

	for _, p := range []string{"/auth/siwe/nonce", "/auth/siwe/verify", "/auth/logout"} {
		req := httptest.NewRequest(http.MethodPost, p, nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code == http.StatusNotFound {
			t.Errorf("path %s returned 404; mount order regressed", p)
		}
	}
}

// When Auth is unset, the /auth routes must NOT be mounted. The default
// router.DefaultConfig() must continue to behave as it always has so the
// existing M2/M3 fixture tests stay valid.
func TestRouter_AuthOptional(t *testing.T) {
	cfg := router.DefaultConfig()
	cfg.Auth = nil

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(file), "..", "handlers", "fixtures")
	store, err := handlers.NewStoreFromDir(dir)
	if err != nil {
		t.Fatalf("load fixture store: %v", err)
	}
	h := router.NewWithConfig(
		cfg,
		handlers.NewAgentHandlersWithStore(store),
		handlers.NewJobHandlersWithStore(store),
	)
	req := httptest.NewRequest(http.MethodPost, "/auth/siwe/nonce", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("auth unset must 404 on /auth/*: got %d", rec.Code)
	}
}
