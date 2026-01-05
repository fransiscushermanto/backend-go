package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID  `json:"id"`
	AppID           uuid.UUID  `json:"app_id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	IsEmailVerified bool       `json:"is_email_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type UserAuthProvider struct {
	UserID         uuid.UUID    `json:"user_id"`
	AppID          uuid.UUID    `json:"app_id"`
	Provider       AuthProvider `json:"provider"`
	ProviderUserID *string      `json:"provider_user_id"`
	Password       string       `json:"password"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type CreateUserRequest struct {
	Provider      AuthProvider `json:"provider" validate:"required,oneof=local google"`
	ProviderToken string       `json:"provider_token" validate:"required_unless=Provider local,omitempty"`
	AppID         uuid.UUID    `json:"app_id" validate:"required"`
	Name          string       `json:"name" validate:"required,min=3,max=100"`
	Email         string       `json:"email" validate:"required,email"`
	Password      string       `json:"password" validate:"required_if=Provider local,omitempty,password-pattern"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type UserResponse struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	IsEmailVerified bool      `json:"is_email_verified"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:              u.ID,
		Name:            u.Name,
		Email:           u.Email,
		IsEmailVerified: u.IsEmailVerified,
	}
}
