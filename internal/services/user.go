package services

import (
	"context"
	"fmt"
	"time"

	"github.com/fransiscushermanto/backend/internal/constants"
	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=user.go -destination=user_mock.go -package=services
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User, auth *models.UserAuthProvider) error
	GetUserAuthenticationByProvider(ctx context.Context, appID, userID string, provider models.AuthProvider) (*models.UserAuthProvider, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	GetAllUsersByAppID(ctx context.Context, appID string) ([]*models.User, error)
	GetAppUserByID(ctx context.Context, appID string, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, appID, email string) (*models.User, error)
}

type UserService struct {
	repo       UserRepository
	appService *AppService
}

func NewUserService(repo UserRepository, appService *AppService) *UserService {
	return &UserService{repo: repo, appService: appService}
}

func (s *UserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	if req.Email == "" || req.Name == "" || req.AppID == "" || req.Provider == "" {
		return nil, utils.ErrBadRequest
	}

	if req.Provider == models.AuthProviderLocal && req.Password == "" {
		return nil, utils.ErrBadRequest
	}

	existingUser, err := s.repo.GetUserByEmail(ctx, req.AppID, req.Email)
	if err != nil {
		utils.Log().Error().Err(err).Msg("Failed to check existing user by email")
		return nil, utils.ErrInternalServerError
	}

	if existingUser != nil {
		utils.Log().Error().Err(err).Msg("User already exists")
		return nil, utils.ErrBadRequest
	}

	user := &models.User{
		ID:        uuid.New().String(),
		AppID:     req.AppID,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now().Format(constants.TimeFormatISO),
		UpdatedAt: time.Now().Format(constants.TimeFormatISO),
	}

	userAuthentication := &models.UserAuthProvider{
		UserID:    user.ID,
		AppID:     user.AppID,
		Provider:  req.Provider,
		CreatedAt: time.Now().Format(constants.TimeFormatISO),
		UpdatedAt: time.Now().Format(constants.TimeFormatISO),
	}

	if req.Provider == models.AuthProviderLocal {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost+2)
		if err != nil {
			utils.Log().Error().Err(err).Msg("Failed to hash password")
			return nil, utils.ErrInternalServerError
		}

		userAuthentication.Password = string(hashedPassword)
	}

	opCtx, cancel := utils.ContextWithTimeout(5 * time.Second)
	defer cancel()

	if err := s.repo.CreateUser(opCtx, user, userAuthentication); err != nil {
		utils.Log().Error().Err(err).Msg("Failed to create user in repository")
		return nil, utils.ErrInternalServerError
	}

	return user, nil
}

func (s *UserService) GetAppUsers(ctx context.Context, appId string) ([]*models.User, error) {
	optCtx, cancel := utils.ContextWithTimeout(5 * time.Second)
	defer cancel()

	users, err := s.repo.GetAllUsersByAppID(optCtx, appId)

	if err != nil {
		return nil, fmt.Errorf("failed to get app users from repository: %w", err)
	}

	return users, nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]*models.User, error) {
	opCtx, cancel := utils.ContextWithTimeout(5 * time.Second)
	defer cancel()

	users, err := s.repo.GetAllUsers(opCtx)

	if err != nil {
		return nil, fmt.Errorf("failed to get users from repository: %w", err)
	}

	return users, nil
}

func (s *UserService) GetUser(ctx context.Context, appID string, id string) (*models.User, error) {
	opCtx, cancel := utils.ContextWithTimeout(5 * time.Second)
	defer cancel()

	if !utils.IsValidUUID(appID) || !utils.IsValidUUID(id) {
		return nil, utils.ErrNotFound
	}

	user, err := s.repo.GetAppUserByID(opCtx, appID, id)

	if err != nil {
		return nil, fmt.Errorf("failed to get user from repository: %w", err)
	}

	if user == nil {
		return nil, utils.ErrNotFound
	}

	return user, nil
}
