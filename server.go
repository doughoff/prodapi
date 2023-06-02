package main

type Server struct {
	//Router *chi.Mux
	// You can add DB, config, etc. here
}

func CreateNewServer() *Server {
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
