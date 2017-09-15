package api

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/jacsmith21/lukabox/domain"
	"github.com/jacsmith21/lukabox/mock"
)

func TestOpen(t *testing.T) {
	bAPI := BoxAPI{}
	bSvc := mock.BoxService{}
	bAPI.BoxService = &bSvc

	uAPI := UserAPI{}
	uSvc := mock.UserService{}
	uAPI.UserService = &uSvc

	tests := []*test{
		{"/users/1/box/open", "PUT", `{"compId": 1, "time": "2012-11-01T22:08:41+00:00"}`, map[string]string{"Content-Type": "application/json"}, http.StatusCreated, ""},
	}

	bSvc.InsertOpenEventFn = func(openEvent *domain.OpenEvent) error {
		return nil
	}

	uSvc.UserByIDFn = func(id int) (*domain.User, error) {
		return &domain.User{ID: id, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}, nil
	}

	r := chi.NewRouter()
	r.Route("/users/{userId}/box", func(r chi.Router) {
		r.Use(uAPI.UserCtx)
		r.Use(bAPI.OpenEventRequestCtx)
		r.Put("/open", bAPI.Open)
	})

	runTests(t, r, tests)
}
