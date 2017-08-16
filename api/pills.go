package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/jacsmith21/lukabox/domain"
)

//PillAPI the services used
type PillAPI struct {
	PillService domain.PillService
}

//Pills returns the pills associated with the user
func (s *PillAPI) Pills(w http.ResponseWriter, r *http.Request) {
	var err error
	var pills []*domain.Pill

	//_, claims, _ := jwtauth.FromContext(r.Context())
	user := r.Context().Value("user").(*domain.User)

	if pills, err = s.PillService.Pills(user.ID); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	if err = render.RenderList(w, r, NewPillListResponse(pills)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
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

// PillResponse respose structure
type PillResponse struct {
	Pill *domain.Pill
}

// Render implementation
func (rd *PillResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
