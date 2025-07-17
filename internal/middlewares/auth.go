package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	authService *services.AuthService
}

func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, utils.ErrorResponsePayload{
				StatusCode: http.StatusUnauthorized,
				Message:    utils.StringPointer("Authorization header required"),
			})
			return
		}

		// Check Bearer format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondWithError(w, utils.ErrorResponsePayload{
				StatusCode: http.StatusUnauthorized,
				Message:    utils.StringPointer("Invalid authorization header format"),
			})
			return
		}

		tokenString := parts[1]

		// Verify token
		token, err := m.authService.VerifyAccessToken(r.Context(), tokenString)
		if err != nil {
			utils.RespondWithError(w, utils.ErrorResponsePayload{
				StatusCode: http.StatusUnauthorized,
				Message:    utils.StringPointer("Invalid or expired token"),
			})
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.RespondWithError(w, utils.ErrorResponsePayload{
				StatusCode: http.StatusUnauthorized,
				Message:    utils.StringPointer("Invalid token claims"),
			})
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
		ctx = context.WithValue(ctx, "app_id", claims["app_id"])
		ctx = context.WithValue(ctx, "token_type", claims["type"])
		ctx = context.WithValue(ctx, "jti", claims["jti"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
