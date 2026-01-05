package services

import (
	"github.com/fransiscushermanto/backend/internal/config"
	"github.com/fransiscushermanto/backend/internal/services/app"
	"github.com/fransiscushermanto/backend/internal/services/auth"
	"github.com/fransiscushermanto/backend/internal/services/user"
)

type AppService = app.AppService
type AppRepository = app.AppRepository

type UserService = user.UserService
type UserRepository = user.UserRepository

type AuthService = auth.AuthService
type AuthRepository = auth.AuthRepository

func NewAppService(repo app.AppRepository, prefixApiKey string, secretKey string) *app.AppService {
	return app.NewAppService(repo, prefixApiKey, secretKey)
}

func NewUserService(repo user.UserRepository, appService *app.AppService) *user.UserService {
	return user.NewUserService(repo, appService)
}

func NewAuthService(repo auth.AuthRepository, userRepository user.UserRepository, userService *user.UserService, keys *config.CryptoKeys) *auth.AuthService {
	return auth.NewAuthService(repo, userRepository, userService, keys)
}
