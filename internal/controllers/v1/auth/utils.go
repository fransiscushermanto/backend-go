package auth

import (
	"net/http"
	"net/url"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

func isValidResponseType(queryResponseType string) bool {
	return queryResponseType == string(models.AuthResponseJSON) || queryResponseType == string(models.AuthResponseRedirect) || queryResponseType == string(models.AuthResponseCallback)
}

type dependentValues struct {
	callbackURL *string
	redirectURL *string
}

type AuthQueryParams struct {
	CookieDomain string
	CallbackUrl  string
	RedirectUrl  string
	SetCookie    bool
	AppID        string
	ResponseType string
}

func extractAuthQueryParams(queryParams url.Values) *AuthQueryParams {

	return &AuthQueryParams{
		AppID:        queryParams.Get("app_id"),
		CookieDomain: queryParams.Get("cookie_domain"),
		CallbackUrl:  queryParams.Get("callback_url"),
		RedirectUrl:  queryParams.Get("redirect_url"),
		SetCookie:    queryParams.Get("set_cookie") == "true",
		ResponseType: queryParams.Get("response_type"),
	}
}

func verifyResponseTypeDependentValueExist(responseType models.AuthResponseType, dependentValues dependentValues) (bool, *string) {
	if responseType == models.AuthResponseCallback && *dependentValues.callbackURL == "" {
		return false, utils.StringPointer("callback_url is required when response_type = callback")
	}

	if responseType == models.AuthResponseRedirect && *dependentValues.redirectURL == "" {
		return false, utils.StringPointer("redirect_url is required when response_type = redirect")
	}

	return true, nil
}

func handleValidationError(w http.ResponseWriter, err error) {
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
}

func setAuthCookies(w http.ResponseWriter, response *models.RegisterResponse, domain string) {
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
func handleCallbackOrRedirectResponse(w http.ResponseWriter, r *http.Request, responseType models.AuthResponseType, response any) {
	var url string

	switch responseType {
	case models.AuthResponseRedirect:
		if resp, ok := response.(*struct {
			RedirectURL string
		}); ok {
			url = resp.RedirectURL
		}
	case models.AuthResponseCallback:
		if resp, ok := response.(*struct {
			CallbackURL string
		}); ok {
			url = resp.CallbackURL
		}
	}

	http.Redirect(w, r, url, http.StatusFound)
}

// Handle JSON response
func handleJSONResponse(w http.ResponseWriter, response any) {
	// Don't expose tokens in JSON if cookies are set
	utils.RespondWithSuccess(w, http.StatusCreated, response, nil)
}

func isValidOtherAuthProvider(provider models.AuthProvider) bool {
	switch provider {
	case models.AuthProviderGoogle:
		return true
	default:
		return false
	}

}
