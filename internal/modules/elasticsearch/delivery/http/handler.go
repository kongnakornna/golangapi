package http

import (
	"errors"
	"net/http"
	"strings"

	"icmongolang/internal/modules/elasticsearch/presenter"
	"icmongolang/internal/modules/elasticsearch/usecase"
	"icmongolang/pkg/logger"

	"github.com/go-chi/render"
)

// ElasticsearchHandler handles HTTP requests for the Elasticsearch module.
type ElasticsearchHandler struct {
	uc     usecase.ElasticsearchUseCase
	logger logger.Logger
}

// NewElasticsearchHandler creates a new handler.
func NewElasticsearchHandler(uc usecase.ElasticsearchUseCase, log logger.Logger) *ElasticsearchHandler {
	return &ElasticsearchHandler{uc: uc, logger: log}
}

// Health godoc
// @Summary      Elasticsearch cluster health
// @Description  Returns Elasticsearch cluster health status
// @Tags         elasticsearch
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} errResponse
// @Router       /elasticsearch/health [get]
func (h *ElasticsearchHandler) Health(w http.ResponseWriter, r *http.Request) {
	status, err := h.uc.Health(r.Context())
	if err != nil {
		h.internalError(w, r, "health", err)
		return
	}
	render.JSON(w, r, status)
}

// CreateIndex godoc
// @Summary      Create vector index
// @Description  Creates (or ensures) the dense_vector index in Elasticsearch
// @Tags         elasticsearch
// @Accept       json
// @Produce      json
// @Param        request body presenter.CreateIndexRequest true "Index dims"
// @Success      200 {object} map[string]string
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /elasticsearch/create-index [post]
func (h *ElasticsearchHandler) CreateIndex(w http.ResponseWriter, r *http.Request) {
	var req presenter.CreateIndexRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	if err := h.uc.CreateIndex(r.Context(), &req); err != nil {
		h.internalError(w, r, "create-index", err)
		return
	}
	render.JSON(w, r, map[string]string{"status": "ok", "message": "index ready"})
}

// EmbedAndIndex godoc
// @Summary      Embed content via LLM and index into Elasticsearch
// @Description  Calls the configured LLM to create an embedding and stores the vector in Elasticsearch
// @Tags         elasticsearch
// @Accept       json
// @Produce      json
// @Param        request body presenter.EmbeddingRequest true "Content to embed"
// @Success      200 {object} presenter.EmbeddingResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /elasticsearch/embeddings [post]
func (h *ElasticsearchHandler) EmbedAndIndex(w http.ResponseWriter, r *http.Request) {
	var req presenter.EmbeddingRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	if strings.TrimSpace(req.Content) == "" {
		render.Render(w, r, ErrInvalidRequest(errors.New("content is required")))
		return
	}
	resp, err := h.uc.EmbedAndIndex(r.Context(), &req)
	if err != nil {
		h.internalError(w, r, "embed-and-index", err)
		return
	}
	render.JSON(w, r, resp)
}

// IndexVector godoc
// @Summary      Index a pre-computed vector
// @Description  Stores an already-computed embedding directly in Elasticsearch
// @Tags         elasticsearch
// @Accept       json
// @Produce      json
// @Param        request body presenter.IndexVectorRequest true "Vector to index"
// @Success      200 {object} presenter.EmbeddingResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /elasticsearch/vectors [post]
func (h *ElasticsearchHandler) IndexVector(w http.ResponseWriter, r *http.Request) {
	var req presenter.IndexVectorRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	if strings.TrimSpace(req.Content) == "" || len(req.Embedding) == 0 {
		render.Render(w, r, ErrInvalidRequest(errors.New("content and embedding are required")))
		return
	}
	resp, err := h.uc.IndexVector(r.Context(), &req)
	if err != nil {
		h.internalError(w, r, "index-vector", err)
		return
	}
	render.JSON(w, r, resp)
}

// SearchVector godoc
// @Summary      Semantic vector search
// @Description  Embeds the query via the LLM and returns the most similar documents using cosine similarity
// @Tags         elasticsearch
// @Accept       json
// @Produce      json
// @Param        request body presenter.SearchRequest true "Search query"
// @Success      200 {object} presenter.SearchResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /elasticsearch/search [post]
func (h *ElasticsearchHandler) SearchVector(w http.ResponseWriter, r *http.Request) {
	var req presenter.SearchRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	if strings.TrimSpace(req.Query) == "" {
		render.Render(w, r, ErrInvalidRequest(errors.New("query is required")))
		return
	}
	resp, err := h.uc.SearchVector(r.Context(), &req)
	if err != nil {
		h.internalError(w, r, "search", err)
		return
	}
	render.JSON(w, r, resp)
}

// --- error helpers (same as influxdb module) ---

// internalError logs the full error server-side and returns a generic message
// to the client so backend/cluster internals are never leaked.
func (h *ElasticsearchHandler) internalError(w http.ResponseWriter, r *http.Request, op string, err error) {
	h.logger.Errorf("elasticsearch %s: %v", op, err)
	render.Render(w, r, &errResponse{HTTPStatusCode: http.StatusInternalServerError, ErrorText: "internal server error"})
}

func ErrInvalidRequest(err error) render.Renderer {
	return &errResponse{HTTPStatusCode: http.StatusBadRequest, ErrorText: err.Error()}
}

type errResponse struct {
	HTTPStatusCode int    `json:"-"`
	ErrorText      string `json:"error"`
}

func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
