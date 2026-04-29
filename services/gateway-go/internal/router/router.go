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

	"github.com/0xDevNinja/titular/services/gateway-go/internal/auth"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/graph"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/handlers"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/middleware"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/openapi"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/sse"
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
// TrustedProxies is forwarded to gin's SetTrustedProxies.
//
// SECURITY: Gin's default is to trust ALL proxies, which makes c.ClientIP()
// honour any X-Forwarded-For header — trivially spoofable, and a direct
// rate-limit / log-attribution bypass. We therefore default to trusting NO
// proxies (TrustedProxies == nil → SetTrustedProxies([]string{})), so
// c.ClientIP() falls back to c.Request.RemoteAddr, which is unspoofable.
// Operators running behind a known L7 proxy MUST set this explicitly via
// GATEWAY_TRUSTED_PROXIES.
//
// Auth, when non-nil, is the SIWE auth handler bundle. The router mounts
// /auth/siwe/nonce, /auth/siwe/verify and /auth/logout when present. When
// nil, the auth endpoints are simply absent — useful in tests that don't
// want to spin up Redis.
//
// SIWS, when non-nil, is the Sign-In-With-Solana handler bundle. The
// router mounts /auth/siws/nonce and /auth/siws/verify when present. The
// SIWS path reuses the SIWE /auth/logout endpoint via the shared JWT
// signer and session store, so a SIWS-minted token logs out through the
// same handler as a SIWE-minted one.
type Config struct {
	Logger         zerolog.Logger
	Service        string
	CORS           middleware.CORSConfig
	RateLimit      middleware.RateLimitConfig
	TrustedProxies []string
	Auth           *auth.Handlers
	SIWS           *auth.SIWSHandlers

	// API, when non-nil, is the Postgres-backed REST API handler bundle
	// introduced in M4 (#88). Mounted under /api/v1; absent when the
	// gateway runs without a database connection (e.g. in tests that only
	// exercise the chassis).
	API *handlers.API

	// GraphQL, when non-nil, is the GraphQL handler introduced in
	// M4 (#89). Mounted at /graphql for queries / mutations and
	// /graphql (websocket upgrade) for subscriptions. Absent when the
	// gateway is started without a Postgres pool — Query and
	// Subscription resolvers both need a Store, so a missing API
	// implies no GraphQL surface either.
	GraphQL *graph.Handler

	// SSE, when non-nil, is the Server-Sent Events handler introduced
	// in M4 (#90). Mounted at /events for one-way NATS-fanout streams
	// over plain HTTP. Absent when the gateway is started without
	// NATS — the multiplexer requires a NATS connection, so a missing
	// GATEWAY_NATS_URL implies no SSE surface either.
	SSE *sse.Handler

	// OpenAPI, when non-nil, is the OpenAPI 3.0 spec bundle introduced
	// in M4 (#91). Mounting is unconditional for the raw spec endpoints
	// (/swagger/doc.json, /swagger/doc.yaml) so SDK generators can
	// always fetch the document; the Swagger UI viewer is gated by
	// the bundle's own EnableUI flag (operator-set via
	// GATEWAY_SWAGGER_UI). Absent in tests that don't care about
	// documentation surface.
	OpenAPI *openapi.Spec
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

// Built bundles the constructed handler with any lifecycle hooks owned by
// the router (today: the rate-limiter sweeper goroutine). Callers that care
// about clean shutdown should use NewWithConfigLifecycle and invoke Stop on
// SIGINT/SIGTERM. The simpler NewWithConfig discards the hook for callers
// who do not.
type Built struct {
	Handler http.Handler
	stop    func()
}

// Stop releases any lifecycle resources (currently: the rate-limit sweeper).
// Safe to call multiple times.
func (b Built) Stop() {
	if b.stop != nil {
		b.stop()
	}
}

// NewWithConfig builds the gateway router with full control over middleware
// configuration. The middleware order, from outer to inner, is:
//
//  1. RequestID    — outermost so every downstream record (including panic
//     recovery) carries the id.
//  2. Recovery     — converts panics from the rest of the chain into a 500
//     and logs with the request id assigned in (1).
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
	return NewWithConfigLifecycle(cfg, agentHandlers, jobHandlers).Handler
}

// NewWithConfigLifecycle is the lifecycle-aware variant. The returned Built
// carries a Stop() method that shuts down the rate-limit sweeper (and any
// future background goroutines). cmd/gateway wires this into its
// SIGINT/SIGTERM path.
func NewWithConfigLifecycle(
	cfg Config,
	agentHandlers *handlers.AgentHandlers,
	jobHandlers *handlers.JobHandlers,
) Built {
	engine := gin.New()

	// SECURITY: When TrustedProxies is nil (env unset), trust NO proxies so
	// c.ClientIP() ignores spoofable X-Forwarded-For headers and falls back
	// to c.Request.RemoteAddr. Operators behind a known L7 proxy must set
	// GATEWAY_TRUSTED_PROXIES explicitly (see service README).
	proxies := cfg.TrustedProxies
	if proxies == nil {
		proxies = []string{}
	}
	_ = engine.SetTrustedProxies(proxies)

	logger := cfg.Logger
	if reflect.DeepEqual(logger, zerolog.Logger{}) {
		logger = log.Logger
	}

	limiter := middleware.NewRateLimiter(cfg.RateLimit)

	engine.Use(middleware.RequestID())
	engine.Use(middleware.Recovery(logger))
	engine.Use(middleware.Log(middleware.LogConfig{
		Logger:    logger,
		Service:   cfg.Service,
		SkipPaths: []string{"/healthz"},
	}))
	engine.Use(middleware.CORS(cfg.CORS))
	engine.Use(limiter.Handler)

	engine.GET("/healthz", HealthZ)

	// SIWE/SIWS auth — mounted before /v1 so the routes can never be
	// accidentally shadowed by a /v1/* catch-all. The /auth group inherits
	// the rate limiter (#86 OWASP gate) but, by design, NOT RequireAuth —
	// the verify endpoints are what bootstrap the session in the first
	// place.
	//
	// SIWS reuses the SIWE Logout handler (cfg.Auth.Logout) because the
	// JWT signer + Redis session store are shared across both paths.
	if cfg.Auth != nil || cfg.SIWS != nil {
		authGroup := engine.Group("/auth")
		if cfg.Auth != nil {
			authGroup.POST("/siwe/nonce", cfg.Auth.Nonce)
			authGroup.POST("/siwe/verify", cfg.Auth.Verify)
			authGroup.POST("/logout", cfg.Auth.Logout)
		}
		if cfg.SIWS != nil {
			authGroup.POST("/siws/nonce", cfg.SIWS.Nonce)
			authGroup.POST("/siws/verify", cfg.SIWS.Verify)
		}
	}

	// Mount chi for the legacy /v1 fixture endpoints. We bind under a
	// catch-all so chi sees the full path, and route inside chi exactly as
	// before. This keeps the existing handler signatures and tests intact.
	v1 := buildV1Chi(agentHandlers, jobHandlers)
	engine.Any("/v1/*proxyPath", gin.WrapH(v1))

	// M4 (#88) — Postgres-backed REST API. Mounted under /api/v1 to keep it
	// distinct from the fixture-backed /v1 surface above. Read-only by
	// design, so it lives outside the auth wall: SDK consumers, browser
	// clients and ops dashboards all need to hit these without holding a
	// SIWE session. If a deployment wants to gate them, wrap the group in
	// auth.RequireAuth(cfg.Auth) before calling Register.
	if cfg.API != nil {
		api := engine.Group("/api/v1")
		cfg.API.Register(api)
	}

	// M4 (#89) — GraphQL surface. Mounted at the engine root rather
	// than under /api/v1 because:
	//   1. The websocket upgrade path uses GET /graphql, which would
	//      collide with a REST resource of the same name; a top-level
	//      mount avoids that.
	//   2. SDK consumers conventionally hit `/graphql` rather than
	//      `/api/v1/graphql`, so the unprefixed path is the one
	//      developers expect to find.
	// As with /api/v1 the surface is read-only by design; deployments
	// that want auth gating should attach their own middleware to the
	// returned group before invoking Register.
	if cfg.GraphQL != nil {
		gqlGroup := engine.Group("")
		cfg.GraphQL.Register(gqlGroup)
	}

	// M4 (#90) — SSE multiplexer. Mounted at the engine root so the
	// canonical path is `/events`, matching the SDK's expected URL
	// shape. Inherits the rate-limiter so a single client cannot open
	// arbitrarily many concurrent connections; auth is intentionally
	// absent because the surface is read-only and parallels the
	// public-read GraphQL subscription transport (#89).
	if cfg.SSE != nil {
		sseGroup := engine.Group("")
		cfg.SSE.Register(sseGroup)
	}

	// M4 (#91) — OpenAPI 3.0 spec + optional Swagger UI viewer.
	// Mounted at the engine root so the canonical paths are
	// /swagger/doc.{json,yaml} regardless of which other M4 surfaces
	// are wired. Inherits the rate-limiter (one shared budget for the
	// whole gateway). The viewer HTML is dev-only, gated server-side
	// by the spec's EnableUI flag — see internal/openapi for the env
	// var binding.
	if cfg.OpenAPI != nil {
		oasGroup := engine.Group("")
		cfg.OpenAPI.Register(oasGroup)
	}

	// Legacy /health probe retained for backwards compatibility with M2/M3
	// integration scripts. New consumers should use /healthz.
	engine.GET("/health", HealthZ)

	return Built{Handler: engine, stop: limiter.Stop}
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
