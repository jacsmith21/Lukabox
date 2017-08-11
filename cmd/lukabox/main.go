package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/jacsmith21/lukabox/core/api"
)

var tokenAuth *jwtauth.JwtAuth

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
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
	r.Post("/login", api.Login)

	r.Route("/users", func(r chi.Router) {
		r.Get("/", api.Users)
		r.Put("/", api.CreateUser)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(api.UserCtx)
			r.Get("/", api.GetUser)
			r.Post("/", api.UpdateUser)
		})
	})

	r.Route("/pills", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(Authenticator)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome :) it worked!"))
		})
	})

	http.ListenAndServe(":3001", r)
}

//Authenticator authenticates the user credentials
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if token == nil || !token.Valid {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}
