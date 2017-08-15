package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jacsmith21/lukabox/domain"
	"github.com/jacsmith21/lukabox/mock"
)

func TestListUsers(t *testing.T) {
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

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ua.Users)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if !us.UsersInvoked {
		t.Fatal("expected Users to be invoked")
	}
}
