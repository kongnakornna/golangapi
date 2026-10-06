-- PDPA processed events for idempotency
CREATE TABLE IF NOT EXISTS pdpa_processed_events (
    event_id     UUID NOT NULL,
    handler_name VARCHAR(128) NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at   TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (event_id, handler_name)
);

CREATE INDEX IF NOT EXISTS idx_pdpa_processed_events_expires_at
    ON pdpa_processed_events(expires_at);