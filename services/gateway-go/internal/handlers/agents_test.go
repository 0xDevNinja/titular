package handlers_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/handlers"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/router"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/types"
)

// testServer builds an httptest.Server wired to the real fixture files.
func testServer(t *testing.T) *httptest.Server {
	t.Helper()
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(file), "fixtures")

	store, err := handlers.NewStoreFromDir(dir)
	if err != nil {
		t.Fatalf("load fixture store: %v", err)
	}

	ah := handlers.NewAgentHandlersWithStore(store)
	jh := handlers.NewJobHandlersWithStore(store)
	srv := httptest.NewServer(router.NewWithHandlers(ah, jh))
	t.Cleanup(srv.Close)
	return srv
}

// get is a small helper that issues a GET and returns the response.
func get(t *testing.T, srv *httptest.Server, path string) *http.Response {
	t.Helper()
	resp, err := http.Get(srv.URL + path) //nolint:noctx // test helper only
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	return resp
}

// decodeJSON decodes resp.Body into out.
func decodeJSON(t *testing.T, resp *http.Response, out interface{}) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
}

// cursorFor returns the base64 cursor token for a given numeric offset,
// matching the internal encodeCursor format.
func cursorFor(offset int) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("offset:%d", offset)))
}

// ---------------------------------------------------------------------------
// GET /v1/agents
// ---------------------------------------------------------------------------

func TestListAgents_DefaultPage(t *testing.T) {
	srv := testServer(t)

	resp := get(t, srv, "/v1/agents")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body types.ListAgentsResponse
	decodeJSON(t, resp, &body)

	// Fixtures contain exactly 5 agents; default limit=20 so all fit on one page.
	if len(body.Agents) != 5 {
		t.Fatalf("expected 5 agents, got %d", len(body.Agents))
	}
	if body.NextCursor != "" {
		t.Fatalf("expected no next_cursor, got %q", body.NextCursor)
	}
}

func TestListAgents_LimitAndCursorPagination(t *testing.T) {
	srv := testServer(t)

	// Page 1: limit=3
	resp1 := get(t, srv, "/v1/agents?limit=3")
	if resp1.StatusCode != http.StatusOK {
		t.Fatalf("page 1: expected 200, got %d", resp1.StatusCode)
	}

	var page1 types.ListAgentsResponse
	decodeJSON(t, resp1, &page1)

	if len(page1.Agents) != 3 {
		t.Fatalf("page 1: expected 3 agents, got %d", len(page1.Agents))
	}
	if page1.NextCursor == "" {
		t.Fatal("page 1: expected next_cursor to be set")
	}
	if page1.NextCursor != cursorFor(3) {
		t.Fatalf("page 1: unexpected cursor %q", page1.NextCursor)
	}

	// Page 2: use cursor from page 1
	resp2 := get(t, srv, "/v1/agents?limit=3&cursor="+page1.NextCursor)
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("page 2: expected 200, got %d", resp2.StatusCode)
	}

	var page2 types.ListAgentsResponse
	decodeJSON(t, resp2, &page2)

	// 5 agents total, 3 on page 1 → 2 on page 2
	if len(page2.Agents) != 2 {
		t.Fatalf("page 2: expected 2 agents, got %d", len(page2.Agents))
	}
	if page2.NextCursor != "" {
		t.Fatalf("page 2: expected no next_cursor, got %q", page2.NextCursor)
	}

	// All IDs across both pages are distinct.
	seen := map[string]bool{}
	for _, a := range append(page1.Agents, page2.Agents...) {
		if seen[a.ID] {
			t.Fatalf("duplicate agent id %q across pages", a.ID)
		}
		seen[a.ID] = true
	}
	if len(seen) != 5 {
		t.Fatalf("expected 5 unique agents across 2 pages, got %d", len(seen))
	}
}

func TestListAgents_CursorBeyondEnd(t *testing.T) {
	srv := testServer(t)

	resp := get(t, srv, "/v1/agents?cursor="+cursorFor(100))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body types.ListAgentsResponse
	decodeJSON(t, resp, &body)

	if len(body.Agents) != 0 {
		t.Fatalf("expected 0 agents, got %d", len(body.Agents))
	}
}

func TestListAgents_InvalidLimit(t *testing.T) {
	srv := testServer(t)

	cases := []struct {
		name  string
		param string
	}{
		{"zero", "limit=0"},
		{"over_max", "limit=101"},
		{"non_numeric", "limit=abc"},
		{"negative", "limit=-1"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			resp := get(t, srv, "/v1/agents?"+tc.param)
			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", resp.StatusCode)
			}
			var env types.ErrorEnvelope
			decodeJSON(t, resp, &env)
			if env.Error.Code != types.ErrBadRequest {
				t.Errorf("expected code %q, got %q", types.ErrBadRequest, env.Error.Code)
			}
		})
	}
}

func TestListAgents_InvalidCursor(t *testing.T) {
	srv := testServer(t)

	cases := []struct {
		name   string
		cursor string
	}{
		{"not_base64", "notvalidbase64!!!"},
		{"wrong_format", base64.StdEncoding.EncodeToString([]byte("badformat"))},
		{"negative_offset", base64.StdEncoding.EncodeToString([]byte("offset:-5"))},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			resp := get(t, srv, "/v1/agents?cursor="+tc.cursor)
			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", resp.StatusCode)
			}
			var env types.ErrorEnvelope
			decodeJSON(t, resp, &env)
			if env.Error.Code != types.ErrBadRequest {
				t.Errorf("expected code %q, got %q", types.ErrBadRequest, env.Error.Code)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// GET /v1/agents/:id
// ---------------------------------------------------------------------------

func TestGetAgent_ByID(t *testing.T) {
	srv := testServer(t)

	resp := get(t, srv, "/v1/agents/agent-001")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var agent types.Agent
	decodeJSON(t, resp, &agent)

	if agent.ID != "agent-001" {
		t.Fatalf("expected id %q, got %q", "agent-001", agent.ID)
	}
	if agent.Slug != "sigma-7" {
		t.Fatalf("expected slug %q, got %q", "sigma-7", agent.Slug)
	}
}

func TestGetAgent_BySlug(t *testing.T) {
	srv := testServer(t)

	resp := get(t, srv, "/v1/agents/nova-prime")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var agent types.Agent
	decodeJSON(t, resp, &agent)

	if agent.ID != "agent-002" {
		t.Fatalf("expected id %q, got %q", "agent-002", agent.ID)
	}
	if agent.Slug != "nova-prime" {
		t.Fatalf("expected slug %q, got %q", "nova-prime", agent.Slug)
	}
}

func TestGetAgent_GraduatedFields(t *testing.T) {
	srv := testServer(t)

	// nova-prime (agent-002) is graduated in fixtures
	resp := get(t, srv, "/v1/agents/nova-prime")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var agent types.Agent
	decodeJSON(t, resp, &agent)

	if !agent.Curve.Graduated {
		t.Fatal("expected graduated=true for nova-prime")
	}
}

func TestGetAgent_NotFound(t *testing.T) {
	srv := testServer(t)

	resp := get(t, srv, "/v1/agents/does-not-exist")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}

	var env types.ErrorEnvelope
	decodeJSON(t, resp, &env)

	if env.Error.Code != types.ErrNotFound {
		t.Fatalf("expected code %q, got %q", types.ErrNotFound, env.Error.Code)
	}
	if env.Error.Message == "" {
		t.Fatal("expected non-empty error message")
	}
}

// ---------------------------------------------------------------------------
// GET /v1/agents/:id/trades
// ---------------------------------------------------------------------------

func TestListAgentTrades_FilterByAgent(t *testing.T) {
	srv := testServer(t)

	resp := get(t, srv, "/v1/agents/agent-001/trades")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body types.ListTradesResponse
	decodeJSON(t, resp, &body)

	if len(body.Trades) == 0 {
		t.Fatal("expected trades for agent-001, got none")
	}
	for _, tr := range body.Trades {
		if tr.AgentID != "agent-001" {
			t.Fatalf("unexpected agent_id %q in trades for agent-001", tr.AgentID)
		}
	}
}

func TestListAgentTrades_FilterByAgentSlug(t *testing.T) {
	srv := testServer(t)

	// Access trades via slug
	resp := get(t, srv, "/v1/agents/flux-dao/trades")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body types.ListTradesResponse
	decodeJSON(t, resp, &body)

	if len(body.Trades) == 0 {
		t.Fatal("expected trades for flux-dao, got none")
	}
	for _, tr := range body.Trades {
		if tr.AgentID != "agent-003" {
			t.Fatalf("unexpected agent_id %q for slug flux-dao", tr.AgentID)
		}
	}
}

func TestListAgentTrades_NoMixBetweenAgents(t *testing.T) {
	srv := testServer(t)

	// Verify trades for agent-001 and agent-002 are disjoint.
	resp1 := get(t, srv, "/v1/agents/agent-001/trades")
	resp2 := get(t, srv, "/v1/agents/agent-002/trades")

	var body1, body2 types.ListTradesResponse
	decodeJSON(t, resp1, &body1)
	decodeJSON(t, resp2, &body2)

	hashes1 := map[string]bool{}
	for _, tr := range body1.Trades {
		hashes1[tr.TxHash] = true
	}
	for _, tr := range body2.Trades {
		if hashes1[tr.TxHash] {
			t.Fatalf("tx_hash %q appears in both agent-001 and agent-002 trades", tr.TxHash)
		}
	}
}

func TestListAgentTrades_Pagination(t *testing.T) {
	srv := testServer(t)

	// agent-001 has 8 trades in fixtures; fetch with limit=5
	resp1 := get(t, srv, "/v1/agents/agent-001/trades?limit=5")
	if resp1.StatusCode != http.StatusOK {
		t.Fatalf("page 1: expected 200, got %d", resp1.StatusCode)
	}

	var page1 types.ListTradesResponse
	decodeJSON(t, resp1, &page1)

	if len(page1.Trades) != 5 {
		t.Fatalf("page 1: expected 5 trades, got %d", len(page1.Trades))
	}
	if page1.NextCursor == "" {
		t.Fatal("page 1: expected next_cursor")
	}

	// Page 2
	resp2 := get(t, srv, "/v1/agents/agent-001/trades?limit=5&cursor="+page1.NextCursor)
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("page 2: expected 200, got %d", resp2.StatusCode)
	}

	var page2 types.ListTradesResponse
	decodeJSON(t, resp2, &page2)

	// 8 trades total, 5 on page 1 → 3 on page 2
	if len(page2.Trades) != 3 {
		t.Fatalf("page 2: expected 3 trades, got %d", len(page2.Trades))
	}
	if page2.NextCursor != "" {
		t.Fatalf("page 2: expected no next_cursor, got %q", page2.NextCursor)
	}

	// No duplicate tx hashes across pages.
	seen := map[string]bool{}
	for _, tr := range append(page1.Trades, page2.Trades...) {
		if seen[tr.TxHash] {
			t.Fatalf("duplicate tx_hash %q across pages", tr.TxHash)
		}
		seen[tr.TxHash] = true
	}
}

func TestListAgentTrades_AgentNotFound(t *testing.T) {
	srv := testServer(t)

	resp := get(t, srv, "/v1/agents/ghost/trades")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}

	var env types.ErrorEnvelope
	decodeJSON(t, resp, &env)

	if env.Error.Code != types.ErrNotFound {
		t.Fatalf("expected code %q, got %q", types.ErrNotFound, env.Error.Code)
	}
}

// ---------------------------------------------------------------------------
// Trades bad query params
// ---------------------------------------------------------------------------

func TestListAgentTrades_InvalidLimit(t *testing.T) {
	srv := testServer(t)

	cases := []struct {
		name  string
		param string
	}{
		{"zero", "limit=0"},
		{"over_max", "limit=101"},
		{"non_numeric", "limit=abc"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			resp := get(t, srv, "/v1/agents/agent-001/trades?"+tc.param)
			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", resp.StatusCode)
			}
			var env types.ErrorEnvelope
			decodeJSON(t, resp, &env)
			if env.Error.Code != types.ErrBadRequest {
				t.Errorf("expected code %q, got %q", types.ErrBadRequest, env.Error.Code)
			}
		})
	}
}

func TestListAgentTrades_InvalidCursor(t *testing.T) {
	srv := testServer(t)

	resp := get(t, srv, "/v1/agents/agent-001/trades?cursor=!!!notbase64")
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
	var env types.ErrorEnvelope
	decodeJSON(t, resp, &env)
	if env.Error.Code != types.ErrBadRequest {
		t.Fatalf("expected code %q, got %q", types.ErrBadRequest, env.Error.Code)
	}
}

func TestListAgentTrades_CursorBeyondEnd(t *testing.T) {
	srv := testServer(t)

	resp := get(t, srv, "/v1/agents/agent-001/trades?cursor="+cursorFor(500))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var body types.ListTradesResponse
	decodeJSON(t, resp, &body)
	if len(body.Trades) != 0 {
		t.Fatalf("expected 0 trades, got %d", len(body.Trades))
	}
}

// ---------------------------------------------------------------------------
// NewAgentHandlers (default singleton store)
// ---------------------------------------------------------------------------

func TestNewAgentHandlers_LoadsFixtures(t *testing.T) {
	ah, err := handlers.NewAgentHandlers()
	if err != nil {
		t.Fatalf("NewAgentHandlers: %v", err)
	}

	srv := httptest.NewServer(router.New(ah))
	t.Cleanup(srv.Close)

	resp := get(t, srv, "/v1/agents")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var body types.ListAgentsResponse
	decodeJSON(t, resp, &body)
	if len(body.Agents) == 0 {
		t.Fatal("expected agents from default store, got none")
	}
}

// ---------------------------------------------------------------------------
// Error envelope shape
// ---------------------------------------------------------------------------

func TestErrorEnvelope_Shape(t *testing.T) {
	srv := testServer(t)

	cases := []struct {
		name       string
		path       string
		wantStatus int
		wantCode   types.ErrorCode
	}{
		{"agent not found", "/v1/agents/unknown-id", http.StatusNotFound, types.ErrNotFound},
		{"trades agent not found", "/v1/agents/unknown-id/trades", http.StatusNotFound, types.ErrNotFound},
		{"bad limit", "/v1/agents?limit=999", http.StatusBadRequest, types.ErrBadRequest},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			resp := get(t, srv, tc.path)
			if resp.StatusCode != tc.wantStatus {
				t.Fatalf("expected %d, got %d", tc.wantStatus, resp.StatusCode)
			}

			var raw map[string]json.RawMessage
			decodeJSON(t, resp, &raw)

			errField, ok := raw["error"]
			if !ok {
				t.Fatal("response missing top-level 'error' field")
			}

			var inner map[string]json.RawMessage
			if err := json.Unmarshal(errField, &inner); err != nil {
				t.Fatalf("error field is not an object: %v", err)
			}

			for _, required := range []string{"code", "message"} {
				if _, ok := inner[required]; !ok {
					t.Errorf("error object missing %q field", required)
				}
			}

			var env types.ErrorEnvelope
			if err := json.Unmarshal(errField, &env.Error); err != nil {
				t.Fatalf("unmarshal error inner: %v", err)
			}
			if env.Error.Code != tc.wantCode {
				t.Errorf("expected code %q, got %q", tc.wantCode, env.Error.Code)
			}
		})
	}
}
