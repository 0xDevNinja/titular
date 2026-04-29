// Package router builds the gateway HTTP entry point.
//
// The outer engine is Gin. M2/M3 endpoints (the agent and job fixtures) still
// live behind chi handlers because they use chi.URLParam to read path
// parameters; chi is mounted as an http.Handler under /v1/* so those handlers
// continue to work without modification while new endpoints introduced from
// M4 onwards can be wired natively into Gin.
package router

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/handlers"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/middleware"
)

// Config controls how the gateway router is constructed.
//
// Logger is the zerolog instance handed to the Log/Recovery middlewares. When
// zero-valued, the package-level zerolog logger is used so day-one callers do
// not have to wire one in.
//
// Service identifies the binary in structured log lines.
//
// CORS / RateLimit are passed through verbatim to the corresponding
// middlewares. See middleware.CORSConfig and middleware.RateLimitConfig.
//
// TrustedProxies, when non-nil, is forwarded to gin's SetTrustedProxies.
// Leaving it nil applies gin's default of trusting all proxies; deployments
// that terminate TLS at a known LB should set this explicitly.
type Config struct {
	Logger         zerolog.Logger
	Service        string
	CORS           middleware.CORSConfig
	RateLimit      middleware.RateLimitConfig
	TrustedProxies []string
}

// DefaultConfig returns a Config suitable for local development.
func DefaultConfig() Config {
	return Config{
		Logger:    log.Logger,
		Service:   "gateway",
		CORS:      middleware.DefaultCORSConfig(),
		RateLimit: middleware.RateLimitConfig{}, // disabled by default
	}
}

// New builds and returns the HTTP handler. agentHandlers must be non-nil; it
// is injected so tests can provide a fixture-backed instance without the
// singleton loader.
func New(agentHandlers *handlers.AgentHandlers) http.Handler {
	jh, err := handlers.NewJobHandlers()
	if err != nil {
		// Fatal startup: fixture files missing or malformed.
		panic(fmt.Sprintf("load job handlers: %v", err))
	}
	return NewWithHandlers(agentHandlers, jh)
}

// NewWithHandlers builds and returns the HTTP handler with explicit handler
// injection. Preferred in tests so fixtures are controlled by the caller. It
// uses DefaultConfig; callers needing custom CORS/rate-limit settings should
// use NewWithConfig.
func NewWithHandlers(agentHandlers *handlers.AgentHandlers, jobHandlers *handlers.JobHandlers) http.Handler {
	return NewWithConfig(DefaultConfig(), agentHandlers, jobHandlers)
}

// NewWithConfig builds the gateway router with full control over middleware
// configuration. The middleware order, from outer to inner, is:
//
//  1. Recovery     — wraps everything so panics surface as 500.
//  2. RequestID    — guarantees the rest of the chain has an id.
//  3. Log          — observes the resolved status, including 4xx/5xx flips.
//  4. CORS         — emits headers / handles preflight before rate-limit so
//     legitimate preflight requests cannot be limited away.
//  5. RateLimit    — final guard; per-key limiter applied to all routes.
//
// Routes:
//
//   - GET /healthz            — chassis liveness probe (M4 #85).
//   - GET /health             — legacy probe retained for M2/M3 callers.
//   - /v1/*                   — chi sub-router for M2/M3 fixture endpoints.
func NewWithConfig(
	cfg Config,
	agentHandlers *handlers.AgentHandlers,
	jobHandlers *handlers.JobHandlers,
) http.Handler {
	engine := gin.New()

	// Trust nothing by default beyond what the caller specified.
	if cfg.TrustedProxies != nil {
		_ = engine.SetTrustedProxies(cfg.TrustedProxies)
	}

	logger := cfg.Logger
	if reflect.DeepEqual(logger, zerolog.Logger{}) {
		logger = log.Logger
	}

	engine.Use(middleware.Recovery(logger))
	engine.Use(middleware.RequestID())
	engine.Use(middleware.Log(middleware.LogConfig{
		Logger:    logger,
		Service:   cfg.Service,
		SkipPaths: []string{"/healthz"},
	}))
	engine.Use(middleware.CORS(cfg.CORS))
	engine.Use(middleware.RateLimit(cfg.RateLimit))

	engine.GET("/healthz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.String(http.StatusOK, `{"status":"ok"}`)
	})

	// Mount chi for the legacy /v1 fixture endpoints. We bind under a
	// catch-all so chi sees the full path, and route inside chi exactly as
	// before. This keeps the existing handler signatures and tests intact.
	v1 := buildV1Chi(agentHandlers, jobHandlers)
	engine.Any("/v1/*proxyPath", gin.WrapH(v1))

	// Legacy /health probe retained for backwards compatibility with M2/M3
	// integration scripts. New consumers should use /healthz.
	engine.GET("/health", func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.String(http.StatusOK, `{"status":"ok"}`)
	})

	return engine
}

// buildV1Chi assembles the chi router that serves the M2 (agents) and M3
// (jobs) fixture endpoints. It is intentionally free of middleware — the
// outer Gin chain has already injected the request id into the underlying
// request.Context, so chi handlers can continue calling
// middleware.GetRequestID(r.Context()) when emitting errors.
func buildV1Chi(
	agentHandlers *handlers.AgentHandlers,
	jobHandlers *handlers.JobHandlers,
) http.Handler {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		// M2 — agents
		r.Get("/agents", agentHandlers.ListAgents)
		r.Get("/agents/{id}", agentHandlers.GetAgent)
		r.Get("/agents/{id}/trades", agentHandlers.ListAgentTrades)

		// M3 — ACP jobs
		r.Get("/jobs", jobHandlers.ListJobs)
		r.Post("/jobs/prepare", jobHandlers.PrepareJob)
		r.Get("/jobs/{id}", jobHandlers.GetJob)
	})

	return r
}
