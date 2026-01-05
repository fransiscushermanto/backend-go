package repositories

import (
	"github.com/fransiscushermanto/backend/internal/repositories/app"
	"github.com/fransiscushermanto/backend/internal/repositories/auth"
	"github.com/fransiscushermanto/backend/internal/repositories/user"
	"github.com/fransiscushermanto/backend/internal/utils"
)

type Repositories struct {
	App *app.AppRepository
}

func NewAppRepository(database *utils.Database, lockTimeout *int) *app.AppRepository {
	return app.NewAppRepository(database, lockTimeout)
}

func NewAuthRepository(database *utils.Database) *auth.AuthRepository {
	return auth.NewAuthRepository(database)
}

func NewUserRepository(database *utils.Database) *user.UserRepository {
	return user.NewUserRepository(database)
}
