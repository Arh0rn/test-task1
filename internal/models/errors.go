package models

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("email or password is incorrect")
	ErrValidation         = errors.New("invalid email or password, password must be at least 8 characters long")

	//ErrUserInvalid  = rest_errors.New("user invalid")

)
