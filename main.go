package main

import (
	"fmt"
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

	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

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
		r.Use(jwtauth.Authenticator)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("welcome :) it worked! id: %v", claims["id"])))
		})
	})

	http.ListenAndServe(":3001", r)
}
