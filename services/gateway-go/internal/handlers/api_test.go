package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/handlers"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/store"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/types"
)

func init() { gin.SetMode(gin.TestMode) }

// ---------------------------------------------------------------------------
// fakeStore — in-memory store for handler tests.
//
// We deliberately avoid a mocking library: the surface is five methods, and
// the only cross-cutting concern (cursor encode/decode) is shared with the
// production store via the exported store.EncodeCursor / DecodeCursor.
// ---------------------------------------------------------------------------

type fakeStore struct {
	agents []store.Agent
	trades []store.Trade
	jobs   []store.Job
	stats  store.Stats

	// Programmable error returns for negative tests.
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
	cur, err := store.DecodeCursor(p.Cursor)
	if err != nil {
		return store.Page[store.Agent]{}, err
	}
	limit := store.ClampLimit(p.Limit)

	out := make([]store.Agent, 0, len(f.agents))
	for _, a := range f.agents {
		if p.Kind != "" && a.Kind != p.Kind {
			continue
		}
		if !cur.CreatedAt.IsZero() {
			if a.CreatedAt.After(cur.CreatedAt) {
				continue
			}
			if a.CreatedAt.Equal(cur.CreatedAt) && a.ID >= cur.ID {
				continue
			}
		}
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool {
		if !out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].CreatedAt.After(out[j].CreatedAt)
		}
		return out[i].ID > out[j].ID
	})
	page := store.Page[store.Agent]{Data: out}
	if len(out) > limit {
		last := out[limit-1]
		page.Data = out[:limit]
		page.NextCursor = store.EncodeCursor(store.Cursor{ID: last.ID, CreatedAt: last.CreatedAt})
	}
	return page, nil
}

func (f *fakeStore) GetAgent(_ context.Context, idOrSlug string) (store.Agent, error) {
	if f.getAgentErr != nil {
		return store.Agent{}, f.getAgentErr
	}
	for _, a := range f.agents {
		// Numeric primary key.
		if v := strings.TrimSpace(idOrSlug); v != "" {
			if v == itoa(a.ID) {
				return a, nil
			}
			if v == string(a.Kind)+":"+a.AgentID && a.AgentID != "" {
				return a, nil
			}
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
	limit := store.ClampLimit(p.Limit)
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
	page := store.Page[store.Trade]{Data: out}
	if len(out) > limit {
		page.Data = out[:limit]
	}
	return page, nil
}

func (f *fakeStore) ListJobs(_ context.Context, p store.ListJobsParams) (store.Page[store.Job], error) {
	if f.listJobsErr != nil {
		return store.Page[store.Job]{}, f.listJobsErr
	}
	limit := store.ClampLimit(p.Limit)
	out := make([]store.Job, 0, len(f.jobs))
	for _, j := range f.jobs {
		if p.Status != "" && j.Phase != p.Status {
			continue
		}
		out = append(out, j)
	}
	page := store.Page[store.Job]{Data: out}
	if len(out) > limit {
		page.Data = out[:limit]
	}
	return page, nil
}

func (f *fakeStore) Stats(_ context.Context) (store.Stats, error) {
	if f.statsErr != nil {
		return store.Stats{}, f.statsErr
	}
	return f.stats, nil
}

// itoa avoids pulling strconv into the fake-store closures; it's used in
// fewer than ten spots so the saving on imports is real.
func itoa(n int64) string {
	const digits = "0123456789"
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
		buf[i] = digits[n%10]
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// newAPI wraps a fakeStore in a Gin engine with the API mounted under
// /api/v1.
func newAPI(t *testing.T, fs *fakeStore) http.Handler {
	t.Helper()
	api := handlers.NewAPI(fs)
	r := gin.New()
	api.Register(r.Group("/api/v1"))
	return r
}

func mustGet(t *testing.T, h http.Handler, target string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// ---------------------------------------------------------------------------
// /api/v1/agents
// ---------------------------------------------------------------------------

func TestAPIListAgents_DefaultLimit(t *testing.T) {
	fs := &fakeStore{agents: makeAgents(3)}
	rec := mustGet(t, newAPI(t, fs), "/api/v1/agents")

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}
	var page store.Page[store.Agent]
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(page.Data) != 3 {
		t.Fatalf("agents: got %d, want 3", len(page.Data))
	}
	if page.NextCursor != "" {
		t.Fatalf("next_cursor: got %q, want empty", page.NextCursor)
	}
}

func TestAPIListAgents_Pagination(t *testing.T) {
	fs := &fakeStore{agents: makeAgents(5)}
	rec := mustGet(t, newAPI(t, fs), "/api/v1/agents?limit=2")
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d", rec.Code)
	}
	var page store.Page[store.Agent]
	_ = json.Unmarshal(rec.Body.Bytes(), &page)
	if len(page.Data) != 2 {
		t.Fatalf("page1 len: got %d", len(page.Data))
	}
	if page.NextCursor == "" {
		t.Fatal("expected next_cursor on full page")
	}
}

func TestAPIListAgents_FilterByKind(t *testing.T) {
	fs := &fakeStore{agents: append(makeAgents(2), store.Agent{
		ID: 99, Kind: store.AgentKindACP, AgentID: "1",
		CreatedAt: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
	})}
	rec := mustGet(t, newAPI(t, fs), "/api/v1/agents?kind=acp")
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d", rec.Code)
	}
	var page store.Page[store.Agent]
	_ = json.Unmarshal(rec.Body.Bytes(), &page)
	if len(page.Data) != 1 || page.Data[0].ID != 99 {
		t.Fatalf("filtered: got %+v", page.Data)
	}
}

func TestAPIListAgents_RejectsBadKind(t *testing.T) {
	rec := mustGet(t, newAPI(t, &fakeStore{}), "/api/v1/agents?kind=bogus")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
	mustErrCode(t, rec, types.ErrBadRequest)
}

func TestAPIListAgents_RejectsBadLimit(t *testing.T) {
	tests := []string{"-1", "0", "201", "abc"}
	for _, raw := range tests {
		raw := raw
		t.Run(raw, func(t *testing.T) {
			rec := mustGet(t, newAPI(t, &fakeStore{}), "/api/v1/agents?limit="+raw)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("limit=%q: got %d, want 400", raw, rec.Code)
			}
		})
	}
}

func TestAPIListAgents_RejectsUnknownParam(t *testing.T) {
	rec := mustGet(t, newAPI(t, &fakeStore{}), "/api/v1/agents?evil=1")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestAPIListAgents_BadCursor(t *testing.T) {
	rec := mustGet(t, newAPI(t, &fakeStore{}), "/api/v1/agents?cursor=not-base64!")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

// ---------------------------------------------------------------------------
// /api/v1/agents/:id
// ---------------------------------------------------------------------------

func TestAPIGetAgent_NotFound(t *testing.T) {
	rec := mustGet(t, newAPI(t, &fakeStore{}), "/api/v1/agents/999")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: got %d, want 404", rec.Code)
	}
}

func TestAPIGetAgent_ByPK(t *testing.T) {
	fs := &fakeStore{agents: makeAgents(1)}
	rec := mustGet(t, newAPI(t, fs), "/api/v1/agents/1")
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d", rec.Code)
	}
	var a store.Agent
	if err := json.Unmarshal(rec.Body.Bytes(), &a); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if a.ID != 1 {
		t.Fatalf("id: got %d, want 1", a.ID)
	}
}

func TestAPIGetAgent_BySlug(t *testing.T) {
	fs := &fakeStore{agents: []store.Agent{{
		ID: 7, Kind: store.AgentKindACP, AgentID: "42",
	}}}
	rec := mustGet(t, newAPI(t, fs), "/api/v1/agents/acp:42")
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d", rec.Code)
	}
}

// Parity with listAgents/listTrades/listJobs: an unknown query parameter on
// the per-id detail endpoint must produce a 400 rather than be silently
// ignored. Catches client typos like ?include=bar where include isn't real.
func TestAPIGetAgent_UnknownQueryParam(t *testing.T) {
	fs := &fakeStore{agents: makeAgents(1)}
	rec := mustGet(t, newAPI(t, fs), "/api/v1/agents/1?bogus=1")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

// ---------------------------------------------------------------------------
// /api/v1/trades
// ---------------------------------------------------------------------------

func TestAPIListTrades_OK(t *testing.T) {
	fs := &fakeStore{trades: makeTrades(2)}
	rec := mustGet(t, newAPI(t, fs), "/api/v1/trades")
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d", rec.Code)
	}
}

func TestAPIListTrades_BadAgentToken(t *testing.T) {
	rec := mustGet(t, newAPI(t, &fakeStore{}), "/api/v1/trades?agent_token=notahex")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

func TestAPIListTrades_BadFromTimestamp(t *testing.T) {
	rec := mustGet(t, newAPI(t, &fakeStore{}), "/api/v1/trades?from=yesterday")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

// ---------------------------------------------------------------------------
// /api/v1/jobs
// ---------------------------------------------------------------------------

func TestAPIListJobs_OK(t *testing.T) {
	fs := &fakeStore{jobs: makeJobs(3)}
	rec := mustGet(t, newAPI(t, fs), "/api/v1/jobs")
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d", rec.Code)
	}
	var page store.Page[store.Job]
	_ = json.Unmarshal(rec.Body.Bytes(), &page)
	if len(page.Data) != 3 {
		t.Fatalf("len: got %d, want 3", len(page.Data))
	}
}

func TestAPIListJobs_FilterByStatus(t *testing.T) {
	fs := &fakeStore{jobs: []store.Job{
		{ID: 1, Phase: store.JobPhaseCreated},
		{ID: 2, Phase: store.JobPhaseFunded},
	}}
	rec := mustGet(t, newAPI(t, fs), "/api/v1/jobs?status=funded")
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d", rec.Code)
	}
	var page store.Page[store.Job]
	_ = json.Unmarshal(rec.Body.Bytes(), &page)
	if len(page.Data) != 1 || page.Data[0].ID != 2 {
		t.Fatalf("filtered: got %+v", page.Data)
	}
}

func TestAPIListJobs_RejectsBadStatus(t *testing.T) {
	rec := mustGet(t, newAPI(t, &fakeStore{}), "/api/v1/jobs?status=bogus")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
}

// ---------------------------------------------------------------------------
// /api/v1/stats
// ---------------------------------------------------------------------------

func TestAPIStats_OK(t *testing.T) {
	fs := &fakeStore{stats: store.Stats{
		TotalAgents:      3,
		TotalJobs:        7,
		TotalTrades:      11,
		LastBlockIndexed: 1234,
	}}
	rec := mustGet(t, newAPI(t, fs), "/api/v1/stats")
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d", rec.Code)
	}
	var got store.Stats
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got != fs.stats {
		t.Fatalf("stats: got %+v, want %+v", got, fs.stats)
	}
}

func TestAPIStats_StoreErrIsOpaque(t *testing.T) {
	// SECURITY: a raw pg error MUST NOT leak through the response body.
	fs := &fakeStore{statsErr: errors.New(`pq: relation "agents" does not exist`)}
	rec := mustGet(t, newAPI(t, fs), "/api/v1/stats")
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}
	body := rec.Body.String()
	if strings.Contains(body, "relation") || strings.Contains(body, "pq:") {
		t.Fatalf("response leaked internal error: %s", body)
	}
}

// ---------------------------------------------------------------------------
// fixtures
// ---------------------------------------------------------------------------

func makeAgents(n int) []store.Agent {
	now := time.Date(2026, 4, 28, 12, 0, 0, 0, time.UTC)
	out := make([]store.Agent, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, store.Agent{
			ID:          int64(i + 1),
			Kind:        store.AgentKindLaunchpad,
			Creator:     "0x000000000000000000000000000000000000beef",
			BlockNumber: uint64(100 + i),
			TxHash:      "0x" + strings.Repeat("0", 64),
			LogIndex:    uint32(i),
			CreatedAt:   now.Add(time.Duration(-i) * time.Minute),
			UpdatedAt:   now,
		})
	}
	return out
}

func makeTrades(n int) []store.Trade {
	now := time.Date(2026, 4, 28, 12, 0, 0, 0, time.UTC)
	out := make([]store.Trade, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, store.Trade{
			ID:          int64(i + 1),
			Side:        "buy",
			Trader:      "0x" + strings.Repeat("a", 40),
			Curve:       "0x" + strings.Repeat("b", 40),
			Fee:         "0",
			BlockNumber: uint64(50 + i),
			TxHash:      "0x" + strings.Repeat("0", 64),
			LogIndex:    uint32(i),
			CreatedAt:   now,
		})
	}
	return out
}

func makeJobs(n int) []store.Job {
	now := time.Date(2026, 4, 28, 12, 0, 0, 0, time.UTC)
	out := make([]store.Job, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, store.Job{
			ID:          int64(i + 1),
			OnChainID:   itoa(int64(i)),
			Principal:   "0x" + strings.Repeat("c", 40),
			JobType:     0,
			Budget:      "1000",
			Phase:       store.JobPhaseCreated,
			BlockNumber: uint64(200 + i),
			TxHash:      "0x" + strings.Repeat("0", 64),
			LogIndex:    uint32(i),
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}
	return out
}

func mustErrCode(t *testing.T, rec *httptest.ResponseRecorder, want types.ErrorCode) {
	t.Helper()
	var env types.ErrorEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("decode err envelope: %v", err)
	}
	if env.Error.Code != want {
		t.Fatalf("error code: got %q, want %q", env.Error.Code, want)
	}
}
