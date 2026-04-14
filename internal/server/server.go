package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/redajn/task-mgr/internal/config"
	"github.com/redajn/task-mgr/internal/handler"
	"github.com/redajn/task-mgr/internal/repository/redis"
)

func New(taskHandler *handler.TaskHandler, authHandler *handler.AuthHandler, tokenRepo *redis.TokenRepo, cfg config.Config) *http.Server {
	r := chi.NewRouter()

	r.Use(handler.Recoverer)
	r.Use(handler.RequestID)
	r.Use(handler.Logger)
	r.Use(chimiddleware.Compress(5))
	r.Use(chimiddleware.StripSlashes)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Group(func(r chi.Router) {
		r.Use(handler.Auth(tokenStore))

		r.Route("/api/v1", func(r chi.Router) {
			r.Route("/tasks", func(r chi.Router) {
				r.Get("/", taskHandler.List)
				r.Post("/", taskHandler.Create)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", taskHandler.Get)
					r.Patch("/", taskHandler.Update)
					r.Delete("/", taskHandler.Delete)
				})
			})
		})
	})

	return &http.Server{
		Addr:         cfg.Addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}
