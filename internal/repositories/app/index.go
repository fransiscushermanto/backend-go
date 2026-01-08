package app

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/repositories/db"
	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type AppRepository struct {
	db          *utils.Database
	queries     *db.Queries
	lockTimeout *int
}

func NewAppRepository(database *utils.Database, lockTimeout *int) *AppRepository {
	return &AppRepository{
		db:          database,
		queries:     db.New(database.Pool),
		lockTimeout: lockTimeout,
	}
}

type LockTimeoutError struct {
	AppID   string
	Timeout time.Duration
}

func (e *LockTimeoutError) Error() string {
	return fmt.Sprintf("lock acquisition timeout for app %s after %v", e.AppID, e.Timeout)
}

var _ services.AppRepository = (*AppRepository)(nil)

func appLog(method string) *zerolog.Logger {
	l := utils.Log().With().Str("repository", "App").Str("method", method).Logger()
	return &l
}

func (r *AppRepository) RegisterApp(ctx context.Context, app *models.App, appApiKey *models.AppApiKey) error {
	log := appLog("RegisterApp")

	txFn := func(tx pgx.Tx) error {
		qtx := r.queries.WithTx(tx)

		if err := qtx.StoreApp(ctx, db.StoreAppParams{
			ID:   app.ID,
			Name: app.Name,
		}); err != nil {
			log.Error().Err(err).Msg("Failed to insert app into DB")
			return fmt.Errorf("failed to insert app: %w", err)
		}

		if err := qtx.StoreAppApiKey(ctx, db.StoreAppApiKeyParams{
			ID:       appApiKey.ID,
			AppID:    appApiKey.AppID,
			KeyHash:  appApiKey.KeyHash,
			IsActive: appApiKey.IsActive,
		}); err != nil {
			log.Error().Err(err).Msg("Failed to insert app api key into DB")
			return fmt.Errorf("failed to create app api key: %w", err)
		}
		return nil
	}

	return r.db.WithTransaction(ctx, txFn)
}

func (r *AppRepository) RegenerateAppApiKey(ctx context.Context, revokedAt time.Time, appApiKey *models.AppApiKey) error {
	log := appLog("RegenerateAppApiKey")

	maxRetries := 3
	baseDelay := 100 * time.Millisecond

	attemptRegenerateApiKey := func(ctx context.Context, revokedAt time.Time, appApiKey *models.AppApiKey) error {
		lockTimeout := getLockTimeout(r.lockTimeout)
		lockCtx, cancel := context.WithTimeout(ctx, lockTimeout)
		defer cancel()

		txFn := func(tx pgx.Tx) error {
			qtx := r.queries.WithTx(tx)
			// Use lockCtx here to enforce the timeout on acquiring the row lock.
			_, err := qtx.LockAppForUpdate(ctx, appApiKey.AppID)
			if err != nil {
				if err == pgx.ErrNoRows {
					log.Error().Str("app_id", appApiKey.AppID.String()).Msg("App not found")
					return fmt.Errorf("app not found: %s", appApiKey.AppID)
				}
				// Check if the error was specifically a timeout on our lock context.
				if lockCtx.Err() == context.DeadlineExceeded {
					log.Error().Str("app_id", appApiKey.AppID.String()).Msg("Lock acquisition timed out")
					return &LockTimeoutError{AppID: appApiKey.AppID.String(), Timeout: lockTimeout}
				}
				log.Error().Err(err).Msg("Failed to lock app for regeneration")
				return fmt.Errorf("failed to lock app: %w", err)
			}

			// Use the original `ctx` for subsequent operations as the lock is already acquired.
			revokedAppApiKeys, err := qtx.RevokeActiveAppApiKeys(ctx, db.RevokeActiveAppApiKeysParams{
				AppID:     appApiKey.AppID,
				RevokedAt: utils.ToPgTimestamp(revokedAt),
			})

			if err != nil {
				log.Error().Err(err).Msg("Failed to revoke old token from DB")
				return fmt.Errorf("failed to revoke old token: %w", err)
			}

			log.Info().Int64("revoked_keys", revokedAppApiKeys).Msg("Revoked active keys")

			err = qtx.StoreAppApiKey(ctx, db.StoreAppApiKeyParams{
				ID:       appApiKey.ID,
				AppID:    appApiKey.AppID,
				KeyHash:  appApiKey.KeyHash,
				IsActive: appApiKey.IsActive,
			})

			if err != nil {
				log.Error().Err(err).Msg("Failed to insert new app api key into DB")
				return fmt.Errorf("failed to generate new app api key: %w", err)
			}

			return nil
		}

		return r.db.WithTransaction(lockCtx, txFn)
	}

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(math.Pow(2, float64(attempt-1))) * baseDelay
			log.Info().Int("attempt", attempt).Dur("delay", delay).Msg("Retrying after lock timeout")

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		err := attemptRegenerateApiKey(ctx, revokedAt, appApiKey)

		if err == nil {
			if attempt > 0 {
				log.Info().Int("attempt", attempt+1).Msg("Successfully regenerated API key after retries")
			}

			return nil
		}

		if isLockTimeoutError(err) && attempt < maxRetries {
			log.Warn().Err(err).Int("attempt", attempt+1).Int("max_retries", maxRetries).Msg("Lock timeout, will retry")
			continue
		}

		// Non-retryable error or max retries reached
		if attempt >= maxRetries {
			log.Error().Err(err).Int("max_retries", maxRetries).Msg("Failed to regenerate API key after all retries")
			return fmt.Errorf("failed to regenerate API key after %d attempts: %w", maxRetries+1, err)
		}

		// Fast fail for non-timeout errors
		log.Error().Err(err).Msg("Non-retryable error during API key regeneration")
		return err
	}

	return fmt.Errorf("unexpected error: retry loop completed without result")
}

func (r *AppRepository) GetAllApps(ctx context.Context) ([]*models.App, error) {
	dbApps, err := r.queries.GetAllApps(ctx)

	if err != nil {
		utils.Log().Error().Err(err).Msg("Failed to query all apps")
		return nil, fmt.Errorf("failed to get all apps: %w", err)
	}

	apps := make([]*models.App, len(dbApps))
	for i, dbApp := range dbApps {
		apps[i] = &models.App{
			ID:   dbApp.ID,
			Name: dbApp.Name,
		}
	}

	return apps, nil
}

func (r *AppRepository) GetAppById(ctx context.Context, id uuid.UUID) (*models.App, error) {
	result, err := r.queries.GetAppByID(ctx, id)
	app := &models.App{
		ID:        result.ID,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		utils.Log().Error().Err(err).Msg("Failed to query app by id")
		return nil, fmt.Errorf("failed to get app by id: %w", err)
	}

	return app, nil
}

func isLockTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*LockTimeoutError)
	return ok
}

func getLockTimeout(lockTimeout *int) time.Duration {
	if lockTimeout != nil && *lockTimeout > 0 {
		return time.Duration(*lockTimeout) * time.Second
	}
	return 30 * time.Second
}
