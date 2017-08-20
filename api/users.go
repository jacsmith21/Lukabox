package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jacsmith21/lukabox/domain"
	log "github.com/jacsmith21/lukabox/ext/logrus"
	"github.com/jacsmith21/lukabox/stc"
)

//UserAPI the services used
type UserAPI struct {
	UserService domain.UserService
}

// UserCtx is used to create a user context by id
func (a *UserAPI) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithField("method", "UserCtx").Info("starting")

		userID := chi.URLParam(r, "id")
		if userID == "" {
			render.Render(w, r, ErrInvalidRequest(errors.New("paramter id should not be empty")))
			return
		}
		log.WithField("id", userID).Debug("user id from paramter")

		id, err := strconv.Atoi(userID)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		user, err := a.UserService.UserByID(id)
		if err != nil {
			render.Render(w, r, ErrNotFound(err))
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserRequestCtx a user request context generator
func (a *UserAPI) UserRequestCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := &stc.UserRequest{}

		err := render.Bind(r, data)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		user := data.User
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserByID gets a user by id
func (a *UserAPI) UserByID(w http.ResponseWriter, r *http.Request) {
	log.WithField("method", "UserByID").Info("starting")
	user := r.Context().Value("user").(*domain.User)
	if err := render.Render(w, r, stc.NewUserResponse(user)); err != nil {
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

	err = render.RenderList(w, r, stc.NewUserListResponse(users))
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

//CreateUser creates a user
func (a *UserAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*domain.User)

	a.UserService.CreateUser(user)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, stc.NewUserResponse(user)) //TODO change this
}

// UpdateUser updates the user
func (a *UserAPI) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*domain.User)

	data := &stc.UserRequest{User: user}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user = data.User
	a.UserService.UpdateUser(user.ID, user)

	render.Render(w, r, stc.NewUserResponse(user))
}
