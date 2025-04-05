package daos

import (
	"github.com/go-playground/validator/v10"
	"test-task1/internal/models"
)

type UserUpdateDAO struct {
	Name  string `json:"name" validate:"omitempty,gte=3,lte=32" example:"John Doe"`
	Email string `json:"email" validate:"required,email" example:"john.doe@example.com"`
}

func (dao *UserUpdateDAO) ValidateWith(v *validator.Validate) error {
	return v.Struct(dao)
}
func (dao *UserUpdateDAO) ToUserUpdate() *models.UserUpdate {
	return &models.UserUpdate{
		Name:  dao.Name,
		Email: dao.Email,
	}
}

func ToUserUpdateDAO(user *models.UserUpdate) *UserUpdateDAO {
	return &UserUpdateDAO{
		Name:  user.Name,
		Email: user.Email,
	}
}
