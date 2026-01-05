package user

import (
	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type Controller struct {
	userService *services.UserService
}

func NewController(userService *services.UserService) *Controller {
	return &Controller{
		userService: userService,
	}
}

var mValidator *validator.Validate = InitValidator()

func log(method string) *zerolog.Logger {
	l := utils.Log().With().Str("controller", "User").Str("method", method).Logger()
	return &l
}
