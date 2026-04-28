// Package db owns the indexer's Postgres schema and exposes the embedded
// migration set used by callers to bring an empty database up to current.
//
// Migrations are plain .sql files under ./migrations, named
//
//	NNNN_<slug>.{up,down}.sql
//
// where NNNN is a zero-padded sequence number. Each file MUST contain its
// own BEGIN/COMMIT — the loader does not wrap statements in a transaction
// because some migrations (e.g. CREATE INDEX CONCURRENTLY in the future)
// cannot run inside one.
package db

import (
	"embed"
	"fmt"
	"io/fs"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Direction selects up or down migrations.
type Direction int

const (
	// Up applies a migration.
	Up Direction = iota
	// Down reverts a migration.
	Down
)

// Migration is a single ordered SQL script.
type Migration struct {
	// Version is the numeric prefix (e.g. 1 for "0001_init.up.sql").
	Version int
	// Name is the slug between the version and the direction
	// (e.g. "init" for "0001_init.up.sql").
	Name string
	// Direction is Up or Down.
	Direction Direction
	// SQL is the file contents.
	SQL string
}

// migrationFilename matches "NNNN_<slug>.{up,down}.sql" where NNNN is at
// least one digit. The slug must not be empty and may contain underscores
// and ASCII letters/digits.
var migrationFilename = regexp.MustCompile(
	`^(\d+)_([A-Za-z0-9][A-Za-z0-9_]*)\.(up|down)\.sql$`,
)

// LoadMigrations returns all migrations in the requested direction sorted by
// version ascending (Up) or descending (Down). It also validates that:
//
//   - every version has both an up and a down file,
//   - versions are contiguous starting from 1,
//   - the slug matches between an up/down pair.
//
// These invariants catch typos like a renamed slug on only one side of a
// pair, which would otherwise drift silently.
func LoadMigrations(direction Direction) ([]Migration, error) {
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("db.LoadMigrations: read embedded dir: %w", err)
	}

	type pair struct {
		name string
		up   *Migration
		down *Migration
	}
	pairs := map[int]*pair{}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fname := e.Name()
		m := migrationFilename.FindStringSubmatch(fname)
		if m == nil {
			return nil, fmt.Errorf(
				"db.LoadMigrations: invalid migration filename %q "+
					"(expected NNNN_<slug>.{up,down}.sql)", fname)
		}

		version, err := strconv.Atoi(m[1])
		if err != nil {
			return nil, fmt.Errorf("db.LoadMigrations: parse version in %q: %w", fname, err)
		}
		slug := m[2]
		dir := m[3]

		body, err := fs.ReadFile(migrationsFS, "migrations/"+fname)
		if err != nil {
			return nil, fmt.Errorf("db.LoadMigrations: read %q: %w", fname, err)
		}

		mig := &Migration{
			Version: version,
			Name:    slug,
			SQL:     string(body),
		}

		p, ok := pairs[version]
		if !ok {
			p = &pair{name: slug}
			pairs[version] = p
		}

		if p.name != slug {
			return nil, fmt.Errorf(
				"db.LoadMigrations: version %04d has mismatched slugs "+
					"(%q vs %q) — rename both files together",
				version, p.name, slug)
		}

		switch dir {
		case "up":
			if p.up != nil {
				return nil, fmt.Errorf(
					"db.LoadMigrations: duplicate up migration for version %04d", version)
			}
			mig.Direction = Up
			p.up = mig
		case "down":
			if p.down != nil {
				return nil, fmt.Errorf(
					"db.LoadMigrations: duplicate down migration for version %04d", version)
			}
			mig.Direction = Down
			p.down = mig
		}
	}

	// Sort versions and validate contiguity.
	versions := make([]int, 0, len(pairs))
	for v := range pairs {
		versions = append(versions, v)
	}
	sort.Ints(versions)

	for i, v := range versions {
		if i == 0 && v != 1 {
			return nil, fmt.Errorf(
				"db.LoadMigrations: first migration must be version 1, got %d", v)
		}
		if i > 0 && v != versions[i-1]+1 {
			return nil, fmt.Errorf(
				"db.LoadMigrations: non-contiguous versions %d → %d", versions[i-1], v)
		}
		p := pairs[v]
		if p.up == nil {
			return nil, fmt.Errorf(
				"db.LoadMigrations: version %04d (%s) missing up migration", v, p.name)
		}
		if p.down == nil {
			return nil, fmt.Errorf(
				"db.LoadMigrations: version %04d (%s) missing down migration", v, p.name)
		}
	}

	out := make([]Migration, 0, len(versions))
	switch direction {
	case Up:
		for _, v := range versions {
			out = append(out, *pairs[v].up)
		}
	case Down:
		for i := len(versions) - 1; i >= 0; i-- {
			out = append(out, *pairs[versions[i]].down)
		}
	default:
		return nil, fmt.Errorf("db.LoadMigrations: unknown direction %d", direction)
	}

	return out, nil
}

// MustLoadMigrations is the panicking variant for use in init blocks where a
// failure is a programmer error (e.g. malformed embedded files).
func MustLoadMigrations(direction Direction) []Migration {
	ms, err := LoadMigrations(direction)
	if err != nil {
		panic(err)
	}
	return ms
}

// String renders a Migration as "NNNN_<slug>.<direction>".
func (m Migration) String() string {
	dir := "up"
	if m.Direction == Down {
		dir = "down"
	}
	return fmt.Sprintf("%04d_%s.%s", m.Version, m.Name, dir)
}

// Filename returns the on-disk filename of the migration.
func (m Migration) Filename() string {
	return m.String() + ".sql"
}

// Tables enumerates the canonical table names created by the initial schema.
// Exported for use in smoke tests and operator tooling.
var Tables = []string{
	"processed_logs",
	"agents",
	"agent_tokens",
	"trades",
	"jobs",
	"job_events",
	"memos",
}

// Enums enumerates the canonical enum types created by the initial schema.
var Enums = []string{
	"agent_kind",
	"trade_side",
	"job_phase",
	"job_event_kind",
}

// schemaSanityCheck is a defensive guard that returns an error if the
// embedded migrations do not mention every advertised table or enum. This
// makes it impossible to silently rename a table without updating the
// exported lists above.
func schemaSanityCheck() error {
	ups, err := LoadMigrations(Up)
	if err != nil {
		return err
	}
	combined := strings.Builder{}
	for _, m := range ups {
		combined.WriteString(m.SQL)
		combined.WriteString("\n")
	}
	body := combined.String()
	for _, t := range Tables {
		if !strings.Contains(body, "CREATE TABLE "+t) {
			return fmt.Errorf("db: schema sanity: missing CREATE TABLE %s", t)
		}
	}
	for _, e := range Enums {
		if !strings.Contains(body, "CREATE TYPE "+e) {
			return fmt.Errorf("db: schema sanity: missing CREATE TYPE %s", e)
		}
	}
	return nil
}
