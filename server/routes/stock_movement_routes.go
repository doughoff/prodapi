package routes

import (
	"database/sql"
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
	for i := range resultRows {
		resultRows[i] = &dto.StockMovementDTO{
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
	return nil
}

func (r *RouteManager) createStockMovement(c *fiber.Ctx, tx *pgx.Tx) error {
	return nil
}

func (r *RouteManager) updateStockMovement(c *fiber.Ctx, tx *pgx.Tx) error {
	return nil
}

func (r *RouteManager) RegisterStockMovementRoutes() {
	r.app.Get("/stock_movement/", r.dbWrapper.WithTransaction(r.getAllStockMovements))
	r.app.Get("/stock_movement/:id", r.dbWrapper.WithTransaction(r.getStockMovementByID))
	r.app.Post("/stock_movement/", r.dbWrapper.WithTransaction(r.createStockMovement))
	r.app.Put("/stock_movement/", r.dbWrapper.WithTransaction(r.updateStockMovement))
}
