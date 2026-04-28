package publisher

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/nats-io/nats.go"

	"github.com/0xDevNinja/titular/services/indexer-go/internal/decoder"
)

// Publisher fans decoded EVM events out to a NATS JetStream context. It
// implements [decoder.Handler] so the indexer's main loop can hand a
// [decoder.DecodedEvent] off without knowing anything about NATS.
//
// Concurrency: Publisher itself is stateless; safe for concurrent calls
// as long as the underlying [JetStreamContext] is. The standard nats.go
// JetStreamContext is.
//
// Idempotency: every publish carries a `Nats-Msg-Id` header derived from
// `tx_hash:log_index`. JetStream uses that as a per-stream dedup key
// inside its dedup window (configured by [EnsureStream]). This gives
// at-least-once delivery without duplicate messages on operator restart
// or block re-fetch within the window — the indexer can therefore safely
// re-process a confirmed block that crashed mid-flight (the subscriber
// reports highest-confirmed each iteration; the publisher absorbs the
// overlap).
//
// Publishes outside the dedup window will produce duplicates; consumers
// that need stronger guarantees should key on the same (tx_hash,
// log_index) pair carried in the JSON envelope.
type Publisher struct {
	js   JetStreamContext
	opts Options
}

// JetStreamContext is the subset of [nats.JetStreamContext] this package
// needs. Defined as an interface so tests can stub the publish path
// without standing up a server, and so future migration to the new
// `jetstream` API in nats.go does not ripple through callers.
type JetStreamContext interface {
	// Publish blocks until the server acknowledges the message. The
	// returned [nats.PubAck] surfaces the stream sequence and the dedup
	// flag — Publisher does not currently inspect either, but the
	// interface keeps the door open for metrics later.
	Publish(subject string, data []byte, opts ...nats.PubOpt) (*nats.PubAck, error)
}

// Options configures Publisher. Zero value is valid: it uses the package
// defaults (synchronous publish, JSON encoding, no extra labels).
type Options struct {
	// SkipUnmapped, when true, silently drops events whose qualified
	// name has no entry in [eventSubjects]. The default (false) returns
	// an error so a regen that adds a tracked event before its subject
	// is assigned cannot ship to production unnoticed.
	//
	// In tests we usually leave this at false; in production the
	// indexer wraps Publisher with a tolerant handler that converts
	// [ErrUnmappedEvent] into a counter-only no-op.
	SkipUnmapped bool
}

// New constructs a Publisher. js must be non-nil. The Options struct is
// copied; callers cannot mutate the publisher's behaviour after
// construction.
func New(js JetStreamContext, opts Options) (*Publisher, error) {
	if js == nil {
		return nil, errors.New("publisher: nil JetStreamContext")
	}
	return &Publisher{js: js, opts: opts}, nil
}

// ErrUnmappedEvent is returned by [Publisher.OnEvent] when the decoded
// event name has no entry in the subject table and [Options.SkipUnmapped]
// is false. Callers can use errors.Is to differentiate this from a
// transport failure.
var ErrUnmappedEvent = errors.New("publisher: event has no subject mapping")

// envelope is the JSON wire format. We wrap the raw abigen payload
// rather than emitting it bare so consumers can correlate to chain data
// without hunting through the binding-specific fields:
//
//   - Name is the qualified event name (same value used in metrics).
//   - Block / Tx / LogIndex are the chain coordinates.
//   - Payload is the abigen struct exactly as decoded; the wire format
//     is determined by the struct tags on the binding (which abigen
//     emits with `json:"…"` for every field).
//
// We deliberately do not version the envelope: subject names are the
// evolution boundary (we add `tokens.fee_split.v2` rather than mutate
// `tokens.fee_split`). That keeps consumers simpler and matches the
// guidance in [SubjectFor].
type envelope struct {
	Name        string `json:"name"`
	BlockNumber uint64 `json:"block_number"`
	TxHash      string `json:"tx_hash"`
	LogIndex    uint   `json:"log_index"`
	Payload     any    `json:"payload"`
}

// OnEvent satisfies [decoder.Handler]. It maps the event name to a
// subject, JSON-encodes the envelope, and publishes synchronously with a
// dedup ID. Returned errors are wrapped with the event name so an
// operator log line points at the contract that produced the failure.
//
// Behaviour:
//
//   - Unknown event name: returns [ErrUnmappedEvent], unless
//     [Options.SkipUnmapped] is true, in which case OnEvent returns nil
//     and the message is dropped.
//   - Encoding failure: returned verbatim. JSON encoding only fails on
//     unsupported types (channels, funcs); abigen structs are plain old
//     data so this should be impossible in practice but is surfaced
//     rather than swallowed.
//   - JetStream Publish failure (network, no-stream-bound, server-side
//     error): wrapped and returned. The indexer's per-block pipeline
//     aborts the block on a non-nil handler error, which is the
//     intended back-pressure: we'd rather pause than skip events.
func (p *Publisher) OnEvent(event decoder.DecodedEvent) error {
	subject, ok := SubjectFor(event.Name)
	if !ok {
		if p.opts.SkipUnmapped {
			return nil
		}
		return fmt.Errorf("%w: %s", ErrUnmappedEvent, event.Name)
	}

	body, err := json.Marshal(envelope{
		Name:        event.Name,
		BlockNumber: event.Raw.BlockNumber,
		TxHash:      event.Raw.TxHash.Hex(),
		LogIndex:    event.Raw.Index,
		Payload:     event.Payload,
	})
	if err != nil {
		return fmt.Errorf("publisher: encode %s: %w", event.Name, err)
	}

	msgID := DedupID(event.Raw.TxHash.Hex(), event.Raw.Index)
	if _, err := p.js.Publish(subject, body, nats.MsgId(msgID)); err != nil {
		return fmt.Errorf("publisher: publish %s to %s: %w", event.Name, subject, err)
	}
	return nil
}

// DedupID constructs the JetStream `Nats-Msg-Id` header value for a log.
// Format: `<lowercase tx hash>:<log index>`. The format is stable: any
// future change requires an explicit migration because deployed
// environments rely on it for the dedup-window guarantee.
//
// Exported so tests and ad-hoc tools can produce the same key without
// re-deriving the format.
func DedupID(txHash string, logIndex uint) string {
	return txHash + ":" + strconv.FormatUint(uint64(logIndex), 10)
}

// Compile-time assertion: Publisher satisfies decoder.Handler.
var _ decoder.Handler = (*Publisher)(nil)
