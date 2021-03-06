package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-errors/errors"
	"github.com/jacsmith21/lukabox/domain"
	"github.com/jacsmith21/lukabox/ext/log"
	"github.com/jacsmith21/lukabox/ext/render"
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
			render.WithError(err).BadRequest(w, r)
			return
		}

		user := r.Context().Value("user").(*domain.User)

		log.WithField("user", user).Debug("user in validate")
		log.WithField("claims", claims).Debug("claims in validate")

		id := int(claims["id"].(float64))
		if user.ID != id {
			render.Unauthorized(w, r)
			return
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
			render.WithError(err).BadRequest(w, r)
			return
		}
		if !available {
			render.WithMessage("email unavailable").Conflict(w, r)
		}

		next.ServeHTTP(w, r)
	})
}

//Login login handler
func (a *AuthenticationAPI) Login(w http.ResponseWriter, r *http.Request) {
	log.WithField("method", "Login").Info("starting")

	c := &stc.CredentialsRequest{}
	if err := render.Bind(r, c); err != nil {
		render.WithError(err).BadRequest(w, r)
		return
	}

	log.WithField("Credentials", c.Credentials).Debug("credentials")

	authenticated, err := a.AuthenticationService.Authenticate(c.Credentials.Email, c.Credentials.Password)
	if err != nil {
		render.WithError(err).BadRequest(w, r)
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

	credentials := c.Credentials
	if credentials == nil {
		err := errors.New("credentials must be supplied")
		render.WithError(err).BadRequest(w, r)
	}

	user, err := a.UserService.UserByEmail(credentials.Email)
	if err != nil {
		render.WithError(err).BadRequest(w, r)
		return
	}

	claims := jwtauth.Claims{"id": user.ID}
	log.WithField("id", user.ID).Debug("adding id to claims")
	_, tokenString, _ := tokenAuth.Encode(claims)
	token := &domain.Token{Token: tokenString}

	if err := render.Instance(w, r, stc.NewTokenResponse(token)); err != nil {
		render.WithError(err).BadRequest(w, r)
		return
	}
}
