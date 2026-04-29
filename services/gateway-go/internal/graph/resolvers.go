package graph

import (
	"errors"
	"fmt"
	"time"

	"github.com/graphql-go/graphql"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/store"
)

// Resolver bundles the dependencies that every Query/Subscription field
// needs. A single instance is built at startup (see NewSchema) and
// mounted into the schema via closures, so we avoid the global-state
// trap and keep the package testable with fakes.
//
// Bus is optional: when nil, subscription resolvers return an
// already-closed channel so a client connecting to a gateway built
// without NATS sees an immediate end-of-stream rather than a hang.
type Resolver struct {
	Store store.Store
	Bus   SubscriptionBus
}

// argString reads a non-empty string argument, returning ("", false) if
// the key is absent or empty. We never trust the raw map[string]any
// type assertion: if a client sends `{"limit": "10"}` instead of
// `{"limit": 10}`, the argument is silently the wrong shape and the
// resolver should treat it as absent.
func argString(args map[string]any, key string) (string, bool) {
	v, ok := args[key]
	if !ok || v == nil {
		return "", false
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return "", false
	}
	return s, true
}

func argInt(args map[string]any, key string) (int, bool) {
	v, ok := args[key]
	if !ok || v == nil {
		return 0, false
	}
	switch n := v.(type) {
	case int:
		return n, true
	case int32:
		return int(n), true
	case int64:
		return int(n), true
	case float64:
		// JSON numerics arrive as float64.
		return int(n), true
	}
	return 0, false
}

func argTime(args map[string]any, key string) (time.Time, error) {
	v, ok := args[key]
	if !ok || v == nil {
		return time.Time{}, nil
	}
	switch t := v.(type) {
	case time.Time:
		return t, nil
	case string:
		out, err := time.Parse(time.RFC3339Nano, t)
		if err != nil {
			out, err = time.Parse(time.RFC3339, t)
			if err != nil {
				return time.Time{}, fmt.Errorf("%s: must be RFC 3339", key)
			}
		}
		return out, nil
	}
	return time.Time{}, fmt.Errorf("%s: unexpected type %T", key, v)
}

// resolveAgent serves Query.agent. The store accepts either a numeric
// primary key or a kind-qualified slug ("kind:agent_id"); we forward
// both forms unchanged. ErrNotFound becomes nil so the GraphQL response
// is `{"data": {"agent": null}}` rather than an error envelope —
// "absent" is not a server error.
func (r *Resolver) resolveAgent(p graphql.ResolveParams) (any, error) {
	id, ok := argString(p.Args, "id")
	if !ok {
		return nil, errors.New("id is required")
	}
	a, err := r.Store.GetAgent(p.Context, id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, nil
		}
		return nil, sanitiseStoreErr(err)
	}
	return a, nil
}

// resolveAgents serves Query.agents.
func (r *Resolver) resolveAgents(p graphql.ResolveParams) (any, error) {
	limit, _ := argInt(p.Args, "limit")
	cursor, _ := argString(p.Args, "cursor")

	var kind store.AgentKind
	if v, ok := p.Args["kind"]; ok && v != nil {
		// graphql-go decodes enums into their underlying value (we set
		// `Value: store.AgentKind...` on the enum); fall back to a
		// string for safety in case the runtime hands us the raw name.
		switch k := v.(type) {
		case store.AgentKind:
			kind = k
		case string:
			kind = store.AgentKind(k)
		}
		if !kind.Valid() {
			return nil, errors.New("kind: must be one of launchpad, acp")
		}
	}

	page, err := r.Store.ListAgents(p.Context, store.ListAgentsParams{
		Kind:   kind,
		Limit:  limit,
		Cursor: cursor,
	})
	if err != nil {
		return nil, sanitiseStoreErr(err)
	}
	return page, nil
}

// resolveTrades serves Query.trades.
func (r *Resolver) resolveTrades(p graphql.ResolveParams) (any, error) {
	limit, _ := argInt(p.Args, "limit")
	cursor, _ := argString(p.Args, "cursor")
	agentToken, _ := argString(p.Args, "agentToken")

	from, err := argTime(p.Args, "from")
	if err != nil {
		return nil, err
	}
	to, err := argTime(p.Args, "to")
	if err != nil {
		return nil, err
	}

	page, err := r.Store.ListTrades(p.Context, store.ListTradesParams{
		AgentToken: agentToken,
		From:       from,
		To:         to,
		Limit:      limit,
		Cursor:     cursor,
	})
	if err != nil {
		return nil, sanitiseStoreErr(err)
	}
	return page, nil
}

// resolveJobs serves Query.jobs.
func (r *Resolver) resolveJobs(p graphql.ResolveParams) (any, error) {
	limit, _ := argInt(p.Args, "limit")
	cursor, _ := argString(p.Args, "cursor")

	var status store.JobPhase
	if v, ok := p.Args["status"]; ok && v != nil {
		switch s := v.(type) {
		case store.JobPhase:
			status = s
		case string:
			status = store.JobPhase(s)
		}
		if !status.Valid() {
			return nil, errors.New("status: invalid value")
		}
	}

	page, err := r.Store.ListJobs(p.Context, store.ListJobsParams{
		Status: status,
		Limit:  limit,
		Cursor: cursor,
	})
	if err != nil {
		return nil, sanitiseStoreErr(err)
	}
	return page, nil
}

// resolveStats serves Query.stats.
func (r *Resolver) resolveStats(p graphql.ResolveParams) (any, error) {
	st, err := r.Store.Stats(p.Context)
	if err != nil {
		return nil, sanitiseStoreErr(err)
	}
	return st, nil
}

// sanitiseStoreErr maps store-layer errors to caller-safe error messages.
//
// SECURITY: the same rule as the REST API (see handlers.api.go's
// writeAPIStoreErr): only return the verbatim message for errors we
// authored. Postgres errors (e.g. "syntax error at...") MUST become an
// opaque "internal error" in the GraphQL response — leaking those
// strings would map directly to schema enumeration.
func sanitiseStoreErr(err error) error {
	if err == nil {
		return nil
	}
	if isUserVisibleStoreErr(err.Error()) {
		return err
	}
	return errors.New("internal error")
}

// isUserVisibleStoreErr is a copy of the same allow-list used by
// handlers/api.go. We deliberately duplicate it (rather than export it
// from handlers) because this graph package owns the GraphQL surface
// and shouldn't import the REST handler package — that would be a
// layering inversion. Keep the two lists in sync; the test suite asserts
// at least one shared prefix to catch silent drift.
func isUserVisibleStoreErr(msg string) bool {
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
