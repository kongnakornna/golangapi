-- PDPA purposes master data
-- ตารางข้อมูลมาสเตอร์วัตถุประสงค์
CREATE TABLE IF NOT EXISTS pdpa_purposes (
    id               UUID PRIMARY KEY,
    purpose_code     VARCHAR(64) NOT NULL UNIQUE,
    purpose_type     BIGINT NOT NULL,
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

INSERT INTO pdpa_purposes (id, purpose_code, purpose_type, name_en, name_th, consent_required)
VALUES
    (gen_random_uuid(), 'NECESSARY',           1, 'Necessary',           'จำเป็นสำหรับการให้บริการ', TRUE),
    (gen_random_uuid(), 'ACCOUNT_SYSTEM',      2, 'Account System',      'ระบบบัญชี',              TRUE),
    (gen_random_uuid(), 'TRANSACTION_HISTORY', 3, 'Transaction History', 'ประวัติธุรกรรม',          TRUE),
    (gen_random_uuid(), 'USAGE_LOGS',          4, 'Usage Logs',          'บันทึกการใช้งาน',        FALSE),
    (gen_random_uuid(), 'MARKETING',           5, 'Marketing',           'การตลาด',                FALSE),
    (gen_random_uuid(), 'ANALYTICS',           6, 'Analytics',           'การวิเคราะห์',           FALSE)
ON CONFLICT (purpose_code) DO NOTHING;