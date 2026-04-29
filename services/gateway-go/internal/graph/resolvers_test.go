package graph

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/store"
)

// resolver_test exercises the schema end-to-end: a fakeStore feeds the
// Resolver, and we drive ExecuteForTest to assert the GraphQL output
// shape. This is the equivalent of handlers/api_test.go for the
// GraphQL surface.

// fakeStore is a tiny in-memory implementation of store.Store. We only
// implement the parts the resolvers reach; other tests should add to
// it as new fields appear in the schema.
type fakeStore struct {
	agents []store.Agent
	trades []store.Trade
	jobs   []store.Job
	stats  store.Stats

	listAgentsErr error
	getAgentErr   error
	listTradesErr error
	listJobsErr   error
	statsErr      error
}

func (f *fakeStore) ListAgents(_ context.Context, p store.ListAgentsParams) (store.Page[store.Agent], error) {
	if f.listAgentsErr != nil {
		return store.Page[store.Agent]{}, f.listAgentsErr
	}
	out := make([]store.Agent, 0, len(f.agents))
	for _, a := range f.agents {
		if p.Kind != "" && a.Kind != p.Kind {
			continue
		}
		out = append(out, a)
	}
	limit := store.ClampLimit(p.Limit)
	if len(out) > limit {
		out = out[:limit]
	}
	return store.Page[store.Agent]{Data: out}, nil
}

func (f *fakeStore) GetAgent(_ context.Context, id string) (store.Agent, error) {
	if f.getAgentErr != nil {
		return store.Agent{}, f.getAgentErr
	}
	for _, a := range f.agents {
		if id == itoaInt64(a.ID) {
			return a, nil
		}
		if id == string(a.Kind)+":"+a.AgentID && a.AgentID != "" {
			return a, nil
		}
	}
	return store.Agent{}, store.ErrNotFound
}

func (f *fakeStore) ListTrades(_ context.Context, p store.ListTradesParams) (store.Page[store.Trade], error) {
	if f.listTradesErr != nil {
		return store.Page[store.Trade]{}, f.listTradesErr
	}
	if p.AgentToken != "" && !strings.HasPrefix(strings.ToLower(p.AgentToken), "0x") {
		return store.Page[store.Trade]{}, errors.New("agent_token: must start with 0x")
	}
	out := make([]store.Trade, 0, len(f.trades))
	for _, t := range f.trades {
		if !p.From.IsZero() && t.CreatedAt.Before(p.From) {
			continue
		}
		if !p.To.IsZero() && t.CreatedAt.After(p.To) {
			continue
		}
		out = append(out, t)
	}
	return store.Page[store.Trade]{Data: out}, nil
}

func (f *fakeStore) ListJobs(_ context.Context, p store.ListJobsParams) (store.Page[store.Job], error) {
	if f.listJobsErr != nil {
		return store.Page[store.Job]{}, f.listJobsErr
	}
	out := make([]store.Job, 0, len(f.jobs))
	for _, j := range f.jobs {
		if p.Status != "" && j.Phase != p.Status {
			continue
		}
		out = append(out, j)
	}
	return store.Page[store.Job]{Data: out}, nil
}

func (f *fakeStore) Stats(_ context.Context) (store.Stats, error) {
	if f.statsErr != nil {
		return store.Stats{}, f.statsErr
	}
	return f.stats, nil
}

func itoaInt64(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// newTestHandler builds a Handler over a fakeStore. The bus is the
// no-op nilBus; subscription tests instantiate the synchronous in-mem
// bus separately.
func newTestHandler(t *testing.T, fs *fakeStore) *Handler {
	t.Helper()
	r := &Resolver{Store: fs, Bus: NilBus()}
	schema, err := NewSchema(r)
	if err != nil {
		t.Fatalf("NewSchema: %v", err)
	}
	h, err := NewHandler(schema, r)
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}
	return h
}

// fixedTime keeps test data deterministic so the RFC3339 wire form
// matches across runs.
func fixedTime() time.Time {
	return time.Date(2026, 4, 28, 12, 0, 0, 0, time.UTC)
}

func mkAgent(id int64, kind store.AgentKind, agentID string) store.Agent {
	return store.Agent{
		ID:          id,
		Kind:        kind,
		AgentID:     agentID,
		Controller:  "0x" + repeatA(40),
		BlockNumber: uint64(100 + id),
		TxHash:      "0x" + repeatA(64),
		LogIndex:    uint32(id),
		CreatedAt:   fixedTime(),
		UpdatedAt:   fixedTime(),
	}
}

func mkTrade(id int64) store.Trade {
	return store.Trade{
		ID:          id,
		Side:        "buy",
		Trader:      "0x" + repeatA(40),
		Curve:       "0x" + repeatA(40),
		Fee:         "100",
		QuoteIn:     "1000",
		AgentOut:    "999",
		BlockNumber: 50,
		TxHash:      "0x" + repeatA(64),
		LogIndex:    uint32(id),
		CreatedAt:   fixedTime(),
	}
}

func mkJob(id int64, phase store.JobPhase) store.Job {
	return store.Job{
		ID:          id,
		OnChainID:   itoaInt64(id),
		Principal:   "0x" + repeatA(40),
		JobType:     1,
		Budget:      "1000000",
		Phase:       phase,
		BlockNumber: 200,
		TxHash:      "0x" + repeatA(64),
		LogIndex:    uint32(id),
		CreatedAt:   fixedTime(),
		UpdatedAt:   fixedTime(),
	}
}

// ---------------------------------------------------------------------------
// Query.stats
// ---------------------------------------------------------------------------

func Test_Query_Stats(t *testing.T) {
	fs := &fakeStore{stats: store.Stats{
		TotalAgents: 3, TotalJobs: 7, TotalTrades: 11, LastBlockIndexed: 1234,
	}}
	h := newTestHandler(t, fs)
	res := h.ExecuteForTest(context.Background(), `{ stats { totalAgents totalJobs totalTrades lastBlockIndexed } }`, nil)
	if len(res.Errors) > 0 {
		t.Fatalf("errors: %v", res.Errors)
	}
	got := res.Data.(map[string]any)["stats"].(map[string]any)
	if got["totalAgents"] != 3 {
		t.Errorf("totalAgents: got %v", got["totalAgents"])
	}
	if got["lastBlockIndexed"] != 1234 {
		t.Errorf("lastBlockIndexed: got %v", got["lastBlockIndexed"])
	}
}

func Test_Query_Stats_StoreErrIsOpaque(t *testing.T) {
	// SECURITY: a raw pg error MUST NOT leak into the GraphQL errors.
	fs := &fakeStore{statsErr: errors.New(`pq: relation "stats" does not exist`)}
	h := newTestHandler(t, fs)
	res := h.ExecuteForTest(context.Background(), `{ stats { totalAgents } }`, nil)
	if len(res.Errors) == 0 {
		t.Fatal("expected error")
	}
	for _, e := range res.Errors {
		if strings.Contains(e.Message, "pq:") || strings.Contains(e.Message, "relation") {
			t.Errorf("error message leaks internals: %q", e.Message)
		}
	}
}

// ---------------------------------------------------------------------------
// Query.agent
// ---------------------------------------------------------------------------

func Test_Query_Agent_ByPK(t *testing.T) {
	fs := &fakeStore{agents: []store.Agent{mkAgent(1, store.AgentKindACP, "42")}}
	h := newTestHandler(t, fs)
	res := h.ExecuteForTest(context.Background(), `{ agent(id: "1") { id kind agentId } }`, nil)
	if len(res.Errors) > 0 {
		t.Fatalf("errors: %v", res.Errors)
	}
	a := res.Data.(map[string]any)["agent"].(map[string]any)
	if a["id"] != "1" {
		t.Errorf("id: got %v", a["id"])
	}
	if a["kind"] != "acp" {
		t.Errorf("kind: got %v", a["kind"])
	}
	if a["agentId"] != "42" {
		t.Errorf("agentId: got %v", a["agentId"])
	}
}

func Test_Query_Agent_BySlug(t *testing.T) {
	fs := &fakeStore{agents: []store.Agent{mkAgent(7, store.AgentKindACP, "99")}}
	h := newTestHandler(t, fs)
	res := h.ExecuteForTest(context.Background(), `{ agent(id: "acp:99") { id agentId } }`, nil)
	if len(res.Errors) > 0 {
		t.Fatalf("errors: %v", res.Errors)
	}
	a := res.Data.(map[string]any)["agent"].(map[string]any)
	if a["id"] != "7" {
		t.Errorf("id: got %v", a["id"])
	}
}

func Test_Query_Agent_NotFoundIsNullNotError(t *testing.T) {
	h := newTestHandler(t, &fakeStore{})
	res := h.ExecuteForTest(context.Background(), `{ agent(id: "999") { id } }`, nil)
	if len(res.Errors) > 0 {
		t.Fatalf("expected no errors, got %v", res.Errors)
	}
	got := res.Data.(map[string]any)["agent"]
	if got != nil {
		t.Errorf("expected null, got %v", got)
	}
}

// ---------------------------------------------------------------------------
// Query.agents
// ---------------------------------------------------------------------------

func Test_Query_Agents_DefaultLimit(t *testing.T) {
	fs := &fakeStore{agents: []store.Agent{
		mkAgent(1, store.AgentKindACP, "1"),
		mkAgent(2, store.AgentKindACP, "2"),
	}}
	h := newTestHandler(t, fs)
	res := h.ExecuteForTest(context.Background(), `{ agents { data { id } nextCursor } }`, nil)
	if len(res.Errors) > 0 {
		t.Fatalf("errors: %v", res.Errors)
	}
	page := res.Data.(map[string]any)["agents"].(map[string]any)
	data := page["data"].([]any)
	if len(data) != 2 {
		t.Errorf("got %d agents, want 2", len(data))
	}
	if page["nextCursor"] != nil {
		t.Errorf("nextCursor: got %v, want nil", page["nextCursor"])
	}
}

func Test_Query_Agents_FilterByKind(t *testing.T) {
	fs := &fakeStore{agents: []store.Agent{
		mkAgent(1, store.AgentKindLaunchpad, "10"),
		mkAgent(2, store.AgentKindACP, "20"),
	}}
	h := newTestHandler(t, fs)
	res := h.ExecuteForTest(context.Background(), `{ agents(kind: acp) { data { id kind } } }`, nil)
	if len(res.Errors) > 0 {
		t.Fatalf("errors: %v", res.Errors)
	}
	data := res.Data.(map[string]any)["agents"].(map[string]any)["data"].([]any)
	if len(data) != 1 {
		t.Fatalf("got %d, want 1", len(data))
	}
	if data[0].(map[string]any)["kind"] != "acp" {
		t.Errorf("kind: got %v", data[0].(map[string]any)["kind"])
	}
}

func Test_Query_Agents_RejectsBadKind(t *testing.T) {
	h := newTestHandler(t, &fakeStore{})
	res := h.ExecuteForTest(context.Background(), `{ agents(kind: bogus) { data { id } } }`, nil)
	if len(res.Errors) == 0 {
		t.Fatal("expected validation error for bad enum value")
	}
}

// ---------------------------------------------------------------------------
// Query.trades
// ---------------------------------------------------------------------------

func Test_Query_Trades(t *testing.T) {
	fs := &fakeStore{trades: []store.Trade{mkTrade(1), mkTrade(2)}}
	h := newTestHandler(t, fs)
	res := h.ExecuteForTest(context.Background(), `{ trades { data { id side fee } } }`, nil)
	if len(res.Errors) > 0 {
		t.Fatalf("errors: %v", res.Errors)
	}
	data := res.Data.(map[string]any)["trades"].(map[string]any)["data"].([]any)
	if len(data) != 2 {
		t.Fatalf("got %d trades", len(data))
	}
	first := data[0].(map[string]any)
	if first["fee"] != "100" {
		t.Errorf("fee: got %v", first["fee"])
	}
}

func Test_Query_Trades_BadAgentToken_BubblesAsBigQueryError(t *testing.T) {
	h := newTestHandler(t, &fakeStore{})
	// Address scalar enforces validation on input; the fakeStore never
	// sees it.
	res := h.ExecuteForTest(context.Background(), `{ trades(agentToken: "notahex") { data { id } } }`, nil)
	if len(res.Errors) == 0 {
		t.Fatal("expected validation error from Address scalar")
	}
}

// ---------------------------------------------------------------------------
// Query.jobs
// ---------------------------------------------------------------------------

func Test_Query_Jobs_FilterByStatus(t *testing.T) {
	fs := &fakeStore{jobs: []store.Job{
		mkJob(1, store.JobPhaseCreated),
		mkJob(2, store.JobPhaseFunded),
	}}
	h := newTestHandler(t, fs)
	res := h.ExecuteForTest(context.Background(), `{ jobs(status: funded) { data { id phase } } }`, nil)
	if len(res.Errors) > 0 {
		t.Fatalf("errors: %v", res.Errors)
	}
	data := res.Data.(map[string]any)["jobs"].(map[string]any)["data"].([]any)
	if len(data) != 1 {
		t.Fatalf("got %d", len(data))
	}
	if data[0].(map[string]any)["phase"] != "funded" {
		t.Errorf("phase: got %v", data[0].(map[string]any)["phase"])
	}
}

// ---------------------------------------------------------------------------
// BigInt round trip via JSON output
// ---------------------------------------------------------------------------

func Test_BigInt_OnTheWire(t *testing.T) {
	// Trade.fee is a BigInt; ensure the produced JSON is a quoted
	// string, never a bare number.
	fs := &fakeStore{trades: []store.Trade{mkTrade(1)}}
	h := newTestHandler(t, fs)
	res := h.ExecuteForTest(context.Background(), `{ trades { data { fee } } }`, nil)
	if len(res.Errors) > 0 {
		t.Fatalf("errors: %v", res.Errors)
	}
	body, err := json.Marshal(res.Data)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), `"fee":"100"`) {
		t.Errorf("fee not encoded as quoted string: %s", body)
	}
}

// ---------------------------------------------------------------------------
// Sanitisation contract
// ---------------------------------------------------------------------------

func Test_isUserVisibleStoreErr(t *testing.T) {
	// The allow-list MUST stay in sync with handlers/api.go's copy of
	// the same logic. A representative entry from each prefix family
	// pins the contract; full duplication of the handlers test is not
	// needed because both call sites share the prefix shape.
	cases := []struct {
		msg  string
		want bool
	}{
		{"cursor is not valid base64", true},
		{"agent_token: must start with 0x", true},
		{"invalid kind whatever", true},
		{"invalid status whatever", true},
		{"pq: relation does not exist", false},
		{"connection refused", false},
		{"", false},
	}
	for _, c := range cases {
		if got := isUserVisibleStoreErr(c.msg); got != c.want {
			t.Errorf("isUserVisibleStoreErr(%q) = %v, want %v", c.msg, got, c.want)
		}
	}
}
