package publisher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/0xDevNinja/titular/services/indexer-go/internal/decoder"
	"github.com/0xDevNinja/titular/services/indexer-go/internal/observability"
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

	// PublishMsg variant carrying a full [nats.Msg] so the publisher
	// can stamp W3C TraceContext headers (and any future correlation
	// metadata) onto the wire message. Both methods accept the same
	// PubOpt set; the primary semantic difference is the header
	// payload being preserved across the boundary.
	PublishMsg(msg *nats.Msg, opts ...nats.PubOpt) (*nats.PubAck, error)
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
//   - Payload is the decoded event with `*big.Int` fields wrapped as
//     decimal strings (see [encodePayload]); JS consumers (gateway, SDK)
//     would otherwise lose precision past 2^53 on the default
//     `big.Int.MarshalJSON` numeric form.
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
//
// The ctx is forwarded to JetStream via [nats.Context] so a stalled
// server cannot block the indexer's main loop indefinitely. ctx may be
// nil for callers that have no per-block deadline; in that case Publish
// uses the underlying client's default timeout.
func (p *Publisher) OnEvent(ctx context.Context, event decoder.DecodedEvent) error {
	if ctx == nil {
		ctx = context.Background()
	}

	subject, ok := SubjectFor(event.Name)
	if !ok {
		if p.opts.SkipUnmapped {
			return nil
		}
		return fmt.Errorf("%w: %s", ErrUnmappedEvent, event.Name)
	}

	// Start a producer span around the publish. The span is the parent
	// of the NATS-side trace once a downstream consumer extracts the
	// W3C TraceContext we inject into the message header (gateway SSE
	// multiplexer, GraphQL subscription bus). Span name follows the
	// OTel messaging conventions: "<destination> publish".
	tracer := otel.Tracer("titular/indexer/publisher")
	ctx, span := tracer.Start(ctx, subject+" publish",
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			attribute.String("messaging.system", "nats"),
			attribute.String("messaging.destination.name", subject),
			attribute.String("messaging.operation", "publish"),
			// Event name is non-cardinal (one per registered abi event),
			// so it is safe as a span attribute — the dispatch table
			// caps the set at decoder.KnownEvents().
			attribute.String("titular.event.name", event.Name),
			attribute.Int64("titular.block_number", int64(event.Raw.BlockNumber)),
		),
	)
	defer span.End()

	payload, err := wrapPayload(event.Payload)
	if err != nil {
		span.SetStatus(codes.Error, "wrap payload")
		span.RecordError(err)
		return fmt.Errorf("publisher: wrap %s: %w", event.Name, err)
	}

	body, err := json.Marshal(envelope{
		Name:        event.Name,
		BlockNumber: event.Raw.BlockNumber,
		TxHash:      event.Raw.TxHash.Hex(),
		LogIndex:    event.Raw.Index,
		Payload:     payload,
	})
	if err != nil {
		span.SetStatus(codes.Error, "encode envelope")
		span.RecordError(err)
		return fmt.Errorf("publisher: encode %s: %w", event.Name, err)
	}

	msgID := DedupID(event.Raw.TxHash.Hex(), event.Raw.Index)
	msg := &nats.Msg{Subject: subject, Data: body}
	// Inject W3C TraceContext headers so a downstream consumer of this
	// JetStream subject can resume the trace as a child of THIS span.
	// Safe when OTel is disabled — the propagator is the default no-op
	// in that case and writes nothing.
	observability.InjectNATSHeaders(ctx, msg)

	pubOpts := []nats.PubOpt{nats.MsgId(msgID), nats.Context(ctx)}
	if _, err := p.js.PublishMsg(msg, pubOpts...); err != nil {
		span.SetStatus(codes.Error, "publish")
		span.RecordError(err)
		return fmt.Errorf("publisher: publish %s to %s: %w", event.Name, subject, err)
	}
	return nil
}

// DedupID constructs the JetStream `Nats-Msg-Id` header value for a log.
// Format: `<lowercase tx hash>:<log index>`. The tx hash is lowercased so
// two callers that capitalise hex differently (geth's checksummed form
// vs. lowercase) still hash to the same dedup key. The format is stable:
// any future change requires an explicit migration because deployed
// environments rely on it for the dedup-window guarantee.
//
// Exported so tests and ad-hoc tools can produce the same key without
// re-deriving the format.
func DedupID(txHash string, logIndex uint) string {
	return strings.ToLower(txHash) + ":" + strconv.FormatUint(uint64(logIndex), 10)
}

// bigIntStr wraps a [big.Int] for JSON encoding as a quoted decimal
// string. This is the precision-preserving wire form: JS consumers
// (gateway, SDK) cannot represent `uint256` values above 2^53 as JSON
// numbers without loss, so every [big.Int] field on an abigen payload is
// rewritten to a [bigIntStr] in [wrapPayload] before encoding.
//
// The receiver is a pointer so the JSON encoder picks up MarshalJSON via
// the value's address; both `*big.Int` and `big.Int` payload fields are
// converted to `*bigIntStr` in [wrapPayload], so callers never construct
// this type directly.
type bigIntStr big.Int

// MarshalJSON returns the decimal representation of the underlying
// big.Int wrapped in quotes. A nil receiver is encoded as JSON null,
// matching the default `*big.Int` semantic for pointer-typed payload
// fields that the binding may leave nil.
func (b *bigIntStr) MarshalJSON() ([]byte, error) {
	if b == nil {
		return []byte("null"), nil
	}
	v := (*big.Int)(b)
	// `"` + decimal + `"`. Pre-size the slice so the hot path avoids
	// the small-buffer growth dance in bytes.Buffer.
	dec := v.Text(10)
	out := make([]byte, 0, len(dec)+2)
	out = append(out, '"')
	out = append(out, dec...)
	out = append(out, '"')
	return out, nil
}

// bigIntType / bigIntPtrType cache the reflect.Type values used by
// [wrapPayload]. Reflection lookups happen once per OnEvent call, so the
// allocation cost is dominated by the encode itself; caching the types
// keeps the conversion branch O(1) per field.
var (
	bigIntType    = reflect.TypeOf(big.Int{})
	bigIntPtrType = reflect.TypeOf((*big.Int)(nil))
)

// wrapPayload takes the abigen-decoded event struct and returns a value
// suitable for JSON encoding. Every `*big.Int` and `big.Int` field is
// replaced with a [bigIntStr] wrapper so the wire form is a quoted
// decimal string instead of a number; every other field passes through
// untouched.
//
// Implementation:
//
//   - Accepts `*Struct` (the form abigen returns from Parse*) or
//     `Struct` directly. Anything else is returned as-is — the publisher
//     never injects untyped maps, but a future test may.
//   - Walks exported fields once. The output is a `map[string]any`
//     keyed by the struct field name (which matches abigen's default
//     JSON tag, so the wire format is unchanged for non-bigint fields).
//   - The "Raw" types.Log field that abigen embeds on every event
//     struct is dropped — the envelope already carries block/tx/index
//     and the embedded log would just bloat the payload.
func wrapPayload(payload any) (any, error) {
	if payload == nil {
		return nil, nil
	}
	v := reflect.ValueOf(payload)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		// Not an abigen struct — pass through. Callers in tests may
		// hand us a map directly; we trust the caller to have already
		// stringified any big ints in that case.
		return payload, nil
	}

	t := v.Type()
	out := make(map[string]any, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		if field.Name == "Raw" {
			// Skip the abigen-embedded raw log; the envelope already
			// carries the chain coordinates we care about.
			continue
		}
		fv := v.Field(i)
		out[field.Name] = wrapValue(fv)
	}
	return out, nil
}

// wrapValue rewrites a single reflect.Value, returning the underlying
// concrete value for the JSON encoder. big.Int / *big.Int are converted
// to *bigIntStr; everything else is returned as the field's interface
// value. Slices and structs are walked recursively so nested big.Int
// fields (e.g. a tuple-typed event parameter) are also rewritten.
func wrapValue(v reflect.Value) any {
	switch {
	case v.Type() == bigIntPtrType:
		if v.IsNil() {
			return nil
		}
		// Convert *big.Int -> *bigIntStr without copying the limb slice
		// — the underlying storage is identical and we only read it.
		bi := v.Interface().(*big.Int)
		return (*bigIntStr)(bi)
	case v.Type() == bigIntType:
		// Prefer Addr() when the field is addressable (the common case
		// when we descend from a *Struct). Otherwise copy the value
		// into a fresh big.Int — the encoder only reads the value, so
		// the copy cost is paid once at marshal time.
		if v.CanAddr() {
			bi := v.Addr().Interface().(*big.Int)
			return (*bigIntStr)(bi)
		}
		bi := v.Interface().(big.Int)
		return (*bigIntStr)(&bi)
	}

	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			return nil
		}
		return wrapValue(v.Elem())
	case reflect.Struct:
		t := v.Type()
		out := make(map[string]any, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}
			out[f.Name] = wrapValue(v.Field(i))
		}
		return out
	case reflect.Slice, reflect.Array:
		n := v.Len()
		// Byte slices encode as base64 strings under the default
		// json.Marshal — keep that behaviour by handing them through
		// rather than walking element-by-element.
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return v.Interface()
		}
		out := make([]any, n)
		for i := 0; i < n; i++ {
			out[i] = wrapValue(v.Index(i))
		}
		return out
	}
	return v.Interface()
}

// Compile-time assertion: Publisher satisfies decoder.Handler.
var _ decoder.Handler = (*Publisher)(nil)
