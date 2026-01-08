package auth

import (
	"context"

	"github.com/fransiscushermanto/backend/internal/models"
)

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest, options AuthOptions) (*models.RegisterResponse, error) {
	registerLog := log("Register")

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

	tokens, err := s.GenerateUserAuthTokens(ctx, user)
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
