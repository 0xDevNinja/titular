-- Seed data for the gateway-go k6 load test.
--
-- This file is run once against a freshly-migrated indexer database.
-- It populates enough rows to make every endpoint hit a non-empty
-- result set and to guarantee that /api/v1/agents/:id requests in the
-- by-id scenario hit existing primary keys, not 404s.
--
-- Seeding is intentionally synthetic — we do NOT replay real chain
-- data. The load test exercises the gateway's *read-path* throughput,
-- and that is dominated by SQL plan reuse + JSON encoding, both of
-- which are insensitive to whether the bytea values came from anvil
-- or from generate_series.
--
-- Row counts (matching the K6_AGENT_MAX_ID default):
--
--   * 100 agents (50 launchpad + 50 acp)
--   *   1 agent_token row per launchpad agent (50 total)
--   * 200 trades
--   * 100 jobs
--
-- The schema constraints on octet_length(bytea) require us to pad each
-- synthesised address / hash to its declared width; we use lpad with
-- the row index so each row has a distinct value (unique indexes
-- otherwise reject duplicates).

BEGIN;

-- Agents: 50 launchpad + 50 acp. ACP rows must have a non-null
-- agent_id (CHECK constraint). We assign agent_id = i so the unique
-- (kind, agent_id) index is satisfied.

INSERT INTO agents (kind, agent_id, controller, creator, metadata_uri,
                    capabilities, block_number, tx_hash, log_index)
SELECT
    'launchpad'::agent_kind,
    NULL,
    NULL,
    decode(lpad(to_hex(i), 40, '0'), 'hex'),                 -- 20-byte creator
    'ipfs://launchpad/' || i::text,
    NULL,
    1000000 + i,
    decode(lpad(to_hex(i), 64, '0'), 'hex'),                 -- 32-byte tx_hash
    i % 10
FROM generate_series(1, 50) AS s(i);

INSERT INTO agents (kind, agent_id, controller, creator, metadata_uri,
                    capabilities, block_number, tx_hash, log_index)
SELECT
    'acp'::agent_kind,
    i::numeric(78, 0),
    decode(lpad(to_hex(i + 1000), 40, '0'), 'hex'),          -- 20-byte controller (distinct from launchpad)
    NULL,
    'ipfs://acp/' || i::text,
    (i * 7)::numeric(78, 0),                                 -- arbitrary non-zero capabilities bitmap
    2000000 + i,
    decode(lpad(to_hex(i + 100000), 64, '0'), 'hex'),        -- 32-byte tx_hash (distinct from launchpad set)
    i % 10
FROM generate_series(1, 50) AS s(i);

-- agent_tokens: one per launchpad agent. We need these to give trades
-- a foreign key target.

INSERT INTO agent_tokens (agent_pk, token, curve, block_number, tx_hash, log_index)
SELECT
    a.id,
    decode(lpad(to_hex(2000 + a.id), 40, '0'), 'hex'),        -- 20-byte token
    decode(lpad(to_hex(3000 + a.id), 40, '0'), 'hex'),        -- 20-byte curve
    a.block_number,
    decode(lpad(to_hex(200000 + a.id), 64, '0'), 'hex'),
    0
FROM agents a
WHERE a.kind = 'launchpad';

-- Trades: 200 rows alternating buy/sell so both shape constraints get
-- exercised by the listTrades query plan.

INSERT INTO trades (agent_token_pk, side, trader, curve,
                    quote_in, agent_out, agent_in, quote_out, fee,
                    block_number, tx_hash, log_index)
SELECT
    -- Distribute trades across the 50 agent_tokens.
    ((i - 1) % 50) + 1,
    CASE WHEN i % 2 = 0 THEN 'buy'::trade_side ELSE 'sell'::trade_side END,
    decode(lpad(to_hex(4000 + i), 40, '0'), 'hex'),
    -- Pick the curve from the agent_tokens table by the same
    -- modulo we used for agent_token_pk so the (token, curve) pair
    -- is internally consistent.
    decode(lpad(to_hex(3000 + (((i - 1) % 50) + 1)), 40, '0'), 'hex'),
    CASE WHEN i % 2 = 0 THEN (i * 1000)::numeric(78, 0) ELSE NULL END,
    CASE WHEN i % 2 = 0 THEN (i * 100)::numeric(78, 0)  ELSE NULL END,
    CASE WHEN i % 2 = 0 THEN NULL ELSE (i * 100)::numeric(78, 0)  END,
    CASE WHEN i % 2 = 0 THEN NULL ELSE (i * 1000)::numeric(78, 0) END,
    (i * 5)::numeric(78, 0),
    3000000 + i,
    decode(lpad(to_hex(300000 + i), 64, '0'), 'hex'),
    0
FROM generate_series(1, 200) AS s(i);

-- Jobs: 100 rows distributed across all phases so the listJobs status
-- filter has results for every value the API accepts.

INSERT INTO jobs (on_chain_id, agent_pk, agent_id, principal, job_type,
                  budget, phase, block_number, tx_hash, log_index)
SELECT
    i::numeric(78, 0),
    -- Reference an ACP agent (kind='acp' rows have ids 51..100 once
    -- the launchpad inserts above land first because bigserial is
    -- per-table and starts at 1; the ACP block runs 51..100).
    51 + ((i - 1) % 50),
    i::numeric(78, 0),
    decode(lpad(to_hex(5000 + i), 40, '0'), 'hex'),
    (i % 4)::smallint,
    (i * 10000)::numeric(78, 0),
    -- Cycle through every job_phase so listJobs?status=... has hits
    -- regardless of which value the test selects.
    (ARRAY[
        'created'::job_phase,
        'funded'::job_phase,
        'active'::job_phase,
        'completed'::job_phase,
        'cancelled'::job_phase,
        'disputed'::job_phase,
        'released'::job_phase,
        'resolved'::job_phase
    ])[1 + ((i - 1) % 8)],
    4000000 + i,
    decode(lpad(to_hex(400000 + i), 64, '0'), 'hex'),
    0
FROM generate_series(1, 100) AS s(i);

-- Sanity: refresh planner statistics so the first query of the load
-- test does not pay a cold-cache plan-build penalty (which would
-- otherwise show up as a single-digit p99 spike in the warmup window
-- and waste the threshold slack).
ANALYZE agents;
ANALYZE agent_tokens;
ANALYZE trades;
ANALYZE jobs;

COMMIT;
