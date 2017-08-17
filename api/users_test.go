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

func implUserServiceMethods(us *mock.UserService) {
	us.UserByIDFn = func(id int) (*domain.User, error) {
		if id != 1 {
			return nil, errors.New("expected id to be 1")
		}
		user := domain.User{ID: 1, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}
		return &user, nil
	}

	us.UserByEmailFn = func(email string) (*domain.User, error) {
		if email != "jacob.smith@unb.ca" {
			return nil, errors.New("expected email to be jacob.smith@unb.ca")
		}
		user := domain.User{ID: 1, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}
		return &user, nil
	}

	us.UsersFn = func() ([]*domain.User, error) {
		users := []*domain.User{
			{ID: 1, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false},
			{ID: 2, Email: "j.a.smith@live.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false},
			{ID: 3, Email: "jacobsmithunb@gmail.com", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false},
		}
		return users, nil
	}
	us.CreateUserFn = func(u1 *domain.User) error {
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

	us.UpdateUserFn = func(id int, u1 *domain.User) error {
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

func TestUserCtx(t *testing.T) {
	var us mock.UserService
	var ua UserAPI
	ua.UserService = &us
	implUserServiceMethods(&us)

	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Route("/users/{id}", func(r chi.Router) {
		r.Use(ua.UserCtx)
		r.Get("/", func(w http.ResponseWriter, request *http.Request) {
			w.Write([]byte("This is a test!"))
		})
	})
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !us.UserByIDInvoked {
		t.Fatal("expected UsersByID to be invoked")
	}
}

func TestUserRequestCtx(t *testing.T) {
	var us mock.UserService
	var ua UserAPI
	ua.UserService = &us
	implUserServiceMethods(&us)

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
		r.Use(ua.UserRequestCtx)
		r.Put("/", func(w http.ResponseWriter, request *http.Request) {
			w.Write([]byte("This is a test!"))
		})
	})
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v\nbody: %v", status, http.StatusOK, w.Body)
	}
}

func TestUsersByID(t *testing.T) {
	var us mock.UserService
	var ua UserAPI
	ua.UserService = &us
	implUserServiceMethods(&us)

	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Route("/users/{id}", func(r chi.Router) {
		r.Use(ua.UserCtx)
		r.Get("/", ua.UserByID)
	})
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "{\"id\":1,\"password\":\"password\",\"email\":\"jacob.smith@unb.ca\",\"firstName\":\"Jacob\",\"lastName\":\"Smith\",\"archived\":false}"
	if body := w.Body.String(); strings.Trim(body, "\n") != expected {
		t.Errorf("expected body to be: \n%s\nnot\n%s", expected, body)
	}

	if !us.UserByIDInvoked {
		t.Fatal("expected UserByID to be invoked")
	}
}

func TestUsers(t *testing.T) {
	var us mock.UserService
	var ua UserAPI
	ua.UserService = &us
	implUserServiceMethods(&us)

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(ua.Users)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !us.UsersInvoked {
		t.Fatal("expected Users to be invoked")
	}
}

func TestCreateUser(t *testing.T) {
	var us mock.UserService
	var ua UserAPI
	ua.UserService = &us
	implUserServiceMethods(&us)

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
		r.Use(ua.UserRequestCtx)
		r.Put("/", ua.CreateUser)
	})
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v.\nBody: %v", status, http.StatusOK, w.Body)
	}

	if !us.CreateUserInvoked {
		t.Fatal("expected CreateUser to be invoked")
	}
}

func TestUpdateUser(t *testing.T) {
	var us mock.UserService
	var ua UserAPI
	ua.UserService = &us
	implUserServiceMethods(&us)

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
	r.Route("/users/{id}", func(r chi.Router) {
		r.Use(ua.UserCtx)
		r.Post("/", ua.UpdateUser)
	})
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !us.UserByIDInvoked {
		t.Fatal("expected UsersByID to be invoked")
	}
}
