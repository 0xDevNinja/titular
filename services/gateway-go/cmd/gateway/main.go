// Command gateway is the HTTP entry point for the Titular API gateway.
//
// Behaviour is controlled entirely via environment variables so the binary is
// suitable for direct deployment in containerised environments without a
// config file.
//
//	GATEWAY_ADDR                  listen address; default ":8080"
//	GATEWAY_SERVICE               service tag emitted in structured logs; default "gateway"
//	GATEWAY_CORS_ORIGINS          comma-separated list of permitted CORS origins;
//	                              empty disables CORS, "*" enables wildcard
//	GATEWAY_CORS_CREDENTIALS      "true" to advertise Access-Control-Allow-Credentials
//	GATEWAY_CORS_REFLECT_ORIGINS  "true" to allow the unsafe wildcard+credentials
//	                              combination (echoes the request Origin); leave
//	                              unset in production
//	GATEWAY_RATE_LIMIT_RPS        per-ip token-bucket refill rate; 0 disables limiting
//	GATEWAY_RATE_LIMIT_BURST      per-ip burst capacity; 0 disables limiting
//	GATEWAY_TRUSTED_PROXIES       comma-separated CIDRs forwarded to
//	                              gin.SetTrustedProxies. SECURITY: leaving this
//	                              UNSET trusts NO proxies (unspoofable client
//	                              IP). Operators behind a known L7 proxy MUST
//	                              set this explicitly.
//	GATEWAY_JWT_SECRET            base64-encoded HMAC key (>=32 bytes after
//	                              decode). REQUIRED to enable SIWE auth; when
//	                              unset the /auth endpoints are not mounted.
//	                              Short or non-base64 values fail startup.
//	GATEWAY_JWT_ISSUER            iss claim minted into and required on every
//	                              JWT. Defaults to "titular-gateway".
//	GATEWAY_JWT_TTL               session lifetime as a Go duration (e.g.
//	                              "24h"). Defaults to 24h.
//	GATEWAY_REDIS_URL             redis connection URL (redis://host:port/db).
//	                              REQUIRED when JWT_SECRET is set.
//	GATEWAY_SIWE_DOMAIN           the value the SIWE message MUST declare in
//	                              its `domain` line. Pinned server-side to
//	                              prevent cross-site replay.
//	GATEWAY_SIWE_CHAIN_ID         the chain id the SIWE message MUST declare.
//	                              Defaults to 84532 (Base Sepolia).
//	GATEWAY_SIWS_DOMAIN           the value the SIWS (Sign-In With Solana)
//	                              message MUST declare in its header line.
//	                              Pinned server-side to prevent cross-site
//	                              replay. Falls back to GATEWAY_SIWE_DOMAIN
//	                              when unset so single-domain deployments
//	                              don't need to configure twice.
//	GATEWAY_SIWS_CLUSTER          the Solana cluster the SIWS message MUST
//	                              declare. One of "devnet", "mainnet-beta",
//	                              "testnet". Defaults to "devnet". Setting
//	                              this enables the /auth/siws/* endpoints
//	                              (provided JWT_SECRET + REDIS_URL are also
//	                              set).
package main

import (
	"context"
	"encoding/base64"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/auth"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/handlers"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/middleware"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/router"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	gin.SetMode(gin.ReleaseMode)

	addr := envOr("GATEWAY_ADDR", ":8080")
	service := envOr("GATEWAY_SERVICE", "gateway")

	agentHandlers, err := handlers.NewAgentHandlers()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialise agent handlers")
	}
	jobHandlers, err := handlers.NewJobHandlers()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialise job handlers")
	}

	corsCfg := buildCORSConfig()
	if err := corsCfg.Validate(); err != nil {
		log.Fatal().Err(err).Msg("invalid CORS configuration")
	}

	authHandlers, redisClient, err := buildAuthHandlers(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("invalid auth configuration")
	}
	siwsHandlers, err := buildSIWSHandlers(authHandlers, redisClient)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid SIWS configuration")
	}

	cfg := router.Config{
		Logger:         log.Logger,
		Service:        service,
		CORS:           corsCfg,
		RateLimit:      buildRateLimitConfig(),
		TrustedProxies: parseList(os.Getenv("GATEWAY_TRUSTED_PROXIES")),
		Auth:           authHandlers,
		SIWS:           siwsHandlers,
	}

	built := router.NewWithConfigLifecycle(cfg, agentHandlers, jobHandlers)
	srv := &http.Server{
		Addr:         addr,
		Handler:      built.Handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
		BaseContext: func(_ net.Listener) context.Context {
			return context.Background()
		},
	}

	// Graceful shutdown on SIGINT / SIGTERM.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info().Str("addr", addr).Str("service", service).Msg("gateway listening")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	<-quit
	log.Info().Msg("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("graceful shutdown failed")
	}

	// Stop background goroutines owned by the router (rate-limit sweeper).
	built.Stop()

	// Close the Redis client we own. The router does not — see
	// buildAuthHandlers — so we are responsible for it on shutdown.
	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			log.Warn().Err(err).Msg("redis close")
		}
	}
}

// buildAuthHandlers parses the SIWE/JWT/Redis env knobs and returns either a
// fully-wired *auth.Handlers (plus the Redis client we now own) or nil to
// indicate that auth is disabled. We refuse to start when the env is partly
// configured — half-configured auth is the dangerous state.
func buildAuthHandlers(ctx context.Context) (*auth.Handlers, *redis.Client, error) {
	rawSecret := strings.TrimSpace(os.Getenv("GATEWAY_JWT_SECRET"))
	if rawSecret == "" {
		// Explicitly opt-out: SIWE not configured. Routes are not mounted.
		return nil, nil, nil
	}

	secret, err := base64.StdEncoding.DecodeString(rawSecret)
	if err != nil {
		return nil, nil, errors.New("GATEWAY_JWT_SECRET must be base64-encoded")
	}

	issuer := envOr("GATEWAY_JWT_ISSUER", "titular-gateway")

	ttl := 24 * time.Hour
	if v := strings.TrimSpace(os.Getenv("GATEWAY_JWT_TTL")); v != "" {
		parsed, perr := time.ParseDuration(v)
		if perr != nil {
			return nil, nil, errors.New("GATEWAY_JWT_TTL must be a Go duration")
		}
		ttl = parsed
	}

	signer, err := auth.NewJWTSigner(auth.SignerConfig{
		Secret: secret,
		Issuer: issuer,
		TTL:    ttl,
	})
	if err != nil {
		return nil, nil, err
	}

	redisURL := strings.TrimSpace(os.Getenv("GATEWAY_REDIS_URL"))
	if redisURL == "" {
		return nil, nil, errors.New("GATEWAY_REDIS_URL required when GATEWAY_JWT_SECRET is set")
	}
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, nil, errors.New("GATEWAY_REDIS_URL invalid")
	}
	client := redis.NewClient(opts)
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, nil, errors.New("redis ping failed")
	}

	domain := strings.TrimSpace(os.Getenv("GATEWAY_SIWE_DOMAIN"))
	if domain == "" {
		_ = client.Close()
		return nil, nil, errors.New("GATEWAY_SIWE_DOMAIN required when GATEWAY_JWT_SECRET is set")
	}

	chainID := 84532 // Base Sepolia default
	if v := strings.TrimSpace(os.Getenv("GATEWAY_SIWE_CHAIN_ID")); v != "" {
		parsed, perr := strconv.Atoi(v)
		if perr != nil || parsed <= 0 {
			_ = client.Close()
			return nil, nil, errors.New("GATEWAY_SIWE_CHAIN_ID must be a positive integer")
		}
		chainID = parsed
	}

	h, err := auth.NewHandlers(auth.HandlerConfig{
		Store:   auth.NewRedisStore(client),
		Signer:  signer,
		Domain:  domain,
		ChainID: chainID,
	})
	if err != nil {
		_ = client.Close()
		return nil, nil, err
	}
	return h, client, nil
}

// buildSIWSHandlers assembles the Sign-In-With-Solana handler bundle on top
// of the existing SIWE infrastructure. We deliberately reuse the signer
// and session store so:
//
//   - There is only one JWT secret in memory (no fan-out of HMAC keys).
//   - Logout is a single endpoint that revokes either flow's tokens.
//   - The nonce namespace is shared, which is correct: the chain binding
//     happens via the signed message, not the nonce.
//
// SIWS is enabled when SIWE is enabled AND GATEWAY_SIWS_CLUSTER is set.
// We accept GATEWAY_SIWS_DOMAIN as an explicit override for deployments
// where the SIWS-facing domain differs from the SIWE-facing domain;
// otherwise we inherit GATEWAY_SIWE_DOMAIN.
func buildSIWSHandlers(siwe *auth.Handlers, client *redis.Client) (*auth.SIWSHandlers, error) {
	if siwe == nil || client == nil {
		// SIWE not configured — SIWS cannot be either, since they share
		// the signer/store. Caller treats nil as "feature off".
		return nil, nil
	}
	cluster := strings.TrimSpace(os.Getenv("GATEWAY_SIWS_CLUSTER"))
	if cluster == "" {
		// SIWS opt-in: explicitly require the operator to declare a
		// cluster. We refuse to default to mainnet — getting that
		// silently wrong is the worst possible failure mode.
		return nil, nil
	}

	domain := strings.TrimSpace(os.Getenv("GATEWAY_SIWS_DOMAIN"))
	if domain == "" {
		domain = strings.TrimSpace(os.Getenv("GATEWAY_SIWE_DOMAIN"))
	}
	if domain == "" {
		return nil, errors.New("GATEWAY_SIWS_DOMAIN or GATEWAY_SIWE_DOMAIN required when SIWS enabled")
	}

	return auth.NewSIWSHandlers(auth.SIWSHandlerConfig{
		Store:   siwe.Store(),
		Signer:  siwe.Signer(),
		Domain:  domain,
		Cluster: cluster,
	})
}

// envOr returns the value of the named env var, falling back to def when
// unset or empty.
func envOr(name, def string) string {
	if v := strings.TrimSpace(os.Getenv(name)); v != "" {
		return v
	}
	return def
}

// parseList splits a comma-separated env string into a slice of trimmed,
// non-empty values. Returns nil when the input is empty so callers can
// distinguish "unset" from "empty list".
func parseList(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// buildCORSConfig pulls CORS settings from the environment.
func buildCORSConfig() middleware.CORSConfig {
	cfg := middleware.DefaultCORSConfig()
	cfg.AllowedOrigins = parseList(os.Getenv("GATEWAY_CORS_ORIGINS"))
	if v := os.Getenv("GATEWAY_CORS_CREDENTIALS"); strings.EqualFold(v, "true") {
		cfg.AllowCredentials = true
	}
	if v := os.Getenv("GATEWAY_CORS_REFLECT_ORIGINS"); strings.EqualFold(v, "true") {
		cfg.AllowReflectedOrigins = true
	}
	return cfg
}

// buildRateLimitConfig pulls rate-limit settings from the environment. Either
// of RPS or BURST being unset/zero disables limiting (the middleware itself
// no-ops in that case).
func buildRateLimitConfig() middleware.RateLimitConfig {
	rps, _ := strconv.ParseFloat(os.Getenv("GATEWAY_RATE_LIMIT_RPS"), 64)
	burst, _ := strconv.Atoi(os.Getenv("GATEWAY_RATE_LIMIT_BURST"))
	return middleware.RateLimitConfig{
		RPS:   rps,
		Burst: burst,
	}
}
