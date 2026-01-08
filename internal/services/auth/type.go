package auth

import (
	"context"
	"crypto/ecdsa"
	"errors"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/services/user"
	"github.com/google/uuid"
)

type AuthRepository interface {
	StoreRefreshToken(ctx context.Context, token *models.RefreshToken) error
	GetRefreshTokenByJTI(ctx context.Context, appID uuid.UUID, jti string) (*models.RefreshToken, error)
	GetUserActiveRefreshTokens(ctx context.Context, appID uuid.UUID, userID *uuid.UUID, jti *string) (*[]models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, appID, userID uuid.UUID) error
	StoreResetPasswordToken(ctx context.Context, token *models.ResetPasswordToken) error
	RevokeResetPasswordToken(ctx context.Context, appID, userID uuid.UUID) error
}

type AuthService struct {
	repo           AuthRepository
	userRepository user.UserRepository
	userService    *user.UserService
	privateKey     *ecdsa.PrivateKey
	publicKey      *ecdsa.PublicKey
}

type AuthOptions struct {
	CallbackURL string
	RedirectURL string
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected token signing method")
	ErrTokenRevoked            = errors.New("token has been revoked")
	ErrTokenNotFound           = errors.New("token not found in storage")
	ErrTokenMismatch           = errors.New("stored token does not match provided token")
	ErrMissingRequiredClaim    = errors.New("token is missing a required claim")
)
