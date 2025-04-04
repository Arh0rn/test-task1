package usersController

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
	"test-task1/internal/controller/restapi/errors"
	"test-task1/internal/models"
)

type UserService interface {
	SignUp(*models.SignUpInput) (*models.User, error)
	GetAll() ([]*models.User, error)
	GetValidator() *validator.Validate
}

type UserController struct {
	service UserService
}

func New(service UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//TODO: add errUserAlreadyExists error handling
	var signUpInputDao SignUpInputDAO
	if err := json.NewDecoder(r.Body).Decode(&signUpInputDao); err != nil {
		rest_errors.HandleError(w, rest_errors.ErrBadRequest, http.StatusBadRequest)
	}

	v := c.service.GetValidator()
	//validate input
	if err := signUpInputDao.ValidateWith(v); err != nil {
		rest_errors.HandleError(w, err, http.StatusBadRequest)
		return
	}

	singUpInput := signUpInputDao.ToSignUpInput()

	user, err := c.service.SignUp(singUpInput)
	if err != nil {
		rest_errors.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		rest_errors.HandleError(w, rest_errors.ErrInternalServer, http.StatusInternalServerError)
		return
	}
}

func (c *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := c.service.GetAll()
	if err != nil {
		rest_errors.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	//TODO: return js object not a array
	if err := json.NewEncoder(w).Encode(users); err != nil {
		rest_errors.HandleError(w, rest_errors.ErrInternalServer, http.StatusInternalServerError)
		return
	}
}
