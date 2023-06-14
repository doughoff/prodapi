package dto

import (
	"github.com/hoffax/prodapi/postgres"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type UserDTO struct {
	ID        pgtype.UUID     `json:"id"`
	Status    postgres.Status `json:"status"`
	Email     string          `json:"email"`
	Name      string          `json:"name"`
	Roles     []string        `json:"roles"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

func ToUserDTO(user *postgres.User) *UserDTO {
	return &UserDTO{
		ID:        user.ID,
		Status:    user.Status,
		Email:     user.Email,
		Name:      user.Name,
		Roles:     user.Roles,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}
}
