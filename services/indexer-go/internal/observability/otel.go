// Package observability sets up OpenTelemetry tracing and metrics for the
// indexer service.
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
package observability

import (
	"context"
	"errors"
	"fmt"
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

// Config controls what Init wires up. See package doc for env var semantics.
type Config struct {
	ServiceName       string
	ServiceVersion    string
	DefaultSampleRate float64
}

// ShutdownFunc tears down the SDK, flushing in-flight spans / metrics.
// Always non-nil so callers can `defer shutdown(ctx)` without a nil-check.
type ShutdownFunc func(context.Context) error

// noopShutdown is the placeholder returned in disabled mode.
func noopShutdown(context.Context) error { return nil }

// Init wires up tracing and metrics. See package doc for behaviour.
func Init(ctx context.Context, cfg Config) (ShutdownFunc, error) {
	endpoint := strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if endpoint == "" {
		return noopShutdown, nil
	}

	res, err := buildResource(ctx, cfg)
	if err != nil {
		return noopShutdown, fmt.Errorf("observability: build resource: %w", err)
	}

	insecure := strings.EqualFold(strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_INSECURE")), "true")

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

	metricOpts := []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(stripScheme(endpoint))}
	if insecure {
		metricOpts = append(metricOpts, otlpmetricgrpc.WithInsecure())
	}
	metricExporter, err := otlpmetricgrpc.New(ctx, metricOpts...)
	if err != nil {
		_ = tp.Shutdown(ctx)
		return noopShutdown, fmt.Errorf("observability: metric exporter: %w", err)
	}
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
	)

	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	shutdown := func(ctx context.Context) error {
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

// buildResource assembles the SDK resource with service identity.
func buildResource(ctx context.Context, cfg Config) (*resource.Resource, error) {
	name := strings.TrimSpace(os.Getenv("OTEL_SERVICE_NAME"))
	if name == "" {
		name = cfg.ServiceName
	}
	if name == "" {
		name = "indexer"
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
// falls back to def when unset / malformed. Out-of-range values are
// rejected (defaulted) rather than clamped so an operator who set "10"
// instead of "0.1" sees zero traces and looks at the env.
func resolveSampleRate(def float64) float64 {
	raw := strings.TrimSpace(os.Getenv("OTEL_TRACES_SAMPLER_ARG"))
	if raw != "" {
		v, err := strconv.ParseFloat(raw, 64)
		if err == nil && v >= 0 && v <= 1 {
			return v
		}
	}
	if def > 0 && def <= 1 {
		return def
	}
	return 0.1
}

// stripScheme drops a leading http:// or https:// prefix from the endpoint
// because OTLP gRPC takes a host:port pair and would otherwise treat the
// scheme as part of the authority.
func stripScheme(endpoint string) string {
	for _, prefix := range []string{"http://", "https://", "grpc://"} {
		if strings.HasPrefix(endpoint, prefix) {
			return strings.TrimPrefix(endpoint, prefix)
		}
	}
	return endpoint
}
