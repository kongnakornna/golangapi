# PDPA Module Documentation

Folder and file structure for `internal/modules/pdpa`.

## Tree

```
pdpa/
├── application/
│   └── dto/                       (empty)
├── domain/
│   ├── entity/
│   │   ├── audit_trail.go
│   │   ├── consent_log.go
│   │   ├── dsar_request.go
│   │   ├── pdpa_policy.go
│   │   ├── pdpa_purpose.go
│   │   ├── pdpa_request_response.go
│   │   └── user_account_status.go
│   ├── errors/
│   │   └── errors.go
│   ├── event/
│   │   └── event.go
│   ├── repository/
│   │   ├── account_status_repo.go
│   │   ├── audit_repo.go
│   │   ├── cache.go
│   │   ├── consent_repo.go
│   │   └── dsar_repo.go
│   ├── service/
│   │   ├── blockchain_service.go
│   │   └── deletion_policy_service.go
│   └── valueobject/
│       ├── account_status.go
│       ├── consent_purpose.go
│       ├── consent_status.go
│       ├── dsar_status.go
│       └── dsar_type.go
├── infrastructure/
│   ├── es/                        (empty)
│   ├── kafka/                     (empty)
│   ├── localization/
│   │   └── locales/               (empty)
│   ├── messaging/                 (empty)
│   ├── persistence/
│   │   └── models/                (empty)
│   ├── redis/                     (empty)
│   ├── scheduler/                 (empty)
│   └── search/                    (empty)
├── interfaces/
│   ├── handlers/                  (empty)
│   ├── http/                      (empty)
│   ├── localization/
│   │   └── locales/               (empty)
│   ├── routes/                    (empty)
│   ├── websocket/                 (empty)
│   └── ws/                        (empty)
└── repository/                    (empty)
```

## Summary

- **Total folders:** 5 top-level (`application`, `domain`, `infrastructure`, `interfaces`, `repository`)
- **Files found:** 21 Go files, all under `domain/`:
  - `entity/`: 7
  - `errors/`: 1
  - `event/`: 1
  - `repository/`: 5
  - `service/`: 2
  - `valueobject/`: 5
- **Empty folders:** `application/dto`, `repository`, and all `infrastructure/*` and `interfaces/*` subfolders
- **Layer coverage:** `domain` is fully populated; `application`, `infrastructure`, and `interfaces` are scaffolding-only