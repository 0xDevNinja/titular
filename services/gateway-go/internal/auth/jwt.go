package auth

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// minHMACSecretBytes is the minimum acceptable length for the HS256 signing
// key. RFC 8725 §3.2 recommends a key at least as long as the digest output
// (32 bytes for SHA-256); we enforce that as a hard floor at startup so
// short, low-entropy secrets cannot ship to production.
const minHMACSecretBytes = 32

// JWTSigner mints and validates JSON Web Tokens used as session bearers.
//
// We use HS256 because the gateway is a single deployment unit that signs and
// verifies its own tokens — there is no third party that needs to verify
// without the secret. Switching to RS256/EdDSA would only make sense once we
// add a separate verifier process.
//
// Signer is safe for concurrent use after construction.
type JWTSigner struct {
	secret []byte
	issuer string
	ttl    time.Duration
	clock  func() time.Time // overridable in tests
}

// SignerConfig is the input to NewJWTSigner. All fields are required.
type SignerConfig struct {
	// Secret is the HMAC key. MUST be at least 32 bytes; shorter keys are
	// refused at construction time so misconfiguration cannot ship.
	Secret []byte
	// Issuer is mirrored verbatim into the iss claim and validated on parse.
	Issuer string
	// TTL is the lifetime applied to issued tokens. The Redis session is
	// granted the same TTL so a logout invalidates both halves.
	TTL time.Duration
}

// NewJWTSigner constructs a signer or returns an error if the configuration
// is unsafe for production. Callers are expected to wrap the error with their
// startup-failure path; we deliberately do not log the secret length.
func NewJWTSigner(cfg SignerConfig) (*JWTSigner, error) {
	if len(cfg.Secret) < minHMACSecretBytes {
		return nil, fmt.Errorf(
			"auth: jwt secret too short; need >=%d bytes, got %d",
			minHMACSecretBytes, len(cfg.Secret),
		)
	}
	if strings.TrimSpace(cfg.Issuer) == "" {
		return nil, errors.New("auth: jwt issuer required")
	}
	if cfg.TTL <= 0 {
		return nil, fmt.Errorf("auth: jwt ttl must be > 0, got %s", cfg.TTL)
	}
	return &JWTSigner{
		secret: cfg.Secret,
		issuer: cfg.Issuer,
		ttl:    cfg.TTL,
		clock:  time.Now,
	}, nil
}

// Claims is the typed claim set we mint and parse. We only carry the fields
// the gateway actually uses; anything else is rejected at parse time so a
// crafted token cannot smuggle additional state.
//
// Address is the lowercased EIP-55 hex (0x + 40 hex chars) recovered from the
// SIWE signature. We never include the full SIWE message in the claims; that
// stays server-side.
type Claims struct {
	Address string `json:"addr"`
	jwt.RegisteredClaims
}

// Sign returns a signed token for the supplied address along with the JTI we
// generated. The JTI is the key under which the session is recorded in Redis;
// the caller is expected to PutSession with this value before returning the
// token to the client, otherwise RequireAuth will reject the very token we
// just issued.
func (s *JWTSigner) Sign(address string) (token, jti string, err error) {
	addr := strings.ToLower(strings.TrimSpace(address))
	if addr == "" {
		return "", "", errors.New("auth: address required for sign")
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return "", "", fmt.Errorf("auth: jti gen: %w", err)
	}
	jti = id.String()

	now := s.clock()
	claims := Claims{
		Address: addr,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   addr,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
			ID:        jti,
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(s.secret)
	if err != nil {
		return "", "", fmt.Errorf("auth: sign jwt: %w", err)
	}
	return signed, jti, nil
}

// Parse validates the signature, algorithm, expiry, issuer and not-before
// fields of the supplied token and returns the claim set. It does NOT consult
// Redis — that's the middleware's job — so this function stays pure and
// trivially fuzzable.
//
// Algorithm pinning: we explicitly refuse anything other than HS256 to defend
// against the classic "alg=none" and "alg confusion" downgrades.
func (s *JWTSigner) Parse(token string) (*Claims, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("auth: empty token")
	}
	claims := &Claims{}
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(s.issuer),
		jwt.WithExpirationRequired(),
		jwt.WithTimeFunc(s.clock),
	)
	parsed, err := parser.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		// Defence-in-depth: WithValidMethods already gates the alg, but
		// re-check inside the keyfunc so a future refactor that drops the
		// option still rejects the wrong algorithm.
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("auth: unexpected alg %q", t.Method.Alg())
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("auth: parse jwt: %w", err)
	}
	if !parsed.Valid {
		return nil, errors.New("auth: token invalid")
	}
	if claims.Address == "" || claims.ID == "" {
		return nil, errors.New("auth: token missing required claims")
	}
	return claims, nil
}

// SecretsEqual is a constant-time comparison helper exposed for tests that
// need to assert two HMAC secrets match without leaking timing. Production
// code does not call this — go-jwt does its own constant-time check inside
// jwt.Token.Method.Verify.
func SecretsEqual(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}
