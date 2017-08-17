package structure

import (
	"net/http"

	"github.com/jacsmith21/lukabox/domain"
)

//CredentialsRequest a request with credentials
type CredentialsRequest struct {
	*domain.Credentials
}

//Bind post-processing after decode
func (c *CredentialsRequest) Bind(r *http.Request) error {
	return nil
}

//TokenResponse a token response
type TokenResponse struct {
	*domain.Token
}

//Render pre-processing before marshelling
func (tr *TokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

//NewTokenResponse creates a new token response
func NewTokenResponse(token *domain.Token) *TokenResponse {
	resp := &TokenResponse{Token: token}
	return resp
}
