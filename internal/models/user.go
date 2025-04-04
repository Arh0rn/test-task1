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

type SingUpOutput struct {
	ID    int
	Name  string
	Email string
}
