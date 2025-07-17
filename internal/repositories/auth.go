package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/fransiscushermanto/backend/internal/utils"
)

type AuthRepository struct {
	db *Database
}

func NewAuthRepository(db *Database) *AuthRepository {
	return &AuthRepository{db: db}
}

var _ services.AuthRepository = (*AuthRepository)(nil)

func (r *AuthRepository) StoreRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	query := `INSERT INTO core.refresh_tokens (jti, user_id, app_id, token, expires_at, is_active, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, token.JTI, token.UserID, token.AppID, token.Token, token.ExpiresAt, token.IsActive, token.CreatedAt)

	if err != nil {
		utils.Log().Error().Err(err).Msg("Failed to insert refresh token into DB")
		return fmt.Errorf("failed to insert refresh token: %w", err)
	}

	return nil
}

func (r *AuthRepository) GetRefreshToken(ctx context.Context, appID, userID, jti string) (*models.RefreshToken, error) {
	query := `SELECT jti, user_id, app_id, token, expires_at, is_active, created_at FROM core.refresh_tokens WHERE app_id = $1 AND user_id = $2 AND jti = $3 ORDER BY created_at DESC LIMIT 1`
	refreshToken := &models.RefreshToken{}

	err := r.db.QueryRowContext(ctx, query, appID, userID, jti).Scan(&refreshToken.JTI, &refreshToken.UserID, &refreshToken.AppID, &refreshToken.Token, &refreshToken.ExpiresAt, &refreshToken.IsActive, &refreshToken.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		utils.Log().Error().Err(err).Str(appID, "app_id").Str(userID, "user_id").Str(jti, "jti").Msg("Failed to query refresh token")
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return refreshToken, nil
}
