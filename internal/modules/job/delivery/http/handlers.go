package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"icmongolang/config"
	"icmongolang/internal/middleware"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/job"
	"icmongolang/internal/modules/job/presenter"
	"icmongolang/pkg/helpers"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/report"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type jobHandler struct {
	cfg    *config.Config
	jobUC  job.JobUseCaseI
	logger logger.Logger
}

func CreateJobHandler(uc job.JobUseCaseI, cfg *config.Config, logger logger.Logger) job.Handlers {
	return &jobHandler{cfg: cfg, jobUC: uc, logger: logger}
}

// Create godoc
// @Summary Create a job
// @Description Create a new repair job card.
// @Tags job
// @Accept json
// @Produce json
// @Param job body presenter.JobCreate true "Job details"
// @Success 200 {object} responses.SuccessResponse[presenter.JobResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 422 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /job [post]
func (h *jobHandler) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := new(presenter.JobCreate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := utils.ValidateStruct(ctx, req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		model := mapCreateRequest(req, user.ID, uuid.Nil)
		created, err := h.jobUC.Create(ctx, model)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapJobResponse(created)))
	}
}

// GetByID godoc
// @Summary Get a job
// @Description Get a job by ID.
// @Tags job
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} responses.SuccessResponse[presenter.JobResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /job/{id} [get]
func (h *jobHandler) GetByID() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		jobObj, err := h.jobUC.Get(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapJobResponse(jobObj)))
	}
}

// List godoc
// @Summary List jobs
// @Description List jobs with pagination.
// @Tags job
// @Accept json
// @Produce json
// @Param page query int false "Page number" Format(page)
// @Param per_page query int false "Items per page" Format(per_page)
// @Success 200 {object} responses.SuccessResponse[presenter.PaginatedJobResponse]
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /job [get]
func (h *jobHandler) List() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		q := r.URL.Query()

		limit := 10
		offset := 0
		page := 1
		perPage := 10

		pageStr := q.Get("page")
		perPageStr := q.Get("per_page")

		if pageStr != "" || perPageStr != "" {
			if pageStr != "" {
				if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
					page = p
				}
			}
			if perPageStr != "" {
				if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 {
					perPage = pp
				}
			}
			const maxPerPage = 100
			if perPage > maxPerPage {
				perPage = maxPerPage
			}
			limit = perPage
			offset = (page - 1) * perPage
		} else {
			if l := q.Get("limit"); l != "" {
				if lim, err := strconv.Atoi(l); err == nil && lim > 0 {
					limit = lim
				}
			}
			if o := q.Get("offset"); o != "" {
				if off, err := strconv.Atoi(o); err == nil && off >= 0 {
					offset = off
				}
			}
			if limit > 0 {
				perPage = limit
				page = (offset / limit) + 1
			}
		}

		jobs, err := h.jobUC.GetMulti(ctx, limit, offset)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		total, err := h.jobUC.Count(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		totalPages := int(total) / perPage
		if int(total)%perPage != 0 {
			totalPages++
		}
		if totalPages < 1 {
			totalPages = 1
		}

		paginatedRes := &presenter.PaginatedJobResponse{
			Jobs:       mapJobsResponse(jobs),
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
		}
		render.Respond(w, r, responses.CreateSuccessResponse(paginatedRes))
	}
}

// Update godoc
// @Summary Update a job
// @Description Update a job by ID.
// @Tags job
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Param job body presenter.JobUpdate true "Job update fields"
// @Success 200 {object} responses.SuccessResponse[presenter.JobResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /job/{id} [put]
func (h *jobHandler) Update() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		req := new(presenter.JobUpdate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		values := make(map[string]interface{})
		if req.MechanicID != nil {
			values["mechanic_id"] = req.MechanicID.String()
		}
		if req.Symptom != nil {
			values["symptom"] = *req.Symptom
		}
		if req.DiagnosisNote != nil {
			values["diagnosis_note"] = *req.DiagnosisNote
		}
		if req.Mileage != nil {
			values["mileage"] = *req.Mileage
		}
		if req.EstimatedCost != nil {
			values["estimated_cost"] = *req.EstimatedCost
		}
		if req.ActualCost != nil {
			values["actual_cost"] = *req.ActualCost
		}
		if req.Priority != nil {
			values["priority"] = *req.Priority
		}

		updated, err := h.jobUC.Update(ctx, id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapJobResponse(updated)))
	}
}

// Delete godoc
// @Summary Delete a job
// @Description Soft-delete a job by ID.
// @Tags job
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} responses.SuccessResponse[presenter.JobResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /job/{id} [delete]
func (h *jobHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		jobObj, err := h.jobUC.Delete(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapJobResponse(jobObj)))
	}
}

// ChangeStatus godoc
// @Summary Change job status
// @Description Update the status of a job.
// @Tags job
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Param status body presenter.JobStatusChange true "New status"
// @Success 200 {object} responses.SuccessResponse[presenter.JobResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /job/{id}/status [put]
func (h *jobHandler) ChangeStatus() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		req := new(presenter.JobStatusChange)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := utils.ValidateStruct(ctx, req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		updated, err := h.jobUC.ChangeStatus(ctx, id, req.Status, req.Reason)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapJobResponse(updated)))
	}
}

// AddService godoc
// @Summary Add service to job
// @Description Add a service item to a job.
// @Tags job
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Param service body presenter.JobServiceRequest true "Service details"
// @Success 200 {object} responses.SuccessResponse[presenter.JobServiceResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /job/{id}/services [post]
func (h *jobHandler) AddService() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		jobID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		req := new(presenter.JobServiceRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := utils.ValidateStruct(ctx, req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		svc := &models.JobService{
			JobID:        jobID,
			ServiceID:    req.ServiceID,
			Quantity:     req.Quantity,
			UnitPrice:    req.UnitPrice,
			Discount:     req.Discount,
			Note:         req.Note,
			UserID:       user.ID,
			WhitelabelID: uuid.Nil,
		}
		if svc.Quantity == 0 {
			svc.Quantity = 1
		}

		created, err := h.jobUC.AddService(ctx, svc)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapJobServiceResponse(created)))
	}
}

// AddPart godoc
// @Summary Add part to job
// @Description Add a part to a job.
// @Tags job
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Param part body presenter.JobPartRequest true "Part details"
// @Success 200 {object} responses.SuccessResponse[presenter.JobPartResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /job/{id}/parts [post]
func (h *jobHandler) AddPart() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		jobID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		req := new(presenter.JobPartRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := utils.ValidateStruct(ctx, req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		part := &models.JobPartSales{
			JobID:        jobID,
			PartID:       req.PartID,
			Quantity:     req.Quantity,
			UnitPrice:    req.UnitPrice,
			Discount:     req.Discount,
			Note:         req.Note,
			UserID:       user.ID,
			WhitelabelID: uuid.Nil,
		}
		if part.Quantity == 0 {
			part.Quantity = 1
		}

		created, err := h.jobUC.AddPart(ctx, part)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapJobPartResponse(created)))
	}
}

// GetReport godoc
// @Summary Get job report
// @Description Get a full job report with services and parts.
// @Tags job
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} responses.SuccessResponse[presenter.JobResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /job/{id}/report [get]
func (h *jobHandler) GetReport() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		report, err := h.jobUC.GetReport(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapJobResponse(report)))
	}
}

// GetStatusHistory godoc
// @Summary Get job status history
// @Description Get status change history for a job.
// @Tags job
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} responses.SuccessResponse[[]presenter.JobStatusHistoryResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /job/{id}/history [get]
func (h *jobHandler) GetStatusHistory() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		history, err := h.jobUC.GetStatusHistory(r.Context(), jobID)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapHistoryResponse(history)))
	}
}

// GetServices godoc
// @Summary Get job services
// @Description Get all services for a job.
// @Tags job
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} responses.SuccessResponse[[]presenter.JobServiceResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /job/{id}/services [get]
func (h *jobHandler) GetServices() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		services, err := h.jobUC.GetServices(r.Context(), jobID)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapServicesResponse(services)))
	}
}

// GetParts godoc
// @Summary Get job parts
// @Description Get all parts for a job.
// @Tags job
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} responses.SuccessResponse[[]presenter.JobPartResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /job/{id}/parts [get]
func (h *jobHandler) GetParts() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		parts, err := h.jobUC.GetParts(r.Context(), jobID)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapPartsResponse(parts)))
	}
}

// GetPDF godoc
// @Summary Generate Job Card PDF
// @Description Generate job card PDF for a job.
// @Tags Jobs
// @Produce application/pdf
// @Param id path string true "Job ID"
// @Success 200 {file} byte "PDF file"
// @Failure 400 {object} responses.ErrorResponse
// @Router /jobs/{id}/pdf [get]
func (h *jobHandler) GetPDF() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		ctx := r.Context()
		jobObj, err := h.jobUC.GetReport(ctx, id)
		if err != nil {
			h.logger.Errorf("Failed to get job for PDF: %v", err)
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotFound(err)))
			return
		}

		services, _ := h.jobUC.GetServices(ctx, id)
		parts, _ := h.jobUC.GetParts(ctx, id)

		svcItems := make([]report.JobCardService, len(services))
		for i, s := range services {
			svcItems[i] = report.JobCardService{
				LineNo:      i + 1,
				Description: s.ServiceID.String(),
				Quantity:    s.Quantity,
				UnitPrice:   s.UnitPrice,
				TotalPrice:  s.UnitPrice * float64(s.Quantity),
			}
		}

		partItems := make([]report.JobCardPart, len(parts))
		for i, p := range parts {
			partItems[i] = report.JobCardPart{
				LineNo:      i + 1,
				PartCode:    p.PartID.String(),
				Description: orDefaultStrPtr(p.Note),
				Quantity:    p.Quantity,
				UnitPrice:   p.UnitPrice,
				TotalPrice:  p.UnitPrice * float64(p.Quantity),
			}
		}

		var mileage int
		if jobObj.Mileage != nil {
			mileage = *jobObj.Mileage
		}
		var diagnosis string
		if jobObj.DiagnosisNote != nil {
			diagnosis = *jobObj.DiagnosisNote
		}

		data := report.JobCardData{
			Company: report.CompanyInfo{
				Name:    "ICMON Auto Repair",
				Address: "123/4 ถนนสุขุมวิท แขวงคลองเตย เขตคลองเตย กรุงเทพฯ 10110",
				Phone:   "02-123-4567",
				TaxID:   "0123456789012",
			},
			JobNo:        jobObj.JobNo,
			Date:         jobObj.CreatedAt,
			CustomerName: orDefault(jobObj.CustomerID.String()),
			LicensePlate: jobObj.CarID.String(),
			Mileage:      mileage,
			Status:       jobObj.Status,
			Diagnosis:    diagnosis,
			Services:     svcItems,
			Parts:        partItems,
			Mechanic:     orDefault(jobObj.MechanicID.String()),
		}

		pdf, err := report.GeneratePDF(ctx, report.TplJobCard, data)
		if err != nil {
			h.logger.Errorf("Failed to generate job card PDF: %v", err)
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrInternalServer(err)))
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "inline; filename=job_card_"+idStr+".pdf")
		w.Write(pdf)
	}
}

// GetPickingPDF godoc
// @Summary Generate Part Picking PDF
// @Description Generate part picking list PDF for a job.
// @Tags Jobs
// @Produce application/pdf
// @Param id path string true "Job ID"
// @Success 200 {file} byte "PDF file"
// @Failure 400 {object} responses.ErrorResponse
// @Router /jobs/{id}/picking/pdf [get]
func (h *jobHandler) GetPickingPDF() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		ctx := r.Context()
		parts, err := h.jobUC.GetParts(ctx, id)
		if err != nil {
			h.logger.Errorf("Failed to get parts for picking PDF: %v", err)
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotFound(err)))
			return
		}

		items := make([]report.PartPickingItem, len(parts))
		for i, p := range parts {
			items[i] = report.PartPickingItem{
				LineNo:     i + 1,
				PartCode:   p.PartID.String(),
				PartName:   orDefaultStrPtr(p.Note),
				QtyRequest: p.Quantity,
				QtyPicked:  0,
			}
		}

		data := report.PartPickingData{
			Company: report.CompanyInfo{
				Name:    "ICMON Auto Repair",
				Address: "123/4 ถนนสุขุมวิท แขวงคลองเตย เขตคลองเตย กรุงเทพฯ 10110",
				Phone:   "02-123-4567",
				TaxID:   "0123456789012",
			},
			PickingNo:   "PK-" + idStr[:8],
			RequestDate: time.Now().Format("02/01/2006"),
			JobNo:       idStr,
			Items:       items,
		}

		pdf, err := report.GeneratePDF(ctx, report.TplPartPicking, data)
		if err != nil {
			h.logger.Errorf("Failed to generate picking PDF: %v", err)
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrInternalServer(err)))
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "inline; filename=picking_"+idStr+".pdf")
		w.Write(pdf)
	}
}

// GetDeliveryPDF godoc
// @Summary Generate Delivery Sheet PDF
// @Description Generate delivery sheet PDF for a job.
// @Tags Jobs
// @Produce application/pdf
// @Param id path string true "Job ID"
// @Success 200 {file} byte "PDF file"
// @Failure 400 {object} responses.ErrorResponse
// @Router /jobs/{id}/delivery/pdf [get]
func (h *jobHandler) GetDeliveryPDF() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		ctx := r.Context()
		jobObj, err := h.jobUC.GetReport(ctx, id)
		if err != nil {
			h.logger.Errorf("Failed to get job for delivery PDF: %v", err)
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotFound(err)))
			return
		}

		services, _ := h.jobUC.GetServices(ctx, id)
		parts, _ := h.jobUC.GetParts(ctx, id)

		deliveryDate := time.Now()
		if jobObj.UpdatedAt != nil {
			deliveryDate = *jobObj.UpdatedAt
		}

		items := []report.DeliveryItem{}
		for i, s := range services {
			items = append(items, report.DeliveryItem{
				LineNo:      i + 1,
				Description: "บริการ: " + s.ServiceID.String(),
				Quantity:    s.Quantity,
			})
		}
		for i, p := range parts {
			items = append(items, report.DeliveryItem{
				LineNo:      len(services) + i + 1,
				Description: "อะไหล่: " + p.PartID.String(),
				Quantity:    p.Quantity,
				Unit:        "ชิ้น",
			})
		}

		data := report.DeliverySheetData{
			Company: report.CompanyInfo{
				Name:    "ICMON Auto Repair",
				Address: "123/4 ถนนสุขุมวิท แขวงคลองเตย เขตคลองเตย กรุงเทพฯ 10110",
				Phone:   "02-123-4567",
				TaxID:   "0123456789012",
			},
			DeliveryNo:   "DEL-" + idStr[:8],
			Date:         deliveryDate,
			JobNo:        jobObj.JobNo,
			CustomerName: jobObj.CustomerID.String(),
			LicensePlate: jobObj.CarID.String(),
			Items:        items,
			Remark:       "ตรวจสอบความเรียบร้อยก่อนส่งมอบ",
		}

		pdf, err := report.GeneratePDF(ctx, report.TplDeliverySheet, data)
		if err != nil {
			h.logger.Errorf("Failed to generate delivery PDF: %v", err)
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrInternalServer(err)))
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "inline; filename=delivery_"+idStr+".pdf")
		w.Write(pdf)
	}
}

func orDefault(s string) string {
	if s == "" || s == "00000000-0000-0000-0000-000000000000" {
		return "-"
	}
	return s
}

func orDefaultStrPtr(s *string) string {
	if s == nil || *s == "" {
		return "-"
	}
	return *s
}

func mapCreateRequest(req *presenter.JobCreate, userID uuid.UUID, whitelabelID uuid.UUID) *models.Job {
	priority := req.Priority
	if priority == "" {
		priority = "NORMAL"
	}
	return &models.Job{
		JobNo:         req.JobNo,
		CustomerID:    req.CustomerID,
		CarID:         req.CarID,
		MechanicID:    req.MechanicID,
		Symptom:       req.Symptom,
		DiagnosisNote: req.DiagnosisNote,
		Mileage:       req.Mileage,
		EstimatedCost: req.EstimatedCost,
		Priority:      priority,
		UserID:        userID,
		WhitelabelID:  whitelabelID,
	}
}

func mapJobResponse(jobObj *models.Job) *presenter.JobResponse {
	return &presenter.JobResponse{
		ID:            jobObj.ID,
		JobNo:         jobObj.JobNo,
		CustomerID:    jobObj.CustomerID,
		CarID:         jobObj.CarID,
		MechanicID:    jobObj.MechanicID,
		Status:        jobObj.Status,
		StartDate:     jobObj.StartDate,
		EndDate:       jobObj.EndDate,
		Symptom:       jobObj.Symptom,
		DiagnosisNote: jobObj.DiagnosisNote,
		Mileage:       jobObj.Mileage,
		EstimatedCost: jobObj.EstimatedCost,
		ActualCost:    jobObj.ActualCost,
		Priority:      jobObj.Priority,
		UserID:        jobObj.UserID,
		WhitelabelID:  jobObj.WhitelabelID,
		CreatedAt:     helpers.NewLocalTime(jobObj.CreatedAt),
		UpdatedAt:     helpers.NewLocalTimePtr(jobObj.UpdatedAt),
	}
}

func mapJobsResponse(jobs []*models.Job) []*presenter.JobResponse {
	out := make([]*presenter.JobResponse, len(jobs))
	for i, j := range jobs {
		out[i] = mapJobResponse(j)
	}
	return out
}

func mapJobServiceResponse(svc *models.JobService) *presenter.JobServiceResponse {
	return &presenter.JobServiceResponse{
		ID:        svc.ID,
		JobID:     svc.JobID,
		ServiceID: svc.ServiceID,
		Quantity:  svc.Quantity,
		UnitPrice: svc.UnitPrice,
		Discount:  svc.Discount,
		Note:      svc.Note,
		CreatedAt: helpers.NewLocalTime(svc.CreatedAt),
	}
}

func mapJobPartResponse(part *models.JobPartSales) *presenter.JobPartResponse {
	return &presenter.JobPartResponse{
		ID:        part.ID,
		JobID:     part.JobID,
		PartID:    part.PartID,
		Quantity:  part.Quantity,
		UnitPrice: part.UnitPrice,
		Discount:  part.Discount,
		Note:      part.Note,
		CreatedAt: helpers.NewLocalTime(part.CreatedAt),
	}
}

func mapHistoryResponse(history []*models.JobStatusHistory) []*presenter.JobStatusHistoryResponse {
	out := make([]*presenter.JobStatusHistoryResponse, len(history))
	for i, h := range history {
		out[i] = &presenter.JobStatusHistoryResponse{
			ID:         h.ID,
			JobID:      h.JobID,
			FromStatus: h.FromStatus,
			ToStatus:   h.ToStatus,
			ChangedBy:  h.ChangedBy,
			ChangedAt:  h.ChangedAt,
			Reason:     h.Reason,
		}
	}
	return out
}

func mapServicesResponse(services []*models.JobService) []*presenter.JobServiceResponse {
	out := make([]*presenter.JobServiceResponse, len(services))
	for i, s := range services {
		out[i] = mapJobServiceResponse(s)
	}
	return out
}

func mapPartsResponse(parts []*models.JobPartSales) []*presenter.JobPartResponse {
	out := make([]*presenter.JobPartResponse, len(parts))
	for i, p := range parts {
		out[i] = mapJobPartResponse(p)
	}
	return out
}
