// Build-tag-gated integration test for the Postgres store.
//
// Run with:
//
//	go test -tags integration ./internal/store/...
//
// Plain `go test ./...` skips this file. The test boots an ephemeral
// Postgres 16 container via ory/dockertest, applies the indexer schema by
// pointing golang-migrate at the indexer's embedded migration set, seeds a
// few rows, and asserts the store's read paths return the expected wire
// shape.

//go:build integration

package store_test

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/store"
)

// startPG boots an ephemeral Postgres container and returns its DSN. Skips
// (rather than fails) when Docker is not reachable so CI without Docker stays
// green.
func startPG(t *testing.T) string {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Skipf("dockertest: pool: %v", err)
	}
	if err := pool.Client.Ping(); err != nil {
		t.Skipf("dockertest: docker daemon unreachable: %v", err)
	}
	pool.MaxWait = 60 * time.Second

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16-alpine",
		Env: []string{
			"POSTGRES_PASSWORD=titular",
			"POSTGRES_USER=titular",
			"POSTGRES_DB=titular",
			"listen_addresses=*",
		},
		Cmd: []string{"postgres", "-c", "fsync=off", "-c", "synchronous_commit=off"},
	}, func(c *dc.HostConfig) {
		c.AutoRemove = true
		c.RestartPolicy = dc.RestartPolicy{Name: "no"}
	})
	if err != nil {
		t.Skipf("dockertest: run postgres: %v", err)
	}
	t.Cleanup(func() { _ = pool.Purge(resource) })
	_ = resource.Expire(90)

	dsn := fmt.Sprintf("postgres://titular:titular@%s/titular?sslmode=disable",
		resource.GetHostPort("5432/tcp"))

	if err := pool.Retry(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn, err := pgx.Connect(ctx, dsn)
		if err != nil {
			return err
		}
		defer conn.Close(ctx)
		return conn.Ping(ctx)
	}); err != nil {
		t.Fatalf("postgres never became ready: %v", err)
	}
	return dsn
}

// applyIndexerSchema runs the indexer's migrations against dsn by pointing
// golang-migrate at the SQL files on disk. We cannot import the indexer's
// internal/db package because Go's module rules forbid cross-module access
// to internal/, so the file:// source is the cheapest stable handshake. The
// path is resolved relative to this file so it survives `go test` invocation
// from anywhere in the repo.
func applyIndexerSchema(t *testing.T, dsn string) {
	t.Helper()
	_, file, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(file),
		"..", "..", "..", "indexer-go", "internal", "db", "migrations")

	// golang-migrate uses pgx5:// for the pgx/v5 driver.
	pgxDSN := dsn
	if len(pgxDSN) > len("postgres://") && pgxDSN[:len("postgres://")] == "postgres://" {
		pgxDSN = "pgx5://" + pgxDSN[len("postgres://"):]
	}

	m, err := migrate.New("file://"+migrationsDir, pgxDSN)
	if err != nil {
		t.Fatalf("migrate.New: %v (migrationsDir=%s)", err, migrationsDir)
	}
	defer func() { _, _ = m.Close() }()
	if err := m.Up(); err != nil {
		t.Fatalf("apply schema: %v", err)
	}
}

// seedRows inserts a deterministic, minimal row set covering the surface the
// gateway reads.
func seedRows(t *testing.T, dsn string) {
	t.Helper()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("pool: %v", err)
	}
	defer pool.Close()

	addr := func(b byte) []byte {
		out := make([]byte, 20)
		for i := range out {
			out[i] = b
		}
		return out
	}
	tx32 := func(b byte) []byte {
		out := make([]byte, 32)
		for i := range out {
			out[i] = b
		}
		return out
	}

	// Two agents (one launchpad, one ACP).
	if _, err := pool.Exec(ctx, `
		INSERT INTO agents (kind, agent_id, controller, creator, metadata_uri, capabilities, block_number, tx_hash, log_index)
		VALUES
			('launchpad', NULL, NULL, $1, 'ipfs://A', NULL, 100, $2, 0),
			('acp', 42, $3, NULL, 'ipfs://B', 7, 101, $4, 0)`,
		addr(0xaa), tx32(0x01), addr(0xbb), tx32(0x02),
	); err != nil {
		t.Fatalf("seed agents: %v", err)
	}

	// One agent_token + two trades.
	if _, err := pool.Exec(ctx, `
		INSERT INTO agent_tokens (agent_pk, token, curve, block_number, tx_hash, log_index)
		VALUES (1, $1, $2, 100, $3, 1)`,
		addr(0xc1), addr(0xc2), tx32(0x03),
	); err != nil {
		t.Fatalf("seed agent_tokens: %v", err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO trades (agent_token_pk, side, trader, curve, quote_in, agent_out, fee, block_number, tx_hash, log_index)
		VALUES
			(1, 'buy', $1, $2, 1000, 500, 1, 110, $3, 0),
			(1, 'buy', $1, $2, 2000, 800, 1, 111, $4, 0)`,
		addr(0xaa), addr(0xc2), tx32(0x10), tx32(0x11),
	); err != nil {
		t.Fatalf("seed trades: %v", err)
	}

	// Two jobs (different phases).
	if _, err := pool.Exec(ctx, `
		INSERT INTO jobs (on_chain_id, agent_pk, agent_id, principal, job_type, budget, phase, block_number, tx_hash, log_index)
		VALUES
			(1, 2, 42, $1, 0, 5000, 'created', 120, $2, 0),
			(2, 2, 42, $1, 0, 6000, 'funded',  121, $3, 0)`,
		addr(0xdd), tx32(0x20), tx32(0x21),
	); err != nil {
		t.Fatalf("seed jobs: %v", err)
	}
}

func TestPG_EndToEnd(t *testing.T) {
	dsn := startPG(t)
	applyIndexerSchema(t, dsn)
	seedRows(t, dsn)

	ctx := context.Background()
	pg, closeFn, err := store.Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer closeFn()

	// ListAgents — newest first; the ACP row was inserted second so it
	// comes back first under created_at DESC.
	agentsPage, err := pg.ListAgents(ctx, store.ListAgentsParams{Limit: 10})
	if err != nil {
		t.Fatalf("list agents: %v", err)
	}
	if len(agentsPage.Data) != 2 {
		t.Fatalf("agents len: got %d, want 2", len(agentsPage.Data))
	}

	// Filter by kind.
	acpOnly, err := pg.ListAgents(ctx, store.ListAgentsParams{Kind: store.AgentKindACP, Limit: 10})
	if err != nil {
		t.Fatalf("list agents (acp): %v", err)
	}
	if len(acpOnly.Data) != 1 || acpOnly.Data[0].Kind != store.AgentKindACP {
		t.Fatalf("kind filter wrong: got %+v", acpOnly.Data)
	}
	if acpOnly.Data[0].AgentID != "42" {
		t.Fatalf("agent_id: got %q, want %q", acpOnly.Data[0].AgentID, "42")
	}

	// GetAgent by primary key.
	if _, err := pg.GetAgent(ctx, "1"); err != nil {
		t.Fatalf("get agent by pk: %v", err)
	}

	// GetAgent by slug.
	a, err := pg.GetAgent(ctx, "acp:42")
	if err != nil {
		t.Fatalf("get agent by slug: %v", err)
	}
	if a.AgentID != "42" {
		t.Fatalf("slug agent_id: got %q", a.AgentID)
	}

	// GetAgent unknown.
	if _, err := pg.GetAgent(ctx, "acp:9999"); err != store.ErrNotFound {
		t.Fatalf("get unknown agent: got %v, want ErrNotFound", err)
	}

	// ListTrades — both rows for the seeded token.
	tradesPage, err := pg.ListTrades(ctx, store.ListTradesParams{Limit: 10})
	if err != nil {
		t.Fatalf("list trades: %v", err)
	}
	if len(tradesPage.Data) != 2 {
		t.Fatalf("trades: got %d, want 2", len(tradesPage.Data))
	}

	// ListTrades filtered by agent_token. The seed used 0xc1 repeated 20
	// times for the token address.
	filterTrades, err := pg.ListTrades(ctx, store.ListTradesParams{
		AgentToken: "0x" + repeat("c1", 1) + repeat("c1", 19),
		Limit:      10,
	})
	if err != nil {
		t.Fatalf("list trades (token): %v", err)
	}
	// All seeded trades belong to the seeded token, so the count matches.
	if len(filterTrades.Data) != 2 {
		t.Fatalf("trades by token: got %d, want 2", len(filterTrades.Data))
	}

	// ListTrades with bogus token rejected at validation.
	if _, err := pg.ListTrades(ctx, store.ListTradesParams{AgentToken: "not-hex", Limit: 5}); err == nil {
		t.Fatal("bogus agent_token should error")
	}

	// ListJobs — both seeded.
	jobsPage, err := pg.ListJobs(ctx, store.ListJobsParams{Limit: 10})
	if err != nil {
		t.Fatalf("list jobs: %v", err)
	}
	if len(jobsPage.Data) != 2 {
		t.Fatalf("jobs: got %d, want 2", len(jobsPage.Data))
	}

	// Filter by status.
	funded, err := pg.ListJobs(ctx, store.ListJobsParams{Status: store.JobPhaseFunded, Limit: 10})
	if err != nil {
		t.Fatalf("list jobs (funded): %v", err)
	}
	if len(funded.Data) != 1 || funded.Data[0].Phase != store.JobPhaseFunded {
		t.Fatalf("funded filter wrong: got %+v", funded.Data)
	}

	// Stats — counts and last_block_indexed = max(trades, jobs) = 121.
	st, err := pg.Stats(ctx)
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
	if st.TotalAgents != 2 || st.TotalJobs != 2 || st.TotalTrades != 2 {
		t.Fatalf("counts: got %+v", st)
	}
	if st.LastBlockIndexed != 121 {
		t.Fatalf("last_block_indexed: got %d, want 121", st.LastBlockIndexed)
	}
}

// repeat returns s concatenated n times. Tiny so we keep it inline rather
// than pulling strings.Repeat (which would be the only stdlib reason to
// import strings in this file).
func repeat(s string, n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}
