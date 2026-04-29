package graph

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/store"
)

// NewSchema builds the executable GraphQL schema. The schema is
// stateless across requests: arguments, context and the resolver bundle
// (which carries Store + Bus) are the only inputs to a field resolver.
//
// The shape is documented in schema.graphqls; build.go is the
// authoritative implementation. We keep the SDL file for SDK tooling
// (graphql-codegen, Apollo CLI) and assert equivalence in
// schema_test.go via introspection.
//
// Subscriptions
//
// graphql-go does not implement subscriptions through the synchronous
// graphql.Do(...) entrypoint. We expose Subscription resolvers as
// regular field resolvers so introspection sees the same root object;
// the actual streaming path runs through the websocket handler in
// subscription_handler.go which calls Resolver.runSubscription
// directly.
func NewSchema(r *Resolver) (graphql.Schema, error) {
	if r == nil {
		return graphql.Schema{}, errors.New("graph: nil resolver")
	}
	if r.Store == nil {
		return graphql.Schema{}, errors.New("graph: nil store")
	}
	if r.Bus == nil {
		// NilBus prefer-nothing-is-better-than-nil-deref: subscription
		// resolvers must always have a bus to talk to.
		r.Bus = NilBus()
	}

	agentPage := pageType(
		"AgentPage",
		agentType,
		func(p graphql.ResolveParams) (any, error) {
			page := p.Source.(store.Page[store.Agent])
			return page.Data, nil
		},
		func(p graphql.ResolveParams) (any, error) {
			page := p.Source.(store.Page[store.Agent])
			if page.NextCursor == "" {
				return nil, nil
			}
			return page.NextCursor, nil
		},
	)
	tradePage := pageType(
		"TradePage",
		tradeType,
		func(p graphql.ResolveParams) (any, error) {
			page := p.Source.(store.Page[store.Trade])
			return page.Data, nil
		},
		func(p graphql.ResolveParams) (any, error) {
			page := p.Source.(store.Page[store.Trade])
			if page.NextCursor == "" {
				return nil, nil
			}
			return page.NextCursor, nil
		},
	)
	jobPage := pageType(
		"JobPage",
		jobType,
		func(p graphql.ResolveParams) (any, error) {
			page := p.Source.(store.Page[store.Job])
			return page.Data, nil
		},
		func(p graphql.ResolveParams) (any, error) {
			page := p.Source.(store.Page[store.Job])
			if page.NextCursor == "" {
				return nil, nil
			}
			return page.NextCursor, nil
		},
	)

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"agent": {
				Type: agentType,
				Args: graphql.FieldConfigArgument{
					"id": {Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: r.resolveAgent,
			},
			"agents": {
				Type: graphql.NewNonNull(agentPage),
				Args: graphql.FieldConfigArgument{
					"cursor": {Type: graphql.String},
					"limit":  {Type: graphql.Int},
					"kind":   {Type: agentKindEnum},
				},
				Resolve: r.resolveAgents,
			},
			"trades": {
				Type: graphql.NewNonNull(tradePage),
				Args: graphql.FieldConfigArgument{
					"cursor":     {Type: graphql.String},
					"limit":      {Type: graphql.Int},
					"agentToken": {Type: addressScalar},
					"from":       {Type: timeScalar},
					"to":         {Type: timeScalar},
				},
				Resolve: r.resolveTrades,
			},
			"jobs": {
				Type: graphql.NewNonNull(jobPage),
				Args: graphql.FieldConfigArgument{
					"cursor": {Type: graphql.String},
					"limit":  {Type: graphql.Int},
					"status": {Type: jobStatusEnum},
				},
				Resolve: r.resolveJobs,
			},
			"stats": {
				Type:    graphql.NewNonNull(statsType),
				Resolve: r.resolveStats,
			},
		},
	})

	subscriptionType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Subscription",
		Fields: graphql.Fields{
			"newTrade": {
				Type: graphql.NewNonNull(tradeType),
				Args: graphql.FieldConfigArgument{
					"agentToken": {Type: addressScalar},
				},
				// Resolver runs when an upstream is already pushing
				// values onto p.Source — the websocket handler invokes
				// Run* helpers below to set up the channel.
				Resolve: passthroughResolver,
			},
			"jobUpdated": {
				Type: graphql.NewNonNull(jobType),
				Args: graphql.FieldConfigArgument{
					"id": {Type: graphql.ID},
				},
				Resolve: passthroughResolver,
			},
			"agentRegistered": {
				Type:    graphql.NewNonNull(agentType),
				Resolve: passthroughResolver,
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:        queryType,
		Subscription: subscriptionType,
	})
}

// passthroughResolver returns the source unchanged. graphql-go uses the
// same resolver entrypoint for queries and subscriptions; on the
// subscription path the source is the value pushed onto the channel by
// our websocket handler, so a no-op resolver is the correct default.
func passthroughResolver(p graphql.ResolveParams) (any, error) {
	return p.Source, nil
}

// RunNewTrade sets up a Subscription.newTrade stream. It returns a
// channel of store.Trade values that the websocket transport forwards
// to the client; the channel is closed when ctx is done or the bus
// reports end-of-stream.
//
// Filter handling:
//
//   - When agentToken is empty, every trades.* envelope is forwarded.
//   - When agentToken is set, the filter compares the bonding-curve
//     token address (case-insensitive) against the publisher payload's
//     `Curve` field. Trades for a different token are dropped at the
//     gateway, not pushed to the client.
func (r *Resolver) RunNewTrade(ctx context.Context, agentToken string) (<-chan store.Trade, error) {
	want := strings.ToLower(strings.TrimSpace(agentToken))
	var filter SubjectFilter
	if want != "" {
		filter = func(_ string, env *NATSEnvelope) bool {
			// The publisher emits the curve address inside the payload
			// JSON; we peek at the bytes here rather than fully
			// decoding twice (the typed decode happens after the
			// filter passes). A naive substring check is safe because
			// 0x-prefixed lowercase hex never collides with quoting
			// metacharacters in JSON.
			needle := []byte(`"Curve":"` + want + `"`)
			return bytesContainsCaseInsensitive(env.Payload, needle)
		}
	}
	return r.Bus.SubscribeTrades(ctx, filter)
}

// RunJobUpdated sets up a Subscription.jobUpdated stream. id, when
// non-empty, filters server-side to a specific on-chain job id; an
// empty id forwards every jobs.* publish.
func (r *Resolver) RunJobUpdated(ctx context.Context, id string) (<-chan store.Job, error) {
	want := strings.TrimSpace(id)
	var filter SubjectFilter
	if want != "" {
		filter = func(_ string, env *NATSEnvelope) bool {
			needle := []byte(`"JobId":"` + want + `"`)
			return bytesContainsCaseInsensitive(env.Payload, needle)
		}
	}
	return r.Bus.SubscribeJobs(ctx, filter)
}

// RunAgentRegistered sets up a Subscription.agentRegistered stream.
// Only agents.registered envelopes are forwarded — other agents.*
// events (metadata updates, score postings, controller transfers) are
// dropped here so this subscription stays narrow. Clients that want
// the broader stream can use the planned `agentEvent` field added in
// a later phase.
func (r *Resolver) RunAgentRegistered(ctx context.Context) (<-chan store.Agent, error) {
	filter := func(subject string, _ *NATSEnvelope) bool {
		return subject == "agents.registered"
	}
	return r.Bus.SubscribeAgents(ctx, filter)
}

// bytesContainsCaseInsensitive is a non-allocating byte-level search.
// Used by the SubjectFilter closures to avoid full JSON decode on
// every envelope just to check one field. The needle is expected
// lowercase; we lowercase each haystack byte on the fly. Marshalling
// rules for the gateway's BigInt/Address scalars guarantee lowercase
// canonical form on the wire, so this is sufficient for our filters.
func bytesContainsCaseInsensitive(haystack, needle []byte) bool {
	if len(needle) == 0 {
		return true
	}
	if len(haystack) < len(needle) {
		return false
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		ok := true
		for j := 0; j < len(needle); j++ {
			h := haystack[i+j]
			if h >= 'A' && h <= 'Z' {
				h += 'a' - 'A'
			}
			if h != needle[j] {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}

// errSchemaBuild is emitted when the schema constructor finds an
// inconsistency between SDL and runtime types; surfaced verbatim from
// graphql-go but with a stable prefix so callers can branch.
var errSchemaBuild = fmt.Errorf("graph: schema build")
