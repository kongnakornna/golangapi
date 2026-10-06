-- PDPA audit trails table
-- ตาราง audit trail PDPA
CREATE TABLE IF NOT EXISTS pdpa_audit_trails (
    id         UUID PRIMARY KEY,
    user_id    UUID,
    action     VARCHAR(64) NOT NULL,
    details    JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pdpa_audit_user_id ON pdpa_audit_trails(user_id);
CREATE INDEX IF NOT EXISTS idx_pdpa_audit_action ON pdpa_audit_trails(action);
CREATE INDEX IF NOT EXISTS idx_pdpa_audit_created_at ON pdpa_audit_trails(created_at);
CREATE INDEX IF NOT EXISTS idx_pdpa_audit_details ON pdpa_audit_trails USING GIN (details);