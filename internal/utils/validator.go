package utils

import (
	"context"
	"fmt"
	"log"
	"reflect"
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

func ValidateBodyRequest(req interface{}) error {
	v := reflect.ValueOf(req)
	t := reflect.TypeOf(req)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ErrBadRequest
		}

		v = v.Elem()
		t = t.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ErrBadRequest
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Get validate tag
		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue // No validation rules
		}

		// Check if this field should be validated for emptiness
		shouldValidate, err := shouldValidateField(validateTag, v, t)
		if err != nil {
			return err
		}

		if !shouldValidate {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			if field.String() == "" {
				return ErrBadRequest
			}
		case reflect.Ptr:
			elem := field.Elem()
			if (elem.Kind() == reflect.String && elem.String() == "") || field.IsNil() {
				return ErrBadRequest
			}
		}

	}

	return nil
}

func shouldValidateField(validateTag string, structValue reflect.Value, structType reflect.Type) (bool, error) {
	rules := strings.Split(validateTag, ",")

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)

		switch {
		case rule == "required":
			return true, nil

		case strings.HasPrefix(rule, "required_unless="):
			// Parse: required_unless=Provider local
			condition := strings.TrimPrefix(rule, "required_unless=")
			parts := strings.Split(condition, " ")
			if len(parts) != 2 {
				continue
			}

			fieldName := parts[0]
			expectedValue := parts[1]

			// Get the referenced field value
			refField, err := getFieldByName(structValue, structType, fieldName)
			if err != nil {
				return false, err
			}

			// Check if condition is met (if true, don't validate)
			if getFieldStringValue(refField) == expectedValue {
				return false, nil // Don't validate because condition is met
			}
			return true, nil // Validate because condition is not met

		case strings.HasPrefix(rule, "required_if="):
			// Parse: required_if=Provider local
			condition := strings.TrimPrefix(rule, "required_if=")
			parts := strings.Split(condition, " ")
			if len(parts) != 2 {
				continue
			}

			fieldName := parts[0]
			expectedValue := parts[1]

			// Get the referenced field value
			refField, err := getFieldByName(structValue, structType, fieldName)
			if err != nil {
				return false, err
			}

			// Check if condition is met (if true, validate)
			if getFieldStringValue(refField) == expectedValue {
				return true, nil // Validate because condition is met
			}
			return false, nil // Don't validate because condition is not met

		case rule == "omitempty":
			return false, nil // Don't validate empty values
		}
	}

	return false, nil // Default: don't validate
}

func getFieldByName(structValue reflect.Value, structType reflect.Type, fieldName string) (reflect.Value, error) {
	for i := 0; i < structValue.NumField(); i++ {
		if structType.Field(i).Name == fieldName {
			return structValue.Field(i), nil
		}
	}
	return reflect.Value{}, fmt.Errorf("field %s not found", fieldName)
}

func getFieldStringValue(field reflect.Value) string {
	switch field.Kind() {
	case reflect.String:
		return field.String()
	case reflect.Ptr:
		if field.IsNil() {
			return ""
		}
		elem := field.Elem()
		if elem.Kind() == reflect.String {
			return elem.String()
		}
	}
	return field.String() // Try to convert to string
}
