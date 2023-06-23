package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/hoffax/prodapi/postgres"
	"github.com/hoffax/prodapi/server/dto"
	"github.com/hoffax/prodapi/server/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type getAllStockMovementsQuery struct {
	StatusOptions []string  `query:"status"`
	Search        string    `query:"search"`
	StartDate     time.Time `query:"startDate"`
	Limit         int32     `query:"limit"`
	Offset        int32     `query:"offset"`
}

func (r *RouteManager) getAllStockMovements(c *fiber.Ctx, tx *pgx.Tx) error {
	params := new(getAllStockMovementsQuery)
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

	stockMovements, err := r.db.GetStockMovements(c.Context(), *tx, &postgres.GetStockMovementsParams{
		StatusOptions: statusOptions,
		Search:        pgtype.Text(sql.NullString{String: params.Search, Valid: true}),
		StartDate:     pgtype.Date{Time: params.StartDate, Valid: true},
		PageLimit:     params.Limit,
		PageOffset:    params.Offset,
	})
	if err != nil {
		return err
	}

	var totalCount int64

	if len(stockMovements) > 0 {
		totalCount = stockMovements[0].FullCount
	}

	resultRows := make([]*dto.StockMovementDTO, len(stockMovements))
	resultMap := make(map[pgtype.UUID]*dto.StockMovementDTO)
	stockMovementIDS := make([]pgtype.UUID, 0)
	for i := range resultRows {
		resultRow := &dto.StockMovementDTO{
			ID:                  stockMovements[i].ID,
			Status:              stockMovements[i].Status,
			Type:                stockMovements[i].Type,
			Date:                stockMovements[i].Date.Time,
			EntityID:            stockMovements[i].EntityID,
			EntityName:          stockMovements[i].EntityName.String,
			DocumentNumber:      stockMovements[i].DocumentNumber.String,
			CreatedByUserID:     stockMovements[i].CreatedByUserID,
			CreateByUserName:    stockMovements[i].CreateByUserName.String,
			CancelledByUserID:   stockMovements[i].CancelledByUserID,
			CancelledByUserName: stockMovements[i].CancelledByUserName.String,
			Items:               make([]*dto.StockMovementItemDTO, 0),
		}
		resultRows[i] = resultRow
		stockMovementIDS = append(stockMovementIDS, resultRow.ID)
		resultMap[resultRow.ID] = resultRow
	}

	items, err := r.db.GetStockMovementItems(c.Context(), *tx, stockMovementIDS)
	if err != nil {
		return err
	}

	for _, item := range items {
		if sm, ok := resultMap[item.StockMovementID]; ok {
			sm.Total += item.Price * item.Quantity / 1000
		}
	}

	return c.JSON(struct {
		TotalCount int64                   `json:"totalCount"`
		Items      []*dto.StockMovementDTO `json:"items"`
	}{
		TotalCount: totalCount,
		Items:      resultRows,
	})
}

func (r *RouteManager) getStockMovementByID(c *fiber.Ctx, tx *pgx.Tx) error {
	idParam := c.Params("id")
	movementID := pgtype.UUID{}
	if err := movementID.Scan(idParam); err != nil {
		return types.NewInvalidParamsError("invalid uuid on id url param")
	}

	movementDTO, err := r.getMovementByID(c, tx, movementID)
	if err != nil {
		return err
	}

	return c.JSON(movementDTO)
}

type CreateMovementBody struct {
	Type           postgres.MovementType `json:"type" validate:"required"`
	Date           time.Time             `json:"date" validate:"required"`
	EntityID       pgtype.UUID           `json:"entityId"`
	DocumentNumber string                `json:"documentNumber"`
	Items          []*struct {
		ProductID *pgtype.UUID `json:"productId" validate:"required"`
		Quantity  float64      `json:"quantity" validate:"required"`
		Price     int64        `json:"price" validate:"required"`
		Batch     string       `json:"batch"`
	} `json:"items"`
}

func (r *RouteManager) createStockMovement(c *fiber.Ctx, tx *pgx.Tx) error {
	userID, err := r.getCurrentUserId(c)
	if err != nil {
		return types.NewInvalidParamsError("invalid userId")
	}
	body := new(CreateMovementBody)
	if err := c.BodyParser(body); err != nil {
		fmt.Printf("err: %+v\n", err)
		return types.NewInvalidBodyError()
	}
	if err := r.validate.Struct(body); err != nil {
		return err
	}

	if !body.Type.Valid() {
		return types.NewInvalidParamsError("invalid value for movement type")
	}

	movementID, err := r.db.CreateStockMovement(c.Context(), *tx, &postgres.CreateStockMovementParams{
		Type:            body.Type,
		EntityID:        body.EntityID,
		Date:            pgtype.Date{Time: body.Date, Valid: true},
		DocumentNumber:  pgtype.Text(sql.NullString{String: body.DocumentNumber, Valid: len(body.DocumentNumber) > 0}),
		CreatedByUserID: *userID,
	})
	if err != nil {
		return err
	}

	movItemParams := make([]*postgres.CreateStockMovementItemsParams, 0)
	for _, item := range body.Items {
		batch := sql.NullString{Valid: false}
		if item.Batch != "" {
			batch = sql.NullString{String: item.Batch, Valid: true}
		}

		movItemParams = append(movItemParams, &postgres.CreateStockMovementItemsParams{
			StockMovementID: movementID,
			ProductID:       *item.ProductID,
			Quantity:        int64(item.Quantity * 1000),
			Price:           item.Price,
			Batch:           pgtype.Text(batch),
		})
	}

	_, err = r.db.CreateStockMovementItems(c.Context(), *tx, movItemParams)
	if err != nil {
		return err
	}

	createdMovement, err := r.getMovementByID(c, tx, movementID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(createdMovement)
}

type UpdateMovementBody struct {
	Status         postgres.Status `json:"status" validate:"required"`
	Date           time.Time       `json:"date" validate:"required"`
	EntityID       pgtype.UUID     `json:"entityId" validate:"required"`
	DocumentNumber string          `json:"documentNumber"`
}

func (r *RouteManager) updateStockMovement(c *fiber.Ctx, tx *pgx.Tx) error {
	body := new(UpdateMovementBody)
	if err := c.BodyParser(body); err != nil {
		return types.NewInvalidBodyError()
	}
	if err := r.validate.Struct(body); err != nil {
		return err
	}

	idParam := c.Params("id")
	movementID := pgtype.UUID{}
	if err := movementID.Scan(idParam); err != nil {
		return types.NewInvalidParamsError("invalid uuid on id url param")
	}

	err := r.db.UpdateStockMovement(c.Context(), *tx, &postgres.UpdateStockMovementParams{
		ID:             movementID,
		Status:         body.Status,
		Date:           pgtype.Date{Time: body.Date, Valid: true},
		EntityID:       body.EntityID,
		DocumentNumber: pgtype.Text(sql.NullString{String: body.DocumentNumber, Valid: len(body.DocumentNumber) > 0}),
	})
	if err != nil {
		return err
	}

	updatedMovement, err := r.getMovementByID(c, tx, movementID)
	if err != nil {
		return err
	}

	return c.JSON(updatedMovement)
}

func (r *RouteManager) getMovementByID(c *fiber.Ctx, tx *pgx.Tx, uuid pgtype.UUID) (*dto.StockMovementDTO, error) {
	movement, err := r.db.GetStockMovementByID(c.Context(), *tx, uuid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, types.NewNotFoundError()
		}
		return nil, err
	}

	movementDTO := dto.ToStockMovementDTO(movement)

	items, err := r.db.GetStockMovementItems(c.Context(), *tx, []pgtype.UUID{movement.ID})
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		itemDTO := dto.ToStockMovementItemDTO(item)
		movementDTO.Total += itemDTO.Total
		movementDTO.Items = append(movementDTO.Items, itemDTO)
	}

	return movementDTO, nil
}

func (r *RouteManager) RegisterStockMovementRoutes() {
	r.app.Get("/stock_movements/", r.dbWrapper.WithTransaction(r.getAllStockMovements))
	r.app.Get("/stock_movements/:id", r.dbWrapper.WithTransaction(r.getStockMovementByID))
	r.app.Post("/stock_movements/", r.dbWrapper.WithTransaction(r.createStockMovement))
	r.app.Put("/stock_movements/:id", r.dbWrapper.WithTransaction(r.updateStockMovement))
}
