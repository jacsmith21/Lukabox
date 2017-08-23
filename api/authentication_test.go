package api

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/jacsmith21/lukabox/domain"
	"github.com/jacsmith21/lukabox/mock"
)

func TestRequestValidator(t *testing.T) {
	aApi := AuthenticationAPI{}
	aSvc := mock.AuthenticationService{}
	uSvc := mock.UserService{}
	aApi.AuthenticationService = &aSvc
	aApi.UserService = &uSvc

	uApi := UserAPI{}
	uApi.UserService = &uSvc

	tests := []*test{
		{"/users/1", "GET", "", map[string]string{"Authorization": "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.tjVEMiS5O2yNzclwLdaZ-FuzrhyqOT7UwM9Hfc0ZQ8Q"}, http.StatusOK, "This is a test!"},
	}

	uSvc.UserByIDFn = func(id int) (*domain.User, error) {
		return &domain.User{ID: id, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}, nil
	}

	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	r := chi.NewRouter()
	r.Route("/users/{userId}", func(r chi.Router) {
		r.Use(uApi.UserCtx)
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(aApi.RequestValidator)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("This is a test!"))
		})
	})

	runTests(t, r, tests)
}

func TestSignUpValidator(t *testing.T) {
	aApi := AuthenticationAPI{}
	aSvc := mock.AuthenticationService{}
	uSvc := mock.UserService{}
	aApi.AuthenticationService = &aSvc
	aApi.UserService = &uSvc

	uApi := UserAPI{}
	uApi.UserService = &uSvc

	tests := []*test{
		{"/users", "GET", `{"email":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lastName":"Smith"}`, map[string]string{"Content-Type": "application/json"}, http.StatusOK, "This is a test!"},
	}

	uSvc.UserByEmailFn = func(email string) (*domain.User, error) {
		return &domain.User{ID: 1, Email: email, Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}, nil
	}

	aSvc.EmailAvailableFn = func(email string) (bool, error) {
		return true, nil
	}

	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		r.Use(uApi.UserRequestCtx)
		r.Use(aApi.SignUpValidator)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("This is a test!"))
		})
	})

	runTests(t, r, tests)
}

func TestLogin(t *testing.T) {
	aApi := AuthenticationAPI{}
	aSvc := mock.AuthenticationService{}
	uSvc := mock.UserService{}
	aApi.AuthenticationService = &aSvc
	aApi.UserService = &uSvc

	tests := []*test{
		{"/login", "POST", `{"email":"jacob.smith@unb.ca","password":"password"}`, map[string]string{"Content-Type": "application/json"}, http.StatusOK, "This is a test!"},
	}

	aSvc.AuthenticateFn = func(email string, password string) (bool, error) {
		return true, nil
	}

	r := chi.NewRouter()
	r.Post("/login", aApi.Login)

	runTests(t, r, tests)
}
