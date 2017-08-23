package api

//ErrResponse represents error responses
import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrResponse error response stc
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	Message string `json:"message,omitempty"` // application-level error message, for debugging
}

// Render pre-processing before the ErrResponse is marshalled
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrBadRequest bad request response
func ErrBadRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		Message:        err.Error(),
	}
}

// ErrNotFound not found ErrBadRequest
func ErrNotFound(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusNotFound,
		Message:        err.Error(),
	}
}

// ErrInternalServerError not found ErrBadRequest
func ErrInternalServerError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		Message:        err.Error(),
	}
}

// ErrConflict conflict error
func ErrConflict(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusConflict,
		Message:        err.Error(),
	}
}

// ErrUnauthorized 401 error
var ErrUnauthorized = &ErrResponse{HTTPStatusCode: http.StatusUnauthorized, Message: "Unauthorized"}
