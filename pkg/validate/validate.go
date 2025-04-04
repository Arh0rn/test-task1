package validate

import "github.com/go-playground/validator/v10"

// TODO: validator errors not user friendly, need to add custom error messages in future
func New() *validator.Validate {
	var validate *validator.Validate
	validate = validator.New()
	return validate
}
