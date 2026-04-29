// Package observability sets up OpenTelemetry tracing and metrics for the
// gateway service.
//
// The package is deliberately small: a single Init entry point that opens
// an OTLP gRPC exporter (when configured), installs the SDK as the global
// TracerProvider / MeterProvider, and returns a shutdown closure the caller
// must invoke on graceful shutdown to flush in-flight spans / metrics.
//
// # Disabled mode
//
// When OTEL_EXPORTER_OTLP_ENDPOINT is unset (or empty), Init is a no-op:
// it returns a nil-shutdown placeholder and leaves the global providers at
// their default no-op implementations. This is the production-safe stance
// because:
//
//   - Operators who do not run a collector must NOT lose startup just because
//     they forgot to configure tracing — observability is opt-in.
//   - The global TracerProvider returned by otel.Tracer in no-op mode emits
//     spans that drop on the floor with no allocation overhead, so call
//     sites can use otel.Tracer(...) unconditionally without a feature flag.
//
// # Configuration env
//
//   - OTEL_EXPORTER_OTLP_ENDPOINT       gRPC endpoint of the collector (e.g.
//                                       "localhost:4317"). When unset, the SDK
//                                       is not installed and the package is
//                                       a no-op (see above).
//   - OTEL_EXPORTER_OTLP_INSECURE       "true" to disable TLS on the exporter
//                                       transport. Defaults to false (TLS).
//                                       Tests and local docker-compose runs
//                                       typically set this to "true".
//   - OTEL_SERVICE_NAME                 logical service name on the resource.
//                                       Defaults to the value of cfg.ServiceName,
//                                       which the binary supplies.
//   - OTEL_SERVICE_VERSION              semver of the running binary; embedded
//                                       on the resource so a collector can
//                                       pivot dashboards on deploy.
//   - OTEL_TRACES_SAMPLER_ARG           parent-based head sampler ratio in
//                                       [0,1]. Defaults to 0.1 (10%) which is
//                                       the production posture; local dev
//                                       sets this to 1.0 to keep every span.
//   - OTEL_RESOURCE_ATTRIBUTES          honoured natively by the SDK's
//                                       resource.Default() merge, no extra
//                                       wiring needed here.
//
// The OTLP exporter package also reads the standard collection of OTEL_*
// vars (timeout, headers, compression) directly via the SDK; this package
// only owns the policy knobs that don't have an upstream binding.
package observability

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

// Config controls what Init wires up.
//
// ServiceName is the logical service name that lands on the resource as
// service.name; the OTEL_SERVICE_NAME env var, when set, takes precedence.
// ServiceVersion is mirrored as service.version; OTEL_SERVICE_VERSION
// overrides.
//
// DefaultSampleRate is the head-sampler ratio used when
// OTEL_TRACES_SAMPLER_ARG is unset. The constant 0.1 (10%) is the
// production posture; tests and local dev that want every span set
// DefaultSampleRate to 1.0 (or set OTEL_TRACES_SAMPLER_ARG=1.0) explicitly.
type Config struct {
	ServiceName       string
	ServiceVersion    string
	DefaultSampleRate float64
}

// ShutdownFunc tears down the SDK, flushing in-flight spans / metrics.
// Always non-nil (a no-op when the SDK was not installed) so callers can
// `defer shutdown(ctx)` without a nil-check.
type ShutdownFunc func(context.Context) error

// noopShutdown is the placeholder returned in disabled mode.
func noopShutdown(context.Context) error { return nil }

// Init wires up tracing and metrics.
//
// When OTEL_EXPORTER_OTLP_ENDPOINT is unset, Init returns a no-op
// shutdown and does NOT touch the global TracerProvider / MeterProvider —
// downstream call sites that grab otel.Tracer("…") will receive the
// default no-op tracer, which is exactly what we want when observability
// is disabled.
//
// When the endpoint is set, Init:
//
//  1. Builds a resource carrying service.name, service.version, and any
//     OTEL_RESOURCE_ATTRIBUTES the operator set.
//  2. Opens an OTLP gRPC trace exporter and installs a batch span
//     processor in front of it.
//  3. Opens an OTLP gRPC metric exporter, mounted via a periodic reader.
//  4. Configures a parent-based head sampler with the configured ratio.
//  5. Sets the W3C TraceContext + Baggage propagators as the global
//     defaults so cross-service propagation just works.
//
// The returned shutdown closure flushes both providers in series and
// MUST be invoked on graceful shutdown — failure to do so drops every
// span / metric still in the batch buffer.
//
// Init never returns a partial install: if any step fails the function
// rolls back what it built and returns the original error. That keeps
// the caller free of "did the metric exporter actually start?" guard
// logic.
func Init(ctx context.Context, cfg Config) (ShutdownFunc, error) {
	endpoint := strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if endpoint == "" {
		// Observability disabled. Leave global providers at their no-op
		// defaults so otel.Tracer(...) call sites still compile and run.
		return noopShutdown, nil
	}

	res, err := buildResource(ctx, cfg)
	if err != nil {
		return noopShutdown, fmt.Errorf("observability: build resource: %w", err)
	}

	insecure := strings.EqualFold(strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_INSECURE")), "true")

	// --- Tracing -----------------------------------------------------------
	traceOpts := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(stripScheme(endpoint))}
	if insecure {
		traceOpts = append(traceOpts, otlptracegrpc.WithInsecure())
	}
	traceExporter, err := otlptracegrpc.New(ctx, traceOpts...)
	if err != nil {
		return noopShutdown, fmt.Errorf("observability: trace exporter: %w", err)
	}

	sampleRate := resolveSampleRate(cfg.DefaultSampleRate)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(
			sdktrace.TraceIDRatioBased(sampleRate),
		)),
	)

	// --- Metrics -----------------------------------------------------------
	metricOpts := []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(stripScheme(endpoint))}
	if insecure {
		metricOpts = append(metricOpts, otlpmetricgrpc.WithInsecure())
	}
	metricExporter, err := otlpmetricgrpc.New(ctx, metricOpts...)
	if err != nil {
		// Roll back the trace provider we just built so the global state
		// is unchanged on a partial failure.
		_ = tp.Shutdown(ctx)
		return noopShutdown, fmt.Errorf("observability: metric exporter: %w", err)
	}
	// Build the MeterProvider with the OTLP periodic reader, plus —
	// when the operator has already wired the Prometheus surface via
	// AttachPrometheusReader (M4 #94) — the Prometheus pull reader.
	// Both readers feed off the same instruments, so a single
	// otel.Meter("…").Int64Counter("foo") shows up in both pipelines.
	mpOpts := []sdkmetric.Option{
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
	}
	if promReader, ok := takeAttachedPrometheusReader(); ok {
		mpOpts = append(mpOpts, sdkmetric.WithReader(promReader))
	}
	mp := sdkmetric.NewMeterProvider(mpOpts...)

	// --- Globals -----------------------------------------------------------
	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	shutdown := func(ctx context.Context) error {
		// Cap shutdown so a wedged collector cannot pin SIGTERM
		// indefinitely. 10s mirrors the upstream OTLP exporter's default
		// timeout and matches the gateway's own srv.Shutdown timeout
		// (15s) with a margin.
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		var errs []error
		if err := tp.Shutdown(shutdownCtx); err != nil {
			errs = append(errs, fmt.Errorf("trace provider shutdown: %w", err))
		}
		if err := mp.Shutdown(shutdownCtx); err != nil {
			errs = append(errs, fmt.Errorf("meter provider shutdown: %w", err))
		}
		return errors.Join(errs...)
	}
	return shutdown, nil
}

// buildResource assembles the SDK resource with service identity. It
// merges (in order) the SDK default resource (telemetry.sdk.* and any
// OTEL_RESOURCE_ATTRIBUTES the operator set), then the explicit
// service.name / service.version we own. Later merges win, so an
// operator-supplied service.name in OTEL_RESOURCE_ATTRIBUTES will be
// overridden by OTEL_SERVICE_NAME — matching the upstream SDK precedence.
func buildResource(ctx context.Context, cfg Config) (*resource.Resource, error) {
	name := strings.TrimSpace(os.Getenv("OTEL_SERVICE_NAME"))
	if name == "" {
		name = cfg.ServiceName
	}
	if name == "" {
		name = "gateway"
	}
	version := strings.TrimSpace(os.Getenv("OTEL_SERVICE_VERSION"))
	if version == "" {
		version = cfg.ServiceVersion
	}

	attrs := []resource.Option{
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithProcessRuntimeName(),
		resource.WithProcessRuntimeVersion(),
		resource.WithAttributes(
			semconv.ServiceName(name),
		),
	}
	if version != "" {
		attrs = append(attrs, resource.WithAttributes(semconv.ServiceVersion(version)))
	}
	return resource.New(ctx, attrs...)
}

// resolveSampleRate parses OTEL_TRACES_SAMPLER_ARG as a float in [0,1] and
// falls back to def when unset / malformed. We refuse out-of-range values
// rather than silently clamping so an operator who set "10" instead of
// "0.1" sees the misconfiguration in the startup log.
//
// The default ratio (when both env and cfg.DefaultSampleRate are zero) is
// 0.1 — the production posture documented on Config.DefaultSampleRate.
//
// On a malformed / out-of-range value we emit a stderr warning so the
// operator sees the typo at boot rather than wondering, an hour later
// at the dashboard, why their "10" sample arg yielded the default 10%.
func resolveSampleRate(def float64) float64 {
	defaultRate := 0.1
	if def > 0 && def <= 1 {
		defaultRate = def
	}
	raw := strings.TrimSpace(os.Getenv("OTEL_TRACES_SAMPLER_ARG"))
	if raw != "" {
		v, err := strconv.ParseFloat(raw, 64)
		if err == nil && v >= 0 && v <= 1 {
			return v
		}
		log.Printf("warn: invalid OTEL_TRACES_SAMPLER_ARG=%q, using default %.2f", raw, defaultRate)
	}
	return defaultRate
}

// stripScheme drops a leading http:// or https:// prefix from the endpoint
// because OTLP gRPC takes a host:port pair and would otherwise treat the
// scheme as part of the authority. We accept either form because operator
// docs (and most collector READMEs) write the URL in the http:// form.
func stripScheme(endpoint string) string {
	for _, prefix := range []string{"http://", "https://", "grpc://"} {
		if strings.HasPrefix(endpoint, prefix) {
			return strings.TrimPrefix(endpoint, prefix)
		}
	}
	return endpoint
}
