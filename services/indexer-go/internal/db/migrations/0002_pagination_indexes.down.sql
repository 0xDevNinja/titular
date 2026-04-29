-- 0002_pagination_indexes.down.sql
-- Reverses 0002_pagination_indexes.up.sql by dropping the composite
-- pagination indexes in reverse creation order.

DROP INDEX IF EXISTS jobs_created_at_idx;
DROP INDEX IF EXISTS trades_created_at_idx;
DROP INDEX IF EXISTS agents_created_at_idx;
