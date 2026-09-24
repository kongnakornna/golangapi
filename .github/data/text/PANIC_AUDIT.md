# Panic Audit & Fix — icmongolang

วันที่: 2026-08-23 | Policy: strict no-panic (non-test code) | Plan: `docs/plan/2026-08-23-panic-audit-and-fix.md`

## สรุปผู้บริหาร

- พบ panic **16 จุด / 6 ไฟล์** (สแกน `*.go` ทั้ง repo ไม่รวม `vendor/`) — แก้ครบทั้งหมด
- HTTP layer มี chi `middleware.Recoverer` ครอบอยู่แล้ว (`internal/delivery/rest/router.go:69`, `internal/server/handlers.go:169`) — panic ใน handler ไม่ kill process แต่เป็น 500 generic; policy ใหม่ตัด panic ตั้งแต่ต้นทาง
- Policy ใหม่: runtime error → return error / startup fatal → log + exit / observability fail → log + continue

## Inventory

| ไฟล์ | Line เดิม | ประเภท | Risk | วิธีแก้ | สถานะ |
|---|---|---|---|---|---|
| `cmd/root.go` | 61, 67 | dead code (`log.Fatalf` exit ก่อน panic ทำงาน) | low | ลบ panic | ✅ |
| `internal/server/server.go` | 160 | crash-risk (`ListenAndServe` fail ใน goroutine) | high | `logger.Errorf` + `os.Exit(1)` | ✅ |
| `internal/server/server.go` | 173 | crash-risk (`Shutdown` fail → ข้าม cleanup MQTT/Influx) | medium | `logger.Errorf` เท่านั้น | ✅ |
| `internal/middleware/prometheus.go` | 45 | boot-crash (`init()` register metric fail) | medium | log warning + skip metric | ✅ |
| `internal/modules/iot/delivery/http/routes.go` | 11, 14, 17 | nil-guard (programmer error ตอน wiring) | low | signature → `return error` | ✅ |
| `internal/modules/mqtt/delivery/http/routes.go` | 16, 19, 22, 38 | nil-guard (รวม middleware คืน nil หลัง map route ไปบางส่วน) | low | validate ก่อน map + `return error` | ✅ |
| `pkg/transaction/manager.go` | 56, 87, 257 | intentional re-panic หลัง rollback | high | rollback + return error (named return) | ✅ |

## Policy

เกณฑ์เขียนโค้ดใหม่:

1. **ห้าม `panic()` ใน non-test code** ทุกกรณี
2. Runtime error → return error wrap ด้วย `%w`
3. Startup fatal (config / port / dependency จำเป็นล่ม) → `log.Fatalf` หรือ `logger.Errorf` + `os.Exit(1)`
4. Observability (metrics / log sink) ล่ม → log warning + continue — **ห้าม** kill service
5. Panic ใน transaction → rollback + return error (ไม่ re-panic) — caller handle เอง

## รายละเอียดรายจุด (Before → After)

### cmd/root.go

```go
// Before                          // After
log.Fatalf("LoadConfig: %v", err)  log.Fatalf("LoadConfig: %v", err)
panic(err) // unreachable          // (ลบ)
```

เหตุผล: `log.Fatalf` เรียก `os.Exit(1)` แล้ว บรรทัด panic เป็น dead code

### internal/server/server.go — Start()

```go
// Before                              // After
if err := srv.server.ListenAndServe(); if err := srv.server.ListenAndServe();
    err != http.ErrServerClosed {          err != http.ErrServerClosed {
    panic(err)                             srv.logger.Errorf("❌ server error: %v", err)
}                                          os.Exit(1)
                                       }

if err := srv.server.Shutdown(ctx);    if err := srv.server.Shutdown(ctx);
    err != nil {                           err != nil {
    panic(err) // ข้าม cleanup             srv.logger.Errorf("❌ server shutdown error: %v", err)
}                                      }
```

เหตุผล: port conflict = startup fatal ควร exit ชัดเจน; shutdown fail ห้ามข้ามการ disconnect MQTT / close InfluxDB

### internal/middleware/prometheus.go — init()

```go
// Before                                   // After
if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
    panic(err)                                  log.Printf("⚠️ prometheus register failed: %v", err)
}
```

เหตุผล: metric เสียจุดเดียวไม่ควร crash service ทั้งตัว; `init()` return error ไม่ได้

### internal/modules/iot/delivery/http/routes.go — MapMQTT3Routes

```go
// Before                                    // After
func MapMQTT3Routes(...) {                   func MapMQTT3Routes(...) error {
    if router == nil {                           if router == nil {
        panic("MapMQTTRoutes: router is nil")        return errors.New("MapMQTT3Routes: router is nil")
    }                                            }
    ...                                          ...
}                                                return nil
                                             }
```

เหตุผล: wiring error ตรวจได้ตอน boot — return error ให้ caller (`handlers.go New()`) fail-fast พร้อม context; แก้ message สะกดผิด `"MapMQTTRoutes"` → `"MapMQTT3Routes"`

### internal/modules/mqtt/delivery/http/routes.go — MapMQTTRoutes

```go
// Before                                     // After
func MapMQTTRoutes(...) {                     func MapMQTTRoutes(...) error {
    once.Do(func() {                              var err error
        if router == nil {                        once.Do(func() { err = mapMQTTRoutesOnce(router, h, mw) })
            panic("...router is nil")             return err
        }                                     }
        ...
        r.Route("/mqtt", ...)                 // helper: validate ครบก่อน แล้วค่อย map route
        // middleware check ทีหลัง            func mapMQTTRoutesOnce(...) error {
        if verifier == nil || ... {               if router == nil { return errors.New("mapMQTTRoutesOnce: router is nil") }
            panic("...middleware...")             ...
        }                                         rateLimit := mw.RateLimit(); verifier := mw.Verifier(true); ...
    })                                            if rateLimit == nil || verifier == nil || ... {
}                                                     return errors.New("mapMQTTRoutesOnce: one or more middleware functions returned nil")
                                                  }
                                                  router.Route("/mqtt", ...)
                                                  return nil
                                              }
```

เหตุผล: เดิม middleware-nil panic โดนหลัง `r.Use(mw.RateLimit())` + register route public ไปแล้ว → state ค้าง; ตอนนี้ validate ครบก่อนแตะ router และ `sync.Once` capture error ได้

### pkg/transaction/manager.go — 3 recover blocks

```go
// Before                                // After
defer func() {                           defer func() {
    if r := recover(); r != nil {            if r := recover(); r != nil {
        tx.Rollback()                            tx.Rollback()
        panic(r) // re-panic                     err = fmt.Errorf("panic in transaction: %v", r)
    }                                        }
}();                                     }()
// (signature เปลี่ยนเป็น named return `(err error)`)
```

ใช้กับ: `ExecuteWithOptions`, `RunInTransaction`, `WithTransaction`
เหตุผล: re-panic ทำให้ business panic กระจายถึง chi Recoverer → 500 generic + stack trace; rollback เสร็จแล้วควรคืน error ให้ caller handle ตาม flow ปกติ

## Behavior Change Notes

- **Panic ใน transaction** เดิม → chi Recoverer จับ → 500 generic; ตอนนี้ → error `panic in transaction: ...` คืน caller — code path ที่เคยพึ่ง recover ต้อง handle error เอง
- **MapMQTT3Routes / MapMQTTRoutes** เปลี่ยน signature → compile error จับ call site ที่พลาดได้หมด (caller ปัจจุบันมีจุดเดียว: `internal/server/handlers.go`)
- Error message ของ mqtt routes เปลี่ยน prefix เป็น `mapMQTTRoutesOnce:` — ถ้ามี log parsing ผูก string เดิม (`MapMQTTRoutes: ...`) ต้องปรับ

## การป้องกันซ้ำ

1. golangci-lint `forbidigo` rule (แนะนำ):

```yaml
linters-settings:
  forbidigo:
    forbid:
      - pattern: 'panic\('
        msg: "use error return / log.Fatal instead (see docs/PANIC_AUDIT.md)"
```

(exclude `*_test.go`)

2. Review checklist: PR ที่เพิ่ม `panic(` ต้อง justify เหตุผลใน description
