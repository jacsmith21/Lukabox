package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/jacsmith21/lukabox/domain"
	"github.com/jacsmith21/lukabox/mock"
)

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
	uApi := UserAPI{}
	uSvc := mock.UserService{}
	uApi.UserService = &uSvc

	userCtxTests := []*test{
		{"/users/1", "GET", "", nil, http.StatusOK, "This is a test!"},
		{"/users/3", "GET", "", nil, http.StatusInternalServerError, `{"message":"test error"}`},
		{"/users/4", "GET", "", nil, http.StatusNotFound, `{"message":"user not found"}`},
		{"/users/ahh", "GET", "", nil, http.StatusBadRequest, `{"message":"unable to parse parameter id"}`},
	}

	uSvc.UserByIDFn = func(id int) (*domain.User, error) {
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

	r := chi.NewRouter()
	r.Route("/users/{userId}", func(r chi.Router) {
		r.Use(uApi.UserCtx)
		r.Get("/", func(w http.ResponseWriter, request *http.Request) {
			w.Write([]byte("This is a test!"))
		})
	})

	runTests(t, r, userCtxTests)
}

func TestUserRequestCtx(t *testing.T) {
	uApi := UserAPI{}
	uSvc := mock.UserService{}
	uApi.UserService = &uSvc

	userRequestCtxTests := []*test{
		{"/users", "PUT", `{"email":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lastName":"Smith"}`, map[string]string{"Content-Type": "application/json"}, http.StatusOK, "This is a test!"},
		{"/users", "PUT", `{"whatisthis":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lastName":"Smith"}`, map[string]string{"Content-Type": "application/json"}, http.StatusOK, "This is a test!"},
	}

	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		r.Use(uApi.UserRequestCtx)
		r.Put("/", func(w http.ResponseWriter, request *http.Request) {
			w.Write([]byte("This is a test!"))
		})
	})

	runTests(t, r, userRequestCtxTests)
}

func TestUserByID(t *testing.T) {
	uApi := UserAPI{}
	uSvc := mock.UserService{}
	uApi.UserService = &uSvc

	userByIDTests := []*test{
		{"/users/1", "GET", "", nil, http.StatusOK, `{"id":1,"password":"password","email":"jacob.smith@unb.ca","firstName":"Jacob","lastName":"Smith","archived":false}`},
		{"/users/2", "GET", "", nil, http.StatusOK, `{"id":2,"password":"password","email":"jacob.smith@unb.ca","firstName":"Jacob","lastName":"Smith","archived":false}`},
		{"/users/3", "GET", "", nil, http.StatusInternalServerError, `{"message":"test error"}`},
		{"/users/4", "GET", "", nil, http.StatusNotFound, `{"message":"user not found"}`},
	}

	uSvc.UserByIDFn = func(id int) (*domain.User, error) {
		return &domain.User{ID: id, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}, nil
	}

	r := chi.NewRouter()
	r.Route("/users/{userId}", func(r chi.Router) {
		r.Use(uApi.UserCtx)
		r.Get("/", uApi.UserByID)
	})

	runTests(t, r, userByIDTests)
}

func TestUsers(t *testing.T) {
	uApi := UserAPI{}
	uSvc := mock.UserService{}
	uApi.UserService = &uSvc

	tests := []*test{
		{"/users", "GET", "", nil, http.StatusOK, `[{"id":1,"password":"password","email":"jacob.smith@unb.ca","firstName":"Jacob","lastName":"Smith","archived":false}]`},
		{"/users", "GET", "", nil, http.StatusInternalServerError, `{"message":"test error"}`},
		{"/users", "GET", "", nil, http.StatusOK, "[]"},
	}

	count := 0
	uSvc.UsersFn = func() ([]*domain.User, error) {
		count++
		if count == 1 {
			users := []*domain.User{
				{ID: 1, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false},
			}
			return users, nil
		} else if count == 2 {
			return nil, nil
		}
		return nil, errors.New("test error")
	}

	r := chi.NewRouter()
	r.Get("/users", uApi.Users)

	runTests(t, r, tests)
}

func TestCreateUser(t *testing.T) {
	uApi := UserAPI{}
	uSvc := mock.UserService{}
	uApi.UserService = &uSvc

	createUserTests := []*test{
		{"/users", "PUT", `{"email":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lastName":"Smith"}`, map[string]string{"Content-Type": "application/json"}, http.StatusCreated, ""},
		{"/users", "PUT", `{"email":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lastName":"Smith"}`, map[string]string{"Content-Type": "application/json"}, http.StatusInternalServerError, `{"message":"test error"}`},
	}

	count := 0
	uSvc.InsertUserFn = func(user *domain.User) error {
		count++
		if count == 1 {
			return nil
		}
		return errors.New("test error")
	}

	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		r.Use(uApi.UserRequestCtx)
		r.Put("/", uApi.CreateUser)
	})

	runTests(t, r, createUserTests)
}

func TestUpdateUser(t *testing.T) {
	uApi := UserAPI{}
	uSvc := mock.UserService{}
	uApi.UserService = &uSvc

	updateUserTests := []*test{
		{"/users/1", "POST", `{"ID":1,"email":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lastName":"Smith","Archived":false}`, map[string]string{"Content-Type": "application/json"}, http.StatusOK, ""},
		{"/users/1", "POST", `{"ID":"1","email":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lastName":"Smith","Archived":false}`, map[string]string{"Content-Type": "application/json"}, http.StatusBadRequest, `{"message":"json: cannot unmarshal string into Go struct field UserRequest.id of type int"}`},
		{"/users/1", "POST", `{"ID":1,"email":"jacob.smith@unb.ca","password":"password","firstName":"Jacob","lastName":"Smith","Archived":"false"}`, map[string]string{"Content-Type": "application/json"}, http.StatusBadRequest, `{"message":"json: cannot unmarshal string into Go struct field UserRequest.archived of type bool"}`},
	}

	uSvc.UserByIDFn = func(id int) (*domain.User, error) {
		return &domain.User{ID: id, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}, nil
	}

	uSvc.UpdateUserFn = func(id int, user *domain.User) error {
		if id != 1 {
			return errors.New("expected id to be 1")
		}
		return nil
	}

	r := chi.NewRouter()
	r.Route("/users/{userId}", func(r chi.Router) {
		r.Use(uApi.UserCtx)
		r.Post("/", uApi.UpdateUser)
	})

	runTests(t, r, updateUserTests)
}
