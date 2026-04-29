package decoder_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"

	contractabi "github.com/0xDevNinja/titular/services/indexer-go/internal/abi"
	"github.com/0xDevNinja/titular/services/indexer-go/internal/decoder"
)

// Test_DispatchAndHandle_EmitsDispatchSpan verifies the per-event-decoder
// dispatch span: every successful DispatchAndHandle call MUST produce a
// span named "decoder.dispatch <event>" with the qualified event name
// and chain coordinates as attributes. This is what lets an operator
// pivot from a publish failure to the originating block in a tracing
// backend.
func Test_DispatchAndHandle_EmitsDispatchSpan(t *testing.T) {
	rec := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(rec))
	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() {
		_ = tp.Shutdown(context.Background())
		otel.SetTracerProvider(prev)
	})

	parsed := mustAbi(t, contractabi.HookRegistryMetaData)
	tID := topicID(t, contractabi.HookRegistryMetaData, "HookRegistered")
	log := types.Log{
		Topics:      []common.Hash{tID, addrTopic(common.HexToAddress("0x9999"))},
		Data:        encodeData(t, parsed, "HookRegistered", "test-hook"),
		BlockNumber: 42,
		TxHash:      common.HexToHash("0xabc"),
		Index:       7,
	}

	var handlerSpan trace.SpanContext
	h := decoder.HandlerFunc(func(ctx context.Context, _ decoder.DecodedEvent) error {
		// Capture the active span context so we can assert the handler
		// runs WITHIN the dispatch span.
		handlerSpan = trace.SpanContextFromContext(ctx)
		return nil
	})

	if _, err := decoder.DispatchAndHandle(context.Background(), log, h); err != nil {
		t.Fatalf("DispatchAndHandle: %v", err)
	}

	spans := rec.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	span := spans[0]

	// Span name must include the qualified event name so trace search by
	// event resolves cleanly.
	wantName := "decoder.dispatch HookRegistry.HookRegistered"
	if span.Name() != wantName {
		t.Errorf("span name = %q, want %q", span.Name(), wantName)
	}

	// The handler must have run inside the dispatch span — its captured
	// SpanContext should match the dispatch span's id, otherwise the
	// downstream NATS publish span would not nest under it.
	if !handlerSpan.IsValid() {
		t.Fatal("handler did not see an active span")
	}
	if handlerSpan.SpanID() != span.SpanContext().SpanID() {
		t.Errorf("handler span id = %s, want %s",
			handlerSpan.SpanID(), span.SpanContext().SpanID())
	}

	// Attributes must carry the chain coordinates without leaking other
	// fields. We do NOT assert event payload values land on the span —
	// the dispatch span deliberately stops at metadata so a future
	// payload field with PII (signed message body, etc.) cannot leak.
	got := map[string]string{}
	for _, kv := range span.Attributes() {
		got[string(kv.Key)] = kv.Value.Emit()
	}
	if got["titular.event.name"] != "HookRegistry.HookRegistered" {
		t.Errorf("titular.event.name = %q", got["titular.event.name"])
	}
	if got["titular.block_number"] != "42" {
		t.Errorf("titular.block_number = %q, want 42", got["titular.block_number"])
	}
	if got["titular.log_index"] != "7" {
		t.Errorf("titular.log_index = %q, want 7", got["titular.log_index"])
	}
}
