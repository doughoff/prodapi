package routes

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/hoffax/prodapi/postgres"
	"github.com/hoffax/prodapi/server/dto"
	"github.com/hoffax/prodapi/server/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type getAllRecipeQuery struct {
	StatusOptions []string `query:"status"`
	Search        string   `query:"search"`
	Limit         int32    `query:"limit"`
	Offset        int32    `query:"offset"`
}

func (r *RouteManager) getAllRecipes(c *fiber.Ctx, tx *pgx.Tx) error {
	params := new(getAllRecipeQuery)
	if err := c.QueryParser(params); err != nil {
		return types.NewInvalidParamsError("invalid query params")
	}

	statusOptions := make([]postgres.Status, len(params.StatusOptions))
	for i, status := range params.StatusOptions {
		statusOptions[i] = postgres.Status(status)
		if !statusOptions[i].Valid() {
			return types.NewInvalidParamsError("invalid status option")
		}
	}

	recipes, err := r.db.GetRecipes(c.Context(), *tx, &postgres.GetRecipesParams{
		Search:        pgtype.Text(sql.NullString{String: params.Search, Valid: true}),
		StatusOptions: statusOptions,
		PageLimit:     params.Limit,
		PageOffset:    params.Offset,
	})
	if err != nil {
		return err
	}

	var totalCount int64
	if len(recipes) > 0 {
		totalCount = recipes[0].FullCount
	}
	recipeIds := make([]pgtype.UUID, 0)
	resultRecipes := make([]*dto.RecipeDTO, len(recipes))
	hashMap := make(map[pgtype.UUID]*dto.RecipeDTO)
	for i := range resultRecipes {
		resultRow := &dto.RecipeDTO{
			RecipeID:          recipes[i].RecipeID,
			RecipeGroupID:     recipes[i].RecipeGroupID,
			Status:            recipes[i].Status,
			Name:              recipes[i].Name,
			ProductID:         recipes[i].ProductID,
			ProductName:       recipes[i].ProductName,
			ProductUnit:       recipes[i].ProductUnit,
			ProducedQuantity:  float64(recipes[i].ProducedQuantity) / 1000,
			Revision:          recipes[i].Revision,
			IsCurrent:         recipes[i].IsCurrent,
			CreatedByUserID:   recipes[i].CreatedByUserID,
			CreatedByUserName: recipes[i].CreatedByUserName,
			CreatedAt:         recipes[i].CreatedAt.Time,
			Ingredients:       make([]*dto.RecipeIngredientDTO, 0),
		}

		recipeIds = append(recipeIds, resultRow.RecipeID)
		resultRecipes[i] = resultRow
		hashMap[recipes[i].RecipeID] = resultRow
	}

	ingredients, err := r.db.GetRecipeIngredients(c.Context(), *tx, recipeIds)
	if err != nil {
		return err
	}

	for _, ing := range ingredients {
		if _, ok := hashMap[ing.RecipeID]; ok {
			r := hashMap[ing.RecipeID]
			r.Ingredients = append(r.Ingredients, dto.ToRecipeIngredientDTO(ing))
		}
	}

	return c.Status(fiber.StatusOK).JSON(struct {
		TotalCount int64            `json:"totalCount"`
		Items      []*dto.RecipeDTO `json:"items"`
	}{
		TotalCount: totalCount,
		Items:      resultRecipes,
	})
}

func (r *RouteManager) getRecipeByID(c *fiber.Ctx, tx *pgx.Tx) error {
	idParam := c.Params("id")
	recipeID := pgtype.UUID{}
	if err := recipeID.Scan(idParam); err != nil {
		return types.NewInvalidParamsError("invalid uuid on id url param")
	}

	recipeDTO, err := r.getRecipeAndIngredientsByID(c, tx, recipeID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(recipeDTO)
}

func (r *RouteManager) getRecipeByGroupID(c *fiber.Ctx, tx *pgx.Tx) error {
	idParam := c.Params("id")
	recipeGroupID := pgtype.UUID{}
	if err := recipeGroupID.Scan(idParam); err != nil {
		return types.NewInvalidParamsError("invalid uuid on id url param")
	}

	recipes, err := r.db.GetRecipesByGroupID(c.Context(), *tx, recipeGroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.NewNotFoundError()
		}
		return err
	}
	if err != nil {
		return err
	}
	recipeIds := make([]pgtype.UUID, 0)
	resultRecipes := make([]*dto.RecipeDTO, len(recipes))
	hashMap := make(map[pgtype.UUID]*dto.RecipeDTO)
	for i := range resultRecipes {
		resultRow := &dto.RecipeDTO{
			RecipeID:          recipes[i].RecipeID,
			RecipeGroupID:     recipes[i].RecipeGroupID,
			Status:            recipes[i].Status,
			Name:              recipes[i].Name,
			ProductID:         recipes[i].ProductID,
			ProductName:       recipes[i].ProductName,
			ProductUnit:       recipes[i].ProductUnit,
			ProducedQuantity:  float64(recipes[i].ProducedQuantity) / 1000,
			Revision:          recipes[i].Revision,
			IsCurrent:         recipes[i].IsCurrent,
			CreatedByUserID:   recipes[i].CreatedByUserID,
			CreatedByUserName: recipes[i].CreatedByUserName,
			CreatedAt:         recipes[i].CreatedAt.Time,
			Ingredients:       make([]*dto.RecipeIngredientDTO, 0),
		}

		recipeIds = append(recipeIds, resultRow.RecipeID)
		resultRecipes[i] = resultRow
		hashMap[recipes[i].RecipeID] = resultRow
	}

	ingredients, err := r.db.GetRecipeIngredients(c.Context(), *tx, recipeIds)
	if err != nil {
		return err
	}

	for _, ing := range ingredients {
		if _, ok := hashMap[ing.RecipeID]; ok {
			IngSlice := &hashMap[ing.RecipeID].Ingredients
			*IngSlice = append(*IngSlice, dto.ToRecipeIngredientDTO(ing))
		}
	}

	return c.Status(fiber.StatusOK).JSON(resultRecipes)
}

type CreateRecipeBody struct {
	Name             string      `json:"name" validate:"required,gte=3,lte=255"`
	ProductID        pgtype.UUID `json:"productId" validate:"required"`
	ProducedQuantity float64     `json:"producedQuantity" validate:"required,gte=1"`
	Ingredients      []struct {
		ProductID *pgtype.UUID `json:"productId" validate:"required"`
		Quantity  float64      `json:"quantity" validate:"required"`
	} `json:"ingredients" validate:"required,dive,required"`
}

func (r *RouteManager) createRecipe(c *fiber.Ctx, tx *pgx.Tx) error {
	userID, err := r.getCurrentUserId(c)
	if err != nil {
		return types.NewInvalidParamsError("invalid userId")
	}
	body := &CreateRecipeBody{}
	if err := c.BodyParser(body); err != nil {
		return types.NewInvalidParamsError("invalid body")
	}

	if err := r.validate.Struct(body); err != nil {
		return err
	}

	if !body.ProductID.Valid {
		return types.NewInvalidParamsError("invalid productID on reciped body")
	}

	recipeID, err := r.db.CreateRecipe(c.Context(), *tx, &postgres.CreateRecipeParams{
		Name:             body.Name,
		ProductID:        body.ProductID,
		ProducedQuantity: int64(body.ProducedQuantity * 1000),
		CreatedByUserID:  *userID,
	})
	if err != nil {
		return err
	}

	ingParams := make([]*postgres.CreateRecipeIngredientsParams, 0)
	for _, ing := range body.Ingredients {
		if ing.ProductID == nil {
			return types.NewInvalidParamsError("invalid product id")
		}
		ingParams = append(ingParams, &postgres.CreateRecipeIngredientsParams{
			RecipeID:  recipeID,
			ProductID: *ing.ProductID,
			Quantity:  int64(ing.Quantity * 1000),
		})
	}

	if _, err := r.db.CreateRecipeIngredients(c.Context(), *tx, ingParams); err != nil {
		return err
	}

	recipeDTO, err := r.getRecipeAndIngredientsByID(c, tx, recipeID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(recipeDTO)
}

type CreateRecipeRevisionBody struct {
	Name             string      `json:"name" validate:"required,gte=3,lte=255"`
	ProductID        pgtype.UUID `json:"productId" validate:"required"`
	ProducedQuantity float64     `json:"producedQuantity" validate:"required,gte=1"`
	Ingredients      []struct {
		ProductID *pgtype.UUID `json:"productId" validate:"required"`
		Quantity  float64      `json:"quantity" validate:"required"`
	} `json:"ingredients" validate:"required,dive,required"`
}

func (r *RouteManager) createRecipeRevision(c *fiber.Ctx, tx *pgx.Tx) error {
	userID, err := r.getCurrentUserId(c)
	if err != nil {
		return types.NewInvalidParamsError("invalid userId")
	}
	body := &CreateRecipeRevisionBody{}
	if err := c.BodyParser(body); err != nil {
		return types.NewInvalidParamsError("invalid body")
	}
	if err := r.validate.Struct(body); err != nil {
		return err
	}

	idParam := c.Params("id")
	recipeID := pgtype.UUID{}
	if err := recipeID.Scan(idParam); err != nil {
		return types.NewInvalidParamsError("invalid uuid on id url param")
	}

	prevRecipe, err := r.db.GetRecipeByID(c.Context(), *tx, recipeID)
	if err != nil {
		return err
	}

	if !prevRecipe.IsCurrent {
		return types.NewInvalidParamsError("recipe is not current")
	}

	_, err = r.db.SetCurrentFalse(c.Context(), *tx, recipeID)
	if err != nil {
		return err
	}

	recipeRevisionID, err := r.db.CreateRecipeRevision(c.Context(), *tx, &postgres.CreateRecipeRevisionParams{
		Name:             body.Name,
		Revision:         prevRecipe.Revision + 1,
		ProductID:        body.ProductID,
		ProducedQuantity: int64(body.ProducedQuantity * 1000),
		RecipeGroupID:    prevRecipe.RecipeGroupID,
		CreatedByUserID:  *userID,
	})
	if err != nil {
		return err
	}

	ingParams := make([]*postgres.CreateRecipeIngredientsParams, 0)
	for _, ing := range body.Ingredients {
		if ing.ProductID == nil {
			return types.NewInvalidParamsError("invalid product id")
		}
		ingParams = append(ingParams, &postgres.CreateRecipeIngredientsParams{
			RecipeID:  recipeRevisionID,
			ProductID: *ing.ProductID,
			Quantity:  int64(ing.Quantity * 1000),
		})
	}

	if _, err := r.db.CreateRecipeIngredients(c.Context(), *tx, ingParams); err != nil {
		return err
	}

	recipeDTO, err := r.getRecipeAndIngredientsByID(c, tx, recipeRevisionID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(recipeDTO)
}

type UpdateRecipeBody struct {
	Status postgres.Status `json:"status" validate:"required"`
}

func (r *RouteManager) updateRecipeStatus(c *fiber.Ctx, tx *pgx.Tx) error {
	body := &UpdateRecipeBody{}
	if err := c.BodyParser(body); err != nil {
		return types.NewInvalidParamsError("invalid body")
	}
	if err := r.validate.Struct(body); err != nil {
		return err
	}

	idParam := c.Params("id")
	groupID := pgtype.UUID{}
	if err := groupID.Scan(idParam); err != nil {
		return types.NewInvalidParamsError("invalid uuid on id url param")
	}

	if !body.Status.Valid() {
		return types.NewInvalidParamsError("invalid status")
	}

	err := r.db.SetRecipeStatusByGroupID(c.Context(), *tx, &postgres.SetRecipeStatusByGroupIDParams{
		Status:        body.Status,
		RecipeGroupID: groupID,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (r *RouteManager) getRecipeAndIngredientsByID(c *fiber.Ctx, tx *pgx.Tx, recipeID pgtype.UUID) (*dto.RecipeDTO, error) {
	recipe, err := r.db.GetRecipeByID(c.Context(), *tx, recipeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, types.NewNotFoundError()
		}
		return nil, err
	}

	recipeDTO := &dto.RecipeDTO{
		RecipeID:          recipe.RecipeID,
		RecipeGroupID:     recipe.RecipeGroupID,
		Status:            recipe.Status,
		Name:              recipe.Name,
		ProductID:         recipe.ProductID,
		ProductName:       recipe.ProductName,
		ProductUnit:       recipe.ProductUnit,
		ProducedQuantity:  float64(recipe.ProducedQuantity) / 1000,
		Revision:          recipe.Revision,
		IsCurrent:         recipe.IsCurrent,
		CreatedByUserID:   recipe.CreatedByUserID,
		CreatedByUserName: recipe.CreatedByUserName,
		CreatedAt:         recipe.CreatedAt.Time,
		Ingredients:       make([]*dto.RecipeIngredientDTO, 0),
	}

	ingredients, err := r.db.GetRecipeIngredients(c.Context(), *tx, []pgtype.UUID{recipe.RecipeID})
	if err != nil {
		return nil, err
	}

	for _, ing := range ingredients {
		recipeDTO.Ingredients = append(recipeDTO.Ingredients, dto.ToRecipeIngredientDTO(ing))
	}

	return recipeDTO, nil
}

func (r *RouteManager) RegisterRecipeRoutes() {
	r.app.Get("/recipes", r.dbWrapper.WithTransaction(r.getAllRecipes))
	r.app.Get("/recipes/:id", r.dbWrapper.WithTransaction(r.getRecipeByID))
	r.app.Get("/recipes_group/:id", r.dbWrapper.WithTransaction(r.getRecipeByGroupID))
	r.app.Post("/recipes", r.dbWrapper.WithTransaction(r.createRecipe))
	r.app.Put("/recipes/:id", r.dbWrapper.WithTransaction(r.createRecipeRevision))
	r.app.Patch("/recipes_group/:id/status", r.dbWrapper.WithTransaction(r.updateRecipeStatus))

}
