package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/handlers"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/router"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/types"
)

// jobsTestServer builds an httptest.Server wired to the real fixture files,
// using explicit handler injection so tests are hermetic.
func jobsTestServer(t *testing.T) *httptest.Server {
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

// postJSON issues a POST with a JSON body and returns the response.
func postJSON(t *testing.T, srv *httptest.Server, path string, body interface{}) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	resp, err := http.Post(srv.URL+path, "application/json", bytes.NewReader(b)) //nolint:noctx
	if err != nil {
		t.Fatalf("POST %s: %v", path, err)
	}
	return resp
}

// ---------------------------------------------------------------------------
// GET /v1/jobs
// ---------------------------------------------------------------------------

func TestListJobs_DefaultPage(t *testing.T) {
	srv := jobsTestServer(t)

	resp := get(t, srv, "/v1/jobs")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body types.ListJobsResponse
	decodeJSON(t, resp, &body)

	// Fixtures contain exactly 7 jobs; default limit=20 so all fit on one page.
	if len(body.Jobs) != 7 {
		t.Fatalf("expected 7 jobs, got %d", len(body.Jobs))
	}
	if body.NextCursor != "" {
		t.Fatalf("expected no next_cursor, got %q", body.NextCursor)
	}
}

func TestListJobs_Pagination(t *testing.T) {
	srv := jobsTestServer(t)

	// Page 1: limit=4
	resp1 := get(t, srv, "/v1/jobs?limit=4")
	if resp1.StatusCode != http.StatusOK {
		t.Fatalf("page 1: expected 200, got %d", resp1.StatusCode)
	}

	var page1 types.ListJobsResponse
	decodeJSON(t, resp1, &page1)

	if len(page1.Jobs) != 4 {
		t.Fatalf("page 1: expected 4 jobs, got %d", len(page1.Jobs))
	}
	if page1.NextCursor == "" {
		t.Fatal("page 1: expected next_cursor to be set")
	}
	if page1.NextCursor != cursorFor(4) {
		t.Fatalf("page 1: unexpected cursor %q", page1.NextCursor)
	}

	// Page 2: use cursor from page 1
	resp2 := get(t, srv, "/v1/jobs?limit=4&cursor="+page1.NextCursor)
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("page 2: expected 200, got %d", resp2.StatusCode)
	}

	var page2 types.ListJobsResponse
	decodeJSON(t, resp2, &page2)

	// 7 jobs total, 4 on page 1 → 3 on page 2
	if len(page2.Jobs) != 3 {
		t.Fatalf("page 2: expected 3 jobs, got %d", len(page2.Jobs))
	}
	if page2.NextCursor != "" {
		t.Fatalf("page 2: expected no next_cursor, got %q", page2.NextCursor)
	}

	// All IDs across both pages are distinct.
	seen := map[string]bool{}
	for _, j := range append(page1.Jobs, page2.Jobs...) {
		if seen[j.ID] {
			t.Fatalf("duplicate job id %q across pages", j.ID)
		}
		seen[j.ID] = true
	}
	if len(seen) != 7 {
		t.Fatalf("expected 7 unique jobs across 2 pages, got %d", len(seen))
	}
}

func TestListJobs_CursorBeyondEnd(t *testing.T) {
	srv := jobsTestServer(t)

	resp := get(t, srv, "/v1/jobs?cursor="+cursorFor(100))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body types.ListJobsResponse
	decodeJSON(t, resp, &body)

	if len(body.Jobs) != 0 {
		t.Fatalf("expected 0 jobs, got %d", len(body.Jobs))
	}
}

func TestListJobs_FilterByAgentID(t *testing.T) {
	srv := jobsTestServer(t)

	resp := get(t, srv, "/v1/jobs?agent_id=agent-001")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body types.ListJobsResponse
	decodeJSON(t, resp, &body)

	// agent-001 has jobs job-001, job-003, job-007 in fixtures
	if len(body.Jobs) != 3 {
		t.Fatalf("expected 3 jobs for agent-001, got %d", len(body.Jobs))
	}
	for _, j := range body.Jobs {
		if j.AgentID != "agent-001" {
			t.Fatalf("unexpected agent_id %q in jobs for agent-001", j.AgentID)
		}
	}
}

func TestListJobs_FilterByPhase(t *testing.T) {
	srv := jobsTestServer(t)

	resp := get(t, srv, "/v1/jobs?phase=active")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body types.ListJobsResponse
	decodeJSON(t, resp, &body)

	// Only job-002 is in the "active" phase in fixtures
	if len(body.Jobs) != 1 {
		t.Fatalf("expected 1 active job, got %d", len(body.Jobs))
	}
	if body.Jobs[0].Phase != types.JobPhaseActive {
		t.Fatalf("expected phase %q, got %q", types.JobPhaseActive, body.Jobs[0].Phase)
	}
}

func TestListJobs_FilterByAgentIDAndPhase(t *testing.T) {
	srv := jobsTestServer(t)

	resp := get(t, srv, "/v1/jobs?agent_id=agent-001&phase=released")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body types.ListJobsResponse
	decodeJSON(t, resp, &body)

	// job-001 is agent-001 + released
	if len(body.Jobs) != 1 {
		t.Fatalf("expected 1 job for agent-001 phase=released, got %d", len(body.Jobs))
	}
	if body.Jobs[0].ID != "job-001" {
		t.Fatalf("expected job-001, got %q", body.Jobs[0].ID)
	}
}

func TestListJobs_InvalidPhase(t *testing.T) {
	srv := jobsTestServer(t)

	resp := get(t, srv, "/v1/jobs?phase=unknown-phase")
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}

	var env types.ErrorEnvelope
	decodeJSON(t, resp, &env)

	if env.Error.Code != types.ErrBadRequest {
		t.Fatalf("expected code %q, got %q", types.ErrBadRequest, env.Error.Code)
	}
}

func TestListJobs_InvalidLimit(t *testing.T) {
	srv := jobsTestServer(t)

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
			resp := get(t, srv, "/v1/jobs?"+tc.param)
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

func TestListJobs_InvalidCursor(t *testing.T) {
	srv := jobsTestServer(t)

	resp := get(t, srv, "/v1/jobs?cursor=!!!notbase64")
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
	var env types.ErrorEnvelope
	decodeJSON(t, resp, &env)
	if env.Error.Code != types.ErrBadRequest {
		t.Fatalf("expected code %q, got %q", types.ErrBadRequest, env.Error.Code)
	}
}

// ---------------------------------------------------------------------------
// GET /v1/jobs/:id
// ---------------------------------------------------------------------------

func TestGetJob_ByID(t *testing.T) {
	srv := jobsTestServer(t)

	resp := get(t, srv, "/v1/jobs/job-001")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var job types.Job
	decodeJSON(t, resp, &job)

	if job.ID != "job-001" {
		t.Fatalf("expected id %q, got %q", "job-001", job.ID)
	}
	if job.Phase != types.JobPhaseReleased {
		t.Fatalf("expected phase %q, got %q", types.JobPhaseReleased, job.Phase)
	}
	if job.AgentID != "agent-001" {
		t.Fatalf("expected agent_id %q, got %q", "agent-001", job.AgentID)
	}
}

func TestGetJob_Fields(t *testing.T) {
	srv := jobsTestServer(t)

	// job-002 has all optional fields set
	resp := get(t, srv, "/v1/jobs/job-002")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var job types.Job
	decodeJSON(t, resp, &job)

	if job.EscrowAddress == "" {
		t.Fatal("expected escrow_address to be set for job-002")
	}
	if job.Amount == "" {
		t.Fatal("expected amount to be set for job-002")
	}
	if job.Deadline == nil {
		t.Fatal("expected deadline to be set for job-002")
	}
}

func TestGetJob_CreatedPhase_NoEscrow(t *testing.T) {
	srv := jobsTestServer(t)

	// job-005 is in "created" phase, no escrow yet
	resp := get(t, srv, "/v1/jobs/job-005")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var job types.Job
	decodeJSON(t, resp, &job)

	if job.Phase != types.JobPhaseCreated {
		t.Fatalf("expected phase %q, got %q", types.JobPhaseCreated, job.Phase)
	}
	if job.EscrowAddress != "" {
		t.Fatalf("expected no escrow_address for created-phase job, got %q", job.EscrowAddress)
	}
}

func TestGetJob_NotFound(t *testing.T) {
	srv := jobsTestServer(t)

	resp := get(t, srv, "/v1/jobs/does-not-exist")
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
// POST /v1/jobs/prepare
// ---------------------------------------------------------------------------

func TestPrepareJob_Valid(t *testing.T) {
	srv := jobsTestServer(t)

	req := types.PrepareJobRequest{
		AgentID:         "agent-001",
		MetadataURI:     "ipfs://QmTestPayload",
		Amount:          "1000000000000000000",
		Token:           "0x0000000000000000000000000000000000000000",
		DeadlineSeconds: 1800000000,
	}

	resp := postJSON(t, srv, "/v1/jobs/prepare", req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body types.PrepareJobResponse
	decodeJSON(t, resp, &body)

	if body.To == "" {
		t.Fatal("expected non-empty 'to' address")
	}
	if body.Calldata == "" {
		t.Fatal("expected non-empty calldata")
	}
	if body.Value == "" {
		t.Fatal("expected non-empty value")
	}
	// Native ETH: value should equal amount
	if body.Value != req.Amount {
		t.Fatalf("expected value %q, got %q", req.Amount, body.Value)
	}
}

func TestPrepareJob_TokenPayment_ValueZero(t *testing.T) {
	srv := jobsTestServer(t)

	req := types.PrepareJobRequest{
		AgentID:         "agent-002",
		MetadataURI:     "ipfs://QmTokenPayload",
		Amount:          "500000000000000000",
		Token:           "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		DeadlineSeconds: 1800000000,
	}

	resp := postJSON(t, srv, "/v1/jobs/prepare", req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body types.PrepareJobResponse
	decodeJSON(t, resp, &body)

	// ERC-20 token: value should be "0"
	if body.Value != "0" {
		t.Fatalf("expected value %q for token payment, got %q", "0", body.Value)
	}
}

func TestPrepareJob_CalldataHasSelector(t *testing.T) {
	srv := jobsTestServer(t)

	req := types.PrepareJobRequest{
		AgentID:         "agent-001",
		MetadataURI:     "ipfs://QmSelectorTest",
		Amount:          "100000000000000000",
		DeadlineSeconds: 1800000000,
	}

	resp := postJSON(t, srv, "/v1/jobs/prepare", req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body types.PrepareJobResponse
	decodeJSON(t, resp, &body)

	// Calldata must start with the 0x prefix and stub selector 0x1c4b3e91
	if len(body.Calldata) < 10 {
		t.Fatalf("calldata too short: %q", body.Calldata)
	}
	if body.Calldata[:10] != "0x1c4b3e91" {
		t.Fatalf("expected calldata to start with selector 0x1c4b3e91, got %q", body.Calldata[:10])
	}
}

func TestPrepareJob_MissingAgentID(t *testing.T) {
	srv := jobsTestServer(t)

	req := types.PrepareJobRequest{
		MetadataURI:     "ipfs://QmTest",
		Amount:          "1000000000000000000",
		DeadlineSeconds: 1800000000,
	}

	resp := postJSON(t, srv, "/v1/jobs/prepare", req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}

	var env types.ErrorEnvelope
	decodeJSON(t, resp, &env)
	if env.Error.Code != types.ErrBadRequest {
		t.Fatalf("expected code %q, got %q", types.ErrBadRequest, env.Error.Code)
	}
}

func TestPrepareJob_MissingMetadataURI(t *testing.T) {
	srv := jobsTestServer(t)

	req := types.PrepareJobRequest{
		AgentID:         "agent-001",
		Amount:          "1000000000000000000",
		DeadlineSeconds: 1800000000,
	}

	resp := postJSON(t, srv, "/v1/jobs/prepare", req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}

	var env types.ErrorEnvelope
	decodeJSON(t, resp, &env)
	if env.Error.Code != types.ErrBadRequest {
		t.Fatalf("expected code %q, got %q", types.ErrBadRequest, env.Error.Code)
	}
}

func TestPrepareJob_MissingAmount(t *testing.T) {
	srv := jobsTestServer(t)

	req := types.PrepareJobRequest{
		AgentID:         "agent-001",
		MetadataURI:     "ipfs://QmTest",
		DeadlineSeconds: 1800000000,
	}

	resp := postJSON(t, srv, "/v1/jobs/prepare", req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}

	var env types.ErrorEnvelope
	decodeJSON(t, resp, &env)
	if env.Error.Code != types.ErrBadRequest {
		t.Fatalf("expected code %q, got %q", types.ErrBadRequest, env.Error.Code)
	}
}

func TestPrepareJob_ZeroDeadline(t *testing.T) {
	srv := jobsTestServer(t)

	req := types.PrepareJobRequest{
		AgentID:         "agent-001",
		MetadataURI:     "ipfs://QmTest",
		Amount:          "1000000000000000000",
		DeadlineSeconds: 0,
	}

	resp := postJSON(t, srv, "/v1/jobs/prepare", req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}

	var env types.ErrorEnvelope
	decodeJSON(t, resp, &env)
	if env.Error.Code != types.ErrBadRequest {
		t.Fatalf("expected code %q, got %q", types.ErrBadRequest, env.Error.Code)
	}
}

func TestPrepareJob_InvalidBody(t *testing.T) {
	srv := jobsTestServer(t)

	// Send raw non-JSON body
	resp, err := http.Post(srv.URL+"/v1/jobs/prepare", "application/json", bytes.NewBufferString("not-json")) //nolint:noctx
	if err != nil {
		t.Fatalf("POST: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}

	var env types.ErrorEnvelope
	decodeJSON(t, resp, &env)
	if env.Error.Code != types.ErrBadRequest {
		t.Fatalf("expected code %q, got %q", types.ErrBadRequest, env.Error.Code)
	}
}

// ---------------------------------------------------------------------------
// NewJobHandlers (default singleton store)
// ---------------------------------------------------------------------------

func TestNewJobHandlers_LoadsFixtures(t *testing.T) {
	jh, err := handlers.NewJobHandlers()
	if err != nil {
		t.Fatalf("NewJobHandlers: %v", err)
	}

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(file), "fixtures")
	store, err := handlers.NewStoreFromDir(dir)
	if err != nil {
		t.Fatalf("load fixture store: %v", err)
	}
	ah := handlers.NewAgentHandlersWithStore(store)

	srv := httptest.NewServer(router.NewWithHandlers(ah, jh))
	t.Cleanup(srv.Close)

	resp := get(t, srv, "/v1/jobs")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var body types.ListJobsResponse
	decodeJSON(t, resp, &body)
	if len(body.Jobs) == 0 {
		t.Fatal("expected jobs from default store, got none")
	}
}
