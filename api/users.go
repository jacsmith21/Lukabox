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

		userID := chi.URLParam(r, "userId")
		log.WithField("id", userID).Debug("user id from paramter")

		id, err := strconv.Atoi(userID)
		if err != nil {
			log.WithError(err).Debugf("unable to parse %s", userID)
			render.Render(w, r, ErrBadRequest(errors.New("unable to parse parameter id")))
			return
		}

		user, err := a.UserService.UserByID(id)
		if err != nil {
			log.WithError(err).Errorf("error fetching user with id %d", id)
			render.Render(w, r, ErrInternalServerError(err))
			return
		}
		if user == nil {
			log.Debugf("no user found with id %d", id)
			render.Render(w, r, ErrNotFound(errors.New("user not found")))
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserRequestCtx a user request context generator
func (a *UserAPI) UserRequestCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRequest := &stc.UserRequest{}

		err := render.Bind(r, userRequest)
		if err != nil {
			log.WithError(err).Error("error binding user request")
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		user := userRequest.User
		log.WithField("user", user).Debug("user from user request")

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserByID gets a user by id
func (a *UserAPI) UserByID(w http.ResponseWriter, r *http.Request) {
	log.WithField("method", "UserByID").Info("starting")
	user := r.Context().Value("user").(*domain.User)
	if err := render.Render(w, r, stc.NewUserResponse(user)); err != nil {
		log.WithError(err).Error("unable to render user response")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}
}

// Users lists the users using the RenderList function
func (a *UserAPI) Users(w http.ResponseWriter, r *http.Request) {
	users, err := a.UserService.Users()
	if err != nil {
		log.WithError(err).Error("error fetching users")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if err := render.RenderList(w, r, stc.NewUserListResponse(users)); err != nil {
		log.WithError(err).Error("error rendering user list response")
		render.Render(w, r, ErrInternalServerError(errors.New("error creating response")))
		return
	}
}

//CreateUser creates a user
func (a *UserAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*domain.User)

	a.UserService.ValidateUser(user)

	err := a.UserService.InsertUser(user)
	if err != nil {
		log.WithError(err).Error("error inserting user")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, stc.NewUserResponse(user)) //TODO change this
}

// UpdateUser updates the user
func (a *UserAPI) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*domain.User)

	data := &stc.UserRequest{User: user}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	user = data.User
	a.UserService.UpdateUser(user.ID, user)

	render.Render(w, r, stc.NewUserResponse(user))
}
