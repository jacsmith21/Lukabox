package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
	"github.com/jacsmith21/lukabox/domain"
	"github.com/jacsmith21/lukabox/ext/log"
	"github.com/jacsmith21/lukabox/ext/render"
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
		if userID == "" {
			render.WithMessage("user id must be supplied").BadRequest(w, r)
			return
		}
		log.WithField("id", userID).Debug("user id from paramter")

		id, err := strconv.Atoi(userID)
		if err != nil {
			log.WithError(err).Debugf("unable to parse %s", userID)
			render.WithMessage("unable to parse parameter id").BadRequest(w, r)
			return
		}

		user, err := a.UserService.UserByID(id)
		if err != nil {
			log.WithError(err).Errorf("error fetching user with id %d", id)
			render.WithError(err).InternalServerError(w, r)
			return
		}
		if user == nil {
			log.Debugf("no user found with id %d", id)
			render.WithMessage("user not found").NotFound(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserRequestCtx a user request context generator
func (a *UserAPI) UserRequestCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithField("method", "UserRequestCtx").Info("starting")
		userRequest := &stc.UserRequest{}

		err := render.Bind(r, userRequest)
		if err != nil {
			log.WithError(err).Error("error binding user request")
			render.WithError(err).InternalServerError(w, r)
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
	if err := render.Instance(w, r, stc.NewUserResponse(user)); err != nil {
		log.WithError(err).Error("unable to render user response")
		render.WithError(err).InternalServerError(w, r)
		return
	}
}

// Users lists the users using the RenderList function
func (a *UserAPI) Users(w http.ResponseWriter, r *http.Request) {
	log.WithField("method", "Users").Info("starting")
	users, err := a.UserService.Users()
	if err != nil {
		log.WithError(err).Error("error fetching users")
		render.WithError(err).InternalServerError(w, r)
		return
	}

	if err := render.List(w, r, stc.NewUserListResponse(users)); err != nil {
		log.WithError(err).Error("error rendering user list response")
		render.WithMessage("error creating response").InternalServerError(w, r)
		return
	}
}

//CreateUser creates a user
func (a *UserAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
	log.WithField("method", "CreateUser").Info("starting")
	user := r.Context().Value("user").(*domain.User)

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		log.WithError(err).Debug("user wasn't validated")
		render.WithError(err).BadRequest(w, r)
		return
	}

	if err := a.UserService.InsertUser(user); err != nil {
		log.WithError(err).Error("error inserting user")
		render.WithError(err).InternalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UpdateUser updates the user
func (a *UserAPI) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*domain.User)

	data := &stc.UserRequest{User: user}
	if err := render.Bind(r, data); err != nil {
		render.WithError(err).BadRequest(w, r)
		return
	}

	user = data.User
	a.UserService.UpdateUser(user.ID, user)
}
