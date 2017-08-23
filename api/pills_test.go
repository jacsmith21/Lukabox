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

var PApi PillAPI
var PSvc mock.PillService

func initPillAPI() {
	PSvc = mock.PillService{}
	PApi.PillService = &PSvc
}

func TestPillCtx(t *testing.T) {
	initPillAPI()

	tests := []*test{
		{"/pills/1", "", "", http.StatusOK, "This is a test!"},
		{"/pills/3", "", "", http.StatusNotFound, `{"message":"pill not found"}`},
		{"/pills/bad", "", "", http.StatusBadRequest, `{"message":"unable to parse parameter id"}`},
	}

	d := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	PSvc.PillFn = func(id int) (*domain.Pill, error) {
		if id == 1 {
			pill := domain.Pill{PillID: 1, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}, TimesOfDay: []time.Time{d}, Archived: false}
			return &pill, nil
		} else if id == 2 {
			pill := domain.Pill{PillID: 2, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}, TimesOfDay: []time.Time{d}, Archived: false}
			return &pill, nil
		}
		return nil, nil
	}

	r := chi.NewRouter()
	r.Route("/pills/{pillId}", func(r chi.Router) {
		r.Use(PApi.PillCtx)
		r.Get("/", func(w http.ResponseWriter, request *http.Request) {
			w.Write([]byte("This is a test!"))
		})
	})

	runTests(t, r, tests)
}

func TestPills(t *testing.T) {
	initPillAPI()
	initUserAPI()

	tests := []*test{
		{"/users/1/pills", "", "", http.StatusOK, `[{"pillId":1,"id":1,"name":"DoxyPoxy","daysOfWeek":[1],"timesOfDay":["2009-11-10T23:00:00Z"],"archived":false}]`},
		{"/users/2/pills", "", "", http.StatusOK, "[]"},
	}

	d := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	PSvc.PillsFn = func(id int) ([]*domain.Pill, error) {
		if id != 1 {
			return nil, nil
		}
		pills := []*domain.Pill{
			{PillID: 1, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1}, TimesOfDay: []time.Time{d}, Archived: false},
		}
		return pills, nil
	}

	r := chi.NewRouter()
	r.Route("/users/{userId}", func(r chi.Router) {
		r.Use(UApi.UserCtx)
		r.Get("/pills", PApi.Pills)
	})

	runTests(t, r, tests)
}

func TestUpdatePill(t *testing.T) {
	initPillAPI()
	initUserAPI()

	var tests = []*test{
		{"/users/1/pills/1", "POST", `{"pillId":1,"id":1,"name":"DoxyPoxy","daysOfWeek":[1],"timesOfDay":["2009-11-10T23:00:00Z"],"archived":false}`, http.StatusOK, ""},
		{"/users/1/pills/2", "POST", `{"pillId":1,"id":1,"name":"DoxyPoxy", "archived":false}`, http.StatusBadRequest, `{"message":"updated pill id must match the parameter pill id"}`},
		{"/users/2/pills/1", "POST", `{"pillId":1,"id":1,"name":"DoxyPoxy", "archived":false}`, http.StatusBadRequest, `{"message":"parameter pill user id should match the parameter user ID"}`},
	}

	d := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	PSvc.PillFn = func(id int) (*domain.Pill, error) {
		return &domain.Pill{PillID: id, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}, TimesOfDay: []time.Time{d}, Archived: false}, nil
	}

	PSvc.UpdatePillFn = func(id int, pill *domain.Pill) error {
		if id != 1 {
			return errors.New("pill not found")
		}
		return nil
	}

	r := chi.NewRouter()
	r.Route("/users/{userId}", func(r chi.Router) {
		r.Use(UApi.UserCtx)
		r.Route("/pills/{pillId}", func(r chi.Router) {
			r.Use(PApi.PillCtx)
			r.Post("/", PApi.UpdatePill)
		})
	})

	runTests(t, r, tests)
}
