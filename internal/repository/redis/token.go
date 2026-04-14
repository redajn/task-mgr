package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	redisClient "github.com/redis/go-redis/v9"

	"github.com/redajn/task-mgr/internal/domain"
)

const (
	tokenTTL    = 24 * time.Hour
	tokenPrefix = "session:"
)

type TokenRepo struct {
	client *redisClient.Client
}

func NewTokenRepo(client *redisClient.Client) *TokenRepo {
	return &TokenRepo{client: client}
}

func Generate() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func (s *TokenRepo) Save(ctx context.Context, token string, info domain.TokenInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("marshal token info: %w", err)
	}

	key := tokenPrefix + token
	if err := s.client.Set(ctx, key, data, tokenTTL).Err(); err != nil {
		return fmt.Errorf("save token: %w", err)
	}
	return nil
}

func (s *TokenRepo) Get(ctx context.Context, token string) (domain.TokenInfo, error) {
	key := tokenPrefix + token
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redisClient.Nil {
			return domain.TokenInfo{}, domain.ErrInvalidToken
		}
		return domain.TokenInfo{}, fmt.Errorf("get token:%w", err)
	}

	var info domain.TokenInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return domain.TokenInfo{}, fmt.Errorf("unmarshal token info: %w", err)
	}
	return info, nil
}

func (s *TokenRepo) Delete(ctx context.Context, token string) error {
	key := tokenPrefix + token
	if err := s.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("delete token: %w", err)
	}
	return nil
}
