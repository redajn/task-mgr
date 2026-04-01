package service

import (
	"context"
	"fmt"

	"github.com/redajn/task-mgr/internal/domain"
	"github.com/redajn/task-mgr/internal/token"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users  domain.UserRepository
	tokens *token.Store
}

func NewAuthService(users domain.UserRepository, tokens *token.Store) *AuthService {
	return &AuthService{users: users, tokens: tokens}
}

func (s *AuthService) Register(ctx context.Context, input domain.RegisterInput) (domain.User, error) {
	if input.Email == "" {
		return domain.User{}, fmt.Errorf("email is required")
	}
	if len(input.Password) < 8 {
		return domain.User{}, fmt.Errorf("password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.users.Create(ctx, input.Email, string(hash))
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, input domain.LoginInput) (string, error) {
	user, err := s.users.GetByEmail(ctx, input.Email)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err := token.Generate()
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	if err := s.tokens.Save(ctx, token, domain.TokenInfo{
		UserID: user.ID,
		Email:  user.Email,
	}); err != nil {
		return "", fmt.Errorf("save token: %w", err)
	}

	return token, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.tokens.Delete(ctx, token)
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (domain.TokenInfo, error) {
	return s.tokens.Get(ctx, token)
}
