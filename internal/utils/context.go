package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ContextKey string

const (
	UserIDContextKey    ContextKey = "user_id"
	AppIDContextKey     ContextKey = "app_id"
	TokenTypeContextKey ContextKey = "token_type"
	JTIContextKey       ContextKey = "jti"
)

func ContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

func GetUserIDFromContext(ctx context.Context) (*uuid.UUID, error) {
	strUserID, ok := ctx.Value(UserIDContextKey).(string)

	userID, errUserID := uuid.Parse(strUserID)

	if !ok {
		return nil, fmt.Errorf("missing %s in context", string(UserIDContextKey))
	}

	if errUserID != nil {
		return nil, fmt.Errorf("invalid %s in context", string(UserIDContextKey))
	}

	return &userID, nil
}

func GetAppIDFromContext(ctx context.Context) (*uuid.UUID, error) {
	strAppId, ok := ctx.Value(AppIDContextKey).(string)

	appID, errAppID := uuid.Parse(strAppId)

	if !ok {
		return nil, fmt.Errorf("missing %s in context", string(AppIDContextKey))
	}

	if errAppID != nil {
		return nil, fmt.Errorf("invalid %s in context", string(AppIDContextKey))
	}

	return &appID, nil
}

func GetTokenTypeFromContext(ctx context.Context) (string, error) {
	tokenType, ok := ctx.Value(TokenTypeContextKey).(string)
	if !ok {
		return "", fmt.Errorf("missing %s in context", string(TokenTypeContextKey))
	}
	return tokenType, nil
}

func GetJtiFromContext(ctx context.Context) (string, error) {
	jti, ok := ctx.Value(JTIContextKey).(string)
	if !ok {
		return "", fmt.Errorf("missing %s in context", string(JTIContextKey))
	}
	return jti, nil
}

func DebugContextValue(ctx context.Context, key ContextKey) {
	value := ctx.Value(key)
	if value == nil {
		fmt.Printf("%s: <not set>\n", key)
	} else {
		fmt.Printf("%s: %v (type: %T)\n", key, value, value)
	}
}
