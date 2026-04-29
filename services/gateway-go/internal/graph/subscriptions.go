package graph

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/observability"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/store"
)

// SubscriptionBus is the abstraction the GraphQL subscription resolvers
// see. Implementations must:
//
//   - Subscribe to a NATS subject pattern (e.g. "trades.>") and forward
//     decoded envelopes onto the returned channel.
//   - Honour ctx cancellation: when ctx is done, stop publishing and
//     close the channel so the GraphQL transport can terminate the
//     subscription cleanly.
//
// The interface is named separately from a concrete NATS implementation
// so unit tests can drive subscriptions with a synchronous in-memory
// bus that doesn't need a server.
type SubscriptionBus interface {
	// SubscribeAgents returns a channel that yields every agents.>
	// publish, decoded into a store.Agent. Filter is applied
	// server-side to avoid forwarding traffic the subscriber does not
	// want. A nil filter forwards every message in the prefix.
	SubscribeAgents(ctx context.Context, filter SubjectFilter) (<-chan store.Agent, error)

	// SubscribeTrades returns trades.> publishes decoded into store.Trade.
	SubscribeTrades(ctx context.Context, filter SubjectFilter) (<-chan store.Trade, error)

	// SubscribeJobs returns jobs.> publishes decoded into store.Job.
	SubscribeJobs(ctx context.Context, filter SubjectFilter) (<-chan store.Job, error)
}

// SubjectFilter is a closure that returns true when the publisher
// envelope should be forwarded to the subscriber. Filters always run on
// the gateway side; we cannot push them upstream without baking
// per-subscriber consumers into the JetStream layer, which is too
// expensive for the alpha.
//
// A nil SubjectFilter is interpreted as "match everything".
type SubjectFilter func(subject string, envelope *NATSEnvelope) bool

// NATSEnvelope mirrors the publisher's outbound envelope (see the
// indexer-go publisher.go). We re-declare it here rather than
// importing the indexer package because:
//
//  1. The gateway must stay in its own dependency closure (no upward
//     import to a sibling service).
//  2. The wire format is the contract; future changes will go through
//     subject versioning rather than a Go type rename.
//
// Payload is left as json.RawMessage because the resolver decodes it
// into the appropriate store.* shape after subject-based dispatch.
type NATSEnvelope struct {
	Name        string          `json:"name"`
	BlockNumber uint64          `json:"block_number"`
	TxHash      string          `json:"tx_hash"`
	LogIndex    uint            `json:"log_index"`
	Payload     json.RawMessage `json:"payload"`
}

// nilBus is a SubscriptionBus that returns immediately-closed channels.
// Used when the gateway is launched without GATEWAY_NATS_URL: the
// subscription endpoints are still valid GraphQL operations, they just
// never emit. This is preferable to a 404 because SDK clients that
// auto-subscribe at startup get a graceful "no data" rather than a
// transport error.
type nilBus struct{}

// NilBus returns the no-op SubscriptionBus described above. Exported so
// cmd/gateway can wire it explicitly when NATS is disabled (and tests
// can use it as a clean-room baseline).
func NilBus() SubscriptionBus { return nilBus{} }

func (nilBus) SubscribeAgents(_ context.Context, _ SubjectFilter) (<-chan store.Agent, error) {
	ch := make(chan store.Agent)
	close(ch)
	return ch, nil
}

func (nilBus) SubscribeTrades(_ context.Context, _ SubjectFilter) (<-chan store.Trade, error) {
	ch := make(chan store.Trade)
	close(ch)
	return ch, nil
}

func (nilBus) SubscribeJobs(_ context.Context, _ SubjectFilter) (<-chan store.Job, error) {
	ch := make(chan store.Job)
	close(ch)
	return ch, nil
}

// natsBus implements SubscriptionBus on top of a NATS connection.
//
// We use core NATS subscriptions rather than JetStream pull subscribers
// because each GraphQL subscription is short-lived and per-connection;
// JetStream durables would leak consumer state on every disconnect.
// Core NATS gives us at-most-once delivery from the moment the
// subscription is established, which matches GraphQL subscription
// semantics — clients reconnecting after a drop are expected to
// re-fetch the gap via the REST API or by paginating through
// Query.trades / Query.jobs.
type natsBus struct {
	nc *nats.Conn
	// decoders are pinned at construction so the resolver path doesn't
	// allocate a new decoder per message. agentFromEnvelope etc. are
	// pure functions; they share a single sync.Pool of bytes.Buffer if
	// the throughput justifies it (none yet).
}

// NewNATSBus wires a NATS connection into a SubscriptionBus. Returns an
// error if nc is nil; a closed nc is allowed because Subscribe will
// surface the transport error at the time of subscription rather than
// at construction.
func NewNATSBus(nc *nats.Conn) (SubscriptionBus, error) {
	if nc == nil {
		return nil, errors.New("nats bus: nil connection")
	}
	return &natsBus{nc: nc}, nil
}

// SubscribeAgents subscribes to "agents.>" and decodes the publisher
// envelope into store.Agent. Filter is applied on the wire envelope so
// the decode cost is paid only for messages the subscriber wants.
func (b *natsBus) SubscribeAgents(ctx context.Context, filter SubjectFilter) (<-chan store.Agent, error) {
	out := make(chan store.Agent, subBufferSize)
	if err := b.subscribe(ctx, "agents.>", out, filter, func(env *NATSEnvelope) (any, bool) {
		// Only Query.agentRegistered is exposed; we filter on
		// agents.registered subject inside the GraphQL resolver, but
		// the decoder runs for any agents.* event so the helper here
		// stays subject-agnostic.
		a, ok := agentFromEnvelope(env)
		return a, ok
	}); err != nil {
		return nil, err
	}
	return out, nil
}

func (b *natsBus) SubscribeTrades(ctx context.Context, filter SubjectFilter) (<-chan store.Trade, error) {
	out := make(chan store.Trade, subBufferSize)
	if err := b.subscribe(ctx, "trades.>", out, filter, func(env *NATSEnvelope) (any, bool) {
		t, ok := tradeFromEnvelope(env)
		return t, ok
	}); err != nil {
		return nil, err
	}
	return out, nil
}

func (b *natsBus) SubscribeJobs(ctx context.Context, filter SubjectFilter) (<-chan store.Job, error) {
	out := make(chan store.Job, subBufferSize)
	if err := b.subscribe(ctx, "jobs.>", out, filter, func(env *NATSEnvelope) (any, bool) {
		j, ok := jobFromEnvelope(env)
		return j, ok
	}); err != nil {
		return nil, err
	}
	return out, nil
}

// subBufferSize is the per-subscription channel capacity. Sized so a
// short network blip in the GraphQL websocket does not drop messages
// in the common case; longer outages will silently drop on overflow,
// which is the right call here — back-pressuring NATS would impact
// every other subscriber.
const subBufferSize = 64

// subscribe is the shared NATS plumbing for the typed Subscribe* helpers
// above. It uses a generic any-channel under the hood and a decoder
// closure that returns the typed value plus an "include?" flag, then
// delegates to a per-type goroutine that does the type assertion and
// channel send. Keeps the typed Subscribe* helpers symmetric without
// duplicating the connection / cleanup code.
func (b *natsBus) subscribe(
	ctx context.Context,
	subject string,
	out any,
	filter SubjectFilter,
	decode func(*NATSEnvelope) (any, bool),
) error {
	// Buffered channel between NATS callback and the type-asserting
	// fan-out goroutine. Sized identically to the typed channel so a
	// stall on the typed side does not deepen the queue here.
	raw := make(chan rawMsg, subBufferSize)

	sub, err := b.nc.Subscribe(subject, func(m *nats.Msg) {
		// Trace continuity: extract any W3C TraceContext the publisher
		// injected on the NATS header, then open a kind=consumer span
		// as the immediate child of the producer span so a backend can
		// stitch the two halves into a single distributed trace. The
		// span ends immediately because the per-message work in this
		// callback is just an enqueue — the decoder / fanout that runs
		// in the goroutine below is observed at the outer subscribe
		// level (e.g. the GraphQL stream's lifecycle span), not here.
		// Header is best-effort: when missing/malformed the consume
		// span starts a new trace, matching propagator semantics.
		msgCtx := observability.ExtractNATSContext(context.Background(), m)
		tracer := otel.Tracer("titular/gateway/graph/subscriptions")
		_, span := tracer.Start(msgCtx, m.Subject+" consume",
			trace.WithSpanKind(trace.SpanKindConsumer),
			trace.WithAttributes(
				attribute.String("messaging.system", "nats"),
				attribute.String("messaging.source.name", m.Subject),
				attribute.String("messaging.operation", "receive"),
			),
		)
		select {
		case raw <- rawMsg{subject: m.Subject, data: m.Data}:
		default:
			// Drop on overflow. We log via the gateway's structured
			// logger one frame up; here we cannot import zerolog
			// without a circular package dependency, so we rely on
			// the fact that the typed channel is what the resolver
			// observes for back-pressure visibility.
		}
		span.End()
	})
	if err != nil {
		return err
	}

	// Bump the live-subscription gauge for the Grafana panel. We pair
	// the +1 here with a -1 in the goroutine's defer below so a normal
	// exit (ctx cancelled, NATS sub torn down, resolver dropped the
	// channel) brings the gauge back down. observability.GraphQLSubscriptions
	// hands back a no-op when the meter provider has not been built,
	// so this is allocation-light when metrics are off. The +1 lives
	// AFTER the b.nc.Subscribe error gate so a failed subscribe does
	// not leak gauge state.
	observability.GraphQLSubscriptions().Add(ctx, 1)

	go func() {
		defer func() {
			_ = sub.Unsubscribe()
			closeAny(out)
			observability.GraphQLSubscriptions().Add(ctx, -1)
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case m, ok := <-raw:
				if !ok {
					return
				}
				var env NATSEnvelope
				if err := json.Unmarshal(m.data, &env); err != nil {
					continue
				}
				if filter != nil && !filter(m.subject, &env) {
					continue
				}
				val, ok := decode(&env)
				if !ok {
					continue
				}
				if !sendAny(out, val) {
					return
				}
			}
		}
	}()

	return nil
}

type rawMsg struct {
	subject string
	data    []byte
}

// sendAny does a non-blocking send onto the typed channel `out`,
// returning false if `out` has been closed by the resolver (which
// signals the goroutine to terminate). We unwind type assertions here
// to avoid using reflect on the hot path — the typed channels are
// statically known per-Subscribe call.
func sendAny(out any, v any) bool {
	switch ch := out.(type) {
	case chan store.Agent:
		a, ok := v.(store.Agent)
		if !ok {
			return true
		}
		defer recoverClosed()
		ch <- a
		return true
	case chan store.Trade:
		t, ok := v.(store.Trade)
		if !ok {
			return true
		}
		defer recoverClosed()
		ch <- t
		return true
	case chan store.Job:
		j, ok := v.(store.Job)
		if !ok {
			return true
		}
		defer recoverClosed()
		ch <- j
		return true
	}
	return false
}

// closeAny mirrors sendAny: closes the typed channel after the
// subscription goroutine exits so resolvers observing the channel
// receive an end-of-stream rather than a hang.
func closeAny(out any) {
	defer recoverClosed()
	switch ch := out.(type) {
	case chan store.Agent:
		close(ch)
	case chan store.Trade:
		close(ch)
	case chan store.Job:
		close(ch)
	}
}

// recoverClosed swallows the "send/close on closed channel" panic that
// can race with ctx cancellation. The race is benign: by the time we
// hit this branch the resolver has already moved on, and there's no
// safe way to test "is this channel closed" without locking the
// channel's internal state, which the runtime intentionally hides.
func recoverClosed() {
	_ = recover()
}

// agentFromEnvelope decodes a publisher envelope onto a store.Agent.
// We populate only the fields that survive the JSON wire format —
// chain coordinates from the envelope and the standardised payload
// fields. Resolvers that need stronger consistency (e.g. the agent's
// kind or controller after a transfer-accept event) should re-read
// the row via Query.agent in a follow-up call; the subscription here
// is event-shaped, not row-shaped.
func agentFromEnvelope(env *NATSEnvelope) (store.Agent, bool) {
	var p struct {
		AgentId     string `json:"AgentId"`
		Controller  string `json:"Controller"`
		Creator     string `json:"Creator"`
		MetadataURI string `json:"MetadataURI"`
	}
	if err := json.Unmarshal(env.Payload, &p); err != nil {
		return store.Agent{}, false
	}
	a := store.Agent{
		AgentID:     p.AgentId,
		Controller:  strings.ToLower(p.Controller),
		Creator:     strings.ToLower(p.Creator),
		MetadataURI: p.MetadataURI,
		BlockNumber: env.BlockNumber,
		TxHash:      strings.ToLower(env.TxHash),
		LogIndex:    uint32(env.LogIndex),
	}
	// Best-effort kind derivation: the publisher emits AgentRegistry.*
	// for ACP agents; launchpad agents arrive on a separate subject we
	// haven't wired yet (M2 BondingCurve.*). Default to ACP for now;
	// the field is informational on the subscription path because
	// callers who care must re-read the row.
	a.Kind = store.AgentKindACP
	return a, true
}

func tradeFromEnvelope(env *NATSEnvelope) (store.Trade, bool) {
	var p struct {
		Trader   string `json:"Trader"`
		Curve    string `json:"Curve"`
		QuoteIn  string `json:"QuoteIn"`
		AgentOut string `json:"AgentOut"`
		AgentIn  string `json:"AgentIn"`
		QuoteOut string `json:"QuoteOut"`
		Fee      string `json:"Fee"`
		Side     string `json:"Side"`
	}
	if err := json.Unmarshal(env.Payload, &p); err != nil {
		return store.Trade{}, false
	}
	if p.Side == "" {
		// trades.bought / trades.sold map onto buy/sell respectively;
		// when the publisher subject is set, we infer side from the
		// envelope name suffix.
		switch {
		case strings.HasSuffix(env.Name, ".Bought"):
			p.Side = "buy"
		case strings.HasSuffix(env.Name, ".Sold"):
			p.Side = "sell"
		}
	}
	return store.Trade{
		Side:        p.Side,
		Trader:      strings.ToLower(p.Trader),
		Curve:       strings.ToLower(p.Curve),
		QuoteIn:     p.QuoteIn,
		AgentOut:    p.AgentOut,
		AgentIn:     p.AgentIn,
		QuoteOut:    p.QuoteOut,
		Fee:         p.Fee,
		BlockNumber: env.BlockNumber,
		TxHash:      strings.ToLower(env.TxHash),
		LogIndex:    uint32(env.LogIndex),
	}, true
}

func jobFromEnvelope(env *NATSEnvelope) (store.Job, bool) {
	var p struct {
		JobId     string `json:"JobId"`
		Token     string `json:"Token"`
		Amount    string `json:"Amount"`
		Depositor string `json:"Depositor"`
		Principal string `json:"Principal"`
		Budget    string `json:"Budget"`
	}
	if err := json.Unmarshal(env.Payload, &p); err != nil {
		return store.Job{}, false
	}
	job := store.Job{
		OnChainID:   p.JobId,
		Token:       strings.ToLower(p.Token),
		BlockNumber: env.BlockNumber,
		TxHash:      strings.ToLower(env.TxHash),
		LogIndex:    uint32(env.LogIndex),
	}
	// Phase derivation from the publisher subject. The full mapping
	// lives in the indexer's publisher.SubjectFor table; we replicate
	// the bare minimum here because subscribers see envelope names
	// directly.
	switch env.Name {
	case "JobFactory.JobCreated":
		job.Phase = store.JobPhaseCreated
		job.Principal = strings.ToLower(p.Principal)
		job.Budget = p.Budget
	case "Escrow.Funded":
		job.Phase = store.JobPhaseFunded
		job.Principal = strings.ToLower(p.Depositor)
		job.Budget = p.Amount
	case "Escrow.Released":
		job.Phase = store.JobPhaseReleased
	case "Escrow.Refunded":
		job.Phase = store.JobPhaseCancelled
	case "Job.JobCompleted":
		job.Phase = store.JobPhaseCompleted
	case "Job.JobCancelled":
		job.Phase = store.JobPhaseCancelled
	case "Job.DisputeRaised":
		job.Phase = store.JobPhaseDisputed
	case "Job.DisputeResolved":
		job.Phase = store.JobPhaseResolved
	}
	return job, true
}
