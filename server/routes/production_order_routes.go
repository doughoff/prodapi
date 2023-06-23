package routes

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hoffax/prodapi/postgres"
	"github.com/hoffax/prodapi/server/dto"
	"github.com/hoffax/prodapi/server/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// production order should be like this :
/*
	# Production Order
	  - id
 	  - recipe //etc, other fields.

	  -> Cycles
		- All production Order Cycles (based on recipe/targetAmount)

	  -> Movements
		- All stock movements "from the cycles" (use the stockMovement Query to bind things together)


	// TODO: Probably is needed to refactor a little the stock movements part to make something "reusable" for this.
*/
type getAllProductionOrdersQuery struct {
	StatusOptions []string  `query:"status"`
	Search        string    `query:"search"`
	StartDate     time.Time `query:"startDate"`
	Limit         int32     `query:"limit"`
	Offset        int32     `query:"offset"`
}

func (r *RouteManager) getAllProductionOrders(c *fiber.Ctx, tx *pgx.Tx) error {
	params := new(getAllProductionOrdersQuery)
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

	productionOrders, err := r.db.GetProductionOrders(c.Context(), *tx, &postgres.GetProductionOrdersParams{
		StatusOptions: statusOptions,
		Search:        pgtype.Text(sql.NullString{String: params.Search, Valid: true}),
		StartDate:     pgtype.Timestamp{Time: params.StartDate, Valid: true},
		PageLimit:     params.Limit,
		PageOffset:    params.Offset,
	})
	if err != nil {
		return err
	}

	var totalCount int64
	if len(productionOrders) > 0 {
		totalCount = productionOrders[0].FullCount
	}

	resultRows := make([]*dto.ProductionOrderDTO, len(productionOrders))
	for i := range resultRows {
		resultRow := &dto.ProductionOrderDTO{
			ID:                  productionOrders[i].ID,
			Status:              productionOrders[i].Status,
			ProductionStep:      productionOrders[i].ProductionStep,
			Code:                productionOrders[i].Code.String,
			Cycles:              float64(productionOrders[i].Cycles / 1000),
			TargetAmount:        float64((productionOrders[i].Cycles * productionOrders[i].ProducedQuantity.Int64) / 1000 / 1000),
			Output:              float64(productionOrders[i].Output.Int64 / 1000),
			RecipeID:            productionOrders[i].RecipeID,
			RecipeName:          productionOrders[i].RecipeName.String,
			CreatedByUserID:     productionOrders[i].CreatedByUserID,
			CreateByUserName:    productionOrders[i].CreateByUserName.String,
			CancelledByUserID:   productionOrders[i].CancelledByUserID,
			CancelledByUserName: productionOrders[i].CancelledByUserName.String,
			CreatedAt:           productionOrders[i].CreatedAt.Time,
			UpdatedAt:           productionOrders[i].UpdatedAt.Time,
		}
		resultRows[i] = resultRow
	}

	return c.JSON(struct {
		TotalCount int64                     `json:"totalCount"`
		Items      []*dto.ProductionOrderDTO `json:"items"`
	}{
		TotalCount: totalCount,
		Items:      resultRows,
	})
}

func (r *RouteManager) getProductionOrderByID(c *fiber.Ctx, tx *pgx.Tx) error {
	idParam := c.Params("id")
	productionOrderID := pgtype.UUID{}
	if err := productionOrderID.Scan(idParam); err != nil {
		return types.NewInvalidParamsError("invalid uuid on id url param")
	}

	productionOrder, err := r.db.GetProductionOrderByID(c.Context(), *tx, productionOrderID)
	if err != nil {
		return err
	}

	productionOrderDTO := dto.ProductionOrderDTO{
		ID:                  productionOrder.ID,
		Status:              productionOrder.Status,
		ProductionStep:      productionOrder.ProductionStep,
		Code:                productionOrder.Code.String,
		Cycles:              float64(productionOrder.Cycles) / 1000,
		TargetAmount:        float64(productionOrder.Cycles*productionOrder.ProducedQuantity.Int64) / 1000 / 1000,
		Output:              float64(productionOrder.Output.Int64) / 1000,
		RecipeID:            productionOrder.RecipeID,
		RecipeName:          productionOrder.RecipeName.String,
		CreatedByUserID:     productionOrder.CreatedByUserID,
		CreateByUserName:    productionOrder.CreateByUserName.String,
		CancelledByUserID:   productionOrder.CancelledByUserID,
		CancelledByUserName: productionOrder.CancelledByUserName.String,
		CreatedAt:           productionOrder.CreatedAt.Time,
		UpdatedAt:           productionOrder.UpdatedAt.Time,
	}

	orderCycles, err := r.db.GetProductionOrderCycles(c.Context(), *tx, []pgtype.UUID{productionOrderID})
	if err != nil {
		return err
	}

	cycles := make([]*dto.ProductionOrderCycleDTO, len(orderCycles))
	for i := range orderCycles {
		cycles[i] = &dto.ProductionOrderCycleDTO{
			ID:             orderCycles[i].ID,
			Factor:         float64(orderCycles[i].Factor) / 1000,
			ProductionStep: orderCycles[i].ProductionStep,
			CompletedAt:    &orderCycles[i].CompletedAt.Time,
		}
	}

	productionOrderMovements, err := r.db.GetProductionOrderMovements(c.Context(), *tx, productionOrderID)
	if err != nil {
		return err
	}

	movementToCycleMap := make(map[pgtype.UUID]pgtype.UUID)
	for _, movement := range productionOrderMovements {
		movementToCycleMap[movement.StockMovementID] = movement.ProductionOrderCycleID
	}

	movementIDS := make([]pgtype.UUID, len(productionOrderMovements))
	for i := range productionOrderMovements {
		movementIDS[i] = productionOrderMovements[i].StockMovementID
	}

	stockMovements, err := r.db.GetStockMovementsByIDS(c.Context(), *tx, movementIDS)
	if err != nil {
		return err
	}

	movements := make([]*dto.ProductionOrderMovementDTO, len(stockMovements))
	movementsMap := make(map[pgtype.UUID]*dto.ProductionOrderMovementDTO)
	stockMovementIDS := make([]pgtype.UUID, 0)
	for i := range movements {
		//get cycleID
		cycleID, ok := movementToCycleMap[stockMovements[i].ID]
		if !ok {
			continue
		}

		mov := &dto.ProductionOrderMovementDTO{
			CycleID: cycleID,
			StockMovementDTO: dto.StockMovementDTO{
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
			},
		}
		movements[i] = mov
		stockMovementIDS = append(stockMovementIDS, mov.ID)
		movementsMap[mov.ID] = mov
	}

	items, err := r.db.GetStockMovementItems(c.Context(), *tx, stockMovementIDS)
	if err != nil {
		return err
	}

	for _, item := range items {
		if sm, ok := movementsMap[item.StockMovementID]; ok {
			itemDTO := dto.ToStockMovementItemDTO(item)
			sm.Total += itemDTO.Total
			sm.Items = append(sm.Items, itemDTO)
		}
	}

	productionOrderDTO.Movements = movements
	productionOrderDTO.CycleInstances = cycles

	return c.JSON(productionOrderDTO)
}

type CreateProductionOrder struct {
	RecipeID pgtype.UUID `json:"recipe_id" validate:"required"`
	Cycles   float64     `json:"cycles" validate:"required"`
}

func (r *RouteManager) createProductionOrder(c *fiber.Ctx, tx *pgx.Tx) error {
	userID, err := r.getCurrentUserId(c)
	if err != nil {
		return types.NewInvalidParamsError("invalid userId")
	}
	body := new(CreateProductionOrder)
	if err := c.BodyParser(body); err != nil {
		fmt.Printf("err: %+v\n", err)
		return types.NewInvalidBodyError()
	}
	if err := r.validate.Struct(body); err != nil {
		return err
	}

	productionOrderID, err := r.db.CreateProductionOrder(c.Context(), *tx, &postgres.CreateProductionOrderParams{
		RecipeID:        body.RecipeID,
		Cycles:          int64(body.Cycles * 1000),
		CreatedByUserID: *userID,
	})
	if err != nil {
		return err
	}

	cycles := make([]*postgres.CreateProductionOrderCyclesParams, 0)
	remainingCycles := int64(body.Cycles * 1000)
	order := int64(1)
	for remainingCycles > 0 {
		factor := int64(1000)
		if remainingCycles < 1000 {
			factor = remainingCycles
		}

		cycle := &postgres.CreateProductionOrderCyclesParams{
			ProductionOrderID: productionOrderID,
			Factor:            factor,
			CycleOrder:        order,
		}
		cycles = append(cycles, cycle)
		remainingCycles -= 1000
		order++
	}

	if _, err := r.db.CreateProductionOrderCycles(c.Context(), *tx, cycles); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(map[string]interface{}{
		"id": productionOrderID,
	})
}

type UpdateProductionOrder struct {
	Status         postgres.Status         `json:"status" validate:"required"`
	ProductionStep postgres.ProductionStep `json:"production_step" validate:"required"`
}

func (r *RouteManager) updateProductionOrder(c *fiber.Ctx, tx *pgx.Tx) error {
	body := new(UpdateProductionOrder)
	if err := c.BodyParser(body); err != nil {
		return types.NewInvalidBodyError()
	}
	if err := r.validate.Struct(body); err != nil {
		return err
	}

	idParam := c.Params("id")
	productionOrderID := pgtype.UUID{}
	if err := productionOrderID.Scan(idParam); err != nil {
		return types.NewInvalidParamsError("invalid uuid on id url param")
	}

	if !body.Status.Valid() {
		return types.NewInvalidParamsError("invalid value for Status")
	}

	if !body.ProductionStep.Valid() {
		return types.NewInvalidParamsError("invalid value for ProductionStep")
	}

	err := r.db.UpdateProductionOrder(c.Context(), *tx, &postgres.UpdateProductionOrderParams{
		Status:         body.Status,
		ProductionStep: body.ProductionStep,
		ID:             productionOrderID,
	})
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

type CreateCycleMovement struct {
	CycleID pgtype.UUID `json:"cycleId" validate:"required"`
	Items   []*struct {
		ProductID *pgtype.UUID `json:"productId" validate:"required"`
		Quantity  float64      `json:"quantity" validate:"required"`
		Price     int64        `json:"price" validate:"required"`
		Batch     string       `json:"batch"`
	} `json:"items"`
}

func (r *RouteManager) CreateOrderCycleMovement(c *fiber.Ctx, tx *pgx.Tx) error {
	userID, err := r.getCurrentUserId(c)
	if err != nil {
		return types.NewInvalidParamsError("invalid userId")
	}
	body := new(CreateCycleMovement)
	if err := c.BodyParser(body); err != nil {
		return types.NewInvalidBodyError()
	}
	if err := r.validate.Struct(body); err != nil {
		return err
	}

	idParam := c.Params("id")
	productionOrderID := pgtype.UUID{}
	if err := productionOrderID.Scan(idParam); err != nil {
		return types.NewInvalidParamsError("invalid uuid on id url param")
	}

	productionOrder, err := r.db.GetProductionOrderByID(c.Context(), *tx, productionOrderID)
	if err != nil {
		return err
	}

	if productionOrder.Status != postgres.StatusACTIVE {
		return types.NewInvalidParamsError("Cannot change a production order that is not active")
	}

	if productionOrder.ProductionStep == postgres.ProductionStepCOMPLETED {
		return types.NewInvalidParamsError("Cannot change a production order that is completed")
	}

	if productionOrder.ProductionStep == postgres.ProductionStepPENDING {
		err = r.db.UpdateProductionOrder(c.Context(), *tx, &postgres.UpdateProductionOrderParams{
			ID:             productionOrderID,
			Status:         postgres.StatusACTIVE,
			ProductionStep: postgres.ProductionStepINPROGRESS,
		})
		if err != nil {
			return err
		}
	}

	cycle, err := r.db.GetProductionOrderCycleByID(c.Context(), *tx, body.CycleID)
	if err != nil {
		return err
	}

	if cycle.ProductionOrderID != productionOrderID {
		return types.NewInvalidParamsError("Cycle does not belong to production order")
	}

	if cycle.ProductionStep == postgres.ProductionStepCOMPLETED {
		return types.NewInvalidParamsError("Cannot change a cycle that is completed")
	}

	if cycle.ProductionStep == postgres.ProductionStepPENDING {
		err = r.db.UpdateProductionOrderCycle(c.Context(), *tx, &postgres.UpdateProductionOrderCycleParams{
			ID:             body.CycleID,
			ProductionStep: postgres.ProductionStepINPROGRESS,
			CompletedAt:    pgtype.Timestamp{Valid: false},
		})
		if err != nil {
			return err
		}
	}

	movementID, err := r.db.CreateStockMovement(c.Context(), *tx, &postgres.CreateStockMovementParams{
		Type:            postgres.MovementTypePRODUCTIONIN,
		Date:            pgtype.Date{Time: time.Now(), Valid: true},
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

	_, err = r.db.CreateOrderCycleMovement(c.Context(), *tx, &postgres.CreateOrderCycleMovementParams{
		CycleID:    body.CycleID,
		MovementID: movementID,
	})

	return c.SendStatus(fiber.StatusNoContent)
}

func (r *RouteManager) CompleteOrderCycle(c *fiber.Ctx, tx *pgx.Tx) error {
	productionOrderIDParam := c.Params("id")
	productionOrderID := pgtype.UUID{}
	if err := productionOrderID.Scan(productionOrderIDParam); err != nil {
		return types.NewInvalidParamsError("invalid uuid on id url param")
	}

	cycleIDParam := c.Params("id")
	cycleID := pgtype.UUID{}
	if err := cycleID.Scan(cycleIDParam); err != nil {
		return types.NewInvalidParamsError("invalid uuid on id url param")
	}

	productionOrder, err := r.db.GetProductionOrderByID(c.Context(), *tx, productionOrderID)
	if err != nil {
		return err
	}

	if productionOrder.Status != postgres.StatusACTIVE {
		return types.NewInvalidParamsError("Cannot change a production order that is not active")
	}

	if productionOrder.ProductionStep == postgres.ProductionStepCOMPLETED {
		return types.NewInvalidParamsError("Cannot change a production order that is completed")
	}

	if productionOrder.ProductionStep == postgres.ProductionStepPENDING {
		err = r.db.UpdateProductionOrder(c.Context(), *tx, &postgres.UpdateProductionOrderParams{
			ID:             productionOrderID,
			Status:         postgres.StatusACTIVE,
			ProductionStep: postgres.ProductionStepINPROGRESS,
		})
		if err != nil {
			return err
		}
	}

	cycle, err := r.db.GetProductionOrderCycleByID(c.Context(), *tx, cycleID)
	if err != nil {
		return err
	}

	if cycle.ProductionOrderID != productionOrderID {
		return types.NewInvalidParamsError("Cycle does not belong to production order")
	}

	err = r.db.UpdateProductionOrderCycle(c.Context(), *tx, &postgres.UpdateProductionOrderCycleParams{
		ID:             cycleID,
		ProductionStep: postgres.ProductionStepCOMPLETED,
		CompletedAt:    pgtype.Timestamp{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (r *RouteManager) RegisterProductionOrderRoutes() {
	r.app.Get("/production_orders", r.dbWrapper.WithTransaction(r.getAllProductionOrders))
	r.app.Get("/production_orders/:id", r.dbWrapper.WithTransaction(r.getProductionOrderByID))
	r.app.Post("/production_orders", r.dbWrapper.WithTransaction(r.createProductionOrder))
	r.app.Put("/production_orders/:id", r.dbWrapper.WithTransaction(r.updateProductionOrder))
	r.app.Post("/production_orders/:id/movements", r.dbWrapper.WithTransaction(r.CreateOrderCycleMovement))
	r.app.Post("/production_orders/:id/complete_cycle/:cycle_id", r.dbWrapper.WithTransaction(r.CompleteOrderCycle))
}
