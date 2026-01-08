package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/fransiscushermanto/backend/internal/models"
	authService "github.com/fransiscushermanto/backend/internal/services/auth"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/google/uuid"
)

func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	var responseType models.AuthResponseType

	registerLog := log("Register")

	queryParams := r.URL.Query()
	params := extractAuthQueryParams(queryParams)

	if params.CookieDomain == "" && params.SetCookie {
		// get host without port
		params.CookieDomain = "*." + strings.Split(r.Host, ":")[0]
	}

	appID, err := uuid.Parse(params.AppID)

	if err != nil {
		registerLog.Error().Err(err).Msg("Missing or Invalid app_id")
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    utils.StringPointer("app_id not found"),
		})
		return
	}

	req.AppID = appID

	if params.ResponseType == "" {
		registerLog.Error().Msg("Missing response_type")
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("response_type is required"),
		})
		return
	}

	if isValidResponseType(params.ResponseType) {
		responseType = models.AuthResponseType(params.ResponseType)
	} else {
		registerLog.Error().Str("response_type", params.ResponseType).Msg("Invalid Response Type")
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("Invalid response type"),
		})
		return
	}

	if isValid, message := verifyResponseTypeDependentValueExist(responseType, dependentValues{
		callbackURL: &params.CallbackUrl,
		redirectURL: &params.RedirectUrl,
	}); !isValid {
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    message,
		})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		registerLog.Error().Err(err).Msg("Invalid JSON")
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("Invalid Payload Request"),
		})
		return
	}

	// TODO: remove this check when implemented code for other providers
	if req.Provider != models.AuthProviderLocal {
		registerLog.Error().Str("provider", string(req.Provider)).Msg("Accessing not implemented provider")
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusNotImplemented,
		})
		return
	}

	if err := utils.ValidateBodyRequest(req); err != nil {
		registerLog.Error().Err(err).Msg("Invalid Payload Request")
		utils.RespondWithError(w, models.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    utils.StringPointer("Invalid Payload Request"),
		})
		return
	}

	if err := mValidator.Struct(req); err != nil {
		registerLog.Error().Err(err).Msg("Validation error")
		handleValidationError(w, err)
		return
	}

	registerResponse, err := c.authService.Register(r.Context(), &req, authService.AuthOptions{
		CallbackURL: params.CallbackUrl,
		RedirectURL: params.RedirectUrl,
	})

	if err != nil {
		registerLog.Error().Err(err).Msg("Failed to register user")

		if responseType == models.AuthResponseCallback || responseType == models.AuthResponseRedirect {
			// Either callbackUrl is provided or responseType is redirect
			handleCallbackOrRedirectResponse(w, r, responseType, registerResponse)
			return
		}

		errConfig := models.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    utils.StringPointer("Failed to create user"),
		}

		if errors.Is(err, utils.ErrBadRequest) {
			errConfig.StatusCode = http.StatusBadRequest
			errConfig.Message = utils.StringPointer("Please provide valid user data")
		}

		// If you had specific data or meta to include with this error, you'd add it to errConfig here
		utils.RespondWithError(w, errConfig)
		return
	}

	if params.SetCookie {
		setAuthCookies(w, registerResponse, params.CookieDomain)
	}

	switch {
	case responseType == models.AuthResponseRedirect:
	case responseType == models.AuthResponseCallback:
		// Either callbackUrl is provided or responseType is redirect
		handleCallbackOrRedirectResponse(w, r, responseType, registerResponse)

	default:
		handleJSONResponse(w, registerResponse)
	}

}
