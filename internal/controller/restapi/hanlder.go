package restapi

import (
	"net/http"
	usersController "test-task1/internal/controller/restapi/controllers/users"
	"test-task1/internal/controller/restapi/middlewares"
	"test-task1/pkg/config"
)

type Handler struct {
	UserController usersController.UserController
}

func NewHandler(userController *usersController.UserController) *Handler {
	return &Handler{
		UserController: *userController,
	}
}

func (h *Handler) InitRoutes(cfg *config.HTTPServer) *http.Handler {
	mainStack := middlewares.CreateMiddlewareStack(
		middlewares.LoggerMiddleware,
	)

	baseRouter := http.NewServeMux()
	authorizedRouter := http.NewServeMux()

	baseRouter.HandleFunc("POST /users", h.UserController.SignUp)
	baseRouter.HandleFunc("POST /login", h.UserController.Login)

	authorizedRouter.HandleFunc("GET /users", h.UserController.GetAll)

	baseRouter.Handle("/", middlewares.AuthMiddleware(cfg.JWTSecret)(authorizedRouter))

	router := mainStack(baseRouter)
	return &router
}
