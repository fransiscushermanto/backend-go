package main

import (
	"log"

	"github.com/fransiscushermanto/backend/internal/config"
	"github.com/fransiscushermanto/backend/internal/repositories"
	"github.com/fransiscushermanto/backend/internal/server"
	"github.com/fransiscushermanto/backend/internal/utils"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	utils.SetLogLevel(cfg.LogLevel)

	db, err := repositories.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		utils.Log().Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			utils.Log().Error().Err(closeErr).Msg("Failed to close database connection")
		} else {
			utils.Log().Info().Msg("Database connection closed")
		}
	}()

	apiServer := server.NewAPIServer(cfg, db)
	if err := apiServer.Run(); err != nil {
		utils.Log().Fatal().Err(err).Msg("API server stopped with error")
	}

	utils.Log().Info().Msg("API server shutdown complete")
}
