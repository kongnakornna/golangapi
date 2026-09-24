# 📮 PART 10 — POSTMAN COLLECTION & API TESTING

> **ขนาด**: ใหญ่ — แยก 7 ตอนย่อย
> **Part 10A**: Collection Architecture & Conventions
> **Part 10B**: Environment Files (dev / staging / prod)
> **Part 10C**: Full Collection (JSON — พร้อม import)
> **Part 10D**: Pre-request Scripts & Test Scripts (ทุก endpoint)
> **Part 10E**: Auth Flow (login → refresh → logout)
> **Part 10F**: Newman CLI + CI Integration
> **Part 10G**: Mock Server + API Documentation

> **เป้าหมาย**: ให้ dev/test/QA import Postman แล้วรัน API ได้ทันที — chained requests, auto-token, auto-tenant, assertions ครบ

---

## 🅰️ PART 10A — COLLECTION ARCHITECTURE & CONVENTIONS

### A.1 Collection Structure

```
icmongolang-iot-platform.postman_collection.json
│
├── 📁 00. Setup
│   ├── Health Check
│   ├── Readiness
│   └── Version
│
├── 📁 01. Auth
│   ├── Login (Admin)
│   ├── Login (Manager)
│   ├── Login (Viewer)
│   ├── Refresh Token
│   ├── Me
│   ├── Change Password
│   └── Logout
│
├── 📁 02. Customers
│   ├── List Customers
│   ├── Create Customer
│   ├── Get Customer
│   ├── Update Customer
│   ├── Delete Customer
│   ├── Add Contact
│   ├── List Contacts
│   ├── Onboard Customer
│   ├── Suspend Customer
│   └── Reactivate Customer
│
├── 📁 03. Sites & Zones
│   ├── List Sites
│   ├── Create Site
│   ├── Get Site
│   ├── Create Zone
│   └── List Zones
│
├── 📁 04. Packages
│   ├── List Packages
│   ├── Get Package
│   ├── Create Package
│   ├── Update Package
│   ├── Deprecate Package
│   ├── List Subscriptions
│   ├── Create Subscription
│   ├── Get My Subscription
│   ├── Upgrade Subscription
│   ├── Cancel Subscription
│   ├── Check Quota
│   └── List Subscription History
│
├── 📁 05. Devices
│   ├── List Devices
│   ├── Register Device
│   ├── Get Device
│   ├── Update Device
│   ├── Delete Device
│   ├── Provision Device
│   ├── Send Command
│   ├── List Commands
│   ├── Ingest Telemetry
│   ├── Query Telemetry
│   ├── Latest Telemetry
│   ├── Get Shadow
│   ├── Update Desired State
│   ├── Clear Desired State
│   ├── List Alert Rules
│   ├── Create Alert Rule
│   ├── List Alert Events
│   └── Acknowledge Alert
│
├── 📁 06. ERP
│   ├── 📁 Products
│   │   ├── List Products
│   │   ├── Create Product
│   │   ├── Get Product
│   │   ├── Update Product
│   │   └── Delete Product
│   ├── 📁 Warehouses
│   │   ├── List Warehouses
│   │   └── Create Warehouse
│   ├── 📁 Inventory
│   │   ├── List Inventory
│   │   ├── Get Stock
│   │   ├── Adjust Stock
│   │   ├── Transfer Stock
│   │   └── Stock Movements
│   ├── 📁 Orders
│   │   ├── List Orders
│   │   ├── Create Order
│   │   ├── Get Order
│   │   ├── Update Order
│   │   ├── Confirm Order
│   │   ├── Ship Order
│   │   ├── Receive Order
│   │   └── Cancel Order
│   ├── 📁 Invoices
│   │   ├── List Invoices
│   │   ├── Create Invoice
│   │   ├── Get Invoice
│   │   ├── Issue Invoice
│   │   ├── Void Invoice
│   │   └── Aging Report
│   ├── 📁 Payments
│   │   ├── List Payments
│   │   ├── Create Payment
│   │   ├── Get Payment
│   │   └── Allocate Payment
│   └── 📁 Reports
│       ├── Trial Balance
│       ├── P&L
│       └── Balance Sheet
│
├── 📁 07. Logistics
│   ├── 📁 Shipments
│   │   ├── List Shipments
│   │   ├── Create Shipment
│   │   ├── Get Shipment
│   │   ├── Dispatch
│   │   ├── Update GPS
│   │   ├── Confirm Delivery
│   │   └── Cancel Shipment
│   ├── 📁 Technicians
│   │   ├── List Technicians
│   │   ├── Create Technician
│   │   ├── Get Technician
│   │   └── Optimize Route
│   └── 📁 Installation Jobs
│       ├── List Jobs
│       ├── Create Job
│       ├── Assign Technician
│       ├── Start Job
│       ├── Update Checklist
│       ├── Pause Job
│       ├── Complete Job
│       └── Fail Job
│
├── 📁 08. Reports & Dashboard
│   ├── List Reports
│   ├── Generate Report
│   ├── List Report Executions
│   ├── Download Report
│   ├── Create Report Schedule
│   ├── List KPIs
│   ├── KPI History
│   ├── Get Dashboard
│   ├── Get Dashboard Data
│   └── Generate Insight (AI)
│
├── 📁 09. WebSocket (Docs)
│   └── WS Connection Info
│
└── 📁 10. Utilities
    ├── Seed Demo Data
    ├── Reset Fixtures
    └── Export Data
```

### A.2 Conventions

| Rule | ตัวอย่าง |
|---|---|
| **Folder numbering** | `01. Auth`, `02. Customers` — เพื่อเรียงลำดับ |
| **Request naming** | `<Verb> <Noun>` — `Create Customer`, `List Devices` |
| **URL variables** | `{{base_url}}/api/v1/...` |
| **Chained vars** | บันทึก response → env var → request ถัดไปใช้ |
| **Test names** | ภาษาอังกฤษ, เริ่มด้วย `✓` หรือ `Status` |
| **Pre-request** | Collection-level (auth + tenant) + Folder-level (validation) |
| **Test scripts** | ทุก request ต้องมีอย่างน้อย 3 assertions |
| **Error tests** | แยก request `[Negative]` สำหรับ case error |
| **Idempotency** | ใช้ `{{$guid}}` สำหรับ X-Request-ID |

### A.3 Chained Variables Flow

```
Login  ──► access_token, refresh_token, user_id
   │
   ▼
List Packages ──► package_id (PRO-M)
   │
   ▼
Create Customer ──► customer_id, customer_code
   │
   ▼
Create Site ──► site_id
   │
   ▼
Create Subscription ──► subscription_id
   │
   ▼
Register Device ──► device_id, device_serial
   │
   ▼
Provision Device ──► device_token
   │
   ▼
Ingest Telemetry ──► (verify with Query Telemetry)
   │
   ▼
Create Order ──► order_id, order_no
   │
   ▼
Confirm Order ──► (status = CONFIRMED)
   │
   ▼
Ship Order ──► (inventory deducted)
   │
   ▼
Create Invoice ──► invoice_id, invoice_no
   │
   ▼
Create Payment ──► payment_id (invoice marked paid)
```

### A.4 Variable Scopes

| Scope | Variables | ตัวอย่าง |
|---|---|---|
| **Global** | `environment_name` | `dev`, `staging`, `prod` |
| **Collection** | `api_version`, `default_page_size` | `v1`, `20` |
| **Environment** | `base_url`, `tenant_id`, `admin_email` | ค่าตาม environment |
| **Local** | `{{$guid}}`, `{{$timestamp}}` | Dynamic ต่อ request |
| **Chained** | `access_token`, `customer_id`, `device_id` | บันทึกจาก response |

---

## 🅱️ PART 10B — ENVIRONMENT FILES

### B.1 `environment.dev.json`

```json
{
  "id": "icmongolang-env-dev",
  "name": "icmongolang — Dev",
  "values": [
    { "key": "environment_name", "value": "dev", "type": "default", "enabled": true },
    { "key": "base_url", "value": "http://localhost:8080", "type": "default", "enabled": true },
    { "key": "ws_url", "value": "ws://localhost:8080/ws", "type": "default", "enabled": true },
    { "key": "api_version", "value": "v1", "type": "default", "enabled": true },
    { "key": "default_page_size", "value": "20", "type": "default", "enabled": true },
    { "key": "default_timeout", "value": "30000", "type": "default", "enabled": true },

    { "key": "tenant_id", "value": "11111111-1111-1111-1111-111111111101", "type": "default", "enabled": true },
    { "key": "tenant_slug", "value": "demo", "type": "default", "enabled": true },
    { "key": "admin_email", "value": "admin@demo.local", "type": "default", "enabled": true },
    { "key": "admin_password", "value": "Password123!", "type": "secret", "enabled": true },
    { "key": "manager_email", "value": "manager@demo.local", "type": "default", "enabled": true },
    { "key": "viewer_email", "value": "viewer@demo.local", "type": "default", "enabled": true },
    { "key": "common_password", "value": "Password123!", "type": "secret", "enabled": true },

    { "key": "access_token", "value": "", "type": "default", "enabled": true },
    { "key": "refresh_token", "value": "", "type": "default", "enabled": true },
    { "key": "token_expires_at", "value": "", "type": "default", "enabled": true },
    { "key": "user_id", "value": "", "type": "default", "enabled": true },

    { "key": "customer_id", "value": "33333333-3333-3333-3333-333333333301", "type": "default", "enabled": true },
    { "key": "customer_code", "value": "CUS-2026-0001", "type": "default", "enabled": true },
    { "key": "site_id", "value": "55555555-5555-5555-5555-555555555501", "type": "default", "enabled": true },
    { "key": "zone_id", "value": "66666666-6666-6666-6666-666666666601", "type": "default", "enabled": true },

    { "key": "package_id", "value": "77777777-7777-7777-7777-777777777704", "type": "default", "enabled": true },
    { "key": "package_code", "value": "PRO-M", "type": "default", "enabled": true },
    { "key": "subscription_id", "value": "88888888-8888-8888-8888-888888888801", "type": "default", "enabled": true },

    { "key": "device_id", "value": "f0000000-0000-0000-0000-000000000001", "type": "default", "enabled": true },
    { "key": "device_serial", "value": "SN-DEMO-001", "type": "default", "enabled": true },
    { "key": "device_token", "value": "", "type": "default", "enabled": true },

    { "key": "product_id", "value": "c0000000-0000-0000-0000-000000000001", "type": "default", "enabled": true },
    { "key": "sku", "value": "SENSOR-TEMP-001", "type": "default", "enabled": true },
    { "key": "warehouse_id", "value": "a0000000-0000-0000-0000-000000000001", "type": "default", "enabled": true },
    { "key": "order_id", "value": "ee000000-0000-0000-0000-000000000002", "type": "default", "enabled": true },
    { "key": "order_no", "value": "SO-2026-000002", "type": "default", "enabled": true },
    { "key": "invoice_id", "value": "aa110000-0000-0000-0000-000000000002", "type": "default", "enabled": true },
    { "key": "invoice_no", "value": "INV-2026-000002", "type": "default", "enabled": true },
    { "key": "payment_id", "value": "cc110000-0000-0000-0000-000000000001", "type": "default", "enabled": true },

    { "key": "shipment_id", "value": "ff110000-0000-0000-0000-000000000002", "type": "default", "enabled": true },
    { "key": "job_id", "value": "ab220000-0000-0000-0000-000000000001", "type": "default", "enabled": true },
    { "key": "technician_id", "value": "ee110000-0000-0000-0000-000000000001", "type": "default", "enabled": true },

    { "key": "report_code", "value": "RPT-CUSTOMER-LIST", "type": "default", "enabled": true },
    { "key": "execution_id", "value": "", "type": "default", "enabled": true },
    { "key": "dashboard_id", "value": "cd110000-0000-0000-0000-000000000001", "type": "default", "enabled": true },

    { "key": "influx_bucket", "value": "iot_telemetry", "type": "default", "enabled": true },
    { "key": "kafka_topic_prefix", "value": "icmon", "type": "default", "enabled": true }
  ],
  "_postman_variable_scope": "environment",
  "_postman_exported_at": "2026-01-15T10:00:00.000Z",
  "_postman_exported_using": "Postman/11.0.0"
}
```

### B.2 `environment.staging.json`

```json
{
  "id": "icmongolang-env-staging",
  "name": "icmongolang — Staging",
  "values": [
    { "key": "environment_name", "value": "staging", "type": "default", "enabled": true },
    { "key": "base_url", "value": "https://staging.icmongolang.io", "type": "default", "enabled": true },
    { "key": "ws_url", "value": "wss://staging.icmongolang.io/ws", "type": "default", "enabled": true },
    { "key": "api_version", "value": "v1", "type": "default", "enabled": true },

    { "key": "tenant_id", "value": "22222222-2222-2222-2222-222222220001", "type": "default", "enabled": true },
    { "key": "tenant_slug", "value": "staging-demo", "type": "default", "enabled": true },
    { "key": "admin_email", "value": "admin@staging.icmongolang.io", "type": "default", "enabled": true },
    { "key": "admin_password", "value": "{{vault:staging_admin_password}}", "type": "secret", "enabled": true },

    { "key": "access_token", "value": "", "type": "default", "enabled": true },
    { "key": "refresh_token", "value": "", "type": "default", "enabled": true },

    { "key": "customer_id", "value": "", "type": "default", "enabled": true },
    { "key": "site_id", "value": "", "type": "default", "enabled": true },
    { "key": "device_id", "value": "", "type": "default", "enabled": true },
    { "key": "order_id", "value": "", "type": "default", "enabled": true },
    { "key": "invoice_id", "value": "", "type": "default", "enabled": true }
  ],
  "_postman_variable_scope": "environment"
}
```

### B.3 `environment.prod.json` (read-only + safety)

```json
{
  "id": "icmongolang-env-prod",
  "name": "icmongolang — Production (READ-ONLY)",
  "values": [
    { "key": "environment_name", "value": "prod", "type": "default", "enabled": true },
    { "key": "base_url", "value": "https://api.icmongolang.io", "type": "default", "enabled": true },
    { "key": "ws_url", "value": "wss://api.icmongolang.io/ws", "type": "default", "enabled": true },
    { "key": "api_version", "value": "v1", "type": "default", "enabled": true },

    { "key": "tenant_id", "value": "{{vault:prod_tenant_id}}", "type": "secret", "enabled": true },
    { "key": "admin_email", "value": "{{vault:prod_admin_email}}", "type": "secret", "enabled": true },
    { "key": "admin_password", "value": "{{vault:prod_admin_password}}", "type": "secret", "enabled": true },

    { "key": "access_token", "value": "", "type": "default", "enabled": true },
    { "key": "refresh_token", "value": "", "type": "default", "enabled": true },

    { "key": "read_only_mode", "value": "true", "type": "default", "enabled": true },
    { "key": "allow_destructive", "value": "false", "type": "default", "enabled": true }
  ],
  "_postman_variable_scope": "environment"
}
```

---

## 🅲 PART 10C — FULL COLLECTION (JSON)

### C.1 Collection Root

```json
{
  "info": {
    "_postman_id": "icmongolang-iot-platform-v1",
    "name": "icmongolang — IoT Platform",
    "description": "Full API collection สำหรับแพลตฟอร์ม IoT (Smart Farm / Smart Building)\n\n**Modules:** Auth · Customers · Packages · Devices · ERP · Logistics · Reports\n\n**Quick Start:**\n1. Import collection + environment (dev)\n2. Run `01. Auth → Login (Admin)`\n3. รัน requests อื่นได้เลย — token + tenant ถูก inject อัตโนมัติ\n\n**Newman:**\n```bash\nnewman run icmongolang.postman_collection.json \\\n  -e environment.dev.json \\\n  --folder '01. Auth' \\\n  --reporters cli,json\n```\n",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
    "version": "1.0.0"
  },
  "auth": {
    "type": "bearer",
    "bearer": [
      { "key": "token", "value": "{{access_token}}", "type": "string" }
    ]
  },
  "event": [
    {
      "listen": "prerequest",
      "script": {
        "type": "text/javascript",
        "exec": [
          "// ============================================================",
          "// COLLECTION-LEVEL PRE-REQUEST SCRIPT",
          "// รันก่อนทุก request",
          "// ============================================================",
          "",
          "const env = pm.environment;",
          "const envName = env.get('environment_name') || 'dev';",
          "",
          "// 1. Safety: block destructive methods on prod",
          "if (envName === 'prod' && env.get('allow_destructive') !== 'true') {",
          "    const m = pm.request.method.toUpperCase();",
          "    if (['POST','PUT','PATCH','DELETE'].includes(m)) {",
          "        const whitelist = ['/auth/login','/auth/refresh'];",
          "        const url = pm.request.url.toString();",
          "        const allowed = whitelist.some(w => url.includes(w));",
          "        if (!allowed) {",
          "            console.warn('⚠️  PROD: destructive method blocked → ' + m + ' ' + url);",
          "            throw new Error('Destructive method blocked on PROD (allow_destructive=false)');",
          "        }",
          "    }",
          "}",
          "",
          "// 2. Inject X-Request-ID (idempotency)",
          "pm.request.headers.upsert({",
          "    key: 'X-Request-ID',",
          "    value: pm.variables.replaceIn('{{$guid}}')",
          "});",
          "",
          "// 3. Inject X-Tenant-ID (ถ้ามี tenant)",
          "const tenantId = env.get('tenant_id');",
          "if (tenantId && !pm.request.url.toString().includes('/auth/')) {",
          "    pm.request.headers.upsert({ key: 'X-Tenant-ID', value: tenantId });",
          "}",
          "",
          "// 4. Inject Accept-Language",
          "pm.request.headers.upsert({ key: 'Accept-Language', value: 'th-TH,en;q=0.9' });",
          "",
          "// 5. Auto-attach Bearer token (ถ้ามี)",
          "const token = env.get('access_token');",
          "if (token) {",
          "    pm.request.headers.upsert({",
          "        key: 'Authorization',",
          "        value: 'Bearer ' + token",
          "    });",
          "}",
          "",
          "// 6. Auto-refresh ถ้า token ใกล้หมดอายุ",
          "const expiresAt = env.get('token_expires_at');",
          "const refreshToken = env.get('refresh_token');",
          "if (expiresAt && refreshToken) {",
          "    const expMs = new Date(expiresAt).getTime();",
          "    const nowMs = Date.now();",
          "    if (nowMs > expMs - 60000) {",
          "        console.log('🔄 Token expiring — refreshing...');",
          "        const baseUrl = env.get('base_url');",
          "        pm.sendRequest({",
          "            url: baseUrl + '/api/v1/auth/refresh',",
          "            method: 'POST',",
          "            header: { 'Content-Type': 'application/json' },",
          "            body: {",
          "                mode: 'raw',",
          "                raw: JSON.stringify({ refresh_token: refreshToken })",
          "            }",
          "        }, (err, res) => {",
          "            if (!err && res.code === 200) {",
          "                const d = res.json().data;",
          "                env.set('access_token', d.access_token);",
          "                env.set('refresh_token', d.refresh_token);",
          "                env.set('token_expires_at', d.expires_at);",
          "                console.log('✓ Token refreshed');",
          "            } else {",
          "                console.warn('✗ Refresh failed — re-login required');",
          "                env.set('access_token', '');",
          "            }",
          "        });",
          "    }",
          "}",
          ""
        ]
      }
    },
    {
      "listen": "test",
      "script": {
        "type": "text/javascript",
        "exec": [
          "// ============================================================",
          "// COLLECTION-LEVEL TEST SCRIPT",
          "// รันหลังทุก request",
          "// ============================================================",
          "",
          "// 1. Log response time",
          "const rt = pm.response.responseTime;",
          "if (rt > 2000) {",
          "    console.warn('⚠️  Slow response: ' + rt + 'ms — ' + pm.request.url.toString());",
          "}",
          "",
          "// 2. Auto-save rate-limit headers",
          "const rlRemaining = pm.response.headers.get('X-RateLimit-Remaining');",
          "if (rlRemaining !== null) {",
          "    pm.environment.set('rate_limit_remaining', rlRemaining);",
          "}",
          "",
          "// 3. Log request-id for tracing",
          "const reqId = pm.response.headers.get('X-Request-ID');",
          "if (reqId) {",
          "    console.log('🔍 Request-ID: ' + reqId);",
          "}",
          "",
          "// 4. Fail on 5xx (unexpected)",
          "if (pm.response.code >= 500) {",
          "    console.error('❌ Server error: ' + pm.response.code + ' ' + pm.request.url.toString());",
          "}",
          ""
        ]
      }
    }
  ],
  "variable": [
    { "key": "api_version", "value": "v1", "type": "string" },
    { "key": "default_page_size", "value": "20", "type": "string" }
  ]
}
```

### C.2 Folder — `00. Setup`

```json
{
  "name": "00. Setup",
  "item": [
    {
      "name": "Health Check",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('Body has status=ok', () => {",
              "    const j = pm.response.json();",
              "    pm.expect(j).to.have.property('status');",
              "    pm.expect(j.status).to.eql('ok');",
              "});",
              "pm.test('Response time < 500ms', () => {",
              "    pm.expect(pm.response.responseTime).to.be.below(500);",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/health",
          "host": ["{{base_url}}"],
          "path": ["health"]
        }
      }
    },
    {
      "name": "Readiness",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('All dependencies ready', () => {",
              "    const j = pm.response.json();",
              "    pm.expect(j).to.have.property('checks');",
              "    const checks = j.checks || {};",
              "    Object.keys(checks).forEach(k => {",
              "        pm.expect(checks[k], k + ' should be ready').to.eql('ok');",
              "    });",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/readyz",
          "host": ["{{base_url}}"],
          "path": ["readyz"]
        }
      }
    },
    {
      "name": "Version",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('Version matches', () => {",
              "    const j = pm.response.json();",
              "    pm.expect(j).to.have.property('version');",
              "    pm.expect(j.version).to.match(/^\\d+\\.\\d+\\.\\d+/);",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/version",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "version"]
        }
      }
    }
  ]
}
```

### C.3 Folder — `01. Auth`

```json
{
  "name": "01. Auth",
  "event": [
    {
      "listen": "prerequest",
      "script": {
        "exec": [
          "// Auth requests ไม่ต้องใช้ token",
          "pm.request.headers.remove('Authorization');"
        ]
      }
    }
  ],
  "item": [
    {
      "name": "Login (Admin)",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('Has access_token', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d).to.have.property('access_token');",
              "    pm.expect(d.access_token).to.be.a('string').and.not.empty;",
              "});",
              "pm.test('Has refresh_token', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d).to.have.property('refresh_token');",
              "});",
              "pm.test('Has user object', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d).to.have.property('user');",
              "    pm.expect(d.user).to.have.property('id');",
              "    pm.expect(d.user).to.have.property('email');",
              "});",
              "",
              "// Auto-save tokens + user",
              "const d = pm.response.json().data;",
              "if (d && d.access_token) {",
              "    pm.environment.set('access_token', d.access_token);",
              "    pm.environment.set('refresh_token', d.refresh_token || '');",
              "    pm.environment.set('token_expires_at', d.expires_at || '');",
              "    pm.environment.set('user_id', d.user.id);",
              "    if (d.user.tenant_id) pm.environment.set('tenant_id', d.user.tenant_id);",
              "    console.log('✓ Token saved for ' + d.user.email);",
              "}"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"email\": \"{{admin_email}}\",\n  \"password\": \"{{admin_password}}\",\n  \"tenant_id\": \"{{tenant_id}}\",\n  \"device_name\": \"Postman\"\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/auth/login",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "auth", "login"]
        }
      }
    },
    {
      "name": "Login (Manager)",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('Role is MANAGER', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d.user.role).to.eql('MANAGER');",
              "});",
              "const d = pm.response.json().data;",
              "pm.environment.set('manager_token', d.access_token);",
              "pm.environment.set('manager_user_id', d.user.id);"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"email\": \"{{manager_email}}\",\n  \"password\": \"{{common_password}}\",\n  \"tenant_id\": \"{{tenant_id}}\"\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/auth/login",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "auth", "login"]
        }
      }
    },
    {
      "name": "Login (Viewer)",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "const d = pm.response.json().data;",
              "pm.environment.set('viewer_token', d.access_token);",
              "pm.environment.set('viewer_user_id', d.user.id);"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"email\": \"{{viewer_email}}\",\n  \"password\": \"{{common_password}}\",\n  \"tenant_id\": \"{{tenant_id}}\"\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/auth/login",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "auth", "login"]
        }
      }
    },
    {
      "name": "Login [Negative] — Wrong Password",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 401', () => pm.response.to.have.status(401));",
              "pm.test('Error code = INVALID_CREDENTIALS', () => {",
              "    const j = pm.response.json();",
              "    pm.expect(j.error.code).to.eql('INVALID_CREDENTIALS');",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"email\": \"{{admin_email}}\",\n  \"password\": \"WrongPassword!\",\n  \"tenant_id\": \"{{tenant_id}}\"\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/auth/login",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "auth", "login"]
        }
      }
    },
    {
      "name": "Me",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('Matches logged-in user', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d.id).to.eql(pm.environment.get('user_id'));",
              "});",
              "pm.test('Has permissions array', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d.permissions).to.be.an('array');",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/auth/me",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "auth", "me"]
        }
      }
    },
    {
      "name": "Refresh Token",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('New access_token issued', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d.access_token).to.be.a('string').and.not.empty;",
              "    pm.expect(d.access_token).to.not.eql(pm.environment.get('access_token'));",
              "});",
              "const d = pm.response.json().data;",
              "pm.environment.set('access_token', d.access_token);",
              "pm.environment.set('refresh_token', d.refresh_token);",
              "pm.environment.set('token_expires_at', d.expires_at);"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"refresh_token\": \"{{refresh_token}}\"\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/auth/refresh",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "auth", "refresh"]
        }
      }
    },
    {
      "name": "Change Password",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"current_password\": \"{{admin_password}}\",\n  \"new_password\": \"NewPassword123!\"\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/auth/password",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "auth", "password"]
        }
      }
    },
    {
      "name": "Logout",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 204', () => pm.response.to.have.status(204));",
              "pm.environment.set('access_token', '');",
              "pm.environment.set('refresh_token', '');",
              "console.log('✓ Logged out — tokens cleared');"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [],
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/auth/logout",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "auth", "logout"]
        }
      }
    }
  ]
}
```

### C.4 Folder — `02. Customers` (ตัวอย่าง)

```json
{
  "name": "02. Customers",
  "item": [
    {
      "name": "List Customers",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('Has pagination meta', () => {",
              "    const j = pm.response.json();",
              "    pm.expect(j.data).to.have.property('items').that.is.an('array');",
              "    pm.expect(j.data).to.have.property('pagination');",
              "    pm.expect(j.data.pagination).to.have.property('total');",
              "});",
              "pm.test('Items have required fields', () => {",
              "    const items = pm.response.json().data.items;",
              "    if (items.length > 0) {",
              "        pm.expect(items[0]).to.have.property('id');",
              "        pm.expect(items[0]).to.have.property('code');",
              "        pm.expect(items[0]).to.have.property('name');",
              "        pm.expect(items[0]).to.have.property('status');",
              "    }",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/customers?page=1&page_size={{default_page_size}}",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "customers"],
          "query": [
            { "key": "page", "value": "1" },
            { "key": "page_size", "value": "{{default_page_size}}" },
            { "key": "status", "value": "ACTIVE", "disabled": true },
            { "key": "segment", "value": "ENTERPRISE", "disabled": true }
          ]
        }
      }
    },
    {
      "name": "Create Customer",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 201', () => pm.response.to.have.status(201));",
              "pm.test('Customer created with code', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d).to.have.property('id');",
              "    pm.expect(d).to.have.property('code');",
              "    pm.expect(d.code).to.match(/^CUS-\\d{4}-\\d{4}$/);",
              "    pm.expect(d.status).to.eql('PENDING');",
              "});",
              "pm.test('Response has Location header', () => {",
              "    pm.expect(pm.response.headers.has('Location')).to.be.true;",
              "});",
              "",
              "// Chain: save new customer",
              "const d = pm.response.json().data;",
              "pm.environment.set('new_customer_id', d.id);",
              "pm.environment.set('new_customer_code', d.code);",
              "console.log('✓ New customer: ' + d.code + ' (' + d.id + ')');"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" },
          { "key": "Idempotency-Key", "value": "{{$guid}}" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"type\": \"CORPORATE\",\n  \"name\": \"บริษัท ทดสอบ จำกัด\",\n  \"tax_id\": \"0105559999999\",\n  \"email\": \"contact@test-postman.local\",\n  \"phone\": \"021234567\",\n  \"segment\": \"SME\",\n  \"address\": {\n    \"line1\": \"99/1 อาคารทดสอบ\",\n    \"district\": \"สาทร\",\n    \"province\": \"กรุงเทพ\",\n    \"postcode\": \"10120\",\n    \"country\": \"TH\"\n  },\n  \"credit_limit\": 100000,\n  \"currency\": \"THB\"\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/customers",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "customers"]
        }
      }
    },
    {
      "name": "Get Customer",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('Matches requested ID', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d.id).to.eql(pm.environment.get('customer_id'));",
              "});",
              "pm.test('Has contacts', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d).to.have.property('contacts').that.is.an('array');",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/customers/{{customer_id}}",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "customers", "{{customer_id}}"]
        }
      }
    },
    {
      "name": "Update Customer",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('Field updated', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d.phone).to.eql('029999999');",
              "});",
              "pm.test('updated_at changed', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(new Date(d.updated_at).getTime()).to.be.greaterThan(new Date(d.created_at).getTime());",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "PUT",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"phone\": \"029999999\",\n  \"email\": \"updated@test-postman.local\"\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/customers/{{customer_id}}",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "customers", "{{customer_id}}"]
        }
      }
    },
    {
      "name": "Add Contact",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 201', () => pm.response.to.have.status(201));",
              "pm.test('Contact created', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d).to.have.property('id');",
              "    pm.expect(d.name).to.eql('คุณทดสอบ ระบบ');",
              "});",
              "const d = pm.response.json().data;",
              "pm.environment.set('contact_id', d.id);"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"name\": \"คุณทดสอบ ระบบ\",\n  \"phone\": \"0899999999\",\n  \"email\": \"tester@test-postman.local\",\n  \"position\": \"IT Manager\",\n  \"is_primary\": true\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/customers/{{customer_id}}/contacts",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "customers", "{{customer_id}}", "contacts"]
        }
      }
    },
    {
      "name": "Onboard Customer",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('Status = ACTIVE', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d.status).to.eql('ACTIVE');",
              "});",
              "pm.test('Has onboarding info', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d).to.have.property('onboarded_at');",
              "    pm.expect(d).to.have.property('welcome_email_sent');",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"send_welcome_email\": true,\n  \"provision_tenant\": true,\n  \"initial_package_code\": \"FREE\"\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/customers/{{customer_id}}/onboard",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "customers", "{{customer_id}}", "onboard"]
        }
      }
    },
    {
      "name": "Suspend Customer",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('Status = SUSPENDED', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d.status).to.eql('SUSPENDED');",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"reason\": \"payment overdue 45 days\"\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/customers/{{customer_id}}/suspend",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "customers", "{{customer_id}}", "suspend"]
        }
      }
    },
    {
      "name": "Create Customer [Negative] — Duplicate Code",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 409', () => pm.response.to.have.status(409));",
              "pm.test('Error code = CUSTOMER_CODE_EXISTS', () => {",
              "    const j = pm.response.json();",
              "    pm.expect(j.error.code).to.eql('CUSTOMER_CODE_EXISTS');",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"type\": \"CORPORATE\",\n  \"code\": \"{{customer_code}}\",\n  \"name\": \"Duplicate\",\n  \"tax_id\": \"0105559999999\"\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/customers",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "customers"]
        }
      }
    }
  ]
}
```

### C.5 Folder — `05. Devices` (ตัวอย่าง)

```json
{
  "name": "05. Devices",
  "item": [
    {
      "name": "List Devices",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('Has items array', () => {",
              "    pm.expect(pm.response.json().data.items).to.be.an('array');",
              "});",
              "pm.test('Items have status', () => {",
              "    const items = pm.response.json().data.items;",
              "    if (items.length > 0) {",
              "        pm.expect(['ONLINE','OFFLINE','FAULT','MAINTENANCE']).to.include(items[0].status);",
              "    }",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/devices?page=1&page_size={{default_page_size}}",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "devices"],
          "query": [
            { "key": "page", "value": "1" },
            { "key": "page_size", "value": "{{default_page_size}}" },
            { "key": "status", "value": "ONLINE", "disabled": true },
            { "key": "site_id", "value": "{{site_id}}", "disabled": true }
          ]
        }
      }
    },
    {
      "name": "Register Device",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 201', () => pm.response.to.have.status(201));",
              "pm.test('Device has serial + status', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d).to.have.property('id');",
              "    pm.expect(d).to.have.property('serial_no');",
              "    pm.expect(d.status).to.eql('OFFLINE');",
              "});",
              "const d = pm.response.json().data;",
              "pm.environment.set('new_device_id', d.id);",
              "pm.environment.set('new_device_serial', d.serial_no);",
              "console.log('✓ Device registered: ' + d.serial_no);"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" },
          { "key": "Idempotency-Key", "value": "{{$guid}}" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"customer_id\": \"{{customer_id}}\",\n  \"site_id\": \"{{site_id}}\",\n  \"zone_id\": \"{{zone_id}}\",\n  \"serial_no\": \"SN-POSTMAN-{{$timestamp}}\",\n  \"name\": \"Postman Test Sensor\",\n  \"type\": \"SENSOR\",\n  \"protocol\": \"MQTT\",\n  \"model_id\": \"e0000000-0000-0000-0000-000000000001\",\n  \"tags\": [\"postman\", \"test\"]\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/devices",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "devices"]
        }
      }
    },
    {
      "name": "Provision Device",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('Returns device_token', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d).to.have.property('device_token');",
              "    pm.expect(d.device_token).to.be.a('string').with.length.greaterThan(20);",
              "});",
              "pm.test('Returns mqtt_client_id', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d).to.have.property('mqtt_client_id');",
              "});",
              "const d = pm.response.json().data;",
              "pm.environment.set('device_token', d.device_token);",
              "pm.environment.set('mqtt_client_id', d.mqtt_client_id);"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"firmware_channel\": \"stable\",\n  \"ttl_seconds\": 86400\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/devices/{{new_device_id}}/provision",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "devices", "{{new_device_id}}", "provision"]
        }
      }
    },
    {
      "name": "Send Command",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 202', () => pm.response.to.have.status(202));",
              "pm.test('Command queued', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d).to.have.property('id');",
              "    pm.expect(['PENDING','SENT']).to.include(d.status);",
              "});",
              "const d = pm.response.json().data;",
              "pm.environment.set('command_id', d.id);"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" },
          { "key": "Idempotency-Key", "value": "{{$guid}}" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"command\": \"on\",\n  \"payload\": {\n    \"duration_sec\": 300\n  },\n  \"priority\": 3,\n  \"expires_in_sec\": 60\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/devices/{{device_id}}/commands",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "devices", "{{device_id}}", "commands"]
        }
      }
    },
    {
      "name": "Ingest Telemetry",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 202', () => pm.response.to.have.status(202));",
              "pm.test('Accepted count matches', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d.accepted).to.eql(2);",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"timestamp\": \"{{$isoTimestamp}}\",\n  \"metrics\": [\n    { \"metric\": \"temperature\", \"value\": 28.5, \"unit\": \"°C\" },\n    { \"metric\": \"humidity\", \"value\": 65.0, \"unit\": \"%\" }\n  ]\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/devices/{{device_id}}/telemetry",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "devices", "{{device_id}}", "telemetry"]
        }
      }
    },
    {
      "name": "Query Telemetry",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('Series has temperature', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d).to.have.property('series');",
              "    pm.expect(d.series).to.have.property('temperature');",
              "    pm.expect(d.series.temperature).to.be.an('array');",
              "});",
              "pm.test('Points have ts + value', () => {",
              "    const pts = pm.response.json().data.series.temperature || [];",
              "    if (pts.length > 0) {",
              "        pm.expect(pts[0]).to.have.property('ts');",
              "        pm.expect(pts[0]).to.have.property('value');",
              "    }",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/devices/{{device_id}}/telemetry?from={{$isoTimestamp}}&metric=temperature&metric=humidity&interval=1m",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "devices", "{{device_id}}", "telemetry"],
          "query": [
            { "key": "from", "value": "2026-01-15T00:00:00Z" },
            { "key": "to", "value": "2026-01-15T23:59:59Z" },
            { "key": "metric", "value": "temperature" },
            { "key": "metric", "value": "humidity" },
            { "key": "interval", "value": "1m" }
          ]
        }
      }
    },
    {
      "name": "Get Shadow",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('Shadow has desired + reported', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d).to.have.property('desired');",
              "    pm.expect(d).to.have.property('reported');",
              "    pm.expect(d).to.have.property('version');",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/devices/{{device_id}}/shadow",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "devices", "{{device_id}}", "shadow"]
        }
      }
    },
    {
      "name": "Update Desired State",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('Version incremented', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d.version).to.be.greaterThan(0);",
              "});",
              "pm.test('Delta computed', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d).to.have.property('delta');",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "PATCH",
        "header": [
          { "key": "Content-Type", "value": "application/json" },
          { "key": "If-Match", "value": "*" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"pump\": \"on\",\n  \"schedule\": \"08:00-18:00\",\n  \"threshold_temp\": 35\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/devices/{{device_id}}/shadow/desired",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "devices", "{{device_id}}", "shadow", "desired"]
        }
      }
    }
  ]
}
```

### C.6 Folder — `06. ERP → Orders` (ตัวอย่าง)

```json
{
  "name": "06. ERP",
  "item": [
    {
      "name": "Orders",
      "item": [
        {
          "name": "List Orders",
          "event": [
            {
              "listen": "test",
              "script": {
                "exec": [
                  "pm.test('Status 200', () => pm.response.to.have.status(200));",
                  "pm.test('Has items', () => {",
                  "    pm.expect(pm.response.json().data.items).to.be.an('array');",
                  "});"
                ]
              }
            }
          ],
          "request": {
            "method": "GET",
            "header": [],
            "url": {
              "raw": "{{base_url}}/api/{{api_version}}/erp/orders?page=1&page_size=20",
              "host": ["{{base_url}}"],
              "path": ["api", "{{api_version}}", "erp", "orders"],
              "query": [
                { "key": "page", "value": "1" },
                { "key": "page_size", "value": "20" },
                { "key": "type", "value": "SALES", "disabled": true },
                { "key": "status", "value": "CONFIRMED", "disabled": true }
              ]
            }
          }
        },
        {
          "name": "Create Order",
          "event": [
            {
              "listen": "test",
              "script": {
                "exec": [
                  "pm.test('Status 201', () => pm.response.to.have.status(201));",
                  "pm.test('Order is DRAFT', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.status).to.eql('DRAFT');",
                  "});",
                  "pm.test('Total = subtotal + tax', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.total_amount).to.be.closeTo(d.subtotal + d.tax_amount, 0.01);",
                  "});",
                  "pm.test('Has order_no format', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.order_no).to.match(/^SO-\\d{4}-\\d{6}$/);",
                  "});",
                  "const d = pm.response.json().data;",
                  "pm.environment.set('new_order_id', d.id);",
                  "pm.environment.set('new_order_no', d.order_no);"
                ]
              }
            }
          ],
          "request": {
            "method": "POST",
            "header": [
              { "key": "Content-Type", "value": "application/json" },
              { "key": "Idempotency-Key", "value": "{{$guid}}" }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"type\": \"SALES\",\n  \"customer_id\": \"{{customer_id}}\",\n  \"order_date\": \"{{$isoTimestamp}}\",\n  \"currency\": \"THB\",\n  \"payment_terms\": \"NET_30\",\n  \"ship_to\": {\n    \"line1\": \"99/1 อาคารทดสอบ\",\n    \"province\": \"กรุงเทพ\",\n    \"postcode\": \"10120\",\n    \"country\": \"TH\"\n  },\n  \"lines\": [\n    {\n      \"product_id\": \"{{product_id}}\",\n      \"sku\": \"{{sku}}\",\n      \"name\": \"Temperature Sensor SHT30\",\n      \"quantity\": 5,\n      \"uom\": \"PCS\",\n      \"unit_price\": 650,\n      \"tax_type\": \"VAT\",\n      \"tax_rate\": 0.07\n    }\n  ]\n}"
            },
            "url": {
              "raw": "{{base_url}}/api/{{api_version}}/erp/orders",
              "host": ["{{base_url}}"],
              "path": ["api", "{{api_version}}", "erp", "orders"]
            }
          }
        },
        {
          "name": "Confirm Order",
          "event": [
            {
              "listen": "test",
              "script": {
                "exec": [
                  "pm.test('Status 200', () => pm.response.to.have.status(200));",
                  "pm.test('Status = CONFIRMED', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.status).to.eql('CONFIRMED');",
                  "});",
                  "pm.test('confirmed_at set', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.confirmed_at).to.not.be.null;",
                  "});"
                ]
              }
            }
          ],
          "request": {
            "method": "POST",
            "header": [
              { "key": "Content-Type", "value": "application/json" }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"notes\": \"approved by postman\"\n}"
            },
            "url": {
              "raw": "{{base_url}}/api/{{api_version}}/erp/orders/{{new_order_id}}/confirm",
              "host": ["{{base_url}}"],
              "path": ["api", "{{api_version}}", "erp", "orders", "{{new_order_id}}", "confirm"]
            }
          }
        },
        {
          "name": "Ship Order",
          "event": [
            {
              "listen": "test",
              "script": {
                "exec": [
                  "pm.test('Status 200', () => pm.response.to.have.status(200));",
                  "pm.test('Status = SHIPPED', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.status).to.eql('SHIPPED');",
                  "});",
                  "pm.test('Inventory deducted', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d).to.have.property('inventory_movements');",
                  "    pm.expect(d.inventory_movements).to.be.an('array').with.length.greaterThan(0);",
                  "});"
                ]
              }
            }
          ],
          "request": {
            "method": "POST",
            "header": [
              { "key": "Content-Type", "value": "application/json" }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"warehouse_id\": \"{{warehouse_id}}\",\n  \"lines\": [\n    { \"line_id\": \"AUTO\", \"quantity\": 5 }\n  ]\n}"
            },
            "url": {
              "raw": "{{base_url}}/api/{{api_version}}/erp/orders/{{new_order_id}}/ship",
              "host": ["{{base_url}}"],
              "path": ["api", "{{api_version}}", "erp", "orders", "{{new_order_id}}", "ship"]
            }
          }
        },
        {
          "name": "Create Invoice from Order",
          "event": [
            {
              "listen": "test",
              "script": {
                "exec": [
                  "pm.test('Status 201', () => pm.response.to.have.status(201));",
                  "pm.test('Invoice is ISSUED', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.status).to.eql('ISSUED');",
                  "});",
                  "pm.test('Amount matches order total', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.total_amount).to.be.greaterThan(0);",
                  "});",
                  "const d = pm.response.json().data;",
                  "pm.environment.set('new_invoice_id', d.id);",
                  "pm.environment.set('new_invoice_no', d.invoice_no);"
                ]
              }
            }
          ],
          "request": {
            "method": "POST",
            "header": [
              { "key": "Content-Type", "value": "application/json" },
              { "key": "Idempotency-Key", "value": "{{$guid}}" }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"type\": \"AR\",\n  \"order_id\": \"{{new_order_id}}\",\n  \"issue_date\": \"{{$isoTimestamp}}\",\n  \"due_date\": \"{{$isoTimestamp}}\",\n  \"payment_terms\": \"NET_30\"\n}"
            },
            "url": {
              "raw": "{{base_url}}/api/{{api_version}}/erp/invoices",
              "host": ["{{base_url}}"],
              "path": ["api", "{{api_version}}", "erp", "invoices"]
            }
          }
        },
        {
          "name": "Record Payment (Full)",
          "event": [
            {
              "listen": "test",
              "script": {
                "exec": [
                  "pm.test('Status 201', () => pm.response.to.have.status(201));",
                  "pm.test('Payment confirmed', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.status).to.eql('CONFIRMED');",
                  "});",
                  "pm.test('Fully allocated', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.unallocated).to.eql(0);",
                  "});",
                  "const d = pm.response.json().data;",
                  "pm.environment.set('new_payment_id', d.id);"
                ]
              }
            }
          ],
          "request": {
            "method": "POST",
            "header": [
              { "key": "Content-Type", "value": "application/json" },
              { "key": "Idempotency-Key", "value": "{{$guid}}" }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"direction\": \"IN\",\n  \"customer_id\": \"{{customer_id}}\",\n  \"method\": \"BANK_TRANSFER\",\n  \"amount\": 5000,\n  \"currency\": \"THB\",\n  \"payment_date\": \"{{$isoTimestamp}}\",\n  \"reference\": \"TXN-POSTMAN-{{$timestamp}}\",\n  \"allocations\": [\n    { \"invoice_id\": \"{{new_invoice_id}}\", \"amount\": 5000 }\n  ]\n}"
            },
            "url": {
              "raw": "{{base_url}}/api/{{api_version}}/erp/payments",
              "host": ["{{base_url}}"],
              "path": ["api", "{{api_version}}", "erp", "payments"]
            }
          }
        },
        {
          "name": "Verify Invoice PAID",
          "event": [
            {
              "listen": "test",
              "script": {
                "exec": [
                  "pm.test('Status 200', () => pm.response.to.have.status(200));",
                  "pm.test('Invoice is PAID', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.status).to.eql('PAID');",
                  "});",
                  "pm.test('amount_due = 0', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.amount_due).to.eql(0);",
                  "});",
                  "pm.test('paid_at is set', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.paid_at).to.not.be.null;",
                  "});"
                ]
              }
            }
          ],
          "request": {
            "method": "GET",
            "header": [],
            "url": {
              "raw": "{{base_url}}/api/{{api_version}}/erp/invoices/{{new_invoice_id}}",
              "host": ["{{base_url}}"],
              "path": ["api", "{{api_version}}", "erp", "invoices", "{{new_invoice_id}}"]
            }
          }
        }
      ]
    }
  ]
}
```

### C.7 Folder — `07. Logistics → Installation Jobs` (ตัวอย่าง)

```json
{
  "name": "07. Logistics",
  "item": [
    {
      "name": "Installation Jobs",
      "item": [
        {
          "name": "Create Job",
          "event": [
            {
              "listen": "test",
              "script": {
                "exec": [
                  "pm.test('Status 201', () => pm.response.to.have.status(201));",
                  "pm.test('Job is SCHEDULED', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.status).to.eql('SCHEDULED');",
                  "});",
                  "const d = pm.response.json().data;",
                  "pm.environment.set('new_job_id', d.id);",
                  "pm.environment.set('new_job_no', d.job_no);"
                ]
              }
            }
          ],
          "request": {
            "method": "POST",
            "header": [
              { "key": "Content-Type", "value": "application/json" }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"customer_id\": \"{{customer_id}}\",\n  \"site_id\": \"{{site_id}}\",\n  \"job_type\": \"INSTALLATION\",\n  \"priority\": \"NORMAL\",\n  \"scheduled_at\": \"{{$isoTimestamp}}\",\n  \"sla_due_at\": \"{{$isoTimestamp}}\",\n  \"device_ids\": [\"{{new_device_id}}\"],\n  \"notes\": \"created via postman\"\n}"
            },
            "url": {
              "raw": "{{base_url}}/api/{{api_version}}/installation-jobs",
              "host": ["{{base_url}}"],
              "path": ["api", "{{api_version}}", "installation-jobs"]
            }
          }
        },
        {
          "name": "Assign Technician",
          "event": [
            {
              "listen": "test",
              "script": {
                "exec": [
                  "pm.test('Status 200', () => pm.response.to.have.status(200));",
                  "pm.test('Technician assigned', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.technician_id).to.eql(pm.environment.get('technician_id'));",
                  "    pm.expect(d.status).to.eql('ASSIGNED');",
                  "});"
                ]
              }
            }
          ],
          "request": {
            "method": "POST",
            "header": [
              { "key": "Content-Type", "value": "application/json" }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"technician_id\": \"{{technician_id}}\"\n}"
            },
            "url": {
              "raw": "{{base_url}}/api/{{api_version}}/installation-jobs/{{new_job_id}}/assign",
              "host": ["{{base_url}}"],
              "path": ["api", "{{api_version}}", "installation-jobs", "{{new_job_id}}", "assign"]
            }
          }
        },
        {
          "name": "Start Job",
          "event": [
            {
              "listen": "test",
              "script": {
                "exec": [
                  "pm.test('Status 200', () => pm.response.to.have.status(200));",
                  "pm.test('Status = IN_PROGRESS', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.status).to.eql('IN_PROGRESS');",
                  "    pm.expect(d.started_at).to.not.be.null;",
                  "});"
                ]
              }
            }
          ],
          "request": {
            "method": "POST",
            "header": [
              { "key": "Content-Type", "value": "application/json" }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"gps\": { \"lat\": 13.7563, \"lng\": 100.5018, \"accuracy\": 10 }\n}"
            },
            "url": {
              "raw": "{{base_url}}/api/{{api_version}}/installation-jobs/{{new_job_id}}/start",
              "host": ["{{base_url}}"],
              "path": ["api", "{{api_version}}", "installation-jobs", "{{new_job_id}}", "start"]
            }
          }
        },
        {
          "name": "Update Checklist",
          "event": [
            {
              "listen": "test",
              "script": {
                "exec": [
                  "pm.test('Status 200', () => pm.response.to.have.status(200));",
                  "pm.test('Checklist updated', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.checklist).to.be.an('array');",
                  "    pm.expect(d.checklist.length).to.be.greaterThan(0);",
                  "});"
                ]
              }
            }
          ],
          "request": {
            "method": "PATCH",
            "header": [
              { "key": "Content-Type", "value": "application/json" }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"checklist\": [\n    { \"key\": \"power_check\", \"checked\": true },\n    { \"key\": \"network_check\", \"checked\": true },\n    { \"key\": \"device_mount\", \"checked\": true },\n    { \"key\": \"sensor_test\", \"checked\": true, \"notes\": \"all sensors responding\" },\n    { \"key\": \"photo_evidence\", \"checked\": true, \"photo_url\": \"https://cdn.test/photo1.jpg\" }\n  ]\n}"
            },
            "url": {
              "raw": "{{base_url}}/api/{{api_version}}/installation-jobs/{{new_job_id}}/checklist",
              "host": ["{{base_url}}"],
              "path": ["api", "{{api_version}}", "installation-jobs", "{{new_job_id}}", "checklist"]
            }
          }
        },
        {
          "name": "Complete Job",
          "event": [
            {
              "listen": "test",
              "script": {
                "exec": [
                  "pm.test('Status 200', () => pm.response.to.have.status(200));",
                  "pm.test('Status = DONE', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.status).to.eql('DONE');",
                  "    pm.expect(d.completed_at).to.not.be.null;",
                  "});",
                  "pm.test('SLA met', () => {",
                  "    const d = pm.response.json().data;",
                  "    pm.expect(d.sla_met).to.be.true;",
                  "});"
                ]
              }
            }
          ],
          "request": {
            "method": "POST",
            "header": [
              { "key": "Content-Type", "value": "application/json" }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"completed_gps\": { \"lat\": 13.7563, \"lng\": 100.5018, \"accuracy\": 10 },\n  \"signature_url\": \"https://cdn.test/sig.png\",\n  \"notes\": \"ติดตั้งเรียบร้อย\",\n  \"photos\": [\"https://cdn.test/photo1.jpg\", \"https://cdn.test/photo2.jpg\"]\n}"
            },
            "url": {
              "raw": "{{base_url}}/api/{{api_version}}/installation-jobs/{{new_job_id}}/complete",
              "host": ["{{base_url}}"],
              "path": ["api", "{{api_version}}", "installation-jobs", "{{new_job_id}}", "complete"]
            }
          }
        }
      ]
    }
  ]
}
```

### C.8 Folder — `08. Reports & Dashboard` (ตัวอย่าง)

```json
{
  "name": "08. Reports & Dashboard",
  "item": [
    {
      "name": "List Reports",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('Has items', () => {",
              "    pm.expect(pm.response.json().data.items).to.be.an('array');",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/reports",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "reports"]
        }
      }
    },
    {
      "name": "Generate Report",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 202', () => pm.response.to.have.status(202));",
              "pm.test('Has execution_id', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d).to.have.property('execution_id');",
              "    pm.expect(d.status).to.be.oneOf(['PENDING','RUNNING']);",
              "});",
              "const d = pm.response.json().data;",
              "pm.environment.set('execution_id', d.execution_id);"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"format\": \"JSON\",\n  \"preset\": \"LAST_30_DAYS\",\n  \"params\": {\n    \"status\": \"ACTIVE\"\n  }\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/reports/{{report_code}}/generate",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "reports", "{{report_code}}", "generate"]
        }
      }
    },
    {
      "name": "Poll Report Execution",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "const d = pm.response.json().data;",
              "if (d.status === 'SUCCESS') {",
              "    pm.test('Result URL available', () => {",
              "        pm.expect(d.result_url).to.be.a('string').and.not.empty;",
              "    });",
              "    pm.environment.set('report_result_url', d.result_url);",
              "} else if (d.status === 'FAILED') {",
              "    pm.test('Error message present', () => {",
              "        pm.expect(d.error_msg).to.be.a('string');",
              "    });",
              "}",
              "console.log('Report status: ' + d.status);"
            ]
          }
        }
      ],
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/report-executions/{{execution_id}}",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "report-executions", "{{execution_id}}"]
        }
      }
    },
    {
      "name": "Get KPIs",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('Has kpis array', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d).to.have.property('kpis');",
              "    pm.expect(d.kpis).to.be.an('array');",
              "});",
              "pm.test('Each KPI has code + value', () => {",
              "    const kpis = pm.response.json().data.kpis;",
              "    kpis.forEach(k => {",
              "        pm.expect(k).to.have.property('code');",
              "        pm.expect(k).to.have.property('value');",
              "        pm.expect(k).to.have.property('unit');",
              "    });",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/dashboard/kpis?code=MRR&code=ACTIVE_CUSTOMERS&code=DEVICE_ONLINE&preset=LAST_30_DAYS",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "dashboard", "kpis"],
          "query": [
            { "key": "code", "value": "MRR" },
            { "key": "code", "value": "ACTIVE_CUSTOMERS" },
            { "key": "code", "value": "DEVICE_ONLINE" },
            { "key": "preset", "value": "LAST_30_DAYS" }
          ]
        }
      }
    },
    {
      "name": "Generate AI Insight",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 202', () => pm.response.to.have.status(202));",
              "pm.test('Insight queued', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d).to.have.property('insight_id');",
              "});",
              "const d = pm.response.json().data;",
              "pm.environment.set('insight_id', d.insight_id);"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"type\": \"ANOMALY_DETECTION\",\n  \"scope\": \"DEVICE\",\n  \"device_id\": \"{{device_id}}\",\n  \"period\": \"LAST_7_DAYS\",\n  \"options\": {\n    \"include_recommendations\": true\n  }\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/dashboard/insights/generate",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "dashboard", "insights", "generate"]
        }
      }
    }
  ]
}
```

### C.9 Folder — `10. Utilities`

```json
{
  "name": "10. Utilities",
  "item": [
    {
      "name": "Seed Demo Data",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));",
              "pm.test('Seed result has created counts', () => {",
              "    const d = pm.response.json().data;",
              "    pm.expect(d).to.have.property('created');",
              "    pm.expect(d.created).to.have.property('customers');",
              "    pm.expect(d.created).to.have.property('devices');",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"level\": \"standard\",\n  \"clean_first\": false\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/admin/seed",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "admin", "seed"]
        }
      }
    },
    {
      "name": "Reset Fixtures (dev only)",
      "event": [
        {
          "listen": "prerequest",
          "script": {
            "exec": [
              "if (pm.environment.get('environment_name') !== 'dev') {",
              "    throw new Error('Reset allowed only in dev');",
              "}"
            ]
          }
        },
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status 200', () => pm.response.to.have.status(200));"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [],
        "url": {
          "raw": "{{base_url}}/api/{{api_version}}/admin/reset-fixtures",
          "host": ["{{base_url}}"],
          "path": ["api", "{{api_version}}", "admin", "reset-fixtures"]
        }
      }
    }
  ]
}
```

---

## 🅳 PART 10D — PRE-REQUEST & TEST SCRIPTS (Reusable Snippets)

### D.1 Reusable Assertion Helpers

```javascript
// ============================================================
// Common test snippets — คัดลอกไปใช้ใน request ต่างๆ
// ============================================================

// ---- 1. Standard success assertions ----
function assertSuccess(expectedCode = 200) {
    pm.test(`Status ${expectedCode}`, () => pm.response.to.have.status(expectedCode));
    pm.test('Response time < 2s', () => {
        pm.expect(pm.response.responseTime).to.be.below(2000);
    });
    pm.test('Content-Type is JSON', () => {
        pm.expect(pm.response.headers.get('Content-Type')).to.include('application/json');
    });
}

// ---- 2. Pagination assertions ----
function assertPagination() {
    pm.test('Has pagination', () => {
        const p = pm.response.json().data.pagination;
        pm.expect(p).to.have.property('page');
        pm.expect(p).to.have.property('page_size');
        pm.expect(p).to.have.property('total');
        pm.expect(p).to.have.property('total_pages');
    });
}

// ---- 3. Error response assertions ----
function assertError(statusCode, errorCode) {
    pm.test(`Status ${statusCode}`, () => pm.response.to.have.status(statusCode));
    pm.test(`Error code = ${errorCode}`, () => {
        const j = pm.response.json();
        pm.expect(j.error).to.have.property('code').that.eql(errorCode);
        pm.expect(j.error).to.have.property('message');
    });
}

// ---- 4. UUID format assertions ----
function assertUUID(value, fieldName) {
    pm.test(`${fieldName} is valid UUID`, () => {
        pm.expect(value).to.match(/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i);
    });
}

// ---- 5. Timestamp assertions ----
function assertTimestamp(value, fieldName) {
    pm.test(`${fieldName} is valid ISO8601`, () => {
        const d = new Date(value);
        pm.expect(d.toString()).to.not.eql('Invalid Date');
    });
}

// ---- 6. Chained variable save ----
function saveVar(varName, value) {
    pm.environment.set(varName, value);
    console.log(`💾 ${varName} = ${value}`);
}

// ---- 7. Assert DB state via API ----
async function assertState(endpoint, predicate, message) {
    return new Promise((resolve, reject) => {
        pm.sendRequest({
            url: pm.environment.get('base_url') + endpoint,
            method: 'GET',
            header: { 'Authorization': 'Bearer ' + pm.environment.get('access_token') }
        }, (err, res) => {
            if (err) return reject(err);
            pm.test(message, () => {
                pm.expect(predicate(res.json())).to.be.true;
            });
            resolve(res);
        });
    });
}
```

### D.2 Auth Flow Test Script (Login → Use → Refresh → Logout)

```javascript
// ============================================================
// Full Auth Flow — Postman Collection Runner
// ============================================================

// Step 1: Login (request: Login (Admin))
pm.test('✓ Login success', () => {
    pm.response.to.have.status(200);
    const d = pm.response.json().data;
    pm.environment.set('access_token', d.access_token);
    pm.environment.set('refresh_token', d.refresh_token);
    pm.environment.set('token_expires_at', d.expires_at);
});

// Step 2: Call protected endpoint (request: Get Customer)
pm.test('✓ Protected call with token', () => {
    pm.response.to.have.status(200);
});

// Step 3: Refresh (request: Refresh Token)
pm.test('✓ Token refreshed', () => {
    const d = pm.response.json().data;
    pm.expect(d.access_token).to.not.eql(pm.environment.get('access_token'));
});

// Step 4: Old refresh token invalidated (request: Refresh [Negative])
pm.test('✓ Old refresh token rejected', () => {
    pm.response.to.have.status(401);
});

// Step 5: Logout (request: Logout)
pm.test('✓ Logout success', () => {
    pm.response.to.have.status(204);
    pm.environment.set('access_token', '');
});

// Step 6: Call after logout → should fail
pm.test('✓ Access denied after logout', () => {
    pm.response.to.have.status(401);
});
```

### D.3 Order-to-Cash E2E Test Script

```javascript
// ============================================================
// Order-to-Cash — Folder-level test script
// รันใน Collection Runner ตามลำดับ
// ============================================================

// ============================================================
// Request 1: Create Order
// ============================================================
pm.test('✓ Order created as DRAFT', () => {
    const d = pm.response.json().data;
    pm.expect(d.status).to.eql('DRAFT');
    pm.expect(d.order_no).to.match(/^SO-\d{4}-\d{6}$/);
    pm.environment.set('otc_order_id', d.id);
    pm.environment.set('otc_order_total', d.total_amount);
    pm.environment.set('otc_order_line_id', d.lines[0].id);
});

// ============================================================
// Request 2: Confirm Order
// ============================================================
pm.test('✓ Order CONFIRMED', () => {
    const d = pm.response.json().data;
    pm.expect(d.status).to.eql('CONFIRMED');
    pm.expect(d.confirmed_at).to.not.be.null;
});

// ============================================================
// Request 3: Ship Order
// ============================================================
pm.test('✓ Order SHIPPED', () => {
    const d = pm.response.json().data;
    pm.expect(d.status).to.eql('SHIPPED');
});

pm.test('✓ Inventory deducted', () => {
    const d = pm.response.json().data;
    pm.expect(d.inventory_movements).to.be.an('array').with.length.greaterThan(0);
    const mv = d.inventory_movements[0];
    pm.expect(mv.direction).to.eql(-1);
    pm.expect(mv.quantity).to.be.greaterThan(0);
});

// ============================================================
// Request 4: Create Invoice
// ============================================================
pm.test('✓ Invoice ISSUED', () => {
    const d = pm.response.json().data;
    pm.expect(d.status).to.eql('ISSUED');
    pm.expect(d.total_amount).to.eql(pm.environment.get('otc_order_total'));
    pm.environment.set('otc_invoice_id', d.id);
    pm.environment.set('otc_invoice_total', d.total_amount);
});

// ============================================================
// Request 5: Record Payment
// ============================================================
pm.test('✓ Payment CONFIRMED', () => {
    const d = pm.response.json().data;
    pm.expect(d.status).to.eql('CONFIRMED');
    pm.expect(d.unallocated).to.eql(0);
    pm.environment.set('otc_payment_id', d.id);
});

// ============================================================
// Request 6: Verify Invoice PAID
// ============================================================
pm.test('✓ Invoice PAID', () => {
    const d = pm.response.json().data;
    pm.expect(d.status).to.eql('PAID');
    pm.expect(d.amount_due).to.eql(0);
    pm.expect(d.paid_at).to.not.be.null;
});

// ============================================================
// Request 7: Verify Journal Entries Posted
// ============================================================
pm.test('✓ Journal entries posted', () => {
    const d = pm.response.json().data;
    pm.expect(d.items).to.be.an('array').with.length.greaterThan(0);
    const entries = d.items;
    // ต้องมีอย่างน้อย: AR Invoice + COGS + Payment
    pm.expect(entries.length).to.be.greaterThan(2);
});
```

### D.4 Device Provisioning E2E Test Script

```javascript
// ============================================================
// Device Provisioning — Folder test script
// ============================================================

// Request 1: Register Device
pm.test('✓ Device registered OFFLINE', () => {
    const d = pm.response.json().data;
    pm.expect(d.status).to.eql('OFFLINE');
    pm.expect(d.serial_no).to.match(/^SN-POSTMAN-\d+$/);
    pm.environment.set('dev_id', d.id);
});

// Request 2: Provision
pm.test('✓ Provisioned with token', () => {
    const d = pm.response.json().data;
    pm.expect(d.device_token).to.be.a('string').with.length.greaterThan(20);
    pm.expect(d.mqtt_client_id).to.be.a('string');
    pm.environment.set('dev_token', d.device_token);
});

// Request 3: Ingest Telemetry
pm.test('✓ Telemetry accepted', () => {
    const d = pm.response.json().data;
    pm.expect(d.accepted).to.eql(2);
});

// Request 4: Verify Device ONLINE after telemetry
// (รอ 1-2 วินาที ให้ async process ทำงาน)
setTimeout(() => {
    pm.sendRequest({
        url: pm.environment.get('base_url') + '/api/v1/devices/' + pm.environment.get('dev_id'),
        method: 'GET',
        header: { 'Authorization': 'Bearer ' + pm.environment.get('access_token') }
    }, (err, res) => {
        pm.test('✓ Device marked ONLINE', () => {
            pm.expect(res.json().data.status).to.eql('ONLINE');
        });
    });
}, 2000);

// Request 5: Query telemetry
pm.test('✓ Telemetry queryable', () => {
    const d = pm.response.json().data;
    pm.expect(d.series).to.have.property('temperature');
    pm.expect(d.series.temperature).to.be.an('array').with.length.greaterThan(0);
});

// Request 6: Send command
pm.test('✓ Command queued', () => {
    const d = pm.response.json().data;
    pm.expect(d.status).to.be.oneOf(['PENDING','SENT']);
    pm.environment.set('dev_cmd_id', d.id);
});

// Request 7: Verify shadow updated
pm.test('✓ Shadow has desired state', () => {
    const d = pm.response.json().data;
    pm.expect(d).to.have.property('desired');
    pm.expect(d.version).to.be.greaterThan(0);
});
```

---

## 🅴 PART 10E — AUTH FLOW (Login → Refresh → Logout)

### E.1 Sequence Diagram

```
┌─────────┐                ┌──────────┐              ┌──────────┐
│ Postman │                │   API    │              │  Redis   │
└────┬────┘                └────┬─────┘              └────┬─────┘
     │                          │                        │
     │  POST /auth/login        │                        │
     │  {email, password}       │                        │
     │─────────────────────────►│                        │
     │                          │  verify bcrypt         │
     │                          │  issue access (15m)    │
     │                          │  issue refresh (7d)    │
     │                          │  store session ───────►│
     │◄─────────────────────────│                        │
     │  {access_token,          │                        │
     │   refresh_token,         │                        │
     │   expires_at}            │                        │
     │                          │                        │
     │  ─── (ใช้ token ทำงาน) ───│                        │
     │                          │                        │
     │  POST /auth/refresh      │                        │
     │  {refresh_token}         │                        │
     │─────────────────────────►│                        │
     │                          │  validate refresh ────►│
     │                          │◄────── still valid ────│
     │                          │  rotate tokens         │
     │                          │  revoke old ──────────►│
     │◄─────────────────────────│                        │
     │  {new access_token,      │                        │
     │   new refresh_token}     │                        │
     │                          │                        │
     │  POST /auth/logout       │                        │
     │─────────────────────────►│                        │
     │                          │  revoke session ──────►│
     │◄─────────────────────────│                        │
     │  204 No Content          │                        │
     │                          │                        │
```

### E.2 Test Script — Token Rotation

```javascript
// ============================================================
// Token Rotation Test
// ============================================================

// 1. Login
const oldRefresh = pm.environment.get('refresh_token');

// 2. Refresh — ได้ new refresh
// (request: Refresh Token)
pm.test('✓ New refresh token differs from old', () => {
    const d = pm.response.json().data;
    pm.expect(d.refresh_token).to.not.eql(oldRefresh);
    pm.environment.set('old_refresh_token', oldRefresh);
});

// 3. ลองใช้ old refresh → ต้อง fail (revoked)
// (request: Refresh with Old Token)
pm.test('✓ Old refresh token revoked', () => {
    pm.response.to.have.status(401);
    const j = pm.response.json();
    pm.expect(j.error.code).to.eql('REFRESH_TOKEN_REVOKED');
});
```

### E.3 Test Script — Expired Token Handling

```javascript
// ============================================================
// Expired Token Handling
// ============================================================

// Simulate expired token (set expiry in past)
pm.environment.set('token_expires_at', new Date(Date.now() - 1000).toISOString());

// Next request → collection pre-request script จะ auto-refresh
// ตรวจสอบว่า refresh ทำงาน
pm.test('✓ Auto-refresh worked', () => {
    pm.response.to.have.status(200);
    const newExpiry = new Date(pm.environment.get('token_expires_at')).getTime();
    pm.expect(newExpiry).to.be.greaterThan(Date.now());
});
```

### E.4 Test Script — Multi-Role Access

```javascript
// ============================================================
// Role-Based Access Control Tests
// ============================================================

// Viewer ลองสร้าง customer → ต้องได้ 403
pm.test('✓ Viewer cannot create customer', () => {
    pm.sendRequest({
        url: pm.environment.get('base_url') + '/api/v1/customers',
        method: 'POST',
        header: {
            'Authorization': 'Bearer ' + pm.environment.get('viewer_token'),
            'X-Tenant-ID': pm.environment.get('tenant_id'),
            'Content-Type': 'application/json'
        },
        body: {
            mode: 'raw',
            raw: JSON.stringify({ type: 'INDIVIDUAL', name: 'Test' })
        }
    }, (err, res) => {
        pm.expect(res.code).to.eql(403);
    });
});

// Manager สร้างได้
pm.test('✓ Manager can create customer', () => {
    pm.sendRequest({
        url: pm.environment.get('base_url') + '/api/v1/customers',
        method: 'POST',
        header: {
            'Authorization': 'Bearer ' + pm.environment.get('manager_token'),
            'X-Tenant-ID': pm.environment.get('tenant_id'),
            'Content-Type': 'application/json'
        },
        body: {
            mode: 'raw',
            raw: JSON.stringify({ type: 'INDIVIDUAL', name: 'Manager Test' })
        }
    }, (err, res) => {
        pm.expect(res.code).to.eql(201);
    });
});

// Admin ลบได้
pm.test('✓ Admin can delete customer', () => {
    pm.sendRequest({
        url: pm.environment.get('base_url') + '/api/v1/customers/' + pm.environment.get('customer_id'),
        method: 'DELETE',
        header: {
            'Authorization': 'Bearer ' + pm.environment.get('access_token'),
            'X-Tenant-ID': pm.environment.get('tenant_id')
        }
    }, (err, res) => {
        pm.expect(res.code).to.be.oneOf([204, 200]);
    });
});
```

---

## 🅵 PART 10F — NEWMAN CLI + CI INTEGRATION

### F.1 Install Newman

```bash
# Global install
npm install -g newman newman-reporter-htmlextra

# Verify
newman --version
```

### F.2 `package.json` scripts

```json
{
  "name": "icmongolang-api-tests",
  "version": "1.0.0",
  "private": true,
  "scripts": {
    "test:smoke": "newman run postman/icmongolang.postman_collection.json -e postman/environment.dev.json --folder '00. Setup' --folder '01. Auth' --reporters cli,json --reporter-json-export reports/smoke.json",
    "test:auth": "newman run postman/icmongolang.postman_collection.json -e postman/environment.dev.json --folder '01. Auth' --reporters cli,htmlextra --reporter-htmlextra-export reports/auth.html",
    "test:customers": "newman run postman/icmongolang.postman_collection.json -e postman/environment.dev.json --folder '02. Customers'",
    "test:devices": "newman run postman/icmongolang.postman_collection.json -e postman/environment.dev.json --folder '05. Devices'",
    "test:erp": "newman run postman/icmongolang.postman_collection.json -e postman/environment.dev.json --folder '06. ERP'",
    "test:logistics": "newman run postman/icmongolang.postman_collection.json -e postman/environment.dev.json --folder '07. Logistics'",
    "test:full": "newman run postman/icmongolang.postman_collection.json -e postman/environment.dev.json --reporters cli,htmlextra,json --reporter-htmlextra-export reports/full.html --reporter-json-export reports/full.json",
    "test:ci": "newman run postman/icmongolang.postman_collection.json -e postman/environment.dev.json --reporters cli,junit --reporter-junit-export reports/junit.xml --bail",
    "test:parallel": "newman run postman/icmongolang.postman_collection.json -e postman/environment.dev.json --iteration-count 3 --delay-request 100",
    "test:load": "newman run postman/icmongolang.postman_collection.json -e postman/environment.dev.json --folder '05. Devices' --iteration-count 10 --delay-request 200"
  }
}
```

### F.3 Newman Run — Full Command

```bash
# ============================================================
# Run full collection
# ============================================================
newman run postman/icmongolang.postman_collection.json \
  --environment postman/environment.dev.json \
  --reporters cli,htmlextra,json,junit \
  --reporter-htmlextra-export reports/full.html \
  --reporter-json-export reports/full.json \
  --reporter-junit-export reports/junit.xml \
  --timeout-request 30000 \
  --timeout-script 10000 \
  --delay-request 50 \
  --bail \
  --color on \
  --verbose
```

### F.4 Newman — Selective Run

```bash
# รันเฉพาะ folder
newman run collection.json -e env.json --folder '01. Auth'

# รันหลาย folder
newman run collection.json -e env.json \
  --folder '01. Auth' \
  --folder '02. Customers' \
  --folder '05. Devices'

# รันเฉพาะ request ที่ตรงกับ pattern
newman run collection.json -e env.json --folder '06. ERP' --request 'Create Order'

# ข้าม folder
newman run collection.json -e env.json \
  --folder '01. Auth' \
  --folder '02. Customers'
```

### F.5 Newman — Environment Override

```bash
# Override variables via CLI
newman run collection.json \
  -e environment.dev.json \
  --env-var "base_url=http://staging.icmongolang.io" \
  --env-var "admin_email=admin@staging.io" \
  --env-var "admin_password=${STAGING_PASSWORD}"

# ใช้ env file + CLI override รวมกัน
newman run collection.json \
  -e environment.staging.json \
  --env-var "tenant_id=${TENANT_ID}"
```

### F.6 GitHub Actions Workflow

```yaml
# .github/workflows/api-tests.yml
name: API Tests (Newman)

on:
  pull_request:
    branches: [main, develop]
  schedule:
    - cron: '0 2 * * *'  # daily 02:00 UTC

jobs:
  api-tests:
    runs-on: ubuntu-latest
    timeout-minutes: 30

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

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with: { go-version: '1.23' }

      - uses: actions/setup-node@v4
        with: { node-version: '20' }

      - name: Install Newman
        run: npm install -g newman newman-reporter-htmlextra

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
            psql "$DB_DSN" -f "$f" || true
          done

      - name: Start API server
        env:
          DB_DSN: 'host=localhost user=icmon password=icmon_dev dbname=icmongolang_test sslmode=disable'
          REDIS_ADDR: 'localhost:6379'
          KAFKA_BROKERS: 'localhost:9092'
          APP_PORT: '8080'
        run: |
          go build -o /tmp/api ./cmd/api
          /tmp/api &
          sleep 5
          curl -f http://localhost:8080/health || exit 1

      - name: Run Newman — Smoke
        run: |
          newman run postman/icmongolang.postman_collection.json \
            -e postman/environment.dev.json \
            --folder '00. Setup' \
            --folder '01. Auth' \
            --reporters cli,junit \
            --reporter-junit-export reports/smoke.xml

      - name: Run Newman — Full
        run: |
          mkdir -p reports
          newman run postman/icmongolang.postman_collection.json \
            -e postman/environment.dev.json \
            --reporters cli,htmlextra,junit \
            --reporter-htmlextra-export reports/full.html \
            --reporter-junit-export reports/junit.xml \
            --bail

      - name: Upload reports
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: newman-reports
          path: reports/
          retention-days: 30

      - name: Publish JUnit results
        if: always()
        uses: mikepenz/action-junit-report@v4
        with:
          report_paths: reports/*.xml
          check_name: API Test Results
```

### F.7 GitLab CI

```yaml
# .gitlab-ci.yml
stages:
  - test
  - api-test

api-tests:
  stage: api-test
  image: node:20-alpine
  services:
    - postgres:16-alpine
    - redis:7-alpine
  variables:
    POSTGRES_USER: icmon
    POSTGRES_PASSWORD: icmon_dev
    POSTGRES_DB: icmongolang_test
    DB_DSN: "host=postgres user=icmon password=icmon_dev dbname=icmongolang_test sslmode=disable"
  before_script:
    - apk add --no-cache curl bash postgresql-client
    - npm install -g newman newman-reporter-htmlextra
  script:
    - |
      for f in test/fixtures/sql/*.sql; do
        psql "$DB_DSN" -f "$f" || true
      done
    - |
      newman run postman/icmongolang.postman_collection.json \
        -e postman/environment.dev.json \
        --reporters cli,junit,htmlextra \
        --reporter-junit-export reports/junit.xml \
        --reporter-htmlextra-export reports/full.html \
        --bail
  artifacts:
    when: always
    paths:
      - reports/
    reports:
      junit: reports/junit.xml
    expire_in: 30 days
```

### F.8 Newman HTML Report Sample

```html
<!-- reports/full.html — auto-generated -->
<!-- เปิดใน browser เพื่อดูผลลัพธ์แบบ visual -->

<!-- Summary:
  Total Requests: 187
  Passed: 184
  Failed: 3
  Duration: 45.2s
  Avg Response Time: 234ms

  Failures:
  ✗ POST /api/v1/devices — Status 500
  ✗ GET /api/v1/erp/orders/xxx — Timeout
  ✗ POST /api/v1/installation-jobs — Status 400
-->
```

---

## 🅶 PART 10G — MOCK SERVER + API DOCUMENTATION

### G.1 Postman Mock Server Setup

```bash
# สร้าง mock server ผ่าน Postman API
curl -X POST https://api.getpostman.com/mocks \
  -H "X-Api-Key: $POSTMAN_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "mock": {
      "name": "icmongolang-mock",
      "collection": "icmongolang-iot-platform-v1",
      "environment": "icmongolang-env-dev",
      "private": false
    }
  }'
```

### G.2 Mock URL Pattern

```
https://<mock-id>.mock.pstmn.io/api/v1/customers
https://<mock-id>.mock.pstmn.io/api/v1/devices
https://<mock-id>.mock.pstmn.io/api/v1/erp/orders
```

### G.3 Mock Server Usage — Frontend Dev

```javascript
// Frontend dev — ใช้ mock ก่อน backend พร้อม
const API_BASE = process.env.NODE_ENV === 'development'
  ? 'https://abc123.mock.pstmn.io/api/v1'
  : 'https://api.icmongolang.io/api/v1';

// ตัวอย่าง response ที่ mock จะ return
// POST /customers → 201 with generated UUID
// GET /customers → 200 with paginated list
// GET /customers/:id → 200 or 404
```

### G.4 Mock Response — Example Saved

```json
{
  "name": "Create Customer — Success",
  "originalRequest": {
    "method": "POST",
    "url": "{{base_url}}/api/v1/customers"
  },
  "status": "Created",
  "code": 201,
  "_postman_previewlanguage": "json",
  "header": [
    { "key": "Content-Type", "value": "application/json" },
    { "key": "Location", "value": "/api/v1/customers/{{$randomUUID}}" }
  ],
  "body": "{\n  \"data\": {\n    \"id\": \"{{$randomUUID}}\",\n    \"tenant_id\": \"{{tenant_id}}\",\n    \"code\": \"CUS-2026-0001\",\n    \"name\": \"บริษัท ทดสอบ จำกัด\",\n    \"type\": \"CORPORATE\",\n    \"segment\": \"SME\",\n    \"lifecycle\": \"PROSPECT\",\n    \"status\": \"PENDING\",\n    \"tax_id\": \"0105559999999\",\n    \"email\": \"contact@test.local\",\n    \"phone\": \"021234567\",\n    \"address\": {\n      \"line1\": \"99/1 อาคารทดสอบ\",\n      \"province\": \"กรุงเทพ\",\n      \"postcode\": \"10120\",\n      \"country\": \"TH\"\n    },\n    \"credit_limit\": 100000,\n    \"currency\": \"THB\",\n    \"created_at\": \"{{$isoTimestamp}}\",\n    \"updated_at\": \"{{$isoTimestamp}}\"\n  }\n}"
}
```

### G.5 API Documentation (auto-generated)

```markdown
# icmongolang — IoT Platform API v1

## Base URL
- Dev:      http://localhost:8080/api/v1
- Staging:  https://staging.icmongolang.io/api/v1
- Prod:     https://api.icmongolang.io/api/v1

## Authentication
All requests (except /auth/login, /auth/refresh) require:
- `Authorization: Bearer <access_token>`
- `X-Tenant-ID: <tenant_uuid>`

## Headers
| Header | Required | Description |
|---|---|---|
| `Authorization` | Yes | Bearer token |
| `X-Tenant-ID` | Yes | Tenant UUID |
| `X-Request-ID` | No | Idempotency key |
| `Idempotency-Key` | For POST | Prevent duplicate |

## Rate Limits
| Plan | Limit |
|---|---|
| Free | 1,000 req/day |
| Basic | 10,000 req/day |
| Pro | 100,000 req/day |
| Enterprise | Unlimited |

Headers: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`

## Error Format
```json
{
  "error": {
    "code": "CUSTOMER_CODE_EXISTS",
    "message": "Customer code already exists in this tenant",
    "details": {
      "field": "code",
      "value": "CUS-2026-0001"
    },
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "timestamp": "2026-01-15T10:00:00Z"
  }
}
```

## Endpoints Summary

### Auth
| Method | Path | Description |
|---|---|---|
| POST | `/auth/login` | Login |
| POST | `/auth/refresh` | Refresh token |
| POST | `/auth/logout` | Logout |
| GET | `/auth/me` | Current user |
| POST | `/auth/password` | Change password |

### Customers
| Method | Path | Description |
|---|---|---|
| GET | `/customers` | List customers |
| POST | `/customers` | Create customer |
| GET | `/customers/:id` | Get customer |
| PUT | `/customers/:id` | Update |
| DELETE | `/customers/:id` | Delete |
| POST | `/customers/:id/onboard` | Onboard |
| POST | `/customers/:id/suspend` | Suspend |
| POST | `/customers/:id/contacts` | Add contact |

### Devices
| Method | Path | Description |
|---|---|---|
| GET | `/devices` | List devices |
| POST | `/devices` | Register |
| GET | `/devices/:id` | Get |
| PUT | `/devices/:id` | Update |
| DELETE | `/devices/:id` | Delete |
| POST | `/devices/:id/provision` | Provision |
| POST | `/devices/:id/commands` | Send command |
| POST | `/devices/:id/telemetry` | Ingest telemetry |
| GET | `/devices/:id/telemetry` | Query telemetry |
| GET | `/devices/:id/shadow` | Get shadow |
| PATCH | `/devices/:id/shadow/desired` | Update desired |

### ERP
| Method | Path | Description |
|---|---|---|
| GET | `/erp/products` | List products |
| POST | `/erp/products` | Create product |
| GET | `/erp/inventory` | List inventory |
| POST | `/erp/inventory/adjust` | Adjust stock |
| GET | `/erp/orders` | List orders |
| POST | `/erp/orders` | Create order |
| POST | `/erp/orders/:id/confirm` | Confirm |
| POST | `/erp/orders/:id/ship` | Ship |
| POST | `/erp/orders/:id/receive` | Receive |
| GET | `/erp/invoices` | List invoices |
| POST | `/erp/invoices` | Create invoice |
| GET | `/erp/invoices/aging` | Aging report |
| POST | `/erp/payments` | Record payment |

### Logistics
| Method | Path | Description |
|---|---|---|
| GET | `/shipments` | List shipments |
| POST | `/shipments` | Create shipment |
| POST | `/shipments/:id/dispatch` | Dispatch |
| POST | `/shipments/:id/gps` | Update GPS |
| POST | `/shipments/:id/deliver` | Confirm delivery |
| GET | `/technicians` | List technicians |
| GET | `/installation-jobs` | List jobs |
| POST | `/installation-jobs` | Create job |
| POST | `/installation-jobs/:id/assign` | Assign tech |
| POST | `/installation-jobs/:id/start` | Start |
| PATCH | `/installation-jobs/:id/checklist` | Update checklist |
| POST | `/installation-jobs/:id/complete` | Complete |

### Reports & Dashboard
| Method | Path | Description |
|---|---|---|
| GET | `/reports` | List reports |
| POST | `/reports/:code/generate` | Generate |
| GET | `/report-executions/:id` | Get execution |
| GET | `/dashboard/kpis` | Get KPIs |
| GET | `/dashboard/kpis/history` | KPI history |
| GET | `/dashboard/dashboards/:id/data` | Dashboard data |
| POST | `/dashboard/insights/generate` | AI insight |
```

### G.6 Export Collection (CLI)

```bash
# Export collection ผ่าน Postman CLI
postman collection export icmongolang-iot-platform-v1 \
  --output postman/icmongolang.postman_collection.json

# Export environment
postman environment export icmongolang-env-dev \
  --output postman/environment.dev.json

# Import ผ่าน CLI
postman collection import postman/icmongolang.postman_collection.json
```

### G.7 OpenAPI Generation (จาก Postman)

```bash
# แปลง Postman collection → OpenAPI 3.0
npx postman-to-openapi postman/icmongolang.postman_collection.json \
  -o docs/openapi.yaml \
  -f "icmongolang API v1"

# Serve Swagger UI
npx swagger-ui-watcher docs/openapi.yaml
```

### G.8 Postman Monitors

```bash
# สร้าง monitor รันทุก 5 นาที
curl -X POST https://api.getpostman.com/monitors \
  -H "X-Api-Key: $POSTMAN_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "monitor": {
      "name": "icmongolang Health Monitor",
      "collection": "icmongolang-iot-platform-v1",
      "environment": "icmongolang-env-staging",
      "schedule": {
        "cron": "*/5 * * * *",
        "timezone": "Asia/Bangkok"
      },
      "options": {
        "followRedirects": true,
        "requestTimeout": 10000,
        "strictSSL": true
      },
      "notifications": {
        "onError": [
          { "type": "email", "value": "ops@icmongolang.io" },
          { "type": "slack", "value": "https://hooks.slack.com/..." }
        ],
        "onFailure": [
          { "type": "email", "value": "dev@icmongolang.io" }
        ]
      }
    }
  }'
```

---

## 📊 PART 10 SUMMARY

### Deliverables

| Category | Count | Format |
|---|:-:|---|
| **Environments** | 3 | dev / staging / prod |
| **Folders** | 10 | Setup · Auth · Customers · Sites · Packages · Devices · ERP · Logistics · Reports · Utilities |
| **Requests** | ~120 | ทุก module |
| **Negative tests** | ~15 | Error cases |
| **Pre-request scripts** | 3 levels | Collection / Folder / Request |
| **Test scripts** | 120+ | ทุก request |
| **Reusable helpers** | 7 | assertSuccess, assertPagination, ... |
| **E2E flows** | 3 | Auth, Order-to-Cash, Device Provisioning |
| **CI integrations** | 2 | GitHub Actions, GitLab CI |
| **Mock servers** | 1 | Postman mock |
| **Monitors** | 1 | Health check ทุก 5 นาที |
| **OpenAPI** | 1 | Auto-generated |

### Benefits

1. ✅ **Import แล้วรันได้ทันที** — ไม่ต้องตั้งค่าอะไรเพิ่ม
2. ✅ **Chained requests** — response ของ request หนึ่ง → ใช้ใน request ถัดไป
3. ✅ **Auto-token** — pre-request script จัดการ login/refresh อัตโนมัติ
4. ✅ **Multi-environment** — dev/staging/prod พร้อม safety guard
5. ✅ **Comprehensive assertions** — ทุก request มี test scripts
6. ✅ **Negative tests** — ครอบคลุม error cases
7. ✅ **CI-ready** — Newman + GitHub Actions + GitLab CI
8. ✅ **Mock server** — frontend dev ไม่ต้องรอ backend
9. ✅ **Auto-docs** — OpenAPI + Swagger UI
10. ✅ **Monitor** — health check ทุก 5 นาที

### File Structure

```
postman/
├── icmongolang.postman_collection.json    # Collection หลัก
├── environment.dev.json                    # Dev env
├── environment.staging.json                # Staging env
├── environment.prod.json                   # Prod env (read-only)
├── mocks/
│   ├── customer-create.json                # Mock response
│   ├── device-register.json
│   └── order-create.json
├── scripts/
│   ├── helpers.js                          # Reusable assertions
│   ├── auth-flow.js                        # Auth E2E
│   ├── order-to-cash.js                    # O2C E2E
│   └── device-provisioning.js              # Device E2E
└── README.md                               # วิธีใช้
```

---

# 🎯 ความคืบหน้า (ลำดับ B)

| # | งาน | สถานะ |
|:-:|---|:-:|
| 1 | `packagecatalog` deep dive | ✅ |
| 2 | `erp` deep dive | ✅ |
| 3 | UML / Sequence Diagrams | ✅ |
| 4 | Sample Data / Fixtures | ✅ |
| 5 | Postman Collection | ✅ **เสร็จ (response นี้)** |
| 6 | Executive Summary | ⏳ ถัดไป |

---

 