package store

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Querier is the subset of *pgxpool.Pool we depend on. Splitting the
// interface keeps tests free of the full pool surface (Acquire, Stat,
// Reset, …) and makes it possible to mock with a single-method fake.
type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// PG is a read-only Postgres-backed Store.
//
// The pool is owned by the caller (cmd/gateway) so the same pool can be
// reused for any future read-only query layers (e.g. websocket subscriptions
// in M5). PG never issues writes — every method goes through SELECT only.
type PG struct {
	q Querier
}

// NewPG wires a PG store onto an existing pool.
func NewPG(q Querier) *PG { return &PG{q: q} }

// Connect opens a new pgxpool.Pool against dsn and returns a wired PG store
// plus the pool's Close hook so the caller can defer it. We expose the close
// separately rather than rolling it into PG so a single pool can back many
// stores in the future without one closing it out from under another.
func Connect(ctx context.Context, dsn string) (*PG, func(), error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("parse dsn: %w", err)
	}
	// Modest pool — the gateway is read-mostly. Operators can tune via DSN
	// parameters if they need more.
	if cfg.MaxConns < 4 {
		cfg.MaxConns = 4
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("connect: %w", err)
	}
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("ping: %w", err)
	}
	return NewPG(pool), pool.Close, nil
}

// ---------------------------------------------------------------------------
// ListAgents
// ---------------------------------------------------------------------------

// agent_id, capabilities are uint256s. We project them into text via
// "(col)::text" because pgx's pgtype.Numeric scan-into-string path requires
// us to read into a Numeric and stringify ourselves; ::text is shorter and
// semantically identical for non-NaN values, which numeric(78,0) cannot have.
const agentSelectCols = `
	id,
	kind,
	(agent_id)::text     AS agent_id_text,
	controller,
	creator,
	metadata_uri,
	(capabilities)::text AS capabilities_text,
	block_number,
	tx_hash,
	log_index,
	created_at,
	updated_at`

// ListAgents returns a page of agents, optionally filtered by kind. Newest
// first, ordered by (created_at DESC, id DESC) so the cursor is monotone.
func (s *PG) ListAgents(ctx context.Context, p ListAgentsParams) (Page[Agent], error) {
	limit := ClampLimit(p.Limit)
	cur, err := DecodeCursor(p.Cursor)
	if err != nil {
		return Page[Agent]{}, err
	}

	// Build the WHERE clause from explicitly-typed params. Every value is a
	// $N placeholder — never string-concatenated — so no caller input ever
	// reaches the SQL parser unescaped.
	var (
		conds []string
		args  []any
	)
	if p.Kind != "" {
		if !p.Kind.Valid() {
			return Page[Agent]{}, fmt.Errorf("invalid kind %q", p.Kind)
		}
		args = append(args, string(p.Kind))
		conds = append(conds, fmt.Sprintf("kind = $%d::agent_kind", len(args)))
	}
	if !cur.CreatedAt.IsZero() {
		args = append(args, cur.CreatedAt, cur.ID)
		conds = append(conds, fmt.Sprintf("(created_at, id) < ($%d, $%d)", len(args)-1, len(args)))
	}
	args = append(args, limit+1) // fetch one extra to know if there's another page

	q := "SELECT" + agentSelectCols + " FROM agents"
	if len(conds) > 0 {
		q += " WHERE " + strings.Join(conds, " AND ")
	}
	q += fmt.Sprintf(" ORDER BY created_at DESC, id DESC LIMIT $%d", len(args))

	rows, err := s.q.Query(ctx, q, args...)
	if err != nil {
		return Page[Agent]{}, err
	}
	defer rows.Close()

	out := make([]Agent, 0, limit)
	for rows.Next() {
		var (
			a            Agent
			agentIDText  *string
			capsText     *string
			controller   []byte
			creator      []byte
			metadataURI  *string
			txHash       []byte
			kindEnum     string
		)
		if err := rows.Scan(
			&a.ID,
			&kindEnum,
			&agentIDText,
			&controller,
			&creator,
			&metadataURI,
			&capsText,
			&a.BlockNumber,
			&txHash,
			&a.LogIndex,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return Page[Agent]{}, err
		}
		a.Kind = AgentKind(kindEnum)
		if agentIDText != nil {
			a.AgentID = *agentIDText
		}
		if capsText != nil {
			a.Capabilities = *capsText
		}
		if metadataURI != nil {
			a.MetadataURI = *metadataURI
		}
		a.Controller = bytesToHex(controller)
		a.Creator = bytesToHex(creator)
		a.TxHash = bytesToHex(txHash)
		out = append(out, a)
	}
	if err := rows.Err(); err != nil {
		return Page[Agent]{}, err
	}

	page := Page[Agent]{Data: out}
	if len(out) > limit {
		last := out[limit-1]
		page.Data = out[:limit]
		page.NextCursor = EncodeCursor(Cursor{ID: last.ID, CreatedAt: last.CreatedAt})
	}
	return page, nil
}

// ---------------------------------------------------------------------------
// GetAgent
// ---------------------------------------------------------------------------

// GetAgent looks up a single agent by primary key OR by the kind:agent_id slug
// (e.g. "acp:42"). Returns ErrNotFound when no row matches.
func (s *PG) GetAgent(ctx context.Context, idOrSlug string) (Agent, error) {
	idOrSlug = strings.TrimSpace(idOrSlug)
	if idOrSlug == "" {
		return Agent{}, ErrNotFound
	}

	var (
		q    string
		args []any
	)

	if kind, agentID, ok := splitSlug(idOrSlug); ok {
		// Composite slug: filter on (kind, agent_id).
		bigID, ok := new(big.Int).SetString(agentID, 10)
		if !ok || bigID.Sign() < 0 {
			return Agent{}, ErrNotFound
		}
		var num pgtype.Numeric
		if err := num.Scan(bigID.String()); err != nil {
			return Agent{}, ErrNotFound
		}
		q = "SELECT" + agentSelectCols +
			" FROM agents WHERE kind = $1::agent_kind AND agent_id = $2 LIMIT 1"
		args = []any{string(kind), num}
	} else {
		// Numeric primary key.
		pk, err := strconv.ParseInt(idOrSlug, 10, 64)
		if err != nil || pk <= 0 {
			return Agent{}, ErrNotFound
		}
		q = "SELECT" + agentSelectCols + " FROM agents WHERE id = $1 LIMIT 1"
		args = []any{pk}
	}

	row := s.q.QueryRow(ctx, q, args...)
	var (
		a            Agent
		agentIDText  *string
		capsText     *string
		controller   []byte
		creator      []byte
		metadataURI  *string
		txHash       []byte
		kindEnum     string
	)
	if err := row.Scan(
		&a.ID,
		&kindEnum,
		&agentIDText,
		&controller,
		&creator,
		&metadataURI,
		&capsText,
		&a.BlockNumber,
		&txHash,
		&a.LogIndex,
		&a.CreatedAt,
		&a.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Agent{}, ErrNotFound
		}
		return Agent{}, err
	}
	a.Kind = AgentKind(kindEnum)
	if agentIDText != nil {
		a.AgentID = *agentIDText
	}
	if capsText != nil {
		a.Capabilities = *capsText
	}
	if metadataURI != nil {
		a.MetadataURI = *metadataURI
	}
	a.Controller = bytesToHex(controller)
	a.Creator = bytesToHex(creator)
	a.TxHash = bytesToHex(txHash)
	return a, nil
}

// ---------------------------------------------------------------------------
// ListTrades
// ---------------------------------------------------------------------------

const tradeSelectCols = `
	t.id,
	t.agent_token_pk,
	t.side,
	t.trader,
	t.curve,
	(t.quote_in)::text   AS quote_in_text,
	(t.agent_out)::text  AS agent_out_text,
	(t.agent_in)::text   AS agent_in_text,
	(t.quote_out)::text  AS quote_out_text,
	(t.fee)::text        AS fee_text,
	t.block_number,
	t.tx_hash,
	t.log_index,
	t.created_at`

// ListTrades returns a page of trades, optionally filtered by agent token
// address and / or a [from, to] created_at range.
func (s *PG) ListTrades(ctx context.Context, p ListTradesParams) (Page[Trade], error) {
	limit := ClampLimit(p.Limit)
	cur, err := DecodeCursor(p.Cursor)
	if err != nil {
		return Page[Trade]{}, err
	}

	var (
		conds []string
		args  []any
	)

	if p.AgentToken != "" {
		raw, err := decodeHexAddress(p.AgentToken)
		if err != nil {
			return Page[Trade]{}, fmt.Errorf("agent_token: %w", err)
		}
		args = append(args, raw)
		conds = append(conds, fmt.Sprintf(
			"t.agent_token_pk = (SELECT id FROM agent_tokens WHERE token = $%d LIMIT 1)",
			len(args),
		))
	}
	if !p.From.IsZero() {
		args = append(args, p.From)
		conds = append(conds, fmt.Sprintf("t.created_at >= $%d", len(args)))
	}
	if !p.To.IsZero() {
		args = append(args, p.To)
		conds = append(conds, fmt.Sprintf("t.created_at <= $%d", len(args)))
	}
	if !cur.CreatedAt.IsZero() {
		args = append(args, cur.CreatedAt, cur.ID)
		conds = append(conds, fmt.Sprintf("(t.created_at, t.id) < ($%d, $%d)", len(args)-1, len(args)))
	}
	args = append(args, limit+1)

	q := "SELECT" + tradeSelectCols + " FROM trades t"
	if len(conds) > 0 {
		q += " WHERE " + strings.Join(conds, " AND ")
	}
	q += fmt.Sprintf(" ORDER BY t.created_at DESC, t.id DESC LIMIT $%d", len(args))

	rows, err := s.q.Query(ctx, q, args...)
	if err != nil {
		return Page[Trade]{}, err
	}
	defer rows.Close()

	out := make([]Trade, 0, limit)
	for rows.Next() {
		var (
			t          Trade
			agentTokPK *int64
			quoteIn    *string
			agentOut   *string
			agentIn    *string
			quoteOut   *string
			fee        *string
			trader     []byte
			curve      []byte
			txHash     []byte
		)
		if err := rows.Scan(
			&t.ID,
			&agentTokPK,
			&t.Side,
			&trader,
			&curve,
			&quoteIn,
			&agentOut,
			&agentIn,
			&quoteOut,
			&fee,
			&t.BlockNumber,
			&txHash,
			&t.LogIndex,
			&t.CreatedAt,
		); err != nil {
			return Page[Trade]{}, err
		}
		t.AgentTokenPK = agentTokPK
		t.Trader = bytesToHex(trader)
		t.Curve = bytesToHex(curve)
		t.TxHash = bytesToHex(txHash)
		t.QuoteIn = derefString(quoteIn)
		t.AgentOut = derefString(agentOut)
		t.AgentIn = derefString(agentIn)
		t.QuoteOut = derefString(quoteOut)
		t.Fee = derefString(fee)
		out = append(out, t)
	}
	if err := rows.Err(); err != nil {
		return Page[Trade]{}, err
	}

	page := Page[Trade]{Data: out}
	if len(out) > limit {
		last := out[limit-1]
		page.Data = out[:limit]
		page.NextCursor = EncodeCursor(Cursor{ID: last.ID, CreatedAt: last.CreatedAt})
	}
	return page, nil
}

// ---------------------------------------------------------------------------
// ListJobs
// ---------------------------------------------------------------------------

const jobSelectCols = `
	id,
	(on_chain_id)::text  AS on_chain_id_text,
	clone_address,
	agent_pk,
	(agent_id)::text     AS agent_id_text,
	principal,
	job_type,
	token,
	(budget)::text       AS budget_text,
	deadline,
	phase,
	result_uri,
	block_number,
	tx_hash,
	log_index,
	created_at,
	updated_at`

// ListJobs returns a page of jobs, optionally filtered by phase.
func (s *PG) ListJobs(ctx context.Context, p ListJobsParams) (Page[Job], error) {
	limit := ClampLimit(p.Limit)
	cur, err := DecodeCursor(p.Cursor)
	if err != nil {
		return Page[Job]{}, err
	}

	var (
		conds []string
		args  []any
	)
	if p.Status != "" {
		if !p.Status.Valid() {
			return Page[Job]{}, fmt.Errorf("invalid status %q", p.Status)
		}
		args = append(args, string(p.Status))
		conds = append(conds, fmt.Sprintf("phase = $%d::job_phase", len(args)))
	}
	if !cur.CreatedAt.IsZero() {
		args = append(args, cur.CreatedAt, cur.ID)
		conds = append(conds, fmt.Sprintf("(created_at, id) < ($%d, $%d)", len(args)-1, len(args)))
	}
	args = append(args, limit+1)

	q := "SELECT" + jobSelectCols + " FROM jobs"
	if len(conds) > 0 {
		q += " WHERE " + strings.Join(conds, " AND ")
	}
	q += fmt.Sprintf(" ORDER BY created_at DESC, id DESC LIMIT $%d", len(args))

	rows, err := s.q.Query(ctx, q, args...)
	if err != nil {
		return Page[Job]{}, err
	}
	defer rows.Close()

	out := make([]Job, 0, limit)
	for rows.Next() {
		var (
			j           Job
			onChainID   *string
			agentIDText *string
			budgetText  *string
			cloneAddr   []byte
			agentPK     *int64
			principal   []byte
			token       []byte
			deadline    *time.Time
			resultURI   *string
			txHash      []byte
			phase       string
		)
		if err := rows.Scan(
			&j.ID,
			&onChainID,
			&cloneAddr,
			&agentPK,
			&agentIDText,
			&principal,
			&j.JobType,
			&token,
			&budgetText,
			&deadline,
			&phase,
			&resultURI,
			&j.BlockNumber,
			&txHash,
			&j.LogIndex,
			&j.CreatedAt,
			&j.UpdatedAt,
		); err != nil {
			return Page[Job]{}, err
		}
		j.OnChainID = derefString(onChainID)
		j.AgentID = derefString(agentIDText)
		j.Budget = derefString(budgetText)
		j.CloneAddress = bytesToHex(cloneAddr)
		j.Principal = bytesToHex(principal)
		j.Token = bytesToHex(token)
		j.TxHash = bytesToHex(txHash)
		j.AgentPK = agentPK
		j.Deadline = deadline
		j.ResultURI = derefString(resultURI)
		j.Phase = JobPhase(phase)
		out = append(out, j)
	}
	if err := rows.Err(); err != nil {
		return Page[Job]{}, err
	}

	page := Page[Job]{Data: out}
	if len(out) > limit {
		last := out[limit-1]
		page.Data = out[:limit]
		page.NextCursor = EncodeCursor(Cursor{ID: last.ID, CreatedAt: last.CreatedAt})
	}
	return page, nil
}

// ---------------------------------------------------------------------------
// Stats
// ---------------------------------------------------------------------------

// Stats returns aggregate counts and the highest indexed block. We compute
// last_block_indexed as GREATEST(MAX(trades), MAX(jobs)) — the indexer track
// will gain a checkpoint table later, at which point this can be a single
// SELECT against that.
//
// PERF: the three COUNT(*) subqueries do a parallel seq scan per call. This
// is fine for v1 — agents/trades/jobs are all small (<<1M rows during the
// alpha) and the endpoint is uncached but called rarely. Once any of these
// tables crosses ~1M rows, swap the COUNT(*) for the planner's row estimate:
//
//   SELECT reltuples::bigint
//     FROM pg_class WHERE relname IN ('agents','jobs','trades')
//
// which is O(1) at the cost of being slightly stale (refreshed by autovacuum).
func (s *PG) Stats(ctx context.Context) (Stats, error) {
	const q = `
		SELECT
			(SELECT count(*) FROM agents)                        AS total_agents,
			(SELECT count(*) FROM jobs)                          AS total_jobs,
			(SELECT count(*) FROM trades)                        AS total_trades,
			COALESCE(GREATEST(
				(SELECT max(block_number) FROM trades),
				(SELECT max(block_number) FROM jobs)
			), 0) AS last_block_indexed`
	var st Stats
	row := s.q.QueryRow(ctx, q)
	if err := row.Scan(&st.TotalAgents, &st.TotalJobs, &st.TotalTrades, &st.LastBlockIndexed); err != nil {
		return Stats{}, err
	}
	return st, nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// bytesToHex returns the lowercased 0x-prefixed hex encoding of b. Empty
// input maps to "" so callers can use omitempty in JSON tags.
func bytesToHex(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return "0x" + hex.EncodeToString(b)
}

// decodeHexAddress parses a 0x-prefixed 20-byte hex string. Inputs are lowered
// before validation so callers can pass mixed-case addresses; we never
// preserve case in the outbound JSON.
func decodeHexAddress(s string) ([]byte, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if !strings.HasPrefix(s, "0x") {
		return nil, errors.New("must start with 0x")
	}
	b, err := hex.DecodeString(s[2:])
	if err != nil {
		return nil, fmt.Errorf("invalid hex: %w", err)
	}
	if len(b) != 20 {
		return nil, fmt.Errorf("must be 20 bytes (got %d)", len(b))
	}
	return b, nil
}

// splitSlug parses "<kind>:<agent_id>" into its components. Returns ok=false
// when the string is not a slug (e.g. a numeric primary key).
func splitSlug(s string) (AgentKind, string, bool) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	k := AgentKind(strings.ToLower(parts[0]))
	if !k.Valid() {
		return "", "", false
	}
	return k, strings.TrimSpace(parts[1]), true
}

// derefString returns *p or "" if p is nil.
func derefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// ensureCommandTagFresh is a tiny compile-time guard: it forces the pgconn
// import so the build catches a future renaming. Cheap, no runtime cost.
var _ pgconn.CommandTag
