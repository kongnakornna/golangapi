package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// errUnauthorized ใช้ภายใน handler helper
var errUnauthorized = errors.New("unauthorized")

// Handlers รวม handler ทั้งหมด
type Handlers struct {
	Consent ConsentHandler
	DSAR    DSARHandler
	Account AccountHandler
	Policy  PolicyHandler
	Admin   AdminHandler
}

// DSARHandler describes the DSAR endpoints used by the router.
type DSARHandler interface {
	Submit(http.ResponseWriter, *http.Request)
	List(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	Process(http.ResponseWriter, *http.Request)
}

// AccountHandler describes the account-related endpoints used by the router.
type AccountHandler interface {
	Suspend(http.ResponseWriter, *http.Request)
	Terminate(http.ResponseWriter, *http.Request)
	ConfirmDeletion(http.ResponseWriter, *http.Request)
}

// AdminHandler describes the admin endpoints used by the router.
type AdminHandler interface {
	GetReport(http.ResponseWriter, *http.Request)
}

// PolicyHandler describes the policy endpoints used by the router.
type PolicyHandler interface {
	GetActive(http.ResponseWriter, *http.Request)
	ListVersions(http.ResponseWriter, *http.Request)
	Publish(http.ResponseWriter, *http.Request)
}

// AuthMiddleware interface เพื่อหลีกเลี่ยง cyclic import กับ auth module
type AuthMiddleware func(http.Handler) http.Handler

// RegisterRoutes ลงทะเบียน routes ทั้งหมด
// RegisterRoutes registers all PDPA routes
func RegisterRoutes(r chi.Router, h *Handlers, auth AuthMiddleware, adminOnly AuthMiddleware) {
	r.Route("/api/v1/pdpa", func(r chi.Router) {
		r.Use(auth)

		// Consent
		r.Post("/consent", h.Consent.RecordConsent)
		r.Delete("/consent", h.Consent.RevokeConsent)
		r.Get("/consent/history", h.Consent.GetHistory)

		// DSAR
		r.Post("/dsar", h.DSAR.Submit)
		r.Get("/dsar", h.DSAR.List)
		r.Get("/dsar/{id}", h.DSAR.Get)

		// Policy (public read)
		r.Get("/policy", h.Policy.GetActive)
		r.Get("/policy/versions", h.Policy.ListVersions)

		// Admin-only
		r.Group(func(r chi.Router) {
			r.Use(adminOnly)

			r.Post("/dsar/{id}/process", h.DSAR.Process)

			r.Post("/account/suspend", h.Account.Suspend)
			r.Post("/account/terminate", h.Account.Terminate)
			r.Post("/account/confirm-deletion", h.Account.ConfirmDeletion)

			r.Post("/policy", h.Policy.Publish)

			r.Get("/admin/reports", h.Admin.GetReport)
		})
	})
}

// ConsentHandler describes the consent endpoints used by the router.
type ConsentHandler interface {
	RecordConsent(http.ResponseWriter, *http.Request)
	RevokeConsent(http.ResponseWriter, *http.Request)
	GetHistory(http.ResponseWriter, *http.Request)
}