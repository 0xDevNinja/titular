package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/graph"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/router"
)

// router_graphql_test exercises the chassis side of the GraphQL mount:
// is /graphql wired when cfg.GraphQL is non-nil, and is it absent
// otherwise? The schema-level behaviour is covered by graph/*_test.go.

// stubGraphQLStore is the minimum-viable Store needed to construct the
// GraphQL Handler. It returns empty pages for every list call and
// ErrNotFound for the per-id lookup; the only thing the router test
// cares about is that the endpoint is mounted.
type stubGraphQLStore struct{ stubAPIStore }

func TestRouter_GraphQLMounted(t *testing.T) {
	resolver := &graph.Resolver{Store: stubGraphQLStore{}, Bus: graph.NilBus()}
	schema, err := graph.NewSchema(resolver)
	if err != nil {
		t.Fatalf("NewSchema: %v", err)
	}
	h, err := graph.NewHandler(schema, resolver)
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}

	cfg := router.DefaultConfig()
	cfg.GraphQL = h
	srv := fixtureRouter(t, cfg)

	body := []byte(`{"query":"{ stats { totalAgents } }"}`)
	req := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, body=%s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	stats := resp.Data["stats"].(map[string]any)
	if _, ok := stats["totalAgents"]; !ok {
		t.Errorf("totalAgents missing: %+v", stats)
	}
}

func TestRouter_GraphQLAbsentWhenNil(t *testing.T) {
	srv := fixtureRouter(t, router.DefaultConfig())
	req := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewBufferString(`{"query":"{stats{totalAgents}}"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 when GraphQL not mounted, got %d", rec.Code)
	}
}
