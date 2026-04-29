package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TraceIDHeader is the response header carrying the active span's trace
// id back to the client. Operators correlating a 5xx with their tracing
// backend (Tempo / Honeycomb / Datadog) read it out of the response;
// the value is the W3C-format hex trace id, never sensitive.
const TraceIDHeader = "X-Trace-Id"

// requestIDSpanAttr is the canonical span attribute used to correlate a
// span with the gateway's per-request log id. Mirrors the request_id
// field emitted by middleware.Log so an operator can pivot from a log
// line to the full trace and back.
const requestIDSpanAttr = "gateway.request_id"

// OTelMiddleware returns the canonical OTel inbound HTTP middleware
// (otelgin) configured for the gateway. Mounted in the router chain
// AFTER RequestID and Recovery so:
//
//   - the request_id is already on the gin context when otelgin starts
//     the span and the TraceCorrelation middleware below stamps it on,
//   - panics rendered as 500 by Recovery still carry a span (otelgin's
//     defer span.End fires after Recovery's, since it sits OUTSIDE in
//     the chain).
//
// service is forwarded to otelgin as the server name; we pass the
// gateway's own Service tag so spans emitted from this binary share
// the same service.name resource attribute applied at SDK boot.
//
// When the global TracerProvider is the default no-op (i.e. OTel was
// not initialised because OTEL_EXPORTER_OTLP_ENDPOINT is unset),
// otelgin still runs but every span is a no-op — cheap, no allocation
// pressure, and importantly, downstream tracer.Start callers can rely
// on a non-nil span context being on the request.
func OTelMiddleware(service string) gin.HandlerFunc {
	return otelgin.Middleware(service)
}

// TraceCorrelation returns a Gin middleware that bridges the gateway's
// existing request_id correlation onto the active OTel span:
//
//   - sets gateway.request_id as a span attribute so trace search by
//     request id works,
//   - mirrors the active trace id into the X-Trace-Id response header
//     so a client receiving a 5xx can paste the id straight into a
//     tracing backend and find the trace.
//
// Mount AFTER OTelMiddleware so the active span is already on the
// request context when this runs. No-ops cleanly when no span is
// recording (the default no-op tracer path), so this can sit in the
// chain unconditionally.
func TraceCorrelation() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		if span.IsRecording() {
			if rid := GetRequestIDFromGin(c); rid != "" {
				span.SetAttributes(attribute.String(requestIDSpanAttr, rid))
			}
		}
		// Set the response header even if the span is not recording —
		// the trace id of an unsampled span is still a valid lookup
		// key in most backends, and clients should receive a stable
		// header shape regardless of sampling decision.
		if sc := span.SpanContext(); sc.IsValid() {
			c.Header(TraceIDHeader, sc.TraceID().String())
		}
		c.Next()
	}
}
