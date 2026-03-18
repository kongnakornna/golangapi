package handler

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"

	authdto "github.com/vadxq/go-rest-starter/internal/core/auth/dto"
	authservice "github.com/vadxq/go-rest-starter/internal/core/auth/service"
	httpx "github.com/vadxq/go-rest-starter/internal/transport/httpx"
	apperrors "github.com/vadxq/go-rest-starter/pkg/errors"
	"github.com/vadxq/go-rest-starter/pkg/logger"
)

// AuthHandler จัดการคำขอ HTTP ที่เกี่ยวข้องกับการรับรองตัวตน
type AuthHandler struct {
    authService authservice.AuthService
    logger      logger.Logger
    validator   *validator.Validate
}

// NewAuthHandler สร้างอินสแตนซ์ AuthHandler ใหม่
func NewAuthHandler(as authservice.AuthService, log logger.Logger, v *validator.Validate) *AuthHandler {
    return &AuthHandler{
        authService: as,
        logger:      log,
        validator:   v,
    }
}

// Login จัดการคำขอเข้าสู่ระบบของผู้ใช้
// @Summary ผู้ใช้เข้าสู่ระบบ
// @Description เข้าสู่ระบบโดยใช้อีเมลและรหัสผ่าน เพื่อรับโทเค็นการเข้าถึง
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "ข้อมูลคำขอเข้าสู่ระบบ"
// @Success 200 {object} httpx.Response{data=authdto.LoginResponse}
// @Failure 400,401,500 {object} httpx.Response{data=httpx.ErrorDetail}
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req authdto.LoginRequest

    if err := httpx.BindJSON(r, &req, func(v interface{}) error {
        return h.validator.Struct(v)
    }); err != nil {
        httpx.Error(w, r, err)
        return
    }

    response, err := h.authService.Login(r.Context(), req)
    if err != nil {
        httpx.Error(w, r, err)
        return
    }

    httpx.JSON(w, r, http.StatusOK, response)
}

// RefreshToken จัดการคำขอรีเฟรชโทเค็น
// @Summary รีเฟรชโทเค็นการเข้าถึง
// @Description ใช้โทเค็นรีเฟรชเพื่อขอรับโทเค็นการเข้าถึงใหม่
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.RefreshTokenRequest true "ข้อมูลคำขอรีเฟรชโทเค็น"
// @Success 200 {object} httpx.Response{data=authdto.TokenResponse}
// @Failure 400,401,500 {object} httpx.Response{data=httpx.ErrorDetail}
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
    var req authdto.RefreshTokenRequest

    if err := httpx.BindJSON(r, &req, func(v interface{}) error {
        return h.validator.Struct(v)
    }); err != nil {
        httpx.Error(w, r, err)
        return
    }

    response, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
    if err != nil {
        httpx.Error(w, r, err)
        return
    }

    httpx.JSON(w, r, http.StatusOK, response)
}

// Logout จัดการคำขอออกจากระบบของผู้ใช้
// @Summary ผู้ใช้ออกจากระบบ
// @Description ทำให้โทเค็นการเข้าถึงปัจจุบันของผู้ใช้หมดอายุ
// @Tags auth
// @Accept json
// @Produce json
// @Success 204 {object} nil
// @Failure 401,500 {object} httpx.Response{data=httpx.ErrorDetail}
// @Router /api/v1/auth/logout [post]
// @Security BearerAuth
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    // ดึงโทเค็นการเข้าถึงจากส่วนหัว Authorization
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        httpx.Error(w, r, apperrors.UnauthorizedError("ไม่มีโทเค็นอนุญาต", nil))
        return
    }

    // แยกคำนำหน้า Bearer และโทเค็น
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
        httpx.Error(w, r, apperrors.UnauthorizedError("รูปแบบการอนุญาตไม่ถูกต้อง", nil))
        return
    }

    accessToken := parts[1]

    // เรียกใช้บริการเพื่อดำเนินการออกจากระบบ
    err := h.authService.Logout(r.Context(), accessToken)
    if err != nil {
        httpx.Error(w, r, err)
        return
    }

    // สำเร็จการออกจากระบบ ส่งคืนรหัสสถานะ 204
    httpx.JSON(w, r, http.StatusNoContent, nil)
}

/*
**คำอธิบายเพิ่มเติม**
ไฟล์ `auth_handler.go` นี้เป็นส่วนหนึ่งของเลเยอร์การขนส่ง (Transport Layer) ในโครงสร้าง Clean Architecture โดยมีหน้าที่หลักคือ:
1. รับคำขอ HTTP ที่เกี่ยวข้องกับการรับรองตัวตน
2. ตรวจสอบความถูกต้องเบื้องต้นของข้อมูลที่ส่งมา (เช่น ใช้ validator)
3. เรียกใช้ `AuthService` ที่เกี่ยวข้องเพื่อดำเนินการตามตรรกะทางธุรกิจ
4. ส่งผลลัพธ์กลับไปยังผู้เรียกในรูปแบบ JSON ผ่านฟังก์ชัน `httpx.JSON`
การแปลคอมเมนต์และข้อความเป็นภาษาไทยจะช่วยให้ทีมพัฒนาหรือผู้ที่เกี่ยวข้องที่ใช้ภาษาไทยเข้าใจการทำงานของโค้ดส่วนนี้ได้ง่ายขึ้น โดยยังคงโครงสร้างและชื่อฟังก์ชันเดิมไว้
*/