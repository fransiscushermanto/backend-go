package controllers

import (
	appController "github.com/fransiscushermanto/backend/internal/controllers/v1/app"
	authController "github.com/fransiscushermanto/backend/internal/controllers/v1/auth"
	userController "github.com/fransiscushermanto/backend/internal/controllers/v1/user"
	"github.com/fransiscushermanto/backend/internal/services"
)

func NewUserController(userService *services.UserService) *userController.Controller {
	return userController.NewController(userService)
}

func NewAppController(appService *services.AppService, options appController.ControllerOptions) *appController.Controller {
	return appController.NewController(appService, options)
}

func NewAuthController(authService *services.AuthService) *authController.Controller {
	return authController.NewController(authService)
}
