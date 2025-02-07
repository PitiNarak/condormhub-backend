package services

import (
	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) ResetPasswordCreate(email string) (domain.User, error) {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *UserService) ResetPasswordResponse(userID uuid.UUID, password string) error {
	user, err := s.UserRepo.GetUser(userID)
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	err = s.UserRepo.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}
