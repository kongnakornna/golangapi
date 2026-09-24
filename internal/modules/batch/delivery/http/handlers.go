package http

import (
	"encoding/json"
	"net/http"

	"icmongolang/config"
	"icmongolang/internal/modules/batch"
	"icmongolang/internal/modules/batch/presenter"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type batchHandler struct {
	cfg     *config.Config
	batchUC batch.BatchUseCaseI
	logger  logger.Logger
}

// CreateBatchHandler creates a new batch HTTP handler.
// สร้างตัวจัดการ HTTP สำหรับงานแบตช์
func CreateBatchHandler(uc batch.BatchUseCaseI, cfg *config.Config, logger logger.Logger) batch.Handlers {
	return &batchHandler{cfg: cfg, batchUC: uc, logger: logger}
}

// CreateJob godoc
// @Summary Create a batch job
// @Description Create a new batch processing job.
// @Tags batch
// @Accept json
// @Produce json
// @Param job body presenter.BatchJobRequest true "Batch job"
// @Success 200 {object} responses.SuccessResponse[presenter.BatchJobResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/batch/jobs [post]
func (h *batchHandler) CreateJob() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.BatchJobRequest)
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

// GetJob godoc
// @Summary Get a batch job
// @Description Get a batch job by ID.
// @Tags batch
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} responses.SuccessResponse[presenter.BatchJobResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/batch/jobs/{id} [get]
func (h *batchHandler) GetJob() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		job, err := h.batchUC.Get(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(job))
	}
}

// ListJobs godoc
// @Summary List batch jobs
// @Description List all batch jobs with pagination.
// @Tags batch
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} responses.SuccessResponse[[]presenter.BatchJobResponse]
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/batch/jobs [get]
func (h *batchHandler) ListJobs() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jobs, err := h.batchUC.GetMulti(r.Context(), 50, 0)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(jobs))
	}
}

// UpdateJob godoc
// @Summary Update a batch job
// @Description Update a batch job by ID.
// @Tags batch
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Param job body presenter.BatchJobRequest true "Batch job"
// @Success 200 {object} responses.SuccessResponse[presenter.BatchJobResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/batch/jobs/{id} [put]
func (h *batchHandler) UpdateJob() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Respond(w, r, responses.CreateSuccessResponse("ok"))
	}
}

// DeleteJob godoc
// @Summary Delete a batch job
// @Description Delete a batch job by ID.
// @Tags batch
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} responses.SuccessResponse[presenter.BatchJobResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/batch/jobs/{id} [delete]
func (h *batchHandler) DeleteJob() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		job, err := h.batchUC.Delete(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(job))
	}
}

// RunJobNow godoc
// @Summary Run a batch job now
// @Description Trigger immediate execution of a batch job.
// @Tags batch
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} responses.SuccessResponse[string]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/batch/jobs/{id}/run [post]
func (h *batchHandler) RunJobNow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := h.batchUC.RunJob(r.Context(), id); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse("Job triggered"))
	}
}

// GetJobLogs godoc
// @Summary Get batch job logs
// @Description Get logs for a specific batch job.
// @Tags batch
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} responses.SuccessResponse[[]presenter.BatchJobLogResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /api/batch/jobs/{id}/logs [get]
func (h *batchHandler) GetJobLogs() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		logs, err := h.batchUC.GetJobLogs(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(logs))
	}
}
