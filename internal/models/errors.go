package models

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("email or password is incorrect")
	//ErrUserInvalid  = rest_errors.New("user invalid")

)
