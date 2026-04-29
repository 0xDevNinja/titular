package graph

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/store"
)

func init() { gin.SetMode(gin.TestMode) }

// newRouterWithGraphQL mounts a Handler onto a fresh Gin engine and
// returns it as an http.Handler so tests can drive POST /graphql via
// httptest. We deliberately do not wire the parent middleware chain
// here — those have separate tests, and including them would mask
// the GraphQL handler's own behaviour.
func newRouterWithGraphQL(t *testing.T, fs *fakeStore) http.Handler {
	t.Helper()
	h := newTestHandler(t, fs)
	r := gin.New()
	h.Register(r.Group(""))
	return r
}

func postQuery(t *testing.T, h http.Handler, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestHandler_Query_OK(t *testing.T) {
	fs := &fakeStore{stats: store.Stats{TotalAgents: 5}}
	h := newRouterWithGraphQL(t, fs)
	rec := postQuery(t, h, `{"query":"{ stats { totalAgents } }"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d", rec.Code)
	}
	var resp queryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Errors) > 0 {
		t.Fatalf("errors: %v", resp.Errors)
	}
	stats := resp.Data.(map[string]any)["stats"].(map[string]any)
	if stats["totalAgents"] != float64(5) {
		t.Errorf("totalAgents: got %v", stats["totalAgents"])
	}
}

func TestHandler_Query_BadJSON(t *testing.T) {
	h := newRouterWithGraphQL(t, &fakeStore{})
	rec := postQuery(t, h, `not json`)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", rec.Code)
	}
}

func TestHandler_Query_EmptyQuery(t *testing.T) {
	h := newRouterWithGraphQL(t, &fakeStore{})
	rec := postQuery(t, h, `{"query":""}`)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", rec.Code)
	}
}

// TestHandler_Query_RejectsSubscriptionOverPOST guards the contract
// that subscriptions only work over the websocket transport. A client
// sending `subscription { newTrade { id } }` over POST gets 400.
func TestHandler_Query_RejectsSubscriptionOverPOST(t *testing.T) {
	h := newRouterWithGraphQL(t, &fakeStore{})
	rec := postQuery(t, h, `{"query":"subscription { newTrade { id } }"}`)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "websocket") {
		t.Errorf("body: %s", rec.Body.String())
	}
}

// TestHandler_Query_OversizeBody pins the 1 MiB limit. We construct a
// body just over the cap and assert the request is rejected.
func TestHandler_Query_OversizeBody(t *testing.T) {
	h := newRouterWithGraphQL(t, &fakeStore{})
	big := strings.Repeat("x", maxQueryBytes+1)
	body := `{"query":"# ` + big + `\n{ stats { totalAgents } }"}`
	rec := postQuery(t, h, body)
	if rec.Code == http.StatusOK {
		t.Errorf("expected non-200, got %d", rec.Code)
	}
}

// TestHandler_Subscribe_RejectsNonWebsocket asserts a plain GET on the
// subscriptions path returns 400 with a human-friendly message rather
// than hanging or 500ing.
func TestHandler_Subscribe_RejectsNonWebsocket(t *testing.T) {
	h := newRouterWithGraphQL(t, &fakeStore{})
	req := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", rec.Code)
	}
}

// TestHandler_Playground_OptIn ensures the playground is NOT mounted
// by default. SDK consumers expect `GET /graphql` to be the websocket
// upgrade path; a misconfigured production deploy that left the
// playground enabled is a security smell we want to catch.
func TestHandler_Playground_OptIn(t *testing.T) {
	fs := &fakeStore{}
	h := newTestHandler(t, fs)
	r := gin.New()
	h.Register(r.Group(""))
	req := httptest.NewRequest(http.MethodGet, "/graphql/playground", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("playground reachable when disabled: got %d", rec.Code)
	}
}

func TestHandler_Playground_EnabledServesHTML(t *testing.T) {
	fs := &fakeStore{}
	h := newTestHandler(t, fs)
	h.EnablePlayground = true
	r := gin.New()
	h.Register(r.Group(""))
	req := httptest.NewRequest(http.MethodGet, "/graphql/playground", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("status: got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "graphiql") {
		t.Errorf("body does not look like GraphiQL UI")
	}
}

// TestSanitiseGraphQLErr_DoesNotLeak is the explicit guardrail that
// any new internal-error message added to the resolver layer goes
// through the sanitiser. Mirrors the equivalent test in handlers
// package (TestAPIStats_StoreErrIsOpaque) for the GraphQL surface.
func TestSanitiseGraphQLErr_DoesNotLeak(t *testing.T) {
	for _, msg := range []string{
		`pq: relation "agents" does not exist`,
		`unexpected EOF`,
		`canceled context`,
		`some random panic recovered text`,
	} {
		got := sanitiseGraphQLErr(msg)
		if got != "internal error" {
			t.Errorf("sanitiseGraphQLErr(%q) = %q, want \"internal error\"", msg, got)
		}
	}
	// Allow-list survivors:
	for _, msg := range []string{
		"cursor is not valid base64",
		"agent_token: must start with 0x",
		"Cannot query field foo on type Bar",
		"Syntax Error GraphQL request",
	} {
		got := sanitiseGraphQLErr(msg)
		if got != msg {
			t.Errorf("sanitiseGraphQLErr(%q) = %q, want passthrough", msg, got)
		}
	}
}
