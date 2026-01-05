package auth

import (
	"github.com/fransiscushermanto/backend/internal/config"
	"github.com/fransiscushermanto/backend/internal/services/user"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/rs/zerolog"
)

func log(method string) *zerolog.Logger {
	l := utils.Log().With().Str("service", "Auth").Str("method", method).Logger()
	return &l
}

func NewAuthService(repo AuthRepository, userRepository user.UserRepository, userService *user.UserService, keys *config.CryptoKeys) *AuthService {
	if !keys.IsValid() {
		panic("AuthService requires valid keys")
	}

	return &AuthService{
		repo:           repo,
		userRepository: userRepository,
		userService:    userService,
		privateKey:     keys.PrivateKey,
		publicKey:      keys.PublicKey,
	}
}
