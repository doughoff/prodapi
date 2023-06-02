package dto

type UserDTO struct {
	ID        string   `json:"id"`
	Status    string   `json:"status"`
	Roles     []string `json:"roles"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}
