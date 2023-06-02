package main

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"path/filepath"
)

type Server struct {
	//Router *chi.Mux
	// You can add DB, config, etc. here
}

func CreateNewServer() *Server {
	//dbpool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	//if err != nil {
	//	log.Fatalf("Unable to connect to database: \n%+v\n", err)
	//}
	//defer dbpool.Close()
	//
	//dbconfig, err := pgxpool.ParseConfig(os.Getenv("DB_URL"))
	//if err != nil {
	//	log.Fatalf("Could not access dbconfig from pgxpool: \n%+v\n", err)
	//}
	//
	//dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
	//	pgxuuid.Register(conn.TypeMap())
	//	return nil
	//}
	//
	//userRepository := user.NewUserPgRepository(dbpool)
	//
	//users, err := userRepository.All(context.Background())
	//if err != nil {
	//	log.Fatalf("error getting all users %v\n", err)
	//}
	//
	//for _, u := range users {
	//	fmt.Printf("user: %+v\n", u)
	//}

	//db, err := sql.Open("pgx", os.Getenv("DB_URL"))
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	//	os.Exit(1)
	//}
	//defer db.Close()
	//
	//var greeting string
	//err = db.QueryRow("select 'Hello, world!'").Scan(&greeting)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	//	os.Exit(1)
	//}
	//
	//fmt.Println(greeting)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Construct the path to your migrations directory
	migrationsDir := filepath.Join(cwd, "db", "migrations")

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsDir),
		os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("error after...\n %v\n", err)
	}
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Printf("not really an error, just dont need to change anything.")

		} else {
			log.Fatal("hard error")
		}
	}
	//
	err = m.Drop() // reset database
	if err != nil {
		fmt.Printf("error while dropping database\n")
	}

	//userRepo := repository.NewUserRepository()
	//userService := service.NewUserService(userRepo)
	//userHandler := handler.NewUserHandler(userService)

	s := &Server{}
	//s.Router = chi.NewRouter()
	//
	//// Mount all Middleware here
	//s.Router.Use(middleware.Logger)
	//
	//// Mount all handlers here
	//userHandler.RegisterRoutes(s.Router)

	return s
}
