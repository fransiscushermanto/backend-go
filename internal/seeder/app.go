package seeder

import (
	"context"
	"time"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AppSeeder struct {
	db         *utils.Database
	appService *services.AppService
}

func NewAppSeeder(db *utils.Database, appService *services.AppService) *AppSeeder {
	return &AppSeeder{
		db:         db,
		appService: appService,
	}
}

var apps = []models.RegisterAppRequest{
	{Name: "fransiscushermanto"},
	{Name: "bloomify-and-co"},
}

func (s *AppSeeder) Seed(ctx context.Context) error {
	for _, app := range apps {
		registeredApp, err := s.appService.Register(ctx, &app)
		if err != nil {
			// Log the error and continue, so one failure doesn't stop all seeding
			utils.Log().Error().Err(err).Str("app_name", app.Name).Msg("Failed to seed app")
			continue
		}
		utils.Log().Info().Str("app_name", registeredApp.Name).Str("app_id", registeredApp.ID).Msg("App seeded successfully")
	}

	return nil
}

func (s *AppSeeder) SeedApiKeyHash(ctx context.Context) error {
	query := `SELECT a.id, name FROM core.apps a LEFT JOIN core.app_api_keys aak ON a.id = aak.app_id WHERE aak.app_id IS NULL`

	rows, err := s.db.Pool.Query(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	utils.Log().Println("Starting to backfill API keys...")

	for rows.Next() {
		var appIDStr string
		var rawAppName []byte

		if err := rows.Scan(&appIDStr, &rawAppName); err != nil {
			utils.Log().Error().Err(err).Msg("Failed to scan appID")
			continue
		}

		apiKey, err := s.appService.GenerateAPIKey()
		if err != nil {
			utils.Log().Error().Err(err).Msg("Failed to generate API key for app")
			continue
		}

		hashedApiKey, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
		if err != nil {
			utils.Log().Error().Err(err).Msg("Failed to hash API key for app")
			continue
		}

		appID, appIDErr := uuid.Parse(appIDStr)
		if appIDErr != nil {
			utils.Log().Error().Err(err).Msg("Failed to parse appIDStr")
			continue
		}

		appApiKeyID, appApiKeyIDerr := uuid.NewV7()
		if appApiKeyIDerr != nil {
			utils.Log().Error().Err(err).Msg("Failed to generate appApiKeyID")
			continue
		}

		appApiKey := &models.AppApiKey{
			ID:        appApiKeyID,
			AppID:     appID,
			KeyHash:   string(hashedApiKey),
			CreatedAt: time.Now(),
			IsActive:  true,
		}

		_, err = s.db.Pool.Exec(ctx, `INSERT INTO core.app_api_keys (id, app_id, key_hash, created_at, is_active) VALUES ($1, $2, $3, $4, $5)`, appApiKey.ID, appApiKey.AppID, appApiKey.KeyHash, appApiKey.CreatedAt, appApiKey.IsActive)
		if err != nil {
			utils.Log().Error().Err(err).Msg("Failed to update app")
			continue
		}

		appName, err := s.appService.ParseAppName(string(rawAppName))

		if err != nil {
			utils.Log().Error().Err(err).Msg("Failed to parsed appName")
			continue
		}

		utils.Log().Info().Str("app_name", appName).Str("api_key", apiKey).Msg("✅ Backfilled app with NEW API Key")
	}

	utils.Log().Info().Msg("Backfill complete. SAVE THE LOG OUTPUT ABOVE!")
	return nil
}
