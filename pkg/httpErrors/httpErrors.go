package httpErrors

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrorBadRequest                = errors.New("bad_request")
	ErrorNotFound                  = errors.New("not_found")
	ErrorUnauthorized              = errors.New("unauthorized")
	ErrorInternalServerError       = errors.New("internal_server_error")
	ErrorRequestTimeoutError       = errors.New("request_timeout")
	ErrorExistsEmailError          = errors.New("email_exists")
	ErrorExistsUsernameError       = errors.New("username_exists")
	ErrorInvalidJWTToken           = errors.New("invalid_jwt_token")
	ErrorInvalidJWTClaims          = errors.New("invalid_jwt_claims")
	ErrorValidation                = errors.New("validation")
	ErrorWrongPassword             = errors.New("wrong_password")
	ErrorTokenNotFound             = errors.New("token_not_found")
	ErrorInactiveUser              = errors.New("inactive_user")
	ErrorNotEnoughPrivileges       = errors.New("not_enough_privileges")
	ErrorGenToken                  = errors.New("generate_token_error")
	ErrorJson                      = errors.New("error_json_marshal")
	ErrorNotFoundRefreshTokenRedis = errors.New("not_found_refresh_token_redis")
	ErrorUserAlreadyVerified       = errors.New("user_already_verified")
	ErrorUserNotVerified           = errors.New("user_not_verified")
)

// Rest error interface
type ErrRest interface {
	GetErr() error
	GetStatus() int
	GetStatusText() string
	GetMsg() string
	Error() string
}

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err        error  `json:"-"`                                      // low-level runtime error
	Status     int    `json:"status" example:"404"`                   // http response status code
	StatusText string `json:"statusText" example:"not_found"`         // user-level status message
	Msg        string `json:"msg,omitempty" example:"not found user"` // application-level error message, for debugging
}

func (e *ErrResponse) GetErr() error {
	return e.Err
}

func (e *ErrResponse) GetStatus() int {
	return e.Status
}

func (e *ErrResponse) GetStatusText() string {
	return e.StatusText
}

func (e *ErrResponse) GetMsg() string {
	return e.Msg
}

// Error Error() interface method
func (e *ErrResponse) Error() string {
	return fmt.Sprintf("status: %d - statusText: %s - msg: %s - error: %v", e.Status, e.StatusText, e.Msg, e.Err)
}

func ErrBadRequest(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusBadRequest,
		StatusText: ErrorBadRequest.Error(),
		Msg:        err.Error(),
	}
}

func ErrNotFound(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusNotFound,
		StatusText: ErrorNotFound.Error(),
		Msg:        err.Error(),
	}
}

func ErrUnauthorized(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorUnauthorized.Error(),
		Msg:        err.Error(),
	}
}

func ErrInternalServer(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusInternalServerError,
		StatusText: ErrorInternalServerError.Error(),
		Msg:        err.Error(),
	}
}

func ErrValidation(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnprocessableEntity,
		StatusText: ErrorValidation.Error(),
		Msg:        err.Error(),
	}
}

func ErrRequestTimeoutError(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusRequestTimeout,
		StatusText: ErrorRequestTimeoutError.Error(),
		Msg:        err.Error(),
	}
}

func ErrInactiveUser(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusForbidden,
		StatusText: ErrorInactiveUser.Error(),
		Msg:        err.Error(),
	}
}

func ErrNotEnoughPrivileges(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusForbidden,
		StatusText: ErrorNotEnoughPrivileges.Error(),
		Msg:        err.Error(),
	}
}

func ErrInvalidJWTToken(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorInvalidJWTToken.Error(),
		Msg:        err.Error(),
	}
}

func ErrInvalidJWTClaims(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorInvalidJWTClaims.Error(),
		Msg:        err.Error(),
	}
}

func ErrWrongPassword(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorWrongPassword.Error(),
		Msg:        err.Error(),
	}
}

func ErrGenToken(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusBadRequest,
		StatusText: ErrorGenToken.Error(),
		Msg:        err.Error(),
	}
}

func ErrTokenNotFound(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorTokenNotFound.Error(),
		Msg:        err.Error(),
	}
}

func ErrJson(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorJson.Error(),
		Msg:        err.Error(),
	}
}

func ErrNotFoundRefreshTokenRedis(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorNotFoundRefreshTokenRedis.Error(),
		Msg:        err.Error(),
	}
}

func ErrUserAlreadyVerified(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorUserAlreadyVerified.Error(),
		Msg:        err.Error(),
	}
}

func ErrUserNotVerified(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorUserNotVerified.Error(),
		Msg:        err.Error(),
	}
}

// Parser of error string messages ,returns RestError

// ParseErrors now recognizes errors that implement statusCoder.
func ParseErrors(err error) ErrRest {
	// First, check if the error itself implements statusCoder
	if sc, ok := err.(statusCoder); ok {
		return &ErrResponse{
			Err:        err,
			Status:     sc.Status(),
			StatusText: http.StatusText(sc.Status()),
			Msg:        err.Error(),
		}
	}
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrNotFound(err)
	case errors.Is(err, context.DeadlineExceeded):
		return ErrRequestTimeoutError(err)
	case strings.Contains(err.Error(), "SQLSTATE"):
		return parseSqlErrors(err)
	default:
		if restErr, ok := err.(ErrRest); ok {
			return restErr
		}
		return ErrBadRequest(err)
	}
}

func ErrExistsEmail(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusBadRequest,
		StatusText: ErrorExistsEmailError.Error(),
		Msg:        err.Error(),
	}
}

func ErrExistsUsername(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusBadRequest,
		StatusText: ErrorExistsUsernameError.Error(),
		Msg:        err.Error(),
	}
}

// Parser sql error, returns RestError
func parseSqlErrors(err error) ErrRest {
	if strings.Contains(err.Error(), "23505") {
		switch {
		case strings.Contains(err.Error(), "idx_sd_user_username"):
			return ErrExistsUsername(err)
		case strings.Contains(err.Error(), "idx_sd_user_email"):
			return ErrExistsEmail(err)
		default:
			return ErrBadRequest(err)
		}
	}
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusBadRequest,
		StatusText: ErrorBadRequest.Error(),
		Msg:        err.Error(),
	}
}

// NewError creates a new error with an HTTP status code.
// The error message will be used as the response text.
func NewError(status int, msg string) error {
	return &statusError{status: status, msg: msg}
}

type statusError struct {
	status int
	msg    string
}

func (e *statusError) Error() string { return e.msg }
func (e *statusError) Status() int   { return e.status }

// statusCoder interface for errors that carry HTTP status.
type statusCoder interface {
	Status() int
}

// Just for swag
type SuccessResponse[D any] struct {
	Data      D    `json:"data"`
	IsSuccess bool `json:"is_success" example:"true"`
}

// SwaggerSuccessResponse is used only for Swagger documentation to avoid generic types.
type SwaggerSuccessResponse struct {
	Data      interface{} `json:"data"`
	IsSuccess bool        `json:"is_success"`
}
