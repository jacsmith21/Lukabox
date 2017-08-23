package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
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
	headers map[string]string
	status  int
	resBody string
}

func runTests(t *testing.T, r *chi.Mux, tests []*test) {
	for i, test := range tests {
		req, err := http.NewRequest(test.method, test.url, strings.NewReader(test.reqBody))
		if err != nil {
			t.Fatal(err)
		}

		for k, v := range test.headers {
			req.Header.Add(k, v)
		}

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if status := w.Code; status != test.status {
			t.Errorf("handler returned wrong status code: got %v want %v on iteration %d", status, test.status, i)
		}

		body := strings.TrimSpace(w.Body.String())
		if body != test.resBody {
			t.Errorf("handler returned wrong body:\n%v\ninstead of:\n%v\n on iteration %d", body, test.resBody, i)
		}
	}
}

func TestUserCtx(t *testing.T) {
	initUserAPI()

	userCtxTests := []*test{
		{"/users/1", "GET", "", nil, http.StatusOK, "This is a test!"},
		{"/users/3", "GET", "", nil, http.StatusInternalServerError, `{"message":"test error"}`},
		{"/users/4", "GET", "", nil, http.StatusNotFound, `{"message":"user not found"}`},
		{"/users/ahh", "GET", "", nil, http.StatusBadRequest, `{"message":"unable to parse parameter id"}`},
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
		{"/users", "PUT", `{"email":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lastName":"Smith"}`, map[string]string{"Content-Type": "application/json"}, http.StatusOK, "This is a test!"},
		{"/users", "PUT", `{"whatisthis":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lastName":"Smith"}`, map[string]string{"Content-Type": "application/json"}, http.StatusOK, "This is a test!"},
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
		{"/users/1", "GET", "", nil, http.StatusOK, `{"id":1,"password":"password","email":"jacob.smith@unb.ca","firstName":"Jacob","lastName":"Smith","archived":false}`},
		{"/users/2", "GET", "", nil, http.StatusOK, `{"id":2,"password":"password","email":"jacob.smith@unb.ca","firstName":"Jacob","lastName":"Smith","archived":false}`},
		{"/users/3", "GET", "", nil, http.StatusInternalServerError, `{"message":"test error"}`},
		{"/users/4", "GET", "", nil, http.StatusNotFound, `{"message":"user not found"}`},
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
		{"/users", "GET", "", nil, http.StatusOK, `[{"id":1,"password":"password","email":"jacob.smith@unb.ca","firstName":"Jacob","lastName":"Smith","archived":false}]`},
		{"/users", "GET", "", nil, http.StatusInternalServerError, `{"message":"test error"}`},
		{"/users", "GET", "", nil, http.StatusOK, "[]"},
	}

	r := chi.NewRouter()
	r.Get("/users", UApi.Users)

	runTests(t, r, usersTests)
}

func TestCreateUser(t *testing.T) {
	initUserAPI()

	createUserTests := []*test{
		{"/users", "PUT", `{"email":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lastName":"Smith"}`, map[string]string{"Content-Type": "application/json"}, http.StatusCreated, ""},
		{"/users", "PUT", `{"email":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lastName":"Smith"}`, map[string]string{"Content-Type": "application/json"}, http.StatusInternalServerError, `{"message":"test error"}`},
		{"/users", "PUT", `{"eml":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lastName":"Smith"}`, map[string]string{"Content-Type": "application/json"}, http.StatusBadRequest, `{"message":"a user must have an email"}`},
		{"/users", "PUT", `{"email":"jacob.smith@unb.ca","passrd":"password","firstName":"Jacob","lastName":"Smith"}`, map[string]string{"Content-Type": "application/json"}, http.StatusBadRequest, `{"message":"a user must have a password"}`},
		{"/users", "PUT", `{"email":"jacob.smith@unb.ca","password":"password","fitName":"Jacob","lastName":"Smith"}`, map[string]string{"Content-Type": "application/json"}, http.StatusBadRequest, `{"message":"a user must have a first name"}`},
		{"/users", "PUT", `{"email":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lasame":"Smith"}`, map[string]string{"Content-Type": "application/json"}, http.StatusBadRequest, `{"message":"a user must have a last name"}`},
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

	updateUserTests := []*test{
		{"/users/1", "POST", `{"ID":1,"email":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lastName":"Smith","Archived":false}`, map[string]string{"Content-Type": "application/json"}, http.StatusOK, ""},
		{"/users/1", "POST", `{"ID":"1","email":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lastName":"Smith","Archived":false}`, map[string]string{"Content-Type": "application/json"}, http.StatusBadRequest, `{"message":"json: cannot unmarshal string into Go struct field UserRequest.id of type int"}`},
		{"/users/1", "POST", `{"ID":1,"email":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lastName":"Smith","Archived":"false"}`, map[string]string{"Content-Type": "application/json"}, http.StatusBadRequest, `{"message":"json: cannot unmarshal string into Go struct field UserRequest.archived of type bool"}`},
	}

	r := chi.NewRouter()
	r.Route("/users/{userId}", func(r chi.Router) {
		r.Use(UApi.UserCtx)
		r.Post("/", UApi.UpdateUser)
	})

	runTests(t, r, updateUserTests)
}
