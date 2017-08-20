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

var PApi PillAPI
var PSvc mock.PillService

func initPillAPI() {
	PSvc = mock.PillService{}
	PApi.PillService = &PSvc
	implPillServiceMethods()
}

func implPillServiceMethods() {
	PSvc.PillsFn = func(id int) ([]*domain.Pill, error) {
		if id != 1 {
			return nil, errors.New("expected id to be 1")
		}

		pills := []*domain.Pill{
			{PillID: 1, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}, TimesOfDay: []time.Time{time.Now()}, Archived: false},
		}

		return pills, nil
	}
	PSvc.PillFn = func(id int) (*domain.Pill, error) {
		if id != 1 {
			return nil, errors.New("expected id to be 1")
		}
		pill := domain.Pill{PillID: 1, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}, TimesOfDay: []time.Time{time.Now()}, Archived: false}
		return &pill, nil
	}
	PSvc.UpdatePillFn = func(id int, pill *domain.Pill) error {
		if id != 1 {
			return errors.New("pill not found")
		}
		return nil
	}
}

func TestPillCtx(t *testing.T) {
	initPillAPI()

	req, err := http.NewRequest("GET", "/pills/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Route("/pills/{id}", func(r chi.Router) {
		r.Use(PApi.PillCtx)
		r.Get("/", func(w http.ResponseWriter, request *http.Request) {
			w.Write([]byte("This is a test!"))
		})
	})
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !PSvc.PillInvoked {
		t.Fatal("expected Pill to be invoked")
	}
}

func TestPills(t *testing.T) {
	initPillAPI()
	initUserAPI()

	req, err := http.NewRequest("GET", "/users/1/pills", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Route("/users/{id}", func(r chi.Router) {
		r.Use(UApi.UserCtx)
		r.Get("/pills", PApi.Pills)
	})
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !PSvc.PillsInvoked {
		t.Fatal("expected Pills to be invoked")
	}
}

func TestUpdatePill(t *testing.T) {
	initPillAPI()
	initUserAPI()

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
		r.Use(UApi.UserCtx)
		r.Route("/pills/{id}", func(r chi.Router) {
			r.Use(PApi.PillCtx)
			r.Get("/", PApi.UpdatePill)
		})
	})
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v\nbody: %v", status, http.StatusOK, w.Body)
	}

	if !USvc.UserByIDInvoked {
		t.Fatal("expected UsersByID to be invoked")
	}

	if !PSvc.PillInvoked {
		t.Fatal("expected Pill to be invoked")
	}

	if !PSvc.UpdatePillInvoked {
		t.Fatal("expected UpdatePill to be invoked")
	}
}
