package usersController

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"test-task1/internal/controller/restapi/rest_errors"
	"test-task1/internal/models"
)

type UserService interface {
	SignUp(*models.SignUpInput) (*models.User, error)
	Login(email, password string) (string, error)
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

	var signUpInputDao SignUpInputDAO
	if err := json.NewDecoder(r.Body).Decode(&signUpInputDao); err != nil {
		rest_errors.HandleError(w, rest_errors.ErrBadRequest, http.StatusBadRequest)
	}

	v := c.service.GetValidator()
	if err := signUpInputDao.ValidateWith(v); err != nil {
		rest_errors.HandleError(w, err, http.StatusBadRequest)
		return
	}

	singUpInput := signUpInputDao.ToSignUpInput()

	user, err := c.service.SignUp(singUpInput)
	if errors.Is(err, models.ErrUserAlreadyExists) {
		rest_errors.HandleError(w, err, http.StatusConflict) // 409
		return
	}
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

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var LoginDao LoginInputDAO
	if err := json.NewDecoder(r.Body).Decode(&LoginDao); err != nil {
		rest_errors.HandleError(w, rest_errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	v := c.service.GetValidator()
	if err := LoginDao.ValidateWith(v); err != nil {
		rest_errors.HandleError(w, err, http.StatusBadRequest)
		return
	}

	LoginInput := LoginDao.ToLoginInput()

	token, err := c.service.Login(LoginInput.Email, LoginInput.Password)
	if errors.Is(err, models.ErrInvalidCredentials) {
		rest_errors.HandleError(w, err, http.StatusUnauthorized)
		return
	}
	if err != nil {
		rest_errors.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	tokenOutput := ToTokenDAO(token)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tokenOutput); err != nil {
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

	UserListOutput := ToUserListDAO(users)

	if err := json.NewEncoder(w).Encode(UserListOutput); err != nil {
		rest_errors.HandleError(w, rest_errors.ErrInternalServer, http.StatusInternalServerError)
		return
	}
}
