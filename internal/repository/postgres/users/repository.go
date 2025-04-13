package postgresUsersRepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Arh0rn/test-task1/internal/domain"
	"github.com/lib/pq"
	"log/slog"
)

// TODO: make soft delete

type UserRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.SignUpInput) (*domain.User, error) {
	var id int
	slog.DebugContext(ctx, "Creating user in DB", "user", user)
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO users (name, email, password) 
		 VALUES ($1, $2, $3) 
		 RETURNING id`,
		user.Name, user.Email, user.Password,
	).Scan(&id)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			slog.ErrorContext(ctx, "User already exists")
			return nil, domain.ErrUserAlreadyExists
		}
		slog.ErrorContext(ctx, "Failed to create user", "error", err)
		return nil, err
	}
	createdUser := &domain.User{
		ID:       id,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	slog.DebugContext(ctx, "User created", "user", createdUser)
	return createdUser, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	slog.DebugContext(ctx, "Getting user by email", "email", email)
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, password 
		 FROM users 
		 WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.ErrorContext(ctx, "User not found", "email", email)
			return nil, domain.ErrUserNotFound

		}
		slog.ErrorContext(ctx, "Failed to get user by email", "error", err)
		return nil, err
	}
	slog.DebugContext(ctx, "User found", "user", user)
	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	slog.DebugContext(ctx, "Getting all users")
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, name, email, password FROM users")
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get all users", "error", err)
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
			slog.ErrorContext(ctx, "Failed to get all users", "error", err)
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	slog.DebugContext(ctx, "All users retrieved", "users", users)
	return users, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	slog.DebugContext(ctx, "Getting user by ID", "id", id)
	var user domain.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, password 
		 FROM users 
		 WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.ErrorContext(ctx, "User not found", "id", id)
			return nil, domain.ErrUserNotFound
		}
		slog.ErrorContext(ctx, "Failed to get user by ID", "error", err)
		return nil, err
	}

	slog.DebugContext(ctx, "User found", "user", user)
	return &user, nil
}

func (r *UserRepository) DeleteByID(ctx context.Context, id int) error {
	slog.DebugContext(ctx, "Deleting user by ID", "id", id)
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM users WHERE id = $1`,
		id,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to delete user", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to delete user", "error", err)
		return err
	}

	if rowsAffected == 0 {
		slog.ErrorContext(ctx, "User does not exist", "id", id)
		return domain.ErrUserNotFound
	}

	slog.DebugContext(ctx, "User deleted", "id", id)
	return nil
}

func (r *UserRepository) UpdateByID(ctx context.Context, user *domain.UserUpdate, id int) (*domain.UserUpdate, error) {
	slog.DebugContext(ctx, "Updating user by ID", "id", id)
	result, err := r.db.ExecContext(ctx,
		`UPDATE users SET name = $1, email = $2 WHERE id = $3`,
		user.Name, user.Email, id,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to update user", "error", err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to update user", "error", err)
		return nil, err
	}

	if rowsAffected == 0 {
		slog.ErrorContext(ctx, "User does not exist", "id", id)
		return nil, domain.ErrUserNotFound
	}

	slog.DebugContext(ctx, "User updated", "user", user)
	return user, nil
}
