package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/store"
)

type application struct {
	config config
	store  store.Storage
}
type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}
type config struct {
	addr string
	db   dbConfig
}

func (app *application) mount() http.Handler {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)

				r.Get("/", app.GetPostHandler)
				r.Delete("/", app.DeletePostHandler)
				r.Patch("/", app.UpdatePostHandler)
				r.Post("/comments", app.CreatePostCommentHandler)
			})
		})
		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/", app.GetUserHandler)
				r.Put("/follow", app.FollowUserHandler)
				r.Put("/unfollow", app.UnfollowUserHandler)
			})
			r.Group(func(r chi.Router) {

				r.Get("/feed", app.GetUserFeedHandler)
			})
		})

	})

	return r
}

func (app *application) run(mux http.Handler) error {

	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	log.Println("server started at ", app.config.addr)
	return srv.ListenAndServe()

}
