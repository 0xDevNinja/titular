package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ErrCORSWildcardWithCredentials is returned by CORSConfig.Validate when the
// caller has combined an "*" allowed origin with AllowCredentials=true without
// opting in via AllowReflectedOrigins. Reflecting any origin while honouring
// cookies turns CORS into a permissive credentialed reflector — exactly the
// configuration browsers warn about — so we refuse it by default and force
// the operator to acknowledge the risk explicitly.
var ErrCORSWildcardWithCredentials = errors.New(
	"middleware: CORS wildcard origin with AllowCredentials=true is unsafe; " +
		"set AllowReflectedOrigins=true to opt in, or list explicit origins",
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
// AllowedOrigins is not the wildcard, unless AllowReflectedOrigins is set —
// see below).
//
// AllowReflectedOrigins, when true, permits the unsafe combination of
// AllowedOrigins=["*"] + AllowCredentials=true. In that mode the middleware
// echoes the request Origin verbatim and emits Access-Control-Allow-
// Credentials. This is effectively a credentialed reflector and should only
// be enabled in controlled local-dev / preview environments. Default is
// false; CORSConfig.Validate refuses to start with the unsafe combination
// otherwise.
//
// MaxAgeSeconds is advertised in Access-Control-Max-Age so browsers can cache
// preflight results.
type CORSConfig struct {
	AllowedOrigins        []string
	AllowedMethods        []string
	AllowedHeaders        []string
	ExposedHeaders        []string
	AllowCredentials      bool
	AllowReflectedOrigins bool
	MaxAgeSeconds         int
}

// Validate returns a non-nil error when the configuration is unsafe and
// callers have not explicitly opted in to the unsafe combination. Today the
// only check is the wildcard+credentials anti-pattern; future audit findings
// can be added here without touching every call-site.
func (c CORSConfig) Validate() error {
	hasWildcard := false
	for _, o := range c.AllowedOrigins {
		if strings.TrimSpace(o) == "*" {
			hasWildcard = true
			break
		}
	}
	if hasWildcard && c.AllowCredentials && !c.AllowReflectedOrigins {
		return ErrCORSWildcardWithCredentials
	}
	return nil
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
// the CORS spec, never combined with credentials unless the caller has
// explicitly opted in via AllowReflectedOrigins=true (see CORSConfig).
//
// CORS panics if cfg fails Validate() — call NewCORS for a non-panicking
// error path. The panic is intentional: an unsafe CORS config at startup is
// a configuration bug, not a runtime condition the request path should
// recover from.
func CORS(cfg CORSConfig) gin.HandlerFunc {
	h, err := NewCORS(cfg)
	if err != nil {
		panic(err)
	}
	return h
}

// NewCORS is the error-returning constructor. Use this from cmd/* so an
// unsafe configuration produces a graceful Fatal log line instead of a panic
// stack.
func NewCORS(cfg CORSConfig) (gin.HandlerFunc, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

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
			// Validated combo: caller has explicitly opted in via
			// AllowReflectedOrigins. Echo the exact origin (cannot use "*"
			// alongside credentials per the CORS spec).
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
	}, nil
}
