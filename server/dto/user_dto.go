package dto

import (
	"github.com/hoffax/prodapi/postgres"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type UserDTO struct {
	ID        pgtype.UUID
	Status    postgres.Status
	Email     string
	Name      string
	Roles     []string
	CreatedAt time.Time
	UpdatedAt time.Time
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
