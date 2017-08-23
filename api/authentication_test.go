package api

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
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
	initAuthenticationAPI()

	tests := []*test{
		{"/users/1", "GET", "", map[string]string{"Authorization": "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.tjVEMiS5O2yNzclwLdaZ-FuzrhyqOT7UwM9Hfc0ZQ8Q"}, http.StatusOK, "This is a test!"},
	}

	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	r := chi.NewRouter()
	r.Route("/users/{userId}", func(r chi.Router) {
		r.Use(UApi.UserCtx)
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(AApi.RequestValidator)
		r.Get("/", func(w http.ResponseWriter, request *http.Request) {
			w.Write([]byte("This is a test!"))
		})
	})

	runTests(t, r, tests)
}

func TestSignUpValidator(t *testing.T) {
	initUserAPI()
	initAuthenticationAPI()

	tests := []*test{
		{"/users", "PUT", `{"email":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lastName":"Smith"}`, map[string]string{"Content-Type": "application/json"}, http.StatusOK, "This is a test!"},
	}

	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		r.Use(UApi.UserRequestCtx)
		r.Use(AApi.SignUpValidator)
		r.Put("/", UApi.CreateUser)
	})

	runTests(t, r, tests)
}

func TestLogin(t *testing.T) {
	initAuthenticationAPI()
	initUserAPI()

	tests := []*test{
		{"/login", "POST", `{"email":"jacob.smith@unb.ca","password":"password"}`, map[string]string{"Content-Type": "application/json"}, http.StatusOK, "This is a test!"},
	}

	r := chi.NewRouter()
	r.Post("/login", AApi.Login)

	runTests(t, r, tests)
}
