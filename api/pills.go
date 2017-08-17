package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/jacsmith21/lukabox/domain"
	"github.com/jacsmith21/lukabox/structure"
)

//PillAPI the services used
type PillAPI struct {
	PillService domain.PillService
}

//Pills returns the pills associated with the user
func (s *PillAPI) Pills(w http.ResponseWriter, r *http.Request) {
	var err error
	var pills []*domain.Pill

	user := r.Context().Value("user").(*domain.User)

	if pills, err = s.PillService.Pills(user.ID); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	if err = render.RenderList(w, r, structure.NewPillListResponse(pills)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}
