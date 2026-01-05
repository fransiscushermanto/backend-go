package app

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/fransiscushermanto/backend/internal/utils"
)

func (s *AppService) GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s_%s", s.prefixApiKey, hex.EncodeToString(bytes)), nil
}

func (s *AppService) ParseAppName(rawName string) (string, error) {
	decryptedName, err := utils.Decrypt([]byte(s.secretKey), []byte(rawName))

	if err != nil {
		return "", err
	}

	return string(decryptedName), nil
}
