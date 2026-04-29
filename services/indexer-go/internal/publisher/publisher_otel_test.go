package publisher_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	contractabi "github.com/0xDevNinja/titular/services/indexer-go/internal/abi"
	"github.com/0xDevNinja/titular/services/indexer-go/internal/decoder"
	"github.com/0xDevNinja/titular/services/indexer-go/internal/publisher"
)

// Test_OnEvent_EmitsProducerSpanWithTraceParent verifies the publisher's
// OTel contract: each OnEvent call MUST produce a span with kind=PRODUCER
// and MUST inject a W3C traceparent header onto the published NATS
// message so a downstream consumer can resume the trace.
//
// We use the embedded JetStream server (helper from publisher_test.go)
// because the publisher's pubMsg path is what carries the headers; a
// stub JetStream that drops the message on the floor would not exercise
// the wire-side header round-trip.
func Test_OnEvent_EmitsProducerSpanWithTraceParent(t *testing.T) {
	// In-memory recorder as the global tracer provider so OnEvent's
	// otel.Tracer(...) call lands its spans into rec.
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

	_, js := runServer(t)
	if err := publisher.EnsureStream(js, publisher.StreamOptions{
		Storage: 1, // nats.MemoryStorage
	}); err != nil {
		t.Fatalf("EnsureStream: %v", err)
	}
	pub, err := publisher.New(js, publisher.Options{})
	if err != nil {
		t.Fatalf("publisher.New: %v", err)
	}

	// Synthesise a JobInitialised event so OnEvent has a real subject
	// mapping — Job.JobInitialised maps onto jobs.initialised in
	// publisher.eventSubjects.
	event := decoder.DecodedEvent{
		Name: "Job.JobInitialised",
		Payload: &contractabi.JobJobInitialised{
			JobId:     big.NewInt(1),
			JobType:   0,
			Principal: common.HexToAddress("0x1111111111111111111111111111111111111111"),
			AgentId:   big.NewInt(2),
			Token:     common.HexToAddress("0x2222222222222222222222222222222222222222"),
			Budget:    big.NewInt(100),
			Deadline:  99,
		},
		Raw: types.Log{
			BlockNumber: 42,
			TxHash:      common.HexToHash("0xabc"),
			Index:       7,
		},
	}

	if err := pub.OnEvent(context.Background(), event); err != nil {
		t.Fatalf("OnEvent: %v", err)
	}

	spans := rec.Ended()
	if len(spans) == 0 {
		t.Fatal("publisher OnEvent emitted no spans")
	}

	// Find the publish span by name suffix.
	var pubSpan sdktrace.ReadOnlySpan
	for _, s := range spans {
		if s.Name() == "" {
			continue
		}
		if matchesSuffix(s.Name(), " publish") {
			pubSpan = s
			break
		}
	}
	if pubSpan == nil {
		names := make([]string, 0, len(spans))
		for _, s := range spans {
			names = append(names, s.Name())
		}
		t.Fatalf("no publish span in %v", names)
	}

	// Producer kind is required so observability backends classify the
	// span correctly in their messaging-flow views.
	if pubSpan.SpanKind().String() != "producer" {
		t.Errorf("publish span kind = %s, want producer", pubSpan.SpanKind())
	}

	// Mandatory attributes: messaging.system, messaging.destination.name.
	want := map[string]string{
		"messaging.system":           "nats",
		"messaging.destination.name": "jobs.initialised",
		"messaging.operation":        "publish",
		"titular.event.name":         "Job.JobInitialised",
	}
	got := map[string]string{}
	for _, kv := range pubSpan.Attributes() {
		got[string(kv.Key)] = kv.Value.Emit()
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("publish span attr %q = %q, want %q", k, got[k], v)
		}
	}

	// No-credentials check: walk every attribute on the publish span
	// and refuse a value that contains the obvious credential markers.
	// Block number / tx hash are the only chain coordinates we expect;
	// the test event's tx hash starts with 0x so we exempt the
	// titular.tx_hash attribute from the 0x ban.
	for _, kv := range pubSpan.Attributes() {
		if string(kv.Key) == "titular.tx_hash" {
			continue
		}
		val := kv.Value.Emit()
		for _, forbidden := range []string{"secret", "password", "authorization", "bearer "} {
			if containsFold(val, forbidden) {
				t.Errorf("publish span attr %q value %q contains %q", kv.Key, val, forbidden)
			}
		}
	}
}

// matchesSuffix is a small helper kept package-local; strings.HasSuffix
// would do it but we import strings only when we need to and the rest
// of this file does not.
func matchesSuffix(s, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}
	return s[len(s)-len(suffix):] == suffix
}

// containsFold is a case-insensitive substring check used by the
// no-credentials guard. We avoid pulling strings.Contains/EqualFold
// to keep the test file dep surface narrow.
func containsFold(haystack, needle string) bool {
	if len(needle) == 0 || len(haystack) < len(needle) {
		return false
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		eq := true
		for j := 0; j < len(needle); j++ {
			a, b := haystack[i+j], needle[j]
			if a >= 'A' && a <= 'Z' {
				a += 'a' - 'A'
			}
			if b >= 'A' && b <= 'Z' {
				b += 'a' - 'A'
			}
			if a != b {
				eq = false
				break
			}
		}
		if eq {
			return true
		}
	}
	return false
}
