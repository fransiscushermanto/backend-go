package app

import (
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/rs/zerolog"
)

func log(method string) *zerolog.Logger {
	l := utils.Log().With().Str("service", "App").Str("method", method).Logger()
	return &l
}

func NewAppService(repo AppRepository, prefixApiKey string, secretKey string) *AppService {
	return &AppService{repo: repo, prefixApiKey: prefixApiKey, secretKey: secretKey}
}
