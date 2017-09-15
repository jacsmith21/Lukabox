package api

import (
	"context"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/jacsmith21/lukabox/domain"
	"github.com/jacsmith21/lukabox/ext/log"
	"github.com/jacsmith21/lukabox/ext/render"
	"github.com/jacsmith21/lukabox/stc"
)

// BoxAPI the services used
type BoxAPI struct {
	BoxService domain.BoxService
}

// OpenEventRequestCtx OpenEventRequestCtx
func (a *BoxAPI) OpenEventRequestCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithField("method", "OpenEventRequestCtx").Info("starting")
		openEventRequest := &stc.OpenEventRequest{}

		if err := render.Bind(r, openEventRequest); err != nil {
			log.WithError(err).Error("error binding open event req")
			render.WithError(err).BadRequest(w, r)
			return
		}

		openEvent := openEventRequest.OpenEvent
		log.WithField("openEvent", openEvent).Debug("open event from the request")

		ctx := context.WithValue(r.Context(), "open", openEvent)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CloseEventRequestCtx CloseEventRequestCtx
func (a *BoxAPI) CloseEventRequestCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithField("method", "CloseEventRequestCtx").Info("starting")
		closeEventRequest := &stc.CloseEventRequest{}

		if err := render.Bind(r, closeEventRequest); err != nil {
			log.WithError(err).Error("error binding close event req")
			render.WithError(err).BadRequest(w, r)
			return
		}

		closeEvent := closeEventRequest.CloseEvent
		log.WithField("closeEvent", closeEvent).Debug("close event from the request")

		ctx := context.WithValue(r.Context(), "close", closeEvent)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Open open a compartment in a box
func (a *BoxAPI) Open(w http.ResponseWriter, r *http.Request) {
	log.WithField("method", "Open").Info("starting")

	tmp := r.Context().Value("open")
	if tmp == nil {
		render.WithMessage("no open event in context").BadRequest(w, r)
		return
	}
	openEvent := tmp.(*domain.OpenEvent)

	validate := validator.New()
	if err := validate.Struct(openEvent); err != nil {
		render.WithError(err).BadRequest(w, r)
		return
	}

	if err := a.BoxService.InsertOpenEvent(openEvent); err != nil {
		render.WithError(err).InternalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Close open a compartment in a box
func (a *BoxAPI) Close(w http.ResponseWriter, r *http.Request) {
	log.WithField("method", "Close").Info("starting")

	tmp := r.Context().Value("close")
	if tmp == nil {
		render.WithMessage("no close event in context").BadRequest(w, r)
		return
	}
	closeEvent := tmp.(*domain.CloseEvent)

	validate := validator.New()
	if err := validate.Struct(closeEvent); err != nil {
		render.WithError(err).BadRequest(w, r)
		return
	}

	if err := a.BoxService.InsertCloseEvent(closeEvent); err != nil {
		render.WithError(err).InternalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
