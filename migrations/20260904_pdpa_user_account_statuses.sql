-- PDPA user account statuses table
-- ตารางสถานะบัญชีผู้ใช้ PDPA
CREATE TABLE IF NOT EXISTS pdpa_user_account_statuses (
    id                    UUID PRIMARY KEY,
    user_id               UUID NOT NULL UNIQUE,
    status                SMALLINT NOT NULL DEFAULT 1,
    suspended_at          TIMESTAMPTZ,
    terminated_at         TIMESTAMPTZ,
    deletion_confirmed_at TIMESTAMPTZ,
    retention_deadline    TIMESTAMPTZ,
    auto_deleted_at       TIMESTAMPTZ,
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);