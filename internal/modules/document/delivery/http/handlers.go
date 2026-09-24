package http

import (
	"encoding/json"
	"net/http"

	"icmongolang/config"
	"icmongolang/internal/modules/document"
	"icmongolang/internal/modules/document/presenter"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type documentHandler struct {
	cfg    *config.Config
	docUC  document.DocumentUseCaseI
	logger logger.Logger
}

// CreateDocumentHandler creates a new document HTTP handler.
// สร้างตัวจัดการ HTTP สำหรับเอกสาร
func CreateDocumentHandler(uc document.DocumentUseCaseI, cfg *config.Config, logger logger.Logger) document.Handlers {
	return &documentHandler{cfg: cfg, docUC: uc, logger: logger}
}

// Upload godoc
// @Summary Upload a document
// @Description Upload a new document file.
// @Tags documents
// @Accept json
// @Produce json
// @Param document body presenter.DocumentUploadRequest true "Document metadata"
// @Success 200 {object} responses.SuccessResponse[presenter.DocumentResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/documents [post]
func (h *documentHandler) Upload() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.DocumentUploadRequest)
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

// Download godoc
// @Summary Download a document
// @Description Get a document by ID.
// @Tags documents
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} responses.SuccessResponse[presenter.DocumentResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/documents/{id} [get]
func (h *documentHandler) Download() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		doc, err := h.docUC.Get(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(doc))
	}
}

// List godoc
// @Summary List documents
// @Description List all documents with pagination.
// @Tags documents
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} responses.SuccessResponse[[]presenter.DocumentResponse]
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/documents [get]
func (h *documentHandler) List() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		docs, err := h.docUC.GetMulti(r.Context(), 50, 0)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(docs))
	}
}

// Delete godoc
// @Summary Delete a document
// @Description Delete a document by ID.
// @Tags documents
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} responses.SuccessResponse[presenter.DocumentResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/documents/{id} [delete]
func (h *documentHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		doc, err := h.docUC.Delete(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(doc))
	}
}
