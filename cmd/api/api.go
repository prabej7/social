package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prabej7/social/internal/auth"
	"github.com/prabej7/social/internal/mailer"
	"github.com/prabej7/social/internal/ratelimiter"
	"github.com/prabej7/social/internal/store"
	"github.com/prabej7/social/internal/store/cache"

	"go.uber.org/zap"

	"github.com/prabej7/social/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

type application struct {
	config        config
	store         store.Storage
	timeout       time.Duration
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Autheticator
	cache         cache.Storage
	ratelimiter   ratelimiter.Limiter
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type sendGridConfig struct {
	apiKey string
}

type smtpConfig struct {
	password string
}

type mailConfig struct {
	exp            time.Duration
	fromEmail      string
	sendGridConfig sendGridConfig
	smtpConfig     smtpConfig
}

type authConfig struct {
	basic basicConfig
	tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type basicConfig struct {
	username string
	password string
}

type redisConfig struct {
	addr     string
	password string
	db       int
	enabled  bool
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiUrl      string
	mail        mailConfig
	frontendURL string
	auth        authConfig
	redisConfig redisConfig
	rateLimiter ratelimiter.Config
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(app.RateLimiterMiddleware)	

	// r.Use(middleware.Timeout(60 * time.Second))
	r.Use(Timeout(app.timeout))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, http.StatusOK, map[string]string{
			"message": "Welcome to social API!",
		})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		docsURL := "/api/v1/swagger/doc.json"
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Group(func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Route("/post", func(r chi.Router) {
				r.Post("/", validatePayload[CreatePostPayload](app.createPostHandler))

				r.Route("/{id}", func(r chi.Router) {
					r.Use(app.postContextMiddleware)
					r.Get("/", app.getPostByIdHandler)
					r.Patch("/", app.checkPostOwnership("moderator", validatePayload[UpdatePostPayload](app.update)))
					r.Delete("/", app.checkPostOwnership("admin", app.delete))
				})
			})

		})

		r.Group(func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Route("/comment", func(r chi.Router) {
				r.Post("/", validatePayload[CreateCommentPayload](app.createCommentHandler))
			})
		})

		r.Route("/user", func(r chi.Router) {
			r.Get("/activate/{token}", app.activateUserHandler)

			r.With(app.AuthTokenMiddleware).Get("/", app.getUserHandler)
			r.With(app.AuthTokenMiddleware).Route("/{id}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unFollowUserHandler)
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Get("/user/feed", app.getUserFeedHandler)
		})

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", validatePayload[RegisterUserPayload](app.registerUserHandler))
			r.Post("/token", validatePayload[CreateTokenPayload](app.createTokenHandler))
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.BasePath = "/api/v1"
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	shutdown <- err

	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil
}
