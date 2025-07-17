package utils

import (
	"context"
	"time"
)

func ContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

func GetUserIDFromContext(ctx context.Context) (string, error) {
	userId, ok := ctx.Value("user_id").(string)
	if !ok {
		return "", ErrUnauthorized
	}
	return userId, nil
}

func GetAppIDFromContext(ctx context.Context) (string, error) {
	appId, ok := ctx.Value("app_id").(string)
	if !ok {
		return "", ErrUnauthorized
	}
	return appId, nil
}

func GetTokenTypeFromContext(ctx context.Context) (string, error) {
	tokenType, ok := ctx.Value("token_type").(string)
	if !ok {
		return "", ErrUnauthorized
	}
	return tokenType, nil
}

func GetJtiFromContext(ctx context.Context) (string, error) {
	jti, ok := ctx.Value("jti").(string)
	if !ok {
		return "", ErrUnauthorized
	}
	return jti, nil
}
