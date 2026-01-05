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
	ErrUnprocessableEntity  = errors.New("unprocessable entity")
)

func RespondWithError(w http.ResponseWriter, payload models.ApiError) {
	if payload.StatusCode == 401 {
		if payload.Meta == nil {
			payload.Meta = &models.ErrorMeta{
				Code: models.CodeUnauthorized,
			}
		} else if payload.Meta.Code == "" {
			payload.Meta.Code = models.CodeUnauthorized
		}

	}

	apiErr := models.ApiError{
		Message:    payload.Message,
		StatusCode: payload.StatusCode,
		Errors:     payload.Errors,
		Data:       payload.Data,
		Meta:       payload.Meta,
	}

	RespondWithJSON(w, payload.StatusCode, apiErr)
}

func RespondWithValidationError(w http.ResponseWriter, errors map[string]string, data *interface{}, meta *models.ErrorMeta) {
	payload := models.ApiError{
		StatusCode: http.StatusUnprocessableEntity,
		Errors:     &errors,
		Data:       data,
		Meta:       meta,
	}

	RespondWithError(w, payload)
}
