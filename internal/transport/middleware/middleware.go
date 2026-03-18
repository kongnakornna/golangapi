package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	httpx "github.com/kongnakornna/golangapi/internal/transport/httpx"
	apperrors "github.com/kongnakornna/golangapi/pkg/errors"
	"github.com/kongnakornna/golangapi/pkg/logger"
)

// ประเภทคีย์บริบท
type contextKey string

const (
	// คีย์บริบทคำขอ
	reqContextKey contextKey = "request_context" // วัตถุบริบทคำขอ
)

// ReqContext โครงสร้างบริบทคำขอ
type ReqContext struct {
	TraceID    string    // ID ติดตามคำขอ
	RequestID  string    // ID คำขอ
	UserID     uint      // ID ผู้ใช้ (หากรับรองความถูกต้องแล้ว)
	UserRole   string    // บทบาทผู้ใช้ (หากรับรองความถูกต้องแล้ว)
	ClientIP   string    // IP ไคลเอ็นต์
	StartTime  time.Time // เวลาเริ่มต้นคำขอ
	RequestURI string    // URI คำขอ
	Method     string    // วิธีการคำขอ
}

// GetUserIDFromContext รับ ID ผู้ใช้จากบริบทคำขอ
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	reqCtx := GetRequestContext(ctx)
	if reqCtx == nil || reqCtx.UserID == 0 {
		return 0, false
	}
	return reqCtx.UserID, true
}

// GetRequestContext รับบริบทคำขอจาก context.Context
func GetRequestContext(ctx context.Context) *ReqContext {
	if ctx == nil {
		return nil
	}
	if rc, ok := ctx.Value(reqContextKey).(*ReqContext); ok {
		return rc
	}
	return nil
}

// RequestContext มิดเดิลแวร์บริบทคำขอ ตั้งค่าข้อมูลที่เกี่ยวข้องกับคำขอ
func RequestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// สร้างบริบทคำขอ
		reqCtx := &ReqContext{
			RequestID:  r.Header.Get("X-Request-ID"),
			ClientIP:   r.Header.Get("X-Forwarded-For"),
			StartTime:  time.Now(),
			RequestURI: r.RequestURI,
			Method:     r.Method,
		}

		// หากไม่มี ID คำขอ ให้สร้างใหม่
		if reqCtx.RequestID == " {
			reqCtx.RequestID = middleware.GetReqID(r.Context())
		}

		// ตั้งค่า ID ติดตามให้เหมือนกับ ID คำขอ
		reqCtx.TraceID = reqCtx.RequestID

		// หากไม่มี IP ไคลเอ็นต์ ให้ใช้ RemoteAddr
		if reqCtx.ClientIP == " {
			reqCtx.ClientIP = r.RemoteAddr
		}

		// ตั้งค่าส่วนหัวการตอบสนอง
		w.Header().Set("X-Request-ID", reqCtx.RequestID)

		// เพิ่มบริบทคำขอลงใน context
		ctx := context.WithValue(r.Context(), reqContextKey, reqCtx)

		// ดำเนินการประมวลผลคำขอต่อไป
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoggingMiddleware มิดเดิลแวร์บันทึกข้อมูล บันทึกคำขอ
func LoggingMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	if log == nil {
		log = logger.Default()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// รับบริบทคำขอ
			reqCtx := GetRequestContext(r.Context())
			if reqCtx == nil {
				// หากไม่มีบริบทคำขอ ให้สร้างใหม่
				reqCtx = &ReqContext{
					StartTime:  time.Now(),
					RequestURI: r.RequestURI,
					Method:     r.Method,
				}
			}

			// รับขนาดเนื้อหาคำขอ
			var requestSize int64
			if r.ContentLength > 0 {
				requestSize = r.ContentLength
			}

			// ห่อ ResponseWriter เพื่อรับรหัสสถานะ
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// ประมวลผลคำขอ
			next.ServeHTTP(ww, r)

			// คำนวณความหน่วงของการประมวลผลคำขอ
			latency := time.Since(reqCtx.StartTime)

			// สร้างพารามิเตอร์เหตุการณ์บันทึก
			args := []interface{}{
				"method", reqCtx.Method,
				"path", reqCtx.RequestURI,
				"query", r.URL.RawQuery,
				"status", ww.Status(),
				"latency", latency.String(),
				"size", ww.BytesWritten(),
				"req_size", requestSize,
				"ip", reqCtx.ClientIP,
				"user_agent", r.UserAgent(),
				"trace_id", reqCtx.TraceID,
			}

			// เพิ่มข้อมูลผู้ใช้ (ถ้ามี)
			if reqCtx.UserID != 0 {
				args = append(args, "user_id", reqCtx.UserID)
			}

			// บันทึก
			ctxLogger := log.WithContext(r.Context())
			ctxLogger.Info(fmt.Sprintf("%s %s - %d", reqCtx.Method, reqCtx.RequestURI, ww.Status()), args...)
		})
	}
}

// CORSMiddleware จัดการคำขอข้ามโดเมน
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware มิดเดิลแวร์กู้คืน จัดการ panic
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer apperrors.RecoverPanicWithCallback("การจัดการคำขอ HTTP", func(err interface{}) {
			// รับบริบทคำขอ
			reqCtx := GetRequestContext(r.Context())

			// สร้างข้อความแสดงข้อผิดพลาด
			message := "ข้อผิดพลาดภายในเซิร์ฟเวอร์"
			if reqCtx != nil && reqCtx.TraceID != " {
				message = "ข้อผิดพลาดภายในเซิร์ฟเวอร์ โปรดลองอีกครั้งในภายหลัง"
			}

			// ใช้การจัดการตอบสนองข้อผิดพลาดแบบรวมศูนย์
			appErr := apperrors.InternalError(message, fmt.Errorf("%v", err))
			httpx.Error(w, r, appErr)
		})

		next.ServeHTTP(w, r)
	})
}
