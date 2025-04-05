package daos

import "test-task1/internal/models"

type UserOutputDAO struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

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

func ToUserListDAO(users []*models.User) *UserListDAO {
	var userList []UserOutputDAO
	for _, user := range users {
		userList = append(userList, *ToUserOutputDAO(user))
	}
	return &UserListDAO{
		Users: userList,
	}
}
