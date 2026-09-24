# ICMON Project — Reusable Checklist

> Reusable verification & review checklist for the ICMON Go project.
> Use for: smoke-testing the full stack, reviewing a new/modified module, verifying fixes, and pre-deployment checks.
> Companion report: [icmongolang-full-report.md](icmongolang-full-report.md).

---

## A. Full-Project Smoke Test

Run after any infra change or before claiming "everything works".

- [ ] `docker compose up -d` starts clean; `docker compose ps` shows all services `Up` (db, redis, kafka, mqtt, influxdb, elasticsearch, ollama, prometheus, grafana, kibana, api, worker)
- [ ] `GET http://localhost:5000/health` → `{"status":"ok"}`
- [ ] `GET http://localhost:5000/api/ping` → ok + timestamp
- [ ] `GET http://localhost:5000/api/health` → all checks `true` (database, redis, influxdb, elasticsearch, vectordata, llm, mqtt, websocket)
- [ ] `GET /metrics` and `GET /apimetric` respond 200
- [ ] Login works: `POST /api/auth/login` with `admin@icmongolang.local` / `admin1234` (body uses `username`) → token saved to `$env:TEMP\token.txt`
- [ ] `GET /api/user/me` with token → returns user
- [ ] `POST /api/vectordata/health` path (public) → 200
- [ ] Swagger UI loads: `GET /swagger/`
- [ ] Worker consumes tasks: `docker compose logs worker` shows no errors
- [ ] No GORM/`AutoMigrate` errors in `docker compose logs api`

### Database
- [ ] psql reachable: `docker exec icmongolang-db-1 psql -U postgres -d icmongolang -p 5435 -c "select version();"`
- [ ] Extensions exist: `uuid-ossp`, `pgcrypto`, `vector`
- [ ] `vector_embeddings` table + vector index exist

### Subsystem checks
- [ ] MQTT: `GET /api/mqtt/status` → connected
- [ ] ES: `GET /api/elasticsearch/health` and `GET /api/vectordata/health` → ok
- [ ] Ollama: `GET http://localhost:11434/api/tags` lists `nomic-embed-text`
- [ ] Grafana `:3000`, Kibana `:5601`, Prometheus `:9090` respond
- [ ] InfluxDB `:9087` healthy

---

## B. Module Review Checklist (new or changed module)

Reusable — generalize the `scrutinize-vectordata` method to any module.

### B.1 Wiring & routing
- [ ] Routes registered in `internal/server/handlers.go` (active router), not only legacy `internal/delivery/rest/router.go`
- [ ] Correct prefix: mounted on `/api` sub-router OR hardcoded `/api/...` consistently (note the existing inconsistency)
- [ ] Health flag added to `/api/health` if the module needs a connection
- [ ] Public vs JWT vs SuperUser middleware correctly applied (no auth leak on protected routes)

### B.2 CRUD semantics
- [ ] `Create`: validation of required fields (via `pkg/utils` validator)
- [ ] `Update`: **partial update** — does NOT wipe untouched fields; preserves values when only some fields sent
- [ ] `Get`/`Update`/`Delete` on missing id → correct error (see B.3), not silent 200
- [ ] `GetMulti`: pagination (limit/offset) respected; limit has an upper bound
- [ ] Cross-provider/backend behaviors are identical (abstraction layer)

### B.3 Error semantics
- [ ] Missing resource → sentinel `ErrNotFound` → handler maps to **404**
- [ ] Bad input / DB constraint (e.g. invalid UUID) → mapped to 400/404 appropriately, not leaking SQLSTATE/internal detail
- [ ] No `RowsAffected` unchecked in `Update`/`Delete` (SQL backends)
- [ ] Nil/empty values don't produce invalid SQL (e.g. empty vector `'[]'::vector` → 500)

### B.4 Data integrity
- [ ] Seed/init functions **idempotent** (no duplicates on re-run)
- [ ] Metadata/extra fields not silently dropped at any backend
- [ ] No dead config (config read but never used; or used value differs from config)
- [ ] Foreign keys / related resources cleaned up or properly constrained

### B.5 Code quality
- [ ] `go build ./...`, `go vet ./...`, `gofmt` clean
- [ ] Automated tests exist for the module's critical paths (currently only 6 test files in repo)
- [ ] Secrets never logged or committed; no hardcoded creds

---

## C. Vector Data Fixes — Verification Checklist

Track the pending findings from `scrutinize-vectordata.md` (verdict: **fix-then-ship**).

### C.1 BLOCKER-1 — ES uses configured index
- [x] `pkg/vectordb/es_backend.go` no longer has dead `index` field; ops use the configured index
- [x] `pkg/elasticsearch` methods take an index param (or client built from `cfg.VectorDB.Index`)
- [x] **Live test:** `VECTOR_DB_PROVIDER=elasticsearch VECTOR_DB_INDEX=vector_documents_test` CLI seed → data lands in `vector_documents_test`, NOT `vector_embeddings`
- [x] README `§3.2` claim corrected

### C.2 MAJOR-2 — Partial update safe
- [x] `usecase.Update` merges with existing doc: content-only PUT keeps document_id/source_type/metadata
- [x] Metadata-only PUT does NOT hit `vector must have at least 1 dimension` (500)
- [x] **Live test:** PUT `{content:"x"}` then GET → fields intact; PUT `{metadata:{...}}` → 200

### C.3 MAJOR-3 — Error semantics
- [x] pgvector UPDATE/DELETE on missing id → 404 (not 200)
- [x] `GET /documents/{bad-id}` → clean 404, no SQLSTATE 22P02 leak
- [x] Both providers return the same status codes for the same scenario

### C.4 MAJOR-4 — ES metadata
- [ ] `esclient.VectorDoc` has `Metadata` field; index/serialize on create, return on get/search

### C.5 MAJOR-5 — Idempotent seed
- [ ] Run `vectordata seed` twice → no duplicate docs
- [ ] README TC17 behavior holds (seeds 5 sample docs)

### C.6 Minors (optional)
- [ ] `cmd/migrate.go` skips `CREATE EXTENSION "vector"` if provider is ES (or made config-aware)
- [ ] `k`/`limit` capped (e.g. ≤ 100)
- [ ] `Create` accepts embedding-only (content optional when embedding provided)
- [ ] `*_test.go` added under `pkg/vectordb` and `internal/modules/vectordata`

---

## D. Pre-Deployment Checklist

- [ ] `go build ./...`, `go vet ./...`, `gofmt -l .` clean
- [ ] All tests pass: `go test ./...` (excl. vendor)
- [ ] `migrate` ran successfully on target DB; extensions present
- [ ] `/api/health` all green on target environment
- [ ] Smoke test §A passed against target
- [ ] Postman collection (`postman/icmongolang.postman_collection.json`) updated with new endpoints
- [ ] README_* manual updated (env vars, troubleshooting)
- [ ] Version bumped via `scripts/update_version.py`
- [ ] Linux build via `linux/build.sh` succeeds
- [ ] `.gitlab-ci` deploy jobs cover the changed service
- [ ] No secrets in committed config/`.env`; no debug logs at `info`+ leaking data

---

## E. Known Project Quirks (read before touching)

- Report PDFs: 6 routes under `/reports/*/pdf` now run through `middleware.VerifierForReport()` — `Authorization: Bearer <access-token>` (scripted) OR `?token=<download-token>` (plain browser). Download tokens are minted by `POST /reports/token` (full access token required), HS256 with the dedicated `JWT_DOWNLOAD_TOKEN_SECRET` (≥32 bytes), TTL 5 min, scoped per report type + `source` for invoice, re-checked per handler via `downloadScopeOK`. Do NOT reuse the RS256 access keypair for these (a download token would validate as an access token elsewhere).
- SPA client pattern for PDFs (no header-in-URL): `fetch('/reports/invoice/pdf?source=<id>', { headers: { Authorization: 'Bearer ' + token } }).then(r => r.blob()).then(b => window.open(URL.createObjectURL(b)))`.
- `go test ./...` runs with `-mod=vendor`; only `stretchr/testify/assert` is vendored (no `require`) — use `assert` in new tests.
- Legacy router `internal/delivery/rest/router.go` is NOT wired; `/reports/...` is wired in the ACTIVE router `internal/server/handlers.go` (root-mounted, same as document/email/wos). Keep both call sites of `CreateReportHandler` in sync (legacy + active).
- `cmd/migrate.go` (cobra `migrate`) only AutoMigrates a fixed model list — it does NOT create `m_customer`/`t_quotation`/etc. Apply `migrations/20260712_new_modules_schema.sql` + `20260712_seed_data.sql` manually (the seed file uses invalid UUIDs for parts/services rows and fails; insert test data directly).
- Report PDFs need Chromium in the image: `Dockerfile.dev` installs `chromium` + `font-noto-thai` + `font-noto-cjk`. Template `init()` uses `//go:embed` + `template.ParseFS` rooted at `base.html` (each doc overrides `{{define "content"}}`); do NOT revert to `template.New(name).ParseFiles` with relative paths (panics on non-root CWD, and leaves the root template empty).
- Auto-numbering triggers (`generate_job_no`, `generate_customer_code`, `generate_quotation_no`, `generate_po_no`) read the sequence with `SPLIT_PART(code, '-', 3)` — do NOT revert to `SUBSTRING(code FROM N)` (it swallows the dash and yields negative sequence numbers, e.g. `CUST-2026-0000`). Defined in `migrations/20260712_new_modules_schema.sql`, `db/public_DB.sql`, `docs/docs/modules_3.md`.
- Env var convention: field `VectorDB` → `VECTOR_DB_PROVIDER` / `VECTOR_DB_INDEX` / `VECTOR_DB_DIMS` (NOT `VECTORDB_*`).
- Ollama `nomic-embed-text` returns identical embeddings for pure-Thai text (upstream limitation).
- Postgres port 5435, Redis internal 6579, MQTT 1885, InfluxDB 9087 — non-default on purpose.
- `air` on Windows 11 may be blocked by AppLocker → use `.\air.cmd` / `go run github.com/air-verse/air@latest`.
