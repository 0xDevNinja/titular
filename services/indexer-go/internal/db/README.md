# indexer-go / db

Postgres schema for the titular indexer.

## Layout

```
migrations/
  0001_init.up.sql      ← initial schema (agents, agent_tokens, trades,
                          jobs, job_events, memos, processed_logs)
  0001_init.down.sql    ← reverse
migrations.go           ← embedded loader + ordering invariants
migrate.go              ← golang-migrate/v4 wrapper (Up/Down/Steps/Version)
migrations_test.go      ← schema sanity tests
migrations_integration_test.go
                        ← `-tags integration` apply→down→up cycle (dockertest)
```

The CLI that drives the wrapper lives at
`services/indexer-go/cmd/migrate/main.go`.

## Tables

| table            | purpose                                                       |
|------------------|---------------------------------------------------------------|
| `processed_logs` | idempotency guard for `(tx_hash, log_index)` pairs            |
| `agents`         | ACP / launchpad agent identities                              |
| `agent_tokens`   | ERC-20 token + curve for launchpad agents (graduation status) |
| `trades`         | bonding-curve `Bought` / `Sold` events                        |
| `jobs`           | ACP job lifecycle records                                     |
| `job_events`     | append-only audit log of job-level events (incl. escrow)      |
| `memos`          | off-chain signed memos referenced by `reasonHash`             |

## Conventions

- Addresses stored as `bytea(20)` with explicit `octet_length` `CHECK`s.
- `uint256` amounts stored as `numeric(78,0)` (covers `2**256 - 1`).
- Every mutable table has `created_at` and `updated_at`; the `set_updated_at`
  trigger bumps `updated_at` on every `UPDATE`.
- `processed_logs` is the single source of truth for handler idempotency.
  Handlers MUST `IsLogProcessed` before any side-effect and `MarkLogProcessed`
  in the same transaction as the persistence write.

## Adding a migration

1. Pick the next free version number `N`.
2. Create `NNNN_<slug>.up.sql` and `NNNN_<slug>.down.sql` under `migrations/`.
3. **Do NOT** wrap the file in `BEGIN;` / `COMMIT;`. The migration runner
   (`golang-migrate`'s pgx/v5 driver) wraps each file in its own transaction;
   a nested in-file `BEGIN` warns and an in-file `COMMIT` closes the wrapper
   early. If a migration genuinely needs to run outside a transaction (e.g.
   `CREATE INDEX CONCURRENTLY`), put `-- migrate:no-transaction` on the
   first line.
4. Run `go test ./internal/db/...` — tests assert contiguous versions, paired
   up/down files, matching slugs, no nested transactions, and that every
   advertised table is created by some up migration.
5. If you add a new table or enum, also extend `Tables` / `Enums` in
   `migrations.go` so operator tooling stays in sync.

## Applying migrations

The `migrate` binary in `cmd/migrate/` drives the wrapper end-to-end:

```bash
go build -o /tmp/migrate ./cmd/migrate
DATABASE_URL=postgres://titular:titular@localhost:5432/titular?sslmode=disable \
    /tmp/migrate up        # apply every pending migration
/tmp/migrate version       # → "version=1 dirty=false"
/tmp/migrate steps -1      # revert one

# Destructive commands refuse to run without an explicit override. Set
# MIGRATE_ALLOW_DESTRUCTIVE=1 to opt in (CI must NOT export this).
MIGRATE_ALLOW_DESTRUCTIVE=1 /tmp/migrate down       # revert everything
MIGRATE_ALLOW_DESTRUCTIVE=1 /tmp/migrate force 0    # recovery
```

Both the libpq URI form (`postgres://...`) and the migrate-specific
`pgx5://...` form are accepted; the wrapper normalises before handing the
DSN to the driver registry.

Programmatic use (e.g. from the indexer-go `main.go` boot path) uses the
`db.New(dsn)` / `*Migrator` API directly — see `migrate.go`.

## Integration test

The end-to-end cycle (apply → down → re-apply, cross-checking
`schema_migrations` and `pg_class`) lives behind the `integration` build
tag and uses `ory/dockertest` to boot an ephemeral Postgres 16 container:

```bash
go test -tags integration ./internal/db/...
```

Without Docker the tests skip with a clear message. The default
`go test ./...` invocation (used by CI) does not pull Docker in.
