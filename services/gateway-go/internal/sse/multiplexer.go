// Package sse implements a Server-Sent Events fanout from the Titular
// NATS subjects (agents.*, jobs.*, tokens.*, trades.*) to per-connection
// HTTP `text/event-stream` clients.
//
// The package exists alongside the GraphQL subscription transport
// (internal/graph) rather than replacing it: SDK consumers that cannot
// open a websocket — restrictive corporate proxies, polyfilled WS in
// older browsers, or the simpler one-way usage shape that does not need
// bidirectional control frames — get a plain-HTTP event stream over the
// same wire data.
//
// # Architecture
//
//	indexer                 gateway                 client
//	┌──────────┐  publish   ┌──────────┐  GET /events
//	│Publisher │────────────▶│ Multi-  │◀──────────── browser / SDK
//	│(JetStr)  │  agents.>  │ plexer  │  text/event-stream
//	└──────────┘  jobs.>    │         │  ┌──────────┐
//	              tokens.>  │ fan-out │  │subscriber│
//	              trades.>  │ N:M     │  │(per conn)│
//	                        └──────────┘  └──────────┘
//
// The Multiplexer holds exactly one NATS subscription per top-level
// prefix and dispatches each incoming envelope to the set of registered
// subscribers whose subject filter matches. Matching is done once on the
// hub goroutine; subscribers receive the already-encoded SSE frame
// rather than decoding the envelope per-client. This keeps fan-out
// O(matched subscribers) rather than O(connected subscribers).
//
// # Backpressure
//
// Each subscriber owns a bounded buffered channel. When the buffer is
// full the hub drops the message for that subscriber and increments a
// per-subscriber drop counter; we do NOT block the hub goroutine because
// a single slow client must not stall delivery to every other client.
// This mirrors the same posture taken by the GraphQL subscription bus
// (see internal/graph/subscriptions.go subBufferSize = 64).
//
// # Replay
//
// Live tail uses core NATS (at-most-once from subscribe time forward),
// matching the documented contract on GraphQL subscriptions: clients
// reconnecting after a drop are expected to re-fetch the gap via the
// REST API, OR — for SSE specifically — by sending the `Last-Event-ID`
// header. The handler layer (handler.go) reads that header and uses a
// JetStream ephemeral consumer to replay messages with a stream
// sequence > Last-Event-ID before attaching the connection to the live
// hub. The multiplexer itself is unaware of replay; it only ever
// dispatches the live tail.
package sse

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/observability"
)

// SubjectPrefixes is the closed set of top-level prefixes the
// multiplexer subscribes to. Mirrors publisher.SubjectPrefixes in the
// indexer; we intentionally re-declare it here rather than importing
// across services to keep the gateway in its own dependency closure
// (same rationale as the NATSEnvelope re-declaration in
// internal/graph/subscriptions.go).
var SubjectPrefixes = []string{
	"agents.>",
	"jobs.>",
	"tokens.>",
	"trades.>",
}

// Event is the dispatched unit. Subject is the NATS subject the
// publisher wrote to (e.g. "trades.bought"); ID is the per-message
// identifier used as the SSE `id:` line — derived from the JetStream
// stream sequence when available, falling back to a hub-local monotonic
// counter when the multiplexer is wired against a core-NATS-only
// deployment. Data is the publisher's JSON envelope, forwarded verbatim
// so SDK consumers see exactly the same wire shape they would get over
// the GraphQL subscription transport.
type Event struct {
	ID      string
	Subject string
	Data    []byte
}

// Subscriber is the per-connection sink. It is the value returned from
// Multiplexer.Subscribe; the handler layer consumes Events from C and
// is responsible for serialising them onto the HTTP response.
//
// The buffer is sized at construction (Subscribe option) and is not
// resizable: changing it would let a misbehaving client absorb arbitrary
// memory by pretending to be slow. Drops accumulates the number of
// events the multiplexer has discarded for this subscriber since
// registration; the SSE handler logs the running total when the
// connection terminates so an operator can spot consistently slow
// readers.
type Subscriber struct {
	C      <-chan Event
	Drops  *atomic.Uint64
	cancel func()
}

// Close detaches the subscriber from the multiplexer and releases its
// channel. Safe to call multiple times. Always called by the handler on
// connection teardown — leaking a Subscriber into the registry would
// cause the hub to keep delivering to a dead reader, eventually
// saturating the buffer and bumping Drops indefinitely.
func (s *Subscriber) Close() {
	if s.cancel != nil {
		s.cancel()
	}
}

// Filter is the per-subscriber subject predicate. When non-empty, only
// events whose Subject matches at least one entry are dispatched to the
// subscriber. Patterns honour the standard NATS subject grammar:
//
//   - "*"  matches exactly one token (e.g. "trades.*" matches
//     "trades.bought" but NOT "trades.bonding.bought")
//   - ">"  matches one or more tokens (e.g. "trades.>" matches both)
//   - bare strings are exact matches (e.g. "trades.bought")
//
// An empty Filter dispatches every event the multiplexer receives.
type Filter []string

// match reports whether subject matches at least one pattern. An empty
// filter matches everything (the SDK convention for "no filter ⇒ all
// subjects" — same shape as publisher.SubjectFor).
func (f Filter) match(subject string) bool {
	if len(f) == 0 {
		return true
	}
	for _, p := range f {
		if subjectMatches(p, subject) {
			return true
		}
	}
	return false
}

// subjectMatches implements the NATS pattern grammar restricted to `*`
// and `>`. We re-implement it here rather than calling into nats.go's
// internal matcher because the function is unexported. The grammar is
// small enough that a token-by-token walk is clearer than a regex and
// easier to audit for ReDoS-style inputs.
func subjectMatches(pattern, subject string) bool {
	if pattern == "" {
		return false
	}
	if pattern == subject {
		return true
	}
	pTokens := strings.Split(pattern, ".")
	sTokens := strings.Split(subject, ".")
	for i, pt := range pTokens {
		switch pt {
		case ">":
			// `>` is only valid as the last token of a pattern. It
			// matches one or more remaining subject tokens.
			return i < len(sTokens)
		case "*":
			if i >= len(sTokens) {
				return false
			}
		default:
			if i >= len(sTokens) || sTokens[i] != pt {
				return false
			}
		}
	}
	// Pattern consumed; subject must be exactly the same length to
	// match (no trailing tokens allowed when there's no `>`).
	return len(pTokens) == len(sTokens)
}

// Multiplexer is the NATS-to-SSE fanout hub.
//
// Construction is via NewMultiplexer; the zero value is unsafe. Start
// must be called before any Subscribe so the underlying NATS
// subscriptions are live; Stop tears down the hub goroutine and closes
// every active subscriber's channel so handler goroutines exit
// promptly.
type Multiplexer struct {
	nc     *nats.Conn
	logger zerolog.Logger

	// subs is the set of registered subscribers, indexed by an
	// internal id. We use a map rather than a slice because Close on a
	// subscriber needs O(1) removal; the id is monotonic per
	// multiplexer instance.
	subsMu  sync.RWMutex
	subs    map[uint64]*subscriberEntry
	nextID  uint64
	natsSub []*nats.Subscription

	// seq is the hub-local fallback ID counter, used when the
	// JetStream sequence is not available (e.g. an embedded NATS
	// server in tests that never created the stream). Atomic so the
	// dispatch path doesn't take a lock.
	seq atomic.Uint64

	// stopped is set on Stop so concurrent Subscribe calls error out
	// cleanly instead of registering against a dead hub.
	stoppedMu sync.RWMutex
	stopped   bool
}

// subscriberEntry pairs a Subscriber with the goroutine-local pieces
// the hub needs (the writable channel end and the filter). Kept private
// so the public Subscriber struct stays read-only from the handler's
// perspective.
type subscriberEntry struct {
	id     uint64
	filter Filter
	ch     chan Event
	drops  *atomic.Uint64
}

// Config tunes Multiplexer construction.
type Config struct {
	// NATS is the connection to subscribe on. Required.
	NATS *nats.Conn

	// Logger is used for drop logging. Zero value is allowed (will
	// log to the default zerolog instance with the "sse" component
	// tag).
	Logger zerolog.Logger
}

// DefaultBufferSize is the per-subscriber channel capacity. Sized so a
// short network blip on the SSE response writer does not drop messages
// in the common case; longer outages will silently drop on overflow,
// which is the right call here — back-pressuring the hub would impact
// every other subscriber.
//
// Consistent with subBufferSize in the GraphQL subscription bus so the
// two transports degrade at the same operational threshold.
const DefaultBufferSize = 64

// NewMultiplexer constructs a multiplexer. Returns an error if cfg.NATS
// is nil; a closed connection is allowed because Start surfaces the
// transport error at subscription time.
func NewMultiplexer(cfg Config) (*Multiplexer, error) {
	if cfg.NATS == nil {
		return nil, errors.New("sse: nil NATS connection")
	}
	logger := cfg.Logger
	if reflect.DeepEqual(logger, zerolog.Logger{}) {
		logger = zerolog.Nop()
	}
	logger = logger.With().Str("component", "sse").Logger()
	return &Multiplexer{
		nc:     cfg.NATS,
		logger: logger,
		subs:   make(map[uint64]*subscriberEntry),
	}, nil
}

// Start opens the NATS subscriptions for every entry in
// SubjectPrefixes. Called once after construction; calling again is a
// no-op (idempotent so a partial-init recovery path can re-call without
// double-subscribing).
func (m *Multiplexer) Start() error {
	m.stoppedMu.Lock()
	defer m.stoppedMu.Unlock()
	if m.stopped {
		return errors.New("sse: multiplexer is stopped")
	}
	if len(m.natsSub) > 0 {
		return nil
	}
	for _, prefix := range SubjectPrefixes {
		sub, err := m.nc.Subscribe(prefix, m.dispatch)
		if err != nil {
			// Roll back any partial subscriptions before returning so a
			// failed Start leaves the hub in a clean stopped state.
			for _, s := range m.natsSub {
				_ = s.Unsubscribe()
			}
			m.natsSub = nil
			return err
		}
		m.natsSub = append(m.natsSub, sub)
	}
	return nil
}

// Stop closes every NATS subscription and every registered subscriber's
// channel. Safe to call multiple times. After Stop returns, Subscribe
// will return ErrStopped — handler goroutines must observe their
// channel close and exit before the gateway's graceful shutdown
// timeout.
func (m *Multiplexer) Stop() {
	m.stoppedMu.Lock()
	if m.stopped {
		m.stoppedMu.Unlock()
		return
	}
	m.stopped = true
	subs := m.natsSub
	m.natsSub = nil
	m.stoppedMu.Unlock()

	for _, s := range subs {
		_ = s.Unsubscribe()
	}

	// Snapshot the subscriber map under the registry lock, then close
	// each channel. Closing under the lock would deadlock if a
	// dispatch goroutine were still fanning out (it holds the read
	// lock); the snapshot-then-close pattern is goroutine-safe here
	// because dispatch sends are non-blocking — a fanout into a
	// closed channel falls into the default branch of the `select`
	// rather than panicking — and safeClose tolerates a double close.
	m.subsMu.Lock()
	entries := make([]*subscriberEntry, 0, len(m.subs))
	for _, e := range m.subs {
		entries = append(entries, e)
	}
	m.subs = make(map[uint64]*subscriberEntry)
	m.subsMu.Unlock()

	for _, e := range entries {
		safeClose(e.ch)
	}
}

// ErrStopped is returned by Subscribe when the multiplexer has been
// stopped. Handlers should treat this as "service shutting down" and
// terminate the SSE response gracefully.
var ErrStopped = errors.New("sse: multiplexer stopped")

// SubscribeOptions tunes a single Subscribe call.
type SubscribeOptions struct {
	// Filter restricts dispatch to events whose subject matches at
	// least one pattern. nil/empty filter receives every event.
	Filter Filter

	// Buffer overrides DefaultBufferSize for this subscriber. Values
	// <= 0 fall back to the default.
	Buffer int
}

// Subscribe registers a new SSE client and returns a Subscriber whose
// channel will be fed every matching event until Close (or hub Stop) is
// called. The returned Subscriber is owned by the caller; its channel
// is closed on either Close or Stop.
func (m *Multiplexer) Subscribe(opts SubscribeOptions) (*Subscriber, error) {
	m.stoppedMu.RLock()
	if m.stopped {
		m.stoppedMu.RUnlock()
		return nil, ErrStopped
	}
	m.stoppedMu.RUnlock()

	buf := opts.Buffer
	if buf <= 0 {
		buf = DefaultBufferSize
	}
	ch := make(chan Event, buf)
	drops := new(atomic.Uint64)

	m.subsMu.Lock()
	id := m.nextID
	m.nextID++
	entry := &subscriberEntry{
		id:     id,
		filter: opts.Filter,
		ch:     ch,
		drops:  drops,
	}
	m.subs[id] = entry
	m.subsMu.Unlock()

	// Bump the live-connection gauge for the Grafana SSE panel. We
	// pair the +1 here with a -1 in cancel so the gauge reflects the
	// instantaneous fan-out and not a cumulative count. observability
	// hands back a no-op instrument when the meter provider has not
	// been built, so this is allocation-light when metrics are off.
	observability.SSEConnections().Add(context.Background(), 1)

	// once guards the gauge decrement so a Subscriber.Close() racing
	// with multiplexer.Stop() (both call into this cancel via
	// safeClose) does not double-decrement.
	var once sync.Once
	cancel := func() {
		m.subsMu.Lock()
		_, ok := m.subs[id]
		delete(m.subs, id)
		m.subsMu.Unlock()
		if ok {
			safeClose(ch)
		}
		once.Do(func() {
			observability.SSEConnections().Add(context.Background(), -1)
		})
	}

	return &Subscriber{
		C:      ch,
		Drops:  drops,
		cancel: cancel,
	}, nil
}

// dispatch is the NATS subscription callback. Runs on the NATS
// connection's read thread; we keep it cheap by doing the filter check
// inline and a non-blocking send to each matched subscriber. Anything
// blocking here would back-pressure the entire NATS client, which would
// drop messages for OTHER consumers in the gateway (most importantly,
// the GraphQL subscription bus that shares the same connection).
//
// Trace continuity: when the publisher injected a W3C TraceContext on
// the NATS header, we extract it here and open a kind=consumer span as
// the immediate child of the producer span so a backend can stitch the
// two halves into a single distributed trace. Header is best-effort —
// when missing/malformed the consume span starts a new trace, matching
// the propagator's "parent is optional" contract.
func (m *Multiplexer) dispatch(msg *nats.Msg) {
	ctx := observability.ExtractNATSContext(context.Background(), msg)
	tracer := otel.Tracer("titular/gateway/sse")
	ctx, span := tracer.Start(ctx, msg.Subject+" consume",
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			attribute.String("messaging.system", "nats"),
			attribute.String("messaging.source.name", msg.Subject),
			attribute.String("messaging.operation", "receive"),
		),
	)
	defer span.End()
	_ = ctx // future child spans (per-subscriber send) will hang off this.

	id := m.eventID(msg)
	ev := Event{
		ID:      id,
		Subject: msg.Subject,
		Data:    msg.Data,
	}

	// Take the read lock for the snapshot of subscribers. Sends are
	// non-blocking, so we hold the lock for at most O(N) channel
	// pokes; new subscribers added during dispatch do not see this
	// message, which is the same semantic GraphQL subscriptions have.
	m.subsMu.RLock()
	defer m.subsMu.RUnlock()
	for _, entry := range m.subs {
		if !entry.filter.match(msg.Subject) {
			continue
		}
		sendNonBlocking(entry, ev, m.logger, msg.Subject)
	}
}

// sendNonBlocking attempts a non-blocking send onto entry.ch and
// records a drop if the buffer is full. It also recovers from the
// `send on closed channel` panic that can race with Stop / Close
// calling close() on the channel; that race is benign — by the time
// we hit it, the subscriber is being torn down and we just want to
// avoid taking the gateway down with us.
func sendNonBlocking(entry *subscriberEntry, ev Event, logger zerolog.Logger, subject string) {
	defer func() {
		if r := recover(); r != nil {
			// Channel was closed concurrently; treat as a drop and
			// move on. We do NOT log because Close/Stop are
			// expected operations, not errors.
			_ = r
		}
	}()
	select {
	case entry.ch <- ev:
	default:
		entry.drops.Add(1)
		logger.Warn().
			Uint64("subscriber_id", entry.id).
			Str("subject", subject).
			Uint64("drops", entry.drops.Load()).
			Msg("sse subscriber buffer full; dropping event")
	}
}

// eventID resolves the SSE `id:` value for an incoming NATS message.
// JetStream messages carry a stream sequence in their metadata, which
// is the canonical resumable identifier; core-NATS messages do not, so
// we fall back to a hub-local monotonic counter. The fallback is
// non-resumable across multiplexer restarts — the Last-Event-ID replay
// path in the handler handles that case explicitly by refusing replay
// when the current event format is the fallback shape.
//
// Format:
//
//	js:<stream_seq>     when the message has JetStream metadata
//	hub:<local_counter> otherwise
//
// Both forms are stable strings, so the SSE handler can echo them back
// in `id:` lines and decode them on a subsequent Last-Event-ID without
// ambiguity.
func (m *Multiplexer) eventID(msg *nats.Msg) string {
	if md, err := msg.Metadata(); err == nil && md != nil {
		// Stream sequence is monotonic per stream and survives client
		// reconnects; this is the SSE-correct id.
		seq := md.Sequence.Stream
		if seq > 0 {
			return formatJetStreamID(seq)
		}
	}
	return formatHubID(m.seq.Add(1))
}

// formatJetStreamID / formatHubID are exported (lowercase but
// package-visible) so the handler's Last-Event-ID parser can match
// against the same prefixes without re-implementing the format.
func formatJetStreamID(seq uint64) string { return "js:" + uitoa(seq) }
func formatHubID(seq uint64) string       { return "hub:" + uitoa(seq) }

// uitoa is strconv.FormatUint(x, 10) inlined to avoid the import for a
// single call site on a hot path.
func uitoa(x uint64) string {
	if x == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for x > 0 {
		i--
		buf[i] = byte('0' + x%10)
		x /= 10
	}
	return string(buf[i:])
}

// parseJetStreamID extracts the stream sequence from an event id of
// the form "js:<digits>". Returns (seq, true) on success; any other
// shape returns (0, false). Used by the handler when parsing
// Last-Event-ID for replay.
func parseJetStreamID(id string) (uint64, bool) {
	const prefix = "js:"
	if !strings.HasPrefix(id, prefix) {
		return 0, false
	}
	rest := id[len(prefix):]
	if rest == "" {
		return 0, false
	}
	var n uint64
	for i := 0; i < len(rest); i++ {
		c := rest[i]
		if c < '0' || c > '9' {
			return 0, false
		}
		n = n*10 + uint64(c-'0')
	}
	return n, true
}

// safeClose closes ch idempotently. Calling close on an
// already-closed channel panics; we swallow that so Stop and Close can
// race without crashing the gateway.
func safeClose(ch chan Event) {
	defer func() { _ = recover() }()
	close(ch)
}

// SubjectsFromCSV parses a comma-separated subject filter list,
// trimming whitespace and dropping empty entries. Returns nil for an
// empty input so a missing/blank ?subjects= query parameter is
// indistinguishable from an unfiltered subscription. Reused by the
// handler layer.
//
// Patterns whose FIRST token is `>` or `*` are rejected: they would
// let a client subscribe across every registered prefix (or worse,
// every subject NATS knows about) regardless of the multiplexer's
// closed-set posture. The handler maps those rejections to 400; the
// request never reaches the live tail. Patterns that narrow within a
// known prefix (e.g. `agents.*`, `trades.>`) are preserved unchanged.
func SubjectsFromCSV(s string) Filter {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make(Filter, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if err := validateSubjectPattern(p, SubjectPrefixes); err != nil {
			// One bad pattern poisons the whole list. The alternative
			// (silently dropping it) would let a client think a broad
			// `>` is being honoured and miss the more specific entries
			// they typed alongside it; explicit rejection is safer.
			return nil
		}
		out = append(out, p)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// ParseSubjectsCSV is the validating parser the SSE handler calls so
// it can return a 400 rather than silently falling back to the
// unfiltered tail when the client sent a banned pattern. Returns
// (nil, nil) for blank input — same "no filter" semantic as
// SubjectsFromCSV — and (nil, err) for any rejected pattern.
func ParseSubjectsCSV(s string) (Filter, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	out := make(Filter, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if err := validateSubjectPattern(p, SubjectPrefixes); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	if len(out) == 0 {
		return nil, nil
	}
	return out, nil
}

// validateSubjectPattern rejects subject patterns that would let a
// client opt out of the multiplexer's closed-set fanout. The rule is
// narrow: the FIRST token must not be a wildcard. Once the first token
// is a literal that matches one of the registered prefixes (agents,
// jobs, tokens, trades), the pattern is allowed to use `*` and `>` for
// further narrowing.
//
// allowedPrefixes is the registered top-level prefix list (i.e.
// SubjectPrefixes); we extract the leading token of each entry and
// compare against the pattern's leading token. Passing nil disables
// the prefix check (allows any non-wildcard leading token), which is
// only useful for unit tests.
func validateSubjectPattern(pattern string, allowedPrefixes []string) error {
	if pattern == "" {
		return errors.New("sse: empty subject pattern")
	}
	first := pattern
	if idx := strings.Index(pattern, "."); idx >= 0 {
		first = pattern[:idx]
	}
	switch first {
	case ">":
		return errors.New("sse: bare `>` not allowed; subscribe within a registered prefix")
	case "*":
		return errors.New("sse: leading `*` not allowed; subscribe within a registered prefix")
	}
	if len(allowedPrefixes) == 0 {
		return nil
	}
	for _, prefix := range allowedPrefixes {
		// Leading token of each registered prefix (e.g. "agents.>" -> "agents").
		head := prefix
		if idx := strings.Index(prefix, "."); idx >= 0 {
			head = prefix[:idx]
		}
		if head == first {
			return nil
		}
	}
	return errors.New("sse: pattern leading token does not match any registered prefix")
}

