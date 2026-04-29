package observability

import (
	"context"
	"testing"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Test_NATSHeaderRoundtrip covers the publisher → consumer trace
// propagation contract. The indexer publisher injects W3C TraceContext
// onto every outbound message; the gateway's SSE multiplexer or
// GraphQL bus extracts it back. A break here would silently sever
// every cross-service trace.
func Test_NATSHeaderRoundtrip(t *testing.T) {
	prev := otel.GetTextMapPropagator()
	otel.SetTextMapPropagator(propagation.TraceContext{})
	t.Cleanup(func() { otel.SetTextMapPropagator(prev) })

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

// Test_NATSHeaderRoundtrip_NilSafety checks that the helpers do not
// panic on nil messages.
func Test_NATSHeaderRoundtrip_NilSafety(t *testing.T) {
	InjectNATSHeaders(context.Background(), nil)

	parent := context.WithValue(context.Background(), nilSafetyKey{}, "marker")
	got := ExtractNATSContext(parent, nil)
	if got.Value(nilSafetyKey{}) != "marker" {
		t.Error("nil message did not return parent ctx unchanged")
	}
}

type nilSafetyKey struct{}

// Test_natsHeaderCarrier_Keys ensures the propagator can enumerate
// header names; a Keys() that returned nothing would silently break
// extraction.
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
