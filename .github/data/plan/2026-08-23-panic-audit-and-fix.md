# Plan: Panic Audit & Fix — Strict No-Panic Policy

## ภาพรวมการแก้ไข

| # | Layer | ไฟล์ | สิ่งที่แก้ |
|---|-------|------|-----------|
| 1 | middleware | `internal/middleware/prometheus.go` | register metric fail → log warning + skip (ไม่ panic) |
| 2 | cmd | `cmd/root.go` | ลบ `panic(err)` dead code 2 จุด |
| 3 | delivery/http | `internal/modules/iot/delivery/http/routes.go` | `MapMQTT3Routes` เปลี่ยน signature → `return error`, nil guard เป็น error |
| 4 | delivery/http | `internal/modules/mqtt/delivery/http/routes.go` | `MapMQTTRoutes` เปลี่ยน signature → `return error` + validate ก่อน map + แยก helper ให้ test ได้ |
| 5 | server wiring | `internal/server/handlers.go` | caller 2 จุดรับ error จาก mapper → wrap `%w` return |
| 6 | server | `internal/server/server.go` | `ListenAndServe` fail → log + `os.Exit(1)`, `Shutdown` fail → log เท่านั้น |
| 7 | transaction | `pkg/transaction/manager.go` | recover 3 จุด → rollback + return error (named return), ไม่ re-panic |
| 8 | docs | `docs/PANIC_AUDIT.md` | เอกสาร audit: inventory 16 จุด, policy, before/after, การป้องกันซ้ำ |
| 9 | tests | `pkg/transaction/manager_test.go`, `internal/modules/{mqtt,iot}/delivery/http/routes_test.go` | unit test ตาม design |

**Policy:** non-test code ห้ามมี `panic()` — runtime error → return error, startup fatal → `log.Fatal` / `log.Errorf + os.Exit(1)`, observability fail → log + continue

**Unit Tests: เขียน** (decision จาก tb-brainstorming)

---

## Section 1: `internal/middleware/prometheus.go`

**Location:** `internal/middleware/prometheus.go:40-55` (func `init()`)

**From → To:** ตอนนี้ถ้า `prometheus.Register` fail ด้วย error ที่ไม่ใช่ `AlreadyRegisteredError` → `panic(err)` → service crash ตอน boot ทั้งที่ metric เสียไม่กระทบ business logic → เปลี่ยนเป็น log warning แล้ว skip metric นั้น (init() return error ไม่ได้)

**Change detail** — เพิ่ม import `"log"` แล้วแก้ closure ใน `init()`:

```go
func init() {
	// Register metrics and ignore "already registered" errors
	mustRegister := func(c prometheus.Collector) {
		if err := prometheus.Register(c); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				log.Printf("⚠️ prometheus: failed to register %s: %v", metricName(c), err)
			}
		}
	}
	// ... mustRegister(...) calls unchanged
}

// metricName คืนชื่อ collector เพื่อใช้ใน log
func metricName(c prometheus.Collector) string {
	if d, ok := c.(interface{ Describe(chan<- *prometheus.Desc) }); ok {
		ch := make(chan *prometheus.Desc, 1)
		go func() {
			defer close(ch)
			d.Describe(ch)
		}()
		for desc := range ch {
			return desc.String()
		}
	}
	return "<unknown>"
}
```

หมายเหตุ: หาก `metricName` ซับซ้อนเกินไป ใช้ `log.Printf("⚠️ prometheus register failed: %v", err)` อย่างเดียวก็พอ (เลือกแบบเรียบได้ implementer discretion — ห้าม panic ทั้งสองแบบ)

---

## Section 2: `cmd/root.go`

**Location:** `cmd/root.go:56-68` (func `initConfig`)

**From → To:** `log.Fatalf` เรียก `os.Exit(1)` อยู่แล้ว ทำให้ `panic(err)` บรรทัดถัดไปเป็น **dead code ที่ไม่มีวันทำงาน** → ลบ panic 2 จุดทิ้ง

**Change detail:**

```go
func initConfig() {
	cfgViper, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	_, err = config.ParseConfig(cfgViper)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}
}
```

---

## Section 3: `internal/modules/iot/delivery/http/routes.go`

**Location:** `internal/modules/iot/delivery/http/routes.go:9-18` (func `MapMQTT3Routes`)

**From → To:** nil guard 3 จุดเป็น `panic(...)` (และ message สะกดผิดเป็น `"MapMQTTRoutes"` ทั้งที่ function ชื่อ `MapMQTT3Routes`) → signature เป็น `error` + `errors.New` + แก้ message ใช้ชื่อ function ถูกต้อง

**Change detail:**

```go
import (
	"errors"

	"icmongolang/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func MapMQTT3Routes(router chi.Router, h *MQTT3Handler, mw *middleware.MiddlewareManager) error {
	if router == nil {
		return errors.New("MapMQTT3Routes: router is nil")
	}
	if h == nil {
		return errors.New("MapMQTT3Routes: handler is nil")
	}
	if mw == nil {
		return errors.New("MapMQTT3Routes: middleware manager is nil")
	}
	router.Route("/iot", func(r chi.Router) {
		// ... body เดิมทั้งหมด ไม่แก้
	})
	return nil
}
```

---

## Section 4: `internal/modules/mqtt/delivery/http/routes.go`

**Location:** `internal/modules/mqtt/delivery/http/routes.go:11-50`

**From → To:** (a) nil guard 4 จุดเป็น panic → return error (b) check middleware ที่ `mw.Verifier()` ฯลฯ คืน nil กระทำ **หลัง** map route บางส่วนไปแล้ว (line 37-39) → ย้าย validate ให้ครบ **ก่อน** แตะ router (c) package-level `sync.Once` กลืน error ไม่ได้ → แยก unexported helper `mapMQTTRoutesOnce(...) error` ให้ `once.Do` capture error และ test เรียก helper ตรง ๆ ได้

**Change detail:**

```go
var once sync.Once

func MapMQTTRoutes(router chi.Router, h *MQTTHandler, mw *middleware.MiddlewareManager) error {
	var err error
	once.Do(func() { err = mapMQTTRoutesOnce(router, h, mw) })
	return err
}

func mapMQTTRoutesOnce(router chi.Router, h *MQTTHandler, mw *middleware.MiddlewareManager) error {
	if router == nil {
		return errors.New("mapMQTTRoutesOnce: router is nil")
	}
	if h == nil {
		return errors.New("mapMQTTRoutesOnce: handler is nil")
	}
	if mw == nil {
		return errors.New("mapMQTTRoutesOnce: middleware manager is nil")
	}

	rateLimit := mw.RateLimit()
	verifier := mw.Verifier(true)
	authenticator := mw.Authenticator()
	currentUser := mw.CurrentUser()
	activeUser := mw.ActiveUser()
	if rateLimit == nil || verifier == nil || authenticator == nil || currentUser == nil || activeUser == nil {
		return errors.New("mapMQTTRoutesOnce: one or more middleware functions returned nil")
	}

	router.Route("/mqtt", func(r chi.Router) {
		r.Use(rateLimit)
		r.Get("/subscriptions", h.Subscriptions)
		r.Get("/status", h.Status)
		r.Get("/gettopicdata", h.GetTopicData)
		r.Get("/devicecontrol", h.DeviceControl)

		r.Group(func(r chi.Router) {
			r.Use(verifier, authenticator, currentUser, activeUser)
			r.Post("/publish", h.Publish)
			r.Post("/subscribe", h.Subscribe)
			r.Post("/unsubscribe", h.Unsubscribe)
		})
	})
	return nil
}
```

เพิ่ม import `"errors"`. Comment block curl ท้ายไฟล์คงเดิม

---

## Section 5: `internal/server/handlers.go` (caller wiring)

**Location:** `internal/server/handlers.go:281-282` และ `304-305` (enclosing func `New` คืน `(*chi.Mux, error)` อยู่แล้ว — line 108)

**From → To:** call site เดิม ignore ผลลัพธ์ (void) → รับ error แล้ว wrap `%w` return ขึ้นไป (`NewServer` → `cmd serve` จะ fail-fast ตอน boot)

**Change detail:**

```go
	iotHandler := iotHttp.NewMQTT3Handler(iotUC, logger, redisCache)
	if err := iotHttp.MapMQTT3Routes(apiRouter, iotHandler, mw); err != nil {
		return nil, fmt.Errorf("register IoT MQTT3 routes: %w", err)
	}
	logger.Info("✅ MQTT3 routes registered")
```

```go
		} else {
			mqttHandler := mqttHttp.CreateMQTTHandler(mqttUC, logger, redisCache)
			if err := mqttHttp.MapMQTTRoutes(apiRouter, mqttHandler, mw); err != nil {
				return nil, fmt.Errorf("register MQTT routes: %w", err)
			}
			logger.Info("✅ MQTT routes registered with WebSocket, alarm log, cache & InfluxDB")
		}
```

(`fmt` import ตรวจว่ามีอยู่แล้วใน handlers.go — ถ้าไม่มีให้เพิ่ม)

---

## Section 6: `internal/server/server.go`

**Location:** `internal/server/server.go:155-185` (func `Start`)

**From → To:** (a) `ListenAndServe` fail → `panic(err)` crash โดยไม่ log ผ่าน logger และไม่ผ่าน graceful path → `logger.Errorf` + `os.Exit(1)` (startup fatal) (b) `Shutdown` fail → panic ทำให้ cleanup MQTT/InfluxDB ไม่รัน → `logger.Errorf` เท่านั้น แล้ว cleanup ต่อ

**Change detail:**

```go
func (srv *Server) Start() {
	srv.logger.Info("✅ starting server...")
	go func() {
		srv.logger.Infof("✅ Listening on %s", srv.server.Addr)
		if err := srv.server.ListenAndServe(); err != http.ErrServerClosed {
			srv.logger.Errorf("❌ server error: %v", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit
	srv.logger.Infof("✅ Shutting down server... Reason: %s", sig)

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	if err := srv.server.Shutdown(ctx); err != nil {
		srv.logger.Errorf("❌ server shutdown error: %v", err)
	}

	// ... MQTT/Influx cleanup เดิม (line 176-183) ไม่แก้
	srv.logger.Info("✅ Server gracefully stopped")
}
```

`os` import มีอยู่แล้ว (line 7)

---

## Section 7: `pkg/transaction/manager.go`

**Location:** 3 recover blocks — `ExecuteWithOptions` (53-58), `RunInTransaction` (84-89), `WithTransaction` (254-259)

**From → To:** recover → rollback → **re-panic** ทำให้ panic ใน tx กระจายขึ้นไปถึง chi Recoverer (500 ทั่วไป + stack trace) → rollback + **return error** `fmt.Errorf("panic in transaction: %v", r)` — caller เห็น error ตาม flow ปกติ ต้องใช้ named return

### 7.1 `ExecuteWithOptions` (line 39-75)

```go
func (m *GormTransactionManager) ExecuteWithOptions(ctx context.Context, opts *sql.TxOptions, fn TxFunc) (err error) {
	tx := m.db.WithContext(ctx)
	if opts != nil {
		tx = tx.Begin(opts)
	} else {
		tx = tx.Begin()
	}

	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("panic in transaction: %v", r)
		}
	}()

	if fnErr := fn(ctx, tx); fnErr != nil {
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %v (original error: %w)", rbErr, fnErr)
		}
		return fnErr
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
```

หมายเหตุ: inner `err` ของ Commit ยัง shadow ได้ปกติ (block scope) แต่เพื่อความชัดเจนใช้ `if commitErr := ...` ก็ได้ — implementer เลือกได้ ขอแค่ named return `(err error)` ตรง signature

### 7.2 `RunInTransaction` (line 78-97)

```go
func RunInTransaction(db *gorm.DB, fn func(*gorm.DB) error) (err error) {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("panic in transaction: %v", r)
		}
	}()

	if fnErr := fn(tx); fnErr != nil {
		tx.Rollback()
		return fnErr
	}

	return tx.Commit().Error
}
```

### 7.3 `WithTransaction` (line 248-263)

```go
func WithTransaction(ctx context.Context, db *gorm.DB, fn func(context.Context, *gorm.DB) error) (err error) {
	tc, tcErr := NewTransactionContext(ctx, db)
	if tcErr != nil {
		return tcErr
	}

	defer func() {
		if r := recover(); r != nil {
			tc.Rollback()
			err = fmt.Errorf("panic in transaction: %v", r)
		}
	}()

	return tc.Complete(fn(tc.Context(), tc.DB()))
}
```

**Behavior change:** code path ที่เคย panic ใน tx แล้วพึ่ง chi Recoverer → ตอนนี้ caller ได้ error ต้อง handle เอง — ครอบคลุมโดย policy strict ที่ user approve

---

## Section 8: `docs/PANIC_AUDIT.md` (เอกสารใหม่)

ภาษา Thai + EN technical โครงสร้าง:

```markdown
# Panic Audit & Fix — icmongolang

วันที่: 2026-08-23 | Policy: strict no-panic (non-test code)

## สรุปผู้บริหาร
- พบ panic 16 จุด / 6 ไฟล์ (ไม่รวม vendor/) — แก้ครบทั้งหมด
- Policy ใหม่: runtime error → return error / startup fatal → log.Fatal หรือ
  log.Errorf + os.Exit(1) / observability fail → log + continue

## Inventory
| ไฟล์ | Line เดิม | ประเภท | Risk | วิธีแก้ | สถานะ |
|---|---|---|---|---|---|
| cmd/root.go | 61, 67 | dead code (log.Fatalf exit ก่อน panic) | low | ลบ | ✅ |
| internal/server/server.go | 160 | crash-risk (ListenAndServe fail) | high | log + os.Exit(1) | ✅ |
| internal/server/server.go | 173 | crash-risk (Shutdown fail ข้าม cleanup) | medium | log เท่านั้น | ✅ |
| internal/middleware/prometheus.go | 45 | boot-crash (init register fail) | medium | log + skip | ✅ |
| internal/modules/iot/delivery/http/routes.go | 11, 14, 17 | nil-guard | low | return error | ✅ |
| internal/modules/mqtt/delivery/http/routes.go | 16, 19, 22, 38 | nil-guard | low | return error | ✅ |
| pkg/transaction/manager.go | 56, 87, 257 | intentional re-panic | high | rollback + return error | ✅ |

## Policy
เกณฑ์เขียนโค้ดใหม่:
1. ห้าม panic() ใน non-test code ทุกกรณี
2. Runtime error → return error wrap ด้วย %w
3. Startup fatal (config/port/db จำเป็น) → log.Fatalf หรือ logger.Errorf + os.Exit(1)
4. Observability (metrics/log sink) ล่ม → log warning + continue ห้าม kill service
5. Panic ใน transaction → rollback + return error (ไม่ re-panic) — caller handle เอง

## รายละเอียดรายจุด (Before/After)
(แต่ละไฟล์แปะ code before → after สั้น ๆ + เหตุผล — copy จาก plan sections 1-7)

## Behavior Change Notes
- panic ใน transaction เดิมกระจายถึง chi middleware.Recoverer (router.go:69,
  handlers.go:169) → 500 generic — ตอนนี้เป็น error ที่ caller เห็น
- MapMQTT3Routes เดิม panic message สะกดผิด "MapMQTTRoutes" → แก้เป็นชื่อถูกต้อง

## การป้องกันซ้ำ
- แนะนำ golangci-lint forbidigo rule:
  linters-settings:
    forbidigo:
      forbid:
        - pattern: 'panic\('
          msg: "use error return / log.Fatal instead (see docs/PANIC_AUDIT.md)"
  (exclude _test.go)
- Review checklist: PR ที่เพิ่ม panic ต้อง justify ใน description
```

---

## Unit Tests (เขียน)

Convention ตาม repo เดิม: pure unit test ไม่ต้องมี live DB (`database/sql` fake driver + `gorm.io/driver/postgres.Config{Conn}` — Initialize ไม่ query เมื่อ Conn set ตรวจแล้วที่ vendor postgres.go:87-88)

### 9.1 `pkg/transaction/manager_test.go`

Fake driver helper (top of file):

```go
type noopConnector struct{}
func (noopConnector) Connect(context.Context) (driver.Conn, error) { return noopConn{}, nil }
func (noopConnector) Driver() driver.Driver                        { return noopDriver{} }

type noopDriver struct{}
func (noopDriver) Open(string) (driver.Conn, error) { return noopConn{}, nil }

type noopConn struct{}
func (noopConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("not implemented") }
func (noopConn) Close() error                        { return nil }
func (noopConn) Begin() (driver.Tx, error)           { return noopTx{}, nil }

type noopTx struct{}
func (noopTx) Commit() error   { return nil }
func (noopTx) Rollback() error { return nil }

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sql.OpenDB(noopConnector{}),
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test gorm: %v", err)
	}
	return gdb
}
```

Table-driven cases:

| Case | Input | Expect |
|---|---|---|
| `ExecuteWithOptions_fn_panics_returns_error` | fn = `func(){ panic("boom") }` | err != nil, contains `"panic in transaction"` และ contains `"boom"`, **ไม่ crash** (ไม่ต้อง assert.Panics) |
| `ExecuteWithOptions_success_commits` | fn = nil-return | err == nil |
| `ExecuteWithOptions_fn_error_rolls_back` | fn returns `errors.New("biz")` | err == biz error (errors.Is) |
| `Execute_delegates_to_options` | fn panics | err contains `"panic in transaction"` |
| `RunInTransaction_fn_panics_returns_error` | fn panics | err contains `"panic in transaction"` |
| `RunInTransaction_fn_error_rolls_back` | fn returns sentinel | errors.Is sentinel |
| `WithTransaction_panic_returns_error` | fn panics | err contains `"panic in transaction"` |
| `WithTransaction_complete_commit` | fn nil-return | err == nil |

Assert ด้วย `strings.Contains` หรือ `testify/require` (repo มี testify v1.11.1 แล้ว — ใช้ testify ให้ตรง convention ไฟล์ test อื่น)

### 9.2 `internal/modules/iot/delivery/http/routes_test.go`

| Case | Input | Expect |
|---|---|---|
| `nil_router` | `(nil, validHandler?, nil)` — router nil ตรวจก่อน h/mw | err message == `"MapMQTT3Routes: router is nil"` |
| `nil_handler` | `(chi.NewRouter(), nil, nil)` | err message == `"MapMQTT3Routes: handler is nil"` |
| `nil_mw` | `(chi.NewRouter(), &MQTT3Handler{}, nil)` | err message == `"MapMQTT3Routes: middleware manager is nil"` |

หมายเหตุ: `&MQTT3Handler{}` zero value ใช้ได้เพราะ guard ตรวจ pointer nil ก่อน touch method

### 9.3 `internal/modules/mqtt/delivery/http/routes_test.go`

Test `mapMQTTRoutesOnce` ตรง (เลี่ยงปัญหา `sync.Once` cache ระหว่าง test):

| Case | Input | Expect |
|---|---|---|
| `nil_router` | `(nil, nil, nil)` | err == `"mapMQTTRoutesOnce: router is nil"` |
| `nil_handler` | `(chi.NewRouter(), nil, nil)` | err == `"mapMQTTRoutesOnce: handler is nil"` |
| `nil_mw` | `(chi.NewRouter(), &MQTTHandler{}, nil)` | err == `"mapMQTTRoutesOnce: middleware manager is nil"` |

Branch "middleware returned nil" (mw จริงแต่ func คืน nil) ไม่ cover เพราะ `CreateMiddlewareManager` ต้องการ cfg+logger+usersUC เต็ม — note limitation ไว้ใน test comment

---

## Verification

```
go build ./...
go vet ./...
go test ./pkg/transaction/... ./internal/modules/mqtt/delivery/http/... ./internal/modules/iot/delivery/http/...
go test ./...
rg -n "panic\(" --glob "!vendor/**" --glob "!*_test.go"  # expect: only matches in PANIC_AUDIT.md examples / none in .go
```
