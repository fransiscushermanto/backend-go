package app

import (
	"context"
	"fmt"
	"time"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/google/uuid"
)

func (s *AppService) GetApp(ctx context.Context, id string) (*models.App, error) {
	optCtx, cancel := utils.ContextWithTimeout(5 * time.Second)
	defer cancel()

	parsedID, err := uuid.Parse(id)

	if err != nil {
		return nil, fmt.Errorf("failed to get app by name: %w", err)
	}

	app, err := s.repo.GetAppById(optCtx, parsedID)

	if err != nil {
		return nil, fmt.Errorf("failed to get app by name: %w", err)
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
		decryptedName, err := s.ParseAppName(string(rawApp.Name))
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt app name: %w", err)
		}
		apps[i] = &models.AppResponse{
			ID:   rawApp.ID.String(),
			Name: string(decryptedName),
		}
	}

	return apps, nil
}
