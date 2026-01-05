package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"

	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *AuthService) GenerateTokens(ctx context.Context, user *models.User) (*AuthTokens, error) {
	log := log("generateTokens")

	accessTokenExpireTime := time.Now().Add(time.Minute * 30)     // 30 minutes
	refreshTokenExpireTime := time.Now().Add(time.Hour * 24 * 30) // 30 days

	refreshJTI := generateTokenID()
	refreshTokenClaims := jwt.MapClaims{
		"jti":     refreshJTI,
		"user_id": user.ID,
		"app_id":  user.AppID,
		"type":    "refresh",
		"exp":     refreshTokenExpireTime.Unix(),
		"iat":     time.Now().Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodES256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(s.privateKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to sign refresh token")
		return nil, err
	}

	accessJTI := generateTokenID()
	accessTokenClaims := jwt.MapClaims{
		"jti":         accessJTI,
		"user_id":     user.ID,
		"app_id":      user.AppID,
		"type":        "access",
		"exp":         accessTokenExpireTime.Unix(),
		"iat":         time.Now().Unix(),
		"refresh_jti": refreshJTI,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodES256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(s.privateKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to sign access token")
		return &AuthTokens{
			AccessToken:  "",
			RefreshToken: "",
		}, err
	}

	err = s.repo.StoreRefreshToken(ctx, &models.RefreshToken{
		JTI:       refreshJTI,
		AppID:     user.AppID,
		UserID:    user.ID,
		Token:     hashToken(refreshTokenString),
		ExpiresAt: refreshTokenExpireTime,
		CreatedAt: time.Now(),
		IsActive:  true,
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to store refresh token")
		return nil, err
	}

	return &AuthTokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (s *AuthService) VerifyRefreshToken(ctx context.Context, token string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return s.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return nil, fmt.Errorf("%w: jti", ErrMissingRequiredClaim)
	}

	strAppID, ok := claims["app_id"].(string)
	appID, err := uuid.Parse(strAppID)
	if !ok || err != nil {
		return nil, fmt.Errorf("%w: app_id", ErrMissingRequiredClaim)
	}

	storedToken, err := s.repo.GetRefreshTokenByJTI(ctx, appID, jti)
	if err != nil {
		return nil, utils.ErrInternalServerError
	}

	if storedToken == nil {
		return nil, ErrTokenNotFound
	}

	if storedToken.Token != hashToken(token) {
		return nil, ErrTokenMismatch
	}

	if !storedToken.IsActive {
		return nil, ErrTokenRevoked
	}

	return jwtToken, nil
}

func (s *AuthService) VerifyAccessToken(ctx context.Context, token string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return s.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)

	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	refreshJTI, ok := claims["refresh_jti"].(string)

	if !ok {
		return nil, fmt.Errorf("%w: refresh_jti", ErrMissingRequiredClaim)
	}

	strAppID, ok := claims["app_id"].(string)
	appID, err := uuid.Parse(strAppID)

	if !ok || err != nil {
		return nil, fmt.Errorf("%w: app_id", ErrMissingRequiredClaim)
	}

	activeRefreshTokens, err := s.repo.GetUserActiveRefreshTokens(ctx, appID, nil, &refreshJTI)

	if err != nil {
		return nil, utils.ErrInternalServerError
	}

	if len(*activeRefreshTokens) == 0 {
		return nil, ErrTokenRevoked
	}

	return jwtToken, nil
}

func (s *AuthService) GetPublicKey() *ecdsa.PublicKey {
	return s.publicKey
}

func (s *AuthService) GetPrivateKey() *ecdsa.PrivateKey {
	return s.privateKey
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func generateTokenID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

var errorCode = "registration_failed"

func buildCallbackURL(baseURL string, tokens *AuthTokens, user *models.User, success bool) string {
	callbackURL, err := url.Parse(baseURL)
	if err != nil {
		log("buildCallbackURL").Error().Err(err).Msg("Invalid callback URL")
		return baseURL
	}

	query := callbackURL.Query()

	if success {
		query.Set("success", "true")
		query.Set("access_token", tokens.AccessToken)
		query.Set("refresh_token", tokens.RefreshToken)
		query.Set("user_id", user.ID.String())
	} else {
		query.Set("success", "false")
		query.Set("error", errorCode)
	}

	callbackURL.RawQuery = query.Encode()
	return callbackURL.String()
}

func buildRedirectURL(baseURL string, success bool) string {
	redirectURL, err := url.Parse(baseURL)
	if err != nil {
		log("buildRedirectURL").Error().Err(err).Msg("Invalid redirect URL")
		return baseURL
	}

	query := redirectURL.Query()

	if success {
		query.Set("success", "true")
	} else {
		query.Set("error", errorCode)
	}

	redirectURL.RawQuery = query.Encode()
	return redirectURL.String()
}
