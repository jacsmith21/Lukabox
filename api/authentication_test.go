package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/jacsmith21/lukabox/domain"
	"github.com/jacsmith21/lukabox/mock"
)

func implAuthenticationServiceMehods(as *mock.AuthenticationService) {
	as.AuthenticateFn = func(email string, password string) (bool, error) {
		if email != "jacob.smith@unb.ca" {
			return false, errors.New("expected different email")
		}

		if password != "password" {
			return false, errors.New("expected different password")
		}

		return true, nil
	}
	as.EmailAvailableFn = func(email string) (bool, error) {
		return email == "jacob.smith@unb.ca", nil
	}
}

func TestRequestValidator(t *testing.T) {
	var us mock.UserService
	var aa AuthenticationAPI
	var ua UserAPI
	ua.UserService = &us
	implUserServiceMethods(&us)

	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	req.Header.Add("Authorization", "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.tjVEMiS5O2yNzclwLdaZ-FuzrhyqOT7UwM9Hfc0ZQ8Q")

	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Route("/users/{id}", func(r chi.Router) {
		r.Use(ua.UserCtx)
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(aa.RequestValidator)
		r.Get("/", func(w http.ResponseWriter, request *http.Request) {
			w.Write([]byte("This is a test!"))
		})
	})
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestSignUpValidator(t *testing.T) {
	var us mock.UserService
	var as mock.AuthenticationService
	var aa AuthenticationAPI
	var ua UserAPI
	ua.UserService = &us
	aa.AuthenticationService = &as
	implUserServiceMethods(&us)
	implAuthenticationServiceMehods(&as)

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
		r.Use(aa.SignUpValidator)
		r.Put("/", ua.CreateUser)
	})
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v\nbody: %v", status, http.StatusCreated, w.Body)
	}

	if !as.EmailAvailableInvoked {
		t.Fatal("expected Authenticate to be invoked")
	}
}

func TestLogin(t *testing.T) {
	var as mock.AuthenticationService
	var us mock.UserService
	var aa AuthenticationAPI
	aa.AuthenticationService = &as
	aa.UserService = &us
	implAuthenticationServiceMehods(&as)
	implUserServiceMethods(&us)

	cred := domain.Credentials{Email: "jacob.smith@unb.ca", Password: "password"}

	var m []byte
	var err error
	if m, err = json.Marshal(cred); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/login", bytes.NewReader(m))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(aa.Login)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v.\nBody: %v", status, http.StatusOK, w.Body)
	}

	if !as.AuthenticateInvoked {
		t.Fatal("expected Authenticate to be invoked")
	}
}
