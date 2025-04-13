package daos

import (
	"github.com/Arh0rn/test-task1/internal/domain"
	"github.com/go-playground/validator/v10"
)

type UserUpdateDAO struct {
	Name  string `json:"name" validate:"required,gte=3,lte=32" example:"John Doe"`
	Email string `json:"email" validate:"required,email" example:"john.doe@example.com"`
}

func (dao *UserUpdateDAO) ValidateWith(v *validator.Validate) error {
	return v.Struct(dao)
}
func (dao *UserUpdateDAO) ToUserUpdate() *domain.UserUpdate {
	return &domain.UserUpdate{
		Name:  dao.Name,
		Email: dao.Email,
	}
}

func ToUserUpdateDAO(user *domain.UserUpdate) *UserUpdateDAO {
	return &UserUpdateDAO{
		Name:  user.Name,
		Email: user.Email,
	}
}
