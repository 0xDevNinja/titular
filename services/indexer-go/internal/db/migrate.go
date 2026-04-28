// migrate.go wires the embedded migration set into golang-migrate/migrate/v4.
//
// The package owns the SQL files and their ordering invariants (see
// migrations.go). This file owns the driver glue: it returns a configured
// *Migrator that callers use to apply, revert, or inspect the schema of a
// Postgres database identified by a DSN.
//
// API surface (kept intentionally small so the CLI and tests share it):
//
//	New(dsn)              → *Migrator with iofs source + pgx5 database driver
//	(*Migrator).Up()      → apply all pending up migrations
//	(*Migrator).Down()    → revert all applied down migrations
//	(*Migrator).Steps(n)  → apply n (n>0) or revert -n (n<0) migrations
//	(*Migrator).Force(v)  → set version v and clear dirty (recovery only)
//	(*Migrator).Version() → current version, dirty flag
//	(*Migrator).Close()   → release source + database handles
//
// Migrations are loaded from the embedded `migrations/` directory via the
// iofs source driver. The same FS that LoadMigrations validates is what
// golang-migrate consumes, so the on-disk layout is a single source of truth.
package db

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // registers "pgx5" driver
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// MigrationsDir is the embed path that golang-migrate scans. It MUST match
// the //go:embed directive in migrations.go.
const MigrationsDir = "migrations"

// ErrNoChange mirrors migrate.ErrNoChange so callers do not need to import
// the migrate package to handle the "already at target" case.
var ErrNoChange = migrate.ErrNoChange

// Migrator is a thin wrapper over *migrate.Migrate that keeps the source
// reachable for Close.
type Migrator struct {
	m *migrate.Migrate
}

// New returns a Migrator wired to the embedded migration set and the
// Postgres database identified by dsn.
//
// dsn accepts either the libpq URI form (`postgres://user:pw@host:port/db`)
// or the migrate-specific `pgx5://...` form. The function normalises to
// `pgx5://...` because that is what golang-migrate's database registry
// looks up.
//
// The caller MUST call Close when done; otherwise the underlying *sql.DB
// connection pool leaks.
func New(dsn string) (*Migrator, error) {
	if strings.TrimSpace(dsn) == "" {
		return nil, errors.New("db.New: dsn is empty")
	}

	src, err := iofs.New(migrationsFS, MigrationsDir)
	if err != nil {
		return nil, fmt.Errorf("db.New: iofs source: %w", err)
	}

	mig, err := migrate.NewWithSourceInstance("iofs", src, normaliseDSN(dsn))
	if err != nil {
		// Best-effort: close the source so we don't leak the FS handle.
		_ = src.Close()
		return nil, fmt.Errorf("db.New: migrate: %w", err)
	}

	return &Migrator{m: mig}, nil
}

// Up applies every pending migration. Returns nil when there is nothing to
// do — callers do NOT need to special-case ErrNoChange.
func (m *Migrator) Up() error {
	if err := m.m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("db.Up: %w", err)
	}
	return nil
}

// Down reverts every applied migration. Returns nil when the schema is
// already at version 0.
func (m *Migrator) Down() error {
	if err := m.m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("db.Down: %w", err)
	}
	return nil
}

// Steps applies n migrations (forward when n > 0, backward when n < 0).
// n == 0 is rejected so the operator intent is unambiguous.
func (m *Migrator) Steps(n int) error {
	if n == 0 {
		return errors.New("db.Steps: n must be non-zero")
	}
	if err := m.m.Steps(n); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("db.Steps(%d): %w", n, err)
	}
	return nil
}

// Force sets the schema_migrations row to (version, dirty=false) without
// running any SQL. This is a recovery escape hatch — use only when a prior
// failed apply left the table in `dirty=true`.
//
// Negative values are rejected because golang-migrate uses -1 as a sentinel
// for "no migrations applied" and exposing that here would let an operator
// trip the sentinel by accident.
func (m *Migrator) Force(version int) error {
	if version < 0 {
		return fmt.Errorf("db.Force: version must be >= 0, got %d", version)
	}
	if err := m.m.Force(version); err != nil {
		return fmt.Errorf("db.Force(%d): %w", version, err)
	}
	return nil
}

// Version returns the current schema version and whether the last apply
// finished cleanly. dirty == true means a prior apply crashed mid-file and
// the operator must investigate before further migrations are allowed.
//
// When no migrations have ever been applied the migrate package returns
// ErrNilVersion; we translate that to (0, false, nil) so the zero value of
// version means "empty schema" without forcing the caller to import migrate.
func (m *Migrator) Version() (version uint, dirty bool, err error) {
	v, d, err := m.m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, fmt.Errorf("db.Version: %w", err)
	}
	return v, d, nil
}

// Close releases the source FS and the database connection pool. It is
// safe to call on a nil receiver.
//
// Close intentionally aggregates the two underlying errors (source +
// database) rather than returning only the first, because either failure
// is operator-actionable.
func (m *Migrator) Close() error {
	if m == nil || m.m == nil {
		return nil
	}
	srcErr, dbErr := m.m.Close()
	switch {
	case srcErr != nil && dbErr != nil:
		return fmt.Errorf("db.Close: source: %v; database: %w", srcErr, dbErr)
	case srcErr != nil:
		return fmt.Errorf("db.Close: source: %w", srcErr)
	case dbErr != nil:
		return fmt.Errorf("db.Close: database: %w", dbErr)
	}
	return nil
}

// normaliseDSN rewrites the URL scheme so callers may pass either
// `postgres://...` (the standard libpq form) or `pgx5://...` (the form
// golang-migrate's database registry uses). Anything else is returned
// unchanged so future drivers (e.g. CockroachDB) can be plugged in.
func normaliseDSN(dsn string) string {
	for _, prefix := range []string{"postgres://", "postgresql://", "pgx://"} {
		if strings.HasPrefix(dsn, prefix) {
			return "pgx5://" + dsn[len(prefix):]
		}
	}
	return dsn
}
