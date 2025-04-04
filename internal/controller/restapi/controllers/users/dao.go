package usersController

import (
	"github.com/go-playground/validator/v10"
	"test-task1/internal/models"
)

type UserDAO struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dao *UserDAO) ToUser() *models.User {
	return &models.User{
		ID:       dao.ID,
		Name:     dao.Name,
		Email:    dao.Email,
		Password: dao.Password,
	}
}

func ToUserDAO(user *models.User) *UserDAO {
	return &UserDAO{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
}

type SignUpInputDAO struct {
	Name     string `json:"name" validate:"required,gte=3,lte=32" example:"John Doe"`
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required,gte=6,lte=32" example:"P@ssw0rd"`
}

func (dao *SignUpInputDAO) ValidateWith(v *validator.Validate) error {
	return v.Struct(dao)
}

func (dao *SignUpInputDAO) ToSignUpInput() *models.SignUpInput {
	return &models.SignUpInput{
		Name:     dao.Name,
		Email:    dao.Email,
		Password: dao.Password,
	}
}

func ToSignUpInputDAO(input *models.SignUpInput) *SignUpInputDAO {
	return &SignUpInputDAO{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}
}

type SingUpOutputDAO struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (dao *SingUpOutputDAO) ToSignUpOutput() *models.SingUpOutput {
	return &models.SingUpOutput{
		ID:    dao.ID,
		Name:  dao.Name,
		Email: dao.Email,
	}
}

func ToSignUpOutputDAO(user *models.SingUpOutput) *SingUpOutputDAO {
	return &SingUpOutputDAO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
