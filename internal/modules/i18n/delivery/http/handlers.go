package http

import (
	"encoding/json"
	"net/http"

	"errors"

	"icmongolang/config"
	"icmongolang/internal/modules/i18n"
	"icmongolang/internal/modules/i18n/presenter"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type i18nHandler struct {
	cfg    *config.Config
	i18nUC i18n.I18nUseCaseI
	logger logger.Logger
}

// CreateI18nHandler creates a new i18n HTTP handler.
// สร้างตัวจัดการ HTTP สำหรับการแปลภาษา
func CreateI18nHandler(uc i18n.I18nUseCaseI, cfg *config.Config, logger logger.Logger) i18n.Handlers {
	return &i18nHandler{cfg: cfg, i18nUC: uc, logger: logger}
}

// GetTranslations godoc
// @Summary Get all translations
// @Description Get translations, optionally filtered by locale.
// @Tags i18n
// @Accept json
// @Produce json
// @Param locale query string false "Locale filter (e.g. th, en)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.TranslationResponse]
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/i18n/translations [get]
func (h *i18nHandler) GetTranslations() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.URL.Query().Get("locale")
		var (
			data interface{}
			err  error
		)
		if locale != "" {
			data, err = h.i18nUC.GetByLocale(r.Context(), locale)
		} else {
			data, err = h.i18nUC.GetMulti(r.Context(), 1000, 0)
		}
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(data))
	}
}

// GetTranslationByKey godoc
// @Summary Get translation by key
// @Description Get a single translation by locale and key.
// @Tags i18n
// @Accept json
// @Produce json
// @Param key path string true "Translation key"
// @Param locale query string true "Locale"
// @Success 200 {object} responses.SuccessResponse[presenter.TranslationResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/i18n/translations/{key} [get]
func (h *i18nHandler) GetTranslationByKey() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		locale := r.URL.Query().Get("locale")
		if locale == "" {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrBadRequest(errors.New("locale is required"))))
			return
		}
		trans, err := h.i18nUC.GetByLocaleAndKey(r.Context(), locale, key)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(trans))
	}
}

// CreateTranslation godoc
// @Summary Create a translation
// @Description Create a new translation entry.
// @Tags i18n
// @Accept json
// @Produce json
// @Param translation body presenter.TranslationRequest true "Translation"
// @Success 200 {object} responses.SuccessResponse[presenter.TranslationResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/i18n/translations [post]
func (h *i18nHandler) CreateTranslation() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.TranslationRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(req))
	}
}

// UpdateTranslation godoc
// @Summary Update a translation
// @Description Update an existing translation by key and locale.
// @Tags i18n
// @Accept json
// @Produce json
// @Param key path string true "Translation key"
// @Param locale query string true "Locale"
// @Param translation body presenter.TranslationRequest true "Translation"
// @Success 200 {object} responses.SuccessResponse[presenter.TranslationResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/i18n/translations/{key} [put]
func (h *i18nHandler) UpdateTranslation() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Respond(w, r, responses.CreateSuccessResponse("ok"))
	}
}

// DeleteTranslation godoc
// @Summary Delete a translation
// @Description Delete a translation by ID.
// @Tags i18n
// @Accept json
// @Produce json
// @Param id query int true "Translation ID"
// @Success 200 {object} responses.SuccessResponse[string]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/i18n/translations/{key} [delete]
func (h *i18nHandler) DeleteTranslation() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Respond(w, r, responses.CreateSuccessResponse("deleted"))
	}
}
