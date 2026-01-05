package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
)

func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest

	queryParams := r.URL.Query()
	appID := queryParams.Get("app_id")

	parsedAppID, err := uuid.Parse(appID)

	if err != nil {
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("App ID is required"),
		})
		return
	} else {
		req.AppID = parsedAppID
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("Invalid request payload"),
		})
		return
	}

	// TODO: remove this check when implemented code for other providers
	if req.Provider != models.AuthProviderLocal {
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusNotImplemented,
		})
		return
	}

	if err := mValidator.Struct(req); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			utils.RespondWithError(w, models.ApiError{
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		formattedErrors := make(map[string]string)
		for _, fieldErr := range validationErrors {
			formattedErrors[strcase.ToSnake(fieldErr.Field())] = utils.GetValidationErrorMessage(fieldErr, RenderErrorMessage)
		}

		utils.RespondWithValidationError(w, formattedErrors, nil, nil)
		return
	}

	user, err := c.userService.CreateUser(r.Context(), &req)
	if err != nil {
		utils.Log().Error().Err(err).Msg("Service error creating user")
		errConfig := models.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    utils.StringPointer("Failed to create user"),
		}
		if errors.Is(err, utils.ErrBadRequest) {
			errConfig.StatusCode = http.StatusBadRequest
			errConfig.Message = utils.StringPointer("Invalid user data provided")
		}
		// If you had specific data or meta to include with this error, you'd add it to errConfig here
		utils.RespondWithError(w, errConfig)
		return
	}

	utils.RespondWithSuccess(w, http.StatusCreated, user.ToResponse(), nil)
}
