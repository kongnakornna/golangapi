-- ============================================================
-- PDPA Module — Full Schema
-- สร้างตารางทั้งหมดของโมดูล PDPA
-- ============================================================

-- 1) Privacy Policies
CREATE TABLE IF NOT EXISTS pdpa_privacy_policies (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_key   VARCHAR(64) NOT NULL,
    version      VARCHAR(32) NOT NULL,
    title_en     VARCHAR(255) NOT NULL,
    title_th     VARCHAR(255) NOT NULL,
    content_en   TEXT,
    content_th   TEXT,
    status       BIGINT NOT NULL DEFAULT 0,
    effective_at TIMESTAMPTZ,
    published_at TIMESTAMPTZ,
    created_by   UUID,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_pdpa_policy_key    ON pdpa_privacy_policies(policy_key);
CREATE INDEX IF NOT EXISTS idx_pdpa_policy_status ON pdpa_privacy_policies(status);
CREATE UNIQUE INDEX IF NOT EXISTS uq_pdpa_policy_key_version
    ON pdpa_privacy_policies(policy_key, version) WHERE deleted_at IS NULL;

-- 2) Purposes
CREATE TABLE IF NOT EXISTS pdpa_purposes (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    purpose_code     VARCHAR(64) NOT NULL,
    purpose_type     BIGINT NOT NULL DEFAULT 1,
    name_en          VARCHAR(255) NOT NULL,
    name_th          VARCHAR(255) NOT NULL,
    description_en   TEXT,
    description_th   TEXT,
    consent_required BOOLEAN NOT NULL DEFAULT TRUE,
    retention_days   INT,
    status           BIGINT NOT NULL DEFAULT 1,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_pdpa_purpose_code ON pdpa_purposes(purpose_code);
CREATE INDEX IF NOT EXISTS idx_pdpa_purpose_status ON pdpa_purposes(status);

-- 3) Consents
CREATE TABLE IF NOT EXISTS pdpa_consents (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    purpose_code    VARCHAR(64) NOT NULL,
    status          SMALLINT NOT NULL DEFAULT 1,
    consented_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at      TIMESTAMPTZ NOT NULL,
    revoked_at      TIMESTAMPTZ,
    auto_deleted_at TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pdpa_consent_user_id  ON pdpa_consents(user_id);
CREATE INDEX IF NOT EXISTS idx_pdpa_consent_purpose  ON pdpa_consents(purpose_code);
CREATE INDEX IF NOT EXISTS idx_pdpa_consent_status   ON pdpa_consents(status);

-- 4) DSAR Requests
CREATE TABLE IF NOT EXISTS pdpa_dsar_requests (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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

CREATE INDEX IF NOT EXISTS idx_pdpa_user_request_user_id ON pdpa_dsar_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_pdpa_user_request_status  ON pdpa_dsar_requests(status);

-- 5) User Account Status
CREATE TABLE IF NOT EXISTS pdpa_user_account_statuses (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id               UUID NOT NULL,
    status                SMALLINT NOT NULL DEFAULT 1,
    suspended_at          TIMESTAMPTZ,
    terminated_at         TIMESTAMPTZ,
    deletion_confirmed_at TIMESTAMPTZ,
    retention_deadline    TIMESTAMPTZ,
    auto_deleted_at       TIMESTAMPTZ,
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_pdpa_account_user_id ON pdpa_user_account_statuses(user_id);
CREATE INDEX IF NOT EXISTS idx_pdpa_account_status ON pdpa_user_account_statuses(status);

-- 6) Audit Trail
CREATE TABLE IF NOT EXISTS pdpa_audit_trails (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID,
    action     VARCHAR(64) NOT NULL,
    details    JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pdpa_audit_user_id    ON pdpa_audit_trails(user_id);
CREATE INDEX IF NOT EXISTS idx_pdpa_audit_action     ON pdpa_audit_trails(action);
CREATE INDEX IF NOT EXISTS idx_pdpa_audit_created_at ON pdpa_audit_trails(created_at);
CREATE INDEX IF NOT EXISTS idx_pdpa_audit_details    ON pdpa_audit_trails USING GIN (details);

-- 7) Request Responses
CREATE TABLE IF NOT EXISTS pdpa_request_responses (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id        UUID NOT NULL,
    user_id           UUID NOT NULL,
    response_type     BIGINT NOT NULL,
    file_path         TEXT,
    file_hash         VARCHAR(128),
    payload           JSONB,
    status            BIGINT NOT NULL DEFAULT 1,
    delivered_at      TIMESTAMPTZ,
    delivered_channel VARCHAR(32),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_pdpa_response_request_id ON pdpa_request_responses(request_id);
CREATE INDEX IF NOT EXISTS idx_pdpa_response_user_id    ON pdpa_request_responses(user_id);

-- ============================================================
-- Seed Data
-- ============================================================

-- Purposes
INSERT INTO pdpa_purposes (purpose_code, purpose_type, name_en, name_th, consent_required)
VALUES
    ('NECESSARY',           1, 'Necessary',           'จำเป็นสำหรับการให้บริการ', TRUE),
    ('ACCOUNT_SYSTEM',      2, 'Account System',      'ระบบบัญชี',              TRUE),
    ('TRANSACTION_HISTORY', 3, 'Transaction History', 'ประวัติธุรกรรม',          TRUE),
    ('USAGE_LOGS',          4, 'Usage Logs',          'บันทึกการใช้งาน',        FALSE),
    ('MARKETING',           5, 'Marketing',           'การตลาด',                FALSE),
    ('ANALYTICS',           6, 'Analytics',           'การวิเคราะห์',           FALSE)
ON CONFLICT (purpose_code) DO NOTHING;

-- Initial Policy
INSERT INTO pdpa_privacy_policies (
    policy_key, version, title_en, title_th,
    content_en, content_th, status, effective_at, published_at
)
VALUES (
    'main', 'v1.0.0',
    'Privacy Policy', 'นโยบายความเป็นส่วนตัว',
    'Initial privacy policy content.', 'เนื้อหานโยบายเริ่มต้น',
    1, NOW(), NOW()
)
ON CONFLICT DO NOTHING;