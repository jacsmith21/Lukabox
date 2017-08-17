package structure

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/jacsmith21/lukabox/domain"
)

// PillResponse respose structure
type PillResponse struct {
	Pill *domain.Pill
}

// Render implementation
func (rd *PillResponse) Render(w http.ResponseWriter, r *http.Request) error {
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
