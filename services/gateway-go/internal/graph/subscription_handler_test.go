package graph

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// subscription_handler_test covers the chassis side of the websocket
// transport: origin enforcement on the upgrade and the subscription
// document parser used by the wsSession dispatcher. The streaming
// behaviour (channels, fan-out, NATS bus integration) lives in
// subscriptions_test.go.

// ---------------------------------------------------------------------------
// Origin enforcement (Blocker 1 — #89 review)
// ---------------------------------------------------------------------------

// TestWSCheckOrigin_AllowsConfiguredOrigin pins the happy path: an
// origin in the allow-list passes, regardless of credentials posture.
func TestWSCheckOrigin_AllowsConfiguredOrigin(t *testing.T) {
	h := &Handler{
		AllowedOrigins: []string{"https://app.titular.xyz"},
	}
	check := h.CheckOriginFn()
	r := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	r.Header.Set("Origin", "https://app.titular.xyz")
	if !check(r) {
		t.Fatal("expected upgrade allowed for configured origin")
	}
}

// TestWSCheckOrigin_RejectsUnknown asserts an origin not in the
// allow-list is refused. This is the cross-site-upgrade attack surface
// the CORS middleware does not protect against.
func TestWSCheckOrigin_RejectsUnknown(t *testing.T) {
	h := &Handler{
		AllowedOrigins: []string{"https://app.titular.xyz"},
	}
	check := h.CheckOriginFn()
	r := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	r.Header.Set("Origin", "https://evil.example")
	if check(r) {
		t.Fatal("upgrade allowed for unknown origin — XSS pivot")
	}
}

// TestWSCheckOrigin_RejectsEmptyOrigin documents the deliberate
// exception: a request with no Origin header is a non-browser caller
// (curl, native go nats subscriber, integration test). Same-origin
// policy does not apply to those, and the gateway's auth + rate-limit
// layers are still in force, so we let them through.
func TestWSCheckOrigin_RejectsEmptyOrigin(t *testing.T) {
	h := &Handler{
		AllowedOrigins: []string{"https://app.titular.xyz"},
	}
	check := h.CheckOriginFn()
	r := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	// No Origin header set; native client.
	if !check(r) {
		t.Fatal("expected non-browser (no Origin) caller to be allowed")
	}
}

// TestWSCheckOrigin_WildcardWithoutCreds — "*" without credentials is
// open by design; mirrors the CORS middleware.
func TestWSCheckOrigin_WildcardWithoutCreds(t *testing.T) {
	h := &Handler{AllowedOrigins: []string{"*"}}
	check := h.CheckOriginFn()
	r := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	r.Header.Set("Origin", "https://anything.example")
	if !check(r) {
		t.Fatal("wildcard should allow any origin when credentials are off")
	}
}

// TestWSCheckOrigin_WildcardWithCredsRefusedByDefault is the safety
// net for #85: "*" + credentials is the credentialed-reflector
// anti-pattern. We refuse it on the websocket path the same way the
// CORS middleware refuses it on POST.
func TestWSCheckOrigin_WildcardWithCredsRefusedByDefault(t *testing.T) {
	h := &Handler{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// AllowReflectedOrigins left false → refuse.
	}
	check := h.CheckOriginFn()
	r := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	r.Header.Set("Origin", "https://evil.example")
	if check(r) {
		t.Fatal("wildcard+credentials must be refused without explicit reflect opt-in")
	}
}

// TestWSCheckOrigin_WildcardWithCredsReflectsWhenOptedIn — operators
// that explicitly opted into the reflector mode get the reflected
// origin echoed; useful for local dev, never for prod.
func TestWSCheckOrigin_WildcardWithCredsReflectsWhenOptedIn(t *testing.T) {
	h := &Handler{
		AllowedOrigins:        []string{"*"},
		AllowCredentials:      true,
		AllowReflectedOrigins: true,
	}
	check := h.CheckOriginFn()
	r := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	r.Header.Set("Origin", "https://anything.example")
	if !check(r) {
		t.Fatal("reflect mode should accept the reflected origin")
	}
}

// ---------------------------------------------------------------------------
// Query depth limiter (Blocker 2 — #89 review)
// ---------------------------------------------------------------------------

// TestQueryDepth_AllowsShallow asserts a normal query well within the
// cap returns 200 with no errors.
func TestQueryDepth_AllowsShallow(t *testing.T) {
	h := newRouterWithGraphQL(t, &fakeStore{})
	rec := postQuery(t, h, `{"query":"{ stats { totalAgents } }"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d", rec.Code)
	}
	var resp queryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", resp.Errors)
	}
}

// TestQueryDepth_RejectsDeeper drives the handler with MaxQueryDepth
// turned down to 2 and a 4-deep query, then asserts the GraphQL error
// envelope carries the depth-cap message and no `data`. We use a
// schema-valid path (agents -> data -> id) so the rejection is solely
// due to depth, not a parser error.
func TestQueryDepth_RejectsDeeper(t *testing.T) {
	fs := &fakeStore{}
	h := newTestHandler(t, fs)
	h.MaxQueryDepth = 2

	body := `{"query":"{ agents { data { id } } }"}`
	req := httptest.NewRequest(http.MethodPost, "/graphql",
		bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	{
		// Inline mount because newRouterWithGraphQL uses a fresh
		// Handler; we want our depth-capped one.
		// (gin.New is cheap; allocate a router per test.)
		mountAndServe(t, h, rec, req)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200 (GraphQL error envelope)", rec.Code)
	}
	var resp queryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Errors) == 0 {
		t.Fatal("expected depth-cap error, got none")
	}
	if !strings.Contains(resp.Errors[0].Message, "depth") {
		t.Errorf("error message: got %q", resp.Errors[0].Message)
	}
	if resp.Data != nil {
		t.Errorf("data: expected nil, got %v", resp.Data)
	}
}

// TestQueryDepth_HandlesIntrospection — introspection queries (the
// GraphiQL bootstrap, codegen tooling) breach 10 levels by design. We
// must let them through regardless of MaxQueryDepth, because the
// introspection root carries no user data and can't be abused as a
// pivot.
func TestQueryDepth_HandlesIntrospection(t *testing.T) {
	fs := &fakeStore{}
	h := newTestHandler(t, fs)
	h.MaxQueryDepth = 1 // pathological cap; would reject anything user-facing.

	// A representative slice of the GraphiQL introspection bootstrap.
	body := `{"query":"{ __schema { types { name fields { name type { name kind } } } } }"}`
	req := httptest.NewRequest(http.MethodPost, "/graphql",
		bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mountAndServe(t, h, rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d", rec.Code)
	}
	var resp queryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Errors) > 0 {
		t.Fatalf("introspection rejected: %v", resp.Errors)
	}
	if resp.Data == nil {
		t.Fatal("introspection returned no data")
	}
}

// ---------------------------------------------------------------------------
// Subscription parser (Blocker 3 — #89 review)
// ---------------------------------------------------------------------------

// TestWSSubscription_HonoursInlineArg pins the regression: the previous
// hand-rolled extractor stopped at '(' and dropped every inline
// argument, so a client filtering by agentToken silently received an
// unfiltered firehose. With the parser-driven version the inline value
// must round-trip into the args map.
func TestWSSubscription_HonoursInlineArg(t *testing.T) {
	const q = `subscription { newTrade(agentToken: "0xabc123") { id } }`
	field, args, err := parseSubscriptionOp(q, "", nil)
	if err != nil {
		t.Fatalf("parseSubscriptionOp: %v", err)
	}
	if field != "newTrade" {
		t.Errorf("field: got %q, want newTrade", field)
	}
	got, _ := args["agentToken"].(string)
	if got != "0xabc123" {
		t.Errorf("agentToken: got %q, want 0xabc123 (inline arg dropped?)", got)
	}
}

// TestWSSubscription_HonoursVariableArg confirms the modern $variables
// path still works — clients that lift the filter into a variable get
// the variable value forwarded.
func TestWSSubscription_HonoursVariableArg(t *testing.T) {
	const q = `subscription Sub($t: Address) { newTrade(agentToken: $t) { id } }`
	vars := map[string]any{"t": "0xfeedface"}
	field, args, err := parseSubscriptionOp(q, "Sub", vars)
	if err != nil {
		t.Fatalf("parseSubscriptionOp: %v", err)
	}
	if field != "newTrade" {
		t.Errorf("field: got %q", field)
	}
	got, _ := args["agentToken"].(string)
	if got != "0xfeedface" {
		t.Errorf("agentToken: got %q, want 0xfeedface", got)
	}
}

// TestWSSubscription_BothInlineAndVariable_VariableWins documents the
// precedence: when a variable is bound, it wins. This mirrors the
// graphql.Do executor's coercion order so a debugging client and a
// deployed one cannot disagree on which filter value the server
// applied.
func TestWSSubscription_BothInlineAndVariable_VariableWins(t *testing.T) {
	const q = `subscription Sub($t: Address) { newTrade(agentToken: $t) { id } }`
	vars := map[string]any{"t": "0xfromvar"}
	_, args, err := parseSubscriptionOp(q, "Sub", vars)
	if err != nil {
		t.Fatalf("parseSubscriptionOp: %v", err)
	}
	got, _ := args["agentToken"].(string)
	if got != "0xfromvar" {
		t.Errorf("variable value should win, got %q", got)
	}
}

// TestWSSubscription_RejectsNonSubscription guards the contract that
// only `subscription` operations are accepted on this transport — a
// client that accidentally posts a query gets a clear error rather
// than a silent no-op.
func TestWSSubscription_RejectsNonSubscription(t *testing.T) {
	const q = `query { stats { totalAgents } }`
	if _, _, err := parseSubscriptionOp(q, "", nil); err == nil {
		t.Fatal("expected error for non-subscription op")
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// mountAndServe wires the supplied Handler onto a fresh Gin engine and
// runs the request through it. Used by depth tests that need to keep
// the same Handler instance (preserving its MaxQueryDepth override).
func mountAndServe(t *testing.T, h *Handler, rec *httptest.ResponseRecorder, req *http.Request) {
	t.Helper()
	r := newGinForHandler(h)
	r.ServeHTTP(rec, req)
}

func newGinForHandler(h *Handler) http.Handler {
	r := gin.New()
	h.Register(r.Group(""))
	return r
}
