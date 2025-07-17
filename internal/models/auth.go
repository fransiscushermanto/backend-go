package models

type AuthProvider string

const (
	AuthProviderLocal  AuthProvider = "local"
	AuthProviderGoogle AuthProvider = "google"
)

type AuthResponseType string

const (
	AuthResponseRedirect AuthResponseType = "redirect"
	AuthResponseJSON     AuthResponseType = "json"
)

type RegisterRequest struct {
	Provider      AuthProvider `json:"provider" validate:"required,oneof=local google"`
	ProviderToken string       `json:"provider_token" validate:"required_unless=Provider local,omitempty"`
	AppID         string       `json:"app_id" validate:"required"`
	Name          string       `json:"name" validate:"required,min=3,max=100"`
	Email         string       `json:"email" validate:"required,email"`
	Password      string       `json:"password" validate:"required_if=Provider local,omitempty,password-pattern"`
}

// TODO: temporary return the auth_token and refresh_token
type RegisterResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`

	RedirectURL string `json:"redirect_url,omitempty"`
}

type RefreshToken struct {
	JTI       string `json:"jti"`
	UserID    string `json:"user_id"`
	AppID     string `json:"app_id"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}
