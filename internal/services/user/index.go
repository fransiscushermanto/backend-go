package user

import (
	"github.com/fransiscushermanto/backend/internal/services/app"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/rs/zerolog"
)

func log(method string) *zerolog.Logger {
	l := utils.Log().With().Str("service", "User").Str("method", method).Logger()
	return &l
}

func NewUserService(repo UserRepository, appService *app.AppService) *UserService {
	return &UserService{repo: repo, appService: appService}
}
