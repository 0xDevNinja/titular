package auth_test

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mr-tron/base58"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/auth"
)

// signedSIWS returns a canonical SIWS message string, a base58 ed25519
// signature over it, and the signer's pubkey (base58, == Solana address).
// We synthesise rather than hard-coding fixtures so the keys are opaque to
// the rest of the suite.
func signedSIWS(t *testing.T, domain, cluster, nonce string) (msg, sig, addr string, priv ed25519.PrivateKey) {
	t.Helper()
	pub, p, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	priv = p
	addr = base58.Encode(pub)

	issuedAt := time.Now().UTC().Format(time.RFC3339)
	msg = fmt.Sprintf(
		"%s wants you to sign in with your Solana account:\n%s\n\nSign in to Titular\n\nURI: https://%s\nVersion: 1\nChain ID: %s\nNonce: %s\nIssued At: %s",
		domain, addr, domain, cluster, nonce, issuedAt,
	)
	sigBytes := ed25519.Sign(priv, []byte(msg))
	sig = base58.Encode(sigBytes)
	return
}

// newSIWSHandlers wires a fresh SIWS handler bundle backed by miniredis.
func newSIWSHandlers(t *testing.T, domain, cluster string) (*auth.SIWSHandlers, auth.SessionStore, *auth.JWTSigner) {
	t.Helper()
	_, c := newRedis(t)
	store := auth.NewRedisStore(c)
	signer, err := auth.NewJWTSigner(auth.SignerConfig{
		Secret: []byte(strings.Repeat("k", 32)),
		Issuer: "gateway-test",
		TTL:    time.Hour,
	})
	if err != nil {
		t.Fatalf("NewJWTSigner: %v", err)
	}
	h, err := auth.NewSIWSHandlers(auth.SIWSHandlerConfig{
		Store:    store,
		Signer:   signer,
		Domain:   domain,
		Cluster:  cluster,
		NonceTTL: 5 * time.Minute,
	})
	if err != nil {
		t.Fatalf("NewSIWSHandlers: %v", err)
	}
	return h, store, signer
}

func siwsRoute(h *auth.SIWSHandlers) *gin.Engine {
	r := gin.New()
	g := r.Group("/auth")
	g.POST("/siws/nonce", h.Nonce)
	g.POST("/siws/verify", h.Verify)
	return r
}

func TestSIWS_NewHandlers_Validation(t *testing.T) {
	signer, _ := auth.NewJWTSigner(auth.SignerConfig{
		Secret: []byte(strings.Repeat("k", 32)),
		Issuer: "gateway-test",
		TTL:    time.Hour,
	})
	_, c := newRedis(t)
	store := auth.NewRedisStore(c)

	cases := []struct {
		name string
		cfg  auth.SIWSHandlerConfig
	}{
		{"nil store", auth.SIWSHandlerConfig{Signer: signer, Domain: "d", Cluster: auth.ClusterDevnet}},
		{"nil signer", auth.SIWSHandlerConfig{Store: store, Domain: "d", Cluster: auth.ClusterDevnet}},
		{"empty domain", auth.SIWSHandlerConfig{Store: store, Signer: signer, Cluster: auth.ClusterDevnet}},
		{"empty cluster", auth.SIWSHandlerConfig{Store: store, Signer: signer, Domain: "d"}},
		{"bad cluster", auth.SIWSHandlerConfig{Store: store, Signer: signer, Domain: "d", Cluster: "polkadot"}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if _, err := auth.NewSIWSHandlers(tc.cfg); err == nil {
				t.Fatalf("expected error for %s", tc.name)
			}
		})
	}
}

func TestSIWS_Nonce_IssuesAndStores(t *testing.T) {
	h, _, _ := newSIWSHandlers(t, "titular.test", auth.ClusterDevnet)
	r := siwsRoute(h)

	req := httptest.NewRequest(http.MethodPost, "/auth/siws/nonce", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}
	var body struct {
		Nonce     string `json:"nonce"`
		ExpiresAt int64  `json:"expires_at"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body.Nonce) == 0 {
		t.Fatal("empty nonce returned")
	}
	if body.ExpiresAt <= time.Now().Unix() {
		t.Fatalf("expires_at %d is in the past", body.ExpiresAt)
	}
}

func TestSIWS_Verify_HappyPath_IssuesJWTAndConsumesNonce(t *testing.T) {
	h, store, signer := newSIWSHandlers(t, "titular.test", auth.ClusterDevnet)
	r := siwsRoute(h)

	nonce, _ := auth.RandomNonce(16)
	if err := store.PutNonce(t.Context(), nonce, time.Minute); err != nil {
		t.Fatalf("seed nonce: %v", err)
	}

	msg, sig, addr, _ := signedSIWS(t, "titular.test", auth.ClusterDevnet, nonce)

	body, _ := json.Marshal(map[string]string{"message": msg, "signature": sig})
	req := httptest.NewRequest(http.MethodPost, "/auth/siws/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, body=%s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Token   string `json:"token"`
		Address string `json:"address"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Address != addr {
		t.Fatalf("address: got %q, want %q", resp.Address, addr)
	}
	// Token must parse back to a session bound to the same lowercased address.
	claims, err := signer.Parse(resp.Token)
	if err != nil {
		t.Fatalf("parse issued token: %v", err)
	}
	if !strings.EqualFold(claims.Address, addr) {
		t.Fatalf("token addr: got %q, want %q", claims.Address, addr)
	}
	if _, err := store.SessionExists(t.Context(), claims.ID); err != nil {
		t.Fatalf("session not present: %v", err)
	}
	// Nonce must be gone — single-use.
	if err := store.ConsumeNonce(t.Context(), nonce); err == nil {
		t.Fatal("nonce was not consumed by verify")
	}
}

func TestSIWS_Verify_ReplayNonce_RefusedOnSecondCall(t *testing.T) {
	h, store, _ := newSIWSHandlers(t, "titular.test", auth.ClusterDevnet)
	r := siwsRoute(h)

	nonce, _ := auth.RandomNonce(16)
	_ = store.PutNonce(t.Context(), nonce, time.Minute)
	msg, sig, _, _ := signedSIWS(t, "titular.test", auth.ClusterDevnet, nonce)

	body, _ := json.Marshal(map[string]string{"message": msg, "signature": sig})

	req := httptest.NewRequest(http.MethodPost, "/auth/siws/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("first verify: got %d, body=%s", rec.Code, rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodPost, "/auth/siws/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("replay: got %d, want 401, body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "invalid_nonce") {
		t.Fatalf("replay body: %s", rec.Body.String())
	}
}

func TestSIWS_Verify_DomainMismatch(t *testing.T) {
	h, store, _ := newSIWSHandlers(t, "titular.test", auth.ClusterDevnet)
	r := siwsRoute(h)

	nonce, _ := auth.RandomNonce(16)
	_ = store.PutNonce(t.Context(), nonce, time.Minute)

	msg, sig, _, _ := signedSIWS(t, "evil.test", auth.ClusterDevnet, nonce)
	body, _ := json.Marshal(map[string]string{"message": msg, "signature": sig})

	req := httptest.NewRequest(http.MethodPost, "/auth/siws/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "domain_mismatch") {
		t.Fatalf("expected domain_mismatch, got %s", rec.Body.String())
	}
	// Nonce must NOT have been consumed (we reject before ConsumeNonce).
	msgGood, sigGood, _, _ := signedSIWS(t, "titular.test", auth.ClusterDevnet, nonce)
	body, _ = json.Marshal(map[string]string{"message": msgGood, "signature": sigGood})
	req = httptest.NewRequest(http.MethodPost, "/auth/siws/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("retry after domain reject: got %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestSIWS_Verify_ClusterMismatch(t *testing.T) {
	h, store, _ := newSIWSHandlers(t, "titular.test", auth.ClusterDevnet)
	r := siwsRoute(h)

	nonce, _ := auth.RandomNonce(16)
	_ = store.PutNonce(t.Context(), nonce, time.Minute)

	// Sign for mainnet-beta when the gateway is pinned to devnet.
	msg, sig, _, _ := signedSIWS(t, "titular.test", auth.ClusterMainnetBeta, nonce)
	body, _ := json.Marshal(map[string]string{"message": msg, "signature": sig})

	req := httptest.NewRequest(http.MethodPost, "/auth/siws/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "cluster_mismatch") {
		t.Fatalf("expected cluster_mismatch, got %s", rec.Body.String())
	}
}

func TestSIWS_Verify_BadSignature_ConsumesNonceAnyway(t *testing.T) {
	// Once static checks pass, ConsumeNonce runs BEFORE signature verify
	// so a single nonce can never service two attempts — even one with a
	// bad sig. This test pins that behaviour.
	h, store, _ := newSIWSHandlers(t, "titular.test", auth.ClusterDevnet)
	r := siwsRoute(h)

	nonce, _ := auth.RandomNonce(16)
	_ = store.PutNonce(t.Context(), nonce, time.Minute)

	msg, _, _, _ := signedSIWS(t, "titular.test", auth.ClusterDevnet, nonce)
	// Forge a sig: 64 bytes of zeros, base58 encoded — wrong shape but
	// passes the size check after decode (still won't verify).
	badSig := base58.Encode(make([]byte, ed25519.SignatureSize))

	body, _ := json.Marshal(map[string]string{"message": msg, "signature": badSig})
	req := httptest.NewRequest(http.MethodPost, "/auth/siws/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: got %d, body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "invalid_signature") {
		t.Fatalf("expected invalid_signature, got %s", rec.Body.String())
	}
	// Nonce must be gone.
	if err := store.ConsumeNonce(t.Context(), nonce); err == nil {
		t.Fatal("bad-sig verify left the nonce intact; that allows brute force")
	}
}

func TestSIWS_Verify_WrongSignerForeignKey(t *testing.T) {
	// Sig is well-formed but produced by a different key than the message
	// claims. Must fail with invalid_signature.
	h, store, _ := newSIWSHandlers(t, "titular.test", auth.ClusterDevnet)
	r := siwsRoute(h)

	nonce, _ := auth.RandomNonce(16)
	_ = store.PutNonce(t.Context(), nonce, time.Minute)

	msg, _, _, _ := signedSIWS(t, "titular.test", auth.ClusterDevnet, nonce)
	// Re-sign the message bytes with a foreign key.
	_, foreign, _ := ed25519.GenerateKey(rand.Reader)
	foreignSigBytes := ed25519.Sign(foreign, []byte(msg))
	foreignSig := base58.Encode(foreignSigBytes)

	body, _ := json.Marshal(map[string]string{"message": msg, "signature": foreignSig})
	req := httptest.NewRequest(http.MethodPost, "/auth/siws/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: got %d, body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "invalid_signature") {
		t.Fatalf("expected invalid_signature, got %s", rec.Body.String())
	}
}

func TestSIWS_Verify_TamperedMessage(t *testing.T) {
	// Sig is valid for the original message but the message we POST is
	// modified by one byte. ed25519.Verify must fail.
	h, store, _ := newSIWSHandlers(t, "titular.test", auth.ClusterDevnet)
	r := siwsRoute(h)

	nonce, _ := auth.RandomNonce(16)
	_ = store.PutNonce(t.Context(), nonce, time.Minute)

	msg, sig, _, _ := signedSIWS(t, "titular.test", auth.ClusterDevnet, nonce)
	tampered := strings.Replace(msg, "Sign in to Titular", "Sign in to Evil   ", 1)

	body, _ := json.Marshal(map[string]string{"message": tampered, "signature": sig})
	req := httptest.NewRequest(http.MethodPost, "/auth/siws/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: got %d, body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "invalid_signature") {
		t.Fatalf("expected invalid_signature, got %s", rec.Body.String())
	}
}

func TestSIWS_Verify_MalformedJSON(t *testing.T) {
	h, _, _ := newSIWSHandlers(t, "titular.test", auth.ClusterDevnet)
	r := siwsRoute(h)

	req := httptest.NewRequest(http.MethodPost, "/auth/siws/verify", strings.NewReader("{not json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestSIWS_Verify_UnparseableMessage(t *testing.T) {
	h, _, _ := newSIWSHandlers(t, "titular.test", auth.ClusterDevnet)
	r := siwsRoute(h)

	body, _ := json.Marshal(map[string]string{"message": "definitely not siws", "signature": base58.Encode(make([]byte, 64))})
	req := httptest.NewRequest(http.MethodPost, "/auth/siws/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "invalid_message") {
		t.Fatalf("expected invalid_message, got %s", rec.Body.String())
	}
}

func TestSIWS_Verify_InvalidAddress_NotBase58(t *testing.T) {
	h, store, _ := newSIWSHandlers(t, "titular.test", auth.ClusterDevnet)
	r := siwsRoute(h)

	nonce, _ := auth.RandomNonce(16)
	_ = store.PutNonce(t.Context(), nonce, time.Minute)

	// Hand-craft a message with a malformed address (contains '0' which
	// is not a valid base58 char — base58 alphabet skips 0OIl).
	msg := fmt.Sprintf(
		"titular.test wants you to sign in with your Solana account:\n0badaddress\n\nSign in to Titular\n\nURI: https://titular.test\nVersion: 1\nChain ID: %s\nNonce: %s\nIssued At: %s",
		auth.ClusterDevnet, nonce, time.Now().UTC().Format(time.RFC3339),
	)
	body, _ := json.Marshal(map[string]string{"message": msg, "signature": base58.Encode(make([]byte, 64))})
	req := httptest.NewRequest(http.MethodPost, "/auth/siws/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "invalid_address") {
		t.Fatalf("expected invalid_address, got %s", rec.Body.String())
	}
	// Nonce must NOT have been consumed.
	if err := store.ConsumeNonce(t.Context(), nonce); err != nil {
		t.Fatalf("invalid_address path consumed the nonce: %v", err)
	}
}

func TestSIWS_Verify_ExpiredMessage(t *testing.T) {
	h, store, _ := newSIWSHandlers(t, "titular.test", auth.ClusterDevnet)
	r := siwsRoute(h)

	nonce, _ := auth.RandomNonce(16)
	_ = store.PutNonce(t.Context(), nonce, time.Minute)

	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	addr := base58.Encode(pub)
	past := time.Now().Add(-time.Hour).UTC().Format(time.RFC3339)
	msg := fmt.Sprintf(
		"titular.test wants you to sign in with your Solana account:\n%s\n\nSign in to Titular\n\nURI: https://titular.test\nVersion: 1\nChain ID: %s\nNonce: %s\nIssued At: %s\nExpiration Time: %s",
		addr, auth.ClusterDevnet, nonce, past, past,
	)
	sig := base58.Encode(ed25519.Sign(priv, []byte(msg)))
	body, _ := json.Marshal(map[string]string{"message": msg, "signature": sig})
	req := httptest.NewRequest(http.MethodPost, "/auth/siws/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "expired_message") {
		t.Fatalf("expected expired_message, got %s", rec.Body.String())
	}
}

func TestSIWS_Verify_ResponseBodyDoesNotLeakSignerKey(t *testing.T) {
	// Defence-in-depth: the response envelope must contain exactly the
	// documented fields and never serialise the HMAC secret.
	h, store, _ := newSIWSHandlers(t, "titular.test", auth.ClusterDevnet)
	r := siwsRoute(h)

	nonce, _ := auth.RandomNonce(16)
	_ = store.PutNonce(t.Context(), nonce, time.Minute)
	msg, sig, _, _ := signedSIWS(t, "titular.test", auth.ClusterDevnet, nonce)

	body, _ := json.Marshal(map[string]string{"message": msg, "signature": sig})
	req := httptest.NewRequest(http.MethodPost, "/auth/siws/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d", rec.Code)
	}
	if strings.Contains(rec.Body.String(), strings.Repeat("k", 32)) {
		t.Fatal("verify response body leaked the HMAC secret")
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &raw); err != nil {
		t.Fatalf("decode: %v", err)
	}
	allowed := map[string]bool{"token": true, "address": true, "expires_at": true}
	for k := range raw {
		if !allowed[k] {
			t.Errorf("unexpected response field %q", k)
		}
	}
}

// TestSIWS_AuthMiddleware_AcceptsSIWSToken proves the existing RequireAuth
// middleware accepts a JWT minted by the SIWS path. The two flows MUST be
// indistinguishable at the protected-route layer.
func TestSIWS_AuthMiddleware_AcceptsSIWSToken(t *testing.T) {
	h, store, signer := newSIWSHandlers(t, "titular.test", auth.ClusterDevnet)
	r := siwsRoute(h)

	nonce, _ := auth.RandomNonce(16)
	_ = store.PutNonce(t.Context(), nonce, time.Minute)
	msg, sig, addr, _ := signedSIWS(t, "titular.test", auth.ClusterDevnet, nonce)

	body, _ := json.Marshal(map[string]string{"message": msg, "signature": sig})
	req := httptest.NewRequest(http.MethodPost, "/auth/siws/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("verify failed: %d %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Token string `json:"token"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)

	// Build a SIWE Handlers bound to the SAME store and signer so the
	// shared RequireAuth path can be exercised. This mirrors the runtime
	// wiring where one *auth.Handlers governs logout for both flows.
	siwe, err := auth.NewHandlers(auth.HandlerConfig{
		Store:   store,
		Signer:  signer,
		Domain:  "titular.test",
		ChainID: 84532,
	})
	if err != nil {
		t.Fatalf("NewHandlers: %v", err)
	}
	guarded := gin.New()
	guarded.GET("/me", auth.RequireAuth(siwe), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"address": auth.AddressFromContext(c)})
	})

	req = httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+resp.Token)
	rec = httptest.NewRecorder()
	guarded.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("RequireAuth rejected SIWS token: %d %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), addr) {
		t.Fatalf("address not propagated: %s", rec.Body.String())
	}
}

// TestParseSIWSMessage_RoundTrip exercises the parser directly with both
// the SIWE-style "Chain ID:" key and the alternate "Cluster:" key, plus
// optional fields.
func TestParseSIWSMessage_RoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	cases := []struct {
		name string
		raw  string
	}{
		{
			name: "chain id key minimal",
			raw: fmt.Sprintf(
				"titular.test wants you to sign in with your Solana account:\n4Nd1m5dJfMaPjV4z9LhPwbqoGSoMpQ6gPYwJoT3KQA8q\n\nURI: https://titular.test\nVersion: 1\nChain ID: devnet\nNonce: aabbcc\nIssued At: %s",
				now.Format(time.RFC3339),
			),
		},
		{
			name: "cluster key with statement and expiration",
			raw: fmt.Sprintf(
				"titular.test wants you to sign in with your Solana account:\n4Nd1m5dJfMaPjV4z9LhPwbqoGSoMpQ6gPYwJoT3KQA8q\n\nWelcome to Titular\n\nURI: https://titular.test\nVersion: 1\nCluster: mainnet-beta\nNonce: aabbcc\nIssued At: %s\nExpiration Time: %s\nRequest ID: req-1",
				now.Format(time.RFC3339), now.Add(time.Hour).Format(time.RFC3339),
			),
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			m, err := auth.ParseSIWSMessage(tc.raw)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if m.Domain != "titular.test" {
				t.Errorf("domain: %q", m.Domain)
			}
			if m.Nonce != "aabbcc" {
				t.Errorf("nonce: %q", m.Nonce)
			}
			if m.Version != "1" {
				t.Errorf("version: %q", m.Version)
			}
		})
	}
}

func TestParseSIWSMessage_RejectsBadShape(t *testing.T) {
	cases := []struct {
		name string
		raw  string
	}{
		{"empty", ""},
		{"missing header", "not-titular wants...\naddr\n\nURI: x\nVersion: 1\nChain ID: devnet\nNonce: n\nIssued At: 2026-01-01T00:00:00Z"},
		{"missing required", "titular.test wants you to sign in with your Solana account:\naddr\n\nURI: x\nVersion: 1\nNonce: n\nIssued At: 2026-01-01T00:00:00Z"},
		{"bad time", "titular.test wants you to sign in with your Solana account:\naddr\n\nURI: x\nVersion: 1\nChain ID: devnet\nNonce: n\nIssued At: not-a-time"},
		{"unknown field", "titular.test wants you to sign in with your Solana account:\naddr\n\nURI: x\nVersion: 1\nChain ID: devnet\nNonce: n\nIssued At: 2026-01-01T00:00:00Z\nFavourite Colour: blue"},
		{"duplicate field", "titular.test wants you to sign in with your Solana account:\naddr\n\nURI: x\nVersion: 1\nChain ID: devnet\nChain ID: testnet\nNonce: n\nIssued At: 2026-01-01T00:00:00Z"},
		{"control char in nonce", "titular.test wants you to sign in with your Solana account:\naddr\n\nURI: x\nVersion: 1\nChain ID: devnet\nNonce: \x01bad\nIssued At: 2026-01-01T00:00:00Z"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if _, err := auth.ParseSIWSMessage(tc.raw); err == nil {
				t.Fatalf("expected error for %s", tc.name)
			}
		})
	}
}
