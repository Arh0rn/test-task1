package app

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	cfg *config.Config
	ctx context.Context

	db        *sql.DB
	hasher    *hash.Hasher
	validator *validator.Validate

	userRepo       *postgresUsersRepo.UserRepository
	userService    *usersService.UserService
	userController *usersController.UserController

	handler *restapi.Handler
	router  *http.Handler
	server  *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	//TODO: Add structured logger

	log.Printf("Config loaded: %+v\n", cfg)

	db, err := databases.NewPostgresConnection(&cfg.Database)
	if err != nil {
		return nil, err
	}

	//TODO: Add cache in future

	log.Println("Database connection established")
	hasher := hash.New(cfg.HashCost)
	v := validate.New()

	jwtSecret := []byte(cfg.JWTSecret)
	atttl := cfg.AccessTokenTTL

	userRepository := postgresUsersRepo.New(db)
	userService := usersService.New(userRepository, hasher, v, jwtSecret, atttl)
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
		ctx:            ctx,
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

func (a *App) Run() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Starting server on", a.cfg.HTTPServer.Address)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(a.ctx, a.cfg.HTTPServer.ShutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	if err := a.db.Close(); err != nil {
		log.Fatalf("Database connection close error: %v", err)
	}

	log.Println("Server exited gracefully")
	return nil
}
