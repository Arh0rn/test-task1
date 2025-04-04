package usersService

import (
	"github.com/go-playground/validator/v10"
	"test-task1/internal/models"
)

type UserRepository interface {
	Create(user *models.SignUpInput) (*models.User, error)
	GetAll() ([]*models.User, error)
}

type Hasher interface {
	Hash(password string) (string, error)
}

type UserService struct {
	repo      UserRepository
	hasher    Hasher
	validator *validator.Validate
}

func New(repo UserRepository, hasher Hasher, validator *validator.Validate) *UserService {
	return &UserService{repo: repo, hasher: hasher, validator: validator}
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
