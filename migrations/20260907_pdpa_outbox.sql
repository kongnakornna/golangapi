-- PDPA transactional outbox table
-- ตาราง transactional outbox PDPA
CREATE TABLE IF NOT EXISTS pdpa_outbox (
    id            UUID PRIMARY KEY,
    event_id      UUID NOT NULL,
    topic         VARCHAR(128) NOT NULL,
    key           VARCHAR(128),
    payload       JSONB NOT NULL,
    handler       VARCHAR(64) NOT NULL,
    status        SMALLINT NOT NULL DEFAULT 1,
    attempts      INT NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    dispatched_at TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_pdpa_outbox_claim
    ON pdpa_outbox(handler, status, created_at)
    WHERE status = 1 AND deleted_at IS NULL;