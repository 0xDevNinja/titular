package openapi

import (
	"time"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/store"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/types"
)

// This file holds the named response shapes consumed by the swag
// annotation comments on the M4 REST handlers. Swag (1.16) cannot
// resolve Go generics, so the generic store.Page[T] return types are
// re-declared here as explicit named structs. The on-the-wire JSON shape
// is identical — these are documentation-only aliases.
//
// IMPORTANT: keep these in sync with store.Page and the underlying
// element types. The CI drift check (gateway-openapi-check job) catches
// drift in the generated spec, but if you add a new field to e.g.
// store.Agent and forget to import it via this aggregation, the spec
// will silently drop the field. Add a test in api_swagger_types_test.go
// that asserts every json tag round-trips.

// AgentsPage is the response shape of GET /api/v1/agents.
type AgentsPage struct {
	// Data is the page contents in (block_number DESC, id DESC) order.
	Data []store.Agent `json:"data"`
	// NextCursor is the opaque pagination token. Empty on the final page.
	NextCursor string `json:"next_cursor,omitempty"`
}

// TradesPage is the response shape of GET /api/v1/trades.
type TradesPage struct {
	Data       []store.Trade `json:"data"`
	NextCursor string        `json:"next_cursor,omitempty"`
}

// JobsPage is the response shape of GET /api/v1/jobs.
type JobsPage struct {
	Data       []store.Job `json:"data"`
	NextCursor string      `json:"next_cursor,omitempty"`
}

// Stats is re-exported under handlers so the swag annotation can refer
// to it via the same package alias the rest of the M4 read endpoints use.
// The on-the-wire shape is identical to store.Stats.
type Stats = store.Stats

// HealthResponse is the shape of GET /healthz.
type HealthResponse struct {
	// Status is always "ok" when the gateway is serving.
	Status string `json:"status" example:"ok"`
}

// AuthNonceResponse is the shape of POST /auth/{siwe,siws}/nonce.
type AuthNonceResponse struct {
	// Nonce is a fresh single-use entropy token bound to the upcoming
	// /verify call. Hex-encoded random bytes.
	Nonce string `json:"nonce"`
	// ExpiresAt is the unix timestamp at which the nonce will be
	// rejected by /verify even if signed correctly.
	ExpiresAt int64 `json:"expires_at" example:"1714512000"`
}

// AuthVerifyRequest is the shape of POST /auth/siwe/verify.
//
// SIWS uses an identical wire shape (just a different signature codec on
// the body), so the same struct is referenced from both annotations.
type AuthVerifyRequest struct {
	// Message is the EIP-4361 (SIWE) or analogous SIWS message string,
	// verbatim — clients MUST NOT pre-parse and re-serialise it.
	Message string `json:"message"`
	// Signature is hex-encoded for SIWE (0x-prefixed, 65 bytes) and
	// base58 for SIWS (64-byte ed25519 signature).
	Signature string `json:"signature"`
}

// AuthVerifyResponse is the shape of POST /auth/{siwe,siws}/verify.
//
// The string fields here are intentionally polymorphic and MUST stay
// un-patched in the generated OpenAPI spec (no format / pattern):
//
//   - Token is a JWT bearer string, NOT an EVM address. It happens to
//     share the JSON name "token" with store.Job.token (an ERC-20
//     contract address), but the wire shape is completely different.
//   - Address is a lowercase 0x-hex EVM address when the token was
//     issued by /auth/siwe/verify, or a base58-encoded Solana pubkey
//     when issued by /auth/siws/verify. The EVM "address" format
//     (^0x[a-fA-F0-9]{40}$) would reject every valid SIWS response.
//
// scripts/openapi-gen scopes its scalar-format patcher by (schemaName,
// fieldName) precisely to keep this struct's fields as plain strings.
// See TestSpecLeavesAuthVerifyTokenUnformatted /
// TestSpecLeavesAuthVerifyAddressUnformatted for the regression guards.
type AuthVerifyResponse struct {
	// Token is the bearer JWT to attach as Authorization: Bearer <token>.
	Token string `json:"token"`
	// Address is the lowercase 0x-hex EVM address (SIWE) or base58
	// Solana pubkey (SIWS) the token was issued to.
	Address string `json:"address"`
	// ExpiresAt is the unix timestamp the token expires at.
	ExpiresAt int64 `json:"expires_at" example:"1714598400"`
}

// AuthErrorResponse is the canonical error envelope returned by every
// auth endpoint (and by the M4 REST endpoints on validation failure).
type AuthErrorResponse struct {
	Error AuthErrorBody `json:"error"`
}

// AuthErrorBody is the inner envelope.
type AuthErrorBody struct {
	// Code is a stable machine-readable identifier such as
	// "invalid_signature", "domain_mismatch", "rate_limited".
	Code string `json:"code"  example:"invalid_signature"`
	// Message is a short human-readable hint. MUST NOT contain values
	// derived from the request body (signature bytes, recovered
	// address, internal errors).
	Message string `json:"message" example:"signature did not verify"`
}

// timeAlias is the placeholder swag uses to render time.Time fields as
// RFC 3339 strings in the spec. Forcing the dependency keeps the import
// from being trimmed in environments where the rest of the file does
// not pull in time; see init below.
var _ = time.Time{}

// _ is a compile-time guard that the response types still match the
// underlying generic shapes. If store.Page changes its public field
// names, this assertion will fail to compile — forcing a doc update.
var (
	_ = AgentsPage{Data: []store.Agent{}, NextCursor: ""}
	_ = TradesPage{Data: []store.Trade{}, NextCursor: ""}
	_ = JobsPage{Data: []store.Job{}, NextCursor: ""}
	_ = AuthErrorResponse{Error: AuthErrorBody{}}
	_ = types.ErrorEnvelope{} // keep types import live for related shapes
)
