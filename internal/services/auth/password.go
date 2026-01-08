package auth

import (
	"context"
	"time"

	"github.com/fransiscushermanto/backend/internal/constants"
	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/services/user"
	"github.com/golang-jwt/jwt/v5"
)

func (s *AuthService) ForgetPassword(ctx context.Context, req *models.ForgetPasswordRequest) error {
	forgetPasswordLog := log("ForgetPassword")

	user, err := s.userService.GetUser(ctx, *req.AppID, user.UserIdentifier{
		Email: req.Email,
	})

	resetPasswordTokenExpiryTime := constants.DEFAULT_JWT_EXPIRY_HOURS
	resetPasswordTokenJTI := generateTokenID()
	resetPasswordTokenClaims := jwt.MapClaims{
		"jti":     resetPasswordTokenJTI,
		"user_id": user.ID,
		"app_id":  user.AppID,
		"exp":     resetPasswordTokenExpiryTime.Unix(),
		"type":    "reset-password",
		"iat":     time.Now().Unix(),
	}
	resetToken, err := s.GenerateToken(constants.DEFAULT_JWT_SIGNING_METHOD, resetPasswordTokenClaims)
	if err != nil {
		forgetPasswordLog.Error().Err(err).Msg("Failed to generate reset password token")
		return err
	}

	err = s.repo.RevokeResetPasswordToken(ctx, user.AppID, user.ID)

	if err != nil {
		forgetPasswordLog.Error().Err(err).Msg("Failed to revoke reset password token")
		return err
	}

	err = s.repo.StoreResetPasswordToken(ctx, &models.ResetPasswordToken{
		JTI:       resetPasswordTokenJTI,
		AppID:     user.AppID,
		UserID:    user.ID,
		Token:     hashToken(*resetToken),
		ExpiresAt: resetPasswordTokenExpiryTime,
	})

	if err != nil {
		forgetPasswordLog.Error().Err(err).Msg("Failed to store reset password token")
		return err
	}

	forgetPasswordLog.Info().Str("token", *resetToken).Interface("user", user).Msg("Successfully request reset password")
	return nil
}
