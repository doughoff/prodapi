package dto

import (
	"github.com/hoffax/prodapi/postgres"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type ProductionOrderDTO struct {
	ID             pgtype.UUID             `json:"id"`
	Status         postgres.Status         `json:"status"`
	ProductionStep postgres.ProductionStep `json:"productionStep"`
	Code           string                  `json:"code"`
	Cycles         float64                 `json:"cycles"`
	Output         float64                 `json:"output"`
	RecipeID       pgtype.UUID             `json:"recipeId"`
	RecipeName     string                  `json:"recipeName"`
	ProductID      pgtype.UUID             `json:"productId"`
	ProductName    string                  `json:"productName"`

	CreatedAt           time.Time   `json:"createdAt"`
	UpdatedAt           time.Time   `json:"updatedAt"`
	CreatedByUserID     pgtype.UUID `json:"createdByUserId"`
	CreateByUserName    string      `json:"createByUserName"`
	CancelledByUserID   pgtype.UUID `json:"cancelledByUserId"`
	CancelledByUserName string      `json:"cancelledByUserName"`

	TargetAmount float64 `json:"targetAmount"`

	CycleInstances []*ProductionOrderCycleDTO    `json:"cycleInstances"`
	Movements      []*ProductionOrderMovementDTO `json:"movements"`
}

type ProductionOrderCycleDTO struct {
	ID             pgtype.UUID             `json:"id"`
	Factor         float64                 `json:"factor"`
	ProductionStep postgres.ProductionStep `json:"productionStep"`
	CompletedAt    *time.Time              `json:"completedAt"`
}

type ProductionOrderMovementDTO struct {
	StockMovementDTO
	CycleID pgtype.UUID `json:"cycleId"`
}
