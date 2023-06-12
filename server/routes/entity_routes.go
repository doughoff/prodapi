package routes

import (
	"bytes"
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/hoffax/prodapi/postgres"
	"github.com/hoffax/prodapi/server/dto"
	"github.com/hoffax/prodapi/server/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type getAllEntitiesQuery struct {
	StatusOptions []string `query:"status"`
	Search        string   `query:"search"`
	Limit         int32    `query:"limit"`
	Offset        int32    `query:"offset"`
}

func (r *RouteManager) getAllEntities(c *fiber.Ctx, tx *pgx.Tx) error {
	params := new(getAllEntitiesQuery)
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

	entities, err := r.db.GetEntities(c.Context(), *tx, &postgres.GetEntitiesParams{
		PageOffset:    params.Offset,
		PageLimit:     params.Limit,
		Search:        pgtype.Text(sql.NullString{String: params.Search, Valid: true}),
		StatusOptions: statusOptions,
	})
	if err != nil {
		return err
	}

	var totalCount int64
	if len(entities) > 0 {
		totalCount = entities[0].FullCount
	}
	resultRows := make([]*dto.EntityDTO, len(entities))
	for i := range resultRows {
		resultRows[i] = &dto.EntityDTO{
			ID:        entities[i].ID,
			Status:    entities[i].Status,
			Name:      entities[i].Name,
			Ci:        entities[i].Ci,
			Ruc:       entities[i].Ruc,
			CreatedAt: entities[i].CreatedAt.Time,
			UpdatedAt: entities[i].UpdatedAt.Time,
		}
	}

	return c.JSON(struct {
		TotalCount int64            `json:"totalCount"`
		Items      []*dto.EntityDTO `json:"items"`
	}{
		TotalCount: totalCount,
		Items:      resultRows,
	})
}

type createEntityBody struct {
	Name string `json:"name" validate:"required,gte=3,lte=255"`
	RUC  string `json:"ruc"`
	CI   string `json:"ci"`
}

func (r *RouteManager) createEntity(c *fiber.Ctx, tx *pgx.Tx) error {
	body := new(createEntityBody)
	if err := c.BodyParser(body); err != nil {
		return types.NewInvalidBodyError()
	}
	if err := r.validate.Struct(body); err != nil {
		return err
	}

	if len(body.CI) == 0 && len(body.RUC) == 0 {
		return types.NewInvalidParamsError("fields RUC or CI are required")
	}

	createEntityParams := &postgres.CreateEntityParams{
		Name: body.Name,
	}

	if len(body.CI) > 0 {
		createEntityParams.Ci = pgtype.Text(sql.NullString{String: body.CI, Valid: true})
		_, err := r.db.GetEntityByCI(c.Context(), *tx, createEntityParams.Ci)
		if err != nil {
			if err != pgx.ErrNoRows {
				return err
			}
		} else {
			return types.NewInvalidParamsError("entity with CI already exists")
		}
	}

	if len(body.RUC) > 0 {
		createEntityParams.Ruc = pgtype.Text(sql.NullString{String: body.RUC, Valid: true})
		_, err := r.db.GetEntityByRUC(c.Context(), *tx, createEntityParams.Ruc)
		if err != nil {
			if err != pgx.ErrNoRows {
				return err
			}
		} else {
			return types.NewInvalidParamsError("entity with RUC already exists")
		}
	}

	newEntity, err := r.db.CreateEntity(c.Context(), *tx, createEntityParams)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ToEntityDTO(newEntity))
}

type updateEntityBody struct {
	Status postgres.Status `json:"status" validate:"required"`
	Name   string          `json:"name" validate:"required,gte=3,lte=255"`
	RUC    string          `json:"ruc"`
	CI     string          `json:"ci"`
}

func (r *RouteManager) updateEntity(c *fiber.Ctx, tx *pgx.Tx) error {
	idParam := c.Params("id")
	entityID := pgtype.UUID{}
	if err := entityID.Scan(idParam); err != nil {
		return types.NewInvalidParamsError("invalid uuid on id url param")
	}

	body := new(updateEntityBody)
	if err := c.BodyParser(body); err != nil {
		return types.NewInvalidBodyError()
	}

	if err := r.validate.Struct(body); err != nil {
		return err
	}

	if len(body.CI) == 0 && len(body.RUC) == 0 {
		return types.NewInvalidParamsError("fields RUC or CI are required")
	}

	if !body.Status.Valid() {
		return types.NewInvalidParamsError("invalid value for status")
	}

	updateEntityParams := &postgres.UpdateEntityByIDParams{
		ID:     entityID,
		Name:   body.Name,
		Status: body.Status,
	}

	if len(body.CI) > 0 {
		updateEntityParams.Ci = pgtype.Text(sql.NullString{String: body.CI, Valid: true})
		entity, err := r.db.GetEntityByCI(c.Context(), *tx, updateEntityParams.Ci)
		if err != nil {
			if err != pgx.ErrNoRows {
				return err
			}
		} else {
			if !bytes.Equal(entity.ID.Bytes[:], entityID.Bytes[:]) {
				return types.NewInvalidParamsError("entity with CI already exists")
			}
		}
	}

	if len(body.RUC) > 0 {
		updateEntityParams.Ruc = pgtype.Text(sql.NullString{String: body.RUC, Valid: true})
		entity, err := r.db.GetEntityByRUC(c.Context(), *tx, updateEntityParams.Ruc)
		if err != nil {
			if err != pgx.ErrNoRows {
				return err
			}
		} else {
			if !bytes.Equal(entity.ID.Bytes[:], entityID.Bytes[:]) {
				return types.NewInvalidParamsError("entity with RUC already exists")
			}
		}
	}

	newEntity, err := r.db.UpdateEntityByID(c.Context(), *tx, updateEntityParams)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ToEntityDTO(newEntity))
}

func (r *RouteManager) getEntityById(c *fiber.Ctx, tx *pgx.Tx) error {
	idParam := c.Params("id")
	entityID := pgtype.UUID{}
	if err := entityID.Scan(idParam); err != nil {
		return types.NewInvalidParamsError("invalid uuid on id url param")
	}

	entity, err := r.db.GetEntityByID(c.Context(), *tx, entityID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.NewNotFoundError()
		}
		return err
	}

	return c.JSON(dto.ToEntityDTO(entity))
}

func (r *RouteManager) RegisterEntityRoutes() {
	g := r.app.Group("/entities")

	g.Get("/", r.dbWrapper.WithTransaction(r.getAllEntities))
	g.Post("/", r.dbWrapper.WithTransaction(r.createEntity))
	g.Put("/:id", r.dbWrapper.WithTransaction(r.updateEntity))
	g.Get("/:id", r.dbWrapper.WithTransaction(r.getEntityById))

}
