package user

import (
	"context"
	"fmt"
	"time"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/google/uuid"
)

func (s *UserService) GetAppUsers(ctx context.Context, appId uuid.UUID) ([]*models.User, error) {
	getAppUsersLog := log("GetAppUsers")

	optCtx, cancel := utils.ContextWithTimeout(5 * time.Second)
	defer cancel()

	users, err := s.repo.GetAllUsersByAppID(optCtx, appId)

	if err != nil {
		getAppUsersLog.Error().Err(err).Msg("Failed execute repository method GetAllUsersByAppID")
		return nil, fmt.Errorf("failed to get app users from repository: %w", err)
	}

	return users, nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]*models.User, error) {
	getUsersLog := log("GetUsers")

	opCtx, cancel := utils.ContextWithTimeout(5 * time.Second)
	defer cancel()

	users, err := s.repo.GetAllUsers(opCtx)

	if err != nil {
		getUsersLog.Error().Err(err).Msg("Failed to execute repository method GetAllUsers")
		return nil, fmt.Errorf("failed to get users from repository: %w", err)
	}

	return users, nil
}

func (s *UserService) GetUser(ctx context.Context, appID uuid.UUID, identifier UserIdentifier) (*models.User, error) {
	var user *models.User

	getUserLog := log("GetUser")

	opCtx, cancel := utils.ContextWithTimeout(5 * time.Second)
	defer cancel()

	if identifier.ID != nil {
		dbUser, err := s.repo.GetAppUserByID(opCtx, appID, *identifier.ID)
		if err != nil {
			getUserLog.Error().Err(err).Msg("Failed to execute repository method GetAppUserByID")
			return nil, fmt.Errorf("failed to get user from repository: %w", err)
		}

		user = dbUser

	} else if identifier.Email != nil {
		dbUser, err := s.repo.GetUserByEmail(opCtx, appID, *identifier.Email)
		if err != nil {
			getUserLog.Error().Err(err).Msg("Failed to execute repository method GetAppUserByID")
			return nil, fmt.Errorf("failed to get user from repository: %w", err)
		}

		user = dbUser
	} else {
		getUserLog.Panic().Msg("You must provide at least one identifier either id or user")
	}

	if user == nil {
		getUserLog.Error().Msg("No user found")
		return nil, utils.ErrNotFound
	}

	return user, nil
}
