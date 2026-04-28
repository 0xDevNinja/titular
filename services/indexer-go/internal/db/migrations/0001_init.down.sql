-- 0001_init.down.sql
-- Reverse of 0001_init.up.sql.
-- Drops in reverse dependency order so foreign-key constraints unwind cleanly.
--
-- The migration runner wraps this file in its own transaction; do not add an
-- explicit BEGIN/COMMIT here. See the up file for the full rationale.

DROP TABLE IF EXISTS memos;
DROP TABLE IF EXISTS job_events;
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS trades;
DROP TABLE IF EXISTS agent_tokens;
DROP TABLE IF EXISTS agents;
DROP TABLE IF EXISTS processed_logs;

DROP TYPE IF EXISTS job_event_kind;
DROP TYPE IF EXISTS job_phase;
DROP TYPE IF EXISTS trade_side;
DROP TYPE IF EXISTS agent_kind;

DROP FUNCTION IF EXISTS set_updated_at();
