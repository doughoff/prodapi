package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hoffax/prodapi/postgres"
	"github.com/hoffax/prodapi/server/types"
)

type RouteManager struct {
	db        *postgres.Queries
	app       *fiber.App
	dbWrapper *types.DBWrapper
}

func NewRouteManager(app *fiber.App, db *postgres.Queries, dbWrapper *types.DBWrapper) *RouteManager {
	return &RouteManager{
		db:        db,
		app:       app,
		dbWrapper: dbWrapper,
	}
}
