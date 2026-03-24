package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/redajn/task-mgr/internal/handler"
)

func New(taskHandler *handler.TaskHandler) *http.Server {
	r := chi.NewRouter()

	r.Use(handler.Recoverer)
	r.Use(handler.RequestID)
	r.Use(handler.Logger)
	r.Use(chimiddleware.Compress(5))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

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

	return &http.Server{
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}
