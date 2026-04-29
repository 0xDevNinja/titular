// Package store defines the read-only data access layer used by the gateway's
// REST API. The interface is decoupled from the wire types in
// internal/types so the same handlers can be backed by a real Postgres
// connection (cmd/gateway in production) or by an in-memory fake (handler
// unit tests).
//
// Pagination is opaque-cursor based: callers receive a NextCursor string and
// pass it back unmodified. Internally we encode (id, created_at) so the
// indexer's primary-key ordering survives even when multiple rows share the
// same created_at timestamp.
//
// All identifiers exposed to the wire are normalised to lowercase 0x-prefixed
// hex (addresses) or decimal strings (amounts) so clients never need to
// understand the database's compact bytea/numeric encoding.
package store

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ErrNotFound is returned by GetAgent when no row matches the requested id.
//
// Implementations MUST return this sentinel rather than a wrapped pgx.ErrNoRows
// so callers can branch on it without importing pgx.
var ErrNotFound = errors.New("store: not found")

// AgentKind mirrors the indexer's agent_kind enum.
type AgentKind string

const (
	AgentKindLaunchpad AgentKind = "launchpad"
	AgentKindACP       AgentKind = "acp"
)

// Valid reports whether k is one of the known enum values. Used to reject
// untrusted query input before it ever reaches a SQL prepare.
func (k AgentKind) Valid() bool {
	return k == AgentKindLaunchpad || k == AgentKindACP
}

// JobPhase mirrors the indexer's job_phase enum.
type JobPhase string

const (
	JobPhaseCreated   JobPhase = "created"
	JobPhaseFunded    JobPhase = "funded"
	JobPhaseActive    JobPhase = "active"
	JobPhaseCompleted JobPhase = "completed"
	JobPhaseCancelled JobPhase = "cancelled"
	JobPhaseDisputed  JobPhase = "disputed"
	JobPhaseReleased  JobPhase = "released"
	JobPhaseResolved  JobPhase = "resolved"
)

// Valid reports whether p is one of the known enum values.
func (p JobPhase) Valid() bool {
	switch p {
	case JobPhaseCreated, JobPhaseFunded, JobPhaseActive, JobPhaseCompleted,
		JobPhaseCancelled, JobPhaseDisputed, JobPhaseReleased, JobPhaseResolved:
		return true
	}
	return false
}

// Agent is the wire-shape of a registered agent (launchpad or ACP).
//
// Numeric ids are encoded as decimal strings because they originate as
// uint256 on chain and may exceed int64.
type Agent struct {
	ID           int64     `json:"id"`
	Kind         AgentKind `json:"kind"`
	AgentID      string    `json:"agent_id,omitempty"`     // on-chain uint256, decimal
	Controller   string    `json:"controller,omitempty"`   // 0x-hex (ACP)
	Creator      string    `json:"creator,omitempty"`      // 0x-hex (launchpad)
	MetadataURI  string    `json:"metadata_uri,omitempty"` // optional
	Capabilities string    `json:"capabilities,omitempty"` // bitmap as decimal
	BlockNumber  uint64    `json:"block_number"`
	TxHash       string    `json:"tx_hash"` // 0x-hex
	LogIndex     uint32    `json:"log_index"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Trade is the wire-shape of a bonding-curve buy or sell.
type Trade struct {
	ID           int64     `json:"id"`
	AgentTokenPK *int64    `json:"agent_token_pk,omitempty"`
	Side         string    `json:"side"` // "buy" | "sell"
	Trader       string    `json:"trader"`
	Curve        string    `json:"curve"`
	QuoteIn      string    `json:"quote_in,omitempty"`
	AgentOut     string    `json:"agent_out,omitempty"`
	AgentIn      string    `json:"agent_in,omitempty"`
	QuoteOut     string    `json:"quote_out,omitempty"`
	Fee          string    `json:"fee"`
	BlockNumber  uint64    `json:"block_number"`
	TxHash       string    `json:"tx_hash"`
	LogIndex     uint32    `json:"log_index"`
	CreatedAt    time.Time `json:"created_at"`
}

// Job is the wire-shape of an ACP job.
type Job struct {
	ID           int64     `json:"id"`
	OnChainID    string    `json:"on_chain_id"`
	CloneAddress string    `json:"clone_address,omitempty"`
	AgentPK      *int64    `json:"agent_pk,omitempty"`
	AgentID      string    `json:"agent_id,omitempty"` // on-chain uint256
	Principal    string    `json:"principal"`
	JobType      int16     `json:"job_type"`
	Token        string    `json:"token,omitempty"`
	Budget       string    `json:"budget"`
	Deadline     *time.Time `json:"deadline,omitempty"`
	Phase        JobPhase  `json:"phase"`
	ResultURI    string    `json:"result_uri,omitempty"`
	BlockNumber  uint64    `json:"block_number"`
	TxHash       string    `json:"tx_hash"`
	LogIndex     uint32    `json:"log_index"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Stats is the aggregate snapshot exposed by GET /api/v1/stats.
//
// LastBlockIndexed is the highest block_number observed across the trades and
// jobs tables (the two tables that monotonically advance with chain height).
// We deliberately do not source this from a "checkpoint" table because the
// indexer track does not expose one yet; this query is a few-millisecond
// MAX() and can be replaced once a checkpoint exists.
type Stats struct {
	TotalAgents      int64 `json:"total_agents"`
	TotalJobs        int64 `json:"total_jobs"`
	TotalTrades      int64 `json:"total_trades"`
	LastBlockIndexed int64 `json:"last_block_indexed"`
}

// ListAgentsParams is the input to Store.ListAgents.
//
// Limit is clamped server-side; callers should still pass a sane value to keep
// the database from over-fetching. Cursor, when non-empty, must have been
// produced by a previous call to Store.ListAgents — its format is opaque to
// callers and intentionally unstable across releases.
type ListAgentsParams struct {
	Kind   AgentKind // empty = no filter
	Limit  int
	Cursor string
}

// ListTradesParams is the input to Store.ListTrades.
//
// AgentToken filters to trades for a specific agent token (matched by 20-byte
// 0x-prefixed token address, case-insensitive). From / To are inclusive
// bounds on the trade's created_at column; either may be zero to disable.
type ListTradesParams struct {
	AgentToken string // 0x-hex address; empty = no filter
	From       time.Time
	To         time.Time
	Limit      int
	Cursor     string
}

// ListJobsParams is the input to Store.ListJobs.
type ListJobsParams struct {
	Status JobPhase // empty = no filter
	Limit  int
	Cursor string
}

// Page is the standard return shape for a paginated list endpoint. NextCursor
// is empty when the current page is the final page.
type Page[T any] struct {
	Data       []T    `json:"data"`
	NextCursor string `json:"next_cursor,omitempty"`
}

// Store is the read-only data access surface used by the REST API.
//
// All methods accept a context.Context for cancellation; callers should
// propagate the request context so a client disconnect cancels in-flight
// queries (pgx honours context cancellation by closing the conn).
type Store interface {
	ListAgents(ctx context.Context, p ListAgentsParams) (Page[Agent], error)
	GetAgent(ctx context.Context, idOrSlug string) (Agent, error)
	ListTrades(ctx context.Context, p ListTradesParams) (Page[Trade], error)
	ListJobs(ctx context.Context, p ListJobsParams) (Page[Job], error)
	Stats(ctx context.Context) (Stats, error)
}

// ---------------------------------------------------------------------------
// Pagination helpers — exported so the in-memory fake store and the Postgres
// store agree on cursor format. Tests rely on this to round-trip.
// ---------------------------------------------------------------------------

// DefaultLimit is used when the caller does not provide a limit.
const DefaultLimit = 50

// MaxLimit caps the page size that any caller can request. The handler layer
// rejects larger values with 400; the store also clamps defensively so a
// programming error in the handler can never blow up the database.
const MaxLimit = 200

// ClampLimit returns a sane page size: zero / negative becomes DefaultLimit;
// values above MaxLimit are reduced to MaxLimit.
func ClampLimit(n int) int {
	if n <= 0 {
		return DefaultLimit
	}
	if n > MaxLimit {
		return MaxLimit
	}
	return n
}

// Cursor is the structured cursor payload that gets base64-encoded as the
// opaque NextCursor token. Encoding `id` keeps ordering stable when many rows
// share the same CreatedAt; encoding CreatedAt keeps the descending scan
// efficient against the (block_number DESC) / (created_at DESC) indexes the
// indexer schema provides.
type Cursor struct {
	ID        int64     `json:"i"`
	CreatedAt time.Time `json:"t"`
}

// EncodeCursor returns the base64-encoded JSON form of c.
func EncodeCursor(c Cursor) string {
	b, _ := json.Marshal(c) // marshalling a fixed struct cannot fail
	return base64.RawURLEncoding.EncodeToString(b)
}

// DecodeCursor parses an opaque cursor string. Returns an error for unknown
// or malformed input — handlers MUST surface this as 400 to the caller, never
// silently fall back to "first page", because that would let an attacker
// brute-force pagination by submitting empty cursors.
func DecodeCursor(s string) (Cursor, error) {
	if s == "" {
		return Cursor{}, nil
	}
	raw, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return Cursor{}, fmt.Errorf("cursor is not valid base64: %w", err)
	}
	var c Cursor
	if err := json.Unmarshal(raw, &c); err != nil {
		return Cursor{}, fmt.Errorf("cursor payload malformed: %w", err)
	}
	if c.ID < 0 {
		return Cursor{}, errors.New("cursor id must be non-negative")
	}
	return c, nil
}
