package daos

import "github.com/Arh0rn/test-task1/internal/domain"

type UserOutputDAO struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ToUserOutputDAO(user *domain.User) *UserOutputDAO {
	return &UserOutputDAO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

type UserListDAO struct {
	Users []UserOutputDAO `json:"users"`
}

func ToUserListDAO(users []*domain.User) *UserListDAO {
	var userList []UserOutputDAO
	for _, user := range users {
		userList = append(userList, *ToUserOutputDAO(user))
	}
	return &UserListDAO{
		Users: userList,
	}
}
