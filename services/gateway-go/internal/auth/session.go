// Package auth implements the SIWE (Sign-In With Ethereum, EIP-4361)
// authentication flow for the gateway: nonce issuance, signature verification,
// JWT minting, Redis-backed session bookkeeping, and a Gin middleware that
// enforces an active session on protected routes.
//
// The package is intentionally split so the moving parts can be unit tested
// in isolation:
//
//   - session.go    — SessionStore interface + Redis impl + nonce store
//   - jwt.go        — JWT signing + parsing + Claims shape
//   - siwe.go       — HTTP handlers for nonce / verify / logout
//   - middleware.go — RequireAuth Gin middleware
//
// Tests use the in-memory miniredis server so they exercise the same Redis
// command surface (GETDEL, SET NX EX, DEL) as production.
package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Sentinel errors returned by the SessionStore. Handlers map these to canonical
// HTTP error codes; we keep the set small so callers cannot accidentally branch
// on a stringly-typed message.
var (
	// ErrNonceNotFound is returned when ConsumeNonce is called with a value
	// that is absent from the store. Either the nonce expired, was already
	// consumed (single-use), or was never issued. Callers must treat this as
	// an authentication failure and surface a generic "invalid nonce"
	// response without distinguishing the underlying cause.
	ErrNonceNotFound = errors.New("auth: nonce not found")

	// ErrSessionNotFound is returned when SessionExists is asked about a JTI
	// that has been logged out, expired, or never existed. RequireAuth must
	// reject the request when this is observed.
	ErrSessionNotFound = errors.New("auth: session not found")
)

// SessionStore captures the persistence surface required by the SIWE flow.
//
// The interface is deliberately narrow so tests can swap in a fake without
// pulling miniredis, and so future swap-outs (e.g. a clustered KV) can be
// dropped in without touching the handler code.
//
// All methods MUST be safe for concurrent use.
type SessionStore interface {
	// PutNonce persists the supplied nonce with the given TTL. Implementations
	// MUST refuse to overwrite an existing nonce — this is the cryptographic
	// guarantee that a nonce is single-use even under racey clients. We use
	// SET NX EX in the Redis impl; the in-memory test fake mirrors that.
	PutNonce(ctx context.Context, nonce string, ttl time.Duration) error

	// ConsumeNonce atomically reads-and-deletes the nonce. Returns
	// ErrNonceNotFound when the nonce is absent. The atomic GETDEL semantic is
	// what enforces single-use; a non-atomic GET-then-DEL would race two
	// concurrent verifies through a single nonce.
	ConsumeNonce(ctx context.Context, nonce string) error

	// PutSession records that the session keyed by jti is active for the
	// given owner address, with the supplied TTL. We store the address (not
	// the full JWT) so logout/list operations work without re-parsing the
	// token, and so log-poisoning by replaying an expired token is impossible
	// once Redis has expired the entry.
	PutSession(ctx context.Context, jti, address string, ttl time.Duration) error

	// SessionExists reports whether a session with the given jti is still
	// active. Returns ErrSessionNotFound when the session has been logged out
	// or expired. The owner address is returned so RequireAuth can re-bind it
	// onto the gin.Context, defensively re-deriving identity from the
	// trusted server-side store rather than trusting the JWT subject blindly.
	SessionExists(ctx context.Context, jti string) (string, error)

	// DeleteSession invalidates the session for the given jti. Idempotent —
	// deleting a missing session is not an error.
	DeleteSession(ctx context.Context, jti string) error
}

// RandomNonce returns a cryptographically random hex-encoded nonce of the
// supplied byte length. SIWE recommends >= 8 alphanumeric chars; we default to
// 16 bytes (32 hex chars / ~128 bits) which makes collisions and brute force
// uninteresting.
//
// nBytes must be > 0; callers passing a non-positive value get an error rather
// than a silent fallback so tests can't accidentally generate zero-entropy
// nonces.
func RandomNonce(nBytes int) (string, error) {
	if nBytes <= 0 {
		return "", fmt.Errorf("auth: nonce byte length must be > 0, got %d", nBytes)
	}
	buf := make([]byte, nBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("auth: read random: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

// redisStore is the production SessionStore backed by go-redis.
type redisStore struct {
	c *redis.Client
}

// NewRedisStore wraps an existing *redis.Client. We do not own the client's
// lifecycle — callers should Close() it on shutdown.
func NewRedisStore(c *redis.Client) SessionStore {
	return &redisStore{c: c}
}

const (
	// nonceKeyPrefix namespaces nonce keys so they cannot collide with
	// session keys even if both happen to be the same hex string.
	nonceKeyPrefix = "siwe:nonce:"
	// sessionKeyPrefix namespaces session keys.
	sessionKeyPrefix = "siwe:session:"
)

func nonceKey(n string) string   { return nonceKeyPrefix + n }
func sessionKey(j string) string { return sessionKeyPrefix + j }

func (s *redisStore) PutNonce(ctx context.Context, nonce string, ttl time.Duration) error {
	if nonce == "" {
		return errors.New("auth: empty nonce")
	}
	if ttl <= 0 {
		return fmt.Errorf("auth: nonce ttl must be > 0, got %s", ttl)
	}
	// SET NX guarantees we never overwrite an existing nonce. In practice
	// RandomNonce makes a collision astronomically unlikely, but the NX
	// turns "astronomically unlikely" into "impossible".
	ok, err := s.c.SetNX(ctx, nonceKey(nonce), "1", ttl).Result()
	if err != nil {
		return fmt.Errorf("auth: redis SETNX nonce: %w", err)
	}
	if !ok {
		return errors.New("auth: nonce collision")
	}
	return nil
}

func (s *redisStore) ConsumeNonce(ctx context.Context, nonce string) error {
	if nonce == "" {
		return ErrNonceNotFound
	}
	// GETDEL is atomic; this is what makes the nonce single-use under
	// concurrent verifies. miniredis supports it as of v2.30.x.
	res, err := s.c.GetDel(ctx, nonceKey(nonce)).Result()
	if errors.Is(err, redis.Nil) {
		return ErrNonceNotFound
	}
	if err != nil {
		return fmt.Errorf("auth: redis GETDEL nonce: %w", err)
	}
	if res == "" {
		return ErrNonceNotFound
	}
	return nil
}

func (s *redisStore) PutSession(ctx context.Context, jti, address string, ttl time.Duration) error {
	if jti == "" || address == "" {
		return errors.New("auth: empty jti or address")
	}
	if ttl <= 0 {
		return fmt.Errorf("auth: session ttl must be > 0, got %s", ttl)
	}
	if err := s.c.Set(ctx, sessionKey(jti), address, ttl).Err(); err != nil {
		return fmt.Errorf("auth: redis SET session: %w", err)
	}
	return nil
}

func (s *redisStore) SessionExists(ctx context.Context, jti string) (string, error) {
	if jti == "" {
		return "", ErrSessionNotFound
	}
	addr, err := s.c.Get(ctx, sessionKey(jti)).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrSessionNotFound
	}
	if err != nil {
		return "", fmt.Errorf("auth: redis GET session: %w", err)
	}
	return addr, nil
}

func (s *redisStore) DeleteSession(ctx context.Context, jti string) error {
	if jti == "" {
		return nil
	}
	if err := s.c.Del(ctx, sessionKey(jti)).Err(); err != nil {
		return fmt.Errorf("auth: redis DEL session: %w", err)
	}
	return nil
}
