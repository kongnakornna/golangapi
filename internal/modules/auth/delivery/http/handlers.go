package http

import (
	"encoding/json"
	"net"
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
	"github.com/google/uuid"
)

// ============================================================================
// Log type IDs (ตรงกับ seed ใน sd_user_log_type)
// ============================================================================
const (
	LogTypeSignIn  = 1
	LogTypeSignOut = 2
)

// ============================================================================
// authHandler
// ============================================================================
type authHandler struct {
	usersUC users.UserUseCaseI
	logUC   users.LogUseCaseI
	cfg     *config.Config
	logger  logger.Logger
}

func CreateAuthHandler(
	uc users.UserUseCaseI,
	logUC users.LogUseCaseI,
	cfg *config.Config,
	logger logger.Logger,
) auth.Handlers {
	return &authHandler{
		cfg:     cfg,
		usersUC: uc,
		logUC:   logUC,
		logger:  logger,
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func (h *authHandler) secureCookies() bool {
	if h.cfg == nil {
		return false
	}
	env := strings.ToLower(h.cfg.Server.Env)
	return env == "production" || env == "staging"
}

func (h *authHandler) tokenExpirySeconds() (accessSec, refreshSec int) {
	return int(h.cfg.Jwt.AccessTokenExpireDuration * 60),
		int(h.cfg.Jwt.RefreshTokenExpireDuration * 60)
}

func (h *authHandler) respondError(w http.ResponseWriter, r *http.Request, err error) {
	if renderErr := render.Render(w, r, responses.CreateErrorResponse(err)); renderErr != nil {
		h.logger.Errorf("failed to render error response: %v", renderErr)
	}
}

// logTokenMeta — sanitized log (ไม่เปิดเผย token)
func (h *authHandler) logTokenMeta(where, token string, r *http.Request) {
	h.logger.Infof("[%s] token.len=%d hasAuthHdr=%t cookieNames=%v",
		where, len(token), r.Header.Get("Authorization") != "", cookieNames(r))
}

func cookieNames(r *http.Request) []string {
	cs := r.Cookies()
	names := make([]string, 0, len(cs))
	for _, c := range cs {
		names = append(names, c.Name)
	}
	return names
}

func bindRequest(r *http.Request, dst any) error {
	ct := r.Header.Get("Content-Type")
	if strings.Contains(ct, "application/json") {
		dec := json.NewDecoder(r.Body)
		dec.UseNumber()
		return dec.Decode(dst)
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	if len(r.Form) == 0 {
		return nil
	}
	formMap := make(map[string]string, len(r.Form))
	for k := range r.Form {
		formMap[k] = r.Form.Get(k)
	}
	buf, err := json.Marshal(formMap)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(strings.NewReader(string(buf)))
	dec.UseNumber()
	return dec.Decode(dst)
}

// extractClientMeta — ดึง IP/Device/Browser จาก request
func extractClientMeta(r *http.Request) (ip, device, browser string) {
	ip = r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
		if i := strings.IndexByte(ip, ','); i > 0 {
			ip = strings.TrimSpace(ip[:i])
		}
	}
	if ip == "" {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			ip = host
		} else {
			ip = r.RemoteAddr
		}
	}
	browser = r.UserAgent()
	device = r.Header.Get("X-Device-Access")
	if device == "" {
		device = detectDevice(browser)
	}
	return
}

func detectDevice(ua string) string {
	ua = strings.ToLower(ua)
	switch {
	case strings.Contains(ua, "iphone"), strings.Contains(ua, "ipad"):
		return "iOS"
	case strings.Contains(ua, "android"):
		return "Android"
	case strings.Contains(ua, "windows"):
		return "Windows"
	case strings.Contains(ua, "mac os"), strings.Contains(ua, "macintosh"):
		return "macOS"
	case strings.Contains(ua, "linux"):
		return "Linux"
	case strings.Contains(ua, "curl"), strings.Contains(ua, "postman"):
		return "API-Client"
	default:
		return "Unknown"
	}
}

// auditAuth — fire-and-forget audit log
func (h *authHandler) auditAuth(r *http.Request, uid uuid.UUID, logTypeID int, name, detail string) {
	if h.logUC == nil || uid == uuid.Nil {
		return
	}
	ip, device, browser := extractClientMeta(r)
	h.logUC.WriteAuditLogWithMeta(r.Context(), uid, logTypeID, name, detail, ip, device, browser)
}

// setAuthCookies stores access/refresh token in cookies.
func (h *authHandler) setAuthCookies(
	w http.ResponseWriter,
	accessToken, refreshToken string,
	accessExpiresIn, refreshExpiresIn int,
) {
	secure := h.secureCookies()

	http.SetCookie(w, &http.Cookie{
		Name:     middleware.AccessTokenCookieName,
		Value:    accessToken,
		Path:     "/",
		MaxAge:   accessExpiresIn,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.RefreshTokenCookieName,
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   refreshExpiresIn,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *authHandler) clearAuthCookies(w http.ResponseWriter) {
	secure := h.secureCookies()
	epoch := time.Unix(0, 0)

	http.SetCookie(w, &http.Cookie{
		Name:     middleware.AccessTokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  epoch,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.RefreshTokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  epoch,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// ============================================================================
// Handlers
// ============================================================================

// SignIn – POST /api/auth/signin
// @Summary      User login with email
// @Description  Authenticate using email + password. Returns JWT tokens + sets HttpOnly cookies.
// @Description  **Demo credentials**: `demo@example.com` / `Demo@1234`
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body presenter.UserSignIn true "User credentials"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Router       /auth/signin [post]
func (h *authHandler) SignIn() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := new(presenter.UserSignIn)

		if err := bindRequest(r, user); err != nil {
			h.logger.Error("SignIn: failed to decode request")
			h.respondError(w, r, httpErrors.ErrValidation(err))
			return
		}
		if user.Email == "" || user.Password == "" {
			h.respondError(w, r, httpErrors.ErrValidation(nil))
			return
		}
		if err := utils.ValidateStruct(r.Context(), user); err != nil {
			h.respondError(w, r, httpErrors.ErrValidation(err))
			return
		}

		accessToken, refreshToken, err := h.usersUC.SignIn(r.Context(), user.Email, user.Password)
		if err != nil {
			h.logger.Info("SignIn failed")
			h.respondError(w, r, err)
			return
		}

		// 🆕 Audit log
		if u, err := h.usersUC.GetByEmail(r.Context(), user.Email); err == nil && u != nil {
			h.auditAuth(r, u.ID, LogTypeSignIn, "User SignIn", "Sign in via /auth/signin")
		}

		expiresIn, refreshExpiresIn := h.tokenExpirySeconds()
		h.setAuthCookies(w, accessToken, refreshToken, expiresIn, refreshExpiresIn)

		render.Respond(w, r, map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "bearer",
			"expires_in":    expiresIn,
		})
	}
}

// Login – POST /api/auth/login
// @Summary      User login with username
// @Description  Authenticate using username (or email) + password. Returns JWT tokens + sets HttpOnly cookies.
// @Description  **Demo credentials**: `demouser` / `Demo@1234`
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body presenter.SignIn true "User credentials (username + password)"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Router       /auth/login [post]
func (h *authHandler) Login() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := new(presenter.SignIn)

		if err := bindRequest(r, user); err != nil {
			h.logger.Error("Login: failed to decode request")
			h.respondError(w, r, httpErrors.ErrValidation(err))
			return
		}
		if user.Username == "" || user.Password == "" {
			h.respondError(w, r, httpErrors.ErrValidation(nil))
			return
		}
		if err := utils.ValidateStruct(r.Context(), user); err != nil {
			h.respondError(w, r, httpErrors.ErrValidation(err))
			return
		}

		accessToken, refreshToken, err := h.usersUC.SignInByUsername(r.Context(), user.Username, user.Password)
		if err != nil {
			h.logger.Info("Login failed")
			h.respondError(w, r, err)
			return
		}

		// 🆕 Audit log
		if u, err := h.usersUC.GetByEmail(r.Context(), user.Username); err == nil && u != nil {
			h.auditAuth(r, u.ID, LogTypeSignIn, "User SignIn", "Sign in via /auth/login")
		}

		expiresIn, refreshExpiresIn := h.tokenExpirySeconds()
		h.setAuthCookies(w, accessToken, refreshToken, expiresIn, refreshExpiresIn)

		render.Respond(w, r, map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "bearer",
			"expires_in":    expiresIn,
		})
	}
}

// RefreshToken – GET /api/auth/refresh
// @Summary      Refresh access token
// @Description  Use HttpOnly refresh-token cookie to issue a new access token (rotation).
// @Tags         auth
// @Produce      json
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Router       /auth/refresh [get]
func (h *authHandler) RefreshToken() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		refreshToken := middleware.RefreshTokenFromRequest(r)

		h.logTokenMeta("RefreshToken", refreshToken, r)

		if refreshToken == "" {
			h.respondError(w, r, httpErrors.ErrUnauthorized(nil))
			return
		}

		accessToken, newRefresh, err := h.usersUC.Refresh(ctx, refreshToken)
		if err != nil {
			h.respondError(w, r, err)
			return
		}

		accessSec, refreshSec := h.tokenExpirySeconds()
		h.setAuthCookies(w, accessToken, newRefresh, accessSec, refreshSec)

		render.Respond(w, r, map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": newRefresh,
			"token_type":    "bearer",
		})
	}
}

// GetPublicKey – GET /api/auth/publickey
// @Summary      Get RSA public keys
// @Description  Return public keys used to verify access/refresh tokens (RS256).
// @Tags         auth
// @Produce      json
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Router       /auth/publickey [get]
func (h *authHandler) GetPublicKey() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		publicKeyAccessToken, err := jwt.DecodeBase64(h.cfg.Jwt.AccessTokenPublicKey)
		if err != nil {
			h.respondError(w, r, err)
			return
		}
		publicKeyRefreshToken, err := jwt.DecodeBase64(h.cfg.Jwt.RefreshTokenPublicKey)
		if err != nil {
			h.respondError(w, r, err)
			return
		}
		render.Respond(w, r, map[string]string{
			"public_key_access_token":  string(publicKeyAccessToken),
			"public_key_refresh_token": string(publicKeyRefreshToken),
		})
	}
}

// Logout – POST /api/auth/logout
// @Summary      Logout current session
// @Description  Revoke the current refresh token and clear auth cookies.
// @Tags         auth
// @Produce      json
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Router       /auth/logout [post]
func (h *authHandler) Logout() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		refreshToken := middleware.RefreshTokenFromRequest(r)

		h.logTokenMeta("Logout", refreshToken, r)

		if refreshToken == "" {
			h.clearAuthCookies(w)
			render.Respond(w, r, responses.CreateSuccessResponse("logged out"))
			return
		}

		// Parse user ID สำหรับ audit (ก่อน revoke)
		var uid uuid.UUID
		if id, err := h.usersUC.ParseIdFromRefreshToken(ctx, refreshToken); err == nil {
			uid = id
		}

		if err := h.usersUC.Logout(ctx, refreshToken); err != nil {
			h.respondError(w, r, err)
			return
		}

		if uid != uuid.Nil {
			h.auditAuth(r, uid, LogTypeSignOut, "User Signout", "Logout via /auth/logout")
		}

		h.clearAuthCookies(w)
		render.Respond(w, r, responses.CreateSuccessResponse("logged out"))
	}
}

// LogoutAllToken – POST /api/auth/logoutall
// @Summary      Logout all sessions
// @Description  Revoke all refresh tokens for the user.
// @Tags         auth
// @Produce      json
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Router       /auth/logoutall [post]
func (h *authHandler) LogoutAllToken() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		refreshToken := middleware.RefreshTokenFromRequest(r)

		h.logTokenMeta("LogoutAll", refreshToken, r)

		if refreshToken == "" {
			h.clearAuthCookies(w)
			h.respondError(w, r, httpErrors.ErrUnauthorized(nil))
			return
		}

		id, err := h.usersUC.ParseIdFromRefreshToken(ctx, refreshToken)
		if err != nil {
			h.respondError(w, r, err)
			return
		}

		if err := h.usersUC.LogoutAll(ctx, id); err != nil {
			h.respondError(w, r, err)
			return
		}

		h.auditAuth(r, id, LogTypeSignOut, "User Signout", "Logout all sessions")

		h.clearAuthCookies(w)
		render.Respond(w, r, responses.CreateSuccessResponse("logged out from all sessions"))
	}
}

// VerifyEmail – GET /api/auth/verifyemail
// @Summary      Verify email
// @Tags         auth
// @Produce      json
// @Param        code query string true "Verification code"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Router       /auth/verifyemail [get]
func (h *authHandler) VerifyEmail() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		verificationCode := r.URL.Query().Get("code")

		if err := h.usersUC.Verify(ctx, verificationCode); err != nil {
			h.respondError(w, r, err)
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse("Email verified successfully"))
	}
}

// ForgotPassword – POST /api/auth/forgotpassword
// @Summary      Request password reset
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body presenter.ForgotPassword true "Email"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Router       /auth/forgotpassword [post]
func (h *authHandler) ForgotPassword() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		forgotPassword := new(presenter.ForgotPassword)

		if err := bindRequest(r, forgotPassword); err != nil {
			h.respondError(w, r, err)
			return
		}
		if err := utils.ValidateStruct(r.Context(), forgotPassword); err != nil {
			h.respondError(w, r, httpErrors.ErrValidation(err))
			return
		}
		if err := h.usersUC.ForgotPassword(ctx, forgotPassword.Email); err != nil {
			h.respondError(w, r, err)
			return
		}
		render.Respond(w, r,
			responses.CreateSuccessResponse("You will receive a reset email if user with that email exist"))
	}
}

// ResetPassword – PATCH /api/auth/resetpassword
// @Summary      Reset password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        code query string true "Reset token"
// @Param        request body presenter.ResetPassword true "New password"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Router       /auth/resetpassword [patch]
func (h *authHandler) ResetPassword() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		resetToken := r.URL.Query().Get("code")
		resetPassword := new(presenter.ResetPassword)

		if err := bindRequest(r, resetPassword); err != nil {
			h.respondError(w, r, err)
			return
		}
		if err := utils.ValidateStruct(r.Context(), resetPassword); err != nil {
			h.respondError(w, r, httpErrors.ErrValidation(err))
			return
		}
		if err := h.usersUC.ResetPassword(
			ctx, resetToken,
			resetPassword.NewPassword, resetPassword.ConfirmPassword,
		); err != nil {
			h.respondError(w, r, err)
			return
		}
		render.Respond(w, r,
			responses.CreateSuccessResponse("Password data updated successfully, please re-login"))
	}
}