// Package http จัดเตรียม HTTP layer ของโมดูล PDPA
//
// Layout ของ route ทั้งหมด:
//
//	/api/v1/pdpa
//	├── /policy              [GET]    — public read
//	├── /policy/versions     [GET]    — public read
//	├── /consent             [POST]   — auth
//	├── /consent             [DELETE] — auth
//	├── /consent/history     [GET]    — auth
//	├── /dsar                [POST]   — auth
//	├── /dsar                [GET]    — auth
//	├── /dsar/{id}           [GET]    — auth
//	├── /dsar/{id}/process   [POST]   — auth + admin
//	├── /account/suspend     [POST]   — auth + admin
//	├── /account/terminate   [POST]   — auth + admin
//	├── /account/confirm-deletion [POST] — auth + admin
//	├── /policy              [POST]   — auth + admin
//	└── /admin/reports       [GET]    — auth + admin
package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// PathPrefix คือ mount point ของทุก PDPA route
//
// ⚠️ ต้องสอดคล้องกับ swag annotation:
//   - @BasePath ใน cmd/api/main.go = "/api"
//   - @Router    ใน handler แต่ละตัว = "/v1/pdpa/..."
//
// ตัวอย่าง: PathPrefix = "/api/v1/pdpa" + @Router "/v1/pdpa/consent" + @BasePath "/api"
// → Swagger UI แสดง /api/v1/pdpa/consent (ตรงกับ runtime)
//
// ถ้าเปลี่ยน PathPrefix ต้องอัปเดต @Router ใน handler ทุกตัวให้สอดคล้อง
const PathPrefix = "/v1/pdpa"

// ============================================================
// Handler interfaces
// ============================================================

// Handlers รวม handler ทั้งหมดของ PDPA module
type Handlers struct {
	Consent ConsentHandler
	DSAR    DSARHandler
	Account AccountHandler
	Policy  PolicyHandler
	Admin   AdminHandler
}

// ConsentHandler — endpoints เกี่ยวกับ consent
type ConsentHandler interface {
	RecordConsent(http.ResponseWriter, *http.Request)
	RevokeConsent(http.ResponseWriter, *http.Request)
	GetHistory(http.ResponseWriter, *http.Request)
}

// DSARHandler — endpoints เกี่ยวกับ Data Subject Access Request
type DSARHandler interface {
	Submit(http.ResponseWriter, *http.Request)
	List(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	Process(http.ResponseWriter, *http.Request)
}

// AccountHandler — endpoints เกี่ยวกับ user account lifecycle
type AccountHandler interface {
	Suspend(http.ResponseWriter, *http.Request)
	Terminate(http.ResponseWriter, *http.Request)
	ConfirmDeletion(http.ResponseWriter, *http.Request)
}

// PolicyHandler — endpoints เกี่ยวกับ privacy policy
type PolicyHandler interface {
	GetActive(http.ResponseWriter, *http.Request)
	ListVersions(http.ResponseWriter, *http.Request)
	Publish(http.ResponseWriter, *http.Request)
}

// AdminHandler — endpoints สำหรับ admin เท่านั้น
type AdminHandler interface {
	GetReport(http.ResponseWriter, *http.Request)
}

// AuthMiddleware เป็น adapter type เพื่อหลีกเลี่ยง cyclic import
// กับ auth module — signature เหมือน http middleware ทั่วไป
type AuthMiddleware func(http.Handler) http.Handler

// ============================================================
// Route registration
// ============================================================

// RegisterRoutes ลงทะเบียนทุก PDPA route เข้ากับ chi router
//
// Routes แบ่งเป็น 3 ระดับสิทธิ์:
//
//  1. Public read    — ไม่ต้อง login (policy อ่านได้ตาม PDPA)
//  2. Authenticated  — ต้อง login
//  3. Admin-only     — ต้อง login + มีสิทธิ์ admin
//
// ลำดับ middleware ที่ถูก apply:
//
//	Public:     (ไม่มี)
//	Auth:       auth
//	Admin:      auth → adminOnly
func RegisterRoutes(
	r chi.Router,
	h *Handlers,
	auth AuthMiddleware,
	adminOnly AuthMiddleware,
) {
	r.Route(PathPrefix, func(r chi.Router) {

		// ============================================================
		// 1. Public read — ไม่ต้อง authenticate
		// ============================================================
		// Policy เปิดให้อ่านได้ทุกคนตาม พ.ร.บ.คุ้มครองข้อมูลส่วนบุคคล
		r.Get("/policy", h.Policy.GetActive)
		r.Get("/policy/versions", h.Policy.ListVersions)

		// ============================================================
		// 2. Authenticated user routes
		// ============================================================
		r.Group(func(r chi.Router) {
			r.Use(auth)

			// ---------- Consent ----------
			r.Post("/consent", h.Consent.RecordConsent)
			r.Delete("/consent", h.Consent.RevokeConsent)
			r.Get("/consent/history", h.Consent.GetHistory)

			// ---------- DSAR (user) ----------
			r.Post("/dsar", h.DSAR.Submit)
			r.Get("/dsar", h.DSAR.List)
			r.Get("/dsar/{id}", h.DSAR.Get)

			// ========================================================
			// 3. Admin-only routes (nested in auth group)
			// ========================================================
			r.Group(func(r chi.Router) {
				r.Use(adminOnly)

				// DSAR (admin) — process / approve / reject
				r.Post("/dsar/{id}/process", h.DSAR.Process)

				// Account lifecycle (admin)
				r.Post("/account/suspend", h.Account.Suspend)
				r.Post("/account/terminate", h.Account.Terminate)
				r.Post("/account/confirm-deletion", h.Account.ConfirmDeletion)

				// Policy publish (admin)
				r.Post("/policy", h.Policy.Publish)

				// Compliance reports (admin)
				r.Get("/admin/reports", h.Admin.GetReport)
			})
		})
	})
}
