package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/services"
	authTypes "github.com/fransiscushermanto/backend/internal/services/auth"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

type AuthMiddleware struct {
	authService *services.AuthService
}

func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func authMiddlewareLog(method string) *zerolog.Logger {
	l := utils.Log().With().Str("middleware", "Auth").Str("method", method).Logger()
	return &l
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	requireAuthLog := authMiddlewareLog("RequireAuth")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, models.ApiError{
				StatusCode: http.StatusUnauthorized,
				Message:    utils.StringPointer("Authorization header required"),
			})
			return
		}

		// Check Bearer format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondWithError(w, models.ApiError{
				StatusCode: http.StatusUnauthorized,
				Message:    utils.StringPointer("Invalid authorization header format"),
			})
			return
		}

		tokenString := parts[1]

		// Verify token
		token, err := m.authService.VerifyAccessToken(r.Context(), tokenString)

		if err != nil {
			var statusCode int
			var message string
			var errorCode models.ErrorCode

			isInvalidToken := errors.Is(err, authTypes.ErrTokenRevoked) || errors.Is(err, authTypes.ErrMissingRequiredClaim) || errors.Is(err, authTypes.ErrTokenMismatch) || errors.Is(err, authTypes.ErrTokenNotFound) || errors.Is(err, jwt.ErrTokenMalformed)

			// Check specific error types
			if errors.Is(err, jwt.ErrTokenExpired) {
				statusCode = http.StatusUnauthorized
				message = "Token has expired"
				errorCode = models.CodeTokenExpired
			} else if isInvalidToken {
				statusCode = http.StatusUnauthorized
				message = "Invalid token"
				errorCode = models.CodeTokenInvalid
			} else {
				statusCode = http.StatusInternalServerError
				message = "Internal Server Error"
			}

			requireAuthLog.Error().Err(err).Msg("Failed to verify access token")
			utils.RespondWithError(w, models.ApiError{
				StatusCode: statusCode,
				Message:    utils.StringPointer(message),
				Meta: &models.ErrorMeta{
					Code: errorCode,
				},
			})
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.RespondWithError(w, models.ApiError{
				StatusCode: http.StatusUnauthorized,
				Message:    utils.StringPointer("Invalid token claims"),
			})
			return
		}

		ctx := context.WithValue(r.Context(), utils.UserIDContextKey, claims[string(utils.UserIDContextKey)])
		ctx = context.WithValue(ctx, utils.AppIDContextKey, claims[string(utils.AppIDContextKey)])
		ctx = context.WithValue(ctx, utils.TokenTypeContextKey, claims[string(utils.TokenTypeContextKey)])
		ctx = context.WithValue(ctx, utils.JTIContextKey, claims[string(utils.JTIContextKey)])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
