package validate

import "github.com/go-playground/validator/v10"

func New() *validator.Validate {
	var validate *validator.Validate
	validate = validator.New()
	return validate
}
