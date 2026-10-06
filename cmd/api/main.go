// Package main — entry point ของ ICMON API server
//
// ============================================================================
// Swagger / OpenAPI metadata
// ============================================================================
//
// @title                      Go IoT API
// @version                    1.0
// @description                API for ICMON project — IoT monitoring, PDPA compliance, order management
// @description
// @description                ---
// @description                ## 🚀 Quick Start
// @description                ---
// @description                **Demo Credentials** (seeded via `go run cmd/initdata/main.go`):
// @description                ```
// @description                Email:    demo@example.com
// @description                Username: demouser
// @description                Password: Demo@1234
// @description                ```
// @description
// @description                **How to authenticate:**
// @description                1. Call `POST /api/auth/signin` (or `/auth/login`) with demo credentials
// @description                2. Copy `access_token` from the response
// @description                3. Click the 🔓 **Authorize** button at the top
// @description                4. Paste `Bearer <access_token>` (include the word "Bearer")
// @description                5. All protected endpoints will now use the token
// @description
// @description                **Admin Demo** (created from env `FIRST_SUPER_USER_EMAIL`):
// @description                - Superuser can access all `/api/user/*` admin endpoints
// @description
// @description                ---
// @description                ## 📦 Modules
// @description                ---
// @description                - **IoT**: devices, sensor data, alarm, MQTT control
// @description                - **PDPA**: consent, DSAR, account lifecycle, policy, audit
// @description                - **Business**: orders, payments, purchase orders, quotations
// @description                - **Platform**: notifier, auditlog, apimanager, realtime, flowengine
// @termsOfService             https://icmon.example.com/terms
//
// @contact.name               ICMON Support
// @contact.email              support@icmon.example.com
// @contact.url                https://icmon.example.com/support
//
// @license.name               MIT
// @license.url                https://opensource.org/licenses/MIT
//
// @host                       localhost:5000
// @BasePath                   /api
// @schemes                    http https
//
// ============================================================================
// Security definitions
// ============================================================================
//
// @securityDefinitions.oauth2.password  OAuth2Password
// @tokenUrl                             /api/auth/login
// @scope.read                           Grants read access
// @scope.write                          Grants write access
//
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer " followed by your access token. Example: "Bearer eyJhbGciOiJSUzI1NiIs..."
//
// ============================================================================
// Tag descriptions (organize Swagger UI)
// ============================================================================
//
// --- Core ---
// @tag.name         auth
// @tag.description  Authentication — login, refresh, logout, register
//
// @tag.name         users
// @tag.description  User management — CRUD, profile, roles, permissions
//
// @tag.name         items
// @tag.description  Items CRUD (sample module)
//
// --- IoT ---
// @tag.name         iot
// @tag.description  IoT operations — devices, sensor data, alarm status, device control
//
// @tag.name         mqtt
// @tag.description  MQTT operations — publish, subscribe, device control, live topic data
//
// @tag.name         alarm
// @tag.description  Alarm validation — calculate alarm level, title, subject, content
//
// @tag.name         influxdb
// @tag.description  InfluxDB time-series queries — raw data, statistics, write points
//
// --- PDPA ---
// @tag.name         pdpa-consent
// @tag.description  PDPA Consent — record, revoke, list history
//
// @tag.name         pdpa-dsar
// @tag.description  Data Subject Access Requests — submit, list, process, reject
//
// @tag.name         pdpa-account
// @tag.description  User account lifecycle — suspend, terminate, confirm deletion (admin only)
//
// @tag.name         pdpa-policy
// @tag.description  Privacy policy — active version, list versions, publish new (admin only)
//
// @tag.name         pdpa-admin
// @tag.description  PDPA admin — compliance reports
//
// --- API Keys ---
// @tag.name         api-keys
// @tag.description  API Key management — create, list, rotate, revoke (manager/admin only)
//
// --- Business ---
// @tag.name         Orders
// @tag.description  Async order processing via Kafka + WebSocket
//
// @tag.name         Websocket
// @tag.description  Real-time bidirectional communication for orders
//
// @tag.name         customers
// @tag.description  Customer management
//
// @tag.name         cars
// @tag.description  Car management (per customer)
//
// @tag.name         job
// @tag.description  Repair job cards — services, parts, status, report
//
// @tag.name         quotation
// @tag.description  Quotations — create, update, delete, approve
//
// @tag.name         purchase-orders
// @tag.description  Purchase orders — create, confirm, receive, cancel
//
// @tag.name         payments
// @tag.description  Payment recording, refund, cancel, search
//
// @tag.name         receipts
// @tag.description  Receipt generation, cancel, PDF
//
// @tag.name         dashboard
// @tag.description  Business dashboard — stats, revenue, job status, top parts
//
// @tag.name         batch
// @tag.description  Batch job management — create, run, logs
//
// @tag.name         wos
// @tag.description  Web Order System — orders, status update
//
// @tag.name         Reports
// @tag.description  PDF report generation — invoice, receipt, daily sales, inventory
//
// --- Schedule / Control ---
// @tag.name         fullschedule
// @tag.description  Full schedule module — normal/full/batch modes, history, report
//
// @tag.name         fullschedule-master
// @tag.description  Schedule master data — groups, zones, areas, device mapping
//
// @tag.name         control
// @tag.description  Auto-control manager — execute device control, run schedules
//
// --- Platform / Infra ---
// @tag.name         settings
// @tag.description  Settings module — NestJS-compatible endpoints (devices, mqtt, schedule, etc.)
//
// @tag.name         notifier
// @tag.description  Notification dispatcher — email, SMS, Line, Discord, IO
//
// @tag.name         auditlog
// @tag.description  Audit log with hash-chain (tamper-evidence)
//
// @tag.name         apimanager
// @tag.description  API key management with rate-limiting (legacy module)
//
// @tag.name         realtime
// @tag.description  Real-time event broadcasting to WebSocket rooms
//
// @tag.name         flowengine
// @tag.description  Node-RED-like workflow engine — nodes, edges, triggers
//
// @tag.name         elasticsearch
// @tag.description  Elasticsearch vector search — embed, index, semantic search
//
// @tag.name         vectordata
// @tag.description  Abstract vector DB — pgvector or Elasticsearch backend
//
// @tag.name         documents
// @tag.description  Document management — upload, list, download, delete
//
// @tag.name         email
// @tag.description  Email service — send, config, logs
//
// @tag.name         i18n
// @tag.description  Internationalization — translations by locale
//
// @tag.name         websocket
// @tag.description  WebSocket HTTP endpoints — messages history, search, presence, retention
//
// @tag.name         Jobs
// @tag.description  Job PDF generation — job card, delivery sheet, picking list
//
// @tag.name         Quotation
// @tag.description  Quotation PDF generation
//
// ============================================================================
// External docs
// ============================================================================
//
// @externalDocs.description   ICMON API Documentation
// @externalDocs.url           https://icmon.example.com/docs
package main

import (
	"log"

	"icmongolang/cmd"
	"icmongolang/config"
	"icmongolang/pkg/helpers"
)

func main() {
	// ---------- Load config ----------
	v, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(v)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	// ---------- Inject timezone from config ----------
	helpers.SetDefaultTimezone(cfg.Server.Timezone)

	// ---------- Execute cobra command ----------
	// RootCmd + subcommands (serve, migrate, worker, kafka) อยู่ที่ cmd/
	cmd.Execute()
}