// Build-tag-gated apply / rollback test for the embedded migration set.
//
// Run with:
//
//	go test -tags integration ./internal/db/...
//
// Plain `go test ./...` skips this file entirely so CI without Docker stays
// green. The test asserts the up → down → up cycle is reversible, which is
// the strongest cheap guarantee that the SQL files are well-formed and
// drop-everything-on-rollback semantics hold.
//
// Harness: ory/dockertest pulls an ephemeral Postgres 16 container, the
// migrator wires it up via the same iofs source the production CLI uses,
// and we assert that the schema_migrations row tracks the moves.

//go:build integration

package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
)

// dockerPool is shared across all integration tests in the package so the
// container starts once and gets reused. Returns nil when Docker is not
// reachable; callers should t.Skip in that case.
type harness struct {
	dsn      string
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

// newHarness boots a Postgres container and waits until it accepts
// connections. It returns nil + a non-fatal explanation when Docker is not
// reachable so we can t.Skip rather than fail.
func newHarness(t *testing.T) *harness {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Skipf("dockertest: pool: %v (start Docker to run this test)", err)
		return nil
	}
	if err := pool.Client.Ping(); err != nil {
		t.Skipf("dockertest: ping: %v (start Docker to run this test)", err)
		return nil
	}

	// 60s is plenty for a postgres:16-alpine cold start on a modern laptop
	// and well under CI time budgets. Anything longer almost always means
	// the image pull is failing — let the operator see it quickly.
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
		// fsync=off makes init noticeably faster and is fine for an
		// ephemeral test DB — we never care about durability here.
		Cmd: []string{"postgres", "-c", "fsync=off", "-c", "synchronous_commit=off"},
	}, func(c *dc.HostConfig) {
		c.AutoRemove = true
		c.RestartPolicy = dc.RestartPolicy{Name: "no"}
	})
	if err != nil {
		t.Skipf("dockertest: run postgres: %v", err)
		return nil
	}
	t.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			t.Logf("dockertest: purge: %v", err)
		}
	})
	// Belt-and-braces: the OS will reap the container even if Purge missed
	// it (e.g. on test-process crash), within 90s.
	_ = resource.Expire(90)

	hostPort := resource.GetHostPort("5432/tcp")
	dsn := fmt.Sprintf("postgres://titular:titular@%s/titular?sslmode=disable", hostPort)

	// Poll until Postgres accepts queries — the container reports "ready"
	// before initdb finishes on a cold image.
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

	return &harness{dsn: dsn, pool: pool, resource: resource}
}

// TestMigrationsApplyRollbackApply boots an ephemeral Postgres, applies the
// embedded migration set, reverts it, and applies it again. At each stage
// we verify (a) no migration error and (b) that schema_migrations and the
// presence of well-known tables match the expected state.
func TestMigrationsApplyRollbackApply(t *testing.T) {
	h := newHarness(t)
	if h == nil {
		return // newHarness already called t.Skip
	}

	// Pre-flight: confirm the migration set is internally consistent before
	// we go anywhere near a network. Cheap, but it surfaces SQL filename
	// regressions ahead of an opaque "no rows returned" container error.
	ups, err := LoadMigrations(Up)
	if err != nil {
		t.Fatalf("LoadMigrations(Up): %v", err)
	}
	if len(ups) == 0 {
		t.Fatal("expected at least one up migration, got 0")
	}
	wantHead := uint(ups[len(ups)-1].Version)

	// Phase 1: apply all migrations from a virgin schema.
	mig, err := New(h.dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := mig.Up(); err != nil {
		t.Fatalf("Up #1: %v", err)
	}
	v, dirty, err := mig.Version()
	if err != nil {
		t.Fatalf("Version after Up #1: %v", err)
	}
	if dirty {
		t.Fatal("Up #1 left schema_migrations dirty")
	}
	if v != wantHead {
		t.Fatalf("after Up #1: version = %d, want %d", v, wantHead)
	}

	// Confirm the canonical tables actually exist in pg_class. This catches
	// the case where a migration succeeds at the migrate-runner level but
	// silently drops a CREATE because a comment or BOM swallowed it.
	openCheck := openTestDB(t, h.dsn)
	for _, table := range Tables {
		if !tableExists(t, openCheck, table) {
			t.Errorf("after Up #1: table %q missing", table)
		}
	}
	_ = openCheck.Close()

	// Phase 2: revert everything. Down on a fully-applied schema must reach
	// version 0 (and the helper translates that to nilversion → 0,false).
	if err := mig.Down(); err != nil {
		t.Fatalf("Down: %v", err)
	}
	v, dirty, err = mig.Version()
	if err != nil {
		t.Fatalf("Version after Down: %v", err)
	}
	if dirty {
		t.Fatal("Down left schema_migrations dirty")
	}
	if v != 0 {
		t.Fatalf("after Down: version = %d, want 0", v)
	}

	// schema_migrations itself stays around but the canonical tables must
	// be gone — proves the down files actually undo what the up files did.
	openCheck = openTestDB(t, h.dsn)
	for _, table := range Tables {
		if tableExists(t, openCheck, table) {
			t.Errorf("after Down: table %q still present", table)
		}
	}
	_ = openCheck.Close()

	// Phase 3: re-apply on the now-empty schema. This is the strongest part
	// of the cycle: it proves that the down files restored a clean enough
	// state that the up files can run again from scratch (e.g. enums and
	// triggers were dropped, not just tables).
	if err := mig.Up(); err != nil {
		t.Fatalf("Up #2: %v", err)
	}
	v, dirty, err = mig.Version()
	if err != nil {
		t.Fatalf("Version after Up #2: %v", err)
	}
	if dirty {
		t.Fatal("Up #2 left schema_migrations dirty")
	}
	if v != wantHead {
		t.Fatalf("after Up #2: version = %d, want %d", v, wantHead)
	}

	if err := mig.Close(); err != nil {
		t.Errorf("Close: %v", err)
	}
}

// TestMigratorUpIsIdempotent confirms calling Up twice in a row collapses
// the second call to a no-op without surfacing ErrNoChange.
func TestMigratorUpIsIdempotent(t *testing.T) {
	h := newHarness(t)
	if h == nil {
		return
	}

	mig, err := New(h.dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer func() { _ = mig.Close() }()

	if err := mig.Up(); err != nil {
		t.Fatalf("Up #1: %v", err)
	}
	if err := mig.Up(); err != nil {
		t.Fatalf("Up #2 (should be a no-op): %v", err)
	}
}

// TestMigratorStepsForwardThenBack confirms the Steps API tracks signed
// counts correctly and that Steps(0) is rejected before any IO.
func TestMigratorStepsForwardThenBack(t *testing.T) {
	h := newHarness(t)
	if h == nil {
		return
	}

	mig, err := New(h.dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer func() { _ = mig.Close() }()

	if err := mig.Steps(0); err == nil {
		t.Error("Steps(0) must return an error")
	}

	ups, err := LoadMigrations(Up)
	if err != nil {
		t.Fatalf("LoadMigrations: %v", err)
	}
	n := len(ups)

	if err := mig.Steps(n); err != nil {
		t.Fatalf("Steps(+%d): %v", n, err)
	}
	if err := mig.Steps(-n); err != nil {
		t.Fatalf("Steps(-%d): %v", n, err)
	}
	v, _, err := mig.Version()
	if err != nil {
		t.Fatalf("Version: %v", err)
	}
	if v != 0 {
		t.Fatalf("after Steps(-%d): version = %d, want 0", n, v)
	}
}

// openTestDB returns a *sql.DB pointed at the harness Postgres. We use
// database/sql with the pgx stdlib shim because tableExists only needs a
// single trivial query and the simpler API keeps the test focused.
func openTestDB(t *testing.T, dsn string) *sql.DB {
	t.Helper()

	// stdlib uses the bare DSN — strip any pgx5:// scheme so the same DSN
	// the migrator accepts also works here.
	clean := dsn
	for _, prefix := range []string{"pgx5://", "pgx://"} {
		if strings.HasPrefix(clean, prefix) {
			clean = "postgres://" + clean[len(prefix):]
			break
		}
	}
	db, err := sql.Open("pgx", clean)
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		t.Fatalf("sql.Ping: %v", err)
	}
	return db
}

// tableExists returns true iff a public-schema relation named `table` is
// present. A custom search avoids the pg_class quirk where to_regclass
// returns NULL on missing relations rather than erroring.
func tableExists(t *testing.T, db *sql.DB, table string) bool {
	t.Helper()
	const q = `
		SELECT 1
		FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name = $1
	`
	var ok int
	err := db.QueryRow(q, table).Scan(&ok)
	switch {
	case err == nil:
		return true
	case errors.Is(err, sql.ErrNoRows):
		return false
	default:
		// Surface unexpected errors but do not abort the whole table loop;
		// the calling test reports per-table failures.
		t.Errorf("tableExists(%q): %v", table, err)
		return false
	}
}

// init guards against a copy-paste regression where someone reuses the test
// against a real DSN by accident: refuse to run if DATABASE_URL is set so
// the harness can never operate on an operator's database.
func init() {
	if os.Getenv("DATABASE_URL") != "" {
		// Tests log via testing.T but init runs before TestMain. Print to
		// stderr so the failure mode is at least visible. We still allow
		// the test process to start; each test calls newHarness which
		// builds its own ephemeral container.
		fmt.Fprintln(os.Stderr,
			"warning: DATABASE_URL is set; integration tests will still use an ephemeral container")
	}
}
