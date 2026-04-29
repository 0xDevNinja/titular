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
//	GATEWAY_DATABASE_URL          read-only Postgres DSN
//	                              (postgres://user:pw@host:port/db). When
//	                              set, the gateway opens a pgx pool and
//	                              mounts the M4 read-only REST endpoints
//	                              under /api/v1. When unset, /api/v1 is
//	                              not mounted.
//	GATEWAY_NATS_URL              NATS connection URL
//	                              (nats://host:port). When set together
//	                              with GATEWAY_DATABASE_URL, the gateway
//	                              mounts the M4 GraphQL surface at
//	                              /graphql with NATS-backed
//	                              subscriptions; when unset, GraphQL
//	                              subscriptions resolve against an
//	                              immediately-closed channel (queries
//	                              still work). Also enables the SSE
//	                              multiplexer at /events; without this
//	                              env var the SSE surface is not
//	                              mounted at all.
//	GATEWAY_GRAPHQL_PLAYGROUND    "true" enables GET /graphql/playground
//	                              (GraphiQL UI). Default false; do NOT
//	                              enable in production-facing
//	                              deployments.
//	GATEWAY_GRAPHQL_MAX_DEPTH     positive integer; caps the nesting
//	                              depth of incoming GraphQL operations
//	                              before they reach the executor.
//	                              Default 10; introspection queries
//	                              (__schema / __type) bypass the cap.
//	GATEWAY_SSE_MAX_PER_IP        non-negative integer; caps the number
//	                              of concurrent /events connections a
//	                              single client IP may hold. Default 10;
//	                              0 falls back to the default; setting
//	                              to a negative number disables the
//	                              cap (intended for trusted networks
//	                              and tests only).
//	GATEWAY_SWAGGER_UI            "true" mounts the Swagger UI at
//	                              /swagger/ alongside the spec at
//	                              /swagger/doc.json (and /swagger/doc.yaml).
//	                              Default false; do NOT enable in
//	                              production-facing deployments — the UI
//	                              lets anyone browse the schema and
//	                              issue arbitrary requests against the
//	                              read-only surface.
//	GATEWAY_METRICS_ADDR          listen address for the Prometheus
//	                              /metrics endpoint (e.g. ":9090"). When
//	                              unset, the Prometheus surface is not
//	                              mounted and the gateway emits metrics
//	                              only via the OTLP push pipeline. The
//	                              metrics listener is intentionally
//	                              separate from GATEWAY_ADDR so an
//	                              internal Prometheus scraper reaches
//	                              /metrics without going through the
//	                              SIWE auth wall, the rate limiter or
//	                              CORS — and so the public load
//	                              balancer can keep the metrics port
//	                              unexposed.
//
// @title           Titular Gateway API
// @version         0.1.0-alpha
// @description     Read-only REST surface plus SIWE / SIWS authentication for the
// @description     Titular protocol. The gateway proxies a Postgres-backed read
// @description     model populated by the indexer track and offers a parallel
// @description     GraphQL surface (see /graphql) and SSE event stream
// @description     (see /events). All amount and id fields that originate as
// @description     uint256 on chain are encoded as decimal strings on the wire.
// @description     Authentication uses bearer tokens minted by the SIWE / SIWS
// @description     verify endpoints; the read endpoints are public by default.
// @termsOfService  https://github.com/0xDevNinja/titular
//
// @contact.name    Titular maintainers
// @contact.url     https://github.com/0xDevNinja/titular
//
// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT
//
// @host      localhost:8080
// @BasePath  /
// @schemes   http https
//
// @securityDefinitions.apikey BearerAuth
// @in    header
// @name  Authorization
// @description  Bearer token issued by POST /auth/siwe/verify or /auth/siws/verify.
//
// @tag.name  health
// @tag.description  Liveness probes for orchestration.
//
// @tag.name  auth
// @tag.description  SIWE (EVM) and SIWS (Solana) wallet authentication.
//
// @tag.name  agents
// @tag.description  Indexed agent registry (launchpad + ACP).
//
// @tag.name  trades
// @tag.description  Bonding-curve trade history.
//
// @tag.name  jobs
// @tag.description  ACP job lifecycle records.
//
// @tag.name  stats
// @tag.description  Aggregate protocol counters.
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
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/auth"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/graph"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/handlers"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/middleware"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/observability"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/openapi"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/router"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/sse"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/store"

	gatewaydocs "github.com/0xDevNinja/titular/services/gateway-go/docs"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	gin.SetMode(gin.ReleaseMode)

	addr := envOr("GATEWAY_ADDR", ":8080")
	service := envOr("GATEWAY_SERVICE", "gateway")
	metricsAddr := strings.TrimSpace(os.Getenv("GATEWAY_METRICS_ADDR"))

	// Prometheus reader is wired BEFORE OTel.Init so the same
	// MeterProvider feeds both pipelines. AttachPrometheusReader is a
	// no-op when GATEWAY_METRICS_ADDR is unset — we only build the
	// reader when the operator has asked for the surface. The
	// resulting reader, when present, lands on the SDK MeterProvider
	// that Init constructs a few lines below.
	var promShutdown func() error
	if metricsAddr != "" {
		if _, _, perr := observability.AttachPrometheusReader(); perr != nil {
			// Soft-fail: losing /metrics is bad but not fatal —
			// gateway request handling is independent. Log loud
			// enough that the on-call notices.
			log.Warn().Err(perr).Msg("prometheus exporter init failed; /metrics disabled")
			metricsAddr = ""
		}
	}

	// OTel must come up BEFORE any subsystem that may emit a span
	// (the pgx tracer is wired the moment the pool is built; the gin
	// middleware reads otel.GetTracerProvider() per-request). Init is
	// a no-op when OTEL_EXPORTER_OTLP_ENDPOINT is unset, so a gateway
	// without observability still starts cleanly.
	//
	// Soft-fail posture: a flaky collector during deploy must NOT
	// wedge gateway startup. We log the error and continue with the
	// SDK left at no-op; operators who explicitly want hard-fail (e.g.
	// CI runs that should refuse to ship unobserved) can set
	// OTEL_FAIL_ON_INIT=1.
	otelShutdown, err := observability.Init(context.Background(), observability.Config{
		ServiceName: service,
	})
	if err != nil {
		if os.Getenv("OTEL_FAIL_ON_INIT") == "1" {
			log.Fatal().Err(err).Msg("failed to initialise observability")
		}
		log.Warn().Err(err).Msg("otel init failed (continuing)")
		otelShutdown = func(context.Context) error { return nil }
	}

	// If Init was skipped (no OTLP endpoint) but the operator still
	// wanted /metrics, fall back to PrometheusHandler which builds a
	// minimal Prom-only MeterProvider. PrometheusHandler is also
	// idempotent against AttachPrometheusReader, so this same call is
	// safe in the OTLP-enabled path.
	var metricsHandler http.Handler
	if metricsAddr != "" {
		h, sd, herr := observability.PrometheusHandler()
		if herr != nil {
			log.Warn().Err(herr).Msg("prometheus handler init failed; /metrics disabled")
			metricsAddr = ""
		} else {
			metricsHandler = h
			promShutdown = sd
		}
	}

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

	apiHandlers, dbClose, err := buildAPIHandlers(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("invalid database configuration")
	}

	// Single shared NATS connection for both the GraphQL subscription
	// bus (#89) and the SSE multiplexer (#90). One connection avoids
	// duplicate auth handshakes and keeps the connection pool size
	// predictable in shared NATS clusters; the indexer publisher
	// dedup window means consumers can use the same wire stream
	// without coordinating.
	natsConn, natsClose, err := buildNATSConn()
	if err != nil {
		log.Fatal().Err(err).Msg("invalid NATS configuration")
	}

	gqlHandler, err := buildGraphQLHandler(apiHandlers, corsCfg, natsConn)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid graphql configuration")
	}

	sseHandler, sseStop, err := buildSSEHandler(natsConn)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid SSE configuration")
	}

	openAPISpec, err := buildOpenAPISpec()
	if err != nil {
		log.Fatal().Err(err).Msg("invalid OpenAPI configuration")
	}

	cfg := router.Config{
		Logger:         log.Logger,
		Service:        service,
		CORS:           corsCfg,
		RateLimit:      buildRateLimitConfig(),
		TrustedProxies: parseList(os.Getenv("GATEWAY_TRUSTED_PROXIES")),
		Auth:           authHandlers,
		SIWS:           siwsHandlers,
		API:            apiHandlers,
		GraphQL:        gqlHandler,
		SSE:            sseHandler,
		OpenAPI:        openAPISpec,
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

	// Dedicated /metrics listener on GATEWAY_METRICS_ADDR. We use a
	// separate http.Server instead of mounting the handler on the
	// main gin engine because:
	//
	//   - The main engine is fronted by the SIWE auth wall, the
	//     rate limiter, and CORS — none of which a Prometheus
	//     scraper needs (or wants) to traverse.
	//   - The metrics port is typically NOT exposed on the public
	//     load balancer; running it on a different listener makes
	//     that operational separation a hard, file-level setting
	//     instead of a per-route middleware skip.
	//   - A misbehaving handler in the main engine cannot starve
	//     scrapes (and vice versa) because the listeners run on
	//     independent goroutines.
	//
	// Timeouts are tighter than the main server: a Prometheus scrape
	// is a single GET that completes in milliseconds. Read/Write 5s
	// is generous and still well below the default Prom 10s scrape
	// timeout, so a stuck handler reports as a Prom timeout (with
	// metric prom_target_metadata_cache_entries… surfacing the
	// failure) rather than as a gateway resource leak.
	var metricsSrv *http.Server
	if metricsAddr != "" && metricsHandler != nil {
		mux := http.NewServeMux()
		mux.Handle("/metrics", metricsHandler)
		metricsSrv = &http.Server{
			Addr:         metricsAddr,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  30 * time.Second,
		}
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

	if metricsSrv != nil {
		go func() {
			log.Info().Str("addr", metricsAddr).Msg("gateway /metrics listening")
			if err := metricsSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				// Metrics listener crash is not fatal: log and let the
				// main server keep running. An on-call sees the
				// missing /metrics surface via Prometheus's `up{}`
				// gauge dropping to 0.
				log.Error().Err(err).Msg("metrics server error")
			}
		}()
	}

	<-quit
	log.Info().Msg("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("graceful shutdown failed")
	}

	// Tear the metrics listener down on the same shutdown deadline.
	// Doing this AFTER the main server's Shutdown so any final
	// in-flight handlers can still record samples that are then
	// served on the very last scrape.
	if metricsSrv != nil {
		if err := metricsSrv.Shutdown(ctx); err != nil {
			log.Warn().Err(err).Msg("metrics shutdown")
		}
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

	// Close the Postgres pool we own. Same ownership rule as Redis above.
	if dbClose != nil {
		dbClose()
	}

	// Stop the SSE multiplexer BEFORE closing the NATS connection so
	// the shutdown order matches the reverse of construction: the
	// multiplexer's NATS subscriptions need a live connection to
	// unsubscribe cleanly.
	if sseStop != nil {
		sseStop()
	}

	// Close the shared NATS connection. Must run after sseStop above
	// AND after srv.Shutdown — both the GraphQL subscription bus and
	// the SSE multiplexer can issue final NATS calls during their
	// teardown.
	if natsClose != nil {
		natsClose()
	}

	// OTel shutdown last so any final spans / metrics emitted by the
	// teardown above are flushed to the collector. The 10s internal
	// cap inside the shutdown closure means a wedged collector cannot
	// pin the process indefinitely on SIGTERM.
	if otelShutdown != nil {
		if err := otelShutdown(context.Background()); err != nil {
			log.Warn().Err(err).Msg("otel shutdown")
		}
	}
	// Prometheus reader teardown runs after the OTel meter provider
	// shutdown so any last metric flushed during otelShutdown is
	// scrapable on the closing /metrics request. Safe to call when
	// the surface was disabled — promShutdown is nil in that case.
	if promShutdown != nil {
		if err := promShutdown(); err != nil {
			log.Warn().Err(err).Msg("prometheus shutdown")
		}
	}
}

// buildNATSConn opens the single shared NATS connection consumed by
// both the GraphQL subscription bus (#89) and the SSE multiplexer
// (#90). Returns (nil, nil, nil) when GATEWAY_NATS_URL is unset — the
// downstream feature builders treat a nil connection as "feature off"
// and mount their handlers in nilBus / no-SSE form so the gateway
// still serves queries.
//
// The returned close function tears the connection down; callers MUST
// invoke it on graceful shutdown AFTER stopping every subsystem that
// holds an active subscription on the connection.
func buildNATSConn() (*nats.Conn, func(), error) {
	url := strings.TrimSpace(os.Getenv("GATEWAY_NATS_URL"))
	if url == "" {
		return nil, nil, nil
	}
	nc, err := nats.Connect(url, nats.Timeout(5*time.Second))
	if err != nil {
		return nil, nil, errors.New("GATEWAY_NATS_URL: " + err.Error())
	}
	return nc, nc.Close, nil
}

// buildGraphQLHandler wires the GraphQL surface introduced in M4 (#89).
// It is a no-op (returns nil, nil) when the gateway is launched
// without GATEWAY_DATABASE_URL — Query and Subscription resolvers both
// need a Store to be useful, and silently mounting a query-less endpoint
// would mask deployment misconfiguration.
//
// natsConn is optional: when nil, subscription resolvers fall back to
// graph.NilBus() so subscriber clients see a clean end-of-stream
// rather than a hang. Queries are unaffected.
//
// GATEWAY_GRAPHQL_PLAYGROUND defaults to off; setting it to "true" mounts
// the GET /graphql/playground GraphiQL UI for local development. Do
// not enable in production-facing deployments — the UI lets anyone
// browse the schema and execute arbitrary queries against the
// read-only surface.
func buildGraphQLHandler(api *handlers.API, cors middleware.CORSConfig, natsConn *nats.Conn) (*graph.Handler, error) {
	if api == nil {
		// No store -> no GraphQL surface. Returning nil is the
		// "feature off" signal the router config recognises.
		return nil, nil
	}

	resolver := &graph.Resolver{Store: api.Store}

	if natsConn != nil {
		bus, err := graph.NewNATSBus(natsConn)
		if err != nil {
			return nil, err
		}
		resolver.Bus = bus
	}

	schema, err := graph.NewSchema(resolver)
	if err != nil {
		return nil, err
	}
	h, err := graph.NewHandler(schema, resolver)
	if err != nil {
		return nil, err
	}
	if v := strings.TrimSpace(os.Getenv("GATEWAY_GRAPHQL_PLAYGROUND")); strings.EqualFold(v, "true") {
		h.EnablePlayground = true
	}

	// Origin enforcement on the websocket transport. Browsers DO NOT
	// apply CORS preflight to WebSocket upgrades, so the gateway's CORS
	// middleware that protects POST /graphql gives this path no
	// protection. Re-use the same allowed-origins list (and credential
	// posture) so operators only have one knob to tune.
	h.AllowedOrigins = cors.AllowedOrigins
	h.AllowCredentials = cors.AllowCredentials
	h.AllowReflectedOrigins = cors.AllowReflectedOrigins

	// Optional depth cap. Anything > 0 wins; otherwise the handler
	// falls back to its package default (defaultMaxQueryDepth).
	if v := strings.TrimSpace(os.Getenv("GATEWAY_GRAPHQL_MAX_DEPTH")); v != "" {
		if n, perr := strconv.Atoi(v); perr == nil && n > 0 {
			h.MaxQueryDepth = n
		}
	}

	return h, nil
}

// buildSSEHandler wires the SSE multiplexer introduced in M4 (#90).
// Returns (nil, nil, nil) when natsConn is nil; the SSE surface
// requires a live NATS connection to fan messages out, and a NATS-less
// gateway has no events to stream.
//
// The returned stop function shuts the multiplexer's hub goroutine
// down (closes every subscriber's channel, unsubscribes from NATS).
// Callers MUST invoke it BEFORE closing the shared NATS connection so
// the unsubscribes complete on a live transport.
//
// JetStream context is opened lazily — when JetStream isn't available
// on the connection (e.g. an embedded server in tests with JS
// disabled), the handler still serves live tail and silently skips
// Last-Event-ID replay. That parallels the publisher's behaviour
// where EnsureStream is the only mutating call site.
func buildSSEHandler(natsConn *nats.Conn) (*sse.Handler, func(), error) {
	if natsConn == nil {
		return nil, nil, nil
	}
	mux, err := sse.NewMultiplexer(sse.Config{
		NATS:   natsConn,
		Logger: log.Logger,
	})
	if err != nil {
		return nil, nil, err
	}
	if err := mux.Start(); err != nil {
		return nil, nil, err
	}

	// JetStream is best-effort. JetStream() returns an error only when
	// the server has rejected the request (e.g. JS not enabled on the
	// connection's permissions); in that case we serve live tail
	// without replay rather than failing startup.
	js, jsErr := natsConn.JetStream()
	if jsErr != nil {
		log.Warn().Err(jsErr).Msg("sse: jetstream unavailable, replay disabled")
		js = nil
	}

	// MaxPerIP: 0 in HandlerConfig means "use the package default";
	// strconv.Atoi returns 0 for an unset env var, which is exactly the
	// fallthrough we want. Operators wanting to disable the cap pass a
	// negative integer explicitly.
	maxPerIP, _ := strconv.Atoi(strings.TrimSpace(os.Getenv("GATEWAY_SSE_MAX_PER_IP")))

	h, err := sse.NewHandler(sse.HandlerConfig{
		Multiplexer: mux,
		JetStream:   js,
		Logger:      log.Logger,
		MaxPerIP:    maxPerIP,
	})
	if err != nil {
		mux.Stop()
		return nil, nil, err
	}
	return h, mux.Stop, nil
}

// buildAPIHandlers opens a Postgres pool from GATEWAY_DATABASE_URL and wires
// it into the read-only REST API. Returns (nil, nil, nil) when the env var is
// unset — the M4 endpoints are simply not mounted in that case, which keeps
// the binary usable in dev environments without a database (M2/M3 fixture
// surface still works).
func buildAPIHandlers(ctx context.Context) (*handlers.API, func(), error) {
	dsn := strings.TrimSpace(os.Getenv("GATEWAY_DATABASE_URL"))
	if dsn == "" {
		return nil, nil, nil
	}
	pg, closeFn, err := store.Connect(ctx, dsn)
	if err != nil {
		return nil, nil, err
	}
	return handlers.NewAPI(pg), closeFn, nil
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

// buildOpenAPISpec returns the OpenAPI handler bundle, populated from the
// embedded swag-generated spec under services/gateway-go/docs/. The raw
// /swagger/doc.json + /swagger/doc.yaml endpoints are mounted
// unconditionally; the Swagger UI viewer at /swagger/ is gated by the
// GATEWAY_SWAGGER_UI env var (default off).
//
// The spec is embedded at build time so the running gateway always
// serves the document that matches its own source — no filesystem
// lookups, no version skew. Regeneration is a developer step
// (`swag init` from services/gateway-go) and is enforced in CI by the
// gateway-openapi-check drift gate.
func buildOpenAPISpec() (*openapi.Spec, error) {
	enableUI := openapi.EnabledFromEnv(os.Getenv("GATEWAY_SWAGGER_UI"))
	return openapi.NewSpec(gatewaydocs.SpecJSON, gatewaydocs.SpecYAML, enableUI)
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
