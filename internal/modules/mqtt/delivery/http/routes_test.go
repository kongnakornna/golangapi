package http

import (
	"testing"

	"github.com/go-chi/chi/v5"
)

// ทดสอบ mapMQTTRoutesOnce โดยตรง เพื่อเลี่ยง sync.Once cache ข้าม test case
//
// note: branch "middleware functions returned nil" (mw จริงแต่ func คืน nil)
// ไม่ cover — CreateMiddlewareManager ต้องการ cfg + logger + usersUC เต็ม

// ============================================================
// Nil guards
// ============================================================

func TestMapMQTTRoutesOnce_NilRouter(t *testing.T) {
	err := mapMQTTRoutesOnce(nil, nil, nil)

	if err == nil || err.Error() != "mapMQTTRoutesOnce: router is nil" {
		t.Fatalf("expected 'router is nil' error, got %v", err)
	}
}

func TestMapMQTTRoutesOnce_NilHandler(t *testing.T) {
	err := mapMQTTRoutesOnce(chi.NewRouter(), nil, nil)

	if err == nil || err.Error() != "mapMQTTRoutesOnce: handler is nil" {
		t.Fatalf("expected 'handler is nil' error, got %v", err)
	}
}

func TestMapMQTTRoutesOnce_NilMiddlewareManager(t *testing.T) {
	err := mapMQTTRoutesOnce(chi.NewRouter(), &MQTTHandler{}, nil)

	if err == nil || err.Error() != "mapMQTTRoutesOnce: middleware manager is nil" {
		t.Fatalf("expected 'middleware manager is nil' error, got %v", err)
	}
}

// ============================================================
// sync.Once idempotency
// ============================================================

// ทดสอบว่า MapMQTTRoutes เรียกได้มากกว่า 1 ครั้งโดยไม่ panic
// (ผลลัพธ์ของครั้งที่ 2 เป็น nil เสมอ เพราะ sync.Once blocking)
func TestMapMQTTRoutes_SecondCallIsNoop(t *testing.T) {
	// ครั้งแรก — จะ error เพราะ argument ไม่ครบ
	err1 := MapMQTTRoutes(nil, nil, nil)
	if err1 == nil {
		t.Fatalf("expected error on first call with nil args")
	}

	// ครั้งที่ 2 — sync.Once ทำงานแล้ว → คืน nil
	err2 := MapMQTTRoutes(nil, nil, nil)
	if err2 != nil {
		t.Fatalf("expected nil on second call (sync.Once), got %v", err2)
	}
}

// ============================================================
// Nil middleware function guards
// ============================================================

// ทดสอบว่า middleware manager ที่คืน nil function → error ที่ถูกต้อง
//
// (ข้ามถ้าไม่สามารถสร้าง mw ได้ — ต้องมี cfg + logger + usersUC จริง)
//
// func TestMapMQTTRoutesOnce_NilMiddlewareFuncs(t *testing.T) {
// 	mw := &middleware.MiddlewareManager{}
// 	err := mapMQTTRoutesOnce(chi.NewRouter(), &MQTTHandler{}, mw)
// 	if err == nil || err.Error() != "mapMQTTRoutesOnce: one or more middleware functions returned nil" {
// 		t.Fatalf("expected 'middleware functions returned nil', got %v", err)
// 	}
// }
