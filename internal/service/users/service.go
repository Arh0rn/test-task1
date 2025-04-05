package usersService

import (
	"github.com/go-playground/validator/v10"
	"test-task1/internal/models"
	"test-task1/pkg/jwt_token"
	"time"
)

type UserRepository interface {
	Create(user *models.SignUpInput) (*models.User, error)
	GetAll() ([]*models.User, error)
	GetByEmail(email string) (*models.User, error)
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

func (s *UserService) SignUp(userInput *models.SignUpInput) (*models.User, error) {

	hashedPassword, err := s.hasher.Hash(userInput.Password)
	if err != nil {
		return nil, err
	}

	userInput.Password = hashedPassword
	user, err := s.repo.Create(userInput)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	valid := s.hasher.Verify(password, user.Password)
	if !valid {
		return "", models.ErrInvalidCredentials
	}
	token, err := jwt_token.GenerateToken(user.ID, user.Email, s.jwtSecret, s.tokenTTL)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) GetAll() ([]*models.User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetValidator() *validator.Validate {
	return s.validator
}
