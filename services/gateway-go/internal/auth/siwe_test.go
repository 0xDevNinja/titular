package auth_test

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	siwe "github.com/spruceid/siwe-go"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/auth"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// signedSIWE returns the EIP-4361 message string and a hex-encoded signature
// for a freshly generated key. Tests use this rather than fixture vectors so
// the signing identity is opaque to the rest of the suite.
func signedSIWE(t *testing.T, domain string, chainID int, nonce string) (msg string, sig string, addr common.Address, key *ecdsa.PrivateKey) {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	addr = crypto.PubkeyToAddress(key.PublicKey)

	m, err := siwe.InitMessage(
		domain,
		addr.Hex(),
		"https://"+domain,
		nonce,
		map[string]interface{}{
			"chainId":   chainID,
			"statement": "Sign in to Titular",
			"issuedAt":  time.Now().UTC().Format(time.RFC3339),
		},
	)
	if err != nil {
		t.Fatalf("InitMessage: %v", err)
	}
	msg = m.String()

	// EIP-191 personal_sign of the message.
	prefixed := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(msg), msg)
	hash := crypto.Keccak256([]byte(prefixed))
	sigBytes, err := crypto.Sign(hash, key)
	if err != nil {
		t.Fatalf("crypto.Sign: %v", err)
	}
	// Geth returns recovery id 0/1; Ethereum convention is 27/28.
	sigBytes[64] += 27
	sig = "0x" + hex.EncodeToString(sigBytes)
	return
}

// newHandlers wires a fresh handler bundle with miniredis + 32-byte secret.
// Returns the handlers, the underlying store, and the signer for tests that
// need to mint tokens directly.
func newHandlers(t *testing.T, domain string, chainID int) (*auth.Handlers, auth.SessionStore, *auth.JWTSigner) {
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
	h, err := auth.NewHandlers(auth.HandlerConfig{
		Store:    store,
		Signer:   signer,
		Domain:   domain,
		ChainID:  chainID,
		NonceTTL: 5 * time.Minute,
	})
	if err != nil {
		t.Fatalf("NewHandlers: %v", err)
	}
	return h, store, signer
}

// route mounts the auth handlers on an isolated Gin engine for the duration
// of one test; we don't pull in the production router because that would
// couple the auth tests to fixture loading.
func route(h *auth.Handlers) *gin.Engine {
	r := gin.New()
	g := r.Group("/auth")
	g.POST("/siwe/nonce", h.Nonce)
	g.POST("/siwe/verify", h.Verify)
	g.POST("/logout", h.Logout)
	return r
}

func TestNewHandlers_Validation(t *testing.T) {
	signer, _ := auth.NewJWTSigner(auth.SignerConfig{
		Secret: []byte(strings.Repeat("k", 32)),
		Issuer: "gateway-test",
		TTL:    time.Hour,
	})
	_, c := newRedis(t)
	store := auth.NewRedisStore(c)

	cases := []struct {
		name string
		cfg  auth.HandlerConfig
	}{
		{"nil store", auth.HandlerConfig{Signer: signer, Domain: "d", ChainID: 1}},
		{"nil signer", auth.HandlerConfig{Store: store, Domain: "d", ChainID: 1}},
		{"empty domain", auth.HandlerConfig{Store: store, Signer: signer, Domain: "", ChainID: 1}},
		{"zero chain", auth.HandlerConfig{Store: store, Signer: signer, Domain: "d", ChainID: 0}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if _, err := auth.NewHandlers(tc.cfg); err == nil {
				t.Fatalf("expected error for %s", tc.name)
			}
		})
	}
}

func TestNonce_IssuesAndStores(t *testing.T) {
	h, _, _ := newHandlers(t, "titular.test", 84532)
	r := route(h)

	req := httptest.NewRequest(http.MethodPost, "/auth/siwe/nonce", nil)
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

func TestVerify_HappyPath_IssuesJWTAndConsumesNonce(t *testing.T) {
	h, store, signer := newHandlers(t, "titular.test", 84532)
	r := route(h)

	// Issue a nonce out-of-band so the test owns the value.
	nonce, _ := auth.RandomNonce(16)
	if err := store.PutNonce(t.Context(), nonce, time.Minute); err != nil {
		t.Fatalf("seed nonce: %v", err)
	}

	msg, sig, addr, _ := signedSIWE(t, "titular.test", 84532, nonce)

	body, _ := json.Marshal(map[string]string{"message": msg, "signature": sig})
	req := httptest.NewRequest(http.MethodPost, "/auth/siwe/verify", bytes.NewReader(body))
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
	if !strings.EqualFold(resp.Address, addr.Hex()) {
		t.Fatalf("address: got %q, want %q", resp.Address, strings.ToLower(addr.Hex()))
	}
	// Token must parse back to the same JTI + address.
	claims, err := signer.Parse(resp.Token)
	if err != nil {
		t.Fatalf("parse issued token: %v", err)
	}
	if !strings.EqualFold(claims.Address, addr.Hex()) {
		t.Fatalf("token addr: got %q, want %q", claims.Address, strings.ToLower(addr.Hex()))
	}
	// Session must be live in Redis under the JTI.
	if _, err := store.SessionExists(t.Context(), claims.ID); err != nil {
		t.Fatalf("session not present: %v", err)
	}

	// Nonce must be gone — single-use.
	if err := store.ConsumeNonce(t.Context(), nonce); err == nil {
		t.Fatal("nonce was not consumed by verify")
	}
}

func TestVerify_ReplayNonce_RefusedOnSecondCall(t *testing.T) {
	h, store, _ := newHandlers(t, "titular.test", 84532)
	r := route(h)

	nonce, _ := auth.RandomNonce(16)
	_ = store.PutNonce(t.Context(), nonce, time.Minute)
	msg, sig, _, _ := signedSIWE(t, "titular.test", 84532, nonce)

	body, _ := json.Marshal(map[string]string{"message": msg, "signature": sig})

	// First verify: success.
	req := httptest.NewRequest(http.MethodPost, "/auth/siwe/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("first verify: got %d, body=%s", rec.Code, rec.Body.String())
	}

	// Second verify with the same body: must fail with 401 invalid_nonce.
	req = httptest.NewRequest(http.MethodPost, "/auth/siwe/verify", bytes.NewReader(body))
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

func TestVerify_DomainMismatch(t *testing.T) {
	h, store, _ := newHandlers(t, "titular.test", 84532)
	r := route(h)

	nonce, _ := auth.RandomNonce(16)
	_ = store.PutNonce(t.Context(), nonce, time.Minute)

	// Sign for a different domain.
	msg, sig, _, _ := signedSIWE(t, "evil.test", 84532, nonce)
	body, _ := json.Marshal(map[string]string{"message": msg, "signature": sig})

	req := httptest.NewRequest(http.MethodPost, "/auth/siwe/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "domain_mismatch") {
		t.Fatalf("expected domain_mismatch, got %s", rec.Body.String())
	}
	// Critically: the nonce must NOT have been consumed (we reject before
	// ConsumeNonce). Verify by issuing a corrected message and confirming it
	// still works.
	msgGood, sigGood, _, _ := signedSIWE(t, "titular.test", 84532, nonce)
	body, _ = json.Marshal(map[string]string{"message": msgGood, "signature": sigGood})
	req = httptest.NewRequest(http.MethodPost, "/auth/siwe/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("retry after domain reject: got %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestVerify_ChainMismatch(t *testing.T) {
	h, store, _ := newHandlers(t, "titular.test", 84532)
	r := route(h)

	nonce, _ := auth.RandomNonce(16)
	_ = store.PutNonce(t.Context(), nonce, time.Minute)

	// Sign for mainnet (chain 1) instead of Base Sepolia.
	msg, sig, _, _ := signedSIWE(t, "titular.test", 1, nonce)
	body, _ := json.Marshal(map[string]string{"message": msg, "signature": sig})

	req := httptest.NewRequest(http.MethodPost, "/auth/siwe/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "chain_mismatch") {
		t.Fatalf("expected chain_mismatch, got %s", rec.Body.String())
	}
}

func TestVerify_BadSignature_ConsumesNonceAnyway(t *testing.T) {
	// Once we've decided the static checks pass, ConsumeNonce runs BEFORE
	// signature verify so a single nonce can never service two attempts —
	// even if the first attempt uses a bad sig. This test pins that
	// behaviour.
	h, store, _ := newHandlers(t, "titular.test", 84532)
	r := route(h)

	nonce, _ := auth.RandomNonce(16)
	_ = store.PutNonce(t.Context(), nonce, time.Minute)

	msg, _, _, _ := signedSIWE(t, "titular.test", 84532, nonce)
	// Forge a signature: same shape, wrong bytes.
	badSig := "0x" + strings.Repeat("11", 65)

	body, _ := json.Marshal(map[string]string{"message": msg, "signature": badSig})
	req := httptest.NewRequest(http.MethodPost, "/auth/siwe/verify", bytes.NewReader(body))
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

func TestVerify_MalformedJSON(t *testing.T) {
	h, _, _ := newHandlers(t, "titular.test", 84532)
	r := route(h)

	req := httptest.NewRequest(http.MethodPost, "/auth/siwe/verify", strings.NewReader("{not json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestVerify_UnparseableSIWEMessage(t *testing.T) {
	h, _, _ := newHandlers(t, "titular.test", 84532)
	r := route(h)

	body, _ := json.Marshal(map[string]string{"message": "definitely not siwe", "signature": "0x00"})
	req := httptest.NewRequest(http.MethodPost, "/auth/siwe/verify", bytes.NewReader(body))
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

func TestLogout_ValidToken_DeletesSession(t *testing.T) {
	h, store, signer := newHandlers(t, "titular.test", 84532)
	r := route(h)

	tok, jti, err := signer.Sign("0xabc")
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	_ = store.PutSession(t.Context(), jti, "0xabc", time.Hour)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status: got %d, body=%s", rec.Code, rec.Body.String())
	}
	if _, err := store.SessionExists(t.Context(), jti); err == nil {
		t.Fatal("session still present after logout")
	}
}

func TestLogout_NoToken_StillReturns204(t *testing.T) {
	h, _, _ := newHandlers(t, "titular.test", 84532)
	r := route(h)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("no-token logout: got %d, want 204", rec.Code)
	}
}

func TestLogout_BadToken_StillReturns204(t *testing.T) {
	h, _, _ := newHandlers(t, "titular.test", 84532)
	r := route(h)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-jwt")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("bad-token logout: got %d, want 204", rec.Code)
	}
}

// TestVerify_ResponseBodyDoesNotLeakSignerKey is a defence-in-depth sanity
// check that ensures the response envelope contains exactly the documented
// fields and does not accidentally serialise the signer's HMAC key (which
// would be a catastrophic leak).
func TestVerify_ResponseBodyDoesNotLeakSignerKey(t *testing.T) {
	h, store, _ := newHandlers(t, "titular.test", 84532)
	r := route(h)

	nonce, _ := auth.RandomNonce(16)
	_ = store.PutNonce(t.Context(), nonce, time.Minute)
	msg, sig, _, _ := signedSIWE(t, "titular.test", 84532, nonce)

	body, _ := json.Marshal(map[string]string{"message": msg, "signature": sig})
	req := httptest.NewRequest(http.MethodPost, "/auth/siwe/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d", rec.Code)
	}
	if strings.Contains(rec.Body.String(), strings.Repeat("k", 32)) {
		t.Fatal("verify response body leaked the HMAC secret")
	}
	// Reject any unknown top-level field that might have been added.
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
