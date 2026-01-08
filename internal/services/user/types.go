package user

import (
	"context"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/services/app"
	"github.com/google/uuid"
)

//go:generate mockgen -source=user.go -destination=user_mock.go -package=services
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User, auth *models.UserAuthProvider) error
	GetUserAuthenticationByProvider(ctx context.Context, appID, userID uuid.UUID, provider models.AuthProvider) (*models.UserAuthProvider, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	GetAllUsersByAppID(ctx context.Context, appID uuid.UUID) ([]*models.User, error)
	GetAppUserByID(ctx context.Context, appID uuid.UUID, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, appID uuid.UUID, email string) (*models.User, error)
}

type UserService struct {
	repo       UserRepository
	appService *app.AppService
}

type UserIdentifier struct {
	ID    *uuid.UUID
	Email *string
}
