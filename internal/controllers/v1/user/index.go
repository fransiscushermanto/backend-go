package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
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

func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest

	queryParams := r.URL.Query()
	appID := queryParams.Get("app_id")

	if appID == "" {
		utils.RespondWithError(w, utils.ErrorResponsePayload{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("App ID is required"),
		})
		return
	} else {
		req.AppID = appID
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, utils.ErrorResponsePayload{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("Invalid request payload"),
		})
		return
	}

	// TODO: remove this check when implemented code for other providers
	if req.Provider != models.AuthProviderLocal {
		utils.RespondWithError(w, utils.ErrorResponsePayload{
			StatusCode: http.StatusNotImplemented,
		})
		return
	}

	if err := mValidator.Struct(req); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			utils.RespondWithError(w, utils.ErrorResponsePayload{
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
		errConfig := utils.ErrorResponsePayload{
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

func (c *Controller) GetUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	appID := queryParams.Get("app_id")

	var users []*models.User
	var err error

	if appID != "" {
		users, err = c.userService.GetAppUsers(r.Context(), appID)
	} else {
		users, err = c.userService.GetUsers(r.Context())
	}

	userResponse := make([]*models.UserResponse, len(users))

	if err != nil {
		utils.Log().Error().Err(err).Msg("Service error getting users")
		utils.RespondWithError(w, utils.ErrorResponsePayload{
			StatusCode: http.StatusInternalServerError,
			Message:    utils.StringPointer("Failed to retrieve users"),
		})
		return
	}

	for i, user := range users {
		if user == nil {
			utils.Log().Warn().Msgf("User at index %d is nil", i)
			continue
		}

		userResponse[i] = user.ToResponse()
	}

	utils.RespondWithSuccess(w, http.StatusOK, userResponse, nil)
}

func (c *Controller) GetUser(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	appID := queryParams.Get("app_id")
	userID := chi.URLParam(r, "id")

	if userID == "" || appID == "" {
		utils.RespondWithError(w, utils.ErrorResponsePayload{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("User ID or App ID is required"),
		})
		return
	}

	user, err := c.userService.GetUser(r.Context(), appID, userID)
	if err != nil {
		utils.Log().Error().Err(err).Msg("Service error getting user")
		errConfig := utils.ErrorResponsePayload{
			StatusCode: http.StatusInternalServerError,
			Message:    utils.StringPointer("Failed to retrieve user"),
		}
		if errors.Is(err, utils.ErrNotFound) {
			errConfig.StatusCode = http.StatusNotFound
			errConfig.Message = utils.StringPointer("User not found")
		}
		utils.RespondWithError(w, errConfig)
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, user.ToResponse(), nil)
}
