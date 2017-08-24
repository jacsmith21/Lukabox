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
	BoxService  domain.BoxService
	UserService domain.UserService
}

// BoxCtx creates a box context
func (a *BoxAPI) BoxCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithField("method", "BoxCtx").Info("starting")

		boxID := chi.URLParam(r, "boxId")
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

		box, err := a.BoxService.BoxByID(id)
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

// CompCtx compartment context
func (a *BoxAPI) CompCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithField("method", "CompCtx").Info("starting")

		c := r.URL.Query().Get("compartment")
		if c == "" {
			render.Render(w, r, ErrBadRequest(errors.New("box id must be supplied")))
			return
		}
		log.WithField("compartment", c).Debug("compartment number from the query")

		compartment, err := strconv.Atoi(c)
		if err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		ctx := context.WithValue(r.Context(), "compartment", compartment)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Open open a compartment in a box
func (a *BoxAPI) Open(w http.ResponseWriter, r *http.Request) {
	log.WithField("method", "Open").Info("starting")
	box := r.Context().Value("box").(*domain.Box)
	//compartment := r.Context().Value("compartment").(int)

	userID := box.UserID
	_, err := a.UserService.UserByID(userID)
	if err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	//TODO modify model & context so that you pass in an open event
}
