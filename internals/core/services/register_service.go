package services

import (
	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/PitiNarak/condormhub-backend/internals/core/ports"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo ports.UserRepository
}

func NewUserService(UserRepo ports.UserRepository) ports.UserService {
	return &UserService{UserRepo: UserRepo}
}

func (s *UserService) Create(user domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	create_err := s.UserRepo.Create(user)
	if create_err != nil {
		return create_err
	}

	return nil
}
