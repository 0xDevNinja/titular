package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	gatewaydocs "github.com/0xDevNinja/titular/services/gateway-go/docs"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/openapi"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/router"
)

// TestRouter_MountsOpenAPISpecWhenConfigured asserts that wiring an
// openapi.Spec into router.Config exposes /swagger/doc.json. This is
// the integration-level guarantee corresponding to the unit tests in
// internal/openapi: it verifies the router actually invokes
// Spec.Register, which is the only place a regression in mount order
// would surface.
func TestRouter_MountsOpenAPISpecWhenConfigured(t *testing.T) {
	t.Parallel()

	spec, err := openapi.NewSpec(gatewaydocs.SpecJSON, gatewaydocs.SpecYAML, false)
	if err != nil {
		t.Fatalf("new spec: %v", err)
	}

	cfg := router.DefaultConfig()
	cfg.OpenAPI = spec

	h := fixtureRouter(t, cfg)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want 200 once OpenAPI bundle is wired", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got == "" {
		t.Error("missing Content-Type on /swagger/doc.json")
	}
}

// TestRouter_OmitsOpenAPIWhenNotConfigured is the symmetric guarantee
// — if a deployment chooses not to wire the bundle, /swagger/* must be
// absent. This matches the behaviour of the GraphQL / SSE / API
// surfaces and gives operators a single uniform "feature off" knob.
func TestRouter_OmitsOpenAPIWhenNotConfigured(t *testing.T) {
	t.Parallel()

	cfg := router.DefaultConfig()
	cfg.OpenAPI = nil

	h := fixtureRouter(t, cfg)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d; want 404 when OpenAPI is absent", rec.Code)
	}
}
