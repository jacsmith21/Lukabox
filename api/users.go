package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jacsmith21/lukabox/domain"
	log "github.com/jacsmith21/lukabox/ext/logrus"
	"github.com/jacsmith21/lukabox/structure"
)

//UserAPI the services used
type UserAPI struct {
	UserService domain.UserService
}

// UserCtx is used to create a user context by id
func (a *UserAPI) UserCtx(next http.Handler) http.Handler {
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

		user, err = a.UserService.UserByID(id)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserByID gets a user by id
func (a *UserAPI) UserByID(w http.ResponseWriter, r *http.Request) {
	log.WithField("method", "UserByID").Info("starting")
	user := r.Context().Value("user").(*domain.User)
	if err := render.Render(w, r, structure.NewUserResponse(user)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// Users lists the users using the RenderList function
func (a *UserAPI) Users(w http.ResponseWriter, r *http.Request) {
	var users []*domain.User
	var err error

	if users, err = a.UserService.Users(); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	err = render.RenderList(w, r, structure.NewUserListResponse(users))
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

//CreateUser creates a user
func (a *UserAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
	data := &structure.UserRequest{}
	err := render.Bind(r, data)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user := data.User
	a.UserService.CreateUser(user)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, structure.NewUserResponse(user)) //TODO change this
}

// UpdateUser updates the user
func (a *UserAPI) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*domain.User)

	data := &structure.UserRequest{User: user}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user = data.User
	a.UserService.UpdateUser(user.ID, user)

	render.Render(w, r, structure.NewUserResponse(user))
}
