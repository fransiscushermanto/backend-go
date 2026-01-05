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

func (s *UserService) GetUser(ctx context.Context, appID, id uuid.UUID) (*models.User, error) {
	getUserLog := log("GetUser")

	opCtx, cancel := utils.ContextWithTimeout(5 * time.Second)
	defer cancel()

	user, err := s.repo.GetAppUserByID(opCtx, appID, id)

	if err != nil {
		getUserLog.Error().Err(err).Msg("Failed to execute repository method GetAppUserByID")
		return nil, fmt.Errorf("failed to get user from repository: %w", err)
	}

	if user == nil {
		getUserLog.Error().Msg("No user found")
		return nil, utils.ErrNotFound
	}

	return user, nil
}
