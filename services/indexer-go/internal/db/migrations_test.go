package db

import (
	"strings"
	"testing"
)

// TestLoadMigrationsUp confirms the embedded migration set is well-formed
// and that the ordering is strictly ascending by version.
func TestLoadMigrationsUp(t *testing.T) {
	t.Parallel()

	ms, err := LoadMigrations(Up)
	if err != nil {
		t.Fatalf("LoadMigrations(Up): %v", err)
	}
	if len(ms) == 0 {
		t.Fatal("expected at least one up migration, got 0")
	}
	for i, m := range ms {
		if m.Direction != Up {
			t.Errorf("migration %d: direction = %d, want Up", i, m.Direction)
		}
		if i > 0 && m.Version <= ms[i-1].Version {
			t.Errorf("migration %d (%s) version %d not after previous version %d",
				i, m.Filename(), m.Version, ms[i-1].Version)
		}
		if strings.TrimSpace(m.SQL) == "" {
			t.Errorf("migration %s has empty SQL", m.Filename())
		}
	}
}

// TestLoadMigrationsDown confirms the down set is the exact reverse of the
// up set by version.
func TestLoadMigrationsDown(t *testing.T) {
	t.Parallel()

	ups, err := LoadMigrations(Up)
	if err != nil {
		t.Fatalf("LoadMigrations(Up): %v", err)
	}
	downs, err := LoadMigrations(Down)
	if err != nil {
		t.Fatalf("LoadMigrations(Down): %v", err)
	}
	if len(ups) != len(downs) {
		t.Fatalf("up count %d != down count %d", len(ups), len(downs))
	}
	for i := range ups {
		want := ups[len(ups)-1-i]
		got := downs[i]
		if got.Version != want.Version {
			t.Errorf("downs[%d] version = %d, want %d", i, got.Version, want.Version)
		}
		if got.Name != want.Name {
			t.Errorf("downs[%d] name = %q, want %q", i, got.Name, want.Name)
		}
		if got.Direction != Down {
			t.Errorf("downs[%d] direction = %d, want Down", i, got.Direction)
		}
	}
}

// TestInitialMigrationCreatesAllTables guarantees that every table named in
// the issue acceptance criteria appears in the initial schema. If a future
// refactor renames or splits a table, this test forces an explicit decision.
func TestInitialMigrationCreatesAllTables(t *testing.T) {
	t.Parallel()

	required := []string{
		"agents",
		"agent_tokens",
		"trades",
		"jobs",
		"job_events",
		"memos",
	}

	ms, err := LoadMigrations(Up)
	if err != nil {
		t.Fatalf("LoadMigrations: %v", err)
	}

	combined := strings.Builder{}
	for _, m := range ms {
		combined.WriteString(m.SQL)
	}
	sql := combined.String()

	for _, table := range required {
		if !strings.Contains(sql, "CREATE TABLE "+table+" ") &&
			!strings.Contains(sql, "CREATE TABLE "+table+"\n") {
			t.Errorf("initial schema is missing CREATE TABLE %s", table)
		}
	}
}

// TestInitialMigrationCreatesEnumTypes verifies the enum types referenced
// by the initial schema exist before they are used in column definitions.
func TestInitialMigrationCreatesEnumTypes(t *testing.T) {
	t.Parallel()

	required := []string{
		"agent_kind",
		"trade_side",
		"job_phase",
		"job_event_kind",
	}

	ms, err := LoadMigrations(Up)
	if err != nil {
		t.Fatalf("LoadMigrations: %v", err)
	}
	sql := strings.Builder{}
	for _, m := range ms {
		sql.WriteString(m.SQL)
	}
	body := sql.String()
	for _, e := range required {
		if !strings.Contains(body, "CREATE TYPE "+e+" AS ENUM") {
			t.Errorf("initial schema missing CREATE TYPE %s AS ENUM", e)
		}
	}
}

// TestSchemaTablesListMatchesEmbedded asserts the exported Tables slice
// covers every CREATE TABLE in the embedded migrations and vice-versa, so
// downstream tooling that iterates Tables stays in sync.
func TestSchemaTablesListMatchesEmbedded(t *testing.T) {
	t.Parallel()

	if err := schemaSanityCheck(); err != nil {
		t.Fatalf("schemaSanityCheck: %v", err)
	}
}

// TestMigrationFilesAreTransactional ensures every file wraps its body in
// BEGIN/COMMIT so a failed apply does not leave half-applied state.
func TestMigrationFilesAreTransactional(t *testing.T) {
	t.Parallel()

	for _, dir := range []Direction{Up, Down} {
		ms, err := LoadMigrations(dir)
		if err != nil {
			t.Fatalf("LoadMigrations(%d): %v", dir, err)
		}
		for _, m := range ms {
			if !strings.Contains(m.SQL, "BEGIN;") {
				t.Errorf("%s missing BEGIN;", m.Filename())
			}
			if !strings.Contains(m.SQL, "COMMIT;") {
				t.Errorf("%s missing COMMIT;", m.Filename())
			}
		}
	}
}

// TestMigrationStringer confirms the human-readable form is stable, since
// we use it in operator-facing logs.
func TestMigrationStringer(t *testing.T) {
	t.Parallel()

	m := Migration{Version: 7, Name: "session_keys", Direction: Up}
	if got, want := m.String(), "0007_session_keys.up"; got != want {
		t.Errorf("String = %q, want %q", got, want)
	}
	if got, want := m.Filename(), "0007_session_keys.up.sql"; got != want {
		t.Errorf("Filename = %q, want %q", got, want)
	}
	d := Migration{Version: 7, Name: "session_keys", Direction: Down}
	if got, want := d.String(), "0007_session_keys.down"; got != want {
		t.Errorf("String = %q, want %q", got, want)
	}
}

// TestProcessedLogsKey checks that the idempotency table is part of the
// initial schema (used by every handler in handlers/store.go).
func TestProcessedLogsKey(t *testing.T) {
	t.Parallel()

	ms, err := LoadMigrations(Up)
	if err != nil {
		t.Fatalf("LoadMigrations: %v", err)
	}
	body := strings.Builder{}
	for _, m := range ms {
		body.WriteString(m.SQL)
	}
	got := body.String()
	if !strings.Contains(got, "CREATE TABLE processed_logs") {
		t.Fatal("processed_logs table missing from initial schema")
	}
	if !strings.Contains(got, "PRIMARY KEY (tx_hash, log_index)") {
		t.Error("processed_logs must have composite primary key (tx_hash, log_index)")
	}
}
