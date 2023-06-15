package dto

import (
	"github.com/hoffax/prodapi/postgres"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type StockMovementDTO struct {
	ID                  pgtype.UUID           `json:"id"`
	Status              postgres.Status       `json:"status"`
	Type                postgres.MovementType `json:"type"`
	Date                time.Time             `json:"date"`
	EntityID            pgtype.UUID           `json:"entityId"`
	EntityName          string                `json:"entityName"`
	CreatedByUserID     pgtype.UUID           `json:"createdByUserId"`
	CreateByUserName    string                `json:"createByUserName"`
	CancelledByUserID   pgtype.UUID           `json:"cancelledByUserId"`
	CancelledByUserName string                `json:"cancelledByUserName"`
}

func ToStockMovementDTO(stockMovement *postgres.GetStockMovementByIDRow) *StockMovementDTO {
	return &StockMovementDTO{
		ID:                  stockMovement.ID,
		Status:              stockMovement.Status,
		Type:                stockMovement.Type,
		Date:                stockMovement.Date.Time,
		EntityID:            stockMovement.EntityID,
		EntityName:          stockMovement.EntityName.String,
		CreatedByUserID:     stockMovement.CreatedByUserID,
		CreateByUserName:    stockMovement.CreateByUserName.String,
		CancelledByUserID:   stockMovement.CancelledByUserID,
		CancelledByUserName: stockMovement.CancelledByUserName.String,
	}
}
