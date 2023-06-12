package dto

import (
	"github.com/hoffax/prodapi/postgres"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type RecipeDTO struct {
	RecipeID      pgtype.UUID     `json:"recipeId"`
	RecipeGroupID pgtype.UUID     `json:"recipeGroupId"`
	Status        postgres.Status `json:"status"`
	Name          string          `json:"name"`
	Revision      int32           `json:"revision"`
	IsCurrent     bool            `json:"isCurrent"`

	Ingredients []*RecipeIngredientDTO `json:"ingredients"`
	CreatedAt   time.Time              `json:"createdAt"`
}

type RecipeIngredientDTO struct {
	ID          pgtype.UUID `json:"id"`
	ProductID   pgtype.UUID `json:"productId"`
	ProductName string      `json:"productName"`
	RecipeID    pgtype.UUID `json:"recipeId"`
	Quantity    int32       `json:"quantity"`
}

func ToRecipeIngredientDTO(ingredient *postgres.GetRecipeIngredientsRow) *RecipeIngredientDTO {
	return &RecipeIngredientDTO{
		ID:          ingredient.ID,
		ProductID:   ingredient.ProductID,
		ProductName: ingredient.ProductName,
		RecipeID:    ingredient.RecipeID,
		Quantity:    ingredient.Quantity,
	}
}
