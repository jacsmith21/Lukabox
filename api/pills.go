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

		pillID := chi.URLParam(r, "pillId")
		log.WithField("id", pillID).Debug("pill id from parameter")

		id, err := strconv.Atoi(pillID)
		if err != nil {
			render.Render(w, r, ErrBadRequest(errors.New("unable to parse parameter id")))
			return
		}

		pill, err := a.PillService.Pill(id)
		if err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}
		if pill == nil {
			render.Render(w, r, ErrNotFound(errors.New("pill not found")))
			return
		}

		ctx := context.WithValue(r.Context(), "pill", pill)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Pills returns the pills associated with the user
func (a *PillAPI) Pills(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(*domain.User)

	pills, err := a.PillService.Pills(user.ID)
	if err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if err := render.RenderList(w, r, stc.NewPillListResponse(pills)); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}
}

// UpdatePill updates a pill
func (a *PillAPI) UpdatePill(w http.ResponseWriter, r *http.Request) {
	log.WithField("method", "UpdatePill").Warn("starting")
	user := r.Context().Value("user").(*domain.User)
	pill := r.Context().Value("pill").(*domain.Pill)

	if pill.UserID != user.ID {
		err := errors.New("parameter pill user id should match the parameter user ID")
		render.Render(w, r, ErrBadRequest(err))
	}

	data := &stc.PillRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrBadRequest(err))
	}

	p := data.Pill
	if p.PillID == 0 {
		p.PillID = pill.PillID
	}
	if p.UserID == 0 {
		p.UserID = pill.UserID
	}

	if p.PillID != pill.PillID {
		err := errors.New("updated pill id must match the parameter pill id")
		render.Render(w, r, ErrBadRequest(err))
	}
	if p.UserID != user.ID {
		err := errors.New("updated pill user id does not match parameter user id")
		render.Render(w, r, ErrBadRequest(err))
	}

	a.PillService.UpdatePill(p.PillID, p)

}
