package api

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/jacsmith21/lukabox/domain"
	"github.com/jacsmith21/lukabox/mock"
)

func TestBoxCtx(t *testing.T) {
	bAPI := BoxAPI{}
	bSvc := mock.BoxService{}
	bAPI.BoxService = &bSvc

	tests := []*test{
		{"users/1/boxes/1", "GET", "", nil, http.StatusOK, "This is a test!"},
		{"users/1/boxes/3", "GET", "", nil, http.StatusNotFound, `{"message":"pill not found"}`},
		{"users/1/boxes/bad", "GET", "", nil, http.StatusBadRequest, `{"message":"unable to parse parameter id"}`},
	}

	bSvc.BoxFn = func(userID int, id int) (*domain.Box, error) {
		box := domain.Box{UserID: 1, ID: 1}
		return &box, nil
	}

	r := chi.NewRouter()
	r.Route("/user/{userID}/boxes/{boxID}", func(r chi.Router) {
		r.Use(bAPI.BoxCtx)
		r.Get("/", func(w http.ResponseWriter, request *http.Request) {
			w.Write([]byte("This is a test!"))
		})
	})

	runTests(t, r, tests)
}

func TestOpen(t *testing.T) {
	bAPI := BoxAPI{}
	bSvc := mock.BoxService{}
	bAPI.BoxService = &bSvc

	tests := []*test{
		{"/users/1/boxes/1/open", "PUT", "", nil, http.StatusCreated, "tbd"},
	}

	bSvc.BoxFn = func(userID int, id int) (*domain.Box, error) {
		box := domain.Box{UserID: 1, ID: 1}
		return &box, nil
	}

	r := chi.NewRouter()
	r.Route("/users/{userID}/boxes/{boxID}", func(r chi.Router) {
		r.Use(bAPI.BoxCtx)
		r.Put("/open", bAPI.Open)
	})

	runTests(t, r, tests)
}
