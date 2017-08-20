package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/go-errors/errors"
	"github.com/jacsmith21/lukabox/domain"
	log "github.com/jacsmith21/lukabox/ext/logrus"
	"github.com/jacsmith21/lukabox/stc"
)

//AuthenticationAPI the services used
type AuthenticationAPI struct {
	AuthenticationService domain.AuthenticationService
	UserService           domain.UserService
}

// RequestValidator validates the request ie. checks whether the user is allowed to make this request
func (a *AuthenticationAPI) RequestValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		user := r.Context().Value("user").(*domain.User)

		log.WithField("user", user).Debug("user in validate")
		log.WithField("claims", claims).Debug("claims in validate")

		id := int(claims["id"].(float64))
		if user.ID != id {
			render.Render(w, r, ErrUnauthorized)
		}

		next.ServeHTTP(w, r)
	})
}

//SignUpValidator signup handler
func (a *AuthenticationAPI) SignUpValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*domain.User)
		email := user.Email

		available, err := a.AuthenticationService.EmailAvailable(email)
		if err != nil {
			render.Render(w, r, ErrBadRequest(errors.New("email taken")))
			return
		}
		if !available {
			fmt.Fprint(w, a)
		}

		next.ServeHTTP(w, r)
	})
}

//Login login handler
func (a *AuthenticationAPI) Login(w http.ResponseWriter, r *http.Request) {
	log.WithField("method", "Login").Info("starting")
	var authenticated bool
	var err error

	c := &stc.CredentialsRequest{}
	if err = render.Bind(r, c); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	log.WithField("Credentials", c.Credentials).Debug("credentials")

	authenticated, err = a.AuthenticationService.Authenticate(c.Credentials.Email, c.Credentials.Password)
	if err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	log.WithField("authenticated", authenticated).Debug("authentication complete")
	if !authenticated {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Invalid credentials")
		return
	}

	log.Debug("Creating Token")
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	user, _ := a.UserService.UserByEmail(c.Credentials.Email)
	claims := jwtauth.Claims{"id": user.ID}
	log.WithField("id", user.ID).Debug("adding id to claims")
	_, tokenString, _ := tokenAuth.Encode(claims)
	token := &domain.Token{Token: tokenString}

	if err := render.Render(w, r, stc.NewTokenResponse(token)); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}
}
