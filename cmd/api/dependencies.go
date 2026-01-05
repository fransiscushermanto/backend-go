package main

import (
	"github.com/fransiscushermanto/backend/internal/config"
	"github.com/fransiscushermanto/backend/internal/repositories"
	"github.com/fransiscushermanto/backend/internal/server/routes"
	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/fransiscushermanto/backend/internal/utils"
)

func newServiceContainer(cfg *config.AppConfig, db *utils.Database, keys *config.CryptoKeys) *routes.Services {
	// Repositories
	appRepo := repositories.NewAppRepository(db, &cfg.LockTimeout)
	userRepo := repositories.NewUserRepository(db)
	authRepo := repositories.NewAuthRepository(db)

	// Services
	appService := services.NewAppService(appRepo, cfg.PrefixApiKey, cfg.SecretKey)
	userService := services.NewUserService(userRepo, appService)
	authService := services.NewAuthService(authRepo, userRepo, userService, keys)

	return &routes.Services{
		AppService:  appService,
		UserService: userService,
		AuthService: authService,
	}
}
