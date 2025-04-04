package postgresUsersRepo

import (
	"database/sql"
	"test-task1/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.SignUpInput) (*models.User, error) {
	var id int

	err := r.db.QueryRow(
		`INSERT INTO users (name, email, password) 
		 VALUES ($1, $2, $3) 
		 RETURNING id`,
		user.Name, user.Email, user.Password,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	createdUser := &models.User{
		ID:       id,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	return createdUser, nil
}

func (r *UserRepository) GetAll() ([]*models.User, error) {
	rows, err := r.db.Query("SELECT id, name, email, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

//
//func (r *UserRepository) GetByID(id int) (*User, error) {
//
//}
//
//func (r *UserRepository) Delete(users *User) error {
//
//}
