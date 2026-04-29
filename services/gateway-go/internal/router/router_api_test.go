package router_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/handlers"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/router"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/store"
)

// stubAPIStore is a Store that always succeeds. We only use it to confirm
// that the /api/v1 group is mounted (and absent when API is nil); the
// per-endpoint behaviour is exercised by handlers/api_test.go via a richer
// fake.
type stubAPIStore struct{}

func (stubAPIStore) ListAgents(context.Context, store.ListAgentsParams) (store.Page[store.Agent], error) {
	return store.Page[store.Agent]{}, nil
}
func (stubAPIStore) GetAgent(context.Context, string) (store.Agent, error) {
	return store.Agent{}, store.ErrNotFound
}
func (stubAPIStore) ListTrades(context.Context, store.ListTradesParams) (store.Page[store.Trade], error) {
	return store.Page[store.Trade]{}, nil
}
func (stubAPIStore) ListJobs(context.Context, store.ListJobsParams) (store.Page[store.Job], error) {
	return store.Page[store.Job]{}, nil
}
func (stubAPIStore) Stats(context.Context) (store.Stats, error) {
	return store.Stats{}, nil
}

func TestRouter_APIv1Mounted(t *testing.T) {
	cfg := router.DefaultConfig()
	cfg.API = handlers.NewAPI(stubAPIStore{})
	h := fixtureRouter(t, cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
}

func TestRouter_APIv1AbsentWhenNil(t *testing.T) {
	// When cfg.API is nil, GET /api/v1/stats should fall through to a 404.
	h := fixtureRouter(t, router.DefaultConfig())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: got %d, want 404", rec.Code)
	}
}
