package models

type App struct {
	ID        string `json:"id"`
	Name      []byte `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type RegisterAppRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}

type AppResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
