package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

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
		t.Fatal("expected CreateUser to be invoked")
	}
}
