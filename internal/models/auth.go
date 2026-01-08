package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthProvider string

const (
	AuthProviderPasswordless AuthProvider = "passwordless"
	AuthProviderLocal        AuthProvider = "local"
	AuthProviderGoogle       AuthProvider = "google"
)

type AuthResponseType string

const (
	AuthResponseCallback AuthResponseType = "callback"
	AuthResponseRedirect AuthResponseType = "redirect"
	AuthResponseJSON     AuthResponseType = "json"
)

type LoginWithEmailRequest struct {
	Provider AuthProvider `json:"provider" validate:"required,oneof=local"`
	AppID    uuid.UUID    `json:"app_id" validate:"required"`
	Email    string       `json:"email" validate:"required,email"`
	Password string       `json:"password" validate:"required"`
	DeviceID string       `json:"device_id"`
}

type LoginWithPasswordlessRequest struct {
	Provider AuthProvider `json:"provider" validate:"required,oneof=passwordless"`
	AppID    uuid.UUID    `json:"app_id" validate:"required"`
	Email    string       `json:"email" validate:"required,email"`
	DeviceID string       `json:"device_id"`
}

type LoginWithOtherProviderRequest struct {
	Provider      AuthProvider `json:"provider" validate:"required,oneof=google"`
	ProviderToken string       `json:"provider_token" validate:"required"`
	AppID         uuid.UUID    `json:"app_id" validate:"required"`
	DeviceID      string       `json:"device_id"`
}

type RegisterRequest struct {
	Provider      AuthProvider `json:"provider" validate:"required,oneof=local google"`
	ProviderToken string       `json:"provider_token" validate:"required_unless=Provider local,omitempty"`
	AppID         uuid.UUID    `json:"app_id" validate:"required"`
	Name          string       `json:"name" validate:"required,min=3,max=100"`
	Email         string       `json:"email" validate:"required,email"`
	Password      string       `json:"password" validate:"required_if=Provider local,omitempty,password-pattern"`
	DeviceID      string       `json:"device_id"`
}

type ForgetPasswordRequest struct {
	AppID *uuid.UUID `json:"app_id" validate:"required"`
	Email *string    `json:"email" validate:"required,email"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`

	CallbackURL string `json:"callback_url,omitempty"`
	RedirectURL string `json:"redirect_url,omitempty"`
}

type RegisterResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`

	CallbackURL string `json:"callback_url,omitempty"`
	RedirectURL string `json:"redirect_url,omitempty"`
}

type RefreshToken struct {
	JTI       string    `json:"jti"`
	UserID    uuid.UUID `json:"user_id"`
	AppID     uuid.UUID `json:"app_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type ResetPasswordToken struct {
	JTI       string    `json:"jti"`
	UserID    uuid.UUID `json:"user_id"`
	AppID     uuid.UUID `json:"app_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshTokenRequest struct {
	RefreshToken *string `json:"refresh_token" validate:"required"`
	DeviceID     *string `json:"device_id" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
