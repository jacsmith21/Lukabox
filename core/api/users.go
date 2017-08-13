package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jacsmith21/lukabox/core/db"
	"github.com/jacsmith21/lukabox/domain"
)

// UserCtx is used to create a user context by id
func UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *domain.User
		var err error
		var id int

		userID := chi.URLParam(r, "id")
		if userID == "" {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		id, err = strconv.Atoi(userID)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		user, err = db.GetUser(id)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUser gets a user by id
func GetUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*domain.User)
	if err := render.Render(w, r, NewUserResponse(user)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// UserResponse for json
type UserResponse struct {
	*domain.User
}

// Render does pre-processing before a response is marshalled
func (rd *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
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

// Users lists the users using the RenderList function
func Users(w http.ResponseWriter, r *http.Request) {
	var users []*domain.User
	var err error

	if users, err = db.GetUsers(); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	err = render.RenderList(w, r, NewUserListResponse(users))
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

//UserRequest a reuqest to create a user
type UserRequest struct {
	*domain.User
}

// Bind post-processing after decode
func (u *UserRequest) Bind(r *http.Request) error {
	u.User.Archived = false
	return nil
}

//CreateUser creates a user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	data := &UserRequest{}
	err := render.Bind(r, data)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	user := data.User
	db.CreateUser(user)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewUserResponse(user))
}

// UpdateUser updates the user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*domain.User)

	data := &UserRequest{User: user}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user = data.User
	db.UpdateUser(user.ID, user)

	render.Render(w, r, NewUserResponse(user))
}
