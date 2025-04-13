package usersController

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Arh0rn/test-task1/internal/controller/restapi/controllers/users/daos"
	"github.com/Arh0rn/test-task1/internal/controller/restapi/rest_errors"
	"github.com/Arh0rn/test-task1/internal/domain"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

type UserService interface {
	SignUp(context.Context, *domain.SignUpInput) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	GetAll(context.Context) ([]*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
	UpdateByID(ctx context.Context, user *domain.UserUpdate, id int) (*domain.UserUpdate, error)
	DeleteByID(ctx context.Context, id int) error
	GetValidator() *validator.Validate
}

type UserController struct {
	service UserService
}

func New(service UserService) *UserController {
	return &UserController{service: service}
}

// SignUp godoc
// @Summary      Register new user
// @Description  Creates a new user with name, email, and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      daos.SignUpInputDAO  true  "User sign up input"
// @Success      201    {object}  daos.UserOutputDAO
// @Failure      400    {object}  rest_errors.ResponseError
// @Failure      409    {object}  rest_errors.ResponseError
// @Failure      500    {object}  rest_errors.ResponseError
// @Router       /users [post]
func (c *UserController) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var signUpInputDao daos.SignUpInputDAO
	if err := json.NewDecoder(r.Body).Decode(&signUpInputDao); err != nil {
		rest_errors.HandleError(w, rest_errors.ErrBadRequest, http.StatusBadRequest)
	}

	v := c.service.GetValidator()
	if err := signUpInputDao.ValidateWith(v); err != nil {
		rest_errors.HandleError(w, domain.ErrValidation, http.StatusBadRequest)
		return
	}

	singUpInput := signUpInputDao.ToSignUpInput()

	user, err := c.service.SignUp(ctx, singUpInput)
	if errors.Is(err, domain.ErrUserAlreadyExists) {
		rest_errors.HandleError(w, err, http.StatusConflict) // 409
		return
	}
	if err != nil {
		rest_errors.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	userOutput := daos.ToUserOutputDAO(user)

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(userOutput); err != nil {
		rest_errors.HandleError(w, rest_errors.ErrInternalServer, http.StatusInternalServerError)
		return
	}
}

// Login godoc
// @Summary      User login
// @Description  Authenticates a user and returns JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      daos.LoginInputDAO  true  "User login input"
// @Success      200    {object}  daos.TokenDAO
// @Failure      400    {object}  rest_errors.ResponseError
// @Failure      401    {object}  rest_errors.ResponseError
// @Failure      422    {object}  rest_errors.ResponseError
// @Failure      500    {object}  rest_errors.ResponseError
// @Router       /login [post]
func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var LoginDao daos.LoginInputDAO
	if err := json.NewDecoder(r.Body).Decode(&LoginDao); err != nil {
		rest_errors.HandleError(w, rest_errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	v := c.service.GetValidator()
	if err := LoginDao.ValidateWith(v); err != nil {
		rest_errors.HandleError(w, domain.ErrValidation, http.StatusUnprocessableEntity)
		return
	}

	LoginInput := LoginDao.ToLoginInput()

	token, err := c.service.Login(ctx, LoginInput.Email, LoginInput.Password)
	if errors.Is(err, domain.ErrInvalidCredentials) {
		rest_errors.HandleError(w, err, http.StatusUnauthorized)
		return
	}
	if err != nil {
		rest_errors.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	tokenOutput := daos.ToTokenDAO(token)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tokenOutput); err != nil {
		rest_errors.HandleError(w, rest_errors.ErrInternalServer, http.StatusInternalServerError)
		return
	}
}

// GetAll godoc
// @Summary      Get all users
// @Tags         users
// @Security  BearerAuth
// @Produce      json
// @Success      200  {array}   daos.UserOutputDAO
// @Failure      401  {object}  rest_errors.ResponseError
// @Failure      500  {object}  rest_errors.ResponseError
// @Router       /users [get]
func (c *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	users, err := c.service.GetAll(ctx)
	if err != nil {
		rest_errors.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	UserListOutput := daos.ToUserListDAO(users)

	if err := json.NewEncoder(w).Encode(UserListOutput); err != nil {
		rest_errors.HandleError(w, rest_errors.ErrInternalServer, http.StatusInternalServerError)
		return
	}
}

// GetByID godoc
// @Summary      Get user by ID
// @Description  Retrieves a single user by their ID
// @Tags         users
// @Security  BearerAuth
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  daos.UserOutputDAO
// @Failure      400  {object}  rest_errors.ResponseError
// @Failure      401  {object}  rest_errors.ResponseError
// @Failure      404  {object}  rest_errors.ResponseError
// @Failure      500  {object}  rest_errors.ResponseError
// @Router       /users/{id} [get]
func (c *UserController) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		rest_errors.HandleError(w, rest_errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	user, err := c.service.GetByID(ctx, id)
	if errors.Is(err, domain.ErrUserNotFound) {
		rest_errors.HandleError(w, err, http.StatusNotFound)
		return
	}
	if err != nil {
		rest_errors.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	userOutput := daos.ToUserOutputDAO(user)

	if err := json.NewEncoder(w).Encode(userOutput); err != nil {
		rest_errors.HandleError(w, rest_errors.ErrInternalServer, http.StatusInternalServerError)
		return
	}
}

// UpdateByID godoc
// @Summary      Update user by ID
// @Description  Updates user fields like name or email by their ID
// @Tags         users
// @Security  BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      int                   true  "User ID"
// @Param        input body      daos.UserUpdateDAO    true  "User update input"
// @Success      200   {object}  daos.UserUpdateDAO
// @Failure      400   {object}  rest_errors.ResponseError
// @Failure      401   {object}  rest_errors.ResponseError
// @Failure      404   {object}  rest_errors.ResponseError
// @Failure      500   {object}  rest_errors.ResponseError
// @Router       /users/{id} [put]
func (c *UserController) UpdateByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		rest_errors.HandleError(w, rest_errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	var userDao daos.UserUpdateDAO
	if err := json.NewDecoder(r.Body).Decode(&userDao); err != nil {
		rest_errors.HandleError(w, rest_errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	v := c.service.GetValidator()
	if err := userDao.ValidateWith(v); err != nil {
		rest_errors.HandleError(w, domain.ErrValidation, http.StatusBadRequest)
		return
	}

	userUpdate := userDao.ToUserUpdate()

	userUpdateOutput, err := c.service.UpdateByID(ctx, userUpdate, id)
	if errors.Is(err, domain.ErrUserNotFound) {
		rest_errors.HandleError(w, err, http.StatusNotFound)
		return
	}
	if err != nil {
		rest_errors.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	userUpdateDao := daos.ToUserUpdateDAO(userUpdateOutput)

	if err := json.NewEncoder(w).Encode(userUpdateDao); err != nil {
		rest_errors.HandleError(w, rest_errors.ErrInternalServer, http.StatusInternalServerError)
		return
	}
}

// DeleteByID godoc
// @Summary      Delete user by ID
// @Description  Deletes a user from the system by their ID
// @Tags         users
// @Security  BearerAuth
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      204  "No Content"
// @Failure      400  {object}  rest_errors.ResponseError
// @Failure      401  {object}  rest_errors.ResponseError
// @Failure      404  {object}  rest_errors.ResponseError
// @Failure      500  {object}  rest_errors.ResponseError
// @Router       /users/{id} [delete]
func (c *UserController) DeleteByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		rest_errors.HandleError(w, rest_errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	err = c.service.DeleteByID(ctx, id)
	if errors.Is(err, domain.ErrUserNotFound) {
		rest_errors.HandleError(w, err, http.StatusNotFound)
		return
	}
	if err != nil {
		rest_errors.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
