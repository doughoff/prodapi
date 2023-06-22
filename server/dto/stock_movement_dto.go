package dto

import (
	"github.com/hoffax/prodapi/postgres"
	"github.com/jackc/pgx/v5/pgtype"
	"math"
	"time"
)

type StockMovementDTO struct {
	ID                  pgtype.UUID           `json:"id"`
	Status              postgres.Status       `json:"status"`
	Type                postgres.MovementType `json:"type"`
	Date                time.Time             `json:"date"`
	EntityID            pgtype.UUID           `json:"entityId"`
	EntityName          string                `json:"entityName"`
	DocumentNumber      string                `json:"documentNumber"`
	CreatedByUserID     pgtype.UUID           `json:"createdByUserId"`
	CreateByUserName    string                `json:"createByUserName"`
	CancelledByUserID   pgtype.UUID           `json:"cancelledByUserId"`
	CancelledByUserName string                `json:"cancelledByUserName"`
	Total               int64                 `json:"total"`

	Items []*StockMovementItemDTO `json:"items"`
}

type StockMovementItemDTO struct {
	ID          pgtype.UUID `json:"id"`
	ProductID   pgtype.UUID `json:"productId"`
	ProductName string      `json:"productName"`
	Quantity    float64     `json:"quantity"`
	Price       int64       `json:"price"`
	Batch       string      `json:"batch"`
	Total       int64       `json:"total"`
}

func ToStockMovementDTO(stockMovement *postgres.GetStockMovementByIDRow) *StockMovementDTO {
	return &StockMovementDTO{
		ID:                  stockMovement.ID,
		Status:              stockMovement.Status,
		Type:                stockMovement.Type,
		Date:                stockMovement.Date.Time,
		EntityID:            stockMovement.EntityID,
		EntityName:          stockMovement.EntityName.String,
		DocumentNumber:      stockMovement.DocumentNumber.String,
		CreatedByUserID:     stockMovement.CreatedByUserID,
		CreateByUserName:    stockMovement.CreateByUserName.String,
		CancelledByUserID:   stockMovement.CancelledByUserID,
		CancelledByUserName: stockMovement.CancelledByUserName.String,
		Items:               make([]*StockMovementItemDTO, 0),
	}
}

func ToStockMovementItemDTO(item *postgres.GetStockMovementItemsRow) *StockMovementItemDTO {
	return &StockMovementItemDTO{
		ID:          item.ID,
		ProductID:   item.ProductID,
		ProductName: item.ProductName.String,
		Quantity:    float64(item.Quantity) / 1000,
		Price:       item.Price,
		Batch:       item.Batch.String,
		Total:       int64(math.Round(float64(item.Quantity*item.Price) / 1000)),
	}
}
