package main

import (
	"context"
	"flag"

	"github.com/fransiscushermanto/backend/internal/config"
	"github.com/fransiscushermanto/backend/internal/repositories"
	"github.com/fransiscushermanto/backend/internal/seeder"
	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/fransiscushermanto/backend/internal/utils"
)

func main() {
	var (
		seedType = flag.String("type", "all", "Type of seed to run (all, apps)")
	)
	flag.Parse()

	cfg, err := config.LoadConfig()

	if err != nil {
		utils.Log().Fatal().Msgf("Failed to load configuration: %v", err)
	}

	utils.SetLogLevel(cfg.LogLevel)

	db, err := repositories.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		utils.Log().Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	ctx := context.Background()

	switch *seedType {
	case "apps":
		if err := seedApps(ctx, cfg, db); err != nil {
			utils.Log().Fatal().Err(err).Msg("Failed to seed apps")
		}
	case "all":
		if err := seedAll(ctx, cfg, db); err != nil {
			utils.Log().Fatal().Err(err).Msg("Failed to seed all data")
		}
	default:
		utils.Log().Fatal().Err(err).Msg("Invalid seed type specified. Use 'apps' or 'all'.")
	}

	utils.Log().Info().Msg("Seeding completed successfully!")
}

func seedApps(ctx context.Context, cfg *config.AppConfig, db *repositories.Database) error {
	appRepo := repositories.NewAppRepository(db)
	appService := services.NewAppService(appRepo, cfg.SecretKey)
	appSeeder := seeder.NewAppSeeder(db, appService)

	if err := appSeeder.Seed(ctx); err != nil {
		return err
	}

	return nil
}

func seedAll(ctx context.Context, cfg *config.AppConfig, db *repositories.Database) error {
	// Seed apps first
	if err := seedApps(ctx, cfg, db); err != nil {
		return err
	}

	// Add other seeders here as you create them
	// if err := seedUsers(ctx, cfg, db); err != nil {
	//     return err
	// }

	return nil
}
