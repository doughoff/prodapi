package dto

import (
	"github.com/hoffax/prodapi/postgres"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type RecipeDTO struct {
	RecipeID         pgtype.UUID     `json:"recipeId"`
	RecipeGroupID    pgtype.UUID     `json:"recipeGroupId"`
	Status           postgres.Status `json:"status"`
	Name             string          `json:"name"`
	ProductID        pgtype.UUID     `json:"productId"`
	ProductName      string          `json:"productName"`
	ProductUnit      postgres.Unit   `json:"productUnit"`
	ProducedQuantity float64         `json:"producedQuantity"`
	Revision         int32           `json:"revision"`
	IsCurrent        bool            `json:"isCurrent"`

	CreatedByUserID   pgtype.UUID `json:"createdByUserId"`
	CreatedByUserName string      `json:"createdByUserName"`

	Ingredients []*RecipeIngredientDTO `json:"ingredients"`
	CreatedAt   time.Time              `json:"createdAt"`
}

type RecipeIngredientDTO struct {
	ID          pgtype.UUID   `json:"id"`
	ProductID   pgtype.UUID   `json:"productId"`
	ProductName string        `json:"productName"`
	ProductUnit postgres.Unit `json:"productUnit"`
	RecipeID    pgtype.UUID   `json:"recipeId"`
	Quantity    float64       `json:"quantity"`
}

func ToRecipeIngredientDTO(ingredient *postgres.GetRecipeIngredientsRow) *RecipeIngredientDTO {
	return &RecipeIngredientDTO{
		ID:          ingredient.ID,
		ProductID:   ingredient.ProductID,
		ProductName: ingredient.ProductName,
		ProductUnit: ingredient.ProductUnit,
		RecipeID:    ingredient.RecipeID,
		Quantity:    float64(ingredient.Quantity) / 1000,
	}
}
