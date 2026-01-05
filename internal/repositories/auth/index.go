package auth

import (
	"context"
	"fmt"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/repositories/db"
	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type AuthRepository struct {
	db      *utils.Database
	queries *db.Queries
}

func NewAuthRepository(database *utils.Database) *AuthRepository {
	return &AuthRepository{
		db:      database,
		queries: db.New(database.Pool),
	}
}

var _ services.AuthRepository = (*AuthRepository)(nil)

func authLog(method string) *zerolog.Logger {
	l := utils.Log().With().Str("repository", "Auth").Str("method", method).Logger()
	return &l
}

func (r *AuthRepository) StoreRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	err := r.queries.StoreRefreshToken(ctx, db.StoreRefreshTokenParams{
		Jti:       token.JTI,
		UserID:    token.UserID,
		AppID:     token.AppID,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		IsActive:  token.IsActive,
	})

	log := authLog("StoreRefreshToken")

	if err != nil {
		log.Error().Err(err).Msg("Failed to insert refresh token into DB")
		return fmt.Errorf("failed to insert refresh token: %w", err)
	}

	return nil
}

func (r *AuthRepository) GetRefreshTokenByJTI(ctx context.Context, appID uuid.UUID, jti string) (*models.RefreshToken, error) {
	dbRefreshToken, err := r.queries.GetRefreshTokenByJTI(ctx, db.GetRefreshTokenByJTIParams{
		AppID: appID,
		Jti:   jti,
	})

	refreshToken := &models.RefreshToken{
		JTI:       dbRefreshToken.Jti,
		UserID:    dbRefreshToken.UserID,
		AppID:     dbRefreshToken.AppID,
		Token:     dbRefreshToken.Token,
		ExpiresAt: dbRefreshToken.ExpiresAt,
		IsActive:  dbRefreshToken.IsActive,
		CreatedAt: dbRefreshToken.CreatedAt.Time,
	}

	log := authLog("GetRefreshTokenByJTI")

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		log.Error().Err(err).Str(appID.String(), "app_id").Str(jti, "jti").Msg("Failed to query refresh token")
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return refreshToken, nil
}

func (r *AuthRepository) GetUserActiveRefreshTokens(ctx context.Context, appID uuid.UUID, userID *uuid.UUID, jti *string) (*[]models.RefreshToken, error) {
	var dbTokens []db.CoreRefreshToken
	var sqlError error
	log := authLog("GetUserActiveRefreshTokens")

	if userID == nil && jti == nil {
		log.Error().Msg("Missing required parameters user_id and jti")
		return nil, fmt.Errorf("missing required parameters user_id and jti")
	}

	if userID != nil {
		dbTokens, sqlError = r.queries.GetUserActiveRefreshTokensByUserID(ctx, db.GetUserActiveRefreshTokensByUserIDParams{
			AppID:  appID,
			UserID: *userID,
		})
	} else {
		dbTokens, sqlError = r.queries.GetUserActiveRefreshTokensByJTI(ctx, db.GetUserActiveRefreshTokensByJTIParams{
			AppID: appID,
			Jti:   *jti,
		})

	}

	activeRefreshTokens := make([]models.RefreshToken, 0, len(dbTokens))
	for _, t := range dbTokens {
		activeRefreshTokens = append(activeRefreshTokens, models.RefreshToken{
			JTI:       t.Jti,
			UserID:    t.UserID,
			AppID:     t.AppID,
			Token:     t.Token,
			ExpiresAt: t.ExpiresAt,
			IsActive:  t.IsActive,
			CreatedAt: t.CreatedAt.Time,
		})
	}

	if sqlError != nil {
		log.Error().Err(sqlError).Msg("Failed to query user active refresh tokens")
		return nil, fmt.Errorf("failed to get user active refresh tokens: %w", sqlError)
	}

	return &activeRefreshTokens, nil
}

func (r *AuthRepository) RevokeRefreshToken(ctx context.Context, appID, userID uuid.UUID, jtis []string) error {
	if len(jtis) == 0 {
		return nil
	}

	err := r.queries.RevokeRefreshTokens(ctx, db.RevokeRefreshTokensParams{
		AppID:  appID,
		UserID: userID,
		Jtis:   jtis,
	})

	if err != nil {
		return fmt.Errorf("failed to revoke tokens: %w", err)
	}

	return nil
}
