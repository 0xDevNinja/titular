-- 0001_init.up.sql
-- Initial Postgres schema for the titular indexer.
--
-- Tables (in dependency order):
--   1. agents          — registered ACP / launchpad agent identities
--   2. agent_tokens    — ERC-20 tokens linked to launchpad agents
--   3. trades          — bonding-curve buys/sells
--   4. jobs            — ACP job lifecycle records
--   5. job_events      — append-only audit log for job state transitions
--   6. memos           — off-chain memos referenced by content hash
--
-- Cross-cutting:
--   - processed_logs   — idempotency guard for (tx_hash, log_index)
--
-- Conventions:
--   - All on-chain addresses stored as bytea(20) for compact indexing.
--   - All uint256 amounts stored as numeric(78,0) — wide enough for 2**256-1.
--   - All timestamps are timestamptz; `created_at` defaults to now().
--   - Updates bump `updated_at` via trigger (see set_updated_at).

BEGIN;

-- ---------------------------------------------------------------------------
-- Helper: bump updated_at on every row update.
-- ---------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END;
$$;

-- ---------------------------------------------------------------------------
-- processed_logs — idempotency guard.
-- A handler writes one row here after persisting a log so retries skip it.
-- ---------------------------------------------------------------------------
CREATE TABLE processed_logs (
    tx_hash      bytea       NOT NULL,
    log_index    integer     NOT NULL,
    block_number bigint      NOT NULL,
    processed_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (tx_hash, log_index),
    CONSTRAINT processed_logs_tx_hash_len CHECK (octet_length(tx_hash) = 32),
    CONSTRAINT processed_logs_log_index_nonneg CHECK (log_index >= 0)
);

CREATE INDEX processed_logs_block_number_idx ON processed_logs (block_number);

-- ---------------------------------------------------------------------------
-- agents — agent identity (ACP AgentRegistry + launchpad agents).
--
-- `agent_id` is the on-chain uint256 id (ACP) or NULL for launchpad-only
-- agents that predate ACP registration. `kind` distinguishes the source so
-- queries can filter without joining.
-- ---------------------------------------------------------------------------
CREATE TYPE agent_kind AS ENUM ('launchpad', 'acp');

CREATE TABLE agents (
    id            bigserial   PRIMARY KEY,
    kind          agent_kind  NOT NULL,
    agent_id      numeric(78, 0),                  -- on-chain uint256 (NULL for legacy launchpad agents)
    controller    bytea,                           -- current controller (ACP only)
    creator       bytea,                           -- launchpad creator (launchpad only)
    metadata_uri  text,
    capabilities  numeric(78, 0),                  -- ACP capabilities bitmap as uint256
    block_number  bigint      NOT NULL,
    tx_hash       bytea       NOT NULL,
    log_index     integer     NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT agents_controller_len    CHECK (controller IS NULL OR octet_length(controller) = 20),
    CONSTRAINT agents_creator_len       CHECK (creator    IS NULL OR octet_length(creator)    = 20),
    CONSTRAINT agents_tx_hash_len       CHECK (octet_length(tx_hash) = 32),
    CONSTRAINT agents_log_index_nonneg  CHECK (log_index >= 0),
    -- ACP agents must have a non-null agent_id; launchpad agents may not.
    CONSTRAINT agents_acp_has_id        CHECK (kind <> 'acp' OR agent_id IS NOT NULL)
);

-- Lookup ACP agents by their on-chain id (one active row per id per kind).
CREATE UNIQUE INDEX agents_kind_agent_id_uidx
    ON agents (kind, agent_id)
    WHERE agent_id IS NOT NULL;

CREATE INDEX agents_controller_idx ON agents (controller) WHERE controller IS NOT NULL;
CREATE INDEX agents_creator_idx    ON agents (creator)    WHERE creator    IS NOT NULL;

CREATE TRIGGER agents_set_updated_at
    BEFORE UPDATE ON agents
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ---------------------------------------------------------------------------
-- agent_tokens — ERC-20 tokens minted for launchpad agents.
--
-- One row per (token, curve) pair. `pair` is set when the bonding curve
-- graduates to a Uniswap V2 pair. `lp_lock` is the LP-lock contract address.
-- ---------------------------------------------------------------------------
CREATE TABLE agent_tokens (
    id            bigserial   PRIMARY KEY,
    agent_pk      bigint      NOT NULL REFERENCES agents (id) ON DELETE CASCADE,
    token         bytea       NOT NULL,
    curve         bytea       NOT NULL,
    lp_lock       bytea,
    pair          bytea,
    graduated     boolean     NOT NULL DEFAULT FALSE,
    graduated_at  timestamptz,
    block_number  bigint      NOT NULL,
    tx_hash       bytea       NOT NULL,
    log_index     integer     NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT agent_tokens_token_len    CHECK (octet_length(token) = 20),
    CONSTRAINT agent_tokens_curve_len    CHECK (octet_length(curve) = 20),
    CONSTRAINT agent_tokens_lp_lock_len  CHECK (lp_lock IS NULL OR octet_length(lp_lock) = 20),
    CONSTRAINT agent_tokens_pair_len     CHECK (pair    IS NULL OR octet_length(pair)    = 20),
    CONSTRAINT agent_tokens_tx_hash_len  CHECK (octet_length(tx_hash) = 32),
    CONSTRAINT agent_tokens_log_index_nonneg CHECK (log_index >= 0),
    -- A graduated token must have a pair address and a graduation timestamp.
    CONSTRAINT agent_tokens_graduated_consistent
        CHECK (graduated = FALSE OR (pair IS NOT NULL AND graduated_at IS NOT NULL))
);

CREATE UNIQUE INDEX agent_tokens_token_uidx ON agent_tokens (token);
CREATE UNIQUE INDEX agent_tokens_curve_uidx ON agent_tokens (curve);
CREATE INDEX agent_tokens_agent_pk_idx ON agent_tokens (agent_pk);
CREATE INDEX agent_tokens_pair_idx ON agent_tokens (pair) WHERE pair IS NOT NULL;

CREATE TRIGGER agent_tokens_set_updated_at
    BEFORE UPDATE ON agent_tokens
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ---------------------------------------------------------------------------
-- trades — bonding-curve Bought / Sold events.
-- ---------------------------------------------------------------------------
CREATE TYPE trade_side AS ENUM ('buy', 'sell');

CREATE TABLE trades (
    id            bigserial    PRIMARY KEY,
    agent_token_pk bigint      REFERENCES agent_tokens (id) ON DELETE SET NULL,
    side          trade_side   NOT NULL,
    trader        bytea        NOT NULL,
    curve         bytea        NOT NULL,
    quote_in      numeric(78, 0),  -- non-null on buy
    agent_out     numeric(78, 0),  -- non-null on buy
    agent_in      numeric(78, 0),  -- non-null on sell
    quote_out     numeric(78, 0),  -- non-null on sell
    fee           numeric(78, 0) NOT NULL,
    block_number  bigint       NOT NULL,
    tx_hash       bytea        NOT NULL,
    log_index     integer      NOT NULL,
    created_at    timestamptz  NOT NULL DEFAULT now(),
    CONSTRAINT trades_trader_len   CHECK (octet_length(trader) = 20),
    CONSTRAINT trades_curve_len    CHECK (octet_length(curve)  = 20),
    CONSTRAINT trades_tx_hash_len  CHECK (octet_length(tx_hash) = 32),
    CONSTRAINT trades_log_index_nonneg CHECK (log_index >= 0),
    -- A buy has quote_in + agent_out; a sell has agent_in + quote_out. Cross-
    -- direction columns must be NULL so totalling queries cannot be ambiguous.
    CONSTRAINT trades_buy_shape CHECK (
        side <> 'buy' OR (
            quote_in  IS NOT NULL AND agent_out IS NOT NULL
            AND agent_in IS NULL AND quote_out IS NULL
        )
    ),
    CONSTRAINT trades_sell_shape CHECK (
        side <> 'sell' OR (
            agent_in IS NOT NULL AND quote_out IS NOT NULL
            AND quote_in IS NULL AND agent_out IS NULL
        )
    )
);

CREATE UNIQUE INDEX trades_log_uidx ON trades (tx_hash, log_index);
CREATE INDEX trades_curve_block_idx ON trades (curve, block_number DESC);
CREATE INDEX trades_trader_block_idx ON trades (trader, block_number DESC);
CREATE INDEX trades_agent_token_pk_idx ON trades (agent_token_pk) WHERE agent_token_pk IS NOT NULL;

-- ---------------------------------------------------------------------------
-- jobs — ACP job lifecycle.
--
-- `phase` tracks the canonical lifecycle and is enforced via a CHECK on
-- forward transitions in application code (no DB-level transition trigger;
-- replays from chain re-establish state, so DB must accept any valid phase).
-- ---------------------------------------------------------------------------
CREATE TYPE job_phase AS ENUM (
    'created',
    'funded',
    'active',
    'completed',
    'cancelled',
    'disputed',
    'released',
    'resolved'
);

CREATE TABLE jobs (
    id              bigserial    PRIMARY KEY,
    on_chain_id     numeric(78, 0) NOT NULL,    -- uint256 from JobFactory.JobCreated
    clone_address   bytea,                      -- minimal-proxy clone (JobFactory.JobCreated.clone)
    agent_pk        bigint       REFERENCES agents (id) ON DELETE SET NULL,
    agent_id        numeric(78, 0),             -- on-chain agent id (denormalised for fast filter)
    principal       bytea        NOT NULL,
    job_type        smallint     NOT NULL,
    token           bytea,                      -- payment token (zero-address-equivalent NULL = native)
    budget          numeric(78, 0) NOT NULL,
    deadline        timestamptz,
    phase           job_phase    NOT NULL DEFAULT 'created',
    result_uri      text,
    block_number    bigint       NOT NULL,
    tx_hash         bytea        NOT NULL,
    log_index       integer      NOT NULL,
    created_at      timestamptz  NOT NULL DEFAULT now(),
    updated_at      timestamptz  NOT NULL DEFAULT now(),
    CONSTRAINT jobs_principal_len     CHECK (octet_length(principal) = 20),
    CONSTRAINT jobs_clone_len         CHECK (clone_address IS NULL OR octet_length(clone_address) = 20),
    CONSTRAINT jobs_token_len         CHECK (token IS NULL OR octet_length(token) = 20),
    CONSTRAINT jobs_tx_hash_len       CHECK (octet_length(tx_hash) = 32),
    CONSTRAINT jobs_log_index_nonneg  CHECK (log_index >= 0),
    CONSTRAINT jobs_job_type_nonneg   CHECK (job_type >= 0)
);

CREATE UNIQUE INDEX jobs_on_chain_id_uidx ON jobs (on_chain_id);
CREATE INDEX jobs_agent_pk_idx ON jobs (agent_pk) WHERE agent_pk IS NOT NULL;
CREATE INDEX jobs_agent_id_idx ON jobs (agent_id) WHERE agent_id IS NOT NULL;
CREATE INDEX jobs_principal_idx ON jobs (principal);
CREATE INDEX jobs_phase_idx ON jobs (phase);

CREATE TRIGGER jobs_set_updated_at
    BEFORE UPDATE ON jobs
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ---------------------------------------------------------------------------
-- job_events — append-only event log for a job (lifecycle + escrow).
--
-- Mirrors every on-chain event tied to a job so the API can render a
-- timeline and operators can audit phase transitions and escrow movements.
-- ---------------------------------------------------------------------------
CREATE TYPE job_event_kind AS ENUM (
    'created',
    'initialised',
    'funded',
    'accepted',
    'result_submitted',
    'completed',
    'cancelled',
    'disputed',
    'released',
    'refunded',
    'resolved'
);

CREATE TABLE job_events (
    id            bigserial      PRIMARY KEY,
    job_pk        bigint         REFERENCES jobs (id) ON DELETE CASCADE,
    on_chain_id   numeric(78, 0) NOT NULL,
    kind          job_event_kind NOT NULL,
    actor         bytea,
    amount        numeric(78, 0),
    token         bytea,
    payload       jsonb          NOT NULL DEFAULT '{}'::jsonb,
    block_number  bigint         NOT NULL,
    tx_hash       bytea          NOT NULL,
    log_index     integer        NOT NULL,
    created_at    timestamptz    NOT NULL DEFAULT now(),
    CONSTRAINT job_events_actor_len   CHECK (actor IS NULL OR octet_length(actor) = 20),
    CONSTRAINT job_events_token_len   CHECK (token IS NULL OR octet_length(token) = 20),
    CONSTRAINT job_events_tx_hash_len CHECK (octet_length(tx_hash) = 32),
    CONSTRAINT job_events_log_index_nonneg CHECK (log_index >= 0)
);

CREATE UNIQUE INDEX job_events_log_uidx ON job_events (tx_hash, log_index);
CREATE INDEX job_events_job_pk_block_idx ON job_events (job_pk, block_number DESC) WHERE job_pk IS NOT NULL;
CREATE INDEX job_events_on_chain_id_block_idx ON job_events (on_chain_id, block_number DESC);
CREATE INDEX job_events_kind_idx ON job_events (kind);

-- ---------------------------------------------------------------------------
-- memos — off-chain signed memos referenced by content hash.
--
-- The dispute-resolution event `DisputeResolved(jobId, arbiter, payoutBps,
-- reasonHash)` (see docs/security/threat-model-m3.md) emits a `reasonHash`
-- pointing at an off-chain memo. This table stores the resolved memo body
-- so the gateway can render reasoned outcomes without a separate fetch.
-- ---------------------------------------------------------------------------
CREATE TABLE memos (
    id            bigserial    PRIMARY KEY,
    content_hash  bytea        NOT NULL,
    job_pk        bigint       REFERENCES jobs (id) ON DELETE SET NULL,
    on_chain_id   numeric(78, 0),
    author        bytea,
    uri           text,
    body          text,
    signature     bytea,
    created_at    timestamptz  NOT NULL DEFAULT now(),
    CONSTRAINT memos_content_hash_len CHECK (octet_length(content_hash) = 32),
    CONSTRAINT memos_author_len       CHECK (author IS NULL OR octet_length(author) = 20),
    -- Either a uri or an inline body must be present. Storing both is allowed.
    CONSTRAINT memos_has_payload      CHECK (uri IS NOT NULL OR body IS NOT NULL)
);

CREATE UNIQUE INDEX memos_content_hash_uidx ON memos (content_hash);
CREATE INDEX memos_job_pk_idx ON memos (job_pk) WHERE job_pk IS NOT NULL;
CREATE INDEX memos_on_chain_id_idx ON memos (on_chain_id) WHERE on_chain_id IS NOT NULL;

COMMIT;
