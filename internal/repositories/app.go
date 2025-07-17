package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/fransiscushermanto/backend/internal/utils"
)

type AppRepository struct {
	db *Database
}

func NewAppRepository(db *Database) *AppRepository {
	return &AppRepository{
		db: db,
	}
}

var _ services.AppRepository = (*AppRepository)(nil)

func (r *AppRepository) RegisterApp(ctx context.Context, app *models.App) error {
	query := `INSERT INTO core.apps (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, app.ID, app.Name, app.CreatedAt, app.UpdatedAt)

	if err != nil {
		utils.Log().Error().Err(err).Msg("Failed to insert app into DB")
		return fmt.Errorf("failed to insert app: %w", err)
	}

	return nil
}

func (r *AppRepository) GetAllApps(ctx context.Context) ([]*models.App, error) {
	query := `SELECT id, name FROM core.apps ORDER BY created_at DESC`
	apps := []*models.App{}

	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		utils.Log().Error().Err(err).Msg("Failed to query all apps")
		return nil, fmt.Errorf("failed to get all apps: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		app := &models.App{}
		if err := rows.Scan(&app.ID, &app.Name); err != nil {
			utils.Log().Error().Err(err).Msg("Failed to scan app row")
			return nil, fmt.Errorf("failed to scan app row: %w", err)
		}
		apps = append(apps, app)
	}

	return apps, nil
}

func (r *AppRepository) GetAppById(ctx context.Context, id string) (*models.App, error) {
	query := `SELECT id, name FROM core.apps WHERE id = $1`
	app := &models.App{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&app.ID, &app.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		utils.Log().Error().Err(err).Msg("Failed to query app by id")
		return nil, fmt.Errorf("failed to get app by id: %w", err)
	}

	return app, nil
}
