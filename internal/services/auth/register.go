package auth

import (
	"context"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/google/uuid"
)

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest, options AuthOptions) (*models.RegisterResponse, error) {
	registerLog := log("Register")

	if req.Provider == "" || req.AppID == uuid.Nil || req.Email == "" || req.Name == "" {
		registerLog.Error().Msg("provider, appID, email, name may not be empty")
		return nil, utils.ErrBadRequest
	}

	if req.Provider == models.AuthProviderLocal && req.Password == "" || req.Provider != models.AuthProviderLocal && req.ProviderToken == "" {
		registerLog.Error().Msg("either provider is local and password may not empty, or provider is not local and providerToken may not empty")
		return nil, utils.ErrBadRequest
	}

	createUserReq := &models.CreateUserRequest{
		AppID:         req.AppID,
		Name:          req.Name,
		Provider:      req.Provider,
		ProviderToken: req.ProviderToken,
		Email:         req.Email,
		Password:      req.Password,
	}

	user, err := s.userService.CreateUser(ctx, createUserReq)

	if err != nil {
		registerLog.Error().Err(err).Msg("Failed to create user")

		if options.CallbackURL != "" {
			return &models.RegisterResponse{CallbackURL: buildCallbackURL(options.CallbackURL, nil, nil, false)}, err
		}

		if options.RedirectURL != "" {
			return &models.RegisterResponse{RedirectURL: buildRedirectURL(options.RedirectURL, false)}, err
		}

		return nil, err
	}

	tokens, err := s.GenerateTokens(ctx, user)
	if err != nil {
		registerLog.Error().Err(err).Msg("Failed to generate tokens")

		if options.CallbackURL != "" {
			return &models.RegisterResponse{CallbackURL: buildCallbackURL(options.CallbackURL, nil, nil, false)}, err
		}

		if options.RedirectURL != "" {
			return &models.RegisterResponse{RedirectURL: buildRedirectURL(options.RedirectURL, false)}, err
		}

		return nil, err
	}

	registerResponse := &models.RegisterResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	if options.CallbackURL != "" {
		registerResponse.CallbackURL = buildCallbackURL(options.CallbackURL, tokens, user, true)
	}

	if options.RedirectURL != "" {
		registerResponse.RedirectURL = buildRedirectURL(options.RedirectURL, true)
	}

	return registerResponse, nil
}
