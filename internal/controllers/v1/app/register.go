package app

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/go-playground/validator/v10"
)

func (c *Controller) RegisterApp(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterAppRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("Invalid request payload"),
		})
		return
	}

	if err := mValidator.Struct(req); err != nil {
		validatorErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			utils.RespondWithError(w, models.ApiError{
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

	app, err := c.appService.Register(r.Context(), &req)
	if err != nil {
		utils.Log().Error().Err(err).Msg("Service error registering app")
		errConfig := models.ApiError{
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
