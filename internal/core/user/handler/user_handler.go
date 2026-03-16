package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	userdto "github.com/vadxq/go-rest-starter/internal/core/user/dto"
	userservice "github.com/vadxq/go-rest-starter/internal/core/user/service"
	httpx "github.com/vadxq/go-rest-starter/internal/transport/httpx"
	apperrors "github.com/vadxq/go-rest-starter/pkg/errors"
	"github.com/vadxq/go-rest-starter/pkg/logger"
)

// UserHandler จัดการคำขอ HTTP ที่เกี่ยวข้องกับผู้ใช้
type UserHandler struct {
	userService userservice.UserService
	logger      logger.Logger
	validator   *validator.Validate
}

// NewUserHandler สร้างอินสแตนซ์ใหม่ของ UserHandler
func NewUserHandler(us userservice.UserService, log logger.Logger, v *validator.Validate) *UserHandler {
	return &UserHandler{
		userService: us,
		logger:      log,
		validator:   v,
	}
}

// GetUser รับข้อมูลผู้ใช้
// @Summary รับข้อมูลผู้ใช้
// @Description รับข้อมูลผู้ใช้ตาม ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "รหัสผู้ใช้"
// @Success 200 {object} httpx.Response{data=userdto.UserResponse}
// @Failure 400,404,500 {object} httpx.Response{data=httpx.ErrorDetail}
// @Router /api/v1/users/{id} [get]
// @Security BearerAuth
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		httpx.Error(w, r, apperrors.BadRequestError("ขาดพารามิเตอร์ ID", nil))
		return
	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}

	// แปลงเป็น DTO
	response := userdto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	httpx.JSON(w, r, http.StatusOK, response)
}

// CreateUser สร้างผู้ใช้
// @Summary สร้างผู้ใช้
// @Description สร้างผู้ใช้ใหม่
// @Tags users
// @Accept json
// @Produce json
// @Param body body dto.CreateUserInput true "ข้อมูลคำขอสร้างผู้ใช้"
// @Success 201 {object} httpx.Response{data=userdto.UserResponse}
// @Failure 400,500 {object} httpx.Response{data=httpx.ErrorDetail}
// @Router /api/v1/users [post]
// @Security BearerAuth
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input userdto.CreateUserInput

	if err := httpx.BindJSON(r, &input, func(v interface{}) error {
		return h.validator.Struct(v)
	}); err != nil {
		httpx.Error(w, r, err)
		return
	}

	user, err := h.userService.CreateUser(r.Context(), input)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}

	// แปลงเป็น DTO
	response := userdto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	httpx.JSON(w, r, http.StatusCreated, response)
}

// UpdateUser อัปเดตผู้ใช้
// @Summary อัปเดตผู้ใช้
// @Description อัปเดตข้อมูลผู้ใช้ตาม ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "รหัสผู้ใช้"
// @Param body body dto.UpdateUserInput true "ข้อมูลคำขออัปเดตผู้ใช้"
// @Success 200 {object} httpx.Response{data=userdto.UserResponse}
// @Failure 400,404,500 {object} httpx.Response{data=httpx.ErrorDetail}
// @Router /api/v1/users/{id} [put]
// @Security BearerAuth
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		httpx.Error(w, r, apperrors.BadRequestError("ขาดพารามิเตอร์ ID", nil))
		return
	}

	var input userdto.UpdateUserInput
	if err := httpx.BindJSON(r, &input, nil); err != nil {
		httpx.Error(w, r, err)
		return
	}

	user, err := h.userService.UpdateUser(r.Context(), userID, input)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}

	// แปลงเป็น DTO
	response := userdto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	httpx.JSON(w, r, http.StatusOK, response)
}

// DeleteUser ลบผู้ใช้
// @Summary ลบผู้ใช้
// @Description ลบผู้ใช้ตาม ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "รหัสผู้ใช้"
// @Success 204 {object} nil
// @Failure 400,404,500 {object} httpx.Response{data=httpx.ErrorDetail}
// @Router /api/v1/users/{id} [delete]
// @Security BearerAuth
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		httpx.Error(w, r, apperrors.BadRequestError("ขาดพารามิเตอร์ ID", nil))
		return
	}

	err := h.userService.DeleteUser(r.Context(), userID)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}

	httpx.JSON(w, r, http.StatusNoContent, nil)
}

// ListUsers รับรายการผู้ใช้
// @Summary รับรายการผู้ใช้
// @Description รับรายการผู้ใช้แบบแบ่งหน้า
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "หมายเลขหน้า ค่าเริ่มต้นคือ 1" default(1)
// @Param page_size query int false "ขนาดต่อหน้า ค่าเริ่มต้นคือ 10" default(10)
// @Success 200 {object} httpx.Response{data=userdto.ListResponse{data=[]userdto.UserResponse}}
// @Failure 500 {object} httpx.Response{data=httpx.ErrorDetail}
// @Router /api/v1/users [get]
// @Security BearerAuth
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// แยกวิเคราะห์พารามิเตอร์แบ่งหน้า
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page := 1
	pageSize := 10

	if pageStr != "" {
		pageVal, err := strconv.Atoi(pageStr)
		if err == nil && pageVal > 0 {
			page = pageVal
		}
	}

	if pageSizeStr != "" {
		pageSizeVal, err := strconv.Atoi(pageSizeStr)
		if err == nil && pageSizeVal > 0 {
			pageSize = pageSizeVal
		}
	}

	users, total, err := h.userService.ListUsers(r.Context(), page, pageSize)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}

	// แปลงเป็น DTO
	userResponses := make([]userdto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = userdto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	response := userdto.ListResponse{
		Data:  userResponses,
		Total: total,
		Page:  page,
		Size:  pageSize,
	}

	httpx.JSON(w, r, http.StatusOK, response)
}