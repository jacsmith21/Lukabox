package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/jacsmith21/lukabox/core/db"
)

//UserCredentials a reguler user credentials
type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Bind post-processing after decode
func (uc *UserCredentials) Bind(r *http.Request) error {
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
func Login(w http.ResponseWriter, r *http.Request) {
	credentials := &UserCredentials{}
	if err := render.Bind(r, credentials); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	authenticated := db.AuthenticateUser(credentials.Email, credentials.Password)
	if !authenticated {
		w.WriteHeader(http.StatusForbidden)
		fmt.Println("Error logging in")
		fmt.Fprint(w, "Invalid credentials")
		return
	}
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	user, _ := db.GetUserByEmail(credentials.Email)
	claims := jwtauth.Claims{"id": user.ID}
	_, tokenString, _ := tokenAuth.Encode(claims)
	token := &Token{tokenString}
	if err := render.Render(w, r, NewTokenResponse(token)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}
