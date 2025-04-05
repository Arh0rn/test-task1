package daos

import (
	"github.com/go-playground/validator/v10"
	"test-task1/internal/models"
)

type LoginInputDAO struct {
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required" example:"P@ssw0rd"`
}

func (dao *LoginInputDAO) ValidateWith(v *validator.Validate) error {
	return v.Struct(dao)
}

func (dao *LoginInputDAO) ToLoginInput() *models.LoginInput {
	return &models.LoginInput{
		Email:    dao.Email,
		Password: dao.Password,
	}
}
