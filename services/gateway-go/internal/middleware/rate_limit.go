package middleware

import (
	"net/http"
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

// stop terminates the background sweeper. Exposed for tests; production code
// runs the limiter for the lifetime of the process and lets it exit with the
// program.
func (rl *rateLimiter) stop() {
	rl.stopOnce.Do(func() { close(rl.stopCh) })
}

// RateLimit returns a Gin middleware that enforces a per-key token-bucket
// limit. The key defaults to gin's resolved client ip; override KeyFunc to
// rate-limit by user once auth is wired in.
//
// When RPS or Burst is non-positive the middleware is a no-op — useful for
// tests that need a router without the limiter active.
func RateLimit(cfg RateLimitConfig) gin.HandlerFunc {
	if cfg.RPS <= 0 || cfg.Burst <= 0 {
		return func(c *gin.Context) { c.Next() }
	}

	keyFn := cfg.KeyFunc
	if keyFn == nil {
		keyFn = func(c *gin.Context) string { return c.ClientIP() }
	}

	rl := newRateLimiter(cfg)

	return func(c *gin.Context) {
		key := keyFn(c)
		if key == "" {
			// Unable to resolve a key — fail open rather than punish callers.
			c.Next()
			return
		}
		if !rl.allow(key) {
			c.Header("Content-Type", "application/json; charset=utf-8")
			c.Header("Retry-After", "1")
			c.String(http.StatusTooManyRequests, rateLimitBody)
			c.Abort()
			return
		}
		c.Next()
	}
}
