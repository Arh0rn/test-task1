package models

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
	//ErrUserInvalid  = errors.New("user invalid")

)
