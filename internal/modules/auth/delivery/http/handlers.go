package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"icmongolang/config"
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/auth"
	"icmongolang/internal/modules/users"
	"icmongolang/internal/modules/users/presenter"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/jwt"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/render"
)

type authHandler struct {
	usersUC             users.UserUseCaseI
	cfg                 *config.Config
	logger              logger.Logger
	rateLimitMiddleware *middleware.RateLimitMiddleware
}

// CreateAuthHandler creates auth handler instance
func CreateAuthHandler(uc users.UserUseCaseI, cfg *config.Config, logger logger.Logger) auth.Handlers {
	// สร้าง rate limit middleware (กำหนดค่าเองได้)
	rateLimitConfig := middleware.RateLimitConfig{
		// RequestsPerSecond: จำนวนคำขอสูงสุดที่อนุญาตต่อวินาที
		// Maximum number of requests allowed per second
		RequestsPerSecond: 50, // 50 requests per second | 50 คำขอต่อวินาที

		// Burst: จำนวนคำขอสูงสุดที่อนุญาตให้ส่งพร้อมกัน (กระชุ)
		// Maximum number of requests allowed to burst at once
		Burst: 100, // Burst up to 100 | อนุญาตให้ส่งสูงสุด 100 คำขอในครั้งเดียว

		// CleanupInterval: ระยะเวลาในการทำความสะอาดข้อมูล IP ที่ไม่ใช้งาน
		// Interval for cleaning up inactive IP data
		CleanupInterval: 15 * time.Minute, // Clean up every 15 minutes | ทำความสะอาดทุก 15 นาที
	}
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(rateLimitConfig)

	return &authHandler{
		cfg:                 cfg,
		usersUC:             uc,
		logger:              logger,
		rateLimitMiddleware: rateLimitMiddleware, // เพิ่ม
	}
}

// setAuthCookies เก็บ access/refresh token ลง cookie เพื่อให้ browser ส่งไปกับ
// request ถัดไปโดยอัตโนมัติ (ไม่ต้องแนบ Authorization header ด้วยมือ)
func (h *authHandler) setAuthCookies(w http.ResponseWriter, accessToken, refreshToken string, accessExpiresIn, refreshExpiresIn int) {
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.AccessTokenCookieName,
		Value:    accessToken,
		Path:     "/",
		MaxAge:   accessExpiresIn,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.RefreshTokenCookieName,
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   refreshExpiresIn,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// clearAuthCookies ลบ cookie ที่เกี่ยวข้องกับ auth (ใช้ตอน logout)
func (h *authHandler) clearAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.AccessTokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.RefreshTokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// SignInEmail – POST /api/auth/signin (email + password)
// @Summary      User login with email and password
// @Description  Authenticates a user using email and password. Accepts either a JSON body {"email","password"} or multipart/form-data (email, password). Returns access and refresh tokens (also set as HttpOnly cookies). Rate-limited.
// @Tags         auth
// @Accept       json,multipart/form-data
// @Produce      json
// @Param        email    formData string true "User email"
// @Param        password formData string true "User password"
// @Param        Authorization header string false "Bearer token (not required for sign-in)"
// @Success      200  {object}  map[string]interface{}  "Tokens returned successfully (access_token, refresh_token, token_type, expires_in)"
// @Failure      400  {object}  responses.ErrorResponse  "Bad request / validation error"
// @Failure      401  {object}  responses.ErrorResponse  "Unauthorized (wrong email/password)"
// @Failure      404  {object}  responses.ErrorResponse  "User not found"
// @Failure      422  {object}  responses.ErrorResponse  "Unprocessable entity / validation error"
// @Failure      429  {object}  responses.ErrorResponse  "Too many requests (rate-limited)"
// @Router       /auth/signin [post]
func (h *authHandler) SignIn() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// ใช้ rate limit middleware ตรวจสอบก่อน
		// สร้าง handler ชั่วคราวเพื่อตรวจสอบ rate limit
		checkHandler := h.rateLimitMiddleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ถ้าผ่าน rate limit ให้ทำงานปกติ
		}))

		// สร้าง response writer จับ error
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// ตรวจสอบ rate limit
		checkHandler.ServeHTTP(recorder, r)

		// ถ้า status เป็น 429 แสดงว่าโดน rate limit
		if recorder.statusCode == http.StatusTooManyRequests {
			return
		}

		user := new(presenter.UserSignIn)

		contentType := r.Header.Get("Content-Type")

		var email, password string

		if strings.Contains(contentType, "application/json") {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				h.logger.Error("Failed to read request body")
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
			defer r.Body.Close()

			var jsonBody map[string]interface{}
			if err := json.Unmarshal(body, &jsonBody); err != nil {
				h.logger.Error("Failed to parse JSON body")
				render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
				return
			}

			email, _ = jsonBody["email"].(string)
			password, _ = jsonBody["password"].(string)
		} else {
			if err := r.ParseForm(); err != nil {
				h.logger.Error("Failed to parse form")
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
			email = r.FormValue("email")
			password = r.FormValue("password")
		}

		if email == "" || password == "" {
			h.logger.Warn("SignIn: missing email or password")
		}

		user.Email = email
		user.Password = password

		h.logger.Info("SignIn attempt")

		err := utils.ValidateStruct(r.Context(), user)
		if err != nil {
			h.logger.Error("SignIn validation failed")
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		accessToken, refreshToken, err := h.usersUC.SignIn(
			r.Context(),
			user.Email,
			user.Password,
		)
		if err != nil {
			h.logger.Info("SignIn failed")
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		expiresIn := h.cfg.Jwt.AccessTokenExpireDuration * 60
		refreshExpiresIn := h.cfg.Jwt.RefreshTokenExpireDuration * 60

		h.setAuthCookies(w, accessToken, refreshToken, int(expiresIn), int(refreshExpiresIn))

		h.logger.Info("SignIn successful")

		render.Respond(w, r, map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "bearer",
			"expires_in":    expiresIn,
		})
	}
}

// SignInUsername – POST /api/auth/login (username + password)
// @Summary      User login with username and password
// @Description  Authenticates a user using username and password. Accepts either a JSON body {"username","password"} or multipart/form-data (username, password). Returns access and refresh tokens (also set as HttpOnly cookies). Rate-limited.
// @Tags         auth
// @Accept       json,multipart/form-data
// @Produce      json
// @Param        username formData string true "username"
// @Param        password formData string true "password"
// @Param        Authorization header string false "Bearer token (not required for login)"
// @Success      200  {object}  map[string]interface{}  "Tokens returned successfully (access_token, refresh_token, token_type, expires_in)"
// @Failure      400  {object}  responses.ErrorResponse  "Bad request / validation error"
// @Failure      401  {object}  responses.ErrorResponse  "Unauthorized (wrong username/password)"
// @Failure      404  {object}  responses.ErrorResponse  "User not found"
// @Failure      422  {object}  responses.ErrorResponse  "Unprocessable entity / validation error"
// @Failure      429  {object}  responses.ErrorResponse  "Too many requests (rate-limited)"
// @Router       /auth/login [post]
func (h *authHandler) Login() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// ใช้ rate limit middleware ตรวจสอบก่อน
		// สร้าง handler ชั่วคราวเพื่อตรวจสอบ rate limit
		checkHandler := h.rateLimitMiddleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ถ้าผ่าน rate limit ให้ทำงานปกติ
		}))

		// สร้าง response writer จับ error
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// ตรวจสอบ rate limit
		checkHandler.ServeHTTP(recorder, r)

		// ถ้า status เป็น 429 แสดงว่าโดน rate limit
		if recorder.statusCode == http.StatusTooManyRequests {
			return
		}

		user := new(presenter.SignIn)

		contentType := r.Header.Get("Content-Type")

		var username, password string

		if strings.Contains(contentType, "application/json") {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				h.logger.Error("Failed to read request body")
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
			defer r.Body.Close()

			var jsonBody map[string]interface{}
			if err := json.Unmarshal(body, &jsonBody); err != nil {
				h.logger.Error("Failed to parse JSON body")
				render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
				return
			}

			username, _ = jsonBody["username"].(string)
			password, _ = jsonBody["password"].(string)
		} else {
			if err := r.ParseForm(); err != nil {
				h.logger.Error("Failed to parse form")
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
			username = r.FormValue("username")
			password = r.FormValue("password")
		}

		if username == "" || password == "" {
			h.logger.Warn("Login: missing username or password")
		}

		user = &presenter.SignIn{
			Username: username,
			Password: password,
		}

		h.logger.Info("Login attempt")

		err := utils.ValidateStruct(r.Context(), user)
		if err != nil {
			h.logger.Info("Login validation failed")
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		accessToken, refreshToken, err := h.usersUC.SignIn(
			r.Context(),
			user.Username,
			user.Password,
		)
		if err != nil {
			h.logger.Info("Login failed")
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		expiresIn := h.cfg.Jwt.AccessTokenExpireDuration * 60
		refreshExpiresIn := h.cfg.Jwt.RefreshTokenExpireDuration * 60

		h.setAuthCookies(w, accessToken, refreshToken, int(expiresIn), int(refreshExpiresIn))

		h.logger.Info("Login successful")

		render.Respond(w, r, map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "bearer",
			"expires_in":    expiresIn,
		})
	}
}

// responseRecorder สำหรับจับ status code
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	return rr.ResponseWriter.Write(b)
}

// RefreshToken godoc
// @Summary Refresh token
// @Description Get new access token from refresh token.
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} presenter.Token
// @Failure 400	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security BearerAuth
// @Router /auth/refresh [get]
func (h *authHandler) RefreshToken() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		refreshToken := middleware.RefreshTokenFromRequest(r)

		accessToken, refreshToken, err := h.usersUC.Refresh(ctx, refreshToken)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		h.setAuthCookies(w, accessToken, refreshToken,
			int(h.cfg.Jwt.AccessTokenExpireDuration*60),
			int(h.cfg.Jwt.RefreshTokenExpireDuration*60))

		render.Respond(w, r, presenter.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "bearer",
		})
	}
}

// GetPublicKey godoc
// @Summary Get public key
// @Description Get rsa public key to decode token.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} presenter.PublicKey
// @Failure 400	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Router /auth/publickey [get]
func (h *authHandler) GetPublicKey() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		publicKeyAccessToken, err := jwt.DecodeBase64(h.cfg.Jwt.AccessTokenPublicKey)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		publicKeyRefreshToken, err := jwt.DecodeBase64(h.cfg.Jwt.RefreshTokenPublicKey)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, presenter.PublicKey{
			PublicKeyAccessToken:  string(publicKeyAccessToken[:]),
			PublicKeyRefreshToken: string(publicKeyRefreshToken[:]),
		})
	}
}

// Logout godoc
// @Summary Logout
// @Description Logout, remove current refresh token in db.
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Success 200
// @Failure 400	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security BearerAuth
// @Router /auth/logout [get]
func (h *authHandler) Logout() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		refreshToken := middleware.RefreshTokenFromRequest(r)

		err := h.usersUC.Logout(ctx, refreshToken)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		h.clearAuthCookies(w)
	}
}

// LogoutAllToken godoc
// @Summary Logout all session
// @Description Logout all session of this user. This endpoint is available for both GET and POST.
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Success 200
// @Failure 400	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security BearerAuth
// @Router /auth/logoutall [get]
// @Router /auth/logoutall [post]
func (h *authHandler) LogoutAllToken() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		refreshToken := middleware.RefreshTokenFromRequest(r)

		id, err := h.usersUC.ParseIdFromRefreshToken(ctx, refreshToken)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = h.usersUC.LogoutAll(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		h.clearAuthCookies(w)
	}
}

// VerifyEmail godoc
// @Summary Verify user
// @Description Verify user using code from email.
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "offset" Format(code)
// @Success 200 {object} string
// @Failure 400	{object} responses.ErrorResponse
// @Router /auth/verifyemail [get]
func (h *authHandler) VerifyEmail() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		q := r.URL.Query()
		verificationCode := q.Get("code")

		err := h.usersUC.Verify(ctx, verificationCode)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse("Email verified successfully"))
	}
}

// ForgotPassword godoc
// @Summary Forgot password
// @Description Forgot password, code will send to email.
// @Tags auth
// @Accept json
// @Produce json
// @Param forgotPassword body presenter.ForgotPassword true "Forgot Password"
// @Success 200 {object} string
// @Failure 400	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Router /auth/forgotpassword [post]
func (h *authHandler) ForgotPassword() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		forgotPassword := new(presenter.ForgotPassword)

		err := json.NewDecoder(r.Body).Decode(&forgotPassword)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(r.Context(), forgotPassword)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		err = h.usersUC.ForgotPassword(ctx, forgotPassword.Email)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r,
			responses.CreateSuccessResponse("You will receive a reset email if user with that email exist"))
	}
}

// ResetPassword godoc
// @Summary Reset Password
// @Description Reset Password, using code from email.
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "code" Format(code)
// @Param resetPassword body presenter.ResetPassword true "Reset Password"
// @Success 200 {object} string
// @Failure 400	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Router /auth/resetpassword [patch]
func (h *authHandler) ResetPassword() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		q := r.URL.Query()
		resetToken := q.Get("code")

		resetPassword := new(presenter.ResetPassword)

		err := json.NewDecoder(r.Body).Decode(&resetPassword)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(r.Context(), resetPassword)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		err = h.usersUC.ResetPassword(
			ctx,
			resetToken,
			resetPassword.NewPassword,
			resetPassword.ConfirmPassword,
		)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r,
			responses.CreateSuccessResponse("Password data updated successfully, please re-login"))
	}
}
