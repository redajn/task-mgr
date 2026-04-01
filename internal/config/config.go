package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr         string
	AuthAddr     string
	DatabaseURL  string
	RedisAddr    string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		Addr:         getEnv("ADDR", ":8080"),
		AuthAddr:     getEnv("AUTH_ADDR", ":8081"),
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		ReadTimeout:  getDuration("READ_TIMEOUT", 5*time.Second),
		WriteTimeout: getDuration("WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:  getDuration("IDLE_TIMEOUT", 120*time.Second),
	}

	if err := cfg.validate(); err != nil {
		return Config{}, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func (c Config) validate() error {
	var errs []error

	if c.DatabaseURL == "" {
		errs = append(errs, errors.New("DATABASE_URL is required"))
	}
	if c.RedisAddr == "" {
		errs = append(errs, errors.New("REDIS_ADDR is required"))
	}
	if c.ReadTimeout <= 0 {
		errs = append(errs, errors.New("READ_TIMEOUT must be positive"))
	}
	if c.WriteTimeout <= 0 {
		errs = append(errs, errors.New("WRITE_TIMEOUT must be positive"))
	}

	return errors.Join(errs...)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	// Поддерживаем секунды как число: READ_TIMEOUT=5
	if secs, err := strconv.Atoi(val); err == nil {
		return time.Duration(secs) * time.Second
	}
	// И стандартный формат Go: READ_TIMEOUT=5s
	if d, err := time.ParseDuration(val); err == nil {
		return d
	}
	return fallback
}
