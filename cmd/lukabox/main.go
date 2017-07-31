package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jacsmith21/lukabox/core/api"
)

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
		w.Write([]byte("lukabox api server"))
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/", api.Users)
		r.Put("/", api.CreateUser)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(api.UserCtx)
			r.Get("/", api.GetUser)
			r.Post("/", api.UpdateUser)
		})
	})

	http.ListenAndServe(":3000", r)
}
