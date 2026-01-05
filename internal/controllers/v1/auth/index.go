package auth

import (
	authService "github.com/fransiscushermanto/backend/internal/services/auth"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type Controller struct {
	authService *authService.AuthService
}

func NewController(authService *authService.AuthService) *Controller {
	return &Controller{
		authService: authService,
	}
}

func log(method string) *zerolog.Logger {
	l := utils.Log().With().Str("controller", "Auth").Str("method", method).Logger()
	return &l
}

var mValidator *validator.Validate = InitValidator()
