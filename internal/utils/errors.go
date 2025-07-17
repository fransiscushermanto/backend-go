package utils

import (
	"errors"
	"net/http"

	"github.com/fransiscushermanto/backend/internal/models"
)

var (
	ErrNotFound             = errors.New("not found")
	ErrBadRequest           = errors.New("bad request")
	ErrInternalServerError  = errors.New("internal server error")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden")
	ErrValidation           = errors.New("validation error")
	ErrConflict             = errors.New("conflict")
	ErrNotImplemented       = errors.New("not implemented")
	ErrServiceUnavailable   = errors.New("service unavailable")
	ErrTooManyRequests      = errors.New("too many requests")
	ErrMethodNotAllowed     = errors.New("method not allowed")
	ErrUnsupportedMediaType = errors.New("unsupported media type")
)

type ErrorResponsePayload struct {
	StatusCode int                     // HTTP status code (e.g., 400, 404, 500)
	Message    *string                 // Primary human-readable message for the API client
	Errors     *map[string]string      // Optional: detailed field-specific validation errors
	Data       *interface{}            // Optional: any data to include with the error response
	Meta       *map[string]interface{} // Optional: additional metadata for the error
}

func RespondWithError(w http.ResponseWriter, payload ErrorResponsePayload) {
	apiErr := models.ApiError{
		Message:    payload.Message,
		StatusCode: payload.StatusCode,
		Errors:     payload.Errors,
		Data:       payload.Data,
		Meta:       payload.Meta,
	}

	RespondWithJSON(w, payload.StatusCode, apiErr)
}

func RespondWithValidationError(w http.ResponseWriter, errors map[string]string, data *interface{}, meta *map[string]interface{}) {
	payload := ErrorResponsePayload{
		StatusCode: http.StatusUnprocessableEntity,
		Errors:     &errors,
		Data:       data,
		Meta:       meta,
	}

	RespondWithError(w, payload)
}
