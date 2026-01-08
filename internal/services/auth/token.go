package auth

import (
	"context"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenString string) (*models.RefreshTokenResponse, error) {
	log := log("RefreshToken")

	token, err := s.VerifyRefreshToken(ctx, refreshTokenString)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to verify refresh token during refresh")
		return nil, jwt.ErrTokenExpired
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	jti := claims["jti"].(string)
	strUserID := claims["user_id"].(string)
	strAppID := claims["app_id"].(string)

	userID, errUserID := uuid.Parse(strUserID)
	appID, errAppID := uuid.Parse(strAppID)

	if errUserID != nil || errAppID != nil {
		log.Error().Err(err).Str("user_id", strUserID).Str("app_id", strAppID).Msg("Failed to parse userID or appID")
		return nil, utils.ErrInternalServerError
	}

	if err := s.repo.RevokeRefreshToken(ctx, appID, userID); err != nil {
		log.Error().Err(err).Str("jti", jti).Msg("Failed to revoke used refresh token")
		return nil, utils.ErrInternalServerError
	}

	user, err := s.userRepository.GetAppUserByID(ctx, appID, userID)
	if err != nil || user == nil {
		log.Error().Err(err).Str("user_id", userID.String()).Msg("User not found for refresh token")
		return nil, jwt.ErrTokenInvalidClaims
	}

	tokens, err := s.GenerateUserAuthTokens(ctx, user)

	newTokens := &models.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	return newTokens, err
}
