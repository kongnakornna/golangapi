package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"icmongolang/config"
	"icmongolang/internal/middleware"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/customer"
	"icmongolang/internal/modules/customer/presenter"
	"icmongolang/pkg/helpers"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type customerHandler struct {
	cfg        *config.Config
	customerUC customer.CustomerUseCaseI
	carUC      customer.CarUseCaseI
	logger     logger.Logger
}

func CreateCustomerHandler(
	customerUC customer.CustomerUseCaseI,
	carUC customer.CarUseCaseI,
	cfg *config.Config,
	logger logger.Logger,
) customer.Handlers {
	return &customerHandler{
		cfg:        cfg,
		customerUC: customerUC,
		carUC:      carUC,
		logger:     logger,
	}
}

// Create godoc
// @Summary Create Customer
// @Description Create new customer.
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body presenter.CustomerCreate true "Add customer"
// @Success 200 {object} responses.SuccessResponse[presenter.CustomerResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 422 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /customer [post]
func (h *customerHandler) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := new(presenter.CustomerCreate)

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		if err := utils.ValidateStruct(ctx, req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		cust := mapCustomerCreate(req)
		cust.UserID = user.ID

		created, err := h.customerUC.Create(ctx, cust)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapCustomerResponse(created)))
	}
}

// Get godoc
// @Summary Read customer
// @Description Get customer by ID.
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer Id"
// @Success 200 {object} responses.SuccessResponse[presenter.CustomerResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 422 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /customer/{id} [get]
func (h *customerHandler) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		cust, err := h.customerUC.Get(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapCustomerResponse(cust)))
	}
}

// GetMulti godoc
// @Summary Read customers
// @Description Retrieve customers with pagination.
// @Tags customers
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param per_page query int false "Items per page (default 10)"
// @Success 200 {object} responses.SuccessResponse[presenter.PaginatedCustomersResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 422 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /customer [get]
func (h *customerHandler) GetMulti() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		q := r.URL.Query()

		limit := 10
		offset := 0
		page := 1
		perPage := 10

		pageStr := q.Get("page")
		if pageStr == "" {
			pageStr = q.Get("offset")
		}
		perPageStr := q.Get("per_page")
		if perPageStr == "" {
			perPageStr = q.Get("limit")
		}

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

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		var customers []*models.Customer
		var total int64

		if user.IsSuperUser {
			customers, err = h.customerUC.GetMulti(ctx, limit, offset)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
			total, err = h.customerUC.Count(ctx)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
		} else {
			customers, err = h.customerUC.GetMultiByUserID(ctx, user.ID, limit, offset)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
			total, err = h.customerUC.CountByUserID(ctx, user.ID)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
		}

		totalPages := int(total) / perPage
		if int(total)%perPage != 0 {
			totalPages++
		}
		if totalPages < 1 {
			totalPages = 1
		}

		paginatedRes := &presenter.PaginatedCustomersResponse{
			Customers:  mapCustomersResponse(customers),
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
		}
		render.Respond(w, r, responses.CreateSuccessResponse(paginatedRes))
	}
}

// Delete godoc
// @Summary Delete customer
// @Description Delete a customer by ID.
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer Id"
// @Success 200 {object} responses.SuccessResponse[presenter.CustomerResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 422 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /customer/{id} [delete]
func (h *customerHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		deleted, err := h.customerUC.Delete(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapCustomerResponse(deleted)))
	}
}

// Update godoc
// @Summary Update customer
// @Description Update a customer by ID.
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer Id"
// @Param customer body presenter.CustomerUpdate true "Update customer"
// @Success 200 {object} responses.SuccessResponse[presenter.CustomerResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 422 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /customer/{id} [put]
func (h *customerHandler) Update() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		req := new(presenter.CustomerUpdate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		values := make(map[string]interface{})
		if req.CustomerCode != nil {
			values["customer_code"] = *req.CustomerCode
		}
		if req.FullName != nil {
			values["full_name"] = *req.FullName
		}
		if req.DisplayName != nil {
			values["display_name"] = *req.DisplayName
		}
		if req.CustomerType != nil {
			values["customer_type"] = *req.CustomerType
		}
		if req.Status != nil {
			values["status"] = *req.Status
		}
		if req.TaxID != nil {
			values["tax_id"] = *req.TaxID
		}
		if req.Email != nil {
			values["email"] = *req.Email
		}
		if req.PhoneNumber != nil {
			values["phone_number"] = *req.PhoneNumber
		}
		if req.SecondaryPhone != nil {
			values["secondary_phone"] = *req.SecondaryPhone
		}
		if req.Address != nil {
			values["address"] = *req.Address
		}
		if req.Province != nil {
			values["province"] = *req.Province
		}
		if req.City != nil {
			values["city"] = *req.City
		}
		if req.District != nil {
			values["district"] = *req.District
		}
		if req.PostalCode != nil {
			values["postal_code"] = *req.PostalCode
		}
		if req.Country != nil {
			values["country"] = *req.Country
		}
		if req.ContactPerson != nil {
			values["contact_person"] = *req.ContactPerson
		}
		if req.ContactPhone != nil {
			values["contact_phone"] = *req.ContactPhone
		}
		if req.Notes != nil {
			values["notes"] = *req.Notes
		}

		updated, err := h.customerUC.Update(ctx, id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapCustomerResponse(updated)))
	}
}

// ─── Car Handlers ──────────────────────────────────────────────────────────

// CreateCar godoc
// @Summary Create Car
// @Description Add a new car to a customer.
// @Tags cars
// @Accept json
// @Produce json
// @Param car body presenter.CarCreate true "Add car"
// @Success 200 {object} responses.SuccessResponse[presenter.CarResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 422 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /car [post]
func (h *customerHandler) CreateCar() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := new(presenter.CarCreate)

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		if err := utils.ValidateStruct(ctx, req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		car := mapCarCreate(req)
		car.UserID = user.ID

		created, err := h.carUC.Create(ctx, car)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapCarResponse(created)))
	}
}

// GetCar godoc
// @Summary Read car
// @Description Get car by ID.
// @Tags cars
// @Accept json
// @Produce json
// @Param id path string true "Car Id"
// @Success 200 {object} responses.SuccessResponse[presenter.CarResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /car/{id} [get]
func (h *customerHandler) GetCar() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		car, err := h.carUC.Get(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapCarResponse(car)))
	}
}

// ListCars godoc
// @Summary List cars
// @Description List cars by customer ID with pagination.
// @Tags cars
// @Accept json
// @Produce json
// @Param customerId query string false "Filter by customer ID"
// @Param page query int false "Page number (default 1)"
// @Param per_page query int false "Items per page (default 10)"
// @Success 200 {object} responses.SuccessResponse[presenter.PaginatedCarsResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /car [get]
func (h *customerHandler) ListCars() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		q := r.URL.Query()

		limit := 10
		offset := 0
		page := 1
		perPage := 10

		pageStr := q.Get("page")
		if pageStr == "" {
			pageStr = q.Get("offset")
		}
		perPageStr := q.Get("per_page")
		if perPageStr == "" {
			perPageStr = q.Get("limit")
		}

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

		var cars []*models.Car
		var total int64
		var err error

		customerIDStr := q.Get("customerId")
		if customerIDStr != "" {
			var customerID uuid.UUID
			customerID, err = uuid.Parse(customerIDStr)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
				return
			}
			cars, err = h.carUC.GetMultiByCustomerID(ctx, customerID, limit, offset)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
			total, err = h.carUC.CountByCustomerID(ctx, customerID)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
		} else {
			cars, err = h.carUC.GetMulti(ctx, limit, offset)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
			total = int64(len(cars))
		}

		totalPages := int(total) / perPage
		if int(total)%perPage != 0 {
			totalPages++
		}
		if totalPages < 1 {
			totalPages = 1
		}

		paginatedRes := &presenter.PaginatedCarsResponse{
			Cars:       mapCarsResponse(cars),
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
		}
		render.Respond(w, r, responses.CreateSuccessResponse(paginatedRes))
	}
}

// UpdateCar godoc
// @Summary Update car
// @Description Update a car by ID.
// @Tags cars
// @Accept json
// @Produce json
// @Param id path string true "Car Id"
// @Param car body presenter.CarUpdate true "Update car"
// @Success 200 {object} responses.SuccessResponse[presenter.CarResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 422 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /car/{id} [put]
func (h *customerHandler) UpdateCar() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		req := new(presenter.CarUpdate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		values := make(map[string]interface{})
		if req.LicensePlate != nil {
			values["license_plate"] = *req.LicensePlate
		}
		if req.Province != nil {
			values["province"] = *req.Province
		}
		if req.Brand != nil {
			values["brand"] = *req.Brand
		}
		if req.Model != nil {
			values["model"] = *req.Model
		}
		if req.SubModel != nil {
			values["sub_model"] = *req.SubModel
		}
		if req.Year != nil {
			values["year"] = *req.Year
		}
		if req.Color != nil {
			values["color"] = *req.Color
		}
		if req.EngineNumber != nil {
			values["engine_number"] = *req.EngineNumber
		}
		if req.ChassisNumber != nil {
			values["chassis_number"] = *req.ChassisNumber
		}
		if req.FuelType != nil {
			values["fuel_type"] = *req.FuelType
		}
		if req.TransmissionType != nil {
			values["transmission_type"] = *req.TransmissionType
		}
		if req.EngineCC != nil {
			values["engine_cc"] = *req.EngineCC
		}
		if req.SeatingCapacity != nil {
			values["seating_capacity"] = *req.SeatingCapacity
		}
		if req.Mileage != nil {
			values["mileage"] = *req.Mileage
		}
		if req.LastServiceDate != nil {
			values["last_service_date"] = *req.LastServiceDate
		}
		if req.NextServiceMileage != nil {
			values["next_service_mileage"] = *req.NextServiceMileage
		}
		if req.Notes != nil {
			values["notes"] = *req.Notes
		}
		if req.IsActive != nil {
			values["is_active"] = *req.IsActive
		}

		updated, err := h.carUC.Update(ctx, id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapCarResponse(updated)))
	}
}

// DeleteCar godoc
// @Summary Delete car
// @Description Delete a car by ID.
// @Tags cars
// @Accept json
// @Produce json
// @Param id path string true "Car Id"
// @Success 200 {object} responses.SuccessResponse[presenter.CarResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /car/{id} [delete]
func (h *customerHandler) DeleteCar() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		deleted, err := h.carUC.Delete(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapCarResponse(deleted)))
	}
}

// ─── Mappers ──────────────────────────────────────────────────────────────

func mapCustomerCreate(req *presenter.CustomerCreate) *models.Customer {
	return &models.Customer{
		CustomerCode:   req.CustomerCode,
		FullName:       req.FullName,
		DisplayName:    req.DisplayName,
		CustomerType:   req.CustomerType,
		TaxID:          req.TaxID,
		Email:          req.Email,
		PhoneNumber:    req.PhoneNumber,
		SecondaryPhone: req.SecondaryPhone,
		Address:        req.Address,
		Province:       req.Province,
		City:           req.City,
		District:       req.District,
		PostalCode:     req.PostalCode,
		Country:        req.Country,
		ContactPerson:  req.ContactPerson,
		ContactPhone:   req.ContactPhone,
		Notes:          req.Notes,
	}
}

func mapCustomerResponse(cust *models.Customer) *presenter.CustomerResponse {
	return &presenter.CustomerResponse{
		ID:              cust.ID,
		CustomerCode:    cust.CustomerCode,
		FullName:        cust.FullName,
		DisplayName:     cust.DisplayName,
		CustomerType:    cust.CustomerType,
		Status:          cust.Status,
		TaxID:           cust.TaxID,
		Email:           cust.Email,
		PhoneNumber:     cust.PhoneNumber,
		SecondaryPhone:  cust.SecondaryPhone,
		Address:         cust.Address,
		Province:        cust.Province,
		City:            cust.City,
		District:        cust.District,
		PostalCode:      cust.PostalCode,
		Country:         cust.Country,
		ContactPerson:   cust.ContactPerson,
		ContactPhone:    cust.ContactPhone,
		Notes:           cust.Notes,
		LastVisitDate:   cust.LastVisitDate,
		TotalVisitCount: cust.TotalVisitCount,
		TotalSpent:      cust.TotalSpent,
		UserID:          cust.UserID,
		WhitelabelID:    cust.WhitelabelID,
		CreatedAt:       helpers.NewLocalTime(cust.CreatedAt),
		UpdatedAt:       helpers.NewLocalTimePtr(cust.UpdatedAt),
	}
}

func mapCustomersResponse(customers []*models.Customer) []*presenter.CustomerResponse {
	out := make([]*presenter.CustomerResponse, len(customers))
	for i, c := range customers {
		out[i] = mapCustomerResponse(c)
	}
	return out
}

func mapCarCreate(req *presenter.CarCreate) *models.Car {
	return &models.Car{
		CustomerID:         req.CustomerID,
		LicensePlate:       req.LicensePlate,
		Province:           req.Province,
		Brand:              req.Brand,
		Model:              req.Model,
		SubModel:           req.SubModel,
		Year:               req.Year,
		Color:              req.Color,
		EngineNumber:       req.EngineNumber,
		ChassisNumber:      req.ChassisNumber,
		FuelType:           req.FuelType,
		TransmissionType:   req.TransmissionType,
		EngineCC:           req.EngineCC,
		SeatingCapacity:    req.SeatingCapacity,
		Mileage:            req.Mileage,
		LastServiceDate:    req.LastServiceDate,
		NextServiceMileage: req.NextServiceMileage,
		Notes:              req.Notes,
	}
}

func mapCarResponse(car *models.Car) *presenter.CarResponse {
	return &presenter.CarResponse{
		ID:                 car.ID,
		CustomerID:         car.CustomerID,
		LicensePlate:       car.LicensePlate,
		Province:           car.Province,
		Brand:              car.Brand,
		Model:              car.Model,
		SubModel:           car.SubModel,
		Year:               car.Year,
		Color:              car.Color,
		EngineNumber:       car.EngineNumber,
		ChassisNumber:      car.ChassisNumber,
		FuelType:           car.FuelType,
		TransmissionType:   car.TransmissionType,
		EngineCC:           car.EngineCC,
		SeatingCapacity:    car.SeatingCapacity,
		Mileage:            car.Mileage,
		LastServiceDate:    car.LastServiceDate,
		NextServiceMileage: car.NextServiceMileage,
		Notes:              car.Notes,
		IsActive:           car.IsActive,
		UserID:             car.UserID,
		WhitelabelID:       car.WhitelabelID,
		CreatedAt:          helpers.NewLocalTime(car.CreatedAt),
		UpdatedAt:          helpers.NewLocalTimePtr(car.UpdatedAt),
	}
}

func mapCarsResponse(cars []*models.Car) []*presenter.CarResponse {
	out := make([]*presenter.CarResponse, len(cars))
	for i, c := range cars {
		out[i] = mapCarResponse(c)
	}
	return out
}
