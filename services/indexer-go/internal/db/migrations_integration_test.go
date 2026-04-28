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
// Status: stubbed. Neither `github.com/ory/dockertest/v3` nor
// `github.com/testcontainers/testcontainers-go/modules/postgres` is in this
// service's `go.mod`, and adding either pulls a non-trivial dependency tree
// into a code path that runs only behind a build tag. The implementation
// is tracked as a follow-up so the dep choice can be made alongside the
// rest of the indexer integration harness.
//
// TODO(#78): wire this up using dockertest (smaller dep tree than
// testcontainers) and replace the t.Skip with the real cycle.

//go:build integration

package db

import "testing"

// TestMigrationsApplyRollbackApply is the contract this stub will satisfy
// once a Postgres harness is wired in: apply all up migrations, apply all
// down migrations, then apply the up set again, asserting no error at any
// stage and that the schema looks identical at the start and end.
func TestMigrationsApplyRollbackApply(t *testing.T) {
	t.Skip("requires dockertest harness; tracked in issue #78")
}
