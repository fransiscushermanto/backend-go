package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/fransiscushermanto/backend/internal/utils"
)

type UserRepository struct {
	db *Database
}

func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{db: db}
}

var _ services.UserRepository = (*UserRepository)(nil)

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User, auth *models.UserAuthProvider) error {
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		utils.Log().Error().Err(err).Msg("Failed to begin transaction for user creation")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			utils.Log().Error().Err(err).Msg("Rolling back transaction due to error")
		}
	}()

	userQuery := `INSERT INTO core.users (id, app_id, name, email, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.ExecContext(ctx, userQuery, user.ID, user.AppID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		utils.Log().Error().Err(err).Msg("Failed to insert user into DB")
		return fmt.Errorf("failed to create user: %w", err)
	}

	authQuery := `INSERT INTO core.user_auth_providers (user_id, app_id, provider, provider_user_id, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.ExecContext(ctx, authQuery, auth.UserID, auth.AppID, auth.Provider, auth.ProviderUserID, auth.Password, auth.CreatedAt, auth.UpdatedAt)
	if err != nil {
		utils.Log().Error().Err(err).Msg("Failed to insert user authentication into DB")
		return fmt.Errorf("failed to create user authentication: %w", err)
	}

	if err = tx.Commit(); err != nil {
		utils.Log().Error().Err(err).Msg("Failed to commit transaction for user creation")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *UserRepository) GetUserAuthenticationByProvider(ctx context.Context, appID, userID string, provider models.AuthProvider) (*models.UserAuthProvider, error) {
	query := `SELECT user_id, app_id, provider, provider_user_id, password, created_at, updated_at FROM core.user_auth_providers WHERE app_id = $1 AND user_id = $2 AND provider = $3`
	auth := &models.UserAuthProvider{}

	err := r.db.QueryRowContext(ctx, query, appID, userID, provider).Scan(&auth.UserID, &auth.AppID, &auth.Provider, &auth.ProviderUserID, &auth.Password, &auth.CreatedAt, &auth.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		utils.Log().Error().Err(err).Str("appID", appID).Str("userID", userID).Str("provider", string(provider)).Msg("Failed to query user authentication by provider")
		return nil, fmt.Errorf("failed to get user authentication: %w", err)
	}

	return auth, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	query := `SELECT id, app_id, name, email FROM core.users ORDER BY created_at DESC`
	users := []*models.User{}
	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		utils.Log().Error().Err(err).Msg("Failed to query all users")
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.AppID, &user.Name, &user.Email); err != nil {
			utils.Log().Error().Err(err).Msg("Failed to scan user row")
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		utils.Log().Error().Err(err).Msg("Error encountered during rows iteration")
		return nil, fmt.Errorf("error encountered during rows iteration: %w", err)
	}

	return users, nil
}

func (r *UserRepository) GetAllUsersByAppID(ctx context.Context, appID string) ([]*models.User, error) {
	query := `SELECT id, name, email FROM core.users WHERE app_id=$1 ORDER BY created_at DESC`
	users := []*models.User{}
	rows, err := r.db.QueryContext(ctx, query, appID)

	if err != nil {
		utils.Log().Error().Err(err).Msg("Failed to query all users")
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			utils.Log().Error().Err(err).Msg("Failed to scan user row")
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		utils.Log().Error().Err(err).Msg("Error encountered during rows iteration")
		return nil, fmt.Errorf("error encountered during rows iteration: %w", err)
	}

	return users, nil
}

func (r *UserRepository) GetAppUserByID(ctx context.Context, appID string, id string) (*models.User, error) {
	query := `SELECT id, name, email FROM core.users WHERE app_id = $1 AND id = $2`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, appID, id).Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		utils.Log().Error().Err(err).Str("id", id).Msg("Failed to query user by ID")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, appID, email string) (*models.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM core.users WHERE app_id = $1 AND email = $2`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, appID, email).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		utils.Log().Error().Err(err).Str("email", email).Msg("Failed to query user by email")
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}
