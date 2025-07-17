package services

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"time"

	"github.com/fransiscushermanto/backend/internal/constants"
	"github.com/fransiscushermanto/backend/internal/models"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type AuthRepository interface {
	StoreRefreshToken(ctx context.Context, token *models.RefreshToken) error
	GetRefreshToken(ctx context.Context, appID, userID, jti string) (*models.RefreshToken, error)
}

type AuthService struct {
	repo        AuthRepository
	userService *UserService
	privateKey  *ecdsa.PrivateKey
	publicKey   *ecdsa.PublicKey
}

func NewAuthService(repo AuthRepository, userService *UserService, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) *AuthService {
	return &AuthService{
		repo:        repo,
		userService: userService,
		privateKey:  privateKey,
		publicKey:   publicKey,
	}
}

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest, callbackUrl string) (*models.RegisterResponse, error) {
	if req.Provider == "" || req.AppID == "" || req.Email == "" || req.Name == "" {
		return nil, utils.ErrBadRequest
	}

	if req.Provider == models.AuthProviderLocal && req.Password == "" || req.Provider != models.AuthProviderLocal && req.ProviderToken == "" {
		return nil, utils.ErrBadRequest
	}

	createUserReq := &models.CreateUserRequest{
		AppID:         req.AppID,
		Name:          req.Name,
		Provider:      req.Provider,
		ProviderToken: req.ProviderToken,
		Email:         req.Email,
		Password:      req.Password,
	}

	user, err := s.userService.CreateUser(ctx, createUserReq)

	if err != nil {
		utils.Log().Error().Err(err).Msg("Failed to create user")
		return nil, err
	}

	tokens, err := s.generateTokens(ctx, user)
	if err != nil {
		utils.Log().Error().Err(err).Msg("Failed to generate tokens")
		return nil, err
	}

	registerResponse := &models.RegisterResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	if callbackUrl != "" {
		registerResponse.RedirectURL = s.buildCallbackURL(callbackUrl, tokens, user, true)
	}

	return registerResponse, nil
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

func (s *AuthService) generateTokens(ctx context.Context, user *models.User) (*AuthTokens, error) {
	accessTokenExpireTime := time.Now().Add(time.Minute * 15)     // 15 minutes
	refreshTokenExpireTime := time.Now().Add(time.Hour * 24 * 30) // 30 days

	accessJTI := generateTokenID()
	accessTokenClaims := jwt.MapClaims{
		"jti":     accessJTI,
		"user_id": user.ID,
		"app_id":  user.AppID,
		"type":    "access",
		"exp":     accessTokenExpireTime.Unix(),
		"iat":     time.Now().Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodES256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(s.privateKey)
	if err != nil {
		utils.Log().Error().Err(err).Msg("Failed to sign access token")
		return &AuthTokens{
			AccessToken:  "",
			RefreshToken: "",
		}, err
	}

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
		utils.Log().Error().Err(err).Msg("Failed to sign refresh token")
		return nil, err
	}

	err = s.repo.StoreRefreshToken(ctx, &models.RefreshToken{
		JTI:       refreshJTI,
		AppID:     user.AppID,
		UserID:    user.ID,
		Token:     hashToken(refreshTokenString),
		ExpiresAt: refreshTokenExpireTime.Format(constants.TimeFormatISO),
		CreatedAt: time.Now().Format(constants.TimeFormatISO),
		IsActive:  true,
	})

	if err != nil {
		utils.Log().Error().Err(err).Msg("Failed to store refresh token")
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

	if err != nil || !jwtToken.Valid {
		return nil, utils.ErrUnauthorized
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, utils.ErrUnauthorized
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return nil, utils.ErrUnauthorized
	}

	app_id, ok := claims["app_id"].(string)
	if !ok {
		return nil, utils.ErrUnauthorized
	}

	user_id, ok := claims["user_id"].(string)
	if !ok {
		return nil, utils.ErrUnauthorized
	}

	storedToken, err := s.repo.GetRefreshToken(ctx, app_id, user_id, jti)
	if err != nil {
		return nil, utils.ErrUnauthorized
	}

	if storedToken.Token != hashToken(token) {
		return nil, utils.ErrUnauthorized
	}

	if !storedToken.IsActive || !jwtToken.Valid {
		return nil, utils.ErrUnauthorized
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

	if err != nil || !jwtToken.Valid {
		return nil, utils.ErrUnauthorized
	}

	return jwtToken, nil
}

func (s *AuthService) GetPublicKey() *ecdsa.PublicKey {
	return s.publicKey
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

func (s *AuthService) buildCallbackURL(baseURL string, tokens *AuthTokens, user *models.User, success bool) string {
	callbackURL, err := url.Parse(baseURL)
	if err != nil {
		utils.Log().Error().Err(err).Msg("Invalid callback URL")
		return baseURL
	}

	query := callbackURL.Query()

	if success {
		query.Set("success", "true")
		query.Set("access_token", tokens.AccessToken)
		query.Set("refresh_token", tokens.RefreshToken)
		query.Set("user_id", user.ID)
	} else {
		query.Set("success", "false")
		query.Set("error", "registration_failed")
	}

	callbackURL.RawQuery = query.Encode()
	return callbackURL.String()
}
