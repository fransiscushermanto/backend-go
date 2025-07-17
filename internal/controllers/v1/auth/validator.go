package auth

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	_validator *validator.Validate
	once       sync.Once
)

const PASSWORD_TAG_PREFIX = "password"

var Password = struct {
	TooShort       string
	InvalidPattern string
	InvalidType    string
	Unknown        string
}{
	TooShort:       fmt.Sprintf("%s_too_short", PASSWORD_TAG_PREFIX),
	InvalidPattern: fmt.Sprintf("%s_invalid_pattern", PASSWORD_TAG_PREFIX),
	InvalidType:    fmt.Sprintf("%s_type_invalid", PASSWORD_TAG_PREFIX),
	Unknown:        fmt.Sprintf("%s_unknown", PASSWORD_TAG_PREFIX),
}

var errorMessages = map[string]string{
	"required": "%s is a mandatory field.",
	"email":    "%s is not a valid email address.",
	"min":      "%s must be at least %s characters long.", // %s for field, %s for param (e.g., min=3)
	"max":      "%s must be at most %s characters long.",
	"gte":      "%s must be %s or greater.",
	"lte":      "%s must be %s or less.",
	"eqfield":  "%s must match the %s field.",

	Password.TooShort:       "%s must be at least 8 characters long.",
	Password.InvalidPattern: "%s must contain at least one uppercase letter, one digit, and one special character.",
	Password.InvalidType:    "Internal error: %s is not a string type for password validation.",
	Password.Unknown:        "%s failed password validation for an unknown reason.",
}

func InitValidator() *validator.Validate {
	once.Do(func() {
		_validator = validator.New()
		_validator.RegisterValidation("password-pattern", PasswordPattern)
	})

	return _validator
}

func PasswordPattern(fl validator.FieldLevel) bool {
	field := fl.Field()

	if field.Kind().String() != "string" {
		return false
	}

	password := field.String()

	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\-\+=\{\}\[\]:;"'<>,.?\/\\|]`).MatchString(password)

	if !hasUpper || !hasNumber || !hasSpecial {
		return false
	}

	return true
}

func RenderErrorMessage(fe validator.FieldError) (string, bool) {
	if fe.Tag() == "password-pattern" {
		// Use the determined specific failure tag
		specificFailureTag := determinePasswordStrengthFailure(fe.Value())
		// Look up the message using this specificFailureTag
		if msg, ok := errorMessages[specificFailureTag]; ok {
			// Return the formatted message and true (meaning it was handled)
			return fmt.Sprintf(msg, fe.Field()), true
		}
		// If the specific tag for password is not found, fallback to generic password message
		if msg, ok := errorMessages[PASSWORD_TAG_PREFIX]; ok {
			return fmt.Sprintf(msg, fe.Field()), true // Handled
		}
	}
	// Not handled by this specific callback, let GetValidationErrorMessage handle it
	return "", false
}

func determinePasswordStrengthFailure(passwordValue interface{}) string {
	strPassword, ok := passwordValue.(string)
	if !ok {
		// If the value is not a string (which PasswordStrength would also reject)
		return Password.InvalidType
	}

	if len(strPassword) < 8 {
		return Password.TooShort
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(strPassword)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(strPassword)
	hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\-\+=\{\}\[\]:;"'<>,.?\/\\|]`).MatchString(strPassword)

	if !hasUpper || !hasNumber || !hasSpecial {
		return Password.InvalidPattern
	}

	// This should ideally not be reached if PasswordStrength returned false.
	// It's a fallback for unexpected scenarios.
	return Password.Unknown
}
