package seeder

import (
	"context"
	"fmt"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/repositories"
	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/fransiscushermanto/backend/internal/utils"
)

type AppSeeder struct {
	db         *repositories.Database
	appService *services.AppService
}

func NewAppSeeder(db *repositories.Database, appService *services.AppService) *AppSeeder {
	return &AppSeeder{
		db:         db,
		appService: appService,
	}
}

func (s *AppSeeder) Seed(ctx context.Context) error {
	apps := []models.RegisterAppRequest{
		{Name: "fransiscushermanto"},
		{Name: "bloomify-and-co"},
	}

	for _, app := range apps {
		registeredApp, err := s.appService.RegisterApp(ctx, &app)

		if err != nil {
			return fmt.Errorf("failed to seed app %s: %w", app.Name, err)
		}
		utils.Log().Info().Msgf("App seeded: %v\n", registeredApp)
	}

	return nil
}
