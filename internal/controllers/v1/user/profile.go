package user

import (
	"errors"
	"net/http"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/services/user"
	"github.com/fransiscushermanto/backend/internal/utils"
)

func (c *Controller) Profile(w http.ResponseWriter, r *http.Request) {
	utils.DebugContextValue(r.Context(), utils.UserIDContextKey)
	userID, errUserID := utils.GetUserIDFromContext(r.Context())
	appID, errAppID := utils.GetAppIDFromContext(r.Context())

	if errUserID != nil || errAppID != nil {
		utils.Log().Error().Err(errUserID).Err(errAppID).Msg("Context missing user_id")
		utils.RespondWithError(w, models.ApiError{
			StatusCode: 500,
			Message:    utils.StringPointer("Internal server error"),
		})
		return
	}

	user, err := c.userService.GetUser(r.Context(), *appID, user.UserIdentifier{
		ID: userID,
	})

	if err != nil {
		utils.Log().Error().Err(err).Msg("Service error getting user")
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
