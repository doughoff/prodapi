package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/hoffax/prodapi/postgres"
	"github.com/hoffax/prodapi/server/types"
	"github.com/jackc/pgx/v5/pgtype"
)

type RouteManager struct {
	db           *postgres.Queries
	app          *fiber.App
	dbWrapper    *types.DBWrapper
	validate     *validator.Validate
	sessionStore *session.Store
}

func NewRouteManager(app *fiber.App, db *postgres.Queries, dbWrapper *types.DBWrapper, validate *validator.Validate, sessionStore *session.Store) *RouteManager {
	return &RouteManager{
		db:           db,
		app:          app,
		dbWrapper:    dbWrapper,
		validate:     validate,
		sessionStore: sessionStore,
	}
}

func (r *RouteManager) getCurrentUserId(c *fiber.Ctx) (*pgtype.UUID, error) {
	userIDBytes := c.Locals("userId").([]byte)
	var userID pgtype.UUID
	err := userID.UnmarshalJSON(userIDBytes)
	if err != nil {
		return nil, types.NewInvalidParamsError("invalid userId")
	}

	return &userID, nil
}
