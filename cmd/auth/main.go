package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	redisClient "github.com/redis/go-redis/v9"

	"github.com/redajn/task-mgr/internal/config"
	"github.com/redajn/task-mgr/internal/handler"
	"github.com/redajn/task-mgr/internal/repository/postgres"
	"github.com/redajn/task-mgr/internal/repository/redis"
	"github.com/redajn/task-mgr/internal/service"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	rdb := redisClient.NewClient(&redisClient.Options{Addr: cfg.RedisAddr})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		slog.Error("failed to connect ro redis", "error", err)
		os.Exit(1)
	}
	defer rdb.Close()

	tokenRepo := redis.NewTokenRepo(rdb)
	userRepo := postgres.NewUserRepo(pool)
	authService := service.NewAuthService(userRepo, tokenRepo)
	authHandler := handler.NewAuthHandler(authService)

	r := chi.NewRouter()
	r.Use(handler.Recoverer)
	r.Use(handler.RequestID)
	r.Use(handler.Logger)

	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)
	r.Post("/logout", authHandler.Logout)

	srv := &http.Server{
		Addr:         cfg.AuthAddr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("auth service started", "addr", cfg.AuthAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("auth service error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
