# Scrutiny Report — `internal/modules/report` + `pkg/report`

Date: 2026-08-14
Reviewer: opencode (scrutinize skill)
Scope: `internal/modules/report/{handler.go,delivery/http}` and `pkg/report/{generator.go,models.go,bahtthai.go}`.

> **Status update (same day):** wired into the live server per user decision (option "1+2+3" — feasible wiring + defer rest + Chrome in Docker). BLOCKER-1 and MAJOR-2 are resolved; MAJOR-1 partially resolved (2 of 6 report types now backed by real data). Details below.

## Intent

Generate PDF documents (daily sales, inventory, customer list, invoice, credit/debit note) from HTML templates via headless Chrome (chromedp), exposed under `/reports/*/pdf`.

Simpler-alternative pass: **this module is not reachable from any running server.** A repo-wide grep for `MapReportRoute` / `CreateReportHandler` across `cmd/` and `internal/` returns zero call sites. The first question is whether this should exist at all — either wire it (with real data) or delete it. Fixing anything inside an unwired module is premature.

## Findings (severity order)

### BLOCKER-1 — Dead code: no wiring, no route is reachable — ✅ FIXED
- **Finding:** `MapReportRoute` (routes.go:11) and `CreateReportHandler` (handlers.go:22) have no callers anywhere in the build. The active router (`internal/server/handlers.go`) never registers `/reports`.
- **Why it matters:** six endpoints advertised in Swagger comments that return 404 in practice. If a consumer is told reports exist, they hit a wall.
- **Evidence:** `Select-String -Path cmd/**,internal/** -Pattern "MapReportRoute|CreateReportHandler"` → no matches.
- **Fix applied:** wired into `internal/server/handlers.go` via `MapReportRoute(r, ...)` (mounted on the root router, same as document/email/wos) with a new report usecase; `CreateReportHandler(cfg, logger, customerUC, reportUC)` now takes the customer + quotation usecases.

### MAJOR-1 — Placeholder data only; no database access — ✅ PARTIAL
- **Finding:** every handler builds an empty dataset: `Sales: []`, `TotalRev: 0`, `Items: []` (handlers.go:53-59, 84-89, 114-119), a hardcoded company `"ICMON Auto Repair"` (handlers.go:226-233), and invoice documents hardcoded to `"INV-2026-001"` / `"ลูกค้า"` (handlers.go:145-149, 174-181, 205-212).
- **Why it matters:** even if wired, the generated PDFs contain no real data — they are static placeholders, so the feature is non-functional regardless of wiring.
- **Evidence:** `DailySalesPDF` never queries sales; `date` is the only input and it only shifts the timestamp.
- **Fix applied (partial):** `CustomerListPDF` now lists real customers via `customerUC.GetMultiByUserID` + `Count`; `InvoicePDF` takes `?source=<quotation-uuid>` and renders the quotation (parts + services + customer) via the new `report/usecase.BuildInvoice`. The remaining four (daily sales, inventory, credit note, debit note) stay documented placeholders — no sales/inventory/credit entities exist in the data model yet (`items` is generic `{Id,Title,Description,OwnerId}`; no invoice/sales tables). Company info centralized in `internal/modules/report/company.go`.

### MAJOR-2 — Template loading panics on any non-root CWD — ✅ FIXED
- **Finding:** `pkg/report/generator.go:32-43` `init()` runs `template.Must(template.New(name).ParseFiles("pkg/report/templates/base.html", ...))` with **relative paths**. Running the compiled binary from any other working directory → template parse fails → panic at process start.
- **Why it matters:** any deployment that does not `cd` into the project root (e.g. systemd, Docker WORKDIR mismatch, a packaged binary) crashes at boot. Latent today only because the module is unwired.
- **Evidence:** path is relative; nothing chdirs or embeds.
- **Fix applied:** `//go:embed templates/*.html` + `template.ParseFS`; root template is `base.html` so each document's `{{define "content"}}` overrides the base block. Also fixed the dead double-assigned `fileURL` (now `url.URL{Scheme:"file", Path:"/"+ToSlash(absPath)}.String()`). Note: `template.New(name)` would produce an empty root — the fix changed execution to render via the `base.html` root.

### NIT-1 — chromedp runtime dependency and fragile file URL — ✅ ADDRESSED
- **Finding:** `GeneratePDF` spawns headless Chrome with `chromedp.Flag("no-sandbox", true)` (generator.go:99-129); the container image must ship Chrome or every PDF fails at runtime. Also `fileURL` is built twice — line 95 (`"file://" + url.PathEscape(absPath)`) is dead, immediately overwritten at line 97 (`"file:///" + absPath`) without escaping spaces/backslashes.
- **Why it matters:** brittle PDF pipeline; the dead line suggests the URL handling was never actually exercised.
- **Fix applied:** `fileURL` uses the escaped `url.URL` form; `Dockerfile.dev` now installs `chromium` + `fontconfig` + `font-noto-thai` + `font-noto-cjk` + `ttf-freefont` so Thai PDFs render in the container.

### NIT-2 — Invalid `date` query param silently ignored — ✅ FIXED
- **Finding:** `time.Parse("2006-01-02", dateStr)` error is swallowed; bad dates fall back to `time.Now()` (handlers.go:45-48).
- **Why it matters:** a client typing a malformed date gets today's report without any indication the filter was dropped.
- **Fix applied:** malformed `date` now returns a 400.

### NIT-3 — Auth + PDF is unusable from a plain browser
- **Finding:** routes are JWT-guarded (routes.go:14-17), but a browser `<a href="/reports/invoice/pdf">` cannot attach an `Authorization` header — the PDF link dead-ends for humans.
- **Why it matters:** the endpoints are designed for browser PDFs but only usable by scripted clients with tokens.
- **Suggested change:** if wired, serve via a session-aware flow or an authenticated download endpoint that sets the header client-side. **Still open.**
- **RESOLVED (2026-08-14):** implemented all three options.
  1. **Documented the fetch+blob pattern** for SPA clients: `fetch('/reports/invoice/pdf?source=<id>', { headers: { Authorization: 'Bearer <access-token>' } })` → `response.blob()` → `URL.createObjectURL()` → open in a new tab. No `Authorization` header ever appears in the URL.
  2. **`?token=` support:** the 6 PDF routes now run through the new `middleware.VerifierForReport()` (internal/middleware/jwtauth.go), which accepts the `Authorization` header as before OR a short-lived `?token=` query param — so `<a href="/reports/invoice/pdf?source=<id>&token=<download-token>">` works directly in a browser.
  3. **Short-lived download-token minting:** `POST /reports/token` (full access token required — download tokens cannot mint more tokens) returns a 5-minute, report-scoped HS256 token signed with a dedicated secret (`JWT_DOWNLOAD_TOKEN_SECRET`, ≥32 bytes; NOT the RS256 access keypair, so it can never validate as an access token elsewhere). Tokens carry `report_type` (+ `source` for the invoice report) and each PDF handler re-checks scope via `downloadScopeOK` (internal/modules/report/delivery/http/handlers.go) — a token minted for one report cannot fetch another. Primitives: `jwt.CreateDownloadTokenHS256`/`ParseDownloadTokenHS256` (pkg/jwt/token.go) with unit tests (pkg/jwt/token_test.go).

## Verification notes

- Confirmed zero call sites for wiring; `pkg/report` compiles (templates resolved during `go build`'s `init` only when CWD is the repo root — the build itself succeeded, consistent with the CWD finding).
- `bahtthai.go` (amount-in-Thai-text) is self-contained and reusable — the only part worth keeping regardless of the module decision.
- **Live verification (2026-08-14):** after applying `migrations/20260712_new_modules_schema.sql` (the `migrate` cobra command AutoMigrates only a fixed model list and does NOT create `m_customer`/`t_quotation`), `GET /reports/customer-list/pdf` and `GET /reports/invoice/pdf?source=<uuid>` both return valid `%PDF-` documents. Text extraction confirms the company header, invoice no `QT-2026-0001`, 3 line items (9000 + 5000 + 15000), tax 1960, net 29960, and both customers with codes/phones/spent.

## Verdict

**rework** — the module is unwired and non-functional-as-written (no data, CWD-dependent templates). The single biggest reason: fix nothing until a product decision is made on whether `/reports/*/pdf` should exist; if yes, wire it with real data and `go:embed` templates.

> **Status update:** per the product decision the module is now wired (BLOCKER-1), templates are embedded (MAJOR-2), and customer-list + invoice-from-quotation are backed by real data (MAJOR-1 partial). NIT-3 (browser-accessible PDF auth) is implemented via `?token=` download tokens + fetch+blob docs. Remaining: daily-sales/inventory/credit/debit notes await sales/inventory/credit entities.
