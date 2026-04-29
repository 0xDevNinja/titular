package observability

import (
	"context"
	"strings"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Test_Init_NoEndpointIsNoop verifies disabled-mode contract: when
// OTEL_EXPORTER_OTLP_ENDPOINT is unset, Init returns a no-op shutdown
// and DOES NOT install an SDK on the global TracerProvider. The
// indexer's main loop relies on this — a misconfigured env must not
// crash startup.
func Test_Init_NoEndpointIsNoop(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")

	before := otel.GetTracerProvider()
	shutdown, err := Init(context.Background(), Config{ServiceName: "test-idx"})
	if err != nil {
		t.Fatalf("Init returned error in disabled mode: %v", err)
	}
	if shutdown == nil {
		t.Fatal("Init returned nil shutdown")
	}
	if otel.GetTracerProvider() != before {
		t.Fatal("Init mutated the global TracerProvider in disabled mode")
	}
	if err := shutdown(context.Background()); err != nil {
		t.Errorf("noop shutdown returned: %v", err)
	}
}

// Test_resolveSampleRate covers the same parsing matrix as the gateway
// twin. Duplicated rather than shared because each service is its own
// Go module and we deliberately keep the observability packages free
// of cross-module imports.
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
				t.Errorf("got %v, want %v", got, tc.expect)
			}
		})
	}
}

// Test_stripScheme covers OTLP gRPC URL normalisation.
func Test_stripScheme(t *testing.T) {
	cases := map[string]string{
		"localhost:4317":         "localhost:4317",
		"http://localhost:4317":  "localhost:4317",
		"https://collector:4317": "collector:4317",
		"grpc://collector:4317":  "collector:4317",
	}
	for in, want := range cases {
		if got := stripScheme(in); got != want {
			t.Errorf("stripScheme(%q) = %q, want %q", in, got, want)
		}
	}
}

// Test_buildResource_HonoursEnvOverride asserts env precedence over
// cfg fields, matching the upstream SDK and the documented env table.
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

// Test_NoCredentialsLeak_InResourceAttributes guards against a future
// contributor accidentally putting the chain DSN, RPC URL, or
// indexer-side secret onto the resource. Runs on every emitted span
// in production, so a leak here would surface 100% of the time.
func Test_NoCredentialsLeak_InResourceAttributes(t *testing.T) {
	res, err := buildResource(context.Background(), Config{
		ServiceName:    "indexer",
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
				t.Errorf("resource attribute key %q contains %q", kv.Key, f)
			}
			if strings.Contains(val, f) {
				t.Errorf("resource attribute %q value contains %q", kv.Key, f)
			}
		}
	}
}

// Test_TestRecorder_CapturesHierarchy verifies the in-memory recorder
// helper that other tests in this module use.
func Test_TestRecorder_CapturesHierarchy(t *testing.T) {
	rec := installTestRecorder(t)

	tracer := otel.Tracer("test")
	ctx, parent := tracer.Start(context.Background(), "parent")
	_, child := tracer.Start(ctx, "child")
	child.End()
	parent.End()

	got := rec.Ended()
	if len(got) != 2 {
		t.Fatalf("expected 2 spans, got %d", len(got))
	}
	if got[0].Parent().SpanID() != got[1].SpanContext().SpanID() {
		t.Errorf("hierarchy broken")
	}
	var _ trace.Span = parent
}
