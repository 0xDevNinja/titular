package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// installTestRecorder installs an in-memory span recorder as the global
// TracerProvider and registers a t.Cleanup that restores the previous
// provider. The returned recorder collects every span End()-ed during
// the test, so the test body can inspect span hierarchy / attributes
// without standing up a collector.
//
// Unexported because the indexer's other packages (publisher, decoder)
// have their own test files that drive observability via the global
// provider — they do not need to know about the recorder type.
func installTestRecorder(t interface {
	Helper()
	Cleanup(func())
}) *tracetest.SpanRecorder {
	t.Helper()
	rec := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(rec))
	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() {
		_ = tp.Shutdown(context.Background())
		otel.SetTracerProvider(prev)
	})
	return rec
}
