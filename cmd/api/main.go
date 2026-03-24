package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/redajn/task-mgr/internal/config"
	"github.com/redajn/task-mgr/internal/handler"
	"github.com/redajn/task-mgr/internal/repository/postgres"
	"github.com/redajn/task-mgr/internal/server"
	"github.com/redajn/task-mgr/internal/service"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connerct to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	taskRepo := postgres.NewTaskRepo(pool)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	srv := server.New(taskHandler)
	srv.Addr = cfg.Addr

	go func() {
		slog.Info("server started", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// graceful shutsown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting sown server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("forced shotdown", "error", err)
	}
	slog.Info("server stopped")
}
