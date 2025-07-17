package app

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/fransiscushermanto/backend/internal/utils"
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

var mValidator *validator.Validate

func (c *Controller) RegisterApp(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterAppRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, utils.ErrorResponsePayload{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("Invalid request payload"),
		})
		return
	}

	if err := mValidator.Struct(req); err != nil {
		validatorErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			utils.RespondWithError(w, utils.ErrorResponsePayload{
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		formattedErrors := make(map[string]string)
		for _, fieldErr := range validatorErrors {
			formattedErrors[fieldErr.Field()] = utils.GetValidationErrorMessage(fieldErr)
		}

		utils.RespondWithValidationError(w, formattedErrors, nil, nil)
	}

	app, err := c.appService.RegisterApp(r.Context(), &req)
	if err != nil {
		utils.Log().Error().Err(err).Msg("Service error registering app")
		errConfig := utils.ErrorResponsePayload{
			StatusCode: http.StatusInternalServerError,
			Message:    utils.StringPointer("Failed to register app"),
		}

		if errors.Is(err, utils.ErrBadRequest) {
			errConfig.StatusCode = http.StatusBadRequest
			errConfig.Message = utils.StringPointer("Invalid app data")
		}

		utils.RespondWithError(w, errConfig)
		return
	}

	utils.RespondWithSuccess(w, http.StatusCreated, app, nil)
}

func (c *Controller) GetApps(w http.ResponseWriter, r *http.Request) {
	apps, err := c.appService.GetApps(r.Context())

	if err != nil {
		utils.Log().Error().Err(err).Msg("Service error getting apps")

		utils.RespondWithError(w, utils.ErrorResponsePayload{
			StatusCode: http.StatusInternalServerError,
			Message:    utils.StringPointer("Failed to get apps"),
		})
	}

	utils.RespondWithSuccess(w, http.StatusOK, apps, nil)
}
