-- Seed initial privacy policy v1.0.0
INSERT INTO pdpa_privacy_policies (
    id, policy_key, version, title_en, title_th, content_en, content_th,
    status, effective_at, published_at
)
VALUES (
    gen_random_uuid(), 'main', 'v1.0.0',
    'Privacy Policy', 'นโยบายความเป็นส่วนตัว',
    'Initial privacy policy content.', 'เนื้อหานโยบายเริ่มต้น',
    1, NOW(), NOW()
)
ON CONFLICT DO NOTHING;