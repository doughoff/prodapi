package dto

import (
	"github.com/hoffax/prodapi/postgres"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type ProductDTO struct {
	ID               pgtype.UUID     `json:"id"`
	Status           postgres.Status `json:"status"`
	Name             string          `json:"name"`
	Barcode          string          `json:"barcode"`
	Unit             postgres.Unit   `json:"unit"`
	BatchControl     bool            `json:"batchControl"`
	ConversionFactor int64           `json:"conversionFactor"`
	Stock            float64         `json:"stock"`
	AverageCost      int64           `json:"averageCost"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`
}

func ToProductDTO(product *postgres.Product) *ProductDTO {
	return &ProductDTO{
		ID:               product.ID,
		Status:           product.Status,
		Name:             product.Name,
		Barcode:          product.Barcode,
		Unit:             product.Unit,
		BatchControl:     product.BatchControl,
		ConversionFactor: product.ConversionFactor,
		CreatedAt:        product.CreatedAt.Time,
		UpdatedAt:        product.UpdatedAt.Time,
		Stock:            1,
		AverageCost:      1500,
	}
}
