-- 0002_pagination_indexes.up.sql
-- Composite indexes that back the gateway's keyset pagination over the three
-- list endpoints introduced in M4 (issue #88):
--
--   GET /api/v1/agents      ORDER BY (created_at DESC, id DESC)
--   GET /api/v1/trades      ORDER BY (created_at DESC, id DESC)
--   GET /api/v1/jobs        ORDER BY (created_at DESC, id DESC)
--
-- Without these, each page beyond the first degenerates to a full sort over
-- the table — fine while seeded, painful once any of these tables grows past
-- a few hundred thousand rows. The (created_at, id) tie-breaker matches the
-- handler's cursor encode/decode shape exactly so the planner can serve the
-- cursor predicate `(created_at, id) < ($1, $2)` as an index range scan.
--
-- Transactions:
--   golang-migrate wraps each file in its own transaction. Do NOT add an
--   explicit BEGIN/COMMIT here — see 0001_init.up.sql for the rationale.
--   These are plain CREATE INDEX (not CONCURRENTLY); on small dev/test
--   databases this is fine. If a future deploy needs CONCURRENTLY against a
--   live DB, split into a no-transaction migration with the
--   `-- migrate:no-transaction` directive on line 1.

CREATE INDEX agents_created_at_idx ON agents (created_at DESC, id DESC);
CREATE INDEX trades_created_at_idx ON trades (created_at DESC, id DESC);
CREATE INDEX jobs_created_at_idx   ON jobs   (created_at DESC, id DESC);
