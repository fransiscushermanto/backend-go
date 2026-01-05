package app

import (
	"context"
	"fmt"
	"time"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *AppService) Register(ctx context.Context, req *models.RegisterAppRequest) (*models.RegisterAppResponse, error) {
	registerLog := log("Register")

	if req.Name == "" {
		return nil, utils.ErrBadRequest
	}

	name, err := utils.Encrypt([]byte(s.secretKey), []byte(req.Name))

	if err != nil {
		return nil, fmt.Errorf("failed to encrypt app name: %w", err)
	}

	apiKey, err := s.GenerateAPIKey()
	if err != nil {
		registerLog.Error().Err(err).Msg("Failed to generate API key for app")
		return nil, utils.ErrInternalServerError
	}

	hashedApiKey, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	if err != nil {
		registerLog.Error().Err(err).Msg("Failed to hash API key for app")
		return nil, utils.ErrInternalServerError
	}

	appID, err := uuid.NewV7()

	if err != nil {
		fmt.Println("Error generating appID:", err)
		return nil, utils.ErrInternalServerError
	}

	app := &models.App{
		ID:        appID,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	appAPIKeyID, err := uuid.NewV7()

	if err != nil {
		fmt.Println("Error generating appAPIKeyID:", err)
		return nil, utils.ErrInternalServerError
	}

	appApiKey := &models.AppApiKey{
		ID:        appAPIKeyID,
		AppID:     app.ID,
		KeyHash:   string(hashedApiKey),
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	opCtx, cancel := utils.ContextWithTimeout(5 * time.Second)
	defer cancel()

	if err := s.repo.RegisterApp(opCtx, app, appApiKey); err != nil {
		registerLog.Error().Err(err).Msg("Failed to execute method RegisterApp")
		return nil, utils.ErrInternalServerError
	}

	return &models.RegisterAppResponse{
		ID:     app.ID.String(),
		Name:   req.Name,
		APIKey: apiKey,
	}, nil
}
