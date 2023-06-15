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
			CreatedByUserID:     stockMovements[i].CreatedByUserID,
			CreateByUserName:    stockMovements[i].CreateByUserName.String,
			CancelledByUserID:   stockMovements[i].CancelledByUserID,
			CancelledByUserName: stockMovements[i].CancelledByUserName.String,
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
			sm.Total += item.Price * item.Quantity
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
	Name string `json:"name" validate:"required,min=1,max=255"`
}

func (r *RouteManager) createStockMovement(c *fiber.Ctx, tx *pgx.Tx) error {
	return nil
}

func (r *RouteManager) updateStockMovement(c *fiber.Ctx, tx *pgx.Tx) error {
	return nil
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
	r.app.Get("/stock_movement/", r.dbWrapper.WithTransaction(r.getAllStockMovements))
	r.app.Get("/stock_movement/:id", r.dbWrapper.WithTransaction(r.getStockMovementByID))
	r.app.Post("/stock_movement/", r.dbWrapper.WithTransaction(r.createStockMovement))
	r.app.Put("/stock_movement/", r.dbWrapper.WithTransaction(r.updateStockMovement))
}
