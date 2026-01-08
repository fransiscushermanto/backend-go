package auth

import (
	"context"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) LoginWithEmail(ctx context.Context, req *models.LoginWithEmailRequest, options AuthOptions) (*models.LoginResponse, error) {
	loginWithEmailLog := log("LoginWithEmail")

	user, err := s.userRepository.GetUserByEmail(ctx, req.AppID, req.Email)

	errUnauthorized := utils.ValidationError{
		Fields: []utils.FieldError{
			{Field: "email", Message: "Please enter valid credentials"},
			{Field: "password", Message: "Please enter valid credentials"},
		},
	}

	if user == nil || err != nil {
		loginWithEmailLog.Error().Err(err).Msg("User not found")
		return nil, errUnauthorized
	}

	userAuth, err := s.userRepository.GetUserAuthenticationByProvider(ctx, user.AppID, user.ID, req.Provider)

	if err != nil {
		loginWithEmailLog.Error().Err(err).Msg("User Auth not found")
		return nil, errUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(userAuth.Password), []byte(req.Password))

	if err != nil {
		loginWithEmailLog.Error().Err(err).Msg("Password not matched")
		return nil, errUnauthorized
	}

	err = s.repo.RevokeRefreshToken(ctx, user.AppID, user.ID)

	if err != nil {
		loginWithEmailLog.Error().Err(err).Msg("Failed to execute RevokeRefreshToken")
		return nil, utils.ErrInternalServerError
	}

	tokens, err := s.GenerateUserAuthTokens(ctx, user)

	if err != nil {
		loginWithEmailLog.Error().Err(err).Msg("Failed to generate tokens")

		if options.CallbackURL != "" {
			return &models.LoginResponse{CallbackURL: buildCallbackURL(options.CallbackURL, nil, nil, false)}, err
		}

		if options.RedirectURL != "" {
			return &models.LoginResponse{RedirectURL: buildRedirectURL(options.RedirectURL, false)}, err
		}

		return nil, err
	}

	loginResponse := &models.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	if options.CallbackURL != "" {
		loginResponse.CallbackURL = buildCallbackURL(options.CallbackURL, tokens, user, true)
	}

	if options.RedirectURL != "" {
		loginResponse.RedirectURL = buildRedirectURL(options.RedirectURL, true)
	}

	return loginResponse, nil
}
