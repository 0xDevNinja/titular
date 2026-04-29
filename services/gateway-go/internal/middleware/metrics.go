package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// httpMetricsMeterName is the otel.Meter scope used for the gateway's
// hand-instrumented HTTP histogram. Keeping it under the gateway-go
// module path lands it on the otel_scope_name label so dashboards can
// distinguish gateway samples from indexer samples that share the same
// instrument name.
const httpMetricsMeterName = "github.com/0xDevNinja/titular/services/gateway-go/middleware/http"

// httpServerDurationName is the canonical OTel HTTP server duration
// metric. Keeping the name aligned with the OTel HTTP semantic
// convention means future drop-in replacements (otelhttp / otelgin
// metrics) emit on the same series and dashboards do not break.
const httpServerDurationName = "http_server_request_duration_seconds"

// httpServerDurationBuckets are the explicit histogram boundaries used
// by the HTTP server duration instrument. The 200ms boundary is
// deliberate — it matches the M4 SLO from the milestone brief — and
// the surrounding distribution lets Grafana's
// histogram_quantile(0.50/0.95/0.99) panels plot meaningful curves at
// the 1k-RPS load target without losing granularity at the long tail.
//
// All values are MILLISECONDS at definition time and divided by 1000
// at recording so the exposition is in seconds (matching the
// semconv-suffixed `_seconds` metric name and the dashboard rate()
// queries).
var httpServerDurationBuckets = []float64{
	0.005, // 5ms
	0.010, // 10ms
	0.025, // 25ms
	0.050, // 50ms
	0.100, // 100ms
	0.200, // 200ms — SLO P99 boundary
	0.500, // 500ms
	1.000, // 1s
	2.500, // 2.5s
}

// HTTPMetrics returns a gin handler that records the wall-clock
// duration of every served request into the
// `http_server_request_duration_seconds` histogram on the global OTel
// meter provider. Mounted alongside OTelMiddleware (which is
// trace-only) in router.go so the gateway's /metrics surface actually
// carries the histogram data the M4 Grafana dashboard charts.
//
// Labels (kept low-cardinality on purpose — see the No-PII contract in
// internal/observability/metrics.go):
//
//   - http_route                 the gin route template (c.FullPath()),
//                                NOT the raw URL — `/agents/:id` instead
//                                of `/agents/0xabc…`. Falls back to "unknown"
//                                when the request did not match a route
//                                (404 catch-all) so the cardinality stays
//                                bounded by the registered route table.
//   - http_request_method        upper-case HTTP verb. The semconv set
//                                is fixed and small.
//   - http_response_status_code  integer status code as a string. Bounded
//                                by the HTTP spec.
//
// Label keys mirror the OTel HTTP semantic convention so Grafana
// dashboards under ops/grafana/dashboards/gateway.json find the series
// they expect (e.g. `http_response_status_code=~"5.."` for the 5xx
// rate panel).
//
// On the no-op MeterProvider (observability disabled or boot order
// hasn't reached Init yet) this middleware is allocation-light: the
// returned histogram is itself a no-op and Record drops the sample on
// the floor without taking a global lock. We still call FullPath /
// c.Writer.Status / time.Since so the path is identical with or
// without exporters wired.
func HTTPMetrics() gin.HandlerFunc {
	meter := otel.GetMeterProvider().Meter(httpMetricsMeterName)
	hist, err := meter.Float64Histogram(
		httpServerDurationName,
		metric.WithUnit("s"),
		metric.WithDescription("Duration of HTTP server requests in seconds."),
		metric.WithExplicitBucketBoundaries(httpServerDurationBuckets...),
	)
	if err != nil {
		// A genuine instrument-creation failure should not take the
		// gateway out — surface it via a no-op so requests still flow
		// and the missing samples become the operator's signal that
		// something is wrong with the meter. Test_GatewayInstrumentsBuilt
		// in observability/metrics_test.go is the canary for the happy
		// path.
		return func(c *gin.Context) { c.Next() }
	}

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		elapsed := time.Since(start).Seconds()

		// FullPath() returns the matched route template (e.g. "/agents/:id");
		// when the request did not match (e.g. 404), it returns "" — we
		// fold that into a fixed "unknown" bucket so unhandled paths
		// never expand the cardinality of the route label.
		route := c.FullPath()
		if route == "" {
			route = "unknown"
		}
		method := c.Request.Method
		status := strconv.Itoa(c.Writer.Status())

		hist.Record(c.Request.Context(), elapsed,
			metric.WithAttributes(
				attribute.String("http_route", route),
				attribute.String("http_request_method", method),
				attribute.String("http_response_status_code", status),
			),
		)
	}
}

