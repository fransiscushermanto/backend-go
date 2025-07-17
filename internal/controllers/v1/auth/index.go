package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

type Controller struct {
	authService *services.AuthService
}

func NewController(authService *services.AuthService) *Controller {
	return &Controller{
		authService: authService,
	}
}

var mValidator *validator.Validate = InitValidator()

func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	var responseType models.AuthResponseType

	queryParams := r.URL.Query()
	cookieDomain := queryParams.Get("cookie_domain")
	callbackUrl := queryParams.Get("callback_url")
	isSetCookie := queryParams.Get("set_cookie") == "true"
	appID := queryParams.Get("app_id")

	if cookieDomain == "" && isSetCookie {
		// get host without port
		cookieDomain = "*." + strings.Split(r.Host, ":")[0]
	}

	if queryParams.Has("response_type") && queryParams.Get("response_type") != "" {
		queryResponseType := queryParams.Get("response_type")

		if queryResponseType == string(models.AuthResponseJSON) || queryResponseType == string(models.AuthResponseRedirect) {
			responseType = models.AuthResponseType(queryResponseType)
		} else {
			utils.RespondWithError(w, utils.ErrorResponsePayload{
				StatusCode: http.StatusBadRequest,
				Message:    utils.StringPointer("Invalid response type"),
			})
			return
		}
	}

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

	registerResponse, err := c.authService.Register(r.Context(), &req, callbackUrl)
	if err != nil {
		utils.Log().Error().Err(err).Msg("Service error registering user")

		errConfig := utils.ErrorResponsePayload{
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

	if isSetCookie {
		c.setAuthCookies(w, registerResponse, cookieDomain)
	}

	switch {
	case callbackUrl != "" && responseType != models.AuthResponseJSON:
		// Either callbackUrl is provided or responseType is redirect
		c.handleRedirectResponse(w, r, registerResponse)
	default:
		c.handleJSONResponse(w, registerResponse)
	}

}

func (c *Controller) setAuthCookies(w http.ResponseWriter, response *models.RegisterResponse, domain string) {
	// Set access token cookie
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    response.AccessToken,
		Path:     "/",
		Domain:   domain,
		MaxAge:   15 * 60, // 15 minutes
		HttpOnly: true,
		Secure:   true, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
	}

	// Set refresh token cookie
	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    response.RefreshToken,
		Path:     "/",
		Domain:   domain,
		MaxAge:   30 * 24 * 60 * 60, // 30 days
		HttpOnly: true,
		Secure:   true, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)
}

// Handle redirect response
func (c *Controller) handleRedirectResponse(w http.ResponseWriter, r *http.Request, response *models.RegisterResponse) {
	if response.RedirectURL != "" {
		http.Redirect(w, r, response.RedirectURL, http.StatusFound)
		return
	}

	// Fallback to JSON if no redirect URL
	c.handleJSONResponse(w, response)
}

// Handle JSON response
func (c *Controller) handleJSONResponse(w http.ResponseWriter, response *models.RegisterResponse) {
	// Don't expose tokens in JSON if cookies are set
	utils.RespondWithSuccess(w, http.StatusCreated, response, nil)
}
