package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
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
	implUserServiceMethods()
}

func implUserServiceMethods() {
	USvc.UserByIDFn = func(id int) (*domain.User, error) {
		if id == 1 {
			user := domain.User{ID: 1, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}
			return &user, nil
		} else if id == 2 {
			user := domain.User{ID: 2, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}
			return &user, nil
		} else if id == 3 {
			return nil, errors.New("test error")
		}
		return nil, nil
	}

	USvc.UserByEmailFn = func(email string) (*domain.User, error) {
		if email != "jacob.smith@unb.ca" {
			return nil, errors.New("expected email to be jacob.smith@unb.ca")
		}
		user := domain.User{ID: 1, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}
		return &user, nil
	}

	USvc.UsersFn = func() ([]*domain.User, error) {
		users := []*domain.User{
			{ID: 1, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false},
		}
		return users, nil
	}
	USvc.CreateUserFn = func(u1 *domain.User) error {
		u2 := domain.User{Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith"}
		if !reflect.DeepEqual(u1, u2) {
			var err error
			var u1m []byte
			var u2m []byte

			if u1m, err = json.Marshal(u1); err != nil {
				return errors.New("error marshaling given user")
			}

			if u2m, err = json.Marshal(u2); err != nil {
				return errors.New("error marshaling fake user")
			}

			return fmt.Errorf("expected user %v to be equal to %v", u2m, u1m)
		}
		return nil
	}
	USvc.UpdateUserFn = func(id int, u1 *domain.User) error {
		u2 := domain.User{ID: 1, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}
		if !reflect.DeepEqual(u1, u2) {
			var err error
			var u1m []byte
			var u2m []byte

			if u1m, err = json.Marshal(u1); err != nil {
				return errors.New("error marshaling given user")
			}

			if u2m, err = json.Marshal(u2); err != nil {
				return errors.New("error marshaling fake user")
			}

			return fmt.Errorf("expected user %v to be equal to %v", u2m, u1m)
		}

		if id != 1 {
			return errors.New("expected id to be 1")
		}

		return nil
	}
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
	}

	r := chi.NewRouter()
	r.Get("/users", UApi.Users)

	runTests(t, r, usersTests)
}

func TestCreateUser(t *testing.T) {
	initUserAPI()

	user := domain.User{Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith"}

	var m []byte
	var err error
	if m, err = json.Marshal(user); err != nil {
		t.Fatal("error marshaling test user")
	}

	req, err := http.NewRequest("PUT", "/users", bytes.NewReader(m))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		r.Use(UApi.UserRequestCtx)
		r.Put("/", UApi.CreateUser)
	})
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v.\nBody: %v", status, http.StatusOK, w.Body)
	}

	if !USvc.CreateUserInvoked {
		t.Fatal("expected CreateUser to be invoked")
	}
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

	if !USvc.UserByIDInvoked {
		t.Fatal("expected UsersByID to be invoked")
	}
}
