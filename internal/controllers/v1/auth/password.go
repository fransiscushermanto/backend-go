package auth

import (
	"encoding/json"
	"net/http"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/google/uuid"
)

func (c *Controller) ForgetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ForgetPasswordRequest

	forgetPasswordLog := log("ForgetPassword")
	queryParams := r.URL.Query()
	params := extractAuthQueryParams(queryParams)

	appID, err := uuid.Parse(params.AppID)
	if err != nil {
		forgetPasswordLog.Error().Err(err).Msg("Missing or Invalid app_id")
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusForbidden,
			Message:    utils.StringPointer("Forbidden"),
		})
		return
	}
	req.AppID = &appID

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		forgetPasswordLog.Error().Err(err).Msg("Invalid JSON")
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("Invalid request payload"),
		})
		return
	}

	if err := utils.ValidateBodyRequest(req); err != nil {
		forgetPasswordLog.Error().Err(err).Msg("Missing required key payload")
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("Invalid request payload"),
		})
		return
	}

	if err := mValidator.Struct(req); err != nil {
		forgetPasswordLog.Error().Err(err).Msg("Validation error")
		handleValidationError(w, err)
		return
	}

	err = c.authService.ForgetPassword(r.Context(), &req)

	if err != nil {
		forgetPasswordLog.Error().Err(err).Msg("Failed to request forget password")
	}

	utils.RespondWithSuccess(w, http.StatusNoContent, nil, nil)
}
