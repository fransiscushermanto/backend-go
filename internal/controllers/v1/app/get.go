package app

import (
	"net/http"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/utils"
)

func (c *Controller) GetApps(w http.ResponseWriter, r *http.Request) {
	apps, err := c.appService.GetApps(r.Context())

	if err != nil {
		utils.Log().Error().Err(err).Msg("Service error getting apps")

		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    utils.StringPointer("Failed to get apps"),
		})
	}

	utils.RespondWithSuccess(w, http.StatusOK, apps, nil)
}
