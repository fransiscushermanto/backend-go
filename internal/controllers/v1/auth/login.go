package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/services/auth"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/google/uuid"
)

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var loginWithEmailReq models.LoginWithEmailRequest
	var loginWithPasswordlessReq models.LoginWithPasswordlessRequest
	var loginWithOtherProviderReq models.LoginWithOtherProviderRequest
	var responseType models.AuthResponseType

	loginLog := log("Login")
	queryParams := r.URL.Query()

	params := extractAuthQueryParams(queryParams)

	if params.CookieDomain == "" && params.SetCookie {
		params.CookieDomain = "*." + strings.Split(r.Host, ":")[0]
	}

	appID, err := uuid.Parse(params.AppID)

	if err != nil {
		loginLog.Error().Str("app_id", params.AppID).Msg("Missing or Invalid app_id")
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    utils.StringPointer("app_id not found"),
		})
		return
	}

	if params.ResponseType == "" {
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("response_type is required"),
		})
		return
	}

	if isValidResponseType(params.ResponseType) {
		responseType = models.AuthResponseType(params.ResponseType)
	} else {
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("Invalid response type"),
		})
	}

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		loginLog.Error().Err(err).Msg("Failed to read the body")

		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    utils.StringPointer("Internal Server Error"),
		})
		return
	}
	defer r.Body.Close()

	errLoginWithPasswordlessReq := json.Unmarshal(bodyBytes, &loginWithPasswordlessReq)
	errLoginWithEmailReq := json.Unmarshal(bodyBytes, &loginWithEmailReq)
	errLoginWithOtherProviderReq := json.Unmarshal(bodyBytes, &loginWithOtherProviderReq)

	if errLoginWithEmailReq != nil && errLoginWithOtherProviderReq != nil && errLoginWithPasswordlessReq != nil {
		loginLog.Error().Err(errLoginWithPasswordlessReq).Err(errLoginWithEmailReq).Err(errLoginWithOtherProviderReq).Msg("Invalid Body JSON")

		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("Invalid request payload"),
		})
		return
	}

	loginResponse := &models.LoginResponse{}
	authOptions := auth.AuthOptions{
		CallbackURL: params.CallbackUrl,
		RedirectURL: params.RedirectUrl,
	}

	if loginWithEmailReq.Provider == models.AuthProviderLocal {
		loginWithEmailReq.AppID = appID

		if err := mValidator.Struct(loginWithEmailReq); err != nil {
			loginLog.Error().Err(err).Msg("Local Auth Missing Payload")
			handleValidationError(w, err)
			return
		}

		res, err := c.authService.LoginWithEmail(r.Context(), &loginWithEmailReq, authOptions)

		if err != nil {
			loginLog.Error().Err(err).Msg("Failed to login with email")

			errConfig := models.ApiError{
				StatusCode: http.StatusInternalServerError,
				Message:    utils.StringPointer("Something went wrong"),
			}

			var validationErrors utils.ValidationError

			if errors.As(err, &validationErrors) {
				errConfig.StatusCode = http.StatusUnauthorized
				errConfig.Message = nil
				errConfig.Meta = &models.ErrorMeta{
					Code: models.CodeInvalidCredentials,
				}

				fieldErrors := make(map[string]string)
				for _, fieldErr := range validationErrors.Fields {
					fieldErrors[fieldErr.Field] = fieldErr.Message
				}

				errConfig.Errors = &fieldErrors
			}

			utils.RespondWithError(w, errConfig)
			return
		}

		loginResponse = res

	} else if loginWithPasswordlessReq.Provider == models.AuthProviderPasswordless {
		loginWithPasswordlessReq.AppID = appID

		if err := mValidator.Struct(loginWithPasswordlessReq); err != nil {
			loginLog.Error().Err(err).Msg("Passwordless Auth Missing Payload")
			handleValidationError(w, err)
			return
		}

	} else if isValidOtherAuthProvider(loginWithOtherProviderReq.Provider) {
		loginWithOtherProviderReq.AppID = appID

		// other than local and passwordless it's asume as other provider
		if err := mValidator.Struct(loginWithOtherProviderReq); err != nil {
			loginLog.Error().Err(err).Msg("Other Provider Auth  Missing Payload")
			handleValidationError(w, err)
			return
		}

	} else {
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("Invalid request payload"),
		})
		return
	}

	switch {
	case responseType == models.AuthResponseRedirect:
	case responseType == models.AuthResponseCallback:
		// Either callbackUrl is provided or responseType is redirect
		handleCallbackOrRedirectResponse(w, r, responseType, loginResponse)

	default:
		handleJSONResponse(w, loginResponse)
	}
}
