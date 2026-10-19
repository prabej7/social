package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prabej7/social/internal/store"
)

type application struct {
	config  config
	store   store.Storage
	timeout time.Duration
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
	env  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(Timeout(app.timeout))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, http.StatusOK, map[string]string{
			"message": "Welcome to social API!",
		})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/post", func(r chi.Router) {
			r.Post("/", validatePayload[CreatePostPayload](app.createPostHandler))

			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.postContextMiddleware)
				r.Get("/", app.getPostByIdHandler)
				r.Patch("/", validatePayload[UpdatePostPayload](app.update))
				r.Delete("/", app.delete)
			})
		})

		r.Route("/comment", func(r chi.Router) {
			r.Post("/", validatePayload[CreateCommentPayload](app.createCommentHandler))
		})

		r.Route("/user", func(r chi.Router) {
			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)
				r.Get("/", app.getUserHandler)
				r.Put("/follow", validatePayload[FollowUser](app.followUserHandler))
				r.Put("/unfollow", validatePayload[FollowUser](app.unFollowUserHandler))
			})
		})

		r.Group(func(r chi.Router) {
			r.Get("/user/feed", app.getUserFeedHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	return srv.ListenAndServe()
}
