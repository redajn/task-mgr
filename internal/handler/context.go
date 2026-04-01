package handler

import (
	"context"

	"github.com/redajn/task-mgr/internal/domain"
)

type contextKey string

const requestIDKey contextKey = "request_id"
const tokenInfoKey contextKey = "token_info"

func contextWithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func RequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}

func contextWithTokenInfo(ctx context.Context, info domain.TokenInfo) context.Context {
	return context.WithValue(ctx, tokenInfoKey, info)
}

func TokenInfoFromContext(ctx context.Context) (domain.TokenInfo, bool) {
	info, ok := ctx.Value(tokenInfoKey).(domain.TokenInfo)
	return info, ok
}
