package auth_test

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/auth"
)

// makeSigner returns a signer with a 32-byte stub key suitable for tests
// where we don't care about key rotation.
func makeSigner(t *testing.T) *auth.JWTSigner {
	t.Helper()
	s, err := auth.NewJWTSigner(auth.SignerConfig{
		Secret: []byte(strings.Repeat("k", 32)),
		Issuer: "gateway-test",
		TTL:    time.Hour,
	})
	if err != nil {
		t.Fatalf("NewJWTSigner: %v", err)
	}
	return s
}

func TestNewJWTSigner_RejectsShortSecret(t *testing.T) {
	_, err := auth.NewJWTSigner(auth.SignerConfig{
		Secret: []byte("too-short"),
		Issuer: "gateway-test",
		TTL:    time.Hour,
	})
	if err == nil {
		t.Fatal("expected error for short secret, got nil")
	}
	if !strings.Contains(err.Error(), "too short") {
		t.Fatalf("error should mention length: %v", err)
	}
}

func TestNewJWTSigner_RejectsMissingIssuer(t *testing.T) {
	_, err := auth.NewJWTSigner(auth.SignerConfig{
		Secret: []byte(strings.Repeat("k", 32)),
		Issuer: "  ",
		TTL:    time.Hour,
	})
	if err == nil {
		t.Fatal("expected error for empty issuer, got nil")
	}
}

func TestNewJWTSigner_RejectsBadTTL(t *testing.T) {
	_, err := auth.NewJWTSigner(auth.SignerConfig{
		Secret: []byte(strings.Repeat("k", 32)),
		Issuer: "gateway-test",
		TTL:    0,
	})
	if err == nil {
		t.Fatal("expected error for zero ttl, got nil")
	}
}

func TestSigner_RoundTrip(t *testing.T) {
	s := makeSigner(t)
	tok, jti, err := s.Sign("0xABC")
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if tok == "" || jti == "" {
		t.Fatal("Sign returned empty values")
	}

	claims, err := s.Parse(tok)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if claims.Address != "0xabc" { // Sign lowercases addresses
		t.Fatalf("address: got %q, want lower-cased 0xabc", claims.Address)
	}
	if claims.ID != jti {
		t.Fatalf("jti: got %q, want %q", claims.ID, jti)
	}
	if claims.Issuer != "gateway-test" {
		t.Fatalf("issuer: got %q", claims.Issuer)
	}
}

func TestSigner_RejectsAlgNone(t *testing.T) {
	s := makeSigner(t)

	// Hand-roll an alg=none token whose claims look fine and feed it to Parse.
	// If Parse honours it we have classic CVE-2015-9235.
	claims := &auth.Claims{
		Address: "0xabc",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "gateway-test",
			Subject:   "0xabc",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			ID:        "j",
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	signed, err := tok.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("sign none: %v", err)
	}
	if _, err := s.Parse(signed); err == nil {
		t.Fatal("alg=none token was accepted; classic alg-confusion bypass")
	}
}

func TestSigner_RejectsExpired(t *testing.T) {
	s, err := auth.NewJWTSigner(auth.SignerConfig{
		Secret: []byte(strings.Repeat("k", 32)),
		Issuer: "gateway-test",
		TTL:    -time.Minute, // negative TTL would fail at construction
	})
	if err == nil {
		t.Fatalf("ttl<0 must fail at construction; got signer %v", s)
	}

	// Build a signer with a tiny TTL and let it expire.
	tiny, err := auth.NewJWTSigner(auth.SignerConfig{
		Secret: []byte(strings.Repeat("k", 32)),
		Issuer: "gateway-test",
		TTL:    time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewJWTSigner tiny: %v", err)
	}
	tok, _, err := tiny.Sign("0xabc")
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	time.Sleep(20 * time.Millisecond)
	if _, err := tiny.Parse(tok); err == nil {
		t.Fatal("expired token was accepted")
	}
}

func TestSigner_RejectsForeignSecret(t *testing.T) {
	a, err := auth.NewJWTSigner(auth.SignerConfig{
		Secret: []byte(strings.Repeat("a", 32)),
		Issuer: "gateway-test",
		TTL:    time.Hour,
	})
	if err != nil {
		t.Fatalf("signer A: %v", err)
	}
	b, err := auth.NewJWTSigner(auth.SignerConfig{
		Secret: []byte(strings.Repeat("b", 32)),
		Issuer: "gateway-test",
		TTL:    time.Hour,
	})
	if err != nil {
		t.Fatalf("signer B: %v", err)
	}
	tok, _, err := a.Sign("0xabc")
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if _, err := b.Parse(tok); err == nil {
		t.Fatal("token signed with secret A was accepted by signer B")
	}
}

func TestSigner_RejectsForeignIssuer(t *testing.T) {
	a, err := auth.NewJWTSigner(auth.SignerConfig{
		Secret: []byte(strings.Repeat("k", 32)),
		Issuer: "issuer-a",
		TTL:    time.Hour,
	})
	if err != nil {
		t.Fatalf("signer A: %v", err)
	}
	b, err := auth.NewJWTSigner(auth.SignerConfig{
		Secret: []byte(strings.Repeat("k", 32)),
		Issuer: "issuer-b",
		TTL:    time.Hour,
	})
	if err != nil {
		t.Fatalf("signer B: %v", err)
	}
	tok, _, err := a.Sign("0xabc")
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if _, err := b.Parse(tok); err == nil {
		t.Fatal("token from issuer-a accepted by issuer-b")
	}
}

func TestSigner_Parse_RejectsEmpty(t *testing.T) {
	s := makeSigner(t)
	if _, err := s.Parse(""); err == nil {
		t.Fatal("empty token accepted")
	}
	if _, err := s.Parse("   "); err == nil {
		t.Fatal("whitespace token accepted")
	}
	if _, err := s.Parse("not.a.jwt"); err == nil {
		t.Fatal("garbage token accepted")
	}
}

func TestSigner_Sign_RejectsEmptyAddress(t *testing.T) {
	s := makeSigner(t)
	if _, _, err := s.Sign("   "); err == nil {
		t.Fatal("empty address accepted")
	}
}

func TestSecretsEqual(t *testing.T) {
	if !auth.SecretsEqual([]byte("aaa"), []byte("aaa")) {
		t.Fatal("equal secrets reported unequal")
	}
	if auth.SecretsEqual([]byte("aaa"), []byte("aab")) {
		t.Fatal("unequal secrets reported equal")
	}
	if auth.SecretsEqual([]byte("aaa"), []byte("aaaa")) {
		t.Fatal("different-length secrets reported equal")
	}
}
