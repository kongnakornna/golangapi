// Package http จัดเตรียม HTTP layer สำหรับโมดูล keys (API keys)
package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"

	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/auth/domain"
	"icmongolang/internal/modules/keys"
	"icmongolang/internal/modules/keys/presenter"
	"icmongolang/internal/modules/keys/usecase"
	pkgapikey "icmongolang/pkg/apikey"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/responses"
)

type handler struct {
	uc usecase.APIKeyUseCase
}

func NewHandler(uc usecase.APIKeyUseCase) keys.Handlers {
	return &handler{uc: uc}
}

// ─── actorFromRequest ─────────────────────────────────────────
// ⭐ รองรับทั้ง claims.ID และ claims.UserID
func actorFromRequest(r *http.Request) (usecase.Actor, error) {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		return usecase.Actor{}, httpErrors.ErrUnauthorized(nil)
	}

	rawID := claims.ID
	if rawID == "" {
		rawID = claims.UserID
	}

	id, err := uuid.Parse(rawID)
	if err != nil {
		return usecase.Actor{}, httpErrors.ErrInvalidJWTClaims(err)
	}

	return usecase.Actor{
		ID:      id,
		IsAdmin: claims.Role == domain.RoleAdmin,
	}, nil
}

// ─── POST /api-keys ───────────────────────────────────────────
func (h *handler) CreateAPIKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor, err := actorFromRequest(r)
		if err != nil {
			_ = render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		var req presenter.CreateAPIKeyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			_ = render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		key, err := h.uc.Create(r.Context(), req.Name, req.Description, req.ExpiresIn, actor)
		if err != nil {
			_ = render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		resp := presenter.APIKeyWithSecretResponse{
			APIKeyResponse: toResponse(key),
			APIKey:         key.PlainKey,
		}

		render.Status(r, http.StatusCreated)
		render.Respond(w, r, responses.CreateSuccessResponse(resp))
	}
}

// ─── GET /api-keys ────────────────────────────────────────────
func (h *handler) ListAPIKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor, err := actorFromRequest(r)
		if err != nil {
			_ = render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		limit := 20
		if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
			limit = l
		}
		if limit > 100 {
			limit = 100
		}
		offset := 0
		if o, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && o >= 0 {
			offset = o
		}

		items, total, err := h.uc.List(r.Context(), limit, offset, actor)
		if err != nil {
			_ = render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		resp := presenter.ListAPIKeysResponse{
			Items:  toResponseList(items),
			Total:  total,
			Limit:  limit,
			Offset: offset,
		}
		render.Respond(w, r, responses.CreateSuccessResponse(resp))
	}
}

// ─── GET /api-keys/{id} ───────────────────────────────────────
func (h *handler) GetAPIKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor, err := actorFromRequest(r)
		if err != nil {
			_ = render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			_ = render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		key, err := h.uc.GetByID(r.Context(), id, actor)
		if err != nil {
			_ = render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(toResponse(key)))
	}
}

// ─── POST /api-keys/{id}/rotate ───────────────────────────────
func (h *handler) RotateAPIKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor, err := actorFromRequest(r)
		if err != nil {
			_ = render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			_ = render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		key, err := h.uc.Rotate(r.Context(), id, actor)
		if err != nil {
			_ = render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		resp := presenter.APIKeyWithSecretResponse{
			APIKeyResponse: toResponse(key),
			APIKey:         key.PlainKey,
		}
		render.Respond(w, r, responses.CreateSuccessResponse(resp))
	}
}

// ─── DELETE /api-keys/{id} ────────────────────────────────────
func (h *handler) RevokeAPIKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor, err := actorFromRequest(r)
		if err != nil {
			_ = render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			_ = render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		if err := h.uc.Revoke(r.Context(), id, actor); err != nil {
			_ = render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse("API key revoked"))
	}
}

// ─── helpers ──────────────────────────────────────────────────
func toResponse(k *pkgapikey.Key) presenter.APIKeyResponse {
	var expiresAt, lastUsedAt *string
	if k.ExpiresAt != nil {
		s := k.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")
		expiresAt = &s
	}
	if k.LastUsedAt != nil {
		s := k.LastUsedAt.Format("2006-01-02T15:04:05Z07:00")
		lastUsedAt = &s
	}
	return presenter.APIKeyResponse{
		ID:          k.ID.String(),
		Name:        k.Name,
		Description: k.Description,
		Prefix:      k.Prefix,
		LastFour:    k.LastFour,
		ExpiresAt:   expiresAt,
		LastUsedAt:  lastUsedAt,
		IsActive:    k.IsActive,
		CreatedAt:   k.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   k.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toResponseList(items []*pkgapikey.Key) []presenter.APIKeyResponse {
	resp := make([]presenter.APIKeyResponse, len(items))
	for i, k := range items {
		resp[i] = toResponse(k)
	}
	return resp
}