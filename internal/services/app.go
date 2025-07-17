package services

import (
	"context"
	"fmt"
	"time"

	"github.com/fransiscushermanto/backend/internal/constants"
	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/google/uuid"
)

//go:generate mockgen -source=app.go -destination=app_mock.go -package=services
type AppRepository interface {
	RegisterApp(ctx context.Context, app *models.App) error
	GetAllApps(ctx context.Context) ([]*models.App, error)
	GetAppById(ctx context.Context, id string) (*models.App, error)
}

type AppService struct {
	repo      AppRepository
	secretKey string
}

func NewAppService(repo AppRepository, secretKey string) *AppService {
	return &AppService{repo: repo, secretKey: secretKey}
}

func (s *AppService) RegisterApp(ctx context.Context, req *models.RegisterAppRequest) (*models.App, error) {
	if req.Name == "" {
		return nil, utils.ErrBadRequest
	}

	name, err := utils.Encrypt([]byte(s.secretKey), []byte(req.Name))

	if err != nil {
		return nil, fmt.Errorf("failed to encrypt app name: %w", err)
	}

	app := &models.App{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now().Format(constants.TimeFormatISO),
		UpdatedAt: time.Now().Format(constants.TimeFormatISO),
	}

	opCtx, cancel := utils.ContextWithTimeout(5 * time.Second)
	defer cancel()

	if err := s.repo.RegisterApp(opCtx, app); err != nil {
		return nil, fmt.Errorf("failed to register app: %w", err)
	}

	return app, nil
}

func (s *AppService) GetApps(ctx context.Context) ([]*models.AppResponse, error) {
	optCtx, cancel := utils.ContextWithTimeout(5 * time.Second)
	defer cancel()

	rawApps, err := s.repo.GetAllApps(optCtx)

	if err != nil {
		return nil, fmt.Errorf("failed to get all apps: %w", err)
	}

	apps := make([]*models.AppResponse, len(rawApps))

	for i, rawApp := range rawApps {
		decryptedName, err := utils.Decrypt([]byte(s.secretKey), []byte(rawApp.Name))
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt app name: %w", err)
		}
		apps[i] = &models.AppResponse{
			ID:   rawApp.ID,
			Name: string(decryptedName),
		}
	}

	return apps, nil
}

func (s *AppService) GetApp(ctx context.Context, id string) (*models.App, error) {
	optCtx, cancel := utils.ContextWithTimeout(5 * time.Second)
	defer cancel()

	app, err := s.repo.GetAppById(optCtx, id)

	if err != nil {
		return nil, fmt.Errorf("failed to get app by name: %w", err)
	}

	return app, nil
}
