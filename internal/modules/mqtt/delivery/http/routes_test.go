package http

import (
	"testing"

	"github.com/go-chi/chi/v5"
)

// ทดสอบ mapMQTTRoutesOnce ตรง เพื่อเลี่ยง sync.Once cache ข้าม test case
// note: branch "middleware functions returned nil" (mw จริงแต่ func คืน nil)
// ไม่ cover — CreateMiddlewareManager ต้องการ cfg + logger + usersUC เต็ม

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
