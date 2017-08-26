package api

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/jacsmith21/lukabox/domain"
	"github.com/jacsmith21/lukabox/mock"
)

func TestPillCtx(t *testing.T) {
	pAPI := PillAPI{}
	pSvc := mock.PillService{}
	pAPI.PillService = &pSvc

	tests := []*test{
		{"/pills/1", "", "GET", nil, http.StatusOK, "This is a test!"},
		{"/pills/3", "", "GET", nil, http.StatusNotFound, `{"message":"pill not found"}`},
		{"/pills/bad", "", "GET", nil, http.StatusBadRequest, `{"message":"unable to parse parameter id"}`},
	}

	d := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	pSvc.PillFn = func(id int) (*domain.Pill, error) {
		if id == 1 {
			pill := domain.Pill{ID: 1, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}, TimesOfDay: []time.Time{d}, Archived: false}
			return &pill, nil
		} else if id == 2 {
			pill := domain.Pill{ID: 2, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}, TimesOfDay: []time.Time{d}, Archived: false}
			return &pill, nil
		}
		return nil, nil
	}

	r := chi.NewRouter()
	r.Route("/pills/{pillId}", func(r chi.Router) {
		r.Use(pAPI.PillCtx)
		r.Get("/", func(w http.ResponseWriter, request *http.Request) {
			w.Write([]byte("This is a test!"))
		})
	})

	runTests(t, r, tests)
}

func TestPills(t *testing.T) {
	pAPI := PillAPI{}
	pSvc := mock.PillService{}
	pAPI.PillService = &pSvc

	uAPI := UserAPI{}
	uSvc := mock.UserService{}
	uAPI.UserService = &uSvc

	tests := []*test{
		{"/users/1/pills", "", "GET", nil, http.StatusOK, `[{"pillId":1,"id":1,"name":"DoxyPoxy","daysOfWeek":[1],"timesOfDay":["2009-11-10T23:00:00Z"],"archived":false}]`},
		{"/users/2/pills", "", "GET", nil, http.StatusOK, "[]"},
	}

	d := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	pSvc.PillsFn = func(id int) ([]*domain.Pill, error) {
		if id != 1 {
			return nil, nil
		}
		pills := []*domain.Pill{
			{ID: 1, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1}, TimesOfDay: []time.Time{d}, Archived: false},
		}
		return pills, nil
	}

	uSvc.UserByIDFn = func(id int) (*domain.User, error) {
		return &domain.User{ID: id, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}, nil
	}

	r := chi.NewRouter()
	r.Route("/users/{userId}", func(r chi.Router) {
		r.Use(uAPI.UserCtx)
		r.Get("/pills", pAPI.Pills)
	})

	runTests(t, r, tests)
}

func TestUpdatePill(t *testing.T) {
	pAPI := PillAPI{}
	pSvc := mock.PillService{}
	pAPI.PillService = &pSvc

	uAPI := UserAPI{}
	uSvc := mock.UserService{}
	uAPI.UserService = &uSvc

	var tests = []*test{
		{"/users/1/pills/1", "POST", `{"pillId":1,"id":1,"name":"DoxyPoxy","daysOfWeek":[1],"timesOfDay":["2009-11-10T23:00:00Z"],"archived":false}`, map[string]string{"Content-Type": "application/json"}, http.StatusOK, ""},
		{"/users/1/pills/2", "POST", `{"pillId":1,"id":1,"name":"DoxyPoxy", "archived":false}`, map[string]string{"Content-Type": "application/json"}, http.StatusBadRequest, `{"message":"updated pill id must match the parameter pill id"}`},
		{"/users/2/pills/1", "POST", `{"pillId":1,"id":1,"name":"DoxyPoxy", "archived":false}`, map[string]string{"Content-Type": "application/json"}, http.StatusBadRequest, `{"message":"parameter pill user id should match the parameter user ID"}`},
	}

	d := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	pSvc.PillFn = func(id int) (*domain.Pill, error) {
		return &domain.Pill{ID: id, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}, TimesOfDay: []time.Time{d}, Archived: false}, nil
	}

	uSvc.UserByIDFn = func(id int) (*domain.User, error) {
		return &domain.User{ID: id, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}, nil
	}

	pSvc.UpdatePillFn = func(id int, pill *domain.Pill) error {
		if id != 1 {
			return errors.New("pill not found")
		}
		return nil
	}

	r := chi.NewRouter()
	r.Route("/users/{userId}", func(r chi.Router) {
		r.Use(uAPI.UserCtx)
		r.Route("/pills/{pillId}", func(r chi.Router) {
			r.Use(pAPI.PillCtx)
			r.Post("/", pAPI.UpdatePill)
		})
	})

	runTests(t, r, tests)
}
