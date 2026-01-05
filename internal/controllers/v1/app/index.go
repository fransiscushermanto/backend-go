package app

import (
	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/go-playground/validator/v10"
)

type ControllerOptions struct {
	SecretKey string
}

type Controller struct {
	appService *services.AppService
	options    ControllerOptions
}

func NewController(appService *services.AppService, options ControllerOptions) *Controller {
	return &Controller{
		appService: appService,
		options:    options,
	}
}

var mValidator *validator.Validate = InitValidator()
