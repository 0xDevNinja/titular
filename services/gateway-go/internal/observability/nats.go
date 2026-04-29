// NATS tracing helpers for the gateway.
//
// There is no canonical otel-contrib instrumentation for nats.go yet
// (see https://github.com/open-telemetry/opentelemetry-go-contrib/issues/2362),
// so we provide a hand-rolled propagator adapter that the SSE multiplexer
// and the GraphQL subscription bus use to pull a trace context out of
// every fanout message. The format is W3C TraceContext, matching what
// the indexer publisher injects.
package observability

import (
	"context"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
)

// natsHeaderCarrier adapts nats.Header to propagation.TextMapCarrier.
// nats.Header is a type alias for net/http.Header in the v1 client, so
// the case-insensitive contract on Get/Set holds.
type natsHeaderCarrier struct {
	h nats.Header
}

// Get returns the first value for key.
func (c natsHeaderCarrier) Get(key string) string { return c.h.Get(key) }

// Set replaces any existing values for key with v.
func (c natsHeaderCarrier) Set(key, v string) { c.h.Set(key, v) }

// Keys returns every header key present on the message; allocates a
// fresh slice per call to match the propagator's mutability contract.
func (c natsHeaderCarrier) Keys() []string {
	out := make([]string, 0, len(c.h))
	for k := range c.h {
		out = append(out, k)
	}
	return out
}

// ExtractNATSContext returns a context derived from parent that carries
// any trace context present in msg. Falls back to parent unchanged when
// the message has no header or the header is malformed (W3C
// TraceContext "best effort" semantics).
//
// Safe on a nil message: returns parent.
func ExtractNATSContext(parent context.Context, msg *nats.Msg) context.Context {
	if msg == nil || msg.Header == nil {
		return parent
	}
	return otel.GetTextMapPropagator().Extract(parent, natsHeaderCarrier{h: msg.Header})
}

// InjectNATSHeaders writes the active trace context from ctx into msg's
// header so a downstream consumer can resume the trace. Safe on a nil
// message; allocates the header map on demand.
func InjectNATSHeaders(ctx context.Context, msg *nats.Msg) {
	if msg == nil {
		return
	}
	if msg.Header == nil {
		msg.Header = nats.Header{}
	}
	otel.GetTextMapPropagator().Inject(ctx, natsHeaderCarrier{h: msg.Header})
}
