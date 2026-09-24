package http

// ไทย: แพ็คเกจ http สำหรับ handler ของ user
// English: http package for user handlers

import (
	"encoding/json"
	// ไทย: นำเข้า encoding/json สำหรับ decode/encode JSON
	// English: import encoding/json for JSON decoding/encoding
	"net/http"
	// ไทย: นำเข้า net/http สำหรับ HTTP types
	// English: import net/http for HTTP types
	"strconv"
	// ไทย: นำเข้า strconv สำหรับแปลง string เป็น int
	// English: import strconv for string to int conversion

	"icmongolang/config"
	// ไทย: นำเข้า config struct
	// English: import config struct
	"icmongolang/internal/middleware"
	// ไทย: middleware สำหรับดึง user จาก context
	// English: middleware to extract user from context
	"icmongolang/internal/models"
	// ไทย: model ของ user
	// English: user model
	"icmongolang/internal/modules/users"
	// ไทย: use case interface ของ user
	// English: user use case interface
	"icmongolang/internal/users/presenter"
	// ไทย: presenter structs สำหรับ request/response
	// English: presenter structs for request/response
	"icmongolang/pkg/httpErrors"
	// ไทย: helper สำหรับ error handling
	// English: helpers for error handling
	"icmongolang/pkg/logger"
	// ไทย: logger interface
	// English: logger interface
	"icmongolang/pkg/responses"
	// ไทย: standard response structs
	// English: standard response structs
	"icmongolang/pkg/utils"
	// ไทย: utilities (เช่น validation)
	// English: utilities (e.g., validation)

	"github.com/go-chi/chi/v5"
	// ไทย: router chi
	// English: chi router
	"github.com/go-chi/render"
	// ไทย: helper สำหรับ HTTP responses
	// English: helper for HTTP responses
	"github.com/google/uuid"
	// ไทย: สำหรับ UUID parsing
	// English: for UUID parsing
)

// ไทย: struct userHandler มี dependencies
// English: userHandler struct holds dependencies
type userHandler struct {
	cfg     *config.Config
	// ไทย: config ของแอพ
	// English: app configuration
	usersUC users.UserUseCaseI
	// ไทย: user use case interface
	// English: user use case interface
	logger  logger.Logger
	// ไทย: logger instance
	// English: logger instance
}

// ไทย: constructor สำหรับ userHandler (factory)
// English: constructor for userHandler (factory)
func CreateUserHandler(uc users.UserUseCaseI, cfg *config.Config, logger logger.Logger) users.Handlers {
	return &userHandler{cfg: cfg, usersUC: uc, logger: logger}
	// ไทย: คืนค่า pointer ไปยัง userHandler
	// English: return pointer to userHandler
}

// Create godoc
// @Summary Create User
// @Description Create new user.
// @Tags users
// @Accept json
// @Produce json
// @Param user body presenter.UserCreate true "Add user"
// @Success 200 {object} responses.SuccessResponse[presenter.UserResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user [post]
// ไทย: Handler สำหรับสร้าง user ใหม่
// English: Handler for creating a new user
func (h *userHandler) Create() func(w http.ResponseWriter, r *http.Request) {
	// ไทย: คืนฟังก์ชันที่รับ http.ResponseWriter และ *http.Request
	// English: return a function that takes http.ResponseWriter and *http.Request
	return func(w http.ResponseWriter, r *http.Request) {
		// ไทย: สร้าง pointer ของ UserCreate
		// English: create a pointer to UserCreate
		user := new(presenter.UserCreate)

		// ไทย: decode JSON body มาใส่ user (ผิด: ควร Decode(user) ไม่ใช่ &user)
		// English: decode JSON body into user (wrong: should be Decode(user) not &user)
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			// ไทย: ถ้า decode ผิด ส่ง error response
			// English: if decode fails, send error response
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		// ไทย: validate struct user
		// English: validate user struct
		err = utils.ValidateStruct(r.Context(), user)
		if err != nil {
			// ไทย: ถ้า validate ผิด ส่ง validation error
			// English: if validation fails, send validation error
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		// ไทย: เรียก use case เพื่อสร้าง user
		// English: call use case to create user
		newUser, err := h.usersUC.CreateUser(
			r.Context(),
			mapModel(user),
			user.ConfirmPassword,
		)
		if err != nil {
			// ไทย: ถ้า use case คืน error ส่งกลับ
			// English: if use case returns error, send it
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		// ไทย: แปลง model เป็น response และ dereference (เสี่ยง nil)
		// English: convert model to response and dereference (risks nil)
		userResponse := *mapModelResponse(newUser)

		// ไทย: ส่ง success response พร้อมข้อมูล user
		// English: send success response with user data
		render.Respond(w, r, responses.CreateSuccessResponse(userResponse))
	}
}

// Get godoc
// @Summary Read user
// @Description Get user by ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Success 200 {object} responses.SuccessResponse[presenter.UserResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user/{id} [get]
// ไทย: Handler สำหรับดึง user ตาม id
// English: Handler for getting user by id
func (h *userHandler) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// ไทย: แปลง id จาก path parameter เป็น UUID
		// English: parse id from path parameter to UUID
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			// ไทย: ถ้า id ไม่ถูกต้อง ส่ง validation error
			// English: if id invalid, send validation error
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		// ไทย: เรียก use case เพื่อหา user
		// English: call use case to find user
		user, err := h.usersUC.Get(r.Context(), id)
		if err != nil {
			// ไทย: ถ้า error ส่งกลับ
			// English: if error, send it
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		// ไทย: ส่ง success response พร้อม user ที่แปลงแล้ว
		// English: send success response with converted user
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(user)))
	}
}

// GetMulti godoc
// @Summary Read Users
// @Description Retrieve users.
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "limit" Format(limit)
// @Param offset query int false "offset" Format(offset)
// @Success 200 {object} responses.SuccessResponse[[]presenter.UserResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user [get]
// ไทย: Handler สำหรับดึง users หลายคน (paginated)
// English: Handler for getting multiple users (paginated)
func (h *userHandler) GetMulti() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		// ไทย: ดึง query parameters
		// English: get query parameters

		limit, _ := strconv.Atoi(q.Get("limit"))
		// ไทย: แปลง limit เป็น int (ถ้า error จะได้ 0)
		// English: convert limit to int (error gives 0)
		offset, _ := strconv.Atoi(q.Get("offset"))
		// ไทย: แปลง offset เป็น int
		// English: convert offset to int

		users, err := h.usersUC.GetMulti(r.Context(), limit, offset)
		// ไทย: เรียก use case เพื่อดึง users
		// English: call use case to get users
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		// ไทย: ส่ง success response พร้อม slice ของ user responses
		// English: send success response with slice of user responses
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelsResponse(users)))
	}
}

// Delete godoc
// @Summary Delete user
// @Description Delete an user by ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Success 200 {object} responses.SuccessResponse[presenter.UserResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user/{id} [delete]
// ไทย: Handler สำหรับลบ user ตาม id
// English: Handler for deleting user by id
func (h *userHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		user, err := h.usersUC.Delete(r.Context(), id)
		// ไทย: เรียก use case เพื่อลบ user
		// English: call use case to delete user
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(user)))
		// ไทย: ส่ง user ที่ถูกลบกลับไป (optional)
		// English: return the deleted user (optional)
	}
}

// Update godoc
// @Summary Update user
// @Description Update an user by ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Param user body presenter.UserUpdate true "Update user"
// @Success 200 {object} responses.SuccessResponse[presenter.UserResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user/{id} [put]
// ไทย: Handler สำหรับอัปเดตข้อมูล user (admin)
// English: Handler for updating user info (admin)
func (h *userHandler) Update() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		user := new(presenter.UserUpdate)
		// ไทย: สร้าง pointer ของ UserUpdate
		// English: create pointer to UserUpdate

		err = json.NewDecoder(r.Body).Decode(&user)
		// ไทย: decode JSON (ผิดอีกแล้ว: ควร Decode(user))
		// English: decode JSON (wrong again: should be Decode(user))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(r.Context(), user)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		// ไทย: สร้าง map สำหรับ field ที่ต้องการอัปเดต
		// English: build map of fields to update
		values := make(map[string]interface{})
		if user.Name != "" {
			values["name"] = user.Name
		}

		updatedUser, err := h.usersUC.Update(r.Context(), id, values)
		// ไทย: เรียก use case อัปเดต
		// English: call use case to update
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
	}
}

// UpdatePassword godoc
// @Summary Update password user
// @Description Update password user by ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Param user body presenter.UserUpdatePassword true "Update user"
// @Success 200 {object} responses.SuccessResponse[presenter.UserResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user/{id}/updatepass [patch]
// ไทย: Handler สำหรับอัปเดตรหัสผ่าน (admin)
// English: Handler for updating password (admin)
func (h *userHandler) UpdatePassword() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		user := new(presenter.UserUpdatePassword)
		err = json.NewDecoder(r.Body).Decode(&user) // ผิดอีก: ควร Decode(user)
		// ไทย: decode request body (ผิดอีก)
		// English: decode request body (wrong again)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(r.Context(), user)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		updatedUser, err := h.usersUC.UpdatePassword(
			r.Context(),
			id,
			user.OldPassword,
			user.NewPassword,
			user.ConfirmPassword,
		)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
	}
}

// Me godoc
// @Summary Read user me
// @Description Get user me.
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} responses.SuccessResponse[presenter.UserResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user/me [get]
// ไทย: Handler สำหรับดึงข้อมูลของตัวเอง (จาก token)
// English: Handler for getting own user info (from token)
func (h *userHandler) Me() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := middleware.GetUserFromCtx(ctx)
		// ไทย: ดึง user object จาก context (ใส่โดย middleware auth)
		// English: extract user object from context (set by auth middleware)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(user)))
	}
}

// UpdateMe godoc
// @Summary Update user me
// @Description Update user me.
// @Tags users
// @Accept json
// @Produce json
// @Param user body presenter.UserUpdate true "Update user"
// @Success 200 {object} responses.SuccessResponse[presenter.UserResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user/me [put]
// ไทย: Handler สำหรับอัปเดตข้อมูลของตัวเอง
// English: Handler for updating own user info
func (h *userHandler) UpdateMe() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		userUpdate := new(presenter.UserUpdate)
		err = json.NewDecoder(r.Body).Decode(&userUpdate) // ผิดอีก
		// ไทย: decode request body (ผิดอีก)
		// English: decode request body (wrong again)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(r.Context(), userUpdate)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		values := make(map[string]interface{})
		if userUpdate.Name != "" {
			values["name"] = userUpdate.Name
		}

		updatedUser, err := h.usersUC.Update(r.Context(), user.Id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
	}
}

// UpdatePasswordMe godoc
// @Summary Update password user me
// @Description Update password user me.
// @Tags users
// @Accept json
// @Produce json
// @Param user body presenter.UserUpdatePassword true "Update user"
// @Success 200 {object} responses.SuccessResponse[presenter.UserResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user/me/updatepass [patch]
// ไทย: Handler สำหรับอัปเดตรหัสผ่านของตัวเอง
// English: Handler for updating own password
func (h *userHandler) UpdatePasswordMe() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		userUpdate := new(presenter.UserUpdatePassword)
		err = json.NewDecoder(r.Body).Decode(&userUpdate) // ผิดอีก
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(r.Context(), userUpdate)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		updatedUser, err := h.usersUC.UpdatePassword(
			r.Context(),
			user.Id,
			userUpdate.OldPassword,
			userUpdate.NewPassword,
			userUpdate.ConfirmPassword,
		)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
	}
}

// LogoutAllAdmin godoc
// @Summary Logout all of user
// @Description Logout all session of user with id.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Success 200
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /user/{id}/logoutall [get]
// ไทย: Handler สำหรับ logout ทุก session ของ user (admin)
// English: Handler for logging out all sessions of a user (admin)
func (h *userHandler) LogoutAllAdmin() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		err = h.usersUC.LogoutAll(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}
		// ไทย: ไม่มี response body (success เท่านั้น)
		// English: no response body (only success status)
	}
}

// ไทย: แปลง UserCreate presenter เป็น models.User
// English: convert UserCreate presenter to models.User
func mapModel(exp *presenter.UserCreate) *models.User {
	return &models.User{
		Name:        exp.Name,
		Email:       exp.Email,
		Password:    exp.Password,
		IsActive:    true,
		IsSuperUser: false,
	}
}

// ไทย: แปลง models.User เป็น UserResponse presenter
// English: convert models.User to UserResponse presenter
func mapModelResponse(exp *models.User) *presenter.UserResponse {
	return &presenter.UserResponse{
		Id:          exp.Id,
		Name:        exp.Name,
		Email:       exp.Email,
		CreatedAt:   exp.CreatedAt,
		UpdatedAt:   exp.UpdatedAt,
		IsActive:    exp.IsActive,
		IsSuperUser: exp.IsSuperUser,
		Verified:    exp.Verified,
	}
}

// ไทย: แปลง slice ของ models.User เป็น slice ของ UserResponse
// English: convert slice of models.User to slice of UserResponse
func mapModelsResponse(exp []*models.User) []*presenter.UserResponse {
	out := make([]*presenter.UserResponse, len(exp))
	for i, user := range exp {
		out[i] = mapModelResponse(user)
	}
	return out
}