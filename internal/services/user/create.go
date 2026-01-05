package user

import (
	"context"
	"time"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	createUserLog := log("CreateUser")

	if req.Email == "" || req.Name == "" || req.AppID == uuid.Nil || req.Provider == "" {
		createUserLog.Error().Msg("provider, appID, email, name may not be empty")
		return nil, utils.ErrBadRequest
	}

	if req.Provider != models.AuthProviderLocal {
		createUserLog.Error().Msg("currently only support 'local' provider")
		return nil, utils.ErrBadRequest
	}

	existingUser, err := s.repo.GetUserByEmail(ctx, req.AppID, req.Email)
	if err != nil {
		createUserLog.Error().Err(err).Msg("Failed to check existing user by email")
		return nil, utils.ErrInternalServerError
	}

	if existingUser != nil {
		createUserLog.Error().Err(err).Msg("User already exists")
		return nil, utils.ErrBadRequest
	}

	userID, err := uuid.NewV7()

	if err != nil {
		createUserLog.Error().Err(err).Msg("Failed to generate uuid V7 for userID")
		return nil, utils.ErrInternalServerError
	}

	user := &models.User{
		ID:              userID,
		AppID:           req.AppID,
		Name:            req.Name,
		Email:           req.Email,
		EmailVerifiedAt: nil,
	}

	userAuthentication := &models.UserAuthProvider{
		UserID:         user.ID,
		AppID:          user.AppID,
		Provider:       req.Provider,
		ProviderUserID: nil,
	}

	if req.Provider == models.AuthProviderLocal {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost+2)
		if err != nil {
			createUserLog.Error().Err(err).Msg("Failed to hash password")
			return nil, utils.ErrInternalServerError
		}

		userAuthentication.Password = string(hashedPassword)
	}

	opCtx, cancel := utils.ContextWithTimeout(5 * time.Second)
	defer cancel()

	if err := s.repo.CreateUser(opCtx, user, userAuthentication); err != nil {
		createUserLog.Error().Err(err).Msg("Failed to create user in repository")
		return nil, utils.ErrInternalServerError
	}

	return user, nil
}
