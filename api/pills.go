package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jacsmith21/lukabox/domain"
	"github.com/jacsmith21/lukabox/stc"
)

//PillAPI the services used
type PillAPI struct {
	PillService domain.PillService
}

// PillCtx is used to create a user context by id
func (a *PillAPI) PillCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithField("method", "PillCtx").Info("starting")

		pillID := chi.URLParam(r, "id")
		if pillID == "" {
			render.Render(w, r, ErrInvalidRequest(errors.New("pillID cannot be empty")))
			return
		}
		log.WithField("id", pillID).Debug("pill id from parameter")

		id, err := strconv.Atoi(pillID)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		pill, err := a.PillService.Pill(id)
		if err != nil {
			render.Render(w, r, ErrNotFound(err))
			return
		}

		ctx := context.WithValue(r.Context(), "pill", pill)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Pills returns the pills associated with the user
func (a *PillAPI) Pills(w http.ResponseWriter, r *http.Request) {
	var err error
	var pills []*domain.Pill

	user := r.Context().Value("user").(*domain.User)

	if pills, err = a.PillService.Pills(user.ID); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	if err = render.RenderList(w, r, stc.NewPillListResponse(pills)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// UpdatePill updates a pill
func (a *PillAPI) UpdatePill(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*domain.User)
	pill := r.Context().Value("pill").(*domain.Pill)

	if pill.UserID != user.ID {
		err := errors.New("pill UserID should match the parameter user ID")
		render.Render(w, r, ErrInvalidRequest(err))
	}

	data := &stc.PillRequest{Pill: pill}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
	}

	pill = data.Pill
	if pill.UserID != user.ID {
		err := errors.New("pill UserID does not match parameter user ID")
		render.Render(w, r, ErrInvalidRequest(err))
	}

	a.PillService.UpdatePill(pill.PillID, pill)

}
