-- PDPA DSAR requests table
-- ตารางคำร้อง DSAR PDPA
CREATE TABLE IF NOT EXISTS pdpa_dsar_requests (
    id               UUID PRIMARY KEY,
    user_id          UUID NOT NULL,
    request_type     SMALLINT NOT NULL,
    status           SMALLINT NOT NULL DEFAULT 1,
    requested_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at     TIMESTAMPTZ,
    data_payload     BYTEA,
    rejection_reason TEXT,
    otp_code         VARCHAR(16),
    otp_expired_at   TIMESTAMPTZ,
    ip_address       VARCHAR(45),
    user_agent       TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pdpa_dsar_user_id ON pdpa_dsar_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_pdpa_dsar_status ON pdpa_dsar_requests(status);
CREATE INDEX IF NOT EXISTS idx_pdpa_dsar_rate_limit
    ON pdpa_dsar_requests(user_id, request_type, requested_at);