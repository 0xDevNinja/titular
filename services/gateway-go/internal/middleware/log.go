package middleware

import (
	"io"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// LogConfig configures the structured request logger.
//
// Logger is the zerolog instance lines are emitted on. When zero-valued, a
// logger writing to io.Discard is used so the middleware never panics when
// invoked from tests that forgot to wire one in.
//
// Service is recorded on every line under the "service" field so that logs
// from multiple binaries can be correlated centrally.
//
// SkipPaths is a small set of exact path strings that should not be logged
// (typically health probe endpoints that would otherwise drown out real
// traffic).
//
// Level, when set to a non-NoLevel value, overrides the level of the supplied
// Logger for this middleware only — useful when the request log should be
// noisier (Debug) or quieter (Warn) than the rest of the binary. Leaving
// Level at its zero value (zerolog.NoLevel) inherits the parent logger's
// level.
type LogConfig struct {
	Logger    zerolog.Logger
	Service   string
	SkipPaths []string
	Level     zerolog.Level
}

// Log returns a Gin middleware that emits one structured log line per request
// once the handler chain completes.
//
// Fields:
//
//   - service     — caller-supplied identifier, e.g. "gateway"
//   - request_id  — value injected by RequestID
//   - method      — HTTP method
//   - path        — URL path (query string is intentionally omitted)
//   - status      — final response status
//   - latency_ms  — wall-clock time the request spent inside the chain
//   - client_ip   — gin's resolved client ip (respects trusted proxies)
//   - bytes       — bytes written to the response
//   - errors      — joined error strings from c.Errors, when non-empty
//
// The level is INFO for status < 500, WARN for status in [400, 500), and ERROR
// for status >= 500. We never include the query string or request body to
// avoid accidentally logging tokens, signatures or PII.
func Log(cfg LogConfig) gin.HandlerFunc {
	skip := make(map[string]struct{}, len(cfg.SkipPaths))
	for _, p := range cfg.SkipPaths {
		skip[p] = struct{}{}
	}

	logger := cfg.Logger
	// zerolog.Logger contains an unexported []byte slice, so it is not
	// directly comparable. Use reflect.DeepEqual against a zero value to
	// detect "caller forgot to wire one in" and substitute a discard sink.
	if reflect.DeepEqual(logger, zerolog.Logger{}) {
		logger = zerolog.New(io.Discard)
	}

	// Per-middleware level override. zerolog.NoLevel is the zero value, which
	// leaves the parent logger's level untouched.
	if cfg.Level != zerolog.NoLevel {
		logger = logger.Level(cfg.Level)
	}

	service := cfg.Service
	if service == "" {
		service = "gateway"
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		if _, skipped := skip[path]; skipped {
			return
		}

		status := c.Writer.Status()
		var event *zerolog.Event
		switch {
		case status >= 500:
			event = logger.Error()
		case status >= 400:
			event = logger.Warn()
		default:
			event = logger.Info()
		}

		event = event.
			Str("service", service).
			Str("request_id", GetRequestIDFromGin(c)).
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", status).
			Int64("latency_ms", time.Since(start).Milliseconds()).
			Str("client_ip", c.ClientIP()).
			Int("bytes", c.Writer.Size())

		if errs := c.Errors.String(); errs != "" {
			event = event.Str("errors", errs)
		}

		event.Msg("request")
	}
}
