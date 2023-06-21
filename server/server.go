package server

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/hoffax/prodapi/postgres"
	"github.com/hoffax/prodapi/server/config"
	"github.com/hoffax/prodapi/server/middleware"
	"github.com/hoffax/prodapi/server/routes"
	"github.com/hoffax/prodapi/server/types"
	"log"
	"time"
)

func Serve() {
	conn := config.NewPgxConn()
	err := conn.Ping(context.Background())
	if err != nil {
		fmt.Printf("err ping db")
	}

	validate := validator.New()
	db := postgres.New()
	dbWrapper := &types.DBWrapper{DB: conn}
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.FiberCustomErrorHandler,
	})

	store := session.New(session.Config{
		Expiration: 72 * time.Hour,
	})
	app.Use(logger.New())
	//app.Get("/metrics", monitor.New())

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://127.0.0.1:5173",
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))

	//memoryStore := memory.New(memory.Config{
	//	GCInterval: 5 * time.Hour,
	//})
	app.Use(middleware.AuthMiddleware(store))

	routeManager := routes.NewRouteManager(app, db, dbWrapper, validate, store)
	routeManager.RegisterAuthRoutes()
	routeManager.RegisterEntityRoutes()
	routeManager.RegisterUserRoutes()
	routeManager.RegisterProductRoutes()
	routeManager.RegisterRecipeRoutes()
	routeManager.RegisterStockMovementRoutes()

	err = app.Listen(":3088")
	if err != nil {
		log.Fatalf("err listen: %v", err)
	}
}
