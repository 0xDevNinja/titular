// Package middleware provides Gin-compatible HTTP middleware for the gateway.
//
// The middlewares are intentionally small and composable; each is exposed as a
// gin.HandlerFunc so callers can attach them in any order via engine.Use(...).
package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDHeader is the canonical header name used to carry the request id on
// both the inbound and outbound HTTP message.
const RequestIDHeader = "X-Request-Id"

// RequestIDContextKey is the gin.Context key under which the resolved request
// id is stored. Downstream middleware and handlers should retrieve it via
// GetRequestID rather than reading the key directly.
const RequestIDContextKey = "request_id"

// requestIDCtxKey is the unexported type used to namespace the value when it is
// stored on the underlying request.Context. Using a distinct unexported type
// avoids collisions with other packages that share the same string key.
type requestIDCtxKey struct{}

// RequestID returns a Gin middleware that:
//
//   - reuses an inbound X-Request-Id header when present (and well-formed), or
//   - generates a fresh UUIDv4 when absent or malformed,
//   - stores the value on gin.Context under RequestIDContextKey,
//   - mirrors the value into the underlying request.Context so non-Gin handlers
//     mounted via WrapH can read it,
//   - sets the X-Request-Id response header before the handler runs so it is
//     present even if the handler aborts.
//
// An "inbound" id is accepted only if it parses as a valid UUID. This prevents
// clients from injecting log-poisoning values such as newlines or oversized
// strings.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(RequestIDHeader)
		if !isValidRequestID(id) {
			id = uuid.NewString()
		}

		c.Set(RequestIDContextKey, id)

		// Mirror onto the request.Context so chi/std-lib handlers behind
		// gin.WrapH can read the id via GetRequestID(r.Context()).
		ctx := context.WithValue(c.Request.Context(), requestIDCtxKey{}, id)
		c.Request = c.Request.WithContext(ctx)

		c.Header(RequestIDHeader, id)
		c.Next()
	}
}

// GetRequestID extracts the request id from a context, returning an empty
// string when none is present. It is safe to call from any handler regardless
// of whether the request originated through Gin or a wrapped http.Handler.
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if id, ok := ctx.Value(requestIDCtxKey{}).(string); ok {
		return id
	}
	return ""
}

// GetRequestIDFromGin extracts the request id from a *gin.Context.
func GetRequestIDFromGin(c *gin.Context) string {
	if v, ok := c.Get(RequestIDContextKey); ok {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}

// isValidRequestID reports whether the supplied string is a syntactically valid
// UUID. We deliberately do not accept arbitrary opaque ids — see RequestID.
func isValidRequestID(s string) bool {
	if s == "" {
		return false
	}
	_, err := uuid.Parse(s)
	return err == nil
}

// Compile-time assertion that the middleware satisfies gin.HandlerFunc and is
// usable as a standard http middleware shim if ever needed in tests.
var _ http.Handler = (http.HandlerFunc)(nil)
