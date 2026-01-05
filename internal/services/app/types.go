package app

import (
	"context"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/google/uuid"
)

type AppService struct {
	repo         AppRepository
	secretKey    string
	prefixApiKey string
}

//go:generate mockgen -source=app.go -destination=app_mock.go -package=services
type AppRepository interface {
	RegisterApp(ctx context.Context, app *models.App, appApiKey *models.AppApiKey) error
	GetAllApps(ctx context.Context) ([]*models.App, error)
	GetAppById(ctx context.Context, id uuid.UUID) (*models.App, error)
}
