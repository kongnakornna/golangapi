package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"icmongolang/config"
	"icmongolang/internal/middleware"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/users"
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

type userHandler struct {
	cfg     *config.Config
	usersUC users.UserUseCaseI
	logger  logger.Logger
}

func CreateUserHandler(uc users.UserUseCaseI, cfg *config.Config, logger logger.Logger) users.Handlers {
	return &userHandler{cfg: cfg, usersUC: uc, logger: logger}
}

// SignInEmail – POST /api/auth/signin (email + password)
// @Summary      User login with email and password
// @Description  Authenticates a user using email and password. Returns access and refresh tokens.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body presenter.UserSignIn true "User credentials"
// @Success      200  {object}  map[string]interface{}  "Tokens returned successfully"
// @Failure      400  {object}  responses.ErrorResponse  "Bad request / validation error"
// @Failure      401  {object}  responses.ErrorResponse  "Unauthorized (wrong email/password)"
// @Failure      404  {object}  responses.ErrorResponse  "User not found"
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

		render.Respond(w, r, responses.CreateSuccessResponse(presenter.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
		}))
	}
}

// SignInUsername – POST /api/auth/login (username + password)
// @Summary      User login with username and password
// @Description  Authenticates a user using username and password. Returns access and refresh tokens.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body presenter.SignIn true "User credentials (username + password)"
// @Success      200  {object}  map[string]interface{}  "Tokens returned successfully"
// @Failure      400  {object}  responses.ErrorResponse  "Bad request / validation error"
// @Failure      401  {object}  responses.ErrorResponse  "Unauthorized (wrong username/password)"
// @Failure      404  {object}  responses.ErrorResponse  "User not found"
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

// Register – POST /register (public, no auth)
// @Summary      สมัครสมาชิกผู้ใช้ใหม่
// @Description  ลงทะเบียนผู้ใช้ทั่วไป (ไม่ต้องใช้ token)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body presenter.UserCreate true "ข้อมูลสมัคร"
// @Success      200  {object}  map[string]interface{}
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
		userResp := mapModelResponse(newUser)
		if userResp == nil {
			render.Render(w, r, responses.CreateErrorResponse(errors.New("internal server error")))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(*userResp))
	}
}

// Create – POST /user (admin only)
// @Security      BearerAuth
// @Summary       Create user (admin only)
// @Tags          users
// @Accept        json
// @Produce       json
// @Param         user body presenter.UserCreate true "User creation data"
// @Success       200  {object}  map[string]interface{}
// @Failure       400  {object}  responses.ErrorResponse
// @Failure       403  {object}  responses.ErrorResponse
// @Router        /user [post]
func (h *userHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.UserCreate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.logger.Errorf("decode request body error: %v", err)
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		h.logger.Infof("request body decoded: email=%s, role_id=%d", req.Email, req.RoleID)

		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			h.logger.Warnf("validation failed: %v", err)
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		h.logger.Info("validation passed")

		userModel := mapModel(req)
		if userModel == nil {
			h.logger.Error("mapModel returned nil")
			render.Render(w, r, responses.CreateErrorResponse(errors.New("internal mapping error")))
			return
		}
		h.logger.Infof("mapped to model: email=%s, role_id=%d, status=%d", userModel.Email, userModel.RoleID, userModel.Status)

		newUser, err := h.usersUC.CreateUser(r.Context(), userModel, req.ConfirmPassword)
		if err != nil {
			h.logger.Errorf("CreateUser error: %v", err)
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		h.logger.Infof("user created successfully with id=%s", newUser.ID)

		responseData := mapModelResponse(newUser)
		if responseData == nil {
			h.logger.Error("mapModelResponse returned nil")
			render.Render(w, r, responses.CreateErrorResponse(errors.New("internal response error")))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(responseData))
	}
}

// Get – GET /user/{id}
// @Security      BearerAuth
// @Summary       Get user by ID
// @Tags          users
// @Produce       json
// @Param         id path string true "User ID"
// @Success       200  {object}  map[string]interface{}
// @Failure       400  {object}  responses.ErrorResponse
// @Failure       404  {object}  responses.ErrorResponse
// @Router        /user/{id} [get]
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

// GetMulti – GET /user?limit=10&offset=0
// @Security      BearerAuth
// @Summary       List users (admin only)
// @Tags          users
// @Produce       json
// @Param         limit query int false "Page size" default(20)
// @Param         offset query int false "Offset" default(0)
// @Param         email query string false "Filter by email"
// @Param         username query string false "Filter by username"
// @Param         status query int false "Filter by status"
// @Param         role_id query int false "Filter by role ID"
// @Success       200  {object}  map[string]interface{}
// @Failure       400  {object}  responses.ErrorResponse
// @Router        /user [get]
func (h *userHandler) GetMulti() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		limit, _ := strconv.Atoi(q.Get("limit"))
		if limit <= 0 {
			limit = 20
		}
		offset, _ := strconv.Atoi(q.Get("offset"))
		if offset < 0 {
			offset = 0
		}

		filters := make(map[string]interface{})
		if email := q.Get("email"); email != "" {
			filters["email"] = email
		}
		if username := q.Get("username"); username != "" {
			filters["username"] = username
		}
		if fullname := q.Get("fullname"); fullname != "" {
			filters["fullname"] = fullname
		}
		if statusStr := q.Get("status"); statusStr != "" {
			if status, err := strconv.ParseInt(statusStr, 10, 16); err == nil {
				filters["status"] = int16(status)
			}
		}
		if roleStr := q.Get("role_id"); roleStr != "" {
			if roleID, err := strconv.Atoi(roleStr); err == nil {
				filters["role_id"] = roleID
			}
		}
		if verifiedStr := q.Get("verified"); verifiedStr != "" {
			filters["verified"] = verifiedStr == "true"
		}
		if sortBy := q.Get("sort_by"); sortBy != "" {
			filters["sort_by"] = sortBy
		}
		if sortOrder := q.Get("sort_order"); sortOrder != "" {
			filters["sort_order"] = sortOrder
		}

		users, total, err := h.usersUC.GetMultiWithTotal(r.Context(), limit, offset, filters)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		response := map[string]interface{}{
			"payload": mapModelsResponse(users),
			"pagination": map[string]interface{}{
				"limit":  limit,
				"offset": offset,
				"total":  total,
				"next_offset": func() int {
					if offset+limit < int(total) {
						return offset + limit
					}
					return -1
				}(),
			},
		}
		render.Respond(w, r, responses.CreateSuccessResponse(response))
	}
}

// Delete – DELETE /user/{id} (admin only)
// @Security      BearerAuth
// @Summary       Delete user (admin only)
// @Tags          users
// @Param         id path string true "User ID"
// @Success       200  {object}  map[string]interface{}
// @Failure       400  {object}  responses.ErrorResponse
// @Failure       404  {object}  responses.ErrorResponse
// @Router        /user/{id} [delete]
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
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(user)))
	}
}

// Update – PUT /user/{id} (admin only)
// @Security      BearerAuth
// @Summary       Update user (admin only)
// @Tags          users
// @Accept        json
// @Produce       json
// @Param         id path string true "User ID"
// @Param         user body presenter.UserUpdate true "Fields to update"
// @Success       200  {object}  map[string]interface{}
// @Failure       400  {object}  responses.ErrorResponse
// @Router        /user/{id} [put]
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
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
	}
}

// UpdatePassword – PATCH /user/{id}/updatepass (admin only)
// @Security      BearerAuth
// @Summary       Update user password (admin only)
// @Tags          users
// @Accept        json
// @Produce       json
// @Param         id path string true "User ID"
// @Param         password body presenter.UserUpdatePassword true "Password data"
// @Success       200  {object}  map[string]interface{}
// @Failure       400  {object}  responses.ErrorResponse
// @Router        /user/{id}/updatepass [patch]
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

// UpdateRole – PATCH /user/{id}/role (admin only)
// @Security      BearerAuth
// @Summary       Update user role (admin only)
// @Tags          users
// @Accept        json
// @Produce       json
// @Param         id path string true "User ID"
// @Param         role body presenter.UserUpdateRole true "Role ID"
// @Success       200  {object}  map[string]interface{}
// @Failure       400  {object}  responses.ErrorResponse
// @Router        /user/{id}/role [patch]
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

// Me – GET /user/me
// @Security      BearerAuth
// @Summary       Get current user profile (with role title and permissions)
// @Tags          users
// @Produce       json
// @Success       200  {object}  map[string]interface{}
// @Failure       401  {object}  responses.ErrorResponse
// @Router        /user/me [get]
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
// @Security      BearerAuth
// @Summary       Update current user profile
// @Tags          users
// @Accept        json
// @Produce       json
// @Param         user body presenter.UserUpdate true "Fields to update"
// @Success       200  {object}  map[string]interface{}
// @Failure       400  {object}  responses.ErrorResponse
// @Router        /user/me [put]
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
// @Security      BearerAuth
// @Summary       Change current user password
// @Tags          users
// @Accept        json
// @Produce       json
// @Param         password body presenter.UserUpdatePassword true "Password data"
// @Success       200  {object}  map[string]interface{}
// @Failure       400  {object}  responses.ErrorResponse
// @Router        /user/me/updatepass [patch]
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

// LogoutAllAdmin – GET /user/{id}/logoutall (admin only)
// @Security      BearerAuth
// @Summary       Force logout all sessions of a user (admin only)
// @Tags          users
// @Param         id path string true "User ID"
// @Success       200  {object}  map[string]interface{}
// @Failure       400  {object}  responses.ErrorResponse
// @Router        /user/{id}/logoutall [get]
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
// @Security      BearerAuth
// @Summary       Get user profile with role title and permissions
// @Tags          users
// @Produce       json
// @Param         id path string true "User ID"
// @Success       200  {object}  map[string]interface{}
// @Failure       400  {object}  responses.ErrorResponse
// @Failure       404  {object}  responses.ErrorResponse
// @Router        /user/profile/{id} [get]
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

// ListUsers – GET /user/list (admin only)
// @Security      BearerAuth
// @Summary       List users with keyword/status/sort/pagination (admin only)
// @Tags          users
// @Produce       json
// @Param         keyword query string false "Keyword (username or email)"
// @Param         status query string false "Status csv (e.g. 0,1)"
// @Param         active_status query string false "Active status csv"
// @Param         sort query string false "Sort field-ASC|DESC (e.g. username-ASC)"
// @Param         page query int false "Page number" default(1)
// @Param         page_size query int false "Page size" default(20)
// @Success       200  {object}  map[string]interface{}
// @Failure       400  {object}  responses.ErrorResponse
// @Router        /user/list [get]
func (h *userHandler) ListUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		f := users.UserFilter{Keyword: q.Get("keyword")}

		if v := q.Get("status"); v != "" {
			for _, part := range strings.Split(v, ",") {
				if n, err := strconv.ParseInt(strings.TrimSpace(part), 10, 16); err == nil {
					f.Statuses = append(f.Statuses, int16(n))
				}
			}
		}
		if v := q.Get("active_status"); v != "" {
			for _, part := range strings.Split(v, ",") {
				if n, err := strconv.ParseInt(strings.TrimSpace(part), 10, 16); err == nil {
					f.ActiveStatuses = append(f.ActiveStatuses, int16(n))
				}
			}
		}
		if v := q.Get("sort"); v != "" {
			if idx := strings.LastIndex(v, "-"); idx > 0 {
				f.SortField = v[:idx]
				f.SortOrder = v[idx+1:]
			} else {
				f.SortField = v
			}
		}

		page, _ := strconv.Atoi(q.Get("page"))
		if page < 1 {
			page = 1
		}
		pageSize, _ := strconv.Atoi(q.Get("page_size"))
		if pageSize < 1 {
			pageSize = 20
		}
		if pageSize > 100 {
			pageSize = 100
		}
		f.Limit = pageSize
		f.Offset = (page - 1) * pageSize

		list, total, err := h.usersUC.ListUsers(r.Context(), f)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		response := map[string]interface{}{
			"payload": mapModelsResponse(list),
			"pagination": map[string]interface{}{
				"page":      page,
				"page_size": pageSize,
				"total":     total,
			},
		}
		render.Respond(w, r, responses.CreateSuccessResponse(response))
	}
}

// Statistics – GET /user/statistics (admin only)
// @Security      BearerAuth
// @Summary       User statistics aggregated by status and role (admin only)
// @Tags          users
// @Produce       json
// @Success       200  {object}  map[string]interface{}
// @Router        /user/statistics [get]
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

// NotifyList – GET /user/notify/{channel} (admin only)
// @Security      BearerAuth
// @Summary       List active users by notification channel (admin only)
// @Tags          users
// @Produce       json
// @Param         channel path string true "Notification channel (email|sms|line)"
// @Success       200  {object}  map[string]interface{}
// @Failure       400  {object}  responses.ErrorResponse
// @Router        /user/notify/{channel} [get]
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

// UpdateActiveStatus – PATCH /user/{id}/activestatus (admin only)
// @Security      BearerAuth
// @Summary       Toggle user active status (admin only)
// @Tags          users
// @Accept        json
// @Produce       json
// @Param         id path string true "User ID"
// @Param         body body presenter.UserUpdateActiveStatus true "Active status"
// @Success       200  {object}  map[string]interface{}
// @Failure       400  {object}  responses.ErrorResponse
// @Failure       404  {object}  responses.ErrorResponse
// @Router        /user/{id}/activestatus [patch]
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
// @Security      BearerAuth
// @Summary       Upload current user avatar (jpg/jpeg/png)
// @Tags          users
// @Accept        multipart/form-data
// @Produce       json
// @Param         file formData file true "Avatar file"
// @Success       200  {object}  map[string]interface{}
// @Failure       400  {object}  responses.ErrorResponse
// @Router        /user/me/avatar [post]
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

// ------------------------------
// Mapping functions
// ------------------------------

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
		CreatedAt:    helpers.NewLocalTime(user.CreatedDate),
		UpdatedAt:    helpers.NewLocalTime(user.UpdatedDate),
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
	base := mapModelResponse(d.User)
	if base == nil {
		return nil
	}
	resp := &presenter.UserProfileResponse{UserResponse: *base, RoleTitle: d.RoleTitle}
	resp.Permissions = make([]presenter.PermissionFlag, 0, len(d.Permissions))
	for _, p := range d.Permissions {
		item := presenter.PermissionFlag{
			Name: p.Name, Insert: p.Insert, Update: p.Update,
			Delete: p.Delete, Select: p.Select, Log: p.Log,
			Config: p.Config, Truncate: p.Truncate,
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
