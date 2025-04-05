package postgresUsersRepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"test-task1/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.SignUpInput) (*models.User, error) {
	var id int

	err := r.db.QueryRowContext(ctx,
		`INSERT INTO users (name, email, password) 
		 VALUES ($1, $2, $3) 
		 RETURNING id`,
		user.Name, user.Email, user.Password,
	).Scan(&id)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, models.ErrUserAlreadyExists
		}
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

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, password 
		 FROM users 
		 WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, name, email, password FROM users")
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

func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, password 
		 FROM users 
		 WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) DeleteByID(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM users WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) UpdateByID(ctx context.Context, user *models.UserUpdate, id int) (*models.UserUpdate, error) {
	result, err := r.db.ExecContext(ctx,
		`UPDATE users SET name = $1, email = $2 WHERE id = $3`,
		user.Name, user.Email, id,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, models.ErrUserNotFound
	}

	return user, nil
}
