package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"icmongolang/internal/modules/vectordata/presenter"
	"icmongolang/internal/modules/vectordata/usecase"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/vectordb"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// VectorDataHandler handles HTTP requests for the VectorData module.
type VectorDataHandler struct {
	uc     usecase.VectorDataUseCase
	logger logger.Logger
}

// NewVectorDataHandler creates a new handler.
func NewVectorDataHandler(uc usecase.VectorDataUseCase, log logger.Logger) *VectorDataHandler {
	return &VectorDataHandler{uc: uc, logger: log}
}

// Health godoc
// @Summary      Vector database status
// @Description  Returns vector database provider and connectivity status
// @Tags         vectordata
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} errResponse
// @Router       /vectordata/health [get]
func (h *VectorDataHandler) Health(w http.ResponseWriter, r *http.Request) {
	status, err := h.uc.Health(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, status)
}

// CreateIndex godoc
// @Summary      Create vector index
// @Description  Creates (or ensures) the vector index/table in the configured backend
// @Tags         vectordata
// @Accept       json
// @Produce      json
// @Param        request body presenter.CreateIndexRequest true "Index dims"
// @Success      200 {object} map[string]string
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /vectordata/create-index [post]
func (h *VectorDataHandler) CreateIndex(w http.ResponseWriter, r *http.Request) {
	var req presenter.CreateIndexRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	if err := h.uc.CreateIndex(r.Context(), &req); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, map[string]string{"status": "ok", "message": "index ready"})
}

// Create godoc
// @Summary      Create a vector document
// @Description  Embeds content via LLM (or uses a provided embedding) and stores the vector
// @Tags         vectordata
// @Accept       json
// @Produce      json
// @Param        request body presenter.CreateDocumentRequest true "Document to index"
// @Success      200 {object} presenter.DocumentResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /vectordata/documents [post]
func (h *VectorDataHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req presenter.CreateDocumentRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	if strings.TrimSpace(req.Content) == "" && len(req.Embedding) == 0 {
		render.Render(w, r, ErrInvalidRequest(errors.New("content is required when embedding is not provided")))
		return
	}
	resp, err := h.uc.Create(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// Get godoc
// @Summary      Get a vector document by id
// @Description  Returns a stored document by its id
// @Tags         vectordata
// @Produce      json
// @Param        id path string true "Document id"
// @Success      200 {object} presenter.DocumentResponse
// @Failure      404 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /vectordata/documents/{id} [get]
func (h *VectorDataHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	resp, err := h.uc.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, vectordb.ErrNotFound) {
			render.Render(w, r, ErrNotFound(err))
		} else {
			render.Render(w, r, ErrInternal(err))
		}
		return
	}
	render.JSON(w, r, resp)
}

// List godoc
// @Summary      List vector documents
// @Description  Returns a paginated list of stored documents
// @Tags         vectordata
// @Produce      json
// @Param        limit query int false "Page size (default 10)"
// @Param        offset query int false "Offset (default 0)"
// @Success      200 {object} presenter.ListResponse
// @Failure      500 {object} errResponse
// @Router       /vectordata/documents [get]
func (h *VectorDataHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	resp, err := h.uc.List(r.Context(), limit, offset)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// Update godoc
// @Summary      Update a vector document
// @Description  Replaces content/embedding of an existing document by id
// @Tags         vectordata
// @Accept       json
// @Produce      json
// @Param        id path string true "Document id"
// @Param        request body presenter.UpdateDocumentRequest true "Updated fields"
// @Success      200 {object} presenter.DocumentResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /vectordata/documents/{id} [put]
func (h *VectorDataHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req presenter.UpdateDocumentRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	resp, err := h.uc.Update(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, vectordb.ErrNotFound) {
			render.Render(w, r, ErrNotFound(err))
		} else {
			render.Render(w, r, ErrInternal(err))
		}
		return
	}
	render.JSON(w, r, resp)
}

// Delete godoc
// @Summary      Delete a vector document
// @Description  Removes a document by id
// @Tags         vectordata
// @Produce      json
// @Param        id path string true "Document id"
// @Success      200 {object} map[string]string
// @Failure      404 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /vectordata/documents/{id} [delete]
func (h *VectorDataHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.uc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, vectordb.ErrNotFound) {
			render.Render(w, r, ErrNotFound(err))
		} else {
			render.Render(w, r, ErrInternal(err))
		}
		return
	}
	render.JSON(w, r, map[string]string{"status": "ok", "message": "deleted"})
}

// Search godoc
// @Summary      Semantic vector search
// @Description  Embeds the query via the LLM and returns the most similar documents
// @Tags         vectordata
// @Accept       json
// @Produce      json
// @Param        request body presenter.SearchRequest true "Search query"
// @Success      200 {object} presenter.SearchResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /vectordata/search [post]
func (h *VectorDataHandler) Search(w http.ResponseWriter, r *http.Request) {
	var req presenter.SearchRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	if strings.TrimSpace(req.Query) == "" {
		render.Render(w, r, ErrInvalidRequest(errors.New("query is required")))
		return
	}
	resp, err := h.uc.Search(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// Seed godoc
// @Summary      Seed sample vector documents
// @Description  Embeds and indexes a set of sample documents (IoT sensor manual + alarm policy)
// @Tags         vectordata
// @Accept       json
// @Produce      json
// @Param        request body presenter.SeedRequest true "Number of sample docs to seed"
// @Success      200 {object} presenter.SeedResponse
// @Failure      500 {object} errResponse
// @Router       /vectordata/seed [post]
func (h *VectorDataHandler) Seed(w http.ResponseWriter, r *http.Request) {
	var req presenter.SeedRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	resp, err := h.uc.Seed(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// --- error helpers (same as other modules) ---

func ErrInvalidRequest(err error) render.Renderer {
	return &errResponse{HTTPStatusCode: http.StatusBadRequest, ErrorText: err.Error()}
}
func ErrInternal(err error) render.Renderer {
	return &errResponse{HTTPStatusCode: http.StatusInternalServerError, ErrorText: err.Error()}
}
func ErrNotFound(err error) render.Renderer {
	return &errResponse{HTTPStatusCode: http.StatusNotFound, ErrorText: err.Error()}
}

type errResponse struct {
	HTTPStatusCode int    `json:"-"`
	ErrorText      string `json:"error"`
}

func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
