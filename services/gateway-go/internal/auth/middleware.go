package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Gin context keys used by RequireAuth to publish the resolved identity to
// downstream handlers. Tests and handlers should use the typed accessors
// (AddressFromContext, JTIFromContext) rather than reading these directly so
// that the storage shape can evolve without rippling through callers.
const (
	ContextKeyAddress = "auth.address"
	ContextKeyJTI     = "auth.jti"
)

// RequireAuth returns a Gin middleware that enforces a valid SIWE-issued JWT
// AND an active Redis session.
//
// The double-check (signature + session) is what makes server-side logout
// effective: a JWT alone would remain valid until exp regardless of logout,
// because HS256 has no built-in revocation. By gating every protected request
// on the Redis session, deleting that session immediately invalidates the
// token. This is also what makes the threat model survive a leaked JWT — an
// operator can flush Redis to evict every active session.
//
// Failure modes are surfaced as a single generic 401 envelope; we never
// distinguish "no token", "bad signature", "wrong issuer" or "session
// missing" in the response code, only in server-side logs at debug level.
func RequireAuth(s *Handlers) gin.HandlerFunc {
	return func(c *gin.Context) {
		tok := bearerFromHeader(c.GetHeader("Authorization"))
		if tok == "" {
			writeError(c, http.StatusUnauthorized, "unauthorized", "authentication required")
			return
		}
		claims, err := s.cfg.Signer.Parse(tok)
		if err != nil {
			writeError(c, http.StatusUnauthorized, "unauthorized", "authentication required")
			return
		}
		// Re-derive identity from the trusted store, not from the token —
		// defence in depth against a forged-but-valid-looking token (which
		// is impossible under HS256 + secret pin, but cheap to defend
		// against and keeps the audit trail honest).
		addr, err := s.cfg.Store.SessionExists(c.Request.Context(), claims.ID)
		if err != nil {
			writeError(c, http.StatusUnauthorized, "unauthorized", "authentication required")
			return
		}
		c.Set(ContextKeyAddress, addr)
		c.Set(ContextKeyJTI, claims.ID)
		c.Next()
	}
}

// AddressFromContext returns the authenticated EIP-55 lowercased address that
// RequireAuth attached to the context, or an empty string when the request
// was not authenticated. Safe to call on any *gin.Context.
func AddressFromContext(c *gin.Context) string {
	if v, ok := c.Get(ContextKeyAddress); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// JTIFromContext returns the JWT id of the active session, or an empty string
// when the request was not authenticated.
func JTIFromContext(c *gin.Context) string {
	if v, ok := c.Get(ContextKeyJTI); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
