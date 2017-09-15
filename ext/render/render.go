package render

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/jacsmith21/lukabox/ext/log"
)

// Bind Bind
func Bind(r *http.Request, v render.Binder) error {
	return render.Bind(r, v)
}

// Instance Instance
func Instance(w http.ResponseWriter, r *http.Request, v render.Renderer) error {
	return render.Render(w, r, v)
}

// List List
func List(w http.ResponseWriter, r *http.Request, l []render.Renderer) error {
	return render.RenderList(w, r, l)
}

// Renderer interface
type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request) error
}

// ErrRenderer ErrRenderer
type ErrRenderer struct {
	HTTPStateCode int    `json:"-"`
	Message       string `json:"message,omitempty"`
}

// Render Renderer implementation
func (ren *ErrRenderer) Render(w http.ResponseWriter, r *http.Request) error {
	code := ren.HTTPStateCode
	if code == 0 {
		log.Debug("no status code provided")
	} else {
		render.Status(r, code)
	}
	return nil
}

// WithError renders with error
func WithError(err error) *ErrRenderer {
	renderer := &ErrRenderer{Message: err.Error()}
	return renderer
}

// WithMessage renders with message
func WithMessage(message string) *ErrRenderer {
	renderer := &ErrRenderer{Message: message}
	return renderer
}

// BadRequest renders a bas request
func (ren *ErrRenderer) BadRequest(w http.ResponseWriter, r *http.Request) {
	ren.HTTPStateCode = http.StatusBadRequest
	render.Render(w, r, ren)
}

// NotFound renders a not found
func (ren *ErrRenderer) NotFound(w http.ResponseWriter, r *http.Request) {
	ren.HTTPStateCode = http.StatusNotFound
	render.Render(w, r, ren)
}

// InternalServerError renders an internal server error
func (ren *ErrRenderer) InternalServerError(w http.ResponseWriter, r *http.Request) {
	ren.HTTPStateCode = http.StatusInternalServerError
	render.Render(w, r, ren)
}

// Conflict renders an internal server error
func (ren *ErrRenderer) Conflict(w http.ResponseWriter, r *http.Request) {
	ren.HTTPStateCode = http.StatusConflict
	render.Render(w, r, ren)
}

// Unauthorized Unauthorized
func Unauthorized(w http.ResponseWriter, r *http.Request) {
	ren := &ErrRenderer{}
	ren.HTTPStateCode = http.StatusUnauthorized
	ren.Message = "unauthorized"
	render.Render(w, r, ren)
}
