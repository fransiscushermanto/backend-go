package user

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

type UserRepository struct {
	db      *utils.Database
	queries *db.Queries
}

func NewUserRepository(database *utils.Database) *UserRepository {
	return &UserRepository{
		db:      database,
		queries: db.New(database.Pool),
	}
}

var _ services.UserRepository = (*UserRepository)(nil)

func userLog(method string) *zerolog.Logger {
	l := utils.Log().With().Str("repository", "User").Str("method", method).Logger()
	return &l
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User, auth *models.UserAuthProvider) error {
	log := userLog("CreateUser")

	txFn := func(tx pgx.Tx) error {
		qtx := r.queries.WithTx(tx)

		if err := qtx.StoreUser(ctx, db.StoreUserParams{
			ID:    user.ID,
			AppID: user.AppID,
			Name:  user.Name,
			Email: user.Email,
		}); err != nil {
			log.Error().Err(err).Msg("Failed to insert user into DB")
			return fmt.Errorf("failed to create user: %w", err)
		}

		if err := qtx.StoreUserAuthProvider(ctx, db.StoreUserAuthProviderParams{
			AppID:          auth.AppID,
			UserID:         auth.UserID,
			Provider:       string(auth.Provider),
			ProviderUserID: auth.ProviderUserID,
			Password:       &auth.Password,
		}); err != nil {
			log.Error().Err(err).Msg("Failed to insert user authentication into DB")
			return fmt.Errorf("failed to create user authentication: %w", err)
		}
		return nil
	}

	return r.db.WithTransaction(ctx, txFn)
}

func (r *UserRepository) GetUserAuthenticationByProvider(ctx context.Context, appID, userID uuid.UUID, provider models.AuthProvider) (*models.UserAuthProvider, error) {
	log := userLog("GetUserAuthenticationByProvider")

	dbUserAuth, err := r.queries.GetUserAuthenticationByProvider(ctx, db.GetUserAuthenticationByProviderParams{
		AppID:    appID,
		UserID:   userID,
		Provider: string(provider),
	})

	userAuth := &models.UserAuthProvider{
		UserID:         dbUserAuth.AppID,
		AppID:          dbUserAuth.AppID,
		Provider:       models.AuthProvider(dbUserAuth.Provider),
		ProviderUserID: dbUserAuth.ProviderUserID,
		Password:       *dbUserAuth.Password,
		CreatedAt:      dbUserAuth.CreatedAt,
		UpdatedAt:      dbUserAuth.UpdatedAt,
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		log.Error().Err(err).Str("appID", appID.String()).Str("userID", userID.String()).Str("provider", string(provider)).Msg("Failed to query user authentication by provider")
		return nil, fmt.Errorf("failed to get user authentication: %w", err)
	}

	return userAuth, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	log := userLog("GetAllUsers")

	dbUsers, err := r.queries.GetAllUsers(ctx)

	if err != nil {
		log.Error().Err(err).Msg("Failed to query all users")
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	users := make([]*models.User, len(dbUsers))

	for i, dbUser := range dbUsers {
		users[i] = &models.User{
			ID:    dbUser.ID,
			AppID: dbUser.AppID,
			Name:  dbUser.Name,
			Email: dbUser.Email,
		}
	}

	return users, nil
}

func (r *UserRepository) GetAllUsersByAppID(ctx context.Context, appID uuid.UUID) ([]*models.User, error) {
	log := userLog("GetAllUsersByAppID")

	dbUsers, err := r.queries.GetAllUsersByAppID(ctx, appID)

	users := make([]*models.User, len(dbUsers))

	if err != nil {
		log.Error().Err(err).Msg("Failed to query all users")
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	for i, dbUser := range dbUsers {
		users[i] = &models.User{
			ID:    dbUser.ID,
			Name:  dbUser.Name,
			Email: dbUser.Email,
		}
	}

	return users, nil
}

func (r *UserRepository) GetAppUserByID(ctx context.Context, appID, id uuid.UUID) (*models.User, error) {
	log := userLog("GetAppUserByID")

	dbUser, err := r.queries.GetAppUserByID(ctx, db.GetAppUserByIDParams{
		AppID: appID,
		ID:    id,
	})

	user := &models.User{
		ID:              dbUser.ID,
		AppID:           dbUser.AppID,
		Name:            dbUser.Name,
		Email:           dbUser.Email,
		IsEmailVerified: dbUser.IsEmailVerified,
		EmailVerifiedAt: &dbUser.EmailVerifiedAt.Time,
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		log.Error().Err(err).Str("id", id.String()).Msg("Failed to query user by ID")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, appID uuid.UUID, email string) (*models.User, error) {
	log := userLog("GetUserByEmail")

	dbUser, err := r.queries.GetUserByEmail(ctx, db.GetUserByEmailParams{
		AppID: appID,
		Email: email,
	})

	user := &models.User{
		ID:              dbUser.ID,
		AppID:           dbUser.AppID,
		Name:            dbUser.Name,
		Email:           dbUser.Email,
		IsEmailVerified: dbUser.IsEmailVerified,
		EmailVerifiedAt: &dbUser.EmailVerifiedAt.Time,
		CreatedAt:       dbUser.CreatedAt,
		UpdatedAt:       dbUser.UpdatedAt,
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		log.Error().Err(err).Str("email", email).Msg("Failed to query user by email")
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}
