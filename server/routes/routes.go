package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/memory"
	"github.com/hoffax/prodapi/postgres"
	"github.com/hoffax/prodapi/server/types"
)

type RouteManager struct {
	db           *postgres.Queries
	app          *fiber.App
	dbWrapper    *types.DBWrapper
	validate     *validator.Validate
	sessionStore *memory.Storage
}

func NewRouteManager(app *fiber.App, db *postgres.Queries, dbWrapper *types.DBWrapper, validate *validator.Validate, sessionStore *memory.Storage) *RouteManager {
	return &RouteManager{
		db:           db,
		app:          app,
		dbWrapper:    dbWrapper,
		validate:     validate,
		sessionStore: sessionStore,
	}
}
