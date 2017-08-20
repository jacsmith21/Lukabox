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

var ASvc mock.AuthenticationService
var AApi AuthenticationAPI

func initAuthenticationAPI() {
	ASvc = mock.AuthenticationService{}
	AApi.AuthenticationService = &ASvc
	implAuthenticationServiceMethods()

	USvc = mock.UserService{}
	AApi.UserService = &USvc
	implUserServiceMethods()
}

func implAuthenticationServiceMethods() {
	ASvc.AuthenticateFn = func(email string, password string) (bool, error) {
		if email != "jacob.smith@unb.ca" {
			return false, errors.New("expected different email")
		}

		if password != "password" {
			return false, errors.New("expected different password")
		}

		return true, nil
	}
	ASvc.EmailAvailableFn = func(email string) (bool, error) {
		return email == "jacob.smith@unb.ca", nil
	}
}

func TestRequestValidator(t *testing.T) {
	initUserAPI()

	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	req.Header.Add("Authorization", "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.tjVEMiS5O2yNzclwLdaZ-FuzrhyqOT7UwM9Hfc0ZQ8Q")

	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Route("/users/{userId}", func(r chi.Router) {
		r.Use(UApi.UserCtx)
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(AApi.RequestValidator)
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
	initAuthenticationAPI()

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
		r.Use(AApi.SignUpValidator)
		r.Put("/", UApi.CreateUser)
	})
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v\nbody: %v", status, http.StatusCreated, w.Body)
	}

	if !ASvc.EmailAvailableInvoked {
		t.Fatal("expected Authenticate to be invoked")
	}
}

func TestLogin(t *testing.T) {
	initAuthenticationAPI()
	initUserAPI()

	cred := domain.Credentials{Email: "jacob.smith@unb.ca", Password: "password"}

	m, err := json.Marshal(cred)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/login", bytes.NewReader(m))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(AApi.Login)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v.\nBody: %v", status, http.StatusOK, w.Body)
	}

	if !ASvc.AuthenticateInvoked {
		t.Fatal("expected Authenticate to be invoked")
	}
}
