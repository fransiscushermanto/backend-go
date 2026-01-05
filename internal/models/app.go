package models

import (
	"time"

	"github.com/google/uuid"
)

type App struct {
	ID        uuid.UUID `json:"id"`
	Name      []byte    `json:"name"`
	CreatedAt time.Time `json:"created_at" time_format:"2006-01-02T15:04:05Z"`
	UpdatedAt time.Time `json:"updated_at" time_format:"2006-01-02T15:04:05Z"`
}

type AppApiKey struct {
	ID         uuid.UUID `json:"id"`
	AppID      uuid.UUID `json:"app_id"`
	KeyHash    string    `json:"key_hash"`
	CreatedAt  time.Time `json:"created_at" time_format:"2006-01-02T15:04:05Z"`
	IsActive   bool      `json:"is_active"`
	LastUsedAt string    `json:"last_used_at"`
}

type UpdateAppApiKey struct {
	AppID      string    `json:"app_id"`
	APIKeyHash string    `json:"api_key_hash"`
	CreatedAt  time.Time `json:"created_at" time_format:"2006-01-02T15:04:05Z"`
}

type RegisterAppRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}

type RegisterAppResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	APIKey string `json:"api_key"`
}

type AppResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
