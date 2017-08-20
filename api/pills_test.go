package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/jacsmith21/lukabox/domain"
	"github.com/jacsmith21/lukabox/mock"
)

func implPillServiceMethods(us *mock.PillService) {
	us.PillsFn = func(id int) ([]*domain.Pill, error) {
		if id != 1 {
			return nil, errors.New("expected id to be 1")
		}

		pills := []*domain.Pill{
			{PillID: 1, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}, TimesOfDay: []time.Time{time.Now()}, Archived: false},
		}

		return pills, nil
	}
	us.PillFn = func(id int) (*domain.Pill, error) {
		if id != 1 {
			return nil, errors.New("expected id to be 1")
		}
		pill := domain.Pill{PillID: 1, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}, TimesOfDay: []time.Time{time.Now()}, Archived: false}
		return &pill, nil
	}
	us.UpdatePillFn = func(id int, pill *domain.Pill) error {
		if id != 1 {
			return errors.New("pill not found")
		}
		return nil
	}
}

func TestPillCtx(t *testing.T) {
	var ps mock.PillService
	var pa PillAPI
	pa.PillService = &ps
	implPillServiceMethods(&ps)

	req, err := http.NewRequest("GET", "/pills/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Route("/pills/{id}", func(r chi.Router) {
		r.Use(pa.PillCtx)
		r.Get("/", func(w http.ResponseWriter, request *http.Request) {
			w.Write([]byte("This is a test!"))
		})
	})
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !ps.PillInvoked {
		t.Fatal("expected Pill to be invoked")
	}
}

func TestPills(t *testing.T) {
	var ps mock.PillService
	var pa PillAPI
	pa.PillService = &ps
	implPillServiceMethods(&ps)

	var us mock.UserService
	var ua UserAPI
	ua.UserService = &us
	implUserServiceMethods(&us)

	req, err := http.NewRequest("GET", "/users/1/pills", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Route("/users/{id}", func(r chi.Router) {
		r.Use(ua.UserCtx)
		r.Get("/pills", pa.Pills)
	})
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !ps.PillsInvoked {
		t.Fatal("expected Pills to be invoked")
	}
}

func TestUpdatePill(t *testing.T) {
	var ps mock.PillService
	var pa PillAPI
	pa.PillService = &ps
	implPillServiceMethods(&ps)

	var us mock.UserService
	var ua UserAPI
	ua.UserService = &us
	implUserServiceMethods(&us)

	pill := domain.Pill{PillID: 1, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}, TimesOfDay: []time.Time{time.Now()}, Archived: false}
	m, err := json.Marshal(pill)
	if err != nil {
		t.Fatal("error marshaling test user")
	}

	req, err := http.NewRequest("GET", "/users/1/pills/1", bytes.NewReader(m))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Route("/users/{id}", func(r chi.Router) {
		r.Use(ua.UserCtx)
		r.Route("/pills/{id}", func(r chi.Router) {
			r.Use(pa.PillCtx)
			r.Get("/", pa.UpdatePill)
		})
	})
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v\nbody: %v", status, http.StatusOK, w.Body)
	}

	if !us.UserByIDInvoked {
		t.Fatal("expected UsersByID to be invoked")
	}

	if !ps.PillInvoked {
		t.Fatal("expected Pill to be invoked")
	}

	if !ps.UpdatePillInvoked {
		t.Fatal("expected UpdatePill to be invoked")
	}
}
