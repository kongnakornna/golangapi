# Module 9: pkg/logger (Structured Logging)

## สำหรับโฟลเดอร์ `internal/pkg/logger/`

ไฟล์ที่เกี่ยวข้อง:
- `internal/pkg/logger/zap_logger.go`

---

## หลักการ (Concept)

### คืออะไร?
Logger คือระบบบันทึกเหตุการณ์ (logging) ที่มีโครงสร้าง (structured) ช่วยในการ debug, monitoring, และวิเคราะห์ปัญหา โดยบันทึกเป็น key-value pairs (JSON) ทำให้เครื่องมือเช่น Loki, ELK, หรือ Datadog สามารถค้นหาและวิเคราะห์ได้ง่าย

### มีกี่แบบ?
1. **Standard log (`log`)** – ง่าย แต่ไม่มี structure
2. **Structured logging (`slog`, `zap`, `zerolog`)** – JSON output, มีระดับ severity
3. **Levels**: Debug, Info, Warn, Error, Fatal, Panic
4. **Contextual logging** – เพิ่ม field เช่น trace_id, user_id ลงใน log ทุกบรรทัด

### ใช้อย่างไร / นำไปใช้กรณีไหน
- ใช้บันทึกการทำงานปกติ (Info), คำเตือน (Warn), ข้อผิดพลาด (Error)
- ใช้ Debug สำหรับ development
- ใช้ Fatal เมื่อแอปไม่สามารถทำงานต่อได้

### ทำไมต้องใช้
- แทนที่ `fmt.Println` และ `log.Println` ที่ไม่มี structure
- ช่วยในการ追踪 request ผ่าน trace_id
- วิเคราะห์ performance ด้วย duration fields

### ประโยชน์ที่ได้รับ
- ค้นหา logs ตาม trace_id, user_id, path ได้ง่าย
- รวม logs จากหลาย instance เข้าด้วยกัน
- ลดเวลา debugging

### ข้อควรระวัง
- อย่า log sensitive data (password, token, credit card)
- ระดับ log ใน production ควรเป็น Info หรือ Warn (ไม่ใช่ Debug)
- ระวัง performance ของการ log จำนวนมาก (ใช้ sampling)

### ข้อดี
- structured, ค้นหาได้, รองรับการวิเคราะห์

### ข้อเสีย
- JSON logs อ่านยากกว่า plain text (แต่เครื่องมือจัดการได้)
- เพิ่ม overhead เล็กน้อย

### ข้อห้าม
- ห้าม log password หรือ secret
- ห้ามใช้ `fmt.Println` ใน production
- ห้าม ignore error จาก logger (ไม่ค่อยเกิด)

---

## โค้ดที่รันได้จริง

### ไฟล์ `internal/pkg/logger/zap_logger.go`

```go
// Package logger provides a structured logging wrapper around uber-go/zap.
// ----------------------------------------------------------------
// แพ็คเกจ logger ให้บริการ logging แบบมีโครงสร้าง (wrapper ของ uber-go/zap)
package logger

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is the global logger instance.
// ----------------------------------------------------------------
// Log คือ instance logger ระดับโลก
var Log *zap.Logger

// Config holds logger configuration.
// ----------------------------------------------------------------
// Config เก็บค่ากำหนด logger
type Config struct {
	Level      string // debug, info, warn, error
	Env        string // development, production
	OutputPath string // stdout, stderr, or file path
}

// Init initializes the global logger based on config.
// ----------------------------------------------------------------
// Init เริ่มต้น logger ระดับโลกตาม config
func Init(cfg Config) error {
	var zapCfg zap.Config
	if cfg.Env == "production" {
		zapCfg = zap.NewProductionConfig()
		zapCfg.EncoderConfig.TimeKey = "timestamp"
		zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapCfg.EncoderConfig.CallerKey = "caller"
	} else {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorEncoder
	}
	// Set log level
	// ตั้งค่าระดับ log
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}
	zapCfg.Level = zap.NewAtomicLevelAt(level)
	
	// Set output path
	// ตั้งค่า output path
	if cfg.OutputPath != "" {
		zapCfg.OutputPaths = []string{cfg.OutputPath}
	} else {
		zapCfg.OutputPaths = []string{"stdout"}
	}
	zapCfg.ErrorOutputPaths = []string{"stderr"}
	
	var err error
	Log, err = zapCfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}
	return nil
}

// Sync flushes any buffered log entries.
// Should be called before application exit.
// ----------------------------------------------------------------
// Sync ล้าง buffer log ควรเรียกก่อนปิดแอปพลิเคชัน
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

// WithTraceID returns a logger with trace_id field.
// ----------------------------------------------------------------
// WithTraceID คืน logger ที่มีฟิลด์ trace_id
func WithTraceID(traceID string) *zap.Logger {
	return Log.With(zap.String("trace_id", traceID))
}

// WithUserID returns a logger with user_id field.
// ----------------------------------------------------------------
// WithUserID คืน logger ที่มีฟิลด์ user_id
func WithUserID(userID uint) *zap.Logger {
	return Log.With(zap.Uint("user_id", userID))
}

// WithRequest returns a logger with HTTP request fields.
// ----------------------------------------------------------------
// WithRequest คืน logger ที่มีฟิลด์ของ HTTP request
func WithRequest(method, path, remoteAddr string) *zap.Logger {
	return Log.With(
		zap.String("method", method),
		zap.String("path", path),
		zap.String("remote_addr", remoteAddr),
	)
}

// Debug logs a debug message with optional fields.
// ----------------------------------------------------------------
// Debug บันทึกข้อความระดับ debug พร้อม fields เพิ่มเติม
func Debug(msg string, fields ...zap.Field) {
	if Log == nil {
		return
	}
	Log.Debug(msg, fields...)
}

// Info logs an info message.
// ----------------------------------------------------------------
// Info บันทึกข้อความระดับ info
func Info(msg string, fields ...zap.Field) {
	if Log == nil {
		return
	}
	Log.Info(msg, fields...)
}

// Warn logs a warning message.
// ----------------------------------------------------------------
// Warn บันทึกข้อความระดับ warning
func Warn(msg string, fields ...zap.Field) {
	if Log == nil {
		return
	}
	Log.Warn(msg, fields...)
}

// Error logs an error message.
// ----------------------------------------------------------------
// Error บันทึกข้อความระดับ error
func Error(msg string, fields ...zap.Field) {
	if Log == nil {
		return
	}
	Log.Error(msg, fields...)
}

// Fatal logs a fatal message and exits.
// ----------------------------------------------------------------
// Fatal บันทึกข้อความระดับ fatal และจบการทำงาน
func Fatal(msg string, fields ...zap.Field) {
	if Log == nil {
		os.Exit(1)
	}
	Log.Fatal(msg, fields...)
}

// Panic logs a panic message and panics.
// ----------------------------------------------------------------
// Panic บันทึกข้อความระดับ panic และ panic
func Panic(msg string, fields ...zap.Field) {
	if Log == nil {
		panic(msg)
	}
	Log.Panic(msg, fields...)
}

// Debugf formats and logs a debug message.
// ----------------------------------------------------------------
// Debugf จัดรูปแบบและบันทึกข้อความระดับ debug
func Debugf(template string, args ...interface{}) {
	if Log == nil {
		return
	}
	Log.Sugar().Debugf(template, args...)
}

// Infof formats and logs an info message.
// ----------------------------------------------------------------
// Infof จัดรูปแบบและบันทึกข้อความระดับ info
func Infof(template string, args ...interface{}) {
	if Log == nil {
		return
	}
	Log.Sugar().Infof(template, args...)
}

// Warnf formats and logs a warning message.
// ----------------------------------------------------------------
// Warnf จัดรูปแบบและบันทึกข้อความระดับ warning
func Warnf(template string, args ...interface{}) {
	if Log == nil {
		return
	}
	Log.Sugar().Warnf(template, args...)
}

// Errorf formats and logs an error message.
// ----------------------------------------------------------------
// Errorf จัดรูปแบบและบันทึกข้อความระดับ error
func Errorf(template string, args ...interface{}) {
	if Log == nil {
		return
	}
	Log.Sugar().Errorf(template, args...)
}

// Fatalf formats and logs a fatal message and exits.
// ----------------------------------------------------------------
// Fatalf จัดรูปแบบและบันทึกข้อความระดับ fatal และจบการทำงาน
func Fatalf(template string, args ...interface{}) {
	if Log == nil {
		os.Exit(1)
	}
	Log.Sugar().Fatalf(template, args...)
}
```

### ตัวอย่างการใช้งานใน middleware (request logger)

```go
// ใน middleware/logger.go (ปรับปรุงให้ใช้ structured logger)
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		traceID := uuid.New().String()
		
		// Add trace ID to context
		ctx := context.WithValue(r.Context(), "trace_id", traceID)
		r = r.WithContext(ctx)
		
		// Create logger with request fields
		reqLogger := logger.WithRequest(r.Method, r.URL.Path, r.RemoteAddr)
		reqLogger = reqLogger.With(zap.String("trace_id", traceID))
		
		reqLogger.Info("request started")
		
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		
		duration := time.Since(start)
		reqLogger.Info("request completed",
			zap.Int("status", wrapped.statusCode),
			zap.Duration("duration", duration),
		)
	})
}
```

### ตัวอย่างการเริ่มต้นใน `main.go`

```go
func main() {
	// Initialize logger
	err := logger.Init(logger.Config{
		Level: "debug",
		Env:   "development",
	})
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	
	logger.Info("application starting", zap.String("version", "1.0.0"))
	
	// ... rest of application
}
```

### ตัวอย่าง output (JSON)

```json
{
  "level": "info",
  "timestamp": "2025-04-04T10:30:00Z",
  "caller": "middleware/logger.go:25",
  "msg": "request started",
  "method": "POST",
  "path": "/login",
  "remote_addr": "192.168.1.1",
  "trace_id": "abc123"
}
```

---

## วิธีใช้งาน module นี้

1. ติดตั้ง dependency:
   ```bash
   go get go.uber.org/zap
   ```
2. วาง `zap_logger.go` ใน `internal/pkg/logger/`
3. เรียก `logger.Init()` ตอนเริ่มโปรแกรม
4. ใช้ `logger.Info()`, `logger.Error()` ทั่วทั้งโปรเจกต์
5. ใช้ `logger.WithTraceID()` เพื่อเพิ่ม trace context

---

## ตารางสรุประดับ Log

| ระดับ | ใช้เมื่อ | ตัวอย่าง |
|-------|---------|----------|
| Debug | ข้อมูลละเอียดสำหรับ developer | ค่า variable, ข้อมูลการดีบัก |
| Info | เหตุการณ์ปกติที่ควรบันทึก | server start, request received, user login |
| Warn | เหตุการณ์ผิดปกติแต่ยังทำงานต่อ | rate limit ใกล้เต็ม, retry |
| Error | เกิดข้อผิดพลาดที่ต้องแจ้งเตือน | DB connection lost, ส่งอีเมล failed |
| Fatal | ข้อผิดพลาดร้ายแรง แอปไม่สามารถทำงานต่อ | migrate failed, config missing |
| Panic | ข้อผิดพลาดที่ไม่คาดคิด | unexpected nil pointer |

---

## แบบฝึกหัดท้าย module (3 ข้อ)

1. เพิ่ม `WithDuration(startTime time.Time)` helper ที่คำนวณ duration และเพิ่มฟิลด์ `duration_ms`
2. สร้าง middleware `RequestLogger` ที่ใช้ logger ตัวนี้แทนที่การ log แบบเดิม (รวม trace_id, status, duration)
3. ปรับปรุง `Init` ให้รองรับการเขียน log ไปยังไฟล์ด้วยการหมุนเวียน (log rotation) โดยใช้ `lumberjack`

---

## แหล่งอ้างอิง

- [Uber Zap documentation](https://github.com/uber-go/zap)
- [Structured logging best practices](https://betterstack.com/community/guides/logging/logging-best-practices/)
- [Zap performance benchmarks](https://github.com/uber-go/zap#performance)

---

**หมายเหตุ:** module นี้ครบถ้วนสำหรับ `pkg/logger` หากต้องการ module ถัดไป (เช่น `pkg/redis`, `pkg/email`, `pkg/jwt`) โปรดแจ้ง