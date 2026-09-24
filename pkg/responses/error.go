package responses

import (
	"net/http"

	"icmongolang/pkg/httpErrors"

	"github.com/go-chi/render"
)

// Render implements render.Renderer for Response[D].
// It sets the HTTP status code based on the error (if any).
func (e *Response[D]) Render(w http.ResponseWriter, r *http.Request) error {
	if e.Error != nil {
		render.Status(r, e.Error.Status)
	} else {
		render.Status(r, http.StatusOK)
	}
	return nil
}

// CreateErrorResponse builds a standardized error response.
func CreateErrorResponse(err error) render.Renderer {
	parsedErr := httpErrors.ParseErrors(err)

	return &Response[*string]{
		Data: nil,
		Error: &httpErrors.ErrResponse{
			Err:        parsedErr.GetErr(),
			Status:     parsedErr.GetStatus(),
			StatusText: parsedErr.GetStatusText(),
			Msg:        parsedErr.GetMsg(),
		},
		IsSuccess: false,
	}
}

// NewError creates a new error with an HTTP status code.
// The error message will be used as the response text.
func NewError(status int, msg string) error {
	return &httpError{
		status: status,
		msg:    msg,
	}
}

type httpError struct {
	status int
	msg    string
}

func (e *httpError) Error() string { return e.msg }
func (e *httpError) Status() int   { return e.status }

type statusCoder interface {
	Status() int
}

type simpleError struct {
	status int
	msg    string
}

func (e *simpleError) Error() string { return e.msg }
func (e *simpleError) Status() int   { return e.status }
