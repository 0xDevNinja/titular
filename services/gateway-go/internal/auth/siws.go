package auth

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mr-tron/base58"
)

// SIWS (Sign-In With Solana) parallels the SIWE flow but with ed25519
// signatures and base58 Solana addresses. The wire format is intentionally
// kept symmetric to EIP-4361 so a single client implementation can target
// both chains by swapping the cluster/chain identifier.
//
// We do NOT pull in the heavy github.com/gagliardetto/solana-go SDK: the
// auth path needs only base58 + ed25519, both of which are available from
// stdlib (crypto/ed25519, which is constant-time) and a tiny base58
// package (github.com/mr-tron/base58, the canonical Solana base58 lib).
// Adding the full SDK would drag in dozens of indirect deps for zero
// security benefit.

// Cluster identifies a Solana network. We pin it server-side (parallel to
// SIWE's chain id) so a message signed against devnet cannot be replayed
// against mainnet, or vice versa.
const (
	ClusterDevnet      = "devnet"
	ClusterMainnetBeta = "mainnet-beta"
	ClusterTestnet     = "testnet"
)

// solanaPubKeyLen is the byte length of an ed25519 public key, which is
// also the byte length of a Solana wallet address before base58 encoding.
const solanaPubKeyLen = ed25519.PublicKeySize // 32

// SIWSHandlerConfig holds the request-time dependencies for the SIWS
// handlers. It mirrors HandlerConfig so the wiring code in cmd/gateway is
// symmetric.
//
// Domain is the value the SIWS message MUST declare. We pin it server-side
// rather than echoing whatever the message contained so a SIWS message
// signed for some-other-app.example cannot be replayed against this
// gateway.
//
// Cluster is the Solana network the message MUST declare. Devnet messages
// MUST NOT validate against a mainnet deployment, and vice versa.
//
// NonceTTL controls how long an issued nonce stays valid. Defaults to
// DefaultNonceTTL when zero.
//
// Store / Signer are reused from the SIWE flow so a logout against either
// path invalidates the same JWT, and so RequireAuth can stay
// chain-agnostic.
type SIWSHandlerConfig struct {
	Store    SessionStore
	Signer   *JWTSigner
	Domain   string
	Cluster  string
	NonceTTL time.Duration
}

// SIWSHandlers exposes the SIWS HTTP handlers as Gin handler funcs.
// Construct via NewSIWSHandlers — the zero value is unsafe.
type SIWSHandlers struct {
	cfg SIWSHandlerConfig
	now func() time.Time // overridable in tests
}

// NewSIWSHandlers wires the handler dependencies after validating that all
// required fields are present. Failing here at startup is preferable to
// failing on the first request.
func NewSIWSHandlers(cfg SIWSHandlerConfig) (*SIWSHandlers, error) {
	if cfg.Store == nil {
		return nil, errors.New("auth: nil session store")
	}
	if cfg.Signer == nil {
		return nil, errors.New("auth: nil jwt signer")
	}
	if strings.TrimSpace(cfg.Domain) == "" {
		return nil, errors.New("auth: domain required")
	}
	switch cfg.Cluster {
	case ClusterDevnet, ClusterMainnetBeta, ClusterTestnet:
		// ok
	default:
		return nil, fmt.Errorf("auth: cluster must be one of %q/%q/%q, got %q",
			ClusterDevnet, ClusterMainnetBeta, ClusterTestnet, cfg.Cluster)
	}
	if cfg.NonceTTL <= 0 {
		cfg.NonceTTL = DefaultNonceTTL
	}
	return &SIWSHandlers{cfg: cfg, now: time.Now}, nil
}

// Nonce issues a fresh nonce and stores it in the session store with the
// configured TTL. The endpoint is intended to be hit unauthenticated, so
// callers SHOULD wire it behind the existing rate-limit middleware.
//
// The nonce namespace is shared with SIWE: a single nonce issued via
// /auth/siws/nonce can only be consumed once across BOTH paths. This is
// the right call because the nonce is just an entropy token; the chain
// binding happens via the signed message, not the nonce.
// @Summary      Issue SIWS nonce
// @Description  Returns a fresh entropy nonce the client embeds in the
// @Description  SIWS (Sign-In With Solana) message it asks the user to
// @Description  sign. The nonce is single-use across both SIWE and SIWS
// @Description  verify endpoints and expires after a short server-controlled
// @Description  TTL.
// @Tags         auth
// @Produce      json
// @Success      200  {object}  openapi.AuthNonceResponse
// @Failure      500  {object}  openapi.AuthErrorResponse
// @Router       /auth/siws/nonce [post]
func (h *SIWSHandlers) Nonce(c *gin.Context) {
	nonce, err := RandomNonce(nonceByteLen)
	if err != nil {
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

// siwsVerifyRequest is the wire shape of POST /auth/siws/verify. We accept
// the raw SIWS message string verbatim and parse it server-side; clients
// MUST NOT send the parsed fields, otherwise we'd have two sources of
// truth and a reliable channel for clients to confuse the verifier.
//
// Signature is a base58-encoded ed25519 signature (64 bytes). Phantom and
// Backpack both return base58 from signMessage.
type siwsVerifyRequest struct {
	Message   string `json:"message"   binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

// Verify parses, validates, consumes the nonce, recovers the signer and
// (on success) issues the JWT + Redis session. The order matches the SIWE
// flow:
//
//  1. Decode JSON body.
//  2. Parse SIWS message.
//  3. Static checks: domain, cluster, time window, address well-formed.
//  4. Atomic ConsumeNonce — single-use guarantee.
//  5. ed25519 verify (constant-time, stdlib).
//  6. Mint JWT, persist session, return.
//
// On step 5 we DO consume the nonce before signature verification (same
// rationale as SIWE: a successful recovery against a stale-but-not-yet
// deleted nonce cannot mint two sessions).
// @Summary      Verify SIWS message and mint session
// @Description  Validates a SIWS (Sign-In With Solana) message and its
// @Description  ed25519 signature against the previously issued nonce,
// @Description  the gateway's pinned domain, and the configured Solana
// @Description  cluster; on success, returns a bearer JWT and persists
// @Description  the session in Redis.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      openapi.AuthVerifyRequest  true  "SIWS message and base58 signature"
// @Success      200   {object}  openapi.AuthVerifyResponse
// @Failure      400   {object}  openapi.AuthErrorResponse  "invalid_request, invalid_message, domain_mismatch, cluster_mismatch, expired_message, invalid_address"
// @Failure      401   {object}  openapi.AuthErrorResponse  "invalid_nonce or invalid_signature"
// @Failure      500   {object}  openapi.AuthErrorResponse  "internal_error"
// @Router       /auth/siws/verify [post]
func (h *SIWSHandlers) Verify(c *gin.Context) {
	var req siwsVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	msg, err := ParseSIWSMessage(req.Message)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid_message", "could not parse SIWS message")
		return
	}

	// 3a. Domain pin: the message MUST be addressed to us.
	if msg.Domain != h.cfg.Domain {
		writeError(c, http.StatusBadRequest, "domain_mismatch", "message domain mismatch")
		return
	}
	// 3b. Cluster pin: the message MUST be on our network.
	if msg.Cluster != h.cfg.Cluster {
		writeError(c, http.StatusBadRequest, "cluster_mismatch", "message cluster mismatch")
		return
	}
	// 3c. Time window: covers expirationTime (and notBefore if present).
	if err := msg.ValidAt(h.now()); err != nil {
		writeError(c, http.StatusBadRequest, "expired_message", "message outside validity window")
		return
	}

	// 3d. Decode the address eagerly so a malformed address never reaches
	// nonce consumption. base58 of an ed25519 pubkey is exactly 32 bytes.
	pubkey, err := base58.Decode(msg.Address)
	if err != nil || len(pubkey) != solanaPubKeyLen {
		writeError(c, http.StatusBadRequest, "invalid_address", "address is not a valid Solana pubkey")
		return
	}

	// 4. Atomic single-use consume. If this fails the nonce was either
	//    forged, expired, or already used. We cannot — and must not —
	//    distinguish, so the response is generic.
	if err := h.cfg.Store.ConsumeNonce(c.Request.Context(), msg.Nonce); err != nil {
		writeError(c, http.StatusUnauthorized, "invalid_nonce", "invalid or expired nonce")
		return
	}

	// 5. Signature verify. ed25519.Verify is constant-time per the stdlib
	//    contract. We surface a single generic code on failure so the
	//    "valid sig wrong signer" / "invalid sig" distinction is not a
	//    side channel.
	sig, err := base58.Decode(req.Signature)
	if err != nil || len(sig) != ed25519.SignatureSize {
		writeError(c, http.StatusUnauthorized, "invalid_signature", "signature did not verify")
		return
	}
	if !ed25519.Verify(ed25519.PublicKey(pubkey), []byte(req.Message), sig) {
		writeError(c, http.StatusUnauthorized, "invalid_signature", "signature did not verify")
		return
	}

	// 6. Mint + persist. Address is taken FROM the message (which is now
	//    proven to be signed by it). Solana addresses are case-sensitive
	//    base58 — we keep the original casing so claim/session keys round-
	//    trip exactly.
	address := msg.Address
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

// SIWSMessage is the parsed form of a SIWS message. Field semantics mirror
// EIP-4361 with one substitution: ChainID becomes Cluster ("devnet" /
// "mainnet-beta" / "testnet"), reflecting how Solana names its networks.
//
// We hold the raw text (Raw) so the verifier signs exactly the bytes that
// were sent — there is no field-by-field re-serialisation step that could
// canonicalise differences between parser and signer.
type SIWSMessage struct {
	Domain         string
	Address        string
	Statement      string
	URI            string
	Version        string
	Cluster        string
	Nonce          string
	IssuedAt       time.Time
	ExpirationTime time.Time // zero if absent
	NotBefore      time.Time // zero if absent
	RequestID      string    // optional
	Raw            string
}

// ValidAt returns nil iff the message's time window includes t. A message
// without expiration and not-before fields is considered always-valid by
// time alone (the nonce TTL is the real upper bound).
func (m *SIWSMessage) ValidAt(t time.Time) error {
	if !m.NotBefore.IsZero() && t.Before(m.NotBefore) {
		return errors.New("message not yet valid")
	}
	if !m.ExpirationTime.IsZero() && !t.Before(m.ExpirationTime) {
		return errors.New("message expired")
	}
	return nil
}

// ParseSIWSMessage parses an EIP-4361-shaped message that uses Solana
// conventions. The grammar is intentionally strict — we refuse anything
// that does not match the canonical layout so a tolerant parser cannot be
// used as a confusion vector.
//
// Canonical layout (line-separated):
//
//	<domain> wants you to sign in with your Solana account:
//	<address>
//	[empty line]
//	[<statement>]
//	[empty line]
//	URI: <uri>
//	Version: <version>
//	Chain ID: <cluster>
//	Nonce: <nonce>
//	Issued At: <RFC3339 timestamp>
//	[Expiration Time: <RFC3339 timestamp>]
//	[Not Before: <RFC3339 timestamp>]
//	[Request ID: <opaque>]
//
// We tolerate the cluster being declared as either "Chain ID" (matching
// the SIWE field name; what most port-of-SIWE clients emit) or the
// alternate "Cluster" key. We do NOT accept arbitrary trailing fields:
// resources/scopes are explicitly out of scope until we have a use case.
func ParseSIWSMessage(raw string) (*SIWSMessage, error) {
	if raw == "" {
		return nil, errors.New("empty message")
	}
	lines := strings.Split(raw, "\n")
	if len(lines) < 7 {
		return nil, errors.New("message too short")
	}

	m := &SIWSMessage{Raw: raw}

	// Header: "<domain> wants you to sign in with your Solana account:"
	const headerSuffix = " wants you to sign in with your Solana account:"
	header := lines[0]
	if !strings.HasSuffix(header, headerSuffix) {
		return nil, errors.New("missing or malformed header")
	}
	m.Domain = strings.TrimSuffix(header, headerSuffix)
	if m.Domain == "" {
		return nil, errors.New("empty domain")
	}

	// Address line.
	m.Address = strings.TrimSpace(lines[1])
	if m.Address == "" {
		return nil, errors.New("empty address")
	}

	// After the address there is an empty line. If a statement is
	// present, it follows; if absent, the URI block follows directly.
	idx := 2
	if idx >= len(lines) || lines[idx] != "" {
		return nil, errors.New("missing blank line after address")
	}
	idx++
	// Optional statement: a single line followed by a blank line. We
	// detect by looking ahead for the URI marker; if the next line starts
	// with "URI:" there is no statement.
	if idx < len(lines) && !strings.HasPrefix(lines[idx], "URI:") {
		m.Statement = lines[idx]
		idx++
		if idx >= len(lines) || lines[idx] != "" {
			return nil, errors.New("missing blank line after statement")
		}
		idx++
	}

	// From here we expect the keyed block. Each remaining line must be
	// "<key>: <value>" or empty (terminator). Unknown keys are refused.
	seen := map[string]bool{}
	for ; idx < len(lines); idx++ {
		ln := lines[idx]
		if ln == "" {
			continue
		}
		colon := strings.Index(ln, ":")
		if colon < 0 {
			return nil, fmt.Errorf("line %d missing colon", idx)
		}
		key := strings.TrimSpace(ln[:colon])
		val := strings.TrimSpace(ln[colon+1:])
		if seen[key] {
			return nil, fmt.Errorf("duplicate field %q", key)
		}
		seen[key] = true
		switch key {
		case "URI":
			m.URI = val
		case "Version":
			m.Version = val
		case "Chain ID", "Cluster":
			// Either spelling lands in the same field — see grammar
			// note above. Refuse if both appear.
			if m.Cluster != "" {
				return nil, errors.New("cluster declared twice")
			}
			m.Cluster = val
		case "Nonce":
			m.Nonce = val
		case "Issued At":
			t, err := time.Parse(time.RFC3339, val)
			if err != nil {
				return nil, fmt.Errorf("issued at: %w", err)
			}
			m.IssuedAt = t
		case "Expiration Time":
			t, err := time.Parse(time.RFC3339, val)
			if err != nil {
				return nil, fmt.Errorf("expiration time: %w", err)
			}
			m.ExpirationTime = t
		case "Not Before":
			t, err := time.Parse(time.RFC3339, val)
			if err != nil {
				return nil, fmt.Errorf("not before: %w", err)
			}
			m.NotBefore = t
		case "Request ID":
			m.RequestID = val
		default:
			return nil, fmt.Errorf("unknown field %q", key)
		}
	}

	// Required-field invariants.
	if m.URI == "" {
		return nil, errors.New("missing URI")
	}
	if m.Version == "" {
		return nil, errors.New("missing Version")
	}
	if m.Cluster == "" {
		return nil, errors.New("missing Chain ID/Cluster")
	}
	if m.Nonce == "" {
		return nil, errors.New("missing Nonce")
	}
	if m.IssuedAt.IsZero() {
		return nil, errors.New("missing Issued At")
	}
	// Defence-in-depth: refuse a nonce containing anything that isn't a
	// printable ASCII byte. Our own nonces are hex; a crafted message
	// with control characters in the nonce should never round-trip
	// through Redis as a sane key.
	for i := 0; i < len(m.Nonce); i++ {
		c := m.Nonce[i]
		if c < 0x21 || c > 0x7e {
			return nil, errors.New("nonce contains non-printable byte")
		}
	}
	return m, nil
}
