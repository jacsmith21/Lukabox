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
	aAPI := AuthenticationAPI{}
	aSvc := mock.AuthenticationService{}
	uSvc := mock.UserService{}
	aAPI.AuthenticationService = &aSvc
	aAPI.UserService = &uSvc

	uAPI := UserAPI{}
	uAPI.UserService = &uSvc

	tests := []*test{
		{"/users/1", "GET", "", map[string]string{"Authorization": "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.tjVEMiS5O2yNzclwLdaZ-FuzrhyqOT7UwM9Hfc0ZQ8Q"}, http.StatusOK, "This is a test!"},
		{"/users/1", "GET", "", map[string]string{"Authorization": "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Mn0.-ScBrpAXat0bA0Q-kJnL7xnst1-dd_SsIzseTUPT2wE"}, http.StatusUnauthorized, `{"message":"Unauthorized"}`},
		{"/users/1", "GET", "", map[string]string{"Authorization": "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.tjVFMiS5O2yNzclwLdaZ-FuzrhyqOT7UwM9Hfc0ZQ8Q"}, http.StatusBadRequest, `{"message":"signature is invalid"}`},
		{"/users/1", "GET", "", map[string]string{"Authorization": "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0OT7UwM9Hfc0ZQ8Q"}, http.StatusBadRequest, `{"message":"token contains an invalid number of segments"}`},
	}

	uSvc.UserByIDFn = func(id int) (*domain.User, error) {
		return &domain.User{ID: id, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}, nil
	}

	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	r := chi.NewRouter()
	r.Route("/users/{userId}", func(r chi.Router) {
		r.Use(uAPI.UserCtx)
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(aAPI.RequestValidator)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("This is a test!"))
		})
	})

	runTests(t, r, tests)
}

func TestSignUpValidator(t *testing.T) {
	aAPI := AuthenticationAPI{}
	aSvc := mock.AuthenticationService{}
	uSvc := mock.UserService{}
	aAPI.AuthenticationService = &aSvc
	aAPI.UserService = &uSvc

	uAPI := UserAPI{}
	uAPI.UserService = &uSvc

	tests := []*test{
		{"/users", "GET", `{"email":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lastName":"Smith"}`, map[string]string{"Content-Type": "application/json"}, http.StatusOK, "This is a test!"},
		{"/users", "GET", `{"email":"j.a.smith@live.ca","password":"password","firstName":"Jacob","lastName":"Smith"}`, map[string]string{"Content-Type": "application/json"}, http.StatusOK, "This is a test!"},
	}

	uSvc.UserByEmailFn = func(email string) (*domain.User, error) {
		return &domain.User{ID: 1, Email: email, Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}, nil
	}

	aSvc.EmailAvailableFn = func(email string) (bool, error) {
		return true, nil
	}

	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		r.Use(uAPI.UserRequestCtx)
		r.Use(aAPI.SignUpValidator)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("This is a test!"))
		})
	})

	runTests(t, r, tests)
}

func TestLogin(t *testing.T) {
	aAPI := AuthenticationAPI{}
	aSvc := mock.AuthenticationService{}
	uSvc := mock.UserService{}
	aAPI.AuthenticationService = &aSvc
	aAPI.UserService = &uSvc

	tests := []*test{
		{"/login", "POST", `{"email":"jacob.smith@unb.ca","password":"password"}`, map[string]string{"Content-Type": "application/json"}, http.StatusOK, `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.tjVEMiS5O2yNzclwLdaZ-FuzrhyqOT7UwM9Hfc0ZQ8Q"}`},
		{"/login", "POST", `{"email":"j.a.smith@live.ca","password":"password"}`, map[string]string{"Content-Type": "application/json"}, http.StatusOK, `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Mn0.-ScBrpAXat0bA0Q-kJnL7xnst1-dd_SsIzseTUPT2wE"}`},
	}

	aSvc.AuthenticateFn = func(email string, password string) (bool, error) {
		return true, nil
	}

	count := 0
	uSvc.UserByEmailFn = func(email string) (*domain.User, error) {
		count++
		return &domain.User{ID: count, Email: email, Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}, nil
	}

	r := chi.NewRouter()
	r.Post("/login", aAPI.Login)

	runTests(t, r, tests)
}
