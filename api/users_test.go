package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/jacsmith21/lukabox/domain"
	"github.com/jacsmith21/lukabox/mock"
)

var USvc mock.UserService
var UApi UserAPI

func initUserAPI() {
	USvc = mock.UserService{}
	UApi.UserService = &USvc
}

type test struct {
	url     string
	method  string
	reqBody string
	status  int
	resBody string
}

func runTests(t *testing.T, r *chi.Mux, tests []*test) {
	for _, test := range tests {
		req, err := http.NewRequest(test.method, test.url, strings.NewReader(test.reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if status := w.Code; status != test.status {
			t.Errorf("handler returned wrong status code: got %v want %v for %s %s", status, test.status, test.method, test.url)
		}

		body := strings.TrimSpace(w.Body.String())
		if body != test.resBody {
			t.Errorf("handler returned wrong body:\n%v\ninstead of:\n%v\nfor %v", body, test.resBody, test.url)
		}
	}
}

func TestUserCtx(t *testing.T) {
	initUserAPI()

	userCtxTests := []*test{
		{"/users/1", "GET", "", http.StatusOK, "This is a test!"},
		{"/users/3", "GET", "", http.StatusInternalServerError, "{\"message\":\"test error\"}"},
		{"/users/4", "GET", "", http.StatusNotFound, "{\"message\":\"user not found\"}"},
		{"/users/ahh", "GET", "", http.StatusBadRequest, "{\"message\":\"unable to parse parameter id\"}"},
	}

	r := chi.NewRouter()
	r.Route("/users/{userId}", func(r chi.Router) {
		r.Use(UApi.UserCtx)
		r.Get("/", func(w http.ResponseWriter, request *http.Request) {
			w.Write([]byte("This is a test!"))
		})
	})

	runTests(t, r, userCtxTests)
}

func TestUserRequestCtx(t *testing.T) {
	initUserAPI()

	userRequestCtxTests := []*test{
		{"/users", "PUT", "{\"email\":\"jacob.smith@unb.ca\",\"password\":\"password\",\"firstName\":\"Jacob\",\"lastName\":\"Smith\"}", http.StatusOK, "This is a test!"},
		{"/users", "PUT", "{\"whatisthis\":\"jacob.smith@unb.ca\",\"password\":\"password\",\"firstName\":\"Jacob\",\"lastName\":\"Smith\"}", http.StatusOK, "This is a test!"},
	}

	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		r.Use(UApi.UserRequestCtx)
		r.Put("/", func(w http.ResponseWriter, request *http.Request) {
			w.Write([]byte("This is a test!"))
		})
	})

	runTests(t, r, userRequestCtxTests)
}

func TestUserByID(t *testing.T) {
	initUserAPI()

	userByIDTests := []*test{
		{"/users/1", "GET", "", http.StatusOK, "{\"id\":1,\"password\":\"password\",\"email\":\"jacob.smith@unb.ca\",\"firstName\":\"Jacob\",\"lastName\":\"Smith\",\"archived\":false}"},
		{"/users/2", "GET", "", http.StatusOK, "{\"id\":2,\"password\":\"password\",\"email\":\"jacob.smith@unb.ca\",\"firstName\":\"Jacob\",\"lastName\":\"Smith\",\"archived\":false}"},
		{"/users/3", "GET", "", http.StatusInternalServerError, "{\"message\":\"test error\"}"},
		{"/users/4", "GET", "", http.StatusNotFound, "{\"message\":\"user not found\"}"},
	}

	r := chi.NewRouter()
	r.Route("/users/{userId}", func(r chi.Router) {
		r.Use(UApi.UserCtx)
		r.Get("/", UApi.UserByID)
	})

	runTests(t, r, userByIDTests)
}

func TestUsers(t *testing.T) {
	initUserAPI()

	usersTests := []*test{
		{"/users", "GET", "", http.StatusOK, "[{\"id\":1,\"password\":\"password\",\"email\":\"jacob.smith@unb.ca\",\"firstName\":\"Jacob\",\"lastName\":\"Smith\",\"archived\":false}]"},
		{"/users", "GET", "", http.StatusOK, "[]"},
		{"/users", "GET", "", http.StatusInternalServerError, "{\"message\":\"test error\"}"},
	}

	r := chi.NewRouter()
	r.Get("/users", UApi.Users)

	runTests(t, r, usersTests)
}

func TestCreateUser(t *testing.T) {
	initUserAPI()

	createUserTests := []*test{
		{"/users", "PUT", "{\"email\":\"jacob.smith@unb.ca\",\"password\":\"password\",\"firstName\":\"Jacob\",\"lastName\":\"Smith\"}", http.StatusCreated, ""},
	}

	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		r.Use(UApi.UserRequestCtx)
		r.Put("/", UApi.CreateUser)
	})

	runTests(t, r, createUserTests)
}

func TestUpdateUser(t *testing.T) {
	initUserAPI()

	user := domain.User{ID: 1, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}

	var m []byte
	var err error
	if m, err = json.Marshal(user); err != nil {
		t.Fatal("error marshaling test user")
	}

	req, err := http.NewRequest("POST", "/users/1", bytes.NewReader(m))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Route("/users/{userId}", func(r chi.Router) {
		r.Use(UApi.UserCtx)
		r.Post("/", UApi.UpdateUser)
	})
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
