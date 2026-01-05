package main

import (
	"log"
	"sync"

	"github.com/fransiscushermanto/backend/internal/config"
	"github.com/fransiscushermanto/backend/internal/server"
	"github.com/fransiscushermanto/backend/internal/utils"
)

func main() {
	var wg sync.WaitGroup
	var cfg *config.AppConfig
	var keys *config.CryptoKeys
	var cfgErr, keyErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		cfg, cfgErr = config.LoadConfig()
	}()

	go func() {
		defer wg.Done()
		keys, keyErr = config.LoadCryptoKeys()
	}()

	wg.Wait()

	if cfgErr != nil {
		log.Fatalf("Failed to load configuration: %v", cfgErr)
	}

	if keyErr != nil {
		log.Fatalf("Failed to load crypto keys: %v", keyErr)
	}

	utils.SetLogLevel(cfg.LogLevel)
	utils.SetEnv(cfg.Env)

	utils.Log().Info().
		Str("env", cfg.Env).
		Str("log_level", cfg.LogLevel).
		Int("port", cfg.Port).
		Str("key_type", keys.PublicKeyInfo()).
		Bool("keys_valid", keys.IsValid()).
		Msg("Application configuration loaded")

	db, err := utils.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		utils.Log().Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer func() {
		db.Close()
		utils.Log().Info().Msg("Database connection closed")
	}()

	services := newServiceContainer(cfg, db, keys)
	apiServer := server.NewAPIServer(cfg, services, keys)

	if err := apiServer.Run(); err != nil {
		utils.Log().Fatal().Err(err).Msg("API server stopped with error")
	}

	utils.Log().Info().Msg("API server shutdown complete")
}
