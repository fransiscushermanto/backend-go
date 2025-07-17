package models

type User struct {
	ID        string `json:"id"`
	AppID     string `json:"app_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserAuthProvider struct {
	UserID         string       `json:"user_id"`
	AppID          string       `json:"app_id"`
	Provider       AuthProvider `json:"provider"`
	ProviderUserID string       `json:"provider_user_id"`
	Password       string       `json:"password"`
	CreatedAt      string       `json:"created_at"`
	UpdatedAt      string       `json:"updated_at"`
}

type CreateUserRequest struct {
	Provider      AuthProvider `json:"provider" validate:"required,oneof=local google"`
	ProviderToken string       `json:"provider_token" validate:"required_unless=Provider local,omitempty"`
	AppID         string       `json:"app_id" validate:"required"`
	Name          string       `json:"name" validate:"required,min=3,max=100"`
	Email         string       `json:"email" validate:"required,email"`
	Password      string       `json:"password" validate:"required_if=Provider local,omitempty,password-pattern"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
