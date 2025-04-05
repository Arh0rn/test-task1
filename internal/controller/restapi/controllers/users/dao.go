package usersController

import (
	"github.com/go-playground/validator/v10"
	"test-task1/internal/models"
)

//type UserDAO struct {
//	ID       int    `json:"id"`
//	Name     string `json:"name"`
//	Email    string `json:"email"`
//	Password string `json:"password"`
//}
//
//func (dao *UserDAO) ToUser() *models.User {
//	return &models.User{
//		ID:       dao.ID,
//		Name:     dao.Name,
//		Email:    dao.Email,
//		Password: dao.Password,
//	}
//}
//
//func ToUserDAO(user *models.User) *UserDAO {
//	return &UserDAO{
//		ID:       user.ID,
//		Name:     user.Name,
//		Email:    user.Email,
//		Password: user.Password,
//	}
//}

type SignUpInputDAO struct {
	Name     string `json:"name" validate:"required,gte=3,lte=32" example:"John Doe"`
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required,gte=6,lte=32" example:"P@ssw0rd"`
}

func (dao *SignUpInputDAO) ValidateWith(v *validator.Validate) error {
	return v.Struct(dao)
}

func (dao *SignUpInputDAO) ToSignUpInput() *models.SignUpInput {
	return &models.SignUpInput{
		Name:     dao.Name,
		Email:    dao.Email,
		Password: dao.Password,
	}
}

func ToSignUpInputDAO(input *models.SignUpInput) *SignUpInputDAO {
	return &SignUpInputDAO{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}
}

type LoginInputDAO struct {
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required,gte=6,lte=32" example:"P@ssw0rd"`
}

func (dao *LoginInputDAO) ValidateWith(v *validator.Validate) error {
	return v.Struct(dao)
}

func (dao *LoginInputDAO) ToLoginInput() *models.LoginInput {
	return &models.LoginInput{
		Email:    dao.Email,
		Password: dao.Password,
	}
}

type UserOutputDAO struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

//
//func (dao *UserOutputDAO) ToUserOutput() *models.UserOutput {
//	return &models.UserOutput{
//		ID:    dao.ID,
//		Name:  dao.Name,
//		Email: dao.Email,
//	}
//}

func ToUserOutputDAO(user *models.User) *UserOutputDAO {
	return &UserOutputDAO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

type UserListDAO struct {
	Users []UserOutputDAO `json:"users"`
}

//func (dao *UserListDAO) ToUserList() []*models.UserOutput {
//	var users []*models.UserOutput
//	for _, user := range dao.Users {
//		users = append(users, user.ToUserOutput())
//	}
//	return users
//}

func ToUserListDAO(users []*models.User) *UserListDAO {
	var userList []UserOutputDAO
	for _, user := range users {
		userList = append(userList, *ToUserOutputDAO(user))
	}
	return &UserListDAO{
		Users: userList,
	}
}

type TokenDAO struct {
	Token string `json:"token"`
}

func ToTokenDAO(token string) *TokenDAO {
	return &TokenDAO{
		Token: token,
	}
}
