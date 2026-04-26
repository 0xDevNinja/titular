package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/types"
)

// JobHandlers holds dependencies for the ACP job-related HTTP handlers.
// store is injected so tests can provide a fixture store without touching the
// singleton.
type JobHandlers struct {
	store *fixtureStore
}

// NewJobHandlers creates a handler set backed by the default singleton store.
// Returns an error if the fixture files cannot be loaded.
func NewJobHandlers() (*JobHandlers, error) {
	s, err := defaultStore()
	if err != nil {
		return nil, fmt.Errorf("load fixture store: %w", err)
	}
	return &JobHandlers{store: s}, nil
}

// NewJobHandlersWithStore creates a handler set backed by the provided store.
// Intended for use in tests.
func NewJobHandlersWithStore(s *fixtureStore) *JobHandlers {
	return &JobHandlers{store: s}
}

// ListJobs handles GET /v1/jobs.
//
// Query params:
//   - limit    int     1-100, default 20
//   - cursor   string  opaque base64-encoded offset token
//   - agent_id string  optional filter by agent id
//   - phase    string  optional filter by job phase
func (h *JobHandlers) ListJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, err := parseLimit(r)
	if err != nil {
		writeError(ctx, w, http.StatusBadRequest, types.ErrBadRequest, err.Error())
		return
	}

	offset, err := parseCursor(r)
	if err != nil {
		writeError(ctx, w, http.StatusBadRequest, types.ErrBadRequest, err.Error())
		return
	}

	jobs := h.store.jobs

	// Optional filter: agent_id
	if agentID := strings.TrimSpace(r.URL.Query().Get("agent_id")); agentID != "" {
		jobs = filterJobsByAgent(jobs, agentID)
	}

	// Optional filter: phase
	if phase := strings.TrimSpace(r.URL.Query().Get("phase")); phase != "" {
		if err := validatePhase(phase); err != nil {
			writeError(ctx, w, http.StatusBadRequest, types.ErrBadRequest, err.Error())
			return
		}
		jobs = filterJobsByPhase(jobs, types.JobPhase(phase))
	}

	if offset >= len(jobs) {
		writeJSON(w, http.StatusOK, types.ListJobsResponse{Jobs: []types.Job{}})
		return
	}

	page := jobs[offset:]
	var nextCursor string
	if len(page) > limit {
		page = page[:limit]
		nextCursor = encodeCursor(offset + limit)
	}

	writeJSON(w, http.StatusOK, types.ListJobsResponse{
		Jobs:       page,
		NextCursor: nextCursor,
	})
}

// GetJob handles GET /v1/jobs/:id.
// :id is the gateway-assigned stable identifier (e.g. "job-001").
func (h *JobHandlers) GetJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	job, found := h.findJob(id)
	if !found {
		writeError(ctx, w, http.StatusNotFound, types.ErrNotFound, "job not found")
		return
	}

	writeJSON(w, http.StatusOK, job)
}

// PrepareJob handles POST /v1/jobs/prepare.
// It validates the request and returns unsigned transaction calldata for the
// client to sign and broadcast. No on-chain interaction is performed here; the
// response is deterministic so it is safe to re-call.
func (h *JobHandlers) PrepareJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req types.PrepareJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(ctx, w, http.StatusBadRequest, types.ErrBadRequest, "invalid request body")
		return
	}

	if err := validatePrepareJobRequest(req); err != nil {
		writeError(ctx, w, http.StatusBadRequest, types.ErrBadRequest, err.Error())
		return
	}

	// Build stub calldata. In Stage 2 this will be replaced with real ABI
	// encoding once abigen bindings are generated from the ACP contracts.
	calldata := buildPrepareCalldata(req)
	value := "0"
	if req.Token == "" || req.Token == "0x0000000000000000000000000000000000000000" {
		value = req.Amount
	}

	writeJSON(w, http.StatusOK, types.PrepareJobResponse{
		// Placeholder JobFactory address; will be replaced with the real deployed
		// address once ACP contracts are live on Base Sepolia.
		To:       "0x0000000000000000000000000000000000000000",
		Calldata: calldata,
		Value:    value,
	})
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// findJob looks up a job by its gateway-assigned id.
func (h *JobHandlers) findJob(id string) (types.Job, bool) {
	for _, j := range h.store.jobs {
		if j.ID == id {
			return j, true
		}
	}
	return types.Job{}, false
}

// filterJobsByAgent returns only jobs assigned to the given agent id.
func filterJobsByAgent(jobs []types.Job, agentID string) []types.Job {
	out := make([]types.Job, 0, len(jobs))
	for _, j := range jobs {
		if j.AgentID == agentID {
			out = append(out, j)
		}
	}
	return out
}

// filterJobsByPhase returns only jobs in the given phase.
func filterJobsByPhase(jobs []types.Job, phase types.JobPhase) []types.Job {
	out := make([]types.Job, 0, len(jobs))
	for _, j := range jobs {
		if j.Phase == phase {
			out = append(out, j)
		}
	}
	return out
}

// validatePhase returns an error if the phase string is not a recognised value.
func validatePhase(phase string) error {
	switch types.JobPhase(phase) {
	case types.JobPhaseCreated,
		types.JobPhaseFunded,
		types.JobPhaseActive,
		types.JobPhaseCompleted,
		types.JobPhaseDisputed,
		types.JobPhaseReleased,
		types.JobPhaseResolved:
		return nil
	default:
		return fmt.Errorf("invalid phase %q: must be one of created, funded, active, completed, disputed, released, resolved", phase)
	}
}

// validatePrepareJobRequest returns a descriptive error for any missing or
// structurally invalid field. It never leaks internal details.
func validatePrepareJobRequest(req types.PrepareJobRequest) error {
	if strings.TrimSpace(req.AgentID) == "" {
		return fmt.Errorf("agent_id is required")
	}
	if strings.TrimSpace(req.MetadataURI) == "" {
		return fmt.Errorf("metadata_uri is required")
	}
	if strings.TrimSpace(req.Amount) == "" {
		return fmt.Errorf("amount is required")
	}
	if req.DeadlineSeconds <= 0 {
		return fmt.Errorf("deadline_seconds must be a positive Unix timestamp")
	}
	return nil
}

// buildPrepareCalldata returns a hex-encoded stub calldata for the createJob
// function selector. This is a deterministic placeholder — Stage 2 will replace
// it with real ABI encoding once abigen bindings are available.
//
// The stub encodes the four-byte selector 0x1c4b3e91 (createJob) followed by
// a zero-padded summary of the parameters so callers can verify round-trip
// correctness in tests.
func buildPrepareCalldata(req types.PrepareJobRequest) string {
	// createJob(string agentId, string metadataURI, uint256 amount, address token, uint256 deadline)
	// Stub selector; real ABI will be derived from JobFactory ABI once contracts land.
	return fmt.Sprintf("0x1c4b3e91%064x%064x%064x%064x%064x",
		len(req.AgentID),
		len(req.MetadataURI),
		len(req.Amount),
		len(req.Token),
		req.DeadlineSeconds,
	)
}
