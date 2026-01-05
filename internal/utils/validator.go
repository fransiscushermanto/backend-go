package utils

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
)

type RenderErrorMessageFunc func(fieldError validator.FieldError) (string, bool)

var errorMessages = map[string]string{
	"required": "%s is a required field.",
	"email":    "%s is not a valid email address.",
	"min":      "%s must be at least %s characters long.", // %s for field, %s for param (e.g., min=3)
	"max":      "%s must be at most %s characters long.",
	"gte":      "%s must be %s or greater.",
	"lte":      "%s must be %s or less.",
	"eqfield":  "%s must match the %s field.",
}

type FieldError struct {
	Field   string
	Message string
}

type ValidationError struct {
	Fields []FieldError
}

func (v ValidationError) Error() string {
	return "validation failed"
}

func NewValidationError(fields []FieldError) error {
	return ValidationError{Fields: fields}
}

func GetValidationErrorMessage(fieldError validator.FieldError, renderErrorMessages ...RenderErrorMessageFunc) string {
	fieldName := fieldError.Field() // Gets the field name (e.g., "name", "password")
	tag := fieldError.Tag()         // Gets the validation tag (e.g., "required")
	param := fieldError.Param()     // Gets any parameters associated with the tag (e.g., "8" for min=8)

	if len(renderErrorMessages) > 0 && renderErrorMessages[0] != nil {
		customMessage, handled := renderErrorMessages[0](fieldError)
		if handled {
			return customMessage
		}
	}

	if strings.Contains(tag, "_") {
		tag = strings.Split(tag, "_")[0]
	}

	if msg, ok := errorMessages[tag]; ok {
		switch tag {
		case "min", "max", "gte", "lte":
			return fmt.Sprintf(msg, fieldName, param)
		case "eqfield":
			return fmt.Sprintf(msg, fieldName, param)
		default:
			return fmt.Sprintf(msg, strcase.ToSnake(fieldName))
		}
	}

	// Fallback message if no specific custom message is found for the tag
	log.Printf("Warning: No custom message found for tag '%s'. Using default error.", tag)
	return fmt.Sprintf("Validation error on %s: %s", fieldName, fieldError.Error())
}

func ValidateAppAccess(ctx context.Context, appID *uuid.UUID) error {
	tokenAppID, err := GetAppIDFromContext(ctx)

	if err != nil {
		log.Printf("Error retrieving app ID from context: %v", err)
		return err
	}

	if appID == nil && tokenAppID != appID {
		log.Printf("App ID mismatch: expected %s, got %s", appID, tokenAppID)
		return ErrUnauthorized
	}

	return nil
}
