package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	redisUsersCache "test-task1/internal/cache/redis/users"
	"test-task1/internal/controller/restapi"
	usersController "test-task1/internal/controller/restapi/controllers/users"
	"test-task1/internal/databases"
	postgresUsersRepo "test-task1/internal/repository/postgres/users"
	usersService "test-task1/internal/service/users"
	"test-task1/pkg/config"
	"test-task1/pkg/hash"
	"test-task1/pkg/logger"
	"test-task1/pkg/validate"
)

type App struct {
	cfg *config.Config
	ctx context.Context
	log *slog.Logger

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

	log := logger.InitLogger(cfg.Env)
	slog.SetDefault(log) //No need to inject logger to every layer ^_^

	log.Debug(fmt.Sprintf("%+v", cfg))
	db, err := databases.NewPostgresConnection(&cfg.Database)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to database", "error", err)
		return nil, err
	}

	cache, err := databases.NewRedisClient(&cfg.Cache)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to Redis", "error", err)
		return nil, err
	}

	hasher := hash.New(cfg.HashCost)
	v := validate.New()

	jwtSecret := []byte(cfg.JWTSecret)
	atttl := cfg.AccessTokenTTL

	userRepository := postgresUsersRepo.New(db)
	userCache := redisUsersCache.New(cache, cfg.Cache.TTL)
	userService := usersService.New(userRepository, userCache, hasher, v, jwtSecret, atttl)
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
		log:            log,
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
		a.log.Info("Starting server", "address", a.cfg.HTTPServer.Address)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("Server starting error", "error", err)
		}
	}()
	<-quit
	a.log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(a.ctx, a.cfg.HTTPServer.ShutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Error("Server shutdown error", "error", err)
	}

	if err := a.db.Close(); err != nil {
		a.log.Error("Database connection close error", "error", err)
	}

	a.log.Info("Server exited gracefully")
	return nil
}
