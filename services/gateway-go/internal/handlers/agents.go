package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/types"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

// AgentHandlers holds dependencies for the agent-related HTTP handlers.
// store is injected so tests can provide a fixture store without touching the
// singleton.
type AgentHandlers struct {
	store *fixtureStore
}

// NewAgentHandlers creates a handler set backed by the default singleton store.
// Returns an error if the fixture files cannot be loaded.
func NewAgentHandlers() (*AgentHandlers, error) {
	s, err := defaultStore()
	if err != nil {
		return nil, fmt.Errorf("load fixture store: %w", err)
	}
	return &AgentHandlers{store: s}, nil
}

// NewAgentHandlersWithStore creates a handler set backed by the provided store.
// Intended for use in tests.
func NewAgentHandlersWithStore(s *fixtureStore) *AgentHandlers {
	return &AgentHandlers{store: s}
}

// ListAgents handles GET /v1/agents.
//
// Query params:
//   - limit  int  1-100, default 20
//   - cursor string  opaque base64-encoded offset token
func (h *AgentHandlers) ListAgents(w http.ResponseWriter, r *http.Request) {
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

	agents := h.store.agents
	if offset >= len(agents) {
		writeJSON(w, http.StatusOK, types.ListAgentsResponse{Agents: []types.Agent{}})
		return
	}

	page := agents[offset:]
	var nextCursor string
	if len(page) > limit {
		page = page[:limit]
		nextCursor = encodeCursor(offset + limit)
	}

	writeJSON(w, http.StatusOK, types.ListAgentsResponse{
		Agents:     page,
		NextCursor: nextCursor,
	})
}

// GetAgent handles GET /v1/agents/:id.
// :id may be either the agent id ("agent-001") or the slug ("sigma-7").
func (h *AgentHandlers) GetAgent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	agent, found := h.findAgent(id)
	if !found {
		writeError(ctx, w, http.StatusNotFound, types.ErrNotFound, "agent not found")
		return
	}

	writeJSON(w, http.StatusOK, agent)
}

// ListAgentTrades handles GET /v1/agents/:id/trades.
//
// Query params:
//   - limit  int  1-100, default 20
//   - cursor string  opaque base64-encoded offset token
func (h *AgentHandlers) ListAgentTrades(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	_, found := h.findAgent(id)
	if !found {
		writeError(ctx, w, http.StatusNotFound, types.ErrNotFound, "agent not found")
		return
	}

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

	// Resolve the canonical agent id for filtering (caller may pass slug).
	resolved, _ := h.findAgent(id)
	agentTrades := filterTradesByAgent(h.store.trades, resolved.ID)

	if offset >= len(agentTrades) {
		writeJSON(w, http.StatusOK, types.ListTradesResponse{Trades: []types.Trade{}})
		return
	}

	page := agentTrades[offset:]
	var nextCursor string
	if len(page) > limit {
		page = page[:limit]
		nextCursor = encodeCursor(offset + limit)
	}

	writeJSON(w, http.StatusOK, types.ListTradesResponse{
		Trades:     page,
		NextCursor: nextCursor,
	})
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// findAgent looks up by id first, then slug. Returns (agent, true) on match.
func (h *AgentHandlers) findAgent(idOrSlug string) (types.Agent, bool) {
	for _, a := range h.store.agents {
		if a.ID == idOrSlug || a.Slug == idOrSlug {
			return a, true
		}
	}
	return types.Agent{}, false
}

// filterTradesByAgent returns only trades for the given agent id.
func filterTradesByAgent(trades []types.Trade, agentID string) []types.Trade {
	out := make([]types.Trade, 0, len(trades))
	for _, t := range trades {
		if t.AgentID == agentID {
			out = append(out, t)
		}
	}
	return out
}

// parseLimit reads ?limit=N from the request. Returns defaultLimit if absent.
// Returns an error if the value is present but invalid.
func parseLimit(r *http.Request) (int, error) {
	raw := r.URL.Query().Get("limit")
	if raw == "" {
		return defaultLimit, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 || n > maxLimit {
		return 0, fmt.Errorf("limit must be an integer between 1 and %d", maxLimit)
	}
	return n, nil
}

// parseCursor decodes ?cursor=<base64> into a numeric offset.
// Returns 0 if absent. Returns an error if present but malformed.
func parseCursor(r *http.Request) (int, error) {
	raw := r.URL.Query().Get("cursor")
	if raw == "" {
		return 0, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return 0, errors.New("cursor is not valid base64")
	}
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 || parts[0] != "offset" {
		return 0, errors.New("cursor has an unrecognised format")
	}
	offset, err := strconv.Atoi(parts[1])
	if err != nil || offset < 0 {
		return 0, errors.New("cursor contains an invalid offset")
	}
	return offset, nil
}

// encodeCursor encodes a numeric offset as an opaque base64 cursor token.
func encodeCursor(offset int) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("offset:%d", offset)))
}

// writeJSON encodes v as JSON and writes it with the given status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// Headers already sent; nothing useful we can do.
		_ = err
	}
}

// writeError writes a structured error envelope. It never leaks stack traces or
// file paths.
func writeError(_ context.Context, w http.ResponseWriter, status int, code types.ErrorCode, message string) {
	writeJSON(w, status, types.ErrorEnvelope{
		Error: types.APIError{
			Code:    code,
			Message: message,
		},
	})
}
