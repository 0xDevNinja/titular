package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	siwe "github.com/spruceid/siwe-go"
)

// Default lifetimes. Operators can override the session TTL via SignerConfig.TTL;
// the nonce TTL is fixed because there is no legitimate reason to keep an
// unused nonce alive longer than the user-facing wallet popup is open.
const (
	// DefaultNonceTTL is the lifetime of an unconsumed nonce. Five minutes is
	// the upper bound on a wallet sign prompt; anything longer increases the
	// window in which a captured nonce could be replayed.
	DefaultNonceTTL = 5 * time.Minute

	// nonceByteLen is the entropy of a fresh nonce in bytes. 16 bytes is
	// 128 bits, far above the EIP-4361 recommendation.
	nonceByteLen = 16
)

// HandlerConfig holds the request-time dependencies for the SIWE handlers.
//
// Domain is the value the message MUST declare as the requesting site. We
// pin it server-side rather than echoing whatever the message contained so
// that a SIWE message signed for some-other-app.example cannot be replayed
// against this gateway.
//
// ChainID is the chain the message MUST declare. Mismatched chain IDs are an
// error so a Mainnet-signed message can't be replayed against a Sepolia
// deployment.
//
// NonceTTL controls how long an issued nonce stays valid. Defaults to
// DefaultNonceTTL when zero.
type HandlerConfig struct {
	Store    SessionStore
	Signer   *JWTSigner
	Domain   string
	ChainID  int
	NonceTTL time.Duration
}

// Handlers exposes the SIWE HTTP handlers as Gin handler funcs. Construct via
// NewHandlers — the zero value is unsafe.
type Handlers struct {
	cfg HandlerConfig
	now func() time.Time // overridable in tests for ValidAt assertions
}

// NewHandlers wires the handler dependencies after validating that all the
// required fields are present. Failing here at startup is preferable to
// failing on the first request.
func NewHandlers(cfg HandlerConfig) (*Handlers, error) {
	if cfg.Store == nil {
		return nil, errors.New("auth: nil session store")
	}
	if cfg.Signer == nil {
		return nil, errors.New("auth: nil jwt signer")
	}
	if strings.TrimSpace(cfg.Domain) == "" {
		return nil, errors.New("auth: domain required")
	}
	if cfg.ChainID <= 0 {
		return nil, errors.New("auth: chain id required")
	}
	if cfg.NonceTTL <= 0 {
		cfg.NonceTTL = DefaultNonceTTL
	}
	return &Handlers{cfg: cfg, now: time.Now}, nil
}

// nonceResponse is the public shape of POST /auth/siwe/nonce. We deliberately
// echo back nothing else; the client already knows the domain/chain it's
// targeting.
type nonceResponse struct {
	Nonce     string `json:"nonce"`
	ExpiresAt int64  `json:"expires_at"`
}

// errorBody is the canonical error envelope used across the auth endpoints.
// Codes are stable identifiers ("invalid_signature", "rate_limited", …) so
// clients can branch on them without parsing English; messages are short and
// MUST NOT include any value derived from the request (signature bytes,
// recovered address, internal redis errors).
type errorBody struct {
	Error errorPayload `json:"error"`
}

type errorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeError(c *gin.Context, status int, code, msg string) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.AbortWithStatusJSON(status, errorBody{Error: errorPayload{Code: code, Message: msg}})
}

// Nonce issues a fresh nonce and stores it in the session store with the
// configured TTL. The endpoint is intended to be hit unauthenticated, so
// callers SHOULD wire it behind the existing rate-limit middleware.
func (h *Handlers) Nonce(c *gin.Context) {
	nonce, err := RandomNonce(nonceByteLen)
	if err != nil {
		// Random read failure is an internal fault; surface a generic 500
		// without including the underlying error string.
		writeError(c, http.StatusInternalServerError, "internal_error", "failed to issue nonce")
		return
	}
	if err := h.cfg.Store.PutNonce(c.Request.Context(), nonce, h.cfg.NonceTTL); err != nil {
		writeError(c, http.StatusInternalServerError, "internal_error", "failed to issue nonce")
		return
	}
	c.JSON(http.StatusOK, nonceResponse{
		Nonce:     nonce,
		ExpiresAt: h.now().Add(h.cfg.NonceTTL).Unix(),
	})
}

// verifyRequest is the wire shape of POST /auth/siwe/verify. We accept the raw
// SIWE message string verbatim and let siwe-go handle parsing — we do NOT try
// to reconstruct it from individual fields, which would create a second source
// of truth and a reliable channel for clients to confuse the verifier.
type verifyRequest struct {
	Message   string `json:"message"   binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

// verifyResponse is the success shape of POST /auth/siwe/verify. We return the
// JWT and the lowercased address; the JTI is intentionally NOT returned —
// clients have no business with it, and exposing it grows the logout API
// surface.
type verifyResponse struct {
	Token     string `json:"token"`
	Address   string `json:"address"`
	ExpiresAt int64  `json:"expires_at"`
}

// Verify parses, validates, consumes the nonce, recovers the signer and (on
// success) issues the JWT + Redis session. The order matters: we consume the
// nonce only AFTER static validations pass so a malformed message does not
// cost the caller their nonce, but we do consume it BEFORE sig recovery so a
// successful recovery against a stale-but-not-yet-deleted nonce cannot mint
// two sessions. Concretely:
//
//  1. Decode JSON body.
//  2. Parse SIWE message (siwe-go).
//  3. Static checks: domain, chain id, time window.
//  4. Atomic ConsumeNonce — single-use guarantee.
//  5. ECDSA signature recovery + EIP-55 address match.
//  6. Mint JWT, persist session, return.
func (h *Handlers) Verify(c *gin.Context) {
	var req verifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	msg, err := siwe.ParseMessage(req.Message)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid_message", "could not parse SIWE message")
		return
	}

	// 3a. Domain pin: the message MUST be addressed to us.
	if msg.GetDomain() != h.cfg.Domain {
		writeError(c, http.StatusBadRequest, "domain_mismatch", "message domain mismatch")
		return
	}
	// 3b. Chain pin: the message MUST be on our chain.
	if msg.GetChainID() != h.cfg.ChainID {
		writeError(c, http.StatusBadRequest, "chain_mismatch", "message chain id mismatch")
		return
	}
	// 3c. Time window: covers expirationTime, notBefore.
	if _, err := msg.ValidAt(h.now()); err != nil {
		writeError(c, http.StatusBadRequest, "expired_message", "message outside validity window")
		return
	}

	// 4. Atomic single-use consume. If this fails, the nonce was either
	//    forged, expired, or already used. We cannot distinguish — and we
	//    must not — so the response is generic.
	if err := h.cfg.Store.ConsumeNonce(c.Request.Context(), msg.GetNonce()); err != nil {
		writeError(c, http.StatusUnauthorized, "invalid_nonce", "invalid or expired nonce")
		return
	}

	// 5. Signature recovery. siwe-go.VerifyEIP191 already does a constant-time
	//    byte compare on the recovered address vs the message address; we
	//    don't re-check it here. If recovery fails OR the address doesn't
	//    match, we surface a single generic code so an attacker cannot tell
	//    "valid sig wrong signer" from "invalid sig" — that distinction is
	//    a side channel.
	if _, err := msg.VerifyEIP191(req.Signature); err != nil {
		writeError(c, http.StatusUnauthorized, "invalid_signature", "signature did not verify")
		return
	}

	// 6. Mint + persist. address is taken FROM the message (which is now
	//    proven to be signed by it) and lowercased so claim/session keys are
	//    canonical.
	address := strings.ToLower(msg.GetAddress().Hex())
	token, jti, err := h.cfg.Signer.Sign(address)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "internal_error", "failed to issue session")
		return
	}
	if err := h.cfg.Store.PutSession(c.Request.Context(), jti, address, h.cfg.Signer.ttl); err != nil {
		writeError(c, http.StatusInternalServerError, "internal_error", "failed to persist session")
		return
	}

	c.JSON(http.StatusOK, verifyResponse{
		Token:     token,
		Address:   address,
		ExpiresAt: h.now().Add(h.cfg.Signer.ttl).Unix(),
	})
}

// Logout invalidates the session associated with the bearer token. It is
// idempotent: repeated logouts and logouts of unknown tokens both return 204
// rather than leaking whether the session ever existed.
//
// We require a syntactically valid token (signature, alg, exp) to find the
// JTI we want to delete; we do NOT require SessionExists, because a logout
// against an already-expired session should still report success.
func (h *Handlers) Logout(c *gin.Context) {
	tok := bearerFromHeader(c.GetHeader("Authorization"))
	if tok == "" {
		// No token — nothing to do, but report success so we don't tell
		// anonymous probers anything about how the endpoint validates.
		c.Status(http.StatusNoContent)
		return
	}
	claims, err := h.cfg.Signer.Parse(tok)
	if err != nil {
		// Same response as the no-token case — the response code is the
		// only public signal, and it is invariant.
		c.Status(http.StatusNoContent)
		return
	}
	if err := h.cfg.Store.DeleteSession(c.Request.Context(), claims.ID); err != nil {
		writeError(c, http.StatusInternalServerError, "internal_error", "failed to logout")
		return
	}
	c.Status(http.StatusNoContent)
}

// bearerFromHeader extracts the bearer token from an Authorization header
// value. We accept a single canonical form — "Bearer <token>" — and refuse
// anything else; this is intentionally stricter than RFC 6750 to keep the
// parsing surface small.
func bearerFromHeader(h string) string {
	const prefix = "Bearer "
	h = strings.TrimSpace(h)
	if len(h) <= len(prefix) {
		return ""
	}
	if !strings.EqualFold(h[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(h[len(prefix):])
}
