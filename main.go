package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/jacsmith21/lukabox/api"
	"github.com/jacsmith21/lukabox/db"
)

var tokenAuth *jwtauth.JwtAuth

func main() {
	r := chi.NewRouter()

	//Creating token with secret string
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

	// Creating services
	var userService = db.UserService{}
	var authenticationService = db.AuthenticationService{}

	// Creating apis
	var userAPI api.UserAPI
	var authenticationAPI api.AuthenticationAPI

	// Adding services to apis
	userAPI.UserService = &userService
	authenticationAPI.AuthenticationService = &authenticationService
	authenticationAPI.UserService = &userService

	// The middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	//to stop processing after 60 seconds
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("lukabox api server!"))
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	r.Post("/singup", api.SignUp)
	r.Post("/login", authenticationAPI.Login)

	r.Route("/users", func(r chi.Router) {
		r.Get("/", userAPI.Users)
		r.Put("/", userAPI.CreateUser)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(userAPI.UserCtx)
			r.Get("/", userAPI.UserByID)
			r.Post("/", userAPI.UpdateUser)

			r.Route("/pills", func(r chi.Router) {
				r.Use(jwtauth.Verifier(tokenAuth))
				r.Use(jwtauth.Authenticator)
				r.Get("/", api.PillsByUser)
			})
		})
	})

	http.ListenAndServe(":3001", r)
}
