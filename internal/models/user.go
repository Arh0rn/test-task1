package models

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type SignUpInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type UserOutput struct {
	ID    int
	Name  string
	Email string
}

type UserUpdate struct {
	Name  string
	Email string
}
