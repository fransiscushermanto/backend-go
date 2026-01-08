package auth

import (
	"encoding/json"
	"net/http"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/utils"
)

func (c *Controller) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshTokenRequest

	refreshTokenLog := log("RefreshToken")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		refreshTokenLog.Error().Err(err).Msg("Missing payload")
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("Invalid request payload"),
		})
		return
	}

	refreshTokenLog.Info().Interface("body", req)
	if err := utils.ValidateBodyRequest(req); err != nil {
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("Invalid request payload"),
		})
		return
	}

	tokens, err := c.authService.RefreshToken(r.Context(), *req.RefreshToken)
	if err != nil {
		refreshTokenLog.Error().Err(err).Msg("Provided token is invalid")
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusUnauthorized,
			Message:    utils.StringPointer("Invalid or expired refresh token"),
			Meta:       &models.ErrorMeta{Code: models.CodeTokenInvalid},
		})
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, tokens, nil)
}
