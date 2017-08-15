package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/jacsmith21/lukabox/domain"
	"github.com/jacsmith21/lukabox/mock"
)

func TestUserCtx(t *testing.T) {
	var us mock.UserService
	var ua UserAPI
	ua.UserService = &us

	us.UserByIDFn = func(id int) (*domain.User, error) {
		if id != 1 {
			t.Fatal("expected id to be 1")
		}
		user := domain.User{ID: 1, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false}
		return &user, nil
	}

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

}

func TestUsers(t *testing.T) {
	var us mock.UserService
	var ua UserAPI
	ua.UserService = &us

	us.UsersFn = func() ([]*domain.User, error) {
		users := []*domain.User{
			{ID: 1, Email: "jacob.smith@unb.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false},
			{ID: 2, Email: "j.a.smith@live.ca", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false},
			{ID: 3, Email: "jacobsmithunb@gmail.com", Password: "password", FirstName: "Jacob", LastName: "Smith", Archived: false},
		}
		return users, nil
	}

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
