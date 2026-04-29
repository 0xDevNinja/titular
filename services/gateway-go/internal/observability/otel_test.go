package observability

import (
	"context"
	"os"
	"strings"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// Test_Init_NoEndpointIsNoop verifies that Init returns a no-op shutdown
// and DOES NOT touch the global TracerProvider when
// OTEL_EXPORTER_OTLP_ENDPOINT is unset. This is the production-safe
// stance for operators who have not yet wired a collector — observability
// must be silently disabled rather than crashing startup.
func Test_Init_NoEndpointIsNoop(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")

	// Snapshot the current provider before Init.
	before := otel.GetTracerProvider()

	shutdown, err := Init(context.Background(), Config{ServiceName: "test-gw"})
	if err != nil {
		t.Fatalf("Init returned error in disabled mode: %v", err)
	}
	if shutdown == nil {
		t.Fatal("Init returned nil shutdown in disabled mode")
	}

	// Provider must be unchanged — Init must NOT install an SDK on
	// the global when disabled.
	if otel.GetTracerProvider() != before {
		t.Fatal("Init mutated the global TracerProvider in disabled mode")
	}

	// Shutdown must be safely callable.
	if err := shutdown(context.Background()); err != nil {
		t.Errorf("noop shutdown returned error: %v", err)
	}
}

// Test_resolveSampleRate covers the parser:
//   - unset env returns the default,
//   - parseable env in [0,1] wins,
//   - out-of-range env values fall through to the default rather than
//     silently clamping (an operator who set "10" should see zero spans
//     and look at the env, not get 100% sampling without warning),
//   - malformed env values fall through too.
func Test_resolveSampleRate(t *testing.T) {
	cases := []struct {
		name   string
		env    string
		def    float64
		expect float64
	}{
		{"unset uses default", "", 0.25, 0.25},
		{"valid env wins", "0.5", 0.25, 0.5},
		{"env zero is honoured", "0", 0.25, 0},
		{"env one is honoured", "1", 0.25, 1},
		{"env out of range falls back", "10", 0.25, 0.25},
		{"env negative falls back", "-1", 0.25, 0.25},
		{"env malformed falls back", "abc", 0.25, 0.25},
		{"all-zero defaults to 0.1", "", 0, 0.1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("OTEL_TRACES_SAMPLER_ARG", tc.env)
			got := resolveSampleRate(tc.def)
			if got != tc.expect {
				t.Errorf("resolveSampleRate(%v) with env %q = %v, want %v",
					tc.def, tc.env, got, tc.expect)
			}
		})
	}
}

// Test_stripScheme covers the URL normaliser. OTLP gRPC takes a
// host:port pair, but operators typically write the URL with a scheme;
// the helper drops http://, https:// and grpc:// so the rest of the
// SDK sees only the authority.
func Test_stripScheme(t *testing.T) {
	cases := map[string]string{
		"localhost:4317":         "localhost:4317",
		"http://localhost:4317":  "localhost:4317",
		"https://collector:4317": "collector:4317",
		"grpc://collector:4317":  "collector:4317",
		"":                       "",
	}
	for in, want := range cases {
		if got := stripScheme(in); got != want {
			t.Errorf("stripScheme(%q) = %q, want %q", in, got, want)
		}
	}
}

// Test_buildResource_HonoursEnvOverride verifies the env-vs-cfg
// precedence: OTEL_SERVICE_NAME and OTEL_SERVICE_VERSION take
// precedence over the cfg fields, matching the upstream SDK contract
// and the documented env table in the package doc.
func Test_buildResource_HonoursEnvOverride(t *testing.T) {
	t.Setenv("OTEL_SERVICE_NAME", "from-env")
	t.Setenv("OTEL_SERVICE_VERSION", "9.9.9")

	res, err := buildResource(context.Background(), Config{
		ServiceName:    "from-cfg",
		ServiceVersion: "0.0.1",
	})
	if err != nil {
		t.Fatalf("buildResource: %v", err)
	}

	got := map[string]string{}
	for _, kv := range res.Attributes() {
		got[string(kv.Key)] = kv.Value.AsString()
	}
	if got["service.name"] != "from-env" {
		t.Errorf("service.name = %q, want from-env", got["service.name"])
	}
	if got["service.version"] != "9.9.9" {
		t.Errorf("service.version = %q, want 9.9.9", got["service.version"])
	}
}

// Test_NoCredentialsLeak_InResourceAttributes checks that the resource
// we build never includes the obvious credential-looking fields. This
// is a guard against a future contributor accidentally adding the
// database DSN or the JWT secret to the resource (an easy mistake when
// reusing the cfg.Service value for service.name).
//
// We do NOT enumerate every possible credential string — that would be
// a moving target — but we DO assert that no attribute key or value
// contains the word "secret", "password", or "Authorization", which
// covers the standard families. Trace exporters and the resource go
// onto every emitted span, so a leak here would surface in 100% of
// traffic, making this a high-leverage check.
func Test_NoCredentialsLeak_InResourceAttributes(t *testing.T) {
	res, err := buildResource(context.Background(), Config{
		ServiceName:    "gateway",
		ServiceVersion: "test",
	})
	if err != nil {
		t.Fatalf("buildResource: %v", err)
	}

	forbidden := []string{"secret", "password", "authorization", "bearer ", "0x"}
	for _, kv := range res.Attributes() {
		key := strings.ToLower(string(kv.Key))
		val := strings.ToLower(kv.Value.Emit())
		for _, f := range forbidden {
			if strings.Contains(key, f) {
				t.Errorf("resource attribute key %q contains forbidden substring %q", kv.Key, f)
			}
			if strings.Contains(val, f) {
				t.Errorf("resource attribute %q value contains forbidden substring %q", kv.Key, f)
			}
		}
	}
}

// Test_GlobalTracerProviderDefault_IsNoop documents the package
// invariant that consumers can safely call otel.Tracer(...) before
// Init runs (or in disabled mode) and the resulting tracer is the
// upstream no-op. This is what lets every call site in publisher /
// decoder / SSE skip a "is observability enabled?" guard.
func Test_GlobalTracerProviderDefault_IsNoop(t *testing.T) {
	// Reset env so we don't accidentally pick up a nearby collector.
	if err := os.Unsetenv("OTEL_EXPORTER_OTLP_ENDPOINT"); err != nil {
		t.Fatalf("unsetenv: %v", err)
	}
	tp := otel.GetTracerProvider()
	tracer := tp.Tracer("test")
	_, span := tracer.Start(context.Background(), "noop-span")
	defer span.End()

	// The default provider's tracer must satisfy the no-op trace
	// package contract: spans are valid but not recording.
	if span.IsRecording() {
		t.Errorf("default global tracer is recording; expected no-op")
	}
	// We deliberately do not assert the concrete type because the
	// upstream SDK has migrated the no-op tracer between major versions;
	// the IsRecording check is the documented stable contract.
	_ = noop.NewTracerProvider // import-keep
}

// Test_TestExporter_CapturesSpan boots an in-memory span recorder via
// tracetest and verifies that a manually-started span shows up in the
// expected hierarchy. This is the canonical pattern for the rest of
// the test suite; the helper is exported to package consumers as
// NewTestRecorder for reuse.
func Test_TestExporter_CapturesSpan(t *testing.T) {
	rec, restore := installTestRecorder(t)
	defer restore()

	tracer := otel.Tracer("test")
	ctx, parent := tracer.Start(context.Background(), "parent")
	_, child := tracer.Start(ctx, "child")
	child.End()
	parent.End()

	got := rec.Ended()
	if len(got) != 2 {
		t.Fatalf("expected 2 spans, got %d", len(got))
	}
	// Spans are recorded in End order — child first, then parent.
	if got[0].Name() != "child" || got[1].Name() != "parent" {
		t.Errorf("span order: %s, %s; want child, parent", got[0].Name(), got[1].Name())
	}
	// Hierarchy invariant: child's parent span id must equal parent's span id.
	if got[0].Parent().SpanID() != got[1].SpanContext().SpanID() {
		t.Errorf("child.Parent.SpanID = %s; want %s",
			got[0].Parent().SpanID(), got[1].SpanContext().SpanID())
	}
	// Trace id must be the same across the hierarchy.
	if got[0].SpanContext().TraceID() != got[1].SpanContext().TraceID() {
		t.Errorf("trace ids differ across parent/child")
	}
	// Compile-time use of trace package import to keep the dep stable.
	var _ trace.Span = parent
}
