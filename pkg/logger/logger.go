package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// Logger อินเทอร์เฟซสำหรับตัวบันทึก
type Logger interface {
	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
	With(keysAndValues ...any) Logger
	WithContext(ctx context.Context) Logger
}

// StructuredLogger ตัวบันทึกแบบมีโครงสร้าง
type StructuredLogger struct {
	logger *slog.Logger
	ctx    context.Context
}

// LogConfig การตั้งค่าสำหรับบันทึก
type LogConfig struct {
	Level      string `yaml:"level" json:"level"`             // ระดับของ: debug, info, warn, error
	File       string `yaml:"file" json:"file"`               // ที่อยู่ไฟล์บันทึก ถ้าว่างจะส่งออกไปที่คอนโซล
	Console    bool   `yaml:"console" json:"console"`         // ส่งออกไปที่คอนโซลหรือไม่
	MaxSize    int    `yaml:"max_size" json:"max_size"`       // ขนาดสูงสุดของไฟล์ (MB)
	MaxBackups int    `yaml:"max_backups" json:"max_backups"` // จำนวนไฟล์สำรองที่เก็บไว้
	MaxAge     int    `yaml:"max_age" json:"max_age"`         // จำนวนวันที่เก็บไว้
	Compress   bool   `yaml:"compress" json:"compress"`       // บีบอัดหรือไม่
}

// ContextKey ชนิดของคีย์ใน context
type ContextKey string

const (
	// TraceIDKey คีย์สำหรับ ID การติดตามเส้นทาง (trace)
	TraceIDKey ContextKey = "trace_id"
	// RequestIDKey คีย์สำหรับ ID คำขอ
	RequestIDKey ContextKey = "request_id"
	// UserIDKey คีย์สำหรับ ID ผู้ใช้
	UserIDKey ContextKey = "user_id"
)

// NewLogger สร้างตัวบันทึกใหม่
func NewLogger(config *LogConfig) (*StructuredLogger, error) {
	level := parseLevel(config.Level)

	var writers []io.Writer

	// ส่งออกไปที่คอนโซล
	if config.Console || config.File == "" {
		writers = append(writers, os.Stdout)
	}

	// ส่งออกไปที่ไฟล์
	if config.File != "" {
		file, err := createLogFile(config.File)
		if err != nil {
			return nil, err
		}
		writers = append(writers, file)
	}

	// ถ้าไม่มีการตั้งค่าการส่งออกใด ๆ ให้ส่งออกไปที่คอนโซลโดยปริยาย
	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}

	var writer io.Writer
	if len(writers) == 1 {
		writer = writers[0]
	} else {
		writer = io.MultiWriter(writers...)
	}

	// สร้าง handler
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// จัดรูปแบบเวลา
			if a.Key == slog.TimeKey {
				return slog.String("timestamp", a.Value.Time().Format(time.RFC3339))
			}

			// จัดรูปแบบตำแหน่ง source code
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				// แสดงเฉพาะชื่อไฟล์และหมายเลขบรรทัด
				return slog.String("source", filepath.Base(source.File)+":"+
					func() string {
						return slog.IntValue(source.Line).String()
					}())
			}

			return a
		},
	})

	logger := slog.New(handler)

	return &StructuredLogger{
		logger: logger,
		ctx:    context.Background(),
	}, nil
}

// Default สร้างตัวบันทึกโดยปริยาย
func Default() *StructuredLogger {
	logger, _ := NewLogger(&LogConfig{
		Level:   "info",
		Console: true,
	})
	return logger
}

// parseLevel แปลงระดับ
func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// createLogFile สร้างไฟล์บันทึก
func createLogFile(filename string) (*os.File, error) {
	// เพิ่มวันที่ในชื่อไฟล์
	dir := filepath.Dir(filename)
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	nameWithoutExt := base[:len(base)-len(ext)]

	// สร้างชื่อไฟล์ที่มีวันที่
	dateStr := time.Now().Format("2006-01-02")
	newFilename := filepath.Join(dir, nameWithoutExt+"-"+dateStr+ext)

	// ตรวจสอบให้แน่ใจว่าไดเรกทอรีมีอยู่
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	return os.OpenFile(newFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
}

// Debug บันทึกระดับ debug
func (l *StructuredLogger) Debug(msg string, keysAndValues ...any) {
	l.log(slog.LevelDebug, msg, keysAndValues...)
}

// Info บันทึกระดับ info
func (l *StructuredLogger) Info(msg string, keysAndValues ...any) {
	l.log(slog.LevelInfo, msg, keysAndValues...)
}

// Warn บันทึกระดับ warn
func (l *StructuredLogger) Warn(msg string, keysAndValues ...any) {
	l.log(slog.LevelWarn, msg, keysAndValues...)
}

// Error บันทึกระดับ error
func (l *StructuredLogger) Error(msg string, keysAndValues ...any) {
	l.log(slog.LevelError, msg, keysAndValues...)
}

// log วิธีการภายในสำหรับบันทึก
func (l *StructuredLogger) log(level slog.Level, msg string, keysAndValues ...any) {
	// เพิ่มข้อมูลผู้เรียก
	pc, file, line, ok := runtime.Caller(2)
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		keysAndValues = append(keysAndValues,
			"caller", filepath.Base(file)+":"+slog.IntValue(line).String(),
			"function", filepath.Base(funcName),
		)
	}

	// ดึงข้อมูลจาก context
	if l.ctx != nil {
		if traceID := l.ctx.Value(TraceIDKey); traceID != nil {
			keysAndValues = append(keysAndValues, "trace_id", traceID)
		}
		if requestID := l.ctx.Value(RequestIDKey); requestID != nil {
			keysAndValues = append(keysAndValues, "request_id", requestID)
		}
		if userID := l.ctx.Value(UserIDKey); userID != nil {
			keysAndValues = append(keysAndValues, "user_id", userID)
		}
	}

	l.logger.Log(l.ctx, level, msg, keysAndValues...)
}

// With เพิ่มฟิลด์ลงในตัวบันทึก
func (l *StructuredLogger) With(keysAndValues ...any) Logger {
	return &StructuredLogger{
		logger: l.logger.With(keysAndValues...),
		ctx:    l.ctx,
	}
}

// WithContext กำหนด context ให้กับตัวบันทึก
func (l *StructuredLogger) WithContext(ctx context.Context) Logger {
	return &StructuredLogger{
		logger: l.logger,
		ctx:    ctx,
	}
}

// GetTraceID ดึง trace ID จาก context
func GetTraceID(ctx context.Context) string {
	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		if str, ok := traceID.(string); ok {
			return str
		}
	}
	// ลองดึงจาก middleware ของ Chi
	if reqID := middleware.GetReqID(ctx); reqID != "" {
		return reqID
	}
	return ""
}

// GetRequestID ดึง request ID จาก context
func GetRequestID(ctx context.Context) string {
	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		if str, ok := requestID.(string); ok {
			return str
		}
	}
	return ""
}

// GetUserID ดึง user ID จาก context
func GetUserID(ctx context.Context) string {
	if userID := ctx.Value(UserIDKey); userID != nil {
		if str, ok := userID.(string); ok {
			return str
		}
	}
	return ""
}

// WithTraceID กำหนด trace ID ใน context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// WithRequestID กำหนด request ID ใน context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// WithUserID กำหนด user ID ใน context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// LoggerMiddleware มิดเดิลแวร์สำหรับบันทึก
func LoggerMiddleware(logger Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// ดึงหรือสร้าง request ID
			requestID := middleware.GetReqID(r.Context())
			if requestID == "" {
				requestID = fmt.Sprintf("%d", middleware.NextRequestID())
			}

			// ดึงหรือสร้าง trace ID
			traceID := r.Header.Get("X-Trace-ID")
			if traceID == "" {
				traceID = requestID
			}

			// ตั้งค่า response header
			w.Header().Set("X-Request-ID", requestID)
			w.Header().Set("X-Trace-ID", traceID)

			// สร้าง context ที่มีข้อมูลติดตาม
			ctx := WithRequestID(r.Context(), requestID)
			ctx = WithTraceID(ctx, traceID)

			// สร้างตัวบันทึกที่มี context
			ctxLogger := logger.WithContext(ctx)

			// บันทึกการเริ่มต้นคำขอ
			ctxLogger.Info("request started",
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
				"user_agent", r.UserAgent(),
				"remote_addr", r.RemoteAddr,
			)

			// สร้างตัวบันทึกการตอบสนอง
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// ดำเนินการคำขอ
			next.ServeHTTP(ww, r.WithContext(ctx))

			// บันทึกการเสร็จสิ้นคำขอ
			duration := time.Since(start)
			ctxLogger.Info("request completed",
				"status", ww.Status(),
				"duration_ms", duration.Milliseconds(),
				"bytes", ww.BytesWritten(),
			)

			// บันทึกคำขอที่ช้า
			if duration > 5*time.Second {
				ctxLogger.Warn("slow request detected",
					"duration_ms", duration.Milliseconds(),
					"threshold_ms", 5000,
				)
			}
		})
	}
}

// RecoveryMiddleware มิดเดิลแวร์สำหรับกู้คืนจาก panic
func RecoveryMiddleware(logger Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// รับข้อมูล stack trace
					buf := make([]byte, 4096)
					n := runtime.Stack(buf, false)
					stackTrace := string(buf[:n])

					// บันทึก panic
					logger.WithContext(r.Context()).Error("panic recovered",
						"error", err,
						"stack_trace", stackTrace,
						"method", r.Method,
						"path", r.URL.Path,
					)

					// ส่งกลับข้อผิดพลาด 500
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}