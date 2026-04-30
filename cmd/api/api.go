package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/sahil1si18ec083/Social-media-app-Golang/docs"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/auth"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/mailer"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/store"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type application struct {
	config config
	store  store.Storage
	mailer mailer.Client
	auth   auth.Authenticator
}
type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}
type config struct {
	addr        string
	db          dbConfig
	auth        authConfig
	mail        mailConfig
	env         string
	frontendURL string
}
type authConfig struct {
	token tokenConfig
}
type tokenConfig struct {
	exp       time.Duration
	secretKey string
}
type mailConfig struct {
	sendGrid sendGridConfig
	mailTrap mailTrapConfig

	exp time.Duration
}
type sendGridConfig struct {
	apikey    string
	fromEmail string
}
type mailTrapConfig struct {
	fromEmail string
	host      string
	port      int
	username  string
	password  string
}

func (app *application) mount() http.Handler {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(app.BasicAuthMiddleware)
			r.Get("/swagger/*", httpSwagger.Handler(
				httpSwagger.URL("http://localhost:8080/v1/swagger/doc.json"),
			))
		})

		r.Get("/health", app.healthCheckHandler)

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.RegisterUserHandler)
			// r.Post("/token", app.createTokenHandler)
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.ActivateUserHandler)

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)

				r.Get("/feed", app.GetUserFeedHandler)

				r.Route("/{userID}", func(r chi.Router) {
					r.Get("/", app.GetUserHandler)
					r.Put("/follow", app.FollowUserHandler)
					r.Put("/unfollow", app.UnfollowUserHandler)
				})
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)

			r.Route("/posts", func(r chi.Router) {
				r.Post("/", app.createPostHandler)
				r.Route("/{postID}", func(r chi.Router) {
					r.Use(app.postsContextMiddleware)
					r.Get("/", app.GetPostHandler)
					r.Delete("/", app.checkPostOwnership("admin", app.DeletePostHandler))
					r.Patch("/", app.checkPostOwnership("moderator", app.UpdatePostHandler))
					r.Post("/comments", app.CreatePostCommentHandler)
				})

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
