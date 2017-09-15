package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jacsmith21/lukabox/domain"
	"github.com/jacsmith21/lukabox/ext/log"
	"github.com/jacsmith21/lukabox/ext/render"
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
		if pillID == "" {
			render.WithMessage("pill id must be supplied").BadRequest(w, r)
			return
		}
		log.WithField("id", pillID).Debug("pill id from parameter")

		id, err := strconv.Atoi(pillID)
		if err != nil {
			render.WithMessage("unable to parse parameter id").BadRequest(w, r)
			return
		}

		pill, err := a.PillService.Pill(id)
		if err != nil {
			render.WithError(err).BadRequest(w, r)
			return
		}
		if pill == nil {
			render.WithMessage("pill not found").NotFound(w, r)
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
		render.WithError(err).BadRequest(w, r)
		render.WithError(err).BadRequest(w, r)
		return
	}

	if err := render.List(w, r, stc.NewPillListResponse(pills)); err != nil {
		render.WithError(err).BadRequest(w, r)
		return
	}
}

// UpdatePill updates a pill
func (a *PillAPI) UpdatePill(w http.ResponseWriter, r *http.Request) {
	log.WithField("method", "UpdatePill").Info("starting")
	user := r.Context().Value("user").(*domain.User)
	pill := r.Context().Value("pill").(*domain.Pill)

	if pill.UserID != user.ID {
		err := errors.New("parameter pill user id should match the parameter user ID")
		render.WithError(err).BadRequest(w, r)
		return
	}

	data := &stc.PillRequest{}
	if err := render.Bind(r, data); err != nil {
		render.WithError(err).BadRequest(w, r)
		return
	}

	p := data.Pill
	if p == nil {
		err := errors.New("a pill must be supplied")
		log.WithError(err).Debug("the pill from the request was nil")
		render.WithError(err).BadRequest(w, r)
		return
	}

	if p.ID == 0 {
		p.ID = pill.ID
	}
	if p.UserID == 0 {
		p.UserID = pill.UserID
	}

	if p.ID != pill.ID {
		err := errors.New("updated pill id must match the parameter pill id")
		render.WithError(err).BadRequest(w, r)
	}
	if p.UserID != user.ID {
		err := errors.New("updated pill user id does not match parameter user id")
		render.WithError(err).BadRequest(w, r)
	}

	a.PillService.UpdatePill(p.ID, p)

}
