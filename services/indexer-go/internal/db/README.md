# indexer-go / db

Postgres schema for the titular indexer.

## Layout

```
migrations/
  0001_init.up.sql      ← initial schema (agents, agent_tokens, trades,
                          jobs, job_events, memos, processed_logs)
  0001_init.down.sql    ← reverse
migrations.go           ← embedded loader + ordering invariants
migrations_test.go      ← schema sanity tests
```

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
3. Each file MUST wrap its body in `BEGIN;` / `COMMIT;`.
4. Run `go test ./internal/db/...` — tests assert contiguous versions, paired
   up/down files, matching slugs, transactional bracketing, and that every
   advertised table is created by some up migration.
5. If you add a new table or enum, also extend `Tables` / `Enums` in
   `migrations.go` so operator tooling stays in sync.

## Applying migrations

This package only loads + orders migrations; the apply driver lives elsewhere
(future work — see follow-up issues for runtime apply, schema-version table,
and idempotency-on-apply). To apply by hand against a local dev DB:

```bash
psql "$DATABASE_URL" -f services/indexer-go/internal/db/migrations/0001_init.up.sql
```
