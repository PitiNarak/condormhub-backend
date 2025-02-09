package services

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/pkg/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) ResetPasswordCreate(email string) error {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	err = s.EmailService.SendResetPasswordEmail(user.Email, user.UserName, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) ResetPasswordResponse(token string, password string) error {
	claims, err := utils.DecodeJWT(token, s.Config)
	if err != nil {
		return err
	}
	userIDstr, ok := (*claims)["user_id"].(string)
	if !ok {
		return errors.New("cannot get user_id")
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return err
	}
	user, err := s.UserRepo.GetUser(userID)
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	err = s.UserRepo.UpdateUser(*user)
	if err != nil {
		return err
	}
	return nil
}
