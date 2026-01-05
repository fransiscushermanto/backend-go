package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/fransiscushermanto/backend/internal/config"
	"github.com/fransiscushermanto/backend/internal/repositories"
	"github.com/fransiscushermanto/backend/internal/seeder"
	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/fransiscushermanto/backend/internal/utils"
)

func main() {
	var (
		seedType = flag.String("type", "all", "Type of seed to run (all, apps, apps_api_key)")
	)
	flag.Parse()

	cfg, err := config.LoadConfig()

	if err != nil {
		utils.Log().Fatal().Msgf("Failed to load configuration: %v", err)
	}

	utils.SetLogLevel(cfg.LogLevel)

	db, err := utils.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		utils.Log().Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	ctx := context.Background()

	appSeeder := initAppSeeder(cfg, db)

	switch *seedType {
	case "apps":
		fmt.Println("Start seeding apps")
		if err := appSeeder.Seed(ctx); err != nil {
			utils.Log().Fatal().Err(err).Msg("Failed to seed apps")
		}
	case "apps_api_key":
		fmt.Println("Start seeding apps_api_key")
		if err = appSeeder.SeedApiKeyHash(ctx); err != nil {
			utils.Log().Fatal().Err(err).Msg("Failed to seed apps api key")
		}
	case "all":
		fmt.Println("Start seeding all data")
		if err := seedAll(ctx, appSeeder); err != nil {
			utils.Log().Fatal().Err(err).Msg("Failed to seed all data")
		}
	default:
		utils.Log().Fatal().Err(err).Msg("Invalid seed type specified. Use 'apps' or 'all'.")
	}

	utils.Log().Info().Msg("Seeding completed successfully!")
}

func initAppSeeder(cfg *config.AppConfig, db *utils.Database) *seeder.AppSeeder {
	appRepo := repositories.NewAppRepository(db, &cfg.LockTimeout)
	appService := services.NewAppService(appRepo, cfg.PrefixApiKey, cfg.SecretKey)
	return seeder.NewAppSeeder(db, appService)
}

func seedAll(ctx context.Context, appSeeder *seeder.AppSeeder) error {
	// Seed apps first
	if err := appSeeder.Seed(ctx); err != nil {
		return err
	}

	// Add other seeders here as you create them
	// if err := seedUsers(ctx, cfg, db); err != nil {
	//     return err
	// }

	return nil
}
