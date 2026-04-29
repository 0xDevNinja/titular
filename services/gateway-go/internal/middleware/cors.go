package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig configures the CORS middleware.
//
// AllowedOrigins is the explicit list of permitted origins, e.g.
// "https://app.titular.xyz". A single entry of "*" enables a wildcard match;
// in that mode the response will not echo Access-Control-Allow-Credentials
// because browsers reject the combination.
//
// AllowedMethods/AllowedHeaders are advertised on preflight responses.
// AllowCredentials toggles the credentials header (only honoured when
// AllowedOrigins is not the wildcard).
//
// MaxAgeSeconds is advertised in Access-Control-Max-Age so browsers can cache
// preflight results.
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAgeSeconds    int
}

// DefaultCORSConfig returns a sensible default appropriate for local
// development: no origins allowed (CORS effectively disabled) but with the
// standard methods and headers pre-populated so callers only need to set
// AllowedOrigins to enable cross-origin access.
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedMethods: []string{
			http.MethodGet, http.MethodPost, http.MethodPut,
			http.MethodPatch, http.MethodDelete, http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization",
			RequestIDHeader,
		},
		ExposedHeaders: []string{RequestIDHeader},
		MaxAgeSeconds:  600,
	}
}

// CORS returns a Gin middleware enforcing the supplied CORSConfig.
//
// For preflight requests (OPTIONS with an Access-Control-Request-Method
// header) the middleware writes a 204 and aborts the chain. For other
// requests it sets the Allow-Origin header (when the request origin is
// permitted) and lets the chain continue.
//
// Origins are matched exactly. The wildcard entry "*" is honoured but, per
// the CORS spec, never combined with credentials.
func CORS(cfg CORSConfig) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(cfg.AllowedOrigins))
	wildcard := false
	for _, o := range cfg.AllowedOrigins {
		o = strings.TrimSpace(o)
		if o == "" {
			continue
		}
		if o == "*" {
			wildcard = true
			continue
		}
		allowed[o] = struct{}{}
	}

	allowMethods := strings.Join(cfg.AllowedMethods, ", ")
	allowHeaders := strings.Join(cfg.AllowedHeaders, ", ")
	exposeHeaders := strings.Join(cfg.ExposedHeaders, ", ")
	maxAge := strconv.Itoa(cfg.MaxAgeSeconds)

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Vary header so caches do not collide between origins.
		c.Writer.Header().Add("Vary", "Origin")

		allow := ""
		switch {
		case origin == "":
			// Same-origin / non-browser caller; nothing to advertise.
		case wildcard && !cfg.AllowCredentials:
			allow = "*"
		case wildcard && cfg.AllowCredentials:
			// With credentials we must echo the exact origin, not "*".
			allow = origin
		default:
			if _, ok := allowed[origin]; ok {
				allow = origin
			}
		}

		if allow != "" {
			c.Header("Access-Control-Allow-Origin", allow)
			if cfg.AllowCredentials && allow != "*" {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			if exposeHeaders != "" {
				c.Header("Access-Control-Expose-Headers", exposeHeaders)
			}
		}

		// Handle preflight.
		if c.Request.Method == http.MethodOptions &&
			c.GetHeader("Access-Control-Request-Method") != "" {
			if allow != "" {
				if allowMethods != "" {
					c.Header("Access-Control-Allow-Methods", allowMethods)
				}
				if allowHeaders != "" {
					c.Header("Access-Control-Allow-Headers", allowHeaders)
				}
				if cfg.MaxAgeSeconds > 0 {
					c.Header("Access-Control-Max-Age", maxAge)
				}
			}
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
