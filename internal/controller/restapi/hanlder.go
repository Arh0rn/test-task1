package restapi

import (
	usersController "github.com/Arh0rn/test-task1/internal/controller/restapi/controllers/users"
	"github.com/Arh0rn/test-task1/internal/controller/restapi/middlewares"
	"github.com/Arh0rn/test-task1/internal/controller/restapi/swagger"
	"github.com/Arh0rn/test-task1/pkg/config"
	"net/http"
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
		middlewares.SetCORS,
		middlewares.LoggerMiddleware,
	)

	baseRouter := http.NewServeMux()
	authorizedRouter := http.NewServeMux()

	baseRouter.HandleFunc("GET /swagger/", swagger.Set(cfg))

	baseRouter.HandleFunc("POST /users", h.UserController.SignUp)
	baseRouter.HandleFunc("POST /login", h.UserController.Login)

	authorizedRouter.HandleFunc("GET /users", h.UserController.GetAll)
	authorizedRouter.HandleFunc("GET /users/{id}", h.UserController.GetByID)
	authorizedRouter.HandleFunc("PUT /users/{id}", h.UserController.UpdateByID)
	authorizedRouter.HandleFunc("DELETE /users/{id}", h.UserController.DeleteByID)

	baseRouter.Handle("/", middlewares.AuthMiddleware(cfg.JWTSecret)(authorizedRouter))

	router := mainStack(baseRouter)
	return &router
}
