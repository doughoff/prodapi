package server

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hoffax/prodapi/postgres"
	"github.com/hoffax/prodapi/server/config"
	"github.com/hoffax/prodapi/server/middleware"
	"github.com/hoffax/prodapi/server/routes"
	"github.com/hoffax/prodapi/server/types"
	"log"
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

	app.Use(logger.New())
	//app.Get("/metrics", monitor.New())

	//memoryStore := memory.New(memory.Config{
	//	GCInterval: 5 * time.Second,
	//})
	//app.Use(middleware.AuthMiddleware(memoryStore))

	routeManager := routes.NewRouteManager(app, db, dbWrapper, validate)
	routeManager.RegisterEntityRoutes()
	routeManager.RegisterUserRoutes()

	err = app.Listen(":3088")
	if err != nil {
		log.Fatalf("err listen: %v", err)
	}
}
