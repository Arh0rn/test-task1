package usersService

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"test-task1/internal/domain"
	"test-task1/pkg/jwtoken"
	"time"
)

type UserRepository interface {
	Create(context.Context, *domain.SignUpInput) (*domain.User, error)
	GetAll(context.Context) ([]*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
	UpdateByID(ctx context.Context, user *domain.UserUpdate, id int) (*domain.UserUpdate, error)
	DeleteByID(ctx context.Context, id int) error
}

type UserCache interface {
	Set(context.Context, *domain.User) error
	GetAll(context.Context) ([]*domain.User, error)
	SetAll(context.Context, []*domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	UpdateByID(ctx context.Context, user *domain.UserUpdate, id int) error
	DeleteByID(ctx context.Context, id int) error
}

type Hasher interface {
	Hash(password string) (string, error)
	Verify(password, hashed string) bool
}

type UserService struct {
	repo  UserRepository
	cache UserCache

	hasher    Hasher
	validator *validator.Validate

	jwtSecret []byte
	tokenTTL  time.Duration
}

func New(
	repo UserRepository,
	cache UserCache,
	hasher Hasher,
	validator *validator.Validate,
	jwts []byte,
	tttl time.Duration,
) *UserService {
	return &UserService{
		repo:      repo,
		cache:     cache,
		hasher:    hasher,
		validator: validator,
		jwtSecret: jwts,
		tokenTTL:  tttl,
	}
}

func (s *UserService) SignUp(ctx context.Context, userInput *domain.SignUpInput) (*domain.User, error) {

	hashedPassword, err := s.hasher.Hash(userInput.Password)
	if err != nil {
		return nil, err
	}

	userInput.Password = hashedPassword
	user, err := s.repo.Create(ctx, userInput)
	if err != nil {
		return nil, err
	}

	go func() {
		err = s.cache.Set(context.Background(), user)
	}()

	return user, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if errors.Is(err, domain.ErrUserNotFound) {
		return "", domain.ErrInvalidCredentials

	}
	if err != nil {
		return "", err
	}

	valid := s.hasher.Verify(password, user.Password)
	if !valid {
		return "", domain.ErrInvalidCredentials
	}
	token, err := jwtoken.GenerateToken(user.ID, user.Email, s.jwtSecret, s.tokenTTL)
	if err != nil {
		return "", err
	}

	slog.DebugContext(ctx, "Token generated", "token", token)
	return token, nil
}

func (s *UserService) GetAll(ctx context.Context) ([]*domain.User, error) {
	//TODO: add pagination if needed
	//users, err := s.cache.GetAll(ctx)
	//if err == nil && len(users) > 0 {
	//	return users, nil
	//}

	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		err = s.cache.SetAll(context.Background(), users)
	}()
	return users, nil
}

func (s *UserService) GetByID(ctx context.Context, id int) (*domain.User, error) {
	user, err := s.cache.GetByID(ctx, id)
	if err == nil && user != nil {
		slog.DebugContext(ctx, "User found in cache", "user", user)
		return user, nil
	}

	user, err = s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	go func() {
		err = s.cache.Set(context.Background(), user)
	}()

	return user, nil
}

func (s *UserService) UpdateByID(ctx context.Context, user *domain.UserUpdate, id int) (*domain.UserUpdate, error) {
	user, err := s.repo.UpdateByID(ctx, user, id)
	if err != nil {
		return nil, err
	}
	//Strange decision, but otherwise i need to make double request to db
	//Or add extra logic to the repo.UpdateByID method that allows to decide when to get the user from db
	go func() {
		err = s.cache.UpdateByID(context.Background(), user, id)
		if err == nil {
			slog.DebugContext(ctx, "User updated in cache", "user", user)
			return
		}
		dbUser, err := s.repo.GetByID(context.Background(), id)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to get user by ID", "id", id)
			return
		}
		err = s.cache.Set(context.Background(), dbUser)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to set user in cache", "id", id)
			return
		}
	}()
	return user, nil
}

func (s *UserService) DeleteByID(ctx context.Context, id int) error {
	err := s.repo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}

	go func() {
		err = s.cache.DeleteByID(context.Background(), id)
	}()
	return nil
}

func (s *UserService) GetValidator() *validator.Validate {
	return s.validator
}
