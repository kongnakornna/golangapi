# Fullschedule Swagger Drift Report

**Subject:** `/api` paths under the `fullschedule` module in `docs/swagger.yaml`
**Scope:** 21 paths, 35 operations (swagger lines ~6358–7663)
**Verdict:** **No drift.** Every documented operation matches the authoritative route registrations, security chain, handlers, and presenter DTOs. Two doc-quality observations are noted below (not drift).

---

## 1. Summary

- **21 paths** documented, all mapped to actual handlers in `routes.go`.
- **35 operations** (17 tagged `fullschedule`, 18 tagged `fullschedule-master`).
- **Security:** every operation carries `OAuth2Password` + `BearerAuth` and mounts under the middleware chain `mw.Verifier(true), mw.Authenticator(), mw.CurrentUser(), mw.ActiveUser()` (routes.go line 14). No operation is missing auth.
- **Tags** are consistent with the two route-group functions:
  - `fullschedule` → `MapFullscheduleRoutes` (routes.go line 54+)
  - `fullschedule-master` → `MapMasterRoutes` (routes.go lines 12–50)
- **Response refs** all resolve to actual presenter DTOs (`presenter.*`).
- **Base path:** `/api` (matches the served API prefix).

---

## 2. Methodology

- Compared each swagger path/operation block against:
  - `internal/modules/fullschedule/delivery/http/routes.go` (authoritative route map, 81 lines, EOF)
  - `internal/modules/fullschedule/delivery/http/handlers.go` and `master_handlers.go`
  - `internal/modules/fullschedule/presenter/presenters.go` (258 lines, EOF) and `presenter/master_presenters.go` (195 lines, EOF)
- Cross-checked generated mirrors `docs/docs.go` and `docs/swagger.json` for security/tag consistency (18 master-tag matches each).
- Path start lines (first operation headers): `/fullschedule` 6358, `/{id}` 6421, `/{id}/day-status` 6521, `/{id}/history` 6559, `/{id}/settings` 6599, `/{id}/status` 6668, `/{id}/trigger` 6706, `/areas` 6739, `/areas/{id}` 6810, `/areas/{id}/devices` 6910, `/devices` 6980, `/groups` 7015, `/groups/{id}` 7078, `/history` 7178, `/listschedulepage` 7218, `/report` 7302, `/scheduleall` 7337, `/scheduledevicepage` 7401, `/scope/preview` 7457, `/zones` 7497, `/zones/{id}` 7564.

> Swagger line numbers are snapshots at review time and may shift on regeneration. They are provided for navigation, not as a contract.

---

## 3. Path-by-Path Findings

### 3.1 `/fullschedule` (line 6358) — tag `fullschedule` ✅

- **GET** — list fullschedules with `page`/`per_page` → 200 `SuccessResponse-presenter_PaginatedScheduleResponse`, 401 `ErrorResponse`. Match: `ps.Get("/fullschedule", mh.ListSchedule())`.
- **POST** — create with required body `presenter.ScheduleCreate` → 200 `SuccessResponse-presenter_ScheduleResponse`, 400/401. Match: `ps.Post("/fullschedule", mh.CreateSchedule())`.
- Both consume/produce `application/json`.

### 3.2 `/fullschedule/{id}` (line 6421) — tag `fullschedule` ✅

- **DELETE** — soft delete → 200 `SuccessResponse-presenter_ScheduleResponse`, 400/404. Match: `ps.Delete("/fullschedule/:id", mh.DeleteSchedule())`.
- **GET** — get by id (`id` string required) → 200 `SuccessResponse-presenter_ScheduleResponse`, 400/404. Match: `ps.Get("/fullschedule/:id", mh.GetSchedule())`.
- **PUT** (line 6484) — update with required body `presenter.ScheduleUpdate` → 200 `SuccessResponse-presenter_ScheduleResponse`, 400/404. Match: `ps.Put("/fullschedule/:id", mh.UpdateSchedule())`.

### 3.3 `/fullschedule/{id}/day-status` (line 6521) — tag `fullschedule` ✅

- **PUT** — required body `presenter.ScheduleDayStatusChange` → 200 `SuccessResponse-presenter_ScheduleResponse`, 400/404. Description: “Toggle a single weekday flag on a fullschedule.” Method is PUT (not POST) and matches `mh.DayStatus()`.

### 3.4 `/fullschedule/{id}/history` (line 6559) — tag `fullschedule` ✅

- **GET** — `page`/`per_page` → 200 `SuccessResponse-presenter_PaginatedHistoryResponse`, 400/404. Match: `mh.History()`.

### 3.5 `/fullschedule/{id}/settings` (line 6599) — tag `fullschedule` ✅

- **GET** → 200 `SuccessResponse-array_presenter_SettingResponse`, 400/404. Match: `mh.Setting()`.
- **POST** — required body `presenter.ScheduleSettingRequest` → 200 `SuccessResponse-presenter_SettingResponse`, 400/404. Match: `mh.UpdateSetting()`.

### 3.6 `/fullschedule/{id}/status` (line 6668) — tag `fullschedule` ✅

- **PUT** — required body `presenter.ScheduleStatusChange` → 200 `SuccessResponse-presenter_ScheduleResponse`, 400/404. Method is PUT (not POST). Match: `mh.Status()`.

### 3.7 `/fullschedule/{id}/trigger` (line 6706) — tag `fullschedule` ✅

- **POST** — body `presenter.TriggerRequest` → 200 `SuccessResponse-...TriggerResponse`, 400/404. Description: “Manually trigger a fullschedule” (force execution; normal/full modes only). Match: `mh.Trigger()`.
- Doc observation: the 200 ref is fully qualified (`presenter_TriggerResponse`) rather than the short `presenter_*` form. See §4.

### 3.8 `/fullschedule/areas` (line 6739) — tag `fullschedule-master` ✅

- **GET** — `zone_id`/`sort`/`page`/`per_page` → 200 `SuccessResponse-presenter_PaginatedAreaResponse`, 401. Match: `mh.AreaList()`.
- **POST** — required body `presenter.AreaCreate` → 200 `SuccessResponse-presenter_AreaResponse`, 400. Match: `mh.AreaCreate()`.

### 3.9 `/fullschedule/areas/{id}` (line 6810) — tag `fullschedule-master` ✅

- **DELETE** — soft delete (only if no mapped devices) → 200 `SuccessResponse-presenter_AreaResponse`, 400/404. Match: `mh.AreaDelete()`.
- **GET** → 200 `SuccessResponse-presenter_AreaResponse`, 400/404. Match: `mh.AreaGet()`.
- **PUT** — body `presenter.AreaUpdate` → 200 `SuccessResponse-presenter_AreaResponse`, 400/404. Match: `mh.AreaUpdate()`.

### 3.10 `/fullschedule/areas/{id}/devices` (line 6910) — tag `fullschedule-master` ✅

- **GET** → 200 `SuccessResponse-array_presenter_AreaDeviceResponse`, 400/404. Match: `mh.AreaDevices()`.
- **POST** — required body `presenter.DeviceMapRequest` → 200 `SuccessResponse-presenter_MessageResponse`. Description: replace mapping; `device_ids` auto-resolves scope. Match: `mh.MapAreaDevices()`.

### 3.11 `/fullschedule/devices` (line 6980) — tag `fullschedule` ✅

- **GET** — `keyword`/`page`/`per_page` → 200 `SuccessResponse-array_presenter_IoTDeviceBrief`, 401. Description: list IoT devices (device picker). Match: `mh.Devices()`.
- Doc observation: documents only 200/401 (no 400). See §4.

### 3.12 `/fullschedule/groups` (line 7015) — tag `fullschedule-master` ✅

- **GET** — `page`/`per_page` → 200 `SuccessResponse-presenter_PaginatedGroupResponse`, 401. Match: `mh.GroupList()`.
- **POST** — body `presenter.GroupCreate` → 200 `SuccessResponse-presenter_GroupResponse`. Match: `mh.GroupCreate()`.

### 3.13 `/fullschedule/groups/{id}` (line 7078) — tag `fullschedule-master` ✅

- **DELETE** — soft delete (only if no child zones) → 200 `SuccessResponse-presenter_GroupResponse`, 400/404. Match: `mh.GroupDelete()`.
- **GET** → 200 `SuccessResponse-presenter_GroupResponse`, 400/404. Match: `mh.GroupGet()`.
- **PUT** — body `presenter.GroupUpdate` → 200 `SuccessResponse-presenter_GroupResponse`, 400/404. Match: `mh.GroupUpdate()`.

### 3.14 `/fullschedule/history` (line 7178) — tag `fullschedule` ✅

- **GET** — `page`/`per_page` + optional `status` enum `[success, failed, skipped, processing]` → 200 `SuccessResponse-presenter_PaginatedHistoryResponse`, 401. Description: get fullschedule history (all). Match: `mh.History()`.
- Mirrors `/api/settings/listschedulehistory` behavior.

### 3.15 `/fullschedule/listschedulepage` (line 7218) — tag `fullschedule` ✅

- **GET** — `page`/`pageSize`/`start`/`keyword` + `mode` enum `[normal, full, batch]` + `status` enum `[active, inactive, draft]` + `event_type`/`event_action`/`group_id`/`zone_id`/`area_id`/`sort` → 200 `SuccessResponse-...ListResult`, 401.
- Mirrors `/api/settings/listschedulepage`.
- Doc observation: 200 ref is fully qualified (`...presenter_ListResult`), matching the documented list-page response shape. See §4.

### 3.16 `/fullschedule/report` (line 7302) — tag `fullschedule` ✅

- **GET** — `from`/`to` (RFC3339)/`schedule_id` → 200 `SuccessResponse-array_presenter_ReportResponse`, 401. Match: `mh.Report()`. `ReportResponse` matches `array_presenter_ReportResponse`.

### 3.17 `/fullschedule/scheduleall` (line 7337) — tag `fullschedule` ✅

- **GET** — `start`/`keyword`/`mode`/`status`/`event_type`/`event_action`/`sort` (e.g. `start-DESC`) → 200 `SuccessResponse-array_presenter_ScheduleResponse`, 400/401. Mirrors `/api/settings/scheduleall`.

### 3.18 `/fullschedule/scheduledevicepage` (line 7401) — tag `fullschedule` ✅

- **GET** — `page`/`pageSize`/`keyword`/`schedule_id` (uuid)/`device_id` (integer)/`status`/`sort` → 200 `SuccessResponse-...ListResult`, 400/401. Mirrors `/api/settings/listscheduledevice`.
- Doc observation: 200 ref is fully qualified. See §4.

### 3.19 `/fullschedule/scope/preview` (line 7457) — tag `fullschedule-master` ✅

- **GET** — `group_id`/`zone_id`/`area_id` → 200 `SuccessResponse-presenter_ScopePreviewResponse`, 400/404. Description: preview device count for a scope. Match: `mh.ScopePreview()` (routes.go line 48, inside `MapMasterRoutes`).

### 3.20 `/fullschedule/zones` (line 7497) — tag `fullschedule-master` ✅

- **GET** — `group_id`/`page`/`per_page` → 200 `SuccessResponse-presenter_PaginatedZoneResponse`, 401. Match: `mh.ZoneList()`.
- **POST** — required body `presenter.ZoneCreate` → 200 `SuccessResponse-presenter_ZoneResponse`, 400/401. Match: `mh.ZoneCreate()`.

### 3.21 `/fullschedule/zones/{id}` (line 7564) — tag `fullschedule-master` ✅

- **DELETE** — soft delete (only if no child areas) → 200 `SuccessResponse-presenter_ZoneResponse`, 400/404. Match: `mh.ZoneDelete()`.
- **GET** → 200 `SuccessResponse-presenter_ZoneResponse`, 400/404. Match: `mh.ZoneGet()`.
- **PUT** — body `presenter.ZoneUpdate` (all-optional) → 200 `SuccessResponse-presenter_ZoneResponse`, 400/404. Match: `mh.ZoneUpdate()`.

---

## 4. Doc-Quality Observations (NOT drift)

These are cosmetic/informational and do not indicate a behavioral mismatch. They require no code change.

1. **Fully-qualified 200 refs** — three operations reference the shared success envelope by fully-qualified name while the rest use the short `presenter_*` form:
   - `POST /fullschedule/{id}/trigger` → `...presenter_TriggerResponse`
   - `GET /fullschedule/listschedulepage` → `...presenter_ListResult`
   - `GET /fullschedule/scheduledevicepage` → `...presenter_ListResult`
   - All three still resolve to real definitions; only the spelling differs.

2. **Missing 400 on list endpoints** — GET `/fullschedule/areas`, `/fullschedule/groups`, `/fullschedule/zones`, and `/fullschedule/devices` document only `200` + `401` (no `400`), unlike the sibling list endpoints (`/fullschedule` GET, `/fullschedule/history`, `/fullschedule/report`) which document `400`. This is an inconsistency in documented responses, not a handler gap.

3. **Method emphasis** — the dynamic sub-resources `/{id}/day-status` and `/{id}/status` are correctly documented as **PUT**; when scanning the file, POST appears on adjacent endpoints, so watch for copy/paste confusion while reading.

---

## 5. Notes & Limitations

- Report generated from a static read of the four cited files; no runtime probe was performed.
- Line numbers refer to the current `docs/swagger.yaml` and will move after regeneration.
- Un-scoped swagger content after line 7663 (e.g. `/influx/devicechart`) was not reviewed.
- All 401 responses map to the shared `httpErrors.ErrResponse` (`msg`/`status`/`statusText`) envelope, consistent with the codebase-wide error convention.