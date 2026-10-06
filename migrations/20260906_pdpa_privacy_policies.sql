-- PDPA privacy policies table
-- ตารางนโยบายความเป็นส่วนตัว PDPA
CREATE TABLE IF NOT EXISTS pdpa_privacy_policies (
    id           UUID PRIMARY KEY,
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

CREATE UNIQUE INDEX IF NOT EXISTS uq_pdpa_policy_key_version
    ON pdpa_privacy_policies(policy_key, version) WHERE deleted_at IS NULL;