// Command gateway is the HTTP entry point for the Titular API gateway.
//
// Behaviour is controlled entirely via environment variables so the binary is
// suitable for direct deployment in containerised environments without a
// config file.
//
//	GATEWAY_ADDR              listen address; default ":8080"
//	GATEWAY_SERVICE           service tag emitted in structured logs; default "gateway"
//	GATEWAY_CORS_ORIGINS      comma-separated list of permitted CORS origins;
//	                          empty disables CORS, "*" enables wildcard
//	GATEWAY_CORS_CREDENTIALS  "true" to advertise Access-Control-Allow-Credentials
//	GATEWAY_RATE_LIMIT_RPS    per-ip token-bucket refill rate; 0 disables limiting
//	GATEWAY_RATE_LIMIT_BURST  per-ip burst capacity; 0 disables limiting
//	GATEWAY_TRUSTED_PROXIES   comma-separated CIDRs forwarded to gin.SetTrustedProxies
package main

import (
	"context"
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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

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

	cfg := router.Config{
		Logger:         log.Logger,
		Service:        service,
		CORS:           buildCORSConfig(),
		RateLimit:      buildRateLimitConfig(),
		TrustedProxies: parseList(os.Getenv("GATEWAY_TRUSTED_PROXIES")),
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      router.NewWithConfig(cfg, agentHandlers, jobHandlers),
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
