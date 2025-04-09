package usersService

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"test-task1/internal/models"
	"test-task1/pkg/jwtoken"
	"time"
)

type UserRepository interface {
	Create(context.Context, *models.SignUpInput) (*models.User, error)
	GetAll(context.Context) ([]*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	UpdateByID(ctx context.Context, user *models.UserUpdate, id int) (*models.UserUpdate, error)
	DeleteByID(ctx context.Context, id int) error
}

type Hasher interface {
	Hash(password string) (string, error)
	Verify(password, hashed string) bool
}

type UserService struct {
	repo UserRepository

	hasher    Hasher
	validator *validator.Validate

	jwtSecret []byte
	tokenTTL  time.Duration
}

func New(
	repo UserRepository,
	hasher Hasher,
	validator *validator.Validate,
	jwts []byte,
	tttl time.Duration,
) *UserService {
	return &UserService{
		repo:      repo,
		hasher:    hasher,
		validator: validator,
		jwtSecret: jwts,
		tokenTTL:  tttl,
	}
}

func (s *UserService) SignUp(ctx context.Context, userInput *models.SignUpInput) (*models.User, error) {

	hashedPassword, err := s.hasher.Hash(userInput.Password)
	if err != nil {
		return nil, err
	}

	userInput.Password = hashedPassword
	user, err := s.repo.Create(ctx, userInput)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if errors.Is(err, models.ErrUserNotFound) {
		return "", models.ErrInvalidCredentials

	}
	if err != nil {
		return "", err
	}

	valid := s.hasher.Verify(password, user.Password)
	if !valid {
		return "", models.ErrInvalidCredentials
	}
	token, err := jwtoken.GenerateToken(user.ID, user.Email, s.jwtSecret, s.tokenTTL)
	if err != nil {
		return "", err
	}

	slog.DebugContext(ctx, "Token generated", "token", token)
	return token, nil
}

func (s *UserService) GetAll(ctx context.Context) ([]*models.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateByID(ctx context.Context, user *models.UserUpdate, id int) (*models.UserUpdate, error) {
	user, err := s.repo.UpdateByID(ctx, user, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) DeleteByID(ctx context.Context, id int) error {
	err := s.repo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetValidator() *validator.Validate {
	return s.validator
}
