package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/jacsmith21/lukabox/domain"
	"github.com/jacsmith21/lukabox/mock"
)

var PApi PillAPI
var PSvc mock.PillService
var t time.Time

func initPillAPI() {
	PSvc = mock.PillService{}
	PApi.PillService = &PSvc
	implPillServiceMethods()
}

func implPillServiceMethods() {
	t = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	PSvc.PillsFn = func(id int) ([]*domain.Pill, error) {
		if id != 1 {
			return nil, nil
		}

		pills := []*domain.Pill{
			{PillID: 1, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}, TimesOfDay: []time.Time{t}, Archived: false},
		}

		return pills, nil
	}
	PSvc.PillFn = func(id int) (*domain.Pill, error) {
		if id == 1 {
			pill := domain.Pill{PillID: 1, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}, TimesOfDay: []time.Time{t}, Archived: false}
			return &pill, nil
		} else if id == 2 {
			pill := domain.Pill{PillID: 2, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}, TimesOfDay: []time.Time{t}, Archived: false}
			return &pill, nil
		}
		return nil, nil
	}
	PSvc.UpdatePillFn = func(id int, pill *domain.Pill) error {
		if id != 1 {
			return errors.New("pill not found")
		}
		return nil
	}
}

type test struct {
	url    string
	status int
	body   string
	pill   *domain.Pill
}

var pillCtxTests = []test{
	{"/pills/1", http.StatusOK, "This is a test!", nil},
	{"/pills/3", http.StatusNotFound, "{\"message\":\"pill not found\"}", nil},
	{"/pills/bad", http.StatusBadRequest, "{\"message\":\"unable to parse parameter id\"}", nil},
}

func TestPillCtx(t *testing.T) {
	initPillAPI()

	r := chi.NewRouter()
	r.Route("/pills/{pillId}", func(r chi.Router) {
		r.Use(PApi.PillCtx)
		r.Get("/", func(w http.ResponseWriter, request *http.Request) {
			w.Write([]byte("This is a test!"))
		})
	})

	for _, test := range pillCtxTests {
		req, err := http.NewRequest("GET", test.url, nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if status := w.Code; status != test.status {
			t.Errorf("handler returned wrong status code: got %v want %v for %v", status, test.status, test.url)
		}

		body := strings.TrimSpace(w.Body.String())
		if body != test.body {
			t.Errorf("handler returned wrong body:\n%v\ninstead of:\n%v\nfor %v", body, test.body, test.url)
		}
	}
}

var pillsTests = []test{
	{"/users/1/pills", http.StatusOK, "[{\"pillId\":1,\"id\":1,\"name\":\"DoxyPoxy\",\"daysOfWeek\":[1,2,3,4,5,6,7],\"timesOfDay\":[\"2009-11-10T23:00:00Z\"],\"archived\":false}]", nil},
	{"/users/2/pills", http.StatusOK, "[]", nil},
}

func TestPills(t *testing.T) {
	initPillAPI()
	initUserAPI()

	r := chi.NewRouter()
	r.Route("/users/{userId}", func(r chi.Router) {
		r.Use(UApi.UserCtx)
		r.Get("/pills", PApi.Pills)
	})

	for _, test := range pillsTests {
		req, err := http.NewRequest("GET", test.url, nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if status := w.Code; status != test.status {
			t.Errorf("handler returned wrong status code: got %v want %v", status, test.status)
		}

		body := strings.TrimSpace(w.Body.String())
		if body != test.body {
			t.Errorf("handler returned wrong body:\n%v\ninstead of:\n%v\nfor %v", body, test.body, test.url)
		}
	}
}

var updatePillTests = []test{
	{"/users/1/pills/1", http.StatusOK, "", &domain.Pill{PillID: 1, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}, TimesOfDay: []time.Time{t}, Archived: false}},
	{"/users/1/pills/2", http.StatusOK, "", &domain.Pill{PillID: 1, UserID: 1, Name: "DoxyPoxy", DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}, TimesOfDay: []time.Time{t}, Archived: false}},
}

func TestUpdatePill(t *testing.T) {
	initPillAPI()
	initUserAPI()

	r := chi.NewRouter()
	r.Route("/users/{userId}", func(r chi.Router) {
		r.Use(UApi.UserCtx)
		r.Route("/pills/{pillId}", func(r chi.Router) {
			r.Use(PApi.PillCtx)
			r.Put("/", PApi.UpdatePill)
		})
	})

	for _, test := range updatePillTests {
		mPill, err := json.Marshal(test.pill)
		if err != nil {
			t.Fatal("error marshaling test pill")
		}

		req, err := http.NewRequest("PUT", test.url, bytes.NewReader(mPill))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if status := w.Code; status != test.status {
			t.Errorf("handler returned wrong status code: got %v want %v", status, test.status)
		}

		body := strings.TrimSpace(w.Body.String())
		if body != test.body {
			t.Errorf("handler returned wrong body:\n%v\ninstead of:\n%v\nfor %v", body, test.body, test.url)
		}
	}
}
