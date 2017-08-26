package stc

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/jacsmith21/lukabox/domain"
)

// PillResponse respose stc
type PillResponse struct {
	*domain.Pill
}

// Render implementation
func (pr *PillResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// PillRequest a pill request
type PillRequest struct {
	*domain.Pill
}

// Bind post-processing PillRequest
func (pr *PillRequest) Bind(r *http.Request) error {
	return nil
}

// NewPillListResponse create new pill list response
func NewPillListResponse(pills []*domain.Pill) []render.Renderer {
	list := []render.Renderer{}
	for _, pill := range pills {
		list = append(list, NewPillResponse(pill))
	}
	return list
}

// NewPillResponse create new response
func NewPillResponse(pill *domain.Pill) render.Renderer {
	resp := &PillResponse{Pill: pill}
	return resp
}
