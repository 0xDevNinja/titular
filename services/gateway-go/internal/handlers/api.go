package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/store"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/types"
)

// API holds the dependencies for the Postgres-backed REST endpoints introduced
// in M4. It is mounted under /api/v1 (see router.NewWithConfig) and is
// deliberately separate from the M2/M3 fixture-backed handlers (AgentHandlers,
// JobHandlers) in this same package — those are kept around for backwards
// compatibility with the SDK alpha while the new endpoints stabilise.
type API struct {
	Store store.Store
}

// NewAPI constructs an API wired to the given store. Required so cmd/gateway
// can build the API after wiring the Postgres pool.
func NewAPI(s store.Store) *API { return &API{Store: s} }

// Register mounts the read-only endpoints on the supplied router group.
//
// The caller is expected to have already applied any auth middleware to the
// group; this method is intentionally agnostic so a test can mount the API
// onto a no-auth group while production wiring goes under RequireAuth.
func (a *API) Register(rg *gin.RouterGroup) {
	rg.GET("/agents", a.listAgents)
	rg.GET("/agents/:id", a.getAgent)
	rg.GET("/trades", a.listTrades)
	rg.GET("/jobs", a.listJobs)
	rg.GET("/stats", a.stats)
}

// allowedQueryParams lists the query keys each endpoint understands. We
// reject unknown keys with 400 to surface client bugs early — silently
// ignoring them lets typos turn into mysterious "no filter applied" results.
//
// `agent` (singular) is the per-id detail endpoint. It takes no query
// parameters today, but the entry exists so getAgent can call
// rejectUnknownParams uniformly with the list endpoints.
var allowedQueryParams = map[string]map[string]struct{}{
	"agents": {"cursor": {}, "limit": {}, "kind": {}},
	"agent":  {},
	"trades": {"cursor": {}, "limit": {}, "agent_token": {}, "from": {}, "to": {}},
	"jobs":   {"cursor": {}, "limit": {}, "status": {}},
}

// rejectUnknownParams returns true and writes a 400 envelope when the request
// contains a key not in allowed. Returns false when the request is clean.
func rejectUnknownParams(c *gin.Context, allowed map[string]struct{}) bool {
	for k := range c.Request.URL.Query() {
		if _, ok := allowed[k]; !ok {
			writeAPIError(c, http.StatusBadRequest, types.ErrBadRequest,
				"unknown query parameter: "+k)
			return true
		}
	}
	return false
}

// parseLimitParam reads and validates ?limit. Empty → store.DefaultLimit.
// Out of [1, store.MaxLimit] → 400. Returns the validated value or 0 on error
// (the caller has already responded in the error case).
func parseLimitParam(c *gin.Context) (int, bool) {
	raw := c.Query("limit")
	if raw == "" {
		return store.DefaultLimit, true
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 || n > store.MaxLimit {
		writeAPIError(c, http.StatusBadRequest, types.ErrBadRequest,
			"limit must be an integer between 1 and "+strconv.Itoa(store.MaxLimit))
		return 0, false
	}
	return n, true
}

// listAgents serves GET /api/v1/agents.
//
// @Summary      List agents
// @Description  Returns a paginated, opaque-cursor page of registered agents.
// @Description  Filterable by agent kind. Cursor is base64-encoded JSON the
// @Description  client MUST treat as opaque and pass back unmodified.
// @Tags         agents
// @Produce      json
// @Param        cursor  query     string  false  "opaque pagination cursor; omit for the first page"
// @Param        limit   query     int     false  "page size; 1..200, default 50"  minimum(1)  maximum(200)
// @Param        kind    query     string  false  "agent kind filter"  Enums(launchpad, acp)
// @Success      200     {object}  openapi.AgentsPage
// @Failure      400     {object}  openapi.AuthErrorResponse  "unknown query parameter, or kind/limit out of range"
// @Failure      500     {object}  openapi.AuthErrorResponse  "internal server error"
// @Router       /api/v1/agents [get]
func (a *API) listAgents(c *gin.Context) {
	if rejectUnknownParams(c, allowedQueryParams["agents"]) {
		return
	}
	limit, ok := parseLimitParam(c)
	if !ok {
		return
	}
	var kind store.AgentKind
	if v := c.Query("kind"); v != "" {
		kind = store.AgentKind(v)
		if !kind.Valid() {
			writeAPIError(c, http.StatusBadRequest, types.ErrBadRequest,
				`kind must be one of "launchpad", "acp"`)
			return
		}
	}
	page, err := a.Store.ListAgents(c.Request.Context(), store.ListAgentsParams{
		Kind:   kind,
		Limit:  limit,
		Cursor: c.Query("cursor"),
	})
	if err != nil {
		writeAPIStoreErr(c, err)
		return
	}
	c.JSON(http.StatusOK, page)
}

// getAgent serves GET /api/v1/agents/:id. The :id is either a numeric primary
// key or the kind-qualified slug "kind:agent_id".
//
// @Summary      Get agent
// @Description  Returns the registered agent identified by its primary key
// @Description  (decimal int64) or kind-qualified slug "kind:agent_id".
// @Tags         agents
// @Produce      json
// @Param        id   path      string  true  "agent primary key or kind:agent_id slug"
// @Success      200  {object}  store.Agent
// @Failure      404  {object}  openapi.AuthErrorResponse  "agent not found"
// @Failure      500  {object}  openapi.AuthErrorResponse  "internal server error"
// @Router       /api/v1/agents/{id} [get]
func (a *API) getAgent(c *gin.Context) {
	if rejectUnknownParams(c, allowedQueryParams["agent"]) {
		return
	}
	id := c.Param("id")
	agent, err := a.Store.GetAgent(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeAPIError(c, http.StatusNotFound, types.ErrNotFound, "agent not found")
			return
		}
		writeAPIStoreErr(c, err)
		return
	}
	c.JSON(http.StatusOK, agent)
}

// listTrades serves GET /api/v1/trades.
//
// @Summary      List trades
// @Description  Returns a paginated, opaque-cursor page of bonding-curve
// @Description  trades (buy + sell). Filterable by agent token address and
// @Description  by an inclusive [from, to] window on the indexer-recorded
// @Description  created_at timestamp.
// @Tags         trades
// @Produce      json
// @Param        cursor       query     string  false  "opaque pagination cursor"
// @Param        limit        query     int     false  "page size; 1..200, default 50"  minimum(1)  maximum(200)
// @Param        agent_token  query     string  false  "0x-prefixed 20-byte token address"
// @Param        from         query     string  false  "inclusive RFC 3339 lower bound on created_at"  format(date-time)
// @Param        to           query     string  false  "inclusive RFC 3339 upper bound on created_at"  format(date-time)
// @Success      200          {object}  openapi.TradesPage
// @Failure      400          {object}  openapi.AuthErrorResponse  "unknown query parameter, or invalid time/limit"
// @Failure      500          {object}  openapi.AuthErrorResponse  "internal server error"
// @Router       /api/v1/trades [get]
func (a *API) listTrades(c *gin.Context) {
	if rejectUnknownParams(c, allowedQueryParams["trades"]) {
		return
	}
	limit, ok := parseLimitParam(c)
	if !ok {
		return
	}
	from, ok := parseTimeParam(c, "from")
	if !ok {
		return
	}
	to, ok := parseTimeParam(c, "to")
	if !ok {
		return
	}
	page, err := a.Store.ListTrades(c.Request.Context(), store.ListTradesParams{
		AgentToken: c.Query("agent_token"),
		From:       from,
		To:         to,
		Limit:      limit,
		Cursor:     c.Query("cursor"),
	})
	if err != nil {
		writeAPIStoreErr(c, err)
		return
	}
	c.JSON(http.StatusOK, page)
}

// listJobs serves GET /api/v1/jobs.
//
// @Summary      List ACP jobs
// @Description  Returns a paginated, opaque-cursor page of ACP jobs.
// @Description  Filterable by phase via the `status` query parameter.
// @Tags         jobs
// @Produce      json
// @Param        cursor  query     string  false  "opaque pagination cursor"
// @Param        limit   query     int     false  "page size; 1..200, default 50"  minimum(1)  maximum(200)
// @Param        status  query     string  false  "ACP job phase filter"  Enums(created, funded, active, completed, cancelled, disputed, released, resolved)
// @Success      200     {object}  openapi.JobsPage
// @Failure      400     {object}  openapi.AuthErrorResponse  "unknown query parameter, or invalid status/limit"
// @Failure      500     {object}  openapi.AuthErrorResponse  "internal server error"
// @Router       /api/v1/jobs [get]
func (a *API) listJobs(c *gin.Context) {
	if rejectUnknownParams(c, allowedQueryParams["jobs"]) {
		return
	}
	limit, ok := parseLimitParam(c)
	if !ok {
		return
	}
	var status store.JobPhase
	if v := c.Query("status"); v != "" {
		status = store.JobPhase(v)
		if !status.Valid() {
			writeAPIError(c, http.StatusBadRequest, types.ErrBadRequest,
				"status must be one of created, funded, active, completed, cancelled, disputed, released, resolved")
			return
		}
	}
	page, err := a.Store.ListJobs(c.Request.Context(), store.ListJobsParams{
		Status: status,
		Limit:  limit,
		Cursor: c.Query("cursor"),
	})
	if err != nil {
		writeAPIStoreErr(c, err)
		return
	}
	c.JSON(http.StatusOK, page)
}

// stats serves GET /api/v1/stats.
//
// @Summary      Aggregate protocol stats
// @Description  Returns total agent / job / trade counters and the highest
// @Description  block number observed by the indexer track. Cached to a few
// @Description  seconds upstream; clients SHOULD NOT poll faster than 5s.
// @Tags         stats
// @Produce      json
// @Success      200  {object}  store.Stats
// @Failure      500  {object}  openapi.AuthErrorResponse  "internal server error"
// @Router       /api/v1/stats [get]
func (a *API) stats(c *gin.Context) {
	st, err := a.Store.Stats(c.Request.Context())
	if err != nil {
		writeAPIStoreErr(c, err)
		return
	}
	c.JSON(http.StatusOK, st)
}

// parseTimeParam reads an RFC 3339 timestamp from the named query key. Returns
// (zero, true) when the key is absent. On parse failure it writes 400 and
// returns (zero, false).
func parseTimeParam(c *gin.Context, key string) (time.Time, bool) {
	raw := c.Query(key)
	if raw == "" {
		return time.Time{}, true
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		writeAPIError(c, http.StatusBadRequest, types.ErrBadRequest,
			key+" must be an RFC 3339 timestamp")
		return time.Time{}, false
	}
	return t, true
}

// writeAPIError sends a structured error envelope. We keep the on-the-wire
// shape identical to the legacy chi handlers so SDK consumers don't have to
// branch on which surface emitted the error.
func writeAPIError(c *gin.Context, status int, code types.ErrorCode, msg string) {
	c.AbortWithStatusJSON(status, types.ErrorEnvelope{
		Error: types.APIError{Code: code, Message: msg},
	})
}

// writeAPIStoreErr maps store-layer errors to HTTP. Anything we don't
// specifically recognise becomes 500 with a generic message — never the raw
// error string, because that path can leak SQL fragments or schema names.
func writeAPIStoreErr(c *gin.Context, err error) {
	// Validation-style errors from the store: surface as 400 with the message
	// (these strings are safe — they're written by us, not by pg).
	if isUserVisibleStoreErr(err) {
		writeAPIError(c, http.StatusBadRequest, types.ErrBadRequest, err.Error())
		return
	}
	writeAPIError(c, http.StatusInternalServerError, types.ErrInternalServer, "internal error")
}

// isUserVisibleStoreErr reports whether err originated as a caller-input
// validation failure inside the store. We use a tag list rather than a custom
// type so the store package can stay stdlib-only.
//
// SECURITY: only return true for errors we authored — a Postgres error
// (e.g. "syntax error at ...") must never leak through this gate.
func isUserVisibleStoreErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	switch {
	case startsWith(msg, "cursor "):
		return true
	case startsWith(msg, "agent_token: "):
		return true
	case startsWith(msg, "invalid kind "):
		return true
	case startsWith(msg, "invalid status "):
		return true
	}
	return false
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
