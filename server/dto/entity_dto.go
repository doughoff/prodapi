package dto

import (
	"github.com/hoffax/prodapi/postgres"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type EntityDTO struct {
	ID        pgtype.UUID     `json:"id"`
	Status    postgres.Status `json:"status"`
	Name      string          `json:"name"`
	Ci        pgtype.Text     `json:"ci"`
	Ruc       pgtype.Text     `json:"ruc"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

func ToEntityDTO(entity *postgres.Entity) *EntityDTO {
	return &EntityDTO{
		ID:        entity.ID,
		Status:    entity.Status,
		Name:      entity.Name,
		Ci:        entity.Ci,
		Ruc:       entity.Ruc,
		CreatedAt: entity.CreatedAt.Time,
		UpdatedAt: entity.UpdatedAt.Time,
	}
}
