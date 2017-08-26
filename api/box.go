package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jacsmith21/lukabox/domain"
	log "github.com/jacsmith21/lukabox/ext/logrus"
)

// BoxAPI the services used
type BoxAPI struct {
	BoxService domain.BoxService
}

// BoxCtx creates a box context
func (a *BoxAPI) BoxCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithField("method", "BoxCtx").Info("starting")

		boxID := chi.URLParam(r, "boxID")
		if boxID == "" {
			render.Render(w, r, ErrBadRequest(errors.New("box id must be supplied")))
			return
		}
		log.WithField("id", boxID).Debug("box id from from parameter")

		id, err := strconv.Atoi(boxID)
		if err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		inter := r.Context().Value("user")
		if inter == nil {
			log.Error("no user in box context")
			render.Render(w, r, ErrBadRequest(errors.New("no user in box context")))
			return
		}
		user := inter.(*domain.User)

		box, err := a.BoxService.Box(user.ID, id)
		if err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}
		if box == nil {
			render.Render(w, r, ErrNotFound(errors.New("box not found")))
			return
		}

		ctx := context.WithValue(r.Context(), "box", box)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Open open a compartment in a box
func (a *BoxAPI) Open(w http.ResponseWriter, r *http.Request) {
	log.WithField("method", "Open").Info("starting")
	//box := r.Context().Value("box").(*domain.Box)

	//TODO modify model & context so that you pass in an open event
}
