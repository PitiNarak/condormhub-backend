package services

import (
	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/PitiNarak/condormhub-backend/internals/core/ports"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo ports.UserRepository
}

func NewUserService(UserRepo ports.UserRepository) ports.UserService {
	return &UserService{UserRepo: UserRepo}
}

func (s *UserService) Create(user *domain.User) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)
	create_err := s.UserRepo.Create(user)
	if create_err != nil {
		return nil, create_err
	}

	return user, nil
}

func (s *UserService) VerifyUser(userID uuid.UUID) error {
	user, err := s.UserRepo.GetUser(userID)
	if err != nil || user.ID == uuid.Nil {
		return err
	}

	user.IsVerified = true
	return s.UserRepo.UpdateUser(user)
}
