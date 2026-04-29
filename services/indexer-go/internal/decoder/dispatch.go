// Package decoder turns raw EVM logs into typed Go event structs and routes
// them to the right handler.
//
// Layering inside the indexer (M4):
//
//	subscriber  -> emits ConfirmedBlock for each finalised height
//	(caller)    -> fetches the block's logs via eth_getLogs
//	decoder     -> THIS PACKAGE: per-log topic0 lookup -> typed struct
//	handler     -> downstream consumer (NATS publisher in #82, then
//	               persistence in later phases)
//
// The decoder is intentionally thin and deterministic:
//
//   - One [Decoder] per Solidity event listed in the indexer's pinned event
//     surface (`internal/abi/abi_test.go` is the source of truth).
//   - All decoders share a single dispatch table keyed by topic0 (the keccak256
//     hash of the canonical event signature). Lookup is therefore O(1) and the
//     caller never has to know which contract emitted the log.
//   - Decoding is offline. Each decoder constructs a per-call abigen filterer
//     bound to the zero address with a nil backend — abigen's Parse* methods
//     only require the embedded ABI, never an RPC, so unit tests can drive
//     synthetic logs without a chain.
//
// Topic0 collision is impossible across the M3 ACP set we currently track:
// every entry in `abi_test.go`'s TestACPEventSurface has a distinct hash. If
// future regenerations introduce a duplicate (e.g. an OZ-style overload), the
// dispatch table will surface it at registration time via Register's panic.
package decoder

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ErrUnknownEvent is returned by [Dispatch] when a log's topic0 has no
// registered decoder. Callers (the NATS publisher in #82, the per-event
// handlers later) typically log-and-skip rather than fail the block: a
// Titular contract may emit non-tracked events (role grants, OZ proxy
// versioning, pause toggles) that are deliberately not in the dispatch map.
var ErrUnknownEvent = errors.New("decoder: unknown event")

// ErrNoTopics is returned when a log carries an empty topics slice. EVM logs
// with at least one indexed parameter — i.e. every event we care about — must
// have topic0 set. An empty Topics slice indicates either a malformed RPC
// response or an anonymous event (which we never emit).
var ErrNoTopics = errors.New("decoder: log has no topics")

// DecodedEvent is the canonical envelope returned by [Dispatch]. It carries
// the contract-qualified event name, the underlying typed struct from the
// abigen bindings, and a back-reference to the raw log so downstream code
// can re-derive block/tx/index metadata without threading them separately.
//
// Name uses the form "<contract>.<event>" (e.g. "Job.ResultSubmitted") so
// duplicate event names across contracts (Job.JobCancelled vs the unused
// JobFactory equivalent, were it to exist) remain unambiguous in logs and
// metrics. The [Handler] interface receives this envelope unchanged.
type DecodedEvent struct {
	// Name is the contract-qualified event name. Stable and safe to use as a
	// metric label or NATS subject component.
	Name string

	// Payload is the abigen-generated event struct (e.g. *abi.JobJobInitialised,
	// *abi.EscrowFunded). Callers type-assert against the expected concrete
	// type; the decoder package itself never inspects it.
	Payload any

	// Raw is the original log. Carrying it makes (BlockNumber, TxHash,
	// LogIndex) accessible to downstream idempotency keys without a second
	// pointer dance.
	Raw types.Log
}

// Decoder converts a single raw log into a typed event struct. The returned
// `any` is one of the abigen-generated event types declared in
// `internal/abi/*.go`. Implementations MUST NOT perform RPC calls — they
// rely solely on the log's topics and data.
type Decoder func(log types.Log) (any, error)

// Handler is the downstream consumer of decoded events. The NATS publisher
// (#82) and the database writer (later) both implement this interface so the
// dispatch layer can fan out without knowing what comes next.
//
// Implementations MUST be safe for synchronous calls from a single goroutine
// (the indexer's main loop). A non-nil error returned from OnEvent terminates
// the per-block processing pipeline; transient downstream failures should be
// retried inside the implementation rather than surfaced.
//
// The ctx is the per-block context: handlers MUST honour ctx.Done() so a
// stalled downstream (e.g. an unreachable NATS server) does not pin the
// indexer's main loop indefinitely. Callers cancel the ctx when shutting
// down or when they want to abandon a block.
type Handler interface {
	OnEvent(ctx context.Context, event DecodedEvent) error
}

// HandlerFunc adapts a plain function to [Handler]. Useful for tests and for
// composing small ad-hoc handlers without a named type.
type HandlerFunc func(ctx context.Context, event DecodedEvent) error

// OnEvent satisfies the Handler interface by calling the underlying func.
func (f HandlerFunc) OnEvent(ctx context.Context, event DecodedEvent) error {
	return f(ctx, event)
}

// registryEntry pairs the decoder for a topic0 with the qualified event name
// it produces, so [Dispatch] does not have to thread the name separately.
type registryEntry struct {
	name    string
	decoder Decoder
}

// registry is the singleton dispatch table populated by init() in decoders.go.
// It is treated as read-only after package initialisation; the mutex only
// guards the brief construction window so tests that import the package in
// parallel cannot observe a torn map (Go's race detector flags this even when
// the writes happen during init from a single goroutine, because tests can
// add custom decoders via [Register] in their own init).
var (
	registryMu sync.RWMutex
	registry   = make(map[common.Hash]registryEntry)
)

// Register installs a decoder for the given topic0. It panics if topic0 is
// already registered, surfacing accidental duplication at startup rather
// than as a silent first-write-wins behaviour.
//
// Register is intended to be called from package init in decoders.go. Tests
// MAY call it (e.g. to inject a fake decoder for an unknown topic) but MUST
// reverse the change with [Unregister] in t.Cleanup, otherwise the test
// would leak global state into sibling tests.
func Register(topic0 common.Hash, name string, decoder Decoder) {
	if decoder == nil {
		panic(fmt.Sprintf("decoder.Register: nil decoder for %s", name))
	}
	if name == "" {
		panic(fmt.Sprintf("decoder.Register: empty name for topic0 %s", topic0.Hex()))
	}

	registryMu.Lock()
	defer registryMu.Unlock()
	if existing, ok := registry[topic0]; ok {
		panic(fmt.Sprintf(
			"decoder.Register: topic0 %s already registered as %q, refusing to overwrite with %q",
			topic0.Hex(), existing.name, name,
		))
	}
	registry[topic0] = registryEntry{name: name, decoder: decoder}
}

// Unregister removes a topic0 from the dispatch table. Tests use this to undo
// [Register] in their cleanup; production code never calls it.
func Unregister(topic0 common.Hash) {
	registryMu.Lock()
	defer registryMu.Unlock()
	delete(registry, topic0)
}

// KnownTopics returns the list of topic0 hashes the dispatch table currently
// recognises. Order is not stable across calls; the result is intended for
// tests and the abigen-check CI step, not as a production filter.
func KnownTopics() []common.Hash {
	registryMu.RLock()
	defer registryMu.RUnlock()
	out := make([]common.Hash, 0, len(registry))
	for h := range registry {
		out = append(out, h)
	}
	return out
}

// KnownEvents returns the list of contract-qualified event names the
// dispatch table currently recognises (e.g. "Job.JobInitialised"). Order is
// not stable. Used by sibling packages (the publisher's subject table) to
// assert that every decoded event has a downstream mapping.
func KnownEvents() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()
	out := make([]string, 0, len(registry))
	for _, e := range registry {
		out = append(out, e.name)
	}
	return out
}

// Dispatch decodes a single log against the registered decoders and returns
// the contract-qualified event name plus the typed payload.
//
// Behaviour:
//
//   - If the log has no topics, returns [ErrNoTopics]. This is a malformed
//     log; the caller should drop the block or surface the RPC error.
//   - If topic0 is not in the dispatch table, returns ([..., ErrUnknownEvent]).
//     The error is wrapped so callers can use errors.Is to choose between
//     log-and-skip and hard-fail.
//   - If the underlying abigen Parse* call fails, returns the wrapped error
//     verbatim. ABI drift between the deployed contract and the regenerated
//     bindings is the most common cause; the abigen-check CI job is meant to
//     catch this before merge.
func Dispatch(log types.Log) (DecodedEvent, error) {
	if len(log.Topics) == 0 {
		return DecodedEvent{}, ErrNoTopics
	}
	topic0 := log.Topics[0]

	registryMu.RLock()
	entry, ok := registry[topic0]
	registryMu.RUnlock()
	if !ok {
		return DecodedEvent{}, fmt.Errorf("%w: topic0=%s", ErrUnknownEvent, topic0.Hex())
	}

	payload, err := entry.decoder(log)
	if err != nil {
		return DecodedEvent{}, fmt.Errorf("decode %s: %w", entry.name, err)
	}

	return DecodedEvent{
		Name:    entry.name,
		Payload: payload,
		Raw:     log,
	}, nil
}

// DispatchAndHandle is the convenience entry point for callers that want to
// route a single log straight to a [Handler]. It decodes, then calls
// handler.OnEvent. Logs whose topic0 is unknown are silently skipped and
// reported as a (false, nil) return — production callers (the per-block
// pipeline) treat unknown events as expected, since not every event in a
// Titular contract is part of the tracked surface.
//
// Returns:
//
//	(true,  nil)        — log decoded and handler succeeded.
//	(false, nil)        — log skipped because topic0 is unknown OR has no topics.
//	(_,     non-nil)    — decode error or handler error; caller decides whether
//	                      to abort the block.
func DispatchAndHandle(ctx context.Context, log types.Log, handler Handler) (bool, error) {
	if handler == nil {
		return false, errors.New("decoder.DispatchAndHandle: handler is nil")
	}
	event, err := Dispatch(log)
	switch {
	case errors.Is(err, ErrNoTopics), errors.Is(err, ErrUnknownEvent):
		return false, nil
	case err != nil:
		return false, err
	}
	if err := handler.OnEvent(ctx, event); err != nil {
		return false, fmt.Errorf("handler %s: %w", event.Name, err)
	}
	return true, nil
}
