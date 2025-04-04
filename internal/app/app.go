package app

import (
	"context"
	"database/sql"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"test-task1/internal/controller/restapi"
	usersController "test-task1/internal/controller/restapi/controllers/users"
	"test-task1/internal/databases"
	postgresUsersRepo "test-task1/internal/repository/postgres/users"
	usersService "test-task1/internal/service/users"
	"test-task1/pkg/config"
	"test-task1/pkg/hash"
	"test-task1/pkg/validate"
)

type App struct {
	cfg       *config.Config
	db        *sql.DB
	hasher    *hash.Hasher
	validator *validator.Validate

	userRepo       *postgresUsersRepo.UserRepository
	userService    *usersService.UserService
	userController *usersController.UserController
	handler        *restapi.Handler
	router         *http.Handler
	server         *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	log.Printf("Config loaded: %+v\n", cfg)

	db, err := databases.NewPostgresConnection(&cfg.Database)
	if err != nil {
		return nil, err
	}

	log.Println("Database connection established")
	hasher := hash.New(cfg.HashCost)
	v := validate.New()

	//TODO: add jwt
	//jwtSecret := []byte(cfg.JWTSecret)

	userRepository := postgresUsersRepo.New(db)
	userService := usersService.New(userRepository, hasher, v)
	userController := usersController.New(userService)
	handler := restapi.NewHandler(userController)
	router := handler.InitRoutes(&cfg.HTTPServer)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      *router,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	app := &App{
		cfg:            cfg,
		db:             db,
		hasher:         hasher,
		validator:      v,
		userRepo:       userRepository,
		userService:    userService,
		userController: userController,
		handler:        handler,
		router:         router,
		server:         srv,
	}

	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	defer a.db.Close()
	log.Println("Starting server...")
	err := a.server.ListenAndServe()
	if err != nil {
		log.Printf("Server error: %v\n", err)
		return err
	}
	return nil
}
