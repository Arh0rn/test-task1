package daos

import (
	"github.com/go-playground/validator/v10"
	"test-task1/internal/domain"
)

type SignUpInputDAO struct {
	Name     string `json:"name" validate:"required,gte=3,lte=32" example:"John Doe"`
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required,gte=6,lte=32" example:"P@ssw0rd"`
}

func (dao *SignUpInputDAO) ValidateWith(v *validator.Validate) error {
	return v.Struct(dao)
}

func (dao *SignUpInputDAO) ToSignUpInput() *domain.SignUpInput {
	return &domain.SignUpInput{
		Name:     dao.Name,
		Email:    dao.Email,
		Password: dao.Password,
	}
}
func ToSignUpInputDAO(input *domain.SignUpInput) *SignUpInputDAO {
	return &SignUpInputDAO{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}
}
