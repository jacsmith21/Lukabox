package stc

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/jacsmith21/lukabox/domain"
)

// UserResponse for json
type UserResponse struct {
	*domain.User
}

// Render does pre-processing before a response is marshalled
func (rd *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

//UserRequest a user request
type UserRequest struct {
	*domain.User
}

// Bind post-processing after decode
func (u *UserRequest) Bind(r *http.Request) error {
	u.User.Archived = false
	return nil
}

// NewUserResponse ceates a new user reponse
func NewUserResponse(user *domain.User) *UserResponse {
	resp := &UserResponse{User: user}
	return resp
}

// NewUserListResponse creates a new renderer list of reponses
func NewUserListResponse(users []*domain.User) []render.Renderer {
	list := []render.Renderer{}
	for _, user := range users {
		list = append(list, NewUserResponse(user))
	}
	return list
}
