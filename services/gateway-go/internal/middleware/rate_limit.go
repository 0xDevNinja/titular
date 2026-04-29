package middleware

import (
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimitConfig configures the per-client token-bucket rate limiter.
//
// RPS is the steady-state rate in requests per second. Burst is the maximum
// number of requests the bucket can absorb in a short spike. Both must be > 0
// to enable limiting; if either is non-positive the middleware returned by
// RateLimit becomes a no-op (used to disable limiting in tests/dev).
//
// IdleTTL is how long a client's bucket is retained after its last use before
// being eligible for cleanup. Defaults to 10 minutes when zero.
//
// KeyFunc resolves the bucket key from the incoming context. Defaults to the
// gin-resolved client ip, which is the most common shape; callers can override
// it to key on auth subject once SIWE auth lands in #86.
type RateLimitConfig struct {
	RPS     float64
	Burst   int
	IdleTTL time.Duration
	KeyFunc func(*gin.Context) string
}

// rateLimitBody is the public 429 envelope. We never include the bucket key,
// remaining tokens, or any other internal value that could leak the keying
// strategy.
const rateLimitBody = `{"error":{"code":"rate_limited","message":"too many requests"}}`

// rateLimiter holds the per-key buckets together with their last-seen time so
// idle buckets can be evicted by a background sweeper.
type rateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucketEntry
	rps      rate.Limit
	burst    int
	idleTTL  time.Duration
	now      func() time.Time // overridable in tests
	stopOnce sync.Once
	stopCh   chan struct{}
}

type bucketEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func newRateLimiter(cfg RateLimitConfig) *rateLimiter {
	idle := cfg.IdleTTL
	if idle <= 0 {
		idle = 10 * time.Minute
	}
	rl := &rateLimiter{
		buckets: make(map[string]*bucketEntry),
		rps:     rate.Limit(cfg.RPS),
		burst:   cfg.Burst,
		idleTTL: idle,
		now:     time.Now,
		stopCh:  make(chan struct{}),
	}
	go rl.sweepLoop()
	return rl
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	entry, ok := rl.buckets[key]
	if !ok {
		entry = &bucketEntry{limiter: rate.NewLimiter(rl.rps, rl.burst)}
		rl.buckets[key] = entry
	}
	entry.lastSeen = rl.now()
	rl.mu.Unlock()

	return entry.limiter.Allow()
}

// sweepLoop periodically drops buckets that have not been used within idleTTL.
// The interval is half of idleTTL with a floor of one minute to avoid waking
// the goroutine more than once per minute under aggressive configuration.
func (rl *rateLimiter) sweepLoop() {
	interval := rl.idleTTL / 2
	if interval < time.Minute {
		interval = time.Minute
	}
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-rl.stopCh:
			return
		case <-t.C:
			rl.sweep()
		}
	}
}

func (rl *rateLimiter) sweep() {
	cutoff := rl.now().Add(-rl.idleTTL)
	rl.mu.Lock()
	defer rl.mu.Unlock()
	for k, entry := range rl.buckets {
		if entry.lastSeen.Before(cutoff) {
			delete(rl.buckets, k)
		}
	}
}

// Stop terminates the background sweeper goroutine. Safe to call multiple
// times. Wire it into the server's graceful-shutdown path so the goroutine
// exits cleanly rather than leaking until process exit.
func (rl *rateLimiter) Stop() {
	rl.stopOnce.Do(func() { close(rl.stopCh) })
}

// RateLimiter bundles the gin handler with a Stop() hook so callers can shut
// the background sweeper down on graceful shutdown. The zero value is unsafe;
// always construct via RateLimit.
type RateLimiter struct {
	Handler gin.HandlerFunc
	stop    func()
}

// Stop tears down any background goroutines owned by the limiter. Safe to
// call multiple times and on a no-op limiter (RPS/Burst <= 0).
func (r RateLimiter) Stop() {
	if r.stop != nil {
		r.stop()
	}
}

// RateLimit returns a Gin middleware that enforces a per-key token-bucket
// limit. The key defaults to gin's resolved client ip; override KeyFunc to
// rate-limit by user once auth is wired in.
//
// When RPS or Burst is non-positive the middleware is a no-op — useful for
// tests that need a router without the limiter active.
//
// RateLimit DISCARDS the lifecycle Stop hook; callers that care about clean
// shutdown should use NewRateLimiter and call Stop() explicitly. Today only
// cmd/gateway needs this.
func RateLimit(cfg RateLimitConfig) gin.HandlerFunc {
	return NewRateLimiter(cfg).Handler
}

// NewRateLimiter is the lifecycle-aware constructor. Use this when the caller
// wants to call Stop() on graceful shutdown; otherwise prefer RateLimit.
func NewRateLimiter(cfg RateLimitConfig) RateLimiter {
	if cfg.RPS <= 0 || cfg.Burst <= 0 {
		return RateLimiter{Handler: func(c *gin.Context) { c.Next() }}
	}

	keyFn := cfg.KeyFunc
	if keyFn == nil {
		keyFn = func(c *gin.Context) string { return c.ClientIP() }
	}

	rl := newRateLimiter(cfg)

	// Retry-After is advertised in seconds. We round up the bucket refill
	// interval (1/RPS) to the nearest whole second, with a floor of 1 so we
	// never advertise an impossibly aggressive retry budget. Computed once at
	// middleware-construction time.
	retryAfter := strconv.Itoa(retryAfterSeconds(cfg.RPS))

	handler := func(c *gin.Context) {
		key := keyFn(c)
		if key == "" {
			// Unable to resolve a key — fail open rather than punish callers.
			c.Next()
			return
		}
		if !rl.allow(key) {
			c.Header("Content-Type", "application/json; charset=utf-8")
			c.Header("Retry-After", retryAfter)
			c.String(http.StatusTooManyRequests, rateLimitBody)
			c.Abort()
			return
		}
		c.Next()
	}
	return RateLimiter{Handler: handler, stop: rl.Stop}
}

// retryAfterSeconds returns max(1, ceil(1/rps)). Inputs of <= 0 fall back to
// 1 — defence in depth; the caller already gates on RPS > 0.
func retryAfterSeconds(rps float64) int {
	if rps <= 0 {
		return 1
	}
	v := int(math.Ceil(1.0 / rps))
	if v < 1 {
		return 1
	}
	return v
}
