# 🌱 PART 9 / 7 — SAMPLE DATA & FIXTURES

> **ขนาด**: ใหญ่ — แยก 6 ตอนย่อย
> **Part 9A**: Fixture Strategy & Conventions
> **Part 9B**: SQL Seed Data (comprehensive)
> **Part 9C**: Go Fixture Builders (programmatic)
> **Part 9D**: Test Scenarios (integration + e2e)
> **Part 9E**: Docker/CI Integration
> **Part 9F**: Postman Environment & Pre-request Scripts

> **เป้าหมาย**: ให้ dev/test/demo รันได้ทันทีด้วย `make seed` โดยมีข้อมูลครอบคลุม 6 modules

---

## 🅰️ PART 9A — FIXTURE STRATEGY & CONVENTIONS

### A.1 Fixture Strategy

```
test/fixtures/
├── sql/                        # SQL files — idempotent
│   ├── 00_cleanup.sql          # ลบข้อมูลเก่า (dev only)
│   ├── 01_tenants.sql
│   ├── 02_users_auth.sql
│   ├── 03_customers.sql
│   ├── 04_sites_zones.sql
│   ├── 05_packages.sql
│   ├── 06_subscriptions.sql
│   ├── 07_products_warehouses.sql
│   ├── 08_inventory.sql
│   ├── 09_devices.sql
│   ├── 10_orders.sql
│   ├── 11_invoices_payments.sql
│   ├── 12_technicians_logistics.sql
│   ├── 13_installation_jobs.sql
│   ├── 14_kpi_snapshots.sql
│   ├── 15_dashboards.sql
│   └── 16_telemetry_influx.flux  # InfluxDB setup
├── go/
│   ├── builders.go             # Fluent builders
│   ├── datasets.go             # Pre-built scenarios
│   └── loader.go               # Bulk insert
├── postman/
│   ├── collection.json         # API collection
│   ├── environment.dev.json
│   └── environment.prod.json
└── README.md
```

### A.2 Conventions

| Rule | ตัวอย่าง |
|---|---|
| **UUID** | คงที่ (deterministic) — ใช้ prefix ตาม entity |
| **Tenant IDs** | `11111111-1111-1111-1111-111111111101` (demo) |
| **Customer Codes** | `CUS-2026-0001` เป็นต้นไป |
| **Emails** | `demo+customer1@icmongolang.local` |
| **Timestamps** | relative to NOW() — ลบ/บวกวัน |
| **Idempotent** | `INSERT ... ON CONFLICT DO NOTHING` |
| **Marked** | `metadata->>'fixture' = 'true'` |
| **Cleanup** | `DELETE WHERE metadata @> '{"fixture":true}'` |

### A.3 Deterministic UUID Convention

```go
// test/fixtures/go/uuid.go
package fixtures

import "github.com/google/uuid"

// namespace for deterministic UUIDs (UUIDv5)
var FixtureNamespace = uuid.MustParse("f1c7a5b0-0000-0000-0000-000000000000")

// TenantID generates stable UUID per (tenant_slug, entity, key)
func TenantID(slug string) uuid.UUID {
    return uuid.NewSHA1(FixtureNamespace, []byte("tenant:"+slug))
}

func CustomerID(tenantSlug, code string) uuid.UUID {
    return uuid.NewSHA1(FixtureNamespace, []byte("customer:"+tenantSlug+":"+code))
}

func PackageID(code string) uuid.UUID {
    return uuid.NewSHA1(FixtureNamespace, []byte("package:"+code))
}

func DeviceID(tenantSlug, serial string) uuid.UUID {
    return uuid.NewSHA1(FixtureNamespace, []byte("device:"+tenantSlug+":"+serial))
}

func ProductID(sku string) uuid.UUID {
    return uuid.NewSHA1(FixtureNamespace, []byte("product:"+sku))
}
```

### A.4 Fixture Scope Levels

| Level | Purpose | Command |
|---|---|---|
| **minimal** | แค่ที่จำเป็นสำหรับ smoke test | `make seed-minimal` |
| **standard** | ข้อมูลครบทุก module | `make seed` (default) |
| **demo** | ข้อมูลสวย + scenarios | `make seed-demo` |
| **stress** | ข้อมูลเยอะสำหรับ load test | `make seed-stress` |

---

## 🅱️ PART 9B — SQL SEED DATA

### B.1 `test/fixtures/sql/00_cleanup.sql`

```sql
-- ============================================================
-- Cleanup fixtures (dev only)
-- รันก่อน seed เพื่อ reset
-- ============================================================
BEGIN;

-- ลบตามลำดับ FK
DELETE FROM erp_journal_lines WHERE journal_id IN (
    SELECT id FROM erp_journal_entries WHERE tenant_id IN (
        SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
    )
);
DELETE FROM erp_journal_entries WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);

DELETE FROM erp_payment_allocations WHERE payment_id IN (
    SELECT id FROM erp_payments WHERE tenant_id IN (
        SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
    )
);
DELETE FROM erp_payments WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);

DELETE FROM erp_invoice_lines WHERE invoice_id IN (
    SELECT id FROM erp_invoices WHERE tenant_id IN (
        SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
    )
);
DELETE FROM erp_invoices WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);

DELETE FROM erp_order_lines WHERE order_id IN (
    SELECT id FROM erp_orders WHERE tenant_id IN (
        SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
    )
);
DELETE FROM erp_orders WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);

DELETE FROM erp_stock_movements WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);
DELETE FROM erp_stock_reservations WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);
DELETE FROM erp_inventory_items WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);

DELETE FROM erp_products WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);
DELETE FROM erp_warehouses WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);
DELETE FROM erp_suppliers WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);

DELETE FROM device_alert_events WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);
DELETE FROM device_alert_rules WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);
DELETE FROM device_commands WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);
DELETE FROM device_shadows WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);
DELETE FROM device_devices WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);
DELETE FROM device_automations WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);

DELETE FROM iotlogistics_installation_jobs WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);
DELETE FROM iotlogistics_shipment_items WHERE shipment_id IN (
    SELECT id FROM iotlogistics_shipments WHERE tenant_id IN (
        SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
    )
);
DELETE FROM iotlogistics_shipments WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);
DELETE FROM iotlogistics_technicians WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);
DELETE FROM iotlogistics_work_orders WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);
DELETE FROM iotlogistics_maintenance_schedules WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);

DELETE FROM packagecatalog_subscription_history WHERE subscription_id IN (
    SELECT id FROM packagecatalog_subscriptions WHERE tenant_id IN (
        SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
    )
);
DELETE FROM packagecatalog_subscriptions WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);

DELETE FROM customer_contracts WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);
DELETE FROM customer_sites WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);
DELETE FROM customer_contacts WHERE customer_id IN (
    SELECT id FROM customer_customers WHERE tenant_id IN (
        SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
    )
);
DELETE FROM customer_customers WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);

DELETE FROM report_kpi_snapshots WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);
DELETE FROM report_insights WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);
DELETE FROM report_dashboards WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);

DELETE FROM auth_sessions WHERE user_id IN (
    SELECT id FROM users_users WHERE tenant_id IN (
        SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
    )
);
DELETE FROM users_users WHERE tenant_id IN (
    SELECT id FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1')
);

DELETE FROM users_tenants WHERE slug IN ('demo', 'acme', 'farm1');

COMMIT;
```

### B.2 `test/fixtures/sql/01_tenants.sql`

```sql
-- ============================================================
-- Tenants
-- ============================================================
BEGIN;

INSERT INTO users_tenants (id, slug, name, status, plan, metadata, created_at)
VALUES
    ('11111111-1111-1111-1111-111111111101', 'demo',
     'Demo Tenant', 'ACTIVE', 'PRO',
     '{"fixture": true, "locale": "th-TH", "timezone": "Asia/Bangkok"}',
     NOW() - INTERVAL '90 days'),

    ('11111111-1111-1111-1111-111111111102', 'acme',
     'ACME IoT Solutions Co., Ltd.', 'ACTIVE', 'ENTERPRISE',
     '{"fixture": true, "locale": "th-TH", "tax_id": "0105555555555"}',
     NOW() - INTERVAL '180 days'),

    ('11111111-1111-1111-1111-111111111103', 'farm1',
     'สมาร์ทฟาร์ม เชียงใหม่', 'ACTIVE', 'BASIC',
     '{"fixture": true, "locale": "th-TH", "region": "north"}',
     NOW() - INTERVAL '30 days')
ON CONFLICT (id) DO NOTHING;

COMMIT;
```

### B.3 `test/fixtures/sql/02_users_auth.sql`

```sql
-- ============================================================
-- Users + Auth
-- Password ทุก user: "Password123!"
-- ============================================================
BEGIN;

INSERT INTO users_users (id, tenant_id, email, password_hash, full_name, role, status, metadata, created_at)
VALUES
    -- Demo tenant
    ('22222222-2222-2222-2222-222222222201',
     '11111111-1111-1111-1111-111111111101',
     'admin@demo.local',
     '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', -- "Password123!"
     'Demo Admin', 'ADMIN', 'ACTIVE',
     '{"fixture": true}', NOW() - INTERVAL '90 days'),

    ('22222222-2222-2222-2222-222222222202',
     '11111111-1111-1111-1111-111111111101',
     'manager@demo.local',
     '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
     'Demo Manager', 'MANAGER', 'ACTIVE',
     '{"fixture": true}', NOW() - INTERVAL '85 days'),

    ('22222222-2222-2222-2222-222222222203',
     '11111111-1111-1111-1111-111111111101',
     'viewer@demo.local',
     '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
     'Demo Viewer', 'VIEWER', 'ACTIVE',
     '{"fixture": true}', NOW() - INTERVAL '80 days'),

    -- ACME tenant
    ('22222222-2222-2222-2222-222222222211',
     '11111111-1111-1111-1111-111111111102',
     'admin@acme.co.th',
     '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
     'ACME Admin', 'ADMIN', 'ACTIVE',
     '{"fixture": true}', NOW() - INTERVAL '180 days'),

    ('22222222-2222-2222-2222-222222222212',
     '11111111-1111-1111-1111-111111111102',
     'sales@acme.co.th',
     '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
     'ACME Sales', 'MANAGER', 'ACTIVE',
     '{"fixture": true}', NOW() - INTERVAL '170 days'),

    -- Farm1 tenant
    ('22222222-2222-2222-2222-222222222221',
     '11111111-1111-1111-1111-111111111103',
     'owner@farm1.local',
     '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
     'สมชาย ใจดี', 'ADMIN', 'ACTIVE',
     '{"fixture": true}', NOW() - INTERVAL '30 days')
ON CONFLICT (id) DO NOTHING;

COMMIT;
```

### B.4 `test/fixtures/sql/03_customers.sql`

```sql
-- ============================================================
-- Customers (10 ตัว)
-- ============================================================
BEGIN;

INSERT INTO customer_customers
    (id, tenant_id, code, name, type, segment, lifecycle, status,
     tax_id, email, phone, address, credit_limit, currency, metadata, created_by, created_at)
VALUES
    -- Demo tenant (4 customers)
    ('33333333-3333-3333-3333-333333333301',
     '11111111-1111-1111-1111-111111111101',
     'CUS-2026-0001', 'บริษัท ไทยเกษตร จำกัด', 'CORPORATE', 'ENTERPRISE',
     'CUSTOMER', 'ACTIVE',
     '0105550000001', 'contact@thaikaset.co.th', '021234567',
     '{"line1":"99/1 อาคารสาทร","district":"สาทร","province":"กรุงเทพ","postcode":"10120","country":"TH"}',
     500000, 'THB', '{"fixture": true, "industry": "agriculture"}',
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '60 days'),

    ('33333333-3333-3333-3333-333333333302',
     '11111111-1111-1111-1111-111111111101',
     'CUS-2026-0002', 'นายสมชาย ใจดี', 'INDIVIDUAL', 'INDIVIDUAL',
     'CUSTOMER', 'ACTIVE',
     NULL, 'somchai.j@example.com', '0812345678',
     '{"line1":"123 ม.5 ต.สันทราย","district":"สันทราย","province":"เชียงใหม่","postcode":"50210","country":"TH"}',
     0, 'THB', '{"fixture": true}',
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '45 days'),

    ('33333333-3333-3333-3333-333333333303',
     '11111111-1111-1111-1111-111111111101',
     'CUS-2026-0003', 'บริษัท สมาร์ทบิลดิ้ง จก.', 'CORPORATE', 'SME',
     'PROSPECT', 'PENDING',
     '0105550000003', 'info@smartbuilding.co.th', '029876543',
     '{"line1":"45 อาคารเพลินจิต","district":"คลองเตย","province":"กรุงเทพ","postcode":"10110","country":"TH"}',
     200000, 'THB', '{"fixture": true}',
     '22222222-2222-2222-2222-222222222202',
     NOW() - INTERVAL '10 days'),

    ('33333333-3333-3333-3333-333333333304',
     '11111111-1111-1111-1111-111111111101',
     'CUS-2026-0004', 'ชลประทานสยาม', 'GOVERNMENT', 'GOVERNMENT',
     'CUSTOMER', 'ACTIVE',
     '0994000000004', 'contact@irrigation.go.th', '025555555',
     '{"line1":"99 ถ.พญาไท","district":"ราชเทวี","province":"กรุงเทพ","postcode":"10400","country":"TH"}',
     1000000, 'THB', '{"fixture": true, "department": "agriculture"}',
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '90 days'),

    -- ACME tenant (4 customers)
    ('33333333-3333-3333-3333-333333333311',
     '11111111-1111-1111-1111-111111111102',
     'CUS-2026-0011', 'SCG Smart Farm', 'CORPORATE', 'ENTERPRISE',
     'CUSTOMER', 'ACTIVE',
     '0105550000011', 'smartfarm@scg.co.th', '029999999',
     '{"line1":"1 ถ.ปูนซิเมนต์ไทย","district":"บางซื่อ","province":"กรุงเทพ","postcode":"10800","country":"TH"}',
     5000000, 'THB', '{"fixture": true}',
     '22222222-2222-2222-2222-222222222211',
     NOW() - INTERVAL '150 days'),

    ('33333333-3333-3333-3333-333333333312',
     '11111111-1111-1111-1111-111111111102',
     'CUS-2026-0012', 'CP All Distribution', 'CORPORATE', 'ENTERPRISE',
     'CUSTOMER', 'ACTIVE',
     '0105550000012', 'iot@cpall.co.th', '021111111',
     '{"line1":"313 อาคารซีพีทาวเวอร์","district":"สีลม","province":"กรุงเทพ","postcode":"10500","country":"TH"}',
     3000000, 'THB', '{"fixture": true}',
     '22222222-2222-2222-2222-222222222211',
     NOW() - INTERVAL '140 days'),

    ('33333333-3333-3333-3333-333333333313',
     '11111111-1111-1111-1111-111111111102',
     'CUS-2026-0013', 'บริษัท เมืองทองธานี จก.', 'CORPORATE', 'SME',
     'CUSTOMER', 'SUSPENDED',
     '0105550000013', 'contact@mt.co.th', '025555111',
     '{"line1":"47/569 ปากเกร็ด","district":"ปากเกร็ด","province":"นนทบุรี","postcode":"11120","country":"TH"}',
     500000, 'THB', '{"fixture": true, "suspend_reason": "payment overdue 45 days"}',
     '22222222-2222-2222-2222-222222222211',
     NOW() - INTERVAL '100 days'),

    ('33333333-3333-3333-3333-333333333314',
     '11111111-1111-1111-1111-111111111102',
     'CUS-2026-0014', 'ร้านกาแฟ ABC', 'CORPORATE', 'STARTUP',
     'LEAD', 'PENDING',
     NULL, 'abc@coffee.local', '085555555',
     '{"line1":"88 ถ.สุขุมวิท","district":"วัฒนา","province":"กรุงเทพ","postcode":"10110","country":"TH"}',
     0, 'THB', '{"fixture": true, "source": "web_inquiry"}',
     '22222222-2222-2222-2222-222222222212',
     NOW() - INTERVAL '5 days'),

    -- Farm1 tenant (2 customers)
    ('33333333-3333-3333-3333-333333333321',
     '11111111-1111-1111-1111-111111111103',
     'CUS-2026-0021', 'ฟาร์มผักไฮโดร 1', 'CORPORATE', 'SME',
     'CUSTOMER', 'ACTIVE',
     '0505550000021', 'info@hydro1.co.th', '053111111',
     '{"line1":"55 ม.3 ต.สันกำแพง","district":"สันกำแพง","province":"เชียงใหม่","postcode":"50130","country":"TH"}',
     100000, 'THB', '{"fixture": true, "crop": "lettuce"}',
     '22222222-2222-2222-2222-222222222221',
     NOW() - INTERVAL '25 days'),

    ('33333333-3333-3333-3333-333333333322',
     '11111111-1111-1111-1111-111111111103',
     'CUS-2026-0022', 'สวนมะม่วงน้ำดอกไม้', 'INDIVIDUAL', 'INDIVIDUAL',
     'CUSTOMER', 'ACTIVE',
     NULL, 'mango@farm.local', '089999999',
     '{"line1":"22 ม.8 ต.ป่าซาง","district":"ป่าซาง","province":"ลำพูน","postcode":"51120","country":"TH"}',
     0, 'THB', '{"fixture": true, "crop": "mango"}',
     '22222222-2222-2222-2222-222222222221',
     NOW() - INTERVAL '20 days')
ON CONFLICT (id) DO NOTHING;

-- Contacts
INSERT INTO customer_contacts (id, customer_id, name, phone, email, position, is_primary, created_at)
VALUES
    ('44444444-4444-4444-4444-444444444401', '33333333-3333-3333-3333-333333333301',
     'คุณวิชัย รักเกษตร', '0811111111', 'wichai@thaikaset.co.th', 'CEO', TRUE, NOW() - INTERVAL '60 days'),
    ('44444444-4444-4444-4444-444444444402', '33333333-3333-3333-3333-333333333301',
     'คุณมาลี สวยงาม', '0822222222', 'malee@thaikaset.co.th', 'CTO', FALSE, NOW() - INTERVAL '58 days'),
    ('44444444-4444-4444-4444-444444444403', '33333333-3333-3333-3333-333333333311',
     'คุณสมศักดิ์ ชัยชนะ', '0833333333', 'somsak@scg.co.th', 'Director', TRUE, NOW() - INTERVAL '150 days')
ON CONFLICT (id) DO NOTHING;

COMMIT;
```

### B.5 `test/fixtures/sql/04_sites_zones.sql`

```sql
-- ============================================================
-- Sites + Zones (Customer + Device modules)
-- ============================================================
BEGIN;

-- Customer sites
INSERT INTO customer_sites (id, tenant_id, customer_id, type, name, geo_lat, geo_lng, address, area_size, metadata, is_active, created_at)
VALUES
    -- Demo
    ('55555555-5555-5555-5555-555555555501',
     '11111111-1111-1111-1111-111111111101', '33333333-3333-3333-3333-333333333301',
     'FARM', 'ไทยเกษตร ฟาร์มสาขาหลัก', 13.7563, 100.5018,
     '{"line1":"99/1","province":"ปทุมธานี","country":"TH"}', 50000,
     '{"fixture": true, "area_unit": "sqm"}', TRUE, NOW() - INTERVAL '60 days'),

    ('55555555-5555-5555-5555-555555555502',
     '11111111-1111-1111-1111-111111111101', '33333333-3333-3333-3333-333333333302',
     'GREENHOUSE', 'โรงเรือนสมชาย', 18.7883, 98.9853,
     '{"line1":"123 ม.5","province":"เชียงใหม่","country":"TH"}', 1200,
     '{"fixture": true}', TRUE, NOW() - INTERVAL '45 days'),

    ('55555555-5555-5555-5555-555555555503',
     '11111111-1111-1111-1111-111111111101', '33333333-3333-3333-3333-333333333304',
     'FARM', 'สถานีชลประทาน', 14.9799, 102.0977,
     '{"line1":"99","province":"นครราชสีมา","country":"TH"}', 100000,
     '{"fixture": true}', TRUE, NOW() - INTERVAL '90 days'),

    -- ACME
    ('55555555-5555-5555-5555-555555555511',
     '11111111-1111-1111-1111-111111111102', '33333333-3333-3333-3333-333333333311',
     'FARM', 'SCG ไร่อ้อย 1', 15.2287, 104.8574,
     '{"line1":"1","province":"ศรีสะเกษ","country":"TH"}', 200000,
     '{"fixture": true}', TRUE, NOW() - INTERVAL '150 days'),

    ('55555555-5555-5555-5555-555555555512',
     '11111111-1111-1111-1111-111111111102', '33333333-3333-3333-3333-333333333312',
     'WAREHOUSE', 'CP คลังสินค้า', 13.7055, 100.5490,
     '{"line1":"313","province":"กรุงเทพ","country":"TH"}', 15000,
     '{"fixture": true}', TRUE, NOW() - INTERVAL '140 days'),

    ('55555555-5555-5555-5555-555555555513',
     '11111111-1111-1111-1111-111111111102', '33333333-3333-3333-3333-333333333313',
     'BUILDING', 'เมืองทอง คอนโด A', 13.9126, 100.5083,
     '{"line1":"47/569","province":"นนทบุรี","country":"TH"}', 3000,
     '{"fixture": true}', TRUE, NOW() - INTERVAL '100 days'),

    -- Farm1
    ('55555555-5555-5555-5555-555555555521',
     '11111111-1111-1111-1111-111111111103', '33333333-3333-3333-3333-333333333321',
     'GREENHOUSE', 'โรงเรือนไฮโดร 1', 18.8134, 99.1267,
     '{"line1":"55","province":"เชียงใหม่","country":"TH"}', 800,
     '{"fixture": true}', TRUE, NOW() - INTERVAL '25 days'),

    ('55555555-5555-5555-5555-555555555522',
     '11111111-1111-1111-1111-111111111103', '33333333-3333-3333-3333-333333333322',
     'FARM', 'สวนมะม่วง', 18.5239, 98.9534,
     '{"line1":"22","province":"ลำพูน","country":"TH"}', 20000,
     '{"fixture": true}', TRUE, NOW() - INTERVAL '20 days')
ON CONFLICT (id) DO NOTHING;

-- Device sites (สำหรับ device module — link กับ customer site)
INSERT INTO device_sites (id, tenant_id, customer_id, name, type, geo_lat, geo_lng, is_active, metadata, created_at)
SELECT
    s.id, s.tenant_id, s.customer_id, s.name, s.type, s.geo_lat, s.geo_lng, s.is_active,
    s.metadata, s.created_at
FROM customer_sites s
WHERE s.metadata @> '{"fixture": true}'::jsonb
ON CONFLICT (id) DO NOTHING;

-- Zones
INSERT INTO device_zones (id, tenant_id, site_id, name, zone_type, metadata, created_at)
VALUES
    ('66666666-6666-6666-6666-666666666601',
     '11111111-1111-1111-1111-111111111101', '55555555-5555-5555-5555-555555555501',
     'โซน A - ผักใบ', 'ZONE', '{"fixture": true}', NOW() - INTERVAL '60 days'),
    ('66666666-6666-6666-6666-666666666602',
     '11111111-1111-1111-1111-111111111101', '55555555-5555-5555-5555-555555555501',
     'โซน B - ผลไม้', 'ZONE', '{"fixture": true}', NOW() - INTERVAL '60 days'),
    ('66666666-6666-6666-6666-666666666603',
     '11111111-1111-1111-1111-111111111101', '55555555-5555-5555-5555-555555555502',
     'โรงเรือน 1', 'GREENHOUSE', '{"fixture": true}', NOW() - INTERVAL '45 days'),
    ('66666666-6666-6666-6666-666666666611',
     '11111111-1111-1111-1111-111111111102', '55555555-5555-5555-5555-555555555511',
     'แปลง 1', 'ZONE', '{"fixture": true}', NOW() - INTERVAL '150 days'),
    ('66666666-6666-6666-6666-666666666621',
     '11111111-1111-1111-1111-111111111103', '55555555-5555-5555-5555-555555555521',
     'ชั้น 1', 'FLOOR', '{"fixture": true}', NOW() - INTERVAL '25 days')
ON CONFLICT (id) DO NOTHING;

COMMIT;
```

### B.6 `test/fixtures/sql/05_packages.sql`

```sql
-- ============================================================
-- Packages (platform-level + tenant custom)
-- ============================================================
BEGIN;

INSERT INTO packagecatalog_packages
    (id, tenant_id, code, name, description, tier, billing_cycle, currency,
     price_amount, version, price_history, quotas, features, trial_policy,
     is_active, is_public, is_deprecated, sort_order, metadata, tags, created_by, created_at)
VALUES
    -- Platform packages (tenant_id = NULL)
    ('77777777-7777-7777-7777-777777777701', NULL,
     'FREE', 'Free Plan', 'ทดลองใช้ฟรี 5 devices', 'FREE', 'MONTHLY', 'THB', 0, 1,
     '[{"version":1,"price":{"amount":0,"currency":"THB"},"effective_at":"2026-01-01T00:00:00Z","changed_by":"00000000-0000-0000-0000-000000000000","changed_at":"2026-01-01T00:00:00Z"}]',
     '{"max_devices":5,"max_sites":1,"max_users":2,"max_automations":5,"storage_gb":1,"retention_days":7,"api_calls_per_day":1000,"telemetry_points_per_day":1,"max_dashboards":1,"max_reports":5}',
     '[{"code":"DASHBOARD_BASIC","name":"Basic Dashboard","enabled":true}]',
     '{"enabled":false,"days":0,"auto_convert":false,"require_payment":false,"trial_tier":"FREE"}',
     TRUE, TRUE, FALSE, 1, '{"fixture": true}', '["free","entry"]',
     '00000000-0000-0000-0000-000000000000', NOW() - INTERVAL '365 days'),

    ('77777777-7777-7777-7777-777777777702', NULL,
     'BASIC-M', 'Basic Monthly', 'Basic features + AI', 'BASIC', 'MONTHLY', 'THB', 990, 1,
     '[{"version":1,"price":{"amount":990,"currency":"THB"},"effective_at":"2026-01-01T00:00:00Z","changed_by":"00000000-0000-0000-0000-000000000000","changed_at":"2026-01-01T00:00:00Z"}]',
     '{"max_devices":50,"max_sites":5,"max_users":10,"max_automations":50,"storage_gb":10,"retention_days":30,"api_calls_per_day":10000,"telemetry_points_per_day":10,"max_dashboards":5,"max_reports":50}',
     '[{"code":"DASHBOARD_BASIC","name":"Basic Dashboard","enabled":true},{"code":"AI_BASIC","name":"Basic AI","enabled":true},{"code":"API_ACCESS","name":"API Access","enabled":true}]',
     '{"enabled":true,"days":14,"auto_convert":false,"require_payment":false,"trial_tier":"BASIC"}',
     TRUE, TRUE, FALSE, 10, '{"fixture": true}', '["basic","monthly"]',
     '00000000-0000-0000-0000-000000000000', NOW() - INTERVAL '365 days'),

    ('77777777-7777-7777-7777-777777777703', NULL,
     'BASIC-Y', 'Basic Yearly', 'Basic — ประหยัด 2 เดือน', 'BASIC', 'YEARLY', 'THB', 9900, 1,
     '[{"version":1,"price":{"amount":9900,"currency":"THB"},"effective_at":"2026-01-01T00:00:00Z","changed_by":"00000000-0000-0000-0000-000000000000","changed_at":"2026-01-01T00:00:00Z"}]',
     '{"max_devices":50,"max_sites":5,"max_users":10,"max_automations":50,"storage_gb":10,"retention_days":30,"api_calls_per_day":10000,"telemetry_points_per_day":10,"max_dashboards":5,"max_reports":50}',
     '[{"code":"DASHBOARD_BASIC","name":"Basic Dashboard","enabled":true},{"code":"AI_BASIC","name":"Basic AI","enabled":true},{"code":"API_ACCESS","name":"API Access","enabled":true}]',
     '{"enabled":true,"days":14,"auto_convert":false,"require_payment":false,"trial_tier":"BASIC"}',
     TRUE, TRUE, FALSE, 11, '{"fixture": true}', '["basic","yearly"]',
     '00000000-0000-0000-0000-000000000000', NOW() - INTERVAL '365 days'),

    ('77777777-7777-7777-7777-777777777704', NULL,
     'PRO-M', 'Pro Monthly', 'Professional — AI เต็มรูปแบบ', 'PRO', 'MONTHLY', 'THB', 4900, 1,
     '[{"version":1,"price":{"amount":4900,"currency":"THB"},"effective_at":"2026-01-01T00:00:00Z","changed_by":"00000000-0000-0000-0000-000000000000","changed_at":"2026-01-01T00:00:00Z"}]',
     '{"max_devices":500,"max_sites":50,"max_users":100,"max_automations":500,"storage_gb":100,"retention_days":90,"api_calls_per_day":100000,"telemetry_points_per_day":100,"max_dashboards":50,"max_reports":500}',
     '[{"code":"DASHBOARD_ADVANCED","name":"Advanced Dashboard","enabled":true},{"code":"AI_FULL","name":"Full AI","enabled":true},{"code":"API_ACCESS","name":"API Access","enabled":true},{"code":"WS_REALTIME","name":"Realtime WS","enabled":true},{"code":"AUTOMATION_ADVANCED","name":"Advanced Automation","enabled":true}]',
     '{"enabled":true,"days":14,"auto_convert":false,"require_payment":true,"trial_tier":"PRO"}',
     TRUE, TRUE, FALSE, 20, '{"fixture": true}', '["pro","monthly"]',
     '00000000-0000-0000-0000-000000000000', NOW() - INTERVAL '365 days'),

    ('77777777-7777-7777-7777-777777777705', NULL,
     'ENT-M', 'Enterprise Monthly', 'Unlimited + SLA 99.95%', 'ENTERPRISE', 'MONTHLY', 'THB', 19900, 1,
     '[{"version":1,"price":{"amount":19900,"currency":"THB"},"effective_at":"2026-01-01T00:00:00Z","changed_by":"00000000-0000-0000-0000-000000000000","changed_at":"2026-01-01T00:00:00Z"}]',
     '{"max_devices":-1,"max_sites":-1,"max_users":-1,"max_automations":-1,"storage_gb":1000,"retention_days":365,"api_calls_per_day":-1,"telemetry_points_per_day":-1,"max_dashboards":-1,"max_reports":-1}',
     '[{"code":"DASHBOARD_ADVANCED","name":"Advanced Dashboard","enabled":true},{"code":"AI_FULL","name":"Full AI","enabled":true},{"code":"API_ACCESS","name":"API Access","enabled":true},{"code":"WS_REALTIME","name":"Realtime WS","enabled":true},{"code":"WHITELABEL","name":"White-label","enabled":true},{"code":"SSO","name":"Single Sign-On","enabled":true},{"code":"SLA_9995","name":"SLA 99.95%","enabled":true}]',
     '{"enabled":true,"days":30,"auto_convert":false,"require_payment":true,"trial_tier":"ENTERPRISE"}',
     TRUE, TRUE, FALSE, 30, '{"fixture": true}', '["enterprise","monthly"]',
     '00000000-0000-0000-0000-000000000000', NOW() - INTERVAL '365 days'),

    -- Deprecated (for testing)
    ('77777777-7777-7777-7777-777777777706', NULL,
     'LEGACY-BASIC', 'Legacy Basic (Deprecated)', 'ไม่แนะนำให้ใช้', 'BASIC', 'MONTHLY', 'THB', 500, 2,
     '[{"version":1,"price":{"amount":500,"currency":"THB"},"effective_at":"2025-01-01T00:00:00Z","expired_at":"2025-06-01T00:00:00Z","changed_by":"00000000-0000-0000-0000-000000000000","changed_at":"2025-01-01T00:00:00Z"},{"version":2,"price":{"amount":600,"currency":"THB"},"effective_at":"2025-06-01T00:00:00Z","changed_by":"00000000-0000-0000-0000-000000000000","changed_at":"2025-05-15T00:00:00Z"}]',
     '{"max_devices":20,"max_sites":2,"max_users":5,"max_automations":20,"storage_gb":5,"retention_days":14,"api_calls_per_day":5000,"telemetry_points_per_day":5,"max_dashboards":2,"max_reports":10}',
     '[]',
     '{"enabled":false,"days":0,"auto_convert":false,"require_payment":false,"trial_tier":"BASIC"}',
     FALSE, FALSE, TRUE, 99, '{"fixture": true, "deprecation_reason": "superseded by BASIC-M"}', '["legacy"]',
     '00000000-0000-0000-0000-000000000000', NOW() - INTERVAL '500 days')
ON CONFLICT (id) DO NOTHING;

COMMIT;
```

### B.7 `test/fixtures/sql/06_subscriptions.sql`

```sql
-- ============================================================
-- Subscriptions (7 active/sub-states + 1 trial)
-- ============================================================
BEGIN;

INSERT INTO packagecatalog_subscriptions
    (id, tenant_id, customer_id, package_id, status, cycle,
     started_at, expires_at, trial_ends_at, grace_ends_at,
     auto_renew, next_billing_at, price_at_signup, currency, current_usage,
     usage_reset_at, metadata, created_by, created_at, updated_at)
VALUES
    -- Demo customers
    ('88888888-8888-8888-8888-888888888801',
     '11111111-1111-1111-1111-111111111101', '33333333-3333-3333-3333-333333333301',
     '77777777-7777-7777-7777-777777777704', -- PRO-M
     'ACTIVE', 'MONTHLY',
     NOW() - INTERVAL '60 days', NOW() + INTERVAL '15 days',
     NULL, NULL, TRUE, NOW() + INTERVAL '15 days',
     4900, 'THB', '{"devices":23,"sites":1,"users":5,"automations":12,"storage_mb":2048,"api_calls_today":1200,"telemetry_today":5,"dashboards":3,"reports":15}',
     NOW() + INTERVAL '15 days',
     '{"fixture": true}', '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '60 days', NOW() - INTERVAL '1 day'),

    ('88888888-8888-8888-8888-888888888802',
     '11111111-1111-1111-1111-111111111101', '33333333-3333-3333-3333-333333333302',
     '77777777-7777-7777-7777-777777777702', -- BASIC-M
     'ACTIVE', 'MONTHLY',
     NOW() - INTERVAL '45 days', NOW() + INTERVAL '20 days',
     NULL, NULL, TRUE, NOW() + INTERVAL '20 days',
     990, 'THB', '{"devices":8,"sites":1,"users":2,"automations":5,"storage_mb":512,"api_calls_today":300,"telemetry_today":2,"dashboards":2,"reports":5}',
     NOW() + INTERVAL '20 days',
     '{"fixture": true}', '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '45 days', NOW() - INTERVAL '1 day'),

    ('88888888-8888-8888-8888-888888888803',
     '11111111-1111-1111-1111-111111111101', '33333333-3333-3333-3333-333333333304',
     '77777777-7777-7777-7777-777777777705', -- ENT-M
     'ACTIVE', 'MONTHLY',
     NOW() - INTERVAL '90 days', NOW() + INTERVAL '5 days',
     NULL, NULL, TRUE, NOW() + INTERVAL '5 days',
     19900, 'THB', '{"devices":156,"sites":1,"users":25,"automations":80,"storage_mb":45000,"api_calls_today":8500,"telemetry_today":45,"dashboards":12,"reports":120}',
     NOW() + INTERVAL '5 days',
     '{"fixture": true}', '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '90 days', NOW() - INTERVAL '1 day'),

    -- ACME customers (mixed states)
    ('88888888-8888-8888-8888-888888888811',
     '11111111-1111-1111-1111-111111111102', '33333333-3333-3333-3333-333333333311',
     '77777777-7777-7777-7777-777777777705', -- ENT-M
     'ACTIVE', 'MONTHLY',
     NOW() - INTERVAL '150 days', NOW() + INTERVAL '25 days',
     NULL, NULL, TRUE, NOW() + INTERVAL '25 days',
     19900, 'THB', '{"devices":450,"sites":1,"users":45,"automations":200,"storage_mb":180000,"api_calls_today":35000,"telemetry_today":80,"dashboards":25,"reports":300}',
     NOW() + INTERVAL '25 days',
     '{"fixture": true}', '22222222-2222-2222-2222-222222222211',
     NOW() - INTERVAL '150 days', NOW() - INTERVAL '2 days'),

    ('88888888-8888-8888-8888-888888888812',
     '11111111-1111-1111-1111-111111111102', '33333333-3333-3333-3333-333333333312',
     '77777777-7777-7777-7777-777777777704', -- PRO-M
     'ACTIVE', 'MONTHLY',
     NOW() - INTERVAL '140 days', NOW() + INTERVAL '8 days',
     NULL, NULL, TRUE, NOW() + INTERVAL '8 days',
     4900, 'THB', '{"devices":320,"sites":1,"users":30,"automations":120,"storage_mb":75000,"api_calls_today":40000,"telemetry_today":30,"dashboards":15,"reports":150}',
     NOW() + INTERVAL '8 days',
     '{"fixture": true}', '22222222-2222-2222-2222-222222222211',
     NOW() - INTERVAL '140 days', NOW() - INTERVAL '2 days'),

    -- PAST_DUE
    ('88888888-8888-8888-8888-888888888813',
     '11111111-1111-1111-1111-111111111102', '33333333-3333-3333-3333-333333333313',
     '77777777-7777-7777-7777-777777777702', -- BASIC-M
     'PAST_DUE', 'MONTHLY',
     NOW() - INTERVAL '100 days', NOW() - INTERVAL '3 days',
     NULL, NOW() + INTERVAL '4 days',  -- grace ยังเหลือ 4 วัน
     FALSE, NULL, 990, 'THB', '{"devices":5,"sites":1,"users":2}',
     NOW() - INTERVAL '3 days',
     '{"fixture": true, "past_due_reason": "payment_failed 3 times"}',
     '22222222-2222-2222-2222-222222222211',
     NOW() - INTERVAL '100 days', NOW() - INTERVAL '1 day'),

    -- Farm1 - TRIAL
    ('88888888-8888-8888-8888-888888888821',
     '11111111-1111-1111-1111-111111111103', '33333333-3333-3333-3333-333333333321',
     '77777777-7777-7777-7777-777777777704', -- PRO-M (trial)
     'TRIAL', 'MONTHLY',
     NOW() - INTERVAL '8 days', NOW() + INTERVAL '6 days',
     NOW() + INTERVAL '6 days',
     NULL, FALSE, NOW() + INTERVAL '6 days',
     4900, 'THB', '{"devices":12,"sites":1,"users":2,"automations":8}',
     NOW() + INTERVAL '6 days',
     '{"fixture": true, "trial_from": "customer onboarding"}',
     '22222222-2222-2222-2222-222222222221',
     NOW() - INTERVAL '8 days', NOW() - INTERVAL '1 day'),

    ('88888888-8888-8888-8888-888888888822',
     '11111111-1111-1111-1111-111111111103', '33333333-3333-3333-3333-333333333322',
     '77777777-7777-7777-7777-777777777702', -- BASIC-M
     'ACTIVE', 'MONTHLY',
     NOW() - INTERVAL '20 days', NOW() + INTERVAL '10 days',
     NULL, NULL, TRUE, NOW() + INTERVAL '10 days',
     990, 'THB', '{"devices":3,"sites":1,"users":1}',
     NOW() + INTERVAL '10 days',
     '{"fixture": true}', '22222222-2222-2222-2222-222222222221',
     NOW() - INTERVAL '20 days', NOW() - INTERVAL '1 day')
ON CONFLICT (id) DO NOTHING;

-- Subscription history (sample)
INSERT INTO packagecatalog_subscription_history
    (id, subscription_id, action, from_status, to_status, from_package_id, to_package_id, amount, currency, reason, actor_id, created_at)
VALUES
    ('99999999-9999-9999-9999-999999999901',
     '88888888-8888-8888-8888-888888888801',
     'CREATED', NULL, 'ACTIVE', NULL,
     '77777777-7777-7777-7777-777777777704',
     4900, 'THB', 'initial subscription',
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '60 days'),

    ('99999999-9999-9999-9999-999999999902',
     '88888888-8888-8888-8888-888888888801',
     'RENEWED', 'ACTIVE', 'ACTIVE', NULL, NULL,
     4900, 'THB', 'auto-renew',
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '30 days'),

    ('99999999-9999-9999-9999-999999999911',
     '88888888-8888-8888-8888-888888888811',
     'CREATED', NULL, 'ACTIVE', NULL,
     '77777777-7777-7777-7777-777777777705',
     19900, 'THB', 'enterprise contract',
     '22222222-2222-2222-2222-222222222211',
     NOW() - INTERVAL '150 days')
ON CONFLICT (id) DO NOTHING;

COMMIT;
```

### B.8 `test/fixtures/sql/07_products_warehouses.sql`

```sql
-- ============================================================
-- Products + Warehouses + Suppliers
-- ============================================================
BEGIN;

-- Warehouses
INSERT INTO erp_warehouses (id, tenant_id, code, name, type, address, is_active, is_default, metadata, created_at)
VALUES
    ('a0000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111101',
     'WH-DEMO', 'Demo Main Warehouse', 'MAIN',
     '{"line1":"99/1","province":"กรุงเทพ","country":"TH"}', TRUE, TRUE,
     '{"fixture": true}', NOW() - INTERVAL '90 days'),

    ('a0000000-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111101',
     'WH-DEMO-2', 'Demo Transit', 'TRANSIT',
     '{"line1":"99/1","province":"กรุงเทพ","country":"TH"}', TRUE, FALSE,
     '{"fixture": true}', NOW() - INTERVAL '90 days'),

    ('a0000000-0000-0000-0000-000000000011', '11111111-1111-1111-1111-111111111102',
     'WH-ACME', 'ACME Main', 'MAIN',
     '{"line1":"1","province":"กรุงเทพ","country":"TH"}', TRUE, TRUE,
     '{"fixture": true}', NOW() - INTERVAL '180 days'),

    ('a0000000-0000-0000-0000-000000000021', '11111111-1111-1111-1111-111111111103',
     'WH-FARM1', 'Farm1 Store', 'MAIN',
     '{"line1":"55","province":"เชียงใหม่","country":"TH"}', TRUE, TRUE,
     '{"fixture": true}', NOW() - INTERVAL '30 days')
ON CONFLICT (id) DO NOTHING;

-- Suppliers
INSERT INTO erp_suppliers (id, tenant_id, code, name, tax_id, email, phone, address, status, payment_terms, currency, is_active, metadata, created_at)
VALUES
    ('b0000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111101',
     'SUP-001', 'SensorTech Co., Ltd.', '0105550000101', 'sales@sensortech.co.th', '021111111',
     '{"line1":"1 นิคมฯ","province":"ปทุมธานี","country":"TH"}', 'ACTIVE', 'NET_30', 'THB', TRUE,
     '{"fixture": true}', NOW() - INTERVAL '90 days'),

    ('b0000000-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111101',
     'SUP-002', 'Gateway Electronics', '0105550000102', 'contact@gateway.co.th', '022222222',
     '{"line1":"2 ราชเทวี","province":"กรุงเทพ","country":"TH"}', 'ACTIVE', 'NET_45', 'THB', TRUE,
     '{"fixture": true}', NOW() - INTERVAL '80 days'),

    ('b0000000-0000-0000-0000-000000000011', '11111111-1111-1111-1111-111111111102',
     'SUP-ACME-001', 'Acme Components', '0105550000201', 'info@acmecomp.co.th', '023333333',
     '{"line1":"3","province":"กรุงเทพ","country":"TH"}', 'ACTIVE', 'NET_30', 'THB', TRUE,
     '{"fixture": true}', NOW() - INTERVAL '180 days'),

    ('b0000000-0000-0000-0000-000000000021', '11111111-1111-1111-1111-111111111103',
     'SUP-FARM1-001', 'เชียงใหม่ อะไหล่', '0505550000301', 'parts@cm.local', '053444444',
     '{"line1":"4","province":"เชียงใหม่","country":"TH"}', 'ACTIVE', 'NET_15', 'THB', TRUE,
     '{"fixture": true}', NOW() - INTERVAL '30 days')
ON CONFLICT (id) DO NOTHING;

-- Products (50 items — representative 20)
INSERT INTO erp_products (id, tenant_id, sku, name, type, category, uom, cost_price, sell_price, currency, tax_type, tax_rate, track_inventory, reorder_point, reorder_qty, lead_time_days, is_active, is_sellable, is_purchasable, metadata, tags, created_by, created_at)
VALUES
    -- Demo tenant sensors
    ('c0000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111101',
     'SENSOR-TEMP-001', 'Temperature Sensor SHT30', 'GOODS', 'SENSOR', 'PCS', 350, 650, 'THB', 'VAT', 0.07,
     TRUE, 50, 100, 7, TRUE, TRUE, TRUE,
     '{"fixture": true, "brand": "Sensirion", "spec": "±0.3°C accuracy"}', '["sensor","temperature"]',
     '22222222-2222-2222-2222-222222222201', NOW() - INTERVAL '90 days'),

    ('c0000000-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111101',
     'SENSOR-HUM-001', 'Humidity Sensor DHT22', 'GOODS', 'SENSOR', 'PCS', 250, 450, 'THB', 'VAT', 0.07,
     TRUE, 40, 80, 7, TRUE, TRUE, TRUE,
     '{"fixture": true, "brand": "Aosong"}', '["sensor","humidity"]',
     '22222222-2222-2222-2222-222222222201', NOW() - INTERVAL '90 days'),

    ('c0000000-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111101',
     'SENSOR-SOIL-001', 'Soil Moisture Sensor v2', 'GOODS', 'SENSOR', 'PCS', 450, 890, 'THB', 'VAT', 0.07,
     TRUE, 30, 60, 10, TRUE, TRUE, TRUE,
     '{"fixture": true}', '["sensor","soil"]',
     '22222222-2222-2222-2222-222222222201', NOW() - INTERVAL '90 days'),

    ('c0000000-0000-0000-0000-000000000004', '11111111-1111-1111-1111-111111111101',
     'SENSOR-PH-001', 'Soil pH Sensor', 'GOODS', 'SENSOR', 'PCS', 850, 1690, 'THB', 'VAT', 0.07,
     TRUE, 15, 30, 14, TRUE, TRUE, TRUE,
     '{"fixture": true}', '["sensor","ph"]',
     '22222222-2222-2222-2222-222222222201', NOW() - INTERVAL '85 days'),

    ('c0000000-0000-0000-0000-000000000010', '11111111-1111-1111-1111-111111111101',
     'GATEWAY-MQTT-001', 'IoT Gateway MQTT 4G', 'GOODS', 'GATEWAY', 'PCS', 3500, 6900, 'THB', 'VAT', 0.07,
     TRUE, 10, 20, 14, TRUE, TRUE, TRUE,
     '{"fixture": true, "spec": "LTE Cat-4, MQTT/TLS"}', '["gateway","mqtt","4g"]',
     '22222222-2222-2222-2222-222222222201', NOW() - INTERVAL '85 days'),

    ('c0000000-0000-0000-0000-000000000011', '11111111-1111-1111-1111-111111111101',
     'GATEWAY-LORA-001', 'LoRaWAN Gateway 8ch', 'GOODS', 'GATEWAY', 'PCS', 8500, 16900, 'THB', 'VAT', 0.07,
     TRUE, 5, 10, 21, TRUE, TRUE, TRUE,
     '{"fixture": true}', '["gateway","lorawan"]',
     '22222222-2222-2222-2222-222222222201', NOW() - INTERVAL '85 days'),

    ('c0000000-0000-0000-0000-000000000020', '11111111-1111-1111-1111-111111111101',
     'ACTUATOR-PUMP-001', 'Water Pump Relay 220V', 'GOODS', 'ACTUATOR', 'PCS', 850, 1490, 'THB', 'VAT', 0.07,
     TRUE, 20, 40, 7, TRUE, TRUE, TRUE,
     '{"fixture": true}', '["actuator","pump"]',
     '22222222-2222-2222-2222-222222222201', NOW() - INTERVAL '80 days'),

    ('c0000000-0000-0000-0000-000000000021', '11111111-1111-1111-1111-111111111101',
     'ACTUATOR-FAN-001', 'Exhaust Fan Controller', 'GOODS', 'ACTUATOR', 'PCS', 1200, 2200, 'THB', 'VAT', 0.07,
     TRUE, 15, 30, 10, TRUE, TRUE, TRUE,
     '{"fixture": true}', '["actuator","fan"]',
     '22222222-2222-2222-2222-222222222201', NOW() - INTERVAL '80 days'),

    -- Services
    ('c0000000-0000-0000-0000-000000000100', '11111111-1111-1111-1111-111111111101',
     'SVC-INSTALL-001', 'Installation Service (Standard)', 'SERVICE', 'SERVICE', 'JOB', 0, 2500, 'THB', 'VAT', 0.07,
     FALSE, 0, 0, 0, TRUE, TRUE, FALSE,
     '{"fixture": true}', '["service","installation"]',
     '22222222-2222-2222-2222-222222222201', NOW() - INTERVAL '90 days'),

    ('c0000000-0000-0000-0000-000000000101', '11111111-1111-1111-1111-111111111101',
     'SVC-MAINT-001', 'Maintenance Service (Annual)', 'SERVICE', 'SERVICE', 'JOB', 0, 12000, 'THB', 'VAT', 0.07,
     FALSE, 0, 0, 0, TRUE, TRUE, FALSE,
     '{"fixture": true}', '["service","maintenance"]',
     '22222222-2222-2222-2222-222222222201', NOW() - INTERVAL '90 days'),

    -- ACME tenant — sample products
    ('c0000000-0000-0000-0000-000000000201', '11111111-1111-1111-1111-111111111102',
     'SENSOR-TEMP-001', 'Temp Sensor (ACME)', 'GOODS', 'SENSOR', 'PCS', 320, 620, 'THB', 'VAT', 0.07,
     TRUE, 100, 200, 7, TRUE, TRUE, TRUE,
     '{"fixture": true}', '["sensor"]',
     '22222222-2222-2222-2222-222222222211', NOW() - INTERVAL '180 days'),

    ('c0000000-0000-0000-0000-000000000202', '11111111-1111-1111-1111-111111111102',
     'GATEWAY-IND-001', 'Industrial Gateway 5G', 'GOODS', 'GATEWAY', 'PCS', 12000, 24000, 'THB', 'VAT', 0.07,
     TRUE, 5, 15, 30, TRUE, TRUE, TRUE,
     '{"fixture": true}', '["gateway","5g"]',
     '22222222-2222-2222-2222-222222222211', NOW() - INTERVAL '180 days'),

    -- Farm1
    ('c0000000-0000-0000-0000-000000000301', '11111111-1111-1111-1111-111111111103',
     'SENSOR-TEMP-001', 'Temp Sensor Farm', 'GOODS', 'SENSOR', 'PCS', 350, 700, 'THB', 'VAT', 0.07,
     TRUE, 10, 20, 14, TRUE, TRUE, TRUE,
     '{"fixture": true}', '["sensor"]',
     '22222222-2222-2222-2222-222222222221', NOW() - INTERVAL '25 days'),

    ('c0000000-0000-0000-0000-000000000302', '11111111-1111-1111-1111-111111111103',
     'SENSOR-SOIL-001', 'Soil Moisture', 'GOODS', 'SENSOR', 'PCS', 450, 900, 'THB', 'VAT', 0.07,
     TRUE, 10, 20, 14, TRUE, TRUE, TRUE,
     '{"fixture": true}', '["sensor","soil"]',
     '22222222-2222-2222-2222-222222222221', NOW() - INTERVAL '25 days')
ON CONFLICT (id) DO NOTHING;

COMMIT;
```

### B.9 `test/fixtures/sql/08_inventory.sql`

```sql
-- ============================================================
-- Inventory items (with mix of levels)
-- ============================================================
BEGIN;

INSERT INTO erp_inventory_items
    (id, tenant_id, product_id, warehouse_id, qty_on_hand, qty_reserved, avg_cost, last_cost, bin_location, created_at, updated_at)
VALUES
    -- Demo warehouse — healthy stock
    ('d0000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111101',
     'c0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001',
     125, 10, 350, 350, 'A-01-01', NOW() - INTERVAL '90 days', NOW() - INTERVAL '1 day'),

    ('d0000000-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111101',
     'c0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001',
     78, 5, 250, 250, 'A-01-02', NOW() - INTERVAL '90 days', NOW() - INTERVAL '1 day'),

    ('d0000000-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111101',
     'c0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001',
     45, 0, 450, 450, 'A-02-01', NOW() - INTERVAL '90 days', NOW() - INTERVAL '1 day'),

    -- LOW STOCK (below reorder_point=30)
    ('d0000000-0000-0000-0000-000000000004', '11111111-1111-1111-1111-111111111101',
     'c0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000001',
     12, 2, 850, 850, 'A-02-02', NOW() - INTERVAL '85 days', NOW() - INTERVAL '2 days'),

    -- LOW STOCK (below reorder_point=10)
    ('d0000000-0000-0000-0000-000000000010', '11111111-1111-1111-1111-111111111101',
     'c0000000-0000-0000-0000-000000000010', 'a0000000-0000-0000-0000-000000000001',
     6, 2, 3500, 3500, 'B-01-01', NOW() - INTERVAL '85 days', NOW() - INTERVAL '1 day'),

    -- OUT OF STOCK
    ('d0000000-0000-0000-0000-000000000011', '11111111-1111-1111-1111-111111111101',
     'c0000000-0000-0000-0000-000000000011', 'a0000000-0000-0000-0000-000000000001',
     0, 0, 8500, 8500, 'B-01-02', NOW() - INTERVAL '85 days', NOW() - INTERVAL '5 days'),

    ('d0000000-0000-0000-0000-000000000020', '11111111-1111-1111-1111-111111111101',
     'c0000000-0000-0000-0000-000000000020', 'a0000000-0000-0000-0000-000000000001',
     42, 5, 850, 850, 'C-01-01', NOW() - INTERVAL '80 days', NOW() - INTERVAL '1 day'),

    ('d0000000-0000-0000-0000-000000000021', '11111111-1111-1111-1111-111111111101',
     'c0000000-0000-0000-0000-000000000021', 'a0000000-0000-0000-0000-000000000001',
     28, 0, 1200, 1200, 'C-01-02', NOW() - INTERVAL '80 days', NOW() - INTERVAL '1 day'),

    -- Transit warehouse
    ('d0000000-0000-0000-0000-000000000101', '11111111-1111-1111-1111-111111111101',
     'c0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000002',
     20, 0, 350, 350, NULL, NOW() - INTERVAL '30 days', NOW() - INTERVAL '1 day'),

    -- ACME warehouse
    ('d0000000-0000-0000-0000-000000000201', '11111111-1111-1111-1111-111111111102',
     'c0000000-0000-0000-0000-000000000201', 'a0000000-0000-0000-0000-000000000011',
     450, 50, 320, 320, 'M-01', NOW() - INTERVAL '180 days', NOW() - INTERVAL '1 day'),

    ('d0000000-0000-0000-0000-000000000202', '11111111-1111-1111-1111-111111111102',
     'c0000000-0000-0000-0000-000000000202', 'a0000000-0000-0000-0000-000000000011',
     8, 2, 12000, 12000, 'M-02', NOW() - INTERVAL '180 days', NOW() - INTERVAL '1 day'),

    -- Farm1
    ('d0000000-0000-0000-0000-000000000301', '11111111-1111-1111-1111-111111111103',
     'c0000000-0000-0000-0000-000000000301', 'a0000000-0000-0000-0000-000000000021',
     15, 3, 350, 350, 'F-01', NOW() - INTERVAL '25 days', NOW() - INTERVAL '1 day'),

    ('d0000000-0000-0000-0000-000000000302', '11111111-1111-1111-1111-111111111103',
     'c0000000-0000-0000-0000-000000000302', 'a0000000-0000-0000-0000-000000000021',
     8, 0, 450, 450, 'F-02', NOW() - INTERVAL '25 days', NOW() - INTERVAL '1 day')
ON CONFLICT (id) DO NOTHING;

-- Stock movements (sample ledger)
INSERT INTO erp_stock_movements
    (tenant_id, product_id, warehouse_id, type, direction, quantity, uom, unit_cost, total_cost,
     qty_on_hand_after, qty_reserved_after, ref_type, ref_number, occurred_at, created_by)
VALUES
    ('11111111-1111-1111-1111-111111111101', 'c0000000-0000-0000-0000-000000000001',
     'a0000000-0000-0000-0000-000000000001', 'RECEIPT', 1, 100, 'PCS', 350, 35000,
     100, 0, 'ORDER', 'PO-2026-000001', NOW() - INTERVAL '85 days',
     '22222222-2222-2222-2222-222222222201'),

    ('11111111-1111-1111-1111-111111111101', 'c0000000-0000-0000-0000-000000000001',
     'a0000000-0000-0000-0000-000000000001', 'ISSUE', -1, 25, 'PCS', 350, 8750,
     75, 0, 'ORDER', 'SO-2026-000001', NOW() - INTERVAL '60 days',
     '22222222-2222-2222-2222-222222222201'),

    ('11111111-1111-1111-1111-111111111101', 'c0000000-0000-0000-0000-000000000001',
     'a0000000-0000-0000-0000-000000000001', 'RECEIPT', 1, 50, 'PCS', 350, 17500,
     125, 0, 'ORDER', 'PO-2026-000010', NOW() - INTERVAL '30 days',
     '22222222-2222-2222-2222-222222222201'),

    ('11111111-1111-1111-1111-111111111101', 'c0000000-0000-0000-0000-000000000001',
     'a0000000-0000-0000-0000-000000000001', 'RESERVE', 0, 10, 'PCS', 0, 0,
     125, 10, 'ORDER', 'SO-2026-000100', NOW() - INTERVAL '2 days',
     '22222222-2222-2222-2222-222222222201');

COMMIT;
```

### B.10 `test/fixtures/sql/09_devices.sql`

```sql
-- ============================================================
-- Devices + Shadow + Rules + Commands
-- ============================================================
BEGIN;

-- Device models
INSERT INTO device_models (id, vendor, model_no, name, device_type, protocol, capabilities, default_config, firmware_channel, created_at)
VALUES
    ('e0000000-0000-0000-0000-000000000001', 'Sensirion', 'SHT30-DIS',
     'SHT30 Temp/Humidity', 'SENSOR', 'MQTT',
     '[{"type":"SENSOR","metric":"temperature","unit":"°C","min_value":-40,"max_value":125},{"type":"SENSOR","metric":"humidity","unit":"%","min_value":0,"max_value":100}]',
     '{}', 'stable', NOW() - INTERVAL '365 days'),

    ('e0000000-0000-0000-0000-000000000010', 'Espressif', 'ESP32-GW',
     'ESP32 Gateway', 'GATEWAY', 'MQTT',
     '[{"type":"GATEWAY"}]',
     '{}', 'stable', NOW() - INTERVAL '365 days'),

    ('e0000000-0000-0000-0000-000000000020', 'Generic', 'RELAY-220V',
     '220V Relay Actuator', 'ACTUATOR', 'MQTT',
     '[{"type":"ACTUATOR","command":"on"},{"type":"ACTUATOR","command":"off"},{"type":"ACTUATOR","command":"toggle"}]',
     '{}', 'stable', NOW() - INTERVAL '365 days')
ON CONFLICT (id) DO NOTHING;

-- Devices (30 items — representative 15)
INSERT INTO device_devices
    (id, tenant_id, customer_id, site_id, zone_id, model_id,
     serial_no, name, type, protocol, status, mqtt_client_id, device_token_hash,
     firmware_version, last_seen_at, provisioned_at, capabilities, tags, metadata,
     created_by, created_at, updated_at)
VALUES
    -- Demo — ONLINE devices
    ('f0000000-0000-0000-0000-000000000001',
     '11111111-1111-1111-1111-111111111101', '33333333-3333-3333-3333-333333333301',
     '55555555-5555-5555-5555-555555555501', '66666666-6666-6666-6666-666666666601',
     'e0000000-0000-0000-0000-000000000001',
     'SN-DEMO-001', 'โซน A - Temp 1', 'SENSOR', 'MQTT', 'ONLINE',
     'dev-demo-sn-001', '$2a$10$abcdefghijklmnopqrstuvwxyz', '1.0.0',
     NOW() - INTERVAL '2 minutes', NOW() - INTERVAL '60 days',
     '[{"type":"SENSOR","metric":"temperature","unit":"°C"},{"type":"SENSOR","metric":"humidity","unit":"%"}]',
     '["demo","zone-a"]',
     '{"fixture": true}', '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '60 days', NOW() - INTERVAL '2 minutes'),

    ('f0000000-0000-0000-0000-000000000002',
     '11111111-1111-1111-1111-111111111101', '33333333-3333-3333-3333-333333333301',
     '55555555-5555-5555-5555-555555555501', '66666666-6666-6666-6666-666666666601',
     'e0000000-0000-0000-0000-000000000001',
     'SN-DEMO-002', 'โซน A - Temp 2', 'SENSOR', 'MQTT', 'ONLINE',
     'dev-demo-sn-002', '$2a$10$abcdefghijklmnopqrstuvwxyz', '1.0.0',
     NOW() - INTERVAL '3 minutes', NOW() - INTERVAL '60 days',
     '[{"type":"SENSOR","metric":"temperature","unit":"°C"}]',
     '["demo","zone-a"]',
     '{"fixture": true}', '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '60 days', NOW() - INTERVAL '3 minutes'),

    ('f0000000-0000-0000-0000-000000000003',
     '11111111-1111-1111-1111-111111111101', '33333333-3333-3333-3333-333333333301',
     '55555555-5555-5555-5555-555555555501', '66666666-6666-6666-6666-666666666602',
     'e0000000-0000-0000-0000-000000000001',
     'SN-DEMO-003', 'โซน B - Temp', 'SENSOR', 'MQTT', 'ONLINE',
     'dev-demo-sn-003', '$2a$10$abcdefghijklmnopqrstuvwxyz', '1.0.0',
     NOW() - INTERVAL '5 minutes', NOW() - INTERVAL '60 days',
     '[{"type":"SENSOR","metric":"temperature","unit":"°C"}]',
     '["demo","zone-b"]',
     '{"fixture": true}', '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '60 days', NOW() - INTERVAL '5 minutes'),

    -- OFFLINE
    ('f0000000-0000-0000-0000-000000000010',
     '11111111-1111-1111-1111-111111111101', '33333333-3333-3333-3333-333333333301',
     '55555555-5555-5555-5555-555555555501', '66666666-6666-6666-6666-666666666602',
     'e0000000-0000-0000-0000-000000000001',
     'SN-DEMO-010', 'โซน B - Offline', 'SENSOR', 'MQTT', 'OFFLINE',
     'dev-demo-sn-010', '$2a$10$abcdefghijklmnopqrstuvwxyz', '1.0.0',
     NOW() - INTERVAL '2 hours', NOW() - INTERVAL '55 days',
     '[{"type":"SENSOR","metric":"temperature","unit":"°C"}]',
     '["demo"]',
     '{"fixture": true, "offline_reason": "heartbeat_timeout"}',
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '55 days', NOW() - INTERVAL '2 hours'),

    -- FAULT
    ('f0000000-0000-0000-0000-000000000020',
     '11111111-1111-1111-1111-111111111101', '33333333-3333-3333-3333-333333333302',
     '55555555-5555-5555-5555-555555555502', '66666666-6666-6666-6666-666666666603',
     'e0000000-0000-0000-0000-000000000001',
     'SN-DEMO-FAULT-001', 'โรงเรือน 1 - Fault', 'SENSOR', 'MQTT', 'FAULT',
     'dev-demo-sn-fault-001', '$2a$10$abcdefghijklmnopqrstuvwxyz', '1.0.0',
     NOW() - INTERVAL '15 minutes', NOW() - INTERVAL '45 days',
     '[{"type":"SENSOR","metric":"temperature","unit":"°C"}]',
     '["demo","faulty"]',
     '{"fixture": true, "fault_code": "E001", "fault_message": "sensor_read_failure"}',
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '45 days', NOW() - INTERVAL '15 minutes'),

    -- Gateway
    ('f0000000-0000-0000-0000-000000000030',
     '11111111-1111-1111-1111-111111111101', '33333333-3333-3333-3333-333333333301',
     '55555555-5555-5555-5555-555555555501', NULL,
     'e0000000-0000-0000-0000-000000000010',
     'SN-DEMO-GW-001', 'Gateway หลัก', 'GATEWAY', 'MQTT', 'ONLINE',
     'dev-demo-gw-001', '$2a$10$abcdefghijklmnopqrstuvwxyz', '2.1.0',
     NOW() - INTERVAL '1 minute', NOW() - INTERVAL '60 days',
     '[{"type":"GATEWAY"}]',
     '["demo","gateway"]',
     '{"fixture": true}', '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '60 days', NOW() - INTERVAL '1 minute'),

    -- Actuator
    ('f0000000-0000-0000-0000-000000000040',
     '11111111-1111-1111-1111-111111111101', '33333333-3333-3333-3333-333333333302',
     '55555555-5555-5555-5555-555555555502', '66666666-6666-6666-6666-666666666603',
     'e0000000-0000-0000-0000-000000000020',
     'SN-DEMO-ACT-001', 'ปั๊มน้ำ 1', 'ACTUATOR', 'MQTT', 'ONLINE',
     'dev-demo-act-001', '$2a$10$abcdefghijklmnopqrstuvwxyz', '1.0.0',
     NOW() - INTERVAL '4 minutes', NOW() - INTERVAL '45 days',
     '[{"type":"ACTUATOR","command":"on"},{"type":"ACTUATOR","command":"off"}]',
     '["demo","actuator"]',
     '{"fixture": true}', '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '45 days', NOW() - INTERVAL '4 minutes'),

    -- ACME devices
    ('f0000000-0000-0000-0000-000000000101',
     '11111111-1111-1111-1111-111111111102', '33333333-3333-3333-3333-333333333311',
     '55555555-5555-5555-5555-555555555511', '66666666-6666-6666-6666-666666666611',
     'e0000000-0000-0000-0000-000000000001',
     'SN-ACME-001', 'ไร่อ้อย Sensor 1', 'SENSOR', 'MQTT', 'ONLINE',
     'dev-acme-001', '$2a$10$abcdefghijklmnopqrstuvwxyz', '1.0.0',
     NOW() - INTERVAL '1 minute', NOW() - INTERVAL '150 days',
     '[{"type":"SENSOR","metric":"temperature","unit":"°C"}]',
     '["acme"]', '{"fixture": true}',
     '22222222-2222-2222-2222-222222222211',
     NOW() - INTERVAL '150 days', NOW() - INTERVAL '1 minute'),

    ('f0000000-0000-0000-0000-000000000102',
     '11111111-1111-1111-1111-111111111102', '33333333-3333-3333-3333-333333333311',
     '55555555-5555-5555-5555-555555555511', '66666666-6666-6666-6666-666666666611',
     'e0000000-0000-0000-0000-000000000001',
     'SN-ACME-002', 'ไร่อ้อย Sensor 2', 'SENSOR', 'MQTT', 'ONLINE',
     'dev-acme-002', '$2a$10$abcdefghijklmnopqrstuvwxyz', '1.0.0',
     NOW() - INTERVAL '1 minute', NOW() - INTERVAL '150 days',
     '[{"type":"SENSOR","metric":"humidity","unit":"%"}]',
     '["acme"]', '{"fixture": true}',
     '22222222-2222-2222-2222-222222222211',
     NOW() - INTERVAL '150 days', NOW() - INTERVAL '1 minute'),

    -- Farm1
    ('f0000000-0000-0000-0000-000000000201',
     '11111111-1111-1111-1111-111111111103', '33333333-3333-3333-3333-333333333321',
     '55555555-5555-5555-5555-555555555521', '66666666-6666-6666-6666-666666666621',
     'e0000000-0000-0000-0000-000000000001',
     'SN-FARM1-001', 'ไฮโดร ชั้น 1 - Temp', 'SENSOR', 'MQTT', 'ONLINE',
     'dev-farm1-001', '$2a$10$abcdefghijklmnopqrstuvwxyz', '1.0.0',
     NOW() - INTERVAL '30 seconds', NOW() - INTERVAL '25 days',
     '[{"type":"SENSOR","metric":"temperature","unit":"°C"}]',
     '["farm1"]', '{"fixture": true}',
     '22222222-2222-2222-2222-222222222221',
     NOW() - INTERVAL '25 days', NOW() - INTERVAL '30 seconds'),

    ('f0000000-0000-0000-0000-000000000202',
     '11111111-1111-1111-1111-111111111103', '33333333-3333-3333-3333-333333333321',
     '55555555-5555-5555-5555-555555555521', '66666666-6666-6666-6666-666666666621',
     'e0000000-0000-0000-0000-000000000001',
     'SN-FARM1-002', 'ไฮโดร ชั้น 1 - Humidity', 'SENSOR', 'MQTT', 'ONLINE',
     'dev-farm1-002', '$2a$10$abcdefghijklmnopqrstuvwxyz', '1.0.0',
     NOW() - INTERVAL '45 seconds', NOW() - INTERVAL '25 days',
     '[{"type":"SENSOR","metric":"humidity","unit":"%"}]',
     '["farm1"]', '{"fixture": true}',
     '22222222-2222-2222-2222-222222222221',
     NOW() - INTERVAL '25 days', NOW() - INTERVAL '45 seconds')
ON CONFLICT (id) DO NOTHING;

-- Shadows (initial state for demo devices)
INSERT INTO device_shadows (device_id, tenant_id, state, version, created_at, updated_at)
VALUES
    ('f0000000-0000-0000-0000-000000000040', '11111111-1111-1111-1111-111111111101',
     '{"desired":{"pump":"off","schedule":"08:00-18:00"},"reported":{"pump":"off"},"delta":{},"metadata":{"desired":{},"reported":{}},"version":3}',
     3, NOW() - INTERVAL '45 days', NOW() - INTERVAL '1 day'),

    ('f0000000-0000-0000-0000-000000000030', '11111111-1111-1111-1111-111111111101',
     '{"desired":{"reporting_interval":60,"mode":"online"},"reported":{"reporting_interval":60,"mode":"online"},"delta":{},"metadata":{"desired":{},"reported":{}},"version":2}',
     2, NOW() - INTERVAL '60 days', NOW() - INTERVAL '1 day')
ON CONFLICT (device_id) DO NOTHING;

-- Alert rules
INSERT INTO device_alert_rules
    (id, tenant_id, name, device_id, group_id, site_id, metric, operator, threshold, duration_sec,
     severity, actions, is_active, cooldown_sec, created_by, created_at, updated_at)
VALUES
    ('aa000000-0000-0000-0000-000000000001',
     '11111111-1111-1111-1111-111111111101', 'อุณหภูมิสูงเกิน 40°C',
     'f0000000-0000-0000-0000-000000000001', NULL, NULL,
     'temperature', '>', 40, 0, 'CRITICAL',
     '[{"type":"NOTIFY","channel":"email","target":"admin@demo.local"},{"type":"BROADCAST_WS"}]',
     TRUE, 300, '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '55 days', NOW() - INTERVAL '55 days'),

    ('aa000000-0000-0000-0000-000000000002',
     '11111111-1111-1111-1111-111111111101', 'ความชื้นต่ำเกิน 30%',
     NULL, '66666666-6666-6666-6666-666666666601', NULL,
     'humidity', '<', 30, 60, 'WARN',
     '[{"type":"NOTIFY","channel":"email","target":"admin@demo.local"}]',
     TRUE, 600, '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '50 days', NOW() - INTERVAL '50 days'),

    ('aa000000-0000-0000-0000-000000000011',
     '11111111-1111-1111-1111-111111111102', 'ACME Temp Critical',
     NULL, NULL, '55555555-5555-5555-5555-555555555511',
     'temperature', '>', 45, 0, 'CRITICAL',
     '[{"type":"NOTIFY","channel":"email","target":"ops@acme.co.th"}]',
     TRUE, 120, '22222222-2222-2222-2222-222222222211',
     NOW() - INTERVAL '140 days', NOW() - INTERVAL '140 days')
ON CONFLICT (id) DO NOTHING;

-- Alert events (sample)
INSERT INTO device_alert_events
    (tenant_id, rule_id, device_id, metric, value, threshold, operator, severity, message,
     acknowledged, triggered_at)
VALUES
    ('11111111-1111-1111-1111-111111111101', 'aa000000-0000-0000-0000-000000000001',
     'f0000000-0000-0000-0000-000000000001', 'temperature', 42.5, 40.0, '>', 'CRITICAL',
     '[CRITICAL] temperature = 42.50 °C', TRUE, NOW() - INTERVAL '2 days'),

    ('11111111-1111-1111-1111-111111111101', 'aa000000-0000-0000-0000-000000000001',
     'f0000000-0000-0000-0000-000000000001', 'temperature', 41.2, 40.0, '>', 'CRITICAL',
     '[CRITICAL] temperature = 41.20 °C', FALSE, NOW() - INTERVAL '30 minutes'),

    ('11111111-1111-1111-1111-111111111102', 'aa000000-0000-0000-0000-000000000011',
     'f0000000-0000-0000-0000-000000000101', 'temperature', 47.0, 45.0, '>', 'CRITICAL',
     '[CRITICAL] temperature = 47.00 °C', FALSE, NOW() - INTERVAL '15 minutes');

-- Automations
INSERT INTO device_automations
    (id, tenant_id, name, description, trigger, conditions, actions, is_active, priority,
     run_count, last_run_at, created_by, created_at, updated_at)
VALUES
    ('bb000000-0000-0000-0000-000000000001',
     '11111111-1111-1111-1111-111111111101',
     'เปิดปั๊มเมื่อดินแห้ง', 'ถ้าความชื้นดิน < 30% → เปิดปั๊ม 5 นาที',
     '{"type":"TELEMETRY","metric":"soil_moisture"}',
     '[{"field":"soil_moisture","operator":"<","value":30}]',
     '[{"type":"SEND_COMMAND","device_id":"f0000000-0000-0000-0000-000000000040","command":"on","payload":{"duration_sec":300}}]',
     TRUE, 5, 42, NOW() - INTERVAL '3 hours',
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '40 days', NOW() - INTERVAL '3 hours')
ON CONFLICT (id) DO NOTHING;

-- Commands (sample)
INSERT INTO device_commands
    (id, tenant_id, device_id, command, payload, status, priority, issued_by, issued_at, sent_at, acked_at, expires_at, result_payload)
VALUES
    ('cc000000-0000-0000-0000-000000000001',
     '11111111-1111-1111-1111-111111111101', 'f0000000-0000-0000-0000-000000000040',
     'on', '{"duration_sec": 300}', 'ACKED', 3,
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour', NOW() - INTERVAL '59 minutes',
     NOW() - INTERVAL '59 minutes', '{"status": "pump_started"}'),

    ('cc000000-0000-0000-0000-000000000002',
     '11111111-1111-1111-1111-111111111101', 'f0000000-0000-0000-0000-000000000040',
     'off', '{}', 'TIMEOUT', 3,
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '29 minutes',
     NOW() - INTERVAL '29 minutes', NULL)
ON CONFLICT (id) DO NOTHING;

COMMIT;
```

### B.11 `test/fixtures/sql/10_orders.sql`

```sql
-- ============================================================
-- Orders (SO + PO) — 8 orders mixed statuses
-- ============================================================
BEGIN;

INSERT INTO erp_orders
    (id, tenant_id, order_no, type, customer_id, supplier_id, status,
     order_date, confirmed_at, shipped_at, received_at,
     ship_to_address, bill_to_address,
     subtotal, tax_amount, total_amount, currency,
     subtotal, tax_amount, total_amount, currency,
     created_by, approved_by, created_at, updated_at)
VALUES
    -- SO #1 — CLOSED (fully paid)
    ('ee000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111101',
     'SO-2026-000001', 'SALES', '33333333-3333-3333-3333-333333333301', NULL, 'CLOSED',
     NOW() - INTERVAL '60 days', NOW() - INTERVAL '60 days', NOW() - INTERVAL '58 days', NULL,
     '{"province":"ปทุมธานี","country":"TH"}', '{"province":"ปทุมธานี","country":"TH"}',
     65000, 4550, 69550, 'THB',
     '22222222-2222-2222-2222-222222222201', '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '60 days', NOW() - INTERVAL '30 days'),

    -- SO #2 — INVOICED (waiting payment)
    ('ee000000-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111101',
     'SO-2026-000002', 'SALES', '33333333-3333-3333-3333-333333333302', NULL, 'INVOICED',
     NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days', NOW() - INTERVAL '18 days', NULL,
     '{"province":"เชียงใหม่","country":"TH"}', '{"province":"เชียงใหม่","country":"TH"}',
     8900, 623, 9523, 'THB',
     '22222222-2222-2222-2222-222222222201', '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '20 days', NOW() - INTERVAL '5 days'),

    -- SO #3 — CONFIRMED (not shipped)
    ('ee000000-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111101',
     'SO-2026-000003', 'SALES', '33333333-3333-3333-3333-333333333304', NULL, 'CONFIRMED',
     NOW() - INTERVAL '3 days', NOW() - INTERVAL '2 days', NULL, NULL,
     '{"province":"นครราชสีมา","country":"TH"}', '{"province":"นครราชสีมา","country":"TH"}',
     45000, 3150, 48150, 'THB',
     '22222222-2222-2222-2222-222222222201', '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '3 days', NOW() - INTERVAL '2 days'),

    -- SO #4 — DRAFT
    ('ee000000-0000-0000-0000-000000000004', '11111111-1111-1111-1111-111111111101',
     'SO-2026-000004', 'SALES', '33333333-3333-3333-3333-333333333301', NULL, 'DRAFT',
     NOW() - INTERVAL '6 hours', NULL, NULL, NULL,
     '{"province":"ปทุมธานี","country":"TH"}', '{"province":"ปทุมธานี","country":"TH"}',
     2500, 175, 2675, 'THB',
     '22222222-2222-2222-2222-222222222201', NULL,
     NOW() - INTERVAL '6 hours', NOW() - INTERVAL '6 hours'),

    -- SO #5 — CANCELLED
    ('ee000000-0000-0000-0000-000000000005', '11111111-1111-1111-1111-111111111101',
     'SO-2026-000005', 'SALES', '33333333-3333-3333-3333-333333333302', NULL, 'CANCELLED',
     NOW() - INTERVAL '10 days', NULL, NULL, NULL,
     '{"province":"เชียงใหม่","country":"TH"}', '{"province":"เชียงใหม่","country":"TH"}',
     1200, 84, 1284, 'THB',
     '22222222-2222-2222-2222-222222222201', NULL,
     NOW() - INTERVAL '10 days', NOW() - INTERVAL '9 days'),

    -- PO #1 — RECEIVED
    ('ee000000-0000-0000-0000-000000000101', '11111111-1111-1111-1111-111111111101',
     'PO-2026-000001', 'PURCHASE', NULL, 'b0000000-0000-0000-0000-000000000001', 'RECEIVED',
     NOW() - INTERVAL '85 days', NOW() - INTERVAL '85 days', NULL, NOW() - INTERVAL '75 days',
     '{"province":"กรุงเทพ","country":"TH"}', '{"province":"กรุงเทพ","country":"TH"}',
     35000, 2450, 37450, 'THB',
     '22222222-2222-2222-2222-222222222201', '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '85 days', NOW() - INTERVAL '75 days'),

    -- PO #2 — CONFIRMED (awaiting delivery)
    ('ee000000-0000-0000-0000-000000000102', '11111111-1111-1111-1111-111111111101',
     'PO-2026-000002', 'PURCHASE', NULL, 'b0000000-0000-0000-0000-000000000002', 'CONFIRMED',
     NOW() - INTERVAL '5 days', NOW() - INTERVAL '4 days', NULL, NULL,
     '{"province":"กรุงเทพ","country":"TH"}', '{"province":"กรุงเทพ","country":"TH"}',
     68000, 4760, 72760, 'THB',
     '22222222-2222-2222-2222-222222222201', '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '5 days', NOW() - INTERVAL '4 days'),

    -- ACME SO
    ('ee000000-0000-0000-0000-000000000201', '11111111-1111-1111-1111-111111111102',
     'SO-2026-000101', 'SALES', '33333333-3333-3333-3333-333333333311', NULL, 'INVOICED',
     NOW() - INTERVAL '15 days', NOW() - INTERVAL '15 days', NOW() - INTERVAL '13 days', NULL,
     '{"province":"ศรีสะเกษ","country":"TH"}', '{"province":"ศรีสะเกษ","country":"TH"}',
     124000, 8680, 132680, 'THB',
     '22222222-2222-2222-2222-222222222211', '22222222-2222-2222-2222-222222222211',
     NOW() - INTERVAL '15 days', NOW() - INTERVAL '5 days')
ON CONFLICT (id) DO NOTHING;

-- Order lines (representative for first 3)
INSERT INTO erp_order_lines
    (id, order_id, line_no, product_id, sku, name, quantity, uom, unit_price, discount,
     tax_type, tax_rate, subtotal, tax_amount, line_total,
     qty_shipped, qty_received, qty_invoiced, warehouse_id, created_at, updated_at)
VALUES
    -- SO #1 lines
    ('ff000000-0000-0000-0000-000000000001', 'ee000000-0000-0000-0000-000000000001', 1,
     'c0000000-0000-0000-0000-000000000001', 'SENSOR-TEMP-001', 'Temperature Sensor SHT30',
     50, 'PCS', 650, 0, 'VAT', 0.07, 32500, 2275, 34775,
     50, 0, 50, 'a0000000-0000-0000-0000-000000000001',
     NOW() - INTERVAL '60 days', NOW() - INTERVAL '30 days'),

    ('ff000000-0000-0000-0000-000000000002', 'ee000000-0000-0000-0000-000000000001', 2,
     'c0000000-0000-0000-0000-000000000010', 'GATEWAY-MQTT-001', 'IoT Gateway MQTT 4G',
     5, 'PCS', 6500, 0, 'VAT', 0.07, 32500, 2275, 34775,
     5, 0, 5, 'a0000000-0000-0000-0000-000000000001',
     NOW() - INTERVAL '60 days', NOW() - INTERVAL '30 days'),

    -- SO #2 lines
    ('ff000000-0000-0000-0000-000000000010', 'ee000000-0000-0000-0000-000000000002', 1,
     'c0000000-0000-0000-0000-000000000003', 'SENSOR-SOIL-001', 'Soil Moisture Sensor v2',
     10, 'PCS', 890, 0, 'VAT', 0.07, 8900, 623, 9523,
     10, 0, 10, 'a0000000-0000-0000-0000-000000000001',
     NOW() - INTERVAL '20 days', NOW() - INTERVAL '5 days'),

    -- SO #3 lines
    ('ff000000-0000-0000-0000-000000000020', 'ee000000-0000-0000-0000-000000000003', 1,
     'c0000000-0000-0000-0000-000000000001', 'SENSOR-TEMP-001', 'Temperature Sensor SHT30',
     30, 'PCS', 650, 0, 'VAT', 0.07, 19500, 1365, 20865,
     0, 0, 0, 'a0000000-0000-0000-0000-000000000001',
     NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days'),

    ('ff000000-0000-0000-0000-000000000021', 'ee000000-0000-0000-0000-000000000003', 2,
     'c0000000-0000-0000-0000-000000000002', 'SENSOR-HUM-001', 'Humidity Sensor DHT22',
     50, 'PCS', 450, 0, 'VAT', 0.07, 22500, 1575, 24075,
     0, 0, 0, 'a0000000-0000-0000-0000-000000000001',
     NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days'),

    -- PO #1 lines
    ('ff000000-0000-0000-0000-000000000101', 'ee000000-0000-0000-0000-000000000101', 1,
     'c0000000-0000-0000-0000-000000000001', 'SENSOR-TEMP-001', 'Temperature Sensor SHT30',
     100, 'PCS', 350, 0, 'VAT', 0.07, 35000, 2450, 37450,
     0, 100, 100, 'a0000000-0000-0000-0000-000000000001',
     NOW() - INTERVAL '85 days', NOW() - INTERVAL '75 days')
ON CONFLICT (id) DO NOTHING;

COMMIT;
```

### B.12 `test/fixtures/sql/11_invoices_payments.sql`

```sql
-- ============================================================
-- Invoices + Payments + Allocations
-- ============================================================
BEGIN;

-- Invoices
INSERT INTO erp_invoices
    (id, tenant_id, invoice_no, type, customer_id, supplier_id, order_id, status,
     issue_date, due_date, paid_at,
     subtotal, tax_amount, total_amount, amount_paid, amount_due, currency,
     payment_terms, created_by, created_at, updated_at)
VALUES
    -- AR #1 (PAID)
    ('aa110000-0000-0000-0000-000000000001',
     '11111111-1111-1111-1111-111111111101',
     'INV-2026-000001', 'AR', '33333333-3333-3333-3333-333333333301', NULL,
     'ee000000-0000-0000-0000-000000000001', 'PAID',
     NOW() - INTERVAL '58 days', NOW() - INTERVAL '28 days', NOW() - INTERVAL '30 days',
     65000, 4550, 69550, 69550, 0, 'THB', 'NET_30',
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '58 days', NOW() - INTERVAL '30 days'),

    -- AR #2 (ISSUED — waiting)
    ('aa110000-0000-0000-0000-000000000002',
     '11111111-1111-1111-1111-111111111101',
     'INV-2026-000002', 'AR', '33333333-3333-3333-3333-333333333302', NULL,
     'ee000000-0000-0000-0000-000000000002', 'ISSUED',
     NOW() - INTERVAL '5 days', NOW() + INTERVAL '25 days', NULL,
     8900, 623, 9523, 0, 9523, 'THB', 'NET_30',
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days'),

    -- AR #3 (OVERDUE 45 days)
    ('aa110000-0000-0000-0000-000000000003',
     '11111111-1111-1111-1111-111111111102',
     'INV-2026-000003', 'AR', '33333333-3333-3333-3333-333333333313', NULL,
     NULL, 'OVERDUE',
     NOW() - INTERVAL '75 days', NOW() - INTERVAL '45 days', NULL,
     25000, 1750, 26750, 0, 26750, 'THB', 'NET_30',
     '22222222-2222-2222-2222-222222222211',
     NOW() - INTERVAL '75 days', NOW() - INTERVAL '1 day'),

    -- ACME AR (PARTIAL)
    ('aa110000-0000-0000-0000-000000000010',
     '11111111-1111-1111-1111-111111111102',
     'INV-2026-000101', 'AR', '33333333-3333-3333-3333-333333333311', NULL,
     'ee000000-0000-0000-0000-000000000201', 'PARTIAL',
     NOW() - INTERVAL '13 days', NOW() + INTERVAL '17 days', NULL,
     124000, 8680, 132680, 50000, 82680, 'THB', 'NET_30',
     '22222222-2222-2222-2222-222222222211',
     NOW() - INTERVAL '13 days', NOW() - INTERVAL '5 days'),

    -- AP #1 (PAID)
    ('aa110000-0000-0000-0000-000000000101',
     '11111111-1111-1111-1111-111111111101',
     'BILL-2026-000001', 'AP', NULL, 'b0000000-0000-0000-0000-000000000001',
     'ee000000-0000-0000-0000-000000000101', 'PAID',
     NOW() - INTERVAL '75 days', NOW() - INTERVAL '45 days', NOW() - INTERVAL '50 days',
     35000, 2450, 37450, 37450, 0, 'THB', 'NET_30',
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '75 days', NOW() - INTERVAL '50 days')
ON CONFLICT (id) DO NOTHING;

-- Invoice lines
INSERT INTO erp_invoice_lines
    (id, invoice_id, line_no, order_line_id, product_id, sku, name, quantity, uom, unit_price,
     tax_type, tax_rate, subtotal, tax_amount, line_total, created_at)
VALUES
    ('bb110000-0000-0000-0000-000000000001', 'aa110000-0000-0000-0000-000000000001', 1,
     'ff000000-0000-0000-0000-000000000001',
     'c0000000-0000-0000-0000-000000000001', 'SENSOR-TEMP-001', 'Temperature Sensor SHT30',
     50, 'PCS', 650, 'VAT', 0.07, 32500, 2275, 34775, NOW() - INTERVAL '58 days'),

    ('bb110000-0000-0000-0000-000000000002', 'aa110000-0000-0000-0000-000000000001', 2,
     'ff000000-0000-0000-0000-000000000002',
     'c0000000-0000-0000-0000-000000000010', 'GATEWAY-MQTT-001', 'IoT Gateway MQTT 4G',
     5, 'PCS', 6500, 'VAT', 0.07, 32500, 2275, 34775, NOW() - INTERVAL '58 days'),

    ('bb110000-0000-0000-0000-000000000010', 'aa110000-0000-0000-0000-000000000002', 1,
     'ff000000-0000-0000-0000-000000000010',
     'c0000000-0000-0000-0000-000000000003', 'SENSOR-SOIL-001', 'Soil Moisture Sensor v2',
     10, 'PCS', 890, 'VAT', 0.07, 8900, 623, 9523, NOW() - INTERVAL '5 days'),

    ('bb110000-0000-0000-0000-000000000020', 'aa110000-0000-0000-0000-000000000101', 1,
     'ff000000-0000-0000-0000-000000000101',
     'c0000000-0000-0000-0000-000000000001', 'SENSOR-TEMP-001', 'Temperature Sensor SHT30',
     100, 'PCS', 350, 'VAT', 0.07, 35000, 2450, 37450, NOW() - INTERVAL '75 days')
ON CONFLICT (id) DO NOTHING;

-- Payments
INSERT INTO erp_payments
    (id, tenant_id, payment_no, direction, customer_id, supplier_id, method, status,
     amount, allocated, unallocated, currency,
     payment_date, reference, bank_account,
     created_by, created_at, updated_at, confirmed_at)
VALUES
    -- IN #1 (fully allocated)
    ('cc110000-0000-0000-0000-000000000001',
     '11111111-1111-1111-1111-111111111101',
     'PAY-2026-000001', 'IN',
     '33333333-3333-3333-3333-333333333301', NULL,
     'BANK_TRANSFER', 'CONFIRMED',
     69550, 69550, 0, 'THB',
     NOW() - INTERVAL '30 days', 'TXN-2026-ABC-001', 'KBANK-001',
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '30 days', NOW() - INTERVAL '30 days', NOW() - INTERVAL '30 days'),

    -- IN #2 (partial — ACME)
    ('cc110000-0000-0000-0000-000000000002',
     '11111111-1111-1111-1111-111111111102',
     'PAY-2026-000002', 'IN',
     '33333333-3333-3333-3333-333333333311', NULL,
     'BANK_TRANSFER', 'CONFIRMED',
     50000, 50000, 0, 'THB',
     NOW() - INTERVAL '5 days', 'TXN-2026-ACME-001', 'SCB-001',
     '22222222-2222-2222-2222-222222222211',
     NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days'),

    -- OUT #1
    ('cc110000-0000-0000-0000-000000000101',
     '11111111-1111-1111-1111-111111111101',
     'PAY-2026-000101', 'OUT',
     NULL, 'b0000000-0000-0000-0000-000000000001',
     'BANK_TRANSFER', 'CONFIRMED',
     37450, 37450, 0, 'THB',
     NOW() - INTERVAL '50 days', 'TXN-2026-SUP-001', 'KBANK-001',
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '50 days', NOW() - INTERVAL '50 days', NOW() - INTERVAL '50 days')
ON CONFLICT (id) DO NOTHING;

-- Allocations
INSERT INTO erp_payment_allocations (id, payment_id, invoice_id, amount, currency, created_at)
VALUES
    ('dd110000-0000-0000-0000-000000000001',
     'cc110000-0000-0000-0000-000000000001',
     'aa110000-0000-0000-0000-000000000001',
     69550, 'THB', NOW() - INTERVAL '30 days'),

    ('dd110000-0000-0000-0000-000000000002',
     'cc110000-0000-0000-0000-000000000002',
     'aa110000-0000-0000-0000-000000000010',
     50000, 'THB', NOW() - INTERVAL '5 days'),

    ('dd110000-0000-0000-0000-000000000101',
     'cc110000-0000-0000-0000-000000000101',
     'aa110000-0000-0000-0000-000000000101',
     37450, 'THB', NOW() - INTERVAL '50 days')
ON CONFLICT (id) DO NOTHING;

COMMIT;
```

### B.13 `test/fixtures/sql/12_technicians_logistics.sql`

```sql
-- ============================================================
-- Technicians + Shipments
-- ============================================================
BEGIN;

INSERT INTO iotlogistics_technicians
    (id, tenant_id, user_id, code, name, phone, email, level, skills, zone,
     is_available, is_active, max_jobs_per_day, current_load, gps_last, gps_updated_at,
     created_at, updated_at)
VALUES
    ('ee110000-0000-0000-0000-000000000001',
     '11111111-1111-1111-1111-111111111101', '22222222-2222-2222-2222-222222222202',
     'TECH-001', 'ช่างสมปอง ขยันงาน', '0811111111', 'sompong@demo.local',
     'SENIOR',
     '["SMART_FARM","NETWORK","ELECTRICAL"]',
     'BANGKOK',
     TRUE, TRUE, 6, 2,
     '{"lat":13.7563,"lng":100.5018,"accuracy":10,"ts":"2026-01-15T10:00:00Z"}',
     NOW() - INTERVAL '5 minutes',
     NOW() - INTERVAL '90 days', NOW() - INTERVAL '5 minutes'),

    ('ee110000-0000-0000-0000-000000000002',
     '11111111-1111-1111-1111-111111111101', '22222222-2222-2222-2222-222222222203',
     'TECH-002', 'ช่างวิชัย เก่งจริง', '0822222222', 'wichai@demo.local',
     'JUNIOR',
     '["NETWORK"]',
     'BANGKOK',
     TRUE, TRUE, 4, 1,
     '{"lat":13.8500,"lng":100.6000,"accuracy":20,"ts":"2026-01-15T10:00:00Z"}',
     NOW() - INTERVAL '15 minutes',
     NOW() - INTERVAL '60 days', NOW() - INTERVAL '15 minutes'),

    ('ee110000-0000-0000-0000-000000000003',
     '11111111-1111-1111-1111-111111111101', NULL,
     'TECH-003', 'ช่างประเสริฐ มืออาชีพ', '0833333333', 'prasert@demo.local',
     'LEAD',
     '["SMART_FARM","SMART_BUILDING","NETWORK","ELECTRICAL","SOLAR"]',
     'CHIANG_MAI',
     TRUE, TRUE, 8, 0,
     '{"lat":18.7883,"lng":98.9853,"accuracy":10,"ts":"2026-01-15T10:00:00Z"}',
     NOW() - INTERVAL '3 minutes',
     NOW() - INTERVAL '30 days', NOW() - INTERVAL '3 minutes'),

    ('ee110000-0000-0000-0000-000000000011',
     '11111111-1111-1111-1111-111111111102', '22222222-2222-2222-2222-222222222212',
     'ACME-TECH-001', 'ช่าง ACME 1', '0844444444', 'tech1@acme.co.th',
     'SENIOR',
     '["SMART_FARM","NETWORK","ELECTRICAL","SOLAR"]',
     'NORTHEAST',
     TRUE, TRUE, 8, 3,
     '{"lat":15.2287,"lng":104.8574,"accuracy":10,"ts":"2026-01-15T10:00:00Z"}',
     NOW() - INTERVAL '10 minutes',
     NOW() - INTERVAL '180 days', NOW() - INTERVAL '10 minutes')
ON CONFLICT (id) DO NOTHING;

-- Shipments
INSERT INTO iotlogistics_shipments
    (id, tenant_id, shipment_no, customer_id, site_id, from_warehouse_id,
     carrier, tracking_no, status,
     scheduled_at, dispatched_at, delivered_at, delivered_to, signature_url, gps_last,
     metadata, created_by, created_at, updated_at)
VALUES
    -- Delivered + installed
    ('ff110000-0000-0000-0000-000000000001',
     '11111111-1111-1111-1111-111111111101',
     'SHP-2026-000001', '33333333-3333-3333-3333-333333333301',
     '55555555-5555-5555-5555-555555555501',
     'a0000000-0000-0000-0000-000000000001',
     'Kerry Express', 'KER-2026-001', 'INSTALLED',
     NOW() - INTERVAL '55 days', NOW() - INTERVAL '54 days', NOW() - INTERVAL '52 days',
     'คุณวิชัย', 'https://cdn.demo/sig-001.png',
     '{"lat":14.0500,"lng":100.6000,"accuracy":10,"ts":"2026-01-15T10:00:00Z"}',
     '{"fixture": true}', '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '55 days', NOW() - INTERVAL '50 days'),

    -- In transit
    ('ff110000-0000-0000-0000-000000000002',
     '11111111-1111-1111-1111-111111111101',
     'SHP-2026-000002', '33333333-3333-3333-3333-333333333302',
     '55555555-5555-5555-5555-555555555502',
     'a0000000-0000-0000-0000-000000000001',
     'Flash Express', 'FL-2026-002', 'IN_TRANSIT',
     NOW() - INTERVAL '1 day', NOW() - INTERVAL '12 hours', NULL, NULL, NULL,
     '{"lat":17.5000,"lng":99.5000,"accuracy":15,"ts":"2026-01-15T10:00:00Z"}',
     '{"fixture": true}', '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '1 day', NOW() - INTERVAL '12 hours'),

    -- Pending
    ('ff110000-0000-0000-0000-000000000003',
     '11111111-1111-1111-1111-111111111101',
     'SHP-2026-000003', '33333333-3333-3333-3333-333333333304',
     '55555555-5555-5555-5555-555555555503',
     'a0000000-0000-0000-0000-000000000001',
     'Kerry Express', 'KER-2026-003', 'PENDING',
     NOW() + INTERVAL '2 days', NULL, NULL, NULL, NULL, NULL,
     '{"fixture": true}', '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours')
ON CONFLICT (id) DO NOTHING;

-- Shipment items
INSERT INTO iotlogistics_shipment_items
    (id, shipment_id, product_id, serial_no, qty, uom, notes, created_at)
VALUES
    ('ff220000-0000-0000-0000-000000000001',
     'ff110000-0000-0000-0000-000000000001',
     'c0000000-0000-0000-0000-000000000001', 'SN-DEMO-001', 1, 'PCS', NULL, NOW() - INTERVAL '55 days'),

    ('ff220000-0000-0000-0000-000000000002',
     'ff110000-0000-0000-0000-000000000001',
     'c0000000-0000-0000-0000-000000000001', 'SN-DEMO-002', 1, 'PCS', NULL, NOW() - INTERVAL '55 days'),

    ('ff220000-0000-0000-0000-000000000003',
     'ff110000-0000-0000-0000-000000000001',
     'c0000000-0000-0000-0000-000000000010', 'SN-DEMO-GW-001', 1, 'PCS', NULL, NOW() - INTERVAL '55 days')
ON CONFLICT (id) DO NOTHING;

COMMIT;
```

### B.14 `test/fixtures/sql/13_installation_jobs.sql`

```sql
-- ============================================================
-- Installation jobs (mixed statuses)
-- ============================================================
BEGIN;

INSERT INTO iotlogistics_installation_jobs
    (id, tenant_id, job_no, shipment_id, customer_id, site_id, device_ids,
     technician_id, job_type, priority, status,
     scheduled_at, sla_due_at,
     started_at, completed_at, started_gps, completed_gps,
     checklist, photos, signature_url, notes,
     assigned_by, assigned_at, created_by, created_at, updated_at)
VALUES
    -- DONE (installed successfully)
    ('ab220000-0000-0000-0000-000000000001',
     '11111111-1111-1111-1111-111111111101',
     'JOB-2026-000001', 'ff110000-0000-0000-0000-000000000001',
     '33333333-3333-3333-3333-333333333301', '55555555-5555-5555-5555-555555555501',
     '["f0000000-0000-0000-0000-000000000001","f0000000-0000-0000-0000-000000000002","f0000000-0000-0000-0000-000000000030"]',
     'ee110000-0000-0000-0000-000000000001',
     'INSTALLATION', 'HIGH', 'DONE',
     NOW() - INTERVAL '52 days', NOW() - INTERVAL '50 days',
     NOW() - INTERVAL '51 days', NOW() - INTERVAL '51 days' + INTERVAL '2 hours',
     '{"lat":13.7563,"lng":100.5018,"accuracy":10,"ts":"2026-01-15T10:00:00Z"}',
     '{"lat":13.7563,"lng":100.5018,"accuracy":10,"ts":"2026-01-15T12:00:00Z"}',
     '[{"key":"power_check","label":"ตรวจไฟ","required":true,"checked":true},{"key":"network_check","label":"ตรวจเน็ต","required":true,"checked":true},{"key":"device_mount","label":"ยึดอุปกรณ์","required":true,"checked":true},{"key":"sensor_test","label":"ทดสอบ sensor","required":true,"checked":true},{"key":"config_check","label":"ตรวจ config","required":true,"checked":true},{"key":"photo_evidence","label":"ถ่ายภาพ","required":true,"checked":true,"photo_url":"https://cdn.demo/photo1.jpg"},{"key":"customer_demo","label":"demo ให้ลูกค้า","required":false,"checked":true}]',
     '["https://cdn.demo/photo1.jpg","https://cdn.demo/photo2.jpg"]',
     'https://cdn.demo/sig-001.png',
     'ติดตั้งเสร็จเรียบร้อย',
     '22222222-2222-2222-2222-222222222201', NOW() - INTERVAL '52 days',
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '52 days', NOW() - INTERVAL '51 days'),

    -- IN_PROGRESS
    ('ab220000-0000-0000-0000-000000000002',
     '11111111-1111-1111-1111-111111111101',
     'JOB-2026-000002', NULL,
     '33333333-3333-3333-3333-333333333302', '55555555-5555-5555-5555-555555555502',
     '["f0000000-0000-0000-0000-000000000020"]',
     'ee110000-0000-0000-0000-000000000003',
     'REPAIR', 'HIGH', 'IN_PROGRESS',
     NOW() - INTERVAL '4 hours', NOW() + INTERVAL '20 hours',
     NOW() - INTERVAL '1 hour', NULL,
     '{"lat":18.7883,"lng":98.9853,"accuracy":10,"ts":"2026-01-15T10:00:00Z"}',
     NULL,
     '[{"key":"power_check","label":"ตรวจไฟ","required":true,"checked":true},{"key":"device_mount","label":"ยึดอุปกรณ์","required":true,"checked":false},{"key":"sensor_test","label":"ทดสอบ sensor","required":true,"checked":false}]',
     '[]', '',
     'กำลังซ่อม sensor',
     '22222222-2222-2222-2222-222222222202', NOW() - INTERVAL '3 hours',
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '4 hours', NOW() - INTERVAL '1 hour'),

    -- SCHEDULED (unassigned)
    ('ab220000-0000-0000-0000-000000000003',
     '11111111-1111-1111-1111-111111111101',
     'JOB-2026-000003', 'ff110000-0000-0000-0000-000000000002',
     '33333333-3333-3333-3333-333333333302', '55555555-5555-5555-5555-555555555502',
     '[]', NULL,
     'INSTALLATION', 'NORMAL', 'SCHEDULED',
     NOW() + INTERVAL '2 days', NOW() + INTERVAL '4 days',
     NULL, NULL, NULL, NULL,
     '[]', '[]', '', '',
     NULL, NULL,
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),

    -- FAILED (SLA breach)
    ('ab220000-0000-0000-0000-000000000004',
     '11111111-1111-1111-1111-111111111101',
     'JOB-2026-000004', NULL,
     '33333333-3333-3333-3333-333333333304', '55555555-5555-5555-5555-555555555503',
     '[]',
     'ee110000-0000-0000-0000-000000000002',
     'INSPECTION', 'LOW', 'FAILED',
     NOW() - INTERVAL '5 days', NOW() - INTERVAL '3 days',
     NOW() - INTERVAL '5 days', NOW() - INTERVAL '4 days',
     '{"lat":14.9799,"lng":102.0977,"accuracy":15,"ts":"2026-01-15T10:00:00Z"}',
     NULL,
     '[]', '[]', '', 'ไม่สามารถเข้าพื้นที่ได้',
     '22222222-2222-2222-2222-222222222201', NOW() - INTERVAL '5 days',
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '5 days', NOW() - INTERVAL '4 days'),

    -- PAUSED
    ('ab220000-0000-0000-0000-000000000005',
     '11111111-1111-1111-1111-111111111101',
     'JOB-2026-000005', NULL,
     '33333333-3333-3333-3333-333333333301', '55555555-5555-5555-5555-555555555501',
     '[]',
     'ee110000-0000-0000-0000-000000000001',
     'MAINTENANCE', 'NORMAL', 'PAUSED',
     NOW() - INTERVAL '2 hours', NOW() + INTERVAL '70 hours',
     NOW() - INTERVAL '90 minutes', NULL,
     '{"lat":13.7563,"lng":100.5018,"accuracy":10,"ts":"2026-01-15T10:00:00Z"}',
     NULL,
     '[]', '[]', '', 'รอลูกค้าเปิดประตู',
     '22222222-2222-2222-2222-222222222201', NOW() - INTERVAL '2 hours',
     '22222222-2222-2222-2222-222222222201',
     NOW() - INTERVAL '2 hours', NOW() - INTERVAL '90 minutes'),

    -- ACME — DONE
    ('ab220000-0000-0000-0000-000000000011',
     '11111111-1111-1111-1111-111111111102',
     'JOB-2026-000101', NULL,
     '33333333-3333-3333-3333-333333333311', '55555555-5555-5555-5555-555555555511',
     '[]',
     'ee110000-0000-0000-0000-000000000011',
     'INSTALLATION', 'NORMAL', 'DONE',
     NOW() - INTERVAL '140 days', NOW() - INTERVAL '140 days' + INTERVAL '48 hours',
     NOW() - INTERVAL '140 days', NOW() - INTERVAL '140 days' + INTERVAL '3 hours',
     NULL, NULL, '[]', '[]', '', 'ACME installation done',
     '22222222-2222-2222-2222-222222222211', NOW() - INTERVAL '140 days',
     '22222222-2222-2222-2222-222222222211',
     NOW() - INTERVAL '140 days', NOW() - INTERVAL '139 days')
ON CONFLICT (id) DO NOTHING;

COMMIT;
```

### B.15 `test/fixtures/sql/14_kpi_snapshots.sql`

```sql
-- ============================================================
-- KPI Definitions + Snapshots (30 days)
-- ============================================================
BEGIN;

-- KPI definitions (platform-wide)
INSERT INTO report_kpi_definitions (id, tenant_id, code, name, category, unit, is_active, created_at)
VALUES
    (gen_random_uuid(), NULL, 'MRR', 'Monthly Recurring Revenue', 'FINANCIAL', 'THB', TRUE, NOW() - INTERVAL '365 days'),
    (gen_random_uuid(), NULL, 'ACTIVE_CUSTOMERS', 'Active Customers', 'CUSTOMER', 'count', TRUE, NOW() - INTERVAL '365 days'),
    (gen_random_uuid(), NULL, 'NEW_CUSTOMERS', 'New Customers', 'CUSTOMER', 'count', TRUE, NOW() - INTERVAL '365 days'),
    (gen_random_uuid(), NULL, 'CHURN_RATE', 'Churn Rate', 'CUSTOMER', '%', TRUE, NOW() - INTERVAL '365 days'),
    (gen_random_uuid(), NULL, 'DEVICE_UPTIME_PCT', 'Device Uptime', 'IOT', '%', TRUE, NOW() - INTERVAL '365 days'),
    (gen_random_uuid(), NULL, 'DEVICE_ONLINE', 'Devices Online', 'IOT', 'count', TRUE, NOW() - INTERVAL '365 days'),
    (gen_random_uuid(), NULL, 'ALERT_COUNT', 'Alert Count', 'IOT', 'count', TRUE, NOW() - INTERVAL '365 days'),
    (gen_random_uuid(), NULL, 'TICKET_OPEN', 'Open Tickets', 'SERVICE', 'count', TRUE, NOW() - INTERVAL '365 days'),
    (gen_random_uuid(), NULL, 'TICKET_SLA_PCT', 'Ticket SLA Compliance', 'SERVICE', '%', TRUE, NOW() - INTERVAL '365 days'),
    (gen_random_uuid(), NULL, 'INSTALLATION_SLA_PCT', 'Installation SLA Compliance', 'LOGISTICS', '%', TRUE, NOW() - INTERVAL '365 days')
ON CONFLICT (tenant_id, code) DO NOTHING;

-- KPI snapshots — 30 days ของ Demo tenant (sample 90 rows)
-- ใช้ generate_series เพื่อความกระชับ
INSERT INTO report_kpi_snapshots (tenant_id, kpi_code, value, dimensions, dimensions_key, snapshot_date, computed_at, source)
SELECT
    '11111111-1111-1111-1111-111111111101',
    'MRR',
    24000 + (random() * 4000 - 2000)::numeric(20,4),
    '{}'::jsonb,
    '',
    d::date,
    d::date + INTERVAL '2 hours',
    'scheduler'
FROM generate_series(NOW() - INTERVAL '30 days', NOW(), '1 day') AS d
ON CONFLICT (tenant_id, kpi_code, dimensions_key, snapshot_date) DO NOTHING;

INSERT INTO report_kpi_snapshots (tenant_id, kpi_code, value, dimensions, dimensions_key, snapshot_date, computed_at, source)
SELECT
    '11111111-1111-1111-1111-111111111101',
    'ACTIVE_CUSTOMERS',
    3 + (random() * 2)::numeric(20,4),
    '{}'::jsonb,
    '',
    d::date,
    d::date + INTERVAL '2 hours',
    'scheduler'
FROM generate_series(NOW() - INTERVAL '30 days', NOW(), '1 day') AS d
ON CONFLICT (tenant_id, kpi_code, dimensions_key, snapshot_date) DO NOTHING;

INSERT INTO report_kpi_snapshots (tenant_id, kpi_code, value, dimensions, dimensions_key, snapshot_date, computed_at, source)
SELECT
    '11111111-1111-1111-1111-111111111101',
    'DEVICE_UPTIME_PCT',
    (85 + random() * 12)::numeric(20,4),
    '{}'::jsonb,
    '',
    d::date,
    d::date + INTERVAL '2 hours',
    'scheduler'
FROM generate_series(NOW() - INTERVAL '30 days', NOW(), '1 day') AS d
ON CONFLICT (tenant_id, kpi_code, dimensions_key, snapshot_date) DO NOTHING;

INSERT INTO report_kpi_snapshots (tenant_id, kpi_code, value, dimensions, dimensions_key, snapshot_date, computed_at, source)
SELECT
    '11111111-1111-1111-1111-111111111101',
    'DEVICE_ONLINE',
    (15 + random() * 5)::numeric(20,4),
    '{}'::jsonb,
    '',
    d::date,
    d::date + INTERVAL '2 hours',
    'scheduler'
FROM generate_series(NOW() - INTERVAL '30 days', NOW(), '1 day') AS d
ON CONFLICT (tenant_id, kpi_code, dimensions_key, snapshot_date) DO NOTHING;

INSERT INTO report_kpi_snapshots (tenant_id, kpi_code, value, dimensions, dimensions_key, snapshot_date, computed_at, source)
SELECT
    '11111111-1111-1111-1111-111111111101',
    'ALERT_COUNT',
    (random() * 8)::int::numeric(20,4),
    '{}'::jsonb,
    '',
    d::date,
    d::date + INTERVAL '2 hours',
    'scheduler'
FROM generate_series(NOW() - INTERVAL '30 days', NOW(), '1 day') AS d
ON CONFLICT (tenant_id, kpi_code, dimensions_key, snapshot_date) DO NOTHING;

INSERT INTO report_kpi_snapshots (tenant_id, kpi_code, value, dimensions, dimensions_key, snapshot_date, computed_at, source)
SELECT
    '11111111-1111-1111-1111-111111111101',
    'TICKET_SLA_PCT',
    (88 + random() * 10)::numeric(20,4),
    '{}'::jsonb,
    '',
    d::date,
    d::date + INTERVAL '2 hours',
    'scheduler'
FROM generate_series(NOW() - INTERVAL '30 days', NOW(), '1 day') AS d
ON CONFLICT (tenant_id, kpi_code, dimensions_key, snapshot_date) DO NOTHING;

COMMIT;
```

### B.16 `test/fixtures/sql/15_dashboards.sql`

```sql
-- ============================================================
-- Dashboards (3 dashboards per tenant)
-- ============================================================
BEGIN;

INSERT INTO report_dashboards
    (id, tenant_id, code, name, description, owner_id, is_default, is_public, layout, widgets, created_at, updated_at)
VALUES
    -- Demo — default dashboard
    ('cd110000-0000-0000-0000-000000000001',
     '11111111-1111-1111-1111-111111111101',
     'OVERVIEW', 'ภาพรวมระบบ', 'KPIs หลัก + Charts',
     '22222222-2222-2222-2222-222222222201',
     TRUE, TRUE,
     '{"columns":12,"rows":4,"gap":16}',
     '[
        {"id":"w1","type":"KPI_CARD","title":"MRR","kpi_codes":["MRR"],"position":{"order":0,"x":0,"y":0,"w":3,"h":2},"config":{"format":"currency"},"refresh_sec":300},
        {"id":"w2","type":"KPI_CARD","title":"Active Customers","kpi_codes":["ACTIVE_CUSTOMERS"],"position":{"order":1,"x":3,"y":0,"w":3,"h":2},"config":{},"refresh_sec":300},
        {"id":"w3","type":"KPI_CARD","title":"Devices Online","kpi_codes":["DEVICE_ONLINE"],"position":{"order":2,"x":6,"y":0,"w":3,"h":2},"config":{},"refresh_sec":60},
        {"id":"w4","type":"KPI_CARD","title":"Uptime","kpi_codes":["DEVICE_UPTIME_PCT"],"position":{"order":3,"x":9,"y":0,"w":3,"h":2},"config":{"format":"percent"},"refresh_sec":300},
        {"id":"w5","type":"LINE_CHART","title":"MRR Trend","kpi_codes":["MRR"],"position":{"order":4,"x":0,"y":2,"w":8,"h":3},"config":{"period":"last_30_days"},"refresh_sec":600},
        {"id":"w6","type":"PIE_CHART","title":"Alert Distribution","kpi_codes":["ALERT_COUNT"],"position":{"order":5,"x":8,"y":2,"w":4,"h":3},"config":{},"refresh_sec":300}
     ]'::jsonb,
     NOW() - INTERVAL '30 days', NOW() - INTERVAL '1 day'),

    -- Demo — IoT dashboard
    ('cd110000-0000-0000-0000-000000000002',
     '11111111-1111-1111-1111-111111111101',
     'IOT-OPS', 'IoT Operations', 'สถานะอุปกรณ์ + Alerts',
     '22222222-2222-2222-2222-222222222202',
     FALSE, TRUE,
     '{"columns":12,"rows":4,"gap":16}',
     '[
        {"id":"w1","type":"KPI_CARD","title":"Online Devices","kpi_codes":["DEVICE_ONLINE"],"position":{"order":0,"x":0,"y":0,"w":4,"h":2},"config":{},"refresh_sec":30},
        {"id":"w2","type":"KPI_CARD","title":"Alerts Today","kpi_codes":["ALERT_COUNT"],"position":{"order":1,"x":4,"y":0,"w":4,"h":2},"config":{},"refresh_sec":30},
        {"id":"w3","type":"KPI_CARD","title":"Uptime","kpi_codes":["DEVICE_UPTIME_PCT"],"position":{"order":2,"x":8,"y":0,"w":4,"h":2},"config":{"format":"percent"},"refresh_sec":30},
        {"id":"w4","type":"BAR_CHART","title":"Uptime by Day","kpi_codes":["DEVICE_UPTIME_PCT"],"position":{"order":3,"x":0,"y":2,"w":12,"h":3},"config":{"period":"last_7_days"},"refresh_sec":300}
     ]'::jsonb,
     NOW() - INTERVAL '20 days', NOW() - INTERVAL '2 days'),

    -- ACME dashboard
    ('cd110000-0000-0000-0000-000000000011',
     '11111111-1111-1111-1111-111111111102',
     'EXEC-DASH', 'Executive Dashboard', 'สำหรับผู้บริหาร',
     '22222222-2222-2222-2222-222222222211',
     TRUE, TRUE,
     '{"columns":12,"rows":3,"gap":16}',
     '[
        {"id":"w1","type":"KPI_CARD","title":"MRR","kpi_codes":["MRR"],"position":{"order":0,"x":0,"y":0,"w":4,"h":2},"config":{"format":"currency"},"refresh_sec":600},
        {"id":"w2","type":"KPI_CARD","title":"Active Customers","kpi_codes":["ACTIVE_CUSTOMERS"],"position":{"order":1,"x":4,"y":0,"w":4,"h":2},"config":{},"refresh_sec":600},
        {"id":"w3","type":"KPI_CARD","title":"SLA Compliance","kpi_codes":["TICKET_SLA_PCT"],"position":{"order":2,"x":8,"y":0,"w":4,"h":2},"config":{"format":"percent"},"refresh_sec":600}
     ]'::jsonb,
     NOW() - INTERVAL '100 days', NOW() - INTERVAL '5 days')
ON CONFLICT (id) DO NOTHING;

COMMIT;
```

### B.17 `test/fixtures/sql/16_telemetry_influx.flux`

```flux
// ============================================================
// InfluxDB fixture — seed telemetry data
// รันด้วย: influx query --file 16_telemetry_influx.flux
// หรือผ่าน API /api/v2/write
// ============================================================

// หมายเหตุ: Flux ไม่สามารถ insert ได้ — ใช้ line protocol ผ่าน API
// ไฟล์นี้ใช้เป็นตัวอย่าง + query verify หลัง seed

// verify: นับ telemetry points ในช่วง 30 วัน
from(bucket: "iot_telemetry")
    |> range(start: -30d)
    |> filter(fn: (r) => r._measurement == "device_telemetry")
    |> filter(fn: (r) => r.tenant_id == "11111111-1111-1111-1111-111111111101")
    |> filter(fn: (r) => r._field == "value")
    |> group(columns: ["device_id", "metric"])
    |> count()
    |> yield(name: "point_counts")

// verify: latest reading ต่อ device
from(bucket: "iot_telemetry")
    |> range(start: -1h)
    |> filter(fn: (r) => r._measurement == "device_telemetry")
    |> filter(fn: (r) => r._field == "value")
    |> group(columns: ["device_id", "metric"])
    |> last()
    |> yield(name: "latest")
```

### B.18 Line Protocol Seed Script

```bash
#!/usr/bin/env bash
# test/fixtures/sql/16_seed_telemetry.sh
# Seed 30 days of telemetry data into InfluxDB

set -euo pipefail

INFLUX_URL="${INFLUX_URL:-http://localhost:8086}"
INFLUX_TOKEN="${INFLUX_TOKEN:?INFLUX_TOKEN required}"
INFLUX_ORG="${INFLUX_ORG:-icmon}"
INFLUX_BUCKET="${INFLUX_BUCKET:-iot_telemetry}"

TENANT="11111111-1111-1111-1111-111111111101"
DEVICES=(
  "f0000000-0000-0000-0000-000000000001"
  "f0000000-0000-0000-0000-000000000002"
  "f0000000-0000-0000-0000-000000000003"
)

echo "→ Seeding 30 days telemetry (10 points/hour × 3 devices × 30 days = 21,600 points)"

TMPFILE=$(mktemp)
trap "rm -f $TMPFILE" EXIT

for hours_back in $(seq 720 -1 0); do
  TS=$(date -u -d "-${hours_back} hours" +%s)
  TS_NS="${TS}000000000"

  for i in $(seq 0 9); do
    # 10 points per hour (ทุก 6 นาที)
    TS_POINT=$((TS + i * 360))
    TS_POINT_NS="${TS_POINT}000000000"
    POINT_TS="${TS_POINT_NS}"

    for DEV in "${DEVICES[@]}"; do
      TEMP=$(awk -v s="$TS_POINT" "BEGIN {print 25 + 5*sin(s/3600) + 2*rand()}")
      HUM=$(awk -v s="$TS_POINT" "BEGIN {print 60 + 15*cos(s/7200) + 3*rand()}")
      SITE="55555555-5555-5555-5555-555555555501"

      echo "device_telemetry,tenant_id=${TENANT},device_id=${DEV},site_id=${SITE},metric=temperature,unit=°C,quality=GOOD value=${TEMP} ${POINT_TS}" >> "$TMPFILE"
      echo "device_telemetry,tenant_id=${TENANT},device_id=${DEV},site_id=${SITE},metric=humidity,unit=%,quality=GOOD value=${HUM} ${POINT_TS}" >> "$TMPFILE"
    done
  done
done

echo "→ Uploading $(wc -l < $TMPFILE) points..."
curl -sX POST "${INFLUX_URL}/api/v2/write?org=${INFLUX_ORG}&bucket=${INFLUX_BUCKET}&precision=ns" \
  -H "Authorization: Token ${INFLUX_TOKEN}" \
  -H "Content-Type: text/plain; charset=utf-8" \
  --data-binary "@${TMPFILE}"

echo "✓ Telemetry seed done"
```

---

## 🅲 PART 9C — GO FIXTURE BUILDERS

### C.1 `test/fixtures/go/builders.go`

```go
package fixtures

import (
    "time"

    "github.com/google/uuid"

    customerentity "icmongolang/internal/modules/customer/domain/entity"
    customervo "icmongolang/internal/modules/customer/domain/value_object"
    deviceentity "icmongolang/internal/modules/device/domain/entity"
    devicevo "icmongolang/internal/modules/device/domain/value_object"
    erpentity "icmongolang/internal/modules/erp/domain/entity"
    erpvo "icmongolang/internal/modules/erp/domain/value_object"
    pkgentity "icmongolang/internal/modules/packagecatalog/domain/entity"
    pkgvo "icmongolang/internal/modules/packagecatalog/domain/value_object"
)

// ============================================================
// CUSTOMER BUILDER
// ============================================================

type CustomerBuilder struct {
    tenantID uuid.UUID
    code     string
    name     string
    typ      customervo.CustomerType
    segment  customervo.Segment
    taxID    string
    email    string
    phone    string
    address  customervo.Address
}

func NewCustomer() *CustomerBuilder {
    return &CustomerBuilder{
        tenantID: TenantID("demo"),
        code:     "CUS-TEST-0001",
        name:     "Test Customer",
        typ:      customervo.CustomerTypeCorporate,
        segment:  customervo.SegmentSME,
        taxID:    "0105555555555",
        email:    "test@example.local",
        phone:    "0812345678",
        address:  customervo.Address{
            Line1: "1 Test St", Province: "Bangkok",
            Postcode: "10110", Country: "TH",
        },
    }
}

func (b *CustomerBuilder) WithTenant(t string) *CustomerBuilder     { b.tenantID = TenantID(t); return b }
func (b *CustomerBuilder) WithCode(c string) *CustomerBuilder       { b.code = c; return b }
func (b *CustomerBuilder) WithName(n string) *CustomerBuilder       { b.name = n; return b }
func (b *CustomerBuilder) WithType(t customervo.CustomerType) *CustomerBuilder {
    b.typ = t
    return b
}
func (b *CustomerBuilder) WithSegment(s customervo.Segment) *CustomerBuilder {
    b.segment = s
    return b
}
func (b *CustomerBuilder) WithTaxID(t string) *CustomerBuilder { b.taxID = t; return b }
func (b *CustomerBuilder) WithEmail(e string) *CustomerBuilder { b.email = e; return b }

func (b *CustomerBuilder) Build() (*customerentity.Customer, error) {
    c, err := customerentity.NewCustomer(b.tenantID, b.typ, b.name)
    if err != nil { return nil, err }
    c.Code = customervo.CustomerCode(b.code)
    c.LegalName = b.name
    c.Phone = b.phone
    c.Address = b.address
    if b.taxID != "" {
        tid, err := customervo.NewTaxID(b.taxID)
        if err == nil { _ = c.SetTaxID(tid) }
    }
    if b.email != "" { _ = c.SetEmail(b.email) }
    return c, nil
}

// ============================================================
// PACKAGE BUILDER
// ============================================================

type PackageBuilder struct {
    tenantID *uuid.UUID
    code     string
    name     string
    tier     pkgvo.PackageTier
    cycle    pkgvo.BillingCycle
    price    float64
    currency string
    quotas   pkgvo.QuotaLimits
    trial    *pkgvo.TrialPolicy
}

func NewPackage() *PackageBuilder {
    return &PackageBuilder{
        code:     "TEST-M",
        name:     "Test Monthly",
        tier:     pkgvo.TierBasic,
        cycle:    pkgvo.CycleMonthly,
        price:    990,
        currency: "THB",
        quotas: pkgvo.QuotaLimits{
            MaxDevices: 50, MaxSites: 5, MaxUsers: 10,
            MaxAutomations: 50, StorageGB: 10,
            RetentionDays: 30, APICallsPerDay: 10000,
        },
    }
}

func (b *PackageBuilder) WithCode(c string) *PackageBuilder { b.code = c; return b }
func (b *PackageBuilder) WithName(n string) *PackageBuilder { b.name = n; return b }
func (b *PackageBuilder) WithTier(t pkgvo.PackageTier) *PackageBuilder {
    b.tier = t
    return b
}
func (b *PackageBuilder) WithPrice(p float64) *PackageBuilder       { b.price = p; return b }
func (b *PackageBuilder) WithCycle(c pkgvo.BillingCycle) *PackageBuilder {
    b.cycle = c
    return b
}
func (b *PackageBuilder) WithQuotas(q pkgvo.QuotaLimits) *PackageBuilder {
    b.quotas = q
    return b
}
func (b *PackageBuilder) WithTrial(days int, autoConvert bool) *PackageBuilder {
    b.trial = &pkgvo.TrialPolicy{
        Enabled: true, Days: days,
        AutoConvert: autoConvert, TrialTier: b.tier,
    }
    return b
}

func (b *PackageBuilder) Build() (*pkgentity.Package, error) {
    code, err := pkgvo.NewPackageCode(b.code)
    if err != nil { return nil, err }
    price, err := pkgvo.NewMoney(b.price, b.currency)
    if err != nil { return nil, err }
    p, err := pkgentity.NewPackage(
        b.tenantID, code, b.name, b.tier, b.cycle, price, b.quotas,
        uuid.Nil,
    )
    if err != nil { return nil, err }
    if b.trial != nil {
        _ = p.SetTrialPolicy(*b.trial)
    }
    return p, nil
}

// ============================================================
// DEVICE BUILDER
// ============================================================

type DeviceBuilder struct {
    tenantID uuid.UUID
    customerID uuid.UUID
    siteID   uuid.UUID
    serial   string
    name     string
    typ      devicevo.DeviceType
    protocol devicevo.Protocol
    provision bool
}

func NewDevice() *DeviceBuilder {
    return &DeviceBuilder{
        tenantID: TenantID("demo"),
        customerID: CustomerID("demo", "CUS-2026-0001"),
        siteID:   uuid.New(),
        serial:   "SN-TEST-001",
        name:     "Test Device",
        typ:      devicevo.DeviceTypeSensor,
        protocol: devicevo.ProtocolMQTT,
    }
}

func (b *DeviceBuilder) WithSerial(s string) *DeviceBuilder { b.serial = s; return b }
func (b *DeviceBuilder) WithName(n string) *DeviceBuilder   { b.name = n; return b }
func (b *DeviceBuilder) WithType(t devicevo.DeviceType) *DeviceBuilder {
    b.typ = t
    return b
}
func (b *DeviceBuilder) WithProtocol(p devicevo.Protocol) *DeviceBuilder {
    b.protocol = p
    return b
}
func (b *DeviceBuilder) Provisioned() *DeviceBuilder { b.provision = true; return b }

func (b *DeviceBuilder) Build() (*deviceentity.Device, error) {
    serial, err := devicevo.NewDeviceSerial(b.serial)
    if err != nil { return nil, err }
    d, err := deviceentity.NewDevice(
        b.tenantID, b.customerID, b.siteID, uuid.Nil,
        serial, b.typ, b.protocol, b.name,
    )
    if err != nil { return nil, err }
    if b.provision {
        token, _ := devicevo.NewDeviceToken()
        _ = d.Provision(token, "dev-"+b.serial)
    }
    return d, nil
}

// ============================================================
// ORDER BUILDER
// ============================================================

type OrderBuilder struct {
    tenantID uuid.UUID
    otype    erpvo.OrderType
    customerID *uuid.UUID
    supplierID *uuid.UUID
    orderNo  string
    lines    []OrderLineSpec
    currency string
}

type OrderLineSpec struct {
    ProductID   uuid.UUID
    SKU         string
    Name        string
    Quantity    float64
    UnitPrice   float64
    TaxType     erpvo.TaxType
    TaxRate     float64
}

func NewOrder(otype erpvo.OrderType) *OrderBuilder {
    return &OrderBuilder{
        tenantID: TenantID("demo"),
        otype:    otype,
        orderNo:  "TEST-2026-000001",
        currency: "THB",
    }
}

func (b *OrderBuilder) WithOrderNo(no string) *OrderBuilder  { b.orderNo = no; return b }
func (b *OrderBuilder) WithCustomer(id uuid.UUID) *OrderBuilder {
    b.customerID = &id
    return b
}
func (b *OrderBuilder) WithSupplier(id uuid.UUID) *OrderBuilder {
    b.supplierID = &id
    return b
}
func (b *OrderBuilder) AddLine(spec OrderLineSpec) *OrderBuilder {
    b.lines = append(b.lines, spec)
    return b
}

func (b *OrderBuilder) Build() (*erpentity.Order, error) {
    docNo, _ := erpvo.NewDocumentNumber(b.otype.DocTypeFor(), 2026, 1)
    o, err := erpentity.NewOrder(b.tenantID, docNo, b.otype, time.Now(), b.currency, uuid.Nil)
    if err != nil { return nil, err }

    if b.otype.IsSales() && b.customerID != nil {
        _ = o.SetCustomer(*b.customerID, erpvo.Address{Country: "TH"})
    }
    if b.otype.IsPurchase() && b.supplierID != nil {
        _ = o.SetSupplier(*b.supplierID, erpvo.Address{Country: "TH"})
    }

    for _, spec := range b.lines {
        line, err := erpentity.NewOrderLine(
            o.ID, 0, spec.ProductID, spec.SKU, spec.Name, "PCS",
            spec.Quantity, spec.UnitPrice, spec.TaxType, spec.TaxRate,
        )
        if err != nil { return nil, err }
        if err := o.AddLine(line); err != nil { return nil, err }
    }
    return o, nil
}
```

### C.2 `test/fixtures/go/datasets.go`

```go
package fixtures

import (
    "context"

    "github.com/google/uuid"

    pkgentity "icmongolang/internal/modules/packagecatalog/domain/entity"
    erpentity "icmongolang/internal/modules/erp/domain/entity"
    erpvo "icmongolang/internal/modules/erp/domain/value_object"
)

// ============================================================
// Pre-built Datasets
// ============================================================

// StandardPackages — 5 packages พร้อมใช้
func StandardPackages() []*pkgentity.Package {
    specs := []struct {
        code, name string
        tier       pkgvo.PackageTier
        cycle      pkgvo.BillingCycle
        price      float64
        quotas     pkgvo.QuotaLimits
    }{
        {"FREE", "Free", pkgvo.TierFree, pkgvo.CycleMonthly, 0,
            pkgvo.QuotaLimits{MaxDevices: 5, MaxSites: 1, MaxUsers: 2, StorageGB: 1, RetentionDays: 7}},
        {"BASIC-M", "Basic Monthly", pkgvo.TierBasic, pkgvo.CycleMonthly, 990,
            pkgvo.QuotaLimits{MaxDevices: 50, MaxSites: 5, MaxUsers: 10, StorageGB: 10, RetentionDays: 30}},
        {"BASIC-Y", "Basic Yearly", pkgvo.TierBasic, pkgvo.CycleYearly, 9900,
            pkgvo.QuotaLimits{MaxDevices: 50, MaxSites: 5, MaxUsers: 10, StorageGB: 10, RetentionDays: 30}},
        {"PRO-M", "Pro Monthly", pkgvo.TierPro, pkgvo.CycleMonthly, 4900,
            pkgvo.QuotaLimits{MaxDevices: 500, MaxSites: 50, MaxUsers: 100, StorageGB: 100, RetentionDays: 90}},
        {"ENT-M", "Enterprise", pkgvo.TierEnterprise, pkgvo.CycleMonthly, 19900,
            pkgvo.QuotaLimits{MaxDevices: -1, MaxSites: -1, MaxUsers: -1, StorageGB: 1000, RetentionDays: 365}},
    }
    out := make([]*pkgentity.Package, 0, len(specs))
    for _, s := range specs {
        p, err := NewPackage().
            WithCode(s.code).WithName(s.name).
            WithTier(s.tier).WithCycle(s.cycle).
            WithPrice(s.price).WithQuotas(s.quotas).
            Build()
        if err != nil { continue }
        out = append(out, p)
    }
    return out
}

// ============================================================
// Scenario: Full Customer Journey
// ============================================================

type ScenarioCustomerJourney struct {
    TenantID     uuid.UUID
    Customer     interface{}
    Package      *pkgentity.Package
    Order        *erpentity.Order
    Devices      []interface{}
}

// BuildCustomerJourney — สร้าง scenario ครบวงจร
func BuildCustomerJourney(ctx context.Context, deps Dependencies) (*ScenarioCustomerJourney, error) {
    // 1. Create package
    pkg, err := NewPackage().WithCode("TEST-PRO").WithTier(pkgvo.TierPro).Build()
    if err != nil { return nil, err }
    if err := deps.PkgRepo.Save(ctx, pkg); err != nil { return nil, err }

    // 2. Create customer
    cust, err := NewCustomer().WithCode("CUS-JOURNEY-001").Build()
    if err != nil { return nil, err }
    if err := deps.CustRepo.Save(ctx, cust); err != nil { return nil, err }

    // 3. Create order
    order, err := NewOrder(erpvo.OrderTypeSales).
        WithOrderNo("SO-2026-JOURNEY").
        WithCustomer(cust.ID).
        AddLine(OrderLineSpec{
            ProductID: uuid.New(), SKU: "TEST-001",
            Name: "Test Product", Quantity: 10, UnitPrice: 100,
            TaxType: erpvo.TaxTypeVAT, TaxRate: 0.07,
        }).
        Build()
    if err != nil { return nil, err }
    if err := deps.OrderRepo.Save(ctx, order); err != nil { return nil, err }

    return &ScenarioCustomerJourney{
        TenantID: cust.TenantID,
        Customer: cust, Package: pkg, Order: order,
    }, nil
}

// ============================================================
// Common Dependencies holder
// ============================================================

type Dependencies struct {
    CustRepo  interface{ Save(context.Context, interface{}) error }
    PkgRepo   interface{ Save(context.Context, *pkgentity.Package) error }
    OrderRepo interface{ Save(context.Context, *erpentity.Order) error }
}
```

---

## 🅳 PART 9D — TEST SCENARIOS

### D.1 Integration Test Scenarios

```go
// test/integration/scenarios_test.go
//go:build integration

package integration_test

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/suite"

    "icmongolang/test/fixtures"
)

type ScenarioSuite struct {
    suite.Suite
    ctx context.Context
    env *TestEnv
}

func TestScenarioSuite(t *testing.T) {
    suite.Run(t, new(ScenarioSuite))
}

func (s *ScenarioSuite) SetupSuite() {
    s.ctx = context.Background()
    s.env = NewTestEnv(s.T())
    s.Require().NoError(s.env.LoadFixtures())
}

func (s *ScenarioSuite) TearDownSuite() {
    s.env.Cleanup()
}

// ============================================================
// Scenario 1: Customer Onboarding
// ============================================================
func (s *ScenarioSuite) TestScenario_CustomerOnboarding() {
    t := s.T()
    tenantID := fixtures.TenantID("demo")

    // 1. Create customer
    cust, err := fixtures.NewCustomer().WithCode("CUS-SCENARIO-001").Build()
    s.Require().NoError(err)
    s.Require().NoError(s.env.CustomerRepo.Save(s.ctx, cust))

    // 2. Qualify → Activate
    s.Require().NoError(cust.Qualify())
    s.Require().NoError(cust.Activate())
    s.Require().NoError(s.env.CustomerRepo.Save(s.ctx, cust))
    s.Equal("ACTIVE", string(cust.Status))

    // 3. Subscribe to package
    pkg, err := s.env.PackageRepo.FindByCode(s.ctx, tenantID, "PRO-M")
    s.Require().NoError(err)

    sub, err := fixtures.NewSubscription().
        WithTenant(tenantID).
        WithCustomer(cust.ID).
        WithPackage(pkg.ID).
        Build()
    s.Require().NoError(err)
    s.Require().NoError(s.env.SubscriptionRepo.Save(s.ctx, sub))
    s.Equal("ACTIVE", string(sub.Status))

    // 4. Verify quota available
    quota := s.env.QuotaSvc.Check(sub, pkg, "devices")
    s.True(quota.Allowed)

    // 5. Verify Kafka event published
    s.Eventually(func() bool {
        return s.env.KafkaMock.HasEvent("customer.created")
    }, 3*time.Second, 100*time.Millisecond)
}

// ============================================================
// Scenario 2: Order-to-Cash End-to-End
// ============================================================
func (s *ScenarioSuite) TestScenario_OrderToCash() {
    tenantID := fixtures.TenantID("demo")
    customerCode := "CUS-2026-0001"

    // 1. Load existing customer
    cust, err := s.env.CustomerRepo.FindByCode(s.ctx, tenantID, customerCode)
    s.Require().NoError(err)

    // 2. Create product + inventory
    prod, err := fixtures.NewProduct().
        WithSKU("TEST-OTC-001").
        WithSellPrice(500).
        Build()
    s.Require().NoError(err)
    s.Require().NoError(s.env.ProductRepo.Save(s.ctx, prod))

    // Seed inventory 100 units
    whID := fixtures.WarehouseID("demo", "WH-DEMO")
    s.Require().NoError(s.env.InventorySvc.SeedStock(s.ctx, tenantID, prod.ID, whID, 100, 300))

    // 3. Create Sales Order
    order, err := fixtures.NewOrder(erpvo.OrderTypeSales).
        WithOrderNo("SO-2026-OTC-001").
        WithCustomer(cust.ID).
        AddLine(fixtures.OrderLineSpec{
            ProductID: prod.ID, SKU: prod.SKU.String(), Name: prod.Name,
            Quantity: 10, UnitPrice: 500,
            TaxType: erpvo.TaxTypeVAT, TaxRate: 0.07,
        }).
        Build()
    s.Require().NoError(err)
    s.Require().NoError(s.env.OrderRepo.Save(s.ctx, order))
    s.Equal(5350.0, order.TotalAmount) // 5000 + 350

    // 4. Confirm
    s.Require().NoError(order.SubmitForApproval())
    s.Require().NoError(order.Confirm(uuid.New()))
    s.Require().NoError(s.env.OrderRepo.Save(s.ctx, order))
    s.Equal(erpvo.OrderStatusConfirmed, order.Status)

    // 5. Ship
    lineIDs := map[uuid.UUID]float64{order.Lines[0].ID: 10}
    s.Require().NoError(s.env.ShipUC.Execute(s.ctx, ShipOrderInput{
        TenantID: tenantID.String(), OrderID: order.ID.String(),
        Lines: fixtures.UUIDMapToStringMap(lineIDs), ActorID: uuid.New().String(),
    }))

    // Verify inventory deducted
    inv, _ := s.env.InventoryRepo.FindByKey(s.ctx, tenantID, prod.ID, whID)
    s.Equal(90.0, inv.QtyOnHand)

    // 6. Create invoice
    inv2, err := s.env.InvoiceUC.Execute(s.ctx, CreateInvoiceInput{
        TenantID: tenantID.String(), InvoiceType: "AR",
        CustomerID: cust.ID.String(), OrderID: order.ID.String(),
        Lines: []InvoiceLineDTO{{
            ProductID: prod.ID.String(), Quantity: 10, UnitPrice: 500,
            OrderLineID: order.Lines[0].ID.String(),
        }},
        ActorID: uuid.New().String(),
    })
    s.Require().NoError(err)
    s.NotNil(inv2)

    // 7. Record payment
    pay, err := s.env.PaymentUC.Execute(s.ctx, CreatePaymentInput{
        TenantID: tenantID.String(), Direction: "IN",
        CustomerID: cust.ID.String(), Method: "BANK_TRANSFER",
        Amount: inv2.TotalAmount,
        Allocations: []AllocationDTO{{
            InvoiceID: inv2.ID, Amount: inv2.TotalAmount,
        }},
        ActorID: uuid.New().String(),
    })
    s.Require().NoError(err)
    s.True(pay.Allocated > 0)

    // Verify invoice paid
    invReload, _ := s.env.InvoiceRepo.FindByID(s.ctx, tenantID, uuid.MustParse(inv2.ID))
    s.Equal(erpvo.InvoiceStatusPaid, invReload.Status)

    // Verify journal entries posted
    journals, _ := s.env.JournalRepo.List(s.ctx, tenantID, time.Now().Add(-1*time.Hour), time.Now(), 1, 100)
    s.GreaterOrEqual(len(journals), 3) // AR Invoice + COGS + Payment
}

// ============================================================
// Scenario 3: Concurrent shadow update (optimistic lock)
// ============================================================
func (s *ScenarioSuite) TestScenario_ShadowConcurrentUpdate() {
    tenantID := fixtures.TenantID("demo")
    deviceID := fixtures.DeviceID("demo", "SN-DEMO-040")

    // Two users update shadow concurrently
    err1 := make(chan error, 1)
    err2 := make(chan error, 1)

    go func() {
        err1 <- s.env.ShadowUC.UpdateDesired(s.ctx, UpdateShadowDesiredInput{
            TenantID: tenantID.String(), DeviceID: deviceID.String(),
            Desired:  map[string]any{"pump": "on"},
        })
    }()
    go func() {
        err2 <- s.env.ShadowUC.UpdateDesired(s.ctx, UpdateShadowDesiredInput{
            TenantID: tenantID.String(), DeviceID: deviceID.String(),
            Desired:  map[string]any{"fan": "on"},
        })
    }()

    e1 := <-err1
    e2 := <-err2

    // ต้องมีอย่างน้อย 1 succeed (หรือ retry)
    s.True(e1 == nil || e2 == nil, "expected at least one success")
}
```

### D.2 E2E Smoke Script

```bash
#!/usr/bin/env bash
# test/e2e/smoke.sh — end-to-end smoke test

set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
TENANT="${TENANT_ID:-11111111-1111-1111-1111-111111111101}"
EMAIL="${EMAIL:-admin@demo.local}"
PASS="${PASS:-Password123!}"

echo "→ E2E Smoke Test"

# 1. Login
echo "  [1/12] Login..."
TOKEN=$(curl -sf -X POST "${BASE_URL}/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${EMAIL}\",\"password\":\"${PASS}\",\"tenant_id\":\"${TENANT}\"}" \
  | jq -r '.data.access_token')

if [[ -z "$TOKEN" || "$TOKEN" == "null" ]]; then
  echo "  ✗ Login failed"
  exit 1
fi
echo "  ✓ Token received"

AUTH_HEADER="Authorization: Bearer ${TOKEN}"
TENANT_HEADER="X-Tenant-ID: ${TENANT}"

# 2. Health check
echo "  [2/12] Health check..."
curl -sf "${BASE_URL}/health" | jq -e '.status=="ok"' >/dev/null
echo "  ✓ Health OK"

# 3. List packages
echo "  [3/12] List packages..."
PKG_COUNT=$(curl -sf -H "${AUTH_HEADER}" -H "${TENANT_HEADER}" \
  "${BASE_URL}/api/v1/packages" | jq '.data.items | length')
echo "  ✓ ${PKG_COUNT} packages"

# 4. List customers
echo "  [4/12] List customers..."
CUST_COUNT=$(curl -sf -H "${AUTH_HEADER}" -H "${TENANT_HEADER}" \
  "${BASE_URL}/api/v1/customers" | jq '.data.items | length')
echo "  ✓ ${CUST_COUNT} customers"

# 5. Create customer
echo "  [5/12] Create customer..."
NEW_CUST=$(curl -sf -X POST "${BASE_URL}/api/v1/customers" \
  -H "${AUTH_HEADER}" -H "${TENANT_HEADER}" \
  -H "Content-Type: application/json" \
  -d '{
    "type":"CORPORATE","name":"Smoke Test Co.",
    "tax_id":"0105559999999","email":"smoke@test.local",
    "address":{"line1":"1 Test","country":"TH"}
  }')
CUST_ID=$(echo "$NEW_CUST" | jq -r '.data.id')
echo "  ✓ Created customer ${CUST_ID}"

# 6. List devices
echo "  [6/12] List devices..."
DEV_COUNT=$(curl -sf -H "${AUTH_HEADER}" -H "${TENANT_HEADER}" \
  "${BASE_URL}/api/v1/devices" | jq '.data.items | length')
echo "  ✓ ${DEV_COUNT} devices"

# 7. Register device
echo "  [7/12] Register device..."
DEV_RESP=$(curl -sf -X POST "${BASE_URL}/api/v1/devices" \
  -H "${AUTH_HEADER}" -H "${TENANT_HEADER}" \
  -H "Content-Type: application/json" \
  -d "{
    \"customer_id\":\"${CUST_ID}\",
    \"site_id\":\"55555555-5555-5555-5555-555555555501\",
    \"serial_no\":\"SN-SMOKE-$(date +%s)\",
    \"name\":\"Smoke Test Device\",
    \"type\":\"SENSOR\",\"protocol\":\"MQTT\"
  }")
DEV_ID=$(echo "$DEV_RESP" | jq -r '.data.id')
echo "  ✓ Registered device ${DEV_ID}"

# 8. Provision
echo "  [8/12] Provision device..."
curl -sf -X POST "${BASE_URL}/api/v1/devices/${DEV_ID}/provision" \
  -H "${AUTH_HEADER}" -H "${TENANT_HEADER}" | jq -e '.data.device_token' >/dev/null
echo "  ✓ Device provisioned"

# 9. Ingest telemetry
echo "  [9/12] Ingest telemetry..."
curl -sf -X POST "${BASE_URL}/api/v1/devices/${DEV_ID}/telemetry" \
  -H "${AUTH_HEADER}" -H "${TENANT_HEADER}" \
  -H "Content-Type: application/json" \
  -d '{"metrics":[{"metric":"temperature","value":28.5},{"metric":"humidity","value":65}]}' >/dev/null
echo "  ✓ Telemetry ingested"

# 10. Query telemetry
echo "  [10/12] Query telemetry..."
sleep 2
TEL_COUNT=$(curl -sf -H "${AUTH_HEADER}" -H "${TENANT_HEADER}" \
  "${BASE_URL}/api/v1/devices/${DEV_ID}/telemetry?from=$(date -u -d '-1 hour' +%Y-%m-%dT%H:%M:%SZ)&metric=temperature" \
  | jq '.data.series.temperature | length')
echo "  ✓ ${TEL_COUNT} telemetry points"

# 11. KPIs
echo "  [11/12] Fetch KPIs..."
curl -sf -H "${AUTH_HEADER}" -H "${TENANT_HEADER}" \
  "${BASE_URL}/api/v1/dashboard/kpis?code=MRR&code=ACTIVE_CUSTOMERS&preset=LAST_30_DAYS" \
  | jq -e '.data.kpis | length > 0' >/dev/null
echo "  ✓ KPIs returned"

# 12. Generate report
echo "  [12/12] Generate report..."
REPORT_RESP=$(curl -sf -X POST "${BASE_URL}/api/v1/reports/RPT-CUSTOMER-LIST/generate" \
  -H "${AUTH_HEADER}" -H "${TENANT_HEADER}" \
  -H "Content-Type: application/json" \
  -d '{"format":"JSON","preset":"LAST_30_DAYS"}')
REPORT_STATUS=$(echo "$REPORT_RESP" | jq -r '.data.status')
echo "  ✓ Report status: ${REPORT_STATUS}"

echo ""
echo "✓ E2E smoke test PASSED"
```

---

## 🅴 PART 9E — DOCKER / CI INTEGRATION

### E.1 `docker-compose.seed.yml`

```yaml
version: "3.9"

services:
  seed:
    build:
      context: .
      dockerfile: Dockerfile
      args: { TARGET: migrate }
    image: icmongolang/seed:latest
    container_name: icmongolang-seed
    depends_on:
      postgres:    { condition: service_healthy }
      kafka:       { condition: service_healthy }
      influxdb:    { condition: service_healthy }
      elasticsearch: { condition: service_started }
    environment:
      DB_DSN: "host=postgres user=icmon password=icmon_dev dbname=icmongolang port=5432 sslmode=disable TimeZone=Asia/Bangkok"
      KAFKA_BROKERS: "kafka:9092"
      INFLUX_URL: "http://influxdb:8086"
      INFLUX_TOKEN: "dev-token-change-me"
      INFLUX_ORG: "icmon"
      INFLUX_BUCKET: "iot_telemetry"
      SEED_LEVEL: "standard"
    volumes:
      - ./test/fixtures:/fixtures:ro
      - ./test/fixtures/sql:/sql:ro
    entrypoint: ["/bin/sh", "-c"]
    command:
      - |
        set -e
        echo "→ Running fixtures..."
        for f in /sql/*.sql; do
          echo "  → $$f"
          psql "$$DB_DSN" -f "$$f" || true
        done
        echo "→ Seeding telemetry..."
        apk add --no-cache curl bash jq >/dev/null 2>&1 || true
        bash /sql/16_seed_telemetry.sh
        echo "✓ Seed complete"
    networks: [backend]
```

### E.2 `Makefile` targets

```makefile
# ─── SEED ───────────────────────────────────────────────────
.PHONY: seed
seed: ## seed standard fixtures
	@echo "→ Seeding standard fixtures..."
	@for f in test/fixtures/sql/*.sql; do \
		echo "  → $$f"; \
		docker compose exec -T postgres psql -U icmon -d icmongolang < "$$f" || true; \
	done
	@bash test/fixtures/sql/16_seed_telemetry.sh
	@echo "✓ Seed complete"

.PHONY: seed-minimal
seed-minimal: ## seed only tenants + users + packages
	@for f in test/fixtures/sql/00_cleanup.sql test/fixtures/sql/01_tenants.sql \
	         test/fixtures/sql/02_users_auth.sql test/fixtures/sql/05_packages.sql; do \
		docker compose exec -T postgres psql -U icmon -d icmongolang < "$$f" || true; \
	done

.PHONY: seed-demo
seed-demo: seed ## alias — full demo data
	@echo "✓ Demo data ready"
	@echo "   Login: admin@demo.local / Password123!"
	@echo "   Tenant: 11111111-1111-1111-1111-111111111101"

.PHONY: seed-stress
seed-stress: ## seed + 10k devices (load test)
	@bash scripts/seed-stress.sh

.PHONY: seed-reset
seed-reset: ## cleanup + reseed
	@docker compose exec -T postgres psql -U icmon -d icmongolang < test/fixtures/sql/00_cleanup.sql
	@$(MAKE) seed

.PHONY: fixtures-clean
fixtures-clean:
	@docker compose exec -T postgres psql -U icmon -d icmongolang < test/fixtures/sql/00_cleanup.sql
	@echo "✓ Fixtures cleaned"

# ─── TEST ───────────────────────────────────────────────────
.PHONY: test-integration
test-integration: ## integration tests (ต้องมี infra)
	@DB_DSN="host=localhost user=icmon password=icmon_dev dbname=icmongolang_test sslmode=disable" \
	 go test -tags=integration -timeout=15m ./test/integration/...

.PHONY: test-smoke
test-smoke: ## E2E smoke test
	@bash test/e2e/smoke.sh

.PHONY: test-scenarios
test-scenarios: ## specific scenario tests
	@go test -tags=integration -run TestScenarioSuite -v ./test/integration/...
```

### E.3 GitHub Actions Workflow

```yaml
# .github/workflows/seed-and-test.yml
name: Seed & Test

on:
  pull_request:
    branches: [main, develop]

jobs:
  integration:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_USER: icmon
          POSTGRES_PASSWORD: icmon_dev
          POSTGRES_DB: icmongolang_test
        ports: ['5432:5432']
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

      redis:
        image: redis:7-alpine
        ports: ['6379:6379']

      kafka:
        image: bitnami/kafka:3.7
        env:
          KAFKA_CFG_NODE_ID: 0
          KAFKA_CFG_PROCESS_ROLES: controller,broker
          KAFKA_CFG_CONTROLLER_QUORUM_VOTERS: 0@kafka:9093
          KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,CONTROLLER://:9093
          KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
          KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP: CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
          KAFKA_CFG_CONTROLLER_LISTENER_NAMES: CONTROLLER
          KAFKA_CFG_INTER_BROKER_LISTENER_NAME: PLAINTEXT
        ports: ['9092:9092']

      influxdb:
        image: influxdb:2.7-alpine
        env:
          DOCKER_INFLUXDB_INIT_MODE: setup
          DOCKER_INFLUXDB_INIT_USERNAME: admin
          DOCKER_INFLUXDB_INIT_PASSWORD: adminpassword
          DOCKER_INFLUXDB_INIT_ORG: icmon
          DOCKER_INFLUXDB_INIT_BUCKET: iot_telemetry
          DOCKER_INFLUXDB_INIT_ADMIN_TOKEN: test-token
        ports: ['8086:8086']

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.23' }

      - name: Install psql
        run: sudo apt-get install -y postgresql-client

      - name: Run migrations
        env:
          DB_DSN: 'host=localhost user=icmon password=icmon_dev dbname=icmongolang_test sslmode=disable'
        run: go run ./cmd/migrate -action up

      - name: Load fixtures
        env:
          DB_DSN: 'host=localhost user=icmon password=icmon_dev dbname=icmongolang_test sslmode=disable'
        run: |
          for f in test/fixtures/sql/*.sql; do
            echo "→ $f"
            psql "$DB_DSN" -f "$f" || true
          done

      - name: Run integration tests
        env:
          DB_DSN: 'host=localhost user=icmon password=icmon_dev dbname=icmongolang_test sslmode=disable'
          REDIS_ADDR: 'localhost:6379'
          KAFKA_BROKERS: 'localhost:9092'
          INFLUX_URL: 'http://localhost:8086'
          INFLUX_TOKEN: 'test-token'
        run: go test -tags=integration -timeout=15m ./test/integration/...
```

---

## 🅵 PART 9F — POSTMAN ENVIRONMENT

### F.1 `test/fixtures/postman/environment.dev.json`

```json
{
  "id": "icmongolang-dev-env",
  "name": "icmongolang — Dev",
  "values": [
    { "key": "base_url", "value": "http://localhost:8080", "enabled": true },
    { "key": "api_version", "value": "v1", "enabled": true },
    { "key": "tenant_id", "value": "11111111-1111-1111-1111-111111111101", "enabled": true },
    { "key": "tenant_slug", "value": "demo", "enabled": true },
    { "key": "admin_email", "value": "admin@demo.local", "enabled": true },
    { "key": "admin_password", "value": "Password123!", "enabled": true },
    { "key": "manager_email", "value": "manager@demo.local", "enabled": true },
    { "key": "viewer_email", "value": "viewer@demo.local", "enabled": true },
    { "key": "access_token", "value": "", "enabled": true },
    { "key": "refresh_token", "value": "", "enabled": true },
    { "key": "user_id", "value": "", "enabled": true },

    { "key": "customer_id", "value": "33333333-3333-3333-3333-333333333301", "enabled": true },
    { "key": "customer_code", "value": "CUS-2026-0001", "enabled": true },
    { "key": "site_id", "value": "55555555-5555-5555-5555-555555555501", "enabled": true },
    { "key": "zone_id", "value": "66666666-6666-6666-6666-666666666601", "enabled": true },

    { "key": "package_id", "value": "77777777-7777-7777-7777-777777777704", "enabled": true },
    { "key": "package_code", "value": "PRO-M", "enabled": true },
    { "key": "subscription_id", "value": "88888888-8888-8888-8888-888888888801", "enabled": true },

    { "key": "device_id", "value": "f0000000-0000-0000-0000-000000000001", "enabled": true },
    { "key": "device_serial", "value": "SN-DEMO-001", "enabled": true },
    { "key": "device_token", "value": "", "enabled": true },

    { "key": "product_id", "value": "c0000000-0000-0000-0000-000000000001", "enabled": true },
    { "key": "sku", "value": "SENSOR-TEMP-001", "enabled": true },
    { "key": "warehouse_id", "value": "a0000000-0000-0000-0000-000000000001", "enabled": true },
    { "key": "order_id", "value": "ee000000-0000-0000-0000-000000000002", "enabled": true },
    { "key": "order_no", "value": "SO-2026-000002", "enabled": true },
    { "key": "invoice_id", "value": "aa110000-0000-0000-0000-000000000002", "enabled": true },
    { "key": "invoice_no", "value": "INV-2026-000002", "enabled": true },
    { "key": "payment_id", "value": "cc110000-0000-0000-0000-000000000001", "enabled": true },

    { "key": "shipment_id", "value": "ff110000-0000-0000-0000-000000000002", "enabled": true },
    { "key": "job_id", "value": "ab220000-0000-0000-0000-000000000001", "enabled": true },
    { "key": "technician_id", "value": "ee110000-0000-0000-0000-000000000001", "enabled": true },

    { "key": "report_code", "value": "RPT-CUSTOMER-LIST", "enabled": true },
    { "key": "dashboard_id", "value": "cd110000-0000-0000-0000-000000000001", "enabled": true }
  ],
  "_postman_variable_scope": "environment"
}
```

### F.2 Pre-request Script (Collection Level)

```javascript
// ============================================================
// Collection-level pre-request script
// รันก่อนทุก request
// ============================================================

const baseUrl = pm.environment.get('base_url') || 'http://localhost:8080';

// 1. Set dynamic headers
pm.request.headers.upsert({
    key: 'X-Request-ID',
    value: pm.variables.replaceIn('{{$guid}}')
});
pm.request.headers.upsert({
    key: 'X-Tenant-ID',
    value: pm.environment.get('tenant_id')
});

// 2. Auto-attach JWT ถ้ามี
const token = pm.environment.get('access_token');
if (token) {
    pm.request.headers.upsert({
        key: 'Authorization',
        value: 'Bearer ' + token
    });
}

// 3. Auto-refresh token ถ้าใกล้หมดอายุ (มี expires_at)
const expiresAt = pm.environment.get('token_expires_at');
if (expiresAt && Date.now() > new Date(expiresAt).getTime() - 60000) {
    console.log('Token expiring soon — refreshing...');
    pm.sendRequest({
        url: baseUrl + '/api/v1/auth/refresh',
        method: 'POST',
        header: { 'Content-Type': 'application/json' },
        body: {
            mode: 'raw',
            raw: JSON.stringify({
                refresh_token: pm.environment.get('refresh_token')
            })
        }
    }, (err, res) => {
        if (!err && res.code === 200) {
            const data = res.json().data;
            pm.environment.set('access_token', data.access_token);
            pm.environment.set('refresh_token', data.refresh_token);
            pm.environment.set('token_expires_at', data.expires_at);
        }
    });
}
```

### F.3 Login Request + Test Script

```javascript
// ============================================================
// POST /api/v1/auth/login — Test Script
// ============================================================

pm.test('Status 200', () => pm.response.to.have.status(200));

pm.test('Response has access_token', () => {
    const data = pm.response.json().data;
    pm.expect(data).to.have.property('access_token');
    pm.expect(data.access_token).to.be.a('string').and.not.empty;
});

// Auto-save tokens to environment
const data = pm.response.json().data;
if (data && data.access_token) {
    pm.environment.set('access_token', data.access_token);
    pm.environment.set('refresh_token', data.refresh_token || '');
    pm.environment.set('user_id', data.user_id || '');
    pm.environment.set('token_expires_at', data.expires_at || '');
    console.log('✓ Token saved');
}
```

### F.4 Create Order + Chained Variables

```javascript
// ============================================================
// POST /api/v1/erp/orders — Test Script
// ============================================================

pm.test('Order created', () => {
    pm.response.to.have.status(201);
    const data = pm.response.json().data;
    pm.expect(data).to.have.property('id');
    pm.expect(data).to.have.property('order_no');
    pm.expect(data.status).to.equal('DRAFT');

    // Save for next requests
    pm.environment.set('order_id', data.id);
    pm.environment.set('order_no', data.order_no);

    // Save line IDs for chaining
    if (data.lines && data.lines.length > 0) {
        pm.environment.set('order_line_id', data.lines[0].id);
    }

    console.log('✓ Order saved: ' + data.order_no);
});
```

### F.5 Postman Collection Structure

```json
{
  "info": {
    "name": "icmongolang IoT Platform",
    "description": "Full API collection — 6 modules",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "auth": { "type": "bearer", "bearer": [{ "key": "token", "value": "{{access_token}}" }] },
  "event": [
    { "listen": "prerequest", "script": { "exec": ["// Collection pre-request script"] } }
  ],
  "item": [
    {
      "name": "01. Auth",
      "item": [
        { "name": "Login", "request": { "method": "POST", "url": "{{base_url}}/api/v1/auth/login" } },
        { "name": "Refresh", "request": { "method": "POST", "url": "{{base_url}}/api/v1/auth/refresh" } },
        { "name": "Logout", "request": { "method": "POST", "url": "{{base_url}}/api/v1/auth/logout" } },
        { "name": "Me", "request": { "method": "GET", "url": "{{base_url}}/api/v1/auth/me" } }
      ]
    },
    {
      "name": "02. Customers",
      "item": [
        { "name": "List", "request": { "method": "GET", "url": "{{base_url}}/api/v1/customers" } },
        { "name": "Create", "request": { "method": "POST", "url": "{{base_url}}/api/v1/customers" } },
        { "name": "Get", "request": { "method": "GET", "url": "{{base_url}}/api/v1/customers/{{customer_id}}" } },
        { "name": "Update", "request": { "method": "PUT", "url": "{{base_url}}/api/v1/customers/{{customer_id}}" } },
        { "name": "Onboard", "request": { "method": "POST", "url": "{{base_url}}/api/v1/customers/{{customer_id}}/onboard" } },
        { "name": "Add Contact", "request": { "method": "POST", "url": "{{base_url}}/api/v1/customers/{{customer_id}}/contacts" } }
      ]
    },
    {
      "name": "03. Packages",
      "item": [
        { "name": "List", "request": { "method": "GET", "url": "{{base_url}}/api/v1/packages" } },
        { "name": "Subscribe", "request": { "method": "POST", "url": "{{base_url}}/api/v1/subscriptions" } },
        { "name": "My Subscription", "request": { "method": "GET", "url": "{{base_url}}/api/v1/subscriptions/me" } },
        { "name": "Upgrade", "request": { "method": "POST", "url": "{{base_url}}/api/v1/subscriptions/{{subscription_id}}/upgrade" } },
        { "name": "Cancel", "request": { "method": "POST", "url": "{{base_url}}/api/v1/subscriptions/{{subscription_id}}/cancel" } },
        { "name": "Check Quota", "request": { "method": "POST", "url": "{{base_url}}/api/v1/subscriptions/{{subscription_id}}/quota/check" } }
      ]
    },
    {
      "name": "04. Devices",
      "item": [
        { "name": "List", "request": { "method": "GET", "url": "{{base_url}}/api/v1/devices" } },
        { "name": "Register", "request": { "method": "POST", "url": "{{base_url}}/api/v1/devices" } },
        { "name": "Provision", "request": { "method": "POST", "url": "{{base_url}}/api/v1/devices/{{device_id}}/provision" } },
        { "name": "Ingest Telemetry", "request": { "method": "POST", "url": "{{base_url}}/api/v1/devices/{{device_id}}/telemetry" } },
        { "name": "Query Telemetry", "request": { "method": "GET", "url": "{{base_url}}/api/v1/devices/{{device_id}}/telemetry" } },
        { "name": "Get Shadow", "request": { "method": "GET", "url": "{{base_url}}/api/v1/devices/{{device_id}}/shadow" } },
        { "name": "Update Desired", "request": { "method": "PATCH", "url": "{{base_url}}/api/v1/devices/{{device_id}}/shadow/desired" } },
        { "name": "Send Command", "request": { "method": "POST", "url": "{{base_url}}/api/v1/devices/{{device_id}}/commands" } }
      ]
    },
    {
      "name": "05. ERP",
      "item": [
        { "name": "List Products", "request": { "method": "GET", "url": "{{base_url}}/api/v1/erp/products" } },
        { "name": "Create Product", "request": { "method": "POST", "url": "{{base_url}}/api/v1/erp/products" } },
        { "name": "List Inventory", "request": { "method": "GET", "url": "{{base_url}}/api/v1/erp/inventory" } },
        { "name": "Adjust Stock", "request": { "method": "POST", "url": "{{base_url}}/api/v1/erp/inventory/adjust" } },
        { "name": "Transfer Stock", "request": { "method": "POST", "url": "{{base_url}}/api/v1/erp/inventory/transfer" } },
        { "name": "Create Order", "request": { "method": "POST", "url": "{{base_url}}/api/v1/erp/orders" } },
        { "name": "Confirm Order", "request": { "method": "POST", "url": "{{base_url}}/api/v1/erp/orders/{{order_id}}/confirm" } },
        { "name": "Ship Order", "request": { "method": "POST", "url": "{{base_url}}/api/v1/erp/orders/{{order_id}}/ship" } },
        { "name": "Receive Order", "request": { "method": "POST", "url": "{{base_url}}/api/v1/erp/orders/{{order_id}}/receive" } },
        { "name": "Create Invoice", "request": { "method": "POST", "url": "{{base_url}}/api/v1/erp/invoices" } },
        { "name": "Aging Report", "request": { "method": "GET", "url": "{{base_url}}/api/v1/erp/invoices/aging" } },
        { "name": "Record Payment", "request": { "method": "POST", "url": "{{base_url}}/api/v1/erp/payments" } }
      ]
    },
    {
      "name": "06. Logistics",
      "item": [
        { "name": "List Shipments", "request": { "method": "GET", "url": "{{base_url}}/api/v1/shipments" } },
        { "name": "Create Shipment", "request": { "method": "POST", "url": "{{base_url}}/api/v1/shipments" } },
        { "name": "Dispatch", "request": { "method": "POST", "url": "{{base_url}}/api/v1/shipments/{{shipment_id}}/dispatch" } },
        { "name": "Update GPS", "request": { "method": "POST", "url": "{{base_url}}/api/v1/shipments/{{shipment_id}}/gps" } },
        { "name": "Confirm Delivery", "request": { "method": "POST", "url": "{{base_url}}/api/v1/shipments/{{shipment_id}}/deliver" } },
        { "name": "List Technicians", "request": { "method": "GET", "url": "{{base_url}}/api/v1/technicians" } },
        { "name": "Assign Technician", "request": { "method": "POST", "url": "{{base_url}}/api/v1/installation-jobs/{{job_id}}/assign" } },
        { "name": "Start Job", "request": { "method": "POST", "url": "{{base_url}}/api/v1/installation-jobs/{{job_id}}/start" } },
        { "name": "Update Checklist", "request": { "method": "PATCH", "url": "{{base_url}}/api/v1/installation-jobs/{{job_id}}/checklist" } },
        { "name": "Complete Job", "request": { "method": "POST", "url": "{{base_url}}/api/v1/installation-jobs/{{job_id}}/complete" } },
        { "name": "Optimize Route", "request": { "method": "POST", "url": "{{base_url}}/api/v1/technicians/{{technician_id}}/optimize-route" } }
      ]
    },
    {
      "name": "07. Reports",
      "item": [
        { "name": "List Reports", "request": { "method": "GET", "url": "{{base_url}}/api/v1/reports" } },
        { "name": "Generate Report", "request": { "method": "POST", "url": "{{base_url}}/api/v1/reports/{{report_code}}/generate" } },
        { "name": "List Executions", "request": { "method": "GET", "url": "{{base_url}}/api/v1/report-executions" } },
        { "name": "Download", "request": { "method": "GET", "url": "{{base_url}}/api/v1/report-executions/{{execution_id}}/download" } },
        { "name": "Create Schedule", "request": { "method": "POST", "url": "{{base_url}}/api/v1/report-schedules" } },
        { "name": "Get KPIs", "request": { "method": "GET", "url": "{{base_url}}/api/v1/dashboard/kpis" } },
        { "name": "KPI History", "request": { "method": "GET", "url": "{{base_url}}/api/v1/dashboard/kpis/history" } },
        { "name": "Dashboard Data", "request": { "method": "GET", "url": "{{base_url}}/api/v1/dashboard/dashboards/{{dashboard_id}}/data" } },
        { "name": "Generate Insight", "request": { "method": "POST", "url": "{{base_url}}/api/v1/dashboard/insights/generate" } }
      ]
    }
  ]
}
```

---

## 📊 PART 9 SUMMARY

### Fixtures Created

| Category | Count | Format |
|---|:-:|---|
| **SQL files** | 17 | SQL |
| **Tenants** | 3 | demo / acme / farm1 |
| **Users** | 7 | admin / manager / viewer × tenants |
| **Customers** | 10 | CORP / INDIVIDUAL / GOV × lifecycle |
| **Sites + Zones** | 8 + 5 | |
| **Packages** | 6 | FREE/BASIC/PRO/ENT + deprecated |
| **Subscriptions** | 8 | TRIAL/ACTIVE/PAST_DUE × tenants |
| **Products** | 15+ | Sensor/Gateway/Actuator/Service |
| **Warehouses** | 4 | Main/Transit × tenants |
| **Inventory Items** | 13 | healthy / low / out-of-stock |
| **Devices** | 12+ | ONLINE/OFFLINE/FAULT/GATEWAY/ACTUATOR |
| **Orders** | 8 | SO/PO × statuses |
| **Invoices** | 5 | AR/AP × DRAFT/ISSUED/PAID/PARTIAL/OVERDUE |
| **Payments** | 3 | IN/OUT with allocations |
| **Technicians** | 4 | Junior/Senior/Lead |
| **Shipments** | 3 | PENDING/IN_TRANSIT/INSTALLED |
| **Installation Jobs** | 6 | SCHEDULED/ASSIGNED/IN_PROGRESS/PAUSED/DONE/FAILED |
| **KPI Snapshots** | 180 | 6 KPIs × 30 days |
| **Dashboards** | 3 | Overview/IoT/Exec |
| **Alert Rules** | 3 | per tenant |
| **Telemetry** | ~21,600 | 30 days × 3 devices × 2 metrics |

### Benefits

1. ✅ **Deterministic UUIDs** — รันซ้ำได้ ผลเหมือนเดิม
2. ✅ **Idempotent SQL** — `ON CONFLICT DO NOTHING` / `UPDATE` 
3. ✅ **Multi-tenant** — demo / acme / farm1 แยกข้อมูลชัดเจน
4. ✅ **Realistic scenarios** — Order-to-cash, subscribe-to-activate, device-provision
5. ✅ **Edge cases** — deprecated package, overdue invoice, out-of-stock, fault device
6. ✅ **Testable** — Go builders + SQL + Postman + smoke script
7. ✅ **CI-ready** — GitHub Actions workflow
8. ✅ **Time-relative** — `NOW() - INTERVAL` ปรับตัวตามเวลา

---

# 🎯 ความคืบหน้า (ลำดับ B)

| # | งาน | สถานะ |
|:-:|---|:-:|
| 1 | `packagecatalog` deep dive | ✅ |
| 2 | `erp` deep dive | ✅ |
| 3 | UML / Sequence Diagrams | ✅ |
| 4 | Sample Data / Fixtures | ✅ **เสร็จ (response นี้)** |
| 5 | Postman Collection | ⏳ ถัดไป |
| 6 | Executive Summary | ⏳ |

---

# 📋 Response ถัดไป: **PART 10 — Postman Collection & API Testing**

จะมี:
- **Full Postman collection** (JSON พร้อม import)
- **Environments** (dev/staging/prod)
- **Auth flow** — login → refresh → logout
- **Chained requests** — ทุก request ต่อกันได้อัตโนมัติ
- **Pre-request scripts** — auto-inject tenant + token + request ID
- **Test scripts** — assertions ครบทุก endpoint
- **Newman CLI** — รันใน CI
- **Collection Runner** — รันทั้งชุด
- **Mock server** — Postman mock สำหรับ frontend dev
- **API documentation** — auto-generated จาก collection
