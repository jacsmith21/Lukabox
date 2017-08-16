package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/jacsmith21/lukabox/domain"
	log "github.com/jacsmith21/lukabox/ext/logrus"
)

//AuthenticationAPI the services used
type AuthenticationAPI struct {
	AuthenticationService domain.AuthenticationService
	UserService           domain.UserService
}

//CredentialsRequest a request with credentials
type CredentialsRequest struct {
	*domain.Credentials
}

//Bind post-processing after decode
func (c *CredentialsRequest) Bind(r *http.Request) error {
	return nil
}

//SignUp signup handler
func SignUp(w http.ResponseWriter, r *http.Request) {

}

//Token jwt token
type Token struct {
	Token string `json:"token"`
}

//TokenResponse a token response
type TokenResponse struct {
	*Token
}

//Render pre-processing before marshelling
func (tr *TokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

//NewTokenResponse creates a new token response
func NewTokenResponse(token *Token) *TokenResponse {
	resp := &TokenResponse{Token: token}
	return resp
}

//Login login handler
func (aa *AuthenticationAPI) Login(w http.ResponseWriter, r *http.Request) {
	log.WithField("method", "Login").Info("starting")
	var authenticated bool
	var err error

	c := &CredentialsRequest{}
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
	token := &Token{tokenString}

	if err := render.Render(w, r, NewTokenResponse(token)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}
