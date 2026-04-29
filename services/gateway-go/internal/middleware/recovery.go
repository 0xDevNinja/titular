package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// recoveryBody is the public 500 body returned to clients on panic. It is a
// fixed string so we never accidentally surface the panic value or stack trace
// to callers — both are written to the structured log instead.
const recoveryBody = `{"error":{"code":"internal_server_error","message":"an unexpected error occurred"}}`

// Recovery returns a Gin middleware that catches any panic raised downstream,
// logs the panic value together with the goroutine stack at ERROR level, and
// terminates the request with a generic 500 envelope.
//
// The supplied logger is used for the structured panic record. When zero, a
// logger that discards output is used so the middleware is still safe to wire
// into tests.
//
// We deliberately do not call gin's stock CustomRecovery here so the response
// shape stays in our control and the log line carries the request_id from the
// RequestID middleware (when chained after it).
func Recovery(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}

			logger.Error().
				Str("request_id", GetRequestIDFromGin(c)).
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Str("client_ip", c.ClientIP()).
				Interface("panic", rec).
				Bytes("stack", debug.Stack()).
				Msg("recovered from panic")

			// Abort the chain. AbortWithStatusJSON would re-encode each call;
			// we want a fixed wire format independent of gin's renderer.
			c.Header("Content-Type", "application/json; charset=utf-8")
			c.String(http.StatusInternalServerError, recoveryBody)
			c.Abort()
		}()
		c.Next()
	}
}
