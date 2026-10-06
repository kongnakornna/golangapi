package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"icmongolang/config"
	"icmongolang/internal/middleware"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/users"
	usermodels "icmongolang/internal/modules/users/models"
	"icmongolang/internal/modules/users/presenter"
	"icmongolang/pkg/helpers"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// ============================================================================
// userHandler — aggregator ของ UC ทุกตัวในโมดูล users
// ============================================================================
type userHandler struct {
	cfg          *config.Config
	usersUC      users.UserUseCaseI
	accessMenuUC users.AccessMenuUseCaseI
	fileUC       users.FileUseCaseI
	logUC        users.LogUseCaseI
	logTypeUC    users.LogTypeUseCaseI
	menuUC       users.MenuUseCaseI
	moduleUC     users.ModuleUseCaseI
	logger       logger.Logger
}

// CreateUserHandler constructs the aggregator handler.
func CreateUserHandler(
	uc users.UserUseCaseI,
	accessMenuUC users.AccessMenuUseCaseI,
	fileUC users.FileUseCaseI,
	logUC users.LogUseCaseI,
	logTypeUC users.LogTypeUseCaseI,
	menuUC users.MenuUseCaseI,
	moduleUC users.ModuleUseCaseI,
	cfg *config.Config,
	l logger.Logger,
) users.Handlers {
	return &userHandler{
		cfg:          cfg,
		usersUC:      uc,
		accessMenuUC: accessMenuUC,
		fileUC:       fileUC,
		logUC:        logUC,
		logTypeUC:    logTypeUC,
		menuUC:       menuUC,
		moduleUC:     moduleUC,
		logger:       l,
	}
}

// ============================================================================
// AUTH
// ============================================================================

// SignInEmail – POST /api/auth/signin
// @Summary      User login with email and password
// @Description  Authenticates a user using email and password. Returns access and refresh tokens.
// @Tags         auth
// @Accept       json,multipart/form-data
// @Produce      json
// @Param        request body presenter.UserSignIn true "User credentials"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=presenter.Token} "Tokens returned"
// @Failure      400  {object}  responses.ErrorResponse "Bad request / validation error"
// @Failure      401  {object}  responses.ErrorResponse "Unauthorized (wrong email/password)"
// @Router       /auth/signin [post]
func (h *userHandler) SignInEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.UserSignIn)
		contentType := r.Header.Get("Content-Type")

		var err error
		if strings.Contains(contentType, "application/json") {
			err = json.NewDecoder(r.Body).Decode(req)
		} else if strings.Contains(contentType, "multipart/form-data") {
			err = utils.DecodeFormToStruct(r, req)
		} else {
			err = errors.New("unsupported content type")
		}
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		accessToken, refreshToken, err := h.usersUC.SignIn(r.Context(), req.Email, req.Password)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		// audit log (IP/Device/Browser)
		if u, err := h.usersUC.GetByEmail(r.Context(), req.Email); err == nil && u != nil {
			ip, device, browser := extractClientMeta(r)
			h.logUC.WriteAuditLogWithMeta(r.Context(), u.ID, 1, "User SignIn", "Sign in via /auth/signin", ip, device, browser)
		}

		render.Respond(w, r, responses.CreateSuccessResponse(presenter.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
		}))
	}
}

// SignInUsername – POST /api/auth/login
// @Summary      User login with username and password
// @Tags         auth
// @Accept       json,multipart/form-data
// @Produce      json
// @Param        request body presenter.SignIn true "User credentials (username + password)"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=presenter.Token} "Tokens returned"
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Router       /auth/login [post]
func (h *userHandler) SignInUsername() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.SignIn)
		contentType := r.Header.Get("Content-Type")

		var err error
		if strings.Contains(contentType, "application/json") {
			err = json.NewDecoder(r.Body).Decode(req)
		} else if strings.Contains(contentType, "multipart/form-data") {
			err = utils.DecodeFormToStruct(r, req)
		} else {
			err = errors.New("unsupported content type")
		}
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		accessToken, refreshToken, err := h.usersUC.SignInByUsername(r.Context(), req.Username, req.Password)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(presenter.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
		}))
	}
}

// Register – POST /register
// @Summary      Register new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body presenter.UserCreate true "Registration data"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=presenter.UserResponse}
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      422  {object}  responses.ErrorResponse
// @Router       /register [post]
func (h *userHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.UserCreate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		newUser, err := h.usersUC.CreateUser(r.Context(), mapModel(req), req.ConfirmPassword)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		resp := mapModelResponse(newUser)
		if resp == nil {
			render.Render(w, r, responses.CreateErrorResponse(errors.New("internal server error")))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(*resp))
	}
}

// ============================================================================
// USER CRUD (admin)
// ============================================================================

// Create – POST /user
// @Security     BearerAuth
// @Summary      Create user (admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body presenter.UserCreate true "User creation data"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=presenter.UserResponse}
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      403  {object}  responses.ErrorResponse
// @Router       /user [post]
func (h *userHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.UserCreate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.logger.Errorf("decode request body error: %v", err)
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		newUser, err := h.usersUC.CreateUser(r.Context(), mapModel(req), req.ConfirmPassword)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		resp := mapModelResponse(newUser)
		if resp == nil {
			render.Render(w, r, responses.CreateErrorResponse(errors.New("internal server error")))
			return
		}

		// audit
		ip, device, browser := extractClientMeta(r)
		h.logUC.WriteAuditLogWithMeta(r.Context(), newUser.ID, 3, "User Create",
			"Created user "+newUser.Email, ip, device, browser)

		render.Respond(w, r, responses.CreateSuccessResponse(resp))
	}
}

// Get – GET /user/{id}
// @Security     BearerAuth
// @Summary      Get user by ID
// @Tags         users
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=presenter.UserResponse}
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Router       /user/{id} [get]
func (h *userHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		user, err := h.usersUC.Get(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(user)))
	}
}

// GetMulti – GET /user
// @Security     BearerAuth
// @Summary      List users (admin only)
// @Tags         users
// @Produce      json
// @Param        page      query int    false "Page" default(1)
// @Param        per_page  query int    false "Items per page" default(10)
// @Param        email     query string false "Filter email"
// @Param        username  query string false "Filter username"
// @Param        fullname  query string false "Filter fullname"
// @Param        status    query int    false "Filter status"
// @Param        role_id   query int    false "Filter role"
// @Param        verified  query bool   false "Filter verified"
// @Param        sort      query string false "Sort field-ASC|DESC"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=presenter.UserListResponse}
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user [get]
func (h *userHandler) GetMulti() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		page, perPage := parsePageParams(q)
		limit := perPage
		offset := (page - 1) * perPage

		filters := make(map[string]interface{})
		if v := q.Get("email"); v != "" {
			filters["email"] = v
		}
		if v := q.Get("username"); v != "" {
			filters["username"] = v
		}
		if v := q.Get("fullname"); v != "" {
			filters["fullname"] = v
		}
		if v := q.Get("status"); v != "" {
			if n, err := strconvAtoi16(v); err == nil {
				filters["status"] = n
			}
		}
		if v := q.Get("role_id"); v != "" {
			if n, err := strconvAtoi(v); err == nil {
				filters["role_id"] = n
			}
		}
		if v := q.Get("verified"); v != "" {
			filters["verified"] = v == "true"
		}
		if field, order := parseSortParam(q.Get("sort")); field != "" {
			filters["sort_by"] = field
			filters["sort_order"] = strings.ToLower(order)
		} else {
			if v := q.Get("sort_by"); v != "" {
				filters["sort_by"] = v
			}
			if v := q.Get("sort_order"); v != "" {
				filters["sort_order"] = v
			}
		}

		list, total, err := h.usersUC.GetMultiWithTotal(r.Context(), limit, offset, filters)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		payload := presenter.NewUserListResponse(mapModelsResponse(list), page, perPage, total)
		render.Respond(w, r, responses.CreateSuccessResponse(payload))
	}
}

// Delete – DELETE /user/{id}
// @Security     BearerAuth
// @Summary      Delete user (admin only)
// @Tags         users
// @Param        id path string true "User ID"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Router       /user/{id} [delete]
func (h *userHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		user, err := h.usersUC.Delete(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		ip, device, browser := extractClientMeta(r)
		h.logUC.WriteAuditLogWithMeta(r.Context(), id, 5, "User Delete",
			"Deleted user "+user.Email, ip, device, browser)

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(user)))
	}
}

// Update – PUT /user/{id}
// @Security     BearerAuth
// @Summary      Update user (admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        user body presenter.UserUpdate true "Fields to update"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=presenter.UserResponse}
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/{id} [put]
func (h *userHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		req := new(presenter.UserUpdate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		values := make(map[string]interface{})
		if req.Firstname != nil {
			values["firstname"] = *req.Firstname
		}
		if req.Lastname != nil {
			values["lastname"] = *req.Lastname
		}
		if req.Fullname != nil {
			values["fullname"] = *req.Fullname
		}
		if req.MobileNumber != nil {
			values["mobile_number"] = *req.MobileNumber
		}
		if req.PhoneNumber != nil {
			values["phone_number"] = *req.PhoneNumber
		}
		if req.LineID != nil {
			values["lineid"] = *req.LineID
		}
		if req.LocationID != nil {
			values["location_id"] = *req.LocationID
		}
		updatedUser, err := h.usersUC.Update(r.Context(), id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		ip, device, browser := extractClientMeta(r)
		h.logUC.WriteAuditLogWithMeta(r.Context(), id, 4, "User Update",
			"Updated user fields", ip, device, browser)

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
	}
}

// UpdatePassword – PATCH /user/{id}/updatepass
// @Security     BearerAuth
// @Summary      Update user password (admin)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        password body presenter.UserUpdatePassword true "Password data"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/{id}/updatepass [patch]
func (h *userHandler) UpdatePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		req := new(presenter.UserUpdatePassword)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		updatedUser, err := h.usersUC.UpdatePassword(r.Context(), id, req.OldPassword, req.NewPassword, req.ConfirmPassword)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
	}
}

// UpdateRole – PATCH /user/{id}/role
// @Security     BearerAuth
// @Summary      Update user role (admin)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        role body presenter.UserUpdateRole true "Role ID"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=presenter.UserResponse}
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/{id}/role [patch]
func (h *userHandler) UpdateRole() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		req := new(presenter.UserUpdateRole)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		if req.RoleID <= 0 {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(errors.New("invalid role id"))))
			return
		}
		updatedUser, err := h.usersUC.Update(r.Context(), id, map[string]interface{}{"role_id": req.RoleID})
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
	}
}

// ============================================================================
// SELF
// ============================================================================

// Me – GET /user/me
// @Security     BearerAuth
// @Summary      Get current user profile
// @Tags         users
// @Produce      json
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=presenter.UserProfileResponse}
// @Failure      401  {object}  responses.ErrorResponse
// @Router       /user/me [get]
func (h *userHandler) Me() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := middleware.GetUserFromCtx(r.Context())
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		detail, err := h.usersUC.Profile(r.Context(), user.ID)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		resp := mapProfileResponse(detail)
		if resp == nil {
			render.Render(w, r, responses.CreateErrorResponse(errors.New("internal server error")))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(*resp))
	}
}

// UpdateMe – PUT /user/me
// @Security     BearerAuth
// @Summary      Update current user profile
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body presenter.UserUpdate true "Fields to update"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=presenter.UserResponse}
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/me [put]
func (h *userHandler) UpdateMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := middleware.GetUserFromCtx(r.Context())
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		req := new(presenter.UserUpdate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		values := make(map[string]interface{})
		if req.Firstname != nil {
			values["firstname"] = *req.Firstname
		}
		if req.Lastname != nil {
			values["lastname"] = *req.Lastname
		}
		if req.Fullname != nil {
			values["fullname"] = *req.Fullname
		}
		if req.MobileNumber != nil {
			values["mobile_number"] = *req.MobileNumber
		}
		if req.PhoneNumber != nil {
			values["phone_number"] = *req.PhoneNumber
		}
		if req.LineID != nil {
			values["lineid"] = *req.LineID
		}
		if req.LocationID != nil {
			values["location_id"] = *req.LocationID
		}
		updatedUser, err := h.usersUC.Update(r.Context(), user.ID, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
	}
}

// UpdatePasswordMe – PATCH /user/me/updatepass
// @Security     BearerAuth
// @Summary      Change current user password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        password body presenter.UserUpdatePassword true "Password data"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/me/updatepass [patch]
func (h *userHandler) UpdatePasswordMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := middleware.GetUserFromCtx(r.Context())
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		req := new(presenter.UserUpdatePassword)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		updatedUser, err := h.usersUC.UpdatePassword(r.Context(), user.ID, req.OldPassword, req.NewPassword, req.ConfirmPassword)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
	}
}

// LogoutAllAdmin – GET /user/{id}/logoutall
// @Security     BearerAuth
// @Summary      Force logout all sessions
// @Tags         users
// @Param        id path string true "User ID"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/{id}/logoutall [get]
func (h *userHandler) LogoutAllAdmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := h.usersUC.LogoutAll(r.Context(), id); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(struct{}{}))
	}
}

// Profile – GET /user/profile/{id}
// @Security     BearerAuth
// @Summary      Get user profile with role title & permissions
// @Tags         users
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=presenter.UserProfileResponse}
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Router       /user/profile/{id} [get]
func (h *userHandler) Profile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		detail, err := h.usersUC.Profile(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		resp := mapProfileResponse(detail)
		if resp == nil {
			render.Render(w, r, responses.CreateErrorResponse(errors.New("internal server error")))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(*resp))
	}
}

// ListUsers – GET /user/list
// @Security     BearerAuth
// @Summary      List users with keyword/status/sort/pagination
// @Tags         users
// @Produce      json
// @Param        keyword       query string false "Keyword"
// @Param        status        query string false "Status csv"
// @Param        active_status query string false "Active status csv"
// @Param        page          query int    false "Page" default(1)
// @Param        per_page      query int    false "Items per page" default(10)
// @Param        sort          query string false "Sort field-ASC|DESC"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=presenter.UserListResponse}
// @Router       /user/list [get]
func (h *userHandler) ListUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		page, perPage := parsePageParams(q)

		f := users.UserFilter{
			Keyword: q.Get("keyword"),
			Limit:   perPage,
			Offset:  (page - 1) * perPage,
		}
		if v := q.Get("status"); v != "" {
			for _, part := range strings.Split(v, ",") {
				if n, err := strconvAtoi16(strings.TrimSpace(part)); err == nil {
					f.Statuses = append(f.Statuses, n)
				}
			}
		}
		if v := q.Get("active_status"); v != "" {
			for _, part := range strings.Split(v, ",") {
				if n, err := strconvAtoi16(strings.TrimSpace(part)); err == nil {
					f.ActiveStatuses = append(f.ActiveStatuses, n)
				}
			}
		}
		if field, order := parseSortParam(q.Get("sort")); field != "" {
			f.SortField = field
			f.SortOrder = order
		}

		list, total, err := h.usersUC.ListUsers(r.Context(), f)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		payload := presenter.NewUserListResponse(mapModelsResponse(list), page, perPage, total)
		render.Respond(w, r, responses.CreateSuccessResponse(payload))
	}
}

// Statistics – GET /user/statistics
// @Security     BearerAuth
// @Summary      User statistics
// @Tags         users
// @Produce      json
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Router       /user/statistics [get]
func (h *userHandler) Statistics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats, err := h.usersUC.Statistics(r.Context())
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(map[string]interface{}{"payload": stats}))
	}
}

// NotifyList – GET /user/notify/{channel}
// @Security     BearerAuth
// @Summary      List active users by notification channel
// @Tags         users
// @Produce      json
// @Param        channel path string true "Channel" Enums(email,sms,line)
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Router       /user/notify/{channel} [get]
func (h *userHandler) NotifyList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ch := users.NotificationChannel(chi.URLParam(r, "channel"))
		list, err := h.usersUC.ActiveByNotification(r.Context(), ch)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(map[string]interface{}{"payload": mapModelsResponse(list)}))
	}
}

// UpdateActiveStatus – PATCH /user/{id}/activestatus
// @Security     BearerAuth
// @Summary      Toggle user active status
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        body body presenter.UserUpdateActiveStatus true "Active status"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/{id}/activestatus [patch]
func (h *userHandler) UpdateActiveStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		req := new(presenter.UserUpdateActiveStatus)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := h.usersUC.UpdateActiveStatus(r.Context(), id, *req.ActiveStatus); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(struct{}{}))
	}
}

// UploadAvatar – POST /user/me/avatar
// @Security     BearerAuth
// @Summary      Upload avatar (jpg/jpeg/png)
// @Tags         users
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "Avatar file"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/me/avatar [post]
func (h *userHandler) UploadAvatar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := middleware.GetUserFromCtx(r.Context())
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		maxBytes := h.cfg.Upload.MaxAvatarBytes
		if maxBytes <= 0 {
			maxBytes = 2 << 20
		}
		if err := r.ParseMultipartForm(maxBytes + 1024); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		defer file.Close()
		if err := h.usersUC.UpdateAvatar(r.Context(), user.ID, header); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(struct{}{}))
	}
}

// ============================================================================
// ACCESS MENU (handlers)
// ============================================================================

// CreateAccessMenu – POST /user/access-menus
// @Security     BearerAuth
// @Summary      Create access menu mapping
// @Tags         access-menu
// @Accept       json
// @Produce      json
// @Param        body body presenter.AccessMenuCreate true "Access menu data"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=usermodels.SdUserAccessMenu}
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/access-menus [post]
func (h *userHandler) CreateAccessMenu() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.AccessMenuCreate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
	    m := &usermodels.SdUserAccessMenu{
			UserTypeID: &req.UserTypeID,
			RoleID:     &req.RoleID,
			MenuID:     &req.MenuID,
			ParentID:   &req.ParentID,
		}
		out, err := h.accessMenuUC.CreateAccessMenu(r.Context(), m)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(out))
	}
}

// GetAccessMenu – GET /user/access-menus/{id}
// @Security     BearerAuth
// @Summary      Get access menu by ID
// @Tags         access-menu
// @Produce      json
// @Param        id path int true "Access menu ID"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=usermodels.SdUserAccessMenu}
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/access-menus/{id} [get]
func (h *userHandler) GetAccessMenu() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconvAtoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		out, err := h.accessMenuUC.GetAccessMenu(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(out))
	}
}

// UpdateAccessMenu – PUT /user/access-menus/{id}
// @Security     BearerAuth
// @Summary      Update access menu
// @Tags         access-menu
// @Accept       json
// @Produce      json
// @Param        id path int true "Access menu ID"
// @Param        body body presenter.AccessMenuUpdate true "Fields to update"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=usermodels.SdUserAccessMenu}
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/access-menus/{id} [put]
func (h *userHandler) UpdateAccessMenu() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconvAtoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		req := new(presenter.AccessMenuUpdate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		values := map[string]interface{}{}
		if req.UserTypeID != nil {
			values["user_type_id"] = *req.UserTypeID
		}
		if req.RoleID != nil {
			values["role_id"] = *req.RoleID
		}
		if req.MenuID != nil {
			values["menu_id"] = *req.MenuID
		}
		if req.ParentID != nil {
			values["parent_id"] = *req.ParentID
		}
		out, err := h.accessMenuUC.UpdateAccessMenu(r.Context(), id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(out))
	}
}

// DeleteAccessMenu – DELETE /user/access-menus/{id}
// @Security     BearerAuth
// @Summary      Delete access menu
// @Tags         access-menu
// @Param        id path int true "Access menu ID"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/access-menus/{id} [delete]
func (h *userHandler) DeleteAccessMenu() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconvAtoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := h.accessMenuUC.DeleteAccessMenu(r.Context(), id); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(struct{}{}))
	}
}

// ListAccessMenus – GET /user/access-menus
// @Security     BearerAuth
// @Summary      List access menus
// @Tags         access-menu
// @Produce      json
// @Param        page         query int false "Page" default(1)
// @Param        per_page     query int false "Items per page" default(20)
// @Param        user_type_id query int false "Filter user_type_id"
// @Param        role_id      query int false "Filter role_id"
// @Param        menu_id      query int false "Filter menu_id"
// @Param        parent_id    query int false "Filter parent_id"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=presenter.AccessMenuListResponse}
// @Router       /user/access-menus [get]
func (h *userHandler) ListAccessMenus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		page, perPage := parsePageParams(q)
		f := users.AccessMenuFilter{Limit: perPage, Offset: (page - 1) * perPage}
		if v := q.Get("user_type_id"); v != "" {
			if n, err := strconvAtoi(v); err == nil {
				f.UserTypeID = &n
			}
		}
		if v := q.Get("role_id"); v != "" {
			if n, err := strconvAtoi(v); err == nil {
				f.RoleID = &n
			}
		}
		if v := q.Get("menu_id"); v != "" {
			if n, err := strconvAtoi(v); err == nil {
				f.MenuID = &n
			}
		}
		if v := q.Get("parent_id"); v != "" {
			if n, err := strconvAtoi(v); err == nil {
				f.ParentID = &n
			}
		}
		items, total, err := h.accessMenuUC.ListAccessMenus(r.Context(), f)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		list := make([]*presenter.AccessMenuResponse, 0, len(items))
		for _, m := range items {
			list = append(list, &presenter.AccessMenuResponse{
				UserAccessID: m.UserAccessID,
				UserTypeID:   derefInt(m.UserTypeID),
				RoleID:       derefInt(m.RoleID),
				MenuID:       derefInt(m.MenuID),
				ParentID:     derefInt(m.ParentID),
			})
		}
		render.Respond(w, r, responses.CreateSuccessResponse(
			presenter.NewAccessMenuListResponse(list, page, perPage, total),
		))
	}
}

// ============================================================================
// Mapping helpers
// ============================================================================

func mapModel(req *presenter.UserCreate) *models.SdUser {
	var fullname *string
	if req.Fullname != "" {
		fullname = &req.Fullname
	} else if req.Firstname != "" || req.Lastname != "" {
		combined := strings.TrimSpace(req.Firstname + " " + req.Lastname)
		if combined != "" {
			fullname = &combined
		}
	}
	return &models.SdUser{
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		Username:     strings.ToLower(strings.TrimSpace(req.Username)),
		Password:     req.Password,
		RoleID:       req.RoleID,
		Firstname:    stringPtr(req.Firstname),
		Lastname:     stringPtr(req.Lastname),
		Fullname:     fullname,
		MobileNumber: stringPtr(req.MobileNumber),
		PhoneNumber:  stringPtr(req.PhoneNumber),
		LineID:       stringPtr(req.LineID),
		LocationID:   stringPtr(req.LocationID),
		Status:       1,
		IsSuperUser:  false,
		Verified:     false,
	}
}

func mapModelResponse(user *models.SdUser) *presenter.UserResponse {
	if user == nil {
		return nil
	}
	createdAt := helpers.NewLocalTime(user.CreatedDate)
	updatedAt := helpers.NewLocalTime(user.UpdatedDate)
	return &presenter.UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		RoleID:       user.RoleID,
		Firstname:    user.Firstname,
		Lastname:     user.Lastname,
		Fullname:     user.Fullname,
		Nickname:     user.Nickname,
		MobileNumber: user.MobileNumber,
		PhoneNumber:  user.PhoneNumber,
		LineID:       user.LineID,
		LocationID:   user.LocationID,
		Avatar:       user.Avatar,
		Gender:       user.Gender,
		Birthday:     user.Birthday,
		Status:       user.Status,
		IsSuperUser:  user.IsSuperUser,
		Verified:     user.Verified,
		LastSignIn:   helpers.NewLocalTimePtr(&user.Lastsignindate),
		Remark:       user.Remark,
		CreatedAt:    &createdAt,
		UpdatedAt:    &updatedAt,
	}
}

func mapModelsResponse(users []*models.SdUser) []*presenter.UserResponse {
	out := make([]*presenter.UserResponse, len(users))
	for i, u := range users {
		out[i] = mapModelResponse(u)
	}
	return out
}

func mapProfileResponse(d *users.UserProfileDetail) *presenter.UserProfileResponse {
	if d == nil {
		return nil
	}
	base := mapModelResponse(d.User)
	if base == nil {
		return nil
	}
	resp := &presenter.UserProfileResponse{UserResponse: *base, RoleTitle: d.RoleTitle}
	resp.Permissions = make([]presenter.PermissionFlag, 0, len(d.Permissions))
	for _, p := range d.Permissions {
		item := presenter.PermissionFlag{
			Name:     p.Name,
			Insert:   p.Insert,
			Update:   p.Update,
			Delete:   p.Delete,
			Select:   p.Select,
			Log:      p.Log,
			Config:   p.Config,
			Truncate: p.Truncate,
		}
		if p.Detail != nil {
			item.Detail = *p.Detail
		}
		resp.Permissions = append(resp.Permissions, item)
	}
	return resp
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}