package routes

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/hoffax/prodapi/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *RouteManager) getAllEntities(c *fiber.Ctx, tx *pgx.Tx) error {
	// Implementation of the handler goes here...
	entities, err := r.db.GetEntities(c.Context(), *tx, &postgres.GetEntitiesParams{
		PageOffset:    0,
		PageLimit:     100,
		Search:        pgtype.Text(sql.NullString{String: "", Valid: true}),
		StatusOptions: []postgres.Status{postgres.StatusACTIVE, postgres.StatusINACTIVE},
	})
	if err != nil {
		return err
	}

	return c.JSON(entities)
}

func (r *RouteManager) RegisterEntityRoutes() {
	g := r.app.Group("/entities")

	g.Get("/", r.dbWrapper.WithTransaction(r.getAllEntities))
}
