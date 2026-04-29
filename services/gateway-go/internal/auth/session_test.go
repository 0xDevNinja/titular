package auth_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/auth"
)

// newRedis spins up a fresh in-memory miniredis + go-redis client for the
// duration of a test. The cleanup hooks unconditionally close both so a
// failed assertion does not leak a server.
func newRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	t.Helper()
	mr := miniredis.RunT(t)
	c := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = c.Close() })
	return mr, c
}

func TestRandomNonce_LengthAndUniqueness(t *testing.T) {
	a, err := auth.RandomNonce(16)
	if err != nil {
		t.Fatalf("RandomNonce: %v", err)
	}
	b, err := auth.RandomNonce(16)
	if err != nil {
		t.Fatalf("RandomNonce: %v", err)
	}
	// 16 bytes -> 32 hex chars
	if len(a) != 32 || len(b) != 32 {
		t.Fatalf("nonce length: got %d/%d, want 32/32", len(a), len(b))
	}
	if a == b {
		t.Fatal("two RandomNonce calls produced identical output — RNG broken")
	}
}

func TestRandomNonce_RejectsNonPositiveLength(t *testing.T) {
	if _, err := auth.RandomNonce(0); err == nil {
		t.Fatal("expected error for n=0, got nil")
	}
	if _, err := auth.RandomNonce(-1); err == nil {
		t.Fatal("expected error for n=-1, got nil")
	}
}

func TestRedisStore_NonceLifecycle_ConsumeOnceOnly(t *testing.T) {
	_, c := newRedis(t)
	store := auth.NewRedisStore(c)
	ctx := context.Background()

	const nonce = "abc123"
	if err := store.PutNonce(ctx, nonce, time.Minute); err != nil {
		t.Fatalf("PutNonce: %v", err)
	}
	if err := store.ConsumeNonce(ctx, nonce); err != nil {
		t.Fatalf("ConsumeNonce first call: %v", err)
	}
	// Second consume MUST fail — that's the single-use guarantee.
	err := store.ConsumeNonce(ctx, nonce)
	if !errors.Is(err, auth.ErrNonceNotFound) {
		t.Fatalf("ConsumeNonce twice: got %v, want ErrNonceNotFound", err)
	}
}

func TestRedisStore_NonceTTL_ExpiresAndIsAbsent(t *testing.T) {
	mr, c := newRedis(t)
	store := auth.NewRedisStore(c)
	ctx := context.Background()

	if err := store.PutNonce(ctx, "n1", time.Minute); err != nil {
		t.Fatalf("PutNonce: %v", err)
	}
	mr.FastForward(2 * time.Minute)
	if err := store.ConsumeNonce(ctx, "n1"); !errors.Is(err, auth.ErrNonceNotFound) {
		t.Fatalf("expired nonce: got %v, want ErrNonceNotFound", err)
	}
}

func TestRedisStore_PutNonce_RefusesCollision(t *testing.T) {
	_, c := newRedis(t)
	store := auth.NewRedisStore(c)
	ctx := context.Background()

	if err := store.PutNonce(ctx, "dup", time.Minute); err != nil {
		t.Fatalf("first PutNonce: %v", err)
	}
	if err := store.PutNonce(ctx, "dup", time.Minute); err == nil {
		t.Fatal("expected collision error, got nil")
	}
}

func TestRedisStore_PutNonce_RejectsBadInputs(t *testing.T) {
	_, c := newRedis(t)
	store := auth.NewRedisStore(c)
	ctx := context.Background()

	if err := store.PutNonce(ctx, "", time.Minute); err == nil {
		t.Fatal("empty nonce: expected error")
	}
	if err := store.PutNonce(ctx, "x", 0); err == nil {
		t.Fatal("zero ttl: expected error")
	}
	if err := store.PutNonce(ctx, "x", -1); err == nil {
		t.Fatal("negative ttl: expected error")
	}
}

func TestRedisStore_SessionLifecycle(t *testing.T) {
	_, c := newRedis(t)
	store := auth.NewRedisStore(c)
	ctx := context.Background()

	const jti = "jti-1"
	const addr = "0xdeadbeef"
	if err := store.PutSession(ctx, jti, addr, time.Hour); err != nil {
		t.Fatalf("PutSession: %v", err)
	}
	got, err := store.SessionExists(ctx, jti)
	if err != nil {
		t.Fatalf("SessionExists: %v", err)
	}
	if got != addr {
		t.Fatalf("address: got %q, want %q", got, addr)
	}
	if err := store.DeleteSession(ctx, jti); err != nil {
		t.Fatalf("DeleteSession: %v", err)
	}
	if _, err := store.SessionExists(ctx, jti); !errors.Is(err, auth.ErrSessionNotFound) {
		t.Fatalf("after delete: got %v, want ErrSessionNotFound", err)
	}
}

func TestRedisStore_SessionTTL(t *testing.T) {
	mr, c := newRedis(t)
	store := auth.NewRedisStore(c)
	ctx := context.Background()

	if err := store.PutSession(ctx, "jti-2", "0xabc", 30*time.Second); err != nil {
		t.Fatalf("PutSession: %v", err)
	}
	mr.FastForward(time.Minute)
	if _, err := store.SessionExists(ctx, "jti-2"); !errors.Is(err, auth.ErrSessionNotFound) {
		t.Fatalf("expired session: got %v, want ErrSessionNotFound", err)
	}
}

func TestRedisStore_DeleteSession_Idempotent(t *testing.T) {
	_, c := newRedis(t)
	store := auth.NewRedisStore(c)
	ctx := context.Background()

	// Deleting an absent session must not error.
	if err := store.DeleteSession(ctx, "never-existed"); err != nil {
		t.Fatalf("delete absent: %v", err)
	}
	if err := store.DeleteSession(ctx, ""); err != nil {
		t.Fatalf("delete empty: %v", err)
	}
}

func TestRedisStore_NonceNamespaceDoesNotCollideWithSession(t *testing.T) {
	// Even with the same opaque key, nonce and session entries must be
	// independent. This regression-guards against a future refactor that
	// drops the prefix.
	_, c := newRedis(t)
	store := auth.NewRedisStore(c)
	ctx := context.Background()

	const id = "shared-id"
	if err := store.PutNonce(ctx, id, time.Minute); err != nil {
		t.Fatalf("PutNonce: %v", err)
	}
	if err := store.PutSession(ctx, id, "0xfeed", time.Hour); err != nil {
		t.Fatalf("PutSession: %v", err)
	}
	// Consuming the nonce must NOT remove the session.
	if err := store.ConsumeNonce(ctx, id); err != nil {
		t.Fatalf("ConsumeNonce: %v", err)
	}
	addr, err := store.SessionExists(ctx, id)
	if err != nil {
		t.Fatalf("session after nonce consume: %v", err)
	}
	if !strings.EqualFold(addr, "0xfeed") {
		t.Fatalf("session address tampered: %q", addr)
	}
}
