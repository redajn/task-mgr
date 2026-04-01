package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenInfo struct {
	UserID int64
	Email  string
}

type UserRepository interface {
	Create(ctx context.Context, email, passwordHash string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByID(ctx context.Context, id int64) (User, error)
}

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (User, error)
	Login(ctx context.Context, input LoginInput) (string, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (TokenInfo, error)
}
