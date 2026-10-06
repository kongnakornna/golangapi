-- PDPA consents table
-- ตารางความยินยอม PDPA
CREATE TABLE IF NOT EXISTS pdpa_consents (
    id              UUID PRIMARY KEY,
    user_id         UUID NOT NULL,
    purpose         VARCHAR(64) NOT NULL,
    status          SMALLINT NOT NULL DEFAULT 1,
    consented_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at      TIMESTAMPTZ NOT NULL,
    revoked_at      TIMESTAMPTZ,
    auto_deleted_at TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pdpa_consents_user_id ON pdpa_consents(user_id);
CREATE INDEX IF NOT EXISTS idx_pdpa_consents_user_purpose ON pdpa_consents(user_id, purpose);
CREATE INDEX IF NOT EXISTS idx_pdpa_consents_status ON pdpa_consents(status);
CREATE INDEX IF NOT EXISTS idx_pdpa_consents_active
    ON pdpa_consents(user_id, purpose, status) WHERE status = 1;