package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/jacsmith21/lukabox/domain"
	log "github.com/jacsmith21/lukabox/ext/logrus"
	"github.com/jacsmith21/lukabox/structure"
)

//AuthenticationAPI the services used
type AuthenticationAPI struct {
	AuthenticationService domain.AuthenticationService
	UserService           domain.UserService
}

//SignUp signup handler
func SignUp(w http.ResponseWriter, r *http.Request) {

}

//Login login handler
func (aa *AuthenticationAPI) Login(w http.ResponseWriter, r *http.Request) {
	log.WithField("method", "Login").Info("starting")
	var authenticated bool
	var err error

	c := &structure.CredentialsRequest{}
	if err = render.Bind(r, c); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	log.WithField("Credentials", c.Credentials).Debug("credentials")

	authenticated, err = aa.AuthenticationService.Authenticate(c.Credentials.Email, c.Credentials.Password)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	log.WithField("authenticated", authenticated).Debug("authentication complete")
	if !authenticated {
		w.WriteHeader(http.StatusForbidden)
		fmt.Println("Error logging in")
		fmt.Fprint(w, "Invalid credentials")
		return
	}

	log.Debug("Creating Token")
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	user, _ := aa.UserService.UserByEmail(c.Credentials.Email)
	claims := jwtauth.Claims{"id": user.ID}
	log.WithField("id", user.ID).Debug("adding id to claims")
	_, tokenString, _ := tokenAuth.Encode(claims)
	token := &domain.Token{tokenString}

	if err := render.Render(w, r, structure.NewTokenResponse(token)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}
