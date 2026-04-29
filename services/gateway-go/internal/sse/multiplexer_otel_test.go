package sse

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/observability"
)

// Test_Multiplexer_PropagatesTraceContextFromNATSHeader is the cross-service
// trace continuity contract: when a publisher injected a W3C TraceContext
// onto the NATS header (kind=PRODUCER span), the multiplexer's dispatch
// callback MUST extract it and open the consume span as a child of the
// producer — same TraceID, parent SpanID equal to the producer's SpanID.
//
// We use a single in-memory recorder for both halves because the wire path
// is in-process here (one nats-server, one nats.Conn); the assertion that
// matters is the propagator round-trip, not the network shape.
func Test_Multiplexer_PropagatesTraceContextFromNATSHeader(t *testing.T) {
	rec := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(rec))
	prevTP := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() {
		_ = tp.Shutdown(context.Background())
		otel.SetTracerProvider(prevTP)
	})

	prevProp := otel.GetTextMapPropagator()
	otel.SetTextMapPropagator(propagation.TraceContext{})
	t.Cleanup(func() { otel.SetTextMapPropagator(prevProp) })

	nc := runEmbeddedNATS(t)
	mux, err := NewMultiplexer(Config{NATS: nc})
	if err != nil {
		t.Fatalf("NewMultiplexer: %v", err)
	}
	if err := mux.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer mux.Stop()

	sub, err := mux.Subscribe(SubscribeOptions{})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	defer sub.Close()

	// Producer side: open a kind=producer span, inject onto a NATS msg
	// header, publish. This mirrors what the indexer's publisher does
	// inside OnEvent (see services/indexer-go/internal/publisher).
	tracer := otel.Tracer("titular/test/producer")
	pubCtx, pubSpan := tracer.Start(context.Background(), "trades.bought publish",
		trace.WithSpanKind(trace.SpanKindProducer))
	pubTraceID := pubSpan.SpanContext().TraceID()
	pubSpanID := pubSpan.SpanContext().SpanID()

	msg := &nats.Msg{Subject: "trades.bought", Data: []byte(`{"x":1}`)}
	observability.InjectNATSHeaders(pubCtx, msg)
	if err := nc.PublishMsg(msg); err != nil {
		t.Fatalf("PublishMsg: %v", err)
	}
	pubSpan.End()

	// Consumer side: wait for the multiplexer to dispatch the message.
	select {
	case ev, ok := <-sub.C:
		if !ok {
			t.Fatal("subscriber channel closed early")
		}
		if ev.Subject != "trades.bought" {
			t.Fatalf("got subject %q, want trades.bought", ev.Subject)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for fanout")
	}

	// The consume span ended inside dispatch; the producer span ended
	// just above. Both should now be in rec.Ended(). Find each by name.
	var producerROS, consumerROS sdktrace.ReadOnlySpan
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		for _, s := range rec.Ended() {
			switch s.Name() {
			case "trades.bought publish":
				producerROS = s
			case "trades.bought consume":
				consumerROS = s
			}
		}
		if producerROS != nil && consumerROS != nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if producerROS == nil {
		t.Fatal("producer span 'trades.bought publish' missing from recorder")
	}
	if consumerROS == nil {
		t.Fatal("consumer span 'trades.bought consume' missing from recorder")
	}

	// TraceID equality is the cross-service stitch — fails immediately
	// if dispatch dropped the header instead of extracting the parent.
	if got := consumerROS.SpanContext().TraceID(); got != pubTraceID {
		t.Errorf("consumer TraceID = %s, want %s (producer)", got, pubTraceID)
	}
	// Parent linkage: the consume span MUST point at the producer span
	// as its parent, otherwise the backend cannot draw the edge.
	if got := consumerROS.Parent().SpanID(); got != pubSpanID {
		t.Errorf("consumer parent SpanID = %s, want %s (producer)", got, pubSpanID)
	}
	// Kind sanity: producer/consumer pair, not internal/server.
	if k := producerROS.SpanKind(); k != trace.SpanKindProducer {
		t.Errorf("producer kind = %s, want producer", k)
	}
	if k := consumerROS.SpanKind(); k != trace.SpanKindConsumer {
		t.Errorf("consumer kind = %s, want consumer", k)
	}
}
