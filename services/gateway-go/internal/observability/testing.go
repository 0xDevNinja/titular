package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// installTestRecorder installs an in-memory span recorder as the global
// TracerProvider and returns it together with a restore function the
// caller MUST defer. The recorder is the canonical OTel testing seam:
// every span emitted via otel.Tracer(...) lands in rec.Ended() once it
// has been End()-ed, so tests can assert hierarchy and attributes
// without standing up a collector.
//
// This helper is unexported because it is only useful inside the
// observability package's own tests; the rest of the codebase uses
// the global provider directly.
func installTestRecorder(t interface {
	Helper()
	Cleanup(func())
}) (*tracetest.SpanRecorder, func()) {
	t.Helper()
	rec := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(rec))
	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	restore := func() {
		_ = tp.Shutdown(context.Background())
		otel.SetTracerProvider(prev)
	}
	t.Cleanup(restore)
	return rec, restore
}
