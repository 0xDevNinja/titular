package observability

import (
	"context"
	"testing"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Test_NATSHeaderRoundtrip verifies that injecting a trace context onto
// a nats.Msg header and extracting it back yields the same trace id and
// span id. This is the bridge that lets the indexer's publisher span
// be the parent of the gateway's consumer span: without W3C
// TraceContext on the wire, every cross-service trace breaks at the
// NATS hop.
//
// We install the W3C TraceContext propagator explicitly because the
// global default is a no-op until observability.Init runs — tests must
// not rely on package-init side effects from a sibling test.
func Test_NATSHeaderRoundtrip(t *testing.T) {
	prev := otel.GetTextMapPropagator()
	otel.SetTextMapPropagator(propagation.TraceContext{})
	t.Cleanup(func() { otel.SetTextMapPropagator(prev) })

	// Build a context carrying an active span by hand. We use a
	// deterministic SpanContext so the round-trip equality check
	// has known values. The span is "remote" so the propagator
	// treats it as already-active and writes a traceparent header.
	traceID, _ := trace.TraceIDFromHex("0102030405060708090a0b0c0d0e0f10")
	spanID, _ := trace.SpanIDFromHex("1112131415161718")
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
		Remote:     true,
	})
	ctx := trace.ContextWithSpanContext(context.Background(), sc)

	msg := &nats.Msg{Subject: "test.subject"}
	InjectNATSHeaders(ctx, msg)

	// The propagator must have written a `traceparent` header onto
	// the message so a downstream consumer can read it.
	if msg.Header.Get("traceparent") == "" {
		t.Fatal("no traceparent header on injected message")
	}

	got := ExtractNATSContext(context.Background(), msg)
	gotSC := trace.SpanContextFromContext(got)
	if !gotSC.IsValid() {
		t.Fatal("extracted span context is invalid")
	}
	if gotSC.TraceID() != traceID {
		t.Errorf("trace id roundtrip: got %s, want %s", gotSC.TraceID(), traceID)
	}
	if gotSC.SpanID() != spanID {
		t.Errorf("span id roundtrip: got %s, want %s", gotSC.SpanID(), spanID)
	}
	if !gotSC.IsSampled() {
		t.Errorf("sampled flag dropped on roundtrip")
	}
}

// Test_NATSHeaderRoundtrip_NilSafety verifies that the helpers do not
// panic on a nil message and return the parent context unchanged on
// extract. The indexer publisher's hot path constructs the message
// just before InjectNATSHeaders, so the nil branch is unreachable in
// production — but a future caller that forgets is still protected.
func Test_NATSHeaderRoundtrip_NilSafety(t *testing.T) {
	// Inject on nil must be a no-op (no panic).
	InjectNATSHeaders(context.Background(), nil)

	// Extract on nil must return the parent unchanged.
	parent := context.WithValue(context.Background(), nilSafetyKey{}, "marker")
	got := ExtractNATSContext(parent, nil)
	if got.Value(nilSafetyKey{}) != "marker" {
		t.Error("nil message did not return parent ctx unchanged")
	}
}

type nilSafetyKey struct{}

// Test_natsHeaderCarrier_Keys returns header names. The propagator
// queries Keys() to enumerate carriers when constructing extractors;
// returning nothing means extraction silently fails. This guard catches
// future regressions where the carrier shadowed the underlying header
// map.
func Test_natsHeaderCarrier_Keys(t *testing.T) {
	h := nats.Header{}
	h.Set("traceparent", "00-1234-5678-01")
	h.Set("baggage", "x=y")
	c := natsHeaderCarrier{h: h}
	keys := c.Keys()
	if len(keys) != 2 {
		t.Errorf("Keys = %v, want length 2", keys)
	}
}
