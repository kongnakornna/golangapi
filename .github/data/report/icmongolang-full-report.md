# ICMON Go — Full Project Report

> **Version:** 1.0.0
> **Module:** `icmongolang`
> **Go Version:** 1.26
> **License:** github.com/kongnakornna
> **Date:** 2026-08-14
> **Status:** Active development — see [PROJECT_CHECKLIST.md](PROJECT_CHECKLIST.md) for the reusable verification checklist.

---

## Table of Contents

1. [Project Overview](#1-project-overview)
2. [System Architecture](#2-system-architecture)
3. [Tech Stack & Dependencies](#3-tech-stack--dependencies)
4. [Configuration](#4-configuration)
5. [Entry Points & CLI Commands](#5-entry-points--cli-commands)
6. [Directory Structure](#6-directory-structure)
7. [Docker Infrastructure](#7-docker-infrastructure)
8. [Module Breakdown (24 modules)](#8-module-breakdown-24-modules)
9. [API Routes](#9-api-routes)
10. [Middleware & Security](#10-middleware--security)
11. [Monitoring & Metrics](#11-monitoring--metrics)
12. [Database & Migrations](#12-database--migrations)
13. [Testing Status](#13-testing-status)
14. [Known Issues & Findings](#14-known-issues--findings)
15. [Setup & Deployment](#15-setup--deployment)

---

## 1. Project Overview

**ICMON** (Industrial IoT Control & Monitoring) is a Go-based backend system providing:

- Device monitoring and IoT data collection via **MQTT**, **InfluxDB** (time-series), and **PostgreSQL**
- Real-time data streaming via **WebSocket**
- Event-driven processing via **Apache Kafka**
- Background task processing via **Asynq** (Redis-backed task queue)
- User management with RBAC roles
- JWT-based authentication (access + refresh tokens, refresh persisted in Redis)
- Email verification & password reset via SMTP with HTML templates
- PDF report generation (HTML templates + headless Chrome, Thai baht text conversion)
- LLM embedding + semantic vector search (Ollama/OpenAI-compatible via `pkg/llm`)
- Vector database abstraction (`pkg/vectordb`) with **pgvector** (default) and **Elasticsearch** backends
- Prometheus metrics + Grafana dashboards + exporters for every subsystem

**Primary use cases:**
- Industrial IoT sensor data collection and alerting
- Repair job work-orders (ใบรับงานซ่อม), parts, services, quotations
- Purchase orders, payments, receipts
- Online order system (WOS)
- Dashboard analytics

---

## 2. System Architecture

3-layer (plus delivery) architecture per module:

```
Delivery (HTTP handlers) → Usecase → Repository → DB
                                   ↘ helpers/pkg (shared infra)
```

- **Router:** `internal/server/handlers.go` is the active registration point (`cmd/serve.go`), mounting a sub-router at `/api`.
- **Alternate router:** `internal/delivery/rest/router.go` exists but is NOT wired into the serve command (legacy; mounts report + a subset of modules).
- **Generic CRUD pattern:** most modules embed `internal.UseCaseI[M]` → `Create`, `Get`, `GetMulti(limit, offset)`, `Delete`, `Update(id, values map)`.
- **Standalone services:** Kafka order service (`cmd/kafka`), WebSocket service (`cmd/websocket`), Ollama exporter (`cmd/ollama-exporter`), AI chat service (`ai/`).

### Data flow (IoT path)

```
MQTT Broker (mosquitto) ─► pkg/mqtt ─► iot/mqtt module ─► Redis cache
                                    ├► InfluxDB (time-series)
                                    ├► PostgreSQL (device/alarm)
                                    └► WebSocket hub (real-time clients)
```

### Vector search path

```
POST /api/vectordata/documents → usecase → llm.Embed → vectordb (pgvector | elasticsearch)
POST /api/vectordata/search    → usecase → llm.Embed(query) → KNN search → ranked results
```

---

## 3. Tech Stack & Dependencies

| Category | Tech |
|---|---|
| Language | Go 1.26 |
| Router | `chi` |
| ORM | `gorm` |
| Config | `viper` |
| CLI | `cobra` |
| Validation | go-playground/validator |
| Auth | `jwt` (HS256/RS256), bcrypt |
| Logging | `zap` |
| Email | `gomail` + `hermes` |
| Task queue | `asynq` (Redis) |
| Message queue | Redis-backed internal queue (`internal/modules/queue`) |
| Kafka | `sarama` |
| MQTT | `paho` |
| Time-series | InfluxDB 2.x client |
| Vector DB | `pgvector` (PostgreSQL ext) + Elasticsearch dense_vector |
| LLM | OpenAI-compatible client (Ollama `nomic-embed-text` default) |
| PDF | HTML templates + headless Chrome (`chromedp`) |
| Hot-reload | `air` (`.\air.cmd` on Windows 11) |

---

## 4. Configuration

`config/config.default.yml` (overridable via `config.dev.yml` and env vars; env keys derive from struct field names, e.g. `VectorDB` → `VECTOR_DB_PROVIDER`).

| Block | Key settings |
|---|---|
| `server` | AppVersion, Port, Mode, Timezone (Asia/Bangkok), BaseUrl, MigrateOnStart |
| `jwt` | SecretKey, Issuer, access/refresh expiry, private/public keys |
| `firstSuperUser` | Email, Name, Password |
| `logger` | Encoding, Level |
| `postgres` | Host, Port, User, Password, Dbname, ConnectionTimeout |
| `redis` | Addr, Password, Db, pool settings |
| `taskRedis` | Asynq task queue Redis |
| `email` / `smtpEmail` | From, Link, SMTP host/user/pass, TLS/SSL |
| `mqtt` | broker, client_id, username, password, qos |
| `influxdb` | url, token, org, bucket, timeout |
| `kafka` | brokers, topic, group_id, timeout |
| `elasticsearch` | addresses, username, password, index, timeout |
| `llm` | base_url, api_key, model (nomic-embed-text), timeout, dims |
| `vectorDb` | provider (`pgvector` \| `elasticsearch`), index, dims |

---

## 5. Entry Points & CLI Commands

Cobra root command (`cmd/root.go`, Use: `icmongolang`):

| Command | File | Purpose |
|---|---|---|
| `serve` | `cmd/serve.go` | Start HTTP API server (connect Postgres/Redis, optional migrate, start server) |
| `migrate` | `cmd/migrate.go` | GORM AutoMigrate over ~130 models, create extensions (`uuid-ossp`, `pgcrypto`, `vector`), ensure pgvector index |
| `worker` | `cmd/worker.go` | Start Asynq task processor (queues: critical/default) |
| `initdata` | `cmd/initdata.go` | Create super-user if not exists |
| `kafka` | `cmd/kafka/` | Standalone Kafka order service (HTTP :5051, topic `icmon-events`) |
| `vectordata seed` | `cmd/vectordata.go` | Seed sample vector docs (requires LLM service) |

Standalone `main()` binaries: `cmd/api/main.go` (main API), `cmd/websocket/main.go` (WebSocket service, :8080), `cmd/ollama-exporter/main.go` (Ollama Prometheus exporter :9309).

---

## 6. Directory Structure

| Path | Contents |
|---|---|
| `internal/modules/` | 24 business modules (see §8) |
| `internal/models/` | GORM models (12 files) |
| `internal/server/` | Main router + wiring (`handlers.go`) |
| `internal/middleware/` | Auth (Verifier/Authenticator/CurrentUser/ActiveUser/SuperUser), CORS, rate limit, monitor, logger |
| `internal/delivery/rest/` | Legacy alternate router (not wired) |
| `pkg/` | 22 shared infra packages (see §2 exploration) |
| `cmd/` | cobra commands + standalone services |
| `config/` | YAML configs |
| `db/` + `migrations/` | SQL dumps & schema/seed SQL |
| `ai/` | Standalone Go AI chat service (port 8001) |
| `elk/` | ELK stack assets (elk.ps1, logstash pipeline) |
| `monitoring/` | Prometheus config, Grafana provisioning + dashboards |
| `mqtt/` | Mosquitto broker config |
| `postman/` | Postman collection |
| `scripts/` | Version helper |
| `linux/` | Linux build script |
| `test/` | Non-Go review/checklist docs |
| `.gitlab-ci/`, `.github/` | CI/CD pipelines & docs/skills |

---

## 7. Docker Infrastructure

`docker-compose.yml` services (see full inventory in exploration notes):

| Service | Image | Port |
|---|---|---|
| api | Dockerfile.dev (migrate + initdata + air serve) | 5000 |
| worker | Dockerfile.dev (`main.go worker`) | — |
| db | `pgvector/pgvector:pg15` | 5435 |
| redis | `redis:7-alpine` | 6579 (internal) |
| kafka | `apache/kafka:3.7.1` (KRaft) | 9092 |
| kafka-exporter | `danielqsj/kafka-exporter` | 9308 (internal) |
| prometheus | `prom/prometheus:latest` | 9090 |
| redis-exporter | `oliver006/redis_exporter` | 9121 (internal) |
| postgres-exporter | `prometheuscommunity/postgres-exporter` | 9187 (internal) |
| blackbox-exporter | `prom/blackbox-exporter` | 9115 (internal) |
| elasticsearch | `docker.elastic.co/elasticsearch:8.13.4` | 9200 |
| elasticsearch-exporter | `prometheuscommunity/elasticsearch-exporter` | 9114 (internal) |
| logstash | `logstash:8.13.4` | 5001, 5044, 9600 |
| kibana | `kibana:8.13.4` | 5601 |
| mqtt | `eclipse-mosquitto:2` | 1885 |
| mosquitto-exporter | `sapcc/mosquitto-exporter` | 9234 (internal) |
| influxdb | `influxdb:2.7` | 9087 |
| ollama | `ollama/ollama:latest` | 11434 |
| ollama-exporter | Dockerfile.ollama-exporter | 9309 (internal) |
| ai | build ./ai | 8001 |
| grafana | `grafana/grafana:latest` | 3000 |

Named volumes for each stateful service.

---

## 8. Module Breakdown (24 modules)

All JWT-protected unless noted. "Public" = no auth middleware; "rate-limited" = whole group rate-limited.

| Module | Purpose | Key routes |
|---|---|---|
| **alarm** | Sensor alarm status validation + localized messages (default/en/th) | `POST /api/alarm/validate[ /en /th]` |
| **auth** | Auth facade (login, refresh, logout, password reset) → users module | Public: `POST /api/auth/login`, `/signin`, `GET /publickey`, `/verifyemail`, `POST /forgotpassword`, `PATCH /resetpassword`; JWT(refresh): `GET /refresh`, `/logout`, `/logoutall` |
| **batch** | Batch job scheduling/management + logs | `POST|GET /api/batch/jobs`, `GET|PUT|DELETE /api/batch/jobs/{id}`, `POST /{id}/run`, `GET /{id}/logs` |
| **customer** | Customer + car (vehicle) management | `/api/customer` CRUD, `/api/car` CRUD |
| **dashboard** | Analytics stats, revenue, top parts, job status | `GET /api/dashboard/stats`, `/revenue`, `/top-parts`, `/job-status` |
| **document** | File upload/download/delete | `POST|GET /api/documents`, `GET|DELETE /{id}` |
| **elasticsearch** | LLM-embed + ES vector search (direct, non-abstracted) | Public `GET /api/elasticsearch/health`; JWT `POST /create-index`, `/embeddings`, `/vectors`, `/search` |
| **email** | Send email, email logs, config | `POST /api/email/send`, `GET /logs`, `GET|PUT /config` |
| **i18n** | Translation strings by key/locale | `/api/i18n/translations` CRUD |
| **influxdb** | Time-series write/query/charts/stats | `POST /api/influx/write`, `GET /query`, `POST /devicechart`, `/filters`, `/statistics` |
| **iot** | Device mgmt + MQTT v3.1.1 telemetry, charts, export | ~17 public GET routes (`/topic`, `/device`, `/sensercharts`, ...) + JWT `POST /control`, `PUT /devicestatus`, `/updatedeviceconfig`, `DELETE /devicedatacleanup` (whole group rate-limited) |
| **items** | Item/inventory parts | `/api/item/` CRUD |
| **job** | Repair work-orders, status/history, services, parts, PDFs | `/api/job/` CRUD + `PUT /{id}/status`, `GET /{id}/history`, `/report`, `/pdf`, `/picking/pdf`, `/delivery/pdf`, services/parts sub-resources |
| **kafka** | Async order processing + WS broadcast (standalone service) | `POST /orders` (202), `GET /ws`, `GET /health` (port 5051) |
| **mqtt** | MQTT publish/subscribe, live topic data, device control | Public `GET /api/mqtt/subscriptions`, `/status`, `/gettopicdata`, `/devicecontrol`; JWT `POST /publish`, `/subscribe`, `/unsubscribe` (rate-limited) |
| **payment** | Payments, receipts, refunds, cancellation, PDFs | `/api/payments` (search, outstanding, history, refund, cancel, invoice), `/api/receipts` (get, pdf, cancel) |
| **purchaseorder** | PO lifecycle: create-with-details, from-quotation, send/confirm/receive/cancel, history, PDF | `/api/purchase-orders` CRUD + `/from-quotation/{qid}`, `/suggestions/{jobId}`, `/{id}/send|confirm|receive|cancel|pdf|history` |
| **queue** | Redis-backed message queue (pub/sub, delayed, retries, DLQ) — no HTTP routes | (infra only) |
| **quotation** | Quotation CRUD + PDF | `/api/quotation/` CRUD + `/{id}/pdf` |
| **report** | PDF reports (daily sales, inventory, customer list, invoice, credit/debit notes) — no usecase layer | `GET /reports/.../pdf` (mounted on ROOT, no `/api` prefix) |
| **users** | User accounts, auth core, passwords, sessions, super-user bootstrap | Public `POST /api/register`, `/users`, `/signin`, `/login`; JWT `GET|PUT /api/user/me`, `PATCH /me/updatepass`; Admin(SuperUser) user/role/password mgmt |
| **vectordata** | Abstract vector-DB CRUD + embeddings + semantic search (pgvector/ES) | Public `GET /api/vectordata/health`; JWT `POST /create-index`, `/documents`, `/search`, `/seed`, `GET|PUT|DELETE /documents[/{id}]` |
| **websocket** | Real-time hub, rooms, message history | Main API: `GET /api/ws`; Standalone: `GET /ws`, `POST|GET /api/ws/messages`, `GET /rooms`, `GET /rooms/{room}/stats`, `GET /health` |
| **wos** | Web Order System (online orders) | `POST|GET /api/wos/orders`, `GET /{id}`, `PUT /{id}/status` |

---

## 9. API Routes

Global middleware on main router: Recoverer, RequestID, RealIP, URLFormat, Monitoring (Prometheus), Logger, Timeout, JSON render, CORS.

| Prefix | Module |
|---|---|
| `/api` | main sub-router mount |
| `/api/influx` | InfluxDB (only if client present) |
| `/api/elasticsearch` | ES vector (only if es+llm present) |
| `/api/vectordata` | Vector DB abstraction |
| `/api/iot` | MQTT v3 IoT |
| `/api/mqtt` | MQTT control |
| `/api/alarm` | Alarm validation |
| `/api/auth` | Auth |
| `/api/user` | Users |
| `/api/item` | Items |
| `/api/purchase-orders` | Purchase orders |
| `/api/payments`, `/api/receipts` | Payments & receipts |
| `/api/job` | Repair jobs |
| `/api/customer`, `/api/car` | Customers & cars |
| `/api/quotation` | Quotations |
| `/api/ws` | WebSocket |
| `/api/ping`, `/api/health` | Liveness & health |
| `/api/dashboard` | Dashboard (root-mounted) |
| `/api/documents` | Documents (root-mounted) |
| `/api/email` | Email (root-mounted) |
| `/api/batch` | Batch (root-mounted) |
| `/api/i18n` | i18n (root-mounted) |
| `/api/wos` | WOS (root-mounted) |
| `/apimetric`, `/metrics` | Prometheus metrics |
| `/health` | root liveness |
| `/swagger/*` | Swagger UI |

> Note: `/reports/...` (report module) is only wired in the legacy router, NOT in the active server.

### `/api/health` checks
`database`, `redis`, `influxdb`, `elasticsearch`, `vectordata`, `llm`, `mqtt` (requires connected), `websocket` (hardcoded true).

---

## 10. Middleware & Security

- **JWT auth:** `mw.Verifier` / `Authenticator` / `CurrentUser` / `ActiveUser` / `SuperUser`; access token + refresh token (refresh stored in Redis); refresh routes use `Verifier(false)`.
- **RBAC roles:** SUPERADMIN, ADMIN, EDITOR, MONITOR, USER.
- **Password:** bcrypt via `pkg/cryptpass`; crypto-secure random (`pkg/secureRandom`).
- **CORS**, **RequestID**, **Recoverer**, **timeout**, **per-IP rate limiting** (iot/mqtt groups).
- **Response envelope:** `pkg/responses` → `{data, error, is_success}`; errors centralized in `pkg/httpErrors`.

---

## 11. Monitoring & Metrics

- Prometheus scrape config in `monitoring/prometheus.yml`.
- Exporters: kafka, redis, postgres, blackbox, elasticsearch, mosquitto, ollama (`:9309`, custom Go exporter).
- Grafana: provisioned dashboards for golang, redis, postgres, mqtt, kafka, influxdb, elasticsearch, ollama, websocket.
- ELK: logstash pipeline (tcp/beats → ES) + Kibana; `elk/elk.ps1` to start/stop.
- Custom JSON metrics at `/apimetric` + Prometheus native `/metrics`.

---

## 12. Database & Migrations

- **Postgres** (`pgvector/pgvector:pg15`, port 5435, db `icmongolang`) with extensions: `uuid-ossp`, `pgcrypto`, `vector`.
- **GORM AutoMigrate** (~130 models) via `cmd/migrate` with a skip-list for problematic models.
- SQL migrations in `migrations/`: `db.sql`, `icmon.sql`, `public.sql`, `public_db.sql`, `sd_user_access_menu.sql`, kafka/websocket tables, `20260712_new_modules_schema.sql`, `20260712_seed_data.sql`.
- **Redis**: caches (user, refresh tokens) + Asynq task queue + internal message queue.
- **InfluxDB 2.x** (bucket, org) for time-series sensor data.

---

## 13. Testing Status

| Metric | Value |
|---|---|
| `*_test.go` files (excl. vendor) | **6** |
| Go source files | 362 |
| Models files | 12 |

**Coverage is minimal** — most modules have no automated tests; the vectordata module has zero tests (finding NIT-9 in `scrutinize-vectordata.md`). Verification so far has been manual/live via curl + psql + CLI.

---

## 14. Known Issues & Findings

### From `docs/report/scrutinize-vectordata.md` (2026-08-14) — verdict **fix-then-ship**
1. **BLOCKER-1:** ES provider ignores `VECTOR_DB_INDEX` — `esBackend.index` dead code; all ops use `cfg.Elasticsearch.Index` (`pkg/elasticsearch/client.go:58`). Data always lands in `vector_embeddings`.
2. **MAJOR-2:** Partial update broken — metadata-only PUT → 500 (`vecToString(nil)` = `"[]"`); content-only PUT wipes document_id/source_type/metadata.
3. **MAJOR-3:** Cross-provider error semantics diverge (pg missing-id UPDATE/DELETE → 200; ES → 500); `handler.Get` maps every error → 404.
4. **MAJOR-4:** ES backend drops `metadata` (no `Metadata` field on `esclient.VectorDoc`).
5. **MAJOR-5:** Seed not idempotent (duplicates accumulate).
6. **MINOR:** migrate unconditionally creates pgvector extension; no `k`/`limit` cap; `Create` requires `content` even with embedding; no automated tests.

### Upstream limitation (documented in README_VectorData.md)
- Ollama `nomic-embed-text` returns **identical embeddings for pure-Thai text** — a model/tokenizer limitation, not a code bug. English/Thai-with-Latin produce distinct vectors.

### Structural
- `internal/delivery/rest/router.go` is legacy/unwired; `/reports/...` PDF routes only exist there.
- Mounting inconsistency: some modules hardcode `/api/...` and are mounted on root vs. the `/api` sub-router (works, but inconsistent).

---

## 15. Setup & Deployment

**Local dev (Docker):**
```powershell
docker compose up -d               # start infra
docker compose up api worker       # API (air hot-reload) + worker
```
**API:** http://localhost:5000 · **Swagger:** http://localhost:5000/swagger/ · **Grafana:** :3000 · **Kibana:** :5601

**Super-user bootstrap:** `initdata` (default `admin@icmongolang.local` / `admin1234`).

**Database direct:**
```powershell
docker exec icmongolang-db-1 psql -U postgres -d icmongolang -p 5435 -c "<sql>"
```

**Linux build:** `linux/build.sh` · **Windows:** `build.ps1`.
**CI/CD:** `.gitlab-ci/` (build, build-tag, apply-host, deploy-dev/uat/prod per service).
