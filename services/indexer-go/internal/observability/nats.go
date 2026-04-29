// NATS tracing helpers for the indexer publisher.
//
// There is no canonical otel-contrib instrumentation for nats.go (the
// repository tracks an open issue tagged "instrumentation"; see
// https://github.com/open-telemetry/opentelemetry-go-contrib/issues/2362).
// Hand-rolling is straightforward because nats.Msg already carries a
// header map in the v1.51 client API and the upstream W3C TraceContext
// propagator works directly on a propagation.TextMapCarrier interface.
//
// We expose two helpers:
//
//   - InjectNATSHeaders(ctx, *nats.Msg)   for the publish side. Writes the
//     active trace context into the message header so a downstream
//     consumer (gateway SSE multiplexer, GraphQL bus) can pick the trace
//     up unchanged.
//
//   - ExtractNATSContext(parent, *nats.Msg) for the consume side. Reads
//     the W3C trace context (if any) and returns a context that becomes
//     the parent of any span the consumer starts. Falls back to the
//     parent ctx unchanged if no header is present, which is exactly the
//     behaviour the propagator's Extract guarantees.
//
// Both functions tolerate nil messages so they are safe to call on the
// publisher's hot path before the JSON envelope is built.
package observability

import (
	"context"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
)

// natsHeaderCarrier adapts a nats.Header to the propagation.TextMapCarrier
// interface the OTel propagator API expects. The adapter is a thin shim
// — every call delegates to the underlying header map — so allocation
// pressure on the hot publish path is one struct per message.
type natsHeaderCarrier struct {
	h nats.Header
}

// Get returns the first value for key, mirroring net/http.Header.Get.
// nats.Header is a type alias for net/http.Header in the v1 client, so
// the contract is identical: case-insensitive lookup, empty string when
// absent.
func (c natsHeaderCarrier) Get(key string) string { return c.h.Get(key) }

// Set replaces any existing values for key with v.
func (c natsHeaderCarrier) Set(key, v string) { c.h.Set(key, v) }

// Keys returns every header key present on the message. The slice is a
// fresh allocation each call to honour the propagator's contract that
// callers may mutate the result.
func (c natsHeaderCarrier) Keys() []string {
	out := make([]string, 0, len(c.h))
	for k := range c.h {
		out = append(out, k)
	}
	return out
}

// InjectNATSHeaders writes the active trace context from ctx into msg's
// header. Safe on a nil message (no-op) and on a context without a span
// (the propagator simply writes nothing). Allocates the header map on
// demand so callers do not need to pre-initialise nats.Msg.Header.
//
// Callers MUST call this before *handing the message off to* nats.Publish
// (rather than after) — once the wire bytes are produced the headers are
// frozen.
func InjectNATSHeaders(ctx context.Context, msg *nats.Msg) {
	if msg == nil {
		return
	}
	if msg.Header == nil {
		msg.Header = nats.Header{}
	}
	otel.GetTextMapPropagator().Inject(ctx, natsHeaderCarrier{h: msg.Header})
}

// ExtractNATSContext returns a context derived from parent that carries
// any trace context present in msg. When the message has no W3C
// traceparent header, the propagator returns the parent ctx unchanged;
// when the header is malformed it is silently dropped and parent is
// returned, which matches the W3C TraceContext spec's "best effort"
// extraction posture.
//
// Safe on a nil message: returns parent.
func ExtractNATSContext(parent context.Context, msg *nats.Msg) context.Context {
	if msg == nil || msg.Header == nil {
		return parent
	}
	return otel.GetTextMapPropagator().Extract(parent, natsHeaderCarrier{h: msg.Header})
}
