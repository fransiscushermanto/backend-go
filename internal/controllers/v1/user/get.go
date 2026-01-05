package user

import (
	"errors"
	"net/http"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (c *Controller) GetUsers(w http.ResponseWriter, r *http.Request) {
	getUsersLog := log("GetUsers")
	queryParams := r.URL.Query()
	strAppID := queryParams.Get("app_id")

	var users []*models.User
	var err error

	if strAppID != "" {
		appID, errAppID := uuid.Parse(strAppID)

		if errAppID != nil {
			getUsersLog.Error().Err(err).Msg("Invalid appID")
			empty := []*models.UserResponse{}
			utils.RespondWithSuccess(w, http.StatusOK, empty, nil)
			return
		}

		users, err = c.userService.GetAppUsers(r.Context(), appID)
	} else {
		users, err = c.userService.GetUsers(r.Context())
	}

	userResponse := make([]*models.UserResponse, len(users))

	if err != nil {
		getUsersLog.Error().Err(err).Msg("Service error getting users")
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    utils.StringPointer("Failed to retrieve users"),
		})
		return
	}

	for i, user := range users {
		if user == nil {
			getUsersLog.Warn().Msgf("User at index %d is nil", i)
			continue
		}

		userResponse[i] = user.ToResponse()
	}

	utils.RespondWithSuccess(w, http.StatusOK, userResponse, nil)
}

func (c *Controller) GetUser(w http.ResponseWriter, r *http.Request) {
	getUserLog := log("GetUser")

	queryParams := r.URL.Query()
	strAppID := queryParams.Get("app_id")
	strUserID := chi.URLParam(r, "id")

	if strAppID == "" || strUserID == "" {
		getUserLog.Error().Str("appID", strAppID).Str("userID", strUserID).Msg("Missing userID or appID")
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("User ID or App ID is required"),
		})
		return
	}

	appID, errAppID := uuid.Parse(strAppID)
	userID, errUserID := uuid.Parse(strUserID)

	if errAppID != nil || errUserID != nil {
		getUserLog.Error().Str("appID", strAppID).Str("userID", strUserID).Msg("Invalid userID or appID")
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("User ID or App ID is invalid"),
		})
		return
	}

	user, err := c.userService.GetUser(r.Context(), appID, userID)
	if err != nil {
		getUserLog.Error().Err(err).Msg("Service error getting user")
		errConfig := models.ApiError{
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
